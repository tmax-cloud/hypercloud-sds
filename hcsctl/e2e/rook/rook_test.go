package rook

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"

	"log"
)

var (
	//TODO 파일이름, 코드 내 이름 정리 'rook', 'ceph', 'rook-ceph' 통일
	testIDCeph           int
)

var _ = Describe("Test Rook Ceph Module", func() {
	BeforeEach(func() {
		testIDCeph++
		//TODO [Rook] 이름 변경
		log.Printf("========== [TEST][Rook][CASE-#%d] Started ==========\n", testIDCeph)

		// Create testing namespace
		err := createK8sObjectFromYaml(ResourceNamespace, TestNamespaceYamlPath, *hyperSdsHelper.Client)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		// Delete testing namespace
        err := deleteK8sObjectFromYaml(ResourceNamespace, TestNamespaceYamlPath, *hyperSdsHelper.Client)
		Expect(err).ToNot(HaveOccurred())

		Eventually(func() bool {
			namespaceStatus, err := getK8sObjectStatusFromYaml(ResourceNamespace, TestNamespaceYamlPath, *hyperSdsHelper.Client)

			if err != nil || namespaceStatus == "" {
				return true
			}
			if namespaceStatus == string(corev1.NamespaceTerminating) {
				log.Printf("Testing Namespace is still in phase %s\n", namespaceStatus)
				return false
			}
			return false
		}, TimeoutForDeletingNamespace, PollingIntervalForDeletingNamespace).Should(BeTrue())
		log.Printf("========== [TEST][Rook][CASE-#%d] Finished ==========\n", testIDCeph)
	})

	Describe("[[TEST][Rook] Create CephFS PVC Create -> Get -> Delete]", func() {
		It("Create CephFS PVC", func() {
            err := createK8sObjectFromYaml(ResourcePVC, TestCephFsPvcYamlPath, *hyperSdsHelper.Client)
            Expect(err).ToNot(HaveOccurred())

			err = waitPvcGetReadyThenDelete(TestCephFsPvcYamlPath, hyperSdsHelper)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("[[TEST][Rook] Create RBD PVC Create -> Get -> Delete]", func() {
		It("Create RBD PVC", func() {
            err := createK8sObjectFromYaml(ResourcePVC, TestRbdPvcYamlPath, *hyperSdsHelper.Client)
            Expect(err).ToNot(HaveOccurred())

			err = waitPvcGetReadyThenDelete(TestRbdPvcYamlPath, hyperSdsHelper)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})

func waitPvcGetReadyThenDelete(yamlPath string, hyperSdsHelper *HyperHelper) error {
	Eventually(func() bool {
        pvcStatus, err := getK8sObjectStatusFromYaml(ResourcePVC, yamlPath, *hyperSdsHelper.Client)
		if err == nil && pvcStatus == string(corev1.ClaimBound) {
			log.Printf("Testing PVC is created and Bound\n")
			return true
        }
		log.Printf("Testing PVC is still creating in phase %s \n", pvcStatus)

		return false
	}, TimeOutForCreatingPvc, PollingIntervalDefault).Should(BeTrue())

    err = deleteK8sObjectFromYaml(ResourcePVC, yamlPath, *hyperSdsHelper.Client)
	return err
}
