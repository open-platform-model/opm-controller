# Design

## Context

See proposal.md for the motivation. The facts this design rests on, read in the code at `main`:

- The ModuleInstance reconcile writes `status.inventory`, the `lastApplied*` digests and `lastAppliedInputs` only in its deferred status commit, and only when the reconcile applied (`internal/reconcile/moduleinstance.go`, the deferred commit). A process that dies after the apply and before that commit leaves the old values.
- `status.IsNoOp` is false when any `lastApplied*` digest is empty (`internal/status/digests.go`). The render skip needs a recorded `lastAppliedInputs` (`internal/reconcile/inputs.go`, `renderSkip.maySkip`). Clearing these fields therefore forces a full render and apply.
- `status.instanceUUID` is stored before the first write (`applyInstance` stores the identities first), so a crash after the apply keeps the UUID that the objects carry.
- The apply guard reads each object live. An object outside the inventory that is OPM-managed and carries the same UUID label is allowed and taken in, with its UID pinned (`internal/apply/guard.go`; library v1.0.0-beta.10 `opm/k8s/ownership`). Specs `ssa-apply` and `reconcile-loop-assembly` already require that a taken-in object is never deleted and created again.
- The stale set is the old inventory minus the new render (`staleEntries`). Prune deletes a stale object only when its live labels pass `ownership.CanDelete`: `app.kubernetes.io/managed-by: opm-controller` and `module-instance.opmodel.dev/uuid: <status.instanceUUID>`.
- After a restart the platform store is empty. Renders wait with `PlatformNotReady` until the Platform reconciler regenerates and re-enqueues them (specs `platform-module-generation`, `platform-gated-rendering`).
- Leader election is on and `LeaderElectionReleaseOnCancel` is off (`cmd/main.go`). A new pod waits for the old lease to expire, about 15s. The Deployment has one replica and a 10s grace period.
- The `hello` fixture v0.0.14 renders one ConfigMap, has `prune: true`, and is readable on GHCR without a login (anonymous tag list of `ghcr.io/open-platform-model/testing.opmodel.dev/modules/controller/hello`, checked 2026-11-27). The controller's built-in registry routes `testing.opmodel.dev/*` to GHCR.

## Goals / Non-Goals

**Goals:**

- Prove, at the reconcile function and on a deployed controller, the two crash windows: first apply without a status commit, and an update with a pending prune.
- Build the controller pod stop and restart fixture, and run one graceful pod delete during a reconcile as a smoke check.

**Non-Goals:**

- No change to product code, flags or manifests. No fault-injection hook in the binary.
- The parallel-instances stub stays skipped.
- No hard kill (`--grace-period=0 --force`): the staged state already proves the no-commit window, and a forced delete lets two pods overlap.

## Decisions

### The crash state is staged, not timed

A `hello` reconcile ends in seconds, so a kill timed from outside rarely lands between the apply and the commit. The test stops the controller (Deployment scaled to 0), writes the status a crash would leave through the status subresource, and starts the controller again (scaled to 1). Alternatives weighed: a timed hard kill only (where it lands is luck), and a test flag in the binary (test code in the shipped binary). The owner chose staging.

Staged first-apply crash, as a status merge patch:

```yaml
status:
  inventory: null
  lastAppliedSourceDigest: null
  lastAppliedConfigDigest: null
  lastAppliedRenderDigest: null
  lastAppliedInputs: null
```

Staged pending prune: the same patch, except that `inventory.entries` holds the ConfigMap that `hello` renders and one extra entry, `hello-stale`, a ConfigMap the test created with the two ownership labels and the instance's UUID:

```yaml
status:
  inventory:
    entries:
      - {kind: ConfigMap, v: v1, namespace: default, name: <rendered name>, component: <rendered component>}
      - {kind: ConfigMap, v: v1, namespace: default, name: hello-stale, component: <rendered component>}
```

`instanceUUID`, the conditions and `observedGeneration` stay as they are. A crash keeps them.

### What "converged" means

Each spec reads the rendered objects' UIDs before the stop and checks after the restart. The check helper is shared by the three e2e specs, and the envtest cases use the same four checks:

```go
// convergedSince waits until the instance reconciled after t, then checks
// the four outcomes.
Eventually(func(g Gomega) {
    mi := getInstance(g)
    g.Expect(mi.Status.LastAppliedAt.After(t)).To(BeTrue())       // a real reconcile ran after the restart
    g.Expect(ready(mi)).To(BeTrue())
    g.Expect(mi.Status.ObservedGeneration).To(Equal(mi.Generation))
    g.Expect(liveUIDs(g)).To(Equal(uidsBefore))                    // nothing recreated
    g.Expect(inventoryKeys(mi)).To(Equal(renderedKeys))            // inventory equals render
    g.Expect(labeledKeys(g, mi.Status.InstanceUUID)).To(BeSubsetOf(inventoryKeys(mi))) // no orphan
}, 3*time.Minute, 2*time.Second).Should(Succeed())
```

`lastAppliedAt` after the staging time is what tells a reconcile after the restart from the Ready condition that the staged status kept. Events are not asserted: `PlatformNotReady` is expected after a restart, and the owner declined an events check as flaky.

### The pod delete smoke check

The spec changes `spec.values.message` and at once runs `kubectl delete pod -l control-plane=manager -n opm-system`, with no wait. A ModuleInstance reconcile never sets `Reconciling=True` (only a panic, a render timeout or a deletion does), so no condition shows that a reconcile runs, and a `hello` reconcile ends in seconds. The delete may land during the reconcile or after it: the spec is a smoke check of a real process stop after a change, and the staged specs prove the crash windows. It then waits for a new pod to be Ready, and runs the full check with the new message in the ConfigMap.

### Where the tests live

Two seams, both agreed with the owner:

- Envtest: `test/integration/reconcile/state_recovery_test.go`, through `opmreconcile.ReconcileModuleInstance` with the suite's `stubRenderer` (one ConfigMap, `test-module`). A take-in or stale-set regression fails `task dev:test` in seconds.
- E2e: an Ordered Describe in `test/e2e/concurrent_test.go` with its own controller deploy, Platform, `hello` instance and teardown, in the pattern of the other e2e files (the deploy steps and the `LOCAL_REGISTRY` override of `e2e_test.go`, `helloInstanceManifest` from `registry_auth_test.go`). The parallel-instances `Skip` moves to its own Describe in the same file.

## Reconcile phase impact

None changed. The tests drive these phases after a restart: Platform gate (`PlatformNotReady`, then re-enqueue), render (no skip, because `lastAppliedInputs` is cleared), apply guard (take-in of the staged objects), apply (SSA, idempotent), prune (stale set from the staged inventory), status commit (inventory and digests written again).

## Research & Decisions

### Crash window

**Context**: The risk sits between the apply and the status commit, and a timed kill rarely hits it.
**Explored**: The deferred commit and the NoOp and render-skip conditions in `internal/reconcile`; the original design of this test in the archived `integration-test-coverage` change, test 6.2 ("not duplicate resources or skip prune").
**Decision**: Stage the state through the status subresource; add a graceful pod delete as a smoke check.
**Rationale**: Hits the window on every run with no product code.

### Fixture

**Context**: The e2e spec needs a module that resolves in a plain CI run.
**Explored**: `test/fixtures`, the controller's default registry, an anonymous GHCR tag list for `hello`.
**Decision**: `hello` v0.0.14, through `helloInstanceManifest`.
**Rationale**: One ConfigMap, `prune: true`, public, no new env gate.

## Risks / Trade-offs

- [A staged status may differ from what a real crash leaves] → The patch clears only the fields that the deferred commit writes, and the pod delete spec adds a real process stop.
- [A new test fails on the current code] → That is a real bug in resume. The build stops that slice and adds a fix slice for the reconcile code; the spec delta stays as written.
- [The lease wait and the Platform rebuild make the e2e slow or flaky] → Eventually windows of 3 minutes, and the restart helper waits for the new pod's Ready and the Platform's Ready before it checks the instance.
- [Scaling to 0 while a reconcile runs] → The staging starts only after the instance is Ready and the pod is gone (`kubectl wait --for=delete`).

## Slices

### S1: Crash-window cases at the reconcile function

- Delivers: Tests only: two envtest cases that stage a first-apply crash and a pending-prune crash on a ModuleInstance and reconcile again.
- Seam: `test/integration/reconcile/state_recovery_test.go`, calling `opmreconcile.ReconcileModuleInstance` with the suite's `stubRenderer`.
- Depends on: none.
- Done when: both cases pass in `task dev:test`; the first-apply case sees the same `test-module` UID and an inventory of exactly that ConfigMap; the prune case sees the stale ConfigMap deleted and absent from the inventory.

### S2: Controller restart specs on Kind

- Delivers: Tests only: the Ordered restart Describe in `test/e2e/concurrent_test.go` with the stop and restart fixture and three specs (staged first-apply crash, staged pending prune, graceful pod delete), and the parallel-instances `Skip` in its own Describe.
- Seam: `test/e2e/concurrent_test.go` on Kind, through `kubectl` against the deployed controller.
- Depends on: S1.
- Done when: the three specs pass in `task dev:e2e`; `concurrent_test.go` holds exactly one `Skip`, the parallel-instances one.
