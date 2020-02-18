# Hypercloud-storage
Hypercloud Storage에서는 고가용성 스토리지(Rook), 백업(Velero), 데이터 업로드(CDI)를 지원하는 프로젝트입니다.

## Quick Start
1. (TODO: 바이너리 다운로드 내용)
2. (TODO: helm 3.0 설치 내용)
3. 설치 설정값을 지정하는 values.yaml을 수정합니다.
4. 설치 `helm install hypercloud-storage .`
5. `helm test hypercloud-storage`로 잘 설치되었는지 확인할 수 있습니다.
6. 제거 `helm uninstall hypercloud-storage`

## 지원 기능
[고가용성 스토리지(Rook)](docs/rook.md)
[백업(Velero)](docs/velero.md)
[데이터 업로드(CDI)](docs/cdi.md)

## 이슈 관리
[트러블슈팅](docs/troubleshooting.md)
이슈 등록

## 본 프로젝트는 아래와 같은 버전의 환경에서 검증되었습니다.
- Kubernetes: (TODO)
