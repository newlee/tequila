package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTequila(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tequila Suite")
}
