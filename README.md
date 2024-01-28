# Parser

[![License](https://img.shields.io/github/license/FollowTheProcess/parser)](https://github.com/FollowTheProcess/parser)
[![Go Reference](https://pkg.go.dev/badge/github.com/FollowTheProcess/parser.svg)](https://pkg.go.dev/github.com/FollowTheProcess/parser)
[![Go Report Card](https://goreportcard.com/badge/github.com/FollowTheProcess/parser)](https://goreportcard.com/report/github.com/FollowTheProcess/parser)
[![GitHub](https://img.shields.io/github/v/release/FollowTheProcess/parser?logo=github&sort=semver)](https://github.com/FollowTheProcess/parser)
[![CI](https://github.com/FollowTheProcess/parser/workflows/CI/badge.svg)](https://github.com/FollowTheProcess/parser/actions?query=workflow%3ACI)
[![codecov](https://codecov.io/gh/FollowTheProcess/parser/branch/main/graph/badge.svg)](https://codecov.io/gh/FollowTheProcess/parser)

Simple, fast, zero-allocation [combinatorial parsing] with Go

> [!WARNING]
> **Parser is in early development and is not yet ready for use**

## Project Description

`parser` is intended to be a simple, expressive and easy to use API for all your text parsing needs. It aims to be:

- **Fast:** Performant text parsing can be tricky, `parser` aims to be as fast as possible without compromising safety or error handling. Every parser function has a benchmark and has been written with performance in mind, none of them allocate on the heap ⚡️
- **Correct:** You get the correct behaviour at all times, on any valid UTF-8 text. Errors are well handled and reported for easy debugging. 100% test coverage.
- **Intuitive:** Some parser combinator libraries are tricky to wrap your head around, I want `parser` to be super simple to use so that anyone can pick it up and be productive quickly
- **Well Documented:** Every combinator in `parser` has a comprehensive doc comment describing it's entire behaviour, as well as an executable example of its use

## Installation

```shell
go get github.com/FollowTheProcess/parser@latest
```

## Quickstart

### Credits

This package was created with [copier] and the [FollowTheProcess/go_copier] project template.

It is also heavily inspired by [nom], an excellent combinatorial parsing library written in [Rust].

[copier]: https://copier.readthedocs.io/en/stable/
[FollowTheProcess/go_copier]: https://github.com/FollowTheProcess/go_copier
[combinatorial parsing]: https://en.wikipedia.org/wiki/Parser_combinator
[nom]: https://github.com/rust-bakery/nom
[Rust]: https://www.rust-lang.org
