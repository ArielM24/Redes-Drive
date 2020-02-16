package main

import(
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"./drive"
)
const sep = os.PathSeparator
const BUFFERSIZE = 1024

func main() {
	server, errs := net.Listen("tcp",":2000")

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
		go readOptions(connection)
	}
}

func readOptions(conn net.Conn) {
	var op int8
	bufferOption := make([]byte,1)
	conn.Read(bufferOption)
	op = int8(bufferOption[0])
	switch op {
		case 0:
			fmt.Println("Connection finished!")
			conn.Close()
		break
		case 1:
			fmt.Println("Creating folder")
			createFolderOp(conn)
		break
		case 2:
			fmt.Println("2")
		break
		case 3:
			//uploadOp(connection)
		break
		case 4:
			//deleteFileOp()
		break
		default:
			fmt.Println("Other")
		break
		}
}

func createFolderOp(conn net.Conn){
	bufferName := make([]byte,256)
	conn.Read(bufferName)
	folderName := drive.GetStr(string(bufferName))
	result := drive.FillString(drive.MakeDirectories("."+string(drive.Sep)+folderName),256)
	fmt.Println(drive.GetStr(result))
	conn.Write([]byte(result))
}

func sendFileToClient(connection net.Conn) {
	fmt.Println("A client has connected")
	defer connection.Close()
	file, errf := os.Open("Carpeta"+string(sep)+"dummyfile.txt")
	if errf != nil {
		fmt.Println(errf)
		return
	}

	fileInfo, erri := file.Stat()
	if erri != nil {
		fmt.Println(errf)
		return
	}

	fileSize := drive.FillString(strconv.FormatInt(fileInfo.Size(),10),10)
	fileName := drive.FillString(fileInfo.Name(),64)
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
