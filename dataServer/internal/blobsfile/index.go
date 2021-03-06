package blobsfile

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"a4.io/blobstash/pkg/rangedb"
)

// FIXME(tsileo): optimize the index with the benchmark (not worth it if inserting the blob take longer)

// MetaKey and BlobPosKey are used to namespace the DB keys.
const (
	metaKey byte = iota
	blobPosKey
)

// formatKey prepends the prefix byte to the given key.
func formatKey(prefix byte, bkey []byte) []byte {
	res := make([]byte, len(bkey)+1)
	res[0] = prefix
	copy(res[1:], bkey)
	return res
}

// fileIndex holds the position of blobs in BlobsFile.
type blobsIndex struct {
	db   *rangedb.RangeDB
	path string
}

// fileMetaData is a blob entry in the index.
type blobPos struct {
	// bobs-n files
	n int
	// blobs offset/size in the blobs file
	offset   int64
	size     int
	blobSize int // the actual blob size (will be different from size if compression is enabled)
}

// Size returns the blob size (as stored in the BlobsFile).
func (blob *blobPos) Size() int {
	return blob.size
}

// Value serialize a BlobsPos as string.
// (value is encoded as uvarint: n + offset + size + blob size)
func (blob *blobPos) Value() []byte {
	bufTmp := make([]byte, 10)
	var buf bytes.Buffer
	w := binary.PutUvarint(bufTmp[:], uint64(blob.n))
	buf.Write(bufTmp[:w])
	w = binary.PutUvarint(bufTmp[:], uint64(blob.offset))
	buf.Write(bufTmp[:w])
	w = binary.PutUvarint(bufTmp[:], uint64(blob.size))
	buf.Write(bufTmp[:w])
	w = binary.PutUvarint(bufTmp[:], uint64(blob.blobSize))
	buf.Write(bufTmp[:w])
	return buf.Bytes()
}

func decodeBlobPos(data []byte) (blob *blobPos, error error) {
	blob = &blobPos{}
	r := bytes.NewBuffer(data)
	// read blob.n
	ures, err := binary.ReadUvarint(r)
	if err != nil {
		return blob, err
	}
	blob.n = int(ures)

	// read blob.offset
	ures, err = binary.ReadUvarint(r)
	if err != nil {
		return blob, err
	}
	blob.offset = int64(ures)

	// read blob.size
	ures, err = binary.ReadUvarint(r)
	if err != nil {
		return blob, err
	}
	blob.size = int(ures)

	// read blob.blobSize
	ures, err = binary.ReadUvarint(r)
	if err != nil {
		return blob, err
	}
	blob.blobSize = int(ures)

	return blob, nil
}

// newIndex initializes a new index.
func newIndex(path string) (*blobsIndex, error) {
	dbPath := filepath.Join(path, "file-index") //????????????
	//log.Println(dbPath)
	db, err := rangedb.New(dbPath) //rangedb:An event store database in Go.
	return &blobsIndex{db: db, path: dbPath}, err
}

func (index *blobsIndex) formatBlobPosKey(key string) []byte {
	return formatKey(blobPosKey, []byte(key))
}

// Close closes all the open file descriptors.
func (index *blobsIndex) Close() error {
	return index.db.Close()
}

// remove removes the kv file.
func (index *blobsIndex) remove() error {
	return os.RemoveAll(index.path)
}

// setPos creates a new fileMetaData entry in the index for the given hash.
func (index *blobsIndex) setPos(hexHash string, pos *blobPos) error {
	hash, err := hex.DecodeString(hexHash)
	if err != nil {
		return err
	}
	//kv?????????k??? filePosKey+hash v??? n + offset + size + blob size

	return index.db.Set(formatKey(blobPosKey, hash), pos.Value())
}

// deletePos deletes the stored fileMetaData for the given hash.
// func (index *fileIndex) deletePos(hexHash string) error {
//	hash, err := hex.DecodeString(hexHash)
//	if err != nil {
//		return err
//	}
//	return index.db.Delete(formatKey(filePosKey, hash))
//}

// checkPos checks if a fileMetaData exists for the given hash (without decoding it).
func (index *blobsIndex) checkPos(hexHash string) (bool, error) {
	hash, err := hex.DecodeString(hexHash)
	if err != nil {
		return false, err
	}
	data, err := index.db.Get(formatKey(blobPosKey, hash))
	if err != nil {
		return false, fmt.Errorf("error getting BlobPos: %v", err)
	}
	if data == nil || len(data) == 0 {
		return false, nil
	}
	return true, nil
}

// getPos retrieve the stored fileMetaData for the given hash.
func (index *blobsIndex) getPos(hexHash string) (*blobPos, error) {
	hash, err := hex.DecodeString(hexHash)
	if err != nil {
		return nil, err
	}
	data, err := index.db.Get(formatKey(blobPosKey, hash))
	if err != nil {
		return nil, fmt.Errorf("error getting BlobPos: %v", err)
	}
	if data == nil {
		return nil, nil
	}
	//data ??????
	bpos, err := decodeBlobPos(data)
	return bpos, err
}

// setN stores the latest N (blobs-N) to remember the latest BlobsFile opened.
func (index *blobsIndex) setN(n int) error {
	return index.db.Set(formatKey(metaKey, []byte("n")), []byte(strconv.Itoa(n)))
}

// getN retrieves the latest N (blobs-N) stored.
func (index *blobsIndex) getN() (int, error) {
	data, err := index.db.Get(formatKey(metaKey, []byte("n")))
	if err != nil || string(data) == "" {
		return 0, nil
	}
	return strconv.Atoi(string(data))
}
