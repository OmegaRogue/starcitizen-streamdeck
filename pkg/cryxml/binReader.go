package cryxml

import (
	"io"

	"github.com/pkg/errors"
)

type BinReader struct {
}

func (b *BinReader) Read(_ []byte) (n int, err error) {
	//TODO implement me
	panic("implement me")
}

func (b *BinReader) LoadFromBuffer(binBuffer []byte) (*BinNode, error) {
	fileSize := uint(len(binBuffer))
	if fileSize < HeaderSize() {
		return nil, errors.New("Data not BinXml")
	}
	pData, err := b.Create(binBuffer)
	if err != nil {
		return nil, errors.Wrap(err, "Create failed")
	}
	n := pData.BinaryNodes[0]
	return n, nil
}

func (b *BinReader) Create(fileContents []byte) (*BinContext, error) {
	pData := NewBinContext()
	header, err := ByteToType[Header](fileContents, 0)
	if err != nil {
		return nil, errors.Wrap(err, "Cant read header")
	}
	if !header.HasCorrectSignature() {
		return nil, errors.New("Invalid Signature")
	}
	pData.BinaryNodes = []*BinNode{}
	for i := 0; i < int(header.NodeCount); i++ {
		pData.BinaryNodes = append(pData.BinaryNodes, new(BinNode))
	}
	pData.Attributes = make([]*Attribute, header.AttributeCount)
	incr := AttributeSize()
	for aIdx := uint(0); aIdx < uint(header.AttributeCount); aIdx++ {
		attr, err := ByteToType[Attribute](fileContents, uint(header.AttributeTablePosition)+aIdx*incr)
		if err != nil {
			return nil, errors.Wrap(err, "Cant read Attribute")
		}
		pData.Attributes[aIdx] = &attr
	}

	pData.ChildIndices = make([]NodeIndex, header.ChildCount)
	incr = NodeIndexSize()
	for aIdx := uint(0); aIdx < uint(header.ChildCount); aIdx++ {
		idx, err := ByteToIndex(fileContents, uint(header.ChildTablePosition)+aIdx*incr)
		if err != nil {
			return nil, errors.Wrap(err, "Cant read NodeIndex")
		}
		pData.ChildIndices[aIdx] = idx
	}

	pData.Nodes = make([]*Node, header.NodeCount)
	incr = NodeSize()
	for aIdx := uint(0); aIdx < uint(header.NodeCount); aIdx++ {
		node, err := ByteToType[Node](fileContents, uint(header.NodeTablePosition)+aIdx*incr)
		if err != nil {
			return nil, errors.Wrap(err, "Cant read BinNode")
		}
		pData.Nodes[aIdx] = &node
	}

	pData.StringDataLength = uint(header.StringDataSize)
	pData.StringData = fileContents[header.StringDataPosition : header.StringDataPosition+header.StringDataSize+1]

	for i := uint(0); i < uint(header.NodeCount); i++ {
		pData.BinaryNodes[i].Data = pData  // add data space
		pData.BinaryNodes[i].NodeIndex = i // self ref..
	}
	return pData, nil
}

func init() {
	var assertion io.Reader
	assertion = &BinReader{}
	_ = assertion
}
