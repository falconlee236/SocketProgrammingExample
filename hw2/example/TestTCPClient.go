package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	socket, err := net.Dial("tcp", "localhost:18080")
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
		return
	}

	defer socket.Close()
	console := bufio.NewReader(os.Stdin)
	recvBuf := make([]byte, 4096)

	for {
		fmt.Print("socket> ")
		text, _ := console.ReadString('\n')

		if number, err := strconv.Atoi(strings.Trim(string(text), string('\n'))); err == nil {
			fmt.Printf("Echo delayed: %v\n", number)
		}

		if text == "quit\n" {
			break
		}
		socket.Write([]byte(text))

		readBytes, err := socket.Read(recvBuf)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Server closed the socket: %v", socket.RemoteAddr().String())
				return
			}

			fmt.Printf("Failed to receive data: %v", err)
			return
		}

		if readBytes > 0 {
			data := recvBuf[:readBytes]
			fmt.Println(string(data))
		} else {
			fmt.Println("Can you see me?")
		}
	}
}
