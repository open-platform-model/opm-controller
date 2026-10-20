## REMOVED Requirements

### Requirement: Prune not attempted while stalled on DeletionSAMissing
**Reason**: Split into two requirements that each state one behavior.
**Migration**: See "The deletion-cleanup prune pass does not run while the release is stalled on DeletionSAMissing" and "Prune executes only with an enabled prune, entries to remove and a valid client".

### Requirement: PersistentVolumeClaims are kept unless the instance opts out
**Reason**: Split into five requirements that each state one behavior.
**Migration**: See "The prune keeps a PersistentVolumeClaim unless the instance opts out", "The prune result names every claim it kept", "A claim skipped by the ownership guard is not reported as kept", "A claim that cannot be read while it is protected is kept and reported" and "Only the kind PersistentVolumeClaim in the core group is kept".

### Requirement: The opt-out is an optional field that changes no other field
**Reason**: Split into five requirements that each state one behavior.
**Migration**: See "ModuleInstance and ModulePackage accept an optional spec.dataPolicy of Keep or Delete", "spec.dataPolicy has no effect without spec.prune", "spec.dataPolicy governs the forced recreate of spec.rollout.forceConflicts", "The description of spec.dataPolicy names the claims that are never tracked" and "The description of spec.dataPolicy states the forced recreate".

### Requirement: Prune asks the library's delete verdict
**Reason**: Split into seven requirements that each state one behavior.
**Migration**: See "The prune reads each object and asks the library's delete verdict", "The prune deletes an object only on a proceed verdict and decides ownership by no comparison of its own", "A proceed verdict deletes the object unless spec.dataPolicy keeps the claim", "An object that does not exist is a success", "An object that is not the instance's own is left in the cluster and named", "A failed read fails the entry and not the run" and "An object without a UUID label is deleted when OPM manages it".

### Requirement: The prune judges with the recorded identities
**Reason**: Split into four requirements that each state one behavior.
**Migration**: See "The prune judges with the identities the status held when the apply started", "A stale object that carries a recorded identity is deleted when the verdict allows", "The prune judges with the render's identity when no identity was recorded" and "A ModulePackage with no recorded identity asks the verdict once more with no identity".

### Requirement: The prune runs the library's deletion plan
**Reason**: Split into five requirements that each state one behavior.
**Migration**: See "The prune deletes a stale object only when the library's deletion transition names the delete", "The prune sends deletes in the plan order with Foreground propagation and the UID precondition", "The prune runs one plan per identity", "A PersistentVolumeClaim that spec.dataPolicy keeps is not an entry of any plan" and "The prune does not wait for a deleted object".

## ADDED Requirements

### Requirement: The deletion-cleanup prune pass does not run while the release is stalled on DeletionSAMissing
The deletion-cleanup prune pass MUST NOT execute while the release is stalled with reason `DeletionSAMissing`. In that state, the impersonated client cannot be built, and prune with any fallback identity is explicitly disallowed.

#### Scenario: Stalled release does not prune
- **GIVEN** a ModuleRelease stalled with reason `DeletionSAMissing`
- **WHEN** reconcile loops fire during the stall window
- **THEN** no delete API calls are made against any resource in `status.inventory`
- **AND** the inventory remains unchanged across reconciles until recovery (SA restore or orphan-exit)

### Requirement: Prune executes only with an enabled prune, entries to remove and a valid client
As in requirement "The deletion-cleanup prune pass does not run while the release is stalled on DeletionSAMissing": This requirement complements the existing prune-stale-resources contract: prune MUST execute only when (a) `spec.prune=true`, (b) the inventory has entries to remove, and (c) a valid apply/prune client has been obtained. Condition (c) now explicitly excludes the controller's own client as a fallback on the deletion path.

#### Scenario: The controller's own client is not a fallback on the deletion path

- **WHEN** a valid apply/prune client has not been obtained on the deletion path
- **THEN** no prune executes and the controller's own client is not used as a fallback

### Requirement: The prune keeps a PersistentVolumeClaim unless the instance opts out
The prune MUST NOT delete a `PersistentVolumeClaim` of the core API group unless the object being reconciled sets `spec.dataPolicy` to `Delete`. This holds for the prune of stale resources and for the deletion cleanup, which use the same prune. A kept claim MUST be left in the cluster unchanged, MUST NOT be reported as an error, and MUST NOT count as a prune failure.

#### Scenario: Stale claim kept by default
- **GIVEN** a ModuleInstance with `spec.prune=true` and no `spec.dataPolicy`, whose stale set holds PersistentVolumeClaim `media/old-cache` and ConfigMap `media/old-config`, both live and owned by the instance
- **WHEN** the controller prunes the stale set
- **THEN** the ConfigMap is deleted
- **AND** the PersistentVolumeClaim still exists in the cluster
- **AND** the prune result names `media/old-cache` as kept and reports no error

#### Scenario: Stale claim deleted when the instance opts out
- **GIVEN** the same stale set on a ModuleInstance with `spec.prune=true` and `spec.dataPolicy=Delete`
- **WHEN** the controller prunes the stale set
- **THEN** the PersistentVolumeClaim and the ConfigMap are both deleted
- **AND** the prune result names no kept claim

### Requirement: The prune result names every claim it kept
The prune result MUST name every claim it kept, so that the caller can report them. A claim counts as kept only when it exists in the cluster and passes the ownership guard: A claim that is not found in the cluster is treated as already gone, as for every stale resource, and is not reported as kept.
(Requirement: A claim skipped by the ownership guard is not reported as kept) (Requirement: A claim that cannot be read while it is protected is kept and reported)

#### Scenario: A claim that is gone is not reported as kept
- **GIVEN** a stale entry for PersistentVolumeClaim `media/old-cache` that does not exist in the cluster
- **WHEN** the controller prunes the stale set with no `spec.dataPolicy`
- **THEN** the prune succeeds and names no kept claim

### Requirement: A claim skipped by the ownership guard is not reported as kept
The prune result MUST name every claim it kept, so that the caller can report them. A claim counts as kept only when it exists in the cluster and passes the ownership guard: A claim whose live object is not OPM-managed, or carries the UUID of another instance, is skipped by the ownership guard and is not reported as kept.
(Requirement: The prune result names every claim it kept) (Requirement: A claim that cannot be read while it is protected is kept and reported)

#### Scenario: A claim of another owner is skipped, not kept
- **GIVEN** a stale entry for PersistentVolumeClaim `media/shared` whose live object carries the UUID label of another ModuleInstance
- **WHEN** the controller prunes the stale set with no `spec.dataPolicy`
- **THEN** the claim is not deleted and is counted as skipped
- **AND** the prune result does not name it as kept

### Requirement: A claim that cannot be read while it is protected is kept and reported
The prune result MUST name every claim it kept, so that the caller can report them. A claim counts as kept only when it exists in the cluster and passes the ownership guard: A claim that cannot be read while it is protected is kept and reported, and the failed read is not an error of the prune.
(Requirement: The prune result names every claim it kept) (Requirement: A claim skipped by the ownership guard is not reported as kept)

#### Scenario: An unreadable claim does not fail the prune
- **GIVEN** a stale entry for PersistentVolumeClaim `media/old-cache` and an API server that refuses the read of that claim
- **WHEN** the controller prunes the stale set with no `spec.dataPolicy`
- **THEN** the prune returns no error for that entry and names the claim as kept

### Requirement: Only the kind PersistentVolumeClaim in the core group is kept
Only the kind `PersistentVolumeClaim` in the core group MUST be kept. A PersistentVolume, a VolumeSnapshot and a claim-like kind of another API group are pruned like any other resource. Namespaces and CustomResourceDefinitions stay excluded, and `spec.dataPolicy` does not change that.

#### Scenario: Other storage kinds are pruned as before
- **GIVEN** a stale entry of a kind other than core `PersistentVolumeClaim`, live and owned by the instance
- **WHEN** the controller prunes the stale set with no `spec.dataPolicy`
- **THEN** the resource is deleted

### Requirement: ModuleInstance and ModulePackage accept an optional spec.dataPolicy of Keep or Delete
`ModuleInstance` and `ModulePackage` MUST accept an optional field `spec.dataPolicy` with exactly two values, `Keep` and `Delete`. An absent value MUST mean `Keep`. The API server MUST refuse any other value. The CRD MUST NOT set a default, and an object stored before the field existed MUST be valid unchanged.

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

#### Scenario: An existing instance is protected after an operator upgrade
- **GIVEN** a ModuleInstance with `spec.prune: true` that an earlier operator release reconciled, with a PersistentVolumeClaim in `status.inventory`
- **WHEN** the upgraded operator reconciles a render that no longer holds the claim
- **THEN** the claim is kept

### Requirement: spec.dataPolicy has no effect without spec.prune
For pruning and deletion, `spec.dataPolicy` MUST have no effect unless `spec.prune` is true, and the pair `dataPolicy: Delete` without `spec.prune` MUST be admitted, with a field description that says so: with `spec.prune` false or absent the controller prunes nothing, as before. `spec.prune` MUST keep its meaning for every kind other than PersistentVolumeClaim.

#### Scenario: The field without prune deletes nothing
- **GIVEN** a ModuleInstance with `spec.dataPolicy: Delete` and no `spec.prune`, whose render drops a PersistentVolumeClaim
- **WHEN** the controller reconciles
- **THEN** no resource is deleted

### Requirement: spec.dataPolicy governs the forced recreate of spec.rollout.forceConflicts
`spec.dataPolicy` MUST also govern the forced recreate of `spec.rollout.forceConflicts`, whatever `spec.prune` says: with `Keep` or no value the controller does not delete and recreate a PersistentVolumeClaim, and with `Delete` it does.

#### Scenario: The forced recreate follows the field without prune
- **GIVEN** a ModuleInstance with `spec.rollout.forceConflicts: true`, no `spec.prune` and no `spec.dataPolicy`, whose render changes an immutable field of a live PersistentVolumeClaim
- **WHEN** the controller reconciles
- **THEN** the claim is kept with its UID

### Requirement: The description of spec.dataPolicy names the claims that are never tracked
As in requirement "ModuleInstance and ModulePackage accept an optional spec.dataPolicy of Keep or Delete": The field's description, which `kubectl explain` and the resource reference show, MUST say that claims a StatefulSet creates from its `volumeClaimTemplates` are never tracked and never deleted by the operator.

#### Scenario: The description of spec.dataPolicy names claims that are never tracked

- **WHEN** the description of `spec.dataPolicy` is read in `kubectl explain` or the resource reference
- **THEN** it says that claims a StatefulSet creates from its `volumeClaimTemplates` are never tracked and never deleted by the operator

### Requirement: The description of spec.dataPolicy states the forced recreate
As in requirement "The description of spec.dataPolicy names the claims that are never tracked": It MUST also say that with `Keep` the forced recreate of `spec.rollout.forceConflicts` leaves a claim in place and the apply reports a conflict, and that with `Delete` the claim is deleted and recreated. It MUST NOT name the forced recreate as an exception to the protection.

#### Scenario: The description of spec.dataPolicy states the forced recreate

- **WHEN** the description of `spec.dataPolicy` is read in `kubectl explain` or the resource reference
- **THEN** it says that with `Keep` the forced recreate of `spec.rollout.forceConflicts` leaves a claim in place and the apply reports a conflict, that with `Delete` the claim is deleted and recreated, and it does not name the forced recreate as an exception to the protection

### Requirement: The prune reads each object and asks the library's delete verdict
For every entry it may delete, the prune MUST read the live object with the client that would delete it and MUST ask the library's delete verdict (`opm/k8s/ownership`) with that object and an identity of the instance. When the instance has more than one identity to judge with, it MUST ask with each in turn, and the object counts as the instance's own when the verdict says proceed for one of them.
**Source**: enhancement 0012 ("The library owns the whole deletion sequence but only the per-object verdict and the order for apply" and "The apply guard runs on every apply, refusing foreign objects outside the inventory and any object adopted by another instance; a per-object adopt annotation is its only override").

#### Scenario: The prune asks the library's delete verdict with the live object

- **WHEN** a stale entry is pruned
- **THEN** the prune reads the live object with the client that would delete it and asks the library's delete verdict (`opm/k8s/ownership`) with that object and an identity of the instance

### Requirement: The prune deletes an object only on a proceed verdict and decides ownership by no comparison of its own
As in requirement "The prune reads each object and asks the library's delete verdict": It MUST delete the object only then, and it MUST NOT decide ownership with a label or annotation comparison of its own.
**Source**: enhancement 0012 ("The library owns the whole deletion sequence but only the per-object verdict and the order for apply" and "The apply guard runs on every apply, refusing foreign objects outside the inventory and any object adopted by another instance; a per-object adopt annotation is its only override").

#### Scenario: The prune deletes an object only on a proceed verdict

- **WHEN** the verdict for an object is not proceed for any identity
- **THEN** the object is not deleted, and the prune decides ownership with no label or annotation comparison of its own

### Requirement: A proceed verdict deletes the object unless spec.dataPolicy keeps the claim
The prune MUST act on each answer as follows: Proceed: the object is deleted, unless it is a PersistentVolumeClaim that `spec.dataPolicy` keeps. A claim counts as kept only after the verdict said proceed.
(Requirement: An object that does not exist is a success)

#### Scenario: Own object is deleted
- **GIVEN** a stale entry for ConfigMap `team-a/example` whose live object is managed by OPM and carries this instance's UUID label and no adopt annotation
- **WHEN** the controller prunes the stale set
- **THEN** the ConfigMap is deleted and counted as deleted

### Requirement: An object that does not exist is a success
The prune MUST act on each answer as follows: The object does not exist: success, as before.
(Requirement: An object that is not the instance's own is left in the cluster and named)

#### Scenario: An object that does not exist is a success

- **WHEN** the object of a stale entry does not exist
- **THEN** the entry is a success

### Requirement: An object that is not the instance's own is left in the cluster and named
The prune MUST act on each answer as follows: For every identity, the object is not managed by OPM, belongs to another instance, or carries an adopt annotation that names another instance: the object is left in the cluster, counted as skipped and named in the prune result with the library's reason and message. It is not an error.
(Requirement: A failed read fails the entry and not the run)

#### Scenario: Object not managed by OPM is left
- **GIVEN** a stale entry for ConfigMap `team-a/example` whose live object has no `app.kubernetes.io/managed-by` label
- **WHEN** the controller prunes the stale set
- **THEN** the ConfigMap still exists
- **AND** the prune result names it with the reason `not-opm-managed`

#### Scenario: Object of another instance is left
- **GIVEN** a stale entry for ConfigMap `team-a/example` whose live object is managed by OPM and carries the UUID label of another instance
- **WHEN** the controller prunes the stale set with this instance's identity
- **THEN** the ConfigMap still exists
- **AND** the prune result names it with the reason `owner-mismatch`

#### Scenario: Object being adopted by another instance is left
- **GIVEN** a stale entry for ConfigMap `team-a/example` whose live object carries this instance's UUID label and the annotation `opmodel.dev/adopt` with the UUID of another instance
- **WHEN** the controller prunes the stale set
- **THEN** the ConfigMap still exists
- **AND** the prune result names it with the reason `adopted-elsewhere`

### Requirement: A failed read fails the entry and not the run
The prune MUST act on each answer as follows: The read fails with an error other than NotFound: the entry is a failed prune and the remaining entries are still attempted, as before. A PersistentVolumeClaim that cannot be read while `spec.dataPolicy` keeps claims is kept without an error, as before.
(Requirement: A proceed verdict deletes the object unless spec.dataPolicy keeps the claim)

#### Scenario: A failed read fails the entry and not the run
- **GIVEN** a stale set of two ConfigMaps and an API server that refuses the read of the first
- **WHEN** the controller prunes the stale set
- **THEN** the second ConfigMap is deleted
- **AND** the prune returns an error that names the first

### Requirement: An object without a UUID label is deleted when OPM manages it
An object without a UUID label MUST still be deleted when OPM manages it, and every OPM manager label value MUST be accepted, as before.

#### Scenario: Object without a UUID label is deleted
- **GIVEN** a stale entry for ConfigMap `team-a/legacy` whose live object carries `app.kubernetes.io/managed-by=open-platform-model` and no UUID label
- **WHEN** the controller prunes the stale set
- **THEN** the ConfigMap is deleted

#### Scenario: Object still carrying the cli manager identity is deleted
- **GIVEN** a stale entry for ConfigMap `team-a/example` whose live object carries `app.kubernetes.io/managed-by=opm-cli` and this instance's UUID label
- **WHEN** the controller prunes the stale set
- **THEN** the ConfigMap is deleted

### Requirement: The prune judges with the identities the status held when the apply started
The prune of stale resources MUST judge with the identities the status held when the apply of the same reconcile started: `status.instanceUUID`, and `status.previousInstanceUUID` when it is set.
**Source**: owner decisions of 2026-10-08 (prune judges with the identity stored in the instance's record, also after the instance identity changed; the status keeps both identities until the prune succeeded).

#### Scenario: Stale object that carries a third identity
- **GIVEN** a stale ConfigMap whose live UUID label is neither of the recorded identities nor empty
- **WHEN** the controller prunes the stale set
- **THEN** the ConfigMap still exists and the prune result names it with the reason `owner-mismatch`

### Requirement: A stale object that carries a recorded identity is deleted when the verdict allows
As in requirement "The prune judges with the identities the status held when the apply started": A stale object that carries either identity, and that the verdict lets the operator delete under it, MUST be deleted.
**Source**: owner decisions of 2026-10-08 (prune judges with the identity stored in the instance's record, also after the instance identity changed; the status keeps both identities until the prune succeeded).

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

### Requirement: The prune judges with the render's identity when no identity was recorded
When `status.instanceUUID` was empty when the reconcile started, the prune MUST judge with the render's identity, so that a stale object that carries another instance's UUID label is left in the cluster, as before the field existed.

#### Scenario: No recorded identity, ModuleInstance, object of another identity
- **GIVEN** a ModuleInstance with an inventory and an empty `status.instanceUUID`, whose render carries identity `B`, and a stale ConfigMap that is managed by OPM, carries the UUID label `X` and has no adopt annotation
- **WHEN** the reconcile applies and prunes
- **THEN** the ConfigMap still exists and the prune result names it with the reason `owner-mismatch`

### Requirement: A ModulePackage with no recorded identity asks the verdict once more with no identity
As in requirement "The prune judges with the render's identity when no identity was recorded": A ModulePackage, and only a ModulePackage, MUST then ask the verdict once more with no identity: before the field existed its prune compared no UUID label, and it MUST NOT lose a delete it had.

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

### Requirement: The prune deletes a stale object only when the library's deletion transition names the delete
The prune of stale resources MUST delete a stale object only when the library's deletion transition (`opm/k8s/lifecycle`) names that delete as the next action.
**Source**: enhancement 0012 ("The library owns the whole deletion sequence but only the per-object verdict and the order for apply").

#### Scenario: A stale object is deleted only as the next action of the plan

- **WHEN** a stale object is pruned
- **THEN** the delete is sent only when the library's deletion transition (`opm/k8s/lifecycle`) names that delete as the next action

### Requirement: The prune sends deletes in the plan order with Foreground propagation and the UID precondition
As in requirement "The prune deletes a stale object only when the library's deletion transition names the delete": The prune MUST send each delete in the order the plan gives, descending kind weight, and MUST add no ordering of its own. Each delete MUST carry the propagation policy the action names, Foreground, and the action's UID precondition.

#### Scenario: Stale objects are deleted in descending kind weight
- **GIVEN** a stale set that lists a ConfigMap before a Deployment
- **WHEN** the controller prunes the stale set
- **THEN** the DELETE of the Deployment is sent before the DELETE of the ConfigMap

#### Scenario: Stale deletes use Foreground propagation
- **GIVEN** a stale Deployment owned by the instance
- **WHEN** the controller prunes it
- **THEN** the DELETE request carries the propagation policy `Foreground`

### Requirement: The prune runs one plan per identity
As in requirement "The prune deletes a stale object only when the library's deletion transition names the delete": With more than one identity to judge with, the prune MUST run one plan per identity, in the order the identities are asked today: the entries a plan skips because the object belongs to, or is being adopted by, another instance form the plan of the next identity.

#### Scenario: One plan runs per identity

- **WHEN** the instance has more than one identity to judge with
- **THEN** the prune runs one plan per identity, in the order the identities are asked today, and the entries a plan skips because the object belongs to, or is being adopted by, another instance form the plan of the next identity

### Requirement: A PersistentVolumeClaim that spec.dataPolicy keeps is not an entry of any plan
As in requirement "The prune deletes a stale object only when the library's deletion transition names the delete": A PersistentVolumeClaim that `spec.dataPolicy` keeps MUST NOT be an entry of any plan; it is still read and judged, so that it is reported as kept, left behind or gone as before.

#### Scenario: A kept stale claim is not in the plan
- **GIVEN** a ModuleInstance with no `spec.dataPolicy` and a stale PersistentVolumeClaim of its own
- **WHEN** the controller prunes the stale set
- **THEN** no DELETE is sent for the claim and it is reported as kept

### Requirement: The prune does not wait for a deleted object
The prune MUST NOT wait for a deleted object to disappear. A stale entry whose delete the API server accepted leaves `status.inventory`, as before.

#### Scenario: The reconcile does not wait for a stale object
- **GIVEN** a stale Deployment whose delete is accepted and which still exists with a `deletionTimestamp`
- **WHEN** the reconcile ends
- **THEN** `Ready` is True and `status.inventory.entries` does not list the Deployment
