package blobsfile

import (
	"DDSS/tools"
)

type StoreWorker struct {
	Back *BlobsFiles
}
func NewStoreWorker() StoreWorker {
	back, err := New(&Opts{Directory: "/Users/xylophone/code/go/DDSS/dataServer/blobsFile", Compression: Snappy})
	if err != nil{
		tools.Error(err,"Blob error")
	}
	//defer back.Close()
	//defer os.RemoveAll("./tmp_blobsfile_test")
	return StoreWorker{Back: back}
}

func (s *StoreWorker)Close(){
	_ = s.Back.Close()
}

func (s *StoreWorker)PutChunk(hash string,data []byte)error{
	err := s.Back.Put(hash, data)
	return err
}

func (s *StoreWorker)GetChunk(hash string)([]byte,error){
	data,err := s.Back.Get(hash)
	return data,err
}
//func (s *StoreWorker)PutAllChunks(c *chunk.ChunksWorker)(chunksHash []string) {
//	//i := 0
//
//	chunksHash = []string{}
//	for {
//		f,err := c.Chunk()
//		if err == io.EOF {
//			break
//		}
//		if err != nil {
//			tools.PrintError(err)
//			os.Exit(1)
//		}
//		//fmt.Printf("%9d  %9d  %s\n", f.Offset, f.Length,fmt.Sprintf("%x", sha1.Sum(f.Data)))
//		//t := time.Now()
//		err = s.Back.Put(f.Hash, f.Data)
//		check(err)
//
//		chunksHash = append(chunksHash, f.Hash)
//		//fmt.Println(time.Since(t))
//		//i+=len(f.Data)
//		//fmt.Printf("\rStore %d", i)
//	}
//	//defer os.RemoveAll("./tmp_blobsfile_test")
//	return
//}
//
//func check(e error) {
//	if e != nil {
//		tools.PrintError(e)
//	}
//}