package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/gonejack/glogger"
	"log"
	"os"
	"sync"
	"time"
)

var logger = glogger.NewLogger("Service:Kafka")

var consumers = make(map[string]*consumeWrapper)

func main() {
	sarama.Logger = log.New(os.Stdout, "[Sarama] ", log.LstdFlags)

	go consume()

	go produce()

	time.Sleep(time.Hour)
}

func consume() {
	topic := "test"
	groupId := "abc"
	consumerId := fmt.Sprintf("%s|%s", topic, groupId)

	consumer := &consumeWrapper{
		Addrs: []string{
			"127.0.0.1:9092",
			"127.0.0.1:9093",
		},
		Topic:   topic,
		GroupId: groupId,
		//SASL: struct {
		//	Enable   bool
		//	User     string
		//	Password string
		//}{
		//	Enable:   true,
		//	User:     "admin",
		//	Password: "admin",
		//},
		mu: sync.RWMutex{},
	}

	err := consumer.addSubscriber(func(topic string, byts []byte) error {
		logger.Infof("收到 topic=%s msg=%s", topic, byts)

		return nil
	}).start() // 启动消费

	if err == nil {
		consumers[consumerId] = consumer
	} else {
		logger.Errorf("订阅[topic=%s, groupId=%s]出错: %s", topic, groupId, err)
	}
}

func produce() {
	config := sarama.NewConfig()
	config.Producer.Return.Errors = true // 消费者需要会报错
	//config.Net.SASL.Enable = true        // 启用密码验证
	//config.Net.SASL.User = "admin"
	//config.Net.SASL.Password = "admin"

	producer, err := sarama.NewAsyncProducer([]string{
		"127.0.0.1:9092",
		"127.0.0.1:9093",
	}, config)
	if err == nil {
		wrapper := &producerWrapper{
			producer: producer,
			msg:      make(chan *sarama.ProducerMessage, 1),
			close:    make(chan os.Signal),
		}
		wrapper.start()

		go func() {
			topic := "test"
			for {
				time.Sleep(time.Second)

				msg := fmt.Sprintf("%d", time.Now().Unix())
				wrapper.input() <- &sarama.ProducerMessage{
					Topic: "test",
					Value: sarama.ByteEncoder(msg),
				}

				logger.Infof("发送topic=%s, msg=%s", topic, msg)
			}
		}()
	} else {
		logger.Errorf("新建生产者出错: %s", err)
	}
}

type producerWrapper struct {
	producer sarama.AsyncProducer
	msg      chan *sarama.ProducerMessage
	close    chan os.Signal
}

func (pw *producerWrapper) input() chan<- *sarama.ProducerMessage {
	return pw.msg
}
func (pw *producerWrapper) start() {
	go pw.produceMessage()
	go pw.readError()
}
func (pw *producerWrapper) stop() {
	pw.close <- os.Interrupt
}
func (pw *producerWrapper) produceMessage() {
	stopping := false

	for {
		select {
		case <-pw.close:
			stopping = true
		case msg := <-pw.msg:
			if stopping {
				logger.Errorf("生产者已关闭，无法发送消息，请检查关机顺序是否合理，消息[topic=%s, value=%s]", msg.Topic, msg.Value)
			} else {
				pw.producer.Input() <- msg
			}
		}
	}
}
func (pw *producerWrapper) readError() {
	for err := range pw.producer.Errors() {
		logger.Errorf("生产者错误: %s, 消息[topic=%s, value=%s]", err.Err, err.Msg.Topic, err.Msg.Value)
	}
}

// 消费组封装
type consumeWrapper struct {
	Addrs   []string
	Topic   string
	GroupId string

	SASL struct {
		Enable   bool
		User     string
		Password string
	}

	mu sync.RWMutex

	consumer    *cluster.Consumer
	subscribers []func(topic string, msg []byte) error
}

func (cw *consumeWrapper) start() (err error) {
	config := cluster.NewConfig()

	config.Consumer.Return.Errors = true    // 设置对外报错
	config.Net.SASL.Enable = cw.SASL.Enable // 是否启用密码
	config.Net.SASL.User = cw.SASL.User
	config.Net.SASL.Password = cw.SASL.Password

	cw.consumer, err = cluster.NewConsumer(cw.Addrs, cw.GroupId, []string{cw.Topic}, config)

	if err == nil {
		logger.Infof("消费组[topic=%s, groupId=%s]就绪", cw.Topic, cw.GroupId)

		go cw.readMessage()
		go cw.readError()
	} else {
		logger.Errorf("新建消费者[topic=%s, groupId=%s]失败: %s", cw.Topic, cw.GroupId, err)
	}

	return
}
func (cw *consumeWrapper) stop() (err error) {
	err = cw.consumer.Close()

	if err == nil {
		logger.Infof("消费组[topic=%s, groupId=%s]已关闭", cw.Topic, cw.GroupId)
	} else {
		logger.Errorf("关闭消费组[topic=%s, groupId=%s]出错: %s", err)
	}

	return
}
func (cw *consumeWrapper) readMessage() {
	for msg := range cw.consumer.Messages() {
		cw.mu.RLock()

		for _, subscriber := range cw.subscribers {
			err := subscriber(msg.Topic, msg.Value)

			if err != nil {
				logger.Errorf("消费组[topic=%s, groupId=%s]订阅者出错: %s", cw.Topic, cw.GroupId, err)
			}
		}

		cw.mu.RUnlock()
	}
}
func (cw *consumeWrapper) readError() {
	for err := range cw.consumer.Errors() {
		logger.Errorf("消费发生错误: %s", err)
	}
}
func (cw *consumeWrapper) addSubscriber(handler func(topic string, msg []byte) error) *consumeWrapper {
	cw.mu.Lock()
	cw.subscribers = append(cw.subscribers, handler)
	cw.mu.Unlock()

	return cw
}
