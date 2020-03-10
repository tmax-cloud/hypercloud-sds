# Velero (formerly Heptio Ark)
>  VMware회사에서 만들어졌던 Open-source 툴이며 Kuberneter cluster resources와 Persistent Volumes까지 Backup(On-demand, schedule), Restore, Migration 기능을 제공해주고 있습니다.
Reference: [Docs](https://velero.io/docs/v1.2.0/), [Github](https://github.com/vmware-tanzu/velero)

## 구성 
1.  **Client**: kubectl처럼 local에서 실행하며 Velero API를 요청할 수 있게 해주는 CLI입니다.
2.  **Server**: Kubernetes cluster 위에서 뜨고 있는 POD Deployement로 생성되고 Velero CRD를 감시해서 요청을 받고 저리해 주는 것입니다.

  
## 설치
* **Client**
```shell
wget https://github.com/vmware-tanzu/velero/releases/download/v1.2.0/velero-v1.2.0-linux-amd64.tar.gz
tar zxf velero-v1.2.0-linux-amd64.tar.gz
mv velero*/velero /usr/local/bin
```

* **Server**  
  * velero install CLI command로 설치 


  ```shell
  velero install \
   --provider aws \
   --bucket kubevelero \
   --secret-file ./minio.credentials \
   --backup-location-config region=minio,s3ForcePathStyle=true,s3Url=http://192.168.0.152:9000 \
   --plugins velero/velero-plugin-for-aws:v1.0.0 \
   --use-volume-snapshots false \
   --use-restic
   
  # --provider : the cloud provider 이름
  # --bucket : bucket 이름
  # --secret-file : access_key and secret_key 정보를 갖고 있는 파일 경로
  # --backup-location-config : backup storage location 설정
  # --plugins : 사용하는 provider에 따라 image plugin
  # --use-volume-snapshots : PV backup할 때 cloud provider snapshot 기능 사용
  # --user-restic : PV backup 할 때 Restic 사용
  ```

   
  * Helm charts로 설치 ([참고](https://github.com/vmware-tanzu/helm-charts))
  ```shell
  # velero repository 추가
  helm repo add vmware-tanzu https://vmware-tanzu.github.io/helm-charts
  # 설치 charts 찾기
  helm search repo vmware-tanzu
  ```
  
## 삭제 
```shell
kubectl delete namespace/velero clusterrolebinding/velero
kubectl delete crds -l component=velero
```

## 사용 
* **Backup**   
    >  [Note] PV를 백업할 때 `Restic` 사용하고 싶으면 annotation를 먼저 만들어야 합니다.  
    >  $ kubectl annotate {POD_NAME} backup.velero.io/backup-volumes={PV_NAME}
  
  * **On-demand backup** (API를 통해 유저가 원할 때 즉시 백업)
      ```shell
      # 요청:
      velero backup create {BACKUP_NAME} --include-namespaces {NAMESPACE_NAME}
      
      # 확인:
      velero backup get
      velero backup describe {BACKUP_NAME}
      kubectl -n velero get backups
      ```
  * **Schedule backup**  (특정 시간마다 자동으로 백업하도록 설정)
    ```shell
    # 요청:
    velero schedule create {SCHEDULE_NAME} --schedule="@every 1d" --include-namespaces {NAMESPACE_NAME}
  
    # 확인:
    velero get schedule
    ```
* **Restore**
  ```shell
  velero restore create {RESTORE_NAME} --from-backup {BACKUP_NAME}
  ```
* **Delete**
  ```shell
  velero schedule delete {SCHEDULE_NAME}
  velero backup delete --all
  velero restore delete --all
  ```
  
# 추가 
  velero help CLI에서 참고해서 사용하시면 좋겠습니다.
  * Backup
    ```shell
    # Backup entire cluster
    velero backup create NAME
    
    # Backup entire cluster excluding namespaces
    velero backup create NAME --exclude-namespaces testing
    
    # Backup entire cluster excluding resources 
    velero backup create NAME --exclude-resources configmaps
    
    # Backup entire cluster only some resources
    velero backup create NAME --include-resources pods,deployments
    
    # Backup entire namespaces
    velero backup create NAME --include-namespaces testing
    
    # Backup entire namespaces only some resources
    velero backup create NAME --include-namespaces testing --include-resources pods,deployments
    
    # Backup entire namespaces excluding some resources
    velero backup create NAME --include-namespaces testing --exclude-resources pods,deployments
    
    ```
  * Schedule
    ```shell
    velero schedule create NAME --schedule="* * * * *"
    # */Minute */Hour */Day of Month */Month */Day of Week
    
    # Create a backup every 6 hours
    velero schedule create NAME --schedule="0 */6 * * *"
    
    # Create a backup every 6 hours with the @every notation
    velero create schedule NAME --schedule="@every 6h"
    velero create schedule NAME --schedule="@every 1w"
    
    # Create a daily backup of the web namespace
    velero create schedule NAME --schedule="@every 24h" --include-namespaces web
    
    # Create a weekly backup, each living for 90 days (2160 hours)
    velero create schedule NAME --schedule="@every 168h" --ttl 2160h0m0s
    
    ```
  * Restore
    ```shell
    # create a restore named "restore-1" from backup "backup-1"
    velero restore create restore-1 --from-backup backup-1
    
    # create a restore with a default name ("backup-1-<timestamp>") from backup "backup-1"
    velero restore create --from-backup backup-1
     
    # create a restore from the latest successful backup triggered by schedule "schedule-1"
    velero restore create --from-schedule schedule-1
    
    # create a restore from the latest successful OR partially-failed backup triggered by schedule "schedule-1"
    velero restore create --from-schedule schedule-1 --allow-partially-failed
    
    # create a restore for only persistentvolumeclaims and persistentvolumes within a backup
    velero restore create --from-backup backup-2 --include-resources persistentvolumeclaims,persistentvolumes
  
    ```
  