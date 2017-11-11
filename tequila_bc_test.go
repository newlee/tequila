package main_test

import (
	"fmt"
	. "github.com/newlee/tequila"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tequila", func() {
	Context("bc code compare", func() {
		It("problem domain", func() {
			dotFile := "examples/cargo-problem.dot"
			dddModel := ParseProblemModel(dotFile)

			codeDir := "examples/bc-code/html"
			codeModel := ParseCodeProblemModel(codeDir, make([]string, 0))
			Expect(dddModel.Compare(codeModel)).Should(BeNil())
		})

		It("solution domain", func() {
			dotFile := "examples/cargo-bc.dot"
			bcModel := ParseSolutionModel(dotFile)

			Expect(len(bcModel.Layers)).Should(Equal(5))

			codeDir := "examples/bc-code/html"
			codeModel := ParseCodeSolutionModel(codeDir)

			fmt.Println(codeModel)
			Expect(bcModel.Compare(codeModel)).Should(BeNil())
		})
	})
})
