// Package parser implements simple, yet expressive mechanisms for [combinatorial parsing] in Go.
//
// [combinatorial parsing]: https://en.wikipedia.org/wiki/Parser_combinator
package parser

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

// Result is a container for the result of a parsing operation, containing the result of
// applying a single parser, and the remaining unparsed input.
type Result struct {
	Value     string // The parsed value
	Remainder string // The remaining unparsed input
}

// Parser is the core parsing function that all parser functions return, they can be combined and composed
// to parse complex grammars.
//
// Each Parser returns a [Result] and an error. The [Result] contains the parsed section of the input (what you asked for)
// and the remaining unparsed input, useful for passing to the next parser in a chain.
type Parser func(input string) (Result, error)

// Take returns a [Parser] that consumes n utf-8 chars from the input.
//
// If n is less than or equal to 0, or greater than the number of utf-8 chars in the input, an error will be returned along with an empty [Result].
func Take(n int) Parser {
	return func(input string) (Result, error) {
		if n <= 0 {
			return Result{}, fmt.Errorf("Take: n must be a non-zero positive integer, got %d", n)
		}
		if input == "" {
			return Result{}, errors.New("Take: cannot take from empty input")
		}

		runes := 0 // How many runes we've seen
		end := 0   // The starting byte position of the nth rune
		for pos, char := range input {
			runes++
			if runes == n {
				// We've hit our limit, pos is the starting byte of the nth rune
				// but we want to return the entire string *including* the nth rune
				// so the actual end is pos + the byte length of the nth rune
				end = pos + utf8.RuneLen(char)
				break
			}
		}

		if runes < n {
			// We've exhausted the entire input before scanning n runes i.e the input
			// was not long enough
			return Result{}, fmt.Errorf("Take: requested n (%d) chars but input had only %d utf-8 chars", n, runes)
		}

		return Result{
			Value:     input[:end],
			Remainder: input[end:],
		}, nil
	}
}

// Exact returns a [Parser] that consumes an exact, case-sensitive string from the input.
//
// If the string is not present at the beginning of the input, an error will be returned along with an empty [Result].
//
// An empty match string or empty input (i.e. "") will also return an error.
//
// Exact is case-sensitive, if you need a case-insensitive match, use [ExactCaseInsensitive] instead.
func Exact(match string) Parser {
	return func(input string) (Result, error) {
		if input == "" {
			return Result{}, errors.New("Exact: cannot match on empty input")
		}

		if match == "" {
			return Result{}, errors.New("Exact: match must not be empty")
		}

		start := strings.Index(input, match)
		if start != 0 {
			return Result{}, fmt.Errorf("Exact: match (%s) not in input", match)
		}

		return Result{
			Value:     match, // Because we know it's an exact match
			Remainder: input[len(match):],
		}, nil
	}
}

// ExactCaseInsensitive returns a [Parser] that consumes an exact, case-insensitive string from the input.
//
// If the string is not present at the beginning of the input, an error will be returned along with an empty [Result].
//
// An empty match string or empty input (i.e. "") will also return an error.
//
// ExactCaseInsensitive is case-insensitive, if you need a case-sensitive match, use [Exact] instead.
func ExactCaseInsensitive(match string) Parser {
	return func(input string) (Result, error) {
		if input == "" {
			return Result{}, errors.New("ExactCaseInsensitive: cannot match on empty input")
		}

		if match == "" {
			return Result{}, errors.New("ExactCaseInsensitive: match must not be empty")
		}

		// TODO: strings.ToLower() is 100% correct but it allocates a new string each time,
		// let's try and figure out how to do this without allocating if possible
		start := strings.Index(strings.ToLower(input), strings.ToLower(match))
		if start != 0 {
			return Result{}, fmt.Errorf("ExactCaseInsensitive: match (%s) not in input", match)
		}

		return Result{
			Value:     input[:len(match)], // We want to return the original match in it's original case
			Remainder: input[len(match):],
		}, nil
	}
}

// Char returns a [Parser] that consumes a single exact, case-sensitive utf-8 character from the input.
//
// If the first char in the input is not the requested char, an error will be returned along with an empty [Result].
func Char(char rune) Parser {
	return func(input string) (Result, error) {
		if input == "" {
			return Result{}, errors.New("Char: input text is empty")
		}

		// TODO: Should probably handle the decode error (r == RuneError and width == 1)
		// strings aren't *guaranteed* to be valid utf-8. Need to google some test cases to exercise this
		r, width := utf8.DecodeRuneInString(input)
		if r != char {
			return Result{}, fmt.Errorf("Char: requested char (%s) not found in input", string(char))
		}

		return Result{
			Value:     input[:width],
			Remainder: input[width:],
		}, nil
	}
}

// TakeWhile returns a [Parser] that continues consuming characters so long as the predicate returns true,
// the parsing stops as soon as the predicate returns false for a particular character. The last character
// for which the predicate returns true is captured; that is, TakeWhile is inclusive.
//
// TakeWhile can be thought of as the inverse of [TakeUntil].
//
// If the input is empty or predicate == nil, an error will be returned along with an empty [Result].
//
// If the predicate doesn't return false before the entire input is consumed, an error will be returned
// along with an empty [Result].
//
// A predicate that never returns true will leave the input unparsed and return a [Result] who's
// value is an empty string, and who's remainder is the entire input.
func TakeWhile(predicate func(r rune) bool) Parser {
	return func(input string) (Result, error) {
		if input == "" {
			return Result{}, errors.New("TakeWhile: input text is empty")
		}

		if predicate == nil {
			return Result{}, errors.New("TakeWhile: predicate must be a non-nil function")
		}

		end := 0        // Byte position of last rune that the predicate returns true for
		broken := false // Whether the predicate ever returned false so we broke the loop
		for pos, char := range input {
			end = pos
			if !predicate(char) {
				broken = true
				break
			}
		}

		if !broken {
			return Result{}, errors.New("TakeWhile: predicate never returned false")
		}

		return Result{
			Value:     input[:end],
			Remainder: input[end:],
		}, nil
	}
}

// TakeUntil returns a [Parser] that continues taking characters until the predicate returns true,
// the parsing stops as soon as the predicate returns true for a particular character. The last character
// for which the predicate returns false is captured; that is, TakeUntil is inclusive.
//
// TakeUntil can be thought of as the inverse of [TakeWhile].
//
// If the input is empty or predicate == nil, an error will be returned along with an empty [Result].
//
// If the predicate doesn't return true before the entire input is consumed, an error will be returned
// along with an empty [Result].
//
// A predicate that never returns false will leave the entire input unparsed and return a [Result] who's
// value is an empty string, and who's remainder is the entire input.
func TakeUntil(predicate func(r rune) bool) Parser {
	return func(input string) (Result, error) {
		if input == "" {
			return Result{}, errors.New("TakeUntil: input text is empty")
		}

		if predicate == nil {
			return Result{}, errors.New("TakeUntil: predicate must be a non-nil function")
		}

		end := 0        // Byte position of last rune that the predicate returns false for
		broken := false // Whether the predicate ever returned true so we broke the loop
		for pos, char := range input {
			end = pos
			if predicate(char) {
				broken = true
				break
			}
		}

		if !broken {
			return Result{}, errors.New("TakeUntil: predicate never returned true")
		}

		return Result{
			Value:     input[:end],
			Remainder: input[end:],
		}, nil
	}
}

// OneOf returns a [Parser] that recognises one of the provided characters from the start of input.
//
// If the input or chars is empty, an error will be returned along with an empty [Result].
// Likewise if none of the chars was recognised.
func OneOf(chars string) Parser {
	return func(input string) (Result, error) {
		if input == "" {
			return Result{}, errors.New("OneOf: input text is empty")
		}

		if chars == "" {
			return Result{}, errors.New("OneOf: chars must not be empty")
		}

		// TODO: Should probably handle the decode error (inputChar == RuneError and width == 1)
		// strings aren't *guaranteed* to be valid utf-8. Need to google some test cases to exercise this
		inputChar, width := utf8.DecodeRuneInString(input)

		found := false // Whether we've actually found a match
		for _, char := range chars {
			if char == inputChar {
				// Found it
				found = true
				break
			}
		}

		// If we get here and found is still false, the first char in the input didn't match
		// any of our given chars
		if !found {
			return Result{}, fmt.Errorf("OneOf: no requested char (%s) found in input", chars)
		}

		return Result{
			Value:     input[:width],
			Remainder: input[width:],
		}, nil
	}
}
