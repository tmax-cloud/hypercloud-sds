package rook

import "time"

const (
	PollingIntervalDefault              = 3 * time.Second
	PollingIntervalForDeletingNamespace = 10 * time.Second

	TimeoutForDeletingNamespace         = 500 * time.Second *2
	TimeOutForCreatingPvc               = 60 * time.Second *2

    TestNamespaceYamlPath   = "test-manifests/namespace.yaml"
    TestCephFsPvcYamlPath   = "test-manifests/cephfs-pvc.yaml"
    TestRbdPvcYamlPath      = "test-manifests/rbd-pvc.yaml"

    TestNamespaceName   = "test-namespace"

    ResourceNamespace   = "namespaces"
    ResourcePVC         = "persistentvolumeclaims"
    ResourceDaemonSet   = "daemonsets"
    ResourceDeployment  = "deployments"

    APIGroupCore   = ""
    APIGroupApps   = "apps"
)
