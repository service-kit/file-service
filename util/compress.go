package util

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
)

func Compress(in io.Reader) (out io.Reader, err error) {
	buf := new(bytes.Buffer)
	zw := gzip.NewWriter(buf)
	if nil == zw {
		return nil, errors.New("new writer fail")
	}
	defer zw.Close()
	_, err = io.Copy(zw, in)
	if nil != err {
		return nil, err
	}
	return buf, err
}

func Uncompress(in io.Reader) (out io.Reader, err error) {
	return gzip.NewReader(in)
}

type Compresser struct {
	buf *bytes.Buffer
	zw  io.Writer
}

func (c *Compresser) Init() error {
	c.buf = new(bytes.Buffer)
	c.zw = gzip.NewWriter(c.buf)
	return nil
}

func (c *Compresser) Read(out []byte) (int, error) {
	return c.buf.Read(out)
}

func (c *Compresser) Write(in []byte) (int, error) {
	return c.zw.Write(in)
}

func (c *Compresser) ReadFrom(r io.Reader) (n int64, err error) {
	return io.Copy(c.zw, r)
}

func (c *Compresser) WriteTo(w io.Writer) (n int64, err error) {
	return io.Copy(w, c.buf)
}

type Uncompresser struct {
	buf *bytes.Buffer
	zr  io.Reader
}

func (uc *Uncompresser) Init() (err error) {
	uc.buf = new(bytes.Buffer)
	uc.zr, err = gzip.NewReader(uc.buf)
	return
}

func (uc *Uncompresser) Read(out []byte) (int, error) {
	return uc.zr.Read(out)
}

func (uc *Uncompresser) Write(in []byte) (int, error) {
	return uc.buf.Write(in)
}

func (uc *Uncompresser) ReadFrom(r io.Reader) (n int64, err error) {
	return io.Copy(uc.buf, r)
}

func (uc *Uncompresser) WriteTo(w io.Writer) (n int64, err error) {
	return io.Copy(w, uc.zr)
}

func NewCompresser() (*Compresser, error) {
	c := new(Compresser)
	c.Init()
	return c, nil
}

func NewUncompresser() (*Uncompresser, error) {
	uc := new(Uncompresser)
	err := uc.Init()
	if nil != err {
		return nil, err
	}
	return uc, err
}
