# Rest API 디자인 문서

Hypercloud 4 Storage 서비스 설계 진행 내역과 설계시에 참고한 타 클라우드 플랫폼의 storage 서비스 rest api를 조사 및 비교한 자료를 포함한 디자인 문서입니다.

### 테스트 환경 정보
  - Rancher v2.x 
    - 싱글노드 클러스터 설치 - [Vagrant + Virtualbox](https://rancher.com/docs/rancher/v2.x/en/quick-start-guide/deployment/quickstart-vagrant/)
  - Openshift v4.3 
    - 싱글노드 클러스터 설치 - [CodeReady Containers](https://code-ready.github.io/crc/) 

### Storage Resources

스토리지 자원은 아래와 같이 크게 3개의 종류로 분류할 수 있습니다.

- PersistentVolume
- StorageClass
- PersistentVolumeClaim

### Storage User

Storage 자원 사용자는 아래와 같이 크게 두 종류가 있습니다.

- Admin
  - 모든 자원을 관리하는 super user
- General User
  - Admin이 지정해 준 role에 따라 지정 된 resource만 관리하는 일반 사용자

## Hypercloud 4 Rest API

version 4 에서는 이전 버전과 달리 kubernetes api 를 그대로 사용하여 스토리지 서비스를 제공하기로 결정 되어, rancher 와 openshift 기능의 합집합을 제공하는 방향으로 설계 진행 되었습니다.

- label, annotation 을 고칠 수 있는 PATCH method 추가
- 여러개의 리소스를 한꺼번에 지울 수 있는 DELETE method 추가

### PersistentVolume

- pv list: `GET /api/v1/persistentvolumes`
- pv get: `GET /api/v1/persistentvolumes/{name}`
- pv create: `POST /api/v1/persistentvolumes`
- pv edit: `PUT /api/v1/persistentvolumes/{name}`
- pv delete: `DELETE /api/v1/persistentvolumes/{name}`
- pv delete collection: `DELETE /api/v1/persistentvolumes`
- pv update: `PATCH /api/v1/persistentvolumes/{name}`

### StorageClass

- sc list: `GET /apis/storage.k8s.io/v1beta1/storageclasses`
- sc get: `GET /apis/storage.k8s.io/v1beta1/storageclasses/{name}`
- sc create: `POST /apis/storage.k8s.io/v1beta1/storageclasses`
- sc edit: `PUT /apis/storage.k8s.io/v1beta1/storageclasses/{name}`
- sc delete: `DELETE /apis/storage.k8s.io/v1beta1/storageclasses/{name}`
- sc delete collection: `DELETE /apis/storage.k8s.io/v1beta1/storageclasses`
- sc update: `PATCH /apis/storage.k8s.io/v1/storageclasses/{name}`

### PersistentVolumeClaim

- pvc list: `GET /api/v1/namespaces/{namespace}/persistentvolumeclaims`
- pvc list: `GET /api/v1/persistentvolumeclaims`
- pvc get: `GET /api/v1/namespaces/{namespace}/persistentvolumeclaims/{name}`
- pvc create: `POST /api/v1/namespaces/{namespace}/persistentvolumeclaims`
- pvc edit: `PUT /api/v1/namespaces/{namespace}/persistentvolumeclaims/{name}`
- pvc delete: `DELETE /api/v1/namespaces/{namespace}/persistentvolumeclaims/{name}`
- pvc delete collection: `DELETE /api/v1/namespaces/{namespace}/persistentvolumeclaims`
- pvc update: `PATCH /api/v1/namespaces/{namespace}/persistentvolumeclaims/{name}`

### Storage User 권한

- Admin
  - 모든 자원 pv, pvc (volume), sc 관리
- General User
  - 자원 관리 하지 않고, read도 할 수 없음
  - template 으로 앱 생성시에 pvc (volume) size 만 지정 가능


## 플랫폼 별 User 권한 비교

각 플랫폼 사용자별 권한 및 관리 가능한 자원은 아래와 같습니다.

### HyperCloud 
- Admin
  - pv 관리 
  - volume type은 NFS만 가능하도록 고정됨
  - sc 관리 기능 없음
- General User
  - pvc (storage) 관리 
  - sc 선택 가능, pv 지정 기능 없음

### Rancher
- Admin
  - 모든 자원 pv, pvc (volume), sc 관리
- General User
  - pv (Read only)
  - sc (Read only)
  - pvc (volume) 관리

### Openshift
- Admin
  - 모든 자원 pv, pvc (volume), sc 관리
- General User
  - 자원 관리 하지 않고, read도 할 수 없음
  - template 으로 앱 생성시에 pvc (volume) size 만 지정 가능

## 플랫폼 별 Rest API 비교

Rancher, Openshift 는 모두 쿠버네티스 api를 그대로 사용하고 있기 때문에 필수값 등 자세한 api 사용법은 [k8s api reference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/) 를 참고하면 좋습니다. 

- [k8s pvc reference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#persistentvolumeclaim-v1-core)
- [k8s sc reference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#storageclass-v1-storage-k8s-io)
- [k8s pv reference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#persistentvolume-v1-core)

[Openshift api reference](https://docs.openshift.com/container-platform/4.3/rest_api/index.html) 도 기본 스토리지 자원들의 경우엔 k8s api 와 큰차이 없지만, openshift 자체적으로 만든 KubeStorageVersionMigrator 관련 rest api가 존재 합니다.

### PersistentVolume Rest APIs

각 플랫폼 별 사용하고 있는 PersistentVolume (pv) 관련 rest api 는 아래와 같습니다.

#### HyperCloud

- pv list: `GET /v3/_api/pvs`
- pv create: `POST /v3/_api/pvs`
- pv delete: `DELETE /v3/_api/pvs`

#### Rancher

- pv list: `GET /api/v1/persistentvolumes`
- pv get: `GET /api/v1/persistentvolumes/{name}`
- pv create: `POST /api/v1/persistentvolumes`
  - name
  - volume source
  - size
  - accessMode
- pv edit: `PUT /api/v1/persistentvolumes/{name}`
- pv delete: `DELETE /api/v1/persistentvolumes/{name}`
  - 상태가 Bound 이면 삭제 불가

#### Openshift

- 다른 operation 들을 모두 rancher와 동일
- pv update: `PATCH /api/v1/persistentvolumes/{name}`
  - labels, annotation 수정 시에 사용됨 

---

### StorageClass Rest APIs

각 플랫폼 별 사용하고 있는 StorageClass (sc) 관련 rest api 는 아래와 같습니다.

#### HyperCloud

- sc list: `GET /v3/_api/storage-class`

#### Rancher

- sc list: `GET /apis/storage.k8s.io/v1beta1/storageclasses`
- sc create: `POST /apis/storage.k8s.io/v1beta1/storageclasses`
  - name
  - provisioner
- sc edit: `PUT /apis/storage.k8s.io/v1beta1/storageclasses/{name}`
- sc delete: `DELETE /apis/storage.k8s.io/v1beta1/storageclasses/{name}`
- sc get: `GET /apis/storage.k8s.io/v1beta1/storageclasses/{name}`

#### Openshift

- 다른 operation 들을 모두 rancher와 동일
- sc update: `PATCH /apis/storage.k8s.io/v1/storageclasses/{name}`
  - labels, annotation 수정 시에 사용됨 

---

### PersistentVolumeClaim Rest APIs

각 플랫폼 별 사용하고 있는 PersistentVolumeClaim (pvc) 관련 rest api 는 아래와 같습니다.

#### HyperCloud

- storage list:
  - `GET /v3/_api/domains/{domainId}/available-storages`
  - `GET /v3/_api/domains/{domainId}/container-storages`
- storage get:
  - `GET /v3/_api/domains/{domainId}/container-storages/{storageId}`
  - `GET /v3/_api/domains/{domainId}/container-storages/{storageId}/mount-infos`
- storage create: `POST /v3/_api/domains/{domainId}/container-storages`
- storage delete: `DELETE /v3/_api/domains/{domainId}/container-storages`
- storage edit: `PUT /v3/_api/domains/{domainId}/container-storages/{storageId}`

#### Rancher 

- volume list: `GET /api/v1/namespaces/{namespace}/persistentvolumeclaims`
- volume get: `GET /api/v1/namespaces/{namespace}/persistentvolumeclaims/{name}`
- volume create: `POST /api/v1/namespaces/{namespace}/persistentvolumeclaims`
  - name
  - source
  - accessMode

#### Openshift

- 다른 operation 들을 모두 rancher와 동일
- pvc update: `PATCH /api/v1/namespaces/{namespace}/persistentvolumeclaims/{name}`
  - labels, annotation 수정과 expand pvc 시에 사용됨
- pvc edit: `PUT /api/v1/namespaces/{namespace}/persistentvolumeclaims/{name}`
- pvc delete: `DELETE /api/v1/namespaces/{namespace}/persistentvolumeclaims/{name}`
- pvc list: `GET /api/v1/persistentvolumeclaims`
  - namespace를 all projects 로 선택 했을 때 모든 pvc를 볼 수 있음

---

### CSIDriver Rest APIs

Openshift 의 경우에는 storage class 생성할 때 csi driver 조회가 이뤄지고 있습니다. 다른 플랫폼에서는 사용 경우가 UI 노출로는 확인되진 않았습니다. 

#### Openshift

- csidriver list: `GET /apis/storage.k8s.io/v1beta1/csidrivers`
