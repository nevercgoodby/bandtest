package main

import (
	"fmt"
	//"math/rand"
	"time"
	//"io/ioutil"
	"bytes"
	"encoding/binary"
	"net"
)

const (
	Sync_Cycle_Time = 30
)

func EncodeHeartBeatRequest(curConnNum int) *bytes.Buffer {
	buf := new(bytes.Buffer)
	var hbeatreq HeartbeatRequest
	hbeatreq.Length = uint32(binary.Size(hbeatreq))
	hbeatreq.Cmd = uint32(0x00000024)
	hbeatreq.Seq = uint32(1)
	hbeatreq.Loads[6] = uint16(curConnNum)
	hbeatreq.Ports[6] = uint16(BandWidthPort)
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

func sendKeepAliveMsg(conn net.Conn) int {
	conn.Write([]byte("a"))
	return 0
}

func sendMsg(conn net.Conn, msg []byte) error {
	var buf []byte
	buf = make([]byte, 10, 10)
	_, err := conn.Write(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return err
	}
	return nil
}
func connectUnixServer(localpath string) (net.Conn, error) {
	fmt.Println("Hi, I will connect!")
	unixAddr, err := net.ResolveUnixAddr("unix", localpath)
	if err != nil {
		fmt.Println("Fatal error: ", err)
	}
	fmt.Println("unixaddr:", unixAddr)
	uconn, err := net.DialUnix("unix", nil, unixAddr)
	if err != nil {
		fmt.Println("Fatal error: ", err)
		uconn.Close()
		return uconn, err
	}
	return uconn, nil
}

func connectServer(host string, port uint16) (net.Conn, error) {
	fmt.Println("Hi, I will connect!")
	server := fmt.Sprintf("%s:%d", host, port)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Println("Fatal error: ", err)
	}
	fmt.Println("tcpaddr:", tcpAddr.IP, tcpAddr.Port)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("Fatal error: ", err)
		return conn, err
	}
	return conn, nil

}

func KeepAlive() {
	reqChan := make(chan []byte)
	//	resChan := make(chan []byte)
	connChan := make(chan net.Conn)
	pipelineChan := make(chan int)
	//Conn Producer
	go func() {
		for {
			//if need reconnect
			<-pipelineChan
			conn, err := connectServer("172.16.10.208", 9869)
			if err != nil {
				time.Sleep( /*rand.Intn(10) * time.Second*/ 10 * 1e9)
				fmt.Println("Fatal error c:", err)
				//need reconnect
				pipelineChan <- 1
			} else {
				//connect success

				connChan <- conn
			}
			fmt.Println("Okay, I connected!")
		}
		fmt.Println("Connect cen goroutine Die!")
	}()
	//Msg Producer
	go func() {
		for {
			buf := EncodeHeartBeatRequest(curConnNum)
			reqChan <- buf.Bytes()
			//time.Sleep(Sync_Cycle_Time * time.Second)
			time.Sleep(1 * time.Second)
		}
	}()

	//Msg Consumer (send reqest and recive response)
	go func() {
		pipelineChan <- 1

		for conn := range connChan {
			fmt.Println("conn Chan len:", len(connChan))
			for msgBytes := range reqChan {

				//send
				fmt.Println("send len:", len(msgBytes))
				wn, err := conn.Write(msgBytes)
				if err != nil {
					conn.Close()
					//need reconnect
					fmt.Println("Error s, connection closed!")
					pipelineChan <- 1
					break
				}
				fmt.Println("Send status:", wn)
				conn.SetReadDeadline((time.Now().Add(time.Second * 10)))

				//recv
				resBytes := make([]byte, 1024)
				fmt.Println("Hi, I will Read response!")

				//				select {
				//					case rn, err := conn.Read(resBytes):
				//					case
				//				}
				rn, err := conn.Read(resBytes)
				if rn == 0 && err != nil {
					conn.Close()
					//need reconnect
					fmt.Println("Read error! Read len:", rn, err)
					pipelineChan <- 1
					break
				}
				fmt.Println("######Recv status:", rn)
			}
		}
		fmt.Println("Msg Send goruntine quit!")
	}()
}
