## Purpose

Define the condition types and reason constants the operator writes on its
custom resources, and the helpers that set them, so every reconciler reports
Ready, Reconciling, Stalled and source readiness in one shared vocabulary.

## Requirements

### Requirement: Condition type constants
The `internal/status` package MUST define string constants for condition types: `ReadyCondition`, `ReconcilingCondition`, `StalledCondition`, `SourceReadyCondition`.

#### Scenario: Constants available
- **WHEN** code imports `internal/status`
- **THEN** all four condition type constants are available and match the design doc values

### Requirement: Condition helper functions
The `internal/status` package MUST provide helper functions for common condition transitions.

#### Scenario: Mark reconciling
- **WHEN** `MarkReconciling` is called on a ModuleRelease with a reason and message
- **THEN** the `Reconciling` condition is set to `True` and `Ready` is set to `Unknown`

#### Scenario: Mark stalled
- **WHEN** `MarkStalled` is called on a ModuleRelease with a reason and message
- **THEN** the `Stalled` condition is set to `True` and `Ready` is set to `False`

#### Scenario: Mark ready
- **WHEN** `MarkReady` is called on a ModuleRelease with a message
- **THEN** the `Ready` condition is set to `True` and `Reconciling` and `Stalled` are removed

#### Scenario: Mark source ready
- **WHEN** `MarkSourceReady` is called with artifact revision info
- **THEN** the `SourceReady` condition is set to `True`

#### Scenario: Mark source not ready
- **WHEN** `MarkSourceNotReady` is called with a reason
- **THEN** the `SourceReady` condition is set to `False`

### Requirement: Flux condition interface compliance
`ModuleRelease` MUST implement `conditions.Getter` and `conditions.Setter` interfaces from `fluxcd/pkg/runtime/conditions`.

#### Scenario: Interface satisfaction
- **WHEN** a `*ModuleRelease` is passed to Flux condition helpers
- **THEN** the helpers compile and function correctly

### Requirement: The Drifted condition can be Unknown

The status package SHALL define the reason constant `DriftCheckForbidden` and a helper that sets the `Drifted` condition to `Unknown` with a given reason and message. The helper SHALL NOT change the `Ready`, `Reconciling` or `Stalled` conditions. The reasons used with it are `DriftCheckForbidden` (the API server refused the dry-run of drift detection) and `ImpersonationFailed` (the effective ServiceAccount could not be impersonated, so drift detection did not run).

#### Scenario: Mark drift unknown

- **WHEN** the helper is called on a ModuleInstance with `Ready=True`, reason `DriftCheckForbidden` and a message
- **THEN** `Drifted` is `Unknown` with that reason and message
- **AND** `Ready` is still `True`

#### Scenario: A verdict replaces Unknown

- **GIVEN** a ModuleInstance with `Drifted=Unknown`
- **WHEN** `ClearDrifted` is called
- **THEN** the `Drifted` condition is removed

### Requirement: The Ready message counts the objects adopted by another instance
When a ModuleInstance or ModulePackage reconcile renders, ends `Ready=True` (after an apply, a restore or a no-op) and found one or more rendered objects that the apply verdict refuses as `adopted-elsewhere`, the message of the `Ready` condition MUST state how many rendered objects are adopted by another instance and are not applied. When it found none, the message MUST NOT mention such objects. The message MUST NOT carry an enhancement reference.
**Source**: enhancement 0012 ("The apply guard runs on every apply, refusing foreign objects outside the inventory and any object adopted by another instance; a per-object adopt annotation is its only override").

#### Scenario: The count after a let-go
- **GIVEN** a Ready instance, two of whose rendered objects are annotated for another instance
- **WHEN** a reconcile renders with unchanged digests
- **THEN** `Ready` is `True` and its message states that 2 rendered objects are adopted by another instance and not applied

#### Scenario: The count goes away
- **GIVEN** the instance above
- **WHEN** both objects are annotated back for the instance and a reconcile renders
- **THEN** `Ready` is `True` and its message no longer mentions adopted objects

### Requirement: The SkewRefused reason constant is defined
The status package SHALL define the reason constants used on the Ready condition, including `SkewRefused` (Ready=False: the platform's skew policy is `Refuse` and the module requires a newer catalog build than the platform pins).

#### Scenario: Skew refusal reason is available

- **WHEN** a render is refused under the `Refuse` policy
- **THEN** the reconciler marks `Ready=False` with reason `SkewRefused`

### Requirement: The DuplicateIdentities reason constant is defined
The status package SHALL define the reason constants used on the Ready condition, including `DuplicateIdentities` (Ready=False: two or more rendered objects share one Kubernetes apply identity, so the render is refused before apply naming every producing component; enhancement 0015 ("One registration per module")).

#### Scenario: Duplicate-identity reason is available

- **WHEN** a render is refused because two compiled objects share one apply identity
- **THEN** the reconciler marks `Ready=False` with reason `DuplicateIdentities`

### Requirement: The ReconcilePanic reason constant is defined
The status package SHALL define the reason constants used on the Ready condition, including `ReconcilePanic` (Ready=False, Reconciling=True: a ModuleInstance or ModulePackage reconcile panicked and the attempt was recorded as a failure before the panic propagated).

#### Scenario: Reconcile-panic reason is available

- **WHEN** a ModuleInstance or ModulePackage reconcile panics
- **THEN** the reconciler marks `Ready=False` and `Reconciling=True` with reason `ReconcilePanic` and removes `Stalled`

### Requirement: The RenderTimedOut reason constant is defined
The status package SHALL define the reason constants used on the Ready condition, including `RenderTimedOut` (Ready=False, Reconciling=True, not stalled: a ModuleInstance or ModulePackage render did not finish within `--render-timeout`, and the attempt retries on the backoff).

#### Scenario: Render-timeout reason is available

- **WHEN** a ModuleInstance or ModulePackage render does not finish within `--render-timeout`
- **THEN** the reconciler marks `Ready=False` and `Reconciling=True` with reason `RenderTimedOut` and removes `Stalled`

### Requirement: The existing reason constants stay defined
The status package SHALL define the reason constants used on the Ready condition beside the existing `ResolutionFailed`, `RenderFailed`, `PlatformNotReady` and the Platform reasons `Generated`, `BuildFailed`, `GenerateFailed`.

#### Scenario: Reason constants available
- **WHEN** code imports `internal/status`
- **THEN** all reason constants are available as exported string constants

### Requirement: The Healthy condition type and its reasons are defined
The `internal/status` package SHALL define the condition type `HealthyCondition` (`"Healthy"`) and the reasons `RolledOut` (`Healthy=True`), `NotRolledOut` (`Healthy=False`), `ProgressDeadlineExceeded` (`Healthy=False`, a Deployment's rollout is stalled) and `HealthUnknown` (`Healthy=Unknown`).

#### Scenario: The constants are available

- **WHEN** code imports `internal/status`
- **THEN** `HealthyCondition`, `RolledOutReason`, `NotRolledOutReason`, `ProgressDeadlineExceededReason` and `HealthUnknownReason` are exported string constants with those values

### Requirement: One helper sets Healthy and no other condition
The `internal/status` package SHALL provide one helper that sets `Healthy` from a status, a reason and a message, and touches no other condition.

#### Scenario: The helper sets only Healthy

- **WHEN** the helper sets `Healthy=True` with reason `RolledOut` on an object carrying `Ready=True`
- **THEN** `Healthy` is `True` with reason `RolledOut` and every other condition is unchanged

### Requirement: Healthy is independent of Ready
`Healthy` SHALL be independent of `Ready`: no helper that sets `Ready`, `Reconciling` or `Stalled` changes it, except that `MarkManagedExternally` and `MarkSelfManagementRefused` remove it.

#### Scenario: Marking Ready leaves Healthy alone

- **WHEN** `MarkReady`, `MarkReconciling`, `MarkStalled` or `MarkSuspended` is called on an object carrying `Healthy=False`
- **THEN** `Healthy` is unchanged

### Requirement: A refused claim recreate sets the ClaimConflict reason
The status package SHALL define the reason constant `ClaimConflict`. A ModuleInstance or ModulePackage reconcile whose apply refused the forced recreate of a PersistentVolumeClaim MUST set `Ready=False` with reason `ClaimConflict`. The message MUST name the claim as `<namespace>/<name>`, MUST say that the claim was kept, and MUST name the ways out, `spec.dataPolicy: Delete` among them.

#### Scenario: Delete recreates the claim
- **GIVEN** a ModuleInstance with `spec.rollout.forceConflicts: true` and `spec.dataPolicy: Delete`, whose render changes `storageClassName` of a live PersistentVolumeClaim
- **WHEN** the controller reconciles
- **THEN** the claim is deleted and created again with the new `storageClassName`, and `Ready` is `True`

### Requirement: The ClaimConflict message names the refused field when the check refused the claim
A ModuleInstance or ModulePackage reconcile whose apply refused the forced recreate of a PersistentVolumeClaim MUST set `Ready=False` with reason `ClaimConflict`. When the check before the apply refused the claim, the message MUST also name the refused field and the API server's message and MUST say that nothing was applied.

#### Scenario: A claim refused by the check before the apply
- **WHEN** the check before the apply refused the claim
- **THEN** the message names the refused field and the API server's message and says that nothing was applied

### Requirement: The ClaimConflict message after the delete guard kept the claim
A ModuleInstance or ModulePackage reconcile whose apply refused the forced recreate of a PersistentVolumeClaim MUST set `Ready=False` with reason `ClaimConflict`. When the resource manager's delete guard kept the claim (the claim changed after the check), no refusal is at hand and a stage of the apply is already under way: the message MUST NOT name a field and MUST NOT say that nothing was applied.

#### Scenario: A claim kept by the delete guard
- **GIVEN** an apply whose check passed and whose staged apply then reached for the delete of the PersistentVolumeClaim `media/config`
- **WHEN** the controller reports the reconcile
- **THEN** `Ready` is `False` with reason `ClaimConflict` and a message that contains `media/config` and `spec.dataPolicy`
- **AND** the message names no refused field and does not say that nothing was applied

### Requirement: A ClaimConflict reconcile does not stall, prune or record
A ModuleInstance or ModulePackage reconcile whose apply refused the forced recreate of a PersistentVolumeClaim MUST set `Ready=False` with reason `ClaimConflict`. The reconcile MUST NOT set `Stalled`: it MUST retry on the bounded backoff, because the conflict can be resolved on the claim, which changes nothing on the reconciled object. The reconcile MUST NOT prune, MUST NOT record a new inventory and MUST NOT record the rendered digests as applied.

#### Scenario: A refused claim on a ModuleInstance
- **GIVEN** a ModuleInstance with `spec.rollout.forceConflicts: true` and no `spec.dataPolicy`, whose render changes `storageClassName` of a live PersistentVolumeClaim `media/config`
- **WHEN** the controller reconciles
- **THEN** `Ready` is `False` with reason `ClaimConflict` and a message that contains `media/config` and `spec`
- **AND** `Stalled` is not set, and the result asks for a retry after a backoff
- **AND** the claim has the UID it had
- **AND** no other object of the render was changed

#### Scenario: The conflict clears when the render fits the claim again
- **GIVEN** that ModuleInstance with reason `ClaimConflict`
- **WHEN** the render goes back to the claim's `storageClassName` and the controller reconciles
- **THEN** `Ready` is `True`

#### Scenario: A ClaimConflict reconcile prunes and records nothing
- **WHEN** a reconcile ends with reason `ClaimConflict`
- **THEN** it does not prune, does not record a new inventory and does not record the rendered digests as applied

### Requirement: A refused second identity change sets the IdentityChangeUnsettled reason
The status package SHALL define the reason constant `IdentityChangeUnsettled`. A ModuleInstance or ModulePackage reconcile that refuses a second change of the instance identity while the first is not settled MUST set `Ready=False` and `Stalled=True` with that reason. The message MUST NOT carry an enhancement reference.

#### Scenario: The IdentityChangeUnsettled message carries no enhancement reference
- **WHEN** a reconcile refuses a second change of the instance identity while the first is not settled
- **THEN** the message carries no enhancement reference

### Requirement: The IdentityChangeUnsettled message names the way out
The status package SHALL define the reason constant `IdentityChangeUnsettled`. A ModuleInstance or ModulePackage reconcile that refuses a second change of the instance identity while the first is not settled MUST set `Ready=False` and `Stalled=True` with that reason. The message MUST say that an earlier change of the instance identity is not finished, and MUST name the way out: restore the earlier module path, wait until the object is Ready, then change it again.

#### Scenario: The reason on a refused second change
- **GIVEN** a ModuleInstance whose identity change is not settled and whose render carries a third identity
- **WHEN** the controller reconciles
- **THEN** `Ready` is `False` with reason `IdentityChangeUnsettled` and `Stalled` is `True`
- **AND** the message tells the user to restore the earlier module path

### Requirement: A refused apply sets the ApplyRefused reason
The status package SHALL define the reason constant `ApplyRefused`. A ModuleInstance or ModulePackage reconcile that the apply verdict refuses MUST set `Ready=False` with that reason and MUST NOT set `Stalled`. The message MUST NOT carry an enhancement reference.

#### Scenario: A let-go object does not set the reason
- **GIVEN** a reconcile whose only refused object is adopted by another instance
- **WHEN** the reconcile completes
- **THEN** `Ready` is `True`

### Requirement: The ApplyRefused message states the count
The status package SHALL define the reason constant `ApplyRefused`. A ModuleInstance or ModulePackage reconcile that the apply verdict refuses MUST set `Ready=False` with that reason and MUST NOT set `Stalled`. The message MUST state how many objects were refused and that nothing was applied.

#### Scenario: The message states the count and that nothing was applied
- **WHEN** a ModuleInstance or ModulePackage reconcile is refused by the apply verdict
- **THEN** the message states how many objects were refused and that nothing was applied

### Requirement: The ApplyRefused message carries the library messages
The status package SHALL define the reason constant `ApplyRefused`. A ModuleInstance or ModulePackage reconcile that the apply verdict refuses MUST set `Ready=False` with that reason and MUST NOT set `Stalled`. The message MUST carry, for each refused object, the message the library words for the refusal, unchanged (at most ten objects, and fewer when the messages would take the text past 1024 characters; then the number of the rest).

#### Scenario: The reason on a refused apply
- **GIVEN** a ModuleInstance whose render names a live ConfigMap `team-a/settings` that OPM does not manage
- **WHEN** the controller reconciles
- **THEN** `Ready` is `False` with reason `ApplyRefused` and `Stalled` is absent
- **AND** the message contains the library's message for `ConfigMap/team-a/settings`, which names the `opmodel.dev/adopt` annotation and the instance's UUID

### Requirement: A reconcile that skips its render leaves the Ready message as written
A reconcile that skips its render MUST leave the message of the `Ready` condition as the last render wrote it. Source: enhancement 0012 ("The apply guard runs on every apply, refusing foreign objects outside the inventory and any object adopted by another instance; a per-object adopt annotation is its only override").

#### Scenario: A skipped render keeps the Ready message
- **WHEN** a reconcile skips its render
- **THEN** the message of the `Ready` condition is as the last render wrote it

### Requirement: The deletion wait reasons are defined
The `internal/status` package MUST define the reasons `DeletionInProgress` and `DeletionBlocked` for the `Ready` condition of a ModuleInstance or a ModulePackage that is being deleted, and the event reason `DeletionUnconfirmed`. Neither reason MUST appear on an object that is not being deleted.

#### Scenario: The deletion reasons are defined
- **WHEN** code imports `internal/status`
- **THEN** the reasons `DeletionInProgress` and `DeletionBlocked` for the `Ready` condition and the event reason `DeletionUnconfirmed` are defined

#### Scenario: An object that is not being deleted carries neither reason
- **WHEN** a ModuleInstance or a ModulePackage is not being deleted
- **THEN** neither `DeletionInProgress` nor `DeletionBlocked` appears on it

### Requirement: DeletionInProgress is set while the cleanup waits
`DeletionInProgress` MUST be set with `Ready=False` and `Reconciling=True`, never with `Stalled=True`: the cleanup sent every delete and waits for the objects to be gone.

#### Scenario: Waiting is not stalled
- **GIVEN** a deleting ModuleInstance whose deleted Deployment still exists and was deleted one minute ago
- **WHEN** the reconcile ends
- **THEN** `Ready` is False with reason `DeletionInProgress`, `Reconciling` is True and `Stalled` is absent or False

### Requirement: DeletionBlocked is set when an object has not gone for 10 minutes
`DeletionInProgress` MUST be set with `Ready=False` and `Reconciling=True`, never with `Stalled=True`: the cleanup sent every delete and waits for the objects to be gone.
`DeletionBlocked` MUST be set with `Ready=False` and `Stalled=True`: an object has not gone for 10 minutes. Both messages MUST name the remaining objects and their finalizers; the `DeletionBlocked` message MUST also name the ways out.

#### Scenario: Blocked is stalled and names the way out
- **GIVEN** a deleting ModuleInstance whose deleted Deployment has existed for more than 10 minutes
- **WHEN** the reconcile ends
- **THEN** `Ready` is False with reason `DeletionBlocked` and `Stalled` is True
- **AND** the message names the object, its finalizers, and `spec.prune=false` as a way to release the instance

### Requirement: A deletion wait message lists at most ten objects
`DeletionInProgress` MUST be set with `Ready=False` and `Reconciling=True`, never with `Stalled=True`: the cleanup sent every delete and waits for the objects to be gone.
`DeletionBlocked` MUST be set with `Ready=False` and `Stalled=True`: an object has not gone for 10 minutes. Both messages MUST name the remaining objects and their finalizers; the `DeletionBlocked` message MUST also name the ways out. A message MUST list at most ten objects and say how many more remain.

#### Scenario: A long list is cut
- **GIVEN** a deleting ModuleInstance with fifteen deleted objects that still exist
- **WHEN** the reconcile ends
- **THEN** the message names ten of them and says that five more remain

### Requirement: A deletion wait message names the finalizer that holds the object first
A message MUST list at most ten objects and say how many more remain. For each object it MUST name at most three finalizers and count the rest, and it MUST name a finalizer other than `foregroundDeletion` first: `foregroundDeletion` is the garbage collector's and is set again by every repeated delete, so it is rarely what holds the object. When `foregroundDeletion` is the only finalizer, the message MUST say that the object waits for its dependents.

#### Scenario: The finalizer that holds the object is named first
- **GIVEN** a terminating object with the finalizers `foregroundDeletion`, `example.com/a`, `example.com/b`, `example.com/c` and `example.com/d`
- **WHEN** the message is built
- **THEN** it names `example.com/a`, `example.com/b` and `example.com/c` and says that two more finalizers are set

#### Scenario: Only the collector's finalizer
- **GIVEN** a terminating Deployment whose only finalizer is `foregroundDeletion`
- **WHEN** the message is built
- **THEN** it says that the Deployment waits for its dependents to be deleted
