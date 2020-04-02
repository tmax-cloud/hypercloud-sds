package cdi

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
	applyTimeout  = 300 * time.Second *2
	deleteTimeout = 300 * time.Second *2
)

// Apply run `kubectl apply -f *.yaml`
func Apply(inventoryPath string) error {
	glog.Info("Start CDI Apply")

	operatorPath := path.Join(inventoryPath, "cdi", "operator.yaml")
	err := kubectl.Run(os.Stdout, os.Stderr, "apply", "-f", operatorPath)
	if err != nil {
		return err
	}

	crPath := path.Join(inventoryPath, "cdi", "cr.yaml")
	err = kubectl.Run(os.Stdout, os.Stderr, "apply", "-f", crPath)
	if err != nil {
		return err
	}
	glog.Info("Wait for cdi deploy state")
	err = wait.PollImmediate(time.Second, applyTimeout, isDeployed)
	if err != nil {
		return err
	}

	glog.Info("End CDI Apply")
	return nil
}

func isDeployed() (bool, error) {
	var stdout bytes.Buffer
	err := kubectl.Run(&stdout, os.Stderr, "get", "cdis.cdi.kubevirt.io", "cdi", "-o", "jsonpath={.status.phase}")
	if err != nil {
		return false, err
	}
	return stdout.String() == "Deployed", nil
}

// Delete run `kubectl delete -f *.yaml`
func Delete(inventoryPath string) error {
	glog.Info("Start CDI Delete")

	crPath := path.Join(inventoryPath, "cdi", "cr.yaml")
	err := kubectl.Run(os.Stdout, os.Stderr, "delete", "-f", crPath, "--ignore-not-found=true")
	if err != nil {
		return err
	}
	glog.Info("Wait for cdi cr deleting")
	err = wait.PollImmediate(time.Second, deleteTimeout, isDeleted)
	if err != nil {
		return err
	}

	operatorPath := path.Join(inventoryPath, "cdi", "operator.yaml")
	err = kubectl.Run(os.Stdout, os.Stderr, "delete", "-f", operatorPath, "--ignore-not-found=true")
	if err != nil {
		return err
	}

	glog.Info("End CDI Delete")
	return nil
}

func isDeleted() (bool, error) {
	var stdout bytes.Buffer
	err := kubectl.Run(&stdout, os.Stderr, "get", "cdis.cdi.kubevirt.io", "cdi", "-o", "json", "--ignore-not-found=true")
	if err != nil {
		return false, err
	}
	return len(stdout.String()) == 0, nil
}
