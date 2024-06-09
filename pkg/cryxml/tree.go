package cryxml

import (
	"fmt"
	"strings"
)

type Tree struct {
	XmlList []string
}

func (t *Tree) String() string {
	return strings.Join(t.XmlList, "\n")
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
	xml := ""
	attr := t.DoAttribute(node)
	xml = fmt.Sprintf("<%s %s ", node.GetTag(), attr)
	if node.GetChildCount() < 1 {
		xml += "/>"
		t.XmlList = append(t.XmlList, tabs+xml)
	} else {
		xml += ">"
		t.XmlList = append(t.XmlList, tabs+xml)
		for i := 0; i < node.GetChildCount(); i++ {
			childRef := node.GetChild(i)
			t.DoNode(childRef, level+1)
		}
		xml = fmt.Sprintf("</%s>\n", node.GetTag())
		t.XmlList = append(t.XmlList, tabs+xml)
	}
}

func (t *Tree) BuildXml(rootRef *BinNode) {
	t.XmlList = []string{}
	t.DoNode(rootRef, 0)
}
