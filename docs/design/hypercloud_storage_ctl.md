# HyperCloud Storage CTL 디자인 문서
이 문서는 HyperCloud Storage를 설치 및 제어하기 위한 HyperCloud Storage CTL(hcsctl)을 정의합니다. `hcsctl`은 단순한 Go 바이너리입니다. 내부적으로는 yaml파일들을 디플로이하고, 디플로이가 완료될 때 까지 대기하는 역할을 수행합니다. 기존의 CRD와 달리 컨테이너나 Pod으로 생성되지 않으며 minikube와 같이 단순한 바이너리로 배포됩니다.

# CRD와의 차이점
CRD 조작은 Operator의 조정루프에 의해 이루어집니다. 그런데 Operator 자체는 k8s 내부에 Pod으로 생성됩니다. 그 Pod이 또 다른 CRD, CR을 생성하고 조작합니다. 하지만 이런 과정은 클러스터 내부에서 Pod으로 조작하기 보다는 클러스터 외부에서 kubectl을 통해 조작하는것이 나을 수 있습니다. 아래 시나리오가 각각의 차이를 나타내고 있습니다.

## CRD를 사용하는 경우
```
kubectl create -f hcs_crd.yaml
# hcs_cr.yaml 파일에 rook cluster 정보를 기입
kubectl create -f hcs_cr.yaml
# 이제 Operator는 해당 CR을 보고 내부적으로 go client를 통해 rook, cdi를 배포
```

## hcsctl을 사용하는 경우
hcsctl은 내부적으로 kustomize와 kubectl을 사용합니다. 사용자는 kustomize 인벤토리를 만들고, 인스톨시에 해당 디렉토리를 전달합니다. hcsctl은 제공받은 디렉토리에서 순서에 따라 CRD를 포함한 모든 객체를 배포하고 대기합니다.

```
# 사용자는 sample 인벤토리를 복사하고 cluster 정보를 기입합니다.
cp -r ./inventory/sample ./inventory/myenv

# hcsctl install 커맨드로 인스톨을 수행합니다.
# hcsctl은 내부적으로 kustomize를 사용해서 빌드를 수행하고 yaml파일을 순서대로 배포하고 대기합니다.
hcsctl create ./inventory/myenv

# 클러스터가 잘 설치되었나 e2e 테스트를 수행합니다.
hcsctl test

# 다음 커맨드로 상태를 받아올 수 있습니다.
hcsctl stats

# 업그레이드는 먼저 yaml 파일을 수정한 뒤, 다음을 통해 적용합니다.
hcsctl apply ./inventory/myenv
```
