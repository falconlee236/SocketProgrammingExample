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
	"strconv"
	"time"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Print("\nBye bye~\n")
			os.Exit(1)
		}
	}()

	start := time.Now()
	reqNum := 0
	serverPort := "20532"

	listener, _ := net.Listen("tcp", ":"+serverPort)
	fmt.Printf("Server is ready to receive on port %s\n", serverPort)

	defer listener.Close()
	for {
		conn, err := listener.Accept()
		fmt.Printf("Connection request from %s\n", conn.RemoteAddr().String())
		if err != nil {
			fmt.Println(err)
			break
		}

		go ClientHandler(conn, reqNum, start)
	}
}

func ClientHandler(conn net.Conn, reqNum int, start time.Time) {
	defer conn.Close()

	typeBuffer := make([]byte, 1024)
	buffer := make([]byte, 1024)

	for {
		t, _ := conn.Read(typeBuffer)
		if t == 0 {
			return
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
		} else if typeStr == "5" {
			fmt.Print("inside 5")
			conn.Write(nil)
		}
		reqNum++
	}
}
