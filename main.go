package main

func main() {
	minItems := DefaultMinItems
	tree := NewTree(minItems)
	tree.findKey("foo", true)
}
