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
