package main

import (
	"time"

	"github.com/AlexanderKorovaev/TCP/trueTCP/shouter"
)

// проверил что работает ограничение по времени
// и ограничение по размеру данных
// надо добиться правильной работы

func main() {
	srv := shouter.Server{
		Port:         ":2000",
		IdleTimeout:  10 * time.Second,
		MaxReadBytes: 0,
	}
	srv.ListenAndServe()
	//srv.Shutdown()
}
