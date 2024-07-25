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
	"sync"
	"time"

	virtoolv1alpha1 "github.com/bryce-davidson/virtool-operator/api/v1alpha1"
	"github.com/bryce-davidson/virtool-operator/factory"
	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ------------------------------------------------------------------------------

type TestLogSink struct {
	logChan chan string
	buffer  []string
	mu      sync.Mutex
}

func NewTestLogSink() *TestLogSink {
	return &TestLogSink{
		logChan: make(chan string, 100),
	}
}

func (t *TestLogSink) Init(info logr.RuntimeInfo) {}
func (t *TestLogSink) Enabled(level int) bool     { return true }
func (t *TestLogSink) Info(level int, msg string, keysAndValues ...interface{}) {
	logMsg := fmt.Sprintf("%s: %v", msg, keysAndValues)
	t.mu.Lock()
	t.buffer = append(t.buffer, logMsg)
	t.mu.Unlock()
	t.logChan <- logMsg
}

func (t *TestLogSink) Error(err error, msg string, keysAndValues ...interface{}) {
	logMsg := fmt.Sprintf("ERROR: %s: %v", msg, err)
	t.mu.Lock()
	t.buffer = append(t.buffer, logMsg)
	t.mu.Unlock()
	t.logChan <- logMsg
}

func (t *TestLogSink) WithValues(keysAndValues ...interface{}) logr.LogSink { return t }
func (t *TestLogSink) WithName(name string) logr.LogSink                    { return t }
func (t *TestLogSink) String() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return strings.Join(t.buffer, "\n")
}

// ------------------------------------------------------------------------------

var _ = Describe("VirtoolApp Controller", func() {
	const resourceName = "test-resource"
	const namespace = "default"

	ctx := context.Background()
	typeNamespacedName := types.NamespacedName{Name: resourceName, Namespace: namespace}

	BeforeEach(func() {
		// Use the new factory function to create the VirtoolApp
		virtoolApp := factory.NewVirtoolApp(resourceName, namespace)
		Expect(k8sClient.Create(ctx, virtoolApp)).To(Succeed())
	})

	AfterEach(func() {
		cleanupResource(ctx, typeNamespacedName)
	})

	Describe("Basic Reconciliation", func() {
		It("should successfully reconcile the resource", func() {
			controllerReconciler := &VirtoolAppReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			result, err := controllerReconciler.Reconcile(ctx, reconcile.Request{NamespacedName: typeNamespacedName})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(ctrl.Result{}))

			// Verify the reconciled state
			var updatedApp virtoolv1alpha1.VirtoolApp
			Expect(k8sClient.Get(ctx, typeNamespacedName, &updatedApp)).To(Succeed())
		})
	})

	Describe("Container Image Logging", func() {
		BeforeEach(func() {
			createTestPod(ctx, namespace)
		})

		AfterEach(func() {
			cleanupTestPod(ctx, namespace)
		})

		It("should locate and print container images", func() {
			logSink := NewTestLogSink()
			logger := logr.New(logSink)

			reconciler := &VirtoolAppReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
				Log:    logger,
			}

			go func() {
				_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: typeNamespacedName})
				Expect(err).NotTo(HaveOccurred())
			}()

			expectedLogs := []string{
				"Starting reconciliation",
				"Listed pods",
				"Processing pod",
				"Container image tag",
				"test-pod",
				"container1",
				"nginx:1.14.2",
				"container2",
				"busybox:1.28",
				"Reconciliation completed",
			}

			Eventually(func() bool {
				return containsAllLogs(logSink.buffer, expectedLogs)
			}, 10*time.Second, 250*time.Millisecond).Should(BeTrue())
		})
	})
})

func createTestPod(ctx context.Context, namespace string) {
	testPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: namespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{Name: "container1", Image: "nginx:1.14.2"},
				{Name: "container2", Image: "busybox:1.28"},
			},
		},
	}
	Expect(k8sClient.Create(ctx, testPod)).To(Succeed())

	Eventually(func() error {
		return k8sClient.Get(ctx, types.NamespacedName{Name: "test-pod", Namespace: namespace}, &corev1.Pod{})
	}, 10*time.Second, 250*time.Millisecond).Should(Succeed())
}

func cleanupResource(ctx context.Context, namespacedName types.NamespacedName) {
	resource := &virtoolv1alpha1.VirtoolApp{}
	err := k8sClient.Get(ctx, namespacedName, resource)
	if err == nil {
		Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
	}
}

func cleanupTestPod(ctx context.Context, namespace string) {
	testPod := &corev1.Pod{}
	err := k8sClient.Get(ctx, types.NamespacedName{Name: "test-pod", Namespace: namespace}, testPod)
	if err == nil {
		Expect(k8sClient.Delete(ctx, testPod)).To(Succeed())
	}
}

func containsAllLogs(receivedLogs []string, expectedLogs []string) bool {
	if len(expectedLogs) == 0 {
		return true
	}

	receivedMap := make(map[string]struct{}, len(receivedLogs))
	for _, log := range receivedLogs {
		receivedMap[log] = struct{}{}
	}

	for _, expected := range expectedLogs {
		found := false
		for received := range receivedMap {
			if strings.Contains(received, expected) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
