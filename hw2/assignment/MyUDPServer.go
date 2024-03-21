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

	typeBuffer := make([]byte, 1024)
	buffer := make([]byte, 1024)

	for {
		//pconn, _ := net.ListenPacket("udp", ":"+serverPort)
		//fmt.Printf("Server is ready to receive on port %s\n", serverPort)
		count, r_addr, _ := pconn.ReadFrom(typeBuffer) // from client -1
		if count == 0 {
			return
		}
		typeStr := string(typeBuffer[:count-1])
		fmt.Printf("Command %s\n", typeStr)

		if typeStr == "1" {
			t, in_addr, _ := pconn.ReadFrom(buffer)           // from client - 2
			pconn.WriteTo(bytes.ToUpper(buffer[:t]), in_addr) // to client - 3
		} else if typeStr == "2" {
			pconn.WriteTo([]byte(r_addr.String()), r_addr) // to client - 3
		} else if typeStr == "3" {
			pconn.WriteTo([]byte(strconv.Itoa(reqNum)), r_addr) // to client - 3
		} else if typeStr == "4" {
			duration := time.Since(start)
			hour := int(duration.Seconds() / 3600)
			minute := int(duration.Seconds()/60) % 60
			second := int(duration.Seconds()) % 60
			totalRuntime := fmt.Sprintf("%02d:%02d:%02d\n", hour, minute, second)
			pconn.WriteTo([]byte(totalRuntime), r_addr) // to client - 3
		} else if typeStr == "5" {
			pconn.WriteTo([]byte("-1"), r_addr) // to client - 3
		}
		reqNum++
	}
	//go UDPClientHandler(reqNum, start)
}

func UDPClientHandler(reqNum int, start time.Time) {
	serverPort := "20532"
	pconn, _ := net.ListenPacket("udp", ":"+serverPort)
	fmt.Printf("Server is ready to receive on port %s\n", serverPort)
	typeBuffer := make([]byte, 1024)
	buffer := make([]byte, 1024)

	for {
		//pconn, _ := net.ListenPacket("udp", ":"+serverPort)
		//fmt.Printf("Server is ready to receive on port %s\n", serverPort)
		count, r_addr, _ := pconn.ReadFrom(typeBuffer)
		if count == 0 {
			return
		}
		typeStr := string(typeBuffer[:count-1])
		fmt.Printf("Command %s\n", typeStr)

		if typeStr == "1" {
			t, in_addr, _ := pconn.ReadFrom(buffer)
			pconn.WriteTo(bytes.ToUpper(buffer[:t]), in_addr)
		} else if typeStr == "2" {
			pconn.WriteTo([]byte(r_addr.String()), r_addr)
		} else if typeStr == "3" {
			pconn.WriteTo([]byte(strconv.Itoa(reqNum)), r_addr)
		} else if typeStr == "4" {
			duration := time.Since(start)
			hour := int(duration.Seconds() / 3600)
			minute := int(duration.Seconds()/60) % 60
			second := int(duration.Seconds()) % 60
			totalRuntime := fmt.Sprintf("%02d:%02d:%02d\n", hour, minute, second)
			pconn.WriteTo([]byte(totalRuntime), r_addr)
		} else if typeStr == "5" {
			pconn.WriteTo([]byte("-1"), r_addr)
		}
		reqNum++
	}
}
