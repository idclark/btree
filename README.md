# Btree
An in-mem B-Tree structure implemented in Go. Learning project only. 
## Motivation
B-Trees are a very (the most?) prominent data structure leveraged by modern RDMS systems including Postgres, MySql, and Oracle Database. Before jumping into B-Trees specifically it's worth noting some potential shortcomings of binary search trees for databases. 

A binary tree is a structure used for storing data in an ordered way. Each node is identified by a key, a value associated with the key and two pointers (thus is its namesake _binary_ tree) for the child nodes. The left child node must be less than its direct parent, and the right child node must be greater than its direct parent. 

This way, we can easily find elements by looking at each node's value and descend accordingly. If the searched value is smaller than the current value, descend to the left child; if it's great, descend to the right. 

#### On-Disk Shortcomings
The storage engine must commit the tree to disk to make the data durable. Data is stored in page frames. The data is laid out contiguously on the disk and is usually 4KiB in size (depending on CPU). Balanced searched trees aren't great on-disk representations for two reasons:

* Locality. Elements are entered in random order, so there's no guarantee that nodes reside close to one another and may spread across pages, causing excessive disk access. 

* Tree Height. Binary search trees only have two children, thus the height of the tree increases very quickly. For each level we have to compare and descend to the node below, and again, requires additional disk access. 

### Enter B-Trees
B-Tree is a self-balancing tree structure that maintains sorted data and allows searches, generalizing the binary search, and allows for nodes with more than two children. 

B-Trees are better suited for storage systems by solving for the two main shortcomings mentioned in the above section. 
* Each node is the size of a disk page. The locality id increased as keys and values reside next to one another on the same disk page, so fewer disk accesses are required for any given search. 
* The trees's height is smaller by having more children in each node- a concept commonly reffered to as "higher fanout". 
## Example

```go
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
```

## Sources 
- [Alex Petrov: What Every Programmer has to know about Database Storage](https://www.youtube.com/watch?v=e1wbQPbFZdk)
- [Alex Petrov: Database Internals](https://www.oreilly.com/library/view/database-internals/9781492040330/)
  - [Implementing B-Trees](https://www.oreilly.com/library/view/database-internals/9781492040330/ch04.html)
