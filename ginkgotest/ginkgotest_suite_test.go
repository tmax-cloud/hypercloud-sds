package ginkgotest

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGinkgotest(t *testing.T) {
	RegisterFailHandler(Fail)
	_, _ = CreateK8sHelper(t)
	RunSpecs(t, "Ginkgotest Suite")
}

var (
	hyperStorageHelper *HyperHelper
	//config    *restclient.Config
)
