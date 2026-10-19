## REMOVED Requirements

### Requirement: Render warnings are emitted as events on transition
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "A render with advisory findings emits one RenderWarning event per finding", "RenderWarning events are emitted only when the findings change", "RenderWarning text is authored by the operator and keyed on facts".

### Requirement: Events emitted when claims are kept
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "A ClaimsKept event is emitted when claims are kept", "No ClaimsKept event is emitted when nothing was kept", "The ClaimsKept message states the count, the claims and how to delete one".

### Requirement: Events emitted when a claim recreate is refused
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "A ClaimConflict event is emitted when a claim recreate is refused", "The ClaimConflict event names the claim and carries the Ready text".

### Requirement: Events emitted when objects are left behind
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "A LeftBehind event is emitted when objects are left in the cluster", "The LeftBehind event type depends on why objects were left", "The LeftBehind message states the count and the library messages", "A LeftBehind event is emitted only by a committed run", "A failed prune or cleanup emits no LeftBehind event", "No LeftBehind event for an absent object or a kept claim".

### Requirement: Events emitted when an object is adopted by another instance
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "A Warning event is emitted when the count of adopted objects changes", "An AdoptedElsewhere event is emitted when the count of adopted objects changes", "The AdoptedElsewhere event is emitted when the count of adopted objects differs", "The AdoptedElsewhere event is emitted whether the reconcile applies, restores or does nothing", "The AdoptedElsewhere event states the count and the library messages", "No AdoptedElsewhere event while the count is unchanged".

### Requirement: The event of a stalled render failure stays inside the event note limit
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "A stalled render failure event carries at most 1024 characters", "A stalled render event message that does not fit keeps whole findings", "A stalled render event message without findings is cut and counted".

### Requirement: Events emitted while a deletion waits
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "A DeletionInProgress event is emitted when a deletion wait starts", "A DeletionBlocked event is emitted when a deletion wait becomes blocked".

### Requirement: Events emitted when a deletion is released without confirmation
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "A DeletionUnconfirmed event is emitted when a deletion is released unread", "The DeletionUnconfirmed message states the count, the reason and the doubt".

## ADDED Requirements

### Requirement: A render with advisory findings emits one RenderWarning event per finding
When a render succeeds with advisory findings (catalog skew under the `Warn` policy, unhandled optional traits), the reconciler SHALL emit one Warning event per distinct finding with reason `RenderWarning` and action `Render`. A render with no findings SHALL emit none.

#### Scenario: Skew under Warn is reported once

- **WHEN** a ModuleInstance's module requires a newer catalog build than the platform pins and the policy is `Warn`
- **THEN** the instance renders, reaches `Ready=True`, and one Warning event names the path and both versions

### Requirement: RenderWarning events are emitted only when the findings change
The reconciler SHALL emit the `RenderWarning` events of a render only when the object's set of findings changes between reconciles, not on every reconcile.

#### Scenario: Unchanged warnings do not repeat

- **WHEN** the same instance reconciles again with the same warnings
- **THEN** no new warning event is emitted

### Requirement: RenderWarning text is authored by the operator and keyed on facts
The `RenderWarning` event text SHALL be authored by the operator from the render's advisory diagnostic rows; the render carries no message string to pass through. The transition check for `RenderWarning` events SHALL be keyed on the facts those rows carry (for skew, the path and both versions; for an unhandled trait, the component and the trait), so that rewording an event does not present an unchanged finding set as changed.

#### Scenario: Rewording does not re-emit

- **WHEN** the operator's warning wording changes while an object's advisory rows are unchanged between reconciles
- **THEN** no new warning event is emitted for that object

### Requirement: A ClaimsKept event is emitted when claims are kept
The controller MUST emit one `Normal` event with reason `ClaimsKept` when a prune or a deletion cleanup keeps one or more PersistentVolumeClaims. The `action` field MUST be `Prune` on the stale-prune path and `Delete` on the deletion path. On the deletion path the event MUST be emitted before the finalizer is removed.

#### Scenario: Kept claims on prune
- **GIVEN** a ModuleInstance with `spec.prune=true` whose prune kept PersistentVolumeClaim `media/old-cache`
- **WHEN** the reconcile completes
- **THEN** a `Normal` event with reason `ClaimsKept` and action `Prune` is emitted
- **AND** its message contains `media/old-cache` and the count 1

#### Scenario: Kept claims on deletion
- **GIVEN** a ModuleInstance with `spec.prune=true` being deleted, whose inventory holds PersistentVolumeClaims `media/config` and `media/cache`
- **WHEN** the deletion cleanup completes
- **THEN** a `Normal` event with reason `ClaimsKept` and action `Delete` is emitted, naming both claims

### Requirement: No ClaimsKept event is emitted when nothing was kept
The controller MUST NOT emit a `ClaimsKept` event when nothing was kept: when `spec.prune` is not true, when `spec.dataPolicy` is `Delete`, or when the stale set or the inventory holds no live claim owned by the instance.

#### Scenario: Nothing kept, no event
- **GIVEN** a ModuleInstance with `spec.prune=true` and `spec.dataPolicy=Delete` whose prune deleted a stale claim
- **WHEN** the reconcile completes
- **THEN** no event with reason `ClaimsKept` is emitted

### Requirement: The ClaimsKept message states the count, the claims and how to delete one
The `ClaimsKept` event MUST state the count, name the kept claims as `<namespace>/<name>` (at most ten, and fewer when the names would take the message past 1024 characters; then the number of the rest), say that the claims are no longer tracked, and say how to delete a claim. The event text MUST NOT carry an enhancement reference.

#### Scenario: More than ten kept claims
- **GIVEN** a deletion cleanup that keeps twelve claims
- **WHEN** the event is emitted
- **THEN** its message names ten claims and says that two more were kept

### Requirement: A ClaimConflict event is emitted when a claim recreate is refused
The controller MUST emit one `Warning` event with reason `ClaimConflict` and action `Apply` when an apply refused the forced recreate of a PersistentVolumeClaim. The controller MUST NOT also emit an `ApplyFailed` event for that apply.

#### Scenario: Refused claim recreate
- **GIVEN** a ModuleInstance with `spec.rollout.forceConflicts: true` and no `spec.dataPolicy`, whose render changes an immutable field of the live PersistentVolumeClaim `media/config`
- **WHEN** the reconcile completes
- **THEN** a `Warning` event with reason `ClaimConflict` and action `Apply` is emitted
- **AND** its message contains `media/config` and the field `spec`
- **AND** no `ApplyFailed` event is emitted

### Requirement: The ClaimConflict event names the claim and carries the Ready text
The `ClaimConflict` event MUST name the claim as `<namespace>/<name>` and carry the same text as the `Ready` condition; it therefore names the refused field when the check before the apply refused the claim, and no field when the resource manager's delete guard kept it. The event text MUST NOT carry an enhancement reference.

#### Scenario: Claim kept by the delete guard
- **GIVEN** an apply whose staged apply reached for the delete of the PersistentVolumeClaim `media/config` after its check passed
- **WHEN** the reconcile completes
- **THEN** a `Warning` event with reason `ClaimConflict` and action `Apply` is emitted
- **AND** its message contains `media/config`, names no refused field and does not say that nothing was applied

### Requirement: A LeftBehind event is emitted when objects are left in the cluster
The controller MUST emit one event with reason `LeftBehind` when a prune or a deletion cleanup leaves one or more objects in the cluster because the delete verdict skipped them. The `action` field MUST be `Prune` on the stale-prune path and `Delete` on the deletion path.

#### Scenario: An object of another instance left on prune
- **GIVEN** a ModuleInstance whose prune left ConfigMap `team-a/example` behind as another instance's
- **WHEN** the reconcile completes
- **THEN** a `Warning` event with reason `LeftBehind` and action `Prune` is emitted
- **AND** its message contains the library's message for `ConfigMap/team-a/example`

#### Scenario: Nothing left, no event
- **GIVEN** a prune that deleted every stale object
- **WHEN** the reconcile completes
- **THEN** no event with reason `LeftBehind` is emitted

### Requirement: The LeftBehind event type depends on why objects were left
The `LeftBehind` event type MUST be `Normal` when every object left is a Namespace or a CustomResourceDefinition that OPM never deletes, and `Warning` when at least one object was left because it is not managed by OPM, belongs to another instance or is adopted by another instance.

#### Scenario: Only a Namespace left on deletion
- **GIVEN** a ModuleInstance being deleted whose inventory holds a core Namespace and a Deployment of its own
- **WHEN** the deletion cleanup completes
- **THEN** a `Normal` event with reason `LeftBehind` and action `Delete` is emitted, naming the Namespace

### Requirement: The LeftBehind message states the count and the library messages
The `LeftBehind` event MUST state the count and MUST carry, for each object, the message the library words for the skip, unchanged (at most ten objects, and fewer when the messages would take the event past 1024 characters; then the number of the rest), with the objects left for an ownership reason first. The event text MUST NOT carry an enhancement reference.

#### Scenario: A Namespace and an adopted object left on deletion
- **GIVEN** a ModuleInstance being deleted whose inventory holds a core Namespace and a ConfigMap annotated for another instance
- **WHEN** the deletion cleanup completes
- **THEN** a `Warning` event with reason `LeftBehind` and action `Delete` is emitted, naming both objects, the ConfigMap first

### Requirement: A LeftBehind event is emitted only by a committed run
The `LeftBehind` event MUST be emitted only by a run whose result is committed: on the stale-prune path by a prune that returned no error, and on the deletion path by the first reconcile of the cleanup that reaches a release verdict: the one that removes the finalizer, or the one that starts to wait for deleted objects. A later reconcile of the same wait MUST NOT emit it again.

#### Scenario: A waiting deletion reports what it left once
- **GIVEN** a ModuleInstance being deleted whose inventory holds a core Namespace and a Deployment that takes three rechecks to disappear
- **WHEN** the deletion completes
- **THEN** exactly one event with reason `LeftBehind` and action `Delete` was emitted

### Requirement: A failed prune or cleanup emits no LeftBehind event
A prune or a deletion cleanup that failed for another entry MUST NOT emit the `LeftBehind` event, because its retry judges the same entries again.

#### Scenario: A failed prune reports nothing
- **GIVEN** a prune that left one object behind and failed to delete another
- **WHEN** the reconcile ends
- **THEN** no event with reason `LeftBehind` is emitted
- **AND** the reconcile that later prunes with success emits one

### Requirement: No LeftBehind event for an absent object or a kept claim
The controller MUST NOT emit a `LeftBehind` event for an object that was already absent or for a kept PersistentVolumeClaim, which has its own event.

#### Scenario: Absent object and kept claim
- **WHEN** a prune or a deletion cleanup meets an object that was already absent, or keeps a PersistentVolumeClaim
- **THEN** no event with reason `LeftBehind` is emitted for them

### Requirement: A Warning event is emitted when the count of adopted objects changes
The controller MUST emit one `Warning` event when a reconcile that renders finds a number of rendered objects that the apply verdict refuses as `adopted-elsewhere`, that number is not zero, and it differs from the number the `Ready` message stated before the reconcile. Source: enhancement 0012 ("The apply guard runs on every apply, refusing foreign objects outside the inventory and any object adopted by another instance; a per-object adopt annotation is its only override").

#### Scenario: Nothing adopted elsewhere, no event
- **GIVEN** a reconcile in which the verdict allows every rendered object
- **WHEN** the reconcile completes
- **THEN** no event with reason `AdoptedElsewhere` is emitted

### Requirement: An AdoptedElsewhere event is emitted when the count of adopted objects changes
The controller MUST emit one `Warning` event with reason `AdoptedElsewhere` and action `Apply` when the condition of requirement "A Warning event is emitted when the count of adopted objects changes" holds. Source: enhancement 0012 ("The apply guard runs on every apply, refusing foreign objects outside the inventory and any object adopted by another instance; a per-object adopt annotation is its only override").

#### Scenario: An object is let go
- **GIVEN** an inventoried ConfigMap `team-a/settings` whose `opmodel.dev/adopt` annotation names another instance
- **WHEN** a reconcile renders
- **THEN** one `Warning` event with reason `AdoptedElsewhere` and action `Apply` is emitted, with the library's message for `ConfigMap/team-a/settings`

### Requirement: The AdoptedElsewhere event is emitted when the count of adopted objects differs
The controller MUST emit the `AdoptedElsewhere` event when a reconcile that renders finds a number of rendered objects that the apply verdict refuses as `adopted-elsewhere`, that number is not zero, and it differs from the number the `Ready` message stated before the reconcile. Source: enhancement 0012 ("The apply guard runs on every apply, refusing foreign objects outside the inventory and any object adopted by another instance; a per-object adopt annotation is its only override").

#### Scenario: A changed count emits the event
- **WHEN** a reconcile that renders finds a number of rendered objects that the apply verdict refuses as `adopted-elsewhere`, that number is not zero, and it differs from the number the `Ready` message stated before the reconcile
- **THEN** the controller emits the `AdoptedElsewhere` event

### Requirement: The AdoptedElsewhere event is emitted whether the reconcile applies, restores or does nothing
The controller MUST emit the `AdoptedElsewhere` event when the condition of requirement "The AdoptedElsewhere event is emitted when the count of adopted objects differs" holds. This holds whether the reconcile applies, restores or is a no-op. Source: enhancement 0012 ("The apply guard runs on every apply, refusing foreign objects outside the inventory and any object adopted by another instance; a per-object adopt annotation is its only override").

#### Scenario: AdoptedElsewhere event on a no-op reconcile
- **WHEN** a reconcile that renders is a no-op and finds a nonzero number of objects refused as `adopted-elsewhere` that differs from the number the `Ready` message stated before
- **THEN** one `Warning` event with reason `AdoptedElsewhere` and action `Apply` is emitted

### Requirement: The AdoptedElsewhere event states the count and the library messages
The `AdoptedElsewhere` event MUST state the count and MUST carry the library's message for each object, unchanged, within the limits of the `ApplyRefused` event (at most ten objects, and fewer when the messages would take the event past 1024 characters; then the number of the rest). Source: enhancement 0012 ("The apply guard runs on every apply, refusing foreign objects outside the inventory and any object adopted by another instance; a per-object adopt annotation is its only override").

#### Scenario: A second object is let go
- **GIVEN** the instance above, and a second rendered object that a user annotates for another instance
- **WHEN** a reconcile renders
- **THEN** one `AdoptedElsewhere` event is emitted that states two objects

### Requirement: No AdoptedElsewhere event while the count is unchanged
A reconcile that finds the same number as before, or that skips its render, MUST emit no `AdoptedElsewhere` event: the `Ready` message carries the state between events (`status-conditions`, "The Ready message counts the objects adopted by another instance"). Source: enhancement 0012 ("The apply guard runs on every apply, refusing foreign objects outside the inventory and any object adopted by another instance; a per-object adopt annotation is its only override").

#### Scenario: No event while nothing changes
- **GIVEN** the ConfigMap above, no longer in `status.inventory`, still rendered and still annotated
- **WHEN** a later reconcile renders
- **THEN** no event with reason `AdoptedElsewhere` is emitted and the `Ready` message still counts one object

### Requirement: A stalled render failure event carries at most 1024 characters
The `Warning` event the controller emits when a render or a package load stalls a ModuleInstance or a ModulePackage SHALL carry a message of at most 1024 characters, the longest note events.k8s.io/v1 accepts. When the failure's message fits, the event message SHALL be the message of the `Ready` condition. The `Ready` condition SHALL keep the whole message.

#### Scenario: A message that fits

- **WHEN** a ModuleInstance leaves three required values unset
- **THEN** the event message is the message of the `Ready` condition

### Requirement: A stalled render event message that does not fit keeps whole findings
When the failure's message does not fit in 1024 characters and the failure lists findings, the event message SHALL keep the text in front of the findings and whole findings from the first, and SHALL end with `; and <N> more findings`, where N is the number of findings left out.

#### Scenario: Twenty unset required values

- **WHEN** a ModuleInstance leaves twenty required `#config` values of its module unset
- **THEN** the `Ready` condition message names all twenty as `values.<field>`
- **AND** the `Warning` event with reason `RenderFailed` has a message of at most 1024 characters that names the first values in whole findings and ends with `; and <N> more findings`, and the values named and N add up to twenty

### Requirement: A stalled render event message without findings is cut and counted
Any other stalled render event message that does not fit in 1024 characters SHALL be cut between two characters and end with ` ... (<N> more characters)`.

#### Scenario: A long message without findings

- **WHEN** a render stalls with a message of 3000 characters that lists no findings
- **THEN** the event message is at most 1024 characters, starts as the condition message and ends with ` ... (<N> more characters)`

### Requirement: A DeletionInProgress event is emitted when a deletion wait starts
The controller MUST emit one Normal event with reason `DeletionInProgress` and action `Delete` when a deletion cleanup starts to wait for deleted objects. The event MUST carry the message of the condition. The controller MUST NOT repeat it at a later reconcile while `Ready` already carries that reason. A deletion that is blocked at its first reconcile emits the `DeletionBlocked` event only.

#### Scenario: One event when the wait starts
- **GIVEN** a deleting ModuleInstance whose cleanup deleted a Deployment that still exists
- **WHEN** three reconciles run while it still exists
- **THEN** exactly one Normal event with reason `DeletionInProgress` was emitted

### Requirement: A DeletionBlocked event is emitted when a deletion wait becomes blocked
The controller MUST emit one Warning event with reason `DeletionBlocked` and action `Delete` when the wait for deleted objects becomes blocked. The event MUST carry the message of the condition. The controller MUST NOT repeat it at a later reconcile while `Ready` already carries that reason. A deletion that is blocked at its first reconcile emits the `DeletionBlocked` event only.

#### Scenario: One warning when the wait is blocked
- **GIVEN** a deleting ModuleInstance that has reason `DeletionInProgress` and whose Deployment passes 10 minutes of termination
- **WHEN** the next reconciles run
- **THEN** exactly one Warning event with reason `DeletionBlocked` was emitted, and it names the Deployment

### Requirement: A DeletionUnconfirmed event is emitted when a deletion is released unread
The controller MUST emit one Warning event with reason `DeletionUnconfirmed` and action `Delete` in the reconcile that removes the cleanup finalizer without having read every object: after the ServiceAccount was deleted during the wait, and for an inventory of kept claims only whose identity is missing or failed. It MUST be emitted before the finalizer is removed.

#### Scenario: Released after the ServiceAccount went
- **GIVEN** a ModuleInstance with reason `DeletionInProgress` and three inventory objects, whose ServiceAccount no longer exists
- **WHEN** the controller reconciles it
- **THEN** one Warning event with reason `DeletionUnconfirmed` is emitted that names the count 3 and the missing ServiceAccount

#### Scenario: A confirmed deletion emits no such event
- **GIVEN** a deletion whose every deleted object was read as gone
- **WHEN** the finalizer is removed
- **THEN** no event with reason `DeletionUnconfirmed` is emitted

### Requirement: The DeletionUnconfirmed message states the count, the reason and the doubt
The `DeletionUnconfirmed` message MUST state the number of objects that were not read, MUST say why (the ServiceAccount is missing or cannot be impersonated), and MUST say that the objects may still exist. The event text MUST NOT carry an enhancement reference.

#### Scenario: Message of a released deletion
- **WHEN** a deletion is released without having read every object
- **THEN** the `DeletionUnconfirmed` message states the number of objects that were not read, says why (the ServiceAccount is missing or cannot be impersonated) and says that the objects may still exist
