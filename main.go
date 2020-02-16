package main

import(
	"os"
	"fmt"
)


func main() {
	fileName := "/home/ariel/Desktop/o"
	f, e := os.Create(fileName)

	fmt.Println(e)
	//bufRead := make([]byte,1024)

	i, er := f.Write([]byte("hola "))
	fmt.Println(i,er)

	i, er = f.Write([]byte("mundo"))
	fmt.Println(i,er)
}



