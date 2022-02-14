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
		newRoot := NewNode(b, []*Item{}, []*Node{b.root})
		newRoot.split(b.root, 0)
		b.root = newRoot
	}
}

func (b *Tree) Remove(key string) {
	// find path to where deletion should happen
	removeItemIndex, nodeToRemove, ancestorIndexes := b.findKey(key, true)

	if nodeToRemove.isLeaf() {
		nodeToRemove.removeItemFromLeaf(removeItemIndex)
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

func (n *Node) addItem(item *Item, insertionIndex int) int {
	if len(n.items) == insertionIndex {
		n.items = append(n.items, item)
		return insertionIndex
	}
	n.items = append(n.items[:insertionIndex+1], n.items[insertionIndex:]...)
	n.items[insertionIndex] = item
	return insertionIndex
}

func (n *Node) addChild(node *Node, insertionIndex int) {
	if len(n.childNodes) == insertionIndex {
		n.childNodes = append(n.childNodes, node)

	}
	n.childNodes = append(n.childNodes[:insertionIndex+1], n.childNodes[insertionIndex:]...)
	n.childNodes[insertionIndex] = node
}

func (n *Node) split(modifiedNode *Node, insertionIndex int) {
	i := 0
	nodeSize := n.bucket.minItems

	for modifiedNode.isOverPopulated() {
		middleItem := modifiedNode.items[nodeSize]
		var newNode *Node
		if modifiedNode.isLeaf() {
			newNode = NewNode(n.bucket, modifiedNode.items[nodeSize+1:], []*Node{})
			modifiedNode.items = modifiedNode.items[:nodeSize]
		} else {
			newNode = NewNode(n.bucket, modifiedNode.items[nodeSize+1:], modifiedNode.childNodes[i+1:])
			modifiedNode.items = modifiedNode.items[:nodeSize]
			modifiedNode.childNodes = modifiedNode.childNodes[:nodeSize+1]
		}
		n.addItem(middleItem, insertionIndex)
		if len(n.childNodes) == insertionIndex+1 { // If middle of list, then move items forward
			n.childNodes = append(n.childNodes, newNode)
		} else {
			n.childNodes = append(n.childNodes[:insertionIndex+1], n.childNodes[insertionIndex:]...)
			n.childNodes[insertionIndex+1] = newNode
		}

		insertionIndex += 1
		i += 1
		modifiedNode = newNode
	}
}

func (n *Node) rebalanceRemove(unbalancedNodeIndex int) {
	pNode := n
	unbalancedNode := pNode.childNodes[unbalancedNodeIndex]

	// Right rotate
	var leftNode *Node
	if unbalancedNodeIndex != 0 {
		leftNode = pNode.childNodes[unbalancedNodeIndex-1]
		if len(leftNode.items) > n.bucket.minItems {
			rotateRight(leftNode, pNode, unbalancedNode, unbalancedNodeIndex)
			return
		}
	}

	// Left Balance
	var rightNode *Node
	if unbalancedNodeIndex != len(pNode.childNodes)-1 {
		rightNode = pNode.childNodes[unbalancedNodeIndex+1]
		if len(rightNode.items) > n.bucket.minItems {
			rotateLeft(unbalancedNode, pNode, rightNode, unbalancedNodeIndex)
			return
		}
	}

	merge(pNode, unbalancedNodeIndex)
}

func (n *Node) removeItemFromLeaf(index int) {
	n.items = append(n.items[:index], n.items[index+1:]...)
}

func (n *Node) removeItemFromInternal(index int) []int {
	affectedNodes := make([]int, 0)
	affectedNodes = append(affectedNodes, index)

	aNode := n.childNodes[index]
	for !aNode.isLeaf() {
		traversingIndex := len(n.childNodes) - 1
		aNode = n.childNodes[traversingIndex]
		affectedNodes = append(affectedNodes, traversingIndex)
	}

	n.items[index] = aNode.items[len(aNode.items)-1]
	aNode.items = aNode.items[:len(aNode.items)-1]
	return affectedNodes
}

func rotateRight(aNode, pNode, bNode *Node, bNodeIndex int) {
	aNodeItem := aNode.items[len(aNode.items)-1]
	aNode.items = aNode.items[:len(aNode.items)-1]

	// Get item from parent node and assign the aNodeItem item instead
	pNodeItemIndex := bNodeIndex - 1
	if isFirst(bNodeIndex) {
		pNodeItemIndex = 0
	}
	pNodeItem := pNode.items[pNodeItemIndex]
	pNode.items[pNodeItemIndex] = aNodeItem

	// Assign parent item to b and make it first
	bNode.items = append([]*Item{pNodeItem}, bNode.items...)

	// If it's a inner leaf then move children as well.
	if !aNode.isLeaf() {
		childNodeToShift := aNode.childNodes[len(aNode.childNodes)-1]
		aNode.childNodes = aNode.childNodes[:len(aNode.childNodes)-1]
		bNode.childNodes = append([]*Node{childNodeToShift}, bNode.childNodes...)
	}
}

func rotateLeft(aNode, pNode, bNode *Node, bNodeIndex int) {
	bNodeItem := bNode.items[0]
	bNode.items = bNode.items[1:]

	// Get item from parent node and assign the bNodeItem item instead
	pNodeItemIndex := bNodeIndex
	if isLast(bNodeIndex, pNode) {
		pNodeItemIndex = len(pNode.items) - 1
	}
	pNodeItem := pNode.items[pNodeItemIndex]
	pNode.items[pNodeItemIndex] = bNodeItem

	// Assign parent item to a and make it last
	aNode.items = append(aNode.items, pNodeItem)

	// If it's a inner leaf then move children as well.
	if !bNode.isLeaf() {
		childNodeToShift := bNode.childNodes[0]
		bNode.childNodes = bNode.childNodes[1:]
		aNode.childNodes = append(aNode.childNodes, childNodeToShift)
	}
}

func merge(pNode *Node, unbalancedNodeIndex int) {
	unbalancedNode := pNode.childNodes[unbalancedNodeIndex]
	if unbalancedNodeIndex == 0 {

		aNode := unbalancedNode
		bNode := pNode.childNodes[unbalancedNodeIndex+1]

		// Take the item from the parent, remove it and add it to the unbalanced node
		pNodeItem := pNode.items[0]
		pNode.items = pNode.items[1:]
		aNode.items = append(aNode.items, pNodeItem)

		//merge the bNode to aNode and remove it. Handle its child nodes as well.
		aNode.items = append(aNode.items, bNode.items...)
		pNode.childNodes = append(pNode.childNodes[0:1], pNode.childNodes[2:]...)
		if !bNode.isLeaf() {
			aNode.childNodes = append(aNode.childNodes, bNode.childNodes...)
		}
	} else {

		bNode := unbalancedNode
		aNode := pNode.childNodes[unbalancedNodeIndex-1]

		// Take the item from the parent, remove it and add it to the unbalanced node
		pNodeItem := pNode.items[unbalancedNodeIndex-1]
		pNode.items = append(pNode.items[:unbalancedNodeIndex-1], pNode.items[unbalancedNodeIndex:]...)
		aNode.items = append(aNode.items, pNodeItem)

		aNode.items = append(aNode.items, bNode.items...)
		pNode.childNodes = append(pNode.childNodes[:unbalancedNodeIndex], pNode.childNodes[unbalancedNodeIndex+1:]...)
		if !aNode.isLeaf() {
			bNode.childNodes = append(aNode.childNodes, bNode.childNodes...)
		}
	}
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
