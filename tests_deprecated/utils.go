package tests

import (
	"flag"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

// TODO rook, cdi client 를 clientSet 과 더불어 사용, 재사용하기 위해서는 해당 코드에서 flag 등록하는 부분의 변경이 필요한 것으로 보임
func GetClientSet() (*kubernetes.Clientset, *restclient.Config) {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"),
			"(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse() // TODO flag 에 붙이는 방식이 ginkgo 사용에 문제가 되는 것 같음 ginkgo 커맨드의 flag 에 적용되는 듯 보임

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

	return clientset, config
}
