package main

import (
	"fmt"
	//"io/ioutil"
	"bytes"
	"encoding/binary"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func EncodeHeartBeatRequest() *bytes.Buffer {
	buf := new(bytes.Buffer)
	var hbeatreq HeartbeatRequest
	hbeatreq.Length = uint32(binary.Size(hbeatreq))
	hbeatreq.Cmd = uint32(0x00000024)
	hbeatreq.Seq = uint32(1)
	binary.Write(buf, binary.BigEndian, hbeatreq)
	// var msg []byte
	// msg = buf.Bytes()
	// fmt.Println(msg)
	return buf
}

func DecodeHeartBeatResponse(buf bytes.Buffer) int {
	data := make([]byte, 6)
	binary.BigEndian.PutUint16(data, 0x0102)
	binary.BigEndian.PutUint32(data[2:6], 0x03040506)
	length := binary.LittleEndian.Uint16(data[:2])
	cmd := binary.LittleEndian.Uint16(data[:4])
	seq := binary.LittleEndian.Uint16(data[:6])
	fmt.Println("data:", data, length, cmd, seq)
	return 0
}

// Convert uint to net.IP
func inet_ntoa(ipnr int64) net.IP {
	var bytes [4]byte
	bytes[0] = byte(ipnr & 0xFF)
	bytes[1] = byte((ipnr >> 8) & 0xFF)
	bytes[2] = byte((ipnr >> 16) & 0xFF)
	bytes[3] = byte((ipnr >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

// Convert net.IP to int32
func inet_aton2(ipnr string) int32 {
	bits := strings.Split(ipnr, ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int32

	sum += int32(b0) << 24
	sum += int32(b1) << 16
	sum += int32(b2) << 8
	sum += int32(b3)

	return sum
}

func inet_aton(ipnr net.IP) int32 {
	bits := strings.Split(ipnr.String(), ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int32

	sum += int32(b0) << 24
	sum += int32(b1) << 16
	sum += int32(b2) << 8
	sum += int32(b3)

	return sum
}

func get_internal() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops:" + err.Error())
		os.Exit(1)
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				os.Stdout.WriteString(ipnet.IP.String() + "\n")
			}
		}
	}
}

func sendKeepAliveMsg(conn net.Conn) int {
	conn.Write([]byte("a"))
	return 0
}

func sendMsg(conn net.Conn, msg []byte) {
	var buf []byte
	buf = make([]byte, 10, 10)
	_, err := conn.Write(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}
}

func connectServer(host string, port uint16) {
	server := "127.0.0.1:8080"

	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	for {
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
			os.Exit(1)
		}
		defer conn.Close()
		fmt.Println("Connect server success!")
		for ret := 0; ret == 0; {
			ret = sendKeepAliveMsg(conn)
		}
		time.Sleep(3 * 1e9)
	}
}
