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
		_, err = conn.Write([]byte("SET ayush gupta 250000000000000000"))
		if err != nil {
			fmt.Println("Failed to write to server:", err)
			return
		}

		time.Sleep(2 * time.Second)
		conn.Write([]byte("GET ayush"))
		for {
			buf := make([]byte, 2048)

			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println("Failed to read from server:", err)
				return
			}
			fmt.Println(string(buf[:n]))
		}

	}()
	server := NEWServer(listOpts, cache.NewCache())
	server.Start()
}
