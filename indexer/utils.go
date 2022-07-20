package indexer

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"

	"github.com/valyala/gozstd"
)

func Marshal(trie *Trie) ([]byte, error) {
	var buffer bytes.Buffer
	buffer.Reset()
	trie.Walk(nil, func(prefix []byte, item *Item) error {
		buffer.Write(prefix)
		buffer.Write(Itos(item.Pos))
		buffer.Write(Itos(item.Length))
		return nil
	})
	return gozstd.Compress(nil, buffer.Bytes()), nil
}

func Unmarshal(data []byte) (*Trie, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data is empty")
	}
	ds, err := gozstd.Decompress(nil, data)
	if err != nil {
		return nil, err
	}
	trie := NewTrie()
	scanner := bufio.NewScanner(bytes.NewReader(ds))
	scanner.Split(ScanData)
	var line, key []byte
	for scanner.Scan() {
		line = scanner.Bytes()
		key = make([]byte, 32)
		copy(key, line[:32])
		trie.Insert(key, &Item{Pos: Stoi(line[32:40]), Length: Stoi(line[40:48])})
	}
	return trie, nil
}

func Itos(n uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b[:], n)
	return b
}

func Stoi(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func ScanData(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := len(data); i >= 48 {
		// We have a full newline-terminated line.
		return 48, data[0:48], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}
func Froms(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Tos(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
