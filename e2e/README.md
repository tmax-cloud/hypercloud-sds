# e2e test
> Kubernetes & HypercloudStorage 환경에 Hypercloud Storage 가 정상 설치되었는지 확인합니다.
>
> 이미 설치된 환경을 테스트하는 경우는 [###On Existing Kubernetes-Cluster] 를 따라하시고, 환경 설치부터 시작하는 경우는 [###On Non-Existing Kubernetes-Cluster] 를 따라하시기 바랍니다.

## Prerequisites

### For Creating Kubernetes-VM-Cluster
1. 패키지 설치:
  - gcc
  - make
  - vagrant
  - vboxmanage 
  - kubectl (v1.16)
  - minikube ( > v1.7.2)

### For Testing
1. 패키지 설치:
  - go 1.13 [official guide](https://golang.org/doc/install)
  - ginkgo 1.12.0
    - `go get -u github.com/onsi/ginkgo/ginkgo`
  - gomega
    - `go get -u github.com/onsi/gomega/...`
2. GO 환경 설정
  - PATH 설정 : `export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin`
  - GO mod 설정 : `export GO111MODULE=auto`
3. 다음 커맨드로 정상 설치를 확인합니다 : `go version`, `ginkgo version`
  - Error 가 출력되지 않고 설치한 버전이 정상적으로 출력되는지 확인합니다.

## Quick Start
> Hypercloud Storage 가 정상적으로 작동하는지 확인합니다.

### On Existing Kubernetes-Cluster
1. default storageclass 를 지정합니다 : `kubectl patch storageclass {$원하는 StorageClassName} -p '{"metadata"':' {"annotations":{"storageclass.kubernetes.io/is-default-class"':'"true"}}}'`
2. 테스트를 시작합니다 : `./../install/installer.sh test`

### On Non-Existing Kubernetes-Cluster
1. 테스트용 k8s-vm-cluster 환경을 설치합니다 : `BOX_OS=centos NODE_COUNT=2 KUBERNETES_VERSION=1.15.3 ./cluster.sh clusterUp`
  - 다음 환경 변수는 원하는 값으로 설정 가능합니다 :
    - `BOX_OS`
    - `NODE_COUNT`
    - `KUBERNETES_VERSION`
2. Hypercloud Storage 를 해당 cluster에 설치합니다 : `./../install/installer.sh install`
3. 테스트를 시작합니다 : `./../install/installer.sh test`
4. Hypercloud Storage 를 해당 cluster에서 제거합니다 : `./../install/installer.sh uninstall`
5. 테스트용 k8s-vm-cluster 환경을 제거합니다 : `./cluster.sh clusterClean`