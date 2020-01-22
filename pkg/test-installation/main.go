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
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

/**
 * check whether installation was successful
 */
//TODO
// 0. go mod 대신 goDep 사용 고려
// 1. testframework (ginkgo같은것) 사용 및 test assertion 추가 필요 (현재 코드는 test 가 아니라 그냥 client 한 번 돌려보는 것에 불과)
// 2. test 결과 저장 및 print out 필요
// 3. 현재 k8s cluster 기본 세팅 관련 test 와 rook & cdi install 관련 test 가 모두 들어있는데, 따로 (폴더?파일?케이스?) 구분해야 함
// 3.1. 현재 make install 한 후에 make build-test 및 make test 를 진행하고 있는데, make install 은 사실 rook, cdi install 이고, k8s cluster 구성에 대한 test 는 이보다 앞서 진행해야 함
// 3.2. k8s cluster 구성에 대한 test 는 추후 e2e test 로 busybox 를 사용한 pod network 통신 케이스가 반드시 필요함
// 4. 정상 설치 아닐 경우 error message 만 뿜는 것이 아니라 다른 로직 추가 필요

func main() {
	fmt.Printf("\n\n [TEST] Check Rook, Cdi Deployment START \n\n")

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

	// get # nodes
	// TODO node status all-ready checks
	nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d nodes in the cluster\n", len(nodes.Items))

	// get # pods
	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	//TODO 각 deployment 들 list 로 만들고 for 문으로 돌리기
	//check desired deployment exists in rook-ceph namespace
	namespace_rook_ceph := "rook-ceph"
	deployment_rook_ceph_operator := "rook-ceph-operator"
	deployment_csi_cephfsplugin_provisioner := "csi-cephfsplugin-provisioner"
	deployment_csi_rbdplugin_provisioner := "csi-rbdplugin-provisioner"

	// check rook_ceph_operator deployment
	_, err = clientset.AppsV1().Deployments(namespace_rook_ceph).Get(deployment_rook_ceph_operator, metav1.GetOptions{})

	if errors.IsNotFound(err) {
		fmt.Printf("Deployment %s in namespace %s not found\n", deployment_rook_ceph_operator, namespace_rook_ceph)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting Deployment %s in namespace %s: %v\n",
			deployment_rook_ceph_operator, namespace_rook_ceph, statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found deployment_rook_ceph_operator %s in namespace %s\n", deployment_rook_ceph_operator, namespace_rook_ceph)
	}

	// check deployment_csi_cephfsplugin_provisioner deployment
	_, err = clientset.AppsV1().Deployments(namespace_rook_ceph).Get(deployment_csi_cephfsplugin_provisioner, metav1.GetOptions{})

	if errors.IsNotFound(err) {
		fmt.Printf("Deployment %s in namespace %s not found\n", deployment_csi_cephfsplugin_provisioner, namespace_rook_ceph)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting Deployment %s in namespace %s: %v\n",
			deployment_csi_cephfsplugin_provisioner, namespace_rook_ceph, statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found deployment_rook_ceph_operator %s in namespace %s\n", deployment_csi_cephfsplugin_provisioner, namespace_rook_ceph)
	}

	// check deployment_csi_rbdplugin_provisioner deployment
	_, err = clientset.AppsV1().Deployments(namespace_rook_ceph).Get(deployment_csi_rbdplugin_provisioner, metav1.GetOptions{})

	if errors.IsNotFound(err) {
		fmt.Printf("Deployment %s in namespace %s not found\n", deployment_csi_rbdplugin_provisioner, namespace_rook_ceph)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting Deployment %s in namespace %s: %v\n",
			deployment_csi_rbdplugin_provisioner, namespace_rook_ceph, statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found deployment_rook_ceph_operator %s in namespace %s\n", deployment_csi_rbdplugin_provisioner, namespace_rook_ceph)
	}

	// CDI
	// check desired deployment exists in cdi namespace
	namespace_cdi := "cdi"
	deployment_cdi_operator := "cdi-operator"
	deployment_cdi_deployment := "cdi-deployment"
	deployment_cdi_apiserver := "cdi-apiserver"
	deployment_cdi_uploadproxy := "cdi-uploadproxy"

	// check deployment_cdi_operator
	_, err = clientset.AppsV1().Deployments(namespace_cdi).Get(deployment_cdi_operator, metav1.GetOptions{})

	if errors.IsNotFound(err) {
		showError(err)
		fmt.Printf("Deployment %s in namespace %s not found\n", deployment_cdi_operator, namespace_cdi)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting Deployment %s in namespace %s: %v\n",
			deployment_cdi_operator, namespace_cdi, statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found deployment_cdi_operator %s in namespace %s\n", deployment_cdi_operator, namespace_cdi)
	}

	// check deployment_cdi_deployment
	_, err = clientset.AppsV1().Deployments(namespace_cdi).Get(deployment_cdi_deployment, metav1.GetOptions{})

	if errors.IsNotFound(err) {
		fmt.Printf("Deployment %s in namespace %s not found\n", deployment_cdi_deployment, namespace_cdi)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting Deployment %s in namespace %s: %v\n",
			deployment_cdi_deployment, namespace_cdi, statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found deployment_cdi_operator %s in namespace %s\n", deployment_cdi_deployment, namespace_cdi)
	}

	// check deployment_cdi_apiserver
	_, err = clientset.AppsV1().Deployments(namespace_cdi).Get(deployment_cdi_apiserver, metav1.GetOptions{})

	if errors.IsNotFound(err) {
		fmt.Printf("Deployment %s in namespace %s not found\n", deployment_cdi_apiserver, namespace_cdi)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting Deployment %s in namespace %s: %v\n",
			deployment_cdi_apiserver, namespace_cdi, statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found deployment_cdi_operator %s in namespace %s\n", deployment_cdi_apiserver, namespace_cdi)
	}

	// check deployment_cdi_apiserver
	_, err = clientset.AppsV1().Deployments(namespace_cdi).Get(deployment_cdi_uploadproxy, metav1.GetOptions{})

	if errors.IsNotFound(err) {
		fmt.Printf("Deployment %s in namespace %s not found\n", deployment_cdi_uploadproxy, namespace_cdi)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting Deployment %s in namespace %s: %v\n",
			deployment_cdi_uploadproxy, namespace_cdi, statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found deployment_cdi_operator %s in namespace %s\n", deployment_cdi_uploadproxy, namespace_cdi)
	}

	fmt.Printf("\n\n [TEST] Check Rook, Cdi Deployment SUCCESS \n\n")
}

//if err has been occured Warning printout
func showError(err error) {
	fmt.Printf("\n\n [TEST] Check Rook, Cdi Deployment FAILED \n\n")
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
