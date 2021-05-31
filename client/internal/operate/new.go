/**
 * @Author : gaopeihan
 * @ClassName : new.go
 * @Date : 2021/5/30 20:11
 */
package operate

import "DDSS/client/internal/esSearch"

type Operate struct {
	esFile esSearch.ESClient
	esChunks esSearch.ESClient
}

func NewOperate(esFile,esChunks esSearch.ESClient)Operate{
	return Operate{esFile: esFile,esChunks: esChunks}
}
