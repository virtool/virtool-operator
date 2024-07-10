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
	"fmt"
	"strings"
	"time"

	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	virtoolv1alpha1 "github.com/bryce-davidson/virtool-operator/api/v1alpha1"
)

type TestLogSink struct {
	buffer []string
}

func (t *TestLogSink) Init(info logr.RuntimeInfo) {}
func (t *TestLogSink) Enabled(level int) bool     { return true }
func (t *TestLogSink) Info(level int, msg string, keysAndValues ...interface{}) {
	t.buffer = append(t.buffer, fmt.Sprintf("%s: %v", msg, keysAndValues))
}

func (t *TestLogSink) Error(err error, msg string, keysAndValues ...interface{}) {
	t.buffer = append(t.buffer, msg)
}
func (t *TestLogSink) WithValues(keysAndValues ...interface{}) logr.LogSink { return t }
func (t *TestLogSink) WithName(name string) logr.LogSink                    { return t }
func (t *TestLogSink) String() string {
	return strings.Join(t.buffer, "\n")
}

// ------------------------------------------------------------------------------

var _ = Describe("VirtoolApp Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-resource"
		const namespace = "default"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: namespace,
		}
		virtoolapp := &virtoolv1alpha1.VirtoolApp{}

		BeforeEach(func() {
			By("creating the custom resource for the Kind VirtoolApp")
			err := k8sClient.Get(ctx, typeNamespacedName, virtoolapp)
			if err != nil && errors.IsNotFound(err) {
				resource := &virtoolv1alpha1.VirtoolApp{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: namespace,
					},
					// TODO(user): Specify other spec details if needed.
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}

			By("creating a test pod")
			testPod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: namespace,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "container1",
							Image: "nginx:1.14.2",
						},
						{
							Name:  "container2",
							Image: "busybox:1.28",
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, testPod)).To(Succeed())

			// Verify pod creation
			createdPod := &corev1.Pod{}
			Eventually(func() error {
				return k8sClient.Get(ctx, types.NamespacedName{Name: "test-pod", Namespace: namespace}, createdPod)
			}, time.Second*10, time.Millisecond*250).Should(Succeed())
		})

		AfterEach(func() {
			By("Cleanup the specific resource instance VirtoolApp")
			resource := &virtoolv1alpha1.VirtoolApp{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())

			By("Cleanup the test pod")
			testPod := &corev1.Pod{}
			err = k8sClient.Get(ctx, types.NamespacedName{Name: "test-pod", Namespace: namespace}, testPod)
			if err == nil {
				Expect(k8sClient.Delete(ctx, testPod)).To(Succeed())
			}
		})

		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			controllerReconciler := &VirtoolAppReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should locate and print container images", func() {
			logOutput := &TestLogSink{}
			logger := logr.New(logOutput)

			reconciler := &VirtoolAppReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
				Log:    logger,
			}

			_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: typeNamespacedName})
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				fmt.Println("Full log output:", logOutput.String())
				return logOutput.String()
			}, 60*time.Second, time.Second).Should(And(
				ContainSubstring("Starting reconciliation"),
				ContainSubstring("Listed pods"),
				ContainSubstring("Container image tag"),
				ContainSubstring("test-pod"),
				ContainSubstring("container1"),
				ContainSubstring("nginx:1.14.2"),
				ContainSubstring("container2"),
				ContainSubstring("busybox:1.28"),
				ContainSubstring("Reconciliation completed"),
			))

			fmt.Println("Log output:", logOutput.String())
		})
	})
})
