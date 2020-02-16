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
	fmt.Println("Server started! Waiting for connections...")

	for {
		connection, errc := server.Accept()
		if errc != nil {
			fmt.Println("Error: ",errc)
			os.Exit(1)
		}
		fmt.Println("Client connected")
		readOptions(connection)
	}
}

func readOptions(conn net.Conn) {
	var op int8 = 5
	bufferOption := make([]byte,1)
	for op != 0 {
		conn.Read(bufferOption)
		op = int8(bufferOption[0])
		fmt.Println(op)
		switch op {
		case 1:
			createFolderOp(conn)
		break
		case 2:
			downloadOp(conn)
		break
		case 3:
			uploadOp(conn)
		break
		case 4:
			deleteFileOp(conn)
		break
		default:
			fmt.Println("Nothing")
		break
		}
	}
	fmt.Println("Connection finished!")
	conn.Close()
}

func createFolderOp(conn net.Conn) {
	bufferName := make([]byte,256)
	conn.Read(bufferName)
	folderName := drive.GetStr(string(bufferName))
	result := drive.FillString(drive.MakeDirectories(folderName),256)
	fmt.Println(drive.GetStr(result))
	conn.Write([]byte(result))
}

func deleteFileOp(conn net.Conn) {
	bufferName := make([]byte,256)
	conn.Read(bufferName)
	fileName := drive.GetStr(string(bufferName))
	result := drive.FillString(drive.DeleteFile(fileName),256)
	fmt.Println(drive.GetStr(result))
	conn.Write([]byte(result))
}

func uploadOp(conn net.Conn) {
	bufferName := make([]byte,256)
	conn.Read(bufferName)
	fileName := drive.GetStr(string(bufferName))
	fmt.Println("u",fileName)
	r := drive.DownloadFile(conn,fileName)
	if r {
		conn.Write([]byte(drive.FillString("Files upload succesfuly!",256)))
	} else {
		conn.Write([]byte(drive.FillString("Error while uploading files!",256)))
	}
}

func downloadOp(conn net.Conn) {
	bufferName := make([]byte,256)
	conn.Read(bufferName)
	fileName := drive.GetStr(string(bufferName))
	fmt.Println("u",fileName)
	r := drive.UploadFile(conn,fileName)
	if r {
		conn.Write([]byte(drive.FillString("Files donwload succesfuly!",256)))
	} else {
		conn.Write([]byte(drive.FillString("Error while downloading files!",256)))
	}
	
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
