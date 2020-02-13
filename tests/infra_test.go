package tests

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("TEST", func() {
	BeforeEach(func() {
		nodes, err := hyperStorageHelper.Clientset.CoreV1().Nodes().List(metav1.ListOptions{})
		Expect(err).ToNot(HaveOccurred())
		Expect(len(nodes.Items)).NotTo(Equal(0)) // TODO should change to another assertion
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
			namespace = "rook-ceph"
		)

		Context("check if each deployment installed", func() {
			It("[TEST - 03] deployment_rook_ceph_operator", func() {
				out, err := hyperStorageHelper.Clientset.AppsV1().Deployments(namespace).
					Get(DeploymentRookCephOperator, metav1.GetOptions{})

				Expect(err).ToNot(HaveOccurred())
				Expect(out.Status.ReadyReplicas).ShouldNot(BeNumerically("==", 0))
				Expect(out.Status.UnavailableReplicas).Should(BeNumerically("==", 0))

				//TODO fmt -> log 사용 및 관리 필요
				fmt.Printf("Found deployment %s in namespace %s\n", DeploymentRookCephOperator, namespace)
			})
		})
		Context("check if each deployment installed", func() {
			It("[TEST - 04] deployment_csi_cephfsplugin_provisioner", func() {
				out, err := hyperStorageHelper.Clientset.AppsV1().Deployments(namespace).
					Get(DeploymentCsiCephfspluginProvisioner, metav1.GetOptions{})

				Expect(err).ToNot(HaveOccurred())
				Expect(out.Status.ReadyReplicas).ShouldNot(BeNumerically("==", 0))
				Expect(out.Status.UnavailableReplicas).Should(BeNumerically("==", 0))

				fmt.Printf("Found deployment %s in namespace %s\n", DeploymentCsiCephfspluginProvisioner, namespace)
			})
		})

		Context("check if each deployment installed", func() {
			It("[TEST - 05] deployment_csi_rbdplugin_provisioner", func() {
				out, err := hyperStorageHelper.Clientset.AppsV1().Deployments(namespace).
					Get(DeploymentCsiRbdpluginProvisioner, metav1.GetOptions{})

				Expect(err).ToNot(HaveOccurred())
				Expect(out.Status.ReadyReplicas).ShouldNot(BeNumerically("==", 0))
				Expect(out.Status.UnavailableReplicas).Should(BeNumerically("==", 0))

				fmt.Printf("Found deployment %s in namespace %s\n", DeploymentCsiRbdpluginProvisioner, namespace)
			})
		})
	})
	Describe("cdi install check", func() {
		const (
			namespace = "cdi"
		)

		Context("check if each deployment installed", func() {
			It("[TEST - 06] deployment_cdi_operator", func() {
				out, err := hyperStorageHelper.Clientset.AppsV1().Deployments(namespace).
					Get(DeploymentCdiOperator, metav1.GetOptions{})

				Expect(err).ToNot(HaveOccurred())
				Expect(out.Status.ReadyReplicas).ShouldNot(BeNumerically("==", 0))
				Expect(out.Status.UnavailableReplicas).Should(BeNumerically("==", 0))

				fmt.Printf("Found deployment %s in namespace %s\n", DeploymentCdiOperator, namespace)
			})
		})
		Context("check if each deployment installed", func() {
			It("[TEST - 07] deployment_cdi_deployment", func() {
				out, err := hyperStorageHelper.Clientset.AppsV1().Deployments(namespace).
					Get(DeploymentCdiDeployment, metav1.GetOptions{})

				Expect(err).ToNot(HaveOccurred())
				Expect(out.Status.ReadyReplicas).ShouldNot(BeNumerically("==", 0))
				Expect(out.Status.UnavailableReplicas).Should(BeNumerically("==", 0))

				fmt.Printf("Found deployment %s in namespace %s\n", DeploymentCdiDeployment, namespace)
			})
		})

		Context("check if each deployment installed", func() {
			It("[TEST - 08] deployment_cdi_apiserver", func() {
				out, err := hyperStorageHelper.Clientset.AppsV1().Deployments(namespace).
					Get(DeploymentCdiApiserver, metav1.GetOptions{})

				Expect(err).ToNot(HaveOccurred())
				Expect(out.Status.ReadyReplicas).ShouldNot(BeNumerically("==", 0))
				Expect(out.Status.UnavailableReplicas).Should(BeNumerically("==", 0))

				fmt.Printf("Found deployment %s in namespace %s\n", DeploymentCdiApiserver, namespace)
			})
		})

		Context("check if each deployment installed", func() {
			It("[TEST - 09] deployment_cdi_uploadproxy", func() {
				out, err := hyperStorageHelper.Clientset.AppsV1().Deployments(namespace).
					Get(DeploymentCdiUploadproxy, metav1.GetOptions{})

				Expect(err).ToNot(HaveOccurred())
				Expect(out.Status.ReadyReplicas).ShouldNot(BeNumerically("==", 0))
				Expect(out.Status.UnavailableReplicas).Should(BeNumerically("==", 0))

				fmt.Printf("Found deployment %s in namespace %s\n", DeploymentCdiUploadproxy, namespace)
			})
		})
	})
})
