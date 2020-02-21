package tests

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Test Pod Network", func() {
	Describe("[TEST][e2e][0001] Pod Networking", func() {
		const (
			pod1Name = "alpha"
			pod2Name = "beta"
		)
		BeforeEach(func() {
			// create testing namespace
			testingNamespace, err = createNamespace(hyperStorageHelper.Clientset,
				makeNamespaceSpec(Pod2PodTestingNamespacePrefix))
			Expect(err).ToNot(HaveOccurred())
			fmt.Printf("Namespace %s is created for testing.\n", testingNamespace.Name) // TODO fmt 대신 log 사용
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
		})
		It("[TEST - 10] Check ping from one pod to another pod by ip address", func() {
			// create pod busybox named alpha
			pod1, err := createPod(hyperStorageHelper.Clientset, pod1Name, testingNamespace.Name)
			Expect(err).ToNot(HaveOccurred())
			fmt.Printf("pod %s is created\n", pod1Name)

			// create pod busybox named beta
			pod2, err := createPod(hyperStorageHelper.Clientset, pod2Name, testingNamespace.Name)
			Expect(err).ToNot(HaveOccurred())
			fmt.Printf("pod %s is created\n", pod2Name)

			err = waitTimeoutForPodStatus(hyperStorageHelper.Clientset, pod1.Name, pod1.Namespace,
				corev1.PodRunning, TimeoutForCreatingPod)
			Expect(err).ToNot(HaveOccurred())
			err = waitTimeoutForPodStatus(hyperStorageHelper.Clientset, pod2.Name, pod2.Namespace,
				corev1.PodRunning, TimeoutForCreatingPod)
			Expect(err).ToNot(HaveOccurred())

			// after pod created
			pod1Ip, err := getPodIP(hyperStorageHelper.Clientset, pod1.Name, testingNamespace.Name)
			Expect(err).ToNot(HaveOccurred())
			pod2Ip, err := getPodIP(hyperStorageHelper.Clientset, pod2.Name, testingNamespace.Name)
			Expect(err).ToNot(HaveOccurred())

			fmt.Printf("IP of %s is %s\n", pod1Name, pod1Ip)
			fmt.Printf("IP of %s is %s\n", pod2Name, pod2Ip)

			// check each ping test case
			Eventually(func() bool {
				return canPingFromPodToIPAddr(pod1.Name, testingNamespace.Name, pod2Ip,
					hyperStorageHelper.Clientset, HyperStorageConfig())
			}, TimeoutForPing, PollingIntervalForPing).Should(BeTrue())

			Eventually(func() bool {
				return canPingFromPodToIPAddr(pod2.Name, testingNamespace.Name, pod1Ip,
					hyperStorageHelper.Clientset, HyperStorageConfig())
			}, TimeoutForPing, PollingIntervalForPing).Should(BeTrue())

			Eventually(func() bool {
				return canPingFromPodToIPAddr(pod1.Name, testingNamespace.Name, GoogleIPAddress,
					hyperStorageHelper.Clientset, HyperStorageConfig())
			}, TimeoutForPing, PollingIntervalForPing).Should(BeTrue())

			Eventually(func() bool {
				return canPingFromPodToIPAddr(pod2.Name, testingNamespace.Name, GoogleIPAddress,
					hyperStorageHelper.Clientset, HyperStorageConfig())
			}, TimeoutForPing, PollingIntervalForPing).Should(BeTrue())
		})
	})
})
