package main

import(
	"os"
	"path/filepath"
	"io/ioutil"
)


func main() {
	fileName := "/home/ariel/Desktop/A/B/algo.txt"
	filePath, _ := filepath.Split(fileName)
	errm := os.MkdirAll(filePath, os.ModePerm)
	if errm != nil {
		panic(errm)
	}
	newFile, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer newFile.Close()
	ioutil.WriteFile(fileName, []byte("algo"), 0644)
}


