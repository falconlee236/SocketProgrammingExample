/**
 * ChatServer.go
20190532 sangyunLee
 **/

package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
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

	// totalClient number
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
		if err != nil {
			fmt.Println(err)
			break
		}

		// get client's nickname from client
		nicknameBuffer := make([]byte, 40)
		cnt, _ := conn.Read(nicknameBuffer)
		nicknameStr := string(nicknameBuffer[:cnt])

		//check client string already exist
		_, isExist := clientMap[nicknameStr]
		var nicknameRes string = ""
		// nickname handle status code
		var nicknameStatusCode int = 200
		if totalClientNum == 8 {
			nicknameStatusCode = 404
			nicknameRes = fmt.Sprintf("%d\n[chatting room full. cannot connect.]\n", nicknameStatusCode)
		} else if isExist {
			nicknameStatusCode = 404
			nicknameRes = fmt.Sprintf("%d\n[nickname already used by another user. cannot connect.]\n", nicknameStatusCode)
		} else {
			// add client name and connection info
			clientMap[nicknameStr] = conn
			// increase client count
			totalClientNum++
			nicknameRes = fmt.Sprintf("%d\n[welcome %s to CAU net-class chat room at %s.]\n"+
				"[There are %d users in the room.]\n", nicknameStatusCode, nicknameStr, conn.LocalAddr(), totalClientNum)
		}
		conn.Write([]byte(nicknameRes))
		// start client handling
		if nicknameStatusCode == 200 {
			fmt.Printf("[%s has joined from %s.]\n"+
				"[There are %d users in room.]\n\n", nicknameStr, conn.RemoteAddr().String(), totalClientNum)
			go TCPClientHandler(conn, &totalClientNum, &clientMap, nicknameStr)
		}
	}
}

func TCPClientHandler(conn net.Conn, totalClientNum *int, clientMap *map[string]net.Conn, nicknameStr string) {
	//multiple defer function Last in First out
	defer conn.Close()
	// if client request closed, that function called
	defer func(totalClientNum *int, clientMap *map[string]net.Conn, nicknameStr string) {
		// subtract total client number
		*totalClientNum -= 1
		// remove client info
		delete(*clientMap, nicknameStr)
	}(totalClientNum, clientMap, nicknameStr)

	for {
		msgRes := make([]byte, 1024)
		t, _ := conn.Read(msgRes)
		if t == 1 {
			command := msgRes[t-1]
			if command == 5 {
				sendMsg := fmt.Sprintf("[%s left the room.]\n[There are %d users now.]\n\n", nicknameStr, *totalClientNum-1)
				for nickname, otherConn := range *clientMap {
					if nickname == nicknameStr {
						continue
					}
					otherConn.Write([]byte(sendMsg))
				}
				fmt.Printf(sendMsg)
				return
			}
		} else {
			msg := string(msgRes[:t-1])
			for nickname, otherConn := range *clientMap {
				if nickname == nicknameStr {
					continue
				}
				sendMsg := fmt.Sprintf("%s> %s\n", nickname, msg)
				otherConn.Write([]byte(sendMsg))
			}
		}

		//var result string = ""
		//if typeStr == "1" {
		//	count, _ := conn.Read(buffer)
		//	result = string(bytes.ToUpper(buffer[:count]))
		//} else if typeStr == "2" {
		//	addrs := strings.Split(conn.RemoteAddr().String(), ":")
		//	result = fmt.Sprintf("client IP = %s, port = %s\n", addrs[0], addrs[1])
		//} else if typeStr == "3" {
		//	result = fmt.Sprintf("request served = %d\n", reqNum)
		//} else if typeStr == "4" {
		//
		//}
		//conn.Write([]byte(result))
	}
}
