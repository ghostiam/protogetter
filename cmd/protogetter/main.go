package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/ghostiam/protogetter"
)

func main() {
	singlechecker.Main(protogetter.NewAnalyzer(nil))
}
