package logging

//
// Hierarchy
//

type Hierarchy struct {
	root *Node
}

func NewMaxLevelHierarchy() *Hierarchy {
	return &Hierarchy{
		root: NewNode(),
	}
}

func (self *Hierarchy) AllowLevel(name []string, level Level) bool {
	return level <= self.GetMaxLevel(name)
}

func (self *Hierarchy) GetMaxLevel(name []string) Level {
	node := self.root
	for _, segment := range name {
		if child, ok := node.children[segment]; ok {
			node = child
		} else {
			break
		}
	}
	return node.maxLevel
}

func (self *Hierarchy) SetMaxLevel(name []string, level Level) {
	node := self.root
	for _, segment := range name {
		if child, ok := node.children[segment]; ok {
			node = child
		} else {
			child = NewNode()
			node.children[segment] = child
			node = child
		}
	}
	node.maxLevel = level
}

//
// Node
//

type Node struct {
	maxLevel Level
	children map[string]*Node
}

func NewNode() *Node {
	return &Node{
		maxLevel: None,
		children: make(map[string]*Node),
	}
}
