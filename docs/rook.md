# Rook Ceph Cluster

> Rook Ceph는 hypercloud-storage에서의 고가용성 storage 제공을 위해서 설치하는 모듈입니다. 본 프로젝트에서 cephFS(file system), rbd(block storage)를 제공합니다.

## Reference

* [Rook Ceph Storage official page](https://rook.github.io/docs/rook/v1.2/ceph-storage.html)

## Rook Ceph Cluster 설정하는 법
> 본 프로젝트에서는 helm을 사용하여 Rook Ceph Cluster를 설치하며, config.yaml를 수정하는 방식을 통해 설치되는 Ceph Cluster의 사양을 설정할 수 있습니다.
> config.yaml에서 rook-ceph-init, rook-ceph-core 부분이 설치되는 Ceph Cluster의 사양을 기입하는 부분입니다. 따로 수정을 하지 않을 경우, default로 설정되어 있는 값으로 Ceph Cluster가 설치됩니다.
```yaml
### 해당 yaml 파일은 /install/config.yaml에 대한 예제 파일입니다.
rook-ceph-init:
  image:
    repository: rook/ceph
    tag: v1.2.4
rook-ceph-core:
  filestorage:
    metaReplicas: 3
    dataReplicas: 3
    reclaimPolicy: Delete
  blockstorage:
    replicas: 3
    reclaimPolicy: Delete
  cephSpec:
    cephVersion:
      image:  ceph/ceph:v14.2.6
    mon:
      count: 3
    storage:
      useAllNodes: false          ## 각 노드에 deploy되는 osd의 설정을 동일하게 하고 싶은 경우 true로 설정하면 됩니다.
      useAllDevices: false
    # directories:                ## useAllNodes가 true일 경우, directories 필드를 통해 osd가 deploy할 dir를 설정할 수 있습니다.
    # - path: /var/lib/rook       ## 해당 예제에서는 useAllNodes가 false이기 때문에 주석처리 하였습니다.
      nodes:                      ## useAllNodes가 false일 경우, nodes를 통해 각 node에 deploy될 osd를 설정할 수 있습니다. 
      - name: "worker1"           ## worker1 노드에 ceph-osd를 배포합니다. (주의: 'kubectl get nodes'의 hostname과 동일해야합니다.)
        devices:                  ## ceph-osd를 배포할 device를 명시합니다.
        - name: "sdc"             ## 주의: deivce가 노드에 존재해야합니다. ('sudo fdisk -l'로 device 확인)
      - name: "worker2"           ## worker2 노드에 ceph-osd를 배포합니다.
        directories:              ## ceph-osd를 배포할 directory를 명시합니다.
        - path: "/root/testdir00" ## 주의: directory가 노드에 존재해야합니다.
        - path: "/root/testdir01" ## 주의: directory가 노드에 존재해야합니다.
        devices:                  ## ceph-osd를 배포할 device를 명시합니다.
        - name: "sdc"             ## 주의: device가 노드에 존재해야합니다.


### cdi 부분은 생략합니다.
```
### rook-ceph-init
> rook-ceph-init 부분은 Rook Ceph의 operator 관련 사양을 정의하는 부분입니다.
- `image`: image 부분은 설치할 Rook Ceph operator의 image repository와 tag를 입력하는 부분입니다.


### rook-ceph-core
> rook-ceph-core 부분은 Rook Ceph의 Cluster, cephFS storageClass, rbd storageClass 관련 사양을 정의하는 부분입니다.
- `filestorage` (cephFS 설정) 설명은 아래와 같습니다.
  - `metaReplicas`:  cephFS의 메타데이터가 저장되는 pool의 replica 개수
  - `dataReplicas`:  cephFS의 데이터가 저장되는 pool의 replica 개수
  - `reclaimPolicy`:  cephFS storageClass의 reclaim policy (ex. Delete, Retain 등)
- `blockstorage` (rbd 설정) 설명은 아래와 같습니다.
  - `replicas`: rbd pool의 replica 개수
  - `reclaimPolicy`: rbd storageClass의 reclaim policy (ex. Delete, Retain 등)
  
> replica 개수는 대체적으로 3을 권장하며 ceph의 성능을 위해서 2로 설정하는 경우도 있습니다.
> replica 개수가 n개일 때 <strong>적어도 n개 이상의 노드에 osd pod이 존재</strong>해야됩니다.

- `cephSpec` (Ceph Cluster 설정) 설명은 아래와 같습니다.
  - `cephVersion`: Ceph Cluster의 버전을 명시하는 부분입니다.
    - `image`: Ceph Cluster의 이미지 버전을 명시합니다.
  - `mon`: Ceph Cluster의 Monitor 설정입니다.
    - `count`: k8s Cluster에 deploy할 Monitor의 개수를 의미합니다. `count`의 값은 1부터 9 사이의 홀수이어야 합니다. Ceph document에서는 기본적으로 3개의 Monitor를 권장합니다.
  - `storage`: Ceph Cluster의 osd 설정입니다.
    - `useAllNodes`: pod을 deploy할 수 있는 모든 노드에 동일한 조건의 osd pod을 deploy하고 싶으면 `true`로, 각 노드별로 다른 설정의 osd pod을 deploy하고 싶으면 `false`로 설정
    - `useAllDevices`: 모든 device에 대해서 osd pod을 deploy하고 싶으면 `true`로, 그렇지 않으면 `false`로 설정
    - `directories`: `useAllNodes`가 `true`일 경우, 모든 노드에서 osd가 배포될 directory를 설정
    - `nodes`: 각 노드 별로 deploy되는 osd pod의 설정을 다르게 하고 싶은 경우 해당 설정에 osd에 관한 설정을 명시하면 됩니다.
      > useAllNodes를 false로 설정했을 경우에만 nodes 필드를 사용하셔야 합니다.
      - `name`: osd pod이 deploy될 node 이름을 명시합니다. 해당 값은 `kubernetes.io/hostname`과 동일해야 됩니다.
      - `devices`: osd pod을 device 위에 deploy하겠다는 옵션입니다. 해당 옵션에서는 `name`에 device 이름을 명시하면 됩니다. 명시되는 device는 아래의 조건을 만족해야 됩니다. (만족하지 않을 경우, osd pod이 deploy되지 않습니다)
        - deploy되는 노드에 해당 디바이스가 반드시 존재해야 하며 unmount된 상태여야 합니다.
        - 해당 device의 초기화가 제대로 되어 있어야 하며 초기화 제대로 안 되어있을 경우 아래의 명령어를 통해 초기화합니다. (초기화 후, 원래 있었던 데이터는 복구할 수 없습니다)
            ```shell
            # /dev/sdb device에 osd 설치를 원하는 경우
            $ sudo sgdisk --zap-all /dev/sdb
            ```
      - `directories`: osd pod을 directory 위에 deploy하겠다는 옵션입니다. 해당 옵션에서는 `path`에 directory 이름을 명시하면 됩니다. 명시되는 directory는 아래의 조건을 만족해야 됩니다.(만족하지 않을 경우, osd pod이 deploy되지 않습니다)
        - deploy되는 노드에 해당 디렉토리가 반드시 존재해야 됩니다.
        - osd pod을 directory 위에 deploy할 경우, storage 사용에 대한 monitoring이 정상적으로 수행되지 않는 현상을 발견했습니다. 그러므로 <strong>osd pod을 device위에 deploy</strong>하는 것을 권장합니다.

## Rook-Ceph Cluster 제거
> Uninstall 작업 후, 초기화를 위해 다음 작업들을 수행해야합니다.
- k8s Cluster의 <strong>모든 노드</strong>에서 `/var/lib/rook` directory를 삭제합니다.

  ```shell
    $ rm -rf /var/lib/rook
  ```
- <strong>모든 osd</strong>의 backend directory 혹은 device를 삭제합니다.
  * osd의 backend가 directory인 경우 (backend directory 경로 예시: `/mnt/cephdir`)
  
    ```shell
      $ rm -rf /mnt/cephdir
    ```
  * osd의 backend가 device인 경우 (backend device 예시: `sdb`)
  
    ```shell
      # device의 파티션 정보 제거
      $ sgdisk --zap-all /dev/sdb
      
      # device mapper에 남아있는 ceph-volume 정보 제거 (각 노드당 한 번씩만 수행하면 됨)
      $ ls /dev/mapper/ceph-* | xargs -I% -- dmsetup remove %
      
      # /dev에 남아있는 찌꺼기 파일 제거
      $ rm -rf /dev/ceph-*
    ```