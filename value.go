package ini

import (
	"errors"
	"strings"
)

type Value string

func (Value) private()         {}
func (v Value) String() string { return string(v) }

func (v Value) TrimSpace() Value {
	return Value(strings.TrimSpace(string(v)))
}

const (
	quoteSingle = '\''
	quoteDouble = '"'
)

var errUnmatchedSingleQuote = errors.New("unmatched single quote")
var errUnmatchedDoubleQuote = errors.New("unmatched double quote")
var errMissingSectionDelimiter = errors.New("missing section delimiter, a closing \"]\"")
var errMissingKey = errors.New("missing key, found bare \"=\" instead")
var errValueWithoutKey = errors.New("found bare value without key")

var emptyValue = Value("")

func (v Value) Unquote() (Value, error) {
	s := v.TrimSpace()
	switch {
	case len(s) == 0:
		return v, nil
	case s[0] == quoteSingle:
		if i := strings.IndexByte(string(s)[1:], quoteSingle); i < 0 || i != len(s)-2 {
			return emptyValue, errUnmatchedSingleQuote
		} else {
			return s[1 : i+1], nil
		}
	case s[0] == quoteDouble:
		if i := strings.IndexByte(string(s)[1:], quoteDouble); i < 0 || i != len(s)-2 {
			return emptyValue, errUnmatchedDoubleQuote
		} else {
			return s[1 : i+1], nil
		}
	default:
		return v, nil
	}
}

var specialInValue = []byte(`#;"'`)
