package tests

import (
	corev1 "k8s.io/api/core/v1"
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
