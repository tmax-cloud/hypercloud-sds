package main

import (
	"hypercloud-storage/hcsctl/pkg/util"
	"testing"
)

var (
	clusterYAML = "../../../../hack/inventory/minikube/rook/cluster.yaml"
	commonYAML  = "../../../../hack/inventory/minikube/rook/common.yaml"
)

func TestGetCephClusterName(t *testing.T) {
	value, err := util.GetValueFromYamlFile(clusterYAML, "CephCluster", "metadata.name")
	if err != nil {
		t.Fatal(err)
	}

	for _, val := range value {
		t.Log("CephCluster Name: ", val)
	}
}
func TestGetCephClusterMgr(t *testing.T) {
	value, err := util.GetValueFromYamlFile(clusterYAML, "CephCluster", "spec.mgr.modules")
	if err != nil {
		t.Fatal(err)
	}

	for _, val := range value {
		t.Log("CephCluster spec.mgr.modules: ", val)
	}
}

func TestGetCephClusterStorageConfig(t *testing.T) {
	value, err := util.GetValueFromYamlFile(clusterYAML, "CephCluster", "spec.storage.config")
	if err != nil {
		t.Fatal(err)
	}

	for _, val := range value {
		t.Log("CephCluster spec.storage.config: ", val)
	}
}

func TestGetNotFound(t *testing.T) {
	value, err := util.GetValueFromYamlFile(clusterYAML, "CephCluster", "metadata.notFound")
	if err != nil {
		t.Fatal(err)
	}

	for _, val := range value {
		t.Log("CephCluster NotFound: ", val)
	}
}

func TestGetCephCRDName(t *testing.T) {
	value, err := util.GetValueFromYamlFile(commonYAML, "CustomResourceDefinition", "metadata.name")
	if err != nil {
		t.Fatal(err)
	}

	for _, val := range value {
		t.Log("CRD Name: ", val)
	}
}

func TestGetCephCRDShort(t *testing.T) {
	value, err := util.GetValueFromYamlFile(commonYAML, "CustomResourceDefinition", "spec.names.shortNames")
	if err != nil {
		t.Fatal(err)
	}

	for _, val := range value {
		t.Log("CRD shortNames: ", val)
	}
}

func TestGetCephRoleBinding(t *testing.T) {
	value, err := util.GetValueFromYamlFile(commonYAML, "RoleBinding", "metadata.name")
	if err != nil {
		t.Fatal(err)
	}

	for _, val := range value {
		t.Log("RoleBinding Name: ", val)
	}
}

func TestGetCephClusterRole(t *testing.T) {
	value, err := util.GetValueFromYamlFile(commonYAML, "ClusterRole", "rules[1]")
	if err != nil {
		t.Fatal(err)
	}

	for _, val := range value {
		t.Log("ClusterRole rules[1]: ", val)
	}
}

func TestGetCephClusterRoleIndex(t *testing.T) {
	value, err := util.GetValueFromYamlFile(commonYAML, "ClusterRole", "rules[1].verbs")
	if err != nil {
		t.Fatal(err)
	}

	for _, val := range value {
		t.Log("ClusterRole rules[1].verbs: ", val)
	}
}
