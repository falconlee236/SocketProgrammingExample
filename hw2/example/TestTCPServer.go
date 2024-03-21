package main

import (
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

func main() {
	sock, err := net.Listen("tcp", "nsl5.cau.ac.kr:20532")
	if err != nil {
		log.Fatalf("Failed to bind address: %v", err)
	}

	defer sock.Close()

	for {
		conn, err := sock.Accept()
		if err != nil {
			log.Printf("Failed to accept: %v", err)
			continue
		}

		go ClientHandler(conn)
	}
}

func ClientHandler(conn net.Conn) {
	defer conn.Close()

	recvBuf := make([]byte, 4096)

	for {

		readBytes, err := conn.Read(recvBuf)
		if err != nil {
			if err == io.EOF {
				log.Printf("Client closed the socket: %v", conn.RemoteAddr().String())
				return
			}

			log.Printf("Failed to receive data: %v", err)
			return
		}

		if readBytes > 0 {
			data := recvBuf[:readBytes]

			if number, err := strconv.Atoi(strings.Trim(string(data), string('\n'))); err == nil {
				time.Sleep(time.Duration(number) * time.Second)
			}

			conn.Write(data)

			log.Println(string(data))
		} else {
			log.Println("Can you see me?")
		}

	}
}
