/**
 * ChatServer.go
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
)

func main() {
	// signal handling
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Print("\nBye bye~\n")
			os.Exit(1)
		}
	}()

	// totalClient
	totalClientNum := 0
	serverPort := "20532"
	// init client id List
	clientMap := make(map[string]net.Conn, 8)

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

		// get client's nickname from client
		nicknameBuffer := make([]byte, 40)
		cnt, _ := conn.Read(nicknameBuffer)
		nicknameStr := string(nicknameBuffer[:cnt])
		_, isExist := clientMap[nicknameStr]
		var nicknameRes string = ""
		var nicknameStatusCode int = 200
		if len(clientMap) == 8 {
			nicknameStatusCode = 404
			nicknameRes = fmt.Sprintf("%d\nchatting room full. cannot connect\n", nicknameStatusCode)
		} else if isExist {
			nicknameRes = fmt.Sprintf("%d\nnickname already used by another user. cannot connect\n", nicknameStatusCode)
		} else {
			clientMap[nicknameStr] = conn
			totalClientNum++
			nicknameRes = fmt.Sprintf("%d\n[welcome %s to CAU net-class chat room at %s]\n"+
				"[There are %d users in the room]\n", nicknameStatusCode, nicknameStr, conn.LocalAddr(), totalClientNum)
		}
		conn.Write([]byte(nicknameRes))
		// start client handling
		if nicknameStatusCode == 200 {
			go TCPClientHandler(conn, &totalClientNum)
		}
	}
}

func TCPClientHandler(conn net.Conn, totalClientNum *int) {
	defer conn.Close() //multiple defer function Last in First out
	// if client request closed, that function called
	defer func(totalClientNum *int) {
		// subtract total client number
		*totalClientNum -= 1
	}(totalClientNum)

	typeBuffer := make([]byte, 1024)
	reqNum := 0

	// same to hw2 client handling
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

		}
		conn.Write([]byte(result))
		reqNum++
	}
}
