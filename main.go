package main

import (
	"fmt"
)

func main() {
	dotFile := "examples/step2-problem.dot"
	dddModel := Parse(dotFile)

	codeDir := "examples/step2-code/html"
	codeModel := ParseCodeDir(codeDir, make([]string, 0))
	fmt.Println(dddModel.Compare(codeModel))
}
