/**
 * MyUDPServer.go
20190532 Sangyun Lee
 **/

package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
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

	start := time.Now()
	reqNum := 0
	serverPort := "20532"

	pconn, _ := net.ListenPacket("udp", ":"+serverPort)
	fmt.Printf("Server is ready to receive on port %s\n", serverPort)

	go UDPClientHandler(pconn, reqNum, start)
}

func UDPClientHandler(conn net.PacketConn, reqNum int, start time.Time) {

	typeBuffer := make([]byte, 1024)
	buffer := make([]byte, 1024)

	for {
		count, r_addr, _ := conn.ReadFrom(typeBuffer)
		if count == 0 {
			return
		}
		typeStr := string(typeBuffer[:count-1])
		fmt.Printf("Command %s\n", typeStr)

		if typeStr == "1" {
			t, in_addr, _ := conn.ReadFrom(buffer)
			conn.WriteTo(bytes.ToUpper(buffer[:t]), in_addr)
		} else if typeStr == "2" {
			conn.WriteTo([]byte(r_addr.String()), r_addr)
		} else if typeStr == "3" {
			conn.WriteTo([]byte(strconv.Itoa(reqNum)), r_addr)
		} else if typeStr == "4" {
			duration := time.Since(start)
			hour := int(duration.Seconds() / 3600)
			minute := int(duration.Seconds()/60) % 60
			second := int(duration.Seconds()) % 60
			totalRuntime := fmt.Sprintf("%02d:%02d:%02d\n", hour, minute, second)
			conn.WriteTo([]byte(totalRuntime), r_addr)
		} else if typeStr == "5" {
			conn.WriteTo([]byte("-1"), r_addr)
		}
		reqNum++
	}
}
