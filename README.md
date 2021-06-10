# logman

:smile: golang library to organize log files (create by date, GC, etc)

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/moul.io/logman)
[![License](https://img.shields.io/badge/license-Apache--2.0%20%2F%20MIT-%2397ca00.svg)](https://github.com/moul/logman/blob/main/COPYRIGHT)
[![GitHub release](https://img.shields.io/github/release/moul/logman.svg)](https://github.com/moul/logman/releases)
[![Docker Metrics](https://images.microbadger.com/badges/image/moul/logman.svg)](https://microbadger.com/images/moul/logman)
[![Made by Manfred Touron](https://img.shields.io/badge/made%20by-Manfred%20Touron-blue.svg?style=flat)](https://manfred.life/)

[![Go](https://github.com/moul/logman/workflows/Go/badge.svg)](https://github.com/moul/logman/actions?query=workflow%3AGo)
[![Release](https://github.com/moul/logman/workflows/Release/badge.svg)](https://github.com/moul/logman/actions?query=workflow%3ARelease)
[![PR](https://github.com/moul/logman/workflows/PR/badge.svg)](https://github.com/moul/logman/actions?query=workflow%3APR)
[![GolangCI](https://golangci.com/badges/github.com/moul/logman.svg)](https://golangci.com/r/github.com/moul/logman)
[![codecov](https://codecov.io/gh/moul/logman/branch/main/graph/badge.svg)](https://codecov.io/gh/moul/logman)
[![Go Report Card](https://goreportcard.com/badge/moul.io/logman)](https://goreportcard.com/report/moul.io/logman)
[![CodeFactor](https://www.codefactor.io/repository/github/moul/logman/badge)](https://www.codefactor.io/repository/github/moul/logman)

[![Gitpod ready-to-code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/moul/logman)

## Example

[embedmd]:# (example_test.go /import\ / $)
```go
import (
	"fmt"

	"moul.io/logman"
)

func Example() {
	// new log manager
	manager := logman.Manager{
		Path:     "./path/to/dir",
		MaxFiles: 10,
	}

	// cleanup old log files for a specific app name
	err := manager.Flush("my-app")
	checkErr(err)

	// cleanup old log files for any app sharing this log directory
	err = manager.FlushAll()
	checkErr(err)

	// list existing log files
	files, err := manager.Files()
	checkErr(err)
	fmt.Println(files)

	// - create an WriteCloser
	// - automatically delete old log files if it hits a limit
	writer, err := manager.New("my-app")
	checkErr(err)
	defer writer.Close()
	writer.Write([]byte("hello world!\n"))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
```

## Usage

[embedmd]:# (.tmp/godoc.txt txt /TYPES/ $)
```txt
TYPES

type File struct {
	// Full path.
	Path string

	// Size in bytes.
	Size int64

	// Provided name when creating the file.
	Name string

	// Creation date of the file.
	Time time.Time

	// Whether it is the most recent log file for the provided app name or not.
	Latest bool

	// If there were errors when trying to get info about this file.
	Errs error `json:"Errs,omitempty"`
}
    File defines a log file with metadata.

func (f File) String() string
    String implements Stringer.

type Manager struct {
	// Path is the target directory containing the log files.
	// Default is '.'.
	Path string

	// MaxFiles is the maximum number of log files in the directory.
	// If 0, won't automatically GC based on this criteria.
	MaxFiles int
}
    Manager is a configuration object used to create log files with automatic GC
    rules.

func (m Manager) Files() ([]File, error)
    Files returns a list of existing log files.

func (m Manager) Flush(name string) error
    Flush deletes old log files for the specified app name.

func (m Manager) FlushAll() error
    FlushAll deletes old log files for any app name.

func (m Manager) New(name string) (io.WriteCloser, error)
    Create a new log file and perform automatic GC of the old log files if
    needed.

    The created log file will looks like:

        <path/to/log/dir>/<name>-<time>.log

    Depending on the provided configuration of Manager, an automatic GC will be
    run automatically.

```

## Install

### Using go

```sh
go get moul.io/logman
```

### Releases

See https://github.com/moul/logman/releases

## Contribute

![Contribute <3](https://raw.githubusercontent.com/moul/moul/main/contribute.gif)

I really welcome contributions.
Your input is the most precious material.
I'm well aware of that and I thank you in advance.
Everyone is encouraged to look at what they can do on their own scale;
no effort is too small.

Everything on contribution is sum up here: [CONTRIBUTING.md](./CONTRIBUTING.md)

### Dev helpers

Pre-commit script for install: https://pre-commit.com

### Contributors ✨

<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->
[![All Contributors](https://img.shields.io/badge/all_contributors-2-orange.svg)](#contributors)
<!-- ALL-CONTRIBUTORS-BADGE:END -->

Thanks goes to these wonderful people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tr>
    <td align="center"><a href="http://manfred.life"><img src="https://avatars1.githubusercontent.com/u/94029?v=4" width="100px;" alt=""/><br /><sub><b>Manfred Touron</b></sub></a><br /><a href="#maintenance-moul" title="Maintenance">🚧</a> <a href="https://github.com/moul/logman/commits?author=moul" title="Documentation">📖</a> <a href="https://github.com/moul/logman/commits?author=moul" title="Tests">⚠️</a> <a href="https://github.com/moul/logman/commits?author=moul" title="Code">💻</a></td>
    <td align="center"><a href="https://manfred.life/moul-bot"><img src="https://avatars1.githubusercontent.com/u/41326314?v=4" width="100px;" alt=""/><br /><sub><b>moul-bot</b></sub></a><br /><a href="#maintenance-moul-bot" title="Maintenance">🚧</a></td>
  </tr>
</table>

<!-- markdownlint-enable -->
<!-- prettier-ignore-end -->
<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the [all-contributors](https://github.com/all-contributors/all-contributors)
specification. Contributions of any kind welcome!

### Stargazers over time

[![Stargazers over time](https://starchart.cc/moul/logman.svg)](https://starchart.cc/moul/logman)

## License

© 2021   [Manfred Touron](https://manfred.life)

Licensed under the [Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0)
([`LICENSE-APACHE`](LICENSE-APACHE)) or the [MIT license](https://opensource.org/licenses/MIT)
([`LICENSE-MIT`](LICENSE-MIT)), at your option.
See the [`COPYRIGHT`](COPYRIGHT) file for more details.

`SPDX-License-Identifier: (Apache-2.0 OR MIT)`
