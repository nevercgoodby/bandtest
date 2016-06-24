package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	BandServer    string
	HttpServer    string
	SyncCycle     int
	ConnTimeout   int
	MaxTestClient int
}

var (
	curConnNum      int
	BandServerAddr  string
	BandServerPort  uint16
	HttpServerAddr  string
	Sync_Cycle_Time int
	ConnTimeOut     int
	MaxTestClient   int
)

func ConfigRead(config_file string) {
	r, err := os.Open(config_file)
	if err != nil {
		log.Fatalln(err)
	}
	decoder := json.NewDecoder(r)
	var c Config
	err = decoder.Decode(&c)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(c)
	BandServerAddr = c.BandServer
	HttpServerAddr = c.HttpServer
	Sync_Cycle_Time = c.SyncCycle
	ConnTimeOut = c.ConnTimeout
	MaxTestClient = c.MaxTestClient
	p := strings.Split(BandServerAddr, ":")
	if len(p) < 2 {
		log.Fatalln(err)
	}
	port, err := strconv.Atoi(p[1])
	BandServerPort = uint16(port)

}
