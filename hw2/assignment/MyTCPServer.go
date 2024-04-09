/**
 * MyTCPServer.go
20190532 sangyunLee
 **/

package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"time"
)

func main() {
	// signal channel
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Print("\nBye bye~\n")
			os.Exit(1)
		}
	}()

	// start to calculate server time
	start := time.Now()
	reqNum := 0
	serverPort := "20532"

	// wait to client request
	listener, _ := net.Listen("tcp", ":"+serverPort)
	fmt.Printf("Server is ready to receive on port %s\n", serverPort)

	defer listener.Close()
	for {
		// accept client request
		conn, err := listener.Accept()
		fmt.Printf("Connection request from %s\n", conn.RemoteAddr().String())
		if err != nil {
			fmt.Println(err)
			break
		}

		// start client handling
		go TCPClientHandler(conn, &reqNum, start)
	}
}

func TCPClientHandler(conn net.Conn, reqNum *int, start time.Time) {
	defer conn.Close()

	typeBuffer := make([]byte, 1024)

	for {
		// get data from server
		buffer := make([]byte, 1024)
		t, _ := conn.Read(typeBuffer)
		if t == 0 {
			return
		}
		typeStr := string(typeBuffer[:t-1])
		fmt.Printf("Command %s\n", typeStr)

		var result string = ""
		if typeStr == "1" {
			count, _ := conn.Read(buffer)
			result = string(bytes.ToUpper(buffer[:count]))
		} else if typeStr == "2" {
			addrs := strings.Split(conn.RemoteAddr().String(), ":")
			result = fmt.Sprintf("client IP = %s, port = %s\n", addrs[0], addrs[1])
		} else if typeStr == "3" {
			result = fmt.Sprintf("request served = %d\n", *reqNum)
		} else if typeStr == "4" {
			// change to millisecond
			duration := time.Since(start)
			hour := int(duration.Seconds() / 3600)
			minute := int(duration.Seconds()/60) % 60
			second := int(duration.Seconds()) % 60
			result = fmt.Sprintf("run time = %02d:%02d:%02d\n", hour, minute, second)
		}
		// return to server
		conn.Write([]byte(result))
		(*reqNum)++
	}
}
