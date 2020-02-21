# Hypercloud Storage Installer
Hypercloud Storage 인스톨러

## Prerequisites
- Helm (>3.0)
- Kubernetes (>?)

## Quick Start
### Install
1. Helm(>3.0)을 설치합니다.: [Installing Helm](https://helm.sh/docs/intro/install/)
2. 환경 설정 파일인 `config.yaml`을 수정합니다.
3. 다음 커맨드로 설치를 시작합니다: `./installer.sh install`

### Uninstall
1. 다음 커맨드로 제거합니다: `./installer.sh uninstall`

### Minikube 환경에서 사용하기
- install, uninstall 시 --minikube 를 붙여서 사용합니다.
