package testCeph

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// 패키지 변수
var (
	hyperSdsHelper *HyperHelper
	err                error
)

func init() {
	var err error
	_, err = CreateK8sHelper()
	if err != nil {
		panic(err)
	}

	hyperSdsHelper = HyperSdsHelper()
}

func TestGinkgotest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ginkgotest Suite")
}
