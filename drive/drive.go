package drive

import(
	"fmt"
	"os"
	"path/filepath"
	"net"
	"strings"
)

const Sep = os.PathSeparator

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

func UploadFile(conn net.Conn , path string) string {
	paths := Paths(path)
	for _, d := range paths {
		p, _ := filepath.Split(d)
		MakeDirectories("."+string(Sep)+p)
	}

	_, errw := conn.Write([]byte("u"))
	if errw != nil {
		return "Error while uploading file" 
	}

	return "u"
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