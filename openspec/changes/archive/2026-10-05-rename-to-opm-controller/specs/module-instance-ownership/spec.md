## REMOVED Requirements

### Requirement: ModuleInstance carries an ownership marker

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-021); this requirement is restated under the controller's name.

**Migration**: See the requirement "A ModuleInstance carries an ownership marker" in this capability.

### Requirement: The operator skips CLI-owned instances before registering a finalizer

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-021); this requirement is restated under the controller's name.

**Migration**: See the requirement "The controller skips CLI-owned instances before registering a finalizer" in this capability.

### Requirement: The operator acknowledges a CLI-owned instance with a single condition

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-021); this requirement is restated under the controller's name.

**Migration**: See the requirement "The controller acknowledges a CLI-owned instance with a single condition" in this capability.

### Requirement: Ownership handoff falls through to a normal reconcile

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-021); this requirement is restated under the controller's name.

**Migration**: See the requirement "An ownership handoff falls through to a normal reconcile" in this capability.

### Requirement: The operator recognises the instance that deploys it

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-021); this requirement is restated under the controller's name.

**Migration**: See the requirement "The controller recognises the instance that deploys it" in this capability.

### Requirement: The operator never reconciles its own instance

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-021); this requirement is restated under the controller's name.

**Migration**: See the requirement "The controller never reconciles its own instance" in this capability.

### Requirement: The operator releases its finalizer from its own instance without pruning

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-021), and the owner narrowed the requirement on 2026-10-06. The finalizer is no longer stripped from a live own instance: there it stays and the refusal runs as before. It is released only when the own instance is being deleted, because an instance recognised as the controller's own only after a normal reconcile (by its recorded inventory, or after its module path was edited to the controller module) already carries it, and would otherwise stay terminating.

**Migration**: See the requirement "The controller releases its finalizer from a deleting own instance without pruning" in this capability.

## ADDED Requirements

### Requirement: A ModuleInstance carries an ownership marker

`ModuleInstanceSpec` SHALL provide an optional `owner` field of a typed enum with exactly two valid values, `cli` and `controller`, serialized as `owner` with `omitempty`. The field SHALL NOT define a CRD-level default. The API SHALL define exported constants for both values. An absent or empty `owner` SHALL be treated by the controller as controller-managed (see the skip requirement), except on the controller's own instance, which the controller never reconciles (see the requirement that the controller never reconciles its own instance).

Source: 0032:D1:R1, 0032:D1:R6.

#### Scenario: Field accepts the two enum values

- **WHEN** a `ModuleInstance` is created with `spec.owner` set to `cli` or to `controller`
- **THEN** the API server accepts it
- **AND** a value other than `cli` or `controller` is rejected by enum validation

#### Scenario: Field is optional with no default

- **WHEN** a `ModuleInstance` is created with no `spec.owner`
- **THEN** the API server accepts it
- **AND** the stored object's `spec.owner` remains empty (no value is defaulted in)

#### Scenario: Empty owner is reconciled as controller-managed

- **WHEN** a `ModuleInstance` with no `spec.owner` that is not the controller's own instance is reconciled
- **THEN** the controller SHALL register the cleanup finalizer and perform a normal render/apply reconcile (no owner-skip, no `ManagedExternally` acknowledgement)

#### Scenario: Unknown owner values cannot reach the reconciler

- **WHEN** a client attempts to create a `ModuleInstance` with `spec.owner: "future-actor"`
- **THEN** the API server SHALL reject it via enum validation, so the controller never observes an unknown owner value

### Requirement: The controller skips CLI-owned instances before registering a finalizer

When `spec.owner == cli`, the controller MUST return from `Reconcile` without rendering, applying, pruning, or running deletion cleanup, and MUST NOT add the `opmodel.dev/cleanup` finalizer. The owner check MUST occur before finalizer registration so that an instance that has always been CLI-owned never carries the controller's finalizer.

An instance that was controller-owned before its `spec.owner` became `cli` still carries the finalizer. When such an instance is being deleted, the controller MUST remove the `opmodel.dev/cleanup` finalizer, MUST NOT delete any object in the instance's inventory whatever its `spec.prune`, and MUST write no status. On a CLI-owned instance that is not being deleted the controller MUST NOT change the finalizers.

The one exception runs before this gate: a deleting instance recognised as the controller's own that carries the finalizer is released as the requirement "The controller releases its finalizer from a deleting own instance without pruning" says, whatever its `spec.owner`.

Rationale: only the controller removes its finalizer, and the owner-skip gate is the only code that runs for a CLI-owned instance. Without the release the delete never completes.

#### Scenario: CLI-owned instance is not reconciled and gets no finalizer

- **GIVEN** a `ModuleInstance` with `spec.owner == cli` and no `DeletionTimestamp`
- **WHEN** the controller reconciles it
- **THEN** no resources are rendered or applied
- **AND** the `opmodel.dev/cleanup` finalizer is NOT present in `metadata.finalizers`
- **AND** the controller returns without requeueing for work

#### Scenario: Deleting a CLI-owned instance is a no-op for the controller

- **GIVEN** a `ModuleInstance` with `spec.owner == cli`, no `opmodel.dev/cleanup` finalizer and a non-zero `DeletionTimestamp`
- **WHEN** the controller reconciles it
- **THEN** the controller prunes no resources
- **AND** the controller does not block deletion

#### Scenario: Deleting a CLI-owned instance that carries the finalizer releases it

- **GIVEN** a `ModuleInstance` with `spec.owner == cli`, `spec.prune: true`, a non-empty `status.inventory`, the `opmodel.dev/cleanup` finalizer and a non-zero `DeletionTimestamp`
- **WHEN** the controller reconciles it
- **THEN** the finalizer is removed and the apiserver deletes the instance
- **AND** every object in the inventory still exists
- **AND** the controller returns without requeueing

#### Scenario: A live CLI-owned instance keeps a leftover finalizer

- **GIVEN** a `ModuleInstance` with `spec.owner == cli`, the `opmodel.dev/cleanup` finalizer and no `DeletionTimestamp`
- **WHEN** the controller reconciles it
- **THEN** the finalizer stays and the instance is acknowledged with `ManagedExternally`

#### Scenario: Controller-managed instances are unaffected

- **GIVEN** a `ModuleInstance` with `spec.owner` absent, empty, or `controller` that is not the controller's own instance
- **WHEN** the controller reconciles it
- **THEN** the controller proceeds with the normal reconcile (finalizer registration, render, apply, prune, status) exactly as before this change

### Requirement: The controller acknowledges a CLI-owned instance with a single condition

On a non-deleting `ModuleInstance` with `spec.owner == cli`, the controller MUST set `Ready: Unknown` with reason `ManagedExternally` and MUST NOT write any other status field — specifically it MUST NOT set `status.observedGeneration`, and MUST NOT modify `status.inventory`, the `lastApplied*` digests, `instanceUUID`, or any field the CLI writes. The `internal/status` package MUST expose a `ManagedExternally` reason constant and a helper that sets this condition (clearing `Reconciling`, `Stalled` and `Healthy`; the controller judges no health on an instance it does not reconcile, so a `Healthy` condition left from a controller-owned past would describe objects it no longer follows). The acknowledgement MUST be idempotent: reconciling an already-acknowledged instance MUST produce no status change.

#### Scenario: ManagedExternally condition is set

- **GIVEN** a `ModuleInstance` with `spec.owner == cli`
- **WHEN** the controller reconciles it
- **THEN** the `Ready` condition is `Unknown` with reason `ManagedExternally`
- **AND** the `Reconciling`, `Stalled` and `Healthy` conditions are absent

#### Scenario: No observedGeneration and no CLI-written status is touched

- **GIVEN** a `ModuleInstance` with `spec.owner == cli` carrying CLI-written `status.inventory` and `lastApplied*` digests
- **WHEN** the controller reconciles it
- **THEN** `status.observedGeneration` is not set by the controller
- **AND** `status.inventory`, the `lastApplied*` digests, and `instanceUUID` are unchanged

#### Scenario: Re-acknowledgement is a no-op

- **GIVEN** a `ModuleInstance` with `spec.owner == cli` already carrying `Ready: Unknown / ManagedExternally`
- **WHEN** the controller reconciles it again (e.g. after a requeue or an informer resync)
- **THEN** the resulting status patch is empty and no condition transition timestamp changes

### Requirement: An ownership handoff falls through to a normal reconcile

When a `ModuleInstance`'s `spec.owner` changes from `cli` to `controller`, the next reconcile MUST proceed through the normal path: register the finalizer, render, apply, prune, and write the real `Ready` status, overwriting the `ManagedExternally` condition. The controller's own instance is the exception: flipping its owner to `controller` MUST lead to the refusal of the requirement that the controller never reconciles its own instance, never to an adoption.

#### Scenario: Flip to controller adopts the instance

- **GIVEN** a `ModuleInstance` that is not the controller's own instance, previously reconciled with `spec.owner == cli` (carrying `ManagedExternally`, no finalizer)
- **WHEN** `spec.owner` is set to `controller` and the controller reconciles it
- **THEN** the `opmodel.dev/cleanup` finalizer is added
- **AND** the instance is rendered and applied
- **AND** the `Ready` condition reflects the real reconcile outcome (no longer `ManagedExternally`)

#### Scenario: Flip to controller on the controller's own instance is refused

- **GIVEN** the controller's own instance, previously reconciled with `spec.owner == cli` (carrying `ManagedExternally`, no finalizer)
- **WHEN** `spec.owner` is set to `controller` and the controller reconciles it
- **THEN** no finalizer is added and nothing is rendered, applied or pruned
- **AND** `Ready` is `False` with reason `SelfManagementRefused`

### Requirement: The controller recognises the instance that deploys it

The controller SHALL treat a `ModuleInstance` as the controller's own instance when any one of these holds, each read from the object as stored, without rendering the module:

- its `metadata.name` is `opm-controller` and its `metadata.namespace` is `opm-system`;
- its `spec.module.path`, with any `@<major>` suffix removed, is `opmodel.dev/modules/opm_controller`;
- its `status.inventory.entries` holds an entry of group `apiextensions.k8s.io` and kind `CustomResourceDefinition` whose name, after its first `.`, is the controller's own API group `opmodel.dev`.

The coordinates are the ones every CLI install records the controller's instance under; the module path is the one the controller module is published under, in any major. A refusal must be decided before the finalizer is registered, which is why every signal reads stored fields only.

#### Scenario: Fixed coordinates identify the instance

- **GIVEN** a `ModuleInstance` named `opm-controller` in namespace `opm-system` whose module path is any value
- **WHEN** the controller reconciles it
- **THEN** the controller treats it as the controller's own instance

#### Scenario: The controller module identifies the instance in any namespace and any major

- **GIVEN** a `ModuleInstance` named `team-ops` in namespace `platform` whose `spec.module.path` is `opmodel.dev/modules/opm_controller@v0`, `opmodel.dev/modules/opm_controller@v1` or `opmodel.dev/modules/opm_controller`
- **WHEN** the controller reconciles it
- **THEN** the controller treats it as the controller's own instance

#### Scenario: A recorded controller CRD identifies the instance

- **GIVEN** a `ModuleInstance` with another name, namespace and module path whose `status.inventory.entries` holds `CustomResourceDefinition` `moduleinstances.opmodel.dev`
- **WHEN** the controller reconciles it
- **THEN** the controller treats it as the controller's own instance

#### Scenario: Similar names are not the controller's own instance

- **GIVEN** a `ModuleInstance` named `opm-controller` in namespace `default`, whose module path is `opmodel.dev/modules/opm_controller_dashboard@v0`, and whose inventory holds only the `CustomResourceDefinition`s `widgets.example.opmodel.dev` and `widgets.example.opmodel.dev.io`
- **WHEN** the controller reconciles it
- **THEN** the controller does not treat it as the controller's own instance and reconciles it as its `spec.owner` says

### Requirement: The controller never reconciles its own instance

On the controller's own instance, whatever `spec.owner` says, the controller MUST never apply and never prune, and MUST NOT add the `opmodel.dev/cleanup` finalizer. An instance recognised as the controller's own only after a normal reconcile already carries that finalizer; the refusal leaves it in place while the instance is live, and the requirement "The controller releases its finalizer from a deleting own instance without pruning" removes it on delete. A live own instance with `spec.owner: cli` is handled by the owner-skip gate and acknowledged with `ManagedExternally`; an own instance with `spec.owner` absent or `controller` MUST be refused before finalizer registration, before the suspend check, and before any render. On refusal the controller MUST:

- set `Ready=False` and `Stalled=True`, both with reason `SelfManagementRefused`, and remove `Reconciling`, `ModuleResolved`, `Drifted` and `Healthy`, with a message that names the signal that matched and says to set `spec.owner` to `cli`;
- set `status.observedGeneration` to the instance's generation, so a client waiting on that generation reads a final verdict instead of timing out;
- emit a Warning event with reason `SelfManagementRefused` when the instance was not already refused;
- leave `status.inventory`, the `lastApplied*` digests, `instanceUUID`, `lastAttempted*`, history and failure counters unchanged;
- not requeue: the refusal is re-evaluated when the instance changes.

The `internal/status` package MUST expose a `SelfManagementRefused` reason constant. Re-refusing an already-refused instance at the same generation MUST produce an empty status patch.

Rationale: the controller is the one workload whose failure stops every other instance, so re-running install must be able to repair or replace it without any step taken by the running controller. Nothing the controller does to its own instance may stand between a re-run install and the cluster.

#### Scenario: An own instance flipped to controller is refused without a finalizer

- **GIVEN** the controller's own instance with `spec.owner: controller`, `spec.prune: true` and no `opmodel.dev/cleanup` finalizer
- **WHEN** the controller reconciles it
- **THEN** the finalizer is not added
- **AND** nothing is rendered, applied or pruned
- **AND** `Ready` is `False` and `Stalled` is `True`, both with reason `SelfManagementRefused`
- **AND** `status.observedGeneration` equals `metadata.generation`

#### Scenario: A suspended own instance is refused, not suspended

- **GIVEN** the controller's own instance with `spec.owner: controller` and `spec.suspend: true`
- **WHEN** the controller reconciles it
- **THEN** `Ready` is `False` with reason `SelfManagementRefused`, not `Suspended`
- **AND** the `opmodel.dev/cleanup` finalizer is not added

#### Scenario: Refusal clears conditions left by an earlier adoption

- **GIVEN** the controller's own instance with `spec.owner: controller` carrying `ModuleResolved=True`, `Drifted=True` and `Healthy=True` from an earlier reconcile
- **WHEN** the controller refuses it
- **THEN** `ModuleResolved`, `Drifted` and `Healthy` are removed

#### Scenario: An own instance with no owner is refused

- **GIVEN** the controller's own instance with no `spec.owner`
- **WHEN** the controller reconciles it
- **THEN** the controller refuses it exactly as an own instance with `spec.owner: controller`

#### Scenario: Refusal leaves CLI-written status alone

- **GIVEN** the controller's own instance flipped to `spec.owner: controller`, carrying a CLI-written `status.inventory`, `lastApplied*` digests and `instanceUUID`
- **WHEN** the controller refuses it
- **THEN** those fields are unchanged

#### Scenario: Re-refusal is a no-op

- **GIVEN** the controller's own instance already refused at its current generation
- **WHEN** the controller reconciles it again
- **THEN** the status patch is empty, no condition transition time changes, and no event is emitted

#### Scenario: Handing the instance back to the CLI clears the refusal

- **GIVEN** a refused own instance
- **WHEN** its `spec.owner` is set to `cli`
- **THEN** the next reconcile acknowledges it with `Ready=Unknown` reason `ManagedExternally` and removes `Stalled`

### Requirement: The controller releases its finalizer from a deleting own instance without pruning

When the controller's own instance, recognised by any of the signals of the requirement "The controller recognises the instance that deploys it", carries the `opmodel.dev/cleanup` finalizer and is being deleted, whatever its `spec.owner`, the controller MUST remove that finalizer and MUST NOT delete any object in the instance's inventory, whatever its `spec.prune`. On that deleting instance the controller MUST write no status and MUST emit a Warning event with reason `SelfManagementRefused` that says the finalizer was released without pruning. This release MUST run before the owner-skip gate. An own instance that is not being deleted keeps the finalizer and is refused as the requirement "The controller never reconciles its own instance" says.

Rationale: an instance can become recognisable as the controller's own only after a normal reconcile has registered the finalizer: a renamed or forked copy of the controller module once its inventory records a controller CRD, or an instance whose `spec.module.path` is later edited to the controller module. The controller refuses to run that instance's deletion cleanup, so without this release the instance would stay terminating until someone removed the finalizer by hand. The release ignores `spec.owner` because the refusal tells the user to set it to `cli`: if the owner-skip gate ran first, deleting the instance after following that advice would hang.

#### Scenario: Deleting an own instance recognised by its inventory prunes nothing

- **GIVEN** a `ModuleInstance` with `spec.owner: controller`, `spec.prune: true`, a module path other than the controller module's, the `opmodel.dev/cleanup` finalizer, and a `status.inventory` holding `CustomResourceDefinition` `moduleinstances.opmodel.dev` and an OPM-managed `ConfigMap`
- **WHEN** the instance is deleted and the controller reconciles it
- **THEN** the finalizer is removed and the instance goes away
- **AND** the `ConfigMap` still exists
- **AND** a Warning event with reason `SelfManagementRefused` is emitted

#### Scenario: An instance whose module path was edited to the controller module is released on delete

- **GIVEN** a `ModuleInstance` with `spec.owner: controller` and `spec.prune: true` that a normal reconcile gave the `opmodel.dev/cleanup` finalizer and an inventory, whose `spec.module.path` is then edited to `opmodel.dev/modules/opm_controller@v0`
- **WHEN** the instance is deleted and the controller reconciles it
- **THEN** the finalizer is removed and the instance goes away
- **AND** every object in the inventory still exists

#### Scenario: An own instance handed to the CLI after a refusal is released on delete

- **GIVEN** a refused own instance recognised by its inventory that carries the `opmodel.dev/cleanup` finalizer, whose `spec.owner` is then set to `cli`
- **WHEN** the instance is deleted and the controller reconciles it
- **THEN** the finalizer is removed and every object in the inventory still exists
- **AND** the instance's status is unchanged by that reconcile

#### Scenario: A live own instance keeps its finalizer and is refused

- **GIVEN** the controller's own instance with `spec.owner: controller` and the `opmodel.dev/cleanup` finalizer, not being deleted
- **WHEN** the controller reconciles it
- **THEN** the finalizer stays
- **AND** `Ready` is `False` with reason `SelfManagementRefused`
