/**
 * MultiTCPServer.go
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
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Print("\nBye bye~\n")
			os.Exit(1)
		}
	}()

	start := time.Now()
	totalClientNum := 0
	serverPort := "20532"
	clientList := make([]int, 2)

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

		nextClientIdx := findNextClientIdx(clientList)
		if len(clientList) <= nextClientIdx {
			clientList = append(clientList, nextClientIdx+1)
		} else {
			clientList[nextClientIdx] = nextClientIdx + 1
		}
		totalClientNum++

		go TCPClientHandler(conn, start, &totalClientNum, nextClientIdx+1, &clientList)
	}
}

func findNextClientIdx(clientList []int) int {
	i := 0
	for i = 0; i < len(clientList); i++ {
		if clientList[i] == 0 {
			return i
		}
	}
	return i
}

func TCPClientHandler(conn net.Conn, start time.Time,
	totalClientNum *int, clientNum int, clientList *[]int) {
	defer conn.Close() //multiple defer function Last in First out
	defer func(totalClientNum *int, clientNum int, clientList *[]int) {
		(*clientList)[clientNum-1] = 0
		*totalClientNum -= 1
		currentTimeStr := time.Now().Format("15:04:05")
		fmt.Printf("[Time: %s] Client %d disconnected."+" Number of clients connected = %d\n",
			currentTimeStr, clientNum, *totalClientNum)
	}(totalClientNum, clientNum, clientList)

	currentTimeStr := time.Now().Format("15:04:05")
	fmt.Printf("[Time: %s] Client %d connected."+" Number of clients connected = %d\n",
		currentTimeStr, clientNum, *totalClientNum)

	typeBuffer := make([]byte, 1024)
	reqNum := 0

	for {
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
			result = fmt.Sprintf("request served = %d\n", reqNum)
		} else if typeStr == "4" {
			duration := time.Since(start)
			hour := int(duration.Seconds() / 3600)
			minute := int(duration.Seconds()/60) % 60
			second := int(duration.Seconds()) % 60
			result = fmt.Sprintf("run time = %02d:%02d:%02d\n", hour, minute, second)
		}
		conn.Write([]byte(result))
		reqNum++
	}
}
