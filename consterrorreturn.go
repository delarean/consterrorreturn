package consterrorreturn

import (
	"go/ast"
	"go/token"
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
			ifStmt, ok := n.(*ast.IfStmt)
			if !ok {
				return true
			}

			errIdent := extractErrIdent(ifStmt.Cond)
			if errIdent == nil {
				return true
			}

			ast.Inspect(ifStmt.Body, func(node ast.Node) bool {
				retStmt, ok := node.(*ast.ReturnStmt)
				if !ok {
					return true
				}

				for _, retExpr := range retStmt.Results {
					if !IsErrorType(pass.TypesInfo.TypeOf(retExpr)) {
						continue
					}

					if !isAllowedErrorReturn(retExpr, errIdent, pass) {
						pass.Reportf(retExpr.Pos(), sentinelErrMsg)
					}
				}
				return true
			})

			return true
		})
	}
	return nil, nil
}

func extractErrIdent(cond ast.Expr) *ast.Ident {
	binExpr, ok := cond.(*ast.BinaryExpr)
	if !ok || binExpr.Op != token.NEQ {
		return nil
	}

	x, y := binExpr.X, binExpr.Y

	isNil := func(e ast.Expr) bool {
		ident, ok := e.(*ast.Ident)
		return ok && ident.Name == "nil"
	}

	if xIdent, ok := x.(*ast.Ident); ok && !isNil(x) && isNil(y) {
		return xIdent
	}

	if yIdent, ok := y.(*ast.Ident); ok && !isNil(y) && isNil(x) {
		return yIdent
	}

	return nil
}

func isAllowedErrorReturn(retExpr ast.Expr, errIdent *ast.Ident, pass *analysis.Pass) bool {
	switch expr := retExpr.(type) {
	case *ast.Ident:
		return expr.Name == errIdent.Name
	case *ast.CallExpr:
		if sel, ok := expr.Fun.(*ast.SelectorExpr); ok {
			if pkg, ok := sel.X.(*ast.Ident); ok && pkg.Name == "fmt" && sel.Sel.Name == "Errorf" {
				if len(expr.Args) < 2 {
					return false
				}
				formatStr, ok := expr.Args[0].(*ast.BasicLit)
				if !ok || !strings.Contains(formatStr.Value, "%w") {
					return false
				}

				for _, arg := range expr.Args[1:] {
					if argIdent, ok := arg.(*ast.Ident); ok && argIdent.Name == errIdent.Name {
						return true
					}
				}
			}
		}
	}
	return false
}

func IsErrorType(t types.Type) bool {
	if t == nil {
		return false
	}
	errorType := types.Universe.Lookup("error").Type().Underlying().(*types.Interface)
	return types.Implements(t, errorType)
}
