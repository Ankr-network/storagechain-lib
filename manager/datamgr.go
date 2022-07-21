package manager

import (
	"bytes"
	"sync"

	"github.com/Ankr-network/storagechain-lib/indexer"
	lru "github.com/hashicorp/golang-lru"
	"golang.org/x/exp/mmap"
)

var (
	bufferpool = &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
)

type DataNodeMgr struct {
	dataPath string
	cache    *lru.ARCCache
}

func NewDataNodeMgr(size int, path string) *DataNodeMgr {
	cache, err := lru.NewARC(size)
	if err != nil {
		panic(err)
	}
	return &DataNodeMgr{cache: cache, dataPath: path}
}

func (mgr *DataNodeMgr) Get(blockNum string) (*DataNode, error) {
	if v, ok := mgr.cache.Get(blockNum); ok {
		return v.(*DataNode), nil
	}

	bp := bufferpool.Get().(*bytes.Buffer)
	defer bufferpool.Put(bp)

	bp.Reset()
	dn := &DataNode{}
	bp.WriteString(mgr.dataPath)
	bp.WriteString("/b")
	bp.WriteString(blockNum)
	ra, err := mmap.Open(bp.String())
	if err != nil {
		return nil, err
	}
	dn.Reader = ra
	bp.WriteString(mgr.dataPath)
	bp.WriteString("/h")
	bp.WriteString(blockNum)
	trie := indexer.NewTrie()
	err = trie.ReadFromFile(bp.String())
	if err != nil {
		return nil, err
	}
	dn.Header = trie
	dn.Name = blockNum
	mgr.cache.Add(blockNum, dn)
	return dn, nil
}

func (mgr *DataNodeMgr) Set(blockNum string, node *DataNode) {
	mgr.cache.Add(blockNum, node)
}

func (mgr *DataNodeMgr) Remove(blockNum string) {
	mgr.cache.Remove(blockNum)
}

func (mgr *DataNodeMgr) Clear() {
	mgr.cache.Purge()
}

func (mgr *DataNodeMgr) Size() int {
	return mgr.cache.Len()
}

func (mgr *DataNodeMgr) Keys() []string {
	keys := make([]string, 0)
	for _, k := range mgr.cache.Keys() {
		keys = append(keys, k.(string))
	}
	return keys
}

func (mgr *DataNodeMgr) Values() []*DataNode {
	values := make([]*DataNode, 0)
	for _, v := range mgr.cache.Keys() {
		val, _ := mgr.cache.Get(v.(string))
		values = append(values, val.(*DataNode))
	}
	return values
}
