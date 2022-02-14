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

func (b *Tree) Put(key string, value interface{}) {
	// Find the path to the node where insertion happens
	i := newItem(key, value)
	insertionIndex, nodeToInsertIn, ancestorsIndex := b.findKey(i.key, false)
	// Add item to the leaf node
	nodeToInsertIn.addItem(i, insertionIndex)

	ancestors := b.getNodes(ancestorsIndex)
	// Rebalance the nodes all the way up. Start from one node before the last and go up
	for i := len(ancestors) - 2; i >= 0; i-- {
		pnode := ancestors[i]
		node := ancestors[i+1]
		nodeIndex := ancestorsIndex[i+1]
		if node.isOverPopulated() {
			pnode.split(node, nodeIndex)
		}
	}
	// Handle root
	if b.root.isOverPopulated() {
		newRoot := newNode(b, []*Item{}, []*Node{b.root})
		newRoot.split(b.root, 0)
		b.root = newRoot
	}
}

func (b *Tree) Remove(key string) {
	// find path to where deletion should happen
	removeItemIndex, nodeToRemove, ancestorIndexes := b.findKey(key, true)

	if nodeToRemove.isLeaf() {
		nodeToRemove.removeFromLeaf(removeItemIndex)
	} else {
		affectedNodes := nodeToRemove.removeItemFromInternal(removeItemIndex)
		ancestorIndexes = append(ancestorIndexes, affectedNodes...)
	}

	ancestors := b.getNodes(ancestorIndexes)

	// Rebalance the nodes all the way up to root
	for i := len(ancestors) - 2; i >= 0; i-- {
		pnode := ancestors[i]
		node := ancestors[i+1]
		if node.isUnderPopulated() {
			pnode.rebalanceRemove(ancestorIndexes[i+1])
		}
	}
	// if the root has no more items post rebalnce
	if len(b.root.items) == 0 && len(b.root.childNodes) > 0 {
		b.root = ancestors[1]
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

func (b *Tree) getNodes(indexes []int) []*Node {
	nodes := []*Node{b.root}
	child := b.root
	for i := 1; i < len(indexes); i++ {
		child = child.childNodes[indexes[i]]
		nodes = append(nodes, child)
	}
	return nodes
}

func NewNode(bucket *Tree, value []*Item, childNodes []*Node) *Node {
	return &Node{
		bucket,
		value,
		childNodes,
	}
}

func isLast(index int, parentNode *Node) bool {
	return index == len(parentNode.items)
}

func isFirst(index int) bool {
	return index == 0
}

func (n *Node) isOverPopulated() bool {
	return len(n.items) > n.bucket.maxItems
}

func (n *Node) isUnderPopulated() bool {
	return len(n.items) < n.bucket.minItems
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
