package main

import (
	"os"

	"gnd.la/util/pathutil"
)

var (
	goRoot = pathutil.Relative("tmp/dist/go")
	goPath = pathutil.Relative("tmp/go")
)

func init() {
	if err := os.MkdirAll(goPath, 0755); err != nil {
		panic(err)
	}
}
