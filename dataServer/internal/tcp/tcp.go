/**
 * @Author : gaopeihan
 * @ClassName : tcp.go
 * @Date : 2021/5/24 20:43
 */
package tcp

import (
	"DDSS/dataServer/internal/blobsfile"
	"DDSS/tools"
	"fmt"
	"net"
)

type Server struct {
	b blobsfile.StoreWorker
}

func New(b blobsfile.StoreWorker)*Server{
	return &Server{b}
}

func (s *Server)Listen(){
	//建立tcp监听
	listener, e := net.Listen("tcp", ":8080")
	if e != nil{
		tools.Error(e, "net.Listen")
	}

	defer listener.Close()
	fmt.Println("server begin...")

	for {
		//接受客户端请求，建立会话专线Conn

		conn, e := listener.Accept()
		if e != nil{
			tools.Error(e, "listener.Accept")
		}
		fmt.Printf("接收到来自%s的请求\n", conn.RemoteAddr())
		go s.process(conn)
	}
}

