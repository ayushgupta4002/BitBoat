package main

import (
	"context"
	"flag"
	"fmt"
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
	//simulating a demo client
	go func() {
		// this function will run after 5 seconds and will send request to Admin Server to set, get and delete the key
		// we can instead send the request to the client also but then we can only GET and HAS the key

		time.Sleep(5 * time.Second)
		// if the server starts and the initiator is admin then wait for 5 seconds and then send request to admin
		// meanwhile in those 5 seconds I can connect a subscriber to the server to check if the server is working fine
		// JUST FOR SIMULATING THE BEHAVIOUR
		if listOpts.isAdmin {
			sendData(*listenAddress)
		}

	}()

	server := NEWServer(listOpts, cache.NewCache())
	server.Start()
}

// simulating a demo client
func sendData(listenAddress string) {
	for i := 0; i < 2; i++ {
		go func() {
			c, err := client.NewClient(listenAddress, client.ClientOpts{})
			if err != nil {
				log.Fatal("client cannot request", err)
			}
			err = c.Set(context.Background(), []byte(fmt.Sprintf("key_%d", i)), []byte(fmt.Sprintf("Val_%d", i)), 2000000000)
			if err != nil {
				log.Fatal("client cannot set", err)
			}

			_, err = c.Get(context.Background(), []byte(fmt.Sprintf("key_%d", i)))
			if err != nil {
				log.Fatal("client cannot get", err)
			}
			// log.Println(string(val))
			err = c.Delete(context.Background(), []byte(fmt.Sprintf("key_%d", i)))
			if err != nil {
				log.Fatal("client could not delete the key", err)
			}
			// err = c.Has(context.Background(), []byte(fmt.Sprintf("key_%d", i)))
			// if err != nil {
			// 	log.Fatal("client could not delete the key", err)
			// }

			c.Close()
		}()
		time.Sleep(2 * time.Second)

	}
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
