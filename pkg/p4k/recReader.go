package p4k

import (
	"encoding/binary"
	"io"
	"os"

	"github.com/pkg/errors"
)

type recReader struct {
	fileStream *os.File
}

func newRecReader(filename string) (*recReader, error) {
	r := new(recReader)
	err := r.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, "open file")
	}
	return r, nil
}

func (r *recReader) Seek(offset int64, whence int) (int64, error) {
	n, err := r.fileStream.Seek(offset, whence)
	if err != nil {
		return n, errors.Wrap(err, "Seek")
	}
	return n, nil
}

func (r *recReader) Pos() (int64, error) {
	pos, err := r.fileStream.Seek(0, io.SeekCurrent)
	if err != nil {
		return -1, errors.Wrap(err, "Seek pos of record reader")
	}
	return pos, nil
}

func (r *recReader) Length() (int64, error) {
	fi, err := r.fileStream.Stat()
	if err != nil {
		return -1, errors.Wrap(err, "Get Length of file stream")
	}
	return fi.Size(), nil
}

func (r *recReader) Close() error {
	err := r.fileStream.Close()
	if err != nil {
		return errors.Wrap(err, "Close Reader")
	}
	return nil
}

func (r *recReader) Read(p []byte) (int, error) {
	n, err := r.fileStream.Read(p)
	if err != nil {
		return n, errors.Wrap(err, "Read")
	}
	return n, nil
}

func (r *recReader) Open(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return errors.Wrap(err, "Open file")
	}
	r.fileStream = f
	return nil
}

func (r *recReader) IsOpen() bool {
	return r.fileStream != nil
}

func ByteToType[T comparable](reader io.Reader) (*T, error) {
	buf := new(T)
	err := binary.Read(reader, binary.LittleEndian, buf)
	if err != nil {
		return nil, errors.Wrap(err, "Try to Unmarshal binary type")
	}
	return buf, nil
}

func (r *recReader) AdvancePage() error {
	current, err := r.Pos()
	if err != nil {
		return errors.Wrap(err, "Get current pos")
	}
	remainder := current % PageSize
	if remainder <= 0 {
		return nil
	}
	seek := PageSize - remainder
	_, err = r.Seek(seek, io.SeekStart)
	if err != nil {
		return errors.Wrap(err, "Seek to next Page")
	}
	return nil
}
func (r *recReader) BackwardPage() error {
	_, err := r.Seek(-PageSize, io.SeekCurrent)
	if err != nil {
		return errors.Wrap(err, "Seek to previous Page")
	}
	return nil
}
func (r *recReader) GotoLastPage() (int64, error) {
	pos, err := r.Seek(-PageSize, io.SeekEnd)
	if err != nil {
		return -1, errors.Wrap(err, "Seek to last Page")
	}
	return pos, nil
}

func (r *recReader) GetPage() ([]byte, error) {
	buf := make([]byte, PageSize)
	_, err := r.Read(buf)
	if err != nil {
		return nil, errors.Wrap(err, "Read Page")
	}
	return buf, nil
}
