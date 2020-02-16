package drive

import(
	"fmt"
	"os"
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