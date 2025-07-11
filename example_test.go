package parser_test

import (
	"fmt"
	"os"
	"strconv"

	"go.followtheprocess.codes/parser"
)

// RGB represents a colour.
type RGB struct {
	Red   int
	Green int
	Blue  int
}

// fromHex parses a string into a hex digit.
func fromHex(s string) (int, error) {
	hx, err := strconv.ParseInt(s, 16, 64)
	return int(hx), err
}

// hexPair is a parser that converts a hex string into it's integer value.
func hexPair(colour string) (int, string, error) {
	return parser.Map(
		parser.Take(2),
		fromHex,
	)(colour)
}

func Example() {
	// Let's parse this into an RGB
	colour := "#2F14DF"

	// We don't actually care about the #
	_, colour, err := parser.Char('#')(colour)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	// We want 3 hex pairs
	pairs, _, err := parser.Count(hexPair, 3)(colour)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if len(pairs) != 3 {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	rgb := RGB{
		Red:   pairs[0],
		Green: pairs[1],
		Blue:  pairs[2],
	}

	fmt.Printf("%#v\n", rgb)

	// Output: parser_test.RGB{Red:47, Green:20, Blue:223}
}
