package rook

import (
	"bytes"
	"hypercloud-storage/hcsctl/pkg/kubectl"
	"os"
	"path"
	"time"
	"strings"

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

	err := rookApply(inventoryPath, "common.yaml")
	if err != nil {
		return err
	}
	err = rookApply(inventoryPath, "operator.yaml")
	if err != nil {
		return err
	}
	err = waitRookOperator() 
	if err != nil {
		return nil
	}

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
	path := path.Join(inventoryPath, "rook", filename)
	return kubectl.Run(os.Stdout, os.Stderr, "apply", "-f", path)
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

	err := kubectl.Run(&stdout, os.Stderr, "get", "cephclusters.ceph.rook.io", "rook-ceph", "-n", "rook-ceph", "-o", "jsonpath={.status.ceph.health}")
	if err != nil {
		return false, err
	}

	if stdout.String() == "HEALTH_OK" {
		osdDeployments := getDaemonDeployNames("ceph-osd")
	
		/// TODO: Count osd number in cluster.yaml
		return isDaemonReadyAndAvailable(osdDeployments)	
	} 
	
	return false, nil
}

func isMdsCreated() (bool, error) {
	mdsDeployments := getDaemonDeployNames("ceph-mds")

	// TODO: Count mds number in cephfs-fs.yaml
	if len(mdsDeployments) >= 2 {
		return isDaemonReadyAndAvailable(mdsDeployments)
	}

	return false, nil
}

func getDaemonDeployNames(name string) []string {
	var stdout bytes.Buffer
	kubectl.Run(&stdout, os.Stderr, "get", "deployments.apps", "-n", "rook-ceph", "-o", "custom-columns=name:.metadata.name", "--no-headers")

	var targetDeployments [] string

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
		if err != nil {	return false, err }
		err = kubectl.Run(&readyReplicaCount, os.Stderr, "get", "deployments.apps",
							"-n", "rook-ceph", daemon,
							"-o", "jsonpath='{.status.readyReplicas}'")
		if err != nil {	return false, err }
		err = kubectl.Run(&availReplicaCount, os.Stderr, "get", "deployments.apps",
							"-n", "rook-ceph", daemon,
							"-o", "jsonpath='{.status.availableReplicas}'")
		if err != nil {	return false, err }

		// If a replica is not ready or is not available, polling should keep going
		if !(replicaCount.String() == readyReplicaCount.String() && replicaCount.String() == availReplicaCount.String()) {
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
	path := path.Join(inventoryPath, "rook", filename)
	return kubectl.Run(os.Stdout, os.Stderr, "delete", "-f", path, "--ignore-not-found=true", "--wait=true")
}

func waitClusterDelete() error {
	glog.Info("Wait for rook cluster delete")
	return wait.PollImmediate(time.Second, applyTimeout, isDeleted)
}

func isDeleted() (bool, error) {
	var stdout bytes.Buffer
	err := kubectl.Run(&stdout, os.Stderr, "get", "cephclusters.ceph.rook.io", "rook-ceph", "-n", "rook-ceph", "-o", "json", "--ignore-not-found=true")
	if err != nil {
		return false, err
	}
	return len(stdout.String()) == 0, nil
}
