/**
 * @Author : gaopeihan
 * @ClassName : fileEs.go
 * @Date : 2021/5/20 08:34
 */
package esSearch

import (
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
)


//InsertMeta 插入元数据
func (e *ESClient)InsertMeta(meta interface{})(succeed bool){

	dataJSON, err := json.Marshal(meta)
	js := string(dataJSON)
	_, err = e.Client.Index().
		Index(e.index).
		BodyJson(js).
		Do(e.ctx)

	if err != nil {
		panic(err)
		return false
	}
	return true
}

//BulkInsert 打包上传chunks元数据
func (e *ESClient) BulkInsert(chunksMap map[string]int64)(succeed bool)  {
	bulkRequest := e.Client.Bulk()
	for hash,size :=range chunksMap{
		c := ChunkMeta{Hash: hash,Size: size}
		req := elastic.NewBulkIndexRequest().Index(e.index).Doc(c)
		bulkRequest = bulkRequest.Add(req)
	}
	_,err := bulkRequest.Do(e.ctx)
	if err != nil {
		panic(err)
		return false
	}
	return true
}

//SearchMeta 查询元数据是否存在
func (e *ESClient) SearchMeta(hash string,size int64)(hit bool,chunksHash []string) {

	searchSource := elastic.NewSearchSource()
	searchSource.Query(elastic.NewMatchQuery("hash", hash))
	searchService := e.Client.Search().Index(e.index).SearchSource(searchSource)

	searchResult, err := searchService.Do(e.ctx)
	if err != nil {
		fmt.Println("[ProductsES][GetPIds]Error=", err)
		return false,nil
	}
	if len(searchResult.Hits.Hits) == 0 {
		return false,nil
	} else {

		if e.index == "file"{
			for _, hit := range searchResult.Hits.Hits {

				var f FileMeta
				err := json.Unmarshal(hit.Source, &f)
				if err != nil {
					fmt.Println("[Getting FileMeta][Unmarshal] Err=", err)
					return false,nil
				}
				if f.Size == size {
					return true,f.Chunks
				}
			}
		}else {
			for _, hit := range searchResult.Hits.Hits {

				var f ChunkMeta
				err := json.Unmarshal(hit.Source, &f)
				if err != nil {
					fmt.Println("[Getting FileMeta][Unmarshal] Err=", err)
					return false,nil
				}
				if f.Size == size {
					return true,nil
				}
			}
		}
	}

	return false,nil
}