package factory

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	virtoolv1alpha1 "github.com/bryce-davidson/virtool-operator/api/v1alpha1"
)

type VirtoolAppOption func(*virtoolv1alpha1.VirtoolApp)

const (
	defaultVersion       = "1.0.0"
	defaultImage         = "default-image:latest"
	defaultReplicas      = 1
	defaultCPULimit      = "100m"
	defaultMemoryLimit   = "128Mi"
	defaultCPURequest    = "50m"
	defaultMemoryRequest = "64Mi"
)

func defaultResourceRequirements() corev1.ResourceRequirements {
	return corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(defaultCPULimit),
			corev1.ResourceMemory: resource.MustParse(defaultMemoryLimit),
		},
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(defaultCPURequest),
			corev1.ResourceMemory: resource.MustParse(defaultMemoryRequest),
		},
	}
}

func NewVirtoolApp(name, namespace string, opts ...VirtoolAppOption) *virtoolv1alpha1.VirtoolApp {
	v := &virtoolv1alpha1.VirtoolApp{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "cache.example.com/v1alpha1",
			Kind:       "VirtoolApp",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: virtoolv1alpha1.VirtoolAppSpec{
			Version: defaultVersion,
			Components: []virtoolv1alpha1.ComponentSpec{
				{
					Name:      "default",
					Image:     defaultImage,
					Replicas:  defaultReplicas,
					Resources: defaultResourceRequirements(),
				},
			},
		},
		Status: virtoolv1alpha1.VirtoolAppStatus{
			ComponentsStatus: []virtoolv1alpha1.ComponentStatus{},
			Conditions:       []metav1.Condition{},
		},
	}

	for _, opt := range opts {
		opt(v)
	}

	return v
}
