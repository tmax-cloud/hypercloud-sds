package rook

import (
	"bytes"
	"github.com/golang/glog"
	"hypercloud-storage/hcsctl/pkg/kubectl"
	"k8s.io/apimachinery/pkg/util/wait"
	"os"
	"path"
	"time"
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
	err = rookApply(inventoryPath, "cephfs-sc.yaml")
	if err != nil {
		return err
	}
	// TODO: ceph의 정상 상태를 기다리도록 변경
	time.Sleep(10 * time.Second)
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

func waitClusterApply() error {
	glog.Info("Wait for rook cluster apply")
	return wait.PollImmediate(time.Second, applyTimeout, isCreated)
}

func isCreated() (bool, error) {
	var stdout bytes.Buffer
	err := kubectl.Run(&stdout, os.Stderr, "get", "cephclusters.ceph.rook.io", "rook-ceph", "-n", "rook-ceph", "-o", "jsonpath={.status.state}")
	if err != nil {
		return false, err
	}
	return stdout.String() == "Created", nil
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
