package chunk

import (
	"github.com/jotfs/fastcdc-go"
)

const kiB = 1024
const miB = 1024 * kiB

var defaultOpts = fastcdc.Options{
	AverageSize: 1 * miB,
}


type ChunksWorker struct {
	chunker *fastcdc.Chunker
}


func (c *ChunksWorker)Chunk()(fastcdc.Chunk,error ) {

	chunk, err := c.chunker.Next()
	return chunk,err

}


