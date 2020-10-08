# rnm

[![GitHub Action](https://img.shields.io/github/workflow/status/raviqqe/rnm/test?style=flat-square)](https://github.com/raviqqe/rnm/actions)
[![Codecov](https://img.shields.io/codecov/c/github/raviqqe/rnm.svg?style=flat-square)](https://codecov.io/gh/raviqqe/rnm)
[![License](https://img.shields.io/github/license/raviqqe/rnm.svg?style=flat-square)](LICENSE)

Rename all names in your code!

## Install

```
GO111MODULE=on go get -u github.com/raviqqe/rnm
```

## Usage

```
rnm 'foo bar' 'baz qux'
```

For more information, see `rnm --help`.

## Examples

Given a file named `foo_bar.go`:

```go
const FOO_BAR = 42;

type FooBar {
  fooBar int
}

func (f FooBar) fooBar() {
  println("foo bar")
}
```

When you run `rnm 'foo bar' 'baz qux'`, you would see a file named `baz_qux.go` with contents:

```go
const BAZ_QUX = 42;

type BazQux {
  bazQux int
}

func (f BazQux) bazQux() {
  println("baz qux")
}
```

## License

[MIT](LICENSE)
