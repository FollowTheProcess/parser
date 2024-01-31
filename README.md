# Parser

[![License](https://img.shields.io/github/license/FollowTheProcess/parser)](https://github.com/FollowTheProcess/parser)
[![Go Reference](https://pkg.go.dev/badge/github.com/FollowTheProcess/parser.svg)](https://pkg.go.dev/github.com/FollowTheProcess/parser)
[![Go Report Card](https://goreportcard.com/badge/github.com/FollowTheProcess/parser)](https://goreportcard.com/report/github.com/FollowTheProcess/parser)
[![GitHub](https://img.shields.io/github/v/release/FollowTheProcess/parser?logo=github&sort=semver)](https://github.com/FollowTheProcess/parser)
[![CI](https://github.com/FollowTheProcess/parser/workflows/CI/badge.svg)](https://github.com/FollowTheProcess/parser/actions?query=workflow%3ACI)
[![codecov](https://codecov.io/gh/FollowTheProcess/parser/branch/main/graph/badge.svg)](https://codecov.io/gh/FollowTheProcess/parser)

Simple, fast, zero-allocation [combinatorial parsing] with Go

## Project Description

`parser` is intended to be a simple, expressive and easy to use API for all your text parsing needs. It aims to be:

- **Fast:** Performant text parsing can be tricky, `parser` aims to be as fast as possible without compromising safety or error handling. Every parser function has a benchmark and has been written with performance in mind, almost none of them allocate on the heap ⚡️
- **Correct:** You get the correct behaviour at all times, on any valid UTF-8 text. Errors are well handled and reported for easy debugging. 100% test coverage.
- **Intuitive:** Some parser combinator libraries are tricky to wrap your head around, I want `parser` to be super simple to use so that anyone can pick it up and be productive quickly
- **Well Documented:** Every combinator in `parser` has a comprehensive doc comment describing it's entire behaviour, as well as an executable example of its use

## Installation

```shell
go get github.com/FollowTheProcess/parser@latest
```

## Quickstart

Let's borrow the [nom] example and parse a hex colour!

```go
package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/FollowTheProcess/parser"
)

// RGB represents a colour.
type RGB struct {
	Red   int
	Green int
	Blue  int
}

// fromHex parses a string into a hex digit.
func fromHex(s string) (int, error) {
	hx, err := strconv.ParseUint(s, 16, 64)
	return int(hx), err
}

// hexPair is a parser that converts a hex string into it's integer value.
func hexPair(colour string) (int, string, error) {
	return parser.Map(
		parser.Take(2),
		fromHex,
	)(colour)
}

func main() {
	// Let's parse this into an RGB
	colour := "#2F14DF"

	// We don't actually care about the #
	_, colour, err := parser.Char('#')(colour)
	if err != nil {
		log.Fatalln(err)
	}

	// We want 3 hex pairs
	pairs, _, err := parser.Count(hexPair, 3)(colour)
	if err != nil {
		log.Fatalln(err)
	}

	if len(pairs) != 3 {
		log.Fatalln("Not enough pairs")
	}

	rgb := RGB{
		Red:   pairs[0],
		Green: pairs[1],
		Blue:  pairs[2],
	}

	fmt.Printf("%#v\n", rgb) // main.RGB{Red:47, Green:20, Blue:223}
}

```

### Credits

This package was created with [copier] and the [FollowTheProcess/go_copier] project template.

It is also heavily inspired by [nom], an excellent combinatorial parsing library written in [Rust].

[copier]: https://copier.readthedocs.io/en/stable/
[FollowTheProcess/go_copier]: https://github.com/FollowTheProcess/go_copier
[combinatorial parsing]: https://en.wikipedia.org/wiki/Parser_combinator
[nom]: https://github.com/rust-bakery/nom
[Rust]: https://www.rust-lang.org
