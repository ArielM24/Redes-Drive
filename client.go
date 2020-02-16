package main

import(
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"./drive"
)

const BUFFERSIZE = 1024

func main() {
	var op int8
	conn, errc := net.Dial("tcp","192.168.100.8:2000")
	if errc != nil {
		panic(errc)
	}

	defer conn.Close()

	for {
		fmt.Println("Selec an option")
		fmt.Println("0 -> exit")
		fmt.Println("1 -> create folder")
		fmt.Println("2 -> dowload file/folder")
		fmt.Println("3 -> upload file/folder")
		fmt.Println("4 -> delete file/folder")
		fmt.Scanf("%d", &op)
		switch op {
		case 0:
			exitOp(conn)
		break
		case 1:
			createFolderOp(conn)
		break
		case 2:
			fmt.Println("2")
		break
		case 3:
			uploadOp(conn)
		break
		case 4:
			deleteFileOp(conn)
		break
		default:
			fmt.Println("Other")
		break
		}
	}

	fmt.Println("Connected to server, start receiving the file name and file size")
	bufferFileName := make([]byte,64)
	bufferFileSize := make([]byte,10)

	conn.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

	conn.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName),":")

	newFile, errn := os.Create(fileName)

	if errn != nil {
		panic(errn)
	}
	defer newFile.Close()

	var receivedBytes int64

	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			io.CopyN(newFile, conn, (fileSize - receivedBytes))
			conn.Read(make([]byte, (receivedBytes + BUFFERSIZE) - fileSize))
			break
		}
		io.CopyN(newFile, conn, BUFFERSIZE)
		receivedBytes += BUFFERSIZE
	}

	fmt.Println("Received file completely!")
}

func exitOp(conn net.Conn) {
	fmt.Println("See you!")
	conn.Write([]byte{0})
	os.Exit(0)
}

func createFolderOp(conn net.Conn) {
	var folderName string
	fmt.Print("Folder name (use '/' to neested folders):\t")
	fmt.Scanf("%s", &folderName)
	folderName = drive.FillString(strings.Replace(folderName, "/", string(drive.Sep), -1),256)
	conn.Write([]byte{1})
	conn.Write([]byte(folderName))
	bufferResult := make([]byte,256)
	conn.Read(bufferResult)
	fmt.Println(drive.GetStr(string(bufferResult)))
}

func deleteFileOp(conn net.Conn) {
	var fileName string
	fmt.Println("File name (use '/' to neested folders):\t")
	fmt.Scanf("%s", &fileName)
	fileName = drive.FillString(strings.Replace(fileName, "/", string(drive.Sep), -1), 256)
	conn.Write([]byte{4})
	conn.Write([]byte(fileName))
	bufferResult := make([]byte,256)
	conn.Read(bufferResult)
	fmt.Println(drive.GetStr(string(bufferResult)))
}

func uploadOp(conn net.Conn) {
	var fileName string
	fmt.Println("File/Folder path (use '/' to neested folders):\t")
	fmt.Scanf("%s", &fileName)
	fileName = strings.Replace(fileName, "/", string(drive.Sep), -1)
	conn.Write([]byte{3})
	conn.Write([]byte(drive.FillString(fileName,256)))
	drive.UploadFile(conn,fileName)

	bufferResult := make([]byte,256)
	conn.Read(bufferResult)
	fmt.Println(drive.GetStr(string(bufferResult)))
}
