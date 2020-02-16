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
	connection, errc := net.Dial("tcp","192.168.100.8:2000")
	if errc != nil {
		panic(errc)
	}
	defer connection.Close()

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
			fmt.Println("See you!")
			os.Exit(0)
		break
		case 1:
			createFolderOp()
		break
		case 2:
			fmt.Println("2")
		break
		case 3:
			fmt.Println("3")
		break
		case 4:
			deleteFileOp()
		break
		default:
			fmt.Println("Other")
		break
		}
	}

	fmt.Println("Connected to server, start receiving the file name and file size")
	bufferFileName := make([]byte,64)
	bufferFileSize := make([]byte,10)

	connection.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

	connection.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName),":")

	newFile, errn := os.Create(fileName)

	if errn != nil {
		panic(errn)
	}
	defer newFile.Close()

	var receivedBytes int64

	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			io.CopyN(newFile, connection, (fileSize - receivedBytes))
			connection.Read(make([]byte, (receivedBytes + BUFFERSIZE) - fileSize))
			break
		}
		io.CopyN(newFile, connection, BUFFERSIZE)
		receivedBytes += BUFFERSIZE
	}

	fmt.Println("Received file completely!")
}

func createFolderOp() {
	var folderName string
	fmt.Print("Folder name (use '/' to neested folders):\t")
	fmt.Scanf("%s", &folderName)
	folderName = strings.Replace(folderName, "/", string(drive.Sep), -1)
	fmt.Println(drive.MakeDirectories("./"+folderName))
}

func deleteFileOp() {
	var fileName string
	fmt.Println("File name (use '/' to neested folders):\t")
	fmt.Scanf("%s", &fileName)
	fileName = strings.Replace(fileName, "/", string(drive.Sep), -1)
	fmt.Println(drive.DeleteFile(fileName))
}
