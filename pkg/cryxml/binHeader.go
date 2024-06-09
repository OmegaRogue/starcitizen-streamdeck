package cryxml

import "encoding/binary"

type Node struct {
	TagStringOffset     uint32
	ContentStringOffset uint32
	AttributeCount      uint16
	ChildCount          uint16
	ParentIndex         uint32
	FirstAttributeIndex uint32
	FirstChildIndex     uint32
	Pad                 uint32
}

func (n *Node) MySize() uint {
	return uint(binary.Size(n))
}

type Attribute struct {
	KeyStringOffset   uint32
	ValueStringOffset uint32
}

func (a *Attribute) MySize() uint {
	return uint(binary.Size(a))
}

type Header struct {
	SzSignatureRaw         [8]byte
	XMLSize                uint32
	NodeTablePosition      uint32
	NodeCount              uint32
	AttributeTablePosition uint32
	AttributeCount         uint32
	ChildTablePosition     uint32
	ChildCount             uint32
	StringDataPosition     uint32
	StringDataSize         uint32
}

func (h *Header) HasCorrectSignature() bool {
	return h.SzSignature() == "CryXmlB"
}

func (h *Header) SzSignature() string {
	return ToString(h.SzSignatureRaw[0:8], 8)
}

func (h *Header) MySize() uint {
	return uint(binary.Size(h))
}
