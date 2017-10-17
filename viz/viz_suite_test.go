package viz_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestViz(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Viz Suite")
}
