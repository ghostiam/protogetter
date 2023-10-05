package protogetter

import (
	"fmt"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

type Report struct {
	analysis.Range
	From, To     string
	SelectorEdit Edit
}

type Edit struct {
	analysis.Range
	From, To string
}

func (r Report) ToAnalysisDiagnostic() analysis.Diagnostic {
	msg := fmt.Sprintf(msgFormat, r.From, r.To)

	return analysis.Diagnostic{
		Pos:     r.Pos(),
		End:     r.End(),
		Message: msg,
		SuggestedFixes: []analysis.SuggestedFix{
			{
				Message: msg,
				TextEdits: []analysis.TextEdit{
					{
						Pos:     r.SelectorEdit.Pos(),
						End:     r.SelectorEdit.End(),
						NewText: []byte(r.SelectorEdit.To),
					},
				},
			},
		},
	}
}

func (r Report) ToIssue(fset *token.FileSet) Issue {
	msg := fmt.Sprintf(msgFormat, r.From, r.To)
	return Issue{
		Pos:     fset.Position(r.Pos()),
		Message: msg,
		InlineFix: InlineFix{
			StartCol:  fset.Position(r.SelectorEdit.Pos()).Column - 1,
			Length:    len(r.SelectorEdit.From),
			NewString: r.SelectorEdit.To,
		},
	}
}

// Issue is used to integrate with golangci-lint's inline auto fix.
type Issue struct {
	Pos       token.Position
	Message   string
	InlineFix InlineFix
}

type InlineFix struct {
	StartCol  int // zero-based
	Length    int
	NewString string
}
