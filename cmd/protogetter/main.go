package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/ghostiam/protogetter"
)

func main() {
	cfg := &protogetter.Config{
		Mode: protogetter.StandaloneMode,
	}

	singlechecker.Main(protogetter.NewAnalyzer(cfg))
}
