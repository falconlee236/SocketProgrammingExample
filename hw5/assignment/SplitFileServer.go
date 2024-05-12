/*
SplitFileServer.go
20190532 sang yun lee
*/

package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("This program only accepts one argument")
		os.Exit(1)
	}
	//server Port
	serverPort := os.Args[1]

	// listen from client request
	listener, err := net.Listen("tcp", ":"+serverPort)
	if err != nil {
		fmt.Println("fail to start server: ", err)
		return
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println("fail to close listener: ", err)
		}
	}(listener)
	fmt.Println("Waiting for connections...")

	// main loop
	for {
		// accept to client request
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("fail to accept client request: ", err)
			continue
		}

		// get command from client
		commandBuffer := make([]byte, 1024)
		read, err := conn.Read(commandBuffer)
		if err != nil {
			_, err = conn.Write([]byte("fail to read file"))
			if err != nil {
				fmt.Println("Write failed")
				continue
			}
		} else {
			_, err = conn.Write([]byte("ok"))
			if err != nil {
				fmt.Println("Write failed")
				continue
			}
		}
		commandName := string(commandBuffer[:read])

		// put case
		if commandName == "put" {
			// get file Name from client
			fileNameBuffer := make([]byte, 1024)
			read, err = conn.Read(fileNameBuffer)
			if err != nil {
				_, err = conn.Write([]byte("fail to read file"))
				if err != nil {
					fmt.Println("Write failed")
					continue
				}
			} else {
				_, err = conn.Write([]byte("ok"))
				if err != nil {
					fmt.Println("Write failed")
					continue
				}
			}
			fileName := string(fileNameBuffer[:read])
			fmt.Println("Received file: ", fileName)

			// get file size from client
			fileSizeBuffer := make([]byte, 1024)
			read, err = conn.Read(fileSizeBuffer)
			if err != nil {
				_, err = conn.Write([]byte("fail to read file"))
				if err != nil {
					fmt.Println("Write failed")
					continue
				}
			} else {
				_, err = conn.Write([]byte("ok"))
				if err != nil {
					fmt.Println("Write failed")
					continue
				}
			}
			// translate string to int
			fileSize, err := strconv.ParseInt(string(fileSizeBuffer[:read]), 10, 64)
			if err != nil {
				fmt.Println("fail to transfer file size:", err)
				continue
			}

			// create file
			file, err := os.Create(fileName)
			if err != nil {
				fmt.Println("fail to create file:", err)
				continue
			}

			// save file contents
			var receivedBytes int64
			buffer := make([]byte, 1024)
			for {
				n, err := conn.Read(buffer)
				if err != nil {
					if err != io.EOF {
						fmt.Println("fail to read file:", err)
					}
					break
				}
				receivedBytes += int64(n)

				// write content to file
				if _, err := file.Write(buffer[:n]); err != nil {
					fmt.Println("fail to write file:", err)
					break
				}

				// check all content is done
				if receivedBytes >= fileSize {
					break
				}
			}

			fmt.Printf("%s file store sucessful!\n", fileName)
			err = file.Close()
			if err != nil {
				fmt.Println("file closed error")
				continue
			}
		} else if commandName == "get" { // get case
			// get fileName
			fileNameBuffer := make([]byte, 1024)
			read, _ = conn.Read(fileNameBuffer)
			fileName := string(fileNameBuffer[:read])
			// try to open file
			originalFile, err := os.Open(fileName)
			// error occur
			if err != nil {
				_, err = conn.Write([]byte("fail to open file"))
				if err != nil {
					fmt.Println("Write failed")
					continue
				}
			} else {
				// try to get file Stats
				fileInfo, err := originalFile.Stat()
				if err != nil {
					fmt.Println("fail to get file stat:", err)
					return
				}
				// convert file Size string to Integer
				size := strconv.FormatInt(fileInfo.Size(), 10)
				// send file Size to server
				_, err = conn.Write([]byte(size))
				if err != nil {
					fmt.Println("Write failed")
					continue
				}
			}
			fmt.Println("Request from client to send :" + fileName)

			statusBuffer := make([]byte, 1024)
			read, _ = conn.Read(statusBuffer)
			if string(statusBuffer[:read]) != "ok" {
				continue
			}

			// get Buffer from file
			reader := bufio.NewReader(originalFile)
			for {
				// get 1 byte from file
				b, err := reader.ReadByte()
				if err != nil {
					break // finish to reach file end
				}
				_, err = conn.Write([]byte{b})
				if err != nil {
					fmt.Println("Write Error")
					break
				}
			}
			fmt.Printf("%s send successful\n", fileName)
		} else {
			fmt.Println("unknown command: ", commandName)
		}
	}
}
