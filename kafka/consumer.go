package main

import (
	cluster "github.com/bsm/sarama-cluster"
	"sync"
)

// 消费组封装
type consumer struct {
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

func (c *consumer) start() (err error) {
	config := cluster.NewConfig()

	config.Consumer.Return.Errors = true   // 设置对外报错
	config.Net.SASL.Enable = c.SASL.Enable // 是否启用密码
	config.Net.SASL.User = c.SASL.User
	config.Net.SASL.Password = c.SASL.Password

	c.consumer, err = cluster.NewConsumer(c.Addrs, c.GroupId, []string{c.Topic}, config)

	if err == nil {
		logger.Infof("消费组[topic=%s, groupId=%s]就绪", c.Topic, c.GroupId)

		go c.readMessage()
		go c.readError()
	} else {
		logger.Errorf("新建消费者[topic=%s, groupId=%s]失败: %s", c.Topic, c.GroupId, err)
	}

	return
}
func (c *consumer) stop() (err error) {
	err = c.consumer.Close()

	if err == nil {
		logger.Infof("消费组[topic=%s, groupId=%s]已关闭", c.Topic, c.GroupId)
	} else {
		logger.Errorf("关闭消费组[topic=%s, groupId=%s]出错: %s", c.Topic, c.GroupId, err)
	}

	return
}
func (c *consumer) readMessage() {
	for msg := range c.consumer.Messages() {
		c.mu.RLock()
		for _, subscriber := range c.subscribers {
			err := subscriber(msg.Topic, msg.Value)
			if err != nil {
				logger.Errorf("消费组[topic=%s, groupId=%s]订阅者出错: %s", c.Topic, c.GroupId, err)
			}
		}
		c.mu.RUnlock()
		c.consumer.MarkOffset(msg, "")
	}
}
func (c *consumer) readError() {
	for err := range c.consumer.Errors() {
		logger.Errorf("消费发生错误: %s", err)
	}
}
func (c *consumer) addSubscriber(handler func(topic string, msg []byte) error) *consumer {
	c.mu.Lock()
	c.subscribers = append(c.subscribers, handler)
	c.mu.Unlock()

	return c
}
