/*
Copyright 2026.

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

package reconcile_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"cuelang.org/go/cue/cuecontext"
	corev1 "k8s.io/api/core/v1"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/open-platform-model/library/opm/k8s/labels"
	"github.com/open-platform-model/library/opm/k8s/object"

	releasesv1alpha1 "github.com/open-platform-model/opm-controller/api/v1alpha1"
	opmreconcile "github.com/open-platform-model/opm-controller/internal/reconcile"
	"github.com/open-platform-model/opm-controller/internal/render"
	"github.com/open-platform-model/opm-controller/internal/status"
)

// unknownKindRenderResult builds a render result whose single resource has a
// GVK the cluster does not serve, so the apply phase fails with a transient
// (non-stalled) error — the staging hook for blocked→Ready recovery tests.
func unknownKindRenderResult(namespace string) *render.RenderResult {
	cueCtx := cuecontext.New()
	obj := cueCtx.CompileString(fmt.Sprintf(`{
	apiVersion: "test.example.com/v1"
	kind:       "Phantom"
	metadata: {
		name:      "phantom-1"
		namespace: %q
		labels: {
			%q: %q
			%q: %q
		}
	}
}`, namespace,
		labels.ManagedBy, labels.ManagedByController,
		labels.ModuleInstanceUUID, stubInstanceUUID))
	if obj.Err() != nil {
		panic(fmt.Sprintf("compiling phantom stub resource: %v", obj.Err()))
	}

	resource := &object.Resource{
		Value:       obj,
		Instance:    "phantom-1",
		Component:   "phantom",
		Transformer: "kubernetes#simple",
	}

	return &render.RenderResult{
		Resources: []*object.Resource{resource},
	}
}

// stageCrashStatus writes, through the status subresource, the status that a
// controller leaves when it dies after the apply and before its deferred status
// commit: the given inventory (nil for a first apply), and no lastApplied*
// digests and no lastAppliedInputs. instanceUUID, the conditions and
// observedGeneration stay as they are, as a crash keeps them.
func stageCrashStatus(nn types.NamespacedName, inventory *releasesv1alpha1.Inventory) {
	GinkgoHelper()
	var mi releasesv1alpha1.ModuleInstance
	Expect(k8sClient.Get(ctx, nn, &mi)).To(Succeed())
	base := mi.DeepCopy()
	mi.Status.Inventory = inventory
	mi.Status.LastAppliedSourceDigest = ""
	mi.Status.LastAppliedConfigDigest = ""
	mi.Status.LastAppliedRenderDigest = ""
	mi.Status.LastAppliedInputs = nil
	Expect(k8sClient.Status().Patch(ctx, &mi, client.MergeFrom(base))).To(Succeed())

	var staged releasesv1alpha1.ModuleInstance
	Expect(k8sClient.Get(ctx, nn, &staged)).To(Succeed())
	Expect(staged.Status.LastAppliedSourceDigest).To(BeEmpty())
	Expect(staged.Status.LastAppliedConfigDigest).To(BeEmpty())
	Expect(staged.Status.LastAppliedRenderDigest).To(BeEmpty())
	Expect(staged.Status.LastAppliedInputs).To(BeNil())
	if inventory == nil {
		Expect(staged.Status.Inventory).To(BeNil())
	} else {
		Expect(staged.Status.Inventory).NotTo(BeNil())
		Expect(staged.Status.Inventory.Entries).To(HaveLen(len(inventory.Entries)))
	}
}

// expectConverged checks the converged outcome of a reconcile after a crash:
// Ready=True on the current generation, an inventory of exactly the rendered
// ConfigMaps (named as Kind/name), and no ConfigMap with the instance's UUID
// label outside it.
func expectConverged(nn types.NamespacedName, rendered ...string) {
	GinkgoHelper()
	var mi releasesv1alpha1.ModuleInstance
	Expect(k8sClient.Get(ctx, nn, &mi)).To(Succeed())

	ready := apimeta.FindStatusCondition(mi.Status.Conditions, status.ReadyCondition)
	Expect(ready).NotTo(BeNil())
	Expect(ready.Status).To(Equal(metav1.ConditionTrue))
	Expect(mi.Status.ObservedGeneration).To(Equal(mi.Generation), "status.observedGeneration")

	Expect(mi.Status.Inventory).NotTo(BeNil(), "status.inventory is written again")
	Expect(inventoryNames(mi.Status.Inventory)).To(ConsistOf(rendered), "inventory lists exactly the rendered objects")

	Expect(mi.Status.InstanceUUID).To(Equal(stubInstanceUUID))
	var owned corev1.ConfigMapList
	Expect(k8sClient.List(ctx, &owned, client.InNamespace(namespace),
		client.MatchingLabels{labels.ModuleInstanceUUID: mi.Status.InstanceUUID})).To(Succeed())
	var ownedNames []string
	for _, cm := range owned.Items {
		if cm.DeletionTimestamp == nil {
			ownedNames = append(ownedNames, "ConfigMap/"+cm.Name)
		}
	}
	Expect(inventoryNames(mi.Status.Inventory)).To(ContainElements(ownedNames),
		"no object with the instance's UUID label is outside status.inventory")
}

var _ = Describe("Reconcile State Recovery", func() {
	// Validates Stalled → Ready recovery once the identity mismatch clears
	// (design 3.1; resolution failure staged via the suite's error renderers).
	It("should recover from Stalled once the identity matches", func() {
		mrName := "stalled-recover-mr"
		createModuleInstance(mrName)
		nn := types.NamespacedName{Name: mrName, Namespace: namespace}

		params := reconcileParams()
		params.Renderer = identityMismatchRenderer()
		ensureFinalizer(params, nn)

		By("first reconcile stalls on the identity mismatch")
		result, err := opmreconcile.ReconcileModuleInstance(ctx, params, ctrl.Request{NamespacedName: nn})
		Expect(err).NotTo(HaveOccurred(), "stalled failures return nil error")
		Expect(result.RequeueAfter).To(Equal(opmreconcile.StalledRecheckInterval))

		var stalledMI releasesv1alpha1.ModuleInstance
		Expect(k8sClient.Get(ctx, nn, &stalledMI)).To(Succeed())
		stalled := apimeta.FindStatusCondition(stalledMI.Status.Conditions, status.StalledCondition)
		Expect(stalled).NotTo(BeNil())
		Expect(stalled.Status).To(Equal(metav1.ConditionTrue))
		Expect(stalled.Reason).To(Equal(status.ResolutionFailedReason))
		ready := apimeta.FindStatusCondition(stalledMI.Status.Conditions, status.ReadyCondition)
		Expect(ready).NotTo(BeNil())
		Expect(ready.Status).To(Equal(metav1.ConditionFalse))

		By("second reconcile recovers once the identity matches")
		params.Renderer = &stubRenderer{}
		result, err = opmreconcile.ReconcileModuleInstance(ctx, params, ctrl.Request{NamespacedName: nn})
		Expect(err).NotTo(HaveOccurred())
		Expect(result.RequeueAfter).To(BeZero())

		var recovered releasesv1alpha1.ModuleInstance
		Expect(k8sClient.Get(ctx, nn, &recovered)).To(Succeed())
		ready = apimeta.FindStatusCondition(recovered.Status.Conditions, status.ReadyCondition)
		Expect(ready).NotTo(BeNil())
		Expect(ready.Status).To(Equal(metav1.ConditionTrue))
		Expect(apimeta.FindStatusCondition(recovered.Status.Conditions, status.StalledCondition)).To(BeNil(),
			"Stalled condition is cleared on recovery")

		var cm corev1.ConfigMap
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: "test-module", Namespace: namespace}, &cm)).To(Succeed())

		By("history reflects the failure → success transition")
		Expect(len(recovered.Status.History)).To(BeNumerically(">=", 2))
		Expect(recovered.Status.History[0].Phase).To(Equal("complete"))
		Expect(recovered.Status.History[1].Message).To(ContainSubstring("identity mismatch"))

		// Cleanup
		Expect(k8sClient.Delete(ctx, &cm)).To(Succeed())
		cleanupInstance(nn)
	})

	// Validates transient acquisition failure → Ready recovery: a module the
	// registry did not serve retries on the bounded backoff as a non-stalled
	// ResolutionFailed, and recovers once it is served.
	It("should recover from a transient acquisition failure", func() {
		mrName := "acquire-recover-mr"
		createModuleInstance(mrName)
		nn := types.NamespacedName{Name: mrName, Namespace: namespace}

		params := reconcileParams()
		params.Renderer = acquireFailureRenderer()
		ensureFinalizer(params, nn)

		By("first reconcile retries the acquisition on the backoff")
		result, err := opmreconcile.ReconcileModuleInstance(ctx, params, ctrl.Request{NamespacedName: nn})
		Expect(err).NotTo(HaveOccurred(), "transient failures return nil error with backoff")
		Expect(result.RequeueAfter).To(Equal(opmreconcile.ComputeBackoff(1)))

		var failing releasesv1alpha1.ModuleInstance
		Expect(k8sClient.Get(ctx, nn, &failing)).To(Succeed())
		ready := apimeta.FindStatusCondition(failing.Status.Conditions, status.ReadyCondition)
		Expect(ready).NotTo(BeNil())
		Expect(ready.Status).To(Equal(metav1.ConditionFalse))
		Expect(ready.Reason).To(Equal(status.ResolutionFailedReason))
		Expect(apimeta.FindStatusCondition(failing.Status.Conditions, status.StalledCondition)).To(BeNil(),
			"an acquisition failure without a typed cause does not stall")
		Expect(failing.Status.NextRetryAt).NotTo(BeNil())
		Expect(failing.Status.FailureCounters).NotTo(BeNil())
		Expect(failing.Status.FailureCounters.Reconcile).To(Equal(int64(1)))

		By("second reconcile recovers once the module is served")
		params.Renderer = &stubRenderer{}
		result, err = opmreconcile.ReconcileModuleInstance(ctx, params, ctrl.Request{NamespacedName: nn})
		Expect(err).NotTo(HaveOccurred())
		Expect(result.RequeueAfter).To(BeZero())

		var recovered releasesv1alpha1.ModuleInstance
		Expect(k8sClient.Get(ctx, nn, &recovered)).To(Succeed())
		ready = apimeta.FindStatusCondition(recovered.Status.Conditions, status.ReadyCondition)
		Expect(ready).NotTo(BeNil())
		Expect(ready.Status).To(Equal(metav1.ConditionTrue))
		Expect(recovered.Status.NextRetryAt).To(BeNil(), "nextRetryAt is cleared on success")
		if recovered.Status.FailureCounters != nil {
			Expect(recovered.Status.FailureCounters.Reconcile).To(BeZero(), "the reconcile counter resets on success")
		}

		var cm corev1.ConfigMap
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: "test-module", Namespace: namespace}, &cm)).To(Succeed())

		// Cleanup
		Expect(k8sClient.Delete(ctx, &cm)).To(Succeed())
		cleanupInstance(nn)
	})

	// Validates transient-blocked → Ready recovery: a failing apply keeps the
	// instance on the fast exponential backoff (not the 30-minute stalled
	// recheck) until the cause clears (design 3.2 — the OCIRepository
	// SourceReady staging is retired design; an unserved GVK stages the
	// equivalent transient block in the current pipeline).
	It("should recover from a transient blocked state when the cause clears", func() {
		mrName := "blocked-recover-mr"
		createModuleInstance(mrName)
		nn := types.NamespacedName{Name: mrName, Namespace: namespace}

		params := reconcileParams()
		params.Renderer = &stubRenderer{result: unknownKindRenderResult(namespace)}
		ensureFinalizer(params, nn)

		By("first reconcile fails the apply transiently")
		result, err := opmreconcile.ReconcileModuleInstance(ctx, params, ctrl.Request{NamespacedName: nn})
		Expect(err).NotTo(HaveOccurred(), "transient failures return nil error with backoff")
		Expect(result.RequeueAfter).To(BeNumerically(">", 0))
		Expect(result.RequeueAfter).To(BeNumerically("<=", opmreconcile.BackoffMaxDelay),
			"transient failures use the fast exponential backoff, not the stalled recheck")

		var blocked releasesv1alpha1.ModuleInstance
		Expect(k8sClient.Get(ctx, nn, &blocked)).To(Succeed())
		ready := apimeta.FindStatusCondition(blocked.Status.Conditions, status.ReadyCondition)
		Expect(ready).NotTo(BeNil())
		Expect(ready.Status).To(Equal(metav1.ConditionFalse))
		Expect(ready.Reason).To(Equal(status.ApplyFailedReason))
		Expect(apimeta.FindStatusCondition(blocked.Status.Conditions, status.StalledCondition)).To(BeNil(),
			"transient failures do not stall")
		Expect(blocked.Status.NextRetryAt).NotTo(BeNil())

		By("second reconcile recovers once the render output is applicable")
		params.Renderer = &stubRenderer{}
		result, err = opmreconcile.ReconcileModuleInstance(ctx, params, ctrl.Request{NamespacedName: nn})
		Expect(err).NotTo(HaveOccurred())
		Expect(result.RequeueAfter).To(BeZero())

		var recovered releasesv1alpha1.ModuleInstance
		Expect(k8sClient.Get(ctx, nn, &recovered)).To(Succeed())
		ready = apimeta.FindStatusCondition(recovered.Status.Conditions, status.ReadyCondition)
		Expect(ready).NotTo(BeNil())
		Expect(ready.Status).To(Equal(metav1.ConditionTrue))
		Expect(recovered.Status.NextRetryAt).To(BeNil(), "nextRetryAt is cleared on success")

		var cm corev1.ConfigMap
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: "test-module", Namespace: namespace}, &cm)).To(Succeed())

		// Cleanup
		Expect(k8sClient.Delete(ctx, &cm)).To(Succeed())
		cleanupInstance(nn)
	})

	// Validates suspend → unsuspend resumes the full reconcile at the
	// reconcile-loop level (design 3.3).
	It("should resume reconciliation when suspend is cleared", func() {
		mrName := "suspend-recover-mr"
		mr := &releasesv1alpha1.ModuleInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      mrName,
				Namespace: namespace,
			},
			Spec: releasesv1alpha1.ModuleInstanceSpec{
				Suspend: true,
				Module: releasesv1alpha1.ModuleReference{
					Path:    "opmodel.dev/test/module",
					Version: "v0.1.0",
				},
				Prune:  true,
				Values: &releasesv1alpha1.RawValues{},
			},
		}
		mr.Spec.Values.Raw = []byte(`{"message": "hello"}`)
		Expect(k8sClient.Create(ctx, mr)).To(Succeed())
		nn := types.NamespacedName{Name: mrName, Namespace: namespace}

		params := reconcileParams()
		ensureFinalizer(params, nn)

		By("reconcile while suspended returns early without applying")
		result, err := opmreconcile.ReconcileModuleInstance(ctx, params, ctrl.Request{NamespacedName: nn})
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(ctrl.Result{}))

		var suspended releasesv1alpha1.ModuleInstance
		Expect(k8sClient.Get(ctx, nn, &suspended)).To(Succeed())
		ready := apimeta.FindStatusCondition(suspended.Status.Conditions, status.ReadyCondition)
		Expect(ready).NotTo(BeNil())
		Expect(ready.Status).To(Equal(metav1.ConditionFalse))
		Expect(ready.Reason).To(Equal(status.SuspendedReason))

		var cm corev1.ConfigMap
		err = k8sClient.Get(ctx, types.NamespacedName{Name: "test-module", Namespace: namespace}, &cm)
		Expect(err).To(HaveOccurred(), "no resources applied while suspended")

		By("clearing suspend resumes the full reconcile")
		var current releasesv1alpha1.ModuleInstance
		Expect(k8sClient.Get(ctx, nn, &current)).To(Succeed())
		current.Spec.Suspend = false
		Expect(k8sClient.Update(ctx, &current)).To(Succeed())

		result, err = opmreconcile.ReconcileModuleInstance(ctx, params, ctrl.Request{NamespacedName: nn})
		Expect(err).NotTo(HaveOccurred())
		Expect(result.RequeueAfter).To(BeZero())

		var resumed releasesv1alpha1.ModuleInstance
		Expect(k8sClient.Get(ctx, nn, &resumed)).To(Succeed())
		ready = apimeta.FindStatusCondition(resumed.Status.Conditions, status.ReadyCondition)
		Expect(ready).NotTo(BeNil())
		Expect(ready.Status).To(Equal(metav1.ConditionTrue))
		Expect(ready.Reason).NotTo(Equal(status.SuspendedReason))

		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: "test-module", Namespace: namespace}, &cm)).To(Succeed())
		Expect(cm.Data["message"]).To(Equal("hello"))

		// Cleanup
		Expect(k8sClient.Delete(ctx, &cm)).To(Succeed())
		cleanupInstance(nn)
	})
	// A controller that dies after the first apply and before its status
	// commit leaves the rendered objects live and a status with no inventory,
	// no lastApplied* digests and no lastAppliedInputs. The next reconcile
	// takes the objects in with their UIDs unchanged.
	It("should converge after a crash between the first apply and the status commit", func() {
		mrName := "crash-first-apply-mr"
		createModuleInstance(mrName) // Prune: true
		nn := types.NamespacedName{Name: mrName, Namespace: namespace}

		params := reconcileParams()
		ensureFinalizer(params, nn)

		By("first reconcile applies test-module and reaches Ready")
		result, err := opmreconcile.ReconcileModuleInstance(ctx, params, ctrl.Request{NamespacedName: nn})
		Expect(err).NotTo(HaveOccurred())
		Expect(result.RequeueAfter).To(BeZero())
		expectConverged(nn, "ConfigMap/test-module")

		var cm corev1.ConfigMap
		cmKey := types.NamespacedName{Name: "test-module", Namespace: namespace}
		Expect(k8sClient.Get(ctx, cmKey, &cm)).To(Succeed())
		uidBefore := cm.UID
		Expect(uidBefore).NotTo(BeEmpty())

		By("staging the status a crash before the status commit leaves")
		stageCrashStatus(nn, nil)

		By("reconciling again, as the restarted controller does")
		result, err = opmreconcile.ReconcileModuleInstance(ctx, params, ctrl.Request{NamespacedName: nn})
		Expect(err).NotTo(HaveOccurred())

		By("test-module is taken in with its UID unchanged")
		expectConverged(nn, "ConfigMap/test-module")
		Expect(k8sClient.Get(ctx, cmKey, &cm)).To(Succeed())
		Expect(cm.UID).To(Equal(uidBefore), "test-module is not deleted and created again")
		Expect(result.RequeueAfter).To(BeZero())

		// Cleanup
		Expect(k8sClient.Delete(ctx, &cm)).To(Succeed())
		cleanupInstance(nn)
	})

	// A controller that dies with a prune pending leaves an inventory that
	// lists an object the current render no longer produces, and no
	// lastApplied* digests. The next reconcile prunes that object, and the
	// rendered objects keep their UIDs.
	It("should prune a pending stale object after a crash", func() {
		mrName := "crash-pending-prune-mr"
		createModuleInstance(mrName) // Prune: true
		nn := types.NamespacedName{Name: mrName, Namespace: namespace}

		params := reconcileParams()
		ensureFinalizer(params, nn)

		By("first reconcile applies test-module and reaches Ready")
		result, err := opmreconcile.ReconcileModuleInstance(ctx, params, ctrl.Request{NamespacedName: nn})
		Expect(err).NotTo(HaveOccurred())
		Expect(result.RequeueAfter).To(BeZero())
		expectConverged(nn, "ConfigMap/test-module")

		var mi releasesv1alpha1.ModuleInstance
		Expect(k8sClient.Get(ctx, nn, &mi)).To(Succeed())
		Expect(mi.Spec.Prune).To(BeTrue())
		Expect(mi.Status.Inventory.Entries).To(HaveLen(1))
		rendered := mi.Status.Inventory.Entries[0]
		Expect(rendered.Name).To(Equal("test-module"))

		var cm corev1.ConfigMap
		cmKey := types.NamespacedName{Name: "test-module", Namespace: namespace}
		Expect(k8sClient.Get(ctx, cmKey, &cm)).To(Succeed())
		uidBefore := cm.UID

		By("creating stale-object live with the instance's ownership labels")
		staleKey := types.NamespacedName{Name: "stale-object", Namespace: namespace}
		stale := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      staleKey.Name,
				Namespace: staleKey.Namespace,
				Labels: map[string]string{
					labels.ManagedBy:          labels.ManagedByController,
					labels.ModuleInstanceUUID: mi.Status.InstanceUUID,
				},
			},
			Data: map[string]string{"stale": "true"},
		}
		Expect(k8sClient.Create(ctx, stale)).To(Succeed())

		By("staging an inventory of test-module and stale-object with the digests cleared")
		staleEntry := rendered
		staleEntry.Name = staleKey.Name
		stageCrashStatus(nn, &releasesv1alpha1.Inventory{
			Count:   2,
			Entries: []releasesv1alpha1.InventoryEntry{rendered, staleEntry},
		})

		By("reconciling again, as the restarted controller does")
		result, err = opmreconcile.ReconcileModuleInstance(ctx, params, ctrl.Request{NamespacedName: nn})
		Expect(err).NotTo(HaveOccurred())

		By("stale-object is pruned and absent from the inventory")
		expectGone(staleKey, &corev1.ConfigMap{}, "stale-object is deleted")
		expectConverged(nn, "ConfigMap/test-module")
		Expect(k8sClient.Get(ctx, nn, &mi)).To(Succeed())
		Expect(inventoryNames(mi.Status.Inventory)).NotTo(ContainElement("ConfigMap/" + staleKey.Name))

		By("test-module keeps its UID")
		Expect(k8sClient.Get(ctx, cmKey, &cm)).To(Succeed())
		Expect(cm.UID).To(Equal(uidBefore), "test-module is not deleted and created again")
		Expect(result.RequeueAfter).To(BeZero())

		// Cleanup
		Expect(k8sClient.Delete(ctx, &cm)).To(Succeed())
		Expect(client.IgnoreNotFound(k8sClient.Delete(ctx, &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: staleKey.Name, Namespace: staleKey.Namespace},
		}))).To(Succeed())
		cleanupInstance(nn)
	})
})
