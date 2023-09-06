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

var Analyzer = &analysis.Analyzer{
	Name:     "protogolint",
	Doc:      "reports direct reads from proto message fields when getters should be used",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
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
	ins.Preorder(nodeTypes, check(pass))

	return nil, nil
}

func isGeneratedFile(f *ast.File) bool {
	for _, c := range f.Comments {
		if strings.HasPrefix(c.Text(), "Code generated") {
			return true
		}
	}

	return false
}

func check(pass *analysis.Pass) func(ast.Node) {
	skippedPos := map[token.Pos]struct{}{}
	isAlreadyReplaced := map[string]map[int][2]int{} // map[filename][line][start, end]

	return func(n ast.Node) {
		// fmt.Printf("\n>>> check: %s\n", formatNode(n))
		if _, ok := skippedPos[n.Pos()]; ok {
			// fmt.Printf(">>> ignored\n")
			return
		}

		a := NewChecker(pass)

		switch x := n.(type) {
		case *ast.AssignStmt:
			// Skip any assignment to the field.
			for _, lhs := range x.Lhs {
				skippedPos[lhs.Pos()] = struct{}{}
			}

		case *ast.IncDecStmt:
			// Skip any increment/decrement to the field.
			skippedPos[x.X.Pos()] = struct{}{}

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

						skippedPos[ue.X.Pos()] = struct{}{}
					}

					return
				}

				a.Check(x)

			default:
				if !isProtoMessage(pass, x.Fun) {
					return
				}

				a.SetError(fmt.Errorf("CallExpr: not implemented for type: %s (%s)", reflect.TypeOf(f), formatNode(n)))
			}

		case *ast.SelectorExpr:
			if !isProtoMessage(pass, x.X) {
				// If the selector is not on a proto message, skip it.
				return
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

			return
		}

		// If existing in skippedPos, skip it.
		if _, ok := skippedPos[n.Pos()]; ok {
			return
		}

		// If from and to are the same, skip it.
		if result.From == result.To {
			return
		}

		// Calculate the position of the expression in the file and if it has already been replaced, skip it.
		{
			filePos := pass.Fset.Position(n.Pos())
			fileEnd := pass.Fset.Position(n.End())

			arf, ok := isAlreadyReplaced[filePos.Filename]
			if !ok {
				arf = make(map[int][2]int)
				isAlreadyReplaced[filePos.Filename] = arf
			}

			arfl, ok := arf[filePos.Line]
			if !ok {
				arf[filePos.Line] = [2]int{filePos.Offset, fileEnd.Offset}
			} else {
				if arfl[0] <= filePos.Offset && fileEnd.Offset <= arfl[1] {
					return
				}
				arf[filePos.Line] = [2]int{filePos.Offset, fileEnd.Offset}
			}
		}

		pass.Report(analysis.Diagnostic{
			Pos:     n.Pos(),
			End:     n.End(),
			Message: fmt.Sprintf(`proto field read without getter: %q should be %q`, result.From, result.To),
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

		if c.methodExists(x.X, x.Sel.Name) {
			// If the method has already been called, leave it as is.
			c.write(x.Sel.Name)
			return
		}

		// If getter exists, use it.
		if c.methodExists(x.X, "Get"+x.Sel.Name) {
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
		c.write(" ")
		c.write(x.Op.String())
		c.write(" ")
		c.Check(x.Y)

	default:
		c.err = fmt.Errorf("checker not implemented for type: %s", reflect.TypeOf(x))
	}
}

func (c *Checker) methodExists(x ast.Expr, name string) bool {
	methods := getMethods(c.info, x)

	for _, m := range methods {
		if m == name {
			return true
		}
	}

	return false
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

func getMethods(info *types.Info, x ast.Expr) []string {
	if info == nil {
		return nil
	}

	t := info.TypeOf(x)
	if t == nil {
		return nil
	}

	ptr, ok := t.Underlying().(*types.Pointer)
	if ok {
		t = ptr.Elem()
	}

	named, ok := t.(*types.Named)
	if !ok {
		return nil
	}

	methods := make([]string, 0, named.NumMethods())
	for i := 0; i < named.NumMethods(); i++ {
		methods = append(methods, named.Method(i).Name())
	}

	return methods
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
