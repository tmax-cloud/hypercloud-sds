## Cluster Tuning: 메타데이터 디바이스 분리

해당 문서는 Cluster에서 OSD의 성능 향상을 위해 Write-Ahead Logging (WAL) 및 DB 디바이스를 성능이 좋은 디스크(ex. SSD) 로 분리하는 방법을 다룹니다.

### Prerequisites
- OSD 별로 데이터 저장을 위한 하나의 HDD 와 메타데이터 저장을 위한 하나의 SSD 가 필요합니다.
  - SSD 는 전체 디바이스를 사용할 수도 있고, 파티션을 나눠서 사용할 수도 있습니다.
- 본 문서는 <strong>하나의 노드에 두 개의 HDD 와 하나의 SSD 가 있는 예제 환경</strong>을 기반으로 설명합니다.
  - 한 노드에 두 개의 OSD 를 deploy 하며, 각 OSD는 하나의 HDD 디바이스와 하나의 SSD 파티션을 사용합니다.

### 방법

#### 1. 파티션 생성하기
SSD 같은 경우 HDD보다 디바이스 수가 적을 수 있으므로, 실질적인 성능 이득을 보기 위해 HDD 수에 맞게 파티션을 나누는 작업이 필요할 수 있습니다.

해당 예시에서는 두 개의 6TB, 2TB HDD, 240G SSD가 있을 때, OSD 2개 배포를 위해 SSD를 용량 비율에 따라 파티셔닝합니다. (각각 180G, 60G로)

- SSD 디바이스를 확인하고 파티션 작업 모드에 들어갑니다.
  - Command: fdisk -l

  ```shell
    root@master1:~# fdisk -l
    ...
    # /dev/sdg가 SSD 디바이스임을 확인
    Disk /dev/sdg: 240 GiB, 257698037760 bytes, 503316480 sectors
    Units: sectors of 1 * 512 = 512 bytes
    Sector size (logical/physical): 512 bytes / 512 bytes
    I/O size (minimum/optimal): 512 bytes / 512 bytes

  ```
  - Command: fdisk /dev/sdg

  ```shell
    root@master1:~# fdisk /dev/sdg

    Welcome to fdisk (util-linux 2.31.1).
    Changes will remain in memory only, until you decide to write them.
    Be careful before using the write command.

    Device does not contain a recognized partition table.
    Created a new DOS disklabel with disk identifier 0xe2f42834.
  ```
- g를 입력하여 GPT 테이블을 생성합니다.

  ```shell
    Command (m for help): g
    Created a new GPT disklabel (GUID: 106525D1-D2FA-1740-8235-7FE38DC1B560).
  ```
- n을 입력하여 파티션을 생성합니다.(예시에서는 n을 2번 입력하여 2개의 파티션을 생성합니다)
  - 파티션 생성시 원하는 용량을 입력해줍니다.

  ```shell
    Command (m for help): n
    Partition number (1-128, default 1):
    First sector (2048-503316446, default 2048):
    Last sector, +sectors or +size{K,M,G,T,P} (2048-503316446, default 503316446): +180G

    Created a new partition 1 of type 'Linux filesystem' and of size 180 GiB.

    Command (m for help): n
    Partition number (2-128, default 2):
    First sector (377489408-503316446, default 377489408):
    Last sector, +sectors or +size{K,M,G,T,P} (377489408-503316446, default 503316446): +60G

    Created a new partition 2 of type 'Linux filesystem' and of size 60 GiB.
  ```
- w을 입력하여 파티션 작업 내용을 반영합니다.
  ```shell
    Command (m for help): w

    The partition table has been altered.
    Calling ioctl() to re-read partition table.
    Syncing disks.
  ```

#### 2. Cluster 에 메타데이터 디바이스를 분리하도록 cluster.yaml 파일을 수정합니다.
- 아래 예시는 하나의 노드에 대한 OSD 배포 설정입니다. (튜닝 설정 이외의 `cluster.yaml` 설정에 대한 자세한 사항은 [cluster.yaml 수정 가이드](/docs/ceph-cluster-setting.md) 참고)
  - OSD 0은 데이터를 /dev/sde (HDD), 메타데이터를 /dev/sdg1 (SSD partition) 에 저장합니다.
  - OSD 1은 데이터를 /dev/sdf (HDD), 메타데이터를 /dev/sdg2 (SSD partition) 에 저장합니다.
```yaml
...
    nodes:
      - name: "master1"
        devices:
        - name: "sde"
          config:
            metadataDevice: "sdg1"
        - name: "sdf"
          config:
          metadataDevice: "sdg2"
...
```
#### 3. 배포 결과 확인
<!--- - 배포는 [installer 사용가이드](/docs/installer.md)를 참고하여 진행합니다.--->
- 배포 결과, 다음과 같이 SSD partition(sdg1,sdg2)에 lvm(ceph--block--dbs-*)이 생성된 것을 확인할 수 있습니다.
```shell
root@master1:/vagrant/hypercloud-rook-ceph# lsblk
NAME                                                                                                MAJ:MIN RM  SIZE RO TYPE MOUNTPOINT
...

sde                                                                                                   8:64   0    6T  0 disk
└─ceph--block--e0c23d2a--a6b7--44a6--9a70--99816a2e6190-osd--block--02251c26--bf20--4348--b835--e10f92936a4c
                                                                                                    253:2    0    6T  0 lvm
sdf                                                                                                   8:80   0    2T  0 disk
└─ceph--block--6b4b100a--10f0--49ab--aa46--3b156463945a-osd--block--fc259112--7694--4e61--a75f--74364f0f8e89
                                                                                                    253:0    0    2T  0 lvm
sdg                                                                                                   8:96   0  240G  0 disk
├─sdg1                                                                                                8:97   0  180G  0 part
│ └─ceph--block--dbs--ade73662--711a--4cc3--a06d--2f9ecd4d5ff7-osd--block--db--24d7f541--6b25--40fd--a7bf--f44091b3d705
│                                                                                                   253:3    0  180G  0 lvm
└─sdg2                                                                                                8:98   0   60G  0 part
  └─ceph--block--dbs--755a8476--a86d--471f--b23b--538f10e31ac0-osd--block--db--17fcc64b--2dfa--4b1f--8d0e--016b5b3221b1
                                                                                                    253:1    0   60G  0 lvm
```
