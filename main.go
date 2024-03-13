package main

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/singlechecker"
	"golang.org/x/tools/go/ast/inspector"
)

func main() {
	singlechecker.Main(a)
}

var a = &analysis.Analyzer{
	Name:     `loopvar`,
	Run:      run,
	Doc:      `something`,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(p *analysis.Pass) (any, error) {
	inspect := p.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	inspect.Preorder([]ast.Node{&ast.RangeStmt{}}, func(n ast.Node) {
		f := n.(*ast.RangeStmt)
		if f.Tok != token.DEFINE {
			return
		}

		r := newReplacer(f.Key, f.Value)
		for _, stmt := range f.Body.List {
			assn, ok := stmt.(*ast.AssignStmt)
			if ok && assn.Tok == token.DEFINE {
				r.handleAssignment(assn)
				continue
			}

			ast.Inspect(stmt, func(n ast.Node) bool {
				if _, ok := n.(*ast.BlockStmt); ok {
					// Don't recurse into block stmts?
					// return false
				}
				ident, ok := n.(*ast.Ident)
				if ok {
					r.handleIdent(ident)
				}
				return true
			})
		}

		for _, diag := range r.diagsByVar {
			p.Report(diag)
		}
	})

	return nil, nil
}

type replacer struct {
	// rangeVars is the set of identifiers from the range statement
	rangeVars map[string]struct{}
	// map from an identifier to the identifier we'll replace it with (which will be one of the loop
	// vars)
	identReplacements map[string]string
	// map from range var name to diagnostic
	diagsByVar map[string]analysis.Diagnostic
}

func newReplacer(rangeVars ...ast.Node) *replacer {
	rangeVarNames := map[string]struct{}{}
	for _, v := range rangeVars {
		if v == nil {
			continue
		}
		ident, ok := v.(*ast.Ident)
		if ok {
			rangeVarNames[ident.Name] = struct{}{}
		}
	}
	return &replacer{
		rangeVars:         rangeVarNames,
		identReplacements: make(map[string]string),
		diagsByVar:        make(map[string]analysis.Diagnostic),
	}
}

func (r *replacer) shadowsRangeVar(ident string) bool {
	_, ok := r.rangeVars[ident]
	return ok
}

func (r *replacer) handleIdent(ident *ast.Ident) {
	if ident == nil {
		return
	}

	toReplace, ok := r.identReplacements[ident.Name]
	if !ok {
		return
	}

	diag, ok := r.diagsByVar[ident.Name]
	if !ok {
		// This would be odd if we had a replacement.
		return
	}

	diag.SuggestedFixes = append(diag.SuggestedFixes, analysis.SuggestedFix{
		Message: `baz`,
		TextEdits: []analysis.TextEdit{{
			Pos:     ident.Pos(),
			End:     ident.End(),
			NewText: []byte(toReplace),
		}},
	})

	r.diagsByVar[ident.Name] = diag
}

func (r *replacer) handleAssignment(a *ast.AssignStmt) {
	if len(a.Rhs) != 1 {
		// TODO we should also handle multiple-assignments
		return
	}

	rhsIdent, ok := a.Rhs[0].(*ast.Ident)
	if !ok || !r.shadowsRangeVar(rhsIdent.Name) {
		// If the RHS isn't an identifier or it doesn't shadow a range var, don't do anything.
		return
	}

	// Okay, we found somebody capturing the range var. We'll report this and suggest removing the
	// assignment.
	diag := analysis.Diagnostic{
		Pos:     a.Pos(),
		End:     a.End(),
		Message: `found a for-range loop shadowed var`,
		SuggestedFixes: []analysis.SuggestedFix{{
			Message: `foobar`,
			TextEdits: []analysis.TextEdit{{
				Pos:     a.Pos(),
				End:     a.End(),
				NewText: nil,
			}},
		}},
	}

	if a.Lhs[0] == nil {
		return
	}

	lhsIdent, ok := a.Lhs[0].(*ast.Ident)
	if !ok {
		// Somehow we don't have an identifier on the left-hand side.
		return
	}

	// map the diagnostic by the LHS name. That way when we encounter usages of the LHS var, we can
	// index back into this to add more suggested fixes.
	r.diagsByVar[lhsIdent.Name] = diag

	// If the lefthand side's name is the same as a range var, we just need to remove the
	// assignment, no need to replace anything else.
	if lhsIdent.Name == rhsIdent.Name {
		return
	}

	// Otherwise, a new variable shadows the range var. Every subsequent usage of this var should be
	// replaced with the range var.
	r.identReplacements[lhsIdent.Name] = rhsIdent.Name
}
