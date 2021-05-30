package metafile

import (
	"DDSS/tools"
	"os"
)

func NewReadFile(path string) (r ReadFile,err error){
	ff, err := os.Open(path)
	if err != nil{
		tools.Error(err,"打开文件错误")
	}
	r = ReadFile{}
	r.F = ff
	err = r.GetFileStat()
	if err != nil{
		return ReadFile{}, err
	}
	return r,nil
}

