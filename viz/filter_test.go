package viz_test

import (
	. "github.com/newlee/tequila/viz"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Viz", func() {
	Context("Filter by regexp", func() {
		It("filter with white list", func() {
			filter := NewRegexpFilter()
			filter.AddReg(".*Hello.*")
			filter.AddReg("^World.*")
			filter.AddReg("^FOO$")

			matchedString := []string{"Hi Hello", "Hi Hello w", "World H", "World", "FOO"}
			for _, s := range matchedString {
				Expect(filter.Match(s)).Should(BeTrue())
			}

			unMatchedString := []string{"Hi Hell", "Hi Hell w", "Hi World", "FOO B", "B FOO", "BAR"}
			for _, s := range unMatchedString {
				Expect(filter.Match(s)).Should(BeFalse())
			}
		})

		It("filter with black list", func() {
			filter := NewRegexpFilter()
			filter.AddReg(".*")
			filter.AddReg("- ^Hello$")
			filter.AddReg("- ^((?!World).)*FOO((?!World).)*$")

			unMatchedString := []string{"Hi Hello", "Hello w", "World H", "World FOO", "FOO World"}
			for _, s := range unMatchedString {
				Expect(filter.Match(s)).Should(BeTrue())
			}

			matchedString := []string{"Hello", "W FOO Q"}
			for _, s := range matchedString {
				Expect(filter.Match(s)).Should(BeFalse())
			}
		})

		It("filter with excludes", func() {
			filter := NewRegexpFilter()
			filter.AddReg(".*")
			filter.AddExclude("Hello")

			unMatchedString := []string{"Hi Hello", "Hello w", "World H", "World FOO", "FOO World"}
			for _, s := range unMatchedString {
				Expect(filter.Match(s)).Should(BeTrue())
			}

			matchedString := []string{"Hello"}
			for _, s := range matchedString {
				Expect(filter.Match(s)).Should(BeFalse())
			}

			for _, s := range matchedString {
				Expect(filter.UnMatch(s)).Should(BeFalse())
			}
			for _, s := range matchedString {
				Expect(filter.NotMatch(s)).Should(BeFalse())
			}
		})
	})
})
