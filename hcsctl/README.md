# hcsctl: hypercloud storage ctl
hcsctl은 hypercloud storage의 설치, 제거 및 관리를 제공합니다.

# Install
## Prerequisite
- kubectl (> 1.15.0)
- Existing Kubernetes Cluster

## 설치
TODO: hcsctl 바이너리 업로드해두고 다운로드 링크 여기에 걸기

- 바이너리 생성 가이드
  - prerequisite
    - go 관련 환경 설정
      - `export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin` 
      - `export GO111MODULE=auto`
    - go 관련 바이너리 다운로드
      - `go get -u github.com/onsi/ginkgo/ginkgo`
      - `go get github.com/markbates/pkger/cmd/pkger`
    - 패키지 설치
      - gcc
      - make
  - `cd hypercloud-storage/hcsctl`
  - `make build`
  - 이제 hypercloud-storage/hcsctl/build 디렉토리에 생성된 `hcsctl` 과 `hcsctl.test` 바이너리를 사용하실 수 있습니다.

## 지원 기능 목록
- 자세한 사항은 `hcsctl help` 를 참고하세요.
- create-inventory
  - ex) hcsctl create-inventory myInventory
  - `hcsctl install` 을 위한 정해진 형식의 yaml 파일을 담은 디렉토리 `./myInventory` 를 생성합니다.
  - `./myInventory/rook/*.yaml` 은 rook-ceph 설치에 사용되는 yaml 파일이며 `./myInventory/cdi/*.yaml` 은 cdi 설치에 사용되는 yaml 파일입니다.
  - 생성된 모든 yaml 파일들은 sample 제공용 파일이므로, 파일명을 제외한 파일 내용은 사용자가 원하는대로 수정 후 사용하시면 됩니다.  
- install
  - ex) hcsctl install myInventory
- uninstall
  - ex) hcsctl uninstall myInventory
- ceph {status/exec}
  - ex) hcsctl ceph status
  - ex) hcsctl ceph exec ceph osd status
  - ex) hcsctl ceph exec ceph df
  
# Quick Start
TODO: 자세하게 작성

```shell
hcsctl install my_inventory
hcsctl uninstall my_inventory
# e2e 테스트
hcsctl.test
```