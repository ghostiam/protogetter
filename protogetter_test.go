package protogetter_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/ghostiam/protogetter"
)

func Test(t *testing.T) {
	cfg := &protogetter.Config{}

	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, protogetter.NewAnalyzer(cfg))

	analysistest.Run(t, testdata, protogetter.NewAnalyzer(cfg), "./proto/...")
}
