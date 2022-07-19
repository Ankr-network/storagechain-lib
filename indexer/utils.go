package indexer

import (
	"encoding/binary"
	"unsafe"
)

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
func froms(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func tos(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
