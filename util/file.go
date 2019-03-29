package util

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func GetFileModTime(path string) (int64, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, errors.New("open file err!!!")
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		log.Println("stat fileinfo error")
		return 0, err
	}
	return fi.ModTime().Unix(), nil
}

func SaveFile(file string, buf io.Reader) {
	f, e := os.OpenFile(file, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0777)
	if nil != e {
		return
	}
	defer f.Close()
	i, e := io.Copy(f, buf)
	if nil != e {
		return
	}
	fmt.Println("Save File ", file, " Size:", i)
}
