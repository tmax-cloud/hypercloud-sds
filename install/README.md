# Hypercloud Storage Installer
Hypercloud Storage 인스톨러

## Prerequisites
- Helm (v3.0)
- Kubernetes (v1.15 ~ v1.17)
- Kubectl (v1.16)

## Quick Start
### Install
> Hypercloud Storage 를 설치합니다.

1. Helm(v3.0)을 설치합니다.: [Installing Helm](https://helm.sh/docs/intro/install/)
2. 환경 설정 파일인 `config.yaml`을 수정합니다.
  - 원하는 환경에 맞게 수정합니다.
3. 다음 커맨드로 설치를 시작합니다: `./installer.sh install`

### Uninstall
> Hypercloud Storage 를 제거합니다.

1. 다음 커맨드로 제거합니다: `./installer.sh uninstall`

### Test
> Hypercloud Storage 의 정상 설치를 확인합니다.

1. 다음과 같은 [Prerequisite](./../e2e/README.md) 을 필요로 합니다
2. 다음 커맨드로 테스트를 시작합니다: `./installer.sh test`

### Minikube 환경에서 사용하기
- install, uninstall 시 `--minikube` 를 붙여서 사용합니다:
  - `./installer.sh install --minikube`
