package ini

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
)

type SectionMap map[string]string

func (s SectionMap) Has(key string) bool {
	if s == nil {
		return false
	}
	_, ok := s[key]
	return ok

}

type File struct {
	Global  SectionMap
	Section map[string]SectionMap
}

func ReadFile(filename string) (*File, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return Read(bytes.NewReader(b))
}

type DuplicateSectionError struct {
	Section string
}

func (d DuplicateSectionError) Error() string {
	return fmt.Sprintf("ini: duplicate section %q", d.Section)
}

type DuplicateKeyError struct {
	Section, Key string
}

func (m DuplicateKeyError) Error() string {
	if m.Section == "" {
		return fmt.Sprintf("ini: duplicate global key %q", m.Key)
	} else {
		return fmt.Sprintf("ini: duplicate key %q in section %q", m.Key, m.Section)
	}
}

func Read(r io.Reader) (*File, error) {
	section := ""
	key := ""
	hasSections := false
	needValue := false

	f := &File{
		Global:  SectionMap{},
		Section: map[string]SectionMap{},
	}

	d := New(r)

	var err error
	var t Token
	for t, err = d.Token(); err == nil; t, err = d.Token() {

		switch token := t.(type) {
		case Value:
			value, err := token.Unquote()
			if err != nil {
				return nil, err
			}
			if hasSections {
				f.Section[section][key] = value.String()
			} else {
				f.Global[key] = value.String()
			}
			needValue = false
		case Key:
			key = token.String()
			if hasSections {
				if f.Section[section].Has(key) {
					return nil, DuplicateKeyError{Key: key, Section: section}
				}
			} else {
				if f.Global.Has(key) {
					return nil, DuplicateKeyError{Key: key}
				}
			}
			needValue = true
		case Section:
			section = token.String()
			if _, ok := f.Section[section]; !ok {
				f.Section[section] = make(SectionMap)
				hasSections = true
			} else {
				return nil, DuplicateSectionError{Section: section}
			}
		case Comment:
			// comments are ignored
		default:
			panic("ini: unknown token, update me!")
		}

	}
	if err == io.EOF {
		if needValue {
			panic("ini: we should at least always get an empty value")
		}
		err = nil
	}
	if err != nil {
		return nil, err
	}

	return f, nil
}
