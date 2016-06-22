package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

var curConnNum int = 0

/* test server Goroutine */
func BandwidthTest(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 10)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			return
		}
		//send reply
		/*
			_, err = conn.Write(buf)
			if err != nil {
				fmt.Println("Error send reply:", err.Error())
				return
			}
		*/
	}
}

/* initial listener and run */
func BandwidthServer() {
	listener, err := net.Listen("tcp", "0.0.0.0:8088")
	if err != nil {
		fmt.Println("error listening:", err.Error())
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("running ...\n")

	connChan := make(chan net.Conn)
	connChangChan := make(chan int)

	go func() {
		for connChange := range connChangChan {
			curConnNum += connChange
		}
	}()

	go func() {
		for _ = range time.Tick(1e9 * 30) {
			fmt.Printf("Cur conn num: %v\n", curConnNum)
		}
	}()
	for i := 0; i < MaxConnNum; i++ {
		go func() {
			for conn := range connChan {
				connChangChan <- 1
				BandwidthTest(conn)
				connChangChan <- -1
			}
		}()
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			println("Error accept:", err.Error())
			return
		}
		connChan <- conn
	}
}
