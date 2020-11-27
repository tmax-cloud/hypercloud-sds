## Rook Ceph Version Upgrade 1.3.6 to 1.4.2

> 본 문서는 rook ceph v1.3.6에서 rook ceph v1.4.2로 업그레이드하는 방법을 다룹니다.

### Kubernetes 및 HyperSDS 버전 설명

- Kubernetes 버전은 1.17 이상이어야 합니다.
- HyperSDS는 `release-1.3`에서 rook ceph v1.3.6을 사용합니다.
- 본 문서는 HyperSDS release-1.3이 설치되어 있는 환경에서 rook과 ceph, ceph-csi의 버전을 업그레이드하는 방법을 다룹니다.

### 주의사항

- Rook ceph 버전 업그레이드를 수행할 경우, 예상치 못한 issue가 발생할 수도 있으므로, 중요한 데이터는 백업 후 업그레이드를 수행하는 것을 권장합니다.
- 업그레이드를 진행하는 동안에는, ceph volume 생성, 변경, 이용 등의 작업을 중단하시는 것을 권장합니다.
- <strong> 반드시 각 단계가 완료된 후, 다음 단계를 진행해주시길 바랍니다.</strong>

### Rook Ceph Upgrade
> Rook ceph version -> Ceph csi version -> Ceph version 순서로 업그레이드를 진행합니다.


#### 1. Rook ceph cluster 상태 확인

- 업그레이드 시작 전에 rook ceph cluster 상태를 확인합니다.
	- Rook ceph 관련 Pod들이 모두 정상적으로 수행되고 있는지 확인합니다.
	- Ceph command를 통해 rook ceph cluster의 상태를 확인합니다.
		- Ceph cluster의 heath status가 HEALTH_OK인지 확인합니다.
		- Ceph mon들이 모두 quorum에 포함되어 있는지 확인합니다.
		- Ceph mgr이 active 상태인지 확인합니다.
		- 모든 ceph osd들이 up & in 상태인지 확인합니다.
		- 모든 pg가 active + clean 상태인지 확인합니다.
		- Ceph mds가 active 상태인지 확인합니다.
		
	```shell
    # Rook ceph cluster를 구성하는 pod들이 모두 RUNNING 상태인지 확인합니다.
    $ kubectl -n rook-ceph get pods -o wide
	
	# Ceph command인 "ceph status"를 사용하여, ceph cluster의 상태를 확인합니다.
    # Rook ceph toolbox pod를 통해서 ceph command를 입력할 수 있습니다.
	# Rook ceph toolbox pod의 이름은 환경마다 다르므로, 해당 환경에서의 toolbox pod의 이름을 꼭 확인해주세요.
	$ kubectl exec -it -n rook-ceph rook-ceph-tools-847744596b-cdv9m  -- ceph status
	cluster:
      id:     14f3596b-7429-45a5-9061-1ea33de648e1
      health: HEALTH_OK

    services:
      mon: 1 daemons, quorum a (age 18m)
      mgr: a(active, since 18m)
      mds: myfs:1 {0=myfs-b=up:active} 1 up:standby-replay
      osd: 3 osds: 3 up (since 17m), 3 in (since 17m)

    data:
      pools:   3 pools, 96 pgs
      objects: 22 objects, 2.2 KiB
      usage:   3.0 GiB used, 747 GiB / 750 GiB avail
      pgs:     96 active+clean
	

#### 2. Alpha version의 snapshot 제거
> Rook ceph v1.4부터는 ceph-csi v3.0 이상이 권장되므로, ceph-csi v2.0에서 사용하던 v1alpha1 API version의 snapshot이 호환되지 않습니다. 
> 따라서 v1alpha1 snapshot 및 snapshot CRD가 있다면 제거하는 작업이 필요합니다.

- STEP A) Snapshot 관련 CRD들의 API version이 v1alpha1인지 확인합니다.

	```shell
	# VolumeSnapshotClass의 API version이 v1alpha1인지 확인합니다.
	$ kubectl get crd volumesnapshotclasses.snapshot.storage.k8s.io -o yaml |grep v1alpha1
      - name: v1alpha1
      - v1alpha1
	# VolumeSnapshotContent의 API version이 v1alpha1인지 확인합니다.
	$ kubectl get crd volumesnapshotcontents.snapshot.storage.k8s.io -o yaml |grep v1alpha1
      - name: v1alpha1
      - v1alpha1
	# VolumeSnapshot의 API version이 v1alpha1인지 확인합니다.
	$ kubectl get crd volumesnapshots.snapshot.storage.k8s.io -o yaml |grep v1alpha1
      - name: v1alpha1
      - v1alpha1
	```
	
- STEP B) Snapshot API version이 v1alpha1이라면 존재하는 snapshot 관련 오브젝트들을 제거합니다.

	```shell
	# API version이 v1alpha1였다면, snapshot이 존재하는지 확인합니다.
	$ kubectl get volumesnapshot
	NAME               AGE
	rbd-pvc-snapshot   22s
	
	# 존재한다면, 해당 snapshot을 제거해줍니다.
	$ kubectl delete volumesnapshot rbd-pvc-snapshot
	volumesnapshot.snapshot.storage.k8s.io "rbd-pvc-snapshot" deleted
	
	# 위의 과정을 volumesnapshotcontent, volumesnapshotclass에 대해서도 수행해줍니다.
	```
	
- STEP C) Snapshot API version이 v1alpha1이라면 snapshot 관련 CRD들을 제거합니다.

	```shell
	# API version이 v1alpha1이라면 VolumeSnapshotClass CRD를 제거합니다.
	$ kubectl delete crd volumesnapshotclasses.snapshot.storage.k8s.io 
	customresourcedefinition.apiextensions.k8s.io "volumesnapshotclasses.snapshot.storage.k8s.io" deleted
	
	# API version이 v1alpha1이라면 VolumeSnapshotContent CRD를 제거합니다.
	$ kubectl delete crd volumesnapshotcontents.snapshot.storage.k8s.io 
	customresourcedefinition.apiextensions.k8s.io "volumesnapshotcontents.snapshot.storage.k8s.io" deleted
	
	# API version이 v1alpha1이라면 VolumeSnapshot CRD를 제거합니다.
	$ kubectl delete crd volumesnapshots.snapshot.storage.k8s.io
	customresourcedefinition.apiextensions.k8s.io "volumesnapshots.snapshot.storage.k8s.io" deleted
	```


#### 3. Beta version의 snapshot 설치
> API version v1beta1의 snapshot 관련 CRD들과 controller를 설치합니다.

- STEP A) v1beta1 API version의 snapshot CRD들을 설치합니다.

	```shell
	# v1beta1 API version의 VolumeSnapshot, VolumeSnapshotClass, VolumeSnapshotContent CRD를 설치합니다.
	$ kubectl create -f 1_snapshot_crds.yaml
	customresourcedefinition.apiextensions.k8s.io/volumesnapshots.snapshot.storage.k8s.io created
	customresourcedefinition.apiextensions.k8s.io/volumesnapshotcontents.snapshot.storage.k8s.io created
	customresourcedefinition.apiextensions.k8s.io/volumesnapshotclasses.snapshot.storage.k8s.io created
	```

- STEP B) v1beta1 snapshot을 관리하는 controller 및 해당 controller의 RBAC 관련 오브젝트들을 설치합니다.
	
	```shell
	# Snapshot controller의 ServiceAccount, ClusterRole 등 RBAC 관련 오브젝트들을 설치합니다.
	$ kubectl create -f 2_snapshot_controller_rbac.yaml
	serviceaccount/snapshot-controller created
	clusterrole.rbac.authorization.k8s.io/snapshot-controller-runner created
	clusterrolebinding.rbac.authorization.k8s.io/snapshot-controller-role created
	role.rbac.authorization.k8s.io/snapshot-controller-leaderelection created
	rolebinding.rbac.authorization.k8s.io/snapshot-controller-leaderelection created
	
	# Snapshot controller를 설치합니다.
	$ kubectl create -f 3_snapshot_controller.yaml
	statefulset.apps/snapshot-controller created
	```


#### 4. Rook v1.4 기준의 CRD와 RBAC 적용
> Rook v1.3에서 쓰이던 CRD와 RBAC 관련 오브젝트들을 Rook v1.4 기준에 맞게 업그레이드합니다.

- STEP A) Rook v1.4에서 변경된 ClusterRole들을 제거하고, v1.4 기준의 CRD와 RBAC 관련 오브젝트들을 설치합니다.

	```shell
	# Rook v1.3의 ClusterRole을 제거합니다.
	$ kubectl delete -f 4_rook_rbac_v1_3.yaml
	clusterrole.rbac.authorization.k8s.io "cephfs-csi-nodeplugin-rules" deleted
	clusterrole.rbac.authorization.k8s.io "cephfs-external-provisioner-runner-rules" deleted
	clusterrole.rbac.authorization.k8s.io "rbd-csi-nodeplugin-rules" deleted
	...
	
	# Rook v1.4의 ServiceAccount, ClusterRole 등 RBAC 관련 오브젝트들을 설치합니다.
	$ kubectl create -f 5_rook_rbac_v1_4.yaml
	clusterrole.rbac.authorization.k8s.io/rook-ceph-global created
	clusterrole.rbac.authorization.k8s.io/rook-ceph-cluster-mgmt created
	clusterrole.rbac.authorization.k8s.io/rook-ceph-mgr-cluster created
	...
	
	# Rook v1.4에서 사용되는 CRD들을 설치합니다.
	$ kubectl create -f 6_rook_crds_v1_4.yaml
	customresourcedefinition.apiextensions.k8s.io/cephrbdmirrors.ceph.rook.io created
	customresourcedefinition.apiextensions.k8s.io/cephobjectrealms.ceph.rook.io created
	customresourcedefinition.apiextensions.k8s.io/cephobjectzonegroups.ceph.rook.io created
	...
	
	```


#### 5. CSI version 업그레이드
> CSI container image version을 ConfigMap에 직접 명시한 경우, image version을 수정해줘야 합니다.

> <strong> 인벤토리의 `rook/operator.yaml`에서 CSI image version을 직접 기입하지 않았다면 다음 단계로 넘어가시기 바랍니다.</strong>

> <strong> 혹은 별도로 `rook-ceph-operator-config` ConfigMap에 CSI image version을 직접 기입하지 않았다면 다음 단계로 넘어가시기 바랍니다.</strong>

- STEP A) rook-ceph-operator-config ConfigMap에 이전 version의 CSI version들이 기입돼있는지 확인합니다.

	```shell
	# Rook v1.3.6에서 기본값으로 사용하는 CSI version은 아래와 같습니다.
	# 만약 아무 값도 출력되지 않는다면, CSI version을 직접 명시하지 않은 경우입니다. 다음 단계로 넘어가시기 바랍니다.
	$ kubectl -n rook-ceph get configmaps rook-ceph-operator-config -o yaml | grep -v f:ROOK_CSI_ | grep -v \"R | grep 'ROOK_CSI_.*_IMAGE'
	  ROOK_CSI_ATTACHER_IMAGE: quay.io/k8scsi/csi-attacher:v2.1.0
	  ROOK_CSI_CEPH_IMAGE: quay.io/cephcsi/cephcsi:v2.1.2
	  ROOK_CSI_PROVISIONER_IMAGE: quay.io/k8scsi/csi-provisioner:v1.4.0
	  ROOK_CSI_REGISTRAR_IMAGE: quay.io/k8scsi/csi-node-driver-registrar:v1.2.0
	  ROOK_CSI_RESIZER_IMAGE: quay.io/k8scsi/csi-resizer:v0.4.0
	  ROOK_CSI_SNAPSHOTTER_IMAGE: quay.io/k8scsi/csi-snapshotter:v1.2.2
	```

- STEP B) rook-ceph-operator-config ConfigMap에서 CSI version을 최신으로 변경합니다.
	- 현 단계에서는 rook operator yaml 파일의 내용만 수정하고, 실제 적용은 다음 단계인 6번에서 진행합니다.
	- 업그레이드하고자 하는 ceph-csi 및 kubernetes-csi version은 다음과 같습니다.
		- quay.io/cephcsi/cephcsi:v3.1.0
		- quay.io/k8scsi/csi-attacher:v2.1.0
		- quay.io/k8scsi/csi-node-driver-registrar:v1.2.0
		- quay.io/k8scsi/csi-provisioner:v1.6.0
		- quay.io/k8scsi/csi-resizer:v0.4.0
		- quay.io/k8scsi/csi-snapshotter:v2.1.1

	```shell
	# CSI version을 ceph-csi v3.1.0에 맞춰 변경합니다.
	# Ceph cluster 설치에 사용한 인벤토리의 rook/operator.yaml 파일에서 CSI image를 변경해줍니다.
	# (예제에서는 /home/k8s/hypercloud-sds/hcsctl/my_inventory를 설치에 사용한 인벤토리라고 가정)
	# rook-ceph-operator-config ConfigMap의 image만 변경하여 적용해주시기 바랍니다.
	# (ROOK_CSI_*_IMAGE)
	$ INVENTORY=/home/k8s/hypercloud-sds/hcsctl/my_inventory
	$ cat $INVENTORY/rook/operator.yaml
	
	...
	
	kind: ConfigMap
	apiVersion: v1
	metadata:
	  name: rook-ceph-operator-config
	  # should be in the namespace of the operator
	  namespace: rook-ceph
	data:
	
	  ...
	  
	  # The default version of CSI supported by Rook will be started. To change the version
	  # of the CSI driver to something other than what is officially supported, change
	  # these images to the desired release of the CSI driver.
	  ROOK_CSI_CEPH_IMAGE: "quay.io/cephcsi/cephcsi:v2.1.2"                        ## v3.1.0으로 변경합니다.
	  ROOK_CSI_REGISTRAR_IMAGE: "quay.io/k8scsi/csi-node-driver-registrar:v1.2.0"
	  ROOK_CSI_RESIZER_IMAGE: "quay.io/k8scsi/csi-resizer:v0.4.0"
	  ROOK_CSI_PROVISIONER_IMAGE: "quay.io/k8scsi/csi-provisioner:v1.4.0"          ## v1.6.0으로 변경합니다.
	  ROOK_CSI_SNAPSHOTTER_IMAGE: "quay.io/k8scsi/csi-snapshotter:v1.2.2"          ## v2.1.1로 변경합니다.
	  ROOK_CSI_ATTACHER_IMAGE: "quay.io/k8scsi/csi-attacher:v2.1.0"
	  
	...
	
	# Rook operator version을 업그레이드할 때 CSI version도 정상적으로 업그레이드되었는지 확인해야 합니다.
	# 해당 커맨드는 6번에서 설명합니다.
	```


#### 6. Rook operator version 업그레이드

- STEP A) 현재 ceph cluster를 구성하는 Deployment들의 rook version이 v1.3.6임을 확인합니다.

	```shell
	# ceph cluster를 구성하는 Deployment들이 rook v1.3.6임을 확인합니다.
	# 모든 Deployment에 대하여 req/upd/avl: 1/1/1 값이 뜨면 정상적으로 작동하는 상태입니다.
	$ kubectl -n rook-ceph get deployments -l rook_cluster=rook-ceph -o jsonpath='{range .items[*]}{.metadata.name}{"  \treq/upd/avl: "}{.spec.replicas}{"/"}{.status.updatedReplicas}{"/"}{.status.readyReplicas}{"  \trook-version="}{.metadata.labels.rook-version}{"\n"}{end}'
	rook-ceph-mds-myfs-a    req/upd/avl: 1/1/1      rook-version=v1.3.6
	rook-ceph-mds-myfs-b    req/upd/avl: 1/1/1      rook-version=v1.3.6
	rook-ceph-mgr-a         req/upd/avl: 1/1/1      rook-version=v1.3.6
	rook-ceph-mon-a         req/upd/avl: 1/1/1      rook-version=v1.3.6
	rook-ceph-osd-0         req/upd/avl: 1/1/1      rook-version=v1.3.6
	rook-ceph-osd-1         req/upd/avl: 1/1/1      rook-version=v1.3.6
	rook-ceph-osd-2         req/upd/avl: 1/1/1      rook-version=v1.3.6
	```

- STEP B) Rook ceph operator version을 v1.4.2로 업그레이드합니다.
	
	```shell
	# rook operator의 version을 v1.4.2로 업그레이드합니다.
	# Ceph cluster를 설치하는 데 사용한 인벤토리의 rook/operator.yaml 파일에서 Rook image를 변경해줍니다.
	# (예제에서는 /home/k8s/hypercloud-sds/hcsctl/my_inventory를 설치에 사용한 인벤토리라고 가정)
	# rook-ceph-operator Deployment의 image를 rook/ceph:v1.4.2 값으로 변경하여 적용해주시기 바랍니다.
	# (spec.template.spec.containers[0].image)
	$ INVENTORY=/home/k8s/hypercloud-sds/hcsctl/my_inventory
	$ cat $INVENTORY/rook/operator.yaml
	
	...
	
	apiVersion: apps/v1
	kind: Deployment
	metadata:
	  name: rook-ceph-operator
	  namespace: rook-ceph
	  labels:
		operator: rook
		storage-backend: ceph
	spec:
	
	...
	
	      containers:
		  - name: rook-ceph-operator
			image: rook/ceph:v1.3.6    ## rook/ceph:v1.4.2 로 변경합니다.
			args: ["ceph", "operator"]

	...
	
	# 변경된 Rook version으로 재적용합니다.
	$ kubectl apply -f $INVENTORY/rook/operator.yaml
	configmap/rook-ceph-operator-config unchanged
	deployment.apps/rook-ceph-operator configured
	
	# 위에서 확인했던 ceph cluster Deployment들이 모두 rook v1.4.2 version으로 업그레이드 되어 정상 작동하기까지 기다립니다.
	# req/upd/avl 값이 1/1/1, rook-version 값이 v1.4.2가 되면 정상적으로 업그레이드된 상태입니다.
	$ watch --exec kubectl -n rook-ceph get deployments -l rook_cluster=rook-ceph -o jsonpath='{range .items[*]}{.metadata.name}{"  \treq/upd/avl: "}{.spec.replicas}{"/"}{.status.updatedReplicas}{"/"}{.status.readyReplicas}{"  \trook-version="}{.metadata.labels.rook-version}{"\n"}{end}'
	
	# 아래는 업그레이드 중일 때 출력되는 예시입니다. 
	rook-ceph-mds-myfs-a    req/upd/avl: 1/1/1      rook-version=v1.3.6
	rook-ceph-mds-myfs-b    req/upd/avl: 1/1/1      rook-version=v1.3.6
	rook-ceph-mgr-a         req/upd/avl: 1//        rook-version=v1.4.2
	rook-ceph-mon-a         req/upd/avl: 1/1/1      rook-version=v1.4.2
	rook-ceph-osd-0         req/upd/avl: 1/1/1      rook-version=v1.3.6
	rook-ceph-osd-1         req/upd/avl: 1/1/1      rook-version=v1.3.6
	rook-ceph-osd-2         req/upd/avl: 1/1/1      rook-version=v1.3.6
	
	# 만약 CSI version을 직접 업그레이드하였다면 아래 커맨드를 통해 변경된 CSI version으로 정상적으로 업그레이드되었는지 확인합니다.
	$ kubectl -n rook-ceph get pod -o jsonpath='{range .items[*]}{range .spec.containers[*]}{.image}{"\n"}' -l 'app in (csi-rbdplugin,csi-rbdplugin-provisioner,csi-cephfsplugin,csi-cephfsplugin-provisioner)' | sort | uniq
	
	quay.io/cephcsi/cephcsi:v3.1.0
	quay.io/k8scsi/csi-attacher:v2.1.0
	quay.io/k8scsi/csi-node-driver-registrar:v1.2.0
	quay.io/k8scsi/csi-provisioner:v1.6.0
	quay.io/k8scsi/csi-resizer:v0.4.0
	quay.io/k8scsi/csi-snapshotter:v2.1.1
	```

- STEP C) Toolbox Deployment의 rook version을 v1.4.2로 업그레이드합니다.

	```shell
	# Toolbox의 image version을 v1.4.2로 업그레이드합니다.
	# (spec.template.spec.containers[0].image)
	$ cat $INVENTORY/rook/toolbox.yaml
	...
	  template:
		metadata:
		  labels:
			app: rook-ceph-tools
		spec:
		  dnsPolicy: ClusterFirstWithHostNet
		  containers:
		  - name: rook-ceph-tools
			image: rook/ceph:v1.3.6   ## rook/ceph:v1.4.2 로 변경합니다.
			command: ["/tini"]
	...

	# 변경된 Rook version으로 재적용합니다.
	$ kubectl apply -f $INVENTORY/rook/toolbox.yaml
	```


#### 7. Ceph version 업그레이드
> Rook v1.4.2는 Ceph 15.2.4(Octopus)를 배포합니다. Ceph major version을 업그레이드하는 것을 권장합니다.

- STEP A) 현재 ceph cluster를 구성하는 Deployment들의 ceph version이 v14.2.9임을 확인합니다.

	```shell
	# ceph cluster를 구성하는 Deployment들이 ceph v14.2.9임을 확인합니다.
	# 모든 Deployment에 대하여 req/upd/avl: 1/1/1 값이 뜨면 정상적으로 작동하는 상태입니다.
	$ kubectl -n rook-ceph get deployments -l rook_cluster=rook-ceph -o jsonpath='{range .items[*]}{.metadata.name}{"  \treq/upd/avl: "}{.spec.replicas}{"/"}{.status.updatedReplicas}{"/"}{.status.readyReplicas}{"  \tceph-version="}{.metadata.labels.ceph-version}{"\n"}{end}'
	rook-ceph-mds-myfs-a    req/upd/avl: 1/1/1      ceph-version=14.2.9-0
	rook-ceph-mds-myfs-b    req/upd/avl: 1/1/1      ceph-version=14.2.9-0
	rook-ceph-mgr-a         req/upd/avl: 1/1/1      ceph-version=14.2.9-0
	rook-ceph-mon-a         req/upd/avl: 1/1/1      ceph-version=14.2.9-0
	rook-ceph-osd-0         req/upd/avl: 1/1/1      ceph-version=14.2.9-0
	rook-ceph-osd-1         req/upd/avl: 1/1/1      ceph-version=14.2.9-0
	rook-ceph-osd-2         req/upd/avl: 1/1/1      ceph-version=14.2.9-0
	```
	
- STEP B) Ceph cluster의 Ceph version을 v15.2.4로 업그레이드합니다.

	```shell
	# Ceph cluster를 설치하는 데 사용한 인벤토리의 rook/cluster.yaml 파일에서,
	# image 값을 ceph/ceph:v15.2.4 값으로 변경하여 적용해주시기 바랍니다. (spec.cephVersion.image)
	# (예제에서는 /home/k8s/hypercloud-sds/hcsctl/my_inventory를 설치에 사용한 인벤토리라고 가정)
	$ INVENTORY=/home/k8s/hypercloud-sds/hcsctl/my_inventory
	$ cat $INVENTORY/rook/cluster.yaml
	kind: CephCluster
	metadata:
	  name: rook-ceph
	  namespace: rook-ceph
	spec:
	  cephVersion:
		image: ceph/ceph:v14.2.9	## ceph/ceph:v15.2.4 로 바꿔줍니다.
		allowUnsupported: true
	  dataDirHostPath: /var/lib/rook
	...
	
	# 변경된 Ceph version으로 재적용합니다.
	$ kubectl apply -f $INVENTORY/rook/cluster.yaml
	cephcluster.ceph.rook.io/rook-ceph configured
	
	# 위에서 확인했던 ceph cluster Deployment들이 모두 ceph v15.2.4 version으로 업그레이드 되어 정상 작동하기까지 기다립니다.
	$ watch --exec kubectl -n rook-ceph get deployments -l rook_cluster=rook-ceph -o jsonpath='{range .items[*]}{.metadata.name}{"  \treq/upd/avl: "}{.spec.replicas}{"/"}{.status.updatedReplicas}{"/"}{.status.readyReplicas}{"  \tceph-version="}{.metadata.labels.ceph-version}{"\n"}{end}'
	
	# 아래는 업그레이드 중일 때 출력되는 예시입니다.
	rook-ceph-mds-myfs-a    req/upd/avl: 1/1/1      ceph-version=14.2.9-0
	rook-ceph-mds-myfs-b    req/upd/avl: 1/1/1      ceph-version=14.2.9-0
	rook-ceph-mgr-a         req/upd/avl: 1//        ceph-version=15.2.4-0
	rook-ceph-mon-a         req/upd/avl: 1/1/1      ceph-version=15.2.4-0
	rook-ceph-osd-0         req/upd/avl: 1/1/1      ceph-version=14.2.9-0
	rook-ceph-osd-1         req/upd/avl: 1/1/1      ceph-version=14.2.9-0
	rook-ceph-osd-2         req/upd/avl: 1/1/1      ceph-version=14.2.9-0
	```


### 업그레이드 완료 후의 image version
- Rook version
	- rook/ceph:v1.4.2
	- Rook Ceph Upgrade 항목의 6번에 Rook version 확인하는 방법이 있습니다.
- Ceph version
	- ceph/ceph:v15.2.4
	- Rook Ceph Upgrade 항목의 7번에 Ceph version 확인하는 방법이 있습니다.
- CSI version
	- quay.io/cephcsi/cephcsi:v3.1.0
	- quay.io/k8scsi/csi-attacher:v2.1.0
	- quay.io/k8scsi/csi-node-driver-registrar:v1.2.0
	- quay.io/k8scsi/csi-provisioner:v1.6.0
	- quay.io/k8scsi/csi-resizer:v0.4.0
	- quay.io/k8scsi/csi-snapshotter:v2.1.1
	- Rook Ceph Upgrade 항목의 6번에 CSI version 확인하는 방법이 있습니다.
