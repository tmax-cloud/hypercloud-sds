package tests

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//TODO infra check, rook installation check, cdi installation check, pod2pod network check 모두 다른 _test.go 파일로 분리
var _ = Describe("TEST", func() {

	BeforeSuite(func() { //TODO BeforeEach 로 clientSet 사용 시,
		hyperStorageHelper = HyperStorageHelper()
		// flag.go 의 Happens only if flags are declared with identical names 에러 발생
		//clientset, config = HyperStorageHelper().Clientset, HyperStorageConfig()
	})

	// TODO Temporary codes for nodes & pods check
	Describe("nodes & pods check", func() {
		Context("Number of Nodes", func() {
			It("[TEST - 01] should not be zero", func() {
				nodes, err := hyperStorageHelper.Clientset.CoreV1().Nodes().List(metav1.ListOptions{})
				Expect(err).ToNot(HaveOccurred())
				Expect(len(nodes.Items)).NotTo(Equal(0)) // TODO should change to another assertion
			})
		})

		Context("Number of Pods", func() {
			It("[TEST - 02] should not be zero ", func() {
				pods, err := hyperStorageHelper.Clientset.CoreV1().Pods("").List(metav1.ListOptions{})
				Expect(err).ToNot(HaveOccurred())
				Expect(len(pods.Items)).NotTo(Equal(0)) // TODO should change to another assertion
			})
		})
	})

	Describe("rook install check", func() {
		const (
			namespace                            = "rook-ceph"
			deploymentRookCephOperator           = "rook-ceph-operator"
			deploymentCsiCephfspluginProvisioner = "csi-cephfsplugin-provisioner"
			deploymentCsiRbdpluginProvisioner    = "csi-rbdplugin-provisioner"
		)

		Context("check if each deployment installed", func() {
			It("[TEST - 03] deployment_rook_ceph_operator", func() {
				out, err := hyperStorageHelper.Clientset.AppsV1().Deployments(namespace).
					Get(deploymentRookCephOperator, metav1.GetOptions{})

				Expect(err).ToNot(HaveOccurred())
				Expect(out.Status.ReadyReplicas).ShouldNot(BeNumerically("==", 0))
				Expect(out.Status.UnavailableReplicas).Should(BeNumerically("==", 0))

				//TODO fmt -> log 사용 및 관리 필요
				fmt.Printf("Found deployment %s in namespace %s\n", deploymentRookCephOperator, namespace)
			})
		})
		Context("check if each deployment installed", func() {
			It("[TEST - 04] deployment_csi_cephfsplugin_provisioner", func() {
				out, err := hyperStorageHelper.Clientset.AppsV1().Deployments(namespace).
					Get(deploymentCsiCephfspluginProvisioner, metav1.GetOptions{})

				Expect(err).ToNot(HaveOccurred())
				Expect(out.Status.ReadyReplicas).ShouldNot(BeNumerically("==", 0))
				Expect(out.Status.UnavailableReplicas).Should(BeNumerically("==", 0))

				fmt.Printf("Found deployment %s in namespace %s\n", deploymentCsiCephfspluginProvisioner, namespace)
			})
		})

		Context("check if each deployment installed", func() {
			It("[TEST - 05] deployment_csi_rbdplugin_provisioner", func() {
				out, err := hyperStorageHelper.Clientset.AppsV1().Deployments(namespace).
					Get(deploymentCsiRbdpluginProvisioner, metav1.GetOptions{})

				Expect(err).ToNot(HaveOccurred())
				Expect(out.Status.ReadyReplicas).ShouldNot(BeNumerically("==", 0))
				Expect(out.Status.UnavailableReplicas).Should(BeNumerically("==", 0))

				fmt.Printf("Found deployment %s in namespace %s\n", deploymentCsiRbdpluginProvisioner, namespace)
			})
		})
	})
	Describe("cdi install check", func() {
		const (
			namespace                = "cdi"
			deploymentCdiOperator    = "cdi-operator"
			deploymentCdiDeployment  = "cdi-deployment"
			deploymentCdiApiserver   = "cdi-apiserver"
			deploymentCdiUploadproxy = "cdi-uploadproxy"
		)

		Context("check if each deployment installed", func() {
			It("[TEST - 06] deployment_cdi_operator", func() {
				out, err := hyperStorageHelper.Clientset.AppsV1().Deployments(namespace).
					Get(deploymentCdiOperator, metav1.GetOptions{})

				Expect(err).ToNot(HaveOccurred())
				Expect(out.Status.ReadyReplicas).ShouldNot(BeNumerically("==", 0))
				Expect(out.Status.UnavailableReplicas).Should(BeNumerically("==", 0))

				fmt.Printf("Found deployment %s in namespace %s\n", deploymentCdiOperator, namespace)
			})
		})
		Context("check if each deployment installed", func() {
			It("[TEST - 07] deployment_cdi_deployment", func() {
				out, err := hyperStorageHelper.Clientset.AppsV1().Deployments(namespace).
					Get(deploymentCdiDeployment, metav1.GetOptions{})

				Expect(err).ToNot(HaveOccurred())
				Expect(out.Status.ReadyReplicas).ShouldNot(BeNumerically("==", 0))
				Expect(out.Status.UnavailableReplicas).Should(BeNumerically("==", 0))

				fmt.Printf("Found deployment %s in namespace %s\n", deploymentCdiDeployment, namespace)
			})
		})

		Context("check if each deployment installed", func() {
			It("[TEST - 08] deployment_cdi_apiserver", func() {
				out, err := hyperStorageHelper.Clientset.AppsV1().Deployments(namespace).
					Get(deploymentCdiApiserver, metav1.GetOptions{})

				Expect(err).ToNot(HaveOccurred())
				Expect(out.Status.ReadyReplicas).ShouldNot(BeNumerically("==", 0))
				Expect(out.Status.UnavailableReplicas).Should(BeNumerically("==", 0))

				fmt.Printf("Found deployment %s in namespace %s\n", deploymentCdiApiserver, namespace)
			})
		})

		Context("check if each deployment installed", func() {
			It("[TEST - 09] deployment_cdi_uploadproxy", func() {
				out, err := hyperStorageHelper.Clientset.AppsV1().Deployments(namespace).
					Get(deploymentCdiUploadproxy, metav1.GetOptions{})

				Expect(err).ToNot(HaveOccurred())
				Expect(out.Status.ReadyReplicas).ShouldNot(BeNumerically("==", 0))
				Expect(out.Status.UnavailableReplicas).Should(BeNumerically("==", 0))

				fmt.Printf("Found deployment %s in namespace %s\n", deploymentCdiUploadproxy, namespace)
			})
		})
	})
})
