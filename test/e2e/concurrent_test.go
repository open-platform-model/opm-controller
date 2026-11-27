//go:build e2e
// +build e2e

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

package e2e

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	releasesv1alpha1 "github.com/open-platform-model/opm-controller/api/v1alpha1"
	"github.com/open-platform-model/opm-controller/test/utils"
)

var _ = Describe("Concurrent Reconciliation", func() {
	// TODO (design 6.1): Validates multiple ModuleInstances against a shared
	// source reconcile correctly under concurrency. Requires: deployed
	// controller, shared OCIRepository, create 3 ModuleInstances referencing it
	// → Eventually all report Ready=True → assert each has its own distinct
	// managed resources and no cross-release contamination in inventories.
	It("should reconcile multiple ModuleInstances from a shared source", func() {
		Skip("TODO: requires concurrent reconcile fixture across parallel ModuleInstances")
	})
})

// restartInstanceNS and restartInstanceName name the hello fixture instance
// that the restart specs drive.
const (
	restartInstanceNS   = "default"
	restartInstanceName = "hello"
	// managerDeployment is the controller Deployment in the manager namespace.
	managerDeployment = "opm-controller-manager"
	// managerSelector selects the controller pods.
	managerSelector = "control-plane=manager"
	// ownerUUIDLabel is the label that carries the instance's UUID on every
	// object the controller applies for it.
	ownerUUIDLabel = "module-instance.opmodel.dev/uuid"
	// managedByLabel marks an object as applied by the controller.
	managedByLabel = "app.kubernetes.io/managed-by"
	// managedByController is the managedByLabel value of the controller.
	managedByController = "opm-controller"
)

// This spec proves that a deployed controller that stops between its apply
// and its status commit leaves a state that the next controller converges. The
// crash state is staged: the controller is scaled to 0, the status that a
// crash leaves is written through the status subresource, and the controller
// is scaled to 1 again. A third spec deletes the controller pod while a
// reconcile runs, as a smoke check of a real process stop.
//
// It is self-contained (own controller deploy, Platform, hello instance and
// teardown). The hello fixture resolves from GHCR through the controller's
// built-in registry, so it needs no login and no local registry.
var _ = Describe("Controller restart", Ordered, func() {
	var (
		projectDir string
		// rendered is the render of the hello instance, read from its
		// inventory once it is first Ready: one ConfigMap.
		rendered []releasesv1alpha1.InventoryEntry
	)

	BeforeAll(func() {
		var err error
		projectDir, err = utils.GetProjectDir()
		Expect(err).NotTo(HaveOccurred())

		By("creating the manager namespace")
		// Ignore an already-exists error from a prior spec in the same suite.
		_, _ = utils.Run(exec.Command("kubectl", "create", "ns", namespace))

		By("installing CRDs")
		_, err = utils.Run(exec.Command("task", "controller:crds"))
		Expect(err).NotTo(HaveOccurred(), "Failed to install CRDs")

		By("deploying the manager")
		_, err = utils.Run(exec.Command("task", "controller:deploy", fmt.Sprintf("IMG=%s", managerImage)))
		Expect(err).NotTo(HaveOccurred(), "Failed to deploy the manager")

		// Local-dev override, as in e2e_test.go: resolve from a registry
		// reachable inside Kind. CI leaves it unset and the controller uses its
		// built-in default.
		if localRegistry := os.Getenv("LOCAL_REGISTRY"); localRegistry != "" {
			By("overriding controller registry for local dev")
			patch := fmt.Sprintf(
				`[{"op":"add","path":"/spec/template/spec/containers/0/args/-","value":"--registry=%s"}]`, localRegistry)
			_, err = utils.Run(exec.Command("kubectl", "-n", namespace, "patch", "deployment",
				managerDeployment, "--type=json", "-p", patch))
			Expect(err).NotTo(HaveOccurred(), "Failed to override controller registry")
			_, err = utils.Run(exec.Command("kubectl", "-n", namespace, "rollout", "status",
				"deployment/"+managerDeployment, "--timeout=120s"))
			Expect(err).NotTo(HaveOccurred(), "registry-override rollout did not settle")
		}
		waitAvailable(namespace, managerDeployment, 3*time.Minute)

		By("applying the cluster Platform")
		_, err = utils.Run(exec.Command("kubectl", "apply", "-f",
			filepath.Join(projectDir, "config/samples/opmodel.dev_v1alpha1_platform.yaml")))
		Expect(err).NotTo(HaveOccurred(), "Failed to apply the Platform")
		waitPlatformReady()

		By("applying the hello ModuleInstance")
		applyStdin(helloInstanceManifest(projectDir))
		Eventually(func(g Gomega) {
			mi := getRestartInstance(g)
			g.Expect(instanceReady(mi)).To(BeTrue(), "hello ModuleInstance not Ready yet")
			g.Expect(mi.Status.ObservedGeneration).To(Equal(mi.Generation))
			g.Expect(mi.Status.Inventory).NotTo(BeNil())
		}, 6*time.Minute, 5*time.Second).Should(Succeed())

		mi := getRestartInstance(Default)
		Expect(mi.Spec.Prune).To(BeTrue(), "the hello fixture prunes")
		Expect(mi.Status.InstanceUUID).NotTo(BeEmpty())
		rendered = mi.Status.Inventory.Entries
		Expect(rendered).To(HaveLen(1), "hello renders one ConfigMap")
		Expect(rendered[0].Kind).To(Equal("ConfigMap"))
		Expect(rendered[0].Namespace).To(Equal(restartInstanceNS))
	})

	AfterAll(func() {
		// The ModuleInstance prunes through its applier ServiceAccount, so the
		// instance goes first and the RBAC after it (see podinfo_test.go).
		By("removing the hello ModuleInstance")
		_, _ = utils.Run(exec.Command("kubectl", "-n", restartInstanceNS, "delete", "moduleinstance",
			restartInstanceName, "--ignore-not-found", "--wait=false"))
		if _, err := utils.Run(exec.Command("kubectl", "-n", restartInstanceNS, "wait", "--for=delete",
			"moduleinstance/"+restartInstanceName, "--timeout=2m")); err != nil {
			_, _ = utils.Run(exec.Command("kubectl", "-n", restartInstanceNS, "patch", "moduleinstance",
				restartInstanceName, "--type=merge", "-p", `{"metadata":{"finalizers":null}}`))
		}
		_, _ = utils.Run(exec.Command("kubectl", "delete", "--ignore-not-found", "--wait=false", "-f",
			filepath.Join(projectDir, "test/fixtures/modules/hello/moduleinstance.yaml")))
		_, _ = utils.Run(exec.Command("kubectl", "-n", restartInstanceNS, "delete", "configmap",
			"hello-stale", "--ignore-not-found"))

		By("removing the cluster Platform and the manager")
		_, _ = utils.Run(exec.Command("kubectl", "delete", "--ignore-not-found", "--wait=false", "-f",
			filepath.Join(projectDir, "config/samples/opmodel.dev_v1alpha1_platform.yaml")))
		_, _ = utils.Run(exec.Command("task", "controller:uninstall", "IGNORE_NOT_FOUND=true"))
		_, _ = utils.Run(exec.Command("task", "controller:crds:remove", "IGNORE_NOT_FOUND=true"))
	})

	AfterEach(func() {
		if !CurrentSpecReport().Failed() {
			return
		}
		By("dumping diagnostics for the failed spec")
		for _, args := range [][]string{
			{"-n", restartInstanceNS, "get", "moduleinstance", restartInstanceName, "-o", "yaml"},
			{"-n", restartInstanceNS, "get", "configmaps", "--show-labels"},
			{"get", "platform", "cluster", "-o", "yaml"},
			{"-n", namespace, "get", "pods", "-o", "wide"},
			{"-n", namespace, "get", "events", "--sort-by=.lastTimestamp"},
			{"-n", namespace, "logs", "-l", managerSelector, "--tail=200"},
		} {
			if out, err := utils.Run(exec.Command("kubectl", args...)); err == nil {
				_, _ = fmt.Fprintf(GinkgoWriter, "--- kubectl %s ---\n%s\n", strings.Join(args, " "), out)
			}
		}
	})

	It("should converge after a crash between the first apply and the status commit", func() {
		uidsBefore := liveUIDs(Default, rendered)

		By("stopping the controller")
		stopController()

		By("staging the status a crash after the first apply leaves")
		stageInstanceStatus(`{"status":{"inventory":null,` +
			`"lastAppliedSourceDigest":null,"lastAppliedConfigDigest":null,` +
			`"lastAppliedRenderDigest":null,"lastAppliedInputs":null}}`)
		staged := getRestartInstance(Default)
		Expect(staged.Status.Inventory).To(BeNil())
		Expect(staged.Status.LastAppliedInputs).To(BeNil())
		Expect(staged.Status.LastAppliedRenderDigest).To(BeEmpty())
		stagedAt := time.Now()

		By("starting the controller again")
		startController()

		By("checking the instance converged with the objects taken in")
		expectRestartConverged(stagedAt, uidsBefore, rendered)
	})

	It("should prune a pending stale object after a crash", func() {
		uidsBefore := liveUIDs(Default, rendered)
		mi := getRestartInstance(Default)
		uuid := mi.Status.InstanceUUID

		By("stopping the controller")
		stopController()

		By("creating hello-stale with the instance's ownership labels")
		applyStdin(fmt.Appendf(nil, `apiVersion: v1
kind: ConfigMap
metadata:
  name: hello-stale
  namespace: %s
  labels:
    %s: %s
    %s: %s
data:
  message: stale
`, restartInstanceNS, managedByLabel, managedByController, ownerUUIDLabel, uuid))

		By("staging an inventory that still lists hello-stale")
		stale := releasesv1alpha1.InventoryEntry{
			Kind: "ConfigMap", Version: "v1", Namespace: restartInstanceNS,
			Name: "hello-stale", Component: rendered[0].Component,
		}
		entries, err := json.Marshal([]releasesv1alpha1.InventoryEntry{rendered[0], stale})
		Expect(err).NotTo(HaveOccurred())
		stageInstanceStatus(fmt.Sprintf(`{"status":{`+
			`"inventory":{"revision":null,"digest":null,"count":2,"entries":%s},`+
			`"lastAppliedSourceDigest":null,"lastAppliedConfigDigest":null,`+
			`"lastAppliedRenderDigest":null,"lastAppliedInputs":null}}`, entries))
		staged := getRestartInstance(Default)
		Expect(inventoryKeys(staged)).To(ConsistOf(entryKey(rendered[0]), entryKey(stale)))
		Expect(staged.Status.LastAppliedInputs).To(BeNil())
		stagedAt := time.Now()

		By("starting the controller again")
		startController()

		By("checking hello-stale is pruned and the instance converged")
		Eventually(func(g Gomega) {
			out, err := utils.Run(exec.Command("kubectl", "-n", restartInstanceNS, "get", "configmap",
				"hello-stale", "--ignore-not-found", "-o", "name"))
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(out).To(BeEmpty(), "hello-stale is deleted")
		}, 3*time.Minute, 2*time.Second).Should(Succeed())
		expectRestartConverged(stagedAt, uidsBefore, rendered)
	})

	It("should converge correctly after a controller restart mid-reconcile", func() {
		const newMessage = "hello after a controller restart"
		uidsBefore := liveUIDs(Default, rendered)
		oldPod := managerPod(Default)

		By("changing spec.values.message")
		changedAt := time.Now()
		_, err := utils.Run(exec.Command("kubectl", "-n", restartInstanceNS, "patch", "moduleinstance",
			restartInstanceName, "--type=merge",
			"-p", fmt.Sprintf(`{"spec":{"values":{"message":%q}}}`, newMessage)))
		Expect(err).NotTo(HaveOccurred())

		// No wait: a ModuleInstance reconcile sets no condition that shows it
		// runs, so the delete may land during the reconcile or after it.
		By("deleting the controller pod at once")
		_, err = utils.Run(exec.Command("kubectl", "delete", "pod", "-l", managerSelector, "-n", namespace))
		Expect(err).NotTo(HaveOccurred())

		By("waiting for a new controller pod to be Ready")
		Eventually(func(g Gomega) {
			pod := managerPod(g)
			g.Expect(pod).NotTo(Equal(oldPod), "the old pod is still the controller pod")
		}, 3*time.Minute, 2*time.Second).Should(Succeed())

		By("checking the new controller applied the new values and converged")
		expectRestartConverged(changedAt, uidsBefore, rendered)
		out, err := utils.Run(exec.Command("kubectl", "-n", rendered[0].Namespace, "get", "configmap",
			rendered[0].Name, "-o", "jsonpath={.data.message}"))
		Expect(err).NotTo(HaveOccurred())
		Expect(out).To(Equal(newMessage))
	})
})

// getRestartInstance reads the hello ModuleInstance.
func getRestartInstance(g Gomega) releasesv1alpha1.ModuleInstance {
	out, err := utils.Run(exec.Command("kubectl", "-n", restartInstanceNS, "get", "moduleinstance",
		restartInstanceName, "-o", "json"))
	g.Expect(err).NotTo(HaveOccurred())
	var mi releasesv1alpha1.ModuleInstance
	g.Expect(json.Unmarshal([]byte(out), &mi)).To(Succeed())
	return mi
}

func instanceReady(mi releasesv1alpha1.ModuleInstance) bool {
	return apimeta.IsStatusConditionTrue(mi.Status.Conditions, "Ready")
}

// stageInstanceStatus writes a merge patch to the hello instance's status
// subresource. The controller must be stopped, or it reconciles the staged
// status away at once.
func stageInstanceStatus(patch string) {
	GinkgoHelper()
	_, err := utils.Run(exec.Command("kubectl", "-n", restartInstanceNS, "patch", "moduleinstance",
		restartInstanceName, "--subresource=status", "--type=merge", "-p", patch))
	Expect(err).NotTo(HaveOccurred(), "staging the status failed")
}

// stopController scales the controller to 0 and waits until its pod is gone.
func stopController() {
	GinkgoHelper()
	_, err := utils.Run(exec.Command("kubectl", "-n", namespace, "scale", "deployment",
		managerDeployment, "--replicas=0"))
	Expect(err).NotTo(HaveOccurred())
	_, err = utils.Run(exec.Command("kubectl", "-n", namespace, "wait", "--for=delete", "pod",
		"-l", managerSelector, "--timeout=2m"))
	Expect(err).NotTo(HaveOccurred(), "the controller pod did not stop")
	Eventually(func(g Gomega) {
		out, err := utils.Run(exec.Command("kubectl", "-n", namespace, "get", "pods",
			"-l", managerSelector, "-o", "name"))
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(strings.TrimSpace(out)).To(BeEmpty(), "a controller pod still exists")
	}, 2*time.Minute, 2*time.Second).Should(Succeed())
}

// startController scales the controller to 1, waits for the rollout and the
// pod Ready, and for the Platform Ready again.
func startController() {
	GinkgoHelper()
	_, err := utils.Run(exec.Command("kubectl", "-n", namespace, "scale", "deployment",
		managerDeployment, "--replicas=1"))
	Expect(err).NotTo(HaveOccurred())
	_, err = utils.Run(exec.Command("kubectl", "-n", namespace, "rollout", "status",
		"deployment/"+managerDeployment, "--timeout=3m"))
	Expect(err).NotTo(HaveOccurred(), "the controller rollout did not settle")
	Eventually(func(g Gomega) { managerPod(g) }, 3*time.Minute, 2*time.Second).Should(Succeed())
	waitPlatformReady()
}

// managerPod returns the one controller pod that is not terminating and is
// Ready.
func managerPod(g Gomega) string {
	out, err := utils.Run(exec.Command("kubectl", "-n", namespace, "get", "pods",
		"-l", managerSelector, "-o", "json"))
	g.Expect(err).NotTo(HaveOccurred())
	var list struct {
		Items []struct {
			Metadata struct {
				Name              string       `json:"name"`
				DeletionTimestamp *metav1.Time `json:"deletionTimestamp"`
			} `json:"metadata"`
			Status struct {
				Conditions []struct {
					Type   string `json:"type"`
					Status string `json:"status"`
				} `json:"conditions"`
			} `json:"status"`
		} `json:"items"`
	}
	g.Expect(json.Unmarshal([]byte(out), &list)).To(Succeed())
	var ready []string
	for _, p := range list.Items {
		if p.Metadata.DeletionTimestamp != nil {
			continue
		}
		for _, c := range p.Status.Conditions {
			if c.Type == "Ready" && c.Status == "True" {
				ready = append(ready, p.Metadata.Name)
			}
		}
	}
	g.Expect(ready).To(HaveLen(1), "expected one Ready controller pod")
	return ready[0]
}

func waitPlatformReady() {
	GinkgoHelper()
	Eventually(func(g Gomega) {
		out, err := utils.Run(exec.Command("kubectl", "get", "platform", "cluster",
			"-o", "jsonpath={.status.conditions[?(@.type=='Ready')].status}"))
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(out).To(Equal("True"), "Platform not Ready yet")
	}, 5*time.Minute, 5*time.Second).Should(Succeed())
}

// entryKey names an inventory entry as kind/namespace/name.
func entryKey(e releasesv1alpha1.InventoryEntry) string {
	return e.Kind + "/" + e.Namespace + "/" + e.Name
}

func inventoryKeys(mi releasesv1alpha1.ModuleInstance) []string {
	if mi.Status.Inventory == nil {
		return nil
	}
	keys := make([]string, 0, len(mi.Status.Inventory.Entries))
	for _, e := range mi.Status.Inventory.Entries {
		keys = append(keys, entryKey(e))
	}
	return keys
}

// liveUIDs reads the UID of each rendered object, keyed by entryKey.
func liveUIDs(g Gomega, entries []releasesv1alpha1.InventoryEntry) map[string]string {
	uids := make(map[string]string, len(entries))
	for _, e := range entries {
		resource := strings.ToLower(e.Kind)
		if e.Group != "" {
			resource += "." + e.Group
		}
		args := []string{"get", resource, e.Name, "-o", "jsonpath={.metadata.uid}"}
		if e.Namespace != "" {
			args = append([]string{"-n", e.Namespace}, args...)
		}
		out, err := utils.Run(exec.Command("kubectl", args...))
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(out).NotTo(BeEmpty(), "no UID for %s", entryKey(e))
		uids[entryKey(e)] = out
	}
	return uids
}

// labeledKeys lists the ConfigMaps, in every namespace, that carry the
// instance's UUID label and are not being deleted. The hello fixture renders
// only ConfigMaps, and its applier may write no other kind.
func labeledKeys(g Gomega, uuid string) []string {
	out, err := utils.Run(exec.Command("kubectl", "get", "configmaps", "-A",
		"-l", ownerUUIDLabel+"="+uuid, "-o", "json"))
	g.Expect(err).NotTo(HaveOccurred())
	var list struct {
		Items []struct {
			Metadata struct {
				Name              string       `json:"name"`
				Namespace         string       `json:"namespace"`
				DeletionTimestamp *metav1.Time `json:"deletionTimestamp"`
			} `json:"metadata"`
		} `json:"items"`
	}
	g.Expect(json.Unmarshal([]byte(out), &list)).To(Succeed())
	keys := make([]string, 0, len(list.Items))
	for _, cm := range list.Items {
		if cm.Metadata.DeletionTimestamp == nil {
			keys = append(keys, "ConfigMap/"+cm.Metadata.Namespace+"/"+cm.Metadata.Name)
		}
	}
	return keys
}

// expectRestartConverged waits until the instance applied at or after since, then
// checks the converged outcome: Ready=True on the current generation, every
// rendered object with the UID it had before the stop, an inventory of
// exactly the rendered objects, and no object with the instance's UUID label
// outside the inventory.
func expectRestartConverged(since time.Time, uidsBefore map[string]string,
	rendered []releasesv1alpha1.InventoryEntry) {
	GinkgoHelper()
	// lastAppliedAt has a resolution of one second, so an apply in the same
	// second as since counts.
	since = since.Truncate(time.Second)
	renderedKeys := make([]string, 0, len(rendered))
	for _, e := range rendered {
		renderedKeys = append(renderedKeys, entryKey(e))
	}
	Eventually(func(g Gomega) {
		mi := getRestartInstance(g)
		g.Expect(mi.Status.LastAppliedAt).NotTo(BeNil())
		g.Expect(mi.Status.LastAppliedAt.Time.Before(since)).To(BeFalse(),
			"no apply after %s yet (lastAppliedAt %s)", since, mi.Status.LastAppliedAt)
		g.Expect(instanceReady(mi)).To(BeTrue(), "Ready=True")
		g.Expect(mi.Status.ObservedGeneration).To(Equal(mi.Generation), "status.observedGeneration")
		g.Expect(liveUIDs(g, rendered)).To(Equal(uidsBefore), "no rendered object is deleted and created again")
		g.Expect(inventoryKeys(mi)).To(ConsistOf(renderedKeys), "status.inventory lists exactly the rendered objects")
		g.Expect(inventoryKeys(mi)).To(ContainElements(labeledKeys(g, mi.Status.InstanceUUID)),
			"no object with the instance's UUID label is outside status.inventory")
	}, 3*time.Minute, 2*time.Second).Should(Succeed())
}
