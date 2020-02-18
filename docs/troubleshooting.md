# Containerized Data Importer Issues

> Containerized Data Importer (CDI) 는 pv 관리를 위한 k8s 의 add-on 으로써 kubevirt 로 vm 을 생성할 때, vm 에 mount 시킬 pvc 에 image 등의 data를 담아 생성할 수 있는 기능을 제공합니다.

## Troubleshooting Issues

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
- cdi 모듈 자체 버그
  - namespace with resourceQuota 문제 ([v1.12에서 해결](https://github.com/kubevirt/containerized-data-importer/releases/tag/v1.12.0))
  - ingress with null http 문제 ([v1.12에서 해결](https://github.com/kubevirt/containerized-data-importer/releases/tag/v1.12.0))

----------

## 이슈 분류

- DataVolume 생성 시 해당 namespace 에 importer-pod 이 임시로 생성되며 data-import 후 삭제됩니다.
- **이슈 케이스는 importer-pod 의 phase 에 따라 분류하였습니다.**

----------

### Issue #1
> namespace, resourcequota 문제

#### 상황

- dv 생성 요청 시 importer-pod 이 아예 생성되지 않은 경우 :
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
  - cdi version 을 유지하고 싶은 경우에는 namespace 에 resourceQuota 에 대응되는 default resource 를 명시한 limitRange 객체를 생성해줍니다.
    - [limitRange 생성 방법](https://kubernetes.io/docs/tasks/administer-cluster/manage-resources/memory-default-namespace/)

----------

### Issue #2
> image, registry 문제

#### 상황

- dv 의 source url 을 http 가 아닌 registry 로 적었으며, dv 생성 요청 시 해당 namespace 에 importer pod 은 생성되었으나 pod 의 Status 가 Error -> CrashLoopBackOff 를 반복하며 계속 restart 하는 경우

#### 테스트

- cdi-deployment pod 의 log 와 importer pod 의 log 확인 후 저장
- sample pod 생성하여 같은 문제가 발생하는지 확인
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

### Issue #3
> network, storage 문제

#### 상황

- dv 생성 요청 시 importer pod 이 Pending 혹은 ContainerCreating 상태로 stuck 된 경우

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
