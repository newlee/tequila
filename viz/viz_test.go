package viz_test

import (
	. "github.com/newlee/tequila/viz"

	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Viz", func() {
	Context("Parse all graph dot files", func() {
		It("bc code", func() {
			codeDir := "../examples/bc-code/html"
			result := ParseCodeDir(codeDir)
			fmt.Println(result)
			Expect(len(result.NodeList)).Should(Equal(12))
		})
	})
})
