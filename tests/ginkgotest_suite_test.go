package tests

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	hyperStorageHelper *HyperHelper
)

func TestGinkgotest(t *testing.T) {
	RegisterFailHandler(Fail)
	_, err := CreateK8sHelper(t)
	if err != nil {
		panic(err)
	}

	RunSpecs(t, "Ginkgotest Suite")
}
