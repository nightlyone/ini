// Package ini parses *.ini files
package ini

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"unicode"
)

type Tokenizer struct {
	scanner   *bufio.Scanner
	line      []byte
	needValue bool
	pos       int64
}

func New(r io.Reader) *Tokenizer {
	return &Tokenizer{
		scanner: bufio.NewScanner(r),
	}
}

type Token interface {
	private()
	String() string
}

type Section string

func (Section) private()         {}
func (s Section) String() string { return string(s) }

type Comment string

func (Comment) private()         {}
func (c Comment) String() string { return string(c) }

func (c Comment) TrimSpace() Comment {
	return Comment(strings.TrimSpace(string(c)))
}

type Key string

func (Key) private()         {}
func (k Key) String() string { return string(k) }

func (t *Tokenizer) Token() (token Token, err error) {
	if len(t.line) == 0 {
		if t.needValue {
			t.needValue = false
			return emptyValue, nil
		}
		for t.scanner.Scan() {
			t.pos++
			t.line = t.scanner.Bytes()
			t.line = bytes.TrimLeftFunc(t.line, unicode.IsSpace)
			if len(t.line) > 0 {
				break
			}
		}
		if err = t.scanner.Err(); err != nil {
			return nil, err
		}
		if len(t.line) == 0 {
			return nil, io.EOF
		}
	}
	if t.needValue {
		// TODO(nightlyone): Parse until specialInValue and handle comments or quotes there
		token = Value(string(t.line))
		t.line = nil
		t.needValue = false
		return token, nil

	}
	switch t.line[0] {
	case ';', '#':
		token = Comment(string(t.line[1:]))
		t.line = nil
		return token, nil
	case '[':
		end := bytes.IndexByte(t.line[1:], ']')
		if end < 0 {
			t.line = nil
			return nil, errMissingSectionDelimiter

		}
		token = Section(string(t.line[1 : end+1]))
		t.line = bytes.TrimLeftFunc(t.line[end+2:], unicode.IsSpace)
		return token, nil
	case '=':
		t.line = nil
		return nil, errMissingKey

	}
	if i := bytes.IndexByte(t.line, '='); i < 0 {
		t.line = nil
		return nil, errValueWithoutKey
	} else {
		key, value := t.line[0:i], t.line[i+1:]
		key = bytes.TrimSpace(key)
		t.line = value
		t.needValue = true
		return Key(string(key)), nil
	}
}

func (d *Tokenizer) Pos() int64 { return d.pos + 1 }
