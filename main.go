package main

import (
	"fmt"
	//"io/ioutil"
	"math/rand"

	"log"
	"net/http"
	"os"
)

const (
	MaxConnNum = 5
)

var ServerPort int

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func httpServer(addr string) {
	Router()
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func main() {
	buf := EncodeHeartBeatRequest(curConnNum)
	fmt.Println(buf.Bytes())
	go KeepAlive()
	go BandwidthServer()
	for i := 0; i < 5; i++ {
		fmt.Printf("%d ", rand.Int()%16)
	}
	fmt.Println()

	//init Bandwidth test server and start listen
	//go BandwidthServer()

	//get internal ip address
	//get_internal()

	/*
		Port := flag.Int("Server Port", 8088, "Bandwidth Port")
		flag.Parse()
		ServerPort = *Port
		httpServer(":8080")
	*/
}
