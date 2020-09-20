package main

import (
	"errors"
	"fmt"
	"io"
	"time"
)

func main() {
	Pipe()
}

func Pipe() {
	pipeReader, pipeWriter := io.Pipe()
	go PipeWrite(pipeWriter)
	go PipeRead(pipeReader)
	time.Sleep(1e7)
}

func PipeWrite(pipeWriter *io.PipeWriter) {
	var (
		i   = 0
		err error
		n int
	)
	data := []byte("Go语言学习园地")
	for _, err = pipeWriter.Write(data); err == nil; n, err = pipeWriter.Write(data) {
		i++
		if i == 3 {
			pipeWriter.CloseWithError(errors.New("输出3次后结束"))
		}
	}
	fmt.Println("close 后输出的字节数：", n, " error：",  err)
}

func PipeRead(pipeReader *io.PipeReader) {
	var (
		err error
		n   int
	)
	data := make([]byte, 1024)
	for n, err = pipeReader.Read(data); err == nil; n, err = pipeReader.Read(data) {
		fmt.Printf("%s\n", data[:n])
	}
	fmt.Println("writer 端 closewitherror 后：", err)
}
