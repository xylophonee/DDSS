package tools

import (
	"encoding/json"
)

func StructToJson(v struct{})(jsonByte []byte,err error){

	jsonByte,err = json.Marshal(v)
	return
}

