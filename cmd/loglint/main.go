package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/pollumna/loglint/analyzer"
)

func main() {
	logAnalyzer := analyzer.Analyzer

	singlechecker.Main(logAnalyzer)
}
