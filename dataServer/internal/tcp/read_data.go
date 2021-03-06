/**
 * @Author : gaopeihan
 * @ClassName : read_data.go
 * @Date : 2021/5/25 10:18
 */
package tcp

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

func (s *Server) readData(r *bufio.Reader) ([]byte, error) {
	Datalen, e := readLen(r)
	if e != nil {
		return []byte(""), e
	}
	d := make([]byte, Datalen)
	_, e = io.ReadFull(r, d)
	if e != nil {
		return []byte(""), e
	}
	return d, nil
}

func readLen(r *bufio.Reader) (int, error) {
	tmp, e := r.ReadString(' ')
	if e != nil {
		return 0, e
	}
	l, e := strconv.Atoi(strings.TrimSpace(tmp))
	if e != nil {
		return 0, e
	}
	return l, nil
}

func sendResponse(lens int, err error, conn net.Conn) error {
	if err != nil {
		errString := err.Error()
		tmp := fmt.Sprintf("-%d ", len(errString)) + errString
		_, e := conn.Write([]byte(tmp))
		return e
	}
	vlen := len(strconv.Itoa(lens))
	message := fmt.Sprintf("%d ", vlen) + string(rune(lens))
	_, e := conn.Write([]byte(message))
	return e
}