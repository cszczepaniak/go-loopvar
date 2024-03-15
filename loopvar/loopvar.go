package loopvar

import (
	"go/ast"
	"go/format"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
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

				r.handleAssignment(f.Body.List[i+1:], assn, nextPos)
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

		seenDiags := make(map[*analysis.Diagnostic]struct{}, len(r.diagsByVar))
		for _, diag := range r.diagsByVar {
			_, ok := seenDiags[diag]
			if ok {
				continue
			}
			seenDiags[diag] = struct{}{}
			p.Report(*diag)
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
	diagsByVar map[string]*analysis.Diagnostic
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
		diagsByVar:        make(map[string]*analysis.Diagnostic),
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

func (r *replacer) handleAssignment(restOfLoop []ast.Stmt, a *ast.AssignStmt, nextPos token.Pos) error {
	var newLhs, newRhs []ast.Expr

	// Preemptively assign this in case we need it :pokerface:
	diag := &analysis.Diagnostic{
		Pos:     a.Pos(),
		End:     a.End(),
		Message: `found unnecessary loop variable capture`,
		SuggestedFixes: []analysis.SuggestedFix{{
			Message: `remove the unnecessary assignments`,
			TextEdits: []analysis.TextEdit{{
				Pos: a.Pos(),
				// Setting nextPos and nil here for now. If we're not removing the entire
				// assignment, we'll need to update End to not consume the newlines leading up to
				// the next statement, and update NewText to simply replace it with the assignments
				// we're not removing.
				End:     nextPos,
				NewText: nil,
			}},
		}},
	}

	for i, rhs := range a.Rhs {
		lhs := a.Lhs[i]

		rhsIdent, ok := rhs.(*ast.Ident)
		if !ok || !r.shadowsRangeVar(rhsIdent.Name) {
			// One of the RHS items was not an identifier. This could be something like:
			// a, b, c := a, fn(123), c
			newLhs = append(newLhs, lhs)
			newRhs = append(newRhs, rhs)
			continue
		}

		// If the RHS does shadow a loop var, we need to map the RHS name to the LHS name so we can
		// subsequently replace usages of the variable.

		// LHS should never be nil.
		if lhs == nil {
			continue
		}

		lhsIdent, ok := lhs.(*ast.Ident)
		if !ok {
			// Somehow we don't have an identifier on the left-hand side.
			continue
		}

		// We have a loop variable being assigned to a different identifier. One last check: we have
		// to make sure that the new binding is never assigned to. If it is, we should keep the
		// assignment.
		if isMutated(restOfLoop, lhsIdent) {
			newLhs = append(newLhs, lhs)
			newRhs = append(newRhs, rhs)
			continue
		}

		// map the diagnostic by the LHS name. That way when we encounter usages of the LHS var, we can
		// index back into this to add more suggested fixes.
		r.diagsByVar[lhsIdent.Name] = diag

		// If the lefthand side's name is the same as a range var, we just need to remove the
		// assignment, no need to replace anything else. Otherwise, a new variable shadows the range
		// var. Every subsequent usage of this var should be replaced with the range var.
		if lhsIdent.Name != rhsIdent.Name {
			r.identReplacements[lhsIdent.Name] = rhsIdent.Name
		}
	}

	if len(newLhs) > 0 {
		allUnderscore := true
		for _, lhs := range newLhs {
			if !isUnderscoreIdent(lhs) {
				allUnderscore = false
				break
			}
		}

		var tok token.Token
		if allUnderscore {
			// If everything on the LHS, we _must_ use `=` as our token, or we'll produce code that
			// doesn't compile. That is, the following doesn't compile:
			//   _, _ := 1, 2
			// ...But this does:
			//   _, _ = 1, 2
			tok = token.ASSIGN
		} else {
			tok = a.Tok
		}

		// In this case, we're not completely removing the assignment. We need to synthesize the new
		// content.
		newAssign := &ast.AssignStmt{
			Lhs: newLhs,
			Tok: tok,
			Rhs: newRhs,
		}

		sb := &strings.Builder{}
		err := format.Node(sb, token.NewFileSet(), newAssign)
		if err != nil {
			return err
		}

		// In this case, we're not deleting the entire line. Reset the End token to the current end
		// of the assignment.
		diag.SuggestedFixes[0].TextEdits[0].End = a.End()
		diag.SuggestedFixes[0].TextEdits[0].NewText = []byte(sb.String())
	}

	return nil
}

func isMutated(stmts []ast.Stmt, ident *ast.Ident) (mut bool) {
	for _, stmt := range stmts {
		foundMutation := false
		ast.Inspect(stmt, func(n ast.Node) bool {
			switch tn := n.(type) {
			case *ast.IncDecStmt:
				if isMatchingIdent(tn.X, ident) {
					foundMutation = true
					return false
				}
			case *ast.AssignStmt:
				switch tn.Tok {
				case token.ASSIGN,
					token.ADD_ASSIGN,
					token.SUB_ASSIGN,
					token.MUL_ASSIGN,
					token.QUO_ASSIGN,
					token.REM_ASSIGN,
					token.AND_ASSIGN,
					token.OR_ASSIGN,
					token.XOR_ASSIGN,
					token.SHL_ASSIGN,
					token.SHR_ASSIGN,
					token.AND_NOT_ASSIGN:

					for _, lhs := range tn.Lhs {
						if isMatchingIdent(lhs, ident) {
							foundMutation = true
							return false
						}
					}
				default:
					return true
				}
			}

			return true
		})
		if foundMutation {
			return true
		}
	}
	return false
}

func isMatchingIdent(n ast.Node, exp *ast.Ident) bool {
	if n == nil {
		return false
	}
	ident, ok := n.(*ast.Ident)
	return ok && ident.Name == exp.Name
}

func isUnderscoreIdent(n ast.Node) bool {
	return isMatchingIdent(n, &ast.Ident{
		Name: `_`,
	})
}
