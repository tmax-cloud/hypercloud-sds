package testCeph

import (
    "k8s.io/client-go/dynamic"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
    "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
    "k8s.io/apimachinery/pkg/runtime/schema"

    "io/ioutil"
    "context"
)

/*func makeNamespaceSpec(namespacePrefix string) *corev1.Namespace {
	namespaceSpec := &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: namespacePrefix,
		},
	}

	return namespaceSpec
}*/

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
	if (resourceType == "namespaces" || resourceType == "persistentvolumeclaims") {
		apiGroup = ""
	} else if (resourceType == "deployments"  || resourceType == "daemonsets") {
		apiGroup = "apps"
	}

	gvs := schema.GroupVersionResource{Group: apiGroup, Version: "v1", Resource: resourceType}

	if resourceType == "namespaces" {
		_, err := client.Resource(gvs).Create(context.TODO(), &resourceSpec, metav1.CreateOptions{})
        if err != nil {
            return err
        }
	} else {
		_, err := client.Resource(gvs).Namespace("test-namespace").Create(context.TODO(), &resourceSpec, metav1.CreateOptions{})
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
	if (resourceType == "namespaces" || resourceType == "persistentvolumeclaims") {
		apiGroup = ""
	} else if (resourceType == "deployments"  || resourceType == "daemonsets") {
		apiGroup = "apps"
	}

	gvs := schema.GroupVersionResource{Group: apiGroup, Version: "v1", Resource: resourceType}

    deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions {
		PropagationPolicy: &deletePolicy,
	}

    if resourceType == "namespaces" {
        err := client.Resource(gvs).Delete(context.TODO(), resourceSpec.GetName(), deleteOptions)
        if err != nil {
            return err
        }
    } else {
        err := client.Resource(gvs).Namespace("test-namespace").Delete(context.TODO(), resourceSpec.GetName(), deleteOptions)
        if err != nil {
            return err
        }
    }

    return nil
}

func getK8sObjectStatusFromYaml(resourceType, yamlPath string, client dynamic.Interface) (string, error) {
	resourceSpec, err := createSpecFromYaml(yamlPath)
	if err != nil {
		return "",err
	}

	var apiGroup string
	if (resourceType == "namespaces" || resourceType == "persistentvolumeclaims") {
		apiGroup = ""
	} else if (resourceType == "deployments"  || resourceType == "daemonsets") {
		apiGroup = "apps"
	}

	gvs := schema.GroupVersionResource{Group: apiGroup, Version: "v1", Resource: resourceType}

    var status string
    if resourceType == "namespaces" {
        result, err := client.Resource(gvs).Get(context.TODO(), resourceSpec.GetName(), metav1.GetOptions{})
        if err != nil {
            return "",err
        }
        status, _, err = unstructured.NestedString(result.Object, "status", "phase")
        return status,err

    } else if resourceType == "persistentvolumeclaims" {
        result, err := client.Resource(gvs).Namespace("test-namespace").Get(context.TODO(), resourceSpec.GetName(), metav1.GetOptions{})
        if err != nil {
            return "",err
        }
        status, _, err = unstructured.NestedString(result.Object, "status", "phase")
        return status,err
    }

    // TODO: add other resources status routine
    return "",nil
}

/*func getNamespaceStatus(client *dynamic.Interface) (string, error) {
	gvs := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "namespace"}

    result, err := *client.Resource(gvs).Get(context.TODO(), "test-namespace", metav1.GetOptions{})

    if err != nil {
        return "",err
    }

    status, _, err := unstructured.NestedString(result.Object, "status", "phase")

    return status,err

}

func getPvcStatus(client *dynamic.Interface) (string, error) {
	gvs := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "namespace"}

    result, err := *client.Resource(gvs).Get(context.TODO(), "test-namespace", metav1.GetOptions{})

    if err != nil {
        return "",err
    }

    status, _, err := unstructured.NestedString(result.Object, "status", "phase")

    return status,err

}


func makePvcInStorageClassSpec(name string, namespace string, size string,
	storageClassName string, accessMode corev1.PersistentVolumeAccessMode) *corev1.PersistentVolumeClaim {
	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{accessMode},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(size),
				},
			},
			StorageClassName: &storageClassName,
		},
	}
}

func createPvcInStorageClass(clientset *kubernetes.Clientset, pvcSpec *corev1.PersistentVolumeClaim) (
	*corev1.PersistentVolumeClaim, error) {
	pvc, err := clientset.CoreV1().PersistentVolumeClaims(pvcSpec.Namespace).Create(pvcSpec)
	return pvc, err
}
*/
