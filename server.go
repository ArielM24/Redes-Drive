package main

import(
	"fmt"
	"net"
	"os"
	"./drive"
)

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
		case 5:
			lookOp(conn)
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
	r := drive.UploadFile(conn,fileName)
	if r {
		conn.Write([]byte(drive.FillString("Files donwload succesfuly!",256)))
	} else {
		conn.Write([]byte(drive.FillString("Error while downloading files!",256)))
	}
	
}

func lookOp(conn net.Conn){
	bufferName := make([]byte,256)
	conn.Read(bufferName)
	fileName := drive.GetStr(string(bufferName))
	drive.LookFiles(conn,fileName)
}
