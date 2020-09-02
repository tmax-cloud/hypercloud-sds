package testCeph

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"strconv"
	"time"
)

var (
	//TODO 파일이름, 코드 내 이름 정리 'rook', 'ceph', 'rook-ceph' 통일
	testingNamespaceRook *corev1.Namespace
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
)

var _ = Describe("Test Rook Ceph Module", func() {
	BeforeEach(func() {
		testIDRook++
		//TODO [Rook] 이름 변경
		log.Printf("========== [TEST][Rook][CASE-#%d] Started ==========\n", testIDRook)

		// Create testing namespace
		testingNamespaceRook, err = createNamespace(hyperStorageHelper.Clientset,
			makeNamespaceSpec(RookCephNamespacePrefix))
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		// Delete testing namespace
		err = hyperStorageHelper.Clientset.CoreV1().Namespaces().
			Delete(testingNamespaceRook.Name, &metav1.DeleteOptions{})
		Expect(err).ToNot(HaveOccurred())

		Eventually(func() bool {
			ns, err := hyperStorageHelper.Clientset.CoreV1().Namespaces().
				Get(testingNamespaceRook.Name, metav1.GetOptions{})
			if err != nil || errors.IsNotFound(err) {
				return true
			}
			if ns.Status.Phase == corev1.NamespaceTerminating {
				log.Printf("Namespace %s is still in phase %s\n", testingNamespaceRook.Name, ns.Status.Phase)
				return false
			}
			return false
		}, TimeoutForDeletingNamespace, PollingIntervalForDeletingNamespace).Should(BeTrue())
		log.Printf("========== [TEST][Rook][CASE-#%d] Finished ==========\n", testIDRook)
	})

	Describe("[[TEST][Rook] Create Cephfs PVC Create -> Get -> Delete]", func() {
		It("Create Cephfs PVC", func() {
			testPvcNameRook = CephPvcNamePrefix + strconv.Itoa(testIDRook)

			// Create PVC
			pvc, err := createPvcInStorageClass(hyperStorageHelper.Clientset,
				makePvcInStorageClassSpec(testPvcNameRook, testingNamespaceRook.Name, VolumeSize, CephFsSc,
					corev1.ReadWriteMany))
			Expect(err).ToNot(HaveOccurred())

			err = waitPvcGetReadyThenDelete(pvc, hyperStorageHelper)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("[[TEST][Rook] Create RBD PVC Create -> Get -> Delete]", func() {
		It("Create RBD PVC", func() {
			testPvcNameRook = CephPvcNamePrefix + strconv.Itoa(testIDRook)

			pvc, err := createPvcInStorageClass(hyperStorageHelper.Clientset,
				makePvcInStorageClassSpec(testPvcNameRook, testingNamespaceRook.Name, VolumeSize, RbdSc,
					corev1.ReadWriteOnce))
			Expect(err).ToNot(HaveOccurred())

			err = waitPvcGetReadyThenDelete(pvc, hyperStorageHelper)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})

func waitPvcGetReadyThenDelete(pvc *corev1.PersistentVolumeClaim, hyperStorageHelper *HyperHelper) error {
	Eventually(func() bool {
		pvcOut, err := hyperStorageHelper.Clientset.CoreV1().PersistentVolumeClaims(testingNamespaceRook.Name).
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
}
