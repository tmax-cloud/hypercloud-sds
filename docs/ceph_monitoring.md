# Rook Ceph Monitoring
> 본 문서에서는 prometheus operator를 사용하여, rook의 ceph cluster를 monitoring 하는 방법을 설명합니다.

## Rook Ceph Monitoring Service OverView
- hyperCloud storage에서는 prometheus를 사용하여, rook으로 구축한 ceph cluster를 monitoring하고자 합니다.
- rook 공식 github 문서와 같이 prometheus operator를 사용하여 ceph monitoring service를 구축합니다.
- 설치되는 prometheus operator의 버전은 `0.37`입니다.
- 설치되는 prometheus operator 관련 cr과 crd는 hyperCloud storage에서 제공하는 rook ceph 버전(`rook/ceph:v1.2.5`)을 따릅니다.

## Rook Ceph Monitoring Service build 과정
> 아래의 과정은 hyperCloud storage가 정상적으로 설치된 후 진행되어야 합니다.
### Prometheus operator deploy
```shell
# 아래의 yaml 파일은 /docs/examples에 존재합니다.
$ kubectl apply -f prometheus-bundle.yaml

# Prometheus operator의 생성을 확인합니다.
$ kubectl get pod
NAME                                   READY   STATUS    RESTARTS   AGE
prometheus-operator-5b855c4d9d-rt5bw   1/1     Running   0          18s
```

### Prometheus CR deploy
- ServiceMonitor와 Prometheus CR를 배포합니다.
    - ServiceMonitor CR를 통해 prometheus가 monitoring하고자 하는 k8s service를 정의할 수 있습니다.
	- Prometheus CR를 통해 prometheus 관련 instances 설정 및 사용하고자 하는 serviceMonitor를 정의할 수 있습니다.
```shell
# ServiceMonitor와 Prometheus에 대한 yaml deploy
# 아래의 yaml파일들은 /docs/examples에 존재합니다.
$ kubectl create -f ceph-service-monitor.yaml
$ kubectl create -f ceph-prometheus.yaml
$ kubectl create -f ceph-prometheus-service.yaml

# prometheus pod 생성 확인
$ kubectl -n rook-ceph get pod prometheus-rook-prometheus-0
NAME                           READY   STATUS    RESTARTS   AGE
prometheus-rook-prometheus-0   3/3     Running   1          51s

# Prometheys Web Console에 접근
$ echo "http://$(kubectl -n rook-ceph -o jsonpath={.status.hostIP} get pod prometheus-rook-prometheus-0):30900
http://192.168.50.90:30900
```
## 참조
- https://github.com/rook/rook/blob/v1.2.5/Documentation/ceph-monitoring.md
- https://github.com/coreos/prometheus-operator/tree/release-0.37