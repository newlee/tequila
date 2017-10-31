package dot_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDot(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dot Suite")
}
