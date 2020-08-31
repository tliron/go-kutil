package ard

//
// Node
//

type Node struct {
	Data interface{}

	container *Node
	key       string
}

func NewNode(data interface{}) *Node {
	return &Node{data, nil, ""}
}

var NoNode = &Node{nil, nil, ""}

func (self *Node) Get(key string) *Node {
	if self != NoNode {
		if data_, ok := self.Data.(StringMap); ok {
			if value, ok := data_[key]; ok {
				return &Node{value, self, key}
			}
		} else if data_, ok := self.Data.(Map); ok {
			if value, ok := data_[key]; ok {
				return &Node{value, self, key}
			}
		}
	}
	return NoNode
}

func (self *Node) Put(key string, value interface{}) bool {
	if self != NoNode {
		if data_, ok := self.Data.(StringMap); ok {
			data_[key] = value
			return true
		} else if data_, ok := self.Data.(Map); ok {
			data_[key] = value
			return true
		}
	}
	return false
}

func (self *Node) Append(value interface{}) bool {
	if self != NoNode {
		if data_, ok := self.Data.(List); ok {
			self.container.Put(self.key, append(data_, value))
			return true
		}
	}
	return false
}

func (self *Node) String(allowNil bool) (string, bool) {
	if self != NoNode {
		if allowNil && (self.Data == nil) {
			return "", true
		}
		value, ok := self.Data.(string)
		return value, ok
	}
	return "", false

}

func (self *Node) Integer(allowNil bool) (int64, bool) {
	if self != NoNode {
		if allowNil && (self.Data == nil) {
			return 0, true
		}
		switch value := self.Data.(type) {
		case int64:
			return value, true
		case int32:
			return int64(value), true
		case int16:
			return int64(value), true
		case int8:
			return int64(value), true
		case int:
			return int64(value), true
		}
	}
	return 0, false
}

func (self *Node) UnsignedInteger(allowNil bool) (uint64, bool) {
	if self != NoNode {
		if allowNil && (self.Data == nil) {
			return 0, true
		}
		switch value := self.Data.(type) {
		case uint64:
			return value, true
		case uint32:
			return uint64(value), true
		case uint16:
			return uint64(value), true
		case uint8:
			return uint64(value), true
		case uint:
			return uint64(value), true
		}
	}
	return 0, false
}

func (self *Node) Float(allowNil bool) (float64, bool) {
	if self != NoNode {
		if allowNil && (self.Data == nil) {
			return 0.0, true
		}
		switch value := self.Data.(type) {
		case float64:
			return value, true
		case float32:
			return float64(value), true
		}
	}
	return 0.0, false
}

func (self *Node) Boolean(allowNil bool) (bool, bool) {
	if self != NoNode {
		if allowNil && (self.Data == nil) {
			return false, true
		}
		if value, ok := self.Data.(bool); ok {
			return value, true
		}
	}
	return false, false
}

func (self *Node) StringMap(allowNil bool) (StringMap, bool) {
	if self != NoNode {
		if allowNil && (self.Data == nil) {
			return make(StringMap), true
		}
		value, ok := self.Data.(StringMap)
		return value, ok
	}
	return nil, false
}

func (self *Node) Map(allowNil bool) (Map, bool) {
	if self != NoNode {
		if allowNil && (self.Data == nil) {
			return make(Map), true
		}
		value, ok := self.Data.(Map)
		return value, ok
	}
	return nil, false
}

func (self *Node) List(allowNil bool) (List, bool) {
	if self != NoNode {
		if allowNil && (self.Data == nil) {
			return nil, true
		}
		value, ok := self.Data.(List)
		return value, ok
	}
	return nil, false
}
