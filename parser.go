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

// TakeWhileBetween returns a [Parser] that recognises the longest (lower <= len <= upper) sequence
// of utf-8 characters for which the predicate returns true.
//
// Any of the following conditions will return an error:
//   - input is empty
//   - input is not valid utf-8
//   - predicate is nil
//   - lower < 0
//   - lower > upper
//   - predicate never returns true
//   - predicate matched some chars but less than lower limit
func TakeWhileBetween(lower, upper int, predicate func(r rune) bool) Parser[string] {
	return func(input string) (string, string, error) {
		if input == "" {
			return "", "", errors.New("TakeWhileBetween: input text is empty")
		}

		if !utf8.ValidString(input) {
			return "", "", errors.New("TakeWhileBetween: input not valid utf-8")
		}

		if predicate == nil {
			return "", "", errors.New("TakeWhileBetween: predicate must be a non-nil function")
		}

		if lower < 0 {
			return "", "", fmt.Errorf("TakeWhileBetween: lower limit (%d) not allowed, must be positive integer", lower)
		}

		if lower > upper {
			return "", "", fmt.Errorf("TakeWhileBetween: invalid range, lower (%d) must be < upper (%d)", lower, upper)
		}

		// Does the predicate ever return true? Quick failure case
		if i := strings.IndexFunc(input, predicate); i == -1 {
			return "", "", errors.New("TakeWhileBetween: predicate matched no chars in input")
		}

		index := -1 // Index of last char for which predicate returns true
		for pos, char := range input {
			if !predicate(char) {
				break
			}
			// Add the byte width of the char in question because the next char is the
			// first one where predicate(char) == false, that's where we want to cut
			// the string
			index = pos + utf8.RuneLen(char)
		}

		// If we have an index, our job now is to return whichever is longest out of
		// the sequence for which the predicate returned true, or the entire input
		// up to the upper limit of chars

		startToIndex := input[:index]
		n := utf8.RuneCountInString(startToIndex)

		if n < lower {
			// The number of chars for which the predicate returned true is less
			// than our lower limit, which is an error
			return "", "", fmt.Errorf("TakeWhileBetween: predicate matched only %d chars (%s), below lower limit (%d)", n, startToIndex, lower)
		}

		if n > upper {
			// The sequence of chars for which the predicate returned true is
			// longer than our upper limit, so cut if off at upper utf-8 chars
			runes := 0  // How many runes we've scanned through
			cutOff := 0 // Index of where to cut the string off
			for pos, char := range startToIndex {
				runes++
				if runes == upper {
					// Add the byte width of the char in question because we want to
					// include it in the slice and pos is just the starting byte position
					cutOff = pos + utf8.RuneLen(char)
				}
			}

			return input[:cutOff], input[cutOff:], nil
		}

		// If we get here, we know that the number of utf-8 chars for which the predicate
		// returned true is already less than our upper limit, so we can just use the
		// index from earlier
		return input[:index], input[index:], nil
	}
}

// TakeTo returns a [Parser] that consumes characters until it first hits an exact string.
//
// If the input is empty or the exact string is not in the input, an error will be returned.
//
// The value will contain everything from the start of the input up to the first occurrence of
// match, and the remainder will contain the match and everything thereafter.
func TakeTo(match string) Parser[string] {
	return func(input string) (string, string, error) {
		if input == "" {
			return "", "", errors.New("TakeTo: input text is empty")
		}

		if !utf8.ValidString(input) {
			return "", "", errors.New("TakeTo: input not valid utf-8")
		}

		if match == "" {
			return "", "", errors.New("TakeTo: match must not be empty")
		}

		start := strings.Index(input, match)
		if start == -1 {
			return "", "", fmt.Errorf("TakeTo: match (%s) not in input", match)
		}

		return input[:start], input[start:], nil
	}
}

// OneOf returns a [Parser] that recognises one of the provided characters from the start of input.
//
// If you want to match anything other than the provided char set, use [NoneOf].
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

// NoneOf returns a [Parser] that recognises any char other than any of the provided characters
// from the start of input.
//
// It can be considered as the opposite to [OneOf].
//
// If the input or chars is empty, an error will be returned.
// Likewise if one of the chars was recognised.
func NoneOf(chars string) Parser[string] {
	return func(input string) (string, string, error) {
		if input == "" {
			return "", "", errors.New("NoneOf: input text is empty")
		}

		if chars == "" {
			return "", "", errors.New("NoneOf: chars must not be empty")
		}

		r, width := utf8.DecodeRuneInString(input)
		if r == utf8.RuneError {
			return "", "", errors.New("NoneOf: input not valid utf-8")
		}

		found := false
		for _, char := range chars {
			if char == r {
				// Found one that's not a match
				found = true
				break
			}
		}

		// If we get here and found is true, the first char in the input matched one
		// of the requested chars, which for NoneOf is bad
		if found {
			return "", "", fmt.Errorf("NoneOf: found match (%s) in input", string(r))
		}

		return input[:width], input[width:], nil
	}
}

// AnyOf returns a [Parser] that continues taking characters so long as they are contained in the
// passed in set of chars.
//
// Parsing stops at the first occurrence of a character not contained in the argument and the
// offending character is not included in the parsed value, but will be in the remainder.
//
// AnyOf is the opposite to [NotAnyOf].
//
// If the input or chars is empty, an error will be returned.
// Likewise if none of the chars are present at the start of the input.
func AnyOf(chars string) Parser[string] {
	return func(input string) (string, string, error) {
		if input == "" {
			return "", "", errors.New("AnyOf: input text is empty")
		}

		if chars == "" {
			return "", "", errors.New("AnyOf: chars must not be empty")
		}

		if !utf8.ValidString(input) {
			return "", "", errors.New("AnyOf: input not valid utf-8")
		}

		end := 0 // The end of the matching sequence
		for pos, char := range input {
			if !strings.ContainsRune(chars, char) {
				end = pos
				break
			}
		}

		// If we've broken the loop but end is still 0, there were no matches
		// in the entire input
		if end == 0 {
			return "", "", fmt.Errorf("AnyOf: no match for any char in (%s) found in input", chars)
		}

		return input[:end], input[end:], nil
	}
}

// NotAnyOf returns a [Parser] that continues taking characters so long as they are not contained
// in the passed in set of chars.
//
// Parsing stops at the first occurrence of a character contained in the argument and the
// offending character is not included in the parsed value, but will be in the remainder.
//
// NotAnyOf is the opposite of [AnyOf].
//
// If the input or chars is empty, an error will be returned.
// Likewise if any of the chars are present at the start of the input.
func NotAnyOf(chars string) Parser[string] {
	return func(input string) (string, string, error) {
		if input == "" {
			return "", "", errors.New("NotAnyOf: input text is empty")
		}

		if chars == "" {
			return "", "", errors.New("NotAnyOf: chars must not be empty")
		}

		if !utf8.ValidString(input) {
			return "", "", errors.New("NotAnyOf: input not valid utf-8")
		}

		end := 0 // The end of the matching sequence
		for pos, char := range input {
			if strings.ContainsRune(chars, char) {
				end = pos
				break
			}
		}

		// If we've broken the loop but end is still 0, there were no matches
		// in the entire input
		if end == 0 {
			return "", "", fmt.Errorf("NotAnyOf: match found for char in (%s)", chars)
		}

		return input[:end], input[end:], nil
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
