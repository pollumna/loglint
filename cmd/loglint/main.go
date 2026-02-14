package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"loglint/analyzer"
)

func main() {
	logAnalyzer := analyzer.Analyzer

	singlechecker.Main(logAnalyzer)
}
