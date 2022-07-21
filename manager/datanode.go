package manager

import (
	"crypto/sha256"
	"hash"
	"sync"

	"github.com/Ankr-network/storagechain-lib/indexer"
	"golang.org/x/exp/mmap"
)

var (
	hasherPool = &sync.Pool{
		New: func() interface{} {
			return sha256.New()
		},
	}
)

type DataNode struct {
	Name   string
	Reader *mmap.ReaderAt
	Header *indexer.Trie
}

func (dn *DataNode) Get(key string) ([]byte, error) {
	hasher := hasherPool.Get().(hash.Hash)
	defer hasherPool.Put(hasher)
	hasher.Reset()
	hasher.Write(indexer.Froms(key))
	item := dn.Header.Get(hasher.Sum(nil))
	rs := make([]byte, item.Length)
	dn.Reader.ReadAt(rs, int64(item.Pos))
	return rs, nil
}
