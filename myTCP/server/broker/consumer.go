package broker

import (
	"log"
	"sync"

	"github.com/streadway/amqp"
)

// шаблонный метод и инверсия зависимостей
//https://refactoring.guru/ru/design-patterns/template-method/go/example
//https://disk.yandex.ru/client/disk/Teory

// Consumer разбор сообщений
func Consumer() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
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
