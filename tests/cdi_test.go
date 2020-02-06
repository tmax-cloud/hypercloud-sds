package tests

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cdiv1 "kubevirt.io/containerized-data-importer/pkg/apis/core/v1alpha1"
	dvutils "kubevirt.io/containerized-data-importer/tests/utils"
	"time"
)

// TODO refactoring & assertion
var _ = Describe("Test CDI Module", func() {
	Describe("[[TEST][e2e][0002] DataVolume Create -> Get -> Delete", func() {
		It("Create -> Get -> Delete", func() {
			const (
				namespacePrefix                     = "test-0002-create-dv-"
				dvName                              = "test-0002-dv"
				dvSize                              = "5Gi"
				timeOutForCreatingDv                = 500 * time.Second
				timeoutForDeletingDv                = 300 * time.Second
				pollingIntervalForDeletingDv        = 5 * time.Second
				timeoutForDeletingNamespace         = 500 * time.Second
				pollingIntervalForDeletingNamespace = 10 * time.Second
			)
			// create ns
			nsSpec := makeNamespaceSpec(namespacePrefix)
			generatedNs, err := createNamespace(hyperStorageHelper.Clientset, nsSpec)
			Expect(err).ToNot(HaveOccurred())
			fmt.Printf("namespace %s is created\n", generatedNs.Name) // TODO fmt 대신 log 사용

			// create dv
			dv, err := dvutils.CreateDataVolumeFromDefinition(hyperStorageHelper.CdiClientset, generatedNs.Name,
				makeDataVolumeSpec(dvName, dvSize))
			Expect(err).ToNot(HaveOccurred())
			fmt.Printf("dv %s is creating\n", dvName)

			// wait dv until succeeded
			fmt.Println("wait until dv phase become succeeded for timeout")
			err = dvutils.WaitForDataVolumePhaseWithTimeout(hyperStorageHelper.CdiClientset, generatedNs.Name,
				cdiv1.Succeeded, dv.Name, timeOutForCreatingDv)
			Expect(err).ToNot(HaveOccurred())
			fmt.Printf("dv %s is created\n", dvName)

			// get dv and check
			out, err := hyperStorageHelper.CdiClientset.CdiV1alpha1().DataVolumes(dv.Namespace).
				Get(dv.Name, metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(out.Name).To(Equal(dvName))
			Expect(out.Status.Phase).To(Equal(cdiv1.Succeeded))

			// delete dv
			err = hyperStorageHelper.CdiClientset.CdiV1alpha1().DataVolumes(dv.Namespace).
				Delete(dv.Name, &metav1.DeleteOptions{})
			Expect(err).ToNot(HaveOccurred())
			Eventually(func() bool {
				_, err := hyperStorageHelper.CdiClientset.CdiV1alpha1().DataVolumes(dv.Namespace).
					Get(generatedNs.Name, metav1.GetOptions{})
				if err != nil || errors.IsNotFound(err) {
					return true
				}
				// TODO dv get 한 status.phase 보고 좀 더 자세하게 wait 할 지 말 지 결정
				fmt.Printf("DataVolume %s is still deleting...\n", dvName)
				return false
			}, timeoutForDeletingDv, pollingIntervalForDeletingDv).Should(BeTrue())

			//delete ns
			//TODO MUST DELETE NAMESPACE regardless of any above assertions !!!!!!!!!!!!!!
			// 현재는 위의 모든 assertion 이 true 일 때만, 아래 라인을 타고 namespace 를 delete 하고 있음
			err = hyperStorageHelper.Clientset.CoreV1().Namespaces().Delete(generatedNs.Name, &metav1.DeleteOptions{})
			Expect(err).ToNot(HaveOccurred())
			Eventually(func() bool {
				ns, err := hyperStorageHelper.Clientset.CoreV1().Namespaces().Get(generatedNs.Name, metav1.GetOptions{})
				if err != nil || errors.IsNotFound(err) {
					return true
				}

				if ns.Status.Phase == corev1.NamespaceTerminating {
					fmt.Printf("Namespace %s is still in phase %s\n", generatedNs.Name, ns.Status.Phase)
					return false
				}
				return false
			}, timeoutForDeletingNamespace, pollingIntervalForDeletingNamespace).Should(BeTrue())
			fmt.Printf("\n[TEST0005][FINISHED] DataVolume Create => Get => Delete\n")
		})
	})
})
