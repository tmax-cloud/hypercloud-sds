package tests

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var _ = Describe("TEST", func() {
	var (
		clientset *kubernetes.Clientset
	)

	BeforeSuite(func() {
		clientset = getClientSet() // Before each 시 flag.go 의 Happens only if flags are declared with identical names 에러 발생
	})

	// TODO Temporary codes for nodes & pods check
	Describe("How many", func() {
		Context("Nodes", func() {
			It("should not be zero", func() {
				nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
				if err != nil {
					panic(err.Error())
				}
				fmt.Printf("There are %d nodes in the cluster\n", len(nodes.Items))

				Expect(len(nodes.Items)).NotTo(Equal(0)) // TODO should change to another assertion
			})
		})

		Context("Pods", func() {
			It("should not be zero ", func() {
				pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
				if err != nil {
					panic(err.Error())
				}
				fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

				Expect(len(pods.Items)).NotTo(Equal(0)) // TODO should change to another assertion
			})
		})
	})

	Describe("rook install check", func() {
		var (
			namespace                               = "rook-ceph"
			deployment_rook_ceph_operator           = "rook-ceph-operator"
			deployment_csi_cephfsplugin_provisioner = "csi-cephfsplugin-provisioner"
			deployment_csi_rbdplugin_provisioner    = "csi-rbdplugin-provisioner"
		)

		Context("check if each deployment installed", func() {
			It("deployment_rook_ceph_operator", func() {
				out, err := clientset.AppsV1().Deployments(namespace).Get(deployment_rook_ceph_operator, metav1.GetOptions{})

				Expect(errors.IsNotFound(err)).NotTo(Equal(true))
				Expect(out.Status.ReadyReplicas).ShouldNot(BeNumerically("==", 0))
				Expect(out.Status.UnavailableReplicas).Should(BeNumerically("==", 0))

				fmt.Printf("Found deployment %s in namespace %s\n", deployment_rook_ceph_operator, namespace)
			})
		})
		Context("check if each deployment installed", func() {
			It("deployment_csi_cephfsplugin_provisioner", func() {
				out, err := clientset.AppsV1().Deployments(namespace).Get(deployment_csi_cephfsplugin_provisioner, metav1.GetOptions{})

				Expect(errors.IsNotFound(err)).NotTo(Equal(true))
				Expect(out.Status.ReadyReplicas).ShouldNot(BeNumerically("==", 0))
				Expect(out.Status.UnavailableReplicas).Should(BeNumerically("==", 0))

				fmt.Printf("Found deployment %s in namespace %s\n", deployment_csi_cephfsplugin_provisioner, namespace)
			})
		})

		Context("check if each deployment installed", func() {
			It("deployment_csi_rbdplugin_provisioner", func() {
				out, err := clientset.AppsV1().Deployments(namespace).Get(deployment_csi_rbdplugin_provisioner, metav1.GetOptions{})

				Expect(errors.IsNotFound(err)).NotTo(Equal(true))
				Expect(out.Status.ReadyReplicas).ShouldNot(BeNumerically("==", 0))
				Expect(out.Status.UnavailableReplicas).Should(BeNumerically("==", 0))

				fmt.Printf("Found deployment %s in namespace %s\n", deployment_csi_rbdplugin_provisioner, namespace)
			})
		})
	})
	Describe("cdi install check", func() {
		var (
			namespace                  = "cdi"
			deployment_cdi_operator    = "cdi-operator"
			deployment_cdi_deployment  = "cdi-deployment"
			deployment_cdi_apiserver   = "cdi-apiserver"
			deployment_cdi_uploadproxy = "cdi-uploadproxy"
		)

		Context("check if each deployment installed", func() {
			It("deployment_cdi_operator", func() {
				out, err := clientset.AppsV1().Deployments(namespace).Get(deployment_cdi_operator, metav1.GetOptions{})

				Expect(errors.IsNotFound(err)).NotTo(Equal(true))
				Expect(out.Status.ReadyReplicas).ShouldNot(BeNumerically("==", 0))
				Expect(out.Status.UnavailableReplicas).Should(BeNumerically("==", 0))

				fmt.Printf("Found deployment %s in namespace %s\n", deployment_cdi_operator, namespace)
			})
		})
		Context("check if each deployment installed", func() {
			It("deployment_cdi_deployment", func() {
				out, err := clientset.AppsV1().Deployments(namespace).Get(deployment_cdi_deployment, metav1.GetOptions{})

				Expect(errors.IsNotFound(err)).NotTo(Equal(true))
				Expect(out.Status.ReadyReplicas).ShouldNot(BeNumerically("==", 0))
				Expect(out.Status.UnavailableReplicas).Should(BeNumerically("==", 0))

				fmt.Printf("Found deployment %s in namespace %s\n", deployment_cdi_deployment, namespace)
			})
		})

		Context("check if each deployment installed", func() {
			It("deployment_cdi_apiserver", func() {
				out, err := clientset.AppsV1().Deployments(namespace).Get(deployment_cdi_apiserver, metav1.GetOptions{})

				Expect(errors.IsNotFound(err)).NotTo(Equal(true))
				Expect(out.Status.ReadyReplicas).ShouldNot(BeNumerically("==", 0))
				Expect(out.Status.UnavailableReplicas).Should(BeNumerically("==", 0))

				fmt.Printf("Found deployment %s in namespace %s\n", deployment_cdi_apiserver, namespace)
			})
		})

		Context("check if each deployment installed", func() {
			It("deployment_cdi_uploadproxy", func() {
				out, err := clientset.AppsV1().Deployments(namespace).Get(deployment_cdi_uploadproxy, metav1.GetOptions{})

				Expect(errors.IsNotFound(err)).NotTo(Equal(true))
				Expect(out.Status.ReadyReplicas).ShouldNot(BeNumerically("==", 0))
				Expect(out.Status.UnavailableReplicas).Should(BeNumerically("==", 0))

				fmt.Printf("Found deployment %s in namespace %s\n", deployment_cdi_uploadproxy, namespace)
			})
		})
	})
})
