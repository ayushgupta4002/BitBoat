package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/ayushgupta4002/bitboat/cache"
	"github.com/ayushgupta4002/bitboat/client"
)

func main() {

	var listenAddress = flag.String("listenaddr", "localhost:8080", "server listen address")
	var adminAddress = flag.String("adminaddr", "", "admin listen address")
	flag.Parse()
	listOpts := ServerOpts{
		listenAddr: *listenAddress,
		isAdmin:    len(*adminAddress) == 0,
		adminAddr:  *adminAddress,
	}

	go func() {
		time.Sleep(2 * time.Second)
		c, err := client.NewClient(*listenAddress, client.ClientOpts{})
		if err != nil {
			log.Fatal("client cannot request", err)
		}
		for i := 0; i < 5; i++ {
			err := c.Set(context.Background(), []byte("ayush"), []byte("gupta"), 200000000)
			if err != nil {
				log.Fatal("client cannot request", err)
				continue
			}
		}
		val, err := c.Get(context.Background(), []byte("ayush"))
		if err != nil {
			log.Fatal("client cannot request", err)
		}
		log.Println(string(val))

		c.Close()

	}()

	server := NEWServer(listOpts, cache.NewCache())
	server.Start()
}

// go func() {
// 	time.Sleep(5 * time.Second)
// 	conn, err := net.Dial("tcp", *listenAddress)
// 	if err != nil {
// 		log.Fatal("server cannot request")
// 	}
// 	_, err = conn.Write([]byte("SET ayush gupta 250000000000000000"))
// 	if err != nil {
// 		fmt.Println("Failed to write to server:", err)
// 		return
// 	}

// time.Sleep(2 * time.Second)
// conn.Write([]byte("GET ayush"))
// for {
// 	buf := make([]byte, 2048)

// 	n, err := conn.Read(buf)
// 	if err != nil {
// 		fmt.Println("Failed to read from server:", err)
// 		return
// 	}
// 	fmt.Println(string(buf[:n]))
// }

// }()
