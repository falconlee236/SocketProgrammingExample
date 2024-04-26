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
	"strings"
)

func main() {
	// signal handling
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Print("\ngg~\n")
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
		//accept response message
		msgRes := make([]byte, 1024)
		t, _ := conn.Read(msgRes)
		if t == 0 {
			return
		}
		// if that message is ls, ping, quit command or invalid command
		if t == 1 {
			command := msgRes[t-1]
			// if command is invalid
			if command == 0 {
				fmt.Println("Invalid command received from client")
			} else if command == 1 { // ls command
				sendMsg := ""
				// iterate all map structure
				for nickname, otherConn := range *clientMap {
					// split connection to ip, port
					ip, port, _ := net.SplitHostPort(otherConn.RemoteAddr().String())
					// add info
					sendMsg += fmt.Sprintf("<%s, %s, %s>\n", nickname, ip, port)
				}
				conn.Write([]byte(sendMsg))
			} else if command == 4 { // ping command
				// return ping byte back
				conn.Write(msgRes[:t])
			} else if command == 5 { // quit command
				sendMsg := fmt.Sprintf("[%s left the room.]\n[There are %d users now.]\n\n", nicknameStr, *totalClientNum-1)
				// send msg to other client
				for nickname, otherConn := range *clientMap {
					if nickname == nicknameStr {
						continue
					}
					otherConn.Write([]byte(sendMsg))
				}
				// send msg to server
				fmt.Printf(sendMsg)
				// finish client
				return
			}
		} else if msgRes[0] == 0 { // invalid command, within space in msg
			fmt.Println("Invalid command: " + string(msgRes[1:]))
		} else if msgRes[0] == 2 || msgRes[0] == 3 { // valid command, command secret, except
			msgArr := strings.SplitN(string(msgRes[1:t]), " ", 2)
			// command parameter error
			if len(msgArr) != 2 {
				fmt.Println("Invalid command: " + string(msgRes[1:]))
				continue
			}
			// split nickname, msg
			commandNickname := msgArr[0]
			commandMsg := msgArr[1]
			if strings.Contains(strings.ToLower(commandMsg), strings.ToLower("I hate professor")) {
				sendMsg := fmt.Sprintf("[%s is disconnected.]\n"+
					"[There are %d users in the chat room.]\n", nicknameStr, *totalClientNum-1)
				for _, otherConn := range *clientMap {
					otherConn.Write([]byte(sendMsg))
				}
				fmt.Print(sendMsg)
				return
			}
			// secret command
			if msgRes[0] == 2 {
				// get nickname connection info
				secretConn, isExist := (*clientMap)[commandNickname]
				// if that nickname does not exist in map
				if !isExist {
					fmt.Println("Invalid command: " + string(msgRes[1:]))
					continue
				}
				sendMsg := fmt.Sprintf("from: %s> %s\n", nicknameStr, commandMsg)
				secretConn.Write([]byte(sendMsg))
			} else if msgRes[0] == 3 { // sxcept command
				for nickname, otherConn := range *clientMap {
					// find except nickname, don't send msg to that nickname
					if nickname == commandNickname {
						continue
					}
					sendMsg := fmt.Sprintf("from: %s> %s\n", nicknameStr, commandMsg)
					otherConn.Write([]byte(sendMsg))
				}
			}
		} else { // otherwise
			msg := string(msgRes[:t-1])
			if strings.Contains(strings.ToLower(msg), strings.ToLower("I hate professor")) {
				sendMsg := fmt.Sprintf("[%s is disconnected.]\n"+
					"[There are %d users in the chat room.]\n", nicknameStr, *totalClientNum-1)
				for _, otherConn := range *clientMap {
					otherConn.Write([]byte(sendMsg))
				}
				fmt.Print(sendMsg)
				return
			}
			for nickname, otherConn := range *clientMap {
				if nickname == nicknameStr {
					continue
				}
				sendMsg := fmt.Sprintf("%s> %s\n", nicknameStr, msg)
				otherConn.Write([]byte(sendMsg))
			}
		}
	}
}
