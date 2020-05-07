package util

import (
	"bytes"
	"hypercloud-storage/hcsctl/pkg/kubectl"
)

const (
	// Namespace in string
	Namespace = "Namespace"
	// Deployment in string
	Deployment = "Deployment"
	// CustomResourceDefinition in string
	CustomResourceDefinition = "CustomResourceDefinition"
)

// IsCrdAvailable check if given CRD name is available
func IsCrdAvailable(crdKind string) func() (bool, error) {
	return func() (bool, error) {
		var stdout, stderr bytes.Buffer
		err := kubectl.Run(&stdout, &stderr, "get", crdKind)

		return err == nil, nil
	}
}

// IsCrDeployed check if CRD is deployed
func IsCrDeployed(crdName, crName string) func() (bool, error) {
	return func() (bool, error) {
		var stdout, stderr bytes.Buffer

		err := kubectl.Run(&stdout, &stderr, "get", crdName, crName,
			"-o", "jsonpath={.status.phase}")
		if err != nil {
			return false, err
		}

		return stdout.String() == "Deployed", nil
	}
}

// IsCrDeleted check if CRD is deleted
func IsCrDeleted(crdName, crName string) func() (bool, error) {
	return func() (bool, error) {
		isCrdExist, _ := IsCrdAvailable(crdName)()

		if !isCrdExist {
			return true, nil
		}

		var stdout, stderr bytes.Buffer
		err := kubectl.Run(&stdout, &stderr, "get", crdName, crName,
			"-o", "json", "--ignore-not-found=true")

		if err != nil {
			return false, err
		}

		return stdout.String() == "", nil
	}
}

// IsDeploymentCreated check if deployment created
func IsDeploymentCreated(namespace, deployment string) func() (bool, error) {
	return func() (bool, error) {
		var stdout, stderr bytes.Buffer
		err := kubectl.Run(&stdout, &stderr, "get", "-n", namespace,
			"deployments.apps", deployment, "-ojsonpath={.status.readyReplicas}")

		if err != nil {
			return false, err
		}

		readyReplicas := stdout.String()

		stdout.Reset()
		stderr.Reset()
		err = kubectl.Run(&stdout, &stderr, "get", "-n", namespace,
			"deployments.apps", deployment, "-ojsonpath={.status.replicas}")

		if err != nil {
			return false, err
		}

		replicas := stdout.String()

		return readyReplicas == replicas, nil
	}
}
