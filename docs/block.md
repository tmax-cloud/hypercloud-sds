# RBD (block storage) 가이드

## 설치 yaml 수정 방법

### block_pool.yaml

```yaml
apiVersion: ceph.rook.io/v1
kind: CephBlockPool
metadata:
  name: replicapool
  namespace: rook-ceph
spec:
  # The failure domain will spread the replicas of the data across different failure zones (osd, host)
  failureDomain: host
  replicated:
    # set the replica size
    size: 3
    # if requireSafeReplicaSize is true, Disallow setting pool with replica 1
    requireSafeReplicaSize: true
    # gives a hint (%) to Ceph in terms of expected consumption of the total cluster capacity of a given pool
    #targetSizeRatio: .5
  #crushRoot: my-root
  # The Ceph CRUSH device class associated with the CRUSH replicated rule
  #deviceClass: my-class
  compressionMode: none
  annotations:
```

#### CephBlockPool 설정 방법

- `failureDomain`: data의 replica를 어떻게 배치할 것인가에 대한 설정입니다. `host` 또는 `osd`가 값으로 올 수 있습니다. `failureDomain`을 host로 설정 했을 경우 데이터의 replica들은 서로 다른 host(node)에 배치되게 됩니다.
- `replicated.size`: pool에서의 replicated size에 대한 설정입니다. 대체적으로 3을 권장하며 ceph의 성능을 위해서 2로 설정하는 경우도 있습니다.
    - `replicated.size`를 1로 설정하고 싶으신 경우, `replicated.requireSafeReplicaSize`의 값을 `false`로 변경해야 합니다.
    - `failureDomain`를 host로 설정하고 replicated size를 n으로 설정했을 경우에는 <strong>적어도 n개 이상의 노드에 osd pod가 존재</strong>해야 됩니다.
    - `failureDomain`를 osd로 설정하고 replicated size를 n으로 설정했을 경우에는 ceph cluster에 <strong>적어도 n개 이상의 osd pod가 존재</strong>해야 합니다.

### block_sc.yaml

```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
   name: rook-ceph-block
provisioner: rook-ceph.rbd.csi.ceph.com
parameters:
    clusterID: rook-ceph
    pool: replicapool
    imageFormat: "2"
    imageFeatures: layering

    # The secrets contain Ceph admin credentials. These are generated automatically by the operator
    # in the same namespace as the cluster.
    csi.storage.k8s.io/provisioner-secret-name: rook-csi-rbd-provisioner
    csi.storage.k8s.io/provisioner-secret-namespace: rook-ceph
    csi.storage.k8s.io/controller-expand-secret-name: rook-csi-rbd-provisioner
    csi.storage.k8s.io/controller-expand-secret-namespace: rook-ceph
    csi.storage.k8s.io/node-stage-secret-name: rook-csi-rbd-node
    csi.storage.k8s.io/node-stage-secret-namespace: rook-ceph
    # Specify the filesystem type of the volume. If not specified, csi-provisioner
    # will set default as `ext4`.
    csi.storage.k8s.io/fstype: ext4
allowVolumeExpansion: true
reclaimPolicy: Delete
```

#### CephBlockPool의 StorageClass 설정 방법

- `clusterID`: rook cluster가 동작하는 namespace
- `csi.storage.k8s.io/fstype`: rbd를 pod의 파일 경로에 마운트할때, rbd를 해당 파일시스템으로 초기화합니다. 설정하지 않으면 default로 ext4를 사용합니다.
- `reclaimPolicy`: 동적으로 생성한 pv의 회수 정책입니다. 기본으로 `Delete`로 동작하지만, 삭제하고 싶지 않으면 `Retain`으로 설정하면 됩니다.

## RBD 사용 예시

hcsctl로 생성한 inventory에 `block_pool.yaml`파일을 목적에 맞게 수정하시고 `$ hcsctl install {$inventory_name}`을 수행하시면 BlockPool과 StorageClass가 생성됩니다.

docs/examples 폴더에 있는 `block-wordpress.yaml`과 `block-mysql.yaml`을 배포하여 block storage 사용을 테스트 할 수 있습니다.

```shell
# mysql과 wordpress 배포
# docs/examples 폴더에 존재하는 yaml 파일입니다.
$ kubectl apply -f block-mysql.yaml
$ kubectl apply -f block-wordpress.yam

# minikube를 사용중인 경우 다음 명령어로 Wordpress의 URL을 확인할 수 있습니다
# 출력된 url로 접근하면 wordpress가 잘 배포되었음을 확인할 수 있습니다.
$ echo http://$(minikube ip):$(kubectl get service wordpress -o jsonpath='{.spec.ports[0].nodePort}'

# toolbox을 통해 rbd du -p {$poolName}을 입력하면 pod에 할당된 RBD image의 사용 및 할당량을 확인할 수 있습니다.
$ rbd du -p replicapool
warning: fast-diff map is not enabled for csi-vol-4af65546-4717-11ea-96fa-ea47d19dc5d1. operation may be slow.
warning: fast-diff map is not enabled for csi-vol-4f0cd143-4717-11ea-96fa-ea47d19dc5d1. operation may be slow.
NAME                                               PROVISIONED  USED
csi-vol-4af65546-4717-11ea-96fa-ea47d19dc5d1       1 GiB        184 MiB
csi-vol-4f0cd143-4717-11ea-96fa-ea47d19dc5d1       1 GiB        84 MiB
<TOTAL>                                            2 GiB        268 MiB
```

## RBD 사용시 주의 사항

- PVC 생성시 `spec.volumeMode.Block`을 명시하지 않으면, block volume 생성 후 자동으로 Storage Class에 명시한 filesystem으로 포맷됩니다.
pod에 attach될 때에는 pod 배포 시 명시한 `volumeMounts.mountPath` 디렉토리에 volume이 filesystem 형태로 mount되며, 해당 방법에서는 accessMode로 `RWO`는 지원하나 `RWX`를 지원하지 않습니다.

- PVC 생성시 `spec.volumeMode.Block`을 명시할 경우, block volume은 자동으로 filesystem으로 포맷되지 않으며, pod에 attach될 때에는 pod 배포 시 명시한 `volumeDevices.devicePath` 경로에 raw block device 형태로 attach됩니다.
해당 방법에서는 accessMode로 `RWO`, `RWX`를 모두 지원합니다.

- RBD 사용을 위한 storageClass 작성시 `reclaimPolicy`를 `Retain`으로 설정할 경우, pv를 지워도 RBD image가 ceph cluster에 남아 있습니다. 이 경우에는 직접 `ceph 명령어`를 사용하여 RBD image를 지워야 합니다. RBD image를 지우는 방법은 <strong>[ceph 명령어 메뉴얼](/docs/ceph-command.md)</strong>을 참고하시기 바랍니다.

## References

- <https://rook.io/docs/rook/v1.3/ceph-block.html>
