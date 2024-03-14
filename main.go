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
	Doc:      `Detects loop variables that are unnecessarily captured (as of Go 1.22)`,
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
		for i, stmt := range f.Body.List {
			assn, ok := stmt.(*ast.AssignStmt)
			if ok && assn.Tok == token.DEFINE {
				var nextPos token.Pos
				if i < len(f.Body.List)-1 {
					// This isn't the last statement in the list. Look at the beginning of the next one.
					nextPos = f.Body.List[i+1].Pos()
				} else {
					// It shouldn't be possible for an assignment to be the last thing in the loop.
					// If it were, it'd be unused and the compiler would complain.
					panic(`dev error: should not encounter a shadowed assignment at the end of the list`)
				}

				r.handleAssignment(assn, nextPos)
				continue
			}

			ast.Inspect(stmt, func(n ast.Node) bool {
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
		Message: `replace this variable with the loop variable`,
		TextEdits: []analysis.TextEdit{{
			Pos:     ident.Pos(),
			End:     ident.End(),
			NewText: []byte(toReplace),
		}},
	})

	r.diagsByVar[ident.Name] = diag
}

func (r *replacer) handleAssignment(a *ast.AssignStmt, nextPos token.Pos) {
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
		Message: `found unnecessary loop variable capture`,
		SuggestedFixes: []analysis.SuggestedFix{{
			Message: `remove the assignment`,
			TextEdits: []analysis.TextEdit{{
				Pos: a.Pos(),
				// Use the end of the assignment. We'll try to update it later to remove any
				// trailing whitespace leading up to the next statement in the loop.
				End:     nextPos,
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
