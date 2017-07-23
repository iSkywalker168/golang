package hello

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"io/ioutil"
)

type Callback func(path string, sha1 string, size int64)

func GetSha1(path string) []byte {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return h.Sum(nil)
}

type VisitData struct {
	path string
	info os.FileInfo
	sha1 string
}

func GetFileInfo(root string, filter string, retrieveFileInfo Callback) {
	fileInfos := make(chan VisitData)
	var n sync.WaitGroup
	n.Add(1)
	go walkDir(root, filter, &n, fileInfos)
	go func() {
		n.Wait()
		close(fileInfos)
	}()

	for fileInfo := range fileInfos {
		retrieveFileInfo(fileInfo.path, fileInfo.sha1, fileInfo.info.Size())
	}
}

func walkDir(root string, filter string, n *sync.WaitGroup, fileInfos chan<- VisitData) {
	defer n.Done()
	for _, entry := range dirents(root) {
		if entry.IsDir() {
			n.Add(1)
			subdir := filepath.Join(root, entry.Name())
			go walkDir(subdir, filter, n, fileInfos)
		} else {
			var fileInfo VisitData
			fileInfo.path = filepath.Join(root, entry.Name())

			// support wildcard
			matched, err := filepath.Match(filter, entry.Name())
			if err != nil {
				fmt.Println(err)
			}
			if !matched{
				continue
			}
			fileInfo.info = entry
			fileInfo.sha1 = fmt.Sprintf("%x", GetSha1(fileInfo.path))
			fileInfos <- fileInfo
		}
	}
}

// create 20 go routines at most
var sema = make(chan struct{}, 20)

func dirents(dir string) []os.FileInfo {
	sema <- struct{}{}        // acquire token
	defer func() { <-sema }() // release token

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		return nil
	}
	return entries
}
