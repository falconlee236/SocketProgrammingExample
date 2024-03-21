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
	"strings"
	"time"
)

func main() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("captured %v\n", sig)
			fmt.Print("Bye bye~\n")
			os.Exit(1)
		}
	}()

	serverName := "nsl5.cau.ac.kr"
	serverPort := "20532"

	conn, _ := net.Dial("tcp", serverName+":"+serverPort)

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
		input_option, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		fmt.Printf("%s\n", input_option)
		start := time.Now()
		conn.Write([]byte(input_option))

		if strings.TrimRight(input_option, "\n") == "1" {
			fmt.Printf("Input lowercase sentence: ")
			input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
			start = time.Now()
			conn.Write([]byte(input))
		}
		buffer := make([]byte, 1024)
		conn.Read(buffer)
		duration := time.Since(start)
		if strings.TrimRight(input_option, "\n") == "1" {
			fmt.Printf("\nReply from server: %s", string(buffer))
		} else if strings.TrimRight(input_option, "\n") == "2" {
			recInfo := strings.Split(string(buffer), ":")
			fmt.Printf("\nReply from server: client IP = %s PORT = %s\n", recInfo[0], recInfo[1])
		} else if strings.TrimRight(input_option, "\n") == "3" {
			fmt.Printf("\nReply from server: requests served = %s\n", string(buffer))
		} else if strings.TrimRight(input_option, "\n") == "4" {
			fmt.Printf("\nReply from server: run time = %s", string(buffer))
		} else if strings.TrimRight(input_option, "\n") == "5" {
			fmt.Print("Bye bye~\n")
		}
		fmt.Printf("RTT = %s\n", duration)
	}
}
