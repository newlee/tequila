package main_test

import (
	. "github.com/newlee/tequila"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tequila", func() {
	const subDomainName = "subdomain"
	const aggregateAName = "AggregateRootA"
	const aggregateBName = "AggregateRootB"
	Context("Parse DDD Model", func() {
		It("step1", func() {

			dotFile := "examples/step1-problem.dot"
			ars := Parse(dotFile).SubDomains[subDomainName].ARs

			Expect(len(ars)).Should(Equal(1))
			aggregateA := ars[aggregateAName]
			Expect(len(aggregateA.ChildrenEntities())).Should(Equal(1))
			Expect(len(aggregateA.ChildrenValueObjects())).Should(Equal(1))
			entityB := aggregateA.ChildrenEntities()[0]
			Expect(len(entityB.ChildrenValueObjects())).Should(Equal(1))
		})
		It("step2", func() {

			dotFile := "examples/step2-problem.dot"
			ars := Parse(dotFile).SubDomains[subDomainName].ARs

			Expect(len(ars)).Should(Equal(2))
			aggregateA := ars[aggregateAName]
			Expect(len(aggregateA.ChildrenEntities())).Should(Equal(1))
			Expect(len(aggregateA.ChildrenValueObjects())).Should(Equal(1))
			entityB := aggregateA.ChildrenEntities()[0]
			Expect(len(entityB.ChildrenValueObjects())).Should(Equal(1))

			aggregateB := ars[aggregateBName]
			Expect(len(aggregateB.ChildrenEntities())).Should(Equal(0))
			Expect(len(aggregateB.Refs)).Should(Equal(1))
			Expect(aggregateB.Refs[0]).Should(Equal(aggregateA))
		})
		It("step2 with repository", func() {

			dotFile := "examples/step2-problem.dot"
			model := Parse(dotFile)
			ars := model.SubDomains[subDomainName].ARs
			repos := model.SubDomains[subDomainName].Repos

			Expect(len(repos)).Should(Equal(1))
			Expect(repos["AggregateRootARepo"].For).Should(Equal(ars[aggregateAName]))
		})
		It("step2 with provider interface", func() {

			dotFile := "examples/step2-problem.dot"
			model := Parse(dotFile)
			providers := model.SubDomains[subDomainName].Providers

			Expect(len(providers)).Should(Equal(1))
		})

		It("sub domain", func() {
			dotFile := "examples/subdomain.dot"
			model := Parse(dotFile)

			Expect(len(model.SubDomains)).Should(Equal(2))
			subDomain := model.SubDomains["subdomain1"]
			ars := subDomain.ARs
			aggregateA := ars[aggregateAName]
			Expect(len(aggregateA.ChildrenEntities())).Should(Equal(1))
			Expect(len(aggregateA.ChildrenValueObjects())).Should(Equal(1))
			entityB := aggregateA.ChildrenEntities()[0]
			Expect(len(entityB.ChildrenValueObjects())).Should(Equal(1))

			aggregateB := ars[aggregateBName]
			Expect(len(aggregateB.ChildrenEntities())).Should(Equal(0))
			Expect(len(aggregateB.Refs)).Should(Equal(1))
			Expect(aggregateB.Refs[0]).Should(Equal(aggregateA))

			subDomain = model.SubDomains["subdomain2"]
			ars = subDomain.ARs
			aggregateC := ars["AggregateRootC"]
			Expect(len(aggregateC.ChildrenEntities())).Should(Equal(1))
			Expect(len(aggregateC.ChildrenValueObjects())).Should(Equal(0))
			EntityC := aggregateC.ChildrenEntities()[0]
			Expect(len(EntityC.ChildrenValueObjects())).Should(Equal(0))

		})
	})

	Context("Parse Doxygen dot files", func() {

		It("step1", func() {

			codeDir := "examples/step1-code/html"
			codeArs := ParseCodeDir(codeDir).SubDomains[subDomainName].ARs

			Expect(len(codeArs)).Should(Equal(1))
			Expect(len(codeArs[aggregateAName].ChildrenEntities())).Should(Equal(1))
			Expect(len(codeArs[aggregateAName].ChildrenValueObjects())).Should(Equal(1))
			entityB := codeArs[aggregateAName].ChildrenEntities()[0]
			Expect(len(entityB.ChildrenValueObjects())).Should(Equal(1))
		})
		It("step2", func() {

			codeDir := "examples/step2-code/html"
			codeArs := ParseCodeDir(codeDir).SubDomains[subDomainName].ARs

			Expect(len(codeArs)).Should(Equal(2))
			ara := aggregateAName
			arb := aggregateBName

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
			ars := model.SubDomains[subDomainName].ARs
			repos := model.SubDomains[subDomainName].Repos

			Expect(len(repos)).Should(Equal(1))
			Expect(repos["AggregateRootARepo"].For).Should(Equal(ars[aggregateAName]))
		})

		It("step2 with provider interface", func() {
			codeDir := "examples/step2-code/html"
			model := ParseCodeDir(codeDir)
			providers := model.SubDomains[subDomainName].Providers

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

	Context("Parse Doxygen dot files with java", func() {

		It("step2", func() {

			codeDir := "examples/step2-Java/html"
			codeArs := ParseCodeDir(codeDir).SubDomains[subDomainName].ARs

			Expect(len(codeArs)).Should(Equal(2))
			ara := aggregateAName
			arb := aggregateBName

			Expect(len(codeArs[ara].ChildrenEntities())).Should(Equal(1))
			Expect(len(codeArs[ara].ChildrenValueObjects())).Should(Equal(1))
			entityB := codeArs[ara].ChildrenEntities()[0]
			Expect(len(entityB.ChildrenValueObjects())).Should(Equal(1))

			Expect(len(codeArs[arb].ChildrenEntities())).Should(Equal(0))
			Expect(len(codeArs[arb].Refs)).Should(Equal(1))
			Expect(codeArs[arb].Refs[0]).Should(Equal(codeArs[ara]))
		})

		It("step2 with repository", func() {

			codeDir := "examples/step2-Java/html"
			model := ParseCodeDir(codeDir)
			ars := model.SubDomains[subDomainName].ARs
			repos := model.SubDomains[subDomainName].Repos

			Expect(len(repos)).Should(Equal(1))
			Expect(repos["AggregateRootARepo"].For).Should(Equal(ars[aggregateAName]))
		})

		It("step2 with provider interface", func() {
			codeDir := "examples/step2-Java/html"
			model := ParseCodeDir(codeDir)
			providers := model.SubDomains[subDomainName].Providers

			Expect(len(providers)).Should(Equal(1))
		})
		It("step3 should failded when aggregate ref another entity", func() {
			codeDir := "examples/step2-Java/html"
			model := ParseCodeDir(codeDir)

			Expect(model.Validate()).Should(Equal(true))
		})
	})
})
