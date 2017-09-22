package main_test

import (
	. "github.com/newlee/tequila"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tequila", func() {
	Context("Parse DDD Model", func() {
		It("step1", func() {

			dotFile := "examples/step1-problem.dot"
			ars := Parse(dotFile).ARs

			Expect(len(ars)).Should(Equal(1))
			Expect(len(ars["AggregateRootA"].ChildrenEntities())).Should(Equal(1))
			Expect(len(ars["AggregateRootA"].ChildrenValueObjects())).Should(Equal(1))
			entityB := ars["AggregateRootA"].ChildrenEntities()[0]
			Expect(len(entityB.ChildrenValueObjects())).Should(Equal(1))
		})
		It("step2", func() {

			dotFile := "examples/step2-problem.dot"
			ars := Parse(dotFile).ARs

			Expect(len(ars)).Should(Equal(2))
			Expect(len(ars["AggregateRootA"].ChildrenEntities())).Should(Equal(1))
			Expect(len(ars["AggregateRootA"].ChildrenValueObjects())).Should(Equal(1))
			entityB := ars["AggregateRootA"].ChildrenEntities()[0]
			Expect(len(entityB.ChildrenValueObjects())).Should(Equal(1))

			Expect(len(ars["AggregateRootB"].ChildrenEntities())).Should(Equal(0))
			Expect(len(ars["AggregateRootB"].Refs)).Should(Equal(1))
			Expect(ars["AggregateRootB"].Refs[0]).Should(Equal(ars["AggregateRootA"]))
		})
		It("step2 with repository", func() {

			dotFile := "examples/step2-problem.dot"
			model := Parse(dotFile)
			ars := model.ARs
			repos := model.Repos

			Expect(len(repos)).Should(Equal(1))
			Expect(repos["AggregateRootARepo"].For).Should(Equal(ars["AggregateRootA"]))
		})
		It("step2 with provider interface", func() {

			dotFile := "examples/step2-problem.dot"
			model := Parse(dotFile)
			providers := model.Providers

			Expect(len(providers)).Should(Equal(1))
		})
	})

	Context("Parse Doxygen dot files", func() {

		It("step1", func() {

			codeDir := "examples/step1-code/html"
			codeArs := ParseCodeDir(codeDir).ARs

			Expect(len(codeArs)).Should(Equal(1))
			Expect(len(codeArs["AggregateRootA"].ChildrenEntities())).Should(Equal(1))
			Expect(len(codeArs["AggregateRootA"].ChildrenValueObjects())).Should(Equal(1))
			entityB := codeArs["AggregateRootA"].ChildrenEntities()[0]
			Expect(len(entityB.ChildrenValueObjects())).Should(Equal(1))
		})
		It("step2", func() {

			codeDir := "examples/step2-code/html"
			codeArs := ParseCodeDir(codeDir).ARs

			Expect(len(codeArs)).Should(Equal(2))
			ara := "AggregateRootA"
			arb := "AggregateRootB"
			Expect(len(codeArs[ara].ChildrenEntities())).Should(Equal(1))
			Expect(len(codeArs[ara].ChildrenValueObjects())).Should(Equal(1))
			entityB := codeArs[ara].ChildrenEntities()[0]
			Expect(len(entityB.ChildrenValueObjects())).Should(Equal(1))

			Expect(len(codeArs[arb].ChildrenEntities())).Should(Equal(0))
			Expect(len(codeArs[arb].Refs)).Should(Equal(1))
			Expect(codeArs[arb].Refs[0]).Should(Equal(codeArs[ara]))
		})

		It("step2 with repository", func() {

			codeDir := "examples/step2-code/html"
			model := ParseCodeDir(codeDir)
			ars := model.ARs
			repos := model.Repos

			Expect(len(repos)).Should(Equal(1))
			Expect(repos["AggregateRootARepo"].For).Should(Equal(ars["AggregateRootA"]))
		})

		It("step2 with provider interface", func() {
			codeDir := "examples/step2-code/html"
			model := ParseCodeDir(codeDir)
			providers := model.Providers

			Expect(len(providers)).Should(Equal(1))
		})
		It("step3 should failded when aggregate ref another entity", func() {
			codeDir := "examples/step2-code/html"
			model := ParseCodeDir(codeDir)

			Expect(model.Validate()).Should(Equal(true))

			codeDir = "examples/step3-code/html"
			model = ParseCodeDir(codeDir)

			Expect(model.Validate()).Should(Equal(false))
		})
	})
})
