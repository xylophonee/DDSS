/**
 * @Author : gaopeihan
 * @ClassName : new.go
 * @Date : 2021/5/30 20:10
 */
package tcp

import (
	"bufio"
	"net"
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
