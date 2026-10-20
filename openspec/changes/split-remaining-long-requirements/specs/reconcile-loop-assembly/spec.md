## MODIFIED Requirements

### Requirement: The instance's contract demand is recorded on its status
The reconciler SHALL record on `ModuleInstance.status.requiredContracts` the set of contract FQNs the instance's components declare, as the render reports it, and SHALL write it on every reconcile whose render succeeds. The field SHALL be derived from the rendered instance and SHALL NOT be readable from, or influenced by, anything a module author writes into `spec`.

#### Scenario: A successful render records the demand

- **WHEN** an instance renders successfully
- **THEN** `status.requiredContracts` lists, sorted and without duplicates, every contract FQN its components declare

#### Scenario: A failed render leaves the previous value

- **WHEN** a render fails
- **THEN** `status.requiredContracts` is unchanged from the last successful render

#### Scenario: The demand follows the spec

- **WHEN** an instance's module is changed to one whose components declare a different contract set, and it renders successfully
- **THEN** `status.requiredContracts` reports the new set and no longer reports the contracts only the previous module declared

### Requirement: A reconcile the apply verdict refuses is retried and not stalled
When the apply verdict refuses a reconcile (`ssa-apply`, "An apply is judged by the ownership verdict before its first write"), the reconciler MUST set `Ready=False` with reason `ApplyRefused`, MUST NOT set `Stalled`, MUST count the attempt as a failed apply and MUST retry on the bounded backoff. `status.inventory`, the applied digests, `status.instanceUUID` and `status.previousInstanceUUID` MUST keep their values.

#### Scenario: A refusal is transient
- **GIVEN** a changed render that names a live object OPM does not manage
- **WHEN** the controller reconciles
- **THEN** `Ready` is `False` with reason `ApplyRefused`, `Stalled` is absent, and the reconcile is requeued on the backoff

#### Scenario: The refusal ends when the object is adopted
- **GIVEN** the instance of the scenario above
- **WHEN** a user annotates the object `opmodel.dev/adopt` with the instance's UUID and the retry runs
- **THEN** the render is applied and `Ready` is `True`

#### Scenario: A refused identity change stores nothing
- **GIVEN** a ModuleInstance whose `spec.module.path` changed, and a rendered object that the verdict refuses
- **WHEN** the controller reconciles
- **THEN** `status.instanceUUID` keeps the earlier identity and `status.previousInstanceUUID` stays empty

#### Scenario: An inventoried object being deleted on unchanged digests
- **GIVEN** a Ready ModuleInstance with unchanged digests and an inventoried Deployment that carries a deletion timestamp
- **WHEN** the controller reconciles
- **THEN** `Ready` does not change to `ApplyRefused`

#### Scenario: A rendered object held by another instance on unchanged digests
- **GIVEN** a ModuleInstance with unchanged digests whose render names an object that exists outside its inventory and carries another instance's UUID label and no adopt annotation
- **WHEN** the controller reconciles
- **THEN** nothing is written and `Ready` is `False` with reason `ApplyRefused`

## REMOVED Requirements

### Requirement: Status always patched
**Reason**: Split into ten requirements that each state one behavior.
**Migration**: See "Every reconcile attempt ends with a status patch", "A successful outcome patches the Healthy condition it judged", "A NoOp patch is bounded to a fixed set of fields", "A NoOp patch leaves lastAttempted, history and inventory alone", "A reconcile that skips its render patches only a changed Healthy condition", "requiredContracts is in the NoOp patch set", "A NoOp patch carries requiredContracts because the demand can move without a digest", "lastAppliedVersion is in the NoOp patch set", "lastAppliedInputs is in the NoOp patch set while the drift render interval is positive" and "Storm safety comes from the generation predicate".

### Requirement: Inventory updated only on full success
**Reason**: Split into five requirements that each state one behavior.
**Migration**: See "The inventory is replaced only after a full success", "An expired Job entry leaves the inventory", "An entry the apply verdict refuses as adopted-elsewhere leaves the inventory", "The inventory digest stays the digest of the rendered set" and "A refused or failed reconcile keeps the inventory".

### Requirement: A rendered resource may be withheld from apply
**Reason**: Split into three requirements that each state one behavior.
**Migration**: See "The reconciler evaluates rendered resources and may withhold one", "Withholding does not change the inventory" and "A reconcile that withheld a resource is not ready and not a no-op".

### Requirement: A recovered panic is recorded as a failed attempt
**Reason**: Split into six requirements that each state one behavior.
**Migration**: See "A recovered panic is never reported as a success", "A panic is recorded in Ready and history", "A panic counts one failure", "A panic keeps the applied state", "A panicking attempt counts as FailedTransient for counters, history and metrics only" and "A recovered panic is raised again with its original value".

### Requirement: The applied module version is recorded in plain text
**Reason**: Split into six requirements that each state one behavior.
**Migration**: See "The applied module version is recorded on status.lastAppliedVersion", "The applied version is written with the lastApplied fields on a success", "A NoOp that rendered sets the applied version", "A NoOp fills an empty applied version and a render-skipping reconcile leaves it", "A version that cannot be read does not fail the reconcile" and "The applied version is not an input to no-op detection".

### Requirement: A healthy ModuleInstance is reconciled periodically
**Reason**: Split into five requirements that each state one behavior.
**Migration**: See "The instance reconcile interval comes from a flag", "A reconciled instance requeues after the interval plus a random addition", "The periodic requeue sets no nextRetryAt and is not a failure", "A suspended, CLI-owned or operator instance is not requeued by the interval" and "An unchanged periodic reconcile ends as a skipped render or a NoOp".

### Requirement: The instance identities are stored before the first write of an apply
**Reason**: Split into seven requirements that each state one behavior.
**Migration**: See "An instance carries its identity and the earlier one until a change is settled", "A first identity is stored", "A changed identity keeps the earlier one", "A return to the earlier identity swaps the two", "previousInstanceUUID is cleared only by a full success", "A reconcile that applies nothing writes an empty instanceUUID and leaves both fields alone otherwise" and "An unsettled identity change is applied and pruned even when the render equals the last apply".

## ADDED Requirements

### Requirement: Every reconcile attempt ends with a status patch
The reconciler MUST patch `ModuleRelease.status` at the end of every reconcile attempt, including `NoOp`. On meaningful outcomes (`Applied`, `AppliedAndPruned`, `FailedTransient`, `FailedStalled`), the patch SHALL update conditions, `lastAttempted*`, `lastApplied*` (on success), `inventory` (on success), history, failure counters, and `nextRetryAt`.

#### Scenario: Status updated on failure
- **WHEN** a phase fails
- **THEN** status conditions, `lastAttempted*` fields, `failureCounters`, and
  `nextRetryAt` are updated

### Requirement: A successful outcome patches the Healthy condition it judged
As in requirement "Every reconcile attempt ends with a status patch": On a successful outcome the patch SHALL also carry the `Healthy` condition the reconcile judged (`instance-health`). A failed attempt SHALL leave `Healthy` as it was.

#### Scenario: A failed attempt leaves Healthy as it was
- **WHEN** a reconcile attempt fails
- **THEN** the `Healthy` condition is as it was

### Requirement: A NoOp patch is bounded to a fixed set of fields
As in requirement "Every reconcile attempt ends with a status patch": On `NoOp`, the patch SHALL be bounded to: drift condition (`Drifted`), the `Healthy` condition (`instance-health`), failure counter deltas (incl. drift counter), `nextRetryAt` clearing, `requiredContracts`, and `lastAppliedVersion` and `lastAppliedInputs` when the attempt rendered, and the removal of expired Job entries from `inventory.entries` (`drift-detection`, "A missing object is restored").

#### Scenario: NoOp removes an expired Job entry
- **WHEN** the outcome is `NoOp` and a rendered Job that sets
  `ttlSecondsAfterFinished` does not exist on the cluster
- **THEN** the patch removes that Job from `status.inventory.entries` and sets
  `status.inventory.count` to the number of entries left
- **AND** `status.inventory.digest` and `status.inventory.revision` keep their
  values

### Requirement: A NoOp patch leaves lastAttempted, history and inventory alone
`lastAttempted*` and history MUST NOT be modified on `NoOp`, and `inventory` MUST NOT be modified in any other way: those fields describe meaningful reconcile outcomes. (Requirement: A NoOp patch is bounded to a fixed set of fields)

#### Scenario: NoOp patches drift and counters
- **WHEN** the outcome is `NoOp` and drift detection ran (with or without
  detecting drift)
- **THEN** `status.conditions[Drifted]` reflects the drift detection result
- **AND** `status.failureCounters` reflects phase counter deltas
  (incl. `drift` counter increment on dry-run failure)
- **AND** `status.nextRetryAt` is cleared
- **AND** `status.lastAttempted*` and `status.history` are NOT modified
- **AND** `status.inventory` is NOT modified, except that an expired Job entry
  is removed

### Requirement: A reconcile that skips its render patches only a changed Healthy condition
A reconcile that skips its render because its inputs are unchanged (`render-input-key`) SHALL NOT be an attempt. It judges health (`instance-health`) and MUST NOT patch anything but the `Healthy` condition, which it patches only when the judgement changed it: every other condition it read is already final, and nothing else it could record has moved.

#### Scenario: A skipped render patches nothing
- **WHEN** a reconcile skips its render because its inputs are unchanged and
  its `Healthy` judgement is the one already on status
- **THEN** no status patch is sent
- **AND** `lastAttemptedAt`, history and `status.lastAppliedInputs.renderedAt`
  keep their values

#### Scenario: A skipped render records a changed health judgement
- **WHEN** a reconcile skips its render and the inventory's Deployment has
  rolled out since the last judgement (`Healthy` was `NotRolledOut`)
- **THEN** the only status change is `Healthy=True` with reason `RolledOut`
- **AND** `observedGeneration`, `lastAttemptedAt`, history and
  `status.lastAppliedInputs.renderedAt` keep their values

### Requirement: requiredContracts is in the NoOp patch set
`requiredContracts` SHALL be in the `NoOp` set deliberately, and it is the one field there that does not describe an outcome.

#### Scenario: NoOp refreshes the demand the render just reported
- **WHEN** the outcome is `NoOp` and the render that preceded it reported a
  different contract demand than the status carries
- **THEN** `status.requiredContracts` is updated to what the render reported

### Requirement: A NoOp patch carries requiredContracts because the demand can move without a digest
As in requirement "requiredContracts is in the NoOp patch set": It SHALL describe what the instance demands, which a reconcile can observe as changed without any digest changing: a regenerated platform re-enqueues every instance, and the render that follows may be a no-op for apply while being the only evidence that the demand moved. Bounding it out of the `NoOp` patch would leave the field describing a render the operator has already superseded.

#### Scenario: A regenerated platform moves the demand without a digest change
- **WHEN** a regenerated platform re-enqueues an instance and the render that follows is a no-op for apply but reports a different contract demand
- **THEN** the `NoOp` patch carries the demand the render reported

### Requirement: lastAppliedVersion is in the NoOp patch set
`lastAppliedVersion` SHALL be in the `NoOp` set because a `NoOp` re-proves it. A `NoOp` means the source digest equals the last apply's, and the source digest pins the module version, so the version the render reports is the one last applied, whoever applied it. Writing it fills the field on an object applied before it existed and corrects it after an ownership handback. A reconcile that did not render SHALL leave it as it was.

#### Scenario: A NoOp that rendered writes the version it reported
- **WHEN** the outcome is `NoOp` and the attempt rendered
- **THEN** `status.lastAppliedVersion` is the version the render reports

### Requirement: lastAppliedInputs is in the NoOp patch set while the drift render interval is positive
`lastAppliedInputs` SHALL be in the `NoOp` set for the same reason: a `NoOp` that rendered proves the cluster holds what those inputs produce, so its key and render time are recorded, and the drift render interval counts from it. It SHALL be written on a `NoOp` only while the drift render interval is greater than zero; with the skip disabled nothing reads it. (Requirement: lastAppliedVersion is in the NoOp patch set)

#### Scenario: A NoOp does not write lastAppliedInputs while the drift render interval is zero
- **WHEN** the outcome is `NoOp` and the attempt rendered, and the drift render interval is zero
- **THEN** `status.lastAppliedInputs` is not written, because with the skip disabled nothing reads it

### Requirement: Storm safety comes from the generation predicate
Storm safety SHALL be provided by `WithEventFilter(predicate.GenerationChangedPredicate{})` on the controller. Status subresource patches do not bump `metadata.generation`, so the predicate filters them at the watch boundary.

#### Scenario: NoOp patch does not trigger reconcile loop
- **WHEN** a NoOp patch updates `status` (drift condition or counters)
- **THEN** the resulting watch event is filtered by `GenerationChangedPredicate`
- **AND** no follow-up reconcile is queued from the status patch alone

### Requirement: The inventory is replaced only after a full success
The `status.inventory` MUST only be replaced after a fully successful apply (and prune, if enabled). Two narrower writes are allowed. (Requirement: An expired Job entry leaves the inventory) (Requirement: An entry the apply verdict refuses as adopted-elsewhere leaves the inventory)

#### Scenario: Partial failure preserves inventory
- **WHEN** apply succeeds but prune fails
- **THEN** `status.inventory` remains at the previous successful value

### Requirement: An expired Job entry leaves the inventory
A ModuleInstance reconcile whose digests are unchanged MUST remove the entries of expired Jobs (`drift-detection`, "A missing object is restored") from `status.inventory.entries`, on a `NoOp` and on a restore.

#### Scenario: An expired Job leaves the inventory without an apply
- **WHEN** a reconcile renders with unchanged digests, sends no apply, and a rendered Job that sets `ttlSecondsAfterFinished` does not exist
- **THEN** `status.inventory.entries` no longer lists the Job
- **AND** the next reconcile with unchanged digests is a `NoOp`

### Requirement: An entry the apply verdict refuses as adopted-elsewhere leaves the inventory
A reconcile that renders MUST leave out of the inventory it records, on a `NoOp` as well, the entry of every object the apply verdict refuses as `adopted-elsewhere` (enhancement 0012, "The apply guard runs on every apply, refusing foreign objects outside the inventory and any object adopted by another instance; a per-object adopt annotation is its only override").

#### Scenario: An adopted object leaves the inventory without an apply
- **GIVEN** a Ready ModuleInstance and an inventoried ConfigMap that a user annotates `opmodel.dev/adopt` with another instance's UUID
- **WHEN** a reconcile renders with unchanged digests
- **THEN** it sends no apply, `status.inventory.entries` no longer lists the ConfigMap and `Ready` stays `True`
- **AND** the next reconcile with unchanged digests is a `NoOp`

### Requirement: The inventory digest stays the digest of the rendered set
`status.inventory.digest` MUST stay the digest of the rendered set, so after such a removal it is not the digest of the entries listed. (Requirement: An expired Job entry leaves the inventory) (Requirement: An entry the apply verdict refuses as adopted-elsewhere leaves the inventory)

#### Scenario: A removal leaves the digest of the rendered set
- **WHEN** an entry leaves `status.inventory.entries` by one of the two removals
- **THEN** `status.inventory.digest` is still the digest of the rendered set, not the digest of the entries listed

### Requirement: A refused or failed reconcile keeps the inventory
A reconcile that is refused by the apply verdict, or that fails, MUST NOT change `status.inventory`.

#### Scenario: A refused reconcile keeps the inventory
- **WHEN** the apply verdict refuses a reconcile
- **THEN** `status.inventory` holds the entries it held before

### Requirement: A reconcile that does not render leaves requiredContracts at its previous value
As in requirement "The instance's contract demand is recorded on its status": A reconcile that does not render SHALL leave the field at its previous value. This includes a failed render, a suspended instance and a CLI-owned instance: a stale value over-reports demand, which blocks a deletion that could have proceeded, and the opposite error would let a provider abandon its dependents.

#### Scenario: A suspended instance leaves the previous value
- **WHEN** a suspended instance is reconciled
- **THEN** `status.requiredContracts` is unchanged from the last successful render

### Requirement: The reconciler evaluates rendered resources and may withhold one
The reconciler SHALL evaluate the rendered resources before applying them and MAY withhold an individual resource whose application would break a guarantee the operator holds. A withheld resource SHALL NOT be applied, and every other rendered resource SHALL be applied as usual, so one refused object does not strand the rest of the module.

#### Scenario: The rest of the module still applies

- **WHEN** one rendered resource is withheld
- **THEN** every other rendered resource is applied, and the apply is not reported as failed

#### Scenario: Nothing withheld behaves as before

- **WHEN** no rendered resource is withheld
- **THEN** apply, inventory, prune and the reported outcome are exactly as they were

### Requirement: Withholding does not change the inventory
As in requirement "The reconciler evaluates rendered resources and may withhold one": Withholding SHALL NOT change the inventory. The inventory records which resources the instance owns, and the operator still owns a resource it declined to update; removing it would mark it stale and prune the very object the refusal exists to protect.

#### Scenario: The withheld resource stays owned

- **WHEN** a resource is withheld
- **THEN** the inventory still lists it, and the prune does not delete it

### Requirement: A reconcile that withheld a resource is not ready and not a no-op
A reconcile that withheld a resource SHALL NOT report ready and SHALL NOT be treated as a no-op, because the cluster does not hold what the render produced.

#### Scenario: A withheld resource keeps the instance unconverged

- **WHEN** a reconcile withholds a resource
- **THEN** the instance reports not ready, and the outcome is not a no-op

### Requirement: A recovered panic is never reported as a success
When a ModuleInstance or ModulePackage reconcile panics after its deferred status commit is installed, the reconciler MUST NOT report the attempt as a success.

#### Scenario: A panicking render does not mark a ready instance Ready

- **WHEN** a ModuleInstance that is `Ready=True` from an earlier successful reconcile is reconciled and its renderer panics
- **THEN** the panic propagates out of the reconcile with its original value
- **AND** the stored instance has `Ready=False` with reason `ReconcilePanic`, `Reconciling=True`, no `Stalled` condition, `failureCounters.reconcile` incremented, a failure history entry and no `nextRetryAt`
- **AND** its `inventory` and `lastAppliedRenderDigest` are unchanged from the earlier reconcile

#### Scenario: A panicking render does not mark a ready package Ready

- **WHEN** a ModulePackage that is `Ready=True` from an earlier successful reconcile is reconciled and its renderer panics
- **THEN** the panic propagates out of the reconcile with its original value
- **AND** the stored package has `Ready=False` with reason `ReconcilePanic`, `Reconciling=True`, no `Stalled` condition, `failureCounters.reconcile` incremented, a failure history entry and no `nextRetryAt`
- **AND** its `inventory` and `lastAppliedRenderDigest` are unchanged from the earlier reconcile

### Requirement: A panic is recorded in Ready and history
Before the panic leaves the reconcile function, the deferred commit MUST record the attempt as a transient failure and patch status:
- `Ready=False` with reason `ReconcilePanic` and a message that starts with `reconcile panicked: ` followed by the panic value.
- `Reconciling=True`, and no `Stalled` condition.
- `lastAttempted*` set for this attempt, and a failure history entry carrying the same message.
(Requirement: A panic counts one failure) (Requirement: A panic keeps the applied state)

#### Scenario: A panic is recorded in Ready and history
- **WHEN** a reconcile panics after its deferred status commit is installed
- **THEN** the deferred commit records `Ready=False` with reason `ReconcilePanic` and a message that starts with `reconcile panicked: ` followed by the panic value, `Reconciling=True` and no `Stalled` condition, `lastAttempted*` set for this attempt, and a failure history entry carrying the same message

### Requirement: A panic counts one failure
Before the panic leaves the reconcile function, the deferred commit MUST record the attempt as a transient failure and patch status:
- `failureCounters.reconcile` incremented, as for any `FailedTransient` outcome.
- The drift, apply and prune failure counters left as they were. A phase that was in flight when the panic hit did not succeed, so a panicking attempt never resets a phase counter.
(Requirement: A panic is recorded in Ready and history) (Requirement: A panic keeps the applied state)

#### Scenario: A panic inside a phase does not reset that phase's counter

- **WHEN** a reconcile panics inside drift detection (ModuleInstance) or apply (ModulePackage), after that phase has started
- **THEN** the drift, apply and prune failure counters keep their earlier values and `failureCounters.reconcile` is incremented

### Requirement: A panic keeps the applied state
Before the panic leaves the reconcile function, the deferred commit MUST record the attempt as a transient failure and patch status:
- `nextRetryAt` cleared, because the retry is scheduled by the controller runtime's rate limiter, not by the operator's backoff.
- `lastApplied*` and `inventory` left as they were.
(Requirement: A panic is recorded in Ready and history) (Requirement: A panic counts one failure)

#### Scenario: A panic clears nextRetryAt and keeps the applied state
- **WHEN** a reconcile panics after its deferred status commit is installed
- **THEN** `nextRetryAt` is cleared, and `lastApplied*` and `inventory` are as they were

### Requirement: A panicking attempt counts as FailedTransient for counters, history and metrics only
A panicking attempt SHALL be counted as `FailedTransient` for failure counters, history and reconcile metrics only. It SHALL return no `ctrl.Result`: the requeue scenarios of "Outcome classification" and the `nextRetryAt` rule of the reconcile-backoff capability do not apply to it. Fields the attempt already wrote on the object before the panic (`requiredContracts`, `instanceUUID`, the `Drifted` condition, a ModulePackage's `status.source`) SHALL be committed as on any other failed attempt.

#### Scenario: A panicking attempt returns no ctrl.Result
- **WHEN** a reconcile panics after its deferred status commit is installed
- **THEN** the attempt is counted as `FailedTransient` for failure counters, history and reconcile metrics only, and it returns no `ctrl.Result`

### Requirement: A recovered panic is raised again with its original value
The reconciler MUST then re-panic with the original panic value, so the controller runtime still recovers, logs (with the stack) and counts the panic and requeues the object. The operator SHOULD log the panic value with the object's name and namespace before it re-panics.

#### Scenario: A panic during conversion is recorded and frees the slot

- **WHEN** the export of a ModuleInstance's render result for apply panics while the reconcile holds a render slot
- **THEN** the panic propagates with its original value, the render slot is free again, and the stored instance has `Ready=False` with reason `ReconcilePanic`

### Requirement: The applied module version is recorded on status.lastAppliedVersion
The ModuleInstance and ModulePackage reconcilers SHALL record on `status.lastAppliedVersion` the version of the module whose render they last applied, as plain text beside the `lastApplied*` digests. The value SHALL be the version the rendered instance's source module declares in its `metadata.version` (bare SemVer, for example `0.1.0`), read the same way on both kinds; for a ModuleInstance it is not copied from `spec.module.version`.
**Source**: owner decision j5 of the kernel-plan walkthrough (2026-10-03): "Operator adds an additive status.lastAppliedVersion (plain text) next to the digest." The ADR-008 amendment (library#167) carries only j5's deletion half.

#### Scenario: A successful apply records the module version

- **WHEN** a ModuleInstance whose module declares `metadata.version: "0.1.0"` is reconciled and its apply succeeds
- **THEN** `status.lastAppliedVersion` is `0.1.0`
- **AND** it is written in the same status patch as `lastAppliedAt` and the `lastApplied*` digests

#### Scenario: A ModulePackage records the version of the module its instance renders

- **WHEN** a ModulePackage whose `instance.cue` renders a module declaring `metadata.version: "0.1.0"` is reconciled and its apply succeeds
- **THEN** `status.lastAppliedVersion` is `0.1.0`

### Requirement: The applied version is written with the lastApplied fields on a success
As in requirement "The applied module version is recorded on status.lastAppliedVersion": The field SHALL be written together with the other `lastApplied*` fields, on an attempt whose apply (and prune, when enabled) succeeded. An attempt that fails, panics, is suspended, belongs to a CLI-owned instance, or is the operator's own instance (refused) SHALL leave it as it was.

#### Scenario: A failed apply keeps the previous version

- **WHEN** an instance that recorded `lastAppliedVersion` `0.1.0` is changed to a module at `0.2.0` and the apply fails
- **THEN** `status.lastAppliedVersion` is still `0.1.0`
- **AND** the `lastApplied*` digests are unchanged too

### Requirement: A NoOp that rendered sets the applied version
As in requirement "The applied module version is recorded on status.lastAppliedVersion": When the outcome is `NoOp` and the attempt rendered, the reconciler SHALL set it to the version that render reported, whatever the field held: a `NoOp` means every digest equals the last apply's, and the source digest pins the module version, so the rendered version is the applied one.

#### Scenario: A NoOp that rendered sets the version the render reported
- **WHEN** the outcome is `NoOp` and the attempt rendered
- **THEN** `status.lastAppliedVersion` is the version that render reported, whatever the field held

### Requirement: A NoOp fills an empty applied version and a render-skipping reconcile leaves it
This fills the field on an object applied before it existed and corrects it after an ownership handback, where the cli moved the digests without touching this field. A reconcile that did not render (for example one that skips the render because its inputs are unchanged) SHALL leave the field as it was. (Requirement: A NoOp that rendered sets the applied version)

#### Scenario: A NoOp fills an empty field

- **WHEN** an object that was applied before the field existed has no `lastAppliedVersion` and its next reconcile is a `NoOp`
- **THEN** `status.lastAppliedVersion` is set to the version the render reported
- **AND** `lastAttempted*`, `inventory` and history are not modified

#### Scenario: A NoOp after an ownership handback corrects a stale version

- **WHEN** `status.lastAppliedVersion` holds `9.9.9`, the module renders at `0.1.0`, and the reconcile is a `NoOp` because the `lastApplied*` digests already match that render
- **THEN** `status.lastAppliedVersion` is `0.1.0`

#### Scenario: A reconcile that skips its render leaves the applied version
- **WHEN** a reconcile skips its render because its inputs are unchanged
- **THEN** `status.lastAppliedVersion` is as it was

### Requirement: A version that cannot be read does not fail the reconcile
As in requirement "The applied module version is recorded on status.lastAppliedVersion": A version that cannot be read as a concrete string SHALL NOT fail the reconcile. A successful apply then writes the field empty, so it never names a version the last apply did not render.

#### Scenario: An apply whose render reports no version clears the field

- **WHEN** an instance that recorded `lastAppliedVersion` `0.1.0` renders a changed set whose module version cannot be read, and the apply succeeds
- **THEN** `status.lastAppliedVersion` is absent from the object

### Requirement: The applied version is not an input to no-op detection
As in requirement "The applied module version is recorded on status.lastAppliedVersion": The version SHALL NOT be an input to no-op detection; the source digest already covers it.

#### Scenario: The applied version is not compared by no-op detection
- **WHEN** a reconcile decides whether its outcome is `NoOp`
- **THEN** the recorded version is not an input to that decision, because the source digest already covers it

### Requirement: The instance reconcile interval comes from a flag
The operator SHALL take an instance reconcile interval from the flag `--instance-reconcile-interval` (default 10 minutes). An interval of `0` SHALL disable the periodic requeue. The operator MUST refuse to start with a negative interval. The interval is operator-wide: the ModuleInstance API SHALL NOT gain a field for it. (Requirement: A reconciled instance requeues after the interval plus a random addition)

#### Scenario: The interval is disabled

- **WHEN** the interval is `0` and a ModuleInstance reconcile ends rolled out
- **THEN** the reconcile returns without a requeue

#### Scenario: A negative interval stops the operator

- **WHEN** the operator starts with `--instance-reconcile-interval=-1s`
- **THEN** it logs the invalid value and exits with a non-zero status

### Requirement: A reconciled instance requeues after the interval plus a random addition
As in requirement "The instance reconcile interval comes from a flag": When a ModuleInstance reconcile ends with outcome `NoOp`, `Applied` or `AppliedAndPruned`, or skips its render because its inputs are unchanged, and the health judgement asks for no requeue, the reconcile SHALL requeue after the interval plus a random addition of at most 10 percent of it.

#### Scenario: A rolled-out instance requeues on the interval

- **WHEN** an operator-owned ModuleInstance reconcile ends `Applied` with `Healthy=True` and the interval is 10 minutes
- **THEN** the reconcile requeues after at least 10 and at most 11 minutes
- **AND** `status.nextRetryAt` is unset

#### Scenario: A NoOp and a skipped render requeue on the interval

- **WHEN** a later reconcile of the same unchanged instance ends as a `NoOp` or skips its render
- **THEN** it requeues after at least 10 and at most 11 minutes

### Requirement: The periodic requeue sets no nextRetryAt and is not a failure
As in requirement "A reconciled instance requeues after the interval plus a random addition": When the health judgement asks for a requeue, that requeue SHALL apply unchanged. The periodic requeue SHALL NOT set `status.nextRetryAt` and SHALL NOT count as a failure.

#### Scenario: The periodic requeue is not a failure
- **WHEN** a ModuleInstance reconcile ends `Applied` and the periodic requeue is set
- **THEN** `status.nextRetryAt` is unset and the requeue does not count as a failure

### Requirement: A suspended, CLI-owned or operator instance is not requeued by the interval
As in requirement "The instance reconcile interval comes from a flag": A suspended ModuleInstance, a CLI-owned ModuleInstance and the operator's own instance SHALL NOT be requeued by the interval. A failed reconcile SHALL keep the requeue of its failure class.

#### Scenario: A suspended instance is not requeued

- **WHEN** a ModuleInstance has `spec.suspend: true` and the interval is 10 minutes
- **THEN** its reconcile returns without a requeue

#### Scenario: A CLI-owned instance is not requeued

- **WHEN** a ModuleInstance has `spec.owner: cli` and the interval is 10 minutes
- **THEN** its reconcile returns without a requeue

### Requirement: An unchanged periodic reconcile ends as a skipped render or a NoOp
The periodic reconcile of an instance whose inputs and objects are unchanged SHALL end as a skipped render or as a `NoOp`: it sends no apply, records no history entry and changes no condition.

#### Scenario: A skipped periodic reconcile writes nothing

- **WHEN** the periodic reconcile of an unchanged, rolled-out instance skips its render
- **THEN** no patch and no status patch is sent for the instance

### Requirement: An instance carries its identity and the earlier one until a change is settled
A ModuleInstance and a ModulePackage MUST carry in `status.instanceUUID` the identity of the instance the operator last rendered and began to apply, and in the optional field `status.previousInstanceUUID` the identity it had before, for as long as a change of identity is not settled.
**Source**: owner decision of 2026-10-08 (the status holds the earlier and the new identity until the prune of that reconcile succeeded; then the earlier one is cleared).

#### Scenario: An unsettled change of identity carries both identities
- **GIVEN** a change of identity that is not settled
- **THEN** `status.instanceUUID` holds the identity of the instance the operator last rendered and began to apply, and `status.previousInstanceUUID` holds the identity it had before

### Requirement: A first identity is stored
When a render's identity differs from `status.instanceUUID`, the reconciler MUST store the new state in the object's status before it applies any object, and MUST NOT apply when that write fails:
- `status.instanceUUID` empty: it becomes the render's identity and `status.previousInstanceUUID` stays empty.
(Requirement: A changed identity keeps the earlier one) (Requirement: A return to the earlier identity swaps the two)

#### Scenario: First reconcile records the identity
- **GIVEN** a new ModuleInstance whose render carries identity `A`
- **WHEN** the reconcile applies with success
- **THEN** `status.instanceUUID` is `A` and `status.previousInstanceUUID` is empty

#### Scenario: The identities cannot be stored
- **GIVEN** a change of identity, and an API server that refuses the status write
- **WHEN** the controller reconciles
- **THEN** no object of the render is applied and the reconcile is retried

### Requirement: A changed identity keeps the earlier one
When a render's identity differs from `status.instanceUUID`, the reconciler MUST store the new state in the object's status before it applies any object, and MUST NOT apply when that write fails:
- `status.instanceUUID` is `A`, `status.previousInstanceUUID` is empty and the render carries `B`: `status.instanceUUID` becomes `B` and `status.previousInstanceUUID` becomes `A`.
(Requirement: A first identity is stored) (Requirement: A return to the earlier identity swaps the two)

#### Scenario: Both identities are stored before the apply
- **GIVEN** a ModuleInstance with `status.instanceUUID` `A` whose `spec.module.path` changes, so that its render carries identity `B`
- **WHEN** the apply fails
- **THEN** the stored object has `status.instanceUUID` `B` and `status.previousInstanceUUID` `A`
- **AND** `status.inventory` is unchanged

### Requirement: A return to the earlier identity swaps the two
When a render's identity differs from `status.instanceUUID`, the reconciler MUST store the new state in the object's status before it applies any object, and MUST NOT apply when that write fails:
- `status.instanceUUID` is `B`, `status.previousInstanceUUID` is `A` and the render carries `A`: the two values swap.
(Requirement: A first identity is stored) (Requirement: A changed identity keeps the earlier one)

#### Scenario: Going back before the change is settled
- **GIVEN** a ModuleInstance with `status.instanceUUID` `B` and `status.previousInstanceUUID` `A`, whose `spec.module.path` is set back so that its render carries `A`
- **WHEN** the apply starts
- **THEN** `status.instanceUUID` is `A` and `status.previousInstanceUUID` is `B`

### Requirement: previousInstanceUUID is cleared only by a full success
`status.previousInstanceUUID` MUST be cleared in the status commit of a reconcile whose apply and prune succeeded, together with the new inventory, and MUST NOT be cleared by any other outcome: a failed or refused apply, a failed prune, a reconcile that applies nothing, or a recovered panic.

#### Scenario: Full success settles the change
- **GIVEN** a ModuleInstance with `status.instanceUUID` `B` and `status.previousInstanceUUID` `A`
- **WHEN** a reconcile applies and prunes with success
- **THEN** `status.previousInstanceUUID` is empty and `status.instanceUUID` is `B`

#### Scenario: A refused apply keeps both identities
- **GIVEN** the same change of identity, and an apply that succeeds while the reconcile is refused afterwards because dependents remain
- **WHEN** the reconcile ends
- **THEN** `status.instanceUUID` is `B` and `status.previousInstanceUUID` is `A`

#### Scenario: A panic keeps both identities
- **GIVEN** the same change of identity, and a reconcile that panics after the identities were stored
- **WHEN** the panic was recorded as a failed attempt
- **THEN** `status.instanceUUID` is `B` and `status.previousInstanceUUID` is `A`

### Requirement: A reconcile that applies nothing writes an empty instanceUUID and leaves both fields alone otherwise
As in requirement "An instance carries its identity and the earlier one until a change is settled": A reconcile that applies nothing MUST write `status.instanceUUID` when it is empty and MUST leave both fields alone otherwise.

#### Scenario: An object without the field gains it without an apply
- **GIVEN** a ModuleInstance with an inventory and an empty `status.instanceUUID`, whose render equals what was applied
- **WHEN** the controller reconciles and applies nothing
- **THEN** `status.instanceUUID` holds the render's identity and `status.previousInstanceUUID` is empty

### Requirement: An unsettled identity change is applied and pruned even when the render equals the last apply
While `status.previousInstanceUUID` is set, or a render is about to set it, a reconcile MUST NOT end as one that applies nothing, also when the render equals what was last applied: it MUST apply and prune, so that the objects an unsettled apply relabelled are relabelled back and the change is settled.

#### Scenario: Going back renders what was last applied
- **GIVEN** a ModuleInstance whose last successful apply carried identity `A`, with `status.instanceUUID` `B` and `status.previousInstanceUUID` `A` after an apply of `B` whose prune failed, and whose `spec.module.path` is set back so that its render equals what was last applied
- **WHEN** the controller reconciles
- **THEN** the render is applied, and the objects that carry `B` carry `A` again
- **AND** after the reconcile `status.instanceUUID` is `A`, `status.previousInstanceUUID` is empty and `Ready` is `True`

### Requirement: A reconcile with unchanged digests is not refused for an object being deleted
As in requirement "A reconcile the apply verdict refuses is retried and not stalled": A reconcile whose digests are unchanged MUST NOT refuse because an inventoried object is being deleted, when it writes nothing over that object: a no-op, and the restore of a ModuleInstance, which never writes an inventoried object.

#### Scenario: A no-op over an inventoried object being deleted is not refused
- **GIVEN** a Ready ModuleInstance with unchanged digests and an inventoried object that carries a deletion timestamp
- **WHEN** the controller reconciles and writes nothing over that object
- **THEN** the apply verdict does not refuse the reconcile

### Requirement: A ModulePackage reconcile on matching digests is refused until the object is gone
A ModulePackage reconcile that applies on matching digests (`modulepackage-reconcile-loop`, "No-op detection") applies the rendered set, so it would write that object: it MUST be refused, as `ssa-apply` says, until the object is gone. (Requirement: A reconcile with unchanged digests is not refused for an object being deleted)

#### Scenario: A ModulePackage reconcile on matching digests is refused until the object is gone
- **GIVEN** a ModulePackage with unchanged digests and an inventoried object that carries a deletion timestamp
- **WHEN** the controller reconciles and applies the rendered set
- **THEN** the reconcile is refused until the object is gone
