package dot_test

import (
	. "github.com/newlee/tequila/dot"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dot", func() {
	Context("Parse api dot file", func() {
		It("api doxygen file", func() {
			node := ParseDoxygenFile("test.dot")

			Expect(node.Name).Should(Equal("api::Api"))
			serviceNode := node.DstNodes[0].Node
			Expect(serviceNode.Name).Should(Equal("services::CargoService"))

			var cargoNode *Node
			for _, node := range serviceNode.DstNodes[0].Node.DstNodes {
				if node.Node.Name == "domain::Cargo" {
					cargoNode = node.Node
				}
			}
			Expect(cargoNode).ShouldNot(Equal(nil))
		})
	})
})
