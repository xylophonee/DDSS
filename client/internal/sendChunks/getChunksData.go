/**
 * @Author : gaopeihan
 * @ClassName : getChunksData.go
 * @Date : 2021/5/23 21:10
 */
package sendChunks

import (
	"DDSS/client/internal/chunk"
	"DDSS/client/internal/esSearch"
	"DDSS/tools"
	"github.com/jotfs/fastcdc-go"
	"io"
	"os"
	"sync"
)

func GetChunksData(f *os.File)(chunksHash []string){
	esChunks, _ := esSearch.NewESClient("chunks")
	chunksMap := make(map[string]int64)
	chunksChan := make(chan fastcdc.Chunk)
	errorChan := make(chan error)
	defer close(chunksChan)
	defer close(errorChan)
	defer f.Close()
	var wg sync.WaitGroup
	wg.Add(1)
	go PutDataOperate(&wg,chunksChan,errorChan)
	c ,err:= chunk.NewChunksWorker(f)
	if err != nil {
		tools.PrintError(err)
	}
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
				go esChunks.BulkInsert(to)
			}
		}
	}
	//存储chunks的元数据
	esChunks.BulkInsert(chunksMap)
	wg.Wait()
	return chunksHash
}

func copyMap(from map[string]int64)map[string]int64{
	to := make(map[string]int64)
	if len(from) != 0{
		for h,s:=range from{
			to[h] = s
		}
	}
	return to
}

