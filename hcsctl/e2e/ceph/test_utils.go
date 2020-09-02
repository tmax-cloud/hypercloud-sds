package testCeph

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func makeNamespaceSpec(namespacePrefix string) *corev1.Namespace {
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
}

func createNamespace(clientset *kubernetes.Clientset, nsSpec *corev1.Namespace) (*corev1.Namespace, error) {
	ns, err := clientset.CoreV1().Namespaces().Create(nsSpec)

	return ns, err
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
