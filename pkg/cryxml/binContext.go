package cryxml

type BinContext struct {
	Attributes       []*Attribute
	BinaryNodes      []*BinNode
	ChildIndices     []NodeIndex
	Nodes            []*Node
	StringData       []byte
	StringDataLength uint
}

func (c *BinContext) string(sOffset uint) string {
	return ToString(c.StringData[sOffset:c.StringDataLength], 999)
}

func NewBinContext() *BinContext {
	c := new(BinContext)
	c.Attributes = []*Attribute{}
	c.BinaryNodes = []*BinNode{}
	c.ChildIndices = []NodeIndex{}
	c.Nodes = []*Node{}
	c.StringData = []byte{}
	return c
}

//TODO public string _string(uint sOffset) => Conversions.ToString(StringData.SliceE(sOffset, StringDataLength));
