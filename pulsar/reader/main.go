package main

import (
	"context"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"log"
	"time"
)

func main() {
	client, err := pulsar.NewClient(pulsar.ClientOptions{URL: "pulsar://localhost:6650"})
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	reader, err := client.CreateReader(pulsar.ReaderOptions{
		Topic:          "my-topic",
		StartMessageID: pulsar.EarliestMessageID(),
	})
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	ctx, cancel := context.WithCancel(context.TODO())
	time.AfterFunc(time.Minute, cancel)
	for {
		msg, err := reader.Next(ctx)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(msg.ID(), string(msg.Payload()))
	}
}
