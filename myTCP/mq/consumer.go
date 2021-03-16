package main

import (
	"log"

	"github.com/streadway/amqp"
)

func main() {
	log.Println("Go RabbitMQ Tutorial")
	conn, err := amqp.Dial("amqp://guest:guest@localhost:8080/")
	if err != nil {
		log.Println("Failed Initializing Broker Connection")
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Println(err)
	}
	defer ch.Close()

	if err != nil {
		log.Println(err)
	}

	msgs, err := ch.Consume(
		"TestQueue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Recieved Message: %s\n", d.Body)
		}
	}()

	log.Println("Successfully Connected to our RabbitMQ Instance")
	log.Println(" [*] - Waiting for messages")
	<-forever
}
