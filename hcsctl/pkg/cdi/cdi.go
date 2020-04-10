package cdi

import (
	"bytes"
	"hypercloud-storage/hcsctl/pkg/kubectl"
	"os"
	"path"
	"time"

	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/util/wait"
)

var (
	applyTimeout  = 300 * time.Second * 2
	deleteTimeout = 300 * time.Second * 2
)

// Apply run `kubectl apply -f *.yaml`
func Apply(inventoryPath string) error {
	glog.Info("Start CDI Apply")

	operatorPath := path.Join(inventoryPath, "cdi", "operator.yaml")

	err := kubectl.Run(os.Stdout, os.Stderr, "apply", "-f", operatorPath)
	if err != nil {
		return err
	}

	// TODO: waiting util operator is ready (김현빈 연구원이 구현하고 있습니다.)
	tmp := 10
	time.Sleep(time.Duration(tmp) * time.Second)

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

	var stderr bytes.Buffer
	err := kubectl.Run(os.Stdout, &stderr, "delete", "-f", crPath, "--ignore-not-found=true")

	if !kubectl.CRDAlreadyExists(stderr.String()) {
		glog.Infof("There isn't any remained custom resource already. Don't need to delete.")
	} else if err != nil {
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
	var stdout, stderr bytes.Buffer
	err := kubectl.Run(&stdout, &stderr, "get", "cdis.cdi.kubevirt.io", "cdi", "-o", "json", "--ignore-not-found=true")

	if !kubectl.CRDAlreadyExists(stderr.String()) {
		glog.Infof("There isn't any remained custom resource already. Don't need to delete.")
	} else if err != nil {
		return false, err
	}

	return stdout.String() == "", nil
}
