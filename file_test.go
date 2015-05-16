package ini

import (
	"reflect"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	tests := [...]struct {
		in  string
		out *File
		err error
	}{
		{
			out: &File{
				Global:  SectionMap{},
				Section: map[string]SectionMap{},
			},
		},
		{
			in: "[section]\n",
			out: &File{
				Global: SectionMap{},
				Section: map[string]SectionMap{
					"section": SectionMap{},
				},
			},
		},
		{
			in: "[section]\nkey=value",
			out: &File{
				Global: SectionMap{},
				Section: map[string]SectionMap{
					"section": SectionMap{
						"key": "value",
					},
				},
			},
		},
		{
			in: "key=value",
			out: &File{
				Global: SectionMap{
					"key": "value",
				},
				Section: map[string]SectionMap{},
			},
		},
		{
			in: "key= 'value' ",
			out: &File{
				Global: SectionMap{
					"key": "value",
				},
				Section: map[string]SectionMap{},
			},
		},
		{
			in: "key=",
			out: &File{
				Global: SectionMap{
					"key": "",
				},
				Section: map[string]SectionMap{},
			},
		},
		{
			in: "; commments only",
			out: &File{
				Global:  SectionMap{},
				Section: map[string]SectionMap{},
			},
		},
		{
			in: "[section]; section commments",
			out: &File{
				Global: SectionMap{},
				Section: map[string]SectionMap{
					"section": SectionMap{},
				},
			},
		},
		{
			in:  "[section][section]",
			err: DuplicateSectionError{Section: "section"},
		},
		{
			in:  "global_key=value\nglobal_key=value2",
			err: DuplicateKeyError{Key: "global_key"},
		},
		{
			in:  "[section]\nkey=value\nkey=value2",
			err: DuplicateKeyError{Key: "key", Section: "section"},
		},
	}
	for i, test := range tests {
		got, err := Read(strings.NewReader(test.in))
		want := test.out
		switch {
		case err != nil && test.err == nil:
			t.Errorf("%d: reader unexpected error:  %v, (Content: %+v)", i, err, got)

		case err != nil && test.err != nil && err.Error() != test.err.Error():
			t.Errorf("%d: reader unexpected error:  got %v, want %v", i, err, test.err)

		case err == nil && test.err != nil:
			t.Errorf("%d: reader expected an error: got nil (Content: %+v), want %v", i, got, test.err)
		case err == nil && !reflect.DeepEqual(got, want):
			t.Errorf("%d: reader unexpected content: got = %+v, want = %+v", i, got, want)
		case err == nil && reflect.DeepEqual(got, want):
			if testing.Verbose() {
				t.Logf("%d: reader content: got = %+v", i, got)
			}
		}
	}

}
