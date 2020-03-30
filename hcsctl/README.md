# hcsctl: hypercloud storage ctl
hcsctl은 hypercloud storage의 설치, 제거 및 관리를 제공합니다.

# Install
## Prerequisite
- kubectl (> 1.15.0)

## 설치
TODO: hcsctl 바이너리 업로드해두고 다운로드 링크 여기에 걸기

## 지원 기능 목록
- install
  - ex) hcsctl install myInventory
- uninstall
  - ex) hcsctl uninstall myInventory
- ceph {status/exec}
  - ex) hcsctl ceph status
  - ex) hcsctl ceph exec ceph osd status
  - ex) hcsctl ceph exec ceph df

# Quick Start
TODO: 자세하게 작성

```shell
hcsctl install my_inventory
hcsctl uninstall my_inventory

# e2e 테스트
hcsctl.test
```
