package hello

import (
	"testing"
	"runtime"
	"fmt"
)

func TestSha1(t *testing.T){
	_, filename, _, ok := runtime.Caller(1)
	if !ok{
		t.Error(`failed to invoke runtime.Caller`)
	}
	fmt.Printf(filename)
	sha1 := GetSha1("hello_test.go")
	if 0 == len(sha1){
		t.Error("failed to get sha1 of %v", filename)
	}
}
func retrieveSingleFileInfo(path string, sha1 string, size int64){
	fmt.Printf("%v, %v, %v\n", path, sha1, size)
}

func TestGetFileInfo(t *testing.T) {
	GetFileInfo(".", "*", func (path string, sha1 string, size int64){
		if 0 == len(path) || 0 == len(sha1){
			t.Error("either file path or sha1 is empty")
		}
	})
}