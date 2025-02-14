package parser_test

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
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

			result := parserTest[string]{
				gotValue:      value,
				gotRemainder:  remainder,
				gotErr:        err,
				wantValue:     tt.value,
				wantRemainder: tt.remainder,
				wantErr:       tt.wantErr,
				wantErrMsg:    tt.err,
			}

			testParser(t, result)
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

			result := parserTest[string]{
				gotValue:      value,
				gotRemainder:  remainder,
				gotErr:        err,
				wantValue:     tt.value,
				wantRemainder: tt.remainder,
				wantErr:       tt.wantErr,
				wantErrMsg:    tt.err,
			}

			testParser(t, result)
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

			result := parserTest[string]{
				gotValue:      value,
				gotRemainder:  remainder,
				gotErr:        err,
				wantValue:     tt.value,
				wantRemainder: tt.remainder,
				wantErr:       tt.wantErr,
				wantErrMsg:    tt.err,
			}

			testParser(t, result)
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

			result := parserTest[string]{
				gotValue:      value,
				gotRemainder:  remainder,
				gotErr:        err,
				wantValue:     tt.value,
				wantRemainder: tt.remainder,
				wantErr:       tt.wantErr,
				wantErrMsg:    tt.err,
			}

			testParser(t, result)
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
			name:      "predicate never returns false",
			input:     "123456", // All digits
			value:     "123456",
			remainder: "",
			predicate: unicode.IsDigit, // True for every char in input
			wantErr:   false,
			err:       "",
		},
		{
			name:      "predicate never returns true",
			input:     "abcdef", // All letters
			value:     "",
			remainder: "",
			predicate: unicode.IsDigit, // False for every char in input
			wantErr:   true,
			err:       "TakeWhile: predicate never returned true",
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

			result := parserTest[string]{
				gotValue:      value,
				gotRemainder:  remainder,
				gotErr:        err,
				wantValue:     tt.value,
				wantRemainder: tt.remainder,
				wantErr:       tt.wantErr,
				wantErrMsg:    tt.err,
			}

			testParser(t, result)
		})
	}
}

func TestTakeWhileBetween(t *testing.T) {
	tests := []struct {
		predicate func(r rune) bool // The predicate function that determines whether the parser should continue taking characters
		name      string            // Identifying test case name
		input     string            // Entire input to be parsed
		value     string            // The parsed value
		remainder string            // The remaining unparsed input
		err       string            // The expected error message (if there is one)
		wantErr   bool              // Whether it should have returned an error
		lower     int               // The lower limit
		upper     int               // The upper limit
	}{
		{
			name:      "empty input",
			input:     "",
			lower:     0,   // Doesn't matter
			upper:     0,   // Also doesn't matter
			predicate: nil, // Same
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "TakeWhileBetween: input text is empty",
		},
		{
			name:      "invalid utf-8",
			input:     "\xf8\xa1\xa1\xa1\xa1",
			lower:     0,   // Doesn't matter
			upper:     0,   // Also doesn't matter
			predicate: nil, // Same
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "TakeWhileBetween: input not valid utf-8",
		},
		{
			name:      "emoji",
			input:     "✅🛠️🧠⚡️⚠️😎🪜",
			lower:     6,
			upper:     8,
			predicate: unicode.IsGraphic,
			value:     "✅🛠️🧠⚡️⚠️",
			remainder: "😎🪜",
			wantErr:   false,
			err:       "",
		},
		{
			name:      "fuzz failure", // Now correctly handled!
			input:     "\U0001925e0",
			lower:     9,
			upper:     83,
			predicate: unicode.IsGraphic,
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "TakeWhileBetween: predicate matched only 0 chars (), below lower limit (9)",
		},
		{
			name:      "nil predicate",
			input:     "some valid input",
			lower:     0, // Doesn't matter
			upper:     0, // Also doesn't matter
			predicate: nil,
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "TakeWhileBetween: predicate must be a non-nil function",
		},
		{
			name:      "lower negative",
			input:     "some valid input",
			lower:     -1, // Not valid
			upper:     4,
			predicate: func(r rune) bool { return true }, // Doesn't matter, but can't be nil
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "TakeWhileBetween: lower limit (-1) not allowed, must be positive integer",
		},
		{
			name:      "upper less than lower",
			input:     "some valid input",
			lower:     4,
			upper:     2,
			predicate: func(r rune) bool { return true }, // Doesn't matter, but can't be nil
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "TakeWhileBetween: invalid range, lower (4) must be < upper (2)",
		},
		{
			name:      "nom example 1", // https://docs.rs/nom/latest/nom/bytes/complete/fn.take_while_m_n.html
			input:     "latin123",
			lower:     3,
			upper:     6,
			predicate: unicode.IsLetter,
			value:     "latin",
			remainder: "123",
			wantErr:   false,
			err:       "",
		},
		{
			name:      "nom example 2", // https://docs.rs/nom/latest/nom/bytes/complete/fn.take_while_m_n.html
			input:     "lengthy",
			lower:     3,
			upper:     6,
			predicate: unicode.IsLetter,
			value:     "length",
			remainder: "y",
			wantErr:   false,
			err:       "",
		},
		{
			name:      "nom example 3", // https://docs.rs/nom/latest/nom/bytes/complete/fn.take_while_m_n.html
			input:     "latin",
			lower:     3,
			upper:     6,
			predicate: unicode.IsLetter,
			value:     "latin",
			remainder: "",
			wantErr:   false,
			err:       "",
		},
		{
			name:      "nom example 4", // https://docs.rs/nom/latest/nom/bytes/complete/fn.take_while_m_n.html
			input:     "ed",            // Not long enough
			lower:     3,
			upper:     6,
			predicate: unicode.IsLetter,
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "TakeWhileBetween: predicate matched only 2 chars (ed), below lower limit (3)",
		},
		{
			name:      "nom example 5", // https://docs.rs/nom/latest/nom/bytes/complete/fn.take_while_m_n.html
			input:     "12345",         // Predicate never returns true
			lower:     3,
			upper:     6,
			predicate: unicode.IsLetter,
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "TakeWhileBetween: predicate never returned true",
		},
		{
			name:  "unicode",
			input: "語ç日ð本Ê語",
			lower: 2,
			upper: 4,
			predicate: func(r rune) bool {
				switch r {
				case '語', 'ç', '日':
					return true
				default:
					return false
				}
			},
			value:     "語ç日",
			remainder: "ð本Ê語",
			wantErr:   false,
			err:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, remainder, err := parser.TakeWhileBetween(tt.lower, tt.upper, tt.predicate)(tt.input)

			result := parserTest[string]{
				gotValue:      value,
				gotRemainder:  remainder,
				gotErr:        err,
				wantValue:     tt.value,
				wantRemainder: tt.remainder,
				wantErr:       tt.wantErr,
				wantErrMsg:    tt.err,
			}

			testParser(t, result)
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
			name:      "predicate never returns true",
			input:     "fixed length input",
			value:     "fixed length input",
			remainder: "",
			predicate: func(r rune) bool { return false },
			wantErr:   false,
			err:       "",
		},
		{
			name:      "predicate never returns false",
			input:     "fixed length input",
			value:     "",
			remainder: "",
			predicate: func(r rune) bool { return true },
			wantErr:   true,
			err:       "TakeUntil: predicate never returned false",
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

			result := parserTest[string]{
				gotValue:      value,
				gotRemainder:  remainder,
				gotErr:        err,
				wantValue:     tt.value,
				wantRemainder: tt.remainder,
				wantErr:       tt.wantErr,
				wantErrMsg:    tt.err,
			}

			testParser(t, result)
		})
	}
}

func TestTakeTo(t *testing.T) {
	tests := []struct {
		name      string // Identifying test case name
		input     string // Entire input to be parsed
		match     string // The exact string to stop at
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
			err:       "TakeTo: input text is empty",
		},
		{
			name:      "bad utf8",
			input:     "\xf8\xa1\xa1\xa1\xa1",
			match:     "something",
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "TakeTo: input not valid utf-8",
		},
		{
			name:      "empty input and match",
			input:     "",
			match:     "",
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "TakeTo: input text is empty",
		},
		{
			name:      "empty match",
			input:     "some text",
			match:     "",
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "TakeTo: match must not be empty",
		},
		{
			name:      "no match",
			input:     "a long sentence",
			match:     "not here",
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "TakeTo: match (not here) not in input",
		},
		{
			name:      "simple",
			input:     "lots of stuff KEYWORD more stuff",
			match:     "KEYWORD",
			value:     "lots of stuff ",
			remainder: "KEYWORD more stuff",
			wantErr:   false,
			err:       "",
		},
		{
			name:      "match at end of input",
			input:     "blah blah lots of inputeof",
			match:     "eof",
			value:     "blah blah lots of input",
			remainder: "eof",
			wantErr:   false,
			err:       "",
		},
		{
			name:      "match at start of input",
			input:     "eofblah blah lots of input",
			match:     "eof",
			value:     "",
			remainder: "eofblah blah lots of input",
			wantErr:   false,
			err:       "",
		},
		{
			name:      "multiple matches",
			input:     "blaheof blah eof lots of inputeof",
			match:     "eof",
			value:     "blah",
			remainder: "eof blah eof lots of inputeof",
			wantErr:   false,
			err:       "",
		},
		{
			name:      "match utf8",
			input:     "abcdef語ç日ð本Ê語",
			match:     "語ç日",
			value:     "abcdef",
			remainder: "語ç日ð本Ê語",
			wantErr:   false,
			err:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, remainder, err := parser.TakeTo(tt.match)(tt.input)

			result := parserTest[string]{
				gotValue:      value,
				gotRemainder:  remainder,
				gotErr:        err,
				wantValue:     tt.value,
				wantRemainder: tt.remainder,
				wantErr:       tt.wantErr,
				wantErrMsg:    tt.err,
			}

			testParser(t, result)
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

			result := parserTest[string]{
				gotValue:      value,
				gotRemainder:  remainder,
				gotErr:        err,
				wantValue:     tt.value,
				wantRemainder: tt.remainder,
				wantErr:       tt.wantErr,
				wantErrMsg:    tt.err,
			}

			testParser(t, result)
		})
	}
}

func TestNoneOf(t *testing.T) {
	tests := []struct {
		name      string // Identifying test case name
		input     string // Entire input to be parsed
		chars     string // The chars to match none of
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
			err:       "NoneOf: input text is empty",
		},
		{
			name:      "empty chars",
			input:     "some input",
			chars:     "",
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "NoneOf: chars must not be empty",
		},
		{
			name:      "empty input and chars",
			input:     "",
			chars:     "",
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "NoneOf: input text is empty",
		},
		{
			name:      "bad utf8",
			input:     "\xf8\xa1\xa1\xa1\xa1",
			chars:     "doesn't matter",
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "NoneOf: input not valid utf-8",
		},
		{
			name:      "match a",
			input:     "abcdef",
			chars:     "abc",
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "NoneOf: found match (a) in input",
		},
		{
			name:      "match b",
			input:     "bacdef",
			chars:     "abc",
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "NoneOf: found match (b) in input",
		},
		{
			name:      "match c",
			input:     "cabdef",
			chars:     "abc",
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "NoneOf: found match (c) in input",
		},
		{
			name:      "no match",
			input:     "abcdef",
			chars:     "xyz",
			value:     "a",
			remainder: "bcdef",
			wantErr:   false,
			err:       "",
		},
		{
			name:      "no match unicode",
			input:     "語ç日ð本Ê語",
			chars:     "ç日ð",
			value:     "語",
			remainder: "ç日ð本Ê語",
			wantErr:   false,
			err:       "",
		},
		{
			name:      "no match unicode single",
			input:     "本Ê語",
			chars:     "Ê",
			value:     "本",
			remainder: "Ê語",
			wantErr:   false,
			err:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, remainder, err := parser.NoneOf(tt.chars)(tt.input)

			result := parserTest[string]{
				gotValue:      value,
				gotRemainder:  remainder,
				gotErr:        err,
				wantValue:     tt.value,
				wantRemainder: tt.remainder,
				wantErr:       tt.wantErr,
				wantErrMsg:    tt.err,
			}

			testParser(t, result)
		})
	}
}

func TestAnyOf(t *testing.T) {
	tests := []struct {
		name      string // Identifying test case name
		input     string // Entire input to be parsed
		chars     string // The chars to match any of
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
			err:       "AnyOf: input text is empty",
		},
		{
			name:      "empty chars",
			input:     "some input",
			chars:     "",
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "AnyOf: chars must not be empty",
		},
		{
			name:      "empty input and chars",
			input:     "",
			chars:     "",
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "AnyOf: input text is empty",
		},
		{
			name:      "bad utf8",
			input:     "\xf8\xa1\xa1\xa1\xa1",
			chars:     "doesn't matter",
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "AnyOf: input not valid utf-8",
		},
		{
			name:      "no match",
			input:     "123 is a number",
			chars:     "abcdefg",
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "AnyOf: no match for any char in (abcdefg) found in input",
		},
		{
			name:      "match a number",
			input:     "123 is a number",
			chars:     "1234567890",
			value:     "123",
			remainder: " is a number",
			wantErr:   false,
			err:       "",
		},
		{
			name:      "match a hex digit",
			input:     "BADBABEsomething",
			chars:     "1234567890ABCDEF",
			value:     "BADBABE",
			remainder: "something",
			wantErr:   false,
			err:       "",
		},
		{
			name:      "match with space",
			input:     "DEADBEEF and the rest",
			chars:     "1234567890ABCDEF",
			value:     "DEADBEEF",
			remainder: " and the rest",
			wantErr:   false,
			err:       "",
		},
		{
			name:      "match unicode",
			input:     "語ç日ð本Ê語",
			chars:     "ðç日語",
			value:     "語ç日ð",
			remainder: "本Ê語",
			wantErr:   false,
			err:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, remainder, err := parser.AnyOf(tt.chars)(tt.input)

			result := parserTest[string]{
				gotValue:      value,
				gotRemainder:  remainder,
				gotErr:        err,
				wantValue:     tt.value,
				wantRemainder: tt.remainder,
				wantErr:       tt.wantErr,
				wantErrMsg:    tt.err,
			}

			testParser(t, result)
		})
	}
}

func TestNotAny(t *testing.T) {
	tests := []struct {
		name      string // Identifying test case name
		input     string // Entire input to be parsed
		chars     string // The chars to match none of
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
			err:       "NotAnyOf: input text is empty",
		},
		{
			name:      "empty chars",
			input:     "some input",
			chars:     "",
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "NotAnyOf: chars must not be empty",
		},
		{
			name:      "empty input and chars",
			input:     "",
			chars:     "",
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "NotAnyOf: input text is empty",
		},
		{
			name:      "bad utf8",
			input:     "\xf8\xa1\xa1\xa1\xa1",
			chars:     "doesn't matter",
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "NotAnyOf: input not valid utf-8",
		},
		{
			name:      "match",
			input:     "123 is a number",
			chars:     "123456789",
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "NotAnyOf: match found for char in (123456789)",
		},
		{
			name:      "no match a number",
			input:     "123 is a number",
			chars:     "abcdefghijklmnopqrstuvwxyz",
			value:     "123 ",
			remainder: "is a number",
			wantErr:   false,
			err:       "",
		},
		{
			name:      "no match a hex digit",
			input:     "BADBABEsomething123",
			chars:     "1234567890",
			value:     "BADBABEsomething",
			remainder: "123",
			wantErr:   false,
			err:       "",
		},
		{
			name:      "not any kind of space",
			input:     "Hello, \tWorld!",
			chars:     " \t\r\n",
			value:     "Hello,",
			remainder: " \tWorld!",
			wantErr:   false,
			err:       "",
		},
		{
			name:      "no match unicode",
			input:     "語ç日ð本Ê語",
			chars:     "Ê本",
			value:     "語ç日ð",
			remainder: "本Ê語",
			wantErr:   false,
			err:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, remainder, err := parser.NotAnyOf(tt.chars)(tt.input)

			result := parserTest[string]{
				gotValue:      value,
				gotRemainder:  remainder,
				gotErr:        err,
				wantValue:     tt.value,
				wantRemainder: tt.remainder,
				wantErr:       tt.wantErr,
				wantErrMsg:    tt.err,
			}

			testParser(t, result)
		})
	}
}

func TestOptional(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		match     string
		value     string
		remainder string
		err       string
		wantErr   bool
	}{
		{
			name:      "empty input",
			input:     "",
			match:     "something",
			value:     "",
			remainder: "",
			err:       "Optional: input text is empty",
			wantErr:   true,
		},
		{
			name:      "empty match",
			input:     "some input",
			match:     "",
			value:     "",
			remainder: "",
			err:       "Optional: match must not be empty",
			wantErr:   true,
		},
		{
			name:      "empty input and match",
			input:     "",
			match:     "",
			value:     "",
			remainder: "",
			err:       "Optional: input text is empty",
			wantErr:   true,
		},
		{
			name:      "bad utf-8",
			input:     "\xf8\xa1\xa1\xa1\xa1",
			match:     "something",
			value:     "",
			remainder: "",
			err:       "Optional: input not valid utf-8",
			wantErr:   true,
		},
		{
			name:      "option present",
			input:     "v1.2.3", // v is optional, not an error if it's not there
			match:     "v",
			value:     "v",
			remainder: "1.2.3",
			err:       "",
			wantErr:   false,
		},
		{
			name:      "option present utf-8",
			input:     "語ç日ð本Ê語",
			match:     "語ç日",
			value:     "語ç日",
			remainder: "ð本Ê語",
			err:       "",
			wantErr:   false,
		},
		{
			name:      "option not present",
			input:     "1.2.3", // v is optional, not an error if it's not there
			match:     "v",
			value:     "",
			remainder: "1.2.3",
			err:       "",
			wantErr:   false,
		},
		{
			name:      "option not present utf-8",
			input:     "ð本Ê語",
			match:     "語ç日",
			value:     "",
			remainder: "ð本Ê語",
			err:       "",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, remainder, err := parser.Optional(tt.match)(tt.input)

			result := parserTest[string]{
				gotValue:      value,
				gotRemainder:  remainder,
				gotErr:        err,
				wantValue:     tt.value,
				wantRemainder: tt.remainder,
				wantErr:       tt.wantErr,
				wantErrMsg:    tt.err,
			}

			testParser(t, result)
		})
	}
}

func TestMap(t *testing.T) {
	type test[T1, T2 any] struct {
		name      string               // Identifying test case name
		input     string               // Entire input to be parsed
		p         parser.Parser[T1]    // The parser to have the map applied to
		fn        func(T1) (T2, error) // The function the map will apply to the result of p
		value     T2                   // The parsed value
		remainder string               // The remaining unparsed input
		err       string               // The expected error message (if there is one)
		wantErr   bool
	}

	// Here we're going to make it convert strings to ints
	tests := []test[string, int]{
		{
			name:      "empty input",
			input:     "",
			p:         parser.Char('x'),
			fn:        func(input string) (int, error) { return 0, nil },
			value:     0,
			remainder: "",
			err:       "Map: parser returned error: Char: input text is empty",
			wantErr:   true,
		},
		{
			name:      "nil fn",
			input:     "",
			p:         parser.Char('x'),
			fn:        nil,
			value:     0,
			remainder: "",
			err:       "Map: fn must be a non-nil function",
			wantErr:   true,
		},
		{
			name:      "map fn error",
			input:     "blah blah blah",
			p:         parser.TakeUntil(unicode.IsSpace),
			fn:        func(input string) (int, error) { return 0, errors.New("uh oh") },
			value:     0,
			remainder: "",
			err:       "Map: fn returned error: uh oh",
			wantErr:   true,
		},
		{
			name:      "take 5",
			input:     "Hello, World",
			p:         parser.Take(5),
			fn:        func(input string) (int, error) { return len(input), nil },
			value:     5,
			remainder: ", World",
			err:       "",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, remainder, err := parser.Map(tt.p, tt.fn)(tt.input)

			result := parserTest[int]{
				gotValue:      value,
				gotRemainder:  remainder,
				gotErr:        err,
				wantValue:     tt.value,
				wantRemainder: tt.remainder,
				wantErr:       tt.wantErr,
				wantErrMsg:    tt.err,
			}

			testParser(t, result)
		})
	}
}

func TestTry(t *testing.T) {
	type test[T any] struct {
		value     T
		name      string
		input     string
		remainder string
		err       string
		parsers   []parser.Parser[T]
		wantErr   bool
	}

	tests := []test[string]{
		{
			name:  "empty input",
			input: "",
			parsers: []parser.Parser[string]{
				// Some random parsers
				parser.Take(3),
				parser.Char('X'),
				parser.TakeWhile(unicode.IsLetter),
			},
			value:     "",
			remainder: "",
			wantErr:   true,
			err:       "Try: all parsers failed",
		},
		{
			name:  "digits then symbols",
			input: "123456*&^$£@",
			parsers: []parser.Parser[string]{
				parser.TakeWhile(unicode.IsLetter), // Will fail
				parser.TakeWhile(unicode.IsDigit),  // Should return the output from this one
			},
			value:     "123456",
			remainder: "*&^$£@",
			wantErr:   false,
			err:       "",
		},
		{
			name:  "digits then symbols",
			input: "xyzabc日ð本Ê語",
			parsers: []parser.Parser[string]{
				parser.OneOf("abc"),                // Will fail
				parser.Char('本'),                   // Same
				parser.ExactCaseInsensitive("XyZ"), // Should succeed
			},
			value:     "xyz",
			remainder: "abc日ð本Ê語",
			wantErr:   false,
			err:       "",
		},
		{
			name:  "first is successful",
			input: "hello there",
			parsers: []parser.Parser[string]{
				parser.Take(2),        // Will succeed and return
				parser.Exact("hello"), // Should never get invoked
			},
			value:     "he",
			remainder: "llo there",
			wantErr:   false,
			err:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, remainder, err := parser.Try(tt.parsers...)(tt.input)

			result := parserTest[string]{
				gotValue:      value,
				gotRemainder:  remainder,
				gotErr:        err,
				wantValue:     tt.value,
				wantRemainder: tt.remainder,
				wantErr:       tt.wantErr,
				wantErrMsg:    tt.err,
			}

			testParser(t, result)
		})
	}
}

func TestChain(t *testing.T) {
	type test[T any] struct {
		value     []T
		name      string
		input     string
		remainder string
		err       string
		parsers   []parser.Parser[T]
		wantErr   bool
	}

	tests := []test[string]{
		{
			name:  "empty input",
			input: "",
			parsers: []parser.Parser[string]{
				// Some random parsers
				parser.Take(3),
				parser.Char('X'),
				parser.TakeWhile(unicode.IsLetter),
			},
			value:     nil,
			remainder: "",
			wantErr:   true,
			err:       "Chain: sub parser failed: Take: cannot take from empty input",
		},
		{
			name:  "pairs of chars",
			input: "abcd1234",
			parsers: []parser.Parser[string]{
				// Chain pairs
				parser.Take(2),
				parser.Take(2),
				parser.Take(2),
				parser.Take(2),
			},
			value:     []string{"ab", "cd", "12", "34"},
			remainder: "",
			wantErr:   false,
			err:       "",
		},
		{
			name:  "too many pairs of chars",
			input: "abcd1234",
			parsers: []parser.Parser[string]{
				// Chain pairs
				parser.Take(2),
				parser.Take(2),
				parser.Take(2),
				parser.Take(2),
				parser.Take(2), // One more than there is in the input
			},
			value:     nil,
			remainder: "",
			wantErr:   true,
			err:       "Chain: sub parser failed: Take: cannot take from empty input",
		},
		{
			name:  "combo",
			input: "abcd1234exact \t\n 日ð本Ê語 eof",
			parsers: []parser.Parser[string]{
				parser.TakeWhile(unicode.IsLetter), // Get abcd
				parser.TakeWhile(unicode.IsDigit),  // 1234
				parser.Exact("exact"),              // exact
				parser.TakeWhile(unicode.IsSpace),  // Consume all whitespace
				parser.TakeTo("語"),                 // Take up to this char
				parser.Char('語'),                   // Take that char
				parser.TakeWhile(unicode.IsSpace),  // More space
				parser.Exact("eof"),                // Boom
			},
			value:     []string{"abcd", "1234", "exact", " \t\n ", "日ð本Ê", "語", " ", "eof"},
			remainder: "",
			wantErr:   false,
			err:       "",
		},
		{
			name:  "remainder left",
			input: "abcd1234",
			parsers: []parser.Parser[string]{
				parser.TakeWhile(unicode.IsLetter), // Get abcd
				parser.Take(2),                     // 12
			},
			value:     []string{"abcd", "12"},
			remainder: "34",
			wantErr:   false,
			err:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, remainder, err := parser.Chain(tt.parsers...)(tt.input)

			// Can't use the helper as []string is not comparable

			// Should only error if we wanted one
			if (err != nil) != tt.wantErr {
				t.Fatalf("\nGot error:\t%v\nWanted error:\t%v\n", err, tt.wantErr)
			}

			// If we did get an error, the message should match what we expect
			if err != nil {
				if msg := err.Error(); msg != tt.err {
					t.Fatalf("\nError message:\t%q\nWanted:\t%q\n", msg, tt.err)
				}
			}

			// The value should be as expected
			if !reflect.DeepEqual(value, tt.value) {
				t.Errorf("\nValue:\t%#v\nWanted:\t%#v\n", value, tt.value)
			}

			// Likewise the remainder
			if remainder != tt.remainder {
				t.Errorf("\nRemainder:\t%q\nWanted:\t%q\n", remainder, tt.remainder)
			}
		})
	}
}

func TestCount(t *testing.T) {
	type test[T any] struct {
		p         parser.Parser[T] // The parser to apply
		name      string           // Identifying test case name
		input     string           // Input to the parser
		remainder string           // Expected remainder after parsing
		err       string           // The expected error message, if there was one
		value     []T              // The expected value after parsing
		count     int              // Number of times to apply p to input
		wantErr   bool             // Whether or not we wanted an error
	}

	tests := []test[string]{
		{
			name:      "empty input",
			input:     "",
			p:         parser.Take(2),
			count:     2,
			value:     nil,
			remainder: "",
			wantErr:   true,
			err:       "Count: parser failed: Take: cannot take from empty input",
		},
		{
			name:      "take pairs",
			input:     "123456",
			p:         parser.Take(2),
			count:     3,
			value:     []string{"12", "34", "56"},
			remainder: "",
			wantErr:   false,
			err:       "",
		},
		{
			name:      "input too short",
			input:     "abcabcabc",
			p:         parser.Exact("abc"),
			count:     4,
			value:     nil,
			remainder: "",
			wantErr:   true,
			err:       "Count: parser failed: Exact: cannot match on empty input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, remainder, err := parser.Count(tt.p, tt.count)(tt.input)

			// Can't use the helper as []string is not comparable

			// Should only error if we wanted one
			if (err != nil) != tt.wantErr {
				t.Fatalf("\nGot error:\t%v\nWanted error:\t%v\n", err, tt.wantErr)
			}

			// If we did get an error, the message should match what we expect
			if err != nil {
				if msg := err.Error(); msg != tt.err {
					t.Fatalf("\nError message:\t%q\nWanted:\t%q\n", msg, tt.err)
				}
			}

			// The value should be as expected
			if !reflect.DeepEqual(value, tt.value) {
				t.Errorf("\nValue:\t%#v\nWanted:\t%#v\n", value, tt.value)
			}

			// Likewise the remainder
			if remainder != tt.remainder {
				t.Errorf("\nRemainder:\t%q\nWanted:\t%q\n", remainder, tt.remainder)
			}
		})
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

func ExampleTakeWhileBetween() {
	input := "2F14DF" // A hex colour (minus the #)

	isHexDigit := func(r rune) bool {
		_, err := strconv.ParseUint(string(r), 16, 64)
		return err == nil
	}

	value, remainder, err := parser.TakeWhileBetween(2, 2, isHexDigit)(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Printf("Value: %q\n", value)
	fmt.Printf("Remainder: %q\n", remainder)

	// Output: Value: "2F"
	// Remainder: "14DF"
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

func ExampleTakeTo() {
	input := "lots of stuff KEYWORD more stuff"

	value, remainder, err := parser.TakeTo("KEYWORD")(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Printf("Value: %q\n", value)
	fmt.Printf("Remainder: %q\n", remainder)

	// Output: Value: "lots of stuff "
	// Remainder: "KEYWORD more stuff"
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

func ExampleNoneOf() {
	input := "abcdefg"

	chars := "xyz" // Match anything other than 'x', 'y', or 'z' from input

	value, remainder, err := parser.NoneOf(chars)(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Printf("Value: %q\n", value)
	fmt.Printf("Remainder: %q\n", remainder)

	// Output: Value: "a"
	// Remainder: "bcdefg"
}

func ExampleAnyOf() {
	input := "DEADBEEF and the rest"

	chars := "1234567890ABCDEF" // Any hexadecimal digit

	value, remainder, err := parser.AnyOf(chars)(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Printf("Value: %q\n", value)
	fmt.Printf("Remainder: %q\n", remainder)

	// Output: Value: "DEADBEEF"
	// Remainder: " and the rest"
}

func ExampleNotAnyOf() {
	input := "69 is a number"

	chars := "abcdefghijklmnopqrstuvwxyz" // Parse until we hit any lowercase letter

	value, remainder, err := parser.NotAnyOf(chars)(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Printf("Value: %q\n", value)
	fmt.Printf("Remainder: %q\n", remainder)

	// Output: Value: "69 "
	// Remainder: "is a number"
}

func ExampleOptional() {
	input := "12.6.7-rc.2" // A semver, but could have an optional v

	// Doesn't matter...
	value, remainder, err := parser.Optional("v")(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Printf("Value: %q\n", value)
	fmt.Printf("Remainder: %q\n", remainder)

	// Output: Value: ""
	// Remainder: "12.6.7-rc.2"
}

func ExampleMap() {
	input := "27 <- this is a number" // Let's convert it to an int!

	value, remainder, err := parser.Map(parser.Take(2), strconv.Atoi)(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Printf("Value %[1]d is type %[1]T\n", value)
	fmt.Printf("Remainder: %q\n", remainder)

	// Output: Value 27 is type int
	// Remainder: " <- this is a number"
}

func ExampleTry() {
	input := "xyzabc日ð本Ê語"

	value, remainder, err := parser.Try(
		parser.OneOf("abc"),                // Will fail
		parser.Char('本'),                   // Same
		parser.ExactCaseInsensitive("XyZ"), // Should succeed, this is the output we'll get
	)(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Printf("Value: %q\n", value)
	fmt.Printf("Remainder: %q\n", remainder)

	// Output: Value: "xyz"
	// Remainder: "abc日ð本Ê語"
}

func ExampleChain() {
	input := "1234abcd\t\n日ð本rest..."

	value, remainder, err := parser.Chain(
		// Can do this is a number of ways, but here's one!
		parser.TakeWhile(unicode.IsDigit),
		parser.Exact("abcd"),
		parser.TakeWhile(unicode.IsSpace),
		parser.Char('日'),
		parser.Char('ð'),
		parser.Char('本'),
	)(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Printf("Value: %#v\n", value)
	fmt.Printf("Remainder: %q\n", remainder)

	// Output: Value: []string{"1234", "abcd", "\t\n", "日", "ð", "本"}
	// Remainder: "rest..."
}

func ExampleCount() {
	input := "12345678rest..." // Pairs of digits with a bit on the end

	value, remainder, err := parser.Count(parser.Take(2), 4)(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Printf("Value: %#v\n", value)
	fmt.Printf("Remainder: %q\n", remainder)

	// Output: Value: []string{"12", "34", "56", "78"}
	// Remainder: "rest..."
}

// parserTest is a simple structure to encapsulate everything we need to test about
// the result of applying a parser to some input.
type parserTest[T comparable] struct {
	gotErr        error  // The error the parser returned
	gotValue      T      // The value the parser actually returned
	wantValue     T      // The expected value
	gotRemainder  string // The remainder the parser actually returned
	wantRemainder string // The expected remainder
	wantErrMsg    string // The expected error message, if any
	wantErr       bool   // Whether we wanted an error or not
}

// testParser is a test helper that takes in the results of applying a parser and performs
// all the testing for us so that this code exists in one place, rather than in every test.
func testParser[T comparable](t *testing.T, p parserTest[T]) {
	t.Helper()

	// Should only error if we wanted one
	if (p.gotErr != nil) != p.wantErr {
		t.Fatalf("\nGot error:\t%v\nWanted error:\t%v\n", p.gotErr, p.wantErr)
	}

	// If we did get an error, the message should match what we expect
	if p.gotErr != nil {
		if msg := p.gotErr.Error(); msg != p.wantErrMsg {
			t.Fatalf("\nError message:\t%q\nWanted:\t%q\n", msg, p.wantErrMsg)
		}
	}

	// The value should be as expected
	if p.gotValue != p.wantValue {
		t.Errorf("\nValue:\t%#v\nWanted:\t%#v\n", p.gotValue, p.wantValue)
	}

	// Likewise the remainder
	if p.gotRemainder != p.wantRemainder {
		t.Errorf("\nRemainder:\t%q\nWanted:\t%q\n", p.gotRemainder, p.wantRemainder)
	}
}
