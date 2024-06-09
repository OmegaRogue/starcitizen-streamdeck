package cryxml

type BinNode struct {
	Data      *BinContext
	NodeIndex uint
}

func (n *BinNode) node() *Node {
	return n.Data.Nodes[n.NodeIndex]
}
func (n *BinNode) string(nIndex uint) string {
	return n.Data.string(nIndex)
}

func (n *BinNode) GetValue(key string) string {
	nFirst := n.node().FirstAttributeIndex
	nLast := nFirst + uint32(n.GetNumAttributes())
	for i := nFirst; i < nLast; i++ {
		attrKey := n.string(uint(n.Data.Attributes[i].KeyStringOffset))
		if key != attrKey {
			continue
		}
		attrValue := n.string(uint(n.Data.Attributes[i].ValueStringOffset))
		return attrValue
	}
	return ""
}
func (n *BinNode) GetTag() string {
	return n.string(uint(n.node().TagStringOffset))
}

func (n *BinNode) IsTag(tag string) bool {
	return tag == n.GetTag()
}
func (n *BinNode) HaveAttr(key string) bool {
	return n.GetValue(key) != ""
}
func (n *BinNode) GetNumAttributes() int {
	return int(n.node().AttributeCount)
}
func (n *BinNode) GetAttributeByIndex(index int) (key string, value string, ok bool) {
	pNode := n.node()
	if index < 0 || index >= int(pNode.AttributeCount) {
		return "", "", false
	}
	attr := n.Data.Attributes[int(pNode.FirstAttributeIndex)+index]
	key = n.string(uint(attr.KeyStringOffset))
	value = n.string(uint(attr.ValueStringOffset))
	return key, value, true
}
func (n *BinNode) GetAttr(key string) (value string, ok bool) {
	sValue := n.GetValue(key)
	if sValue != "" {
		return sValue, true
	}
	return "", false
}

func (n *BinNode) GetChildCount() int {
	return int(n.node().ChildCount)
}

func (n *BinNode) GetChild(i int) *BinNode {
	pNode := n.node()
	if i < 0 || i > n.GetChildCount() {
		return nil
	}
	return n.Data.BinaryNodes[n.Data.ChildIndices[int(pNode.FirstChildIndex)+i]]
}

func (n *BinNode) FindChild(tag string) *BinNode {
	pNode := n.node()
	nFirst := pNode.FirstChildIndex
	nAfterLast := pNode.FirstChildIndex + uint32(pNode.ChildCount)
	for i := nFirst; i < nAfterLast; i++ {
		sChildTag := n.string(uint(n.Data.Nodes[n.Data.ChildIndices[i]].TagStringOffset))
		if tag != sChildTag {
			continue
		}
		return n.Data.BinaryNodes[n.Data.ChildIndices[i]]
	}
	return nil
}

func (n *BinNode) GetParent() *BinNode {
	pNode := n.node()
	return n.Data.BinaryNodes[pNode.ParentIndex]
}

func (n *BinNode) GetContent() string {
	return n.string(uint(n.node().ContentStringOffset))
}
