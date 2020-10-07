# rnm

[![GitHub Action](https://img.shields.io/github/workflow/status/raviqqe/rnm/test?style=flat-square)](https://github.com/raviqqe/rnm/actions)
[![Codecov](https://img.shields.io/codecov/c/github/raviqqe/rnm.svg?style=flat-square)](https://codecov.io/gh/raviqqe/rnm)
[![Go Report Card](https://goreportcard.com/badge/github.com/raviqqe/rnm?style=flat-square)](https://goreportcard.com/report/github.com/raviqqe/rnm)
[![License](https://img.shields.io/github/license/raviqqe/rnm.svg?style=flat-square)](LICENSE)

Rename any names in your code!

## Install

```
GO111MODULE=on go get -u github.com/raviqqe/rnm
```

## Usage

```
rnm FooBar BazQux
```

For more information, see `rnm --help`.

## Examples

Given a file named `foo.ts`:

```typescript
const FOO_BAR: number = 42;

export class FooBar {
  public fooBar() {
    console.log("foo bar");
  }
}
```

When you run `rnm FooBar BazQux`, you would see the file with contents:

```typescript
const BAZ_QUX: number = 42;

export class BazQux {
  public bazQux() {
    console.log("baz qux");
  }
}
```

## License

[MIT](LICENSE)
