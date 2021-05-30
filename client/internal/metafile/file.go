package metafile

func NewReadFile(path string) (r ReadFile,err error){
	r = ReadFile{filePath: path}
	err = r.ReadFile()
	err = r.GetFileStat()
	if err != nil{
		return ReadFile{}, err
	}
	return r,nil
}


