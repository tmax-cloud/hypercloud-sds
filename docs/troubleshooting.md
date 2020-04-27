# Rook Ceph Issues



----------

# Containerized Data Importer Issues

> Containerized Data Importer (CDI) 는 pv 관리를 위한 k8s 의 add-on 으로써 kubevirt 로 vm 을 생성할 때, vm 에 mount 시킬 pvc 에 image 등의 data를 담아 생성할 수 있는 기능을 제공합니다.

## Troubleshooting Issues

### 주요 커맨드
> 이슈 상황 파악에 용이하게 쓰일 기본적인 커맨드입니다.

- `kubectl describe cdi`
  - output 의 status.Observed Version, status.Operator Version, status.Target Version 을 통해 배포된 cdi 모듈의 version 을 확인할 수 있습니다.
- `kubectl get pod -n cdi`
  - cdi namespace 에 떠있는 pod 의 목록을 확인합니다.
  - 일반적인 상황에서는 cdi-apiserver, cdi-deployment, cdi-operator, cdi-uploadproxy 의 4 개의 pod 이 떠있습니다.
- `kubectl describe pod {$PodName} -n {$PodNamespace}`
  - {$PodName} 이름을 갖는 pod 이 {$PodNamespace} 에 떠있을 때, 해당 pod 의 정보를 확인할 수 있습니다.
- `kubectl logs {$PodName} {$ContainerName} -n {$PodNamespace}`
  - {$PodName} 이름을 갖는 pod 이 {$PodNamespace} 에 떠있을 때, 해당 pod 의 {$ContainerName} 컨테이너의 로그를 확인할 수 있습니다.
    - (pod 이 single container 만 가지고 있을 경우는 kubectl logs {$PodName} -n {$PodNamespace} 만으로 확인할 수 있습니다.
  - Pod 의 status 가 ContainerCreating 혹은 CrashLoopBackoff 상태일 경우는 해당 커맨드로 로그를 확인하기 어려울 수 있습니다.
- `kubectl get sc`
  - kube cluster 환경에 등록되어있는 storageclass 목록을 확인합니다.
  - storageclassName 뒤에 (default) 라고 적혀져 있는 storageclass 가 현재 default storageclass 입니다.
    - kubernetes cluster 환경에 default sc 가 2개 이상일 경우 문제가 발생할 수 있으니 1개로 지정하기 바랍니다.

### 주요 원인 목록

> cdi 설치 및 사용 중 발생하는 대부분 에러의 원인은 다음 중 하나입니다.

- cdi 모듈의 비정상 설치
  - docker registry 등록 문제
  - image 문제
- storage class 문제
  - provioning 불가
  - mount to pod 불가
- network 문제
  - pod 에서 외부로의 통신 불가
  - pod to pod 통신 불가
  - NetworkPolicy 문제
- namespace 문제
  - resourceQuota 부족
- cdi 모듈 자체 버그
  - namespace with resourceQuota 문제 ([v1.12에서 해결](https://github.com/kubevirt/containerized-data-importer/releases/tag/v1.12.0))
  - ingress with null http 문제 ([v1.12에서 해결](https://github.com/kubevirt/containerized-data-importer/releases/tag/v1.12.0))

----------

## 이슈 분류

- **이슈 케이스는 상황(주로 importer-pod 의 phase)에 따라 분류하였습니다.**
  - DataVolume 생성 시 해당 namespace 에 importer-pod 이 임시로 생성되며 data-import 후 삭제됩니다.
- 에러 상황에 해당하는 이슈를 검색하려면 `ctrl + F`로 `!{$키워드}` 를 검색하세요.
  - 에러 **상황**을 나타내는 키워드를 검색하세요.
  - 키워드 예시) `!importer pod`, `!pending`, `!crashloopbackoff`, `!storageclass`, `!webhook`, `!insecure registry`

----------

## Issue (1)
> 원인 : 불명

> !webhook
> !cdi-apiserver pod
> !apiserver pod
> !version
> !authentication 정보
> !finalizer

#### 상황
- datavolume 생성 혹은 삭제 요청 시 다음과 같은 에러 발생
```
Internal error occurred: failed calling webhook "datavolume-mutate.cdi.kubevirt.io": Post https://cdi-api.cdi.svc:443/datavolume-mutate?timeout=30s
```

#### 테스트
- cdi 에서 관리하는 cr 의 다른 CRUD api 가 모두 같은 에러로 실패하는지 확인합니다.
  - 예) `kubectl get dv -A`, `kubectl describe dv`
  - 실패하지 않고 성공한다면 **해당 이슈가 아닙니다.**
- cdi namespace 에 pod 이 모두 `RUNNING` state 인지 확인합니다.
  - `ERROR` 등의 다른 state 로 떠있는 pod 이 있다면 **해당 이슈가 아닙니다.**

#### 해결방법
- **아직 정확한 원인 및 해결방법을 찾지 못하였습니다.**
  - 에러 상황은 kube-system ns 의 apiserver pod 에서 cdi ns 의 cdi-apiserver pod 으로 통신이 되지 않는 상황이지만 그 원인은 불명확합니다. 다음과 같은 문제가 원인일 수 있습니다:
    - cdi ns 에 걸려있는 `NetworkPolicy` 문제
    - cdi-apiserver pod 이 떠있는 `node 의 network 문제`
    - cdi 설치 시 `버전 미통일 문제`
    - cdi 에서 사용하는 `authentication 정보`가 삭제되었거나 만료된 문제
    - [비슷한 이슈](https://github.com/kubevirt/containerized-data-importer/issues/1117)
- 임시 해결방안은 cdi 모듈 전체를 제거 후 재설치하는 방안이 있습니다.
  - cdi 모듈 전체 제거 시 namespace delete 중 stuck 이 걸릴 경우 다음 링크를 참고하여 제거하시면 됩니다.
    - [namespace 강제 삭제 방법](https://success.docker.com/article/kubernetes-namespace-stuck-in-terminating)

----------

## Issue (2)
> 원인 : namespace, resourcequota

> !cdi-deployment pod
> !limits.cpu
> !importer pod
> !v1.11
> !pvcbound

#### 상황

- **dv 생성 요청 시 importer-pod 이 아예 생성되지 않은 경우**:
  - 해당 namespace 에 `kubectl get pod` 했을 시 pod 이 보이지 않으며, cdi-deployment pod 의 log 를 확인했을 때, 다음과 같은 형태의 에러 메시지가 있는 경우

```
importer pod API create errored: pods importer-XXXXX is forbidden: failed quota: {$Namespace 이름} : must specify limits.cpu, limits.memory
```

#### 테스트

- 해당 namespace 에 sample pod 이 정상 생성되는지 테스트
  - nginx pod 과 같은 임의의 pod 을 `required paramter`만 입력하여 생성 시도하여 정상 생성되는지 확인합니다.
  - 정상 생성된다면 **해당 이슈가 아닙니다.**

- cdi version 확인
  - `kubectl -n cdi describe deployments.apps cdi-deployment`
    - `Labels`에 명시된 `operator.cdi.kubevirt.io`의 버전이 `v1.12.0` 이상의 버전이면 **해당 이슈가 아닙니다.**

- dv 생성 요청한 namespace 에 resourceQuota 객체가 존재하는지 확인
  - `kubectl -n {$Namespace} get resourceQuota`
    - 존재한다면 describe 하여 `limits.cpu, limits.memory, requests.cpu, requests.memory` 에 대응되는 값이 하나라도 존재하는지 확인
    - 존재하지 않는다면 **해당 이슈가 아닙니다.**

#### 해결방법

- 위의 테스트 결과 해당 이슈가 맞다면 해결하는 방법은 다음과 같습니다.
  - cdi version 을 v1.12.0 이상의 버전으로 업그레이드하면 모두 해결됩니다.
  - cdi version 을 유지하고 싶은 경우에는 **datavolume 을 생성할 namespace** 에 resourceQuota 에 대응되는 default resource 를 명시한 limitRange 객체를 생성해줍니다.
    - [limitRange 생성 방법](https://kubernetes.io/docs/tasks/administer-cluster/manage-resources/memory-default-namespace/)

----------
## Issue (3)
> 원인 : namespace, resourcequota

> !importer pod
> !v1.11
> !resourceQuota !resourcequota !quota
> !exceeded !exceeded quota
> !pvcbound

#### 상황

- **dv 생성 요청 시 importer-pod 이 아예 생성되지 않은 경우**:
  - 해당 namespace 에 `kubectl get pod` 했을 시 pod 이 보이지 않으며, cdi-deployment pod 의 log 를 확인했을 때, 다음과 같은 형태의 에러 메시지가 있는 경우

```
import-controller.go:297] error processing pvc "hpcd-ccf03101/hpcd-d03451e2": scratch PVC API create errored: persistentvolumeclaims "hpcd-d03451e2-scratch" is forbidden: exceeded quota: hpcd-ccf03101-quota, requested: requests.storage=20Gi, used: requests.storage=82Gi, limited: requests.storage=100Gi
```

#### 테스트

- dv 생성 요청한 namespace 에 resourceQuota 객체가 존재하는지, 현재 사용량은 어떻게 되는지 확인
  - `kubectl -n {$Namespace} describe resourceQuota`

#### 해결방법

- 위의 테스트 결과 해당 이슈가 맞다면 해결하는 방법은 다음과 같습니다.
  - 단순히 해당 namespace 에 사용가능한 resourceQuota 가 부족한 것이 원인이므로 resourceQuota 를 `kubectl edit` 을 통해 늘려주면 됩니다.
  - **단, 사용량이 충분해보이는데도 생성되지 않는 경우가 cdi 모듈 특성 상 발생할 수 있습니다.**
  - 예) requested: 30Gi, used: 60Gi, limited: 100Gi 인데도 같은 이슈가 발생합니다.
    - cdi 모듈 특성 상 정상적인 작동을 위해서는 남은 disk size 가 requested size 의 **2배 이상**이 있어야 합니다.
      - 이는 cdi 모듈이 data import 를 위하여 임시로 같은 크기의 pvc 를 하나 더 만들어 사용한 뒤 import 가 완료되면 삭제하는 특징을 지니고 있기 때문에 발생합니다.
    - 따라서, requested size : 30Gi 이므로 60Gi 를 필요로 하는데, 남은 size 는 40Gi 이므로 에러가 발생합니다.

----------

## Issue (4)
> 원인 : image, registry

> !importer pod
> !crashloopbackoff
> !namespace
> !docker
> !registry
> !insecure registry
> !configmap
> !networkpolicy

#### 상황

- dv 의 source url 을 http 가 아닌 registry 로 적었으며, **dv 생성 요청 시 해당 namespace 에 importer pod 은 생성되었으나 pod 의 Status 가 Error -> CrashLoopBackOff 를 반복하며 계속 restart 하는 경우**

#### 테스트

- cdi-deployment pod 의 log 와 importer pod 의 log 확인 후 저장
- sample pod 생성하여 같은 문제가 발생하는지 확인
  - 1) busybox 이미지로 pod 을 importer pod 이 생성된 namespace, node 에 생성합니다.
  - 2) 생성된 busybox pod 에 `kubectl exec -it` 로 붙어 해당 registry 로 ping 혹은 curl 이 되는지 확인합니다.
  - 3) 서로 다른 namespace 와 서로 다른 node 에 busybox pod 2개를 생성하여 두 busybox pod 간에 통신이 되는지 확인합니다.
- docker registry 에 정확한 버전의 cdi image 들이 모두 존재하는지 확인
- dv 생성 요청시 명시한 image:tag 가 docker registry 에 있는지 확인 (curl GET 으로 확인)
- 모든 node 에서 docker registry 에 접근 가능한지 확인
  - insecure registry 등록했는지 확인
    - /etc/docker/daemon.json 모든 노드에 있는지 확인
- cdi configmap 에 추가되었는지 확인

#### 해결 방법

- sample pod 생성하여 같은 문제가 발생한다면 network 문제입니다.
  - networkpolicy 등을 확인합니다.
- docker registry 에 cdi image 가 없다면 push 해줍니다.
- dv 생성 요청시 명시한 image:tag 가 docker registry 에 없다면 push 해줍니다.
- 특정 node 에서 docker registry 에 접근 가능하지 않다면 해당 node 의 insecure registry 등록 여부 확인 후 /etc/docker/daemon.json 에 추가해줍니다.
- cdi configmap 에 registry 정보가 없다면 추가해줍니다.([cdi configmap 변경](./cdi.md))

----------

## Issue (5)
> 원인 : network, storage

> !importer pod
> !pending
> !containercreating
> !storageclass
> !cdiconfig
> !stuck

#### 상황

- **dv 생성 요청 시 importer pod 이 Pending 혹은 ContainerCreating 상태로 stuck 된 경우**

#### 테스트

- dv 생성 요청시 storage class 를 명시하였는지 테스트
  - 명시하지 않았다면 `kubectl get sc` 로 default sc 가 1개인지 확인
- dv 생성 요청시 명시한 storage class 가 현재 정상인지 테스트
  - provisoner pod 이 Running Status 인지
  - 해당 sc 로 pvc 생성 시 pv 가 정상적으로 dynamic provisioning 되는지 테스트
  - 해당 sc 의 pvc 를 mount 시킨 pod 이 정상 생성되는지 테스트
- cdiconfig 의 status.scratchSpaceStorageClass 에 적힌 storageClassName이 현재 사용가능한지 테스트
  - `kubectl describe cdiconfig` 후 `status.scratchSpaceStorageClass` 확인
- dv 가 필요로하는 pv, pvc 정상 생성 테스트
  - `kubectl get pv,pvc -n {$Namespace}`
  - dv 와 같은 이름으로 시작하는 pvc 와 dvName-scratch 이름의 pvc **2개**가 모두 생성되고 pv 와 bound 되었는지 확인
- pod 끼리 통신은 잘 되고 있는지 테스트
  - busybox image 로 pod 을 각 노드별로 create 하여 pod 간 ping 이 정상적으로 가는지 확인

#### 해결 방법

- dv 생성 요청시 storage class 를 명시하지 않았는데 default sc 가 없거나 2개 이상인 경우는 1개를 선택하여 default sc 로 만들어줍니다.
  - default sc 설정 :
    - `kubectl patch storageclass {$StorageClassName} -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'`
  - default sc 해제 :
    - `kubectl patch storageclass {$StorageClassName} -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"false"}}}'`
- dv 생성 요청한 sc 가 정상이 아니라면 해당 storageclass 를 확인합니다.
- cdiconfig 에 적힌 storageclass 가 더이상 사용 중이 아닌 storageclass 라면 다른 storageclass 로 변경합니다.([cdiconfig 변경](./cdi.md))
- dvName 으로 시작하는 pvc, pv 가 정상 생성되지 않는다면 storageclass 를 확인합니다.
- pod 끼리 통신이 되지 않는다면 network 를 확인합니다.

----------

## Issue (6)
> 원인 : 비정상 삭제, 재설치

> !cdi-operator pod
> !crd
> !v1.11.0
> !삭제 !uninstall !delete
> !설치 !install !deploy

#### 상황

- **cdi 설치 단계에서 cdi-operator pod 을 제외한 다른 pod 이 아예 생성되지 않는 경우**:
  - cdi-operator pod 의 log 를 확인했을 때, 다음과 같은 형태의 에러 메시지가 있는 경우

```
{"level":"error","ts":1587955224.906415,"logger":"kubebuilder.controller","msg":"Reconciler error","controller":"cdi-operator-controller","request":"/cdi","error":"*v1beta1.CustomResourceDefinition /datavolumes.cdi.kubevirt.io missing last applied config","stacktrace":"kubevirt.io/containerized-data-importer/vendor/github.com/go-logr/zapr.(*zapLogger).Error\n\t/go/src/kubevirt.io/containerized-data-importer/vendor/github.com/go-logr/zapr/zapr.go:128\nkubevirt.io/containerized-data-importer/vendor/sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).processNextWorkItem\n\t/go/src/kubevirt.io/containerized-data-importer/vendor/sigs.k8s.io/controller-runtime/pkg/internal/controller/controller.go:217\nkubevirt.io/containerized-data-importer/vendor/sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).Start.func1\n\t/go/src/kubevirt.io/containerized-data-importer/vendor/sigs.k8s.io/controller-runtime/pkg/internal/controller/controller.go:158\nkubevirt.io/containerized-data-importer/vendor/k8s.io/apimachinery/pkg/util/wait.JitterUntil.func1\n\t/go/src/kubevirt.io/containerized-data-importer/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:133\nkubevirt.io/containerized-data-importer/vendor/k8s.io/apimachinery/pkg/util/wait.JitterUntil\n\t/go/src/kubevirt.io/containerized-data-importer/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:134\nkubevirt.io/containerized-data-importer/vendor/k8s.io/apimachinery/pkg/util/wait.Until\n\t/go/src/kubevirt.io/containerized-data-importer/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:88"}

=> 중요 로그 : "request":"/cdi","error":"*v1beta1.CustomResourceDefinition /datavolumes.cdi.kubevirt.io missing last applied config"
```

#### 테스트

- log 에 적힌 crd 가 kubernetes cluster 에 존재하는지 확인
  - 위의 로그의 예에서는 `datavolumes.cdi.kubevirt.io` 를 확인합니다.
  - `kubectl get crd | grep datavolumes.cdi.kubevirt.io`
    - ```
      datavolumes.cdi.kubevirt.io              {$과거 날짜}
      ```
    - cdi 모듈 install 을 시도한 날짜가 아닌 {$과거 날짜} 가 조회된다면, 해당 k8s cluster 에 과거 cdi 모듈을 install 하였고, 삭제 과정에서 비정상적으로 삭제되었다는 것을 의미합니다.
  - 존재하지 않는다면 **해당 이슈가 아닙니다.**


#### 해결방법

- 위의 테스트 결과 해당 이슈가 맞다면 해결하는 방법은 다음과 같습니다.
  - `kubectl delete -f cdi-cr.yaml` 과 `kubectl delete -f cdi-operator.yaml` 을 통해 cdi 관련 resource 전체를 삭제합니다.
  - 과거 delete 할 당시 미처 삭제되지 않은 cdi 관련 crd 를 `kubectl get crd` 를 통해 조회 후 삭제합니다.
  - `kubectl apply` 를 통해 재설치하면 정상적으로 설치됩니다.
