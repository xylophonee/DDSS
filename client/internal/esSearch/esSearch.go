/**
 * @Author : gaopeihan
 * @ClassName : esSearch.go
 * @Date : 2021/5/23 14:12
 */
package esSearch

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
)

type FileMeta struct {
	Name    string  `json:"name"`
	Size    int64   `json:"size"`
	Ext     string `json:"ext"`
	Hash    string `json:"hash"`
	Chunks  []string `json:"chunks"`
}

type ChunkMeta struct {
	Hash string `json:"hash"`
	Size int64 `json:"size"`
}

type ESClient struct {
	index string
	Client *elastic.Client
	ctx context.Context
}

func NewESClient(index string) (ESClient, error) {

	client, err :=  elastic.NewClient(elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false))

	fmt.Printf("%sES initialized...\n", index)

	return ESClient{index: index,Client: client,ctx: context.Background()}, err

}