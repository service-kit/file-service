package util

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestCompress(t *testing.T) {
	f, err := os.OpenFile("test.wav", os.O_RDONLY, 0777)
	if nil != err {
		t.Error(err)
		return
	}
	defer f.Close()
	buf0 := new(bytes.Buffer)
	io.Copy(buf0, f)
	fmt.Println("befor compress:", buf0.Len())
	data0 := buf0.Bytes()
	out, err := Compress(buf0)
	if nil != err {
		t.Error(err)
		return
	}
	buf1 := new(bytes.Buffer)
	io.Copy(buf1, out)
	fmt.Println("behind compress:", buf1.Len())
	out, err = Uncompress(buf1)
	if nil != err {
		t.Error(err)
		return
	}
	buf2 := new(bytes.Buffer)
	io.Copy(buf2, out)
	fmt.Println("behind uncompress:", buf2.Len())
	if string(buf2.Bytes()) == string(data0) {
		fmt.Println("data is right")
	} else {
		fmt.Println("data is wrong")
	}
}
