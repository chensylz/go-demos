package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func main() {
	conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
	defer conn.Close()
	ch, _ := conn.Channel()
	defer ch.Close()
	//创建交换机
	if err := ch.ExchangeDeclare(
		"rabbit_exchange", //name
		"topic",           //type
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		fmt.Println("ExchangeDeclare Err =", err)
		return
	}
	count := 0
	for {
		count++
		body := fmt.Sprintf("Hello World! %d", count)
		err := ch.Publish(
			"rabbit_exchange",      // exchange
			"rabbit.chensy.chenjp", // routing key
			false,                  // mandatory
			false,                  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		log.Printf(" [x] Sent %s", body)
		if err != nil {
			log.Printf(" [x] err %s", err.Error())
		}
		time.Sleep(1 * time.Second)
	}
}
