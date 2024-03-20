/**
 * TCPServer.go
 **/

package main

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"time"
)

func main() {
	start := time.Now()
	reqNum := 0
	serverPort := "20532"

	listener, _ := net.Listen("tcp", ":"+serverPort)
	fmt.Printf("Server is ready to receive on port %s\n", serverPort)
	typeBuffer := make([]byte, 1024)
	buffer := make([]byte, 1024)

	for {
		conn, _ := listener.Accept()
		fmt.Printf("Connection request from %s\n", conn.RemoteAddr().String())

		t, _ := conn.Read(typeBuffer)
		if t == 0 {
			fmt.Print("Bye bye~\n")
			continue
		}
		typeStr := string(typeBuffer[:t-1])
		fmt.Printf("Command %s\n", typeStr)

		if typeStr == "1" {
			count, _ := conn.Read(buffer)
			conn.Write(bytes.ToUpper(buffer[:count]))
		} else if typeStr == "2" {
			conn.Write([]byte(conn.RemoteAddr().String()))
		} else if typeStr == "3" {
			conn.Write([]byte(strconv.Itoa(reqNum)))
		} else if typeStr == "4" {
			duration := time.Since(start)
			hour := int(duration.Seconds() / 3600)
			minute := int(duration.Seconds()/60) % 60
			second := int(duration.Seconds()) % 60
			totalRuntime := fmt.Sprintf("%02d:%02d:%02d\n", hour, minute, second)
			conn.Write([]byte(totalRuntime))
		}
		reqNum++
		conn.Close()
	}
}
