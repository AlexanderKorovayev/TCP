package broker

import (
	"bytes"
	"log"
	"time"

	"github.com/streadway/amqp"
)

// Consumer разбор сообщений
func Consumer() {
	log.Println("Go RabbitMQ Tutorial")
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
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

	// отвечает за то что бы воркеры не проставивали
	// по умолчанию кролик четные таски отдаёт первому воркеру а нечётные второму воркеру
	// поэтому если один будет часто занят то другой часто свободен
	// и неплохо бы перераспределить нагрузку более честно
	err = ch.Qos(
		1, // prefetch count - не давать одному воркеру больше чем одно
		// сообщение, в таком случае кролик отдаст свободному
		0,     // prefetch size
		false, // global
	)

	if err != nil {
		log.Println(err)
	}

	msgs, err := ch.Consume(
		"Test_Queue", // queue
		"",           // consumer
		false,        // auto-ack - не удаляем автоматически
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			log.Printf("Done")
			d.Ack(false) // сообщение удалится только когда мы его точно обработаем
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
