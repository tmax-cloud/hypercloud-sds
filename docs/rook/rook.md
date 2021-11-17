# Rook Ceph Cluster

> Rook Ceph는 hypercloud-sds에서의 고가용성 storage 제공을 위해서 설치하는 모듈입니다. 본 프로젝트에서 cephFS(file system), rbd(block storage)를 제공합니다.

## Reference

* [Rook Ceph Storage official page](https://rook.github.io/docs/rook/v1.3/ceph-storage.html)

## Rook Ceph Cluster 설정하는 법

* 자세한 방법은 이 [문서](./ceph-cluster-setting.md)를 확인하세요.

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
