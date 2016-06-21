package main

import (
	"flag"
	"fmt"
	//"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	MaxConnNum = 5
	//MaxConnNum
)

var HeartbeatRequest = []interface{}{
	uint32(100),        //len
	uint32(0x00000024), //cmd 80000024 bandtest 0X00000040 #define OPT_P2PT     0X00000040
	uint32(0),          //seq
	uint32(0x01),       //apptype
	[]uint32(),

	//appname [32]uint8,
	//loads [32]uint16,
	//ports [32]uint16
}

type Head struct {
	InodeCount     uint32 //  0:4
	BlockCount     uint32 //  4:8
	Unknown1       uint32 //  8:12
	Unknown2       uint32 // 12:16
	Unknown3       uint32 // 16:20
	FirstBlock     uint32 // 20:24
	BlockSize      uint32 // 24:28
	Unknown4       uint32 // 28:32
	BlocksPerGroup uint32 // 32:36
	Unknown5       uint32 // 36:40
	InodesPerBlock uint32 // 40:44
}

var ServerPort int

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func encodeHeartBeatMsg() {
	data := make([]byte, 6)
	binary.BigEndian.PutUint16(data, 0x1011)
	binary.BigEndian.PutUint32(data[2:6], 0x12131415)
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

func EchoHander(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.RemoteAddr)
	w.Write([]byte("abc\n"))
	return
}

func CmdHander(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.RemoteAddr)
	w.Write([]byte("CMD\n"))
	return
}

func BandwidthHander(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.RemoteAddr)
	w.Write([]byte("Bandwidth\n"))
	return
}

func StatHander(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.RemoteAddr)
	w.Write([]byte("Stat\n"))
	return
}

/* DataHander */
func DataHander(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.RemoteAddr)
	w.Write([]byte("Data\n"))
	return
}

/* DownloadHander */
func DownloadHander(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.RemoteAddr)
	w.Write([]byte("Download\n"))
	return
}

/* UploadApplyHander */
func UploadApplyHander(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	addr := strings.Split(r.Host, ":")
	ip := addr[0]

	bits := strings.Split(ip, ".")
	fmt.Println(bits, bits[0], bits[1], bits[2])
	//var buffer []byte
	buffer := make([]byte, 7, 7)
	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])
	buffer[0] = byte(b0 & 0xFF)
	buffer[1] = byte(b1 & 0xFF)
	buffer[2] = byte(b2 & 0xFF)
	buffer[3] = byte(b3 & 0xFF)
	buffer[4] = byte(ServerPort >> 8)
	buffer[5] = byte(ServerPort)

	fmt.Println(path, "host:", ip, r.Host, r.RemoteAddr, r.URL)
	//w.Write([]byte(buffer))
	w.Write(buffer)
}

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
		_, err = conn.Write(buf)
		if err != nil {
			fmt.Println("Error send reply:", err.Error())
			return
		}
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

	var curConnNum int = 0
	connChan := make(chan net.Conn)
	connChangChan := make(chan int)

	go func() {
		for connChange := range connChangChan {
			curConnNum += connChange
		}
	}()

	go func() {
		for _ = range time.Tick(1e9) {
			fmt.Printf("cur conn num: %v\n", curConnNum)
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

func sendMsg(conn net.Conn, msg []byte) {
	var buf []byte
	buf = make([]byte, 10, 10)
	_, err := conn.Write(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}

}

func sendKeepAliveMsg(conn net.Conn) int {
	conn.Write([]byte("a"))
	return 0
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
		time.Sleep(1 * 1e9)
	}
}

func httpServer(addr string) {
	http.HandleFunc("/cmd", CmdHander)
	http.HandleFunc("/band", BandwidthHander)
	http.HandleFunc("/data", DataHander)
	http.HandleFunc("/down", DownloadHander)
	http.HandleFunc("/echo", EchoHander)
	http.HandleFunc("/stat", StatHander)
	http.HandleFunc("/band/test/up", UploadApplyHander)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func main() {
	m := make(map[string]int)
	m["aaa"] = 42
	fmt.Println(m)
	get_internal()
	Port := flag.Int("Server Port", 8088, "Bandwidth Port")
	flag.Parse()
	ServerPort = *Port
	go BandwidthServer()
	httpServer(":8080")
}
