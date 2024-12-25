package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/ayushgupta4002/bitboat/cache"
)

func main() {
	conn, err := net.Dial("tcp", ":3000")
	if err != nil {
		panic(err)
	}

	// Send a SET command
	_, err = conn.Write([]byte("SET ayush 100 250000000000"))
	if err != nil {
		panic(err)
	}
	// Wait for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	fmt.Println("Waiting for interrupt signal...")
	<-signalChan // Blocks until an interrupt signal is received
	fmt.Println("Shutting down gracefully")

	var listenAddress = flag.String("listenaddr", "localhost:8080", "server listen address")
	var adminAddress = flag.String("adminaddr", "", "admin listen address")
	flag.Parse()
	listOpts := ServerOpts{
		listenAddr: *listenAddress,
		isAdmin:    len(*adminAddress) == 0,
		adminAddr:  *adminAddress,
	}

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
