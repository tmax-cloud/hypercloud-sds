# Containerized Data Importer (CDI)

> Containerized Data Importer (CDI) 는 pv 관리를 위한 k8s 의 add-on 으로써 kubevirt 로 vm 을 생성할 때, vm 에 mount 시킬 pvc 에 image 등의 data를 담아 생성할 수 있는 기능을 제공합니다. <br>

- **CDI 모듈은 dynamic provisioning 이 가능한 storageClass 를 Prerequisite 으로 필요로 합니다.**

## Reference

* [kubevirt/cdi github official page](https://github.com/kubevirt/containerized-data-importer)

## 사용 예 

### DataVolumes

- CDI 모듈은 Operator 패턴으로 배포되어, DataVolume 등의 CustomResource 에 대하여 Create, Delete, Get, Update 등의 작업을 가능하게 합니다.
- DataVolume 은 kubernetes 의 resource 인 PVC 를 추상화한 개념으로 pvc 에 data 를 import 하여 [kubevirt](https://github.com/kubevirt/kubevirt) 와의 통합을 쉽게 합니다.

### Examples
> [kubevirt/cdi github official example page](https://github.com/kubevirt/containerized-data-importer/tree/master/manifests/example) <br>
* 소개된 모든 example 은 해당 환경의 default storage class 를 사용하는 예로 구성되어있으며, 특정 storage class 사용을 원하는 경우, `spec.pvc.storageClassName` 에 원하는 `storage class name` 를 적어서 `kubectl create -f xxx.yaml` 하시면 됩니다.

  * [create datavolume import from http](./examples/datavolume-import-from-http.yaml)
  * [create datavolume import *.tar file from http](./examples/datavolume-import-from-http-archive-type.yaml)
  * [create datavolume import from registry](./examples/datavolume-import-from-registry.yaml)
  * [create datavolume import from registry with block-mode-pvc](./examples/datavolume-import-from-registry-block.yaml)
  * [create datavolume clone from pvc](./examples/datavolume-clone-from-pvc.yaml)

#### DataVolume 생성
- yaml 파일에 적은 spec 으로 datavolume 을 생성합니다.
  - `kubectl create -f xxx.yaml -n {$namespace}`

#### DataVolume 조회
- namespace 에 존재하는 모든 datavolume 을 조회합니다.
  - `kubectl get datavolume -n {$namespace}`
- namespace 에 존재하는 datavolume 중 {$datavolumename} 을 자세히 조회합니다.
  - `kubectl describe datavolume {$datavolumename} -n {$namespace}`
  
#### DataVolume 삭제
- namespace 에 존재하는 datavolume 중 {$datavolumename} 을 삭제합니다.
  - `kubectl delete datavolume {$datavolumename} -n {$namespace}`

## Appendix

### configmap 변경 방법

* 추후 cdi 를 사용하여 dataVolume 을 생성할 때, private repository 로부터 특정 이미지(데이터)를 pull 하여 생성하고자 할 경우 다음과 같이 cdi 모듈 내부에서 사용하는 configmap 에 사용하고자 하는 registry 의 url 등록이 필요합니다.
  * `kubectl patch configmap cdi-insecure-registries -n cdi --type merge -p '{"data":{"mykey": "my-private-registry-host:5000"}}'` 명령어를 입력하여 cdi-insecure-registries 이름을 가진 configmap 에 data 를 추가하는 방식을 통하여 registry 정보를 등록 할 수 있습니다.
  * `{"mykey": "my-private-registry-host:5000"}`를 원하는 key 와 사용하고자 하는 registry url 로 변경하면 되고 위의 명령어를 통하여 여러개의 registry url 등록이 가능합니다. 
    * ex) `kubectl patch configmap cdi-insecure-registries -n cdi --type merge -p '{"data":{"url": "192.168.1.1:5000"}}'`
  * 다음과 같이 `kubectl describe configmap -n cdi cdi-insecure-registries` 명령어를 입력하여 Data 하위에 registry 정보가 추가 된 것을 확인 할 수 있습니다.

    ```{yaml}
    Name:         cdi-insecure-registries
    Namespace:    cdi
    Labels:       cdi.kubevirt.io=
                  operator.cdi.kubevirt.io/createVersion=v1.11.0
    Annotations:  operator.cdi.kubevirt.io/lastAppliedConfiguration:
                    {"kind":"ConfigMap","apiVersion":"v1","metadata":{"name":"cdi-insecure-registries","namespace":"cdi","creationTimestamp":null,"labels":{"c...
        
    Data
    ====
    url:
    ----
    192.168.1.1:5000
    url2:
    ----
    192.168.1.2:5000
    Events:  <none>
    ```

### cdiconfig 변경 방법
> 주의: cdiconfig 와 configmap 은 **다른** resource입니다.

* cdi 모듈 내부적으로 사용하는 storageClass 를 ScratchSpaceStorageClass 라고 하는데, DataVolume 생성 요청 시 해당 sc 로 pvc 생성 및 삭제를 진행합니다. 
  * 이 값은 최초에 `cdi-cr.yaml` 을 apply 할 당시의 k8s cluster 환경의 default storageClassName 을 fetch 하여 입력되며, 추후 이를 변경하고자 할 때는 다음과 같은 방식으로 변경 가능합니다.
    * 이 값이 설정되어 있지 않을 경우, cluster 에 default storageClass 를 설정해주면 cdi 모듈은 자동으로 해당 default storageClass 를 사용하게 됩니다.
    * 이 값이 설정되어 있지 않고 default storageClass 도 설정되어 있지 않으면 datavolume 생성 시 입력한 `spec.pvc.storageClassName` 값을 사용하게 됩니다.
  * 해당 storage class 는 반드시 **dynamic provisioning** 을 지원하는 storageClass 여야 합니다.

* `kubectl patch cdiconfig {$CDIConfigName} --type merge -p '{"spec":{"scratchSpaceStorageClass": "{$storageClassName}"}}'` 명령어를 입력하여 scratchSpaceStoragClass 값을 수정 할 수 있습니다.
  * `{$CDIConfigName}` 는 cdiconfig 이름으로 변경하고, `{$storageClassName}` 는 수정하고자 하는 storageClass 이름으로 변경하면 됩니다.
    * ex) `kubectl patch cdiconfig config --type merge -p '{"spec":{"scratchSpaceStorageClass": "ceph-block-sc"}}'`
