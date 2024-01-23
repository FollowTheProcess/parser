package parser_test

import (
	"fmt"
	"os"
	"testing"
	"unicode"

	"github.com/FollowTheProcess/parser"
)

func TestTake(t *testing.T) {
	tests := []struct {
		name    string        // Identifying test case name
		want    parser.Result // The expected result of parsing
		input   string        // Entire input to be parsed
		err     string        // The expected error message (if there is one)
		n       int           // The number of chars to consume
		wantErr bool          // Whether it should have returned an error
	}{
		{
			name:    "empty input",
			want:    parser.Result{},
			input:   "",
			n:       999,
			wantErr: true,
			err:     "Take: cannot take from empty input",
		},
		{
			name:    "empty input with n zero",
			want:    parser.Result{},
			input:   "",
			n:       0,
			wantErr: true,
			err:     "Take: n must be a non-zero positive integer, got 0",
		},
		{
			name:    "n too large",
			want:    parser.Result{},
			input:   "some stuff here",
			n:       999,
			wantErr: true,
			err:     "Take: requested n (999) chars but input had only 15 utf-8 chars",
		},
		{
			name:    "n negative",
			want:    parser.Result{},
			input:   "some stuff here",
			n:       -1,
			wantErr: true,
			err:     "Take: n must be a non-zero positive integer, got -1",
		},
		{
			name:    "n zero",
			want:    parser.Result{},
			input:   "some stuff here",
			n:       0,
			wantErr: true,
			err:     "Take: n must be a non-zero positive integer, got 0",
		},
		{
			name:    "n 1 more than len",
			want:    parser.Result{},
			input:   "This is an exact length",
			n:       24,
			wantErr: true,
			err:     "Take: requested n (24) chars but input had only 23 utf-8 chars",
		},
		{
			name:    "n 1 more than len utf8",
			want:    parser.Result{},
			input:   "日a本b語ç日ð本Ê語þ日¥本¼語i日©", // This is 20 utf-8 runes
			n:       21,
			wantErr: true,
			err:     "Take: requested n (21) chars but input had only 20 utf-8 chars",
		},
		{
			name:    "bad utf8",
			want:    parser.Result{},
			input:   "\xf8\xa1\xa1\xa1\xa1",
			n:       3,
			wantErr: true,
			err:     "Take: input not valid utf-8",
		},
		{
			name: "simple",
			want: parser.Result{
				Value:     "Hello I am",
				Remainder: " some input",
			},
			input:   "Hello I am some input",
			n:       10,
			wantErr: false,
			err:     "",
		},
		{
			name: "n same as len",
			want: parser.Result{
				Value:     "This is an exact length",
				Remainder: "",
			},
			input:   "This is an exact length",
			n:       23,
			wantErr: false,
			err:     "",
		},
		{
			name: "n same as len utf8",
			want: parser.Result{
				Value:     "日a本b語ç日ð本Ê語þ日¥本¼語i日©",
				Remainder: "",
			},
			input:   "日a本b語ç日ð本Ê語þ日¥本¼語i日©", // This is 20 utf-8 runes
			n:       20,
			wantErr: false,
			err:     "",
		},
		{
			name: "n 1 less than len",
			want: parser.Result{
				Value:     "This is an exact lengt",
				Remainder: "h",
			},
			input:   "This is an exact length",
			n:       22,
			wantErr: false,
			err:     "",
		},
		{
			name: "n 1 less than len utf8",
			want: parser.Result{
				Value:     "日a本b語ç日ð本Ê語þ日¥本¼語i日",
				Remainder: "©",
			},
			input:   "日a本b語ç日ð本Ê語þ日¥本¼語i日©", // This is 20 utf-8 runes
			n:       19,
			wantErr: false,
			err:     "",
		},
		{
			name: "non ascii",
			want: parser.Result{
				Value:     "Hello",
				Remainder: ", 世界",
			},
			input:   "Hello, 世界",
			n:       5,
			wantErr: false,
			err:     "",
		},
		{
			// https://github.com/golang/exp/blob/master/utf8string/string_test.go
			name: "utf8string test",
			want: parser.Result{
				Value:     "日a本b語ç",
				Remainder: "日ð本Ê語þ日¥本¼語i日©",
			},
			input:   "日a本b語ç日ð本Ê語þ日¥本¼語i日©",
			n:       6,
			wantErr: false,
			err:     "",
		},
		{
			name: "emoji",
			want: parser.Result{
				Value:     "😱 ",
				Remainder: "emoji works too",
			},
			input:   "😱 emoji works too",
			n:       2,
			wantErr: false,
			err:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.Take(tt.n)(tt.input)

			// Should only error if we wanted one
			if (err != nil) != tt.wantErr {
				t.Fatalf("\nGot error:\t%v\nWanted error:\t%v\n", err, tt.wantErr)
			}

			// If we did get an error, the message should match what we expect
			if err != nil {
				if msg := err.Error(); msg != tt.err {
					t.Fatalf("\nGot:\t%q\nWanted:\t%q\n", msg, tt.err)
				}
			}

			// The parser result should be what we expected
			if got != tt.want {
				t.Errorf("\nGot:\t%#v\nWanted:\t%#v\n", got, tt.want)
			}
		})
	}
}

func TestExact(t *testing.T) {
	tests := []struct {
		name    string        // Identifying test case name
		want    parser.Result // The expected result of parsing
		input   string        // Entire input to be parsed
		err     string        // The expected error message (if there is one)
		match   string        // The exact string to parser
		wantErr bool          // Whether it should have returned an error
	}{
		{
			name:    "empty input",
			want:    parser.Result{},
			input:   "",
			match:   "something",
			wantErr: true,
			err:     "Exact: cannot match on empty input",
		},
		{
			name:    "bad utf8",
			want:    parser.Result{},
			input:   "\xf8\xa1\xa1\xa1\xa1",
			match:   "something",
			wantErr: true,
			err:     "Exact: input not valid utf-8",
		},
		{
			name:    "empty input and match",
			want:    parser.Result{},
			input:   "",
			match:   "",
			wantErr: true,
			err:     "Exact: cannot match on empty input",
		},
		{
			name:    "empty match",
			want:    parser.Result{},
			input:   "some text",
			match:   "",
			wantErr: true,
			err:     "Exact: match must not be empty",
		},
		{
			name:    "match longer than input",
			want:    parser.Result{},
			input:   "A single sentence",
			match:   "A single sentence but this one is longer so it can't possibly be matched",
			wantErr: true,
			err:     "Exact: match (A single sentence but this one is longer so it can't possibly be matched) not in input",
		},
		{
			name:    "match not found",
			want:    parser.Result{},
			input:   "Nothing to see in here",
			match:   "Found me",
			wantErr: true,
			err:     "Exact: match (Found me) not in input",
		},
		{
			name:    "wrong case",
			want:    parser.Result{},
			input:   "Found me, in a larger sentence",
			match:   "found me", // Note: lower case 'f', not an exact match
			wantErr: true,
			err:     "Exact: match (found me) not in input",
		},
		{
			name: "simple match",
			want: parser.Result{
				Value:     "Found me",
				Remainder: ", in a larger sentence",
			},
			input:   "Found me, in a larger sentence",
			match:   "Found me",
			wantErr: false,
			err:     "",
		},

		{
			name: "utf8 match",
			want: parser.Result{
				Value:     "世",
				Remainder: "界, Hello",
			},
			input:   "世界, Hello",
			match:   "世",
			wantErr: false,
			err:     "",
		},
		{
			// https://github.com/golang/exp/blob/master/utf8string/string_test.go
			name: "utf8string test",
			want: parser.Result{
				Value:     "日a本",
				Remainder: "b語ç日ð本Ê語þ日¥本¼語i日©",
			},
			input:   "日a本b語ç日ð本Ê語þ日¥本¼語i日©",
			match:   "日a本",
			wantErr: false,
			err:     "",
		},
		{
			name: "emoji",
			want: parser.Result{
				Value:     "😱 emoji",
				Remainder: " works too",
			},
			input:   "😱 emoji works too",
			match:   "😱 emoji",
			wantErr: false,
			err:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.Exact(tt.match)(tt.input)

			// Should only error if we wanted one
			if (err != nil) != tt.wantErr {
				t.Fatalf("\nGot error:\t%v\nWanted error:\t%v\n", err, tt.wantErr)
			}

			// If we did get an error, the message should match what we expect
			if err != nil {
				if msg := err.Error(); msg != tt.err {
					t.Fatalf("\nGot:\t%q\nWanted:\t%q\n", msg, tt.err)
				}
			}

			// The parser result should be what we expected
			if got != tt.want {
				t.Errorf("\nGot:\t%#v\nWanted:\t%#v\n", got, tt.want)
			}
		})
	}
}

func TestExactCaseInsensitive(t *testing.T) {
	tests := []struct {
		name    string        // Identifying test case name
		want    parser.Result // The expected result of parsing
		input   string        // Entire input to be parsed
		err     string        // The expected error message (if there is one)
		match   string        // The exact string to parser
		wantErr bool          // Whether it should have returned an error
	}{
		{
			name:    "empty input",
			want:    parser.Result{},
			input:   "",
			match:   "something",
			wantErr: true,
			err:     "ExactCaseInsensitive: cannot match on empty input",
		},
		{
			name:    "bad utf8",
			want:    parser.Result{},
			input:   "\xf8\xa1\xa1\xa1\xa1",
			match:   "something",
			wantErr: true,
			err:     "ExactCaseInsensitive: input not valid utf-8",
		},
		{
			name:    "empty input and match",
			want:    parser.Result{},
			input:   "",
			match:   "",
			wantErr: true,
			err:     "ExactCaseInsensitive: cannot match on empty input",
		},
		{
			name:    "empty match",
			want:    parser.Result{},
			input:   "some text",
			match:   "",
			wantErr: true,
			err:     "ExactCaseInsensitive: match must not be empty",
		},
		{
			name:    "match longer than input",
			want:    parser.Result{},
			input:   "A single sentence",
			match:   "A single sentence but this one is longer so it can't possibly be matched",
			wantErr: true,
			err:     "ExactCaseInsensitive: match (A single sentence but this one is longer so it can't possibly be matched) not in input",
		},
		{
			name:    "match not found",
			want:    parser.Result{},
			input:   "Nothing to see in here",
			match:   "Found me",
			wantErr: true,
			err:     "ExactCaseInsensitive: match (Found me) not in input",
		},
		{
			name: "exact match",
			want: parser.Result{
				Value:     "Found me",
				Remainder: ", in a larger sentence",
			},
			input:   "Found me, in a larger sentence",
			match:   "Found me",
			wantErr: false,
			err:     "",
		},
		{
			name: "case insensitive match",
			want: parser.Result{
				Value:     "Found me",
				Remainder: ", in a larger sentence",
			},
			input:   "Found me, in a larger sentence",
			match:   "found me", // Lower case f, should still match
			wantErr: false,
			err:     "",
		},
		{
			// https://github.com/golang/exp/blob/master/utf8string/string_test.go
			name: "utf8string test",
			want: parser.Result{
				Value:     "日a本",
				Remainder: "b語ç日ð本Ê語þ日¥本¼語i日©",
			},
			input:   "日a本b語ç日ð本Ê語þ日¥本¼語i日©",
			match:   "日a本", // Apparently this is already lower case
			wantErr: false,
			err:     "",
		},
		{
			name: "utf8string test upper case",
			want: parser.Result{
				Value:     "日a本",
				Remainder: "b語ç日ð本Ê語þ日¥本¼語i日©",
			},
			input:   "日a本b語ç日ð本Ê語þ日¥本¼語i日©",
			match:   "日A本", // Upper case now
			wantErr: false,
			err:     "",
		},
		{
			name: "emoji",
			want: parser.Result{
				Value:     "😱 EMOJI",
				Remainder: " WORKS TOO",
			},
			input:   "😱 EMOJI WORKS TOO",
			match:   "😱 emoji",
			wantErr: false,
			err:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.ExactCaseInsensitive(tt.match)(tt.input)

			// Should only error if we wanted one
			if (err != nil) != tt.wantErr {
				t.Fatalf("\nGot error:\t%v\nWanted error:\t%v\n", err, tt.wantErr)
			}

			// If we did get an error, the message should match what we expect
			if err != nil {
				if msg := err.Error(); msg != tt.err {
					t.Fatalf("\nGot:\t%q\nWanted:\t%q\n", msg, tt.err)
				}
			}

			// The parser result should be what we expected
			if got != tt.want {
				t.Errorf("\nGot:\t%#v\nWanted:\t%#v\n", got, tt.want)
			}
		})
	}
}

func TestChar(t *testing.T) {
	tests := []struct {
		name    string        // Identifying test case name
		want    parser.Result // The expected result of parsing
		input   string        // Entire input to be parsed
		err     string        // The expected error message (if there is one)
		char    rune          // The exact char to match
		wantErr bool          // Whether it should have returned an error
	}{
		{
			name:    "empty input",
			want:    parser.Result{},
			input:   "",
			char:    0,
			wantErr: true,
			err:     "Char: input text is empty",
		},
		{
			name:    "bad utf8",
			want:    parser.Result{},
			input:   "\xf8\xa1\xa1\xa1\xa1",
			char:    'x',
			wantErr: true,
			err:     "Char: input not valid utf-8",
		},
		{
			name:    "not found",
			want:    parser.Result{},
			input:   "something",
			char:    'x',
			wantErr: true,
			err:     "Char: requested char (x) not found in input",
		},
		{
			name:    "wrong case",
			want:    parser.Result{},
			input:   "General Kenobi!",
			char:    'g',
			wantErr: true,
			err:     "Char: requested char (g) not found in input",
		},
		{
			name: "found",
			want: parser.Result{
				Value:     "G",
				Remainder: "eneral Kenobi!",
			},
			input:   "General Kenobi!",
			char:    'G',
			wantErr: false,
			err:     "",
		},
		{
			name: "found utf8",
			want: parser.Result{
				Value:     "日",
				Remainder: "a本b語ç日ð本Ê語þ日¥本¼語i日©",
			},
			input:   "日a本b語ç日ð本Ê語þ日¥本¼語i日©",
			char:    '日',
			wantErr: false,
			err:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.Char(tt.char)(tt.input)

			// Should only error if we wanted one
			if (err != nil) != tt.wantErr {
				t.Fatalf("\nGot error:\t%v\nWanted error:\t%v\n", err, tt.wantErr)
			}

			// If we did get an error, the message should match what we expect
			if err != nil {
				if msg := err.Error(); msg != tt.err {
					t.Fatalf("\nGot:\t%q\nWanted:\t%q\n", msg, tt.err)
				}
			}

			// The parser result should be what we expected
			if got != tt.want {
				t.Errorf("\nGot:\t%#v\nWanted:\t%#v\n", got, tt.want)
			}
		})
	}
}

func TestTakeWhile(t *testing.T) {
	tests := []struct {
		predicate func(r rune) bool // The predicate function that determines whether the parser should continue taking characters
		name      string            // Identifying test case name
		want      parser.Result     // The expected result of parsing
		input     string            // Entire input to be parsed
		err       string            // The expected error message (if there is one)
		wantErr   bool              // Whether it should have returned an error
	}{
		{
			name:      "empty input",
			want:      parser.Result{},
			input:     "",
			predicate: nil, // Shouldn't matter as it should never get called
			wantErr:   true,
			err:       "TakeWhile: input text is empty",
		},
		{
			name:      "bad utf8",
			want:      parser.Result{},
			input:     "\xf8\xa1\xa1\xa1\xa1",
			predicate: nil, // Shouldn't matter as it should never get called
			wantErr:   true,
			err:       "TakeWhile: input not valid utf-8",
		},
		{
			name:      "nil predicate", // Good libraries don't panic
			want:      parser.Result{},
			input:     "some input",
			predicate: nil,
			wantErr:   true,
			err:       "TakeWhile: predicate must be a non-nil function",
		},
		{
			name:      "predicate never returns false", // Good libraries don't allow infinite loops
			want:      parser.Result{},
			input:     "fixed length input",
			predicate: func(r rune) bool { return true },
			wantErr:   true,
			err:       "TakeWhile: predicate never returned false",
		},
		{
			name: "predicate never returns true",
			want: parser.Result{
				Value:     "",
				Remainder: "fixed length input",
			},
			input:     "fixed length input",
			predicate: func(r rune) bool { return false },
			wantErr:   false,
			err:       "",
		},
		{
			name: "consume whitespace",
			want: parser.Result{
				Value:     "  \t\t\n\n ",
				Remainder: "end of whitespace",
			},
			input:     "  \t\t\n\n end of whitespace",
			predicate: unicode.IsSpace,
			wantErr:   false,
			err:       "",
		},
		{
			name: "consume non ascii rune",
			want: parser.Result{
				Value:     "本本本",
				Remainder: " b語ç日ð本Ê語",
			},
			input:     "本本本 b語ç日ð本Ê語",
			predicate: func(r rune) bool { return r == '本' },
			wantErr:   false,
			err:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.TakeWhile(tt.predicate)(tt.input)

			// Should only error if we wanted one
			if (err != nil) != tt.wantErr {
				t.Fatalf("\nGot error:\t%v\nWanted error:\t%v\n", err, tt.wantErr)
			}

			// If we did get an error, the message should match what we expect
			if err != nil {
				if msg := err.Error(); msg != tt.err {
					t.Fatalf("\nGot:\t%q\nWanted:\t%q\n", msg, tt.err)
				}
			}

			// The parser result should be what we expected
			if got != tt.want {
				t.Errorf("\nGot:\t%#v\nWanted:\t%#v\n", got, tt.want)
			}
		})
	}
}

func TestTakeUntil(t *testing.T) {
	tests := []struct {
		predicate func(r rune) bool // The predicate function that determines whether the parser should stop taking characters
		name      string            // Identifying test case name
		want      parser.Result     // The expected result of parsing
		input     string            // Entire input to be parsed
		err       string            // The expected error message (if there is one)
		wantErr   bool              // Whether it should have returned an error
	}{
		{
			name:      "empty input",
			want:      parser.Result{},
			input:     "",
			predicate: nil, // Shouldn't matter as it should never get called
			wantErr:   true,
			err:       "TakeUntil: input text is empty",
		},
		{
			name:      "bad utf8",
			want:      parser.Result{},
			input:     "\xf8\xa1\xa1\xa1\xa1",
			predicate: nil, // Shouldn't matter as it should never get called
			wantErr:   true,
			err:       "TakeUntil: input not valid utf-8",
		},
		{
			name:      "nil predicate", // Good libraries don't panic
			want:      parser.Result{},
			input:     "some input",
			predicate: nil,
			wantErr:   true,
			err:       "TakeUntil: predicate must be a non-nil function",
		},
		{
			name:      "predicate never returns true", // Good libraries don't allow infinite loops
			want:      parser.Result{},
			input:     "fixed length input",
			predicate: func(r rune) bool { return false },
			wantErr:   true,
			err:       "TakeUntil: predicate never returned true",
		},
		{
			name: "predicate never returns false",
			want: parser.Result{
				Value:     "",
				Remainder: "fixed length input",
			},
			input:     "fixed length input",
			predicate: func(r rune) bool { return true },
			wantErr:   false,
			err:       "",
		},
		{
			name: "consume until whitespace",
			want: parser.Result{
				Value:     "something",
				Remainder: " <- first whitespace",
			},
			input:     "something <- first whitespace",
			predicate: unicode.IsSpace,
			wantErr:   false,
			err:       "",
		},
		{
			name: "consume until non-ascii",
			want: parser.Result{
				Value:     "abcdef",
				Remainder: "語ç日ð本Ê語",
			},
			input:     "abcdef語ç日ð本Ê語",
			predicate: func(r rune) bool { return r > unicode.MaxASCII },
			wantErr:   false,
			err:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.TakeUntil(tt.predicate)(tt.input)

			// Should only error if we wanted one
			if (err != nil) != tt.wantErr {
				t.Fatalf("\nGot error:\t%v\nWanted error:\t%v\n", err, tt.wantErr)
			}

			// If we did get an error, the message should match what we expect
			if err != nil {
				if msg := err.Error(); msg != tt.err {
					t.Fatalf("\nGot:\t%q\nWanted:\t%q\n", msg, tt.err)
				}
			}

			// The parser result should be what we expected
			if got != tt.want {
				t.Errorf("\nGot:\t%#v\nWanted:\t%#v\n", got, tt.want)
			}
		})
	}
}

func TestOneOf(t *testing.T) {
	tests := []struct {
		name    string        // Identifying test case name
		want    parser.Result // The expected result of parsing
		input   string        // Entire input to be parsed
		chars   string        // The chars to match one of
		err     string        // The expected error message (if there is one)
		wantErr bool          // Whether it should have returned an error
	}{
		{
			name:    "empty input",
			want:    parser.Result{},
			input:   "",
			chars:   "abc", // Doesn't matter
			wantErr: true,
			err:     "OneOf: input text is empty",
		},
		{
			name:    "empty chars",
			want:    parser.Result{},
			input:   "some input",
			chars:   "",
			wantErr: true,
			err:     "OneOf: chars must not be empty",
		},
		{
			name:    "empty input and chars",
			want:    parser.Result{},
			input:   "",
			chars:   "",
			wantErr: true,
			err:     "OneOf: input text is empty",
		},
		{
			name:    "bad utf8",
			want:    parser.Result{},
			input:   "\xf8\xa1\xa1\xa1\xa1",
			chars:   "doesn't matter",
			wantErr: true,
			err:     "OneOf: input not valid utf-8",
		},
		{
			name: "match a",
			want: parser.Result{
				Value:     "a",
				Remainder: "bcdef",
			},
			input:   "abcdef",
			chars:   "abc",
			wantErr: false,
			err:     "",
		},
		{
			name: "match b",
			want: parser.Result{
				Value:     "b",
				Remainder: "acdef",
			},
			input:   "bacdef",
			chars:   "abc",
			wantErr: false,
			err:     "",
		},
		{
			name: "match c",
			want: parser.Result{
				Value:     "c",
				Remainder: "abdef",
			},
			input:   "cabdef",
			chars:   "abc",
			wantErr: false,
			err:     "",
		},
		{
			name: "match utf8 first",
			want: parser.Result{
				Value:     "語",
				Remainder: "ç日ð本Ê語",
			},
			input:   "語ç日ð本Ê語",
			chars:   "語ç日",
			wantErr: false,
			err:     "",
		},
		{
			name: "match utf8 second",
			want: parser.Result{
				Value:     "ç",
				Remainder: "日ð本Ê語",
			},
			input:   "ç日ð本Ê語",
			chars:   "語ç日",
			wantErr: false,
			err:     "",
		},
		{
			name: "match utf8 single",
			want: parser.Result{
				Value:     "本",
				Remainder: "Ê語",
			},
			input:   "本Ê語",
			chars:   "本",
			wantErr: false,
			err:     "",
		},
		{
			name:    "no match",
			want:    parser.Result{},
			input:   "abcdef",
			chars:   "xyz",
			wantErr: true,
			err:     "OneOf: no requested char (xyz) found in input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.OneOf(tt.chars)(tt.input)

			// Should only error if we wanted one
			if (err != nil) != tt.wantErr {
				t.Fatalf("\nGot error:\t%v\nWanted error:\t%v\n", err, tt.wantErr)
			}

			// If we did get an error, the message should match what we expect
			if err != nil {
				if msg := err.Error(); msg != tt.err {
					t.Fatalf("\nGot:\t%q\nWanted:\t%q\n", msg, tt.err)
				}
			}

			// The parser result should be what we expected
			if got != tt.want {
				t.Errorf("\nGot:\t%#v\nWanted:\t%#v\n", got, tt.want)
			}
		})
	}
}

func BenchmarkTake(b *testing.B) {
	input := "Please take some chars from me"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.Take(7)(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkExact(b *testing.B) {
	input := "Hello there mr exact match"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.Exact("Hello")(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkExactCaseInsensitive(b *testing.B) {
	input := "ThIs Is SpOnGeBob CaSe"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.ExactCaseInsensitive("this is")(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkChar(b *testing.B) {
	input := "X marks the spot"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.Char('X')(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTakeWhile(b *testing.B) {
	input := "  \t\t\n\n end of whitespace"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.TakeWhile(unicode.IsSpace)(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTakeUntil(b *testing.B) {
	input := "  \t\t\n\n end of whitespace"
	predicate := func(r rune) bool { return !unicode.IsSpace(r) }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.TakeUntil(predicate)(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkOneOf(b *testing.B) {
	input := "abcdef"
	chars := "abc"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.OneOf(chars)(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func ExampleTake() {
	input := "Hello I am some input for you to parser"

	got, err := parser.Take(10)(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Printf("Taken: %q\n", got.Value)
	fmt.Printf("Remainder: %q\n", got.Remainder)

	// Output: Taken: "Hello I am"
	// Remainder: " some input for you to parser"
}

func ExampleExact() {
	input := "General Kenobi! You are a bold one."

	got, err := parser.Exact("General Kenobi!")(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Printf("Taken: %q\n", got.Value)
	fmt.Printf("Remainder: %q\n", got.Remainder)

	// Output: Taken: "General Kenobi!"
	// Remainder: " You are a bold one."
}

func ExampleExactCaseInsensitive() {
	input := "GENERAL KENOBI! YOU ARE A BOLD ONE."

	got, err := parser.ExactCaseInsensitive("GEnErAl KeNobI!")(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Printf("Taken: %q\n", got.Value)
	fmt.Printf("Remainder: %q\n", got.Remainder)

	// Output: Taken: "GENERAL KENOBI!"
	// Remainder: " YOU ARE A BOLD ONE."
}

func ExampleChar() {
	input := "X marks the spot!"

	got, err := parser.Char('X')(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Printf("Taken: %q\n", got.Value)
	fmt.Printf("Remainder: %q\n", got.Remainder)

	// Output: Taken: "X"
	// Remainder: " marks the spot!"
}

func ExampleTakeWhile() {
	input := "本本本b語ç日ð本Ê語"

	pred := func(r rune) bool { return r == '本' }

	got, err := parser.TakeWhile(pred)(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Printf("Taken: %q\n", got.Value)
	fmt.Printf("Remainder: %q\n", got.Remainder)

	// Output: Taken: "本本本"
	// Remainder: "b語ç日ð本Ê語"
}

func ExampleTakeUntil() {
	input := "something <- first whitespace is here"

	got, err := parser.TakeUntil(unicode.IsSpace)(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Printf("Taken: %q\n", got.Value)
	fmt.Printf("Remainder: %q\n", got.Remainder)

	// Output: Taken: "something"
	// Remainder: " <- first whitespace is here"
}

func ExampleOneOf() {
	input := "abcdefg"

	chars := "abc" // Match any of 'a', 'b', or 'c' from input

	got, err := parser.OneOf(chars)(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Printf("Taken: %q\n", got.Value)
	fmt.Printf("Remainder: %q\n", got.Remainder)

	// Output: Taken: "a"
	// Remainder: "bcdefg"
}
