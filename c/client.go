package main

import(
	"fmt"
	"net"
	"os"
	"strings"
	"./drive"
)

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
		fmt.Println("5 -> look at file list")
		fmt.Scanf("%d", &op)
		switch op {
		case 0:
			exitOp(conn)
		break
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

func downloadOp(conn net.Conn) {
	var fileName string
	fmt.Println("File/Folder path (use '/' to neested folders):\t")
	fmt.Scanf("%s", &fileName)
	fileName = strings.Replace(fileName, "/", string(drive.Sep), -1)
	conn.Write([]byte{2})
	conn.Write([]byte(drive.FillString(fileName,256)))
	drive.DownloadFile(conn,fileName)

	bufferResult := make([]byte,256)
	conn.Read(bufferResult)
	fmt.Println(drive.GetStr(string(bufferResult)))
}

func lookOp(conn net.Conn) {
	var fileName string 
	fmt.Println("File/Folder path (use '/' to neested folders):\t")
	fmt.Scanf("%s", &fileName)
	fileName = strings.Replace(fileName, "/", string(drive.Sep), -1)
	conn.Write([]byte{5})
	conn.Write([]byte(drive.FillString(fileName,256)))

	drive.ShowFiles(conn)
}