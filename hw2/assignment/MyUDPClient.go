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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Print("\nBye bye~\n")
			os.Exit(1)
		}
	}()

	serverName := "nsl2.cau.ac.kr"
	serverPort := "20532"

	pconn, err := net.ListenPacket("udp", ":")
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
		return
	}

	localAddr := pconn.LocalAddr().(*net.UDPAddr)
	fmt.Printf("Client is running on port %d\n", localAddr.Port)

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
		//write to server
		start := time.Now()
		server_addr, _ := net.ResolveUDPAddr("udp", serverName+":"+serverPort)
		pconn.WriteTo([]byte(inputOption), server_addr)

		if strings.TrimRight(inputOption, "\n") == "1" {
			fmt.Printf("Input lowercase sentence: ")
			input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
			start = time.Now()
			pconn.WriteTo([]byte(input), server_addr)
		}
		// read from server
		buffer := make([]byte, 1024)
		read, _, err := pconn.ReadFrom(buffer)
		if err != nil || read == 0 {
			log.Fatalf("Failed to connect to server: %v", err)
			return
		}
		duration := time.Since(start)
		if strings.TrimRight(inputOption, "\n") == "1" {
			fmt.Printf("\nReply from server: %s", string(buffer))
		} else if strings.TrimRight(inputOption, "\n") == "2" {
			recInfo := strings.Split(string(buffer), ":")
			fmt.Printf("\nReply from server: client IP = %s PORT = %s\n", recInfo[0], recInfo[1])
		} else if strings.TrimRight(inputOption, "\n") == "3" {
			fmt.Printf("\nReply from server: requests served = %s\n", string(buffer))
		} else if strings.TrimRight(inputOption, "\n") == "4" {
			fmt.Printf("\nReply from server: run time = %s", string(buffer))
		} else if strings.TrimRight(inputOption, "\n") == "5" {
			fmt.Print("Bye bye~\n")
			return
		}
		fmt.Printf("RTT = %dms\n", duration.Milliseconds())
		pconn.Close()
	}
}
