package main

import (
	"runtime"
	"sync"
	"time"

	"github.com/AlexanderKorovaev/TCP/myTCP/server/core"
)

// попробовать паттерн когда я подкладываю любой брокер сообщений вместо кролика

func main() {
	srv := core.Server{
		Port:        ":2000",
		IdleTimeout: 5 * time.Second,
	}
	var wg sync.WaitGroup
	// ограничим колличество воркеров возможностями процессора
	maxWorkers := runtime.NumCPU()
	wg.Add(maxWorkers)
	go srv.ListenAndServe(&wg, maxWorkers)
	time.Sleep(60 * time.Second)
	srv.Shutdown(&wg)
	wg.Wait()
}
