package cdi

import (
	"bytes"
	"hypercloud-storage/hcsctl/pkg/kubectl"
	"hypercloud-storage/hcsctl/pkg/util"
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
	crPath := path.Join(inventoryPath, "cdi", "cr.yaml")

	cdiNamespace, err := util.GetSingleValueFromYaml(operatorPath,
		"Namespace", "metadata.name")
	if err != nil {
		return err
	}

	cdiOperatorDeploymentName, err := util.GetSingleValueFromYaml(operatorPath,
		"Deployment", "metadata.name")
	if err != nil {
		return err
	}

	crdName, err := util.GetSingleValueFromYaml(operatorPath,
		"CustomResourceDefinition", "metadata.name")
	if err != nil {
		return err
	}

	cdiKindName, err := util.GetSingleValueFromYaml(operatorPath,
		"CustomResourceDefinition", "spec.names.kind")
	if err != nil {
		return err
	}

	cdiName, err := util.GetSingleValueFromYaml(crPath, cdiKindName, "metadata.name")
	if err != nil {
		return err
	}

	err = kubectl.Run(os.Stdout, os.Stderr, "apply", "-f", operatorPath)
	if err != nil {
		return err
	}

	glog.Info("Wait for CDI operator to be created...")

	err = wait.PollImmediate(time.Second, applyTimeout,
		util.IsDeploymentCreated(cdiNamespace, cdiOperatorDeploymentName))
	if err != nil {
		return err
	}

	glog.Info("Wait for CDI CRD to be available...")

	err = wait.PollImmediate(time.Second, applyTimeout, util.IsCrdAvailable(crdName))
	if err != nil {
		return err
	}

	err = kubectl.Run(os.Stdout, os.Stderr, "apply", "-f", crPath)
	if err != nil {
		return err
	}

	glog.Info("Wait for CDI to be deployed...")

	err = wait.PollImmediate(time.Second, applyTimeout, util.IsCrDeployed(crdName, cdiName))
	if err != nil {
		return err
	}

	glog.Info("End CDI Apply")

	return nil
}

// Delete run `kubectl delete -f *.yaml`
func Delete(inventoryPath string) error {
	glog.Info("Start CDI Delete")

	crPath := path.Join(inventoryPath, "cdi", "cr.yaml")
	operatorPath := path.Join(inventoryPath, "cdi", "operator.yaml")

	var stderr bytes.Buffer
	err := kubectl.Run(os.Stdout, &stderr, "delete", "-f", crPath,
		"--ignore-not-found=true")

	if !kubectl.CRDAlreadyExists(stderr.String()) {
		glog.Infof("There isn't any remained custom resource already. Don't need to delete.")
	} else if err != nil {
		return err
	}

	glog.Info("Wait for cdi cr deleting")

	crdName, err := util.GetSingleValueFromYaml(operatorPath,
		"CustomResourceDefinition", "metadata.name")
	if err != nil {
		return err
	}

	cdiKindName, err := util.GetSingleValueFromYaml(operatorPath,
		"CustomResourceDefinition", "spec.names.kind")
	if err != nil {
		return err
	}

	cdiName, err := util.GetSingleValueFromYaml(crPath,
		cdiKindName, "metadata.name")
	if err != nil {
		return err
	}

	err = wait.PollImmediate(time.Second, deleteTimeout,
		util.IsCrDeleted(crdName, cdiName))
	if err != nil {
		return err
	}

	err = kubectl.Run(os.Stdout, os.Stderr, "delete", "-f",
		operatorPath, "--ignore-not-found=true")
	if err != nil {
		return err
	}

	glog.Info("End CDI Delete")

	return nil
}
