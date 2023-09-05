package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/ghostiam/protogolint"
)

func main() {
	singlechecker.Main(protogolint.Analyzer)
}
