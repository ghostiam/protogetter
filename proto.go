package protogolint

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"log"
	"os"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "protogolint",
	Doc:      "reports direct reads from proto message fields when getters should be used",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	ins, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, errors.New("analyzer is not type *inspector.Inspector")
	}

	nodeTypes := []ast.Node{
		(*ast.AssignStmt)(nil),
		(*ast.CallExpr)(nil),
		(*ast.SelectorExpr)(nil),
	}

	ins.Preorder(nodeTypes, check(pass))

	return nil, nil
}

func check(pass *analysis.Pass) func(ast.Node) {
	return func(n ast.Node) {
		ignores := map[token.Pos]struct{}{}

		var oldExpr, newExpr string

		switch x := n.(type) {
		case *ast.AssignStmt:
			for _, lhs := range x.Lhs {
				ignores[lhs.Pos()] = struct{}{}
			}

		case *ast.CallExpr:
			switch f := x.Fun.(type) {
			case *ast.SelectorExpr:
				if !isProtoMessage(pass, f.X) {
					ignores[x.Pos()] = struct{}{}
					return
				}

				// TODO
				// oldExpr, newExpr = handleExpr(pass.TypesInfo, x, x)
				// oldExpr, newExpr = makeFromCallAndSelectorExpr(pass.TypesInfo, x)

			default:
				for _, arg := range x.Args {
					a, ok := arg.(*ast.UnaryExpr)
					if !ok || a.Op != token.AND {
						continue
					}

					ignores[a.X.Pos()] = struct{}{}
				}
			}

			f, ok := x.Fun.(*ast.SelectorExpr)
			if !ok || !isProtoMessage(pass, f.X) {
				for _, arg := range x.Args {
					var a *ast.UnaryExpr
					a, ok = arg.(*ast.UnaryExpr)
					if !ok || a.Op != token.AND {
						continue
					}

					ignores[a.X.Pos()] = struct{}{}
				}

				ignores[x.Pos()] = struct{}{}
				return
			}

		case *ast.SelectorExpr:
			if !isProtoMessage(pass, x.X) {
				return
			}

			oldExpr, newExpr = handleExpr(pass.TypesInfo, x, x.X)
		}

		if _, ok := ignores[n.Pos()]; ok {
			return
		}

		if oldExpr == "" || newExpr == "" ||
			oldExpr == newExpr {
			return
		}
		ignores[n.Pos()] = struct{}{}

		pass.Report(analysis.Diagnostic{
			Pos:     n.Pos(),
			End:     n.End(),
			Message: fmt.Sprintf(`proto field read without getter: %q should be %q`, oldExpr, newExpr),
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message: fmt.Sprintf("%q should be replaced with %q", oldExpr, newExpr),
					TextEdits: []analysis.TextEdit{
						{
							Pos:     n.Pos(),
							End:     n.End(),
							NewText: []byte(newExpr),
						},
					},
				},
			},
		})
	}
}

func handleExpr(info *types.Info, base, child ast.Expr) (string, string) {
	// TODO
	printAST("base", base)
	return "", ""
}

const messageState = "google.golang.org/protobuf/internal/impl.MessageState"

func isProtoMessage(pass *analysis.Pass, expr ast.Expr) bool {
	t := pass.TypesInfo.TypeOf(expr)
	if t == nil {
		return false
	}
	ptr, ok := t.Underlying().(*types.Pointer)
	if !ok {
		return false
	}
	named, ok := ptr.Elem().(*types.Named)
	if !ok {
		return false
	}
	sct, ok := named.Underlying().(*types.Struct)
	if !ok {
		return false
	}
	if sct.NumFields() == 0 {
		return false
	}

	return sct.Field(0).Type().String() == messageState
}

func formatNode(node ast.Node) string {
	buf := new(bytes.Buffer)
	if err := format.Node(buf, token.NewFileSet(), node); err != nil {
		log.Printf("Error formatting expression: %v", err)
		return ""
	}

	return buf.String()
}

func printAST(msg string, node ast.Node) {
	fmt.Printf(">>> %s:\n%s\n\n\n", msg, formatNode(node))
	ast.Fprint(os.Stdout, nil, node, nil)
	fmt.Println("--------------")
}
