package p4k

import (
	"io"
	"time"

	"github.com/pkg/errors"
)

type MyDirectoryEntryRecord struct {
	// central file header signature
	ID [4]byte

	// version made by
	VersionMadeBy uint16

	// version made by
	ExtractVersion uint16

	// general purpose bit flag
	BitFlags uint16

	// compression method
	CompressionMethod uint16

	// last mod file time
	LastModTime uint16

	// last mod file date
	LastModDate uint16

	// crc-32
	CRC32 uint32

	// compressed size
	CompressedSize uint32

	// uncompressed size
	UncompressedSize uint32

	// file name length
	FilenameLength uint16

	// extra field length
	ExtraFieldLength uint16

	// file comment length
	FilecommentLength uint16

	// disk number start
	DiskNumberStart uint16

	// internal file attributes
	IntFileAttr uint16

	// external file attributes
	ExtFileAttr uint32

	// relative offset of local header
	RelOffsetHeader uint32
}

type MyZ64DirectoryEntryExtraRecord struct {
	// (Zip64 ExtraHeader Signature)
	ID uint16

	// Size of this "extra" block
	Size uint16

	// Original uncompressed file size
	UncompressedSize uint64

	// Size of compressed data
	CompressedSize uint64

	// Offset of local header record
	LocalHeaderOffset uint64

	// Number of the disk on which this file starts
	DiskStart uint32
}

type DirectoryEntry struct {
	extraBytes     []byte
	item           *MyDirectoryEntryRecord
	z64Item        *MyZ64DirectoryEntryExtraRecord
	FileModifyDate time.Time
	FileSizeUnComp int64
	FileSizeComp   int64
	Filename       string
	RecordOffset   int64
	IsValid        bool
}

func NewDirectoryEntry(reader *recReader) (*DirectoryEntry, error) {
	de := new(DirectoryEntry)
	if !reader.IsOpen() {
		return nil, errors.New("Reader not open")
	}
	cPos, err := FindSignatureInPage(reader, SignatureCentralDirRecord)
	if err != nil {
		return nil, errors.Wrap(err, "Find Central Dir Record Signature")
	}
	de.RecordOffset = cPos

	if _, err := reader.Seek(cPos, io.SeekStart); err != nil {
		return nil, errors.Wrap(err, "Seek to Central Dir Record")
	}
	de.item, err = ByteToType[MyDirectoryEntryRecord](reader)
	de.IsValid = true
	if de.item.FilenameLength > 0 {
		err := de.ReadFilename(reader)
		if err != nil {
			return nil, errors.Wrap(err, "Read Filename")
		}
	}
	if de.item.ExtraFieldLength > 0 {
		err := de.ReadExtradata(reader)
		if err != nil {
			return nil, errors.Wrap(err, "Read Extra Data")
		}
	}
	de.FileModifyDate = TimestampFromDos(de.item.LastModDate, de.item.LastModTime)

	if de.item.CompressedSize < 0xffffffff {
		de.FileSizeComp = int64(de.item.CompressedSize)
		de.FileSizeUnComp = int64(de.item.UncompressedSize)
	} else {
		de.FileSizeComp = int64(de.z64Item.CompressedSize)
		de.FileSizeUnComp = int64(de.z64Item.UncompressedSize)
	}

	return de, nil
}

func (de *DirectoryEntry) FileHeaderOffset() int64 {
	if de.IsValid {
		return int64(de.z64Item.LocalHeaderOffset)
	}
	return -1
}

func (de *DirectoryEntry) ReadFilename(reader *recReader) error {
	filenameBytes := make([]byte, de.item.FilenameLength)
	_, err := reader.Read(filenameBytes)
	if err != nil {
		return errors.Wrap(err, "Cannot read filename")
	}
	de.Filename = string(filenameBytes)
	return nil
}

func (de *DirectoryEntry) ReadExtradata(reader *recReader) error {
	z64Item, err := ByteToType[MyZ64DirectoryEntryExtraRecord](reader)
	if err != nil {
		return errors.Wrap(err, "Unmarshal z64 Extra data")
	}
	de.z64Item = z64Item
	de.extraBytes = make([]byte, de.item.ExtraFieldLength-Z64DirectoryEntryExtraRecordLength)
	if _, err := reader.Read(de.extraBytes); err != nil {
		return errors.Wrap(err, "Read extra data")
	}
	return nil
}
