/**
 * TCPServer.go
 **/

package main

import (
	"bytes"
	"fmt"
	"net"
	"time"
)

func main() {
	start := time.Now()
	serverPort := "20532"

	listener, _ := net.Listen("tcp", ":"+serverPort)
	fmt.Printf("Server is ready to receive on port %s\n", serverPort)
	typeBuffer := make([]byte, 1024)
	buffer := make([]byte, 1024)

	for {
		conn, _ := listener.Accept()
		fmt.Printf("Connection request from %s\n", conn.RemoteAddr().String())

		t, _ := conn.Read(typeBuffer)
		typeStr := string(typeBuffer[:t - 1])
		fmt.Printf("Command %s\n", typeStr)

		if typeStr == "1" {
			count, _ := conn.Read(buffer)
			conn.Write(bytes.ToUpper(buffer[:count]))
		} else if typeStr == "4" {
			duration := time.Since(start)
			//conn.Write(duration)
			fmt.Println(duration)
		}
		conn.Close()
	}
}
