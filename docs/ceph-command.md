## Ceph 명령어 메뉴얼

### Rook으로 구축한 ceph cluster의 경우, toolbox pod를 통해 ceph 명령어를 수행시킬 수 있습니다.

- <strong>Toolbox pod를 통해 ceph 명령어를 수행하는 방법</strong>
    ```shell
    # toolbox pod의 이름을 찾아냄.
    $ kubectl -n rook-ceph get pod -l "app=rook-ceph-tools"
    NAME                               READY   STATUS    RESTARTS   AGE
    rook-ceph-tools-8648fbb998-dwhs7   1/1     Running   0          112m

    # toolbox pod의 이름으로 toolbox pod에 접근함.
    # kubectl exec -it -n rook-ceph {$toolboxPodName} -- /bin/bash
    # 본인 환경의 toolbox pod의 이름을 사용해야 합니다.
    $ kubectl exec -it -n rook-ceph rook-ceph-tools-8648fbb998-dwhs7 -- /bin/bash

    # 접속한 toolbox pod에서 ceph 명령어를 입력하면 됩니다.
    $ {$cephCommand}
    ```

- <strong>Ceph cluster 전체의 상태 확인</strong>
    ```shell
    $ ceph -s
    cluster:
      id:     9c13e4ad-a4eb-4047-9c57-61149cace075
      health: HEALTH_OK

    services:
      mon: 1 daemons, quorum a (age 14m)
      mgr: a(active, since 13m)
      mds: myfs:1 {0=myfs-b=up:active} 1 up:standby-replay
      osd: 3 osds: 3 up (since 13m), 3 in (since 13m)

    data:
      pools:   3 pools, 24 pgs
      objects: 165 objects, 478 MiB
      usage:   3.9 GiB used, 29 GiB / 33 GiB avail
      pgs:     24 active+clean

    io:
      client:   852 B/s rd, 1 op/s rd, 0 op/s wr
    ```
- <strong>Ceph cluster의 disk 사용량 확인</strong>
    ```shell
    $ ceph df
    RAW STORAGE:
        CLASS     SIZE       AVAIL      USED        RAW USED     %RAW USED
        hdd       33 GiB     29 GiB     942 MiB      3.9 GiB         11.88
        TOTAL     33 GiB     29 GiB     942 MiB      3.9 GiB         11.88

    POOLS:
        POOL              ID     STORED      OBJECTS     USED        %USED     MAX AVAIL
        myfs-metadata      1      50 KiB          24       1 MiB         0        14 GiB
        myfs-data0         2      32 MiB           8      64 MiB      0.23        14 GiB
        replicapool        3     431 MiB         133     865 MiB      3.02        14 GiB
    ```
- <strong>RBD image의 사용량 확인</strong>
    ```shell
    # rbd du -p {$poolName}
    $ rbd du -p replicapool
    warning: fast-diff map is not enabled for csi-vol-4af65546-4717-11ea-96fa-ea47d19dc5d1. operation may be slow.
    warning: fast-diff map is not enabled for csi-vol-4f0cd143-4717-11ea-96fa-ea47d19dc5d1. operation may be slow.
    NAME                                         PROVISIONED USED
    csi-vol-4af65546-4717-11ea-96fa-ea47d19dc5d1       1 GiB 184 MiB
    csi-vol-4f0cd143-4717-11ea-96fa-ea47d19dc5d1       1 GiB  84 MiB
    <TOTAL>                                            2 GiB 268 MiB
    ```
- <strong>RBD image 제거하는 방법</strong>
    - `reclaimPolicy`를 `Retain`으로 설정할 경우, pv를 지워도 RBD image가 ceph cluster에 남게 됩니다. 이 경우에는 직접 `ceph command`를 통해 RBD image를 제거해야 됩니다.  
    ```shell
    # RBD image 이름 확인 방법
    # 해당 작업은 toolbox가 아닌 k8s 운영 환경에서 진행하셔야 합니다.
    # relased된 pv가 바라보고있는 volume handle에서 RBD image의 이름을 구합니다.
    # volume_name=$(kubectl get pv ${pv_name} -o custom-columns=name:.spec.csi.volumeHandle --no-headers); echo "csi-vol-${volume_name#*-*-*-*-*-}"
    $ volume_name=$(kubectl get pv pvc-5a239f29-1845-4d1e-81ff-af1a48bd31b6 -o custom-columns=name:.spec.csi.volumeHandle --no-headers); echo "csi-vol-${volume_name#*-*-*-*-*-}"
    csi-vol-f66401ce-4887-11ea-b951-aada1f93f6fa

    # RBD image list 확인
    # rbd ls {$poolName}
    $ rbd ls replicapool
    csi-vol-e2ce2428-4726-11ea-96fa-ea47d19dc5d1
    csi-vol-e645b9f9-4726-11ea-96fa-ea47d19dc5d1

    # RBD image 제거
    # rbd rm -p {$poolName} {RBDImageName}
    $ rbd rm -p replicapool csi-vol-e2ce2428-4726-11ea-96fa-ea47d19dc5d1
    Removing image: 100% complete...done.

    # RBD image 삭제 확인
    $ rbd ls replicapool
    csi-vol-e645b9f9-4726-11ea-96fa-ea47d19dc5d1
    ```
- <strong> CephFS volume 제거하는 방법</strong>
    - `reclaimPolicy`를 `Retain`으로 설정할 경우, pv를 지워도 cephfs volume이 ceph cluster에 남게 됩니다. 이 경우에는 직접 `ceph command`를 통해 cephFS volume를 제거해야 됩니다.  
    ```shell
    # cephFS volume 이름 확인 방법
    # 해당 작업은 toolbox가 아닌 k8s 운영 환경에서 진행하셔야 합니다.
    # relased된 pv가 바라보고있는 volume handle에서 cephfs volume의 이름을 구합니다.
    # volume_name=$(kubectl get pv ${pv_name} -o custom-columns=name:.spec.csi.volumeHandle --no-headers); echo "csi-vol-${volume_name#*-*-*-*-*-}"
    $ volume_name=$(kubectl get pv pvc-20d77747-e9c8-4338-8ebe-eb04239dbff3 -o custom-columns=name:.spec.csi.volumeHandle --no-headers); echo "csi-vol-${volume_name#*-*-*-*-*-}"
    csi-vol-2f2b9560-48a2-11ea-9a8e-ca0f135cc4c1

    # cephFS volume 삭제
    # ceph fs subvolume rm ${cephFsName} --group_name csi ${csi_volume_name}
    $ ceph fs subvolume rm myfs --group_name csi csi-vol-2f2b9560-48a2-11ea-9a8e-ca0f135cc4c1
    ```
