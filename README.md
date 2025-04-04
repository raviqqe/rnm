# rnm

[![GitHub Action](https://img.shields.io/github/actions/workflow/status/raviqqe/rnm/test.yaml?branch=main&style=flat-square)](https://github.com/raviqqe/rnm/actions)
[![Codecov](https://img.shields.io/codecov/c/github/raviqqe/rnm.svg?style=flat-square)](https://codecov.io/gh/raviqqe/rnm)
[![License](https://img.shields.io/github/license/raviqqe/rnm.svg?style=flat-square)](LICENSE)

Yet another [`fastmod`](https://github.com/facebookincubator/fastmod) alternative.

Replace all occurrences of a name to another name in `{camel,kebab,shout,snake,...}`cases in your codes!

## Features

- Support for different case styles
  - See `rnm --help` to list them all.
- Automatic pluralization
- File renaming
- Massive speed

## Install

```sh
GO111MODULE=on go get -u github.com/raviqqe/rnm
```

## Usage

```sh
rnm 'foo bar' 'baz qux'
```

For more information, see `rnm --help`.

## Examples

Given a file named `foo_bar.go`:

```go
const FOO_BAR = 42

type FooBar struct {
  fooBar int
}

func (f FooBar) fooBar() {
  println("foo bar")
}
```

When you run `rnm 'foo bar' 'baz qux'`, you would see a file named `baz_qux.go` with contents:

```go
const BAZ_QUX = 42

type BazQux struct {
  bazQux int
}

func (f BazQux) bazQux() {
  println("baz qux")
}
```

## License

[MIT](LICENSE)
