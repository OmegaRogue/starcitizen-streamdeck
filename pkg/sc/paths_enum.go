// Code generated by go-enum DO NOT EDIT.
// Version:
// Revision:
// Build Date:
// Built By:

package sc

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const (
	// VersionEPTU is a Version of type EPTU.
	VersionEPTU Version = "EPTU"
	// VersionPTU is a Version of type PTU.
	VersionPTU Version = "PTU"
	// VersionLIVE is a Version of type LIVE.
	VersionLIVE Version = "LIVE"
)

var ErrInvalidVersion = errors.New("not a valid Version")

// String implements the Stringer interface.
func (x Version) String() string {
	return string(x)
}

// String implements the Stringer interface.
func (x Version) IsValid() bool {
	_, err := ParseVersion(string(x))
	return err == nil
}

var _VersionValue = map[string]Version{
	"EPTU": VersionEPTU,
	"eptu": VersionEPTU,
	"PTU":  VersionPTU,
	"ptu":  VersionPTU,
	"LIVE": VersionLIVE,
	"live": VersionLIVE,
}

// ParseVersion attempts to convert a string to a Version.
func ParseVersion(name string) (Version, error) {
	if x, ok := _VersionValue[name]; ok {
		return x, nil
	}
	// Case insensitive parse, do a separate lookup to prevent unnecessary cost of lowercasing a string if we don't need to.
	if x, ok := _VersionValue[strings.ToLower(name)]; ok {
		return x, nil
	}
	return Version(""), fmt.Errorf("%s is %w", name, ErrInvalidVersion)
}

// MarshalText implements the text marshaller method.
func (x Version) MarshalText() ([]byte, error) {
	return []byte(string(x)), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *Version) UnmarshalText(text []byte) error {
	tmp, err := ParseVersion(string(text))
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

func (x Version) MarshalZerologObject(e *zerolog.Event) {
	e.Str("Version", fmt.Sprint(x))
}