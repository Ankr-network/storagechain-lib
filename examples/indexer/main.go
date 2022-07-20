package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/Ankr-network/storagechain-lib/indexer"
	"github.com/sunvim/utils/tools"
)

func main() {

	// testFunc()
	testCompress()
	// testHashLen()

}

func testHashLen() {
	hasher := sha256.New()
	hasher.Write(tools.StringToBytes(strings.Repeat(strconv.Itoa(8), 4)))
	hs := hex.EncodeToString(hasher.Sum(nil))
	fmt.Printf("hash: %s \n len: %d \n", hs, len(hs))
	bs := make([]byte, 64)
	hex.Encode(bs, hasher.Sum(nil))
	fmt.Printf("hash: %s \n len: %d \n", bs, len(bs))
}

type TestItem struct {
	Prefix []byte
	Item   *indexer.Item
}

func testCompress() {
	var key = 888

	hasher := sha256.New()

	items := make([]*TestItem, 0)
	for i := 0; i < 2000; i++ {
		hasher.Reset()
		hasher.Write([]byte(strings.Repeat(strconv.Itoa(i), 4)))
		items = append(items, &TestItem{Prefix: hasher.Sum(nil), Item: &indexer.Item{Pos: uint64(i), Length: uint64(i)}})
	}

	hasher.Reset()
	hasher.Write([]byte(strings.Repeat(strconv.Itoa(key), 4)))
	target := hasher.Sum(nil)
	trie := indexer.NewTrie()
	for _, item := range items {
		trie.Insert(item.Prefix, item.Item)
	}

	item := trie.Get(target)
	fmt.Printf("STX key: %s item: %v\n", hex.EncodeToString(target), item)

	trie.Walk(nil, func(prefix []byte, item *indexer.Item) error {
		if int(item.Pos) == key {
			fmt.Printf("trie1 %s  %v\n", hex.EncodeToString(prefix), item)
		}
		return nil
	})
	fmt.Println(strings.Repeat("-", 80))

	ss, err := trie.Marshal()
	if err != nil {
		fmt.Printf("Marshal error: %v\n", err)
		return
	}
	fmt.Printf("trie size: %d serialize: %d \n", trie.Size(), len(ss))

	t2 := indexer.NewTrie()
	err = t2.Unmarshal(ss)
	if err != nil {
		fmt.Printf("unmarshal error: %v\n", err)
		return
	}
	item = t2.Get(target)
	fmt.Printf("key: %s item: %v\n", hex.EncodeToString(target), item)

}

func testFunc() {
	// Create a new tree.
	trie := indexer.NewTrie()

	// Insert some items.
	trie.Insert([]byte("0x111234"), &indexer.Item{Pos: 1, Length: 2})
	trie.Insert([]byte("0x111241"), &indexer.Item{Pos: 3, Length: 4})
	trie.Insert([]byte("0x222123"), &indexer.Item{Pos: 5, Length: 6})
	trie.Insert([]byte("0x2222223"), &indexer.Item{Pos: 7, Length: 8})

	// Just check if some things are present in the tree.
	key := []byte("0x111234")
	fmt.Printf("%q present? %v\n", key, trie.Match(key))
	fmt.Println(strings.Repeat("-", 80))
	key = []byte("0x111")
	fmt.Printf("Anybody called %q here? %v\n", key, trie.MatchSubtree(key))
	fmt.Println(strings.Repeat("-", 80))

	// Walk the tree.
	trie.Visit(printItem)
	fmt.Println(strings.Repeat("-", 80))

	// Walk a subtree.
	trie.VisitSubtree([]byte("0x111"), printItem)
	fmt.Println(strings.Repeat("-", 80))

	// Modify an item, then fetch it from the tree.
	trie.Set([]byte("0x111222333"), &indexer.Item{Pos: 9, Length: 10})
	key = []byte("0x111222333")
	fmt.Printf("%q: %v\n", key, trie.Get(key))

	fmt.Println("dump trie")
	ds := trie.Dump()
	fmt.Printf("dump: %v\n", ds)

	ss, err := indexer.Marshal(trie)
	if err != nil {
		fmt.Printf("Marshal error: %v\n", err)
		return
	}
	fmt.Printf("msgpack marshal: %v, size: %d \n", ss, len(ss))

	ss, err = json.Marshal(trie)
	if err != nil {
		fmt.Printf("Marshal error: %v\n", err)
		return
	}
	fmt.Printf("json marshal: %v, size: %d \n", ss, len(ss))
	// Walk prefixes.
	prefix := []byte("0x111222333")
	trie.VisitPrefixes(prefix, printItem)
	// "Karel Hynek Macha": 10

	// Delete some items.
	trie.Delete([]byte("0x111234"))
	trie.Delete([]byte("0x111241"))

	// Walk again.
	trie.Visit(printItem)

	// Delete a subtree.
	trie.DeleteSubtree([]byte("0x111"))

	// Print what is left.
	trie.Visit(printItem)

}

func printItem(prefix []byte, item *indexer.Item) error {
	fmt.Printf("%s  %v\n", prefix, item)
	return nil
}
