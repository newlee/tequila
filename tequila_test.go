package main_test

import (
	"fmt"
	. "github.com/newlee/tequila"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tequila", func() {
	Context("Parse DDD Model", func() {
		It("step1", func() {

			dotFile := "examples/step1-problem.dot"
			ars := Parse(dotFile)

			Expect(len(ars)).Should(Equal(1))
			Expect(len(ars["AggregateRootA"].ChildrenEntities())).Should(Equal(1))
			Expect(len(ars["AggregateRootA"].ChildrenValueObjects())).Should(Equal(1))
			entityB := ars["AggregateRootA"].ChildrenEntities()[0]
			Expect(len(entityB.ChildrenValueObjects())).Should(Equal(1))
		})
		It("step2", func() {

			dotFile := "examples/step2-problem.dot"
			ars := Parse(dotFile)

			Expect(len(ars)).Should(Equal(2))
			Expect(len(ars["AggregateRootA"].ChildrenEntities())).Should(Equal(1))
			Expect(len(ars["AggregateRootA"].ChildrenValueObjects())).Should(Equal(1))
			entityB := ars["AggregateRootA"].ChildrenEntities()[0]
			Expect(len(entityB.ChildrenValueObjects())).Should(Equal(1))

			Expect(len(ars["AggregateRootB"].ChildrenEntities())).Should(Equal(0))
			fmt.Println(ars["AggregateRootB"])
		})
	})

	Context("Parse Doxygen dot files", func() {
		It("step1", func() {

			codeDir := "examples/step1-code/html"
			codeArs := ParseCode(codeDir)

			Expect(len(codeArs)).Should(Equal(1))
			Expect(len(codeArs["AggregateRootA"].ChildrenEntities())).Should(Equal(1))
			Expect(len(codeArs["AggregateRootA"].ChildrenValueObjects())).Should(Equal(1))
			entityB := codeArs["AggregateRootA"].ChildrenEntities()[0]
			Expect(len(entityB.ChildrenValueObjects())).Should(Equal(1))
		})
		It("step2", func() {

			codeDir := "examples/step2-code/html"
			codeArs := ParseCode(codeDir)

			Expect(len(codeArs)).Should(Equal(2))
			Expect(len(codeArs["AggregateRootA"].ChildrenEntities())).Should(Equal(1))
			Expect(len(codeArs["AggregateRootA"].ChildrenValueObjects())).Should(Equal(1))
			entityB := codeArs["AggregateRootA"].ChildrenEntities()[0]
			Expect(len(entityB.ChildrenValueObjects())).Should(Equal(1))

			Expect(len(codeArs["AggregateRootB"].ChildrenEntities())).Should(Equal(0))
			fmt.Println(codeArs["AggregateRootB"])
		})
	})
})
