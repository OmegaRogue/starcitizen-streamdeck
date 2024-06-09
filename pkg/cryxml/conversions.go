package cryxml

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	"github.com/pkg/errors"
)

func ByteToIndex(data []byte, offset uint) (NodeIndex, error) {
	var f uint
	size := unsafe.Sizeof(f)
	byt := data[int(offset) : int(offset)+int(size)]
	switch size {
	case 4:
		return NodeIndex(binary.LittleEndian.Uint32(byt)), nil
	case 8:
		return NodeIndex(binary.LittleEndian.Uint64(byt)), nil
	}
	return 0, errors.New("Error unmarshalling binary data")
}
func ByteToType[T any](data []byte, offset uint) (T, error) {
	var out T
	size := binary.Size(out)
	byt := data[int(offset) : int(offset)+size]
	err := binary.Read(bytes.NewReader(byt), binary.LittleEndian, &out)
	if err != nil {
		return out, errors.Wrap(err, "Error unmarshalling binary data")
	}
	return out, nil
}
