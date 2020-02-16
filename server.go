package main

import(
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

const BUFFERSIZE = 1024

func main() {
	server, errs := net.Listen("tcp","192.168.100.8:2000")

	if errs != nil {
		fmt.Println("Error listening: ",errs)
		os.Exit(1)
	}

	defer server.Close()
	fmt.Println("Server started! Waiting for connections...")

	for {
		connection, errc := server.Accept()
		if errc != nil {
			fmt.Println("Error: ",errc)
			os.Exit(1)
		}
		fmt.Println("Client connected")
		go sendFileToClient(connection)
	}
}

func sendFileToClient(connection net.Conn) {
	fmt.Println("A client has connected")
	defer connection.Close()
	file, errf := os.Open("dummyfile.txt")
	if errf != nil {
		fmt.Println(errf)
		return
	}

	fileInfo, erri := file.Stat()
	if erri != nil {
		fmt.Println(errf)
		return
	}

	fileSize := fillString(strconv.FormatInt(fileInfo.Size(),10),10)
	fileName := fillString(fileInfo.Name(),64)
	fmt.Println("Sending file name and file size!")
	connection.Write([]byte(fileSize))
	connection.Write([]byte(fileName))
	sendBuffer := make([]byte, BUFFERSIZE)
	fmt.Println("Start sending file!")
	for {
		_, errF := file.Read(sendBuffer)
		if errF == io.EOF {
			break
		}
		connection.Write(sendBuffer)
	}
	fmt.Println("File has been sent, closing connection!")
	return
}

func fillString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}