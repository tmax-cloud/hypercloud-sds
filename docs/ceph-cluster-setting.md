# cluster.yaml 수정 메뉴얼
본 메뉴얼에서는 cluster.yaml를 수정하는 방법에 대해서 설명합니다. 메뉴얼에서 설명하는 옵션에 대해서만 수정하는 것을 권장합니다.

## 하드웨어 요건 및 권장사항
클러스터 배포시 하드웨어 요건 및 권장사항은 다음과 같습니다.
- 하드웨어 요건
  - `OSD`
    - CPU: 2GHz CPU 2 core 이상 권장
    - Memory: 최소 2GB + 1TB당 1GB 이상의 메모리 권장
      - ex) 2TB disk 기반 OSD -> 4GB(2GB+2GB) memory 이상 권장
    - Disk: 최소 10GB 이상 디스크 권장, 데이터 저장용으로 큰 용량의 디바이스를 사용하는 게 좋음
  - `MON`
    - CPU: 2GHz CPU 1 core 이상 권장
    - Memory: 2GB 메모리 이상 권장
    - Disk: 10GB 디스크 이상 권장
  - `MGR`
    - CPU: 2GHz CPU 1 core 이상 권장
    - Memory: 1GB 메모리 이상 권장
  - `MDS`
    - CPU: 2GHz CPU 4 core 이상 권장
    - Memory: 4GB 메모리 이상 권장
  - <strong>production용으로 배포될 경우 해당 사항들을 반드시 지켜주셔야 ceph가 정상적인 성능 및 안정성을 제공합니다.</strong>
    - <strong>ceph pod들에 다음 값들로 resource를 설정하기 위해 아래의 resource 설정과 [cephFS의 resource 설정](/docs/file.md)</strong>을 참고바랍니다.
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
    image: ceph/ceph:v14.2.9
    allowUnsupported: false
  dataDirHostPath: /var/lib/rook
  skipUpgradeChecks: false
  continueUpgradeAfterChecksEvenIfNotHealthy: false
  mon:
    # set the amount of mons to be started, Recommendation: Use odd numbers (ex. 3, 5)
    count: 1
    allowMultiplePerNode: false
  mgr:
    modules:
    - name: pg_autoscaler
      enabled: true
  dashboard:
    enabled: true
    ssl: true
  monitoring:
    # requires Prometheus to be pre-installed for enabled is true
    enabled: false
    rulesNamespace: rook-ceph
  network:
    # enable host networking
    #provider: host
  rbdMirroring:
    workers: 0
  crashCollector:
    disable: false
  cleanupPolicy:
    deleteDataDirOnHosts: ""
  annotations:
  resources:
# set the requests and limits for osd, mon, mgr
#    osd:
#      limits:
#        cpu: "2"
#        memory: "4096Mi"
#      requests:
#        cpu: "2"
#        memory: "4096Mi"
#    mon:
#      limits:
#        cpu: "1"
#        memory: "2048Mi"
#      requests:
#        cpu: "1"
#        memory: "2048Mi"
#    mgr:
#      limits:
#        cpu: "1"
#        memory: "1024Mi"
#      requests:
#        cpu: "1"
#        memory: "1024Mi"
  removeOSDsIfOutAndSafeToRemove: false
  priorityClassNames:
    all: rook-ceph-default-priority-class
  disruptionManagement:
    managePodBudgets: false
    osdMaintenanceTimeout: 30
    manageMachineDisruptionBudgets: false
    machineDisruptionBudgetNamespace: openshift-machine-api
  storage:
    # set useAllNodes,useAllDevices to false for node-specific config
    useAllNodes: false
    useAllDevices: false
    config:
# Example for node-specific config. It works only when 'useAllNodes' is false.
    nodes:
      - name: "worker1"    # Add worker1 node to ceph-osd. (Caution: check hostname by 'kubectl get nodes')
        devices:           # Add disk of worker1 to ceph-osd.
        - name: "sdb"      # Caution: Disk must exist in worker1 node. (check disk by 'sudo fdisk -l')
        - name: "sdc"
#          config:          
#            metadataDevice: "sdd1"     # Separate metadata device to high-performance device. (ex. SSD)
#      - name: "worker2"         # Add worker2 node to ceph-osd.
#        devices:           # Add disk of worker1 to ceph-osd.
#      - name: "nvme01" # multiple osds can be created on high performance devices
#        config:
#          osdsPerDevice: "5"
```
### `Resource setting`
- `spec.resources`: cluster에 배포될 osd, mgr, mon에 대한 resource를 설정합니다.
  -  <strong>test환경이 아닌 production 환경일 경우 반드시 주석을 풀고 각 데몬에 대한 resource를 설정해주시기를 바랍니다.</strong>
  - `spec.resources.{dameon}.limits`와 `spec.resources.{dameon}.requests`에 설정되는 값들은 반드시 `동일`해야 합니다.
  - `spec.resources.osd`: osd pod에 대한 resource의 request, limit을 설정합니다. 해당 값은 각 osd pod에 모두 동일하게 적용됩니다. 그러므로 <strong>가장 큰 용량을 가진 osd</strong>를 기준으로 <strong>위의 하드웨어 요건</strong>에 따라 해당 값들을 설정해주시기를 바랍니다.
    - ex) 2TB disk 기반 OSD 배포, 배포되는 노드의 CPU 성능은 2GHz
      ```yaml
      osd:
        limits:
          cpu: "2"
          memory: "4096Mi" # 2GB+2GB
        Requests:
          cpu: "2"
          memory: "4096Mi"
      ```
    - ex) 3TB disk 기반 OSD 배포, 배포되는 노드의 CPU 성능은 2.5GHz
      ```yaml
      osd:
        limits:
          cpu: "1.6" # osd는 2GHz CPU core 2개가 필요하므로 CPU 성능이 2.5GHz일 경우는 ((2+2)/2.5)=1.6 이면 됩니다.
          memory: "5120Mi" # 2GB+3GB
        Requests:
          cpu: "1.6"
          memory: "5120Mi"
      ```
  - `spec.resources.mon`: mon pod에 대한 resource의 request, limit을 설정합니다. 위의 osd resource 설정 방식을 참고하여 작성하시면 됩니다.
  - `spec.resources.mgr`: mgr pod에 대한 resource의 request, limit을 설정합니다. 위의 osd resource 설정 방식을 참고하여 작성하시면 됩니다.
### Mon deploy setting
- `spec.mon.count`: kube cluster에 deploy할 mon의 개수를 의미합니다. `count`의 값은 1부터 9 사이의 홀수이어야 합니다. Ceph document에서는 기본적으로 3개의 mon를 권장합니다.
- `spec.mon.allowMultiplePerNode`: `true` 또는 `false`를 값으로 가질 수 있으며, 하나의 노드에 여러 개의 mon를 deploy할 수 있는지에 대한 여부를 결정하는 옵션입니다. `false`로 설정 할 경우, 하나의 node에 여러 개의 mon pod이 deploy될 수 없습니다.

### OSD deploy setting
- `spec.storage.useAllNodes`, `spec.storage.useAllDevices`: <strong> production 환경일 경우, 각 노드마다 올바르게 osd를 배포하기 위해 해당 값들을 `false`로 하길 바랍니다. </strong>
  - 이 두 값이 `true`일 경우, 모든 노드에서 사용가능한 모든 device에 osd 배포를 시도합니다. 
- `spec.storage.nodes`: 각 노드 별로 deploy되는 OSD pod의 설정을 다르게 하고 싶은 경우 해당 설정에 osd에 관한 설정을 명시하면 됩니다.
    - `name`: OSD pod가 deploy될 node 이름을 명시합니다. 해당 값은 `kubernetes.io/hostname`과 동일해야 됩니다.
    - `devices`: OSD pod를 device 위에 deploy하겠다는 옵션입니다. 해당 옵션에서는 `name`에 device 이름(device 파티션 이름도 가능)을 명시하면 됩니다. 명시되는 device는 아래의 조건을 만족해야 됩니다.(만족하지 않을 경우, OSD pod가 deploy되지 않습니다.)
        - <strong>lvm device의 경우 현재 지원되지 않습니다.</strong>
        - Deploy되는 노드에 해당 디바이스가 반드시 존재해야 하며 unmount된 상태여야 합니다.
        - 해당 device의 초기화가 제대로 되어 있어야 하며 초기화 제대로 안 되어있을 경우 아래의 명령어를 통해 초기화합니다.(당연히 초기화후 원래 있었던 데이터는 복구할 수 없습니다.)
            ```shell
            # /dev/sdb device에 osd 설치를 원하는 경우
            DISK="/dev/sdb"
            $ sudo sgdisk --zap-all $DISK
            $ sudo if=/dev/zero of="$DISK" bs=1M count=100 oflag=direct,dsync
            ```
    - 하나의 노드에 두 개 이상의 OSD pod를 deploy하고 싶은 경우 아래와 같이 osd 설정을 추가하면 됩니다.
        ```yaml
        nodes:
         - name: "worker1"   
           devices:          
           - name: "sdc"
           - name: "sdb"  ## "Add sdb"
        ```
    - <strong>OSD의 성능 향상을 위해 Write-Ahead Logging (WAL) 및 DB 디바이스를 성능이 좋은 디스크(ex. SSD, NVMe) 로 분리하는 방법</strong> 은 [cluster tuning 문서](/docs/cluster-tuning.md)를 참고하길 바랍니다.
    - <strong>NVMe SSD를 기반으로 OSD를 생성할 경우</strong> 하나의 OSD만으로는 NVMe SSD의 성능을 모두 활용하지 못합니다. 하나의 NVMe SSD에 여러 OSD를 생성하기 위해서는 아래와 같이 `config.osdsPerDevice`로 생성할 OSD 개수를 명시해주면 됩니다.  
      ```yaml
      nodes:
      - name: "worker2"
        devices:
        - name: "nvme01" # multiple osds can be created on high performance devices
          config:
            osdsPerDevice: "2" # nvme01 device에 OSD를 2개 생성합니다.
      ```
### Ceph Cluster network 설정
- `spec.network.provider`: 해당 주석을 풀고 `host`로 설정할 경우 ceph cluster를 구성하는 pod들은 host network의 대역대에서 ip를 할당받고, 주석을 풀지 않을 경우에는 pod들은 k8s cluster의 대역대에서 ip를 할당받습니다. 본 프로젝트에서는 주석을 풀지 않을 것을 권장합니다.

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
- OSD가 제거되었으므로, 제거한 OSD가 사용하던 디바이스는 <strong>초기화</strong>합니다.