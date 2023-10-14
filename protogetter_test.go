package protogetter_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/ghostiam/protogetter"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, protogetter.NewAnalyzer())

	analysistest.Run(t, testdata, protogetter.NewAnalyzer(), "./proto/...")
}
