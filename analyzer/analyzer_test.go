package analyzer

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestLoglint(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get wd: %s", err)
	}

	testdata := filepath.Join(filepath.Dir(wd), "testdata")

	analysistest.Run(t, testdata, Analyzer,
		"valid",
		"invalid/uppercase",
		"invalid/nonenglish",
		"invalid/specialchars",
		"invalid/sensitive",
	)
}
