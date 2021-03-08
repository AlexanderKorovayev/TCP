package main

import (
	"time"

	"trueTCP/shouter"
)

func main() {
	srv := shouter.Server{
		Addr:         ":2000",
		IdleTimeout:  10 * time.Second,
		MaxReadBytes: 1000,
	}
	go srv.ListenAndServe()
	time.Sleep(10 * time.Second)
	//srv.Shutdown()
	select {}
}
