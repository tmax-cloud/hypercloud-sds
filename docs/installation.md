# HyperCloud Storage 설치 방법

- rook(v1.2.0) 모듈을 설치한 이후 cdi(v1.11.0) 모듈을 설치합니다.

## 방법

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
    - TODO) rook cluster yaml 변경 방법 확인

    - ```{shell}
      ./hack/hypercloud-storage.sh yaml
      Generated yaml in ./_deploy
      ```

5) `kubectl get nodes` 를 입력하여 사용 중인 kube cluster 환경에 정상적으로 접속할 수 있는 상태인지 확인합니다.
    - 단순 테스트를 위한 상황인 경우 `make minikubeUp` 혹은 `make clusterUp`을 하면 k8s cluster 를 만들 수 있습니다.

6) `make install` 을 입력합니다.
    - TODO) default storage 세팅
    - TODO) (cluster_config 에 registry url이 입력된 경우) cdi configmap 변경

- 최종적으로 다음과 같이 출력되는지 확인합니다.

  - ```{shell}
    + kubectl wait -n cdi --for=condition=available cdi.cdi.kubevirt.io/cdi --timeout=300s
    cdi.cdi.kubevirt.io/cdi condition met
    ========================== ok install cdi ==========================
    ```
