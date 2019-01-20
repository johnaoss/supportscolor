# supports-color

## Features

This package provides a simple interface to query the level of support found in a given terminal.

* Respects the `--color` flags, and additional colour specifications.
* Handles TTY / Cygwin colour levels
* Supports popular CI pipelines like TravisCI & CircleCI.

## Roadmap

* Check popular terminal support schemes (urvxt, xterm-256colour flag, etc.)
* Proper unit tests (if possible)
* Remove external dependency

## Usage

```go

import supportscolour

func main() {
    col := supportscolour.GetSupportLevel()
    if col.Has1m {
        // do something with truecolour
    } else if col.Has256 {
        // do something with 256 colours
    } else {
        // print out standard 16 colours.
    }
}

```

## Install

```bash
go get -u github.com/johnaoss/supports-color
```

## Contributing

Due to the nature of this problem, pull requests are appreciated if you find a terminal in which the data reported is inaccurate.
I've currently only had access to a iTerm2 terminal, and as such I've only been able to test on the single platform.

## License

This is licensed with the MIT License, and is found in the LICENSE file.

## Thanks

Very much directly inspired by the npm module of the same name: [github.com/chalk/supports-color](https://github.com/chalk/supports-color) !
