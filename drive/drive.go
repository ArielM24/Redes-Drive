package drive

import(
	"fmt"
	"os"
	"path/filepath"
	"net"
	"strings"
	"strconv"
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

func getFiles(path string, dir bool) []string{
	paths := Paths(path)
	result := make([]string,0)
	for _, p := range paths {
		f, _ := os.Open(p)
		s, _ := f.Stat()
		if s.IsDir() {
			if dir {
				result = append(result, p)
			} 
		} else {
			if !dir {
				result = append(result, p)
			}
		}
	}
	return result
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

func MakePaths(path string) {
	MakeDirectories(path)
	paths := Paths(path)
	for _, d := range paths {
		p, _ := filepath.Split(d)
		MakeDirectories(p)
	}
}

func UploadFile(conn net.Conn, path string) bool {
	files := getFiles(path,false)
	result := true
	uploadPaths(conn, path)
	nf := FillString(strconv.FormatInt(int64(len(files)),10),64)
	conn.Write([]byte(nf))
	fmt.Println("Uploading files...")
	for _, f := range files {
		result = result && upload(conn, f)
	}
	return result
}

func uploadPaths(conn net.Conn, path string) bool {
	paths := getFiles(path,true)
	np := FillString(strconv.FormatInt(int64(len(paths)),10),64)
	conn.Write([]byte(np))
	for _, p := range paths {
		conn.Write([]byte(FillString(p,256)))
	}
	return true
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

	conn.Write([]byte(fileSize))
	conn.Write([]byte(fileName))

	var np, rest, i int64

	fsize := fileInfo.Size()

	np = int64(fsize / BUFFERSIZE)
	rest = fsize % BUFFERSIZE

	var sendBuffer []byte
	for i = 0; i < np; i++ {
		sendBuffer = make([]byte, BUFFERSIZE)
		file.Read(sendBuffer)
		conn.Write(sendBuffer)
	}

	if rest > 0 {
		sendBuffer = make([]byte, rest)
		file.Read(sendBuffer)
		conn.Write(sendBuffer)
	}
	
	return true
}

func DownloadFile(conn net.Conn, file string) bool {
	downloadPaths(conn)

	buffNf := make([]byte,64)

	conn.Read(buffNf)
	nf ,_:= strconv.ParseInt(strings.Trim(string(buffNf),":"),10,64)
	var i int64
	fmt.Println("Dowloading files...")
	for i = 0; i < nf; i++ {
		download(conn)
	}
	return true
}

func downloadPaths(conn net.Conn){
	buffNp := make([]byte,64)
	bufferName := make([]byte,256)

	conn.Read(buffNp)
	np ,_:= strconv.ParseInt(strings.Trim(string(buffNp),":"),10,64)
	var i int64
	for i = 0; i < np; i++ {
		conn.Read(bufferName)
		name := strings.Trim(string(bufferName),":")
		os.MkdirAll(name, os.ModePerm)
	}
}

func download(conn net.Conn) bool {
	r := true
	bufferFileName := make([]byte,256)
	bufferFileSize := make([]byte,64)

	conn.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
	conn.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName),":")
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
	}

	if rest > 0 {
		buffReceived = make([]byte, rest)
		conn.Read(buffReceived)
		newFile.Write(buffReceived)
	}

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

func LookFiles(conn net.Conn, path string) {
	paths := Paths(path)
	np := FillString(strconv.FormatInt(int64(len(paths)),10),64)
	conn.Write([]byte(np))
	for _, p := range paths {
		name := FillString(p,256)
		conn.Write([]byte(name))
	}
}

func ShowFiles(conn net.Conn) {
	buffNp := make([]byte,64)
	conn.Read(buffNp)
	np, _ := strconv.ParseInt(strings.Trim(string(buffNp), ":"), 10, 64)
	bufferName := make([]byte,256)
	var i int64
	if np == 0 {
		fmt.Println("No files to show!")
	} else {
		fmt.Println("Files:")
	}
	for i = 0; i < np; i++ {
		conn.Read(bufferName)
		name := strings.Trim(string(bufferName),":")
		fmt.Println(name)
	}
	fmt.Println()
}