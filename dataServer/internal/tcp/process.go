/**
 * @Author : gaopeihan
 * @ClassName : process.go
 * @Date : 2021/5/24 20:48
 */
package tcp

import (
	"DDSS/tools"
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

func (s *Server)get(conn net.Conn, r *bufio.Reader)error{
	return nil
}

func (s *Server)put(conn net.Conn, r *bufio.Reader)error{
	data,err := s.readData(r)
	if err != nil{
		return err
	}
	fmt.Printf("接收到%d.\n", len(data))
	hash := tools.HashByteToString(tools.CalHash(data))
	err = s.b.PutChunk(hash, data)
	return sendResponse(len(data),err,conn)
}

func (s *Server)del(conn net.Conn, r *bufio.Reader)error{

	return nil
}

func (s *Server) process(conn net.Conn)  {
	if conn == nil{
		log.Panic("invalid TCP connection")
	}
	defer conn.Close()
	//buffer := make([]byte, 1024)
	//n, _ := conn.Read(buffer)
	//fmt.Println(buffer[:n])
	r := bufio.NewReader(conn)

	for {
		op, e := r.ReadByte()
		if e != nil {
			if e != io.EOF {
				log.Println("close connection due to error:", e)
			}
			return
		}
		if op == 'P' {
			e = s.put(conn, r)
		} else if op == 'G' {
			e = s.get(conn, r)
		} else if op == 'D' {
			e = s.del(conn, r)
		} else{
			log.Println("close connection due to invalid operation:", op)
			return
		}
		if e != nil {
			log.Println("close connection due to error:", e)
			return
		}
	}
}