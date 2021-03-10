package main

import (
	"time"

	"github.com/AlexanderKorovaev/TCP/myTCP/core"
)

//1) заменять хэндлеры патерном билдер
// нам приходят данные на вход и на основе этих данных мы применяем выбор обработчика
//2) проработать работу шатдаун
//3) потестить работу многопоточности и применить производитьель потребитель
//4) лишние зарпосы помещать в брокер сообщений

func main() {
	srv := core.Server{
		Port:        ":2000",
		IdleTimeout: 10 * time.Second,
	}
	srv.ListenAndServe()
	//srv.Shutdown()
}
