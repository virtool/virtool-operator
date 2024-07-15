package factory

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	virtoolv1alpha1 "github.com/bryce-davidson/virtool-operator/api/v1alpha1"
)

func pointer(i int32) *int32 {
	copy := i
	return &copy
}

type VirtoolAppOption func(*virtoolv1alpha1.VirtoolApp)

const (
	defaultVersion         = "1.0.0"
	defaultImage           = "default-image:latest"
	defaultReplicas        = 1
	defaultCPULimit        = "100m"
	defaultMemoryLimit     = "128Mi"
	defaultCPURequest      = "50m"
	defaultMemoryRequest   = "64Mi"
	defaultUpdateType      = "RollingUpdate"
	defaultMaxUnavailable  = "25%"
	defaultMaxSurge        = "25%"
	defaultRegistry        = "default-registry.example.com"
	defaultImagePullSecret = "default-pull-secret"
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
			Kind:       "virtoolv1alpha1.VirtoolApp",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: virtoolv1alpha1.VirtoolAppSpec{
			Components: map[string]virtoolv1alpha1.ComponentSpec{
				"default": {
					Version:     defaultVersion,
					Image:       defaultImage,
					Replicas:    pointer(int32(defaultReplicas)),
					UpdateOrder: 0,
					Resources:   defaultResourceRequirements(),
				},
			},
			UpdateStrategy: virtoolv1alpha1.UpdateStrategy{
				Type:           defaultUpdateType,
				MaxUnavailable: &intstr.IntOrString{Type: intstr.String, StrVal: defaultMaxUnavailable},
				MaxSurge:       &intstr.IntOrString{Type: intstr.String, StrVal: defaultMaxSurge},
			},
			GlobalConfig: virtoolv1alpha1.GlobalConfig{
				Registry:        defaultRegistry,
				ImagePullSecret: defaultImagePullSecret,
			},
		},
		Status: virtoolv1alpha1.VirtoolAppStatus{
			ComponentStatus: map[string]virtoolv1alpha1.ComponentStatus{},
			Conditions:      []metav1.Condition{},
		},
	}

	for _, opt := range opts {
		opt(v)
	}

	return v
}
