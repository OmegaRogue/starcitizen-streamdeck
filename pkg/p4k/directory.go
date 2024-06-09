package p4k

import (
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"starcitizen-streamdeck/internal/util"
	"starcitizen-streamdeck/pkg/zip"
)

type Directory struct {
	endOfCentralDirRecord     *EndOfCentralDirRecord
	z64EndOfCentralDirLocator *Z64EndOfCentralDirLocator
	z64EndOfCentralDirRecord  *Z64EndOfCentralDirRecord
}

func NewDirectory() *Directory {
	d := new(Directory)

	return d
}

func GetFile(filename string, file *File) ([]byte, error) {
	f, err := os.Open(filename)
	if os.IsNotExist(err) {
		return []byte{}, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "Open File")
	}
	fi, err := os.Stat(filename)
	if err != nil {
		return nil, errors.Wrap(err, "Stat File")
	}

	pak, err := zip.NewReader(f, fi.Size())
	if err != nil {
		return nil, errors.Wrap(err, "Create Zip Reader")
	}
	entry, ok := lo.Find(pak.File, func(item *zip.File) bool {
		return item.Name == file.Filename
	})
	if !ok {
		return nil, errors.New("File not found")
	}
	rc, err := entry.Open()
	if err != nil {
		return nil, errors.Wrap(err, "Open file in archive")
	}
	buf, err := io.ReadAll(rc)
	if err != nil {
		return nil, errors.Wrap(err, "Read file in archive")
	}
	return buf, nil
}

func (d *Directory) getRecords(reader *recReader) (err error) {
	if _, err := reader.GotoLastPage(); err != nil {
		return errors.Wrap(err, "Go to Last Page")
	}
	d.endOfCentralDirRecord, err = NewEndOfCentralDirRecord(reader)
	if err != nil {
		return errors.Wrap(err, "New EndOfCentralDirRecord")
	}
	if _, err := reader.Seek(d.endOfCentralDirRecord.RecordOffset-PageSize, io.SeekStart); err != nil {
		return errors.Wrap(err, "Seek to Locator")
	}
	d.z64EndOfCentralDirLocator, err = NewZ64EndOfCentralDirLocator(reader)
	if err != nil {
		return errors.Wrap(err, "Create Z64EndOfCentralDirLocator")
	}
	if _, err := reader.Seek(d.z64EndOfCentralDirLocator.Z64EndOfCentralDir(), io.SeekStart); err != nil {
		return errors.Wrap(err, "Seek to Z64EndOfCentralDirRecord")
	}
	d.z64EndOfCentralDirRecord, err = NewZ64EndOfCentralDirRecord(reader)
	if err != nil {
		return errors.Wrap(err, "Create Z64EndOfCentralDirRecord")
	}
	if _, err := reader.Seek(d.z64EndOfCentralDirRecord.Z64StartOfCentralDir(), io.SeekStart); err != nil {
		return errors.Wrap(err, "Seek to Z64StartOfCentralDir")
	}
	return nil
}

func (d *Directory) ScanDirectoryFor(p4kFilename, filename string) (*File, error) {
	log.Trace().Str("p4kFilename", p4kFilename).Str("filename", filename).Msg("Scanning Directory")
	reader, err := newRecReader(p4kFilename)
	if err != nil {
		return nil, errors.Wrap(err, "Open Reader")
	}
	defer util.DiscardErrorOnly(reader.Close())

	if err := d.getRecords(reader); err != nil {
		return nil, errors.Wrap(err, "Get Records")
	}
	for i := int64(0); i < d.z64EndOfCentralDirRecord.NumberOfEntries(); i++ {
		de, err := NewDirectoryEntry(reader)
		if err != nil {
			return nil, errors.Wrap(err, "New Directory Entry")
		}
		if filename == "" || !strings.HasSuffix(strings.ToLower(de.Filename), strings.ToLower(filename)) {
			continue
		}
		log.Trace().Str("filename", de.Filename).Msg("File found")
		p := NewFile(de)
		return p, nil
	}
	return nil, errors.New("File not found")
}

func (d *Directory) ScanDirectoryContaining(p4kFilename, filenamePart string) ([]*File, error) {
	log.Trace().Str("p4kFilename", p4kFilename).Str("filenamePart", filenamePart).Msg("Scanning Directory")
	reader, err := newRecReader(p4kFilename)
	if err != nil {
		return nil, errors.Wrap(err, "Open Reader")
	}
	defer util.DiscardErrorOnly(reader.Close())
	var fileList []*File
	if err := d.getRecords(reader); err != nil {
		return nil, errors.Wrap(err, "Get Records")
	}

	for i := int64(0); i < d.z64EndOfCentralDirRecord.NumberOfEntries(); i++ {
		de, err := NewDirectoryEntry(reader)
		if err != nil {
			return nil, errors.Wrap(err, "New Directory Entry")
		}
		if filenamePart == "" || !strings.Contains(strings.ToLower(de.Filename), strings.ToLower(filenamePart)) {
			continue
		}
		log.Trace().Str("filename", de.Filename).Msg("File found")
		p := NewFile(de)
		fileList = append(fileList, p)
	}
	return fileList, nil
}
