# Containerized Data Importer (CDI)

> Containerized Data Importer (CDI) 는 pv 관리를 위한 k8s 의 add-on 으로써 kubevirt 로 vm 을 생성할 때, vm 에 mount 시킬 pvc 에 image 등의 data를 담아 생성할 수 있는 기능을 제공합니다.

## Reference

* [kubevirt/cdi github official page](https://github.com/kubevirt/containerized-data-importer)

## 사용 예

> [kubevirt/cdi github official example page](https://github.com/kubevirt/containerized-data-importer/tree/master/manifests/example) <br>
* 모든 example 은 해당 환경의 default storage class 를 사용하는 예로 구성되어있으며, 특정 storage class 사용을 원하는 경우, `spec.pvc.storageClassName` 에 원하는 `storage class name` 를 적어서 `kubectl create -f xxx.yaml` 하시면 됩니다.

  * [create datavolume import from http](./examples/datavolume-import-from-http.yaml)
  * [create datavolume import *.tar file from http](./examples/datavolume-import-from-http-archive-type.yaml)
  * [create datavolume import from registry](./examples/datavolume-import-from-registry-image.yaml)
  * [create datavolume import from registry with block-mode-pvc](./examples/datavolume-import-from-registry-image-block.yaml)
  * [create datavolume clone from pvc](./examples/datavolume-clone-from-pvc.yaml)

## configmap 변경 방법

* 추후 cdi 를 사용하여 dataVolume 을 생성할 때, private repository 로부터 특정 이미지(데이터)를 pull 하여 생성하고자 할 경우 다음과 같이 cdi 모듈 내부에서 사용하는 configmap 의 등록이 필요합니다.
  * `kubectl edit configmaps cdi-insecure-registries -n cdi` 명령어를 입력하여 yaml edit 창을 open 합니다.
  * 다음과 같이 apiVersion, kind, metadata 와 같은 indent 로 data.url 란에 registry 정보를 입력합니다.
    * 예) registry 주소가 192.168.1.1:5000 일 때

    ```{yaml}
    apiVersion: v1
    data:
      url: 192.168.1.1:5000
    kind: ConfigMap
    metadata:
      annotations:
        operator.cdi.kubevirt.io/lastAppliedConfiguration: `xxx`
      creationTimestamp: "2019-12-20T04:05:12Z"
      labels:
        cdi.kubevirt.io: ""
        operator.cdi.kubevirt.io/createVersion: v1.11.0
      name: cdi-insecure-registries
      namespace: cdi
      ownerReferences:
      - apiVersion: cdi.kubevirt.io/v1alpha1
        blockOwnerDeletion: true
        controller: true
        kind: CDI
        name: cdi
        uid: 2655e870-69d7-4b3e-bfe1-d4eab78cbca5
      resourceVersion: "46553577"
      selfLink: /api/v1/namespaces/cdi/configmaps/di-insecure-registries
      uid: aebe00c9-5bfc-4911-8170-25e23b144e29
    ```

## cdiconfig 변경 (주의: cdiconfig 와 configmap 은 다른 resource입니다.)

> cdi 모듈 내부적으로 임시로 create 하는 pod 과 pvc 가 존재하는데, <br>
>해당 pvc 는 CDIConfig 의 status.scratchSpaceStorageClass 에 적힌 storageClass 로 provisioning 하여 생성됩니다. <br>
>이 값은 `cdi-cr.yaml` 을 apply 할 당시의 k8s cluster 환경의 default storageClassName 을 fetch 하여 입력되며, <br>
>추후 이를 변경하고자 할 때는 다음과 같은 방식으로 변경 가능합니다.

* `kubectl edit cdiconfig {$CDIConfigName}` 명령으로 CDIconfig 의

```
spec: {}
```

 으로 적힌 부분을

````
spec:
  scratchSpaceStorageClass: {$변경하길 원하는 storageClassName}
````

으로 변경하여 CDIConfig 의 status.scratchSpaceStorageClass 를 변경할 수 있습니다.