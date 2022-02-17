package logging

//
// Hierarchy
//

type Hierarchy struct {
	root *Node
	//lock util.RWLocker
}

func NewHierarchy() *Hierarchy {
	return &Hierarchy{
		root: NewNode(),
		//lock: util.NewDefaultRWLocker(),
	}
}

func (self *Hierarchy) AllowLevel(id []string, level Level) bool {
	//self.lock.RLock()
	maxLevel := self.root.GetMaxLevel(id)
	//self.lock.RUnlock()
	return level <= maxLevel
}

func (self *Hierarchy) SetMaxLevel(id []string, level Level) {
	//self.lock.Lock()
	self.root.SetMaxLevel(id, level)
	//self.lock.Unlock()
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

func (self *Node) GetMaxLevel(id []string) Level {
	node := self
	for _, i := range id {
		if child, ok := node.children[i]; ok {
			node = child
		} else {
			break
		}
	}
	return node.maxLevel
}

func (self *Node) SetMaxLevel(id []string, level Level) {
	node := self
	for _, i := range id {
		if child, ok := node.children[i]; ok {
			node = child
		} else {
			child = NewNode()
			node.children[i] = child
			node = child
		}
	}
	node.maxLevel = level
}
