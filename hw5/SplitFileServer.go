package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

func main() {
	serverPort := "20532"

	// 서버 시작
	listener, err := net.Listen("tcp", ":"+serverPort)
	if err != nil {
		fmt.Println("서버 시작 실패:", err)
		return
	}
	defer listener.Close()
	fmt.Println("서버 시작, 클라이언트 연결 대기중...")

	for {
		// 클라이언트 연결 대기
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("클라이언트 연결 실패:", err)
			continue
		}

		// 클라이언트가 연결되었을 때 파일을 수신하는 함수 호출
		go handleClient(conn)
	}
}

// 클라이언트로부터 파일을 수신하여 저장하는 함수
func handleClient(conn net.Conn) {
	defer conn.Close()

	// 파일 크기 읽기
	fileSizeBuffer := make([]byte, 8)
	cnt, err := conn.Read(fileSizeBuffer)
	if err != nil {
		fmt.Println("파일 크기 읽기 실패:", err)
		return
	}

	// 파일 이름 읽기
	filenameBuffer := make([]byte, 1024)
	_, err = conn.Read(filenameBuffer)
	if err != nil {
		fmt.Println("파일 이름 읽기 실패:", err)
		return
	}
	filename := string(filenameBuffer)
	fmt.Println(filename)

	// 파일 크기 변환
	fileSize, err := strconv.ParseInt(string(fileSizeBuffer[:cnt]), 10, 64)
	if err != nil {
		fmt.Println("파일 크기 변환 실패:", err)
		return
	}

	// 파일 생성
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("파일 생성 실패:", err)
		return
	}
	defer file.Close()

	// 파일 내용 수신하여 저장
	var receivedBytes int64
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println("파일 읽기 실패:", err)
			}
			break
		}
		receivedBytes += int64(n)

		// 파일에 받은 내용 쓰기
		if _, err := file.Write(buffer[:n]); err != nil {
			fmt.Println("파일 쓰기 실패:", err)
			return
		}

		// 파일이 전부 받아졌는지 확인
		if receivedBytes >= fileSize {
			break
		}
	}

	fmt.Printf("%s 파일 전송 완료\n", filename)
}
