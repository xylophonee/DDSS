package metafile

import (
	"DDSS/tools"
	"bytes"
	"io"
	"os"
)

type ReadFile struct {
	F *os.File
	Stat fileStat
}

type fileStat struct {
	Name    string  `json:"name"`
	Size    int64   `json:"size"`
	Hash    string `json:"hash"`
	Chunks  []string `json:"chunks"`
}


//GetFileStat 获取文件的属性 Hash,Name,Size
func (r *ReadFile)GetFileStat() error {
	info, err := r.F.Stat()
	if err != nil{
		return err
	}
	//读取文件
	var n int64 = bytes.MinRead
	if size := info.Size() + bytes.MinRead; size > n {
		n = size
	}
	data,err:=readAll(r.F,n)

	if err != nil {
		tools.PrintError(err)
	}
	_, err = r.F.Seek(0, 0)
	if err != nil {
		tools.Error(err,"移动文件指针出错")
	}
	//计算Hash
	fileHash := tools.HashByteToString(tools.CalHash(data))
	r.Stat = fileStat{Name: info.Name(),Size: info.Size(),Hash: fileHash,Chunks: nil}
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

func readAll(r io.Reader, capacity int64) (b []byte, err error) {
	var buf bytes.Buffer
	// If the buffer overflows, we will get bytes.ErrTooLarge.
	// Return that as an error. Any other panic remains.
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if panicErr, ok := e.(error); ok && panicErr == bytes.ErrTooLarge {
			err = panicErr
		} else {
			panic(e)
		}
	}()
	if int64(int(capacity)) == capacity {
		buf.Grow(int(capacity))
	}
	_, err = buf.ReadFrom(r)
	return buf.Bytes(), err
}