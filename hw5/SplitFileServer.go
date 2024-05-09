/*
SplitFileServer.go
20190532 sang yun lee
*/

package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

func main() {
	//if len(os.Args) != 2 {
	//	fmt.Println("This program only accepts one argument")
	//	os.Exit(1)
	//}
	//serverPort := os.Args[1]

	// server Port
	serverPort := "20532"

	// listen from client request
	listener, err := net.Listen("tcp", ":"+serverPort)
	if err != nil {
		fmt.Println("fail to start server: ", err)
		return
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println("fail to close listener: ", err)
		}
	}(listener)
	fmt.Println("Waiting for connections...")

	// main loop
	for {
		// accept to client request
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("fail to accept client request: ", err)
			continue
		}
		// file handling function with goroutine
		go putFile(conn)
	}
}

// Store file from client request handling function
func putFile(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("fail to close connection: ", err)
		}
	}(conn)

	// get file size from client
	fileSizeBuffer := make([]byte, 8)
	cnt, err := conn.Read(fileSizeBuffer)
	if err != nil {
		fmt.Println("fail to read file size:", err)
		return
	}

	// get file name from client
	filenameBuffer := make([]byte, 1024)
	_, err = conn.Read(filenameBuffer)
	if err != nil {
		fmt.Println("fail to read file name:", err)
		return
	}
	filename := string(filenameBuffer)
	fmt.Println(filename)

	// 파일 크기 변환
	fileSize, err := strconv.ParseInt(string(fileSizeBuffer[:cnt]), 10, 64)
	if err != nil {
		fmt.Println("fail to transfer file size:", err)
		return
	}

	// 파일 생성
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("fail to create file:", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("fail to close file:", err)
		}
	}(file)

	// 파일 내용 수신하여 저장
	var receivedBytes int64
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println("fail to read file:", err)
			}
			break
		}
		receivedBytes += int64(n)

		// 파일에 받은 내용 쓰기
		if _, err := file.Write(buffer[:n]); err != nil {
			fmt.Println("fail to write file:", err)
			return
		}

		// 파일이 전부 받아졌는지 확인
		if receivedBytes >= fileSize {
			break
		}
	}

	fmt.Printf("%s file store sucessful!\n", filename)
	return
}
