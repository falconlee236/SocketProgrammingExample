/**
 * TCPServer.go
 **/

package main

import (
	"bytes"
	"fmt"
	"net"
)

func main() {
	serverPort := "20532"

	listener, _ := net.Listen("tcp", ":"+serverPort)
	fmt.Printf("Server is ready to receive on port %s\n", serverPort)
	typebuffer := make([]byte, 1024)
	buffer := make([]byte, 1024)

	for {
		conn, _ := listener.Accept()
		t, _ := conn.Read(typebuffer)
		fmt.Printf("Connection request from %s\n", conn.RemoteAddr().String())
		count, _ := conn.Read(buffer)
		conn.Write(bytes.ToUpper(buffer[:count]))
		fmt.Printf("Command %s\n", typebuffer[:t])
		conn.Close()
	}
}
