package metafile

import (
	"DDSS/tools"
	"io/ioutil"
	"os"
	"path"
)

type ReadFile struct {
	filePath string
	F *os.File
	Stat fileStat
}

type fileStat struct {
	Name    string  `json:"name"`
	Size    int64   `json:"size"`
	Ext     string `json:"ext"`
	Hash    string `json:"hash"`
	Chunks  []string `json:"chunks"`
}

//ReadFile 打开文件
func (r *ReadFile) ReadFile() error {

	ff, err := os.Open(r.filePath)
	if err != nil {
		return err
	}
	r.F = ff
	return nil
}
//GetFileStat 获取文件的属性
func (r *ReadFile)GetFileStat() error {
	info, err := r.F.Stat()
	if err != nil{
		return err
	}
	fileExt := path.Ext(r.filePath)
	if fileExt == ""{
		fileExt = "none"
	}
	data,err:=ioutil.ReadAll(r.F)

	if err != nil {
		tools.PrintError(err)
	}
	_, _ = r.F.Seek(0, 0)
	fileHash := tools.HashByteToString(tools.CalHash(data))
	r.Stat = fileStat{Name: info.Name(),Ext: fileExt,Size: info.Size(),Hash: fileHash,Chunks: nil}
	return nil
}

//Close 关闭文件
func (r *ReadFile) Close() error {
	err := r.F.Close()
	if  err != nil {
		return err
	}
	return nil
}