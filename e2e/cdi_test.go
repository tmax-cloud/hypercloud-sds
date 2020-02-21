package tests

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cdiv1alpha1 "kubevirt.io/containerized-data-importer/pkg/apis/core/v1alpha1"
	dvutils "kubevirt.io/containerized-data-importer/tests/utils"
)

var (
	testingNamespace *corev1.Namespace
	err              error
	testID           string
	testDvName       string
	scList           *v1.StorageClassList
	scItem           v1.StorageClass
)

var _ = Describe("Test CDI Module", func() {
	BeforeEach(func() {
		// create testing namespace
		testingNamespace, err = createNamespace(hyperStorageHelper.Clientset, makeNamespaceSpec(CdiTestingNamespacePrefix))
		Expect(err).ToNot(HaveOccurred())
		fmt.Printf("Namespace %s is created for testing.\n", testingNamespace.Name) // TODO fmt 대신 log 사용

		// TODO storage-v1 만 확인해도 괜찮은지?
		// TODO RWM, RWO sc 구분하여 변수로 저장 후 다른 테스트에서 사용
		// TODO reclaimPolicy Retain 과 Delete sc 구분
		// TODO dynamic sc 구분
		// TODO 현재 사용 가능 여부 구분
		scList, err = hyperStorageHelper.Clientset.StorageV1().StorageClasses().List(metav1.ListOptions{})

		for _, scItem = range scList.Items {
			fmt.Printf("One of %d existing storageclasses is : %s \n", len(scList.Items), scItem.Name)
		}
	})

	AfterEach(func() {
		//delete testing namespace
		err = hyperStorageHelper.Clientset.CoreV1().Namespaces().Delete(testingNamespace.Name, &metav1.DeleteOptions{})
		Expect(err).ToNot(HaveOccurred())
		Eventually(func() bool {
			ns, err := hyperStorageHelper.Clientset.CoreV1().Namespaces().Get(testingNamespace.Name, metav1.GetOptions{})
			if err != nil || errors.IsNotFound(err) {
				return true
			}
			if ns.Status.Phase == corev1.NamespaceTerminating {
				fmt.Printf("Namespace %s is still in phase %s\n", testingNamespace.Name, ns.Status.Phase)
				return false
			}
			return false
		}, TimeoutForDeletingNamespace, PollingIntervalForDeletingNamespace).Should(BeTrue())
		fmt.Printf("\n %s DataVolume Create => Get => Delete\n", testID)
	})

	Describe("[[TEST][e2e][0002] DataVolume Create -> Get -> Delete]", func() {
		It("Create DataVolume from registry", func() {
			testID = "0002"
			testDvName = DataVolumeNamePrefix + testID // TODO test 별로 매번 이렇게 set하지 않도록 변경
			fmt.Printf("[TEST][e2e][%s] started\n", testID)

			// create dv
			dv, err := dvutils.CreateDataVolumeFromDefinition(hyperStorageHelper.CdiClientset, testingNamespace.Name,
				makeDataVolumeSpec(testDvName, DataVolumeSize, makeDataVolumeSourceRegistry(SampleRegistryURL)))
			Expect(err).ToNot(HaveOccurred())
			fmt.Printf("dv %s is creating\n", testDvName)

			err = waitDataVolumeGetReadyThenDelete(dv, hyperStorageHelper)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("[[TEST][e2e][0003] DataVolume Create -> Get -> Delete]", func() {
		It("Create DataVolume from http", func() {
			testID = "0003"
			testDvName = DataVolumeNamePrefix + testID
			fmt.Printf("[TEST][e2e][%s] started\n", testID)

			// create DV 를 from http
			dv, err := dvutils.CreateDataVolumeFromDefinition(hyperStorageHelper.CdiClientset, testingNamespace.Name,
				makeDataVolumeSpec(testDvName, DataVolumeSize, makeDataVolumeSourceHTTP(SampleHTTPURL)))
			Expect(err).ToNot(HaveOccurred())
			fmt.Printf("dv %s is creating\n", testDvName)

			err = waitDataVolumeGetReadyThenDelete(dv, hyperStorageHelper)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("[[TEST][e2e][0004] DataVolume Create -> Get -> Delete]", func() {
		It("Create DataVolume from pvc", func() {
			testID = "0004"
			testDvName = DataVolumeNamePrefix + testID
			pvcToBeClonedName := "pvctobecloned"
			fmt.Printf("[TEST][e2e][%s] started\n", testID)

			// create pvc-original with some sc
			// TODO block sc 에 대해서도 테스트 추가
			pvc, err := createPvcInStorageClass(hyperStorageHelper.Clientset,
				makePvcInStorageClassSpec(pvcToBeClonedName, testingNamespace.Name, DataVolumeSize, StorageClassCephfs))
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() bool {
				pvcOut, err := hyperStorageHelper.Clientset.CoreV1().PersistentVolumeClaims(testingNamespace.Name).
					Get(pvc.Name, metav1.GetOptions{})
				if err == nil && pvcOut.Status.Phase == corev1.ClaimBound {
					fmt.Printf("Pvc %s is created and Bound\n", pvcOut.Name)
					return true
				}
				fmt.Printf("Pvc %s is still creating in phase %s \n", pvcOut.Name, pvcOut.Status.Phase)
				return false
			}, TimeOutForCreatingPvc, PollingIntervalDefault).Should(BeTrue())

			//clone dv from pvc-original
			dv, err := dvutils.CreateDataVolumeFromDefinition(hyperStorageHelper.CdiClientset, testingNamespace.Name,
				makeDataVolumeSpec(testDvName, DataVolumeSize, makeDataVolumeSourcePVC(testingNamespace.Name, pvc.Name)))
			Expect(err).ToNot(HaveOccurred())
			fmt.Printf("dv %s is creating\n", testDvName)

			err = waitDataVolumeGetReadyThenDelete(dv, hyperStorageHelper)
			Expect(err).ToNot(HaveOccurred())

			// delete pvc-original
			err = hyperStorageHelper.Clientset.CoreV1().PersistentVolumeClaims(testingNamespace.Name).
				Delete(pvc.Name, &metav1.DeleteOptions{})
			Expect(err).ToNot(HaveOccurred())
			// TODO delete pv (optional if reclaimPolicy of sc is retain)
		})
	})
})

func waitDataVolumeGetReadyThenDelete(dv *cdiv1alpha1.DataVolume, hyperStorageHelper *HyperHelper) error {
	// wait dv until succeeded
	fmt.Printf("wait until dv phase become succeeded for timeout %s\n", TimeOutForCreatingDv.String())
	err := dvutils.WaitForDataVolumePhaseWithTimeout(hyperStorageHelper.CdiClientset, dv.Namespace,
		cdiv1alpha1.Succeeded, dv.Name, TimeOutForCreatingDv)
	Expect(err).ToNot(HaveOccurred())
	fmt.Printf("dv %s is created\n", dv.Name)

	// get dv and check
	out, err := hyperStorageHelper.CdiClientset.CdiV1alpha1().DataVolumes(dv.Namespace).
		Get(dv.Name, metav1.GetOptions{})
	Expect(err).ToNot(HaveOccurred())
	Expect(out.Name).To(Equal(dv.Name))
	Expect(out.Status.Phase).To(Equal(cdiv1alpha1.Succeeded))

	// delete dv
	err = hyperStorageHelper.CdiClientset.CdiV1alpha1().DataVolumes(dv.Namespace).
		Delete(dv.Name, &metav1.DeleteOptions{})
	Expect(err).ToNot(HaveOccurred())
	Eventually(func() bool {
		_, err := hyperStorageHelper.CdiClientset.CdiV1alpha1().DataVolumes(dv.Namespace).
			Get(dv.Namespace, metav1.GetOptions{})
		if err != nil || errors.IsNotFound(err) {
			return true
		}
		fmt.Printf("DataVolume %s is still deleting...\n", dv.Name)
		return false
	}, TimeoutForDeletingDv, PollingIntervalDefault).Should(BeTrue())

	return err
}
