package p4k

import (
	"io"

	"github.com/pkg/errors"
)

type MyEndOfCentralDirRecord struct {
	ID                   [4]byte
	DiskNumber           uint16
	DiskNumbersFromStart uint16
	NumEntriesOnDisk     uint16
	TotalNumEntries      uint16
	SizeOfCDir           uint32
	OffsetOfCDir         uint32
	CommentLen           uint16
}

type EndOfCentralDirRecord struct {
	item         *MyEndOfCentralDirRecord
	IsValid      bool
	RecordOffset int64
}

func NewEndOfCentralDirRecord(reader *recReader) (*EndOfCentralDirRecord, error) {
	r := new(EndOfCentralDirRecord)
	r.RecordOffset = -1

	tries := 10
	for !r.IsValid && tries > 0 {
		thisPos, err := reader.Pos()
		if err != nil {
			return nil, errors.Wrap(err, "Error getting position")
		}
		cPos, err := FindSignatureInPage(reader, SignatureEndOfCentralDirRecord)
		if err != nil {
			return nil, errors.Wrap(err, "Error getting signature")
		}
		if cPos >= 0 {
			r.RecordOffset = cPos

			if _, err := reader.Seek(cPos, io.SeekStart); err != nil {
				return nil, errors.Wrap(err, "Seek to cPos")
			}
			item, err := ByteToType[MyEndOfCentralDirRecord](reader)
			if err != nil {
				r.IsValid = false
				r.RecordOffset = -1
				return nil, errors.New("EOF - cannot find EndOfCentralDirRecord")
			}
			r.item = item
			r.IsValid = true
		} else {
			if _, err := reader.Seek(thisPos, io.SeekStart); err != nil {
				return nil, errors.Wrap(err, "Seek to thisPos")
			}
			if err := reader.BackwardPage(); err != nil {
				return nil, errors.Wrap(err, "Seek to cPos")
			}
			tries--
		}
	}
	return r, nil
}
