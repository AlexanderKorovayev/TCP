/*
сервис логирования коннектов
*/
package main

import (
	"github.com/AlexanderKorovaev/TCP/myTCP/broker"
)

func main() {
	rabbit := broker.Rabbit{ConnPath: "amqp://guest:guest@localhost:5672/"}
	rabbit.Consume()
}
