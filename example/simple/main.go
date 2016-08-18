package main

import (
	"log"

	"github.com/Akagi201/redface/server"
	"github.com/mediocregopher/radix.v2/redis"
)

var pongSS = redis.NewRespSimple("PONG")

func pingHandler(args []string) (interface{}, error) {
	return pongSS, nil
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic: %v\n", err)
		}
	}()

	srv, err := server.NewServer("tcp", "127.0.0.1", 6379)
	if err != nil {
		panic(err)
	}

	srv.Handle("ping", pingHandler)

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
