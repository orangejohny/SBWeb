package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	file, _ := os.Create("file.txt")
	defer os.Remove("file.txt")
	defer file.Close()
	file2, _ := os.Open("file.txt")
	defer file2.Close()
	file.Write([]byte{1, 2, 3})
	func(r io.Reader) {
		buf := make([]byte, 10)
		r.Read(buf)
		fmt.Print(buf)
	}(file2)

}
