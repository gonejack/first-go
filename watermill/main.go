// Sources for https://watermill.io/docs/getting-started/
package main

import (
	"context"
	"log"
	"time"

	"github.com/Shopify/sarama"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
)

var (
	brokers = []string{
		"192.168.11.10:9093",
		"192.168.11.10:9094",
		"192.168.11.10:9095",
	}
	topic = "test_topic"
)

func main() {
	config := kafka.DefaultSaramaSubscriberConfig()
	config.Version, _ = sarama.ParseKafkaVersion("0.10.2.0")

	sub, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:               brokers,
			Unmarshaler:           kafka.DefaultMarshaler{},
			OverwriteSaramaConfig: config,
			ConsumerGroup:         "test_group",
		},
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		panic(err)
	}

	output, err := sub.Subscribe(context.Background(), topic)
	if err != nil {
		panic(err)
	}
	go consume(output)

	config.Producer.Return.Successes = true
	pub, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:               brokers,
			Marshaler:             kafka.DefaultMarshaler{},
			OverwriteSaramaConfig: config,
		},
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		panic(err)
	}

	publish(pub)
}

func publish(publisher message.Publisher) {
	for {
		msg := message.NewMessage(watermill.NewUUID(), []byte("Hello, world!"))
		if err := publisher.Publish(topic, msg); err != nil {
			panic(err)
		}
		time.Sleep(time.Second)
	}
}

func consume(messages <-chan *message.Message) {
	for msg := range messages {
		log.Printf("received message: %s, payload: %s", msg.UUID, string(msg.Payload))

		// we need to Acknowledge that we received and processed the message,
		// otherwise, it will be resent over and over again.
		msg.Ack()
	}
}
