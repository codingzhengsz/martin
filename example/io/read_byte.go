package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

func readByte() {
	buffer := bytes.NewBuffer([]byte{'a', 'b'})
	if c, err := buffer.ReadByte(); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("----", c)
	}
	//err := buffer.UnreadByte()
	if c, err := buffer.ReadByte(); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("++++", c)
	}
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
}

func limitedRead() {
	content := "This Is LimitReader Example"
	reader := strings.NewReader(content)
	limitReader := &io.LimitedReader{R: reader, N: 8}
	for limitReader.N > 0 {
		tmp := make([]byte, 3)
		if read, err := limitReader.Read(tmp); err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("read: ",read)
		}
		fmt.Printf("%s\n", tmp)
	}
}

//func main() {
//	limitedRead()
//}
