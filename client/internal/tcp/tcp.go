/**
 * @Author : gaopeihan
 * @ClassName : tcp.go
 * @Date : 2021/5/25 16:53
 */
package tcp

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

type TcpClient struct {
	net.Conn
	r *bufio.Reader
}

type Cmd struct {
	Name  string
	Size   int64
	Data []byte
	Hash string
	Error error
}

func NewTCPClient() *TcpClient {
	c, e := net.Dial("tcp", "127.0.0.1:8080")
	if e != nil {
		panic(e)
	}
	r := bufio.NewReader(c)
	return &TcpClient{c, r}
}

func (c *TcpClient) sendGet(hash string) {
	klen := len(hash)
	c.Write([]byte(fmt.Sprintf("G%d %s", klen, hash)))
}

func (c *TcpClient) sendPut(w []byte)(int,error) {

	n, err := c.Write(w)
	return n,err
}

func (c *TcpClient) sendDel(key string) {
	klen := len(key)
	c.Write([]byte(fmt.Sprintf("D%d %s", klen, key)))
}

func (c *TcpClient) Run(cmd *Cmd) {
	if cmd.Name == "get" {
		c.sendGet(cmd.Hash)
		cmd.Data, cmd.Error = c.recvResponse()
		return
	}
	if cmd.Name == "PUT" {
		sendFirst := []byte(fmt.Sprintf("P%d", cmd.Size))
		s := BytesCombine(sendFirst, cmd.Data)
		fmt.Println(string(sendFirst))
		n, err := c.sendPut(s)
		if n != int(cmd.Size) || err != nil{
			cmd.Error = err
			return
		}
		//_, cmd.Error = c.recvResponse()
		return
	}
	if cmd.Name == "del" {
		c.sendDel(cmd.Hash)
		_, cmd.Error = c.recvResponse()
		return
	}
	panic("unknown cmd name " + cmd.Name)
}

func (c *TcpClient) recvResponse() ([]byte, error) {
	vlen := readLen(c.r)
	if vlen == 0 {
		return nil, nil
	}
	if vlen < 0 {
		err := make([]byte, -vlen)
		_, e := io.ReadFull(c.r, err)
		if e != nil {
			return nil, e
		}
		return nil, errors.New(string(err))
	}
	value := make([]byte, vlen)
	_, e := io.ReadFull(c.r, value)
	if e != nil {
		return nil, e
	}
	return value, nil
}

func readLen(r *bufio.Reader) int {
	tmp, e := r.ReadString(' ')
	if e != nil {
		log.Println(e)
		return 0
	}
	l, e := strconv.Atoi(strings.TrimSpace(tmp))
	if e != nil {
		log.Println(tmp, e)
		return 0
	}
	return l
}

func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(" "))
}