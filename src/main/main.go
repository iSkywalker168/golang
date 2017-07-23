package main

import (
	"fmt"
	"hello"
	"flag"
	"os"
)

func main() {
	rootPtr := flag.String("root", ".", "the root directory")
	filterPtr := flag.String("filter", "*", "the filter to dump only expected files or directory")
	flag.Parse()

	logFile, err := os.Create("log.txt")
	check(err)
	defer logFile.Close()

	hello.GetFileInfo(*rootPtr, *filterPtr, func (path string, sha1 string, size int64){
		ln := fmt.Sprintf("\"%v\", %v, %v\n", path, sha1, size)
		fmt.Printf(ln)
		logFile.WriteString(ln)
	})
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}