/*
 * MyTCPClient.go
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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Print("\nBye bye~\n")
			os.Exit(1)
		}
	}()

	serverName := "localhost"
	serverPort := "20532"

	conn, err := net.Dial("tcp", serverName+":"+serverPort)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
		return
	}

	localAddr := conn.LocalAddr().(*net.TCPAddr)
	fmt.Printf("Client is running on port %d\n", localAddr.Port)
	defer conn.Close()

	for {
		fmt.Printf("<Menu>\n")
		fmt.Printf("1) convert text to UPPER-case\n")
		fmt.Printf("2) get my IP address and port number\n")
		fmt.Printf("3) get server request count\n")
		fmt.Printf("4) get server running time\n")
		fmt.Printf("5) exit\n")
		fmt.Printf("Input option: ")
		inputOption, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		inputNum, _ := strconv.Atoi(strings.TrimRight(inputOption, "\n"))
		if inputNum < 1 || inputNum > 5 {
			fmt.Print("Invalid option\n")
			continue
		}
		start := time.Now()
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
		buffer := make([]byte, 1024)
		read, err := conn.Read(buffer)
		duration := time.Since(start)
		if err != nil || read == 0 {
			log.Fatalf("Failed to connect to server: %v", err)
			return
		}
		fmt.Printf("\nReply from server: %s", string(buffer))
		fmt.Printf("RTT = %dms\n", duration.Milliseconds())
	}
}
