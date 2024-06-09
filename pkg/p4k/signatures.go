package p4k

import (
	"bytes"

	. "github.com/ahmetb/go-linq/v3"
	"github.com/pkg/errors"
)

type Signature [4]byte

var (
	SignatureLocalFileHeader           = Signature{0x50, 0x4B, 0x03, 0x04}
	SignatureLocalFileHeaderCry        = Signature{0x50, 0x4B, 0x03, 0x14}
	SignatureExtraDataRecord           = Signature{0x50, 0x4B, 0x06, 0x08}
	SignatureCentralDirRecord          = Signature{0x50, 0x4B, 0x01, 0x02}
	SignatureDigitalSignature          = Signature{0x50, 0x4B, 0x05, 0x05}
	SignatureZ64EndOfCentralDirRec     = Signature{0x50, 0x4B, 0x06, 0x06}
	SignatureZ64EndOfCentralDirLocator = Signature{0x50, 0x4B, 0x06, 0x07}
	SignatureEndOfCentralDirRecord     = Signature{0x50, 0x4B, 0x05, 0x06}
)

func FindSignatureInPage(reader *recReader, signature Signature) (int64, error) {
	pos, err := reader.Pos()
	if err != nil {
		return -1, errors.Wrap(err, "Get current pos")
	}
	lPage, err := reader.GetPage()
	if err != nil {
		return -1, errors.Wrap(err, "Get current Page")
	}
	for i := 0; i < len(lPage)-4; i++ {
		if From(lPage).Skip(i).Take(4).SequenceEqual(From(signature)) {
			return pos + int64(i), nil
		}
	}
	return -1, errors.New("Signature not found")
}

func FindSignatureInPageBackwards(reader *recReader, signature Signature) (int64, error) {
	pos, err := reader.Pos()
	if err != nil {
		return -1, errors.Wrap(err, "Get current pos")
	}
	lPage, err := reader.GetPage()
	if err != nil {
		return -1, errors.Wrap(err, "Get current Page")
	}

	if i := bytes.LastIndex(lPage, signature[0:4]); i >= 0 {
		return pos + int64(i), nil
	}
	return -1, errors.New("Signature not found")
}
