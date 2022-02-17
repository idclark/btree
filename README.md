# Btree

## Example

`
package main

import "fmt"

func main() {
	minItems := DefaultMinItems
	tree := NewTree(minItems)
	tree.Put("Key1", "1")
	tree.Put("Key2", "2")
	Key1 := tree.Find("Key1")

	fmt.Printf("Returned value is key: %s value: %s \n", Key1.key, Key1.value)

}
`
