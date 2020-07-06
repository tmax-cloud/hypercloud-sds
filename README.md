# Hypercloud-Storage
Hypercloud Storage 는 고가용성 스토리지(Rook) 와 이미지 임포팅(CDI)을 지원하는 프로젝트입니다.
두 모듈의 편리한 설치, 제거, 관리를 위해서 hcsctl 바이너리 파일을 제공하고 있습니다.
hcsctl 최신 바이너리는 프로젝트 [release page](http://192.168.1.150:10080/ck3-4/hypercloud-storage/releases) 에서 확인 가능합니다.

## 주요기능
- 설치를 위한 sample yaml 파일 생성
- hypercloud-storage 설치
- hypercloud-storage 제거
- ceph 명령 수행 및 출력

> 더 자세한 사항은 `hcsctl help` 를 참고하시기 바랍니다.


## 시작하기
- hcsctl 바이너리 다운로드와 자세한 사용법은 [hcsctl 문서](hcsctl/README.md)를 참고하시기 바랍니다.

## Documentation
- [고가용성 스토리지(Rook)](docs/rook.md)
- [데이터 업로드(CDI)](docs/cdi.md)

## 이슈 등록

이슈 등록 전에 먼저 [트러블 슈팅 페이지](docs/troubleshooting.md)를 확인해주시기 바랍니다. 트러블 슈팅 페이지를 통해 해결되지 않는 버그는 [IMS](https://ims.tmaxsoft.com/)나 프로젝트 내 [issue 페이지](http://192.168.1.150:10080/ck3-4/hypercloud-storage/issues)에서 등록할 수 있습니다. IMS 이슈 등록시에는 Product는 HyperCloud, Module은 K8SStorage로 등록하시면 됩니다.

## Compatibility
- 본 프로젝트는 아래와 같은 버전에서 검증되었습니다.
    - Kubernetes
        - v1.17
        - v1.16
        - v1.15
    - OS
        - ubuntu 18.04
        - centos 8.1, 7.7
        - prolinux 7.5

## Contact
- CK 3-4 팀
