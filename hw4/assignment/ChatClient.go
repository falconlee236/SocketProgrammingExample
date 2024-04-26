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
		log.Fatalf("Failed to connect to server\n%v", err)
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Failed to connect to server\n%v", err)
		}
	}(conn)

	// send nickname to server
	_, err = conn.Write([]byte(nickname))
	if err != nil {
		log.Fatalf("Failed to connect to server\n%v", err)
	}
	accessResBuffer := make([]byte, 1024)
	cnt, _ := conn.Read(accessResBuffer)
	accessRes := strings.SplitN(string(accessResBuffer[:cnt]), "\n", 2)
	fmt.Println(accessRes[1])
	if accessRes[0] == "404" {
		os.Exit(1)
	}

	// signal channel
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(conn net.Conn) {
		<-c
		_, err := conn.Write([]byte{commandMap["quit"]})
		if err != nil {
			log.Fatalf("Failed to connect to server\n%v", err)
		}
		fmt.Print("\ngg~\n")
		os.Exit(1)
	}(conn)

	// Listen for incoming messages from the client
	go func(conn net.Conn) {
		for {
			// read from server
			buffer := make([]byte, 1024)
			read, err := conn.Read(buffer)
			if err != nil || read == 0 {
				return
			}
			fmt.Println(string(buffer[:read]))
		}
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
				log.Fatalf("Failed to connect to server\n%v", err)
			}
			_, err := conn.Write([]byte{byteValue})
			if err != nil {
				log.Fatalf("Failed to connect to server\n%v", err)
			}
		} else {
			//start = time.Now()
			_, err := conn.Write([]byte(msgInput))
			if err != nil {
				log.Fatalf("Failed to connect to server\n%v", err)
			}
		}
		fmt.Println()
		//if strings.TrimRight(msgInput, "\n") == "1" {
		//	fmt.Printf("Input lowercase sentence: ")
		//	input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		//	start = time.Now()
		//	conn.Write([]byte(input))
		//} else if strings.TrimRight(msgInput, "\n") == "5" {
		//	fmt.Printf("Bye bye~")
		//	return
		//}
		//duration := time.Since(start)
		//// return microsecond
		//fmt.Printf("\nReply from server: %s", string(buffer))
		//fmt.Printf("RTT = %fms\n", float64(duration.Nanoseconds())/1e+6)
	}
}
