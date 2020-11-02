package main

import (
	"github.com/Shopify/sarama"
	"os"
)

type producer struct {
	producer sarama.AsyncProducer
	msg      chan *sarama.ProducerMessage
	close    chan os.Signal
}

func (p *producer) input() chan<- *sarama.ProducerMessage {
	return p.msg
}
func (p *producer) start() {
	go p.produceMessage()
	go p.readError()
}
func (p *producer) stop() {
	p.close <- os.Interrupt
}
func (p *producer) produceMessage() {
	stopping := false

	for {
		select {
		case <-p.close:
			stopping = true
		case msg := <-p.msg:
			if stopping {
				logger.Errorf("生产者已关闭，无法发送消息，请检查关机顺序是否合理，消息[topic=%s, value=%s]", msg.Topic, msg.Value)
			} else {
				p.producer.Input() <- msg
			}
		}
	}
}
func (p *producer) readError() {
	for err := range p.producer.Errors() {
		logger.Errorf("生产者错误: %s, 消息[topic=%s, value=%s]", err.Err, err.Msg.Topic, err.Msg.Value)
	}
}
