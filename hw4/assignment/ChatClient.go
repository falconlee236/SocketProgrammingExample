/**
 * ChatClient.go
20190532 sangyunLee
 **/

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
)

func main() {
	// system args handling
	if len(os.Args) != 2 {
		fmt.Println("This program must be run as an argument.")
		os.Exit(1)
	}
	// get nickname from system args
	nickname := os.Args[1]
	//define command map
	commandMap := map[string]byte{
		"ls":     0x01,
		"secret": 0x02,
		"except": 0x03,
		"ping":   0x04,
		"quit":   0x05,
	}
	// server info
	serverName := "127.0.0.1"
	serverPort := "20532"

	// connect to server
	conn, err := net.Dial("tcp", serverName+":"+serverPort)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
		return
	}
	defer conn.Close()

	// send nickname to server
	conn.Write([]byte(nickname))
	accessResBuffer := make([]byte, 1024)
	cnt, _ := conn.Read(accessResBuffer)
	accessRes := strings.SplitN(string(accessResBuffer[:cnt]), "\n", 2)
	fmt.Println(accessRes[1])
	if accessRes[0] == "404" {
		return
	}

	// signal channel
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(conn net.Conn) {
		<-c
		conn.Write([]byte{commandMap["quit"]})
		fmt.Print("\ngg~\n")
		os.Exit(1)
	}(conn)

	for {
		msgInput, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		commandIdx := strings.IndexByte(msgInput, '\\')
		// start calculate RTT
		//start := time.Now()
		if commandIdx == 0 {
			commandMsg := strings.TrimRight(msgInput[1:], "\n")
			fmt.Println("Command:", commandMsg)
			byteValue, isExist := commandMap[commandMsg]
			if !isExist {
				fmt.Println("invalid command")
				continue
			}
			fmt.Println(byteValue)
			conn.Write([]byte{byteValue})
		} else {
			//start = time.Now()
			conn.Write([]byte(msgInput))
		}

		//if strings.TrimRight(msgInput, "\n") == "1" {
		//	fmt.Printf("Input lowercase sentence: ")
		//	input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		//	start = time.Now()
		//	conn.Write([]byte(input))
		//} else if strings.TrimRight(msgInput, "\n") == "5" {
		//	fmt.Printf("Bye bye~")
		//	return
		//}
		//// read from server
		//buffer := make([]byte, 1024)
		//read, err := conn.Read(buffer)
		//duration := time.Since(start)
		//if err != nil || read == 0 {
		//	log.Fatalf("Failed to connect to server: %v", err)
		//	return
		//}
		//// return microsecond
		//fmt.Printf("\nReply from server: %s", string(buffer))
		//fmt.Printf("RTT = %fms\n", float64(duration.Nanoseconds())/1e+6)
	}
}
