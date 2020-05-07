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

	"k8s.io/apimachinery/pkg/util/sets"
)

const (
	applyTimeout  = 10 * time.Minute
	deleteTimeout = 10 * time.Minute
)

var (
	// OperatorYaml represents operator.yaml
	OperatorYaml string = "operator.yaml"
	// CrYaml represents cr.yaml
	CrYaml string = "cr.yaml"
	// CdiYamlSet represents required yamls of cdi
	CdiYamlSet = sets.NewString(OperatorYaml, CrYaml)
)

var (
	cdiNamespaceName, cdiOperatorDeploymentName, cdiCrdName, cdiKindName, cdiCrName string

	operatorPath string
	crPath       string
)

// Apply executes `kubectl apply -f *.yaml`
func Apply(inventoryPath string) error {
	glog.Info("[STEP 0 / 5] Start Applying CDI")

	glog.Info("[STEP 1 / 5] Fetch CDI variables from inventory")

	err := setCdiValuesFrom(inventoryPath)
	if err != nil {
		return err
	}

	err = kubectl.Run(os.Stdout, os.Stderr, "apply", "-f", operatorPath)
	if err != nil {
		return err
	}

	glog.Infof("[STEP 2 / 5] Wait up to %s for CDI operator to be created...", applyTimeout.String())

	err = wait.PollImmediate(time.Second, applyTimeout,
		util.IsDeploymentCreated(cdiNamespaceName, cdiOperatorDeploymentName))
	if err != nil {
		return err
	}

	glog.Infof("[STEP 3 / 5] Wait up to %s for CDI CRD to be available...", applyTimeout.String())

	err = wait.PollImmediate(time.Second, applyTimeout, util.IsCrdAvailable(cdiCrdName))
	if err != nil {
		return err
	}

	err = kubectl.Run(os.Stdout, os.Stderr, "apply", "-f", crPath)
	if err != nil {
		return err
	}

	glog.Infof("[STEP 4 / 5] Wait up to %s for CDI CR to be deployed...", applyTimeout.String())

	err = wait.PollImmediate(time.Second, applyTimeout, util.IsCrDeployed(cdiCrdName, cdiCrName))
	if err != nil {
		return err
	}

	glog.Info("[STEP 5 / 5] End Applying CDI")

	return nil
}

// Delete executes `kubectl delete -f *.yaml`
func Delete(inventoryPath string) error {
	glog.Info("[STEP 0 / 4] Start Deleting CDI")

	glog.Info("[STEP 1 / 4] Fetch CDI variables from inventory")

	err := setCdiValuesFrom(inventoryPath)
	if err != nil {
		return err
	}

	var stderr bytes.Buffer

	err = kubectl.Run(os.Stdout, &stderr, "delete", "-f", crPath,
		"--ignore-not-found=true")

	if err != nil && kubectl.CRDAlreadyExists(stderr.String()) {
		return err
	}

	glog.Infof("[STEP 2 / 4] Wait up to %s for CDI CR to be deleted...", deleteTimeout.String())

	err = wait.PollImmediate(time.Second, deleteTimeout,
		util.IsCrDeleted(cdiCrdName, cdiCrName))
	if err != nil {
		return err
	}

	glog.Infof("[STEP 3 / 4] Wait up to %s for CDI operator to be deleted...", deleteTimeout.String())

	err = kubectl.Run(os.Stdout, os.Stderr, "delete", "-f",
		operatorPath, "--ignore-not-found=true")
	if err != nil {
		return err
	}

	glog.Info("[STEP 4 / 4] End Deleting CDI")

	return nil
}

func setCdiValuesFrom(inventoryPath string) error {
	operatorPath = path.Join(inventoryPath, "cdi", OperatorYaml)
	crPath = path.Join(inventoryPath, "cdi", CrYaml)

	var err error

	cdiNamespaceName, err = util.GetUniqueStringValueFromYamlFile(operatorPath,
		util.Namespace, "metadata.name")
	if err != nil {
		return err
	}

	cdiOperatorDeploymentName, err = util.GetUniqueStringValueFromYamlFile(operatorPath,
		util.Deployment, "metadata.name")
	if err != nil {
		return err
	}

	cdiCrdName, err = util.GetUniqueStringValueFromYamlFile(operatorPath,
		util.CustomResourceDefinition, "metadata.name")
	if err != nil {
		return err
	}

	cdiKindName, err = util.GetUniqueStringValueFromYamlFile(operatorPath,
		util.CustomResourceDefinition, "spec.names.kind")
	if err != nil {
		return err
	}

	cdiCrName, err = util.GetUniqueStringValueFromYamlFile(crPath,
		cdiKindName, "metadata.name")
	if err != nil {
		return err
	}

	return nil
}
