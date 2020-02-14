package tests

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("TEST basic infra", func() {
	Describe("nodes & pods check", func() {
		It("[TEST - 01] Number of nodes should not be zero", func() {
			nodes, err := hyperStorageHelper.Clientset.CoreV1().Nodes().List(metav1.ListOptions{})

			Expect(err).ToNot(HaveOccurred())
			Expect(len(nodes.Items)).NotTo(Equal(0))
		})
		It("[TEST - 02] Number of pods should not be zero ", func() {
			pods, err := hyperStorageHelper.Clientset.CoreV1().Pods("").List(metav1.ListOptions{})

			Expect(err).ToNot(HaveOccurred())
			Expect(len(pods.Items)).NotTo(Equal(0))
		})
	})

	Describe("rook install check", func() {
		const (
			namespace = "rook-ceph"
		)
		It("[TEST - 03] deployment_rook_ceph_operator", func() {
			out, err := hyperStorageHelper.Clientset.AppsV1().Deployments(namespace).
				Get(DeploymentRookCephOperator, metav1.GetOptions{})

			Expect(err).ToNot(HaveOccurred())
			expectDeploymentAvailable(namespace, out)
		})
		It("[TEST - 04] deployment_csi_cephfsplugin_provisioner", func() {
			out, err := hyperStorageHelper.Clientset.AppsV1().Deployments(namespace).
				Get(DeploymentCsiCephfspluginProvisioner, metav1.GetOptions{})

			Expect(err).ToNot(HaveOccurred())
			expectDeploymentAvailable(namespace, out)
		})
		It("[TEST - 05] deployment_csi_rbdplugin_provisioner", func() {
			out, err := hyperStorageHelper.Clientset.AppsV1().Deployments(namespace).
				Get(DeploymentCsiRbdpluginProvisioner, metav1.GetOptions{})

			Expect(err).ToNot(HaveOccurred())
			expectDeploymentAvailable(namespace, out)
		})
	})

	Describe("cdi install check", func() {
		const (
			namespace = "cdi"
		)
		It("[TEST - 06] deployment_cdi_operator", func() {
			out, err := hyperStorageHelper.Clientset.AppsV1().Deployments(namespace).
				Get(DeploymentCdiOperator, metav1.GetOptions{})

			Expect(err).ToNot(HaveOccurred())
			expectDeploymentAvailable(namespace, out)
		})
		It("[TEST - 07] deployment_cdi_deployment", func() {
			out, err := hyperStorageHelper.Clientset.AppsV1().Deployments(namespace).
				Get(DeploymentCdiDeployment, metav1.GetOptions{})

			Expect(err).ToNot(HaveOccurred())
			expectDeploymentAvailable(namespace, out)
		})
		It("[TEST - 08] deployment_cdi_apiserver", func() {
			out, err := hyperStorageHelper.Clientset.AppsV1().Deployments(namespace).
				Get(DeploymentCdiApiserver, metav1.GetOptions{})

			Expect(err).ToNot(HaveOccurred())
			expectDeploymentAvailable(namespace, out)
		})
		It("[TEST - 09] deployment_cdi_uploadproxy", func() {
			out, err := hyperStorageHelper.Clientset.AppsV1().Deployments(namespace).
				Get(DeploymentCdiUploadproxy, metav1.GetOptions{})

			Expect(err).ToNot(HaveOccurred())
			expectDeploymentAvailable(namespace, out)
		})
	})
})

func expectDeploymentAvailable(namespace string, deployment *v1.Deployment) {
	Expect(deployment.Status.ReadyReplicas).ShouldNot(BeNumerically("==", 0))
	Expect(deployment.Status.UnavailableReplicas).Should(BeNumerically("==", 0))

	//TODO fmt -> log 사용 및 관리 필요
	fmt.Printf("Found deployment %s in namespace %s\n", deployment.Name, namespace)
}
