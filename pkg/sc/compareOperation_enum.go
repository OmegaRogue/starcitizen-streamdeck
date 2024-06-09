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
	// CompareOperationEmpty is a CompareOperation of type Empty.
	CompareOperationEmpty CompareOperation = ""
	// CompareOperationGreaterThan is a CompareOperation of type GreaterThan.
	CompareOperationGreaterThan CompareOperation = "GreaterThan"
	// CompareOperationNotEquals is a CompareOperation of type NotEquals.
	CompareOperationNotEquals CompareOperation = "NotEquals"
)

var ErrInvalidCompareOperation = errors.New("not a valid CompareOperation")

// String implements the Stringer interface.
func (x CompareOperation) String() string {
	return string(x)
}

// String implements the Stringer interface.
func (x CompareOperation) IsValid() bool {
	_, err := ParseCompareOperation(string(x))
	return err == nil
}

var _CompareOperationValue = map[string]CompareOperation{
	"":            CompareOperationEmpty,
	"GreaterThan": CompareOperationGreaterThan,
	"greaterthan": CompareOperationGreaterThan,
	"NotEquals":   CompareOperationNotEquals,
	"notequals":   CompareOperationNotEquals,
}

// ParseCompareOperation attempts to convert a string to a CompareOperation.
func ParseCompareOperation(name string) (CompareOperation, error) {
	if x, ok := _CompareOperationValue[name]; ok {
		return x, nil
	}
	// Case insensitive parse, do a separate lookup to prevent unnecessary cost of lowercasing a string if we don't need to.
	if x, ok := _CompareOperationValue[strings.ToLower(name)]; ok {
		return x, nil
	}
	return CompareOperation(""), fmt.Errorf("%s is %w", name, ErrInvalidCompareOperation)
}

// MarshalText implements the text marshaller method.
func (x CompareOperation) MarshalText() ([]byte, error) {
	return []byte(string(x)), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *CompareOperation) UnmarshalText(text []byte) error {
	tmp, err := ParseCompareOperation(string(text))
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

func (x CompareOperation) MarshalZerologObject(e *zerolog.Event) {
	e.Str("CompareOperation", fmt.Sprint(x))
}