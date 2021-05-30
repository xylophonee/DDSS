package tools

import (
	"fmt"
	"golang.org/x/crypto/blake2b"
)

func CalHash(data []byte) (hash [32]byte) {

	hash = blake2b.Sum256(data)
	return hash
}

func HashByteToString(h [32]byte) string {
	return fmt.Sprintf("%x", h)
}
