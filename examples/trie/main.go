package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb/leveldb"
	"github.com/ethereum/go-ethereum/trie"
)

func main() {
	// testdb()
	openOldData()
}

func openOldData() {
	ldb, err := leveldb.New("trie.db", 256, 16, "", false)
	if err != nil {
		fmt.Println(err)
	}
	defer ldb.Close()
	db := trie.NewDatabase(ldb)
	t1, err := trie.New(common.HexToHash("09c889feaafd53779755259beaa0ff41c32512c8cac45152af46fae7ebdef210"), db)
	if err != nil {
		fmt.Println(err)
		return
	}

	vals := []struct{ k, v string }{
		{"do", "verb"},
		{"ether", "wookiedoo"},
		{"horse", "stallion"},
		{"shaman", "horse"},
		{"doge", "coin"},
		{"dog", "puppy"},
		{"somethingveryoddindeedthis is", "myothernodedata"},
	}
	for _, val := range vals {
		fmt.Printf("key: %s, val: %s, got: %s\n", val.k, val.v, getString(t1, val.k))
	}

}

func testdb() {

	ldb, err := leveldb.New("trie.db", 256, 16, "", false)
	if err != nil {
		fmt.Println(err)
	}
	defer ldb.Close()
	db := trie.NewDatabase(ldb)
	t1, _ := trie.New(common.Hash{}, db)

	vals := []struct{ k, v string }{
		{"do", "verb"},
		{"ether", "wookiedoo"},
		{"horse", "stallion"},
		{"shaman", "horse"},
		{"doge", "coin"},
		{"dog", "puppy"},
		{"somethingveryoddindeedthis is", "myothernodedata"},
	}
	for _, val := range vals {
		updateString(t1, val.k, val.v)
	}

	exp, _, err := t1.Commit(nil)
	if err != nil {
		fmt.Printf("commit error: %v", err)
		return
	}
	fmt.Printf("root: %x\n", exp)

	db.Commit(exp, false, nil)

	trie2, err := trie.New(common.Hash{}, db)
	if err != nil {
		fmt.Printf("can't recreate trie at %x: %v", exp, err)
		return
	}
	for _, kv := range vals {
		fmt.Printf("trie2 have %s => %s, got: %s\n", kv.k, kv.v, getString(trie2, kv.k))
	}

}

func getString(trie *trie.Trie, k string) []byte {
	return trie.Get([]byte(k))
}

func updateString(trie *trie.Trie, k, v string) {
	trie.Update([]byte(k), []byte(v))
}
