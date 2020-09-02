package testCeph

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"

	"log"
	"time"
)

var (
	//TODO 파일이름, 코드 내 이름 정리 'rook', 'ceph', 'rook-ceph' 통일
	testingNamespaceRook string
	testIDRook           int
	testPvcNameRook      string
)

const (
	RookCephNamespacePrefix = "test-ceph-"
	CephPvcNamePrefix       = "ceph-"

	VolumeSize = "100Mi"
	TimeOutForCreatingPvc = 60 * time.Second *2

	CephFsSc = "csi-cephfs-sc"
	RbdSc    = "rook-ceph-block"

    TestNamespaceYamlPath   = "manifests/namespace.yaml"
    TestCephFsPvcYamlPath   = "manifests/cephfs-pvc.yaml"
    TestRbdPvcYamlPath      = "manifests/rbd-pvc.yaml"
)

var _ = Describe("Test Rook Ceph Module", func() {
	BeforeEach(func() {
		testIDRook++
		//TODO [Rook] 이름 변경
		log.Printf("========== [TEST][Rook][CASE-#%d] Started ==========\n", testIDRook)

		// Create testing namespace
		err := createK8sObjectFromYaml("namespaces", TestNamespaceYamlPath, *hyperSdsHelper.Client)
		Expect(err).ToNot(HaveOccurred())
		log.Printf("========== [TEST][Rook][CASE-#%d] fucking Start Processing===\n", testIDRook)

	})

	AfterEach(func() {
		// Delete testing namespace
		log.Printf("========== [TEST][Rook][CASE-#%d] fucking Finishing Processing===\n", testIDRook)
        err := deleteK8sObjectFromYaml("namespaces", TestNamespaceYamlPath, *hyperSdsHelper.Client)
		Expect(err).ToNot(HaveOccurred())

		Eventually(func() bool {
			namespaceStatus, err := getK8sObjectStatusFromYaml("namespace", TestNamespaceYamlPath, *hyperSdsHelper.Client)
			if err != nil || errors.IsNotFound(err) {
				return true
			}
			if namespaceStatus == string(corev1.NamespaceTerminating) {
				log.Printf("Testing Namespace is still in phase %s\n", namespaceStatus)
				return false
			}
			return false
		}, TimeoutForDeletingNamespace, PollingIntervalForDeletingNamespace).Should(BeTrue())
		log.Printf("========== [TEST][Rook][CASE-#%d] Finished ==========\n", testIDRook)
	})

	Describe("[[TEST][Rook] Create CephFS PVC Create -> Get -> Delete]", func() {
		It("Create Cephfs PVC", func() {

		    log.Printf("========== [TEST][Rook][CASE-#%d] Why fucking does not come here?===\n", testIDRook)
            err := createK8sObjectFromYaml("persistentvolumeclaims", TestCephFsPvcYamlPath, *hyperSdsHelper.Client)
            Expect(err).ToNot(HaveOccurred())

		    log.Printf("========== [TEST][Rook][CASE-#%d] calling wait fucking===\n", testIDRook)
			err = waitPvcGetReadyThenDelete(TestCephFsPvcYamlPath, hyperSdsHelper)
		    log.Printf("========== [TEST][Rook][CASE-#%d] finishing wait fucking===\n", testIDRook)
			Expect(err).ToNot(HaveOccurred())

			//testPvcNameRook = CephPvcNamePrefix + strconv.Itoa(testIDRook)
			// Create PVC
			/*pvc, err := createPvcInStorageClass(hyperStorageHelper.Clientset,
				makePvcInStorageClassSpec(testPvcNameRook, testingNamespaceRook.Name, VolumeSize, CephFsSc,
					corev1.ReadWriteMany))
			Expect(err).ToNot(HaveOccurred())*/
		})
	})

	Describe("[[TEST][Rook] Create RBD PVC Create -> Get -> Delete]", func() {
		It("Create RBD PVC", func() {
		    log.Printf("========== [TEST][Rook][CASE-#%d] RBD, Why fucking does not come here?===\n", testIDRook)
            err := createK8sObjectFromYaml("persistentvolumeclaims", TestRbdPvcYamlPath, *hyperSdsHelper.Client)
            Expect(err).ToNot(HaveOccurred())

		    log.Printf("========== [TEST][Rook][CASE-#%d] RBD, calling wait fucking===\n", testIDRook)
			err = waitPvcGetReadyThenDelete(TestRbdPvcYamlPath, hyperSdsHelper)
			Expect(err).ToNot(HaveOccurred())

			/*testPvcNameRook = CephPvcNamePrefix + strconv.Itoa(testIDRook)

			pvc, err := createPvcInStorageClass(hyperStorageHelper.Clientset,
				makePvcInStorageClassSpec(testPvcNameRook, testingNamespaceRook.Name, VolumeSize, RbdSc,
					corev1.ReadWriteOnce))
			Expect(err).ToNot(HaveOccurred())

			err = waitPvcGetReadyThenDelete(pvc, hyperStorageHelper)
			Expect(err).ToNot(HaveOccurred())
            */
		})
	})
})

func waitPvcGetReadyThenDelete(yamlPath string, hyperSdsHelper *HyperHelper) error {
	Eventually(func() bool {
        pvcStatus, err := getK8sObjectStatusFromYaml("persistentvolumeclaims", yamlPath, *hyperSdsHelper.Client)
		if err == nil && pvcStatus == string(corev1.ClaimBound) {
			log.Printf("Testing Pvc is created and Bound\n")
			return true
        }
		log.Printf("Testing Pvc is still creating in phase %s \n", pvcStatus)
		return false
	}, TimeOutForCreatingPvc, PollingIntervalDefault).Should(BeTrue())

    err = deleteK8sObjectFromYaml("persistentvolumeclaims", yamlPath, *hyperSdsHelper.Client)
	return err

		/*pvcOut, err := hyperStorageHelper.Clientset.CoreV1().PersistentVolumeClaims(testingNamespaceRook.Name).
			Get(pvc.Name, metav1.GetOptions{})
		if err == nil && pvcOut.Status.Phase == corev1.ClaimBound {
			log.Printf("Pvc %s is created and Bound\n", pvcOut.Name)
			return true
		}
		log.Printf("Pvc %s is still creating in phase %s \n", pvcOut.Name, pvcOut.Status.Phase)
		return false
	}, TimeOutForCreatingPvc, PollingIntervalDefault).Should(BeTrue())
    

	err = hyperStorageHelper.Clientset.CoreV1().PersistentVolumeClaims(testingNamespaceRook.Name).
		Delete(pvc.Name, &metav1.DeleteOptions{})

	return err
    */
}
