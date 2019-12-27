# hypercloud-storage
HyperCloud를 위해 제공되는 Java Client와 인스톨 파일을 포함하는 프로젝트입니다.

## 아키텍쳐
![스크린샷__2019-12-27_13-53-40](http://192.168.1.150:10080/ck3-4/hypercloud-storage/wikis/uploads/3b7e3dc851ea0cd809090c00899b04b4/hypercloud-storage_architecture.png)

## 아래 프로젝트를 합칠 예정입니다.
- Rook-Ceph, CSI
- CDI, CDI java client
- snapshot client
- Backup
- 
## 아래와 같은 일이 필요합니다 (helper 요청)
- Vagrant + Kubeadm으로 로컬에 멀티 노드 Kubernetes 구축 자동화
- Rook, CDI, Snapshot 프로젝트 머지
- helm을 통한 패키지 관리
- e2e 테스트 및 유닛 테스트