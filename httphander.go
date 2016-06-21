package main

import (
	"fmt"
	"net/http"
	"strconv"
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

func Router() {
	http.HandleFunc("/cmd", CmdHander)
	http.HandleFunc("/band", BandwidthHander)
	http.HandleFunc("/data", DataHander)
	http.HandleFunc("/down", DownloadHander)
	http.HandleFunc("/echo", EchoHander)
	http.HandleFunc("/stat", StatHander)
	http.HandleFunc("/band/test/up", UploadApplyHander)
}
