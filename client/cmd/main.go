package main

import (
	"DDSS/client/internal/esSearch"
	"DDSS/client/internal/operate"
)

//const (
//
//	ip0 ="127.0.0.1:8080"
//	ip1 ="127.0.0.1:8081"
//	ip2 ="127.0.0.1:8082"
//	ip3 ="127.0.0.1:8083"
//	ip4 ="127.0.0.1:8084"
//	ip5 ="127.0.0.1:8085"
//
//)

func main(){
	//启动服务
	esFile, _ := esSearch.NewESClient("file")
	esChunks, _ := esSearch.NewESClient("chunks")
	operate := operate.NewOperate(esFile,esChunks)
	operate.Put("/Users/xylophone/Downloads/Go语言实战.pdf")
}



//func PUT(path string)  {
//
//	//读取文件和获取文件的信息
//	f, err := metafile.NewReadFile(path)
//	if err != nil {
//		tools.PrintError(err)
//	}
//	defer f.Close()
//	//查询文件是否已经存在
//	t1 := time.Now()
//	esFile, _ := esSearch.NewESClient("file")
//	exist,chunksHash := esFile.SearchMeta(f.Stat.Hash,f.Stat.Size)
//	//不存在则进行分块去重
//	if !exist{
//		//小于4K不分块
//		if f.Stat.Size < 4000{
//			//todo 是否需要存储到chunk es？
//			chunksHash = []string{f.Stat.Hash}
//		}else {
//			chunksHash = sendChunks.GetChunksData(f.F)
//		}
//	}
//	//存储文件元数据
//	f.Stat.Chunks = chunksHash
//	fileMeta := f.Stat
//	esFile.InsertMeta(fileMeta)
//	t := time.Since(t1).Seconds()
//	fmt.Printf("%s 上传完成,耗时%f\n",f.Stat.Name,t)
//}


func Get(hash string,name string){

}

