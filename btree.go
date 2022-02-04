package main

var DefaultMinItems = 128

type Tree struct {
	root     *Node
	minItems int
	maxItems int
}

type Node struct {
	bucket     *Tree
	items      []*Item
	childNodes []*Node
}

type Item struct {
	key   string
	value interface{}
}

func NewTree(minItems int) *Tree {
	return newTreeWithRoot(NewEmptyNode(), minItems)
}

func newTreeWithRoot(root *Node, minItems int) *Tree {
	bucket := &Tree{
		root: root,
	}
	bucket.root.bucket = bucket
	bucket.minItems = minItems
	bucket.maxItems = minItems * 2
	return bucket
}

func NewEmptyNode() *Node {
	return &Node{
		items:      []*Item{},
		childNodes: []*Node{},
	}
}

func (b *Tree) Find(key string) *Item {
	index, containingNode, _ := b.findKey(key, true)
	if index == -1 {
		return nil
	}
	return containingNode.items[index]
}

func (b *Tree) findKey(key string, exact bool) (int, *Node, []int) {
	n := b.root

	ancestorsIndexes := []int{0}
	for true {
		wasFound, index := n.findKey(key)
		if wasFound {
			return index, n, ancestorsIndexes
		} else {
			if n.isLeaf() {
				if exact {
					return -1, nil, nil
				}
				return index, n, ancestorsIndexes
			}
			nextChild := n.childNodes[index]
			ancestorsIndexes = append(ancestorsIndexes, index)
			n = nextChild
		}
	}
	return -1, nil, nil
}

func (n *Node) findKey(key string) (bool, int) {
	for i, existingItem := range n.items {
		if key == existingItem.key {
			return true, i
		}

		if key < existingItem.key {
			return false, i
		}
	}
	return false, len(n.items)
}

func (n *Node) isLeaf() bool {
	return len(n.childNodes) == 0
}

func newItem(key string, value interface{}) *Item {
	return &Item{
		key:   key,
		value: value,
	}
}
