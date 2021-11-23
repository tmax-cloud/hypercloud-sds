# 스토리지 종류에 따른 rbd 볼륨 생성 방법

> 해당 문서는 rbd에서 storage 종류 (hdd/ssd/nvme)에 따라 서로 다른 pool에 볼륨 생성을 하는 방법을 설명합니다.


## HOWTO

1. 구분하고 싶은 storage 종류만큼 rbd pool을 생성합니다.

- 스토리지 종류 개수만큼 CephBlockPool CR을 생성합니다.
  - 아래 예시처럼 hdd/ssd/nvme에 대해서 각각의 CR을 생성할 수 있습니다.

    ```yaml
    apiVersion: ceph.rook.io/v1
    kind: CephBlockPool
    metadata:
      name: replicapool-hdd   ## hdd를 위한 rbd pool
      namespace: rook-ceph
    spec:
      failureDomain: host
      replicated:
        size: 2
    ``` 
    
- 위와 같이 스토리지 별로 CephBlockPool CR을 생성하면 여러 rbd pool이 생성되는 것을 확인할 수 있습니다.

    ```shell
    ## toolbox에서 해당 command 실행
    $ ceph osd pool ls
    
    ...  ## 다른 pool들은 생략
    
    replicapool-hdd
    replicapool-ssd
    replicapool-nvme
    ```
    
2. toolbox에서 여러 pool에 대하여 다른 스토리지를 사용하도록 CRUSH rule을 설정합니다.
 
- 스토리지 종류에 따른 CRUSH rule을 생성합니다.
  - Command: `ceph osd crush rule create-replicated <rule-name> <root> <failure-domain> <class>`

    ```shell
    ## toolbox에서 해당 command 실행
    $ ceph osd crush rule create-replicated hdd_rule default host hdd     ## hdd 스토리지에 대한 CRUSH rule
    $ ceph osd crush rule create-replicated ssd_rule default host ssd     ## ssd 스토리지에 대한 CRUSH rule
    $ ceph osd crush rule create-replicated nvme_rule default host nvme   ## nvme 스토리지에 대한 CRUSH rule
    ```

- CRUSH rule을 pool에 적용합니다.
  - Command: `ceph osd pool set <pool-name> crush_rule <rule-name>`
  - 참고: metadata pool에도 다음 방법으로 CRUSH rule을 적용할 수 있습니다.

    ```shell
    ## toolbox에서 해당 command 실행
    $ ceph osd pool set replicapool-hdd crush_rule hdd_rule    ## replicapool-hdd pool에 hdd에 대한 CRUSH rule을 적용
    set pool 8 crush_rule to hdd_rule
    $ ceph osd pool set replicapool-ssd crush_rule ssd_rule    ## replicapool-ssd pool에 ssd에 대한 CRUSH rule을 적용
    set pool 9 crush_rule to ssd_rule
    $ ceph osd pool set replicapool-nvme crush_rule nvme_rule   ## replicapool-nvme pool에 nvme에 대한 CRUSH rule을 적용
    set pool 10 crush_rule to nvme_rule
    ```
  
3. 스토리지 종류에 따른 storage class를 생성합니다. 

- storage class의 파라미터 중 pool 필드를 스토리지 종류에 맞는 pool로 명시하여 적용합니다.

    ```yaml
    apiVersion: storage.k8s.io/v1
    kind: StorageClass
    metadata:
       name: ceph-block-hdd-sc  ## hdd를 사용할 storage class
    provisioner: rook-ceph.rbd.csi.ceph.com
    parameters:
        clusterID: rook-ceph
    
        pool: replicapool-hdd   ## replicapool-hdd는 hdd_rule이 적용된 pool
    
        imageFormat: "2"
    
        imageFeatures: layering
    
        csi.storage.k8s.io/provisioner-secret-name: rook-csi-rbd-provisioner
        csi.storage.k8s.io/provisioner-secret-namespace: rook-ceph
        csi.storage.k8s.io/node-stage-secret-name: rook-csi-rbd-node
        csi.storage.k8s.io/node-stage-secret-namespace: rook-ceph
        csi.storage.k8s.io/fstype: ext4
    reclaimPolicy: Delete
    ```
    
- 스토리지 종류들에 따라 위 예시대로 storage class들을 생성합니다.

- 사용자는 pvc를 만들 때 원하는 storage class를 선택해서 스토리지 종류대로 rbd를 이용 가능합니다.



## 작업이 제대로 수행되었는지 확인하기

> Storage 종류에 따라 pool 이 잘 생성되었는지는 Ceph Placement Group (pg) 조회를 통해 확인할 수 있습니다.


1. 어떤 osd들이 어떤 Storage 종류에 속하는지를 파악합니다.
- Command: `ceph osd df | awk '{print $1, $2}'`

    ```shell
    ## toolbox에서 해당 command 실행
    $ ceph osd df | awk '{print $1, $2}'
    
    ID CLASS
    2 hdd   ## osd 2는 HDD에 속함
    5 ssd   ## osd 5는 SSD에 속함
    4 hdd   ## osd 4는 HDD에 속함
    1 ssd   ## osd 1은 SSD에 속함
    3 hdd   ## osd 3은 HDD에 속함
    0 ssd   ## osd 0은 SSD에 속함
    ...
    ```

2. pool의 pg들이 Storage 종류에 속하는 osd들을 acting set으로 가지고 있는지 확인합니다.
- Command: `ceph pg dump | awk '{print $1, $19}'`

    ```shell
    ## toolbox에서 해당 command 실행
    $ ceph pg dump | awk '{print $1, $19}'
    
    ## 앞서 확인했듯이 HDD osd는 2,3,4, SSD osd는 0,1,5
    ...
    PG_STAT ACTING_PRIMARY
    9.16 [1,5]   ## pool 9의 pg 0x16은 acting set으로 osd 1,5를 가짐 (SSD osd)
    8.15 [4,3]   ## pool 8의 pg 0x15는 acting set으로 osd 4,3을 가짐 (HDD osd)
    9.17 [0,5]   ## pool 9의 pg 0x17은 acting set으로 osd 0,5를 가짐 (SSD osd)
    8.13 [2,4]   ## pool 8의 pg 0x13은 acting set으로 osd 2,4를 가짐 (HDD osd)
    9.12 [0,1]   ## pool 9의 pg 0x12은 acting set으로 osd 0,1를 가짐 (SSD osd)
    8.11 [3,2]   ## pool 8의 pg 0x11은 acting set으로 osd 3,2를 가짐 (HDD osd)
    ...
    ```
    - 위와 같이 각 pool의 모든 pg에 대해서 Storage 종류에 따른 osd들로 mapping 됐는지 확인하면 됩니다.

