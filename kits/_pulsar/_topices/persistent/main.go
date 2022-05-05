package main

import (
	"context"
	"log"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

func main() {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: "pulsar://localhost:6650",
	})
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	// {persistent|non-persistent}://tenant/namespace/topic
	// 持久化topic
	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: "persistent://chensy/namespace/simple01",
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
