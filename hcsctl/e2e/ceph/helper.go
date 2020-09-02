package testCeph

import (
	"flag"
	"fmt"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"

	"github.com/rook/rook/pkg/util/exec"
	"k8s.io/client-go/dynamic"
)

type HyperHelper struct {
	executor  *exec.CommandExecutor
	Client    *dynamic.Interface
	RunningInCluster bool
}

var (
	uniqueHyperSdsHelper *HyperHelper
	uniqueConfig             *restclient.Config
)

func CreateK8sHelper() (*HyperHelper, error) {
	executor := &exec.CommandExecutor{}

	// TODO home 에서 가져오지 말고 KUBECONFIG 환경변수로부터 가져오기 ("" 일 경우 하드코딩으로)
	var kubeconfig *string
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	if home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"),
			"(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	//flag.Parse() // TODO flag 에 붙이는 방식때문에 **반드시** 단 한 번 이 함수가 불려야 함

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to get clientset. %+v", err)
	}

	h := &HyperHelper{executor: executor, Client: &client}

	//TODO cluster 밖에서 명령 보내는 경우 고려 (현재는 InCluster인 경우만 하고 있음)

	//if strings.Index(config.Host, "//10.") != -1 {
	//	h.RunningInCluster = true
	//}
	h.RunningInCluster = true
	uniqueHyperSdsHelper, uniqueConfig = h, config
	return h, err
}

func HyperSdsHelper() *HyperHelper {
	return uniqueHyperSdsHelper
}

func HyperSdsConfig() *restclient.Config {
	return uniqueConfig
}
