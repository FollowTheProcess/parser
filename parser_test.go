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
		name      string // Identifying test case name
		input     string // Entire input to be parsed
		value     string // The parsed value
		remainder string // The remaining unparsed input
		err       string // The expected error message (if there is one)
		n         int    // The number of chars to consume
		wantErr   bool   // Whether it should have returned an error
	}{
		{
			name:      "empty input",
			input:     "",
			value:     "",
			remainder: "",
			n:         999,
			wantErr:   true,
			err:       "Take: cannot take from empty input",
		},
		{
			name:      "empty input with n zero",
			input:     "",
			value:     "",
			remainder: "",
			n:         0,
			wantErr:   true,
			err:       "Take: n must be a non-zero positive integer, got 0",
		},
		{
			name:      "n too large",
			input:     "some stuff here",
			value:     "",
			remainder: "",
			n:         999,
			wantErr:   true,
			err:       "Take: requested n (999) chars but input had only 15 utf-8 chars",
		},
		{
			name:      "n negative",
			input:     "some stuff here",
			value:     "",
			remainder: "",
			n:         -1,
			wantErr:   true,
			err:       "Take: n must be a non-zero positive integer, got -1",
		},
		{
			name:      "n zero",
			input:     "some stuff here",
			value:     "",
			remainder: "",
			n:         0,
			wantErr:   true,
			err:       "Take: n must be a non-zero positive integer, got 0",
		},
		{
			name:      "n 1 more than len",
			input:     "This is an exact length",
			value:     "",
			remainder: "",
			n:         24,
			wantErr:   true,
			err:       "Take: requested n (24) chars but input had only 23 utf-8 chars",
		},
		{
			name:      "n 1 more than len utf8",
			input:     "日a本b語ç日ð本Ê語þ日¥本¼語i日©", // This is 20 utf-8 runes
			value:     "",
			remainder: "",
			n:         21,
			wantErr:   true,
			err:       "Take: requested n (21) chars but input had only 20 utf-8 chars",
		},
		{
			name:      "bad utf8",
			input:     "\xf8\xa1\xa1\xa1\xa1",
			value:     "",
			remainder: "",
			n:         3,
			wantErr:   true,
			err:       "Take: input not valid utf-8",
		},
		{
			name:      "simple",
			input:     "Hello I am some input",
			value:     "Hello I am",
			remainder: " some input",
			n:         10,
			wantErr:   false,
			err:       "",
		},
		{
			name:      "n same as len",
			input:     "This is an exact length",
			value:     "This is an exact length",
			remainder: "",
			n:         23,
			wantErr:   false,
			err:       "",
		},
		{
			name:      "n same as len utf8",
			input:     "日a本b語ç日ð本Ê語þ日¥本¼語i日©", // This is 20 utf-8 runes
			value:     "日a本b語ç日ð本Ê語þ日¥本¼語i日©",
			remainder: "",
			n:         20,
			wantErr:   false,
			err:       "",
		},
		{
			name:      "n 1 less than len",
			input:     "This is an exact length",
			value:     "This is an exact lengt",
			remainder: "h",
			n:         22,
			wantErr:   false,
			err:       "",
		},
		{
			name:      "n 1 less than len utf8",
			input:     "日a本b語ç日ð本Ê語þ日¥本¼語i日©", // This is 20 utf-8 runes
			value:     "日a本b語ç日ð本Ê語þ日¥本¼語i日",
			remainder: "©",
			n:         19,
			wantErr:   false,
			err:       "",
		},
		{
			name:      "non ascii",
			input:     "Hello, 世界",
			value:     "Hello",
			remainder: ", 世界",
			n:         5,
			wantErr:   false,
			err:       "",
		},
		{
			// https://github.com/golang/exp/blob/master/utf8string/string_test.go
			name:      "utf8string test",
			input:     "日a本b語ç日ð本Ê語þ日¥本¼語i日©",
			value:     "日a本b語ç",
			remainder: "日ð本Ê語þ日¥本¼語i日©",
			n:         6,
			wantErr:   false,
			err:       "",
		},
		{
			name:      "emoji",
			input:     "😱 emoji works too",
			value:     "😱 ",
			remainder: "emoji works too",
			n:         2,
			wantErr:   false,
			err:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, remainder, err := parser.Take(tt.n)(tt.input)

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

			// The value should be as expected
			if value != tt.value {
				t.Errorf("\nGot:\t%q\nWanted:\t%q\n", value, tt.value)
			}

			// Likewise the remainder
			if remainder != tt.remainder {
				t.Errorf("\nGot:\t%q\nWanted:\t%q\n", remainder, tt.remainder)
			}
		})
	}
}

func TestExact(t *testing.T) {
	tests := []struct {
		name      string // Identifying test case name
		input     string // Entire input to be parsed
		match     string // The exact string to parser
		value     string // The parsed value
		remainder string // The remaining unparsed input
		err       string // The expected error message (if there is one)
		wantErr   bool   // Whether it should have returned an error
	}{
		{
			name:      "empty input",
			input:     "",
			value:     "",
			remainder: "",
			match:     "something",
			wantErr:   true,
			err:       "Exact: cannot match on empty input",
		},
		{
			name:      "bad utf8",
			input:     "\xf8\xa1\xa1\xa1\xa1",
			value:     "",
			remainder: "",
			match:     "something",
			wantErr:   true,
			err:       "Exact: input not valid utf-8",
		},
		{
			name:      "empty input and match",
			input:     "",
			value:     "",
			remainder: "",
			match:     "",
			wantErr:   true,
			err:       "Exact: cannot match on empty input",
		},
		{
			name:      "empty match",
			input:     "some text",
			value:     "",
			remainder: "",
			match:     "",
			wantErr:   true,
			err:       "Exact: match must not be empty",
		},
		{
			name:      "match longer than input",
			input:     "A single sentence",
			value:     "",
			remainder: "",
			match:     "A single sentence but this one is longer so it can't possibly be matched",
			wantErr:   true,
			err:       "Exact: match (A single sentence but this one is longer so it can't possibly be matched) not in input",
		},
		{
			name:      "match not found",
			input:     "Nothing to see in here",
			value:     "",
			remainder: "",
			match:     "Found me",
			wantErr:   true,
			err:       "Exact: match (Found me) not in input",
		},
		{
			name:      "wrong case",
			input:     "Found me, in a larger sentence",
			value:     "",
			remainder: "",
			match:     "found me", // Note: lower case 'f', not an exact match
			wantErr:   true,
			err:       "Exact: match (found me) not in input",
		},
		{
			name:      "simple match",
			input:     "Found me, in a larger sentence",
			value:     "Found me",
			remainder: ", in a larger sentence",
			match:     "Found me",
			wantErr:   false,
			err:       "",
		},

		{
			name:      "utf8 match",
			input:     "世界, Hello",
			value:     "世",
			remainder: "界, Hello",
			match:     "世",
			wantErr:   false,
			err:       "",
		},
		{
			// https://github.com/golang/exp/blob/master/utf8string/string_test.go
			name: "utf8string test",

			input:     "日a本b語ç日ð本Ê語þ日¥本¼語i日©",
			value:     "日a本",
			remainder: "b語ç日ð本Ê語þ日¥本¼語i日©",
			match:     "日a本",
			wantErr:   false,
			err:       "",
		},
		{
			name:      "emoji",
			input:     "😱 emoji works too",
			value:     "😱 emoji",
			remainder: " works too",
			match:     "😱 emoji",
			wantErr:   false,
			err:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, remainder, err := parser.Exact(tt.match)(tt.input)

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

			// The value should be as expected
			if value != tt.value {
				t.Errorf("\nGot:\t%q\nWanted:\t%q\n", value, tt.value)
			}

			// Likewise the remainder
			if remainder != tt.remainder {
				t.Errorf("\nGot:\t%q\nWanted:\t%q\n", remainder, tt.remainder)
			}
		})
	}
}

func TestExactCaseInsensitive(t *testing.T) {
	tests := []struct {
		name      string // Identifying test case name
		input     string // Entire input to be parsed
		match     string // The exact string to parser
		value     string // The parsed value
		remainder string // The remaining unparsed input
		err       string // The expected error message (if there is one)
		wantErr   bool   // Whether it should have returned an error
	}{
		{
			name:      "empty input",
			input:     "",
			match:     "something",
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "ExactCaseInsensitive: cannot match on empty input",
		},
		{
			name:      "bad utf8",
			input:     "\xf8\xa1\xa1\xa1\xa1",
			match:     "something",
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "ExactCaseInsensitive: input not valid utf-8",
		},
		{
			name:      "empty input and match",
			input:     "",
			match:     "",
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "ExactCaseInsensitive: cannot match on empty input",
		},
		{
			name:      "empty match",
			input:     "some text",
			match:     "",
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "ExactCaseInsensitive: match must not be empty",
		},
		{
			name:      "match longer than input",
			input:     "A single sentence",
			match:     "A single sentence but this one is longer so it can't possibly be matched",
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "ExactCaseInsensitive: match (A single sentence but this one is longer so it can't possibly be matched) not in input",
		},
		{
			name:      "match not found",
			input:     "Nothing to see in here",
			match:     "Found me",
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "ExactCaseInsensitive: match (Found me) not in input",
		},
		{
			name:      "match same length as input",
			input:     "A single sentence",
			match:     "A single sentence",
			value:     "A single sentence",
			remainder: "",
			wantErr:   false,
			err:       "",
		},

		{
			name:      "exact match",
			input:     "Found me, in a larger sentence",
			match:     "Found me",
			value:     "Found me",
			remainder: ", in a larger sentence",
			wantErr:   false,
			err:       "",
		},
		{
			name:      "case insensitive match",
			input:     "Found me, in a larger sentence",
			match:     "found me", // Lower case f, should still match
			value:     "Found me",
			remainder: ", in a larger sentence",
			wantErr:   false,
			err:       "",
		},
		{
			// https://github.com/golang/exp/blob/master/utf8string/string_test.go
			name:      "utf8string test",
			input:     "日a本b語ç日ð本Ê語þ日¥本¼語i日©",
			match:     "日a本", // Apparently this is already lower case
			value:     "日a本",
			remainder: "b語ç日ð本Ê語þ日¥本¼語i日©",
			wantErr:   false,
			err:       "",
		},
		{
			name:      "utf8string test upper case",
			input:     "日a本b語ç日ð本Ê語þ日¥本¼語i日©",
			match:     "日A本", // Upper case now
			value:     "日a本",
			remainder: "b語ç日ð本Ê語þ日¥本¼語i日©",
			wantErr:   false,
			err:       "",
		},
		{
			name:      "emoji",
			input:     "😱 EMOJI WORKS TOO",
			match:     "😱 emoji",
			value:     "😱 EMOJI",
			remainder: " WORKS TOO",
			wantErr:   false,
			err:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, remainder, err := parser.ExactCaseInsensitive(tt.match)(tt.input)

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

			// The value should be as expected
			if value != tt.value {
				t.Errorf("\nGot:\t%q\nWanted:\t%q\n", value, tt.value)
			}

			// Likewise the remainder
			if remainder != tt.remainder {
				t.Errorf("\nGot:\t%q\nWanted:\t%q\n", remainder, tt.remainder)
			}
		})
	}
}

func TestChar(t *testing.T) {
	tests := []struct {
		name      string // Identifying test case name
		input     string // Entire input to be parsed
		value     string // The parsed value
		remainder string // The remaining unparsed input
		err       string // The expected error message (if there is one)
		char      rune   // The exact char to match
		wantErr   bool   // Whether it should have returned an error
	}{
		{
			name:      "empty input",
			input:     "",
			value:     "",
			remainder: "",
			char:      0,
			wantErr:   true,
			err:       "Char: input text is empty",
		},
		{
			name:      "bad utf8",
			input:     "\xf8\xa1\xa1\xa1\xa1",
			value:     "",
			remainder: "",
			char:      'x',
			wantErr:   true,
			err:       "Char: input not valid utf-8",
		},
		{
			name:      "not found",
			input:     "something",
			value:     "",
			remainder: "",
			char:      'x',
			wantErr:   true,
			err:       "Char: requested char (x) not found in input",
		},
		{
			name:      "wrong case",
			input:     "General Kenobi!",
			value:     "",
			remainder: "",
			char:      'g',
			wantErr:   true,
			err:       "Char: requested char (g) not found in input",
		},
		{
			name:      "found",
			input:     "General Kenobi!",
			value:     "G",
			remainder: "eneral Kenobi!",
			char:      'G',
			wantErr:   false,
			err:       "",
		},
		{
			name:      "found utf8",
			input:     "日a本b語ç日ð本Ê語þ日¥本¼語i日©",
			value:     "日",
			remainder: "a本b語ç日ð本Ê語þ日¥本¼語i日©",
			char:      '日',
			wantErr:   false,
			err:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, remainder, err := parser.Char(tt.char)(tt.input)

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

			// The value should be as expected
			if value != tt.value {
				t.Errorf("\nGot:\t%q\nWanted:\t%q\n", value, tt.value)
			}

			// Likewise the remainder
			if remainder != tt.remainder {
				t.Errorf("\nGot:\t%q\nWanted:\t%q\n", remainder, tt.remainder)
			}
		})
	}
}

func TestTakeWhile(t *testing.T) {
	tests := []struct {
		predicate func(r rune) bool // The predicate function that determines whether the parser should continue taking characters
		name      string            // Identifying test case name
		input     string            // Entire input to be parsed
		value     string            // The parsed value
		remainder string            // The remaining unparsed input
		err       string            // The expected error message (if there is one)
		wantErr   bool              // Whether it should have returned an error
	}{
		{
			name:      "empty input",
			input:     "",
			value:     "",
			remainder: "",
			predicate: nil, // Shouldn't matter as it should never get called
			wantErr:   true,
			err:       "TakeWhile: input text is empty",
		},
		{
			name:      "bad utf8",
			input:     "\xf8\xa1\xa1\xa1\xa1",
			value:     "",
			remainder: "",
			predicate: nil, // Shouldn't matter as it should never get called
			wantErr:   true,
			err:       "TakeWhile: input not valid utf-8",
		},
		{
			name:      "nil predicate", // Good libraries don't panic
			input:     "some input",
			value:     "",
			remainder: "",
			predicate: nil,
			wantErr:   true,
			err:       "TakeWhile: predicate must be a non-nil function",
		},
		{
			name:      "predicate never returns false", // Good libraries don't allow infinite loops
			input:     "fixed length input",
			value:     "",
			remainder: "",
			predicate: func(r rune) bool { return true },
			wantErr:   true,
			err:       "TakeWhile: predicate never returned false",
		},
		{
			name:      "predicate never returns true",
			input:     "fixed length input",
			value:     "",
			remainder: "fixed length input",
			predicate: func(r rune) bool { return false },
			wantErr:   false,
			err:       "",
		},
		{
			name:      "consume whitespace",
			input:     "  \t\t\n\n end of whitespace",
			value:     "  \t\t\n\n ",
			remainder: "end of whitespace",
			predicate: unicode.IsSpace,
			wantErr:   false,
			err:       "",
		},
		{
			name:      "consume non ascii rune",
			input:     "本本本 b語ç日ð本Ê語",
			value:     "本本本",
			remainder: " b語ç日ð本Ê語",
			predicate: func(r rune) bool { return r == '本' },
			wantErr:   false,
			err:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, remainder, err := parser.TakeWhile(tt.predicate)(tt.input)

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

			// The value should be as expected
			if value != tt.value {
				t.Errorf("\nGot:\t%q\nWanted:\t%q\n", value, tt.value)
			}

			// Likewise the remainder
			if remainder != tt.remainder {
				t.Errorf("\nGot:\t%q\nWanted:\t%q\n", remainder, tt.remainder)
			}
		})
	}
}

func TestTakeUntil(t *testing.T) {
	tests := []struct {
		predicate func(r rune) bool // The predicate function that determines whether the parser should stop taking characters
		name      string            // Identifying test case name
		input     string            // Entire input to be parsed
		value     string            // The parsed value
		remainder string            // The remaining unparsed input
		err       string            // The expected error message (if there is one)
		wantErr   bool              // Whether it should have returned an error
	}{
		{
			name:      "empty input",
			input:     "",
			value:     "",
			remainder: "",
			predicate: nil, // Shouldn't matter as it should never get called
			wantErr:   true,
			err:       "TakeUntil: input text is empty",
		},
		{
			name:      "bad utf8",
			input:     "\xf8\xa1\xa1\xa1\xa1",
			value:     "",
			remainder: "",
			predicate: nil, // Shouldn't matter as it should never get called
			wantErr:   true,
			err:       "TakeUntil: input not valid utf-8",
		},
		{
			name:      "nil predicate", // Good libraries don't panic
			input:     "some input",
			value:     "",
			remainder: "",
			predicate: nil,
			wantErr:   true,
			err:       "TakeUntil: predicate must be a non-nil function",
		},
		{
			name:      "predicate never returns true", // Good libraries don't allow infinite loops
			input:     "fixed length input",
			value:     "",
			remainder: "",
			predicate: func(r rune) bool { return false },
			wantErr:   true,
			err:       "TakeUntil: predicate never returned true",
		},
		{
			name:      "predicate never returns false",
			input:     "fixed length input",
			value:     "",
			remainder: "fixed length input",
			predicate: func(r rune) bool { return true },
			wantErr:   false,
			err:       "",
		},
		{
			name:      "consume until whitespace",
			input:     "something <- first whitespace",
			value:     "something",
			remainder: " <- first whitespace",
			predicate: unicode.IsSpace,
			wantErr:   false,
			err:       "",
		},
		{
			name:      "consume until non-ascii",
			input:     "abcdef語ç日ð本Ê語",
			value:     "abcdef",
			remainder: "語ç日ð本Ê語",
			predicate: func(r rune) bool { return r > unicode.MaxASCII },
			wantErr:   false,
			err:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, remainder, err := parser.TakeUntil(tt.predicate)(tt.input)

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

			// The value should be as expected
			if value != tt.value {
				t.Errorf("\nGot:\t%q\nWanted:\t%q\n", value, tt.value)
			}

			// Likewise the remainder
			if remainder != tt.remainder {
				t.Errorf("\nGot:\t%q\nWanted:\t%q\n", remainder, tt.remainder)
			}
		})
	}
}

func TestOneOf(t *testing.T) {
	tests := []struct {
		name      string // Identifying test case name
		input     string // Entire input to be parsed
		chars     string // The chars to match one of
		value     string // The parsed value
		remainder string // The remaining unparsed input
		err       string // The expected error message (if there is one)
		wantErr   bool   // Whether it should have returned an error
	}{
		{
			name:      "empty input",
			input:     "",
			chars:     "abc", // Doesn't matter
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "OneOf: input text is empty",
		},
		{
			name:      "empty chars",
			input:     "some input",
			chars:     "",
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "OneOf: chars must not be empty",
		},
		{
			name:      "empty input and chars",
			input:     "",
			chars:     "",
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "OneOf: input text is empty",
		},
		{
			name:      "bad utf8",
			input:     "\xf8\xa1\xa1\xa1\xa1",
			chars:     "doesn't matter",
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "OneOf: input not valid utf-8",
		},
		{
			name:      "no match",
			input:     "abcdef",
			chars:     "xyz",
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "OneOf: no requested char (xyz) found in input",
		},
		{
			name:      "match a",
			input:     "abcdef",
			chars:     "abc",
			value:     "a",
			remainder: "bcdef",
			wantErr:   false,
			err:       "",
		},
		{
			name:      "match b",
			input:     "bacdef",
			chars:     "abc",
			value:     "b",
			remainder: "acdef",
			wantErr:   false,
			err:       "",
		},
		{
			name:      "match c",
			input:     "cabdef",
			chars:     "abc",
			value:     "c",
			remainder: "abdef",
			wantErr:   false,
			err:       "",
		},
		{
			name:      "match utf8 first",
			input:     "語ç日ð本Ê語",
			chars:     "語ç日",
			value:     "語",
			remainder: "ç日ð本Ê語",
			wantErr:   false,
			err:       "",
		},
		{
			name:      "match utf8 second",
			input:     "ç日ð本Ê語",
			chars:     "語ç日",
			value:     "ç",
			remainder: "日ð本Ê語",
			wantErr:   false,
			err:       "",
		},
		{
			name:      "match utf8 single",
			input:     "本Ê語",
			chars:     "本",
			value:     "本",
			remainder: "Ê語",
			wantErr:   false,
			err:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, remainder, err := parser.OneOf(tt.chars)(tt.input)

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

			// The value should be as expected
			if value != tt.value {
				t.Errorf("\nGot:\t%q\nWanted:\t%q\n", value, tt.value)
			}

			// Likewise the remainder
			if remainder != tt.remainder {
				t.Errorf("\nGot:\t%q\nWanted:\t%q\n", remainder, tt.remainder)
			}
		})
	}
}

func BenchmarkTake(b *testing.B) {
	input := "Please take some chars from me"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := parser.Take(7)(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkExact(b *testing.B) {
	input := "Hello there mr exact match"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := parser.Exact("Hello")(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkExactCaseInsensitive(b *testing.B) {
	input := "ThIs Is SpOnGeBob CaSe"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := parser.ExactCaseInsensitive("this is")(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkChar(b *testing.B) {
	input := "X marks the spot"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := parser.Char('X')(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTakeWhile(b *testing.B) {
	input := "  \t\t\n\n end of whitespace"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := parser.TakeWhile(unicode.IsSpace)(input)
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
		_, _, err := parser.TakeUntil(predicate)(input)
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
		_, _, err := parser.OneOf(chars)(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func ExampleTake() {
	input := "Hello I am some input for you to parser"

	value, remainder, err := parser.Take(10)(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Printf("Value: %q\n", value)
	fmt.Printf("Remainder: %q\n", remainder)

	// Output: Value: "Hello I am"
	// Remainder: " some input for you to parser"
}

func ExampleExact() {
	input := "General Kenobi! You are a bold one."

	value, remainder, err := parser.Exact("General Kenobi!")(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Printf("Value: %q\n", value)
	fmt.Printf("Remainder: %q\n", remainder)

	// Output: Value: "General Kenobi!"
	// Remainder: " You are a bold one."
}

func ExampleExactCaseInsensitive() {
	input := "GENERAL KENOBI! YOU ARE A BOLD ONE."

	value, remainder, err := parser.ExactCaseInsensitive("GEnErAl KeNobI!")(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Printf("Value: %q\n", value)
	fmt.Printf("Remainder: %q\n", remainder)

	// Output: Value: "GENERAL KENOBI!"
	// Remainder: " YOU ARE A BOLD ONE."
}

func ExampleChar() {
	input := "X marks the spot!"

	value, remainder, err := parser.Char('X')(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Printf("Value: %q\n", value)
	fmt.Printf("Remainder: %q\n", remainder)

	// Output: Value: "X"
	// Remainder: " marks the spot!"
}

func ExampleTakeWhile() {
	input := "本本本b語ç日ð本Ê語"

	pred := func(r rune) bool { return r == '本' }

	value, remainder, err := parser.TakeWhile(pred)(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Printf("Value: %q\n", value)
	fmt.Printf("Remainder: %q\n", remainder)

	// Output: Value: "本本本"
	// Remainder: "b語ç日ð本Ê語"
}

func ExampleTakeUntil() {
	input := "something <- first whitespace is here"

	value, remainder, err := parser.TakeUntil(unicode.IsSpace)(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Printf("Value: %q\n", value)
	fmt.Printf("Remainder: %q\n", remainder)

	// Output: Value: "something"
	// Remainder: " <- first whitespace is here"
}

func ExampleOneOf() {
	input := "abcdefg"

	chars := "abc" // Match any of 'a', 'b', or 'c' from input

	value, remainder, err := parser.OneOf(chars)(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Printf("Value: %q\n", value)
	fmt.Printf("Remainder: %q\n", remainder)

	// Output: Value: "a"
	// Remainder: "bcdefg"
}
