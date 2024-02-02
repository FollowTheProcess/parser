package parser_test

// The fuzz tests in here are designed to fully exercise all our error handling, identify any
// cases we haven't handled, and to try and ensure that no parser ever panics.

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/FollowTheProcess/parser"
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

	f.Fuzz(func(t *testing.T, input string, match string) {
		value, remainder, err := parser.Exact(match)(input)
		fuzzParser(t, value, remainder, err)
	})
}

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

func randomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}
