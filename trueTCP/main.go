package main

import (
	"time"

	"github.com/AlexanderKorovaev/TCP/trueTCP/shouter"
)

// проверил что работает ограничение по времени
// и ограничение по размеру данных
// надо добиться правильной работы

// запилю универсальный настраиваемый сервак
// мы сделаем 4 паралельных обработчика а остальные пойдут в брокер сообщений

func main() {
	srv := shouter.Server{
		Port:        ":2000",
		IdleTimeout: 10 * time.Second,
	}
	srv.ListenAndServe()
	//srv.Shutdown()
}
