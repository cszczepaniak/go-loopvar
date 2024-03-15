package loopvar

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestLoopvar(t *testing.T) {
	analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), Analyzer)
}
