package ini

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
)

type TokenList []Token

func (toks TokenList) String() string {
	var b bytes.Buffer
	if len(toks) == 0 {
		return "[]"
	}
	b.WriteByte('[')
	fmt.Fprintf(&b, "%T:%q", toks[0], toks[0])
	for _, tok := range toks[1:] {
		fmt.Fprintf(&b, ",%T:%q", tok, tok)
	}
	b.WriteByte(']')
	return b.String()
}

func tokenizeString(s string) (tokens TokenList, err error) {
	d := New(strings.NewReader(s))
	var t Token
	for t, err = d.Token(); err == nil; t, err = d.Token() {

		tokens = append(tokens, t)
	}
	return tokens, err
}

func TestTokenizer(t *testing.T) {
	tests := [...]struct {
		in  string
		out TokenList
		err error
	}{
		{
			err: io.EOF,
		},
		{
			err: io.EOF,
			in:  "key=value\n",
			out: TokenList{
				Key("key"), Value("value"),
			},
		},
		{
			err: io.EOF,
			in:  "key=value",
			out: TokenList{
				Key("key"), Value("value"),
			},
		},
		{
			err: io.EOF,
			in:  "    key=value",
			out: TokenList{
				Key("key"), Value("value"),
			},
		},
		{
			err: io.EOF,
			in:  "key= value",
			out: TokenList{
				Key("key"), Value(" value"),
			},
		},
		{
			err: io.EOF,
			in:  "key= value ",
			out: TokenList{
				Key("key"), Value(" value "),
			},
		},
		{
			err: io.EOF,
			in:  "key=",
			out: TokenList{
				Key("key"), Value(""),
			},
		},
		{
			err: errMissingKey,
			in:  "=",
			out: TokenList{},
		},
		{
			err: io.EOF,
			in:  "[section]\n",
			out: TokenList{
				Section("section"),
			},
		},
		{
			err: io.EOF,
			in:  "[section]",
			out: TokenList{
				Section("section"),
			},
		},
		{
			err: io.EOF,
			in:  ";comment\n",
			out: TokenList{
				Comment("comment"),
			},
		},
		{
			err: io.EOF,
			in:  "#comment\n",
			out: TokenList{
				Comment("comment"),
			},
		},
		{
			err: io.EOF,
			in:  "[section]\nkey=value",
			out: TokenList{
				Section("section"), Key("key"), Value("value"),
			},
		},
		{ // we are a little bit more flexible here
			err: io.EOF,
			in:  "[section]key=value",
			out: TokenList{
				Section("section"), Key("key"), Value("value"),
			},
		},
		{
			err: errMissingSectionDelimiter,
			in:  "[section",
			out: TokenList{},
		},
		{
			err: errValueWithoutKey,
			in:  "key=\nvalue",
			out: TokenList{Key("key")},
		},
		{
			err: errValueWithoutKey,
			in:  "[section]\nkey=\nvalue",
			out: TokenList{Section("section"), Key("key")},
		},
	}

	for i, test := range tests {
		got, err := tokenizeString(test.in)
		want := test.out
		switch {
		case err != nil && test.err == nil:
			t.Errorf("%d: tokenizer unexpected error:  %v, (Tokens: %s)", i, err, got)

		case err != nil && test.err != nil && err.Error() != test.err.Error():
			t.Errorf("%d: tokenizer unexpected error:  got %v, want %v", i, err, test.err)

		case err == nil && test.err != nil:
			t.Errorf("%d: tokenizer expected an error: got nil (Tokens: %s), want %v", i, got, test.err)
		case err == nil && got.String() != want.String():
			t.Errorf("%d: tokenizer unexpected tokens: got = %s, want = %s", i, got, want)
		case err == nil && got.String() == want.String():
			if testing.Verbose() {
				t.Logf("%d: tokenizer tokens: got = %s, want = %s", i, got, want)
			}
		case err == io.EOF && got.String() != want.String():
			t.Errorf("%d: tokenizer unexpected tokens: got = %s, want = %s", i, got, want)
		case err == io.EOF && got.String() == want.String():
			if testing.Verbose() {
				t.Logf("%d: tokenizer tokens: got = %s", i, got)
			}
		}
	}
}
