package protogolint_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/ghostiam/protogolint"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, protogolint.Analyzer)
}
