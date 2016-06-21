package main

import (
	_ "bytes"
	"fmt"
	"net"
	"net/url"
	_ "os"
	"strconv"
	"time"
)

var (
	DefalutTimeout = 25 * time.Second
	MaxClient      = 70000
	clientNum      = 0
	msg_nums       = 0
)

func sockConn(daemon string, timeout time.Duration) (net.Conn, error) {
	daemonURL, err := url.Parse(daemon)
	//fmt.Printf("daemon url %v %v \n", daemonURL, daemonURL.Scheme)
	if err != nil {
		return nil, fmt.Errorf("could not parse url %q: %v", daemon, err)
	}

	var c net.Conn
	switch daemonURL.Scheme {
	case "unix":
		return net.DialTimeout(daemonURL.Scheme, daemonURL.Path, timeout)
	case "tcp":
		return net.DialTimeout(daemonURL.Scheme, daemonURL.Host, timeout)
	default:
		return c, fmt.Errorf("unknown scheme %v (%s)", daemonURL.Scheme, daemon)
	}
}

func sendData(socket net.Conn, n int, msg_num_chan chan int) {
	buf := make([]byte, 10)
	for {
		_, err := socket.Write([]byte(strconv.Itoa(n)))
		if err != nil {
			fmt.Printf("Error reading:%s\n", err.Error())
			clientNum--
			return
		}
		//send reply
		_, err = socket.Read(buf)
		fmt.Printf("client %v\n", string(buf))
		if err != nil {
			fmt.Printf("Error send reply:%s\n", err.Error())
			clientNum--
			return
		}
		//time.Sleep(1 * time.Second)
		msg_num_chan <- 1
		buf = buf[0:]
	}
}

func connectServerTcp(msg_num_chan chan int) {
	for i := 1; i <= MaxClient; i++ {
		socket, err := sockConn("tcp://127.0.0.1:8088", DefalutTimeout)
		if err != nil {
			fmt.Printf("connect error:%s\n", err)
		} else {
			clientNum++
			go sendData(socket, i, msg_num_chan)
		}

	}
}

func clinet_test() {

	msg_num_chan := make(chan int)
	connectServerTcp(msg_num_chan)

	go func() {
		for msg_num := range msg_num_chan {
			msg_nums += msg_num
		}
	}()

	time.Sleep(10 * time.Second)
	fmt.Printf("rec msg %d \n", msg_nums)
}
