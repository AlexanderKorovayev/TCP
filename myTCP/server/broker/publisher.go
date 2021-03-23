package broker

import (
	"log"

	"github.com/streadway/amqp"
)

// Publish постановка задачи в очередь
func Publish(task []byte) {
	log.Println("start send")
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Printf("error %v", err)
		panic(err)
	}

	// открываем коннект, это более легковестная сущность по сравнению
	// с коннектом, поэтому лучше делать много каналов, чем много конектов
	ch, err := conn.Channel()
	if err != nil {
		log.Println(err)
	}
	defer ch.Close()

	// создаём очередь
	q, err := ch.QueueDeclare(
		"Test_Queue", // name
		true,         // durable - не удалять сообщения если кролик упадёт
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	// информация об очереди
	log.Println(q)

	if err != nil {
		log.Println(err)
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         task,
		},
	)

	if err != nil {
		log.Println(err)
	}
	log.Printf("sended")
}
