## MODIFIED Requirements

### Requirement: The render input key names what a render is a function of

The controller SHALL compute a render input key for a ModuleInstance or a ModulePackage from six parts:

- the source digest: for a ModuleInstance the digest of `spec.module.path@spec.module.version`, for a ModulePackage the resolved Flux artifact digest;
- the config digest: the digest of `spec.values` for a ModuleInstance, the digest of empty input for a ModulePackage;
- the platform package identity, in the string form `Platform.status.packageIdentity` carries;
- the catalog skew policy resolved from `Platform.spec.skewPolicy`, spelled as the API spells it (`Warn` when unset);
- the controller version, as the controller publishes it to `Platform.status.controllerVersion`;
- the library version the controller binary was built with.

The key's digest SHALL be `sha256:<hex>` over a fixed, versioned encoding of the parts in a fixed order, so the same parts always give the same digest and any change to a part gives a different one. A key with an empty part is incomplete: it SHALL never be recorded and never match.

Source: owner decision g3 of the kernel-plan walkthrough (2026-10-02): "Skip render when a pre-render input key (source, values, platform identity, controller/library version, skew) matches."

#### Scenario: Any part changes the digest

- **WHEN** two keys differ in exactly one part
- **THEN** their digests differ

#### Scenario: The same parts give the same digest

- **WHEN** a key is computed twice from the same parts
- **THEN** both digests are equal

#### Scenario: A missing part makes the key incomplete

- **WHEN** the Platform has no `status.packageIdentity`, or no Platform named `cluster` exists, or the binary carries no library version
- **THEN** the key is incomplete and the reconcile renders

### Requirement: A reconcile skips its render when its inputs are unchanged

Before rendering, a ModuleInstance or ModulePackage reconcile SHALL compute the key from the object (for a ModulePackage, after resolving its source), `Platform.status.packageIdentity`, the resolved `Platform.spec.skewPolicy` of the `cluster` Platform, and the running controller and library versions. It SHALL skip the render when, and only when, all of these hold:

- the drift render interval is greater than zero;
- `status.lastAppliedInputs` is set and its `renderedAt` is less than the drift render interval ago and not in the future;
- `Ready` is `True` with reason `ReconciliationSucceeded`;
- `status.observedGeneration` equals `metadata.generation`;
- the key is complete and its digest equals `status.lastAppliedInputs.digest`;
- for a ModulePackage, the source just resolved (ref, artifact revision, digest and URL) equals `status.source`.

A skipped reconcile SHALL NOT take a render slot, lease the platform, fetch the source artifact, render, run drift detection, apply, prune or emit an event. It is not a reconcile attempt: it records no outcome, no history and no `lastAttempted*`, and it SHALL NOT move `renderedAt`. It SHALL judge health over `status.inventory` (`instance-health`), and its only status patch SHALL be a changed `Healthy` condition. A skipped ModuleInstance reconcile SHALL return without a requeue unless health asks for one; a skipped ModulePackage reconcile SHALL requeue after `spec.interval`, as a `NoOp` does, or sooner when health asks for it.

Rationale: a render is the most expensive step of a reconcile, and when every input that shapes it is unchanged it can only reproduce what the last apply wrote; drift is still caught by the render that the drift render interval forces. The `Healthy` condition is the only status a skip writes; it reads the inventory objects and, under impersonation, the ServiceAccount, because a requeue that waits for a rollout must be able to observe it without rendering.

#### Scenario: Unchanged inputs skip the render

- **WHEN** a Ready, rolled-out (`Healthy=True`) ModuleInstance whose recorded key matches its current inputs, rendered 5 minutes ago, is reconciled with a 30-minute interval
- **THEN** the renderer is not called
- **AND** the instance's status is not patched

#### Scenario: A changed input renders

- **WHEN** the Platform's `status.packageIdentity` moves from `gen-3` to `gen-4` and the instance is reconciled within the interval
- **THEN** the instance renders

#### Scenario: A failed attempt renders on the next reconcile

- **WHEN** an instance's last attempt failed (`Ready=False`) and its spec was reverted to the inputs of the recorded key
- **THEN** the next reconcile renders

#### Scenario: A spec edit outside the key renders

- **WHEN** only a field outside the key changes (`spec.prune`, `spec.serviceAccountName`, `spec.rollout`), so `metadata.generation` moves and the key does not
- **THEN** the reconcile renders and the commit writes `status.observedGeneration`

#### Scenario: A new revision with the same digest renders once

- **WHEN** a ModulePackage's Flux source moves to a new revision whose artifact digest is unchanged, within the drift render interval
- **THEN** the package renders, the outcome is `NoOp`, and `status.source.artifactRevision` names the new revision

#### Scenario: A skipped package keeps its interval

- **WHEN** a ModulePackage whose artifact digest and other inputs are unchanged is requeued by its `spec.interval` within the drift render interval
- **THEN** the artifact is not fetched and the package is not rendered
- **AND** the reconcile requeues after `spec.interval`

#### Scenario: A skipped instance that has not rolled out requeues

- **WHEN** a Ready ModuleInstance whose recorded key matches its current inputs is `NotRolledOut` and its render is skipped
- **THEN** the renderer is not called, health is judged, and the reconcile requeues after the health requeue

## REMOVED Requirements

### Requirement: An operator or library upgrade renders every object once

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-019); this requirement is restated under the controller's name.

**Migration**: See the requirement "A controller or library upgrade renders every object once" in this capability.

## ADDED Requirements

### Requirement: A controller or library upgrade renders every object once

Because the controller and library versions are key parts, the first reconcile of each object after either version changes SHALL render. The rule rests on every controller release changing `version.Version`; a development image built without a version bump keeps its key. A change to how the controller computes a stored digest (`lastApplied*` digests, `status.inventory.digest`) SHALL ship with a controller version change or a library version change, so that render finds the stored digests out of date and applies once. Starting the same controller and library again (a restart) SHALL NOT by itself make an object render while its key matches and the interval has not passed.

#### Scenario: A newer controller renders and applies once

- **WHEN** an instance's key was recorded by controller `v1.0.0-beta.9`, the running controller is `v1.0.0-beta.10`, and the render digest the new controller computes differs from `lastAppliedRenderDigest`
- **THEN** the reconcile renders and applies
- **AND** `status.lastAppliedInputs` then carries the new controller's key

#### Scenario: A restart does not render unchanged objects

- **WHEN** the controller restarts, the platform store is still empty, and an instance's key matches and was recorded 10 minutes ago
- **THEN** the instance is not rendered and stays `Ready=True`, without passing through `PlatformNotReady`

#### Scenario: The library inventory digests apply once and then converge

- **WHEN** an instance or package whose `status.inventory.digest` and `lastAppliedRenderDigest` were written in the earlier controller's encoding, with a key recorded by controller `v1.0.0-test`, is reconciled by controller `v1.0.1-test` that computes them with `opm/k8s/inventory`
- **THEN** the reconcile renders once and applies once, and the stored digests become the library's digests of the same entries and render
- **AND** the next reconcile under `v1.0.1-test` skips the render, and a render after the interval ends `NoOp` without a second apply

#### Scenario: A key recorded by the same version keeps the old digests until the interval

- **WHEN** an instance holds digests in the earlier encoding and a key recorded by the running controller version
- **THEN** the reconcile skips its render, because the key holds no digest
- **AND** the first render after the interval applies once and records the new digests
