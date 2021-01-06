# hcsctl: hypercloud storage ctl
hcsctl은 hypercloud storage의 설치, 제거 및 관리 기능을 제공합니다.

## Prerequisite
- kubectl (>= 1.15)
- Kubernetes Cluster
- lvm2 Package (OSD Host)

## 바이너리 다운로드
- release 버전의 hcsctl 바이너리는 다음 경로에서 다운받을 수 있습니다.
    - [hypercloud-sds-release](https://github.com/tmax-cloud/hypercloud-sds/releases)
- master branch 버전의 hcsctl 바이너리를 사용하고 싶으면 다음 가이드를 수행하여 로컬에서 빌드할 수 있습니다.

  ``` shell
  # 바이너리 생성 시 gcc, make, git 패키지가 필요합니다.

  # hypercloud-sds 프로젝트 clone
  $ git clone https://github.com/tmax-cloud/hypercloud-sds.git

  # go 환경 설정
  $ export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin
  $ export GO111MODULE=auto

  # go 관련 바이너리 다운로드
  $ go get -u github.com/onsi/ginkgo/ginkgo
  $ go get github.com/markbates/pkger/cmd/pkger

  # 빌드 방법
  $ cd hypercloud-sds/hcsctl
  $ make build

  # 이제 hypercloud-sds/hcsctl/build 디렉토리에 hcsctl과 hcsctl.test 바이너리가 생성된 것을 확인하실 수 있습니다.
  ```

## 시작하기
- hcsctl 로 hypercloud-sds 설치에 앞서 설치에 필요한 yaml 파일을 생성하고 환경에 맞게 변경합니다.

   ``` shell
   $ hcsctl create-inventory {$inventory_name} [--include-cdi=BOOLEAN]
   # 예) hcsctl create-inventory myInventory --include-cdi=true
   ```

    - hcsctl 은 고가용성 스토리지인 rook-ceph 을 기본적으로 제공하며, 추가로 이미지 임포팅을 위한 cdi 를 옵션으로 제공합니다.
      - inventory 를 create 할 때 `--include-cdi` 플래그를 true 혹은 1 로 설정함으로써 cdi 를 inventory 에 포함시킬 수 있습니다. 
      - `--include-cdi` 플래그를 입력하지 않을 경우 기본값이 false 로 설정되어, inventory 에는 rook-ceph 만 생성됩니다.
    - 생성된 inventory 에는 rook, cdi 두 개의 디렉토리가 생성될 수 있습니다.
      - `./myInventory/rook/*.yaml` 에는 rook-ceph 설치에 사용되는 yaml 파일이 생성 됩니다.
      - `./myInventory/cdi/*.yaml` 에는 cdi 설치에 사용되는 yaml 파일이 생성됩니다. (`--include-cdi` 값을 true 로 설정한 경우)
    - 생성된 모든 yaml 파일들은 sample 제공용 파일이므로, 각 yaml 파일 내용을 반드시 **사용자의 환경에 맞게 수정** 후 사용하셔야 합니다. 
      - **생성된 폴더와 파일명은 절대 수정하시면 안됩니다.**
    - `./myInventory/rook/` 경로 밑에 생성된 yaml 파일을 환경에 맞게 수정하는 가이드 입니다.
        - cluster.yaml: [Rook-Ceph 클러스터 구성 가이드](./../docs/ceph-cluster-setting.md)
        - block_pool.yaml, block_sc.yaml: [Block Storage 가이드](./../docs/block.md)
        - file_system.yaml, file_sc.yaml: [Shared Filesystem 가이드](./../docs/file.md)
    - `./myInventory/cdi/` 경로 밑에 생성된 yaml 파일의 경우, 설치할 cdi version 변경이 필요한 경우에만 `operator.yaml` 파일 내의 `OPERATOR_VERSION`과 container image 버전을 변경하시면 됩니다.

- 테스트 환경용 inventory 설정 방법
    - OSD를 3개 이상 배포할 수 없는 테스트용 환경을 위한 inventory가 프로젝트 내 `hypercloud-sds/hack/test-sample` 경로에 존재합니다. 해당 inventory는 다음과 같은 설정이 되어 있으니 필요하신 경우 해당 inventory를 사용하시기를 바랍니다.
      - osd 3개 미만으로 배포되더라도 정상적으로 설치되도록 설정 (`cluster.yaml`에 `ConfigMap` 추가)
      - `cluster.yaml`의 `spec.storage.useAllNodes`, `spec.storage.useAllDevices`값을 `true`로 설정하여 모든 노드에서 사용 가능한 모든 device에 osd 배포를 시도
      - 배포하는 block, cephfs pool의 replication 개수를 1로 설정

- 환경에 맞게 inventory의 파일을 수정후, hcsctl로 hypercloud-sds 설치합니다.
   ``` shell
   $ hcsctl install {$inventory_name}
   # 예) hcsctl install myInventory
   # 테스트 환경용 inventory 예) hcsctl install ../hack/test-sample
   ```

    - 정상 설치가 완료되면 hypercloud-sds 를 사용하실 수 있습니다.

- hcsctl.test 로 hypercloud-sds 가 정상 설치되었는지 검증합니다.
    ``` shell
    $ hcsctl.test  
    ```
    - hypercloud-sds 정상 사용 가능 여부 확인을 위해, 여러 시나리오 테스트를 수행하게 됩니다.
    - 약 15분 정도의 시간이 소요됩니다.

### Uninstall
- hcsctl 로 설치시 사용한 inventory 를 참고하여 hypercloud-sds 를 제거합니다.
   ``` shell
   $ hcsctl uninstall {$inventory_name}
   # 예) hcsctl uninstall myInventory
   ```
    - 제거 완료 후 출력되는 메시지를 확인하여 초기화 작업이 필요한 경우 추가 작업을 완료 해야합니다. 초기화 방법은 이 [문서](./../docs/rook.md)를 참고하세요.

## Additional features
> 기본 설치, 제거 외에 효율적인 사용을 위해 여러 추가 기능 또한 제공하고 있습니다.

hcsctl 로 ceph 명령어를 수행할 수 있습니다.

``` shell
$ hcsctl ceph status
$ hcsctl ceph exec {$ceph_cmd}

# 상태 점검을 위해서 자주 사용되는 명령어는 아래와 같습니다.
$ hcsctl ceph status
$ hcsctl ceph exec ceph osd status
$ hcsctl ceph exec ceph df
```
