package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net/http"
	"strings"
)

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

	addr := strings.Split(r.Host, ":")
	ip := addr[0]

	http.Redirect(w, r, "http://"+ip+":"+"8080"+"/band/test/up", 302)
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

/* UploadApplyHander
return values
http status 200
content:
Host_len	Uint32	测速服务器host长度
Host_name	Uint8，长度由host_len指定	测速服务器域名
Port	Uint16	测速服务器port
status	Uint8	允许测速状态码，1代表允许测速，0代表不允许测速
*/
func UploadApplyHander(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	addr := strings.Split(r.Host, ":")
	ip := addr[0]
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, len(ip))
	buffer.WriteString(ip)
	binary.Write(buffer, binary.BigEndian, 8088)
	if curConnNum < 10 {
		buffer.WriteString("1")
	} else {
		buffer.WriteString("0")
	}
	/*
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
	*/
	fmt.Println(path, "host:", ip, r.Host, r.RemoteAddr, r.URL)
	//w.Write([]byte(buffer))
	w.Write(buffer.Bytes())
}

func Router() {
	http.HandleFunc("/cmd", CmdHander)
	http.HandleFunc("/band", BandwidthHander)
	http.HandleFunc("/data", DataHander)
	http.HandleFunc("/down", DownloadHander)
	http.HandleFunc("/echo", EchoHander)
	http.HandleFunc("/stat", StatHander)
	http.HandleFunc("/band/test/up", UploadApplyHander)
}
