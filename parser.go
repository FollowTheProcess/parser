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

// Parser is the core parsing function that all parser functions return, they can be combined and composed
// to parse complex grammars.
//
// Each Parser is generic over type T and returns the parsed value from the input, the remaining unparsed input and an error.
type Parser[T any] func(input string) (value T, remainder string, err error)

// Take returns a [Parser] that consumes n utf-8 chars from the input.
//
// If n is less than or equal to 0, or greater than the number of utf-8 chars in the input, an error will be returned.
func Take(n int) Parser[string] {
	return func(input string) (string, string, error) {
		if n <= 0 {
			return "", "", fmt.Errorf("Take: n must be a non-zero positive integer, got %d", n)
		}

		if input == "" {
			return "", "", errors.New("Take: cannot take from empty input")
		}

		if !utf8.ValidString(input) {
			return "", "", errors.New("Take: input not valid utf-8")
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
			return "", "", fmt.Errorf("Take: requested n (%d) chars but input had only %d utf-8 chars", n, runes)
		}

		return input[:end], input[end:], nil
	}
}

// Exact returns a [Parser] that consumes an exact, case-sensitive string from the input.
//
// If the string is not present at the beginning of the input, an error will be returned.
//
// An empty match string or empty input (i.e. "") will also return an error.
//
// Exact is case-sensitive, if you need a case-insensitive match, use [ExactCaseInsensitive] instead.
func Exact(match string) Parser[string] {
	return func(input string) (string, string, error) {
		if input == "" {
			return "", "", errors.New("Exact: cannot match on empty input")
		}

		if !utf8.ValidString(input) {
			return "", "", errors.New("Exact: input not valid utf-8")
		}

		if match == "" {
			return "", "", errors.New("Exact: match must not be empty")
		}

		start := strings.Index(input, match)
		if start != 0 {
			return "", "", fmt.Errorf("Exact: match (%s) not in input", match)
		}

		return match, input[len(match):], nil
	}
}

// ExactCaseInsensitive returns a [Parser] that consumes an exact, case-insensitive string from the input.
//
// If the string is not present at the beginning of the input, an error will be returned.
//
// An empty match string or empty input (i.e. "") will also return an error.
//
// ExactCaseInsensitive is case-insensitive, if you need a case-sensitive match, use [Exact] instead.
func ExactCaseInsensitive(match string) Parser[string] {
	return func(input string) (string, string, error) {
		inputLen := len(input)
		if inputLen == 0 {
			return "", "", errors.New("ExactCaseInsensitive: cannot match on empty input")
		}

		if !utf8.ValidString(input) {
			return "", "", errors.New("ExactCaseInsensitive: input not valid utf-8")
		}

		matchLen := len(match)
		if matchLen == 0 {
			return "", "", errors.New("ExactCaseInsensitive: match must not be empty")
		}

		// Serves two purposes: It's a quick check that we'd never find a match and it guards
		// the input slicing below
		if matchLen > inputLen {
			return "", "", fmt.Errorf("ExactCaseInsensitive: match (%s) not in input", match)
		}

		// The beginning of input where the match string could possibly be
		potentialMatch := input[:matchLen]

		if !strings.EqualFold(potentialMatch, match) {
			return "", "", fmt.Errorf("ExactCaseInsensitive: match (%s) not in input", match)
		}

		return potentialMatch, input[matchLen:], nil
	}
}

// Char returns a [Parser] that consumes a single exact, case-sensitive utf-8 character from the input.
//
// If the first char in the input is not the requested char, an error will be returned.
func Char(char rune) Parser[string] {
	return func(input string) (string, string, error) {
		if input == "" {
			return "", "", errors.New("Char: input text is empty")
		}

		r, width := utf8.DecodeRuneInString(input)
		if r == utf8.RuneError {
			return "", "", errors.New("Char: input not valid utf-8")
		}

		if r != char {
			return "", "", fmt.Errorf("Char: requested char (%s) not found in input", string(char))
		}

		return input[:width], input[width:], nil
	}
}

// TakeWhile returns a [Parser] that continues consuming characters so long as the predicate returns true,
// the parsing stops as soon as the predicate returns false for a particular character. The last character
// for which the predicate returns true is captured; that is, TakeWhile is inclusive.
//
// TakeWhile can be thought of as the inverse of [TakeUntil].
//
// If the input is empty or predicate == nil, an error will be returned.
//
// If the predicate doesn't return false before the entire input is consumed, an error will be returned.
//
// A predicate that never returns true will leave the entire input unparsed and return a
// value that is an empty string, and a remainder that is the entire input.
func TakeWhile(predicate func(r rune) bool) Parser[string] {
	return func(input string) (string, string, error) {
		if input == "" {
			return "", "", errors.New("TakeWhile: input text is empty")
		}

		if !utf8.ValidString(input) {
			return "", "", errors.New("TakeWhile: input not valid utf-8")
		}

		if predicate == nil {
			return "", "", errors.New("TakeWhile: predicate must be a non-nil function")
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
			return "", "", errors.New("TakeWhile: predicate never returned false")
		}

		return input[:end], input[end:], nil
	}
}

// TakeUntil returns a [Parser] that continues taking characters until the predicate returns true,
// the parsing stops as soon as the predicate returns true for a particular character. The last character
// for which the predicate returns false is captured; that is, TakeUntil is inclusive.
//
// TakeUntil can be thought of as the inverse of [TakeWhile].
//
// If the input is empty or predicate == nil, an error will be returned.
//
// If the predicate doesn't return true before the entire input is consumed, an error will be returned
// .
//
// A predicate that never returns false will leave the entire input unparsed and return a
// value that is an empty string, and a remainder that is the entire input.
func TakeUntil(predicate func(r rune) bool) Parser[string] {
	return func(input string) (string, string, error) {
		if input == "" {
			return "", "", errors.New("TakeUntil: input text is empty")
		}

		if !utf8.ValidString(input) {
			return "", "", errors.New("TakeUntil: input not valid utf-8")
		}

		if predicate == nil {
			return "", "", errors.New("TakeUntil: predicate must be a non-nil function")
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
			return "", "", errors.New("TakeUntil: predicate never returned true")
		}

		return input[:end], input[end:], nil
	}
}

// OneOf returns a [Parser] that recognises one of the provided characters from the start of input.
//
// If the input or chars is empty, an error will be returned.
// Likewise if none of the chars was recognised.
func OneOf(chars string) Parser[string] {
	return func(input string) (string, string, error) {
		if input == "" {
			return "", "", errors.New("OneOf: input text is empty")
		}

		if chars == "" {
			return "", "", errors.New("OneOf: chars must not be empty")
		}

		r, width := utf8.DecodeRuneInString(input)
		if r == utf8.RuneError {
			return "", "", errors.New("OneOf: input not valid utf-8")
		}

		found := false // Whether we've actually found a match
		for _, char := range chars {
			if char == r {
				// Found it
				found = true
				break
			}
		}

		// If we get here and found is still false, the first char in the input didn't match
		// any of our given chars
		if !found {
			return "", "", fmt.Errorf("OneOf: no requested char (%s) found in input", chars)
		}

		return input[:width], input[width:], nil
	}
}

// Map returns a [Parser] that applies a function to the result of another parser.
//
// It is particularly useful for parsing a section of string input, then converting
// that captured string to another type.
//
// If the provided parser or the mapping function 'fn' return an error, Map will
// bubble up this error to the caller.
func Map[T1, T2 any](parser Parser[T1], fn func(T1) (T2, error)) Parser[T2] {
	return func(input string) (T2, string, error) {
		var zero T2

		// Note: Since we're applying the function to another parser
		// we don't need to check for empty input or invalid utf-8
		// because the other parser will enforce it's own invariants

		if fn == nil {
			return zero, "", errors.New("Map: fn must be a non-nil function")
		}

		// Apply the parser to the input
		value, remainder, err := parser(input)
		if err != nil {
			return zero, "", fmt.Errorf("Map: parser returned error: %w", err)
		}

		// Now apply the map function to the value returned from that
		newValue, err := fn(value)
		if err != nil {
			return zero, "", fmt.Errorf("Map: fn returned error: %w", err)
		}

		return newValue, remainder, nil
	}
}
