package cryxml

import (
	"fmt"
	"strings"
)

type Tree struct {
	XMLList []string
}

func (t *Tree) String() string {
	return strings.Join(t.XMLList, "\n")
}

func (t *Tree) DoAttribute(node *BinNode) string {
	xml := ""
	for i := 0; i < node.GetNumAttributes(); i++ {
		key, value, _ := node.GetAttributeByIndex(i)
		xml += fmt.Sprintf(" %s=\"%s\"", key, value)
	}
	return xml
}

func (t *Tree) DoNode(node *BinNode, level int) {
	tabs := strings.Repeat("\t", level)
	attr := t.DoAttribute(node)
	xml := fmt.Sprintf("<%s %s ", node.GetTag(), attr)
	if node.GetChildCount() < 1 {
		xml += "/>"
		t.XMLList = append(t.XMLList, tabs+xml)
	} else {
		xml += ">"
		t.XMLList = append(t.XMLList, tabs+xml)
		for i := 0; i < node.GetChildCount(); i++ {
			childRef := node.GetChild(i)
			t.DoNode(childRef, level+1)
		}
		xml = fmt.Sprintf("</%s>\n", node.GetTag())
		t.XMLList = append(t.XMLList, tabs+xml)
	}
}

func (t *Tree) BuildXML(rootRef *BinNode) {
	t.XMLList = []string{}
	t.DoNode(rootRef, 0)
}
