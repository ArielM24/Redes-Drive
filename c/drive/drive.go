package drive

import(
	"fmt"
	"os"
	"path/filepath"
	"net"
	"strings"
	"strconv"
	//"io"
)

const Sep = os.PathSeparator
const BUFFERSIZE = 1024

func ExitOnError(err error) {
	if err != nil {
		fmt.Println("Error",err)
		os.Exit(1)
	}
}

func Paths(path string) []string {
	names := make([]string, 0)

	f, errf := os.Open(path)
	ExitOnError(errf)

	sf, errs := f.Stat()
	ExitOnError(errs)

	if sf.IsDir() {
		names = append(names,f.Name())
		subf, errsf := f.Readdirnames(0)
		ExitOnError(errsf)

		for _, s := range subf {
			names = append(names,Paths(path + string(Sep) + s)...)
		}

	} else {
		names = append(names, path)
	}

	return names
}

func MakeDirectories(path string) string {
	errm := os.MkdirAll(path, os.ModePerm)
	if errm != nil {
		return "Error while creating directories"
	} else {
		return "Directories created succesfully"
	}
}

func DeleteFile(path string) string {
	errr := os.RemoveAll(path)
	if errr != nil {
		return "Error while deleting directories"
	} else {
		return "File deleted succesfully"
	}
}

func MakePaths(path string) string {
	MakeDirectories(path)
	paths := Paths(path)
	for _, d := range paths {
		p, _ := filepath.Split(d)
		MakeDirectories(p)
	}
	return "m"
}

func UploadFile(conn net.Conn, path string) bool {
	paths := Paths(path)
	fmt.Println("U",paths)
	result := true
	for _, p := range paths {
		fi, erri := os.Open(p)
		if erri != nil {
			return false
		}
		fs, errs := fi.Stat()
		if errs != nil {
			return false
		}
		if !fs.IsDir() {
			result = result && upload(conn, p)
		} else {
			uploadPath(conn, p)
		}
		fmt.Println("a")
	}
	return result
}

func uploadPath(conn net.Conn, path string) {
	
}

func upload(conn net.Conn, filePath string) bool{
	file, errf := os.Open(filePath)
	defer file.Close()
	if errf != nil {
		fmt.Println(errf)
		return false
	}

	fileInfo, erri := file.Stat()
	if erri != nil {
		fmt.Println(errf)
		return false
	}

	fileSize := FillString(strconv.FormatInt(fileInfo.Size(),10),64)
	fileName := FillString(filePath,256)

	fmt.Println("u s",fileSize)
	fmt.Println("u n",fileName)
	fmt.Println("Sending file name and file size!")
	conn.Write([]byte(fileSize))
	conn.Write([]byte(fileName))
	fmt.Println("Start sending file!")

	var np, rest, i int64

	fsize := fileInfo.Size()

	np = int64(fsize / BUFFERSIZE)
	rest = fsize % BUFFERSIZE

	var sendBuffer []byte
	for i = 0; i < np; i++ {
		sendBuffer = make([]byte, BUFFERSIZE)
		file.Read(sendBuffer)
		conn.Write(sendBuffer)
		fmt.Println(string(sendBuffer))
	}

	if rest > 0 {
		sendBuffer = make([]byte, rest)
		file.Read(sendBuffer)
		conn.Write(sendBuffer)
		fmt.Println(string(sendBuffer))
	}
	
	fmt.Println("File has been sent, closing connection!")
	return true

}

func DownloadFile(conn net.Conn, file string) bool {
	MakePaths(file)
	paths := Paths(file)
	fmt.Println("D",paths)
	result := true
	for _, p := range paths {
		fmt.Println("aaa ",p)
		result = result && download(conn, p)
	}
	return result
}

func download(conn net.Conn, file string) bool {
	r := true
	fmt.Println("Connected to server, start receiving the file name and file size")
	bufferFileName := make([]byte,256)
	bufferFileSize := make([]byte,64)

	conn.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
	fmt.Println("d s",fileSize)
	conn.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName),":")
	fmt.Println("d n",fileName)
	p, _ := filepath.Split(fileName)
	MakeDirectories(p)
	newFile, errn := os.Create(fileName)

	if errn != nil {
		r = false
		fmt.Println(errn)
	}
	defer newFile.Close()

	var np, rest, i int64

	np = int64(fileSize / BUFFERSIZE)
	rest = fileSize % BUFFERSIZE

	var buffReceived []byte
	for i = 0; i < np; i++ {
		buffReceived = make([]byte, BUFFERSIZE)
		conn.Read(buffReceived)
		newFile.Write(buffReceived)
		fmt.Println(string(buffReceived))
	}

	if rest > 0 {
		buffReceived = make([]byte, rest)
		conn.Read(buffReceived)
		newFile.Write(buffReceived)
		fmt.Println(string(buffReceived))
	}

	fmt.Println("Received file completely!")
	return r
}

func FillString(retunString string, toLength int) string {
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

func GetStr(str string) string{
	return strings.Replace(str, ":", "", -1)
}