package tests

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	testingNamespaceRook *corev1.Namespace
	errorRook            error
	testIDRook           string
	testPvcNameRook      string
)

const (
	RookCephNamespacePrefix = "test-ceph-"
	CephPvcNamePrefix       = "ceph-"

	VolumeSize = "100Mi"

	CephFsSc = "csi-cephfs-sc"
	RbdSc    = "rook-ceph-block"
)

var _ = Describe("Test Rook Ceph Module", func() {

	BeforeEach(func() {
		// Create testing namespace
		testingNamespaceRook, errorRook = createNamespace(hyperStorageHelper.Clientset, makeNamespaceSpec(RookCephNamespacePrefix))
		Expect(errorRook).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		// Delete testing namespace
		errorRook = hyperStorageHelper.Clientset.CoreV1().Namespaces().Delete(testingNamespaceRook.Name, &metav1.DeleteOptions{})
		Expect(errorRook).ToNot(HaveOccurred())

		Eventually(func() bool {
			ns, err := hyperStorageHelper.Clientset.CoreV1().Namespaces().Get(testingNamespaceRook.Name, metav1.GetOptions{})
			if err != nil || errors.IsNotFound(err) {
				return true
			}
			if ns.Status.Phase == corev1.NamespaceTerminating {
				fmt.Printf("Namespace %s is still in phase %s\n", testingNamespaceRook.Name, ns.Status.Phase)
				return false
			}
			return false
		}, TimeoutForDeletingNamespace, PollingIntervalForDeletingNamespace).Should(BeTrue())
		fmt.Printf("\n %s Rook Ceph Storage Create => Get => Delete\n", testIDRook)
	})

	Describe("[[TEST][e2e][0002] Create Cephfs PVC Create -> Get -> Delete]", func() {
		It("Create Cephfs PVC", func() {
			testIDRook = "0002"
			testPvcNameRook = CephPvcNamePrefix + testIDRook
			fmt.Printf("[TEST][e2e][%s] started\n", testIDRook)

			// Create PVC
			pvc, err := createPvcInStorageClass(hyperStorageHelper.Clientset,
				makePvcInStorageClassSpec(testPvcNameRook, testingNamespaceRook.Name, VolumeSize, CephFsSc, corev1.ReadWriteMany))
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() bool {
				pvcOut, err := hyperStorageHelper.Clientset.CoreV1().PersistentVolumeClaims(testingNamespaceRook.Name).
					Get(pvc.Name, metav1.GetOptions{})
				if err == nil && pvcOut.Status.Phase == corev1.ClaimBound {
					fmt.Printf("Pvc %s is created and Bound\n", pvcOut.Name)
					return true
				}
				fmt.Printf("Pvc %s is still creating in phase %s \n", pvcOut.Name, pvcOut.Status.Phase)
				return false
			}, TimeOutForCreatingPvc, PollingIntervalDefault).Should(BeTrue())

			err = hyperStorageHelper.Clientset.CoreV1().PersistentVolumeClaims(testingNamespaceRook.Name).
				Delete(pvc.Name, &metav1.DeleteOptions{})
			Expect(err).ToNot(HaveOccurred())

		})
	})

	Describe("[[TEST][e2e][0003] Create RBD PVC Create -> Get -> Delete]", func() {
		It("Create RBD PVC", func() {
			testIDRook = "0003"
			testPvcNameRook = CephPvcNamePrefix + testIDRook
			fmt.Printf("[TEST][e2e][%s] started\n", testIDRook)

			pvc, err := createPvcInStorageClass(hyperStorageHelper.Clientset,
				makePvcInStorageClassSpec(testPvcNameRook, testingNamespaceRook.Name, VolumeSize, RbdSc, corev1.ReadWriteOnce))
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() bool {
				pvcOut, err := hyperStorageHelper.Clientset.CoreV1().PersistentVolumeClaims(testingNamespaceRook.Name).
					Get(pvc.Name, metav1.GetOptions{})
				if err == nil && pvcOut.Status.Phase == corev1.ClaimBound {
					fmt.Printf("Pvc %s is created and Bound\n", pvcOut.Name)
					return true
				}
				fmt.Printf("Pvc %s is still creating in phase %s \n", pvcOut.Name, pvcOut.Status.Phase)
				return false
			}, TimeOutForCreatingPvc, PollingIntervalDefault).Should(BeTrue())

			err = hyperStorageHelper.Clientset.CoreV1().PersistentVolumeClaims(testingNamespaceRook.Name).
				Delete(pvc.Name, &metav1.DeleteOptions{})
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
