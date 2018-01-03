package viz_test

import (
	. "github.com/newlee/tequila/viz"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"
)

var _ = Describe("Viz", func() {
	Context("Parse all include dependencies", func() {
		It("bc code", func() {
			codeDir := "../examples/bc-code/html"
			result := ParseInclude(codeDir)
			Expect(len(result.NodeList)).Should(Equal(12))
			Expect(len(result.RelationList)).Should(Equal(14))
			var mergeFunc = func(input string) string {
				return strings.Replace(strings.Replace(input, ".h", "", -1), ".cpp", "", -1)
			}
			crossRefs := result.FindCrossRef(mergeFunc)
			Expect(len(crossRefs)).Should(Equal(0), "Cross references: %v", crossRefs)
		})

		It("merge header files", func() {
			codeDir := "../examples/bc-code/html"
			fullGraph := ParseInclude(codeDir)

			result := fullGraph.MergeHeaderFile(MergeHeaderFunc)
			Expect(len(result.NodeList)).Should(Equal(8))
			Expect(len(result.RelationList)).Should(Equal(10))
		})

		It("merge package", func() {
			codeDir := "../examples/bc-code/html"
			fullGraph := ParseInclude(codeDir)

			result := fullGraph.MergeHeaderFile(MergePackageFunc)
			Expect(len(result.NodeList)).Should(Equal(6))
			Expect(len(result.RelationList)).Should(Equal(8))
		})

		It("entry points", func() {
			codeDir := "../examples/bc-code/html"
			fullGraph := ParseInclude(codeDir)

			entryPoints := fullGraph.EntryPoints(MergePackageFunc)
			Expect(len(entryPoints)).Should(Equal(2))
		})

		It("sort by fan-in fan-out", func() {
			codeDir := "../examples/bc-code/html"
			fullGraph := ParseInclude(codeDir)

			fans := fullGraph.SortedByFan(MergeHeaderFunc)
			Expect(len(fans)).Should(Equal(8))
			//fan := fans[0]
			//Expect(fan.Name).Should(Equal("services/service"))
			//Expect(fan.FanIn).Should(Equal(3))
			//Expect(fan.FanOut).Should(Equal(1))
		})
	})
	Context("Parse all collaboration", func() {
		It("bc code", func() {
			codeDir := "../examples/step2-Java/html"
			result := ParseColl(codeDir, "coll__graph.dot")
			Expect(len(result.NodeList)).Should(Equal(13))
			Expect(len(result.RelationList)).Should(Equal(15))
		})
	})
})
