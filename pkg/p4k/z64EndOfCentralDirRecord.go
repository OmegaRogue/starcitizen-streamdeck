package p4k

import (
	"io"

	"github.com/pkg/errors"
)

type MyZ64EndOfCentralDirRecord struct {
	ID                   [4]byte
	SizeOfZ64CDir        uint64
	VersionMadeBy        uint16
	ExtractVersion       uint16
	DiskNumber           uint32
	DiskNumbersFromStart uint32
	NumEntriesOnDisk     uint64
	NumEntriesTotal      uint64
	SizeOfCDir           uint64
	OffsetOfZ64CDir      uint64
}

type Z64EndOfCentralDirRecord struct {
	item         *MyZ64EndOfCentralDirRecord
	IsValid      bool
	RecordOffset int64
}

func NewZ64EndOfCentralDirRecord(reader *recReader) (*Z64EndOfCentralDirRecord, error) {
	r := new(Z64EndOfCentralDirRecord)
	r.RecordOffset = -1

	for !r.IsValid {
		cPos, err := FindSignatureInPageBackwards(reader, SignatureZ64EndOfCentralDirRec)
		if err != nil {
			return nil, errors.Wrap(err, "Error getting signature")
		}
		if cPos < 0 {
			r.IsValid = false
			r.RecordOffset = -1
			return nil, errors.New("EOF - cannot find Z64EndOfCentralDirRecord")
		}

		r.RecordOffset = cPos

		if _, err := reader.Seek(cPos, io.SeekStart); err != nil {
			return nil, errors.Wrap(err, "Seek to cPos")
		}
		item, err := ByteToType[MyZ64EndOfCentralDirRecord](reader)
		if err != nil {
			r.IsValid = false
			r.RecordOffset = -1
			return nil, errors.New("EOF - cannot find Z64EndOfCentralDirRecord")
		}
		r.item = item
		r.IsValid = true
	}

	return r, nil
}
func (i *Z64EndOfCentralDirRecord) Z64StartOfCentralDir() int64 {
	if i.IsValid {
		return int64(i.item.OffsetOfZ64CDir)
	}
	return -1
}
func (i *Z64EndOfCentralDirRecord) NumberOfEntries() int64 {
	if i.IsValid {
		return int64(i.item.NumEntriesTotal)
	}
	return -1
}
