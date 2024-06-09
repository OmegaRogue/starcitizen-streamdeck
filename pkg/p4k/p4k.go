package p4k

import (
	"encoding/binary"
	"reflect"
	"time"

	"github.com/rs/zerolog/log"
)

const PageSize int64 = 0x1000

const (
	DirectoryEntryRecordLength         = 46
	Z64DirectoryEntryExtraRecordLength = 32
	EndOfCentralDirRecordLength        = 22
	Z64EndOfCentralDirLocatorLength    = 20
	Z64EndOfCentralDirRecordLength     = 56
	FileHeaderRecordLength             = 30
	Z64FileHeaderExtraRecordLength     = 32
)

func init() {
	if binary.Size(MyZ64EndOfCentralDirRecord{}) != Z64EndOfCentralDirRecordLength {
		log.Panic().Str("name", reflect.TypeOf(MyZ64EndOfCentralDirRecord{}).Name()).Uint64("recordSize", uint64(binary.Size(MyZ64EndOfCentralDirRecord{}))).Msg("Record size does not match!")
	}
	if binary.Size(MyZ64EndOfCentralDirLocator{}) != Z64EndOfCentralDirLocatorLength {
		log.Panic().Str("name", reflect.TypeOf(MyZ64EndOfCentralDirLocator{}).Name()).Uint64("recordSize", uint64(binary.Size(MyZ64EndOfCentralDirLocator{}))).Msg("Record size does not match!")
	}
	if binary.Size(MyEndOfCentralDirRecord{}) != EndOfCentralDirRecordLength {
		log.Panic().Str("name", reflect.TypeOf(MyEndOfCentralDirRecord{}).Name()).Uint64("recordSize", uint64(binary.Size(MyEndOfCentralDirRecord{}))).Msg("Record size does not match!")
	}
	if binary.Size(MyDirectoryEntryRecord{}) != DirectoryEntryRecordLength {
		log.Panic().Str("name", reflect.TypeOf(MyDirectoryEntryRecord{}).Name()).Uint64("recordSize", uint64(binary.Size(MyDirectoryEntryRecord{}))).Msg("Record size does not match!")
	}
	if binary.Size(MyZ64DirectoryEntryExtraRecord{}) != Z64DirectoryEntryExtraRecordLength {
		log.Panic().Str("name", reflect.TypeOf(MyZ64DirectoryEntryExtraRecord{}).Name()).Uint64("recordSize", uint64(binary.Size(MyZ64DirectoryEntryExtraRecord{}))).Msg("Record size does not match!")
	}
	if binary.Size(MyFileHeaderRecord{}) != FileHeaderRecordLength {
		log.Panic().Str("name", reflect.TypeOf(MyFileHeaderRecord{}).Name()).Uint64("recordSize", uint64(binary.Size(MyFileHeaderRecord{}))).Msg("Record size does not match!")
	}
	if binary.Size(MyZ64FileHeaderExtraRecord{}) != Z64FileHeaderExtraRecordLength {
		log.Panic().Str("name", reflect.TypeOf(MyZ64FileHeaderExtraRecord{}).Name()).Uint64("recordSize", uint64(binary.Size(MyZ64FileHeaderExtraRecord{}))).Msg("Record size does not match!")
	}
}

func TimestampFromDos(date, t uint16) time.Time {
	dateInt := int(date)
	timeInt := int(t)
	year := ((dateInt >> 9) & 0x7f) + 1980
	month := (dateInt >> 5) & 0x0f
	day := dateInt & 0x01f
	hour := (timeInt >> 11) & 0x1f
	minute := (timeInt >> 5) & 0x3f
	sec := (timeInt & 0x01f) * 2
	return time.Date(year, time.Month(month), day, hour, minute, sec, 0, time.UTC)
}
