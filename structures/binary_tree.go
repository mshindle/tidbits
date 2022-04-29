package structures

import (
	"fmt"
	"math/rand"
)

// IntTree is a tree structure containing an integer
type IntTree struct {
	Left  *IntTree
	Value int
	Right *IntTree
}

// NewIntTree returns a new, random binary tree holding the values k, 2k, ..., 10k.
func NewIntTree(k int) *IntTree {
	var t *IntTree
	for _, v := range rand.Perm(10) {
		t = insert(t, (1+v)*k)
	}
	return t
}

func insert(t *IntTree, v int) *IntTree {
	if t == nil {
		return &IntTree{nil, v, nil}
	}
	if v < t.Value {
		t.Left = insert(t.Left, v)
	} else {
		t.Right = insert(t.Right, v)
	}
	return t
}

func (t *IntTree) String() string {
	if t == nil {
		return "()"
	}
	s := ""
	if t.Left != nil {
		s += t.Left.String() + " "
	}
	s += fmt.Sprint(t.Value)
	if t.Right != nil {
		s += " " + t.Right.String()
	}
	return "(" + s + ")"
}

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *IntTree, ch chan int) {
	if t == nil {
		return
	}
	Walk(t.Left, ch)
	ch <- t.Value
	Walk(t.Right, ch)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *IntTree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go Walk(t1, ch1)
	go Walk(t2, ch2)
	for i := 0; i < 10; i++ {
		if <-ch1 != <-ch2 {
			return false
		}
	}
	return true
}
