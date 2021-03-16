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

	// Let's start by opening a channel to our RabbitMQ instance
	// over the connection we have already established
	ch, err := conn.Channel()
	if err != nil {
		log.Println(err)
	}
	defer ch.Close()

	// with this channel open, we can then start to interact
	// with the instance and declare Queues that we can publish and
	// subscribe to
	q, err := ch.QueueDeclare(
		"TestQueue",
		false,
		false,
		false,
		false,
		nil,
	)
	// We can print out the status of our Queue here
	// this will information like the amount of messages on
	// the queue
	log.Println(q)
	// Handle any errors if we were unable to create the queue
	if err != nil {
		log.Println(err)
	}

	// attempt to publish a message to the queue!
	err = ch.Publish(
		"",
		"TestQueue",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("Hello World"),
		},
	)

	if err != nil {
		log.Println(err)
	}
	log.Println("Successfully Published Message to Queue")
}
