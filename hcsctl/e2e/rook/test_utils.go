package rook

import (
    "k8s.io/client-go/dynamic"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
    "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
    "k8s.io/apimachinery/pkg/runtime/schema"

    "io/ioutil"
    "context"
)

func createSpecFromYaml(yamlPath string) (unstructured.Unstructured, error) {
    yamlByteArray, err := ioutil.ReadFile(yamlPath)
    if err != nil {
        return unstructured.Unstructured{}, err
    }

    spec := &unstructured.Unstructured{}
    decode := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode

    _, _, err = decode(yamlByteArray, nil, spec)
    if err != nil {
        return unstructured.Unstructured{}, err
    }

    return *spec, nil
}

/* Create k8s Object with its resource name and yaml path
  yamlPath : Path of yaml manifest file
  resourceType: Resource Type of k8s object 
*/
func createK8sObjectFromYaml(resourceType, yamlPath string, client dynamic.Interface) error {
	resourceSpec, err := createSpecFromYaml(yamlPath)
	if err != nil {
		return err
	}

	var apiGroup string
	if (resourceType == ResourceNamespace || resourceType == ResourcePVC) {
		apiGroup = APIGroupCore
	} else if (resourceType == ResourceDeployment  || resourceType == ResourceDaemonSet) {
		apiGroup = APIGroupApps
	}

	gvs := schema.GroupVersionResource{Group: apiGroup, Version: "v1", Resource: resourceType}

	if resourceType == ResourceNamespace {
		_, err := client.Resource(gvs).Create(context.TODO(), &resourceSpec, metav1.CreateOptions{})
        if err != nil {
            return err
        }
	} else {
		_, err := client.Resource(gvs).Namespace(TestNamespaceName).Create(context.TODO(), &resourceSpec, metav1.CreateOptions{})
        if err != nil {
            return err
        }
	}

	return nil
}

func deleteK8sObjectFromYaml(resourceType, yamlPath string, client dynamic.Interface) error {
	resourceSpec, err := createSpecFromYaml(yamlPath)
	if err != nil {
		return err
	}

	var apiGroup string
	if (resourceType == ResourceNamespace || resourceType == ResourcePVC) {
		apiGroup = APIGroupCore
	} else if (resourceType == ResourceDeployment  || resourceType == ResourceDaemonSet) {
		apiGroup = APIGroupApps
	}

	gvs := schema.GroupVersionResource{Group: apiGroup, Version: "v1", Resource: resourceType}

    deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions {
		PropagationPolicy: &deletePolicy,
	}

    if resourceType == ResourceNamespace {
        err := client.Resource(gvs).Delete(context.TODO(), resourceSpec.GetName(), deleteOptions)
        if err != nil {
            return err
        }
    } else {
        err := client.Resource(gvs).Namespace(TestNamespaceName).Delete(context.TODO(), resourceSpec.GetName(), deleteOptions)
        if err != nil {
            return err
        }
    }

    return nil
}

func getK8sObjectStatusFromYaml(resourceType, yamlPath string, client dynamic.Interface) (string, error) {
	resourceSpec, err := createSpecFromYaml(yamlPath)
	if err != nil {
		return "", err
	}

	var apiGroup string
	if (resourceType == ResourceNamespace || resourceType == ResourcePVC) {
		apiGroup = APIGroupCore
	} else if (resourceType == ResourceDeployment || resourceType == ResourceDaemonSet) {
		apiGroup = APIGroupApps
	}

	gvs := schema.GroupVersionResource{Group: apiGroup, Version: "v1", Resource: resourceType}

    var status string
    if resourceType == ResourceNamespace {
        result, err := client.Resource(gvs).Get(context.TODO(), resourceSpec.GetName(), metav1.GetOptions{})

        if err != nil {
            return "", err
        }
        status, _, err = unstructured.NestedString(result.Object, "status", "phase")
        return status, err

    } else if resourceType == ResourcePVC {
        result, err := client.Resource(gvs).Namespace(TestNamespaceName).Get(context.TODO(), resourceSpec.GetName(), metav1.GetOptions{})
        if err != nil {
            return "", err
        }
        status, _, err = unstructured.NestedString(result.Object, "status", "phase")
        return status,err
    }

    // TODO: add other resources status routine
    return "",nil
}

