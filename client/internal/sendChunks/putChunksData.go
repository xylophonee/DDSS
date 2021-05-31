/**
 * @Author : gaopeihan
 * @ClassName : putChunksData.go
 * @Date : 2021/5/23 21:10
 */
package sendChunks

import (
	"DDSS/client/internal/chunk"
	"DDSS/client/internal/esSearch"
	"DDSS/client/internal/tcp"
	"DDSS/tools"
	"bytes"
	"fmt"
	"github.com/jotfs/fastcdc-go"
	"io"
	"net"
	"os"
	"sync"
)
// PutChunksData 获取块数据并上传至数据服务器
func PutChunksData(f *os.File,esChunks esSearch.ESClient)(chunksHash []string){

	c ,err:= chunk.NewChunksWorker(f)
	if err != nil {
		tools.Error(err,"初始化分块出错")
	}
	chunksMap := make(map[string]int64)
	chunksChan := make(chan fastcdc.Chunk)
	errorChan := make(chan error)
	defer close(chunksChan)
	defer close(errorChan)
	defer f.Close()
	var wg sync.WaitGroup
	wg.Add(1)
	go putDataOperate(&wg,chunksChan,errorChan)
	for {
		c,err := c.Chunk()
		if err == io.EOF {
			errorChan <- err
			break
		}
		exist,_ := esChunks.SearchMeta(c.Hash, int64(c.Length))
		chunksHash = append(chunksHash,c.Hash)

		if !exist{
			_,ok := chunksMap[c.Hash]
			if !ok{
				chunksMap[c.Hash] = int64(c.Length)
				chunksChan <- c
			}
			if len(chunksMap) == 100{
				to := copyMap(chunksMap)
				chunksMap = make(map[string]int64)
				//存储chunks的元数据
				go esChunks.BulkInsert(to)
			}
		}
	}
	//存储chunks的元数据
	if len(chunksMap) != 0{
		esChunks.BulkInsert(chunksMap)
	}
	wg.Wait()
	return chunksHash
}


func putDataOperate(wg *sync.WaitGroup, chunksChan chan fastcdc.Chunk, errorChan chan error){
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

//goOperate 负责发送数据
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

//BytesCombine []byte合并
func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(" "))
}
//copyMap 将一个Map复制到另一个中
func copyMap(from map[string]int64)map[string]int64{
	to := make(map[string]int64)
	if len(from) != 0{
		for h,s:=range from{
			to[h] = s
		}
	}
	return to
}

