/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Note: the example only works with the code within the same release/branch.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/remotecommand"
	"os"
	"path/filepath"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

/**
 * 1. create 2 pod with image busybox
 * 2. check ping is possible with each other
 */

func main() {
	//TODO printf 대신 logger(klog?) 사용
	fmt.Printf("\n\n [TEST] Check Pod 2 Pod Networking START \n\n")

	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// get # pods
	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	//TODO make 2nd pod and ping eachother
	namespace := "test-pod-networking"
	pod_1_name := "alpha"
	pod_2_name := "beta"

	// create namespace
	nsSpec := makeNamespace(namespace)
	_, err = clientset.CoreV1().Namespaces().Create(nsSpec)
	if err != nil {
		showError(err)
		panic(err.Error())
	}
	fmt.Printf("namespace %s is created\n", namespace)

	// create pod busybox named alpha
	err = createPod(clientset, pod_1_name, namespace)
	if err != nil {
		showError(err)
		panic(err.Error())
	}

	// create pod busybox named beta
	err = createPod(clientset, pod_2_name, namespace)
	if err != nil {
		showError(err)
		panic(err.Error())
	}

	//TODO status-running check
	// wait some reasonable sec until phase become running status.
	// while문 같은 for문 사용 필요
	err = waitTimeoutForPodStatus(clientset, pod_1_name, namespace, corev1.PodRunning, time.Second*10)
	if err != nil {
		showError(err)
		panic(err.Error())
	}

	err = waitTimeoutForPodStatus(clientset, pod_2_name, namespace, corev1.PodRunning, time.Second*10)
	if err != nil {
		showError(err)
		panic(err.Error())
	}

	// After Pod Created
	pod_1_ip, err := getPodIp(clientset, pod_1_name, namespace)
	pod_2_ip, err := getPodIp(clientset, pod_2_name, namespace)

	fmt.Printf("IP of pod_1 is %s\n", pod_1_ip)
	fmt.Printf("IP of pod_2 is %s\n", pod_2_ip)

	var result bool
	result = canPingFromPodToIpAddr(pod_1_name, namespace, pod_2_ip, clientset, config)
	if result {
		fmt.Printf("ping from pod %s to pod %s is available\n", pod_1_name, pod_2_name)
	} else {
		fmt.Printf("ping from pod %s to pod %s is Unavailable\n", pod_1_name, pod_2_name)
	}

	//TODO bug ? 간헐적으로 beta -> google.com failed
	time.Sleep(time.Second)

	result = canPingFromPodToIpAddr(pod_2_name, namespace, pod_1_ip, clientset, config)
	if result {
		fmt.Printf("ping from pod %s to pod %s is available\n", pod_2_name, pod_1_name)
	} else {
		fmt.Printf("ping from pod %s to pod %s is Unavailable\n", pod_2_name, pod_1_name)
	}

	time.Sleep(time.Second)

	googleAddress := "google.com"
	result = canPingFromPodToIpAddr(pod_1_name, namespace, googleAddress, clientset, config)
	if result {
		fmt.Printf("ping from pod %s to %s is available\n", pod_1_name, googleAddress)
	} else {
		fmt.Printf("ping from pod %s to %s is Unavailable\n", pod_1_name, googleAddress)
	}

	time.Sleep(time.Second)

	result = canPingFromPodToIpAddr(pod_2_name, namespace, googleAddress, clientset, config)
	if result {
		fmt.Printf("ping from pod %s to %s is available\n", pod_2_name, googleAddress)
	} else {
		fmt.Printf("ping from pod %s to %s is Unavailable\n", pod_2_name, googleAddress)
	}

	//TODO go 버전의 try~finally 같은 걸로 delete 모든 리소스 after test finished (indep.of Success/Failure)
	// 그럼 위의 err 를 모두 저장해야 함? 컴포넌트화 ?

	fmt.Printf("Deleting namespace %s ...\n", namespace)
	//TODO 현재 delete 는 커맨드 날리기만하고 terminating 상태인데도 불구하고 종료함
	// DeletePropagationForeground 는 효과 x
	// watch 를 하는 방법 ?
	err = clientset.CoreV1().Namespaces().Delete(namespace, &metav1.DeleteOptions{})
	if err != nil {
		showError(err)
		panic(err.Error())
	}
	fmt.Printf("\n\n [TEST] Check Pod 2 Pod Networking SUCCESS \n\n")
}

func makeNamespace(namespace string) *corev1.Namespace {
	ns := &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}

	return ns
}

func createPod(clientset *kubernetes.Clientset, podName string, namespace string) error {
	pod1 := makePodSpec(podName, namespace)
	_, err := clientset.CoreV1().Pods(namespace).Create(pod1)

	if err != nil {
		showError(err)
		fmt.Printf("pod %s creating failed in namespace %s \n", podName, namespace)

		return err
	}

	return nil
}

func makePodSpec(podName string, namespace string) *corev1.Pod {
	//TODO need to be clean
	cmd := []string{"sleep", "3600"}

	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Image:           "busybox",
					Name:            "busybox",
					Command:         cmd,
					ImagePullPolicy: corev1.PullIfNotPresent,
				},
			},
			RestartPolicy: corev1.RestartPolicyAlways,
		},
	}

	return pod
}

//TODO 리팩토링 필요
// watch 명령어 확인
func waitTimeoutForPodStatus(clientset *kubernetes.Clientset, podName string, namespace string,
	desiredStatus corev1.PodPhase, timeout time.Duration) error {

	out, err := clientset.CoreV1().Pods(namespace).Get(podName, metav1.GetOptions{})

	if errors.IsNotFound(err) {
		showError(err)
		fmt.Printf("pod %s in namespace %s not found\n", podName, namespace)

		return err
	}

	if out.Status.Phase != desiredStatus {
		fmt.Printf("status of pod %s is now %s\n", podName, out.Status.Phase)
		fmt.Printf("wait for %d %s\n", timeout/time.Second, "second")
		time.Sleep(timeout)
		fmt.Printf("wait finished %d %s\n", timeout/time.Second, "second")

		out, err = clientset.CoreV1().Pods(namespace).Get(podName, metav1.GetOptions{})
	}

	if out.Status.Phase != desiredStatus {
		fmt.Printf("status of pod %s is still %s after waiting %d %s\n", podName, out.Status.Phase, timeout/time.Second, "second")
		fmt.Printf("something wrong\n")

		return errors.FromObject(out) //TODO 변경 필요
	}

	fmt.Printf("status of pod %s is now desired status %s\n", podName, desiredStatus)
	return nil
}

func getPodIp(clientset *kubernetes.Clientset, podName string, namespace string) (string, error) {
	out, err := clientset.CoreV1().Pods(namespace).Get(podName, metav1.GetOptions{})

	if errors.IsNotFound(err) {
		showError(err)
		fmt.Printf("pod %s in namespace %s not found\n", podName, namespace)

		return "", err
	}

	return out.Status.PodIP, nil
}

/////////////////////
// TODO pod2pod network test 더 간단하게 하는 방법
//kubectl exec -n test-pod-networking alpha -- ping -c 4 172.17.0.5 | grep "0% packet loss" >/dev/null; echo $?
//actualCmd := "ping -c 4 " + pod_1_ip + " | grep \"0% packet loss\" > /dev/null; echo $?"
//pingCmd := []string{actualCmd}
//result, err := framework.LookForStringInPodExec(namespace, pod_1_name, pingCmd, "0", 10)
//result, err := framework.LookForStringInPodExec(namespace, pod_1_name, []string{"/bin/ping", "-c", "2", pod_2_ip}, "1", 10)

// 아래 코드는 a4abhishek / Client-Go-Examples 의 github 참고
func canPingFromPodToIpAddr(podName string, namespace string, destinationIpAddress string, clientset *kubernetes.Clientset,
	config *restclient.Config) bool {
	//TODO 커맨드에 ping 명령어 이후 파이프라인(|)이랑 "> /dev/null" 먹지 않아서 조잡하게 코드 짰는데 확인 필요
	command := []string{"/bin/ping", "-c", "2", destinationIpAddress}

	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec")

	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		panic(err)
	}

	parameterCodec := runtime.NewParameterCodec(scheme)
	req.VersionedParams(&corev1.PodExecOptions{
		Command:   command,
		Container: "",
		Stdin:     false,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
	}, parameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		panic(err)
	}

	var stdout, stderr bytes.Buffer
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	})

	if err != nil {
		fmt.Printf("stdout is %s\n", stdout)
		fmt.Printf("stderr is %s\n", stderr)
		showError(err)
		panic(err)
	}

	if !strings.Contains(stdout.String(), "0% packet loss") {
		return false
	}

	return true
}

func showError(err error) {
	fmt.Printf("\n\n [TEST] Check Pod 2 Pod Networking FAILED \n\n")
	fmt.Println("=== WARNING ===\n")
	fmt.Printf("=== ERROR IS === %s ===\n", err)
	fmt.Println("=== WARNING ===\n")
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
