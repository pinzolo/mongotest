# mongotest

[![Build Status](https://travis-ci.org/pinzolo/mongotest.png)](http://travis-ci.org/pinzolo/mongotest)
[![Coverage Status](https://coveralls.io/repos/github/pinzolo/mongotest/badge.svg?branch=master)](https://coveralls.io/github/pinzolo/mongotest?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/pinzolo/mongotest)](https://goreportcard.com/report/github.com/pinzolo/mongotest)
[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/pinzolo/mongotest)
[![license](http://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/pinzolo/mongotest/master/LICENSE)

## Description

mongotest is testing helper for using MongoDB.

## Example

```go
package main

import (
	"fmt"

	"github.com/pinzolo/mongotest"
)

func main() {
	mongotest.Configure(mongotest.Config{
		URL:            "mongodb://root:password@127.0.0.1:27017",
		Database:       "mongotest",
		FixtureRootDir: "testdata",
		FixtureFormat:  mongotest.FixtureFormatJSON,
		PreInsertFuncs: []mongotest.PreInsertFunc{
			mongotest.SimpleConvertTime("users", "created_at"),
		},
	})

	// 1. Read testdata/json/admin_users.json and testdata/json/foo_users.json
	// 2. Merge read data
	// 3. Drop collection and insert read data
	err := mongotest.UseFixture("json/admin_users", "json/foo_users")
	if err != nil {
		panic(err)
	}

	// Count is helper function.
	// mongotest has some useful helper functions.
	n, err := mongotest.Count("users")
	if err != nil {
		panic(err)
	}
	fmt.Println(n)
}
```

## Install

```bash
$ go get github.com/pinzolo/mongotest
```

## Contribution

1. Fork ([https://github.com/pinzolo/mongotest/fork](https://github.com/pinzolo/mongotest/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[pinzolo](https://github.com/pinzolo)
