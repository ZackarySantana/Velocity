package vets

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var CommandNoFmt = &analysis.Analyzer{
	Name: "CommandNoFmt",
	Doc:  "check for fmt.Println or fmt.Print in cli.Command structs",
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
	Run: func(pass *analysis.Pass) (interface{}, error) {
		inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
		nodeFilter := []ast.Node{
			(*ast.CompositeLit)(nil),
		}

		inspect.Preorder(nodeFilter, func(n ast.Node) {
			cl := n.(*ast.CompositeLit)
			if typeIdent, ok := cl.Type.(*ast.SelectorExpr); ok {
				if typeIdent.Sel.Name == "Command" {
					for _, elt := range cl.Elts {
						if kv, ok := elt.(*ast.KeyValueExpr); ok {
							if fn, ok := kv.Value.(*ast.FuncLit); ok {
								checkFuncBody(pass, fn.Body)
							}
						}
					}
				}
			}
		})

		return nil, nil
	},
}

func checkFuncBody(pass *analysis.Pass, body *ast.BlockStmt) {
	for _, stmt := range body.List {
		switch s := stmt.(type) {
		case *ast.ExprStmt:
			if call, ok := s.X.(*ast.CallExpr); ok {
				if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
					if pkg, ok := fun.X.(*ast.Ident); ok && pkg.Name == "fmt" {
						if fun.Sel.Name == "Println" || fun.Sel.Name == "Print" {
							pass.Reportf(call.Pos(), "cli.Command struct contains call to fmt.%s", fun.Sel.Name)
						}
					}
				}
			}
		case *ast.IfStmt:
			checkFuncBody(pass, s.Body)
			if s.Else != nil {
				if elseBody, ok := s.Else.(*ast.BlockStmt); ok {
					checkFuncBody(pass, elseBody)
				}
			}
		case *ast.BlockStmt:
			checkFuncBody(pass, s)
		case *ast.ForStmt:
			checkFuncBody(pass, s.Body)
		case *ast.RangeStmt:
			checkFuncBody(pass, s.Body)
		case *ast.SwitchStmt:
			checkFuncBody(pass, s.Body)
		case *ast.TypeSwitchStmt:
			checkFuncBody(pass, s.Body)
		case *ast.SelectStmt:
			checkFuncBody(pass, s.Body)
		}
	}
}
