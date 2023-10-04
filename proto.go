package protogolint

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"log"
	"os"
	"reflect"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

type Mode int

const (
	StandaloneMode Mode = iota
	GolangciLintMode
)

func NewAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "protogolint",
		Doc:  "Reports direct reads from proto message fields when getters should be used",
		Run: func(pass *analysis.Pass) (any, error) {
			Run(pass, StandaloneMode)
			return nil, nil
		},
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

func Run(pass *analysis.Pass, mode Mode) []Issue {
	nodeTypes := []ast.Node{
		(*ast.AssignStmt)(nil),
		(*ast.CallExpr)(nil),
		(*ast.SelectorExpr)(nil),
		(*ast.IncDecStmt)(nil),
	}

	// Skip generated files.
	var files []*ast.File
	for _, f := range pass.Files {
		if !isGeneratedFile(f) {
			files = append(files, f)

			// ast.Print(pass.Fset, f)
		}
	}

	ins := inspector.New(files)
	checker := check(pass, mode)

	var issues []Issue

	ins.Preorder(nodeTypes, func(node ast.Node) {
		issue := checker(node)
		if issue != nil {
			issues = append(issues, *issue)
		}
	})

	return issues
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

func isGeneratedFile(f *ast.File) bool {
	for _, c := range f.Comments {
		if strings.HasPrefix(c.Text(), "Code generated") {
			return true
		}
	}

	return false
}

func check(pass *analysis.Pass, mode Mode) func(ast.Node) *Issue {
	filter := newPosFilter()

	return func(n ast.Node) *Issue {
		// fmt.Printf("\n>>> check: %s\n", formatNode(n))
		if filter.IsFiltered(n.Pos()) {
			// fmt.Printf(">>> filtered\n")
			return nil
		}

		a := NewChecker(pass)

		switch x := n.(type) {
		case *ast.AssignStmt:
			// Skip any assignment to the field.
			for _, lhs := range x.Lhs {
				filter.AddPos(lhs.Pos())
			}

		case *ast.IncDecStmt:
			// Skip any increment/decrement to the field.
			filter.AddPos(x.X.Pos())

		case *ast.CallExpr:
			switch f := x.Fun.(type) {
			case *ast.SelectorExpr:
				if !isProtoMessage(pass, f.X) {
					for _, arg := range x.Args {
						// Skip all expressions when the function points to a field, for example somefunc(&t).
						// Because this is not direct reading, but most likely writing by pointer (for example like sql.Scan).
						ue, ok := arg.(*ast.UnaryExpr)
						if !ok || ue.Op != token.AND {
							continue
						}

						filter.AddPos(ue.X.Pos())
					}

					return nil
				}

				a.Check(x)

			default:
				if !isProtoMessage(pass, x.Fun) {
					return nil
				}

				a.SetError(fmt.Errorf("CallExpr: not implemented for type: %s (%s)", reflect.TypeOf(f), formatNode(n)))
			}

		case *ast.SelectorExpr:
			if !isProtoMessage(pass, x.X) {
				// If the selector is not on a proto message, skip it.
				return nil
			}

			a.Check(x)

		default:
			a.SetError(fmt.Errorf("not implemented for type: %s (%s)", reflect.TypeOf(x), formatNode(n)))
		}

		result, err := a.Result()

		// fmt.Printf(">>> check: res: %v, err: %v\n", result, err)

		if err != nil {
			pass.Report(analysis.Diagnostic{
				Pos:     n.Pos(),
				End:     n.End(),
				Message: fmt.Sprintf("error: %v", err),
			})

			return nil
		}

		// If existing in filter, skip it.
		if filter.IsFiltered(n.Pos()) {
			return nil
		}

		// If from and to are the same, skip it.
		if result.From == result.To {
			return nil
		}

		// If the expression has already been replaced, skip it.
		if filter.IsAlreadyReplaced(pass, n.Pos(), n.End()) {
			return nil
		}
		// Add the expression to the filter.
		filter.AddAlreadyReplaced(pass, n.Pos(), n.End())

		msg := fmt.Sprintf(`proto field read without getter: %q should be %q`, result.From, result.To)

		switch mode {
		case StandaloneMode:
			pass.Report(analysis.Diagnostic{
				Pos:     n.Pos(),
				End:     n.End(),
				Message: msg,
				SuggestedFixes: []analysis.SuggestedFix{
					{
						Message: fmt.Sprintf("%q should be replaced with %q", result.From, result.To),
						TextEdits: []analysis.TextEdit{
							{
								Pos:     n.Pos(),
								End:     n.End(),
								NewText: []byte(result.To),
							},
						},
					},
				},
			})

		case GolangciLintMode:
			return &Issue{
				Pos:     pass.Fset.Position(n.Pos()),
				Message: msg,
				InlineFix: InlineFix{
					StartCol:  pass.Fset.Position(n.Pos()).Column - 1,
					Length:    len(result.From),
					NewString: result.To,
				},
			}
		}

		return nil
	}
}

type Checker struct {
	info *types.Info

	to   strings.Builder
	from strings.Builder
	err  error
}

// NewChecker creates a new Checker instance.
func NewChecker(pass *analysis.Pass) *Checker {
	return &Checker{
		info: pass.TypesInfo,
	}
}

func (c *Checker) SetError(err error) {
	c.err = err
}

func (c *Checker) Result() (*CheckerResult, error) {
	if c.err != nil {
		return nil, c.err
	}

	return &CheckerResult{
		From: c.from.String(),
		To:   c.to.String(),
	}, nil
}

func (c *Checker) Check(expr ast.Expr) {
	switch x := expr.(type) {
	case *ast.Ident:
		c.write(x.Name)

	case *ast.BasicLit:
		c.write(x.Value)

	case *ast.SelectorExpr:
		c.Check(x.X)
		c.write(".")

		if methodIsExists(c.info, x.X, x.Sel.Name) {
			// If the method has already been called, leave it as is.
			c.write(x.Sel.Name)
			return
		}

		// If getter exists, use it.
		if methodIsExists(c.info, x.X, "Get"+x.Sel.Name) {
			c.writeFrom(x.Sel.Name)
			c.writeTo("Get" + x.Sel.Name + "()")
			return
		}

		// If method does not exist, leave it as is.
		c.write(x.Sel.Name)

	case *ast.CallExpr:
		c.Check(x.Fun)
		c.write("(")
		for i, arg := range x.Args {
			if i > 0 {
				c.write(",")
			}
			c.Check(arg)
		}
		c.write(")")

	case *ast.IndexExpr:
		c.Check(x.X)
		c.write("[")
		c.Check(x.Index)
		c.write("]")

	case *ast.BinaryExpr:
		c.Check(x.X)
		c.write(x.Op.String())
		c.Check(x.Y)

	default:
		c.err = fmt.Errorf("checker not implemented for type: %s", reflect.TypeOf(x))
	}
}

func (c *Checker) write(s string) {
	c.writeTo(s)
	c.writeFrom(s)
}

func (c *Checker) writeTo(s string) {
	c.to.WriteString(s)
}

func (c *Checker) writeFrom(s string) {
	c.from.WriteString(s)
}

// CheckerResult contains source code (from) and suggested change (to)
type CheckerResult struct {
	From string
	To   string
}

type posFilter struct {
	positions       map[token.Pos]struct{}
	alreadyReplaced map[string]map[int][2]int // map[filename][line][start, end]
}

func newPosFilter() *posFilter {
	return &posFilter{
		positions:       make(map[token.Pos]struct{}),
		alreadyReplaced: make(map[string]map[int][2]int),
	}
}

func (f *posFilter) IsFiltered(pos token.Pos) bool {
	_, ok := f.positions[pos]
	return ok
}

func (f *posFilter) AddPos(pos token.Pos) {
	f.positions[pos] = struct{}{}
}

func (f *posFilter) IsAlreadyReplaced(pass *analysis.Pass, pos token.Pos, end token.Pos) bool {
	filePos := pass.Fset.Position(pos)
	fileEnd := pass.Fset.Position(end)

	lines, ok := f.alreadyReplaced[filePos.Filename]
	if !ok {
		return false
	}

	lineRange, ok := lines[filePos.Line]
	if !ok {
		return false
	}

	if lineRange[0] <= filePos.Offset && fileEnd.Offset <= lineRange[1] {
		return true
	}

	return false
}

func (f *posFilter) AddAlreadyReplaced(pass *analysis.Pass, pos token.Pos, end token.Pos) {
	filePos := pass.Fset.Position(pos)
	fileEnd := pass.Fset.Position(end)

	lines, ok := f.alreadyReplaced[filePos.Filename]
	if !ok {
		lines = make(map[int][2]int)
		f.alreadyReplaced[filePos.Filename] = lines
	}

	lineRange, ok := lines[filePos.Line]
	if ok && lineRange[0] <= filePos.Offset && fileEnd.Offset <= lineRange[1] {
		return
	}

	lines[filePos.Line] = [2]int{filePos.Offset, fileEnd.Offset}
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

func methodIsExists(info *types.Info, x ast.Expr, name string) bool {
	if info == nil {
		return false
	}

	t := info.TypeOf(x)
	if t == nil {
		return false
	}

	ptr, ok := t.Underlying().(*types.Pointer)
	if ok {
		t = ptr.Elem()
	}

	named, ok := t.(*types.Named)
	if !ok {
		return false
	}

	for i := 0; i < named.NumMethods(); i++ {
		if named.Method(i).Name() == name {
			return true
		}
	}

	return false
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
