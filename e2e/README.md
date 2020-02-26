# e2e 테스트 수행
e2e를 위해 멀티노드 기반의 Kubernetes 환경을 생성하고 go test를 위한 준비를 합니다. HyperCloud Storage 설치 관련해서는 해당 문서를 참조하세요.

## Prerequisites
1. 필요한 패키지
- gcc, make
- vagrant, virtualbox
- kubectl(>v1.16.0), minikube(>v1.7.2), kustomize(v3.5.4)
- go(>v1.13)

2. go 환경 설정
```shell
go get -u github.com/onsi/ginkgo/ginkgo
go get -u github.com/onsi/gomega/...
export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin
export GO111MODULE=auto
```
## VM 기반의 Kuberentes 환경 구축
```shell
# VM 기반의 멀티노드 Kubernetes 환경 생성
# BOX_OS, NODE_COUNT, KUBERNETES_VERSION를 설정할 수 있습니다.
BOX_OS=centos NODE_COUNT=2 KUBERNETES_VERSION=1.15.3 ./e2e/cluster.sh clusterUp

# 제거
./e2e/cluster.sh clusterClean
```
## e2e 테스트
생성된 HyperCloud Storage가 정상적으로 설치되었는지 확인합니다.

```shell
# default storageclass를 지정합니다
kubectl patch storageclass {$원하는 StorageClassName} -p '{"metadata"':' {"annotations":{"storageclass.kubernetes.io/is-default-class"':'"true"}}}

# 테스트를 시작합니다
./installer.sh test
```
