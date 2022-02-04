package main

import "testing"

func Test_FindNode(t *testing.T) {

	// first define a root wiht two nodes, then create children for each of two original nodes
	mockRoot := NewEmptyNode()
	mockRoot.addItems()
	mockTree := newTreeWithRoot(mockRoot, minItems)

	mockChild := NewEmptyNode()
	mockChild.addItems()
	mockRoot.addChildNode()
}

func (n *Node) addItems(keys ...string) *Node {
	for _, key := range keys {
		n.items = append(n.items, newItem(key, key))
	}
	return n
}
