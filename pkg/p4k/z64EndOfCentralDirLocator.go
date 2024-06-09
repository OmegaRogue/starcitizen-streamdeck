package p4k

import (
	"io"

	"github.com/pkg/errors"
)

type MyZ64EndOfCentralDirLocator struct {
	ID                 [4]byte
	DiskNumber         uint32
	OffsetOfz64EofCDir uint64
	TotalNumEntries    uint32
}

type Z64EndOfCentralDirLocator struct {
	item         *MyZ64EndOfCentralDirLocator
	IsValid      bool
	RecordOffset int64
}

func NewZ64EndOfCentralDirLocator(reader *recReader) (*Z64EndOfCentralDirLocator, error) {
	r := new(Z64EndOfCentralDirLocator)
	r.RecordOffset = -1

	for !r.IsValid {
		cPos, err := FindSignatureInPageBackwards(reader, SignatureZ64EndOfCentralDirLocator)
		if err != nil {
			return nil, errors.Wrap(err, "Error getting signature")
		}
		if cPos < 0 {
			r.IsValid = false
			r.RecordOffset = -1
			return nil, errors.New("EOF - cannot find Z64EndOfCentralDirLocator")
		}

		r.RecordOffset = cPos

		if _, err := reader.Seek(cPos, io.SeekStart); err != nil {
			return nil, errors.Wrap(err, "Seek to cPos")
		}
		item, err := ByteToType[MyZ64EndOfCentralDirLocator](reader)
		if err != nil {
			r.IsValid = false
			r.RecordOffset = -1
			return nil, errors.New("EOF - cannot find Z64EndOfCentralDirLocator")
		}
		r.item = item
		r.IsValid = true
	}

	return r, nil
}
func (i *Z64EndOfCentralDirLocator) Z64EndOfCentralDir() int64 {
	if i.IsValid {
		return int64(i.item.OffsetOfz64EofCDir)
	}
	return -1
}
