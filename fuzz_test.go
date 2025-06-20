package parser_test

// The fuzz tests in here are designed to fully exercise all our error handling, identify any
// cases we haven't handled, and to try and ensure that no parser ever panics.

import (
	"math/rand/v2"
	"reflect"
	"testing"
	"unicode"

	"go.followtheprocess.codes/parser"
)

var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@£$%^&*()_+][';/.,]語ç日ð本Ê語")

var corpus = [...]string{
	"",
	"a normal sentence",
	"日a本b語ç日ð本Ê語þ日¥本¼語i日©",
	"\xf8\xa1\xa1\xa1\xa1",
	"£$%^&*(((())))",
	"91836347287",
	"日ð本Ê語þ日¥本¼語i",
	"✅🛠️🧠⚡️⚠️😎🪜",
	"\n\n\r\n\t   ",
}

func FuzzTake(f *testing.F) {
	for _, item := range corpus {
		f.Add(item, rand.Int())
	}

	f.Fuzz(func(t *testing.T, input string, n int) {
		value, remainder, err := parser.Take(n)(input)
		fuzzParser(t, value, remainder, err)
	})
}

func FuzzExact(f *testing.F) {
	for _, item := range corpus {
		f.Add(item, randomString(5))
	}

	f.Fuzz(func(t *testing.T, input, match string) {
		value, remainder, err := parser.Exact(match)(input)
		fuzzParser(t, value, remainder, err)
	})
}

func FuzzExactCaseInsensitive(f *testing.F) {
	for _, item := range corpus {
		f.Add(item, randomString(5))
	}

	f.Fuzz(func(t *testing.T, input, match string) {
		value, remainder, err := parser.ExactCaseInsensitive(match)(input)
		fuzzParser(t, value, remainder, err)
	})
}

func FuzzChar(f *testing.F) {
	for _, item := range corpus {
		f.Add(item, randomRune())
	}

	f.Fuzz(func(t *testing.T, input string, char rune) {
		value, remainder, err := parser.Char(char)(input)
		fuzzParser(t, value, remainder, err)
	})
}

func FuzzTakeWhile(f *testing.F) {
	for _, item := range corpus {
		f.Add(item)
	}

	f.Fuzz(func(t *testing.T, input string) {
		value, remainder, err := parser.TakeWhile(unicode.IsLetter)(input)
		fuzzParser(t, value, remainder, err)
	})
}

func FuzzTakeUntil(f *testing.F) {
	for _, item := range corpus {
		f.Add(item)
	}

	f.Fuzz(func(t *testing.T, input string) {
		value, remainder, err := parser.TakeUntil(unicode.IsSpace)(input)
		fuzzParser(t, value, remainder, err)
	})
}

func FuzzTakeWhileBetween(f *testing.F) {
	for _, item := range corpus {
		f.Add(item, rand.IntN(10), rand.IntN(10))
	}

	f.Fuzz(func(t *testing.T, input string, lower, upper int) {
		value, remainder, err := parser.TakeWhileBetween(lower, upper, unicode.IsGraphic)(input)
		fuzzParser(t, value, remainder, err)
	})
}

func FuzzTakeTo(f *testing.F) {
	for _, item := range corpus {
		f.Add(item, randomString(5))
	}

	f.Fuzz(func(t *testing.T, input, match string) {
		value, remainder, err := parser.TakeTo(match)(input)
		fuzzParser(t, value, remainder, err)
	})
}

func FuzzOneOf(f *testing.F) {
	for _, item := range corpus {
		f.Add(item, randomString(rand.IntN(10)))
	}

	f.Fuzz(func(t *testing.T, input, chars string) {
		value, remainder, err := parser.OneOf(chars)(input)
		fuzzParser(t, value, remainder, err)
	})
}

func FuzzNoneOf(f *testing.F) {
	for _, item := range corpus {
		f.Add(item, randomString(rand.IntN(10)))
	}

	f.Fuzz(func(t *testing.T, input, chars string) {
		value, remainder, err := parser.NoneOf(chars)(input)
		fuzzParser(t, value, remainder, err)
	})
}

func FuzzAnyOf(f *testing.F) {
	for _, item := range corpus {
		f.Add(item, randomString(rand.IntN(10)))
	}

	f.Fuzz(func(t *testing.T, input, chars string) {
		value, remainder, err := parser.AnyOf(chars)(input)
		fuzzParser(t, value, remainder, err)
	})
}

func FuzzNotAnyOf(f *testing.F) {
	for _, item := range corpus {
		f.Add(item, randomString(rand.IntN(10)))
	}

	f.Fuzz(func(t *testing.T, input, chars string) {
		value, remainder, err := parser.NotAnyOf(chars)(input)
		fuzzParser(t, value, remainder, err)
	})
}

func FuzzOptional(f *testing.F) {
	for _, item := range corpus {
		f.Add(item, randomString(5))
	}

	f.Fuzz(func(t *testing.T, input, match string) {
		value, remainder, err := parser.Optional(match)(input)
		fuzzParser(t, value, remainder, err)
	})
}

// fuzzParser is a helper that asserts empty value and remainders were returned if the
// err was not nil.
func fuzzParser[T any](t *testing.T, value T, remainder string, err error) {
	t.Helper()

	var zero T // The zero value of type T

	// If err is not nil, value and remainder must be empty
	if err != nil {
		if !reflect.DeepEqual(value, zero) {
			t.Errorf("Value: %#v, Wanted: %#v", value, zero)
		}
		if !reflect.DeepEqual(remainder, zero) {
			t.Errorf("Remainder: %#v, Wanted: %#v", remainder, zero)
		}
	}
}

// generate a random utf-8 string of length n.
func randomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = chars[rand.IntN(len(chars))]
	}
	return string(b)
}

// generate a random utf-8 rune.
func randomRune() rune {
	return chars[rand.IntN(len(chars))]
}
