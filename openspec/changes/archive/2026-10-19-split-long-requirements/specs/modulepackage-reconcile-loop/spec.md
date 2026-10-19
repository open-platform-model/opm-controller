## MODIFIED Requirements

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

## REMOVED Requirements

### Requirement: Status always patched
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "The ReleaseReconciler patches the status at the end of every reconcile attempt", "A reconcile that skips its render is not an attempt and patches only a changed Healthy condition", "The patch of a successful outcome and of a NoOp carries the judged Healthy condition".

### Requirement: A ModulePackage records its instance identity
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "A ModulePackage carries the identity of the instance it renders and the identity before an unsettled change", "The identity fields are written, kept and cleared as on a ModuleInstance", "A ModulePackage stored before the fields existed gains the identity", "The prune and the deletion cleanup of a ModulePackage judge with the recorded identities", "A ModulePackage with no recorded identity deletes what it deleted before the field existed".

### Requirement: A ModulePackage guards every apply by ownership
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "A ModulePackage reconcile judges every rendered object with the instance's identity", "A ModulePackage reconcile never asks the verdict with the earlier identity of an unsettled change", "A ModulePackage reconcile never asks the verdict with an empty identity", "A ModulePackage reconcile judges exactly as a ModuleInstance reconcile does under the ssa-apply requirement", "A ModulePackage reconcile judges exactly as a ModuleInstance reconcile does under the reconcile-loop-assembly requirement", "The read is made by the client that applies and a skipped render reads nothing", "A reconcile that renders fails when the verdict could not be asked", "The failure reason is ImpersonationFailed when the ServiceAccount is missing or the read is Forbidden", "The failure reason is ApplyFailed on the backoff except when the ServiceAccount is missing or the read is Forbidden", "A failed ModulePackage reconcile writes nothing and keeps the inventory".

## ADDED Requirements

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
