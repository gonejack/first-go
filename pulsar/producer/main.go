package main

import (
	"context"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"log"
	"time"
)

func main() {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: "pulsar://localhost:6650",
	})
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: "my-topic",
	})
	if err != nil {
		log.Fatal(err)
	}

	defer producer.Close()

	var i = 0
	for i < 100 {
		i++

		message := fmt.Sprintf("这是一条消息: %d", time.Now().Unix())
		_, err = producer.Send(context.Background(), &pulsar.ProducerMessage{
			Payload: []byte(message),
		})

		if err == nil {
			fmt.Println(message)
		} else {
			fmt.Println("Failed to publish message", err)
		}

		time.Sleep(time.Second)
	}
}
