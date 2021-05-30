/**
 * @Author : gaopeihan
 * @ClassName : main.go
 * @Date : 2021/5/25 10:07
 */
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"time"
)

func main(){

	t1 := time.Now()
	time.Sleep(1 * time.Second)
	t := time.Since(t1)
	fmt.Println(t.Seconds())

}

func RetrieveROM(filename string) ([]byte, error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats, statsErr := file.Stat()
	if statsErr != nil {
		return nil, statsErr
	}

	var size int64 = stats.Size()
	bytes := make([]byte, size)

	bufr := bufio.NewReader(file)
	_,err = bufr.Read(bytes)

	return bytes, err
}

func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(" "))
}

func BytesCombine1(pBytes ...[]byte) []byte {
	length := len(pBytes)
	s := make([][]byte, length)
	for index := 0; index < length; index++ {
		s[index] = pBytes[index]
	}
	sep := []byte("")
	return bytes.Join(s, sep)
}

func BytesCombine3(pBytes ...[]byte) []byte {
	var buffer bytes.Buffer
	for index := 0; index < len(pBytes); index++ {
		buffer.Write(pBytes[index])
	}
	return buffer.Bytes()
}
