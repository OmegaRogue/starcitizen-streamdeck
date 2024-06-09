package cryxml

import (
	"unsafe"
)

type NodeIndex uint32

func (n NodeIndex) MySize() uint {
	return uint(unsafe.Sizeof(n))
}

type Sizer interface {
	MySize() uint
}

func NodeIndexSize() uint {
	return new(NodeIndex).MySize()
}
func NodeSize() uint {
	return new(Node).MySize()
}
func AttributeSize() uint {
	return new(Attribute).MySize()
}

func HeaderSize() uint {
	return new(Header).MySize()
}

type NodeRef interface {
	GetValue(key string) string
	GetTag() string
	IsTag(tag string) bool
	HaveAttr(key string) bool
	GetNumAttributes() int
	GetAttributeByIndex(index int) (key string, value string, ok bool)
	GetAttr(key string) (value string, ok bool)
	GetChildCount() int
	GetChild(i int) NodeRef
	FindChild(tag string) NodeRef
	GetParent() NodeRef
	GetContent() string
}

func ToString(data []byte, size uint) string {
	s := ""
	for i := uint(0); i < size; i++ {
		if data[i] != 0 {
			s += string(data[i])
		} else {
			break
		}
	}
	return s
}
