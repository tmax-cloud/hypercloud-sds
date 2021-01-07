## Rook Ceph Object Storage
> 본 문서는 Ceph의 Object Storage인 object store를 사용하는 방법에 대해서 다룹니다.
> 본 문서에서 rook ceph v1.1.6을 기준으로 합니다.

### Ceph Object Store 사용 방법

#### 1. How to Create Object Store
- /docs/examples에 있는 object-store.yaml를 이용하여 `S3 API`를 사용하는 `RGW service`를 deploy 해줍니다.
- object store를 구성하는 metadataPool과 dataPool의 replica size를 변경하고 싶은 경우, `replicated의 size` 필드를 수정해주시면 됩니다.
    - 현재 제공되는 object store의 설정에서는 replica size를 `n`으로 지정했을 경우에는 `n개 이상의 node에 osd`가 deploy되어 있어야 합니다.
    ```yaml
    # object-store.yaml의 pool 설정 부분입니다.
    spec:
      metadataPool:
        failureDomain: host
        replicated:
          size: 3    ## 수정하시면 됩니다.    
      dataPool:
        failureDomain: host
        replicated:
          size: 3    ## 수정하시면 됩니다.
    ```
- object-store.yaml deploy
    ```shell
	# object-store.yaml은 /docs/examples에 있습니다.
	# Object store를 생성합니다.
	$ kubectl create -f object-store.yaml
	
	# Object store의 생성을 확인합니다.
	$ kubectl -n rook-ceph get pod -l app=rook-ceph-rgw
	NAME                                        READY   STATUS    RESTARTS   AGE
    rook-ceph-rgw-my-store-a-6f576cd98f-7wkpc   1/1     Running   0          15s
    ```
- Kubernetes cluster 외부에서 접근하기 위해서는 NodePort를 위한 service를 추가해야 합니다.
	```shell
	# ceph-rgw service를 외부로 노출시키기 위해서 NodePort service deploy합니다.
    # rgw-external.yaml은 /docs/examples에 있습니다.
    $ kubectl create -f rgw-external.yaml
	
	# object store관련 service deploy를 확인하는 방법입니다.
	$ kubectl -n rook-ceph get service rook-ceph-rgw-my-store rook-ceph-rgw-my-store-external
	NAME                              TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)        AGE
    rook-ceph-rgw-my-store            ClusterIP   10.110.129.139   <none>        80/TCP         35m
    rook-ceph-rgw-my-store-external   NodePort    10.109.32.66     <none>        80:30187/TCP   107s
	```

#### 2. How to Create a Bucket
- 사용자가 bucket를 생성할 수 있도록 bucket에 대한 storageClass를 정의한 후, 해당 storageClass를 사용하여 bucket를 생성합니다.
- 아래의 예제에서 사용되는 yaml 파일들은 /docs/examples에 있습니다.
   - Deploy storageClass
	    ```shell
        # storageclass-bucket-delete.yaml은 /docs/examples에 있습니다.
        $ kubectl create -f storageclass-bucket-delete.yaml
        ```
		
    - Deploy Bucket
	    - Object Bucket Claim(OBC)을 deploy하면 bucket이 생성되고, 해당 bucket에 접근하기 위한 정보를 저장하고 있는 configMap과 secret이 생성됩니다.
		    - configMap과 secret는 OBC과 동일한 이름으로 생성됩니다.
			- 생성되는 bucket의 이름은 `spec의 bucketName` 필드를 통해 설정할 수 있습니다.
			    - 제공되는 yaml의 경우, bucketName이 ceph-object로 설정되어 있습니다.
	    ```shell
        # Deploy된 storageClass를 바탕으로, 사용자는 OBC를 사용하여 bucket를 생성할 수 있습니다.
        # object-bucket-claim-delete.yaml은 /docs/examples에 있습니다.
        $ kubectl create -f object-bucket-claim-delete.yaml
		
		# Bucket의 정상적인 생성 여부를 확인합니다.(bound가 나오면 됩니다)
		# kubectl get objectbucketclaims ${OBC_NAME} -o yaml | grep "Phase" | awk '{print $2}'
		$ kubectl get objectbucketclaims ceph-delete-bucket -o yaml | grep "Phase" | awk '{print $2}'
		bound
		
        # OBC deploy후에 operator는 bucket과 bucket의 접근 정보를 가지고 있는 configMap과 secret를 생성합니다.
        $ kubectl get configMap
        NAME                 DATA   AGE
        ceph-delete-bucket   6      5m7s
        $ kubectl get secret
        NAME                  TYPE                                  DATA   AGE
        ceph-delete-bucket    Opaque                                2      5m45s
        ```
#### 3. How to Get Env Path
> Bucket 접근을 위한 인자들을 fetch하는 방법에 대해서 설명합니다.
- access_key 와 access_id를 fetch하는 방법입니다.
    ```shell
    # 아래의 환경 변수는 OBC deploy인해서 생성된 configMap과 secret으로부터 fetch해옵니다.
    # configMap과 secret의 이름은 OBC의 이름과 동일합니다.
    
    # AWS_ACCESS_KEY_ID
    # echo $(kubectl -n default get secret ${OBC_NAME} -o yaml | grep AWS_ACCESS_KEY_ID | awk '{print $2}' | base64 --decode)
    $ echo $(kubectl -n default get secret ceph-delete-bucket -o yaml | grep AWS_ACCESS_KEY_ID | awk '{print $2}' | base64 --decode)
    206G6D4XZJVTHZSBRONX
    
    # AWS_SECRET_ACCESS_KEY
    # echo $(kubectl -n default get secret ${OBC_NAME} -o yaml | grep AWS_SECRET_ACCESS_KEY | awk '{print $2}' | base64 --decode)
    $ echo $(kubectl -n default get secret ceph-delete-bucket -o yaml | grep AWS_SECRET_ACCESS_KEY | awk '{print $2}' | base64 --decode)
    wcjKyZhMvappNj1MuBz5Lzp01zaDWSlJUXY0UxpS
    ```
- endPoint 정보를 fetch하는 방법입니다.
    - kubernetes cluster 외부와 내부에서 object store에 접근하기 위해서 필요한 endPoint를 각각 구합니다.
	```shell
	# kubernetes cluster 내부에서 사용하기 위한 endPoint를 설정하는 방법입니다.
	# 다음의 명령어를 통해 ClusterIP와 PORT를 조합해서 endPoint를 설정합니다.
	$ kubectl -n rook-ceph get service rook-ceph-rgw-my-store
	NAME                     TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)   AGE
    rook-ceph-rgw-my-store   ClusterIP   10.110.129.139   <none>        80/TCP    41m
	# 위 경우 endPoint는 10.110.129.139:80입니다.
	
	# kubernetes cluster 외부에서 사용하기 위한 endPoint를 설정하는 방법입니다.
	# 다음의 명령어를 통해 HostIP와 PORT를 조합해서 endPoint를 설정합니다. HostIP의 경우, k8s cluster를 구성하는 node들의 ip 중 하나의 ip입니다.
	$ kubectl -n rook-ceph get service rook-ceph-rgw-my-store-external
	NAME                              TYPE       CLUSTER-IP     EXTERNAL-IP   PORT(S)        AGE
    rook-ceph-rgw-my-store-external   NodePort   10.109.32.66   <none>        80:30187/TCP   13m
    # 위 경우 endPoint는 192.168.50.90:30187입니다.(hostIP가 192.168.50.90)
	```
	
#### 4. How to Use Bucket
> 본 문서에서는 `s3cmd`를 통해 bucket를 사용하는 경우를 설명합니다.
- 다음은 kubernetes cluster 외부에서 object store를 사용하는 것에 대한 예시입니다.
    - kubernetes cluster 내부에서 사용하기 위해서는 endpoint를 section 3를 참고해서 수정하시면 됩니다.
    ```shell
    # s3cmd put 예시입니다.
    $ echo "Test Ceph Object Store" > testCaseOne
    # s3cmd put ${fileName} --no-ssl --host=${endPoint} --host-bucket= s3://${bucketName} --access_key=${AWS_ACCESS_KEY_ID} --secret_key=${AWS_SECRET_ACCESS_KEY}
    $ s3cmd put ./testCaseOne --no-ssl --host=http://192.168.50.90:30187 --host-bucket= s3://ceph-object --access_key=206G6D4XZJVTHZSBRONX --secret_key=wcjKyZhMvappNj1MuBz5Lzp01zaDWSlJUXY0UxpS
    upload: './testCaseOne' -> 's3://ceph-object/testCaseOne'  [1 of 1]
     23 of 23   100% in    0s  1089.79 B/s  done
    
    # s3cmd get 예시입니다.
    # s3cmd get --no-ssl --host=${endPoint} --host-bucket= s3://${bucketName}/${objectName} --access_key=${AWS_ACCESS_KEY_ID} --secret_key=${AWS_SECRET_ACCESS_KEY}
    $ s3cmd get --no-ssl --host=http://192.168.50.90:30187 --host-bucket= s3://ceph-object/testCaseOne --access_key=206G6D4XZJVTHZSBRONX --secret_key=wcjKyZhMvappNj1MuBz5Lzp01zaDWSlJUXY0UxpS
    download: 's3://ceph-object/testCaseOne' -> './testCaseOne'  [1 of 1]
     23 of 23   100% in    0s   462.01 B/s  done
    $ cat testCaseOne
    Test Ceph Object Store
    ```