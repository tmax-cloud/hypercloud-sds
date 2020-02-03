# HyperCloud Storage 설치 가이드

- rook (v1.2.0) 모듈을 설치한 이후 cdi (v1.11.0) 모듈을 설치합니다.

## 설치 방법

1) 가장 먼저 `Makefile` 이 존재하는 hypercloud-storage 프로젝트의 최상위 디렉토리로 이동합니다.

2) shell 창에 `make help` 를 입력하여 다음과 같이 출력되는지 확인합니다.

    ```{shell}
    Usage: make [Target ...]
      yaml           cluster_config로 부터 설치를 위한 yaml 파일을 생성합니다.
      install        생성된 yaml 파일로 hypercloud-storage를 설치합니다.
      build-test     hypercloud-storage를 test 하기위한 실행파일을 build 합니다.
      test           hypercloud-storage가 잘 설치되었는지 확인합니다. (end-to-end 테스트 수행)
      test-lint      hypercloud-storage 내의 go 소스코드에 정적분석을 실행합니다.
      uninstall      hypercloud-storage를 제거합니다.
      minikubeUp     테스트를 위한 싱글 노드 가상 환경을 만듭니다.
      minikubeClean  싱글 노드 가상 환경을 삭제합니다.
      clusterUp      테스트를 위한 멀티 노드 가상 환경을 생성합니다.
      clusterClean   멀티 노드 가상 환경을 삭제합니다.
    ```

3) (폐쇄망인 경우) 같은 디렉토리의 cluster_config 에 rook, cdi container 이미지가 저장되어있는 private registry url 을 입력합니다.
    - public network 에서 이미지를 받아와도 무방한 경우, 다음과 같이 비워두면 됩니다.
    - `dockerRegistry=`

4) `make yaml` 을 입력하고 다음과 같이 출력되는지 확인합니다.

    - ```{shell}
      ./hack/hypercloud-storage.sh yaml
      Generated yaml in ./_deploy
      ```

    - TODO) rook cluster yaml 변경 방법 확인
    - TODO) 현재 sample 용 cluster-test.yaml 의 csi-cephfs sc 를 default 로 설정 필요, reclaimPolicy 를 Delete 로 변경 필요

5) `kubectl get nodes` 를 입력하여 사용 중인 kube cluster 환경에 정상적으로 접속할 수 있는 상태인지 확인합니다.
    - 단순 install 테스트를 위한 상황인 경우 `make minikubeUp` 혹은 `make clusterUp` 을 하면 테스트용 k8s cluster 를 해당 노드에 띄울 수 있습니다.
    - `minikubeUp` : 노드 1개
    - `clusterUp` : 노드 4개

6) `hypercloud-storage/_deploy/rook/cluster/cluster-test.yaml `파일을 열고 사용하고자하는 **rook cluster 환경에 맞게 환경변수들을 수정합니다.** [rook 가이드 참고](./rook.md)
    * replication 개수 등 기타 환경변수도 수정하고 싶은 경우, _deploy/ 디렉토리 아래의 *.yaml 파일들을 수정하면 아래 7. `make install` 시 반영되어 설치됩니다.

7) `make install` 을 입력합니다.
    - `hypercloud-storage/_deploy/` 아래의 yaml 파일들을 정해진 순서대로 설치 (`kubectl create -f xxx.yaml`) 하며, 모든 deployment 들이 `status: Available` 이 될 때까지 기다립니다.
    - 해당 커맨드가 성공적으로 종료되면, rook 과 cdi 의 사용이 가능합니다.
    - TODO) default storage 세팅 필요
    - TODO) (cluster_config 에 registry url이 입력된 경우) cdi configmap 변경 [cdi 가이드 참고](./cdi.md)

  - 최종적으로 다음과 같이 출력되는지 확인합니다.

     ```{shell}
    + kubectl wait -n cdi --for=condition=available cdi.cdi.kubevirt.io/cdi --timeout=300s
    cdi.cdi.kubevirt.io/cdi condition met
    ========================== ok install cdi ==========================
    ```

----------------------------------

## 정상 설치되었는지 확인하고자 한다면 다음을 실행합니다.

8) `make build-test` 를 입력합니다.

    - 일부 prerequisite tool 들이 설치되지 않은 경우 실패하며 필요한 tool 을 알려줍니다.
    - 하나씩 설치 후 `make build-test` 를 날립니다.
    - 최종적으로 다음과 같이 출력되는지 확인합니다.

    - ```{shell}
      ========================== ok build prerequisites ==========================
      build finished
      ```

9) `make test` 를 입력합니다.

    - 해당 target 은 프로젝트에 있는 test 코드를 lint 후 test 를 실행합니다.
    - stdout 을 통해 어떤 test case 가 성공하고 어떤 test case 가 실패하는지 확인할 수 있습니다.
    - `Failed ` 개수가 0 개이면 정상 설치되었음을 확인할 수 있습니다.