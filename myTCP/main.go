package main

import (
	"runtime"
	"sync"
	"time"

	"github.com/AlexanderKorovaev/TCP/myTCP/core"
)

//1) заменять хэндлеры патерном билдер
// Пока не знаю как решить эту проблему, может потом придумаю
// пока что просто оставляю хендлер возвращающий ошибку и пользователь сам
// должен переопределить его +
//2) проработать работу шатдаун +
//3) потестить работу многопоточности и применить производитьель потребитель
// прочитать статью из конкьюренси и воспроизвести тут
//4) лишние зарпосы помещать в брокер сообщений

func main() {
	srv := core.Server{
		Port:        ":2000",
		IdleTimeout: 10 * time.Second,
	}
	var wg sync.WaitGroup
	// ограничим колличество воркеров возможностями процессора
	maxWorkers := runtime.NumCPU()
	wg.Add(maxWorkers)
	go srv.ListenAndServe(&wg, maxWorkers)
	time.Sleep(5 * time.Second)
	srv.Shutdown(&wg)
	wg.Wait()
}
