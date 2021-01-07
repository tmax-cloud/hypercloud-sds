# RunTime에서의 ceph cluster update하는 방법

> Ceph cluster의 상태가 `HEALTH_OK`일 경우에만 ceph cluster update를 진행하시기를 권장합니다.

## OSD를 추가하는 방법
- ceph cluster에 OSD daemon를 추가하고 싶은 경우, cluster.yaml의 `spec.storage. node[].device`에 osd에 관련한 설정을 추가하고 `kubectl apply -f cluster.yaml`를 진행하시면 됩니다. `spec.storage.node[].device`를 수정하는 방법은 OSD deploy setting section을 참고하시면 됩니다.
  - `kubectl apply -f cluster.yaml` 이후 operator pod에 의해 OSD 생성 작업이 진행되며, osd의 추가 확인은 toolbox pod에서 `ceph osd tree` 명령을 통해 확인할 수 있습니다.
  - toolbox pod 접근 방법은 [Ceph 명령어 메뉴얼](ceph-command.md)를 참고해주세요.
    ```shell
    $ ceph osd tree
    ID CLASS WEIGHT  TYPE NAME              STATUS REWEIGHT PRI-AFF
    -1       0.18359 root default                                   
    -2       0.18359     host ask-b360m-d3h                         
    0   ssd 0.09180         osd.0              up  1.00000 1.00000
    1   ssd 0.09180         osd.1              up  1.00000 1.00000  # osd 1 추가
    ```

## OSD를 제거하는 방법

- **주의사항**
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
- OSD가 제거되었으므로, 제거한 OSD가 사용하던 디바이스는 **초기화**합니다. 초기화 방법은 이 [문서](rook.md)를 참고하세요.