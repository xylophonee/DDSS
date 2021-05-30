package chunk

import (
	"github.com/jotfs/fastcdc-go"
	"os"
)

func NewChunksWorker(f *os.File)(ChunksWorker,error){

	y := ChunkConfig

	chunker, err := fastcdc.NewChunker(f, fastcdc.Options{
		AverageSize:          y.Chunk.AverageSize*kiB,
		MinSize:              y.Chunk.MinSize*kiB,
		MaxSize:              y.Chunk.MaxSize*kiB,
		Normalization:        y.Chunk.Normalization,
		DisableNormalization: y.Chunk.DisableNormalization,
	})


	return ChunksWorker{chunker: chunker},err
}