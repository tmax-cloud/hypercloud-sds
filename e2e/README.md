# e2e test

## Prerequisites

- gcc
- make
- vagrant
- vboxmanage 
- kubectl v1.16
- minikube ( > v1.7.2)
- go setting
    - 패키지 설치 : go 1.13 & ginkgo 1.12.0 & gomega
    - PATH 설정 : export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
      - check `go version`, `ginkgo version`
    - GO111MODULE=auto 설정
