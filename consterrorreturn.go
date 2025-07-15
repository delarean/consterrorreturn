package consterrorreturn

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "consterrorreturn",
	Doc:  "returning sentinel (constant) error instead of propagating original err variable",
	Run:  run,
}

const sentinelErrMsg = "returning sentinel (constant) error instead of propagating original err variable"

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			if retStmt, ok := n.(*ast.ReturnStmt); ok {
				for _, retExpr := range retStmt.Results {
					if insideIfErrorsIs(pass, retStmt) {
						continue
					}

					typ := pass.TypesInfo.TypeOf(retExpr)
					if typ == nil {
						continue
					}
					if isErrorType(typ, pass) {
						switch expr := retExpr.(type) {
						case *ast.SelectorExpr:
							pass.Reportf(expr.Pos(), sentinelErrMsg)
						case *ast.Ident:
							obj := pass.TypesInfo.ObjectOf(expr)
							if obj != nil && isConstant(obj) {
								pass.Reportf(expr.Pos(), sentinelErrMsg)
							}
						}
					}
				}
			}

			if call, ok := n.(*ast.CallExpr); ok {
				if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
					if pkgIdent, ok := sel.X.(*ast.Ident); ok && pkgIdent.Name == "fmt" && sel.Sel.Name == "Errorf" {
						if len(call.Args) >= 2 {
							if formatLit, ok := call.Args[0].(*ast.BasicLit); ok && strings.Contains(formatLit.Value, "%w") {
								wrappedArg := call.Args[1]
								switch expr := wrappedArg.(type) {
								case *ast.SelectorExpr:
									pass.Reportf(expr.Pos(), sentinelErrMsg)
								case *ast.Ident:
									obj := pass.TypesInfo.ObjectOf(expr)
									if obj != nil && isConstant(obj) {
										pass.Reportf(expr.Pos(), sentinelErrMsg)
									}
								}
							}
						}
					}
				}
			}

			return true
		})
	}
	return nil, nil
}

func isErrorType(t types.Type, pass *analysis.Pass) bool {
	errorType := types.Universe.Lookup("error").Type().Underlying().(*types.Interface)
	return types.Implements(t, errorType)
}

func insideIfErrorsIs(pass *analysis.Pass, node ast.Node) bool {
	var insideErrorsIs bool

	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			if ifStmt, ok := n.(*ast.IfStmt); ok {
				if call, ok := ifStmt.Cond.(*ast.CallExpr); ok {
					if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
						if pkgIdent, ok := sel.X.(*ast.Ident); ok && pkgIdent.Name == "errors" &&
							(sel.Sel.Name == "Is" || sel.Sel.Name == "As") {

							if containsNode(ifStmt.Body, node) {
								insideErrorsIs = true
								return false
							}
							if ifStmt.Else != nil && containsNode(ifStmt.Else, node) {
								insideErrorsIs = true
								return false
							}
						}
					}
				}
			}
			return true
		})
	}
	return insideErrorsIs
}

func containsNode(parent, child ast.Node) bool {
	found := false
	ast.Inspect(parent, func(n ast.Node) bool {
		if n == child {
			found = true
			return false
		}
		return !found
	})
	return found
}

func isConstant(obj types.Object) bool {
	_, isConst := obj.(*types.Const)
	return isConst
}
