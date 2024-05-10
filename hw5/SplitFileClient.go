/*
SplitFileClient.go
20190532 sang yun lee
*/
package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	// argument handling
	if len(os.Args) != 3 {
		fmt.Println("Invalid argument.")
		os.Exit(1)
	}
	// get argument from system args
	commandName := os.Args[1]
	fileName := os.Args[2]

	// server Info hardcoding
	firstServerName := "127.0.0.1"
	firstServerPort := "20532"
	secondServerName := "127.0.0.1"
	secondServerPort := "20532"

	// put command case
	if commandName == "put" {
		sendFile(fileName, firstServerName, firstServerPort, 0)
		sendFile(fileName, secondServerName, secondServerPort, 1)
	} else if os.Args[1] == "get" { // get command case

	} else { // otherwise case
		fmt.Println("Invalid argument.")
	}
}

// put file to server
func sendFile(fileName string, serverName string, serverPort string, part int) {

	// connect to TCP server
	conn, err := net.Dial("tcp", serverName+":"+serverPort)
	if err != nil {
		fmt.Println("fail to connect server: ", err)
		return
	}
	defer conn.Close()

	// prepare command string to send command
	commandStr := "put"
	conn.Write([]byte(commandStr))
	commandBuffer := make([]byte, 1024)
	read, _ := conn.Read(commandBuffer)
	// is server response is not ok
	if string(commandBuffer[:read]) != "ok" {
		fmt.Println("fail to receive fileName")
		os.Exit(1)
	}

	// try to open file
	originalFile, err := os.Open(fileName)
	if err != nil {
		return
	}

	// get file extension from fileName
	fileExtension := filepath.Ext(fileName)
	// get file name without extension from fileName
	fileName = fileName[0 : len(fileName)-len(fileExtension)]
	// join that strings
	fileName = fmt.Sprintf("%s-part%d%s", fileName, part+1, fileExtension)
	// send to target File name
	conn.Write([]byte(fileName))
	fileNameBuffer := make([]byte, 1024)
	read, _ = conn.Read(fileNameBuffer)
	if string(fileNameBuffer[:read]) != "ok" {
		fmt.Println("fail to receive command")
		os.Exit(1)
	}

	// try to get file Stats
	fileInfo, err := originalFile.Stat()
	if err != nil {
		fmt.Println("파일 정보 가져오기 실패:", err)
		return
	}

	// convert file Size string to Integer
	size := strconv.FormatInt(fileInfo.Size(), 10)
	// send file Size to server
	conn.Write([]byte(size))
	fileSizeBuffer := make([]byte, 1024)
	read, _ = conn.Read(fileSizeBuffer)
	if string(fileSizeBuffer[:read]) != "ok" {
		fmt.Println("fail to receive fileName")
		os.Exit(1)
	}

	// get Buffer from file
	reader := bufio.NewReader(originalFile)
	var byteCount int
	for {
		// get 1 byte from file
		b, err := reader.ReadByte()
		if err != nil {
			break // finish to reach file end
		}

		// check odd, even byte
		if byteCount%2 == part {
			_, err := conn.Write([]byte{b})
			if err != nil {
				return
			}
		}
		byteCount++
	}

	fmt.Printf("%s send successful\n", fileName)
	return
}

func receiveFile(fileName string, serverName string, serverPort string, part int) {
	// connect to TCP server
	conn, err := net.Dial("tcp", serverName+":"+serverPort)
	if err != nil {
		fmt.Println("fail to connect server: ", err)
		return
	}
	defer conn.Close()

	// prepare command string to send command
	commandStr := "get"
	conn.Write([]byte(commandStr))
	commandBuffer := make([]byte, 1024)
	read, _ := conn.Read(commandBuffer)
	// is server response is not ok
	if string(commandBuffer[:read]) != "ok" {
		fmt.Println("fail to receive command")
		os.Exit(1)
	}

	fileNameBuffer := make([]byte, 1024)
	read, err = conn.Read(fileNameBuffer)
	if err != nil {
		conn.Write([]byte("fail to read file"))
	} else {
		conn.Write([]byte("ok"))
	}
	fileName = string(fileNameBuffer[:read])
	fmt.Println("Received file: ", fileName)

	fileSizeBuffer := make([]byte, 1024)
	read, err = conn.Read(fileSizeBuffer)
	if err != nil {
		conn.Write([]byte("fail to read file"))
	} else {
		conn.Write([]byte("ok"))
	}
	fileSize, err := strconv.ParseInt(string(fileSizeBuffer[:read]), 10, 64)
	if err != nil {
		fmt.Println("fail to transfer file size:", err)
		return
	}

	// 파일 생성
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("fail to create file:", err)
		return
	}

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
			file.Close()
			break
		}
	}

	fmt.Printf("%s file store sucessful!\n", fileName)
}
