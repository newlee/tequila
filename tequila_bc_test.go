package main_test

import (
	. "github.com/newlee/tequila"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"fmt"
)

var _ = Describe("Tequila", func() {
	Context("bc code compare", func() {

		It("problem domain", func() {
			dotFile := "examples/cargo-problem.dot"
			dddModel := Parse(dotFile)

			codeDir := "examples/bc-code/html"
			codeModel := ParseCodeDir(codeDir, make([]string, 0))
			fmt.Println(codeModel.SubDomains["subdomain"].Providers)
			Expect(dddModel.Compare(codeModel)).Should(Equal(true))

		})
	})

})
