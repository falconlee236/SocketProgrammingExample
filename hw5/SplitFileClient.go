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
	"strconv"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Invalid argument.")
		os.Exit(1)
	}
	//commandName := os.Args[1]
	fileName := os.Args[2]
	// 서버의 주소와 포트
	serverName := "127.0.0.1"
	serverPort := "20532"

	// 파일을 두 부분으로 나눔
	part1, part2, err := splitFile(fileName)
	if err != nil {
		fmt.Println("파일 분할 실패:", err)
		return
	}

	// TCP 서버에 연결
	conn, err := net.Dial("tcp", serverName+":"+serverPort)
	if err != nil {
		fmt.Println("서버에 연결 실패:", err)
		return
	}
	defer conn.Close()

	// 첫 번째 부분 파일 전송
	sendFile(conn, part1, fileName+"-part1.txt")

	// 두 번째 부분 파일 전송
	sendFile(conn, part2, fileName+"-part2.txt")

	fmt.Println("파일 전송 완료")
}

// 파일을 TCP 소켓을 통해 서버로 전송하는 함수
func sendFile(conn net.Conn, file *os.File, filename string) {
	// 파일 정보 가져오기
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("파일 정보 가져오기 실패:", err)
		return
	}

	// 파일 크기 전송
	fileSize := make([]byte, 1024)
	size := strconv.FormatInt(fileInfo.Size(), 10)
	copy(fileSize[:], size)
	conn.Write(fileSize)

	// 파일 이름 전송
	conn.Write([]byte(filename))

	// 파일 내용 전송
	bufferSize := 1024
	buffer := make([]byte, bufferSize)
	for {
		bytesRead, err := file.Read(buffer)
		if err != nil {
			break
		}
		conn.Write(buffer[:bytesRead])
	}
}

// 파일을 두 개의 부분으로 분할하는 함수
func splitFile(filename string) (part1, part2 *os.File, err error) {
	// 원본 파일 열기
	originalFile, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}

	// 부분 파일들 생성
	part1, err = os.Create(filename + "-part1.txt")
	if err != nil {
		return nil, nil, err
	}

	part2, err = os.Create(filename + "-part2.txt")
	if err != nil {
		return nil, nil, err
	}

	// 원본 파일을 읽기 위한 버퍼 생성
	reader := bufio.NewReader(originalFile)

	// 바이트 단위로 파일을 읽고, 홀수와 짝수 바이트를 번갈아가며 각각의 파일에 쓰기
	var byteCount int
	for {
		// 파일에서 한 바이트 읽기
		b, err := reader.ReadByte()
		if err != nil {
			break // 파일 끝에 도달하면 종료
		}

		// 홀수 바이트는 part1에, 짝수 바이트는 part2에 쓰기
		if byteCount%2 == 0 {
			_, err := part1.Write([]byte{b}) // part1에 바이트 쓰기
			if err != nil {
				return nil, nil, err
			}
		} else {
			_, err := part2.Write([]byte{b}) // part2에 바이트 쓰기
			if err != nil {
				return nil, nil, err
			}
		}

		byteCount++
	}

	return part1, part2, nil
}
