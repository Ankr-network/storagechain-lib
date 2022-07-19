package main

import (
	"crypto/sha256"
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

}

func testCompress() {
	// Create a new tree.
	trie := indexer.NewTrie()

	hasher := sha256.New()
	for i := 0; i < 100; i++ {
		hasher.Reset()
		hasher.Write(tools.StringToBytes(strings.Repeat(strconv.Itoa(i), 16)))
		trie.Insert(indexer.Prefix(hasher.Sum(nil)),
			&indexer.Item{Pos: uint64(i), Length: uint64(i)})
	}

	hasher.Reset()
	hasher.Write(tools.StringToBytes(strings.Repeat(strconv.Itoa(12), 16)))
	target := indexer.Prefix(hasher.Sum(nil))
	item := trie.Get(target)
	fmt.Printf("item: %v\n", item)

	fmt.Printf("trie size: %d\n", trie.Size())

	// trie.Walk(nil, func(prefix indexer.Prefix, item *indexer.Item) error {
	// 	fmt.Printf("%s  %v\n", prefix, item)
	// 	return nil
	// })
	// fmt.Println(strings.Repeat("-", 80))

	ss, err := trie.Marshal()
	if err != nil {
		fmt.Printf("Marshal error: %v\n", err)
		return
	}
	fmt.Printf("origin marshal size: %d \n", len(ss))

	t2 := indexer.NewTrie()
	err = t2.Unmarshal(ss)
	if err != nil {
		fmt.Printf("unmarshal error: %v\n", err)
		return
	}
	fmt.Printf("t2: %v\n", t2)
	item = t2.Get(target)
	fmt.Printf("item: %v\n", item)

}

func testFunc() {
	// Create a new tree.
	trie := indexer.NewTrie()

	// Insert some items.
	trie.Insert(indexer.Prefix("0x111234"), &indexer.Item{Pos: 1, Length: 2})
	trie.Insert(indexer.Prefix("0x111241"), &indexer.Item{Pos: 3, Length: 4})
	trie.Insert(indexer.Prefix("0x222123"), &indexer.Item{Pos: 5, Length: 6})
	trie.Insert(indexer.Prefix("0x2222223"), &indexer.Item{Pos: 7, Length: 8})

	// Just check if some things are present in the tree.
	key := indexer.Prefix("0x111234")
	fmt.Printf("%q present? %v\n", key, trie.Match(key))
	fmt.Println(strings.Repeat("-", 80))
	key = indexer.Prefix("0x111")
	fmt.Printf("Anybody called %q here? %v\n", key, trie.MatchSubtree(key))
	fmt.Println(strings.Repeat("-", 80))

	// Walk the tree.
	trie.Visit(printItem)
	fmt.Println(strings.Repeat("-", 80))

	// Walk a subtree.
	trie.VisitSubtree(indexer.Prefix("0x111"), printItem)
	fmt.Println(strings.Repeat("-", 80))

	// Modify an item, then fetch it from the tree.
	trie.Set(indexer.Prefix("0x111222333"), &indexer.Item{Pos: 9, Length: 10})
	key = indexer.Prefix("0x111222333")
	fmt.Printf("%q: %v\n", key, trie.Get(key))

	fmt.Println("dump trie")
	ds := trie.Dump()
	fmt.Printf("dump: %v\n", ds)

	ss, err := trie.Marshal()
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
	prefix := indexer.Prefix("0x111222333")
	trie.VisitPrefixes(prefix, printItem)
	// "Karel Hynek Macha": 10

	// Delete some items.
	trie.Delete(indexer.Prefix("0x111234"))
	trie.Delete(indexer.Prefix("0x111241"))

	// Walk again.
	trie.Visit(printItem)

	// Delete a subtree.
	trie.DeleteSubtree(indexer.Prefix("0x111"))

	// Print what is left.
	trie.Visit(printItem)

}

func printItem(prefix indexer.Prefix, item *indexer.Item) error {
	fmt.Printf("%s  %v\n", prefix, item)
	return nil
}
