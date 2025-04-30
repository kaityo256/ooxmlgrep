# ooxmlgrep

`ooxmlgrep` is a command-line tool written in Go for recursively searching text inside `.pptx` (PowerPoint OOXML) files.  
It extracts slide text and performs grep-like matching, with optional highlighting, case-insensitive search, and slide-level reporting.

## Features

- 🔍 Recursively search `.pptx` files for a keyword
- 📄 Print matched slide text, including the slide number and filename
- 🎯 Highlight matched terms in red (ANSI color)
- 🧾 Support for `--ignore-case`, `--number`, and `--version` flags
- ⚡ Fast, dependency-free, cross-platform

## Installation

### From Source (requires Go 1.17+)

```sh
go install github.com/kaityo256/ooxmlgrep@latest
```

Make sure your `$GOPATH/bin` or `$HOME/go/bin` is in your `$PATH`.

### From Prebuilt Binaries

You can download prebuilt binaries for macOS, Linux, and Windows from the [Releases](https://github.com/kaityo256/ooxmlgrep/releases) page.

## Usage

```sh
ooxmlgrep [options] <keyword>
```

### Options

| Flag              | Description                          |
|-------------------|--------------------------------------|
| `-n`, `--number`        | Show slide number only                 |
| `-i`, `--ignore-case`   | Case-insensitive match                |
| `--version`             | Show version number                   |

### Examples

```sh
ooxmlgrep machine        # search for 'machine' (case-sensitive)
ooxmlgrep -i AI          # case-insensitive search for 'AI'
ooxmlgrep -n keyword     # show only filenames and slide numbers
```

## Output Format

By default:

```
./slides/sample.pptx:slide 2:Deep learning is a kind of machine learning.
```

With `-n`:

```
./slides/sample.pptx:slide 2
```

## License

MIT License. See [LICENSE](./LICENSE) for details.

---

© 2025 H. Watanabe
