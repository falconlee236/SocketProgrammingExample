/*
SplitFileClient.go
20190532 sang yun lee
*/
package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

func main() {
	// argument handling
	if len(os.Args) != 3 {
		fmt.Println("Invalid argument.")
		os.Exit(1)
	}
	// get argument from system args
	commandName := os.Args[1]
	fileName := os.Args[2]

	// server Info hardcoding
	firstServerName := "127.0.0.1"
	firstServerPort := "40532"
	secondServerName := "127.0.0.1"
	secondServerPort := "50532"

	var wg sync.WaitGroup

	// put command case
	if commandName == "put" {
		wg.Add(2)
		// sendFile odd byte
		go sendFile(fileName, firstServerName, firstServerPort, 0, &wg)
		// sendFile even byte
		go sendFile(fileName, secondServerName, secondServerPort, 1, &wg)
		wg.Wait()
	} else if os.Args[1] == "get" { // get command case
		wg.Add(2)
		// receiveFile odd byte
		go receiveFile(fileName, firstServerName, firstServerPort, 0, &wg)
		// receiveFile even byte
		go receiveFile(fileName, secondServerName, secondServerPort, 1, &wg)
		wg.Wait()
		// get file extension from fileName
		fileExtension := filepath.Ext(fileName)
		// get file name without extension from fileName
		fileName = fileName[0 : len(fileName)-len(fileExtension)]
		// merged file create
		file, err := os.Create(fileName + "-merged" + fileExtension)
		if err != nil {
			fmt.Println("fail to create file:", err)
			os.Exit(1)
		}
		// set Temp file name
		tmpFileName1 := fmt.Sprintf("%s-part%d%stmp%s", fileName, 1, fileExtension, fileExtension)
		tmpFileName2 := fmt.Sprintf("%s-part%d%stmp%s", fileName, 2, fileExtension, fileExtension)
		// open Temp file
		tmpFile1, _ := os.Open(tmpFileName1)
		tmpFile2, _ := os.Open(tmpFileName2)
		// set close Function and delete tmp file
		defer func(tmpFile1 *os.File) {
			err := tmpFile1.Close()
			if err != nil {
				fmt.Println("File close error")
				os.Exit(1)
			}
			err = os.Remove(tmpFileName1)
			if err != nil {
				fmt.Println("File Remove error")
				os.Exit(1)
			}
		}(tmpFile1)
		defer func(tmpFile2 *os.File) {
			err := tmpFile2.Close()
			if err != nil {
				fmt.Println("File close error")
			}
			err = os.Remove(tmpFileName2)
			if err != nil {
				fmt.Println("File Remove error")
			}
		}(tmpFile2)

		// set Buffer from each tmp file
		reader1 := bufio.NewReader(tmpFile1)
		reader2 := bufio.NewReader(tmpFile2)
		// check even, odd
		var byteCnt int64
		var finishCnt int64
		for {
			// all file is finished, exit
			if finishCnt == 2 {
				break
			}
			if byteCnt%2 == 0 {
				// get 1 byte from odd file
				b, err := reader1.ReadByte()
				if err != nil {
					finishCnt++
					continue // finish to reach file end
				}
				_, err = file.Write([]byte{b})
				if err != nil {
					return
				}
			} else {
				// get 1 byte from even file
				b, err := reader2.ReadByte()
				if err != nil {
					finishCnt++
					continue // finish to reach file end
				}
				_, err = file.Write([]byte{b})
				if err != nil {
					return
				}
			}
			byteCnt++
		}

		fmt.Printf("%s%s file merge sucessful!\n", fileName, fileExtension)
	} else { // otherwise case
		fmt.Println("Invalid argument.")
	}
}

// put file to server
func sendFile(fileName string, serverName string, serverPort string, part int, wg *sync.WaitGroup) {
	defer wg.Done()
	// connect to TCP server
	conn, err := net.Dial("tcp", serverName+":"+serverPort)
	if err != nil {
		fmt.Println("fail to connect server: ", err)
		os.Exit(1)
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("socket close error : ", err)
			os.Exit(1)
		}
	}(conn)

	// prepare command string to send command
	commandStr := "put"
	_, err = conn.Write([]byte(commandStr))
	if err != nil {
		fmt.Println("Write failed")
		os.Exit(1)
	}
	commandBuffer := make([]byte, 1024)
	read, _ := conn.Read(commandBuffer)
	// is server response is not ok
	if string(commandBuffer[:read]) != "ok" {
		fmt.Println("fail to receive fileName")
		os.Exit(1)
	}
	fmt.Println("Request to server to put :" + fileName)

	// try to open file
	originalFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println("File Open failed")
		os.Exit(1)
	}
	defer func(originalFile *os.File) {
		err := originalFile.Close()
		if err != nil {
			fmt.Println("File Closed failed")
			os.Exit(1)
		}
	}(originalFile)

	// get file extension from fileName
	fileExtension := filepath.Ext(fileName)
	// get file name without extension from fileName
	fileName = fileName[0 : len(fileName)-len(fileExtension)]
	// join that strings
	fileName = fmt.Sprintf("%s-part%d%s", fileName, part+1, fileExtension)
	// send to target File name
	_, err = conn.Write([]byte(fileName))
	if err != nil {
		fmt.Println("Write failed")
		os.Exit(1)
	}
	fileNameBuffer := make([]byte, 1024)
	read, _ = conn.Read(fileNameBuffer)
	if string(fileNameBuffer[:read]) != "ok" {
		fmt.Println("fail to receive command")
		os.Exit(1)
	}

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
		os.Exit(1)
	}
	fileSizeBuffer := make([]byte, 1024)
	read, _ = conn.Read(fileSizeBuffer)
	if string(fileSizeBuffer[:read]) != "ok" {
		fmt.Println("fail to receive fileSize")
		os.Exit(1)
	}

	// get Buffer from file
	reader := bufio.NewReader(originalFile)
	var byteCount int
	for {
		// get 1 byte from file
		b, err := reader.ReadByte()
		if err != nil {
			break // finish to reach file end
		}

		// check odd, even byte
		if byteCount%2 == part {
			_, err := conn.Write([]byte{b})
			if err != nil {
				fmt.Println("Write failed")
				os.Exit(1)
			}
		}
		byteCount++
	}
	fmt.Printf("%s send successful\n", fileName)
}

func receiveFile(fileName string, serverName string, serverPort string, part int, wg *sync.WaitGroup) {
	// function finished, notify to wait group
	defer wg.Done()

	fmt.Println("Request to server to get :" + fileName)
	// connect to TCP server
	conn, err := net.Dial("tcp", serverName+":"+serverPort)
	if err != nil {
		fmt.Println("fail to connect server: ", err)
		return
	}

	// prepare command string to send command
	commandStr := "get"
	_, err = conn.Write([]byte(commandStr))
	if err != nil {
		fmt.Println("Write failed")
		os.Exit(1)
	}
	commandBuffer := make([]byte, 1024)
	read, _ := conn.Read(commandBuffer)
	// is server response is not ok
	if string(commandBuffer[:read]) != "ok" {
		fmt.Println(string(commandBuffer[:read]))
		os.Exit(1)
	}

	// get file extension from fileName
	fileExtension := filepath.Ext(fileName)
	// get file name without extension from fileName
	fileName = fileName[0 : len(fileName)-len(fileExtension)]
	// join that strings
	fileName = fmt.Sprintf("%s-part%d%s", fileName, part+1, fileExtension)
	// request file Name
	_, err = conn.Write([]byte(fileName))
	if err != nil {
		fmt.Println("Write failed")
		os.Exit(1)
	}

	// try to get fileSize
	fileSizeBuffer := make([]byte, 1024)
	read, _ = conn.Read(fileSizeBuffer)
	fileSize, err := strconv.ParseInt(string(fileSizeBuffer[:read]), 10, 64)
	if err != nil {
		fmt.Println("fail to transfer file size:", err)
		_, err = conn.Write([]byte("fail to transfer file size"))
		if err != nil {
			fmt.Println("Write failed")
			os.Exit(1)
		}
		os.Exit(1)
	} else {
		_, err = conn.Write([]byte("ok"))
		if err != nil {
			fmt.Println("Write failed")
			os.Exit(1)
		}
	}

	file, err := os.Create(fileName + "tmp" + filepath.Ext(fileName))
	if err != nil {
		fmt.Println("fail to create file:", err)
		return
	}
	// received file contents and saved
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

		// write to file from server contents
		if _, err := file.Write(buffer[:n]); err != nil {
			fmt.Println("fail to write file:", err)
			return
		}

		// check file is read finished
		if receivedBytes >= fileSize {
			break
		}
	}
	fmt.Printf("%s file store sucessful!\n", fileName)
}
