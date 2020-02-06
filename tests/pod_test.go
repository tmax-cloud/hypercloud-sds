package tests

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

var _ = Describe("Test Pod Network", func() {
	Describe("[TEST][e2e][0001] Pod Networking", func() {
		const (
			namespacePrefix                     = "test-pod-networking-"
			pod1Name                            = "alpha"
			pod2Name                            = "beta"
			timeoutForCreatingPod               = time.Second * 60
			timeoutForPing                      = time.Second * 20
			pollingIntervalForPing              = time.Second * 5
			timeoutForDeletingNamespace         = time.Second * 300
			pollingIntervalForDeletingNamespace = time.Second * 5
		)

		It("[TEST - 10] Check ping from one pod to another pod by ip address", func() {
			nsSpec := makeNamespaceSpec(namespacePrefix)

			generatedNs, err := createNamespace(hyperStorageHelper.Clientset, nsSpec)
			Expect(err).ToNot(HaveOccurred())
			fmt.Printf("namespace %s is created\n", generatedNs.Name)

			// create pod busybox named alpha
			pod1, err := createPod(hyperStorageHelper.Clientset, pod1Name, generatedNs.Name)
			Expect(err).ToNot(HaveOccurred())
			fmt.Printf("pod %s is created\n", pod1Name)

			// create pod busybox named beta
			pod2, err := createPod(hyperStorageHelper.Clientset, pod2Name, generatedNs.Name)
			Expect(err).ToNot(HaveOccurred())
			fmt.Printf("pod %s is created\n", pod2Name)

			err = waitTimeoutForPodStatus(hyperStorageHelper.Clientset, pod1.Name, pod1.Namespace,
				corev1.PodRunning, timeoutForCreatingPod)
			Expect(err).ToNot(HaveOccurred())
			err = waitTimeoutForPodStatus(hyperStorageHelper.Clientset, pod2.Name, pod2.Namespace,
				corev1.PodRunning, timeoutForCreatingPod)
			Expect(err).ToNot(HaveOccurred())

			// after pod created
			pod1Ip, err := getPodIP(hyperStorageHelper.Clientset, pod1.Name, generatedNs.Name)
			Expect(err).ToNot(HaveOccurred())
			pod2Ip, err := getPodIP(hyperStorageHelper.Clientset, pod2.Name, generatedNs.Name)
			Expect(err).ToNot(HaveOccurred())

			fmt.Printf("IP of %s is %s\n", pod1Name, pod1Ip)
			fmt.Printf("IP of %s is %s\n", pod2Name, pod2Ip)

			// check each ping test case
			Eventually(func() bool {
				return canPingFromPodToIPAddr(pod1.Name, generatedNs.Name, pod2Ip,
					hyperStorageHelper.Clientset, HyperStorageConfig())
			}, timeoutForPing, pollingIntervalForPing).Should(BeTrue())

			Eventually(func() bool {
				return canPingFromPodToIPAddr(pod2.Name, generatedNs.Name, pod1Ip,
					hyperStorageHelper.Clientset, HyperStorageConfig())
			}, timeoutForPing, pollingIntervalForPing).Should(BeTrue())

			googleAddress := "google.com"
			Eventually(func() bool {
				return canPingFromPodToIPAddr(pod1.Name, generatedNs.Name, googleAddress,
					hyperStorageHelper.Clientset, HyperStorageConfig())
			}, timeoutForPing, pollingIntervalForPing).Should(BeTrue())

			Eventually(func() bool {
				return canPingFromPodToIPAddr(pod2.Name, generatedNs.Name, googleAddress,
					hyperStorageHelper.Clientset, HyperStorageConfig())
			}, timeoutForPing, pollingIntervalForPing).Should(BeTrue())

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
		})
	})
})
