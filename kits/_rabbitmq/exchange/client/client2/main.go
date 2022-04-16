package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func main() {
	conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
	defer conn.Close()
	ch, _ := conn.Channel()
	defer ch.Close()
	//创建交换机
	if err := ch.ExchangeDeclare(
		"rabbit_exchange", // name
		"topic",           // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	); err != nil {
		fmt.Println("ExchangeDeclare Err =", err)
		return
	}
	q, _ := ch.QueueDeclare(
		"",    //队列名
		false, //持久的
		false, // delete when unused
		true,  //独占的
		false,
		nil,
	)
	//交换机绑定队列
	if err := ch.QueueBind(q.Name, "*.chensy.*", "rabbit_exchange", false, nil); err != nil {
		fmt.Println("QueueBind Err =", err)
		return
	}
	msgs, _ := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	<-forever
}
