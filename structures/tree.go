package structures

// TraversalType defines which type of traversal should occur on a tree
type TraversalType int

const (
	PreOrder TraversalType = 1 << iota
	InOrder
	PostOrder
)

// Tree is a binary tree
type Tree struct {
	Root *Node
}

// NewTree creates a new tree object from the given root
func NewTree(root *Node) *Tree {
	return &Tree{Root: root}
}

// String creates a string representation of the tree
func (tr *Tree) String() string {
	return tr.Root.String()
}

// Insert inserts a Node to a Tree without replacement.
func (tr *Tree) Insert(nd *Node) {
	if tr.Root == nd {
		return
	}
	tr.Root = tr.Root.insert(nd)
}

// Traverse moves across the tree executing visit on each node as determined by TraversalType
func (tr *Tree) Traverse(visit Visit, t TraversalType) {
	switch t {
	case PreOrder:
		tr.Root.TraversePre(visit)
	case InOrder:
		tr.Root.TraverseIn(visit)
	case PostOrder:
		tr.Root.TraversePost(visit)
	}

}
