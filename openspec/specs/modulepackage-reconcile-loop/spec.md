## Purpose

Defines the end-to-end reconciliation loop for the Release CRD: phase ordering, triggers, suspend/no-op handling, finalizer cleanup, and status shape.

## Requirements

### Requirement: Full reconcile loop execution
The `ReleaseReconciler` MUST execute phases sequentially: source resolution → artifact fetch → path navigation → CUE load → kind detection → render → apply → prune → status update. A failure in CUE load or render MUST be classified by its type (see `reconcile-backoff`, "Registry fetch failures are transient wherever they occur"): a typed registry fetch failure retries on the backoff, any other failure stalls.

#### Scenario: First successful reconcile (ModuleRelease)
- **WHEN** a Release CR is created with a valid `sourceRef`, the Flux source is ready, `spec.path` contains a valid `release.cue` evaluating to `#ModuleRelease`
- **THEN** the controller resolves the source, fetches the artifact, navigates to path, loads CUE, detects kind, renders resources, applies via SSA, updates status with conditions/digests/inventory/history, and sets `Ready=True`

#### Scenario: Source not ready
- **WHEN** the referenced Flux source exists but is not ready
- **THEN** the controller sets `Ready=False` with reason `SourceNotReady` and requeues with interval

#### Scenario: Render failure
- **WHEN** the CUE package loads but its render fails evaluation for a cause that is neither a resolution-class failure nor a registry fetch failure
- **THEN** the controller sets `Ready=False`, `Stalled=True` with reason `RenderFailed`, and does NOT modify inventory or attempt apply

#### Scenario: Package load failure
- **WHEN** the CUE package fails to load for a cause that is not a registry fetch failure, for example a CUE syntax error
- **THEN** the controller sets `Ready=False`, `Stalled=True` with reason `ResolutionFailed`, requeues on the 30-minute recheck, and does NOT modify inventory or attempt apply

#### Scenario: Registry failure during load or render
- **WHEN** loading or rendering the CUE package fails with a typed registry fetch failure
- **THEN** the controller sets `Ready=False` with reason `ResolutionFailed`, does not set `Stalled=True`, requeues on the exponential backoff, and does NOT modify inventory or attempt apply

#### Scenario: Apply failure
- **WHEN** SSA apply fails
- **THEN** the controller sets `Ready=False` with reason `ApplyFailed`, does NOT prune, does NOT update `lastApplied*` digests, and requeues with backoff

### Requirement: Reconcile triggers
The `ReleaseReconciler` MUST reconcile on three triggers: CR spec changes, source artifact revision changes, and interval-based re-reconciliation. An interval requeue renders only when a render input key part changed, the resolved source differs from `status.source`, the package is not `Ready`, or `--drift-render-interval` has passed since `status.lastAppliedInputs.renderedAt` (`render-input-key`); otherwise it resolves the source and requeues without rendering.

#### Scenario: CR spec change triggers reconcile
- **WHEN** a Release CR's spec is modified (path, sourceRef, prune, etc.)
- **THEN** reconciliation is triggered immediately

#### Scenario: Source revision change triggers reconcile
- **WHEN** the referenced Flux source's `status.artifact` changes (new revision/digest)
- **THEN** all Release CRs referencing that source are enqueued for reconciliation

#### Scenario: Interval-based re-reconciliation
- **WHEN** the interval period elapses since the last successful reconcile
- **THEN** reconciliation is triggered, and it renders and re-applies if needed when an input changed or the drift render interval has passed since `status.lastAppliedInputs.renderedAt`

### Requirement: Suspend check
The `ReleaseReconciler` MUST skip reconciliation when `spec.suspend` is true.

#### Scenario: Suspended release
- **WHEN** `spec.suspend` is true
- **THEN** the controller sets condition reason `Suspended` and returns without requeue

#### Scenario: Resume from suspend
- **WHEN** `spec.suspend` changes from true to false
- **THEN** the controller emits a resume event and proceeds with normal reconciliation

### Requirement: No-op detection
The `ReleaseReconciler` MUST detect no-op reconciliations when source artifact revision, config, render, and inventory digests all match the last applied values.

#### Scenario: All digests match
- **WHEN** source artifact digest, config digest, render digest, and inventory digest all match the last applied values
- **AND** every rendered object that exists is in `status.inventory` and none is adopted by another instance
- **THEN** the controller skips apply and prune, keeps `Ready=True`, and requeues with interval

#### Scenario: All digests match and an object was adopted by another instance
- **GIVEN** a Ready ModulePackage and an inventoried ConfigMap that a user annotates `opmodel.dev/adopt` with another instance's UUID
- **WHEN** a reconcile renders with matching digests
- **THEN** the ConfigMap is not written, `status.inventory` no longer lists it and `Ready` stays `True`

#### Scenario: Matching digests, an object to take in and an inventoried object being deleted
- **GIVEN** a Ready ModulePackage with matching digests, a rendered ConfigMap that exists outside `status.inventory` with the package's adopt annotation, and an inventoried ConfigMap that carries a deletion timestamp
- **WHEN** a reconcile renders
- **THEN** nothing is written and the package reports `Ready=False` with reason `ApplyRefused`
- **AND** once the inventoried ConfigMap is gone, the retry applies the render, creates it again and takes the other in

### Requirement: Source digest from artifact metadata
The Release reconciler MUST derive the source digest from the Flux source artifact's revision and digest, not from CUE module path/version.

#### Scenario: Source digest computation
- **WHEN** the Flux source artifact has revision `main@sha1:abc123` and digest `sha256:def456`
- **THEN** the source digest is computed from the artifact revision and digest

### Requirement: Finalizer and deletion cleanup
The `ReleaseReconciler` MUST register a finalizer on Release CRs and clean up owned resources on deletion.

#### Scenario: Deletion with prune enabled
- **WHEN** a Release CR is deleted and `spec.prune` is true
- **THEN** the controller prunes all inventory entries, then removes the finalizer

#### Scenario: Deletion with prune disabled
- **WHEN** a Release CR is deleted and `spec.prune` is false
- **THEN** the controller removes the finalizer without pruning (orphans resources)

### Requirement: Source status tracking
The `ReleaseStatus` MUST include a `source` field reflecting the resolved Flux artifact metadata (ref, revision, digest, URL).

#### Scenario: Source metadata recorded
- **WHEN** a Flux source is resolved successfully
- **THEN** `status.source` reflects the source reference, artifact revision, artifact digest, and artifact URL

### Requirement: ModulePackage keeps PersistentVolumeClaims by default
The ModulePackage reconciler MUST apply the same protection as the ModuleInstance reconciler: with `spec.prune` true it MUST keep core `PersistentVolumeClaim` resources on the prune of stale resources and on deletion cleanup, unless `spec.dataPolicy` on the ModulePackage is `Delete`. It MUST emit the same `ClaimsKept` event, a kept stale claim MUST leave `status.inventory`, and a kept claim MUST NOT hold the finalizer.

#### Scenario: Stale claim of a ModulePackage kept
- **GIVEN** a ModulePackage with `spec.prune=true` and no `spec.dataPolicy`, whose render drops a PersistentVolumeClaim it applied before
- **WHEN** the controller reconciles
- **THEN** the claim still exists in the cluster and is not in `status.inventory`
- **AND** a `Normal` event with reason `ClaimsKept` is emitted

#### Scenario: Deletion of a ModulePackage keeps its claims
- **GIVEN** a ModulePackage with `spec.prune=true` and no `spec.dataPolicy` that is being deleted, with a PersistentVolumeClaim in its inventory
- **WHEN** the controller reconciles
- **THEN** the claim still exists, every other entry is deleted, and the finalizer is removed once those entries are gone

#### Scenario: A ModulePackage that opts out deletes its claims
- **GIVEN** a ModulePackage with `spec.prune=true` and `spec.dataPolicy=Delete` that is being deleted, with a PersistentVolumeClaim in its inventory
- **WHEN** the controller reconciles
- **THEN** the claim is deleted with the other entries

### Requirement: A ModulePackage keeps claims on a forced recreate
A ModulePackage reconcile MUST pass `spec.dataPolicy` to the apply as it does to the prune. With `spec.rollout.forceConflicts: true` and `spec.dataPolicy` `Keep` or absent, a PersistentVolumeClaim whose update the API server refuses MUST be kept, and the reconcile MUST report `Ready=False` with reason `ClaimConflict` and one `Warning` event with that reason, and retry on its backoff. With `spec.dataPolicy: Delete` the claim MUST be deleted and recreated.

#### Scenario: A refused claim on a ModulePackage
- **GIVEN** a ModulePackage with `spec.rollout.forceConflicts: true` and no `spec.dataPolicy`, whose render changes `storageClassName` of a live PersistentVolumeClaim
- **WHEN** the controller reconciles
- **THEN** the claim has the UID it had
- **AND** `Ready` is `False` with reason `ClaimConflict`, and a `Warning` event with reason `ClaimConflict` names the claim

#### Scenario: Delete recreates the claim of a ModulePackage
- **GIVEN** a ModulePackage with `spec.rollout.forceConflicts: true` and `spec.dataPolicy: Delete`, whose render changes `storageClassName` of a live PersistentVolumeClaim
- **WHEN** the controller reconciles
- **THEN** the claim is deleted and created again, and `Ready` is `True`

### Requirement: A ModulePackage deletes and waits as a ModuleInstance does
The deletion cleanup and the prune of stale resources of a ModulePackage MUST follow the same rules as a ModuleInstance's (capabilities `finalizer-and-deletion` and `prune-stale-resources`): the library's deletion plan, the hold verdict, the wait until deleted objects are gone, the reasons `DeletionInProgress` and `DeletionBlocked`, the release when the deleting identity is lost after every delete was sent, and the release of an inventory of kept claims only.

#### Scenario: A terminating object holds the finalizer of a ModulePackage
- **GIVEN** a ModulePackage with `spec.prune=true` being deleted, whose Deployment still exists with a `deletionTimestamp`
- **WHEN** the controller reconciles it
- **THEN** the finalizer is still present and `Ready` is False with reason `DeletionInProgress`

#### Scenario: A ModulePackage is released when its ServiceAccount goes during the wait
- **GIVEN** a ModulePackage being deleted with reason `DeletionInProgress`, whose ServiceAccount is deleted
- **WHEN** the controller reconciles it
- **THEN** the finalizer is removed and one `Warning` event with reason `DeletionUnconfirmed` is emitted

### Requirement: A reconcile whose digests match applies the objects the verdict allows when one is outside the inventory or another instance adopts an object
The `ReleaseReconciler` MUST detect no-op reconciliations when source artifact revision, config, render, and inventory digests all match the last applied values. A reconcile whose digests match is not a no-op when a rendered object that the apply verdict allows exists outside `status.inventory`, or when an inventoried object is adopted by another instance: it then applies the objects the verdict allows and records the inventory without the adopted object.

#### Scenario: A reconcile whose digests match is not a no-op
- **WHEN** a reconcile whose digests match finds a rendered object that the apply verdict allows outside `status.inventory`, or an inventoried object adopted by another instance
- **THEN** it applies the objects the verdict allows and records the inventory without the adopted object

### Requirement: A reconcile whose digests match applies the rendered set and an inventoried object being deleted refuses it
A reconcile whose digests match is not a no-op when a rendered object that the apply verdict allows exists outside `status.inventory`, or when an inventoried object is adopted by another instance. Such a reconcile SHALL apply the rendered set, so an inventoried object that is being deleted refuses it (`ssa-apply`, "An apply is judged by the ownership verdict before its first write"): nothing is written and `Ready` is `False` with reason `ApplyRefused` until that object is gone.

#### Scenario: Matching digests apply the rendered set
- **WHEN** a reconcile whose digests match finds a rendered object that the apply verdict allows outside `status.inventory`, or an inventoried object adopted by another instance
- **THEN** it applies the rendered set, so an inventoried object that is being deleted refuses it: nothing is written and `Ready` is `False` with reason `ApplyRefused` until that object is gone

### Requirement: The ReleaseReconciler patches the status at the end of every reconcile attempt
The `ReleaseReconciler` MUST patch `Release.status` at the end of every reconcile attempt, including NoOp. The status shape mirrors ModuleRelease: conditions, digests, inventory, history, failure counters, `nextRetryAt`, and `lastAppliedInputs`.

#### Scenario: Status updated on failure
- **WHEN** a phase fails
- **THEN** status conditions, `lastAttempted*` fields, `failureCounters`, and `nextRetryAt` are updated

### Requirement: A reconcile that skips its render is not an attempt and patches only a changed Healthy condition
A reconcile that skips its render because its inputs are unchanged (`render-input-key`) is not an attempt. It judges health (`instance-health`) and MUST NOT patch anything but the `Healthy` condition, which it patches only when the judgement changed it; it never writes the transient `Reconciling` condition the attempt set before the skip. It requeues after `spec.interval`, or sooner when health asks for it.

#### Scenario: A skipped interval patches nothing
- **WHEN** the interval requeue finds the package's inputs unchanged within the drift render interval, and the package is `Healthy=True` with reason `RolledOut` and still rolled out
- **THEN** no status patch is sent and the reconcile requeues after `spec.interval`

#### Scenario: A skipped interval records a changed health judgement
- **WHEN** the interval requeue skips the render and an inventory Deployment has lost its available replicas since the package was `RolledOut`
- **THEN** the only status change is `Healthy=False` with reason `NotRolledOut`, `Reconciling` is not written, and the reconcile requeues after the shorter of the health requeue and `spec.interval`

### Requirement: The patch of a successful outcome and of a NoOp carries the judged Healthy condition
The `ReleaseReconciler` MUST patch `Release.status` at the end of every reconcile attempt, including NoOp. On a successful outcome and on `NoOp`, the patch also carries the judged `Healthy` condition.

#### Scenario: Successful reconcile status
- **WHEN** all phases succeed
- **THEN** `Ready=True`, `lastApplied*` digests are set, inventory is replaced, and a success history entry is recorded

### Requirement: A ModulePackage carries the identity of the instance it renders and the identity before an unsettled change
A ModulePackage MUST carry the identity of the instance it renders in the optional status field `status.instanceUUID`, the UUID label value of its rendered objects, and the identity it had before an unsettled change in the optional field `status.previousInstanceUUID`. Source: enhancement 0012 ("The apply guard runs on every apply, refusing foreign objects outside the inventory and any object adopted by another instance; a per-object adopt annotation is its only override").

#### Scenario: The identity is recorded
- **GIVEN** a new ModulePackage that renders and applies with success
- **WHEN** the status is committed
- **THEN** `status.instanceUUID` holds the UUID label value of the rendered objects

### Requirement: The identity fields are written, kept and cleared as on a ModuleInstance
`status.instanceUUID` and `status.previousInstanceUUID` MUST be written, kept and cleared as on a ModuleInstance, and a second identity change before the first is settled MUST be refused as on a ModuleInstance. Source: enhancement 0012 ("The apply guard runs on every apply, refusing foreign objects outside the inventory and any object adopted by another instance; a per-object adopt annotation is its only override").

#### Scenario: A second identity change is refused
- **WHEN** a ModulePackage has an unsettled identity change and a second identity change comes
- **THEN** it is refused as on a ModuleInstance

### Requirement: A ModulePackage stored before the fields existed gains the identity
A ModulePackage stored before the fields existed MUST be valid unchanged and MUST gain `status.instanceUUID`, with no change to its spec, in its first reconcile that renders: before the apply when that reconcile applies, and in its status commit when it applies nothing. Source: enhancement 0012 ("The apply guard runs on every apply, refusing foreign objects outside the inventory and any object adopted by another instance; a per-object adopt annotation is its only override").

#### Scenario: An existing package gains the field
- **GIVEN** a ModulePackage with an inventory and no `status.instanceUUID`, reconciled by an earlier operator release
- **WHEN** the upgraded operator reconciles it with unchanged inputs
- **THEN** `status.instanceUUID` is set, `status.previousInstanceUUID` is empty and nothing is applied

### Requirement: The prune and the deletion cleanup of a ModulePackage judge with the recorded identities
The prune of stale resources and the deletion cleanup of a ModulePackage MUST judge with the recorded identities as a ModuleInstance's do. Source: enhancement 0012 ("The apply guard runs on every apply, refusing foreign objects outside the inventory and any object adopted by another instance; a per-object adopt annotation is its only override").

#### Scenario: A package does not delete another instance's object
- **GIVEN** a ModulePackage with `status.instanceUUID` `A` and `spec.prune=true`, whose stale set holds a ConfigMap that is managed by OPM and carries the UUID label `B`
- **WHEN** the controller prunes
- **THEN** the ConfigMap still exists

#### Scenario: Deletion of a package judges with the recorded identity
- **GIVEN** a ModulePackage with `status.instanceUUID` `A` being deleted, whose inventory holds a ConfigMap with the UUID label `A` and a Secret with the UUID label `B`
- **WHEN** the deletion cleanup runs
- **THEN** the ConfigMap is deleted, the Secret still exists and the finalizer is removed once the ConfigMap is gone

### Requirement: A ModulePackage with no recorded identity deletes what it deleted before the field existed
While no identity is recorded, a ModulePackage MUST delete what it deleted before the field existed, except an object that carries an adopt annotation naming another identity than the render's, and on deletion any object that carries an adopt annotation. Source: enhancement 0012 ("The apply guard runs on every apply, refusing foreign objects outside the inventory and any object adopted by another instance; a per-object adopt annotation is its only override").

#### Scenario: A package that never rendered is deleted as before
- **GIVEN** a suspended ModulePackage with an inventory and no `status.instanceUUID`
- **WHEN** it is deleted with `spec.prune=true`
- **THEN** every inventory object that is managed by OPM and has no adopt annotation is deleted, whatever its UUID label

#### Scenario: The first render after the upgrade also changes the identity
- **GIVEN** a ModulePackage with an inventory and no `status.instanceUUID`, whose source now renders the instance under another module path, and a stale ConfigMap that is managed by OPM and carries the UUID label of the earlier identity
- **WHEN** the upgraded operator reconciles it
- **THEN** the stale ConfigMap is deleted, as before the field existed

### Requirement: A ModulePackage reconcile judges every rendered object with the instance's identity
A ModulePackage reconcile that renders MUST judge every rendered object with the library's apply verdict before its first write, with the identity that capability `ssa-apply` defines (the instance's identity). Source: enhancement 0012 ("The apply guard runs on every apply, refusing foreign objects outside the inventory and any object adopted by another instance; a per-object adopt annotation is its only override").

#### Scenario: A package is refused on another instance's object
- **GIVEN** a ModulePackage whose render names a ConfigMap that a ModuleInstance holds and that is not in the package's inventory
- **WHEN** the controller reconciles the package
- **THEN** nothing is written, and the package reports `Ready=False` with reason `ApplyRefused` and no `Stalled` condition

#### Scenario: A package adopts an annotated object
- **GIVEN** the ConfigMap of the scenario above annotated `opmodel.dev/adopt` with the package's `status.instanceUUID`
- **WHEN** the retry runs
- **THEN** the render is applied and the ConfigMap is listed in the package's `status.inventory`

### Requirement: A ModulePackage reconcile never asks the verdict with the earlier identity of an unsettled change
A ModulePackage reconcile that renders MUST judge every rendered object with the library's apply verdict before its first write, with the identity that capability `ssa-apply` defines (never the earlier identity of an unsettled change). Source: enhancement 0012 ("The apply guard runs on every apply, refusing foreign objects outside the inventory and any object adopted by another instance; a per-object adopt annotation is its only override").

#### Scenario: The earlier identity is not asked with
- **WHEN** a ModulePackage reconcile that renders judges its rendered objects with the apply verdict
- **THEN** it does not ask with the earlier identity of an unsettled change

### Requirement: A ModulePackage reconcile never asks the verdict with an empty identity
A ModulePackage reconcile that renders MUST judge every rendered object with the library's apply verdict before its first write, with the identity that capability `ssa-apply` defines (never an empty identity, also when the package has none recorded). Source: enhancement 0012 ("The apply guard runs on every apply, refusing foreign objects outside the inventory and any object adopted by another instance; a per-object adopt annotation is its only override").

#### Scenario: No empty identity
- **WHEN** a ModulePackage reconcile that renders judges its rendered objects with the apply verdict and the package has no identity recorded
- **THEN** it does not ask with an empty identity

### Requirement: A ModulePackage reconcile judges exactly as a ModuleInstance reconcile does under the ssa-apply requirement
A ModulePackage reconcile that renders MUST judge every rendered object with the library's apply verdict before its first write, exactly as a ModuleInstance reconcile does (`ssa-apply`, "An apply is judged by the ownership verdict before its first write"). Source: enhancement 0012 ("The apply guard runs on every apply, refusing foreign objects outside the inventory and any object adopted by another instance; a per-object adopt annotation is its only override").

#### Scenario: Judged as a ModuleInstance reconcile does
- **WHEN** a ModulePackage reconcile renders
- **THEN** it judges every rendered object with the library's apply verdict before its first write, exactly as a ModuleInstance reconcile does

### Requirement: A ModulePackage reconcile judges exactly as a ModuleInstance reconcile does under the reconcile-loop-assembly requirement
A ModulePackage reconcile that renders MUST judge every rendered object with the library's apply verdict before its first write, exactly as a ModuleInstance reconcile does (`reconcile-loop-assembly`, "A reconcile the apply verdict refuses is retried and not stalled"). Source: enhancement 0012 ("The apply guard runs on every apply, refusing foreign objects outside the inventory and any object adopted by another instance; a per-object adopt annotation is its only override").

#### Scenario: A ModulePackage reconcile judges as a ModuleInstance reconcile does under the loop assembly
- **WHEN** a ModulePackage reconcile renders
- **THEN** it judges every rendered object with the library's apply verdict before its first write, exactly as a ModuleInstance reconcile does

### Requirement: The read is made by the client that applies and a skipped render reads nothing
A ModulePackage reconcile that renders MUST judge every rendered object with the library's apply verdict before its first write. The read MUST be made by the client that applies. A reconcile that skips its render MUST NOT read the objects. Source: enhancement 0012 ("The apply guard runs on every apply, refusing foreign objects outside the inventory and any object adopted by another instance; a per-object adopt annotation is its only override").

#### Scenario: A skipped render reads nothing
- **GIVEN** a ModulePackage whose render inputs are unchanged inside the drift render interval
- **WHEN** the controller reconciles
- **THEN** no rendered object is read for the verdict

### Requirement: A reconcile that renders fails when the verdict could not be asked
A ModulePackage has no drift check to report a verdict that could not be asked. So a reconcile that renders MUST fail, also when every digest matches, when the client that applies cannot be built or the read of a rendered object fails for a reason other than that the object does not exist. Source: enhancement 0012 ("The apply guard runs on every apply, refusing foreign objects outside the inventory and any object adopted by another instance; a per-object adopt annotation is its only override").

#### Scenario: Matching digests and a verdict that could not be asked
- **WHEN** a ModulePackage reconcile that renders cannot build the client that applies, or the read of a rendered object fails for a reason other than that the object does not exist, also when every digest matches
- **THEN** the reconcile fails

### Requirement: The failure reason is ImpersonationFailed when the ServiceAccount is missing or the read is Forbidden
So a reconcile that renders MUST fail, as listed in requirement "A reconcile that renders fails when the verdict could not be asked": `Stalled=True` with reason `ImpersonationFailed` when the ServiceAccount is missing or the read is Forbidden under an effective ServiceAccount. Source: enhancement 0012 ("The apply guard runs on every apply, refusing foreign objects outside the inventory and any object adopted by another instance; a per-object adopt annotation is its only override").

#### Scenario: Matching digests and an unreadable object
- **GIVEN** a Ready ModulePackage with matching digests whose effective ServiceAccount may no longer get ConfigMaps
- **WHEN** a reconcile renders
- **THEN** nothing is written, `status.inventory` keeps its entries, and the package reports `Stalled=True` with reason `ImpersonationFailed`

#### Scenario: Matching digests and a missing ServiceAccount
- **GIVEN** a Ready ModulePackage with matching digests whose effective ServiceAccount was deleted
- **WHEN** a reconcile renders
- **THEN** the package reports `Stalled=True` with reason `ImpersonationFailed` and is not a `NoOp`

### Requirement: The failure reason is ApplyFailed on the backoff except when the ServiceAccount is missing or the read is Forbidden
So a reconcile that renders MUST fail, as listed in requirement "A reconcile that renders fails when the verdict could not be asked": `Ready=False` with reason `ApplyFailed` on the backoff, except when the ServiceAccount is missing or the read is Forbidden under an effective ServiceAccount. Source: enhancement 0012 ("The apply guard runs on every apply, refusing foreign objects outside the inventory and any object adopted by another instance; a per-object adopt annotation is its only override").

#### Scenario: Another failure to ask the verdict is ApplyFailed
- **WHEN** a reconcile that renders must fail, the ServiceAccount is not missing and the read is not Forbidden under an effective ServiceAccount
- **THEN** the package reports `Ready=False` with reason `ApplyFailed` on the backoff

### Requirement: A failed ModulePackage reconcile writes nothing and keeps the inventory
So a reconcile that renders MUST fail, also when every digest matches, when the client that applies cannot be built or the read of a rendered object fails for a reason other than that the object does not exist. It MUST write nothing and MUST keep `status.inventory`. Source: enhancement 0012 ("The apply guard runs on every apply, refusing foreign objects outside the inventory and any object adopted by another instance; a per-object adopt annotation is its only override").

#### Scenario: A failed reconcile writes nothing
- **WHEN** a ModulePackage reconcile that renders fails because the verdict could not be asked
- **THEN** it writes nothing and keeps `status.inventory`
