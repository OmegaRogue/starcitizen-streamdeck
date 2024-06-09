package p4k

import (
	"bytes"
	"io"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/klauspost/compress/zstd"
	"github.com/pkg/errors"
)

type MyFileHeaderRecord struct {
	ID                [4]byte
	ExtractVersion    uint16
	BitFlags          uint16
	CompressionMethod uint16
	LastModTime       uint16
	LastModDate       uint16
	CRC32             uint32
	CompressedSize    uint32
	UncompressedSize  uint32
	FilenameLength    uint16
	ExtraFieldLength  uint16
}

type MyZ64FileHeaderExtraRecord struct {
	ID                uint16
	Size              uint16
	UncompressedSize  uint64
	CompressedSize    uint64
	LocalHeaderOffset uint64
	DiskStart         uint32
}

type FileHeader struct {
	extraBytes     []byte
	item           *MyFileHeaderRecord
	z64Item        *MyZ64FileHeaderExtraRecord
	IsValid        bool
	Filename       string
	RecordOffset   int64
	FileOffset     int64
	FileSizeComp   int64
	FileSizeUncomp int64
	FileModifyDate time.Time
}

func NewFileHeader(reader *recReader) (*FileHeader, error) {
	fh := new(FileHeader)
	if !reader.IsOpen() {
		return nil, errors.New("Reader not open")
	}
	var cPos int64
	length, err := reader.Length()
	if err != nil {
		return nil, errors.New("Get length failed")
	}
	for ok := true; ok; ok = cPos < length && !fh.IsValid {
		if err := reader.AdvancePage(); err != nil {
			return nil, errors.New("Advance Page failed")
		}
		cPos, err = reader.Pos()
		if err != nil {
			return nil, errors.New("Get pos failed")
		}
		fh.RecordOffset = cPos
		fh.item, err = ByteToType[MyFileHeaderRecord](reader)
		if err != nil {
			return nil, errors.New("Read record failed")
		}
		fh.IsValid = linq.From(fh.item.ID).SequenceEqual(linq.From(SignatureLocalFileHeaderCry))

		if fh.IsValid {
			if fh.item.FilenameLength > 0 {
				if err := fh.ReadFilename(reader); err != nil {
					return nil, errors.New("Read filename failed")
				}
			}
			if fh.item.ExtraFieldLength > 0 {
				if err := fh.ReadExtraData(reader); err != nil {
					return nil, errors.New("Read extra data failed")
				}
			}
			fh.FileModifyDate = TimestampFromDos(fh.item.LastModDate, fh.item.LastModTime)

			if fh.item.CompressedSize < 0xffffffff {
				fh.FileSizeComp = int64(fh.item.CompressedSize)
				fh.FileSizeUncomp = int64(fh.item.UncompressedSize)
			} else {
				fh.FileSizeComp = int64(fh.z64Item.CompressedSize)
				fh.FileSizeUncomp = int64(fh.z64Item.UncompressedSize)
			}
			fh.FileOffset, err = reader.Pos()
			if err != nil {
				return nil, errors.New("Get file offset failed")
			}

			if _, err := reader.Seek(fh.FileSizeComp, io.SeekCurrent); err != nil {
				return nil, errors.New("Seek failed")
			}
		} else {
			fh.RecordOffset = -1
			fh.FileOffset = -1
			fh.FileSizeComp = 0
			fh.FileSizeUncomp = 0
		}
		length, err = reader.Length()
		if err != nil {
			return nil, errors.New("Get length failed")
		}
	}

	if !fh.IsValid {
		if linq.From(fh.item.ID).SequenceEqual(linq.From(SignatureCentralDirRecord)) {
			return nil, errors.New("EOF - found Central Directory header")
		}
		return nil, errors.New("Cannot process fileheader")
	}
	return fh, nil
}

func (fh *FileHeader) GetFile(reader *recReader) ([]byte, error) {
	if !fh.IsValid {
		return nil, errors.New("Header not valid")
	}
	if !reader.IsOpen() {
		return nil, errors.New("Reader not open")
	}
	_, err := reader.Seek(fh.FileOffset, io.SeekStart)
	if err != nil {
		return nil, errors.Wrap(err, "Seek File")
	}
	fileBytes := make([]byte, fh.FileSizeComp)
	if _, err := reader.Read(fileBytes); err != nil {
		return nil, errors.Wrap(err, "Read File")
	}
	var decompFile []byte
	// ZStd
	if fh.item.CompressionMethod == 0x64 {
		read, err := zstd.NewReader(bytes.NewReader(fileBytes))
		if err != nil {
			return nil, errors.Wrap(err, "cannot create decoder")
		}
		decompFile, err = io.ReadAll(read)
		if err != nil {
			return nil, errors.Wrap(err, "cannot decode file")
		}
		return decompFile, nil
	}
	decompFile = fileBytes
	return decompFile, nil
}

func (fh *FileHeader) ReadFilename(reader *recReader) error {
	fileNameBytes := make([]byte, fh.item.FilenameLength)
	_, err := reader.Read(fileNameBytes)
	if err != nil {
		return errors.Wrap(err, "failed reading filename")
	}
	fh.Filename = string(fileNameBytes)
	return nil
}
func (fh *FileHeader) ReadExtraData(reader *recReader) error {
	z64Item, err := ByteToType[MyZ64FileHeaderExtraRecord](reader)
	if err != nil {
		return errors.Wrap(err, "Unmarshal z64 Extra data")
	}
	fh.z64Item = z64Item
	fh.extraBytes = make([]byte, fh.item.ExtraFieldLength-Z64FileHeaderExtraRecordLength)
	if _, err := reader.Read(fh.extraBytes); err != nil {
		return errors.Wrap(err, "Read extra data")
	}
	return nil
}

func (fh *FileHeader) Close() error {
	fh.IsValid = false
	fh.extraBytes = nil
	return nil
}
