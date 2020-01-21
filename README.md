# Hypercloud-storage
HyperCloud를 위해 제공되는 (Java Client 와) 인스톨 파일 그리고 정상 각 모듈의 설치 확인을 위한 테스트 스크립트를 포함하는 프로젝트입니다.

## 기능 별 가이드 문서
- [hypercloud-storage CICD](docs/cicd.md)
- [CDI](docs/cdi.md)

## 아래 프로젝트를 합칠 예정입니다.
- Rook-Ceph/CSI
- CDI
- snapshot
- Backup (velero)

## 아래와 같은 추가 작업이 필요합니다.
- 멀티 노드 Kubernetes Cluster 구축 자동화
  - on Virtual Machine
  - on Physical Node
- Rook, CDI, Snapshot, Backup 프로젝트 merge
  - Rook 프로젝트의 cluster.yaml 설정 방법 수정
- 패키지 관리 방법 결정 후 구축
  - helm
- e2e 테스트 추가
  - test framework 결정
  - 필요한 테스트 시나리오 정리 후 구현
- 폐쇄망 환경 고려

## 현재 지원하는 기능 목록입니다. (20.01.21 기준)
- 4 node Kubernetes Cluster 구축 자동화
  - using Virtual Machine
- Rook, CDI 설치
  - 현재 rook 의 cluster 설정은 테스트용으로 임의로 설정되어있습니다.
- go-client 를 사용하여 정상 설치 확인 명령 날리는 example 코드 (test framework 를 적용하여 단언문을 포함하는 테스트 코드로의 변환이 필요하며, main package 및 main.go 파일을 하나만 두도록 변경해야 합니다.)
  - pkg/test-installation : rook, cdi 관련 deployment 들이 모두 떠있는지 test
  - pkg/test-pod-networking : pod to pod ping, pod to "google.com" ping 정상적으로 가는지 test

## gitlab-ci 파이프라인 관련 정보
- ck3-4 팀환경의 172.22.4.101 (ck34-1) 노드를 사용하고 있습니다.
  - Ubuntu 18.04.3 LTS (GNU/Linux 4.15.0-72-generic x86_64)
- 해당 노드에서 `gitlab-runner` 유저로 pipeline 이 실행됩니다.
- 본 프로젝트의 테스트는 해당 노드에 직접 테스트를 진행하는 것이 아닌 vm 을 생성하여 테스트를 진행한 후 삭제합니다.

## 본 프로젝트는 아래와 같은 버전의 환경에서 검증되었습니다.
- make : GNU 4.1
- vboxmanage : 5.2.34
- vagrant : 2.2.6
- go : 1.13.6
- kubectl : v1.17.0 (Client), v1.17.1 (Server)
- kubernetes go client :  v13.0 (for k8s version v1.15.x)
- k8s-vagrant-multi-node : ..
  