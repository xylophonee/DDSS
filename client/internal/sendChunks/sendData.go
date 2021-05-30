/**
 * @Author : gaopeihan
 * @ClassName : sendData.go
 * @Date : 2021/5/23 20:46
 */
package sendChunks

import (
	"DDSS/client/internal/tcp"
	"DDSS/tools"
	"bytes"
	"fmt"
	"github.com/jotfs/fastcdc-go"
	"net"
	"sync"
)

func PutDataOperate(wg *sync.WaitGroup, chunksChan chan fastcdc.Chunk, errorChan chan error){
	tcpClient := tcp.NewTCPClient()
	exit := false
	for {
		if !exit {
			select {
			case <-errorChan:
				exit = true
			case c := <-chunksChan:
				go goOperate(tcpClient,&c,errorChan)
			}
		} else {
			break
		}
	}

	wg.Done()
}


func goOperate(tcpClient *tcp.TcpClient,c *fastcdc.Chunk,errorChan chan error){
	cmd := tcp.Cmd{Name: "PUT", Hash: c.Hash, Size: int64(c.Length), Data: c.Data, Error: nil}
	tcpClient.Run(&cmd)
	if cmd.Error != nil{
		errorChan <- cmd.Error
	}
}

func SendData(wg *sync.WaitGroup, chunksChan chan fastcdc.Chunk, errorChan chan error) {

	conn, e := net.Dial("tcp", "127.0.0.1:8080")
	if e != nil {
		tools.Error(e, "TCP连接失败")
	}
	defer func() {
		fmt.Println("conn关闭")
		conn.Close()
	}()
	exit := false
	for {
		if !exit {
			select {
			case <-errorChan:
				exit = true
			case c := <-chunksChan:
				sendFirst := []byte(fmt.Sprintf("P%d", len(c.Data)))
				s := BytesCombine(sendFirst, c.Data)
				fmt.Println(string(sendFirst))
				n, e := conn.Write(s)
				if n != len(s) {
					tools.Error(e, "chunkData 读取错误")
					errorChan <- e
				}
			}
		} else {
			break
		}
	}

	wg.Done()

}

func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(" "))
}

// 获取服务端发送来的信息
func onMessageBack(conn net.Conn) {

	buffer := make([]byte, 1024)
	for {
		if conn != nil {
			n, err := conn.Read(buffer)
			if err != nil {
				tools.Error(err, "接收失败")
			}
			replyMsg := buffer[:n]
			fmt.Println("接收到的消息:", string(replyMsg))
		}
	}
}
