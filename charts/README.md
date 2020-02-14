# Hypercloud Storage Helm Charts
Hypercloud Storage를 Helm으로 설치 가능하게 합니다. 현재는 minikube에서만 지원합니다.

# Helm 설치
아래 환경에서 테스트되었습니다.

- Helm v3.1.0
- k8s 1.16.2

```shell
# 설치
minikube ssh "sudo mkdir -p /mnt/sda1/var/lib/rook" && helm install hypercloud-storage charts

# 제거
helm delete hypercloud-storage; minikube ssh "sudo rm -rf /var/lib/rook"; minikube ssh "sudo rm -rf /mnt/sda1/var/lib/rook"
```
