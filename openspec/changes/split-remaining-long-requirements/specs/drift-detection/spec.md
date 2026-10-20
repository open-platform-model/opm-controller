## REMOVED Requirements

### Requirement: Drift runs on no-op reconciles
**Reason**: Split into two requirements that each state one behavior.
**Migration**: See "Drift detection runs on every reconcile that renders" and "A reconcile that skips its render runs no drift detection and drift runs at most once per interval".

### Requirement: A withheld resource is excluded from drift detection
**Reason**: Split into two requirements that each state one behavior.
**Migration**: See "Drift detection skips a resource the reconcile withheld from apply" and "The refusal that withheld a resource carries the signal that drift does not".

### Requirement: A drift check that is refused or has no identity is reported
**Reason**: Split into five requirements that each state one behavior.
**Migration**: See "A Forbidden read or dry-run sets Drifted to Unknown", "Forbidden is every 403 answer", "A ServiceAccount that cannot be impersonated sets Drifted to Unknown", "A refused or unidentified drift check counts as failed and restores nothing" and "A later drift detection replaces the condition and a successful apply removes it".

### Requirement: An object adopted by another instance is excluded from drift detection and from the restore
**Reason**: Split into three requirements that each state one behavior.
**Migration**: See "An object adopted by another instance is neither drift nor restored", "The read the verdict needs is the one read per object before the dry-run" and "A failed read on unchanged digests is a failed drift check that restores nothing".

## MODIFIED Requirements

### Requirement: A missing object is restored
An object that the render produces and that does not exist on the cluster is not drift, and the `Drifted` condition SHALL NOT report it. When a ModuleInstance reconcile renders, finds every digest unchanged, withholds nothing, and its dry-run shows that one or more rendered objects do not exist, the controller SHALL apply those missing objects through the identity that applies the instance.

#### Scenario: A deleted object is created again

- **GIVEN** a Ready ModuleInstance whose ConfigMap `foo` was deleted by hand
- **WHEN** the controller reconciles and renders with unchanged digests
- **THEN** ConfigMap `foo` exists again with the rendered content
- **AND** the outcome is `Applied` and `status.lastAppliedAt` moves

#### Scenario: Only the missing object is applied

- **GIVEN** a Ready ModuleInstance with ConfigMap `foo` deleted and ConfigMap `bar` modified by hand
- **WHEN** the controller reconciles and renders with unchanged digests
- **THEN** `foo` is created again
- **AND** `bar` keeps its modified content and `Drifted=True` reports it

#### Scenario: A finished Job with a TTL stays absent

- **GIVEN** a Ready ModuleInstance whose rendered Job sets `ttlSecondsAfterFinished` and no longer exists
- **WHEN** the controller reconciles and renders with unchanged digests
- **THEN** the Job is not created and the outcome is `NoOp`
- **AND** `status.inventory.entries` no longer lists the Job, and `status.inventory.digest` and `status.inventory.revision` keep their values
- **AND** `Healthy` is `True` with reason `RolledOut` when the other inventory objects are healthy

#### Scenario: A Job with a TTL that was removed before it ran

- **GIVEN** a Ready ModuleInstance whose rendered Job sets `ttlSecondsAfterFinished` and was deleted before any Pod of it ran
- **WHEN** the controller reconciles and renders with unchanged digests
- **THEN** the Job is not created, because the controller cannot tell it from a Job that finished
- **AND** the Job leaves `status.inventory.entries` and does not make the instance unhealthy

#### Scenario: A Job with a TTL that failed and expired

- **GIVEN** a ModuleInstance with `Healthy=False` because its rendered Job with `ttlSecondsAfterFinished` has condition `Failed`
- **WHEN** the cluster removes the Job and the controller reconciles with unchanged digests
- **THEN** the Job is not created, it leaves `status.inventory.entries`, and `Healthy` is judged over the remaining entries

#### Scenario: A missing Job without a TTL is created again

- **GIVEN** a Ready ModuleInstance whose rendered Job sets no `ttlSecondsAfterFinished` and was deleted
- **WHEN** the controller reconciles
- **THEN** the reconcile renders, the Job is created again and the outcome is `Applied`
- **AND** `status.inventory.entries` still lists the Job

#### Scenario: An expired Job beside a deleted object

- **GIVEN** a Ready ModuleInstance whose Job with a TTL expired and whose ConfigMap `foo` was deleted by hand
- **WHEN** the controller reconciles and renders with unchanged digests
- **THEN** `foo` is created again, the Job is not, and the outcome is `Applied`
- **AND** the new `status.inventory` lists `foo` and not the Job, with the digest of the rendered set

#### Scenario: A restore that fails

- **GIVEN** a Ready ModuleInstance whose ConfigMap `foo` was deleted by hand
- **WHEN** the controller reconciles, renders with unchanged digests, and the apply of `foo` fails
- **THEN** `Ready` is `False` with reason `ApplyFailed`, the reconcile requeues on the transient backoff and `status.nextRetryAt` is set

#### Scenario: Nothing is missing

- **GIVEN** a Ready ModuleInstance whose rendered objects all exist
- **WHEN** the controller reconciles and renders with unchanged digests
- **THEN** no apply is sent and the outcome is `NoOp`

#### Scenario: A taken-in object is applied by the restore step

- **GIVEN** a Ready ModuleInstance with unchanged digests, a rendered ConfigMap `foo` that exists outside `status.inventory` with the instance's adopt annotation and other data than the render, and an inventoried ConfigMap `bar` modified by hand
- **WHEN** the controller reconciles and renders
- **THEN** `foo` has the rendered content and is listed in `status.inventory`, and the outcome is `Applied`
- **AND** `bar` keeps its modified content, and `Drifted=True` names `bar` and does not name `foo`

### Requirement: Drift detection failure increments counter
If the SSA dry-run API call fails, or drift detection cannot run because the effective ServiceAccount cannot be impersonated, the controller MUST increment `status.failureCounters.drift`. A failure other than a `Forbidden` answer or a failed impersonation leaves the `Drifted` condition as it was; those two set it to `Unknown` ("A Forbidden read or dry-run sets Drifted to Unknown" and "A ServiceAccount that cannot be impersonated sets Drifted to Unknown").

#### Scenario: Dry-run API failure
- **GIVEN** a ModuleRelease where the API server returns an error other than `Forbidden` during dry-run
- **WHEN** the controller runs Phase 4
- **THEN** `status.failureCounters.drift` is incremented
- **AND** the `Drifted` condition keeps its previous value
- **AND** the reconcile continues to Phase 5 (drift failure is non-blocking)

#### Scenario: Dry-run refused
- **GIVEN** a ModuleInstance where the API server answers `Forbidden` to the dry-run
- **WHEN** the controller runs Phase 4
- **THEN** `status.failureCounters.drift` is incremented
- **AND** `Drifted` is `Unknown` with reason `DriftCheckForbidden`
- **AND** the reconcile continues (drift failure is non-blocking)

## ADDED Requirements

### Requirement: Drift detection runs on every reconcile that renders
Drift detection MUST run on every reconcile that renders, even when digest comparison indicates no-op.

#### Scenario: Drift detected during no-op
- **GIVEN** a ModuleRelease where source, config, and render digests are unchanged
- **AND** a resource has been manually modified on the cluster
- **AND** the reconcile renders (an input changed or the drift render interval passed)
- **WHEN** the controller reconciles
- **THEN** drift is detected and `Drifted=True` is set
- **AND** apply is still skipped (no source/config/render changes)

### Requirement: A reconcile that skips its render runs no drift detection and drift runs at most once per interval
A reconcile that skips its render because its inputs are unchanged (`render-input-key`) has no rendered objects to compare and MUST NOT run drift detection; it leaves the `Drifted` condition and the drift counter as they were.
The drift render interval (`--drift-render-interval`) bounds how long such skips last, so drift detection runs at most once per interval for an object whose inputs do not change, on the first reconcile after the interval.

#### Scenario: A skipped render runs no drift detection
- **GIVEN** a Ready ModuleInstance whose inputs are unchanged and whose last confirming render is younger than the drift render interval
- **AND** a resource has been manually modified on the cluster
- **WHEN** the controller reconciles
- **THEN** no dry-run is sent and the `Drifted` condition is unchanged

#### Scenario: Drift is found once the interval has passed
- **GIVEN** the same instance, reconciled again after the drift render interval
- **WHEN** the controller reconciles
- **THEN** it renders, drift is detected and `Drifted=True` is set

### Requirement: Drift detection skips a resource the reconcile withheld from apply
Drift detection SHALL skip a resource the reconcile withheld from apply, and its difference from live state SHALL NOT set the `Drifted` condition.
Drift reports that the cluster diverged from what the operator asserts. A withheld resource is one the operator is deliberately not asserting, so reporting it as drift would name a difference the operator created on purpose and intends not to close: a condition that never clears, and one that would bury real drift on the same instance behind it.

#### Scenario: A withheld resource does not set Drifted

- **WHEN** a reconcile withholds a resource whose live state differs from the rendered one
- **THEN** the `Drifted` condition is not set by that difference

#### Scenario: Drift returns when the resource stops being withheld

- **WHEN** a previously withheld resource is applied on a later reconcile
- **THEN** it is included in drift detection from that point on

#### Scenario: Real drift is still reported alongside

- **WHEN** an instance has both a withheld resource and another resource that genuinely drifted
- **THEN** the `Drifted` condition is set, reflecting only the genuinely drifted resource

### Requirement: The refusal that withheld a resource carries the signal that drift does not
As in requirement "Drift detection skips a resource the reconcile withheld from apply": The refusal that withheld the resource SHALL carry that signal instead.

#### Scenario: A withheld resource is reported by its refusal

- **WHEN** a reconcile withholds a resource whose live state differs from the rendered one
- **THEN** the refusal that withheld the resource carries that signal instead

### Requirement: The restore step also applies a rendered object that is taken in
As in requirement "A missing object is restored": The same step SHALL apply the rendered objects that are taken in: objects that exist, are not in `status.inventory` and that the apply verdict allows (`reconcile-loop-assembly`, "An allowed object outside the inventory is taken in"), also when nothing is missing. It SHALL apply these two sets and only these.

#### Scenario: A taken-in object is applied when nothing is missing
- **WHEN** a Ready ModuleInstance reconcile renders with unchanged digests, no rendered object is missing, and an existing rendered object outside `status.inventory` is allowed by the apply verdict
- **THEN** the restore step applies that object

### Requirement: The restore step leaves an inventoried object alone and a taken-in object out of the dry-run diff
As in requirement "A missing object is restored": Objects that exist and are in `status.inventory` SHALL NOT be applied by this step, so the `Drifted` condition that the same reconcile computed stays as computed. A taken-in object SHALL be left out of the dry-run diff of the reconcile that takes it in, so `Drifted` does not name an object that the same reconcile applies; from the next render on it is compared like every other inventoried object.

#### Scenario: A taken-in object is left out of the dry-run diff
- **WHEN** a reconcile takes in an object that exists outside `status.inventory`
- **THEN** `Drifted` does not name that object in that reconcile

### Requirement: A reconcile that restores an object is an apply
As in requirement "A missing object is restored": A reconcile that restores an object is an apply: its outcome is `Applied`, it records a history entry, moves `status.lastAppliedAt` and judges health from that moment.
A failed restore SHALL be classified and retried as a failed apply.

#### Scenario: A restore is an apply
- **WHEN** a reconcile restores an object
- **THEN** its outcome is `Applied`, it records a history entry, moves `status.lastAppliedAt` and judges health from that moment

### Requirement: A Job with a TTL is not restored
As in requirement "A missing object is restored": A `batch/v1` Job whose rendered spec sets `ttlSecondsAfterFinished` SHALL NOT be restored: the cluster deletes such a Job after it finished, and to create it again would run it again.

#### Scenario: A rendered Job with a TTL is not restored
- **WHEN** a rendered Job sets `ttlSecondsAfterFinished` and does not exist while every digest is unchanged
- **THEN** the controller does not restore it

### Requirement: An expired Job leaves the inventory
As in requirement "A Job with a TTL is not restored": Such a Job that does not exist while every digest is unchanged is an expired Job, and the controller SHALL treat it as finished: the reconcile SHALL remove it from the entries of `status.inventory`, so that `instance-health` does not read it, and SHALL leave `status.inventory.digest` as the digest of the rendered set, so that the next reconcile with unchanged digests is still a `NoOp`.

#### Scenario: An expired Job is treated as finished
- **WHEN** a rendered Job with a TTL does not exist while every digest is unchanged
- **THEN** the reconcile removes it from the entries of `status.inventory` and leaves `status.inventory.digest` as the digest of the rendered set, so the next reconcile with unchanged digests is still a `NoOp`

### Requirement: The controller records no outcome of a Job, and a changed digest applies the Job again
As in requirement "An expired Job leaves the inventory": The controller records no outcome of a Job. A Job with a TTL that was removed before it ran, and one that failed before the cluster removed it, therefore read as expired too: neither is created again until a digest changes, and the failure is no longer reported on the instance once the Job is gone. A reconcile whose digests changed SHALL apply every rendered object, the Job included, and SHALL record the full rendered inventory.

#### Scenario: A changed digest applies the Job again
- **WHEN** the digests of a ModuleInstance changed and a rendered Job with a TTL does not exist
- **THEN** the reconcile applies every rendered object, the Job included, and records the full rendered inventory

### Requirement: A reconcile that skips its render restores nothing and a missing Job waits less
As in requirement "A missing object is restored": A reconcile that skips its render has no rendered objects and SHALL NOT restore anything; the drift render interval (`--drift-render-interval`) bounds how long a missing object waits. A missing Job waits less: a ModuleInstance reconcile that reads an inventory Job as absent does not skip its render (`render-input-key`), so the Job is restored, or removed from the inventory as expired, on that reconcile.

#### Scenario: A skipped render restores nothing
- **WHEN** a reconcile skips its render while a rendered object is missing
- **THEN** nothing is restored on that reconcile

### Requirement: A failed dry-run leaves the missing set unknown
As in requirement "A missing object is restored": A failed dry-run SHALL leave the missing set unknown, and nothing is restored on that reconcile.

#### Scenario: A failed dry-run restores nothing
- **WHEN** the dry-run of a reconcile fails
- **THEN** the missing set is unknown and nothing is restored on that reconcile

### Requirement: A Forbidden read or dry-run sets Drifted to Unknown
The controller SHALL read each object through the identity of the dry-run before it sends the dry-run. When the API server refuses that read or the dry-run as `Forbidden`, the controller SHALL set the `Drifted` condition to `Unknown` with reason `DriftCheckForbidden` and a message that carries the API server's refusal, which names the identity.

#### Scenario: The ServiceAccount may not dry-run an object

- **GIVEN** a Ready ModuleInstance whose ServiceAccount `deploy-sa` lost the `patch` verb on ConfigMaps
- **WHEN** the controller reconciles and renders with unchanged digests
- **THEN** `Drifted` is `Unknown` with reason `DriftCheckForbidden` and the message names `system:serviceaccount:<namespace>:deploy-sa`
- **AND** `status.failureCounters.drift` is incremented
- **AND** `Ready=True` is preserved and the outcome is `NoOp`

#### Scenario: The ServiceAccount may not read an object

- **GIVEN** a Ready ModuleInstance whose ServiceAccount `deploy-sa` keeps `patch` and lost the `get` verb on ConfigMaps
- **WHEN** the controller reconciles and renders with unchanged digests
- **THEN** `Drifted` is `Unknown` with reason `DriftCheckForbidden` and the message names the refused `get`
- **AND** `Drifted` is not `True`: no verdict is given against an object that was not read

#### Scenario: The ServiceAccount may not create a missing object

- **GIVEN** a Ready ModuleInstance whose ServiceAccount `deploy-sa` lost the `create` verb on ConfigMaps, and whose rendered ConfigMap was deleted
- **WHEN** the controller reconciles and renders with unchanged digests
- **THEN** `Drifted` is `Unknown` with reason `DriftCheckForbidden`
- **AND** the ConfigMap is not created, by any identity

#### Scenario: A restore runs as the ServiceAccount

- **GIVEN** a Ready ModuleInstance with `spec.serviceAccountName=deploy-sa` whose rendered ConfigMap was deleted
- **WHEN** the controller reconciles and renders with unchanged digests
- **THEN** every read, dry-run and apply of the ConfigMap carries the identity of `deploy-sa`
- **AND** the ConfigMap exists again

### Requirement: Forbidden is every 403 answer
As in requirement "A Forbidden read or dry-run sets Drifted to Unknown": `Forbidden` SHALL be every 403 answer, not only a missing RBAC verb of the identity: an admission control that denies the dry-run with 403, and a controller that may not impersonate the ServiceAccount, give the same reason, and the message says which.

#### Scenario: An admission control that denies the dry-run gives DriftCheckForbidden
- **WHEN** an admission control denies the dry-run with 403
- **THEN** `Drifted` is `Unknown` with reason `DriftCheckForbidden` and the message says which

### Requirement: A ServiceAccount that cannot be impersonated sets Drifted to Unknown
When the effective ServiceAccount cannot be impersonated (it does not exist, or the client cannot be built), drift detection SHALL NOT run and the controller SHALL set `Drifted` to `Unknown` with reason `ImpersonationFailed`.

#### Scenario: The ServiceAccount does not exist

- **GIVEN** a Ready ModuleInstance whose ServiceAccount was deleted
- **WHEN** the controller reconciles and renders with unchanged digests
- **THEN** no dry-run is sent, by any identity
- **AND** `Drifted` is `Unknown` with reason `ImpersonationFailed`
- **AND** `status.failureCounters.drift` is incremented and `Ready=True` is preserved

### Requirement: A refused or unidentified drift check counts as failed and restores nothing
As in requirements "A Forbidden read or dry-run sets Drifted to Unknown" and "A ServiceAccount that cannot be impersonated sets Drifted to Unknown": Both cases SHALL count as a failed drift detection, SHALL leave the missing set unknown so that nothing is restored, and SHALL NOT change the `Ready` condition or the outcome of a reconcile with unchanged digests.

#### Scenario: A refused drift check restores nothing
- **WHEN** the dry-run or the read is refused as `Forbidden`
- **THEN** the missing set is unknown, nothing is restored, and the `Ready` condition and the outcome of a reconcile with unchanged digests do not change

### Requirement: A later drift detection replaces the condition and a successful apply removes it
As in requirements "A Forbidden read or dry-run sets Drifted to Unknown" and "A ServiceAccount that cannot be impersonated sets Drifted to Unknown": A later drift detection that succeeds SHALL replace the condition with its verdict, and a successful apply of the rendered set SHALL remove it.

#### Scenario: The refusal ends

- **GIVEN** the same instance after `deploy-sa` was given the `patch` verb again
- **WHEN** the controller reconciles and renders
- **THEN** the `Drifted` condition carries the verdict of the dry-run: removed when nothing drifted, `True` when an object did
- **AND** `status.failureCounters.drift` is zero

### Requirement: An object adopted by another instance is neither drift nor restored
Drift detection MUST compare only the objects the apply verdict allows, without the objects taken in by the same reconcile ("A missing object is restored"). An object the verdict refuses as `adopted-elsewhere` MUST NOT be reported as drifted or as missing, and MUST NOT be restored.
**Source**: enhancement 0012 ("The apply guard runs on every apply, refusing foreign objects outside the inventory and any object adopted by another instance; a per-object adopt annotation is its only override").

#### Scenario: A let-go object is not drift
- **GIVEN** an inventoried ConfigMap whose `opmodel.dev/adopt` annotation names another instance and whose data the other instance changed
- **WHEN** a reconcile renders with unchanged digests
- **THEN** `Drifted` does not name the ConfigMap

#### Scenario: A let-go object is not restored
- **GIVEN** a ConfigMap the instance let go, which the other instance then deleted and created again with its annotation
- **WHEN** the controller reconciles
- **THEN** the controller does not write the ConfigMap

### Requirement: The read the verdict needs is the one read per object before the dry-run
As in requirement "An object adopted by another instance is neither drift nor restored": The read the verdict needs MUST be the one read made per object before the dry-run, so a reconcile that renders sends no more requests per object than before.
**Source**: enhancement 0012 ("The apply guard runs on every apply, refusing foreign objects outside the inventory and any object adopted by another instance; a per-object adopt annotation is its only override").

#### Scenario: One read per object
- **GIVEN** a ModuleInstance with ten rendered objects and unchanged digests
- **WHEN** a reconcile renders
- **THEN** the controller sends one GET of its own per object before the dry-run, as before this requirement

### Requirement: A failed read on unchanged digests is a failed drift check that restores nothing
As in requirement "An object adopted by another instance is neither drift nor restored": A reconcile with unchanged digests whose read of an object fails MUST report that as a failed drift check, as before (`Drifted=Unknown` with reason `DriftCheckForbidden` when the read is Forbidden), MUST restore nothing and MUST let no object go in that reconcile.
**Source**: enhancement 0012 ("The apply guard runs on every apply, refusing foreign objects outside the inventory and any object adopted by another instance; a per-object adopt annotation is its only override").

#### Scenario: A refused read on unchanged digests
- **GIVEN** an effective ServiceAccount that may not get ConfigMaps, and unchanged digests
- **WHEN** the controller reconciles
- **THEN** `Drifted` is `Unknown` with reason `DriftCheckForbidden`, nothing is restored and `status.inventory` keeps its entries
