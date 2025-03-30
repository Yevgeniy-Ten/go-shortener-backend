package analyzer

import "testing"
import "golang.org/x/tools/go/analysis/analysistest"

func TestOsExitAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), OsExitAnalyzer, "./...")
}
