package tests

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubevirt.io/containerized-data-importer/pkg/apis/core/v1alpha1"
)

func makeDataVolumeSpec(name string, size string) *v1alpha1.DataVolume {
	return &v1alpha1.DataVolume{
		TypeMeta: metav1.TypeMeta{
			Kind:       "cdi.kubevirt.io/v1alpha1",
			APIVersion: "DataVolume",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1alpha1.DataVolumeSpec{
			Source: v1alpha1.DataVolumeSource{
				Registry: &v1alpha1.DataVolumeSourceRegistry{ // TODO 변수로 받기
					URL: "docker://kubevirt/fedora-cloud-registry-disk-demo", // TODO 변수로 받기
				},
			},
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
