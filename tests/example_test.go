package tests

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
)

var _ = Describe("TEST", func() {
	var (
		clientset *kubernetes.Clientset
		config    *restclient.Config
	)

	BeforeSuite(func() { //TODO BeforeEach 로 clientSet 사용 시,
		// flag.go 의 Happens only if flags are declared with identical names 에러 발생
		clientset, config = getClientSet()
	})

	// TODO Temporary codes for nodes & pods check
	Describe("nodes & pods check", func() {
		Context("Number of Nodes", func() {
			It("[TEST - 01] should not be zero", func() {
				nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
				Expect(err).ToNot(HaveOccurred())
				Expect(len(nodes.Items)).NotTo(Equal(0)) // TODO should change to another assertion
			})
		})

		Context("Number of Pods", func() {
			It("[TEST - 02] should not be zero ", func() {
				pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
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
				out, err := clientset.AppsV1().Deployments(namespace).Get(deploymentRookCephOperator, metav1.GetOptions{})

				Expect(err).ToNot(HaveOccurred())
				Expect(out.Status.ReadyReplicas).ShouldNot(BeNumerically("==", 0))
				Expect(out.Status.UnavailableReplicas).Should(BeNumerically("==", 0))

				//TODO fmt -> log 사용 및 관리 필요
				fmt.Printf("Found deployment %s in namespace %s\n", deploymentRookCephOperator, namespace)
			})
		})
		Context("check if each deployment installed", func() {
			It("[TEST - 04] deployment_csi_cephfsplugin_provisioner", func() {
				out, err := clientset.AppsV1().Deployments(namespace).Get(deploymentCsiCephfspluginProvisioner, metav1.GetOptions{})

				Expect(err).ToNot(HaveOccurred())
				Expect(out.Status.ReadyReplicas).ShouldNot(BeNumerically("==", 0))
				Expect(out.Status.UnavailableReplicas).Should(BeNumerically("==", 0))

				fmt.Printf("Found deployment %s in namespace %s\n", deploymentCsiCephfspluginProvisioner, namespace)
			})
		})

		Context("check if each deployment installed", func() {
			It("[TEST - 05] deployment_csi_rbdplugin_provisioner", func() {
				out, err := clientset.AppsV1().Deployments(namespace).Get(deploymentCsiRbdpluginProvisioner, metav1.GetOptions{})

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
				out, err := clientset.AppsV1().Deployments(namespace).Get(deploymentCdiOperator, metav1.GetOptions{})

				Expect(err).ToNot(HaveOccurred())
				Expect(out.Status.ReadyReplicas).ShouldNot(BeNumerically("==", 0))
				Expect(out.Status.UnavailableReplicas).Should(BeNumerically("==", 0))

				fmt.Printf("Found deployment %s in namespace %s\n", deploymentCdiOperator, namespace)
			})
		})
		Context("check if each deployment installed", func() {
			It("[TEST - 07] deployment_cdi_deployment", func() {
				out, err := clientset.AppsV1().Deployments(namespace).Get(deploymentCdiDeployment, metav1.GetOptions{})

				Expect(err).ToNot(HaveOccurred())
				Expect(out.Status.ReadyReplicas).ShouldNot(BeNumerically("==", 0))
				Expect(out.Status.UnavailableReplicas).Should(BeNumerically("==", 0))

				fmt.Printf("Found deployment %s in namespace %s\n", deploymentCdiDeployment, namespace)
			})
		})

		Context("check if each deployment installed", func() {
			It("[TEST - 08] deployment_cdi_apiserver", func() {
				out, err := clientset.AppsV1().Deployments(namespace).Get(deploymentCdiApiserver, metav1.GetOptions{})

				Expect(err).ToNot(HaveOccurred())
				Expect(out.Status.ReadyReplicas).ShouldNot(BeNumerically("==", 0))
				Expect(out.Status.UnavailableReplicas).Should(BeNumerically("==", 0))

				fmt.Printf("Found deployment %s in namespace %s\n", deploymentCdiApiserver, namespace)
			})
		})

		Context("check if each deployment installed", func() {
			It("[TEST - 09] deployment_cdi_uploadproxy", func() {
				out, err := clientset.AppsV1().Deployments(namespace).Get(deploymentCdiUploadproxy, metav1.GetOptions{})

				Expect(err).ToNot(HaveOccurred())
				Expect(out.Status.ReadyReplicas).ShouldNot(BeNumerically("==", 0))
				Expect(out.Status.UnavailableReplicas).Should(BeNumerically("==", 0))

				fmt.Printf("Found deployment %s in namespace %s\n", deploymentCdiUploadproxy, namespace)
			})
		})
	})

	Describe("[TEST][e2e][0001] Pod Networking", func() {
		const (
			namespace                           = "test-pod-networking"
			pod1Name                            = "alpha"
			pod2Name                            = "beta"
			timeoutForCreatingPod               = time.Second * 30
			timeoutForPing                      = time.Second * 15
			pollingIntervalForPing              = time.Second * 3
			timeoutForDeletingNamespace         = time.Second * 300
			pollingIntervalForDeletingNamespace = time.Second * 10
		)

		It("[TEST - 10] Check ping from one pod to another pod by ip address", func() {
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

			err = waitTimeoutForPodStatus(clientset, pod1.Name, pod1.Namespace, corev1.PodRunning, timeoutForCreatingPod)
			Expect(err).ToNot(HaveOccurred())
			err = waitTimeoutForPodStatus(clientset, pod2.Name, pod2.Namespace, corev1.PodRunning, timeoutForCreatingPod)
			Expect(err).ToNot(HaveOccurred())

			// after pod created
			pod1Ip, err := getPodIP(clientset, pod1.Name, namespace)
			Expect(err).ToNot(HaveOccurred())
			pod2Ip, err := getPodIP(clientset, pod2.Name, namespace)
			Expect(err).ToNot(HaveOccurred())

			fmt.Printf("IP of %s is %s\n", pod1Name, pod1Ip)
			fmt.Printf("IP of %s is %s\n", pod2Name, pod2Ip)

			// check each ping test case
			Eventually(func() bool {
				return canPingFromPodToIPAddr(pod1.Name, namespace, pod2Ip, clientset, config)
			}, timeoutForPing, pollingIntervalForPing).Should(BeTrue())

			Eventually(func() bool {
				return canPingFromPodToIPAddr(pod2.Name, namespace, pod1Ip, clientset, config)
			}, timeoutForPing, pollingIntervalForPing).Should(BeTrue())

			googleAddress := "google.com"
			Eventually(func() bool {
				return canPingFromPodToIPAddr(pod1.Name, namespace, googleAddress, clientset, config)
			}, timeoutForPing, pollingIntervalForPing).Should(BeTrue())

			Eventually(func() bool {
				return canPingFromPodToIPAddr(pod2.Name, namespace, googleAddress, clientset, config)
			}, timeoutForPing, pollingIntervalForPing).Should(BeTrue())

			//TODO MUST DELETE NAMESPACE regardless of any above assertions !!!!!!!!!!!!!!
			// 현재는 위의 모든 assertion 이 true 일 때만, 아래 라인을 타고 namespace 를 delete 하고 있음
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
			}, timeoutForDeletingNamespace, pollingIntervalForDeletingNamespace).Should(BeTrue())
		})
	})
})
