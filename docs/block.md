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
  # Set the replica size
  replicated:
    size: 3
    # Disallow setting pool with replica 1, this could lead to data loss without recovery.
    # Make sure you're *ABSOLUTELY CERTAIN* that is what you want
    requireSafeReplicaSize: true
  # Ceph CRUSH root location of the rule
  # For reference: https://docs.ceph.com/docs/nautilus/rados/operations/crush-map/#types-and-buckets
  #crushRoot: my-root
  # The Ceph CRUSH device class associated with the CRUSH replicated rule
  # For reference: https://docs.ceph.com/docs/nautilus/rados/operations/crush-map/#device-classes
  #deviceClass: my-class
  # Enables collecting RBD per-image IO statistics by enabling dynamic OSD performance counters. Defaults to false.
  # For reference: https://docs.ceph.com/docs/master/mgr/prometheus/#rbd-io-statistics
  # enableRBDStats: true
  # Set any property on a given pool
  # see https://docs.ceph.com/docs/master/rados/operations/pools/#set-pool-values
  parameters:
    # Inline compression mode for the data pool
    # Further reference: https://docs.ceph.com/docs/nautilus/rados/configuration/bluestore-config-ref/#inline-compression
    compression_mode: none
    # gives a hint (%) to Ceph in terms of expected consumption of the total cluster capacity of a given pool
    # for more info: https://docs.ceph.com/docs/master/rados/operations/placement-groups/#specifying-expected-pool-size
    #target_size_ratio: ".5"
  # A key/value list of annotations
  annotations:
  #  key: value
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
# Change "rook-ceph" provisioner prefix to match the operator namespace if needed
provisioner: rook-ceph.rbd.csi.ceph.com
parameters:
    # clusterID is the namespace where the rook cluster is running
    # If you change this namespace, also change the namespace below where the secret namespaces are defined
    clusterID: rook-ceph

    # If you want to use erasure coded pool with RBD, you need to create
    # two pools. one erasure coded and one replicated.
    # You need to specify the replicated pool here in the `pool` parameter, it is
    # used for the metadata of the images.
    # The erasure coded pool must be set as the `dataPool` parameter below.
    #dataPool: ec-data-pool
    pool: replicapool

    # RBD image format. Defaults to "2".
    imageFormat: "2"

    # RBD image features. Available for imageFormat: "2". CSI RBD currently supports only `layering` feature.
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
    # will set default as `ext4`. Note that `xfs` is not recommended due to potential deadlock
    # in hyperconverged settings where the volume is mounted on the same node as the osds.
    csi.storage.k8s.io/fstype: ext4
# uncomment the following to use rbd-nbd as mounter on supported nodes
# **IMPORTANT**: If you are using rbd-nbd as the mounter, during upgrade you will be hit a ceph-csi
# issue that causes the mount to be disconnected. You will need to follow special upgrade steps
# to restart your application pods. Therefore, this option is not recommended.
#mounter: rbd-nbd
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

- RBD 사용을 위한 storageClass 작성시 `reclaimPolicy`를 `Retain`으로 설정할 경우, pv를 지워도 RBD image가 ceph cluster에 남아 있습니다. 이 경우에는 직접 `ceph 명령어`를 사용하여 RBD image를 지워야 합니다. RBD image를 지우는 방법은 <strong>[ceph 명령어 메뉴얼](ceph-command.md)</strong>을 참고하시기 바랍니다.

## References

- <https://rook.github.io/docs/rook/v1.4/ceph-block.html>
