# CephFS (file storage) 가이드

## 설치 yaml 수정 방법

### file_system.yaml

```yaml
apiVersion: ceph.rook.io/v1
kind: CephFilesystem
metadata:
  name: myfs
  namespace: rook-ceph
spec:
  # The metadata pool spec. Must use replication.
  metadataPool:
    # The failure domain will spread teh replicas of the across different failure zones (osd, host)
    failureDomain: host
    # Set the replica size
    replicated:
      size: 3
      requireSafeReplicaSize: true
    parameters:
      # Inline compression mode for the data pool
      # Further reference: https://docs.ceph.com/docs/nautilus/rados/configuration/bluestore-config-ref/#inline-compression
      compression_mode: none
        # gives a hint (%) to Ceph in terms of expected consumption of the total cluster capacity of a given pool
      # for more info: https://docs.ceph.com/docs/master/rados/operations/placement-groups/#specifying-expected-pool-size
      #target_size_ratio: ".5"
  # The list of data pool specs. Can use replication or erasure coding.
  dataPools:
    - failureDomain: host
      replicated:
        size: 3
        # Disallow setting pool with replica 1, this could lead to data loss without recovery.
        # Make sure you're *ABSOLUTELY CERTAIN* that is what you want
        requireSafeReplicaSize: true
      parameters:
        # Inline compression mode for the data pool
        # Further reference: https://docs.ceph.com/docs/nautilus/rados/configuration/bluestore-config-ref/#inline-compression
        compression_mode: none
          # gives a hint (%) to Ceph in terms of expected consumption of the total cluster capacity of a given pool
        # for more info: https://docs.ceph.com/docs/master/rados/operations/placement-groups/#specifying-expected-pool-size
        #target_size_ratio: ".5"
  # Whether to preserve metadata and data pools on filesystem deletion
  preservePoolsOnDelete: true
  # The metadata service (mds) configuration
  metadataServer:
    # The number of active MDS instances
    activeCount: 1
    # Whether each active MDS instance will have an active standby with a warm metadata cache for faster failover.
    # If false, standbys will be available, but will not have a warm cache.
    activeStandby: true
    # The affinity rules to apply to the mds deployment
    placement:
       podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - weight: 100
            podAffinityTerm:
              labelSelector:
                matchExpressions:
                - key: app
                  operator: In
                  values:
                  - rook-ceph-mds
              # topologyKey: kubernetes.io/hostname will place MDS across different hosts
              topologyKey: kubernetes.io/hostname
          preferredDuringSchedulingIgnoredDuringExecution:
          - weight: 100
            podAffinityTerm:
              labelSelector:
                matchExpressions:
                - key: app
                  operator: In
                  values:
                  - rook-ceph-mds
              # topologyKey: */zone can be used to spread MDS across different AZ
              # Use <topologyKey: failure-domain.beta.kubernetes.io/zone> in k8s cluster if your cluster is v1.16 or lower
              # Use <topologyKey: topology.kubernetes.io/zone>  in k8s cluster is v1.17 or upper
              topologyKey: topology.kubernetes.io/zone
    # A key/value list of annotations
    annotations:
    resources:
    # The requests and limits set here, allow the filesystem MDS Pod(s) to use half of one CPU core and 1 gigabyte of memory
    #  limits:
    #    cpu: "4"
    #    memory: "4096Mi"
    #  requests:
    #    cpu: "4"
    #    memory: "4096Mi"
    priorityClassName: rook-ceph-default-priority-class
```

#### Resource setting
- `spec.metadataServer.resources`: cluster에 배포될 mds pod에 대한 resource를 설정합니다.
  -  **test환경이 아닌 production 환경일 경우 반드시 주석을 풀고 mds에 대한 resource를 설정해주시기를 바랍니다.**
  - `spec.metadataServer.resources.limits`와 `spec.metadataServer.resources.requests`에 설정되는 값들은 반드시 `동일`해야 합니다.
  - [cluster.yaml](/docs/ceph-cluster-setting.md) 파일에서 **MDS의 하드웨어 요건과 Resource setting 방법 파트**를 참고하여 작성해주시면 됩니다.

#### CephFS setting
- CephFS의 경우 metedataPool과 dataPool 두 종류의 pool를 생성하며, 각 pool에 대한 설정을 해야 합니다.
  - `failureDomain`: data의 replica를 어떻게 배치할 것인가에 대한 설정입니다. `host` 또는 `osd`가 값으로 올 수 있습니다. `failureDomain`을 host로 설정 했을 경우 데이터의 replica들은 서로 다른 host(node)에 배치되게 됩니다.
  - `replicated.size`: pool에서의 replicated size에 대한 설정입니다. 대체적으로 3을 권장하며 ceph의 성능을 위해서 2로 설정하는 경우도 있습니다.
    - `replicated.size`를 1로 설정하고 싶으신 경우, `replicated.requireSafeReplicaSize`의 값을 `false`로 변경해야 합니다.
    - `failureDomain`를 host로 설정하고 replicated size를 n으로 설정했을 경우에는 <strong>적어도 n개 이상의 노드에 osd pod가 존재</strong>해야 됩니다.
    - `failureDomain`를 osd로 설정하고 replicated size를 n으로 설정했을 경우에는 ceph cluster에 <strong>적어도 n개 이상의 osd pod가 존재</strong>해야 합니다.

### file_sc.yaml

```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: csi-cephfs-sc
provisioner: rook-ceph.cephfs.csi.ceph.com
parameters:
  # clusterID is the namespace where operator is deployed.
  clusterID: rook-ceph

  # CephFS filesystem name into which the volume shall be created
  fsName: myfs

  # Ceph pool into which the volume shall be created
  # Required for provisionVolume: "true"
  pool: myfs-data0

  # Root path of an existing CephFS volume
  # Required for provisionVolume: "false"
  # rootPath: /absolute/path

  # The secrets contain Ceph admin credentials. These are generated automatically by the operator
  # in the same namespace as the cluster.
  csi.storage.k8s.io/provisioner-secret-name: rook-csi-cephfs-provisioner
  csi.storage.k8s.io/provisioner-secret-namespace: rook-ceph
  csi.storage.k8s.io/controller-expand-secret-name: rook-csi-cephfs-provisioner
  csi.storage.k8s.io/controller-expand-secret-namespace: rook-ceph
  csi.storage.k8s.io/node-stage-secret-name: rook-csi-cephfs-node
  csi.storage.k8s.io/node-stage-secret-namespace: rook-ceph

  # (optional) The driver can use either ceph-fuse (fuse) or ceph kernel client (kernel)
  # If omitted, default volume mounter will be used - this is determined by probing for ceph-fuse
  # or by setting the default mounter explicitly via --volumemounter command-line argument.
  # mounter: kernel
reclaimPolicy: Delete
allowVolumeExpansion: true
mountOptions:
  # uncomment the following line for debugging
  #- debug
```

#### CephFS의 StorageClass 설정 방법

- `clusterID`: rook cluster가 동작하는 namespace
- `reclaimPolicy`: 동적으로 생성한 pv의 회수 정책입니다. 기본으로 `Delete`로 동작하지만, 삭제하고 싶지 않으면 `Retain`으로 설정하면 됩니다.

## CephFS 사용 예시

hcsctl로 생성한 inventory에 `file_system.yaml`을 목적에 맞게 수정하시고 `$ hcsctl install {$inventory_name}`을 수행하시면 myfs 파일시스템과 StorageClass가 생성됩니다.

docs/examples 폴더에 있는 `file-nginx.yaml`을 배포하여 file storage 사용을 테스트 할 수 있습니다.

```shell
# file_nginx.yaml를 통해 하나의 file storage를 공유하는 두 개의 nginx pod를 생성합니다.
# PVC 생성 및 nginx Deployment 생성
$ kubectl apply -f file-nginx.yam

# 배포된 pod을 확인. 서로 다른 이름을 가진 pod이 두개 생성되었습니다.
$ kubectl get pods
NAME                            READY   STATUS    RESTARTS   AGE
cephfs-nginx-64df995589-4j4n6   1/1     Running   0          62s
cephfs-nginx-64df995589-z98nn   1/1     Running   0          62

# 디렉토리가 잘 공유되었나 확인하기 위해 하나의 pod의 공유 디렉토리에 파일을 생성합니다.
# pod의 이름이 다름을 주의하세요.
$ kubectl exec -it cephfs-nginx-64df995589-4j4n6 -- touch /mnt/cephfs/testfil

# 이제 다른 하나의 pod에서 공유 디렉토리에 파일이 잘 생성되어 있나 확인합니다.
$ kubectl exec -it cephfs-nginx-64df995589-z98nn -- ls /mnt/cephf

# 아래와 같이 출력되어야 합니다.
total 4
drwxrwxrwx 2 root root    1 Nov 18 04:34 .
drwxr-xr-x 1 root root 4096 Nov 18 04:33 ..
-rw-r--r-- 1 root root    0 Nov 18 04:34 testfile
```

## References
- <https://rook.github.io/docs/rook/v1.4/ceph-filesystem.html>
