# hypercloud-storage
HyperCloud를 위해 제공되는 Java Client와 인스톨 파일을 포함하는 프로젝트입니다.

## 기능 별 가이드 문서
- [hypercloud-storage CICD](docs/cicd.md)
- [CDI](docs/cdi.md)

## 아래 프로젝트를 합칠 예정입니다.
- Rook-Ceph/CSI
- CDI
- snapshot
- Backup

## 아래와 같은 일이 필요합니다.
- Vagrant + Kubeadm으로 로컬에 멀티 노드 Kubernetes 구축 자동화
- Rook, CDI, Snapshot 프로젝트 머지
- helm을 통한 패키지 관리
- e2e 테스트 및 유닛 테스트
