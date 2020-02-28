package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cdiv1alpha1 "kubevirt.io/containerized-data-importer/pkg/apis/core/v1alpha1"
	dvutils "kubevirt.io/containerized-data-importer/tests/utils"
	"log"
	"strconv"
	"time"
)

var (
	testingNamespaceCDI *corev1.Namespace
	testIDCDI           int
	testDvName          string
	dataVolumeSize      string
)

const (
	CdiTestingNamespacePrefix = "test-cdi-"
	DataVolumeNamePrefix      = "test-dv-"

	DataVolumeSize = "5Gi" // TODO >=5Gi && =<sc 가용량 을 get 하여 설정하도록

	TimeOutForCreatingDv  = 500 * time.Second
	TimeoutForDeletingDv  = 300 * time.Second
	TimeOutForCreatingPvc = 60 * time.Second

	StorageClassCephfs = "csi-cephfs-sc"

	SampleRegistryURL = "docker://kubevirt/fedora-cloud-registry-disk-demo"
	SampleHTTPURL     = "https://download.cirros-cloud.net/contrib/0.3.0/cirros-0.3.0-i386-disk.img"
)

var _ = Describe("Test CDI Module", func() {
	BeforeEach(func() {
		testIDCDI++
		//TODO glog+파일로거로 변경
		log.Printf("========== [TEST][CDI][CASE-#%d] Started ==========\n", testIDCDI)

		// create testing namespace
		testingNamespaceCDI, err = createNamespace(hyperStorageHelper.Clientset, makeNamespaceSpec(CdiTestingNamespacePrefix))
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		//delete testing namespace
		err = hyperStorageHelper.Clientset.CoreV1().Namespaces().Delete(testingNamespaceCDI.Name, &metav1.DeleteOptions{})
		Expect(err).ToNot(HaveOccurred())
		Eventually(func() bool {
			ns, err := hyperStorageHelper.Clientset.CoreV1().Namespaces().Get(testingNamespaceCDI.Name, metav1.GetOptions{})
			if err != nil || errors.IsNotFound(err) {
				return true
			}
			if ns.Status.Phase == corev1.NamespaceTerminating {
				log.Printf("Namespace %s is still in phase %s\n", testingNamespaceCDI.Name, ns.Status.Phase)
				return false
			}
			return false
		}, TimeoutForDeletingNamespace, PollingIntervalForDeletingNamespace).Should(BeTrue())
		log.Printf("========== [TEST][CDI][CASE-#%d] Finished ==========\n", testIDCDI)
	})

	Describe("[[TEST][CDI] DataVolume Create -> Get -> Delete]", func() {
		It("Create DataVolume from registry", func() {
			testDvName = DataVolumeNamePrefix + strconv.Itoa(testIDCDI)
			dataVolumeSize = DataVolumeSize

			// create dv
			dv, err := dvutils.CreateDataVolumeFromDefinition(hyperStorageHelper.CdiClientset, testingNamespaceCDI.Name,
				makeDataVolumeSpec(testDvName, dataVolumeSize, makeDataVolumeSourceRegistry(SampleRegistryURL),
					corev1.ReadWriteMany))
			Expect(err).ToNot(HaveOccurred())
			log.Printf("dv %s is creating\n", testDvName)

			err = waitDataVolumeGetReadyThenDelete(dv, hyperStorageHelper)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("[[TEST][CDI] DataVolume Create -> Get -> Delete]", func() {
		It("Create DataVolume from http", func() {
			testDvName = DataVolumeNamePrefix + strconv.Itoa(testIDCDI)
			dataVolumeSize = DataVolumeSize

			// create DV 를 from http
			dv, err := dvutils.CreateDataVolumeFromDefinition(hyperStorageHelper.CdiClientset, testingNamespaceCDI.Name,
				makeDataVolumeSpec(testDvName, dataVolumeSize, makeDataVolumeSourceHTTP(SampleHTTPURL), corev1.ReadWriteMany))
			Expect(err).ToNot(HaveOccurred())
			log.Printf("dv %s is creating\n", testDvName)

			err = waitDataVolumeGetReadyThenDelete(dv, hyperStorageHelper)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("[[TEST][CDI] DataVolume Create -> Get -> Delete]", func() {
		It("Create DataVolume from pvc", func() {
			testDvName = DataVolumeNamePrefix + strconv.Itoa(testIDCDI)
			pvcToBeClonedName := "pvctobecloned"
			dataVolumeSize = DataVolumeSize

			// create pvc-original with some sc
			// TODO block sc 에 대해서도 테스트 추가
			pvc, err := createPvcInStorageClass(hyperStorageHelper.Clientset,
				makePvcInStorageClassSpec(pvcToBeClonedName, testingNamespaceCDI.Name, dataVolumeSize, StorageClassCephfs,
					corev1.ReadOnlyMany))
			Expect(err).ToNot(HaveOccurred())
			Eventually(func() bool {
				pvcOut, err := hyperStorageHelper.Clientset.CoreV1().PersistentVolumeClaims(testingNamespaceCDI.Name).
					Get(pvc.Name, metav1.GetOptions{})
				if err == nil && pvcOut.Status.Phase == corev1.ClaimBound {
					log.Printf("Pvc %s is created and Bound\n", pvcOut.Name)
					return true
				}
				log.Printf("Pvc %s is still creating in phase %s \n", pvcOut.Name, pvcOut.Status.Phase)
				return false
			}, TimeOutForCreatingPvc, PollingIntervalDefault).Should(BeTrue())

			//clone dv from pvc-original
			dv, err := dvutils.CreateDataVolumeFromDefinition(hyperStorageHelper.CdiClientset, testingNamespaceCDI.Name,
				makeDataVolumeSpec(testDvName, dataVolumeSize,
					makeDataVolumeSourcePVC(testingNamespaceCDI.Name, pvc.Name), corev1.ReadWriteMany))
			Expect(err).ToNot(HaveOccurred())
			log.Printf("dv %s is creating\n", testDvName)

			err = waitDataVolumeGetReadyThenDelete(dv, hyperStorageHelper)
			Expect(err).ToNot(HaveOccurred())

			// delete pvc-original
			err = hyperStorageHelper.Clientset.CoreV1().PersistentVolumeClaims(testingNamespaceCDI.Name).
				Delete(pvc.Name, &metav1.DeleteOptions{})
			Expect(err).ToNot(HaveOccurred())
			// TODO delete pv (optional if reclaimPolicy of sc is retain)
		})
	})
})

func waitDataVolumeGetReadyThenDelete(dv *cdiv1alpha1.DataVolume, hyperStorageHelper *HyperHelper) error {
	// wait dv until succeeded
	log.Printf("wait until dv phase become succeeded for timeout %s\n", TimeOutForCreatingDv.String())
	err := dvutils.WaitForDataVolumePhaseWithTimeout(hyperStorageHelper.CdiClientset, dv.Namespace,
		cdiv1alpha1.Succeeded, dv.Name, TimeOutForCreatingDv)
	Expect(err).ToNot(HaveOccurred())
	log.Printf("dv %s is created\n", dv.Name)

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
		log.Printf("DataVolume %s is still deleting...\n", dv.Name)
		return false
	}, TimeoutForDeletingDv, PollingIntervalDefault).Should(BeTrue())

	return err
}

func makeDataVolumeSpec(name string, size string, source *cdiv1alpha1.DataVolumeSource,
	accessMode corev1.PersistentVolumeAccessMode) *cdiv1alpha1.DataVolume {
	return &cdiv1alpha1.DataVolume{
		TypeMeta: metav1.TypeMeta{
			Kind:       "cdi.kubevirt.io/v1alpha1",
			APIVersion: "DataVolume",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: cdiv1alpha1.DataVolumeSpec{
			Source: *source,
			PVC: &corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{accessMode},
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: resource.MustParse(size),
					},
				},
			},
		},
	}
}

func makeDataVolumeSourceHTTP(url string) *cdiv1alpha1.DataVolumeSource {
	return &cdiv1alpha1.DataVolumeSource{
		HTTP: &cdiv1alpha1.DataVolumeSourceHTTP{URL: url},
	}
}

func makeDataVolumeSourceRegistry(url string) *cdiv1alpha1.DataVolumeSource {
	return &cdiv1alpha1.DataVolumeSource{
		Registry: &cdiv1alpha1.DataVolumeSourceRegistry{URL: url},
	}
}

func makeDataVolumeSourcePVC(namespace string, name string) *cdiv1alpha1.DataVolumeSource {
	return &cdiv1alpha1.DataVolumeSource{
		PVC: &cdiv1alpha1.DataVolumeSourcePVC{
			Namespace: namespace,
			Name:      name,
		},
	}
}
