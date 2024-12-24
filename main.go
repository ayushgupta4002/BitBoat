package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/ayushgupta4002/bitboat/cache"
)

func main() {
	fmt.Println("hello world")
	listOpts := ServerOpts{
		listenAddr: "localhost:8080",
		isAdmin:    true,
	}

	go func() {
		time.Sleep(5 * time.Second)
		conn, err := net.Dial("tcp", "localhost:8080")
		if err != nil {
			log.Fatal("server cannot request")
		}
		defer conn.Close()
		conn.Write([]byte("Set ayush gupta 2500"))
	}()
	server := NEWServer(listOpts, cache.NewCache())
	server.Start()
}
