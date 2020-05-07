package rook

import (
	"bytes"
	"errors"
	"fmt"
	"hypercloud-storage/hcsctl/pkg/kubectl"
	"hypercloud-storage/hcsctl/pkg/util"

	"os"
	"path"
	"strings"
	"time"

	"github.com/golang/glog"

	"k8s.io/apimachinery/pkg/util/wait"

	"k8s.io/apimachinery/pkg/util/sets"
)

const (
	applyTimeout  = 20 * time.Minute
	deleteTimeout = 20 * time.Minute
)

var (
	// CommonYaml represents common.yaml
	CommonYaml string = "common.yaml"
	// OperatorYaml represents operator.yaml
	OperatorYaml string = "operator.yaml"
	// ClusterYaml represents cluster.yaml
	ClusterYaml string = "cluster.yaml"
	// RbdStorageClassYaml represents rbd-sc.yaml
	RbdStorageClassYaml string = "rbd-sc.yaml"
	// CephfsFilesystemYaml represents cephfs-fs.yaml
	CephfsFilesystemYaml string = "cephfs-fs.yaml"
	// CephfsStorageClassYaml represents cephfs-sc.yaml
	CephfsStorageClassYaml string = "cephfs-sc.yaml"
	// ToolboxYaml represents toolbox.yaml
	ToolboxYaml string = "toolbox.yaml"
	// RookYamlSet represents required yamls of rook
	RookYamlSet = sets.NewString(CommonYaml, OperatorYaml, ClusterYaml, RbdStorageClassYaml,
		CephfsFilesystemYaml, CephfsStorageClassYaml, ToolboxYaml)
)

var (
	rookCephNamespaceName, cephclusterCrdName, cephclusterKindName, cephclusterCrName string

	// TODO 함수마다 parameter로 들고 다닐 건지 or 그냥 package variable 로 선언해놓고 공유할 건지
	_inventoryPath string
)

// Apply executes `kubectl apply -f *.yaml`
func Apply(inventoryPath string) error {
	_inventoryPath = inventoryPath

	glog.Info("[STEP 0 / 6] Start Applying Rook-ceph")

	glog.Info("[STEP 1 / 6] Fetch Rook-ceph variables from inventory")

	err := setRookCephValuesFrom(inventoryPath)
	if err != nil {
		return err
	}

	err = rookApply(inventoryPath, CommonYaml)
	if err != nil {
		return err
	}

	err = rookApply(inventoryPath, OperatorYaml)
	if err != nil {
		return err
	}

	glog.Infof("[STEP 2 / 6] Wait up to %s for Rook-ceph operator to be created...", applyTimeout.String())

	err = waitRookOperator()
	if err != nil {
		return err
	}

	glog.Infof("[STEP 3 / 6] Wait up to %s for Cephcluster CRD to be available...", applyTimeout.String())

	err = wait.PollImmediate(time.Second, applyTimeout, util.IsCrdAvailable(cephclusterCrdName))
	if err != nil {
		return err
	}

	err = rookApply(inventoryPath, ClusterYaml)
	if err != nil {
		return err
	}

	glog.Infof("[STEP 4 / 6] Wait up to %s for CephCluster applied...", applyTimeout.String())

	err = waitClusterApply()
	if err != nil {
		return err
	}

	err = rookApply(inventoryPath, RbdStorageClassYaml)
	if err != nil {
		return err
	}

	err = rookApply(inventoryPath, CephfsFilesystemYaml)
	if err != nil {
		return err
	}

	glog.Infof("[STEP 5 / 6] Wait up to %s for CephFS created...", applyTimeout.String())

	err = waitCephFSReady()
	if err != nil {
		return err
	}

	err = rookApply(inventoryPath, CephfsStorageClassYaml)
	if err != nil {
		return err
	}

	err = rookApply(inventoryPath, ToolboxYaml)
	if err != nil {
		return err
	}

	glog.Info("[STEP 6 / 6] End Applying Rook-ceph")

	return nil
}

func rookApply(inventoryPath, filename string) error {
	yamlPath := path.Join(inventoryPath, "rook", filename)
	return kubectl.Run(os.Stdout, os.Stderr, "apply", "-f", yamlPath)
}

func waitRookOperator() error {
	return wait.PollImmediate(time.Second, applyTimeout, isOperatorCreated)
}

func waitClusterApply() error {
	return wait.PollImmediate(time.Second, applyTimeout, isClusterCreated)
}

func waitCephFSReady() error {
	return wait.PollImmediate(time.Second, applyTimeout, isMdsCreated)
}

// TODO inventoryPath 를 parameter 로 받도록
func isOperatorCreated() (bool, error) {
	operatorDeployments, err := getDaemonDeployNames("ceph-operator")
	if err != nil {
		return false, err
	}

	operatorReplicaNumNoType, err := util.GetValueFromYamlFile(path.Join(_inventoryPath, "rook", OperatorYaml),
		util.Deployment, "spec.replicas")
	if err != nil {
		return false, err
	}

	expectedOperatorReplicaNum, isConvertibleToInt := operatorReplicaNumNoType[0].(int)
	if !isConvertibleToInt {
		return false, errors.New("Unable to convert value of " +
			"spec.replicas" + " to int: " + fmt.Sprintf("%v", operatorReplicaNumNoType[0]))
	}

	if len(operatorDeployments) == expectedOperatorReplicaNum {
		return isDaemonReadyAndAvailable(operatorDeployments)
	}

	return false, nil
}

func isClusterCreated() (bool, error) {
	var stdout bytes.Buffer

	err := kubectl.Run(&stdout, os.Stderr, "get", cephclusterCrdName,
		cephclusterCrName, "-n", rookCephNamespaceName, "-o", "jsonpath={.status.ceph.health}")
	if err != nil {
		return false, err
	}

	if stdout.String() == "HEALTH_OK" {
		osdDeployments, err := getDaemonDeployNames("ceph-osd")

		if err != nil {
			return false, err
		}

		// TODO: Count osd number in cluster.yaml
		return isDaemonReadyAndAvailable(osdDeployments)
	}

	return false, nil
}

// TODO inventoryPath 를 parameter 로 받도록
func isMdsCreated() (bool, error) {
	mdsDeployments, err := getDaemonDeployNames("ceph-mds")
	if err != nil {
		return false, err
	}

	mdsActiveNumNoType, err := util.GetValueFromYamlFile(path.Join(_inventoryPath, "rook", CephfsFilesystemYaml),
		"CephFilesystem", "spec.metadataServer.activeCount")
	if err != nil {
		return false, err
	}

	expectedMdsActiveNum, isConvertibleToInt := mdsActiveNumNoType[0].(int)
	if !isConvertibleToInt {
		return false, errors.New("Unable to convert value of " +
			"spec.metadataServer.activeCount" + " to int: " + fmt.Sprintf("%v", mdsActiveNumNoType[0]))
	}

	// since rook-ceph's policy
	expectedMdsActiveNum *= 2

	if len(mdsDeployments) == expectedMdsActiveNum {
		return isDaemonReadyAndAvailable(mdsDeployments)
	}

	return false, nil
}

func getDaemonDeployNames(name string) ([]string, error) {
	var stdout bytes.Buffer

	err := kubectl.Run(&stdout, os.Stderr, "get", "deployments.apps", "-n", rookCephNamespaceName,
		"-o", "custom-columns=name:.metadata.name", "--no-headers")
	if err != nil {
		return nil, err
	}

	var targetDeployments []string

	for _, deployment := range strings.Split(stdout.String(), "\n") {
		if strings.Contains(deployment, name) {
			targetDeployments = append(targetDeployments, deployment)
		}
	}

	return targetDeployments, nil
}

// Check daemon pods are all ready and available
func isDaemonReadyAndAvailable(daemons []string) (bool, error) {
	for _, daemon := range daemons {
		var replicaCount, readyReplicaCount, availReplicaCount bytes.Buffer

		err := kubectl.Run(&replicaCount, os.Stderr, "get", "deployments.apps",
			"-n", rookCephNamespaceName, daemon,
			"-o", "jsonpath='{.status.replicas}'")
		if err != nil {
			return false, err
		}

		err = kubectl.Run(&readyReplicaCount, os.Stderr, "get", "deployments.apps",
			"-n", rookCephNamespaceName, daemon,
			"-o", "jsonpath='{.status.readyReplicas}'")
		if err != nil {
			return false, err
		}

		err = kubectl.Run(&availReplicaCount, os.Stderr, "get", "deployments.apps",
			"-n", rookCephNamespaceName, daemon,
			"-o", "jsonpath='{.status.availableReplicas}'")
		if err != nil {
			return false, err
		}

		// If a replica is not ready or is not available, polling should keep going
		if replicaCount.String() != readyReplicaCount.String() ||
			replicaCount.String() != availReplicaCount.String() {
			return false, nil
		}
	}

	return true, nil
}

// Delete executes `kubectl delete -f *.yaml`
func Delete(inventoryPath string) error {
	glog.Info("[STEP 0 / 4] Start Deleting Rook")

	glog.Info("[STEP 1 / 4] Fetch Rook-ceph variables from inventory")

	err := setRookCephValuesFrom(inventoryPath)
	if err != nil {
		return err
	}

	err = rookDelete(inventoryPath, ToolboxYaml)
	if err != nil {
		return err
	}

	err = rookDelete(inventoryPath, CephfsStorageClassYaml)
	if err != nil {
		return err
	}

	err = rookDelete(inventoryPath, CephfsFilesystemYaml)
	if err != nil {
		return err
	}

	err = rookDelete(inventoryPath, RbdStorageClassYaml)
	if err != nil {
		return err
	}

	err = rookDelete(inventoryPath, ClusterYaml)
	if err != nil {
		return err
	}

	glog.Infof("[STEP 2 / 4] Wait up to %s for rook cluster to be deleted...", deleteTimeout.String())

	err = waitClusterDelete()
	if err != nil {
		return err
	}

	glog.Infof("[STEP 3 / 4] Wait up to %s for rook operator to be deleted...", deleteTimeout.String())

	err = rookDelete(inventoryPath, OperatorYaml)
	if err != nil {
		return err
	}

	err = rookDelete(inventoryPath, CommonYaml)
	if err != nil {
		return err
	}

	glog.Info("[STEP 4 / 4] End Deleting Rook")
	glog.Info("[WARNING] You need to remove /var/lib/rook directory in every nodes. " +
		"Also, you need to reset all devices used by rook-ceph. " +
		"There is the reset manual in the hypercloud-storage gitlab project")

	return nil
}

func rookDelete(inventoryPath, yamlName string) error {
	yamlPath := path.Join(inventoryPath, "rook", yamlName)

	var stderr bytes.Buffer
	err := kubectl.Run(os.Stdout, &stderr, "delete", "-f", yamlPath, "--ignore-not-found=true", "--wait=true")

	if err != nil && kubectl.CRDAlreadyExists(stderr.String()) {
		return err
	}

	return nil
}

func waitClusterDelete() error {
	return wait.PollImmediate(time.Second, deleteTimeout, isDeleted)
}

func isDeleted() (bool, error) {
	var stdout, stderr bytes.Buffer
	err := kubectl.Run(&stdout, &stderr, "get", cephclusterCrdName,
		cephclusterCrName, "-n", rookCephNamespaceName, "-o", "json", "--ignore-not-found=true")

	if err != nil && kubectl.CRDAlreadyExists(stderr.String()) {
		return false, err
	}

	return stdout.String() == "", nil
}

func setRookCephValuesFrom(inventoryPath string) error {
	commonPath := path.Join(inventoryPath, "rook", CommonYaml)
	clusterPath := path.Join(inventoryPath, "rook", ClusterYaml)

	var err error

	rookCephNamespaceName, err = util.GetUniqueStringValueFromYamlFile(commonPath,
		util.Namespace, "metadata.name")
	if err != nil {
		return err
	}

	// TODO should be set by fetching from certain yaml file
	cephclusterCrdName = "cephclusters.ceph.rook.io"
	cephclusterKindName = "CephCluster"

	cephclusterCrName, err = util.GetUniqueStringValueFromYamlFile(clusterPath,
		cephclusterKindName, "metadata.name")
	if err != nil {
		return err
	}

	return nil
}
