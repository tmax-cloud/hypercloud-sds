# cluster.yaml 수정 메뉴얼
본 메뉴얼에서는 cluster.yaml를 수정하는 방법에 대해서 설명합니다. 메뉴얼에서 설명하는 옵션에 대해서만 수정하는 것을 권장합니다.

## 최소요건 및 권장사항
클러스터 배포시 최소요건 및 권장사항은 다음과 같습니다.
- 하드웨어 요건
  - OSD Memory: 1TB당 1GB 정도의 메모리 권장
  - OSD disk: 최소 10GB 이상 디스크 권장, 데이터 저장용으로 큰 용량의 디바이스나 디렉토리를 사용하는 게 좋음
    - 디렉토리로 root filesystem 사용시 kubernetes 상에서 <strong>`disk pressure`</strong> 발생할 수 있으므로 주의
  - Ceph Mon memory: 1GB 정도의 메모리 권장
  - Ceph Mon disk: 10GB 정도의 디스크 권장
- 권장사항
  - 각 노드마다 OSD를 배포하도록 권장 (Taint 걸린 host 없는 걸 확인해야함)
  - 총 OSD 개수는 3개 이상으로 권장
  - CephFS 및 RBD pool 설정 시 Replication 개수 3개 권장

## cluster.yaml 수정 방법
위의 예시는 cluster.yaml 파일을 설명을 위해 수정한 예시입니다. hcsctl 명령으로 만들어진 인벤토리의 rook/cluster.yaml 파일을 환경에 맞게 수정해야 합니다.
```yaml
apiVersion: ceph.rook.io/v1
kind: CephCluster
metadata:
  name: rook-ceph
  namespace: rook-ceph
spec:
  cephVersion:
    image: ceph/ceph:v14.2.4-20190917
    allowUnsupported: true
  dataDirHostPath: /var/lib/rook
  skipUpgradeChecks: false
  mon:
    count: 3
    allowMultiplePerNode: true  # If more than 3 nodes are available, false is recommended.
  dashboard:
    enabled: true
    ssl: true
  monitoring:
    enabled: false  # Require Prometheus to be pre-installed
    rulesNamespace: rook-ceph
  network:
    hostNetwork: false
  rbdMirroring:
    workers: 0
  mgr:
    modules:
    # The pg_autoscaler is only available on nautilus or newer. remove this if testing mimic.
    - name: pg_autoscaler
      enabled: true
  storage:
#   useAllNodes: true      # Apply ceph-osd to all nodes in the same way.
    useAllNodes: false     # Apply ceph-osd to specific nodes.
    useAllDevices: false
    deviceFilter:
    config:
      databaseSizeMB: "1024" # This value can be removed for environments with normal sized disks (100 GB or larger)
      journalSizeMB: "1024"  # This value can be removed for environments with normal sized disks (20 GB or larger)
      osdsPerDevice: "1"   # This value can be overridden at the node or device level
#   directories:
#   - path: /var/lib/rook
    nodes:
      - name: "worker1"    # Add worker1 node to ceph-osd. (Caution: check hostname by 'kubectl get nodes')
        devices:           # Add disk of worker1 to ceph-osd.
        - name: "sdc"      # Caution: Disk must exist in worker1 node. (check disk by 'sudo fdisk -l')
      - name: "worker2"         # Add worker2 node to ceph-osd.
        directories:            # Add a directory of worker2 to ceph-osd.
        - path: "/root/testdir" # Caution: Directory must exist in worker2 node.
        devices:
        - name: "sdc"           # Add disk of worker2 to ceph-osd. disk must exist in the nodes.
```
### Mon deploy setting
- `spec.mon.count`: kube cluster에 deploy할 mon의 개수를 의미합니다. `count`의 값은 1부터 9 사이의 홀수이어야 합니다. Ceph document에서는 기본적으로 3개의 mon를 권장합니다.
- `spec.mon.allowMultiplePerNode`: `true` 또는 `false`를 값으로 가질 수 있으며, 하나의 노드에 여러 개의 mon를 deploy할 수 있는지에 대한 여부를 결정하는 옵션입니다. `false`로 설정 할 경우, 하나의 node에 여러 개의 mon pod이 deploy될 수 없습니다.

### OSD deploy setting
- `spec.storage.useAllNodes`: OSD pod에 deploy관련 설정입니다. Pod를 deploy할 수 있는 모든 노드에 동일한 조건의 OSD pod를 deploy하고 싶으면 `true`로, 각 노드별로 다른 설정의 OSD를 deploy하고 싶으면 `false`로 설정합니다.
  -  `true`일 경우, `spec.storage.directories`: 를 통해 모든 노드에서 OSD가 배포될 directory를 설정할 수 있습니다.
- `spec.storage.nodes`: 각 노드 별로 deploy되는 OSD pod의 설정을 다르게 하고 싶은 경우 해당 설정에 osd에 관한 설정을 명시하면 됩니다.
    - `name`: OSD pod가 deploy될 node 이름을 명시합니다. 해당 값은 `kubernetes.io/hostname`과 동일해야 됩니다.
    - `devices`: OSD pod를 device 위에 deploy하겠다는 옵션입니다. 해당 옵션에서는 `name`에 device 이름을 명시하면 됩니다. 명시되는 device는 아래의 조건을 만족해야 됩니다.(만족하지 않을 경우, OSD pod가 deploy되지 않는다.)
        - Deploy되는 노드에 해당 디바이스가 반드시 존재해야 하며 unmount된 상태여야 합니다.
        - 해당 device의 초기화가 제대로 되어 있어야 하며 초기화 제대로 안 되어있을 경우 아래의 명령어를 통해 초기화합니다.(당연히 초기화후 원래 있었던 데이터는 복구할 수 없습니다.)
            ```shell
            # /dev/sdb device에 osd 설치를 원하는 경우
            DISK="/dev/sdb"
            $ sudo sgdisk --zap-all $DISK
            $ sudo if=/dev/zero of="$DISK" bs=1M count=100 oflag=direct,dsync
            ```
    - `directories`: OSD pod를 directory 위에 deploy하겠다는 옵션입니다. 해당 옵션에서는 `path`에 directory이름을 명시하면 됩니다. 명시되는 directory는 아래의 조건을 만족해야 됩니다.(만족하지 않을 경우, OSD pod가 deploy되지 않습니다.)
        - Deploy되는 노드에 해당 디렉토리가 반드시 존재해야 됩니다.
        - OSD를 directory 위에 deploy할 경우, storage 사용에 대한 monitoring이 정상적으로 수행되지 않는 현상을 발견했습니다. 그러므로 <strong>OSD pod를 device위에 deploy</strong>하는 것을 권장합니다.
    - 하나의 노드에 두 개 이상의 OSD pod를 deploy하고 싶은 경우 아래와 같이 osd 설정을 추가하면 됩니다.
        ```yaml
        nodes:
         - name: "worker1"   
           devices:          
           - name: "sdc"
           - name: "sdb"  ## "Add sdb"
         - name: "worker2"       
           directories:          
           - path: "/root/testdir"
           - path: "/root/testdir2" ## Add "/root/testdir2"
           devices:
           - name: "sdc"
        ```
### Ceph Cluster network 설정
- `spec.network.hostNetwork`: `true`로 설정할 경우 ceph cluster를 구성하는 pod들은 host network의 대역대에서 ip를 할당받고 `false`로 설정할 경우에는 pod들은 k8s cluster의 대역대에서 ip를 할당받습니다. 본 프로젝트에서는 `false`로 설정하는 것을 권장합니다.

## RunTime에서의 ceph cluster update하는 방법
- Ceph cluster의 상태가 `HEALTH_OK`일 경우에만 ceph cluster update를 진행하시기를 권장합니다.
### OSD를 추가하는 방법
- ceph cluster에 OSD daemon를 추가하고 싶은 경우, cluster.yaml의 `spec.storage. node[].device`에 osd에 관련한 설정을 추가하고 `kubectl apply -f cluster.yaml`를 진행하시면 됩니다. `spec.storage.node[].device`를 수정하는 방법은 OSD deploy setting section을 참고하시면 됩니다.
  - `kubectl apply -f cluster.yaml` 이후 operator pod에 의해 OSD 생성 작업이 진행되며, osd의 추가 확인은 toolbox pod에서 `ceph osd tree` 명령을 통해 확인할 수 있습니다. ([toolbox pod 접근 방법](docs/ceph-command.md) )
    ```shell
    $ ceph osd tree
    ID CLASS WEIGHT  TYPE NAME              STATUS REWEIGHT PRI-AFF
    -1       0.18359 root default                                   
    -2       0.18359     host ask-b360m-d3h                         
    0   ssd 0.09180         osd.0              up  1.00000 1.00000
    1   ssd 0.09180         osd.1              up  1.00000 1.00000  # osd 1 추가
    ```
### OSD를 제거하는 방법
- <strong>주의사항</strong>
  - OSD 제거 후에도 pool(blook, cephFS)의 replication 개수를 만족할만큼 Cluster에 OSD들이 존재하는지를 확인 후에 진행해주시기를 바랍니다.
  - node의 OSD 위치 변경을 위해 OSD 제거를 수행하실 경우, 위의  OSD 추가 방법을 통해 OSD 추가 후 해당 방법을 수행하기를 권장드립니다.
- 먼저 제거하고 싶은 OSD의 ID를 확인합니다. OSD의 ID는 OSD의 POD 이름에서 확인할 수 있습니다. (예시: OSD 0 제거)
  - ex) pod name: rook-ceph-osd-0-7d4c749468-4b9nj -> $ID: 0
- toolbox pod에서 다음과 같은 명령을 수행합니다.
    ```shell
    # osd 를 ceph cluster에서 out으로 마크하고, osd에 있던 데이터를 다른 osd로 옮기는 rebalancing 작업을 시작합니다.
    # ceph osd out <$ID>
    $ ceph osd out 0
    marked out osd.0.

    # osd에 있던 데이터가 다른 osd로 정상적으로 옮겨저서 osd를 제거해도 이상이 없을 때까지 기다립니다.
    # while ! ceph osd safe-to-destroy <$ID> ; do sleep 10 ; done
    $ while ! ceph osd safe-to-destroy 0 ; do sleep 10 ; done
    Error EBUSY: OSD(s) 0 have 24 pgs currently mapped to them.
    ...
    OSD(s) 0 are safe to destroy without reducing data durability.
    ```
- cluster.yaml에서 제거할 OSD에 해당하는 설정을 `spec.storage.node[].device`에서 지우신 후, `kubectl apply -f cluster.yaml`를 합니다.
- 제거할 OSD의 deployment를 삭제해줍니다.
    ```shell
    # kubectl delete deployment -n rook-ceph rook-ceph-osd-<$ID>
    $ kubectl delete deployment -n rook-ceph rook-ceph-osd-0
    ```
- 최종적으로 toolbox pod에서 다음 명령을 수행합니다.
    ```shell
    # cluster에서 osd를 최종적으로 제거
    # ceph osd purge <$ID> --yes-i-really-mean-it
    $ ceph osd purge 0 --yes-i-really-mean-it
    purged osd.0

    # ceph osd tree로 남아있는 osd 확인
    $ ceph osd tree
    ID CLASS WEIGHT  TYPE NAME              STATUS REWEIGHT PRI-AFF
    -1       0.09180 root default                                   
    -2       0.09180     host ask-b360m-d3h                         
    1   ssd 0.09180         osd.1              up  1.00000 1.00000
    ```
- OSD가 제거되었으므로, 제거한 OSD가 사용하던 디렉토리나 디바이스는 <strong>삭제하거나 초기화</strong>합니다.
- 삭제한 OSD가 `directories` 기반이고, 더 이상 해당 노드에 `directories` 기반으로 생성된 OSD가 <strong>없는</strong> 경우, 다음 configmap을 삭제해줍니다. (예시: worker1노드 configmap 삭제)
    ```shell
    # 노드별 directory 기반 osd 정보 저장하는 configmap 삭제
    # kubectl delete configmap -n rook-ceph rook-ceph-osd-<$hostname>-config
    $ kubectl delete configmap -n rook-ceph rook-ceph-osd-worker1-config
    ```
