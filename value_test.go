package ini

import "testing"

func TestUnquote(t *testing.T) {

	tests := [...]struct {
		in, out Value
		err     error
	}{
		{},
		{in: ".", out: "."},
		{in: "   .", out: "   ."},
		{in: "..", out: ".."},
		{in: ".'", out: ".'"},
		{in: ".\"", out: ".\""},
		{in: `'`, err: errUnmatchedSingleQuote},
		{in: `"`, err: errUnmatchedDoubleQuote},
		{in: `'"`, err: errUnmatchedSingleQuote},
		{in: `"'`, err: errUnmatchedDoubleQuote},
		{in: `""`},
		{in: `''`},
		{in: `"h"`, out: "h"},
		{in: ` " h " `, out: " h "},
		{in: `'h'`, out: "h"},
	}

	for i, test := range tests {
		got, err := test.in.Unquote()
		want := test.out
		switch {
		case err != nil && test.err == nil:
			t.Errorf("%d: Unquote(%v) unexpected error, want nil, got %v\n", i, test.in, err)
		case err != nil && test.err.Error() != err.Error():
			t.Errorf("%d: Unquote(%v) unexpected error, want '%v', got %v\n", i, test.in, test.err, err)
		case err == nil && test.err != nil:
			t.Errorf("%d: Unquote(%v) expected error, want '%v', got nil\n", i, test.in, test.err)
		case got != want:
			t.Errorf("%d: unexpected value: Unquote(%v) = %q, want %q\n", i, test.in, got, want)
		}
	}
}
