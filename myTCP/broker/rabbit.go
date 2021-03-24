/*
реализация брокера сообщений на основе кролика
*/
package broker

import (
	"log"
	"sync"

	"github.com/streadway/amqp"
)

// Rabbit реализация работы с кроликом
type Rabbit struct {
	ConnPath string
}

// Publish постановка задачи в очередь
func (r *Rabbit) Publish(task []byte) {
	log.Println("start send")
	conn, err := amqp.Dial(r.ConnPath)
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

// Consume разбор сообщений
func (r *Rabbit) Consume() {
	conn, err := amqp.Dial(r.ConnPath)
	if err != nil {
		log.Println("ошибка инициализации коннекта к кролику")
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
	// по умолчанию кролик четные таски отдаёт первому воркеру а нечётные
	// второму воркеру поэтому если один будет часто занят то другой
	// часто свободен и неплохо бы перераспределить нагрузку более честно
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

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for d := range msgs {
			log.Printf("обработан коннект: %s", d.Body)
			d.Ack(false) // сообщение удалится только когда мы его точно обработаем
		}
		wg.Done()
	}()
	wg.Wait()
}
