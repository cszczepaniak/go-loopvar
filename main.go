package main

import (
	"github.com/cszczepaniak/go-loopvar/loopvar"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(loopvar.Analyzer)
}
