package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Suite")
	//defer GinkgoRecover()
	//RegisterFailHandler(tests.CDIFailHandler)
	//RunSpecsWithDefaultAndCustomReporters(t, "Tests Suite", reporters.NewReporters())
}
