/**
 * @Author : gaopeihan
 * @ClassName : put.go
 * @Date : 2021/5/30 20:08
 */
package operate

import (
	"DDSS/client/internal/metafile"
	"DDSS/client/internal/sendChunks"
	"DDSS/tools"
	"fmt"
	"time"
)

func (o *Operate)Put(path string)  {

	//读取文件和获取文件的信息
	f, err := metafile.NewReadFile(path)
	if err != nil {
		tools.PrintError(err)
	}
	defer f.Close()
	t1 := time.Now()
	//查询文件是否已经存在
	exist,chunksHash := o.esFile.SearchMeta(f.Stat.Hash,f.Stat.Size)
	//不存在则进行分块去重
	if !exist{

		//小于4K不分块
		if f.Stat.Size < 4000{
			//todo 是否需要存储到chunk es？
			chunksHash = []string{f.Stat.Hash}
		}else {
			chunksHash = sendChunks.PutChunksData(f.F,o.esChunks)
		}
	}
	//存储文件元数据
	f.Stat.Chunks = chunksHash
	fileMeta := f.Stat
	o.esFile.InsertMeta(fileMeta)
	t := time.Since(t1).Seconds()
	fmt.Printf("%s 上传完成,耗时%f\n",f.Stat.Name,t)
}
