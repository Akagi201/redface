package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"sync"

	"github.com/Akagi201/redface/server"
	"github.com/mediocregopher/radix.v2/redis"
)

var okSS = redis.NewRespSimple("OK")
var pongSS = redis.NewRespSimple("PONG")

// newRespErr is a convenience method for making Resps to wrap errors
func newRespErr(s string) *redis.Resp {
	err := fmt.Errorf("ERR %v", s)
	return redis.NewResp(err)
}

func pingHandler(conn net.Conn, args []string) (interface{}, error) {
	return pongSS, nil
}

var (
	mu        sync.RWMutex
	redisHash = make(map[string]string)
)

func setHandler(conn net.Conn, args []string) (interface{}, error) {
	if len(args) < 2 {
		return newRespErr("missing args"), nil
	}

	mu.Lock()
	defer mu.Unlock()
	redisHash[args[0]] = args[1]

	return okSS, nil
}

func getHandler(conn net.Conn, args []string) (interface{}, error) {
	if len(args) < 1 {
		return newRespErr("missing args"), nil
	}

	mu.RLock()
	defer mu.RUnlock()

	return redis.NewRespSimple(redisHash[args[0]]), nil
}

func quitHandler(conn net.Conn, args []string) (interface{}, error) {
	conn.Close()
	return nil, nil
}

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
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

	srv.Handle("set", setHandler)

	srv.Handle("get", getHandler)

	srv.Handle("quit", quitHandler)

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
