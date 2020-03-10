## ENV
ubuntu 18.10

## iscsi target (server)
```shell
# iscsi target 설치
apt-get install tgt
service tgt status

# iscsi target 생성
# tgtadm -lld iscsi --op new --mode target --tid tid -T targetName
tgtadm -L iscsi -o new -m target -t 1 -T iqn.2019-11.com.tmax:storage

# iscsi target에 volume 등록
tgtadm -L iscsi -o new -m logicalunit -t 1 -l 1 -b /dev/vg/lv

# 모든 시스템에서 initiator가 접근하는 것을 허용
tgtadm -L iscsi -o bind -m target -t 1 -I ALL

# LUN 삭제
tgtadm -L iscsi -o delete -m logicalunit -t 1 -l 1

# iscsi target 삭제
tgtadm -L iscsi -o delete -m target -t 1 --force
```

## iscsi initiator (client)
```shell
# iscsi initiator 설치
apt-get install open-iscsi

# iscsi target 검색
iscsiadm -m discovery -t sendtargets -p 192.168.7.18

# iscsi target 로그인
iscsiadm -m node -T iqn.2019-11.com.tmax:storage -p 192.168.7.18 --login

# iscsi target 로그아웃
iscsiadm -m node -T iqn.2019-11.com.tmax:storage -p 192.168.7.18:3260 -u

# iscsi initiator 삭제
iscsiadm -m node -o delete -T iqn.2019-11.com.tmax:storage
```

## pv 생성
```yml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: pv-iscsi
spec:
  accessModes:
    - ReadWriteOnce
  capacity:
    storage: 100Mi
  iscsi:
    targetPortal: 172.22.4.2:3260
    iqn: iqn.2019-11.com.server:storage.target9
    lun: 9
    fsType: 'ext4'
    readOnly: false

```

## pvc 생성
```yml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: pvc-iscsi
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 50Mi
  storageClassName: ""
```

## pod 생성
```yml
apiVersion: v1
kind: Pod
metadata:
  name: test
spec:
  containers:
  - name: test
    image: nginx
    volumeMounts:
    - mountPath: /mnt/test
      name: pvc
  volumes:
  - name: pvc
    persistentVolumeClaim:
      claimName: pvc-iscsi
```
