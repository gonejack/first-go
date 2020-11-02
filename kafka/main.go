package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/gonejack/glogger"
	"log"
	"os"
	"sync"
	"time"
)

var logger = glogger.NewLogger("Service:Kafka")
var consumers = make(map[string]*consumer)

func init() {
	sarama.Logger = log.New(os.Stdout, "[Sarama] ", log.LstdFlags)
}

func main() {
	go consume()

	//go produce()

	time.Sleep(time.Hour)
}
func consume() {
	topic := "test"
	groupId := "abc"
	consumerId := fmt.Sprintf("%s|%s", topic, groupId)

	consumer := &consumer{
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

	time.Sleep(time.Minute)
	consumer.stop()
}
func produce() {
	config := sarama.NewConfig()
	config.Producer.Return.Errors = true // 消费者需要会报错
	//config.Net.SASL.Enable = true        // 启用密码验证
	//config.Net.SASL.User = "admin"
	//config.Net.SASL.Password = "admin"

	prd, err := sarama.NewAsyncProducer([]string{
		"127.0.0.1:9092",
		"127.0.0.1:9093",
	}, config)
	if err == nil {
		wrapper := &producer{
			producer: prd,
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
