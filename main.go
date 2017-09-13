package main

import (
	"fmt"
)

func main() {
	dotFile := "examples/step2-problem.dot"
	ars := Parse(dotFile)

	codeDir := "examples/step2-code/html"
	codeArs := ParseCode(codeDir)
	for key, _ := range codeArs {
		ar := codeArs[key]
		fmt.Println(ar.Compare(ars[key]))
	}
}
