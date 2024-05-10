/*
SplitFileClient.go
20190532 sang yun lee
*/
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Invalid argument.")
		os.Exit(1)
	}
	//commandName := os.Args[1]
	fileName := os.Args[2]

	firstServerName := "127.0.0.1"
	firstServerPort := "20532"
	secondServerName := "127.0.0.1"
	secondServerPort := "20532"

	if os.Args[1] == "put" {
		sendFile(fileName, firstServerName, firstServerPort, 0)
		sendFile(fileName, secondServerName, secondServerPort, 1)
	} else if os.Args[1] == "get" {

	} else {
		fmt.Println("Invalid argument.")
	}
}

func sendFile(fileName string, serverName string, serverPort string, part int) {

	// TCP 서버에 연결
	conn, err := net.Dial("tcp", serverName+":"+serverPort)
	if err != nil {
		fmt.Println("서버에 연결 실패:", err)
		return
	}
	defer conn.Close()

	originalFile, err := os.Open(fileName)
	if err != nil {
		return
	}

	fileExtension := filepath.Ext(fileName)
	fileName = fileName[0 : len(fileName)-len(fileExtension)]
	fileName = fmt.Sprintf("%s-part%d%s", fileName, part+1, fileExtension)
	conn.Write([]byte(fileName))
	fileNameBuffer := make([]byte, 1024)
	read, _ := conn.Read(fileNameBuffer)
	if string(fileNameBuffer[:read]) != "ok" {
		fmt.Println("fail to receive fileName")
		os.Exit(1)
	}

	fileInfo, err := originalFile.Stat()
	if err != nil {
		fmt.Println("파일 정보 가져오기 실패:", err)
		return
	}

	// 파일 크기 전송
	size := strconv.FormatInt(fileInfo.Size(), 10)
	conn.Write([]byte(size))
	fileSizeBuffer := make([]byte, 1024)
	read, _ = conn.Read(fileSizeBuffer)
	if string(fileSizeBuffer[:read]) != "ok" {
		fmt.Println("fail to receive fileName")
		os.Exit(1)
	}

	reader := bufio.NewReader(originalFile)
	// 바이트 단위로 파일을 읽고, 홀수와 짝수 바이트를 번갈아가며 각각의 파일에 쓰기
	var byteCount int
	for {
		// 파일에서 한 바이트 읽기
		b, err := reader.ReadByte()
		if err != nil {
			break // 파일 끝에 도달하면 종료
		}

		if byteCount%2 == part {
			_, err := conn.Write([]byte{b})
			if err != nil {
				return
			}
		}
		byteCount++
	}

	fmt.Println("파일 전송 완료")
	return
}
