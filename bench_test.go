package parser_test

import (
	"testing"
	"unicode"

	"go.followtheprocess.codes/parser"
)

func BenchmarkTake(b *testing.B) {
	input := "Please take some chars from me"

	for b.Loop() {
		_, _, err := parser.Take(7)(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkExact(b *testing.B) {
	input := "Hello there mr exact match"

	for b.Loop() {
		_, _, err := parser.Exact("Hello")(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkExactCaseInsensitive(b *testing.B) {
	input := "ThIs Is SpOnGeBob CaSe"

	for b.Loop() {
		_, _, err := parser.ExactCaseInsensitive("this is")(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkChar(b *testing.B) {
	input := "X marks the spot"

	for b.Loop() {
		_, _, err := parser.Char('X')(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTakeWhile(b *testing.B) {
	input := "  \t\t\n\n end of whitespace"

	for b.Loop() {
		_, _, err := parser.TakeWhile(unicode.IsSpace)(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTakeWhileBetween(b *testing.B) {
	input := "latin123"

	for b.Loop() {
		_, _, err := parser.TakeWhileBetween(3, 6, unicode.IsLetter)(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTakeUntil(b *testing.B) {
	input := "  \t\t\n\n end of whitespace"
	predicate := func(r rune) bool { return !unicode.IsSpace(r) }

	for b.Loop() {
		_, _, err := parser.TakeUntil(predicate)(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTakeTo(b *testing.B) {
	input := "some words KEYWORD the rest"

	for b.Loop() {
		_, _, err := parser.TakeTo("KEYWORD")(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkOneOf(b *testing.B) {
	input := "abcdef"
	chars := "abc"

	for b.Loop() {
		_, _, err := parser.OneOf(chars)(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkNoneOf(b *testing.B) {
	input := "abcdef"
	chars := "xyz"

	for b.Loop() {
		_, _, err := parser.NoneOf(chars)(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAnyOf(b *testing.B) {
	input := "DEADBEEF and the rest"
	chars := "1234567890ABCDEF"

	for b.Loop() {
		_, _, err := parser.AnyOf(chars)(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkNotAnyOf(b *testing.B) {
	input := "69 is a number"
	chars := "abcdefghijklmnopqrstuvwxyz"

	for b.Loop() {
		_, _, err := parser.NotAnyOf(chars)(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkOptional(b *testing.B) {
	input := "v1.2.3-rc.1+build.123"
	option := "v"

	for b.Loop() {
		_, _, err := parser.Optional(option)(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMap(b *testing.B) {
	input := "Hello, World!"

	mapper := func(input string) (int, error) { return len(input), nil }

	for b.Loop() {
		_, _, err := parser.Map(parser.Take(5), mapper)(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTry(b *testing.B) {
	input := "123456)(*&^%"

	for b.Loop() {
		_, _, err := parser.Try(
			parser.TakeWhile(unicode.IsLetter),
			parser.TakeWhile(unicode.IsDigit),
		)(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMany(b *testing.B) {
	input := "abcd1234eof"

	for b.Loop() {
		_, _, err := parser.Chain(
			parser.Take(4),
			parser.TakeWhile(unicode.IsDigit),
		)(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCount(b *testing.B) {
	input := "abcabcabc"

	for b.Loop() {
		_, _, err := parser.Count(parser.Exact("abc"), 3)(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}
