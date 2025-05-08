//go:build plugin
package main

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
			// Check return statements
			if retStmt, ok := n.(*ast.ReturnStmt); ok {
				for _, retExpr := range retStmt.Results {
					// Skip if inside `if errors.Is(...)` context
					if insideIfErrorsIs(pass, retStmt) {
						continue
					}

					// Check direct returns of constant errors
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

			// Check fmt.Errorf with %w
			if call, ok := n.(*ast.CallExpr); ok {
				if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
					if pkgIdent, ok := sel.X.(*ast.Ident); ok && pkgIdent.Name == "fmt" && sel.Sel.Name == "Errorf" {
						if len(call.Args) >= 2 {
							// Check format string contains %w
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

// Helper: is type error?
func isErrorType(t types.Type, pass *analysis.Pass) bool {
	errorType := types.Universe.Lookup("error").Type().Underlying().(*types.Interface)
	return types.Implements(t, errorType)
}

// Helper: check if node is inside `if errors.Is(...)`
func insideIfErrorsIs(pass *analysis.Pass, node ast.Node) bool {
	var insideErrorsIs bool

	// For each file in the package
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			if ifStmt, ok := n.(*ast.IfStmt); ok {
				// Check if it's an errors.Is/As call
				if call, ok := ifStmt.Cond.(*ast.CallExpr); ok {
					if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
						if pkgIdent, ok := sel.X.(*ast.Ident); ok && pkgIdent.Name == "errors" &&
							(sel.Sel.Name == "Is" || sel.Sel.Name == "As") {

							// Check if our node is in this if statement's body or else clause
							if containsNode(ifStmt.Body, node) {
								insideErrorsIs = true
								return false // Stop traversing
							}
							if ifStmt.Else != nil && containsNode(ifStmt.Else, node) {
								insideErrorsIs = true
								return false // Stop traversing
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

// Helper: check if a parent node contains a child node
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

// Helper: check if object is a constant
func isConstant(obj types.Object) bool {
	_, isConst := obj.(*types.Const)
	return isConst
}
