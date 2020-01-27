package tests

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"time"
)

var _ = Describe("TEST", func() {
	var (
		clientset *kubernetes.Clientset
		config    *restclient.Config
	)

	BeforeSuite(func() {
		clientset, config = getClientSet() // Before each 시 flag.go 의 Happens only if flags are declared with identical names 에러 발생
	})

	// TODO Temporary codes for nodes & pods check
	Describe("How many", func() {
		Context("Nodes", func() {
			It("should not be zero", func() {
				nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
				Expect(err).ToNot(HaveOccurred())
				Expect(len(nodes.Items)).NotTo(Equal(0)) // TODO should change to another assertion
			})
		})

		Context("Pods", func() {
			It("should not be zero ", func() {
				pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
				Expect(err).ToNot(HaveOccurred())
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

				Expect(err).ToNot(HaveOccurred())
				Expect(out.Status.ReadyReplicas).ShouldNot(BeNumerically("==", 0))
				Expect(out.Status.UnavailableReplicas).Should(BeNumerically("==", 0))

				fmt.Printf("Found deployment %s in namespace %s\n", deployment_rook_ceph_operator, namespace)
			})
		})
		Context("check if each deployment installed", func() {
			It("deployment_csi_cephfsplugin_provisioner", func() {
				out, err := clientset.AppsV1().Deployments(namespace).Get(deployment_csi_cephfsplugin_provisioner, metav1.GetOptions{})

				Expect(err).ToNot(HaveOccurred())
				Expect(out.Status.ReadyReplicas).ShouldNot(BeNumerically("==", 0))
				Expect(out.Status.UnavailableReplicas).Should(BeNumerically("==", 0))

				fmt.Printf("Found deployment %s in namespace %s\n", deployment_csi_cephfsplugin_provisioner, namespace)
			})
		})

		Context("check if each deployment installed", func() {
			It("deployment_csi_rbdplugin_provisioner", func() {
				out, err := clientset.AppsV1().Deployments(namespace).Get(deployment_csi_rbdplugin_provisioner, metav1.GetOptions{})

				Expect(err).ToNot(HaveOccurred())
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

				Expect(err).ToNot(HaveOccurred())
				Expect(out.Status.ReadyReplicas).ShouldNot(BeNumerically("==", 0))
				Expect(out.Status.UnavailableReplicas).Should(BeNumerically("==", 0))

				fmt.Printf("Found deployment %s in namespace %s\n", deployment_cdi_operator, namespace)
			})
		})
		Context("check if each deployment installed", func() {
			It("deployment_cdi_deployment", func() {
				out, err := clientset.AppsV1().Deployments(namespace).Get(deployment_cdi_deployment, metav1.GetOptions{})

				Expect(err).ToNot(HaveOccurred())
				Expect(out.Status.ReadyReplicas).ShouldNot(BeNumerically("==", 0))
				Expect(out.Status.UnavailableReplicas).Should(BeNumerically("==", 0))

				fmt.Printf("Found deployment %s in namespace %s\n", deployment_cdi_deployment, namespace)
			})
		})

		Context("check if each deployment installed", func() {
			It("deployment_cdi_apiserver", func() {
				out, err := clientset.AppsV1().Deployments(namespace).Get(deployment_cdi_apiserver, metav1.GetOptions{})

				Expect(err).ToNot(HaveOccurred())
				Expect(out.Status.ReadyReplicas).ShouldNot(BeNumerically("==", 0))
				Expect(out.Status.UnavailableReplicas).Should(BeNumerically("==", 0))

				fmt.Printf("Found deployment %s in namespace %s\n", deployment_cdi_apiserver, namespace)
			})
		})

		Context("check if each deployment installed", func() {
			It("deployment_cdi_uploadproxy", func() {
				out, err := clientset.AppsV1().Deployments(namespace).Get(deployment_cdi_uploadproxy, metav1.GetOptions{})

				Expect(err).ToNot(HaveOccurred())
				Expect(out.Status.ReadyReplicas).ShouldNot(BeNumerically("==", 0))
				Expect(out.Status.UnavailableReplicas).Should(BeNumerically("==", 0))

				fmt.Printf("Found deployment %s in namespace %s\n", deployment_cdi_uploadproxy, namespace)
			})
		})
	})

	Describe("[TEST][e2e][0001] Pod Networking", func() {
		var (
			namespace = "test-pod-networking"
			pod1Name  = "alpha"
			pod2Name  = "beta"
		)

		It("Check ping from one pod to another pod by ip address", func() {
			nsSpec := makeNamespaceSpec(namespace)

			err := createNamespace(clientset, nsSpec)
			Expect(err).ToNot(HaveOccurred())
			fmt.Printf("namespace %s is created\n", namespace)

			// create pod busybox named alpha
			pod1, err := createPod(clientset, pod1Name, namespace)
			Expect(err).ToNot(HaveOccurred())
			fmt.Printf("pod %s is created\n", pod1Name)

			// create pod busybox named beta
			pod2, err := createPod(clientset, pod2Name, namespace)
			Expect(err).ToNot(HaveOccurred())
			fmt.Printf("pod %s is created\n", pod2Name)

			err = waitTimeoutForPodStatus(clientset, pod1.Name, pod1.Namespace, corev1.PodRunning, time.Second*30)
			Expect(err).ToNot(HaveOccurred())
			err = waitTimeoutForPodStatus(clientset, pod2.Name, pod2.Namespace, corev1.PodRunning, time.Second*30)
			Expect(err).ToNot(HaveOccurred())

			// after pod created
			pod1Ip, err := getPodIp(clientset, pod1.Name, namespace)
			Expect(err).ToNot(HaveOccurred())
			pod2Ip, err := getPodIp(clientset, pod2.Name, namespace)
			Expect(err).ToNot(HaveOccurred())

			fmt.Printf("IP of pod_1 is %s\n", pod1Ip)
			fmt.Printf("IP of pod_2 is %s\n", pod2Ip)

			// check each ping test case
			var timeout = time.Second* 15
			var pollingInterval = time.Second * 3
			Eventually(func() bool {
				return canPingFromPodToIpAddr(pod1.Name, namespace, pod2Ip, clientset, config)
			}, timeout, pollingInterval).Should(BeTrue())

			Eventually(func() bool {
				return canPingFromPodToIpAddr(pod2.Name, namespace, pod1Ip, clientset, config)
			}, timeout, pollingInterval).Should(BeTrue())

			googleAddress := "google.com"
			Eventually(func() bool {
				return canPingFromPodToIpAddr(pod1.Name, namespace, googleAddress, clientset, config)
			}, timeout, pollingInterval).Should(BeTrue())

			Eventually(func() bool {
				return canPingFromPodToIpAddr(pod2.Name, namespace, googleAddress, clientset, config)
			}, timeout, pollingInterval).Should(BeTrue())

			//TODO Delete ns regardless of any above assertions !!!!!!!!!!!!!!
			// now delete only if all assertion pass
			timeout = time.Second* 300
			pollingInterval = time.Second * 10

			err = clientset.CoreV1().Namespaces().Delete(namespace, &metav1.DeleteOptions{})
			Expect(err).ToNot(HaveOccurred())
			Eventually(func() bool {
				ns, err := clientset.CoreV1().Namespaces().Get(namespace, metav1.GetOptions{})
				if err != nil || errors.IsNotFound(err) {
					return true
				}

				if ns.Status.Phase == corev1.NamespaceTerminating {
					fmt.Printf("Namespace %s is still in phase %s\n", namespace, ns.Status.Phase)
					return false
				}
				return false
			}, timeout, pollingInterval).Should(BeTrue())
		})
	})
})
