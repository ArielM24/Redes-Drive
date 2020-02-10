package main

import(
	"fmt"
	"os"
	//"path/filepath"
)

const sep = os.PathSeparator

func main() {
	name := "/home/ariel/Documents/redes/p1/a"
	/*fi, _ := os.Stat(name)
	if fi.IsDir() {
		fmt.Println("Dir")
	}else {
		fmt.Println("File")
	}

	f, _ := os.Open(name)
	files, _ := f.Readdirnames(0)
	for _, n := range files {
		fmt.Println(n)
	}*/

	names := paths(name)
	fmt.Println("Paths: ",names)
}

func exitOnError(err error) {
	if err != nil {
		fmt.Println("Error",err)
		os.Exit(1)
	}
}

func paths(path string) []string {
	names := make([]string, 0)

	f, errf := os.Open(path)
	exitOnError(errf)

	sf, errs := f.Stat()
	exitOnError(errs)

	if sf.IsDir() {
		names = append(names,f.Name())
		subf, errsf := f.Readdirnames(0)
		exitOnError(errsf)

		for _, s := range subf {
			names = append(names,paths(path + string(sep) + s)...)
		}

	} else {
		names = append(names, path)
	}

	return names
}