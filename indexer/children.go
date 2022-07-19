package indexer

import (
	"fmt"
	"io"
	"sort"
)

type ChildList interface {
	length() int
	head() *Trie
	add(child *Trie) ChildList
	remove(b byte)
	replace(b byte, child *Trie)
	next(b byte) *Trie
	walk(prefix *Prefix, visitor VisitorFunc) error
	print(w io.Writer, indent int)
	clone() ChildList
	total() int
}

type Tries []*Trie

func (t Tries) Len() int {
	return len(t)
}

func (t Tries) Less(i, j int) bool {
	strings := sort.StringSlice{string(t[i].Prefix), string(t[j].Prefix)}
	return strings.Less(0, 1)
}

func (t Tries) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

type SparseChildList struct {
	Children Tries
}

func newSparseChildList(maxChildrenPerSparseNode int) ChildList {
	return &SparseChildList{
		Children: make(Tries, 0, maxChildrenPerSparseNode),
	}
}

func (list *SparseChildList) length() int {
	return len(list.Children)
}

func (list *SparseChildList) head() *Trie {
	return list.Children[0]
}

func (list *SparseChildList) add(child *Trie) ChildList {
	// Search for an empty spot and insert the child if possible.
	if len(list.Children) != cap(list.Children) {
		list.Children = append(list.Children, child)
		return list
	}

	// Otherwise we have to transform to the dense list type.
	return newDenseChildList(list, child)
}

func (list *SparseChildList) remove(b byte) {
	for i, node := range list.Children {
		if node.Prefix[0] == b {
			list.Children[i] = list.Children[len(list.Children)-1]
			list.Children[len(list.Children)-1] = nil
			list.Children = list.Children[:len(list.Children)-1]
			return
		}
	}

	// This is not supposed to be reached.
	panic("removing non-existent child")
}

func (list *SparseChildList) replace(b byte, child *Trie) {
	// Make a consistency check.
	if p0 := child.Prefix[0]; p0 != b {
		panic(fmt.Errorf("child prefix mismatch: %v != %v", p0, b))
	}

	// Seek the child and replace it.
	for i, node := range list.Children {
		if node.Prefix[0] == b {
			list.Children[i] = child
			return
		}
	}
}

func (list *SparseChildList) next(b byte) *Trie {
	for _, child := range list.Children {
		if child.Prefix[0] == b {
			return child
		}
	}
	return nil
}

func (list *SparseChildList) walk(prefix *Prefix, visitor VisitorFunc) error {

	sort.Sort(list.Children)

	for _, child := range list.Children {
		*prefix = append(*prefix, child.Prefix...)
		if child.Item != nil {
			err := visitor(*prefix, child.Item)
			if err != nil {
				if err == SkipSubtree {
					*prefix = (*prefix)[:len(*prefix)-len(child.Prefix)]
					continue
				}
				*prefix = (*prefix)[:len(*prefix)-len(child.Prefix)]
				return err
			}
		}

		err := child.Children.walk(prefix, visitor)
		*prefix = (*prefix)[:len(*prefix)-len(child.Prefix)]
		if err != nil {
			return err
		}
	}

	return nil
}

func (list *SparseChildList) total() int {
	tot := 0
	for _, child := range list.Children {
		if child != nil {
			tot = tot + child.total()
		}
	}
	return tot
}

func (list *SparseChildList) clone() ChildList {
	clones := make(Tries, len(list.Children), cap(list.Children))
	for i, child := range list.Children {
		clones[i] = child.Clone()
	}

	return &SparseChildList{
		Children: clones,
	}
}

func (list *SparseChildList) print(w io.Writer, indent int) {
	for _, child := range list.Children {
		if child != nil {
			child.print(w, indent)
		}
	}
}

type DenseChildList struct {
	Min         int
	Max         int
	NumChildren int
	HeadIndex   int
	Children    []*Trie
}

func newDenseChildList(list *SparseChildList, child *Trie) ChildList {
	var (
		min int = 255
		max int = 0
	)
	for _, child := range list.Children {
		b := int(child.Prefix[0])
		if b < min {
			min = b
		}
		if b > max {
			max = b
		}
	}

	b := int(child.Prefix[0])
	if b < min {
		min = b
	}
	if b > max {
		max = b
	}

	children := make([]*Trie, max-min+1)
	for _, child := range list.Children {
		children[int(child.Prefix[0])-min] = child
	}
	children[int(child.Prefix[0])-min] = child

	return &DenseChildList{
		Min:         min,
		Max:         max,
		NumChildren: list.length() + 1,
		HeadIndex:   0,
		Children:    children,
	}
}

func (list *DenseChildList) length() int {
	return list.NumChildren
}

func (list *DenseChildList) head() *Trie {
	return list.Children[list.HeadIndex]
}

func (list *DenseChildList) add(child *Trie) ChildList {
	b := int(child.Prefix[0])
	var i int

	switch {
	case list.Min <= b && b <= list.Max:
		if list.Children[b-list.Min] != nil {
			panic("dense child list collision detected")
		}
		i = b - list.Min
		list.Children[i] = child

	case b < list.Min:
		children := make([]*Trie, list.Max-b+1)
		i = 0
		children[i] = child
		copy(children[list.Min-b:], list.Children)
		list.Children = children
		list.Min = b

	default: // b > list.max
		children := make([]*Trie, b-list.Min+1)
		i = b - list.Min
		children[i] = child
		copy(children, list.Children)
		list.Children = children
		list.Max = b
	}

	list.NumChildren++
	if i < list.HeadIndex {
		list.HeadIndex = i
	}
	return list
}

func (list *DenseChildList) remove(b byte) {
	i := int(b) - list.Min
	if list.Children[i] == nil {
		// This is not supposed to be reached.
		panic("removing non-existent child")
	}
	list.NumChildren--
	list.Children[i] = nil

	// Update head index.
	if i == list.HeadIndex {
		for ; i < len(list.Children); i++ {
			if list.Children[i] != nil {
				list.HeadIndex = i
				return
			}
		}
	}
}

func (list *DenseChildList) replace(b byte, child *Trie) {
	// Make a consistency check.
	if p0 := child.Prefix[0]; p0 != b {
		panic(fmt.Errorf("child prefix mismatch: %v != %v", p0, b))
	}

	// Replace the child.
	list.Children[int(b)-list.Min] = child
}

func (list *DenseChildList) next(b byte) *Trie {
	i := int(b)
	if i < list.Min || list.Max < i {
		return nil
	}
	return list.Children[i-list.Min]
}

func (list *DenseChildList) walk(prefix *Prefix, visitor VisitorFunc) error {
	for _, child := range list.Children {
		if child == nil {
			continue
		}
		*prefix = append(*prefix, child.Prefix...)
		if child.Item != nil {
			if err := visitor(*prefix, child.Item); err != nil {
				if err == SkipSubtree {
					*prefix = (*prefix)[:len(*prefix)-len(child.Prefix)]
					continue
				}
				*prefix = (*prefix)[:len(*prefix)-len(child.Prefix)]
				return err
			}
		}

		err := child.Children.walk(prefix, visitor)
		*prefix = (*prefix)[:len(*prefix)-len(child.Prefix)]
		if err != nil {
			return err
		}
	}

	return nil
}

func (list *DenseChildList) print(w io.Writer, indent int) {
	for _, child := range list.Children {
		if child != nil {
			child.print(w, indent)
		}
	}
}

func (list *DenseChildList) clone() ChildList {
	clones := make(Tries, cap(list.Children))

	if list.NumChildren != 0 {
		clonedCount := 0
		for i := list.HeadIndex; i < len(list.Children); i++ {
			child := list.Children[i]
			if child != nil {
				clones[i] = child.Clone()
				clonedCount++
				if clonedCount == list.NumChildren {
					break
				}
			}
		}
	}

	return &DenseChildList{
		Min:         list.Min,
		Max:         list.Max,
		NumChildren: list.NumChildren,
		HeadIndex:   list.HeadIndex,
		Children:    clones,
	}
}

func (list *DenseChildList) total() int {
	tot := 0
	for _, child := range list.Children {
		if child != nil {
			tot = tot + child.total()
		}
	}
	return tot
}
