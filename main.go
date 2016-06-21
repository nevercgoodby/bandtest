package main

import (
	"fmt"
	//"io/ioutil"

	"log"
	"net/http"
	"os"
)

const (
	MaxConnNum = 5
	//MaxConnNum
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
	EncodeHeartBeatRequest()
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
