/*
 * TCPClient.go
 **/

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {

	serverName := "nsl5.cau.ac.kr"
	serverPort := "20532"

	conn, _ := net.Dial("tcp", serverName+":"+serverPort)

	localAddr := conn.LocalAddr().(*net.TCPAddr)
	fmt.Printf("Client is running on port %d\n", localAddr.Port)

	fmt.Printf("<Menu>\n")
	fmt.Printf("1) convert text to UPPER-case\n")
	fmt.Printf("2) get my IP address and port number\n")
	fmt.Printf("3) get server request count\n")
	fmt.Printf("4) get server running time\n")
	fmt.Printf("5) exit\n")
	fmt.Printf("Input option: ")
	input_option, _ := bufio.NewReader(os.Stdin).ReadString('\n')
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
		fmt.Printf("\nReply from server: client IP = %s PORT = %s", recInfo[0], recInfo[1])
	} else if strings.TrimRight(input_option, "\n") == "3" {
		fmt.Printf("\nReply from server: %s", string(buffer))
	} else if strings.TrimRight(input_option, "\n") == "4" {
		fmt.Printf("\nReply from server: run time = %s", string(buffer))
	}
	fmt.Printf("RTT = %s\n", duration)
	conn.Close()
}
