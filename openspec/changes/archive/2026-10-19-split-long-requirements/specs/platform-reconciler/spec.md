## REMOVED Requirements

### Requirement: Generate and build the singleton Platform on reconcile
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "The Platform named cluster is generated, built and recorded in the store", "Only the Platform named cluster is reconciled", "The generated package is a function of the spec and the active claims", "The waking event's content is not an input to the generated package", "A reconcile whose inputs resolve to the held identity skips regeneration".

### Requirement: Surface build outcome on status
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "A successful build sets Ready True with reason Generated", "A build failure sets Ready False with reason BuildFailed", "A write failure sets Ready False with reason GenerateFailed", "A contract collision sets Ready False with reason ContractCollisions", "An over-subscribed contract sets Ready False with reason OverSubscribedContracts", "Comparable predicates set Ready False with reason ComparablePredicates", "When several refusals hold the first reason applies", "Every outcome sets observedGeneration and a failure keeps the last good module".

### Requirement: Single-slot identity-keyed store for concurrent read
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "The store holds at most one generated-module record", "Readers lease the record and superseded directories are removed only when unleased", "The package identity is derived from the generation and the active claims".

### Requirement: Build failures requeue on a bounded interval
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "A build failure requeues after a bounded interval and keeps the last good module", "Each reconcile resolves module files through a source built for it".

### Requirement: Transient failures retry faster than semantic failures
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "Transient and other failures are requeued on different intervals", "A failure is transient when its error chain holds a registry fetch failure", "The failure class picks the requeue interval only".

### Requirement: Subscription mapping and build failures
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "A subscription version is used verbatim as the pin and the stamp", "A subscription that cannot be built surfaces as BuildFailed".

## ADDED Requirements

### Requirement: The Platform named cluster is generated, built and recorded in the store
The operator SHALL reconcile the `Platform` named `cluster` by generating its platform CUE module on the operator's own disk, building it through the kernel's shape-gated platform loader, and recording the generated module (its package identity, directory, built platform) together with the resolved skew policy (`spec.skewPolicy`, `Warn` when unset) in the process-local store.
No second, derived copy of the platform exists.

#### Scenario: Valid platform is generated

- **WHEN** a `Platform` named `cluster` with resolvable pins is applied
- **THEN** the reconciler generates and builds its module
- **AND** the store's record carries the package identity, the module directory and the skew policy

### Requirement: Only the Platform named cluster is reconciled
The reconciler SHALL reconcile only the object named `cluster`; any other name SHALL be ignored without error.

#### Scenario: Non-cluster object ignored

- **WHEN** the reconciler is triggered for an object whose name is not `cluster`
- **THEN** it performs no generation and returns without error

### Requirement: The generated package is a function of the spec and the active claims
The package the reconcile generates SHALL be a function of exactly two inputs: the CR's spec and the set of accepted-and-active `TransformerRegistration` claims, and the reconciler SHALL be woken by a change to either (enhancement 0015, "The platform package regenerates edge-triggered and level-computed, keyed by generation plus the active-claim set; blast radius accepted").

#### Scenario: A claim change wakes the reconciler

- **WHEN** a `TransformerRegistration` becomes active with no edit to the `Platform` CR
- **THEN** the reconciler is woken and records a package built from the new active-claim set

### Requirement: The waking event's content is not an input to the generated package
The package the reconcile generates SHALL be a function of exactly two inputs: the CR's spec and the set of accepted-and-active `TransformerRegistration` claims.
The waking event's content SHALL NOT be an input: the reconcile computes from the current state of both, so a stale, duplicated or reordered event yields the package the current state implies.

#### Scenario: A stale event yields the package the current state implies
- **WHEN** the reconciler is woken by a stale, duplicated or reordered event
- **THEN** the reconcile computes from the current state of the CR's spec and the active claims and yields the package the current state implies

### Requirement: A reconcile whose inputs resolve to the held identity skips regeneration
A reconcile whose inputs resolve to the identity the store already holds, for a module directory that still exists, SHALL skip regeneration and still report the outcome on status; this is what bounds a burst of claim activations to one build rather than one per claim.

#### Scenario: An unchanged tuple regenerates nothing

- **WHEN** the reconciler is woken repeatedly while the CR spec and the active-claim set are unchanged
- **THEN** the held package is not rewritten and the CR still reports `Ready=True` with reason `Generated`

### Requirement: A successful build sets Ready True with reason Generated
The reconciler SHALL record the outcome on the `Platform` status: `Ready=True` with reason `Generated` on success. On success the reconciler SHALL also write the non-gating `ContractsFulfilled` condition from the effective package's inventory.

#### Scenario: Success sets Ready and observedGeneration

- **WHEN** generation and build succeed for generation N and the built inventory is routable and discriminated
- **THEN** `status.conditions` carries `Ready=True` (reason `Generated`), a `ContractsFulfilled` condition, and `status.observedGeneration == N`

### Requirement: A build failure sets Ready False with reason BuildFailed
The reconciler SHALL record the outcome on the `Platform` status: `Ready=False` with reason `BuildFailed` (a dependency did not resolve, the module did not build, or its contract inventory could not be read; the message names the dependency, the registry entry or the field).

#### Scenario: Build failure surfaces structured error

- **WHEN** a pinned build does not exist or an entry's key disagrees with its imported catalog
- **THEN** `status.conditions` carries `Ready=False` (reason `BuildFailed`) with a message naming the path and version or the registry entry

### Requirement: A write failure sets Ready False with reason GenerateFailed
The reconciler SHALL record the outcome on the `Platform` status: `Ready=False` with reason `GenerateFailed` (the module could not be written).

#### Scenario: A module that could not be written surfaces as GenerateFailed
- **WHEN** the generated module could not be written
- **THEN** `status.conditions` carries `Ready=False` with reason `GenerateFailed`

### Requirement: A contract collision sets Ready False with reason ContractCollisions
The reconciler SHALL record the outcome on the `Platform` status: `Ready=False` with reason `ContractCollisions` (the built platform's inventory reports contract keys that more than one enabled registry entry defines; the message names each colliding key and the registry entries defining it).

#### Scenario: A refusal surfaces a structured verdict

- **WHEN** the module builds and its inventory reports a colliding contract key, or is not routable, or is not discriminated
- **THEN** `status.conditions` carries `Ready=False` with reason `ContractCollisions`, `OverSubscribedContracts` or `ComparablePredicates` and a message naming the contracts involved and the defining registry entries, the providing registry entries or the comparable transformers

### Requirement: An over-subscribed contract sets Ready False with reason OverSubscribedContracts
The reconciler SHALL record the outcome on the `Platform` status: `Ready=False` with reason `OverSubscribedContracts` (the built platform's inventory is not routable; the message names each over-subscribed contract, its defining catalog when one is enabled and the registry entries providing it).

#### Scenario: A platform that is not routable surfaces as OverSubscribedContracts
- **WHEN** the built platform's inventory is not routable
- **THEN** `status.conditions` carries `Ready=False` with reason `OverSubscribedContracts` and a message naming each over-subscribed contract, its defining catalog when one is enabled and the registry entries providing it

### Requirement: Comparable predicates set Ready False with reason ComparablePredicates
The reconciler SHALL record the outcome on the `Platform` status: `Ready=False` with reason `ComparablePredicates` (the built platform's inventory is not discriminated; the message names each comparable pair and the contracts it shares).

#### Scenario: A platform that is not discriminated surfaces as ComparablePredicates
- **WHEN** the built platform's inventory is not discriminated
- **THEN** `status.conditions` carries `Ready=False` with reason `ComparablePredicates` and a message naming each comparable pair and the contracts it shares

### Requirement: When several refusals hold the first reason applies
When several refusals hold, the reason SHALL be the first of `ContractCollisions`, `OverSubscribedContracts`, `ComparablePredicates` that applies, and the message SHALL carry every finding.

#### Scenario: Several refusals hold at once
- **WHEN** several refusals hold
- **THEN** the reason is the first of `ContractCollisions`, `OverSubscribedContracts`, `ComparablePredicates` that applies, and the message carries every finding

### Requirement: Every outcome sets observedGeneration and a failure keeps the last good module
`status.observedGeneration` SHALL be set on every outcome.
A failure or a refusal SHALL NOT overwrite a previously recorded good module.

#### Scenario: Failure preserves last-good generated platform

- **WHEN** a previously recorded good module exists and a subsequent reconcile fails or is refused
- **THEN** the store still returns the last-good record and the failure or refusal is reflected only on the Ready condition

### Requirement: The store holds at most one generated-module record
The store SHALL hold at most one current generated-module record, keyed on the package identity it was built for, safe for a single writer and many readers.

#### Scenario: Concurrent reads are safe

- **WHEN** multiple goroutines lease and release the record while the reconciler records a new package
- **THEN** reads return a consistent record or absence without data races, and lease counts stay exact

### Requirement: Readers lease the record and superseded directories are removed only when unleased
Readers SHALL obtain the record through a lease that they release when their render completes; the store SHALL report which identities are leased.
Superseded module directories SHALL be removed only when no lease holds their identity.

#### Scenario: A leased package survives the swap

- **WHEN** a render holds a lease on one identity while a package under another identity is recorded
- **THEN** the leased package's directory is not removed until the lease is released, and is removed by a later reconcile once unleased

### Requirement: The package identity is derived from the generation and the active claims
The identity SHALL be derived from the `Platform` CR's generation and the sorted set of active claims' catalog coordinates, so two packages built from the same inputs carry the same identity and any change to either input yields a different one.
**Source**: enhancement 0015, "The platform package regenerates edge-triggered and level-computed, keyed by generation plus the active-claim set; blast radius accepted" and "The generated-platform record survives the render-parity enhancement and is held per package identity, not per Platform CR generation".

#### Scenario: A spec change replaces the slot

- **WHEN** the Platform spec changes (new generation M) and builds successfully
- **THEN** the store returns the record for the identity carrying generation M and the prior record is no longer current

#### Scenario: A claim change replaces the slot without a spec change

- **WHEN** a claim becomes active while the CR's generation is unchanged
- **THEN** the store returns a record under a different identity than the one it held

### Requirement: A build failure requeues after a bounded interval and keeps the last good module
When closure derivation, generation or the build fails, the `PlatformReconciler` SHALL requeue the `Platform` after a bounded interval rather than waiting for a spec change; no such failure is terminal.
The reconciler SHALL set the failure reason (`BuildFailed` or `GenerateFailed`) and SHALL preserve any previously recorded good module.

#### Scenario: Failure requeues instead of stalling indefinitely

- **WHEN** the build fails for the `cluster` Platform
- **THEN** the reconcile result carries a non-zero `RequeueAfter`, the status is `Ready=False`, and the last-good record, if any, is still held

#### Scenario: Recovery without a spec change

- **WHEN** a Platform is in `BuildFailed` and the registry condition clears with no change to the spec
- **THEN** a subsequent automatic reconcile builds successfully and sets `Ready=True` (reason `Generated`)

### Requirement: Each reconcile resolves module files through a source built for it
Each reconcile SHALL resolve module files for closure derivation through a module-file source built for that reconcile, so that no module-file lookup failure met by an earlier reconcile's closure derivation decides a later reconcile's closure derivation.
A module-file source injected into the reconciler (tests) SHALL be used as given.

#### Scenario: A transient registry failure does not outlive its reconcile

- **WHEN** a reconcile of the `cluster` Platform ends `Ready=False` with reason `BuildFailed` because the registry refused the catalog's module file, and the registry then serves it again, with no change to the Platform spec, the reconciler's configuration or the operator process
- **THEN** the next reconcile fetches the module file from the registry, builds successfully, and sets `Ready=True` (reason `Generated`) at the same generation

### Requirement: Transient and other failures are requeued on different intervals
The reconciler SHALL requeue a transient failure on a short interval and every other failure on the long stalled-recheck interval.
The reconciler SHALL classify by type, never by message text; an unrecognized cause SHALL default to the long interval.

#### Scenario: Semantic or unknown cause retries slowly

- **WHEN** the build fails for a cause that is not a registry fetch failure, for example a generated module that does not build, a stored registry entry without its version, or an unclassifiable cause
- **THEN** the reconcile requeues on the long stalled-recheck interval

### Requirement: A failure is transient when its error chain holds a registry fetch failure
A failure SHALL be transient when its error chain holds the library's typed registry fetch failure (`*oerrors.FetchError`, of any `Kind`) and none of the typed terminal causes, the same classification the ModuleInstance and ModulePackage reconcilers use (`reconcile-backoff`, "Registry fetch failures are transient wherever they occur"), or when it holds `context.DeadlineExceeded`.

#### Scenario: Transient cause retries quickly

- **WHEN** the build fails because a registry fetch failed: the registry is unreachable or times out, does not hold a pinned catalog version, or refuses the credentials
- **THEN** the reconcile requeues on the short interval

### Requirement: The failure class picks the requeue interval only
The reconciler SHALL requeue a transient failure on a short interval and every other failure on the long stalled-recheck interval.
The shared classification picks the Platform's requeue interval only: the Platform's condition stays `Stalled=True` for every failure, unchanged by this rule.
The condition the failure sets is unchanged by its class.

#### Scenario: The failure class does not change the condition
- **WHEN** a failure is transient or not
- **THEN** the Platform's condition stays `Stalled=True`, unchanged by the class of the failure

### Requirement: A subscription version is used verbatim as the pin and the stamp
Each subscription's `version` SHALL be used verbatim as the generated module's pin and as the entry's stamped expected version.

#### Scenario: A subscription version is used as written
- **WHEN** a subscription carries a version
- **THEN** the version is used verbatim as the generated module's pin and as the entry's stamped expected version

### Requirement: A subscription that cannot be built surfaces as BuildFailed
A subscription with an empty version (a stored object predating the required field) SHALL surface as `BuildFailed` naming the path before any registry I/O.
A version absent from the registry SHALL surface as `BuildFailed` naming the path and version.
A pin whose bytes disagree with the stamp SHALL surface as `BuildFailed` naming the registry entry.
All retain the stalled recheck interval.

#### Scenario: Missing version surfaces as BuildFailed

- **WHEN** the stored Platform carries a subscription without a version
- **THEN** the Platform's Ready condition is False with reason `BuildFailed` and a message naming the subscription path

#### Scenario: Named build not published

- **WHEN** a subscription names a version the registry does not have
- **THEN** `BuildFailed` names the path and version

#### Scenario: Identity mismatch surfaces with both values

- **WHEN** a pinned catalog's bytes disagree with the stamped expected version or the entry's key disagrees with the catalog's declared module path
- **THEN** `BuildFailed` carries the conflict naming the registry entry and both values
