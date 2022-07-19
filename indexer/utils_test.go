package indexer

import "testing"

func TestItos(t *testing.T) {
	var n = []uint64{0, 1, 2, 3, 4, 5, 6, 321872, 10}

	for _, v := range n {
		b := Itos(v)
		t.Logf("string: %v", b)
	}

	for _, v := range n {
		b := Stoi(Itos(v))
		t.Logf("%d", b)
	}

}
