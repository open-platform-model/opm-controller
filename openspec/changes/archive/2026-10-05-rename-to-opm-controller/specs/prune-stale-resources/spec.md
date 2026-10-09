## REMOVED Requirements

### Requirement: The opt-out is an optional field that changes no other field

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-021); this requirement is restated under the controller's name.

**Migration**: See the requirement "The opt-out field is optional and changes no other field" in this capability.

## MODIFIED Requirements

### Requirement: The prune judges with the recorded identities
The prune of stale resources MUST judge with the identities the status held when the apply of the same reconcile started: `status.instanceUUID`, and `status.previousInstanceUUID` when it is set. A stale object that carries either identity, and that the verdict lets the controller delete under it, MUST be deleted. Source: owner decisions of 2026-10-08 (prune judges with the identity stored in the instance's record, also after the instance identity changed; the status keeps both identities until the prune succeeded).

When `status.instanceUUID` was empty when the reconcile started, the prune MUST judge with the render's identity, so that a stale object that carries another instance's UUID label is left in the cluster, as before the field existed. A ModulePackage, and only a ModulePackage, MUST then ask the verdict once more with no identity: before the field existed its prune compared no UUID label, and it MUST NOT lose a delete it had.

#### Scenario: Stale object after the instance identity changed
- **GIVEN** a ModuleInstance with `status.instanceUUID` `A` and `spec.prune=true`, whose `spec.module.path` changes so that its render carries identity `B` and no longer holds ConfigMap `team-a/old`, which carries the UUID label `A`
- **WHEN** the reconcile applies and prunes
- **THEN** ConfigMap `team-a/old` is deleted
- **AND** after the reconcile `status.instanceUUID` is `B` and `status.previousInstanceUUID` is empty

#### Scenario: A failed prune leaves nothing orphaned on the retry
- **GIVEN** the same change of identity, a second stale ConfigMap `team-a/older` with the UUID label `A`, and a prune whose DELETE of `team-a/old` fails
- **WHEN** the reconcile ends and the next reconcile runs
- **THEN** after the first reconcile `status.instanceUUID` is `B`, `status.previousInstanceUUID` is `A` and `status.inventory` is unchanged
- **AND** the second reconcile deletes `team-a/old`, and `team-a/older` if it still exists
- **AND** after it `status.previousInstanceUUID` is empty

#### Scenario: A stale object that was already relabelled
- **GIVEN** an identity change from `A` to `B` that is not settled, and a later render of identity `B` that drops Deployment `team-a/app`, which the earlier apply relabelled to `B`
- **WHEN** the reconcile applies and prunes
- **THEN** Deployment `team-a/app` is deleted

#### Scenario: Stale object that carries a third identity
- **GIVEN** a stale ConfigMap whose live UUID label is neither of the recorded identities nor empty
- **WHEN** the controller prunes the stale set
- **THEN** the ConfigMap still exists and the prune result names it with the reason `owner-mismatch`

#### Scenario: No recorded identity, ModuleInstance, object of another identity
- **GIVEN** a ModuleInstance with an inventory and an empty `status.instanceUUID`, whose render carries identity `B`, and a stale ConfigMap that is managed by OPM, carries the UUID label `X` and has no adopt annotation
- **WHEN** the reconcile applies and prunes
- **THEN** the ConfigMap still exists and the prune result names it with the reason `owner-mismatch`

#### Scenario: No recorded identity, ModulePackage, object of an unknown earlier identity
- **GIVEN** a ModulePackage with an inventory and an empty `status.instanceUUID`, whose render carries identity `B`, and a stale ConfigMap that is managed by OPM, carries the UUID label `X` and has no adopt annotation
- **WHEN** the reconcile applies and prunes
- **THEN** the ConfigMap is deleted

#### Scenario: No recorded identity, object annotated for this instance
- **GIVEN** the same object, and a stale ConfigMap that is managed by OPM, carries no UUID label and carries the annotation `opmodel.dev/adopt` with the value `B`
- **WHEN** the reconcile applies and prunes
- **THEN** the ConfigMap is deleted

#### Scenario: No recorded identity, object annotated for another instance
- **GIVEN** the same object, and a stale ConfigMap that carries the annotation `opmodel.dev/adopt` with the value `C`
- **WHEN** the reconcile applies and prunes
- **THEN** the ConfigMap still exists and the prune result names it with the reason `adopted-elsewhere`

## ADDED Requirements

### Requirement: The opt-out field is optional and changes no other field
`ModuleInstance` and `ModulePackage` MUST accept an optional field `spec.dataPolicy` with exactly two values, `Keep` and `Delete`. An absent value MUST mean `Keep`. The API server MUST refuse any other value. The CRD MUST NOT set a default, and an object stored before the field existed MUST be valid unchanged.

For pruning and deletion, `spec.dataPolicy` MUST have no effect unless `spec.prune` is true, and the pair `dataPolicy: Delete` without `spec.prune` MUST be admitted, with a field description that says so: with `spec.prune` false or absent the controller prunes nothing, as before. `spec.prune` MUST keep its meaning for every kind other than PersistentVolumeClaim.

`spec.dataPolicy` MUST also govern the forced recreate of `spec.rollout.forceConflicts`, whatever `spec.prune` says: with `Keep` or no value the controller does not delete and recreate a PersistentVolumeClaim, and with `Delete` it does.

The field's description, which `kubectl explain` and the resource reference show, MUST say that claims a StatefulSet creates from its `volumeClaimTemplates` are never tracked and never deleted by the controller. It MUST also say that with `Keep` the forced recreate of `spec.rollout.forceConflicts` leaves a claim in place and the apply reports a conflict, and that with `Delete` the claim is deleted and recreated. It MUST NOT name the forced recreate as an exception to the protection.

#### Scenario: An object without the field is admitted and protected
- **GIVEN** a ModuleInstance manifest with `spec.prune: true` and no `spec.dataPolicy`
- **WHEN** it is applied to the API server
- **THEN** it is admitted, and the stored object has no `spec.dataPolicy`
- **AND** the controller keeps its PersistentVolumeClaims

#### Scenario: A value outside the enum is refused
- **GIVEN** a ModuleInstance manifest with `spec.dataPolicy: Purge`
- **WHEN** it is applied to the API server
- **THEN** the API server refuses it

#### Scenario: An explicit Keep protects
- **GIVEN** a ModuleInstance with `spec.prune: true` and `spec.dataPolicy: Keep`, whose render drops a PersistentVolumeClaim
- **WHEN** the controller reconciles
- **THEN** the claim is kept

#### Scenario: The field without prune deletes nothing
- **GIVEN** a ModuleInstance with `spec.dataPolicy: Delete` and no `spec.prune`, whose render drops a PersistentVolumeClaim
- **WHEN** the controller reconciles
- **THEN** no resource is deleted

#### Scenario: An existing instance is protected after a controller upgrade
- **GIVEN** a ModuleInstance with `spec.prune: true` that an earlier controller release reconciled, with a PersistentVolumeClaim in `status.inventory`
- **WHEN** the upgraded controller reconciles a render that no longer holds the claim
- **THEN** the claim is kept

#### Scenario: The forced recreate follows the field without prune
- **GIVEN** a ModuleInstance with `spec.rollout.forceConflicts: true`, no `spec.prune` and no `spec.dataPolicy`, whose render changes an immutable field of a live PersistentVolumeClaim
- **WHEN** the controller reconciles
- **THEN** the claim is kept with its UID
