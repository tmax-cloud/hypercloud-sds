## Rook Ceph Version Upgrade Guide

- 본 문서는 Rook Ceph verion upgrade를 위한 문서입니다.
    - Rook Ceph `v1.2.7`를 제공하고 있는 hypercloud-rook-ceph (tag 1.1.1 ~ 1.1.7)에서 Rook Ceph `v1.3.6`를 제공하는 Hypercloud-SDS release 1.3로 upgrade하는 것을 목표로 합니다.
    - Upgrade되는 Rook Ceph version을 기준으로 Ceph CSI와 Ceph version에 대해서 upgrade를 진행합니다.

### 주의사항

- Rook Ceph 버전 upgrade를 수행할 경우, 예상치 못한 issue 가 발생할 수도 있으므로, 중요한 data(Ceph volume)는 backup 후 upgrade를 수행하는 것을 권장합니다.
- Rook Ceph 버전 upgrade하는 동안에는, 일시적으로 Ceph volume에 대한 정상적인 사용이 불가능할 수 있습니다.
- 배포되어 있는 Ceph OSD daemon이 `device 기반`이 아닌 `directory` 기반일 경우 다음 Rook Ceph version upgrade를 진행할 수 없습니다.
- 각 단계에서 인지가 필요한 내용을 주석 및 글 형식으로 기재 하였으므로, 꼭 읽으시는 것을 권장합니다.
- <strong> 반드시 각 단계가 완료된 후, 다음 단계를 진행해주시길 바랍니다.</strong>
    - 각 단계가 완료되기 전에 다음 단계를 진행할 경우 추후 복구가 어려울 수 있습니다.


### Rook Ceph Upgrade

1. Rook Ceph cluster 상태 확인
    - Upgrade 수행 전, Rook Ceph cluster 상태를 확인합니다.
        - Rook Ceph 관련 pod이 모두 정상적으로 수행되고 있는지 확인합니다.
        - Ceph command를 통해 Ceph cluster의 상태를 확인합니다.
            - Ceph cluster의 heath status가 HEALTH_OK인지 확인합니다.
            - Ceph mon들이 모두 quorum에 포함되어 있는지 확인합니다.
            - Ceph mgr이 active 상태인지 확인합니다.
            - 모든 ceph osd들이 up & in 상태인지 확인합니다.
            - 모든 pg가 active + clean 상태인지 확인합니다.
            - Ceph mds가 active 상태인지 확인합니다.
    
    ```shell
    # Rook Ceph cluster를 구성하고 있는 pod 들이 모두 정상적으로 RUNNING 하고 있는지 확인합니다.
    $ kubectl -n rook-ceph get pods -o wide

    # Ceph command인 "ceph status"를 사용하여, ceph cluster의 상태를 확인합니다.
    # Rook Ceph toolbox pod를 통해서 ceph command를 입력할 수 있습니다.
    # TOOLS_POD=$(kubectl -n rook-ceph get pod -l "app=rook-ceph-tools" -o jsonpath='{.items[0].metadata.name}')
    # kubectl -n rook-ceph exec -it $TOOLS_POD -- ceph status
    # Rook Ceph toolbox pod의 이름은 deploy 될 때마다 다르므로, 해당 환경에서의 Rook Ceph toolbox의 이름을 꼭 확인해주세요.
    # Rook Ceph cluster의 health status가 HEALTH_OK이어야 합니다.
    $ kubectl -n rook-ceph exec -it rook-ceph-tools-8648fbb998-7lkmp -- ceph status
    cluster:
      id:     aede3c9c-c390-4bb2-9d12-7a569aee1808
      health: HEALTH_OK
  
    services:
      mon: 1 daemons, quorum a (age 4m)
      mgr: a(active, since 3m)
      mds: myfs:1 {0=myfs-b=up:active} 1 up:standby-replay
      osd: 3 osds: 3 up (since 2m), 3 in (since 2m)
  
    data:
      pools:   3 pools, 24 pgs
      objects: 22 objects, 2.2 KiB
      usage:   3.0 GiB used, 33 GiB / 36 GiB avail
      pgs:     24 active+clean

    io:
      client:   1.2 KiB/s rd, 2 op/s rd, 0 op/s wr
    ```

----------

2. Rook Ceph Operator & Ceph CSI Version Upgrade
    - STEP A) Ceph 관련 pod의 우선 순위를 올리기 위해서 priority class를 추가적으로 설정합니다.
        - 설정된 Priority class는 추후 과정에서 Ceph CSI pod과 Ceph daemon pod에 적용됩니다.
        ```shell
        # Create Rook Cpeh priority class
        $ kubectl apply -f 0_priority.yaml

    - STEP B) 변경된 RBAC, CRD에 대해서 update를 수행합니다.
        - `kubectl apply` 를 적용 시, 다음과 같은 Warning 메시지가 나올 수 있으나, Error 메시지가 출력되지 않고 원하는 resource 가 `configured` 혹은 `created` 되었다는 메시지가 출력되었다면 Warning 메시지는 무시하셔도 됩니다. : `Warning: kubectl apply should be used on resource created by either kubectl create --save-config or kubectl apply`
        ```shell
        # Update RBAC, CRD
        $ kubectl apply -f 1_upgrade-from-v1.2-apply.yaml -f 2_upgrade-from-v1.2-crds.yaml
        ```

    - STEP C) Rook Ceph operator 설정을 update합니다.
        ```shell
        # Apply 되어있는 Rook Ceph operator의 설정이 hypercloud-rook-ceph (tag 1.1.1 ~ tag 1.1.7)의 operator.yaml과 동일한 상황에서 적용해야 합니다.
        $ kubectl apply -f 3_operator_1.3.6_patch.yaml
        ```

    - STEP D) Rook Ceph version upgrade가 완료될 때까지 기다립니다.
        ```shell
        # 다음 커맨드를 입력하여 보이는 모든 rook-version 의 value 가  v1.3.6 으로 출력될 때까지 기다립니다.
        # (총 소요 시간은 환경에 따라 다르지만 대략 10분 이상 걸립니다.)
        $ watch --exec kubectl -n rook-ceph get deployments -l rook_cluster=rook-ceph -o jsonpath='{range .items[*]}{.metadata.name}{"  \treq/upd/avl: "}{.spec.replicas}{"/"}{.status.updatedReplicas}{"/"}{.status.readyReplicas}{"  \trook-version="}{.metadata.labels.rook-version}{"\n"}{end}'
        Every 2.0s: kubectl -n rook-ceph get deployments -l rook_cluster=rook-ceph -o jsonpath={range .items[*]}{.metadata.name}{"  \treq/upd/avl: "}{.spec.repl...  master1: Wed Nov 25 04:52:57 2020

        rook-ceph-crashcollector-master1        req/upd/avl: 1/1/1      rook-version=v1.3.6              ## 모든 pod에 대해서 req/upd/avl가 1/1/1, 
        rook-ceph-crashcollector-worker1        req/upd/avl: 1/1/1      rook-version=v1.3.6              ## rook-version=1.3.6이어야 합니다
        rook-ceph-crashcollector-worker2        req/upd/avl: 1/1/1      rook-version=v1.3.6
        rook-ceph-mds-myfs-a    req/upd/avl: 1/1/1      rook-version=v1.3.6
        rook-ceph-mds-myfs-b    req/upd/avl: 1/1/1      rook-version=v1.2.7                              ## rook-verion=v1.3.6가 되어야 합니다
        rook-ceph-mgr-a         req/upd/avl: 1/1/1      rook-version=v1.3.6                              
        rook-ceph-mon-a         req/upd/avl: 1/1/1      rook-version=v1.3.6
        rook-ceph-osd-0         req/upd/avl: 1/1/1      rook-version=v1.2.7                              ## rook-version=v1.3.6가 되어야 합니다
        rook-ceph-osd-1         req/upd/avl: 1/1/1      rook-version=v1.2.7                              ## rook-version=v1.3.6가 되어야 합니다
        rook-ceph-osd-2         req/upd/avl: 1/1/1      rook-version=v1.2.7                              ## rook-version=v1.3.6가 되어야 합니다

        # 위의 stdout에서 rook-version 이 모두 v1.3.6으로 출력되면, 최종적으로 다음 커맨드를 입력하여 하나의 값만 출력되는지 확인합니다.
        # (정상 upgrade 가 완료되지 않을 경우에는, v1.2.7 과 v1.3.6 두 개의 row 가 출력됩니다.)
        $ kubectl -n rook-ceph get deployment -l rook_cluster=rook-ceph -o jsonpath='{range .items[*]}{"rook-version="}{.metadata.labels.rook-version}{"\n"}{end}' | sort | uniq
        rook-version=v1.3.6
        ```

    - STEP E) Ceph CSI version을 확인합니다.
        - Rook Ceph operator 설정을 update 할 경우, Ceph CSI version 또한 upgrade 됩니다.
        ```shell
        $ kubectl --namespace rook-ceph get pod -o jsonpath='{range .items[*]}{range .spec.containers[*]}{.image}{"\n"}' -l 'app in (csi-rbdplugin,csi-rbdplugin-provisioner,csi-cephfsplugin,csi-cephfsplugin-provisioner)' | sort | uniq
        
        quay.io/cephcsi/cephcsi:v2.1.2
        quay.io/k8scsi/csi-attacher:v2.1.0
        quay.io/k8scsi/csi-node-driver-registrar:v1.2.0
        quay.io/k8scsi/csi-provisioner:v1.4.0
        quay.io/k8scsi/csi-resizer:v0.4.0
        quay.io/k8scsi/csi-snapshotter:v1.2.2
        ```

----------

3. Ceph version upgrade & Config Ceph daemon pod priority class
    > 본 단계에서는 deploy되어 있는 Ceph cluster 구성 시, 사용된 cluster.yaml를 수정하여 적용하는 방법으로 upgrade가 수행됩니다.
    - STEP A) Ceph image name를 수정합니다.
        - Ceph image name를 ceph/ceph:v14.2.9로 변경합니다.
        ```yaml
        ## 반드시 현재 Ceph cluster를 구성한 cluster.yaml를 수정해야합니다.
        ## 수정전 yaml를 backup하시는 것을 권장합니다.
        ... 
        ## cluster.yaml
        spec:
           cephVersion:
              # image: 192.168.50.100:5000/ceph/ceph:v14.2.8 
              image: ceph/ceph:v14.2.9
              allowUnsupported: true
        dataDirHostPath: /var/lib/rook
        ...
        ```
    - STEP B) Priority class 설정을 추가합니다.
        ```yaml
        ## 반드시 현재 Ceph cluster를 구성한 cluster.yaml를 수정해야합니다.
        ## 수정전 yaml를 backup하시는 것을 권장합니다.
        ...
        ## cluster.yaml
        ## Add priorityClassName
        priorityClassNames:
          all: rook-ceph-default-priority-class
 
        storage:
          useAllNodes: false
          useAllDevices: false
        ...
        ```
    - STEP C) 변경된 cluster.yaml를 반영합니다.
        ```shell
        ## example_cluster.yaml은 STEP A)와 STEP B)를 반영한 예시 yaml입니다. 
        $ kubectl apply -f cluster.yaml
        ```
    - STEP D) Ceph daemon 들의 upgrade 진행 상태를 확인합니다.
      - 보이는 모든 ceph-version 의 value 가  14.2.9-0 으로 출력될 때까지 기다립니다.
      - (mon -> mgr -> osd -> mds 순서로 버전 upgrade 가 적용되며, 총 소요 시간은 환경에 따라 다르지만 대략 10분 이상 걸립니다.)
        ```shell
        $ watch --exec kubectl -n rook-ceph get deployments -l rook_cluster=rook-ceph -o jsonpath='{range .items[*]}{.metadata.name}{"  \treq/upd/avl: "}{.spec.replicas}{"/"}{.status.updatedReplicas}{"/"}{.status.readyReplicas}{"  \tceph-version="}{.metadata.labels.ceph-version}{"\n"}{end}'
        Every 2.0s: kubectl -n rook-ceph get deployments -l rook_cluster=rook-ceph -o jsonpath={range .items[*]}{.metadata.name}{"  \treq/upd/avl: "}{.spec.repl...  master1: Wed Nov 25 05:06:59 2020

        rook-ceph-crashcollector-master1        req/upd/avl: 1/1/1      ceph-version=14.2.9-0              ## 모든 pod에 대해서 req/upd/avl이 1/1/1
        rook-ceph-crashcollector-worker1        req/upd/avl: 1/1/1      ceph-version=14.2.9-0              ## ceph-version=14.2.9-0이어야 합니다
        rook-ceph-crashcollector-worker2        req/upd/avl: 1/1/1      ceph-version=14.2.9-0
        rook-ceph-mds-myfs-a    req/upd/avl: 1/1/1      ceph-version=14.2.9-0
        rook-ceph-mds-myfs-b    req/upd/avl: 1/1/1      ceph-version=14.2.8-0                              ## ceph-version=14.2.9-0가 되어야 합니다
        rook-ceph-mgr-a         req/upd/avl: 1/1/1      ceph-version=14.2.9-0
        rook-ceph-mon-a         req/upd/avl: 1/1/1      ceph-version=14.2.9-0
        rook-ceph-osd-0         req/upd/avl: 1/1/1      ceph-version=14.2.8-0                              ## ceph-version=14.2.9-0가 되어야 합니다
        rook-ceph-osd-1         req/upd/avl: 1/1/1      ceph-version=14.2.8-0                              ## ceph-version=14.2.9-0가 되어야 합니다
        rook-ceph-osd-2         req/upd/avl: 1/1/1      ceph-version=14.2.9-0

        # 위의 stdout 에서 ceph-version 이 모두 14.2.9-0 로 출력되면, 최종적으로 다음 커맨드를 입력하여 하나의 값만 출력되는지 확인합니다.
        # (정상 upgrade 가 완료되지 않으면 14.2.8-0 과 14.2.9-0 두 개의 row 가 출력됩니다.)
        $ kubectl -n rook-ceph get deployment -l rook_cluster=rook-ceph -o jsonpath='{range .items[*]}{"ceph-version="}{.metadata.labels.ceph-version}{"\n"}{end}' | sort | uniq
        ceph-version=14.2.9-0
        ```
----------
    
4. Rook Ceph toolbox upgrade
    - rook/ceph:v1.2.7으로 deploy 된 Rook Ceph toolbox를 rook/ceph:v1.3.6으로 재배포합니다.
        - 해당 Rook Ceph toolbox가 rook/ceph:v1.2.7을 사용하는지 확인하는 것을 권장 드립니다.
    ```shell
    # Rook Ceph toolbox deployment를 확인합니다.
    $ kubectl get deployment -n rook-ceph
    csi-cephfsplugin-provisioner       2/2     2            2           35m
                       ...
    rook-ceph-tools                    1/1     1            1           34m
                       ...
    rook-ceph-osd-2                    1/1     1            1           34m
    
    # 기존의 Rook Ceph toolbox deployment를 삭제합니다.
    $ kubectl delete deployment rook-ceph-tools -n rook-ceph
    deployment.extensions "rook-ceph-tools" deleted

    # Rook Ceph toolbox를 재배포 합니다.
    $ kubectl apply -f 4_toolbox.yaml
    
    # Rook Ceph toolbox의 생성 및 동작을 확인합니다.
    $ kubectl get pods -n rook-ceph | grep ceph-tools
    rook-ceph-tools-67788f4dd7-qcwnw                    1/1     Running     0          37s
    $ kubectl -n rook-ceph exec -it rook-ceph-tools-67788f4dd7-qcwnw -- ceph status 
    ```
----------

- Rook Ceph upgrade 완료 후, image version 정보
    - rook/ceph:v1.3.6
    - ceph/ceph:v14.2.9
    - quay.io/cephcsi/cephcsi:v2.1.2
    - quay.io/k8scsi/csi-node-driver-registrar:v1.2.0
    - quay.io/k8scsi/csi-resizer:v0.4.0
    - quay.io/k8scsi/csi-provisioner:v1.4.0
    - quay.io/k8scsi/csi-snapshotter:v1.2.2
    - quay.io/k8scsi/csi-attacher:v2.1.0

### Official Reference
- See the upgrade guide: https://rook.io/docs/rook/v1.3/ceph-upgrade.html