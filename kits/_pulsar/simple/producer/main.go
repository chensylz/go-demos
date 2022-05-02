package main

import (
	"context"
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

	for {
		_, err = producer.Send(context.Background(), &pulsar.ProducerMessage{
			Payload: []byte("hello"),
		})
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second)
	}
}
