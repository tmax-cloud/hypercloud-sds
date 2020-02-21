package tests

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cdiv1alpha1 "kubevirt.io/containerized-data-importer/pkg/apis/core/v1alpha1"
)

func makeDataVolumeSpec(name string, size string, source *cdiv1alpha1.DataVolumeSource) *cdiv1alpha1.DataVolume {
	return &cdiv1alpha1.DataVolume{
		TypeMeta: metav1.TypeMeta{
			Kind:       "cdi.kubevirt.io/v1alpha1",
			APIVersion: "DataVolume",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: cdiv1alpha1.DataVolumeSpec{
			Source: *source,
			PVC: &v1.PersistentVolumeClaimSpec{
				AccessModes: []v1.PersistentVolumeAccessMode{v1.ReadWriteMany}, // TODO 변수로 받기
				Resources: v1.ResourceRequirements{
					Requests: v1.ResourceList{
						v1.ResourceStorage: resource.MustParse(size),
					},
				},
			},
		},
	}
}

func makeDataVolumeSourceHTTP(url string) *cdiv1alpha1.DataVolumeSource {
	return &cdiv1alpha1.DataVolumeSource{
		HTTP: &cdiv1alpha1.DataVolumeSourceHTTP{URL: url},
	}
}

func makeDataVolumeSourceRegistry(url string) *cdiv1alpha1.DataVolumeSource {
	return &cdiv1alpha1.DataVolumeSource{
		Registry: &cdiv1alpha1.DataVolumeSourceRegistry{URL: url},
	}
}

func makeDataVolumeSourcePVC(namespace string, name string) *cdiv1alpha1.DataVolumeSource {
	return &cdiv1alpha1.DataVolumeSource{
		PVC: &cdiv1alpha1.DataVolumeSourcePVC{
			Namespace: namespace,
			Name:      name,
		},
	}
}
