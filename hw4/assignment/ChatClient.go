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
	"time"
)

var currentTime time.Time

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
	serverName := "nsl5.cau.ac.kr"
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

	// remove res status and msg
	accessRes := strings.SplitN(string(accessResBuffer[:cnt]), "\n", 2)
	fmt.Println(accessRes[1])
	// if status is error
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
				fmt.Print("\ngg~\n")
				os.Exit(1)
			}
			// if server send ping command again
			if read == 1 || buffer[0] == commandMap["ping"] {
				duration := time.Since(currentTime)
				// return microsecond
				fmt.Printf("RTT = %fms\n", float64(duration.Nanoseconds())/1e+6)
			} else { // else msg, print msg to client
				fmt.Println(string(buffer[:read]))
			}
		}
	}(conn)

	for {
		// input msg from user
		msgInput, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		if len(msgInput) == 1 && msgInput[0] == '\n'{
			continue
		}
		// find command index
		commandIdx := strings.IndexByte(msgInput, '\\')
		// if command
		if commandIdx == 0 {
			msgInput = strings.TrimRight(msgInput, "\n")
			msgArr := strings.SplitN(msgInput, " ", 2)
			command := msgArr[0][1:]
			// encoding command to 1byte
			byteValue, isExist := commandMap[command]
			// cannot find in command table
			if !isExist ||
				((byteValue == 1 || byteValue == 4 || byteValue == 5) && len(msgArr) != 1 ||
					((byteValue == 2 || byteValue == 3) && len(msgArr) != 2)) {
				fmt.Println("Invalid command")
				continue
			}
			if len(msgArr) > 1 {
				msg := msgArr[1]
				_, err := conn.Write(append([]byte{byteValue, ' '}, []byte(msg)...))
				if err != nil {
					log.Fatalf("Failed to connect to server\n%v", err)
				}
			} else {
				// start calculate start time
				currentTime = time.Now()
				_, err := conn.Write([]byte{byteValue})
				if err != nil {
					log.Fatalf("Failed to connect to server\n%v", err)
				}
			}
		} else {
			_, err := conn.Write([]byte(msgInput))
			if err != nil {
				log.Fatalf("Failed to connect to server\n%v", err)
			}
		}
		fmt.Println()
	}
}
