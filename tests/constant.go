package tests

import "time"

const (
	Pod2PodTestingNamespacePrefix = "test-pod-networking-"
	GoogleIPAddress               = "google.com"

	CdiTestingNamespacePrefix = "test-cdi-"
	DataVolumeNamePrefix      = "test-dv-"

	DataVolumeSize = "5Gi" // TODO 각 테스트별로 변경할 수 있도록

	TimeoutForCreatingPod = 60 * time.Second
	TimeoutForPing        = 20 * time.Second

	TimeOutForCreatingDv        = 500 * time.Second
	TimeoutForDeletingDv        = 300 * time.Second
	TimeOutForCreatingPvc       = 60 * time.Second
	TimeoutForDeletingNamespace = 500 * time.Second

	PollingIntervalForPing              = 5 * time.Second
	PollingIntervalDefault              = 3 * time.Second
	PollingIntervalForDeletingNamespace = 10 * time.Second

	// rook
	DeploymentRookCephOperator           = "rook-ceph-operator"
	DeploymentCsiCephfspluginProvisioner = "csi-cephfsplugin-provisioner"
	DeploymentCsiRbdpluginProvisioner    = "csi-rbdplugin-provisioner"

	ProvisionerCephfs = "rook-ceph.cephfs.csi.ceph.com"
	ProvisionerRbd    = "rook-ceph.rbd.csi.ceph.com"

	StorageClassCephfs = "csi-cephfs"

	//cdi
	DeploymentCdiOperator    = "cdi-operator"
	DeploymentCdiDeployment  = "cdi-deployment"
	DeploymentCdiApiserver   = "cdi-apiserver"
	DeploymentCdiUploadproxy = "cdi-uploadproxy"

	SampleRegistryURL = "docker://kubevirt/fedora-cloud-registry-disk-demo"
	SampleHTTPURL     = "https://download.cirros-cloud.net/contrib/0.3.0/cirros-0.3.0-i386-disk.img"
)
