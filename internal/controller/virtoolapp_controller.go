/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"strings"

	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	virtoolv1alpha1 "github.com/bryce-davidson/virtool-operator/api/v1alpha1"
)

// VirtoolAppReconciler reconciles a VirtoolApp object
type VirtoolAppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=virtool.virtool.ca,resources=virtoolapps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=virtool.virtool.ca,resources=virtoolapps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=virtool.virtool.ca,resources=virtoolapps/finalizers,verbs=update

//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the VirtoolApp object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *VirtoolAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// List all pods in the cluster
	var podList corev1.PodList
	if err := r.List(ctx, &podList, &client.ListOptions{}); err != nil {
		logger.Error(err, "Unable to list pods")
		return ctrl.Result{}, err
	}

	// Iterate through all pods
	for _, pod := range podList.Items {
		// Log the pod name and namespace
		logger.Info("Processing pod", "name", pod.Name, "namespace", pod.Namespace)

		// Iterate through all containers in the pod
		for _, container := range pod.Spec.Containers {
			// Extract the image tag
			imageParts := strings.Split(container.Image, ":")
			imageTag := "latest" // Default tag if not specified
			if len(imageParts) > 1 {
				imageTag = imageParts[len(imageParts)-1]
			}

			// Log the container name and image tag
			logger.Info("Container image tag", "container", container.Name, "image", container.Image, "tag", imageTag)
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *VirtoolAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&virtoolv1alpha1.VirtoolApp{}).
		Complete(r)
}
