package main

import (
	"log"
	"net"

	"github.com/Akagi201/redface/server"
	"github.com/mediocregopher/radix.v2/redis"
)

var pongSS = redis.NewRespSimple("PONG")

func pingHandler(conn net.Conn, args []string) (interface{}, error) {
	return pongSS, nil
}

func quitHandler(conn net.Conn, args []string) (interface{}, error) {
	conn.Close()
	return nil, nil
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic: %v\n", err)
		}
	}()

	srv, err := server.New(6389)
	if err != nil {
		panic(err)
	}

	srv.Handle("ping", pingHandler)

	srv.Handle("quit", quitHandler)

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
