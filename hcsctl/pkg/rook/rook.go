package rook

import (
	"bytes"
	"hypercloud-storage/hcsctl/pkg/kubectl"
	"os"
	"path"
	"strings"
	"time"

	"github.com/golang/glog"

	"k8s.io/apimachinery/pkg/util/wait"
)

var (
	applyTimeout  = 600 * time.Second * 2
	deleteTimeout = 600 * time.Second * 2
)

// Apply run `kubectl apply -f *.yaml`
func Apply(inventoryPath string) error {
	glog.Info("Start Rook Apply")
	// TODO 추가적으로 필요한 rook 관련 yaml 파일 추가
	// TODO yaml 파일명 및 목록 고정해서 상수로 관리
	err := rookApply(inventoryPath, "common.yaml")
	if err != nil {
		return err
	}

	err = rookApply(inventoryPath, "operator.yaml")
	if err != nil {
		return err
	}

	glog.Info("Before operator:", time.Now())

	err = waitRookOperator()
	if err != nil {
		return err
	}

	glog.Info("After operator:", time.Now())

	// TODO
	tmp := 10
	time.Sleep(time.Duration(tmp) * time.Second)

	err = rookApply(inventoryPath, "cluster.yaml")
	if err != nil {
		return err
	}

	err = waitClusterApply()
	if err != nil {
		return err
	}

	err = rookApply(inventoryPath, "rbd-sc.yaml")
	if err != nil {
		return err
	}

	err = rookApply(inventoryPath, "cephfs-fs.yaml")
	if err != nil {
		return err
	}

	err = waitCephFSReady()
	if err != nil {
		return err
	}

	err = rookApply(inventoryPath, "cephfs-sc.yaml")
	if err != nil {
		return err
	}

	err = rookApply(inventoryPath, "toolbox.yaml")
	if err != nil {
		return err
	}

	glog.Info("End Rook Apply")

	return nil
}

func rookApply(inventoryPath, filename string) error {
	yamlPath := path.Join(inventoryPath, "rook", filename)
	return kubectl.Run(os.Stdout, os.Stderr, "apply", "-f", yamlPath)
}

func waitRookOperator() error {
	glog.Info("Wait for rook operator created")
	return wait.PollImmediate(time.Second, applyTimeout, isOperatorCreated)
}

func waitClusterApply() error {
	glog.Info("Wait for CephCluster applied")
	return wait.PollImmediate(time.Second, applyTimeout, isClusterCreated)
}

func waitCephFSReady() error {
	glog.Info("Wait for CephFS created")
	return wait.PollImmediate(time.Second, applyTimeout, isMdsCreated)
}

func isOperatorCreated() (bool, error) {
	operatorDeployments := getDaemonDeployNames("ceph-operator")

	if len(operatorDeployments) == 1 {
		return isDaemonReadyAndAvailable(operatorDeployments)
	}

	return false, nil
}

func isClusterCreated() (bool, error) {
	var stdout bytes.Buffer

	err := kubectl.Run(&stdout, os.Stderr, "get", "cephclusters.ceph.rook.io",
		"rook-ceph", "-n", "rook-ceph", "-o", "jsonpath={.status.ceph.health}")
	if err != nil {
		return false, err
	}

	if stdout.String() == "HEALTH_OK" {
		osdDeployments := getDaemonDeployNames("ceph-osd")

		// TODO: Count osd number in cluster.yaml
		return isDaemonReadyAndAvailable(osdDeployments)
	}

	return false, nil
}

func isMdsCreated() (bool, error) {
	mdsDeployments := getDaemonDeployNames("ceph-mds")

	// TODO: Count mds number in cephfs-fs.yaml
	mdsCount := 2
	if len(mdsDeployments) >= mdsCount {
		return isDaemonReadyAndAvailable(mdsDeployments)
	}

	return false, nil
}

// TODO: Should has return error (ErrorHandling)
func getDaemonDeployNames(name string) []string {
	var stdout bytes.Buffer

	err := kubectl.Run(&stdout, os.Stderr, "get", "deployments.apps", "-n", "rook-ceph",
		"-o", "custom-columns=name:.metadata.name", "--no-headers")
	if err != nil {
		glog.Error(err)
	}

	var targetDeployments []string

	for _, deployment := range strings.Split(stdout.String(), "\n") {
		if strings.Contains(deployment, name) {
			targetDeployments = append(targetDeployments, deployment)
		}
	}

	return targetDeployments
}

// Check daemon pods are all ready and available
func isDaemonReadyAndAvailable(daemons []string) (bool, error) {
	allComplete := true

	for _, daemon := range daemons {
		var replicaCount, readyReplicaCount, availReplicaCount bytes.Buffer

		err := kubectl.Run(&replicaCount, os.Stderr, "get", "deployments.apps",
			"-n", "rook-ceph", daemon,
			"-o", "jsonpath='{.status.replicas}'")
		if err != nil {
			return false, err
		}

		err = kubectl.Run(&readyReplicaCount, os.Stderr, "get", "deployments.apps",
			"-n", "rook-ceph", daemon,
			"-o", "jsonpath='{.status.readyReplicas}'")
		if err != nil {
			return false, err
		}

		err = kubectl.Run(&availReplicaCount, os.Stderr, "get", "deployments.apps",
			"-n", "rook-ceph", daemon,
			"-o", "jsonpath='{.status.availableReplicas}'")
		if err != nil {
			return false, err
		}

		// If a replica is not ready or is not available, polling should keep going
		if replicaCount.String() != readyReplicaCount.String() ||
			replicaCount.String() != availReplicaCount.String() {
			allComplete = false
			break
		}
	}

	return allComplete, nil
}

// Delete run `kubectl delete -f *.yaml`
func Delete(inventoryPath string) error {
	glog.Info("Start Rook Delete")

	err := rookDelete(inventoryPath, "toolbox.yaml")
	if err != nil {
		return err
	}

	err = rookDelete(inventoryPath, "cephfs-sc.yaml")
	if err != nil {
		return err
	}

	err = rookDelete(inventoryPath, "cephfs-fs.yaml")
	if err != nil {
		return err
	}

	err = rookDelete(inventoryPath, "rbd-sc.yaml")
	if err != nil {
		return err
	}

	err = rookDelete(inventoryPath, "cluster.yaml")
	if err != nil {
		return err
	}

	err = waitClusterDelete()
	if err != nil {
		return err
	}

	err = rookDelete(inventoryPath, "operator.yaml")
	if err != nil {
		return err
	}

	err = rookDelete(inventoryPath, "common.yaml")
	if err != nil {
		return err
	}

	glog.Info("End Rook Delete")

	return nil
}

func rookDelete(inventoryPath, filename string) error {
	yamlPath := path.Join(inventoryPath, "rook", filename)

	var stderr bytes.Buffer
	err := kubectl.Run(os.Stdout, &stderr, "delete", "-f", yamlPath, "--ignore-not-found=true", "--wait=true")

	if !kubectl.CRDAlreadyExists(stderr.String()) {
		glog.Infof("There isn't any remained custom resource already. Don't need to delete.")
	} else if err != nil {
		return err
	}

	return nil
}

func waitClusterDelete() error {
	glog.Info("Wait for rook cluster delete")
	return wait.PollImmediate(time.Second, deleteTimeout, isDeleted)
}

func isDeleted() (bool, error) {
	var stdout, stderr bytes.Buffer
	err := kubectl.Run(&stdout, &stderr, "get", "cephclusters.ceph.rook.io",
		"rook-ceph", "-n", "rook-ceph", "-o", "json", "--ignore-not-found=true")

	if !kubectl.CRDAlreadyExists(stderr.String()) {
		glog.Infof("There isn't any remained custom resource already. Don't need to delete.")
	} else if err != nil {
		return false, err
	}

	return stdout.String() == "", nil
}
