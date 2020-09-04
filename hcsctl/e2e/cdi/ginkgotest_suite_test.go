package cdi

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// 패키지 변수
var (
	hyperStorageHelper *HyperHelper
	err                error
)

func init() {
	var err error
	_, err = CreateK8sHelper()
	if err != nil {
		panic(err)
	}

	hyperStorageHelper = HyperStorageHelper()
}

func TestGinkgotest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ginkgotest Suite")
}
