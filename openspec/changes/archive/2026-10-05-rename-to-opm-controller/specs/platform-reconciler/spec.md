## MODIFIED Requirements

### Requirement: Generate and build the singleton Platform on reconcile

The controller SHALL reconcile the `Platform` named `cluster` by generating its platform CUE module on the controller's own disk, building it through the kernel's shape-gated platform loader, and recording the generated module (its package identity, directory, built platform) together with the resolved skew policy (`spec.skewPolicy`, `Warn` when unset) in the process-local store. The reconciler SHALL reconcile only the object named `cluster`; any other name SHALL be ignored without error. No second, derived copy of the platform exists.

The package the reconcile generates SHALL be a function of exactly two inputs — the CR's spec and the set of accepted-and-active `TransformerRegistration` claims — and the reconciler SHALL be woken by a change to either (enhancement 0015 D13). The waking event's content SHALL NOT be an input: the reconcile computes from the current state of both, so a stale, duplicated or reordered event yields the package the current state implies. A reconcile whose inputs resolve to the identity the store already holds, for a module directory that still exists, SHALL skip regeneration and still report the outcome on status; this is what bounds a burst of claim activations to one build rather than one per claim.

#### Scenario: Valid platform is generated

- **WHEN** a `Platform` named `cluster` with resolvable pins is applied
- **THEN** the reconciler generates and builds its module
- **AND** the store's record carries the package identity, the module directory and the skew policy

#### Scenario: Non-cluster object ignored

- **WHEN** the reconciler is triggered for an object whose name is not `cluster`
- **THEN** it performs no generation and returns without error

#### Scenario: A claim change wakes the reconciler

- **WHEN** a `TransformerRegistration` becomes active with no edit to the `Platform` CR
- **THEN** the reconciler is woken and records a package built from the new active-claim set

#### Scenario: An unchanged tuple regenerates nothing

- **WHEN** the reconciler is woken repeatedly while the CR spec and the active-claim set are unchanged
- **THEN** the held package is not rewritten and the CR still reports `Ready=True` with reason `Generated`

### Requirement: Build failures requeue on a bounded interval

When closure derivation, generation or the build fails, the `PlatformReconciler` SHALL requeue the `Platform` after a bounded interval rather than waiting for a spec change; no such failure is terminal. The reconciler SHALL set the failure reason (`BuildFailed` or `GenerateFailed`) and SHALL preserve any previously recorded good module. Each reconcile SHALL resolve module files for closure derivation through a module-file source built for that reconcile, so that no module-file lookup failure met by an earlier reconcile's closure derivation decides a later reconcile's closure derivation. A module-file source injected into the reconciler (tests) SHALL be used as given.

#### Scenario: Failure requeues instead of stalling indefinitely

- **WHEN** the build fails for the `cluster` Platform
- **THEN** the reconcile result carries a non-zero `RequeueAfter`, the status is `Ready=False`, and the last-good record, if any, is still held

#### Scenario: Recovery without a spec change

- **WHEN** a Platform is in `BuildFailed` and the registry condition clears with no change to the spec
- **THEN** a subsequent automatic reconcile builds successfully and sets `Ready=True` (reason `Generated`)

#### Scenario: A transient registry failure does not outlive its reconcile

- **WHEN** a reconcile of the `cluster` Platform ends `Ready=False` with reason `BuildFailed` because the registry refused the catalog's module file, and the registry then serves it again, with no change to the Platform spec, the reconciler's configuration or the controller process
- **THEN** the next reconcile fetches the module file from the registry, builds successfully, and sets `Ready=True` (reason `Generated`) at the same generation

## REMOVED Requirements

### Requirement: Reconciler stamps status.operatorVersion on every outcome

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-021); this requirement is restated under the controller's name.

**Migration**: See the requirement "Reconciler stamps status.controllerVersion on every outcome" in this capability.

## ADDED Requirements

### Requirement: Reconciler stamps status.controllerVersion on every outcome

`PlatformReconciler` SHALL set `status.controllerVersion` to the running controller's version on every status patch it makes, on `Generated`, `BuildFailed` and `GenerateFailed` alike.

#### Scenario: Stamped on success and failure

- **WHEN** the singleton Platform reconciles to any outcome
- **THEN** `status.controllerVersion` equals the running controller's version string
