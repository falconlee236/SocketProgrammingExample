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
	"strconv"
	"strings"
	"time"
)

func main() {
	// system args handling
	if len(os.Args) != 2 {
		fmt.Println("This program must be run as an argument.")
		os.Exit(1)
	}
	// get nickname from system args
	nickname := os.Args[1]

	// signal channel
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Print("\nBye bye~\n")
			os.Exit(1)
		}
	}()
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

	for {
		fmt.Printf("<Menu>\n")
		fmt.Printf("1) convert text to UPPER-case\n")
		fmt.Printf("2) get my IP address and port number\n")
		fmt.Printf("3) get server request count\n")
		fmt.Printf("4) get server running time\n")
		fmt.Printf("5) exit\n")
		// input option string
		fmt.Printf("Input option: ")
		inputOption, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		inputNum, _ := strconv.Atoi(strings.TrimRight(inputOption, "\n"))
		// not number error handling
		if inputNum < 1 || inputNum > 5 {
			fmt.Print("Invalid option\n")
			continue
		}
		// start calculate RTT
		start := time.Now()
		// send to Server
		conn.Write([]byte(inputOption))

		if strings.TrimRight(inputOption, "\n") == "1" {
			fmt.Printf("Input lowercase sentence: ")
			input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
			start = time.Now()
			conn.Write([]byte(input))
		} else if strings.TrimRight(inputOption, "\n") == "5" {
			fmt.Printf("Bye bye~")
			return
		}
		// read from server
		buffer := make([]byte, 1024)
		read, err := conn.Read(buffer)
		duration := time.Since(start)
		if err != nil || read == 0 {
			log.Fatalf("Failed to connect to server: %v", err)
			return
		}
		// return microsecond
		fmt.Printf("\nReply from server: %s", string(buffer))
		fmt.Printf("RTT = %fms\n", float64(duration.Nanoseconds())/1e+6)
	}
}
