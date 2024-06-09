package p4k

import (
	"io"
	"time"

	"github.com/pkg/errors"
)

type File struct {
	fileDateTime       time.Time
	Filename           string
	FileHeaderPosition int64
	CompressedSize     int64
	FileSize           int64
}

func (f *File) FileModifyDate() time.Time {
	return f.fileDateTime
}

func NewFile(dirEntry *DirectoryEntry) *File {
	f := new(File)
	f.Filename = dirEntry.Filename
	f.CompressedSize = dirEntry.FileSizeComp
	f.FileSize = dirEntry.FileSizeUnComp
	f.fileDateTime = dirEntry.FileModifyDate
	f.FileHeaderPosition = dirEntry.FileHeaderOffset()
	return f
}
func (f *File) GetFile(reader *recReader) ([]byte, error) {
	if f.FileHeaderPosition < 0 {
		return nil, errors.New("File Header not found")
	}
	if !reader.IsOpen() {
		return nil, errors.New("Reader not open")
	}
	_, err := reader.Seek(f.FileHeaderPosition, io.SeekStart)
	if err != nil {
		return nil, errors.Wrap(err, "Seek File Header")
	}
	fileHeader, err := NewFileHeader(reader)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read File Header")
	}

	fContent, err := fileHeader.GetFile(reader)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read File")
	}

	return fContent, nil
}
