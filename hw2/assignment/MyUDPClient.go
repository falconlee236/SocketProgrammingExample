/**
 * MyUDPClient.go
20190532 Sangyun Lee
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
	serverName := "localhost"
	serverPort := "30532"

	// connect to server
	pconn, err := net.ListenPacket("udp", ":")
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
		return
	}

	// get server address
	localAddr := pconn.LocalAddr().(*net.UDPAddr)
	fmt.Printf("Client is running on port %d\n", localAddr.Port)

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
		if inputNum < 1 || inputNum > 5 {
			fmt.Print("Invalid option\n")
			continue
		}
		// start calculate RTT
		start := time.Now()

		// get datagram to severName
		server_addr, _ := net.ResolveUDPAddr("udp", serverName+":"+serverPort)
		pconn.WriteTo([]byte(inputOption), server_addr) // to server - 1
		if strings.TrimRight(inputOption, "\n") == "1" {
			fmt.Printf("Input lowercase sentence: ")
			input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
			start = time.Now()
			s_addr, _ := net.ResolveUDPAddr("udp", serverName+":"+serverPort)
			pconn.WriteTo([]byte(input), s_addr) // to server - 2
		} else if strings.TrimRight(inputOption, "\n") == "5" {
			fmt.Println("Bye bye~")
			return
		}

		// read from server
		// set timeout server doesn't reply 10 sec
		pconn.SetReadDeadline(time.Now().Add(10 * time.Second))
		buffer := make([]byte, 1024)
		read, _, err := pconn.ReadFrom(buffer) //from server - 3
		if err != nil || read == 0 {
			log.Fatalf("Failed to connect to server: %v", err)
			return
		}
		duration := time.Since(start)
		fmt.Printf("Reply from server: %s\n", string(buffer))
		fmt.Printf("RTT = %fms\n", float64(duration.Nanoseconds())/1e+6)
	}
}
