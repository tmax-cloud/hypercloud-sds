# hypercloud-etcd

본 문서는 k8s의 etcd를 백업하는 방법에 대해 다룹니다.


## What To Do

일정 주기마다 etcd의 snapshot을 찍어서 hostPath에 저장하는 CronJob을 생성하여 수행시킵니다.


## Caution

- etcd의 port는 2379를 사용함
- Admin의 인증서 파일들은 사전에 backup되어야 함
  - (master node) `/etc/kubernetes/pki/ca.crt`
  - (master node) `/etc/kubernetes/pki/ca.key`


## How To

1. Script를 수행시킵니다. (Command : ./tools/etcd_yaml_change_script.sh {$registry_endpoint} {$master_node_name})
    ```shell
    $ ./tools/etcd_yaml_change_script.sh 192.168.50.90:5000 master1
    ```
    - 이 작업은 `etcd_snapshot.yaml` 파일에 해당 환경의 registry endpoint와 master node의 이름을 설정해줍니다.


2. `etcd_snapshot.yaml` 파일의 `spec: schedule` 필드의 값을 사용하고자 하는 crontab 포맷으로 변경합니다.
    - 5분마다 수행 : `"*/5 * * * *"`
    - 매시 0분마다 수행 : `"0 * * * *"`
    - 매일 오전 6시에 수행 : `"0 6 * * *"`
    - 자세한 crontab 포맷 사용법은 [cron wikipedia](https://en.wikipedia.org/wiki/Cron#Overview)를 참조하세요.


3. etcd의 version을 확인합니다.
    ```shell
    $ sudo cat /etc/kubernetes/manifests/etcd.yaml | grep "image"
    ...
    image: k8s.gcr.io/etcd:3.3.10   # etcd version은 k8s.gcr.io/etcd:3.3.10
    ...
    ```


4. 확인한 etcd image version을 `etcd_snapshot.yaml` 파일의 `spec: jobTemplate: spec: template: spec: containers: image` 필드에 반영합니다.
    ```yaml
    ...
    spec:
      schedule: "*/5 * * * *"
      jobTemplate:
        spec:
          template:
            spec:
              containers:
              - name: etcd-backup
                # image: {$registry_endpoint}/{$etcd_image_version} 
                image: 192.168.50.90:5000/k8s.gcr.io/etcd:3.3.10
    ...
    ```


5. `etcd_snapshot.yaml` 파일을 적용하여 CronJob을 생성합니다.
    ```shell
    $ kubectl apply -f etcd_snapshot.yaml
    ```
    - 현재 etcdctl 인증을 위한 파일들이 위치한 디렉토리 경로는 설정한 master node의 `/etc/kubernetes/pki/etcd`입니다.
    - 현재 etcd backup이 저장되는 디렉토리 경로는 설정한 master node의 `/mnt/backup`입니다.
        ```yaml
        ...
            volumes:
              - name: etcd-certs
                # etcdctl 인증 파일들의 위치는 hostPath /etc/kubernetes/pki/etcd 디렉토리
                hostPath:
                  path: /etc/kubernetes/pki/etcd
                  type: DirectoryOrCreate
              - name: backup
                # etcd backup의 위치는 hostPath /mnt/backup 디렉토리
                hostPath:
                  path: /mnt/backup
                  type: DirectoryOrCreate
        ...
        ```
    - 보다 안전한 backup을 위해 hostPath volume 대신 AWS EBS 등을 사용할 수 있습니다.