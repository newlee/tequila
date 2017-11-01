package main

import (
	"fmt"
)

func main() {
	dotFile := "examples/cargo-problem.dot"
	dddModel := Parse(dotFile)

	codeDir := "examples/bc-code/html"
	codeModel := ParseCodeDir(codeDir, make([]string, 0))
	fmt.Println(dddModel.Compare(codeModel))
}
