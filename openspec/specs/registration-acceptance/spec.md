## Purpose

Define the claim-to-verdict path: what acceptance fetches, what it compares it against, which claims it refuses, and what each refusal names. A claim carries no author-trusted data (enhancement 0015, "Every field of the claim is derived or stamped; the only authored fact is the module's catalog dependency version"), so acceptance re-derives every fact it judges rather than reading it off the CR.

This capability decides one transition — whether a claim is accepted. Activation, the finalizer, the shrink refusal and regeneration on the accepted set are separate capabilities.

## Requirements

### Requirement: Every claim receives a verdict

The controller SHALL reconcile every `TransformerRegistration`, recording the outcome on its status as `conditions`, `accepted` and `observedGeneration`. A refusal SHALL name what failed and the value that failed it, so the claimant can act on it without reading controller logs. A claim SHALL NOT be partially honoured: acceptance is whole or it is a refusal.

#### Scenario: An accepted claim reports acceptance

- **WHEN** a claim passes every check
- **THEN** `status.accepted` is true and a condition records the acceptance

#### Scenario: A refused claim names its reason

- **WHEN** a claim fails any check
- **THEN** `status.accepted` is false and a condition names the failing check and the offending value

### Requirement: A claimed artifact that is not a catalog is refused
Acceptance SHALL acquire the artifact named by `spec.catalog` at `spec.version` and SHALL refuse the claim unless that artifact is a `#Catalog` (enhancement 0015, "The claim names a published catalog artifact, and only that; transformers never ship inside module artifacts and never ride the CR"). The refusal SHALL be structural (the kind gate the acquisition applies) rather than a rule this repo maintains.
(Requirement: An accepted claim is not refused for a transient registry failure)

#### Scenario: A module artifact is refused

- **WHEN** `spec.catalog` names an artifact whose kind is `Module`
- **THEN** the claim is refused, naming the kind found

### Requirement: An unresolvable artifact is refused distinguishably from a wrong-kind artifact
As in requirement "A claimed artifact that is not a catalog is refused": An artifact that cannot be resolved at all SHALL be refused distinguishably from one that resolves and is the wrong kind, since the first is a registry or coordinate problem and the second is an authoring one.
(Requirement: An accepted claim is not refused for a transient registry failure)

#### Scenario: An unresolvable coordinate is refused distinguishably

- **WHEN** `spec.catalog` at `spec.version` resolves to nothing in the registry
- **THEN** the claim is refused with a reason distinct from the wrong-kind refusal, naming the coordinate

### Requirement: An accepted claim is not refused for a transient registry failure
As in requirement "A claimed artifact that is not a catalog is refused": One case is not a refusal: an accepted claim whose acquisition fails with a transient registry failure SHALL keep its verdict, as the requirement "A transient registry failure leaves an accepted claim's verdict as it was" states.

#### Scenario: A transient registry failure does not refuse an accepted claim

- **WHEN** an accepted claim's acquisition of the claimed catalog fails with a transient registry failure
- **THEN** the claim is not refused and keeps its verdict

### Requirement: The derived provider set must equal spec.provides
Acceptance SHALL derive the provider-fulfilled contract set from the acquired catalog and SHALL compare it to `spec.provides` for exact equality, refusing drift in either direction and naming both lists (enhancement 0015, "Every field of the claim is derived or stamped; the only authored fact is the module's catalog dependency version").

#### Scenario: An exactly matching claim passes this check

- **WHEN** `spec.provides` equals the derived set
- **THEN** this check does not refuse the claim, whatever the order the two were produced in

### Requirement: A claim that lists a contract the catalog does not implement or omits one it does is refused
As in requirement "The derived provider set must equal spec.provides": A claim listing a contract the catalog does not implement, and a claim omitting one it does, SHALL both be refused. Acceptance SHALL NOT accept a subset, because partial registration would leave the remainder reported as unfulfilled with no indication that the provider withheld it.

#### Scenario: A claim listing an unimplemented contract is refused

- **WHEN** `spec.provides` names a contract the catalog's transformers do not require with provider fulfilment
- **THEN** the claim is refused, naming both the claimed and the derived list

#### Scenario: A claim omitting an implemented contract is refused

- **WHEN** the catalog implements a provider-fulfilled contract absent from `spec.provides`
- **THEN** the claim is refused, naming both lists

### Requirement: A stray claim is refused
Acceptance SHALL verify that the `ModuleInstance` named by `spec.providerRef` exists and that its `status.inventory` owns this claim, refusing the claim otherwise (the claim-origin check of enhancement 0015, deferred to this slice). `providerRef` is stamped by the renderer and cannot be authored, so a claim its named instance does not own did not come from that instance: a hand-applied stray, which would otherwise point a later health gate at the wrong package.

#### Scenario: A claim its named instance does not own is refused

- **WHEN** the `ModuleInstance` named by `spec.providerRef` has a settled inventory that does not hold this claim
- **THEN** the claim is refused, naming both identities

#### Scenario: A claim naming an instance that does not exist is refused

- **WHEN** `spec.providerRef` names a `ModuleInstance` that is absent
- **THEN** the claim is refused as not rendered output, rather than accepted on the strength of its own `providerRef`

### Requirement: The origin check is grounded in the inventory, not in the claim's labels
As in requirement "A stray claim is refused": The check SHALL be grounded in the inventory rather than in the claim's labels. D11 suggests the owner labels, but measured, a rendered claim carries an instance *name* label and no namespace-bearing or uuid label, while `providerRef` carries a namespace and a name; a label-only check would therefore admit a stray placed by any instance sharing the provider's name in another namespace: the exact substitution the check exists to stop.

#### Scenario: A stray claim from an instance of the same name in another namespace is refused

- **WHEN** a claim is placed by an instance that shares the provider's name in another namespace
- **THEN** the claim is refused, because the inventory of the instance that `providerRef` names does not hold it

### Requirement: The inventory is the ownership record of this repo
As in requirement "A stray claim is refused": The check SHALL be grounded in the inventory rather than in the claim's labels. The inventory is already this repo's ownership record by constitutional principle.

#### Scenario: The origin check reads ownership from the inventory

- **WHEN** the origin of a claim is checked
- **THEN** ownership is read from the `status.inventory` of the instance that `spec.providerRef` names

### Requirement: An instance whose inventory has not settled leaves the claim awaiting a verdict
As in requirement "A stray claim is refused": An instance whose inventory has not settled SHALL leave the claim awaiting a verdict rather than refused: a claim can reach the API server before its owner's status does, and a race is not a verdict.

#### Scenario: A claim racing its provider's inventory is not refused

- **WHEN** the `ModuleInstance` named by `spec.providerRef` exists but has written no inventory yet
- **THEN** no verdict is recorded and the claim is retried

### Requirement: A second claim for a provider catalog is refused naming the claimant
When more than one claim names the same provider catalog, acceptance SHALL accept at most one and SHALL refuse the others naming the claimant that holds it (enhancement 0015, "The rendered CR's name is instance-derived: dot-joined `namespace.name`"). The refusal SHALL identify the competing claim by name, so an operator can see which object to remove. The decision SHALL be stable: which claim holds acceptance SHALL NOT change from one reconcile to the next while both exist.

#### Scenario: The second claim is refused naming the first

- **WHEN** two instances of one provider module each render a claim for the same catalog
- **THEN** one is accepted and the other is refused, naming the accepted claim

#### Scenario: The holder does not change under repeated reconciles

- **WHEN** both claims are reconciled repeatedly with no spec change
- **THEN** the same claim keeps acceptance

### Requirement: A claim whose deletion is blocked still holds its provider catalog
As in requirement "A second claim for a provider catalog is refused naming the claimant": A claim being deleted SHALL NOT compete for the provider, **unless its deletion is blocked**, in which case it SHALL keep holding it. A blocked claim has not left: it is still accepted, still active, and its catalog is still an entry of the generated platform, so a second claim for that catalog is the duplicate this requirement refuses.

#### Scenario: A claim whose deletion is blocked keeps holding its provider

- **WHEN** a claim for the same catalog is judged while the holder's deletion is blocked by its dependents
- **THEN** the new claim is refused, naming the blocked claim

### Requirement: A claim is refused when a contract it provides has another provider
Acceptance SHALL refuse a claim whose `provides` names a contract that an enabled entry of the built platform's registry or another **active** claim already provides, naming the holder and the contract (enhancement 0015 ("One provider per contract stands; provider classes are rejected in this entry"), keeping the exactly-one-provider rule of enhancement 0010).

#### Scenario: A claim is refused against an enabled subscription

- **WHEN** a claim provides a contract an enabled registry subscription's catalog already provides
- **THEN** it is refused, naming the contract and the subscribed catalog by its registry entry (path with major)

#### Scenario: A claim is refused against an active claim

- **WHEN** a claim provides a contract another active claim already provides
- **THEN** it is refused, naming the contract and the holding claim

### Requirement: A contract refusal is distinct from the duplicate-claim refusal
As in requirement "A claim is refused when a contract it provides has another provider": This refusal SHALL be distinct from the duplicate-claim refusal: that one is about two claims naming the same catalog, this one is about two providers of the same contract, which can arrive from different catalogs entirely.

#### Scenario: Two providers of one contract from different catalogs are a contract conflict

- **WHEN** a claim provides a contract that an active claim from a different catalog already provides
- **THEN** the claim is refused for the contract, naming the holder and the contract, and not as a duplicate claim

### Requirement: The provider count of a contract is read from the built platform
As in requirement "A claim is refused when a contract it provides has another provider": Which enabled entries provide a contract SHALL be read from the built platform's own provider count (the contract inventory's providers, which core computes once and the render build refuses on), never recounted by the controller. When the built platform's provider count cannot be read, acceptance SHALL defer the verdict (not refuse), as it does when no platform has been built.

#### Scenario: An unreadable provider count defers the verdict

- **WHEN** the built platform's contract inventory cannot be read
- **THEN** the claim is not refused: it reports `Ready=Unknown` with reason `PlatformNotReady` and is requeued

### Requirement: A provider is identified by its registry entry
As in requirement "A claim is refused when a contract it provides has another provider": A provider SHALL be identified by its registry entry: the catalog module path with its major, the same string a claim's `spec.catalog` and a `spec.registry` key carry. Two majors of one catalog SHALL be two providers.

#### Scenario: Another major of the claim's own catalog is another provider

- **WHEN** a claim for one major of a catalog provides a contract that an enabled registry entry for a different major of the same catalog already provides
- **THEN** it is refused with reason `ContractSubscribed`, naming the contract and the other major's registry entry

### Requirement: An inactive accepted claim does not hold a contract against a competitor
An **inactive** accepted claim SHALL NOT hold a contract against a competitor, because a claim that has never served has no dependents to protect.

#### Scenario: An inactive accepted claim does not block a competitor

- **WHEN** the holder of a contract is accepted but has never activated
- **THEN** a competing claim for that contract is not refused against it

### Requirement: A claim whose deletion is blocked keeps holding its contracts
As in requirement "A claim is refused when a contract it provides has another provider": A claim whose deletion is **blocked** SHALL keep holding its contracts. Its catalog is still supplying them to the generated platform, so admitting a second provider would not replace it: it would over-subscribe the contract and refuse platform generation for every instance in the cluster. Refusing the newcomer names both ends of the problem instead.

#### Scenario: A blocked claim still holds its contracts against another catalog

- **WHEN** a claim from a different catalog providing the same contract is judged while the holder's deletion is blocked
- **THEN** it is refused, naming the blocked claim

### Requirement: A claim's own catalog does not hold a contract against it
A claim's **own** catalog SHALL NOT hold a contract against it. Once a claim is active its catalog is an enabled entry of the very platform the next reconcile judges it against (enhancement 0015, "The platform package regenerates edge-triggered and level-computed, keyed by generation plus the active-claim set; blast radius accepted"). The claim's own entry is the registry entry equal to its `spec.catalog`; another major of the same catalog is not its own entry.

#### Scenario: An active claim is not refused against its own catalog

- **WHEN** an active claim is re-judged against a platform whose registry now carries its own catalog, with transformers stamped the way core stamps them (a major-free module path)
- **THEN** it is not refused for providing its own contracts

### Requirement: Only another provider of a contract conflicts with a claim
As in requirement "A claim's own catalog does not hold a contract against it": So counting the claim's own entry would refuse the claim for providing what it exists to provide, deactivate it, drop its catalog from the next generated package and accept it again, an oscillation rather than a verdict. Only ANOTHER provider of the contract is a conflict, and the claim's own catalog being excused SHALL NOT excuse a different catalog providing the same contract.

#### Scenario: Another provider still refuses a claim whose catalog is in the registry

- **WHEN** a different catalog provides the same contract as an active claim whose own catalog is in the registry
- **THEN** the competing claim is still refused

#### Scenario: An active claim's own entry is not counted against it

- **WHEN** an active claim's catalog is an enabled entry of the platform that the next reconcile judges the claim against
- **THEN** the claim's own entry is not counted, so the claim is not refused, deactivated and accepted again

### Requirement: A newer-build requirement is refused
Acceptance SHALL compare the catalog's committed requirements against the platform's resolved versions, per shared OPM-namespace path, and SHALL refuse the claim when the catalog requires a version greater than the platform's within the same major, or requires a major the platform does not resolve for that path (enhancement 0015, "A build-incompatible registration is refused at acceptance, by comparing committed resolutions under the stability discipline").
**Source**: enhancement 0015 ("A build-incompatible registration is refused at acceptance, by comparing committed resolutions under the stability discipline") and enhancement 0026 ("The shared-path requirement check runs per render against the consumer's pins, and at acceptance against the platform floors" and "A platform may admit several majors of one catalog; each resolution holds exactly one major per catalog") (the shared-path comparison).

#### Scenario: A provider requiring a newer build is refused

- **WHEN** the claimed catalog requires a shared path at a version greater than the platform's resolved version, same major
- **THEN** the claim is refused, naming the path, both versions, and the one-line fix

#### Scenario: A provider requiring an older build is accepted by this check

- **WHEN** every shared requirement is at or below the platform's resolved version within the same major
- **THEN** this check does not refuse the claim

### Requirement: A requirement is compared against the resolved entry of the provider's own major
As in requirement "A newer-build requirement is refused": Two majors of one catalog are two paths: the platform's resolution is every module its generated closure carries, disabled registry entries included, and a provider's requirement SHALL be compared against the resolved entry of the provider's own major. A requirement in a different major SHALL be refused only when no resolved path shares the provider's major, and that refusal SHALL name every major the platform resolves for the path.
**Source**: enhancement 0015 ("A build-incompatible registration is refused at acceptance, by comparing committed resolutions under the stability discipline") and enhancement 0026 ("The shared-path requirement check runs per render against the consumer's pins, and at acceptance against the platform floors" and "A platform may admit several majors of one catalog; each resolution holds exactly one major per catalog") (the shared-path comparison).

#### Scenario: A provider on a different major is refused unconditionally

- **WHEN** the claimed catalog requires a shared path in a different major than every major the platform resolves for that path
- **THEN** the claim is refused without comparing versions, naming every major the platform resolves for that path

#### Scenario: A provider is compared against the resolved entry of its own major

- **WHEN** the platform resolves `opmodel.dev/catalogs/opm@v4` at `v4.2.0` and `opmodel.dev/catalogs/opm@v5` at `v5.0.0` (for example an enabled v4 entry beside a disabled v5 entry), and the claimed catalog requires `opmodel.dev/catalogs/opm@v4` at `v4.1.0`
- **THEN** this check does not refuse the claim, and a catalog requiring `opmodel.dev/catalogs/opm@v5` at `v5.1.0` is refused as requiring a newer build than `v5.0.0`, never as a major mismatch

### Requirement: The build-compatibility verdict does not depend on the order of the platform's resolution
As in requirement "A newer-build requirement is refused": The verdict and its message SHALL NOT depend on the order in which the platform's resolution is read.
**Source**: enhancement 0015 ("A build-incompatible registration is refused at acceptance, by comparing committed resolutions under the stability discipline") and enhancement 0026 ("The shared-path requirement check runs per render against the consumer's pins, and at acceptance against the platform floors" and "A platform may admit several majors of one catalog; each resolution holds exactly one major per catalog") (the shared-path comparison).

#### Scenario: The verdict does not depend on iteration order

- **WHEN** the same claim is judged repeatedly against a platform resolving several majors of a path the claimed catalog requires
- **THEN** every judgement reaches the same verdict with a byte-identical message

### Requirement: The build-incompatibility refusal happens at acceptance and states the fix
As in requirement "A newer-build requirement is refused": The refusal SHALL happen here rather than at render, where the failure would name an unrelated module instance. The refusal message SHALL state that the comparison is conservative (a requirement records what the provider was tidied against, not what it uses) and that lowering the requirement is the author's fix.
**Source**: enhancement 0015 ("A build-incompatible registration is refused at acceptance, by comparing committed resolutions under the stability discipline") and enhancement 0026 ("The shared-path requirement check runs per render against the consumer's pins, and at acceptance against the platform floors" and "A platform may admit several majors of one catalog; each resolution holds exactly one major per catalog") (the shared-path comparison).

#### Scenario: A refusal for a build incompatibility states the conservative comparison and the fix

- **WHEN** a claim is refused because its catalog requires a newer build than the platform resolves
- **THEN** the refusal happens at acceptance, not at render, and its message states that the comparison is conservative and that lowering the requirement is the author's fix

### Requirement: A claim gets the same verdict on both paths of the provider-set derivation
Acceptance SHALL reach the same verdict for a claim whether the library derives the claimed catalog's provider set from the `provides` field core computes, which a catalog built against core `v2.0.0-beta.3` or later carries, or from the deprecated fold over `#transformers`, which answers a catalog built against an older core.

#### Scenario: A claim on a catalog built against an older core is accepted through the fold

- **WHEN** the claim rendered by `backup_provider` names the published `backup` catalog fixture, whose committed core pin predates `ProvidesSince`, and its provider instance's inventory owns the claim
- **THEN** the catalog carries no `provides` field and its derived provider set is exactly `opmodel.dev/catalogs/opm/traits/backup@v1alpha1`
- **AND** the reconciler records `accepted: true` and `Ready=True` with reason `Accepted`, so the exact-equality check did not refuse it

#### Scenario: A claim on a catalog built against a newer core is accepted through the decoded field

- **WHEN** the same claim is judged against the `backup` fixture tree with its core pin moved to exactly `ProvidesSince`
- **THEN** the catalog carries a `provides` field equal to its derived provider set, which is exactly `opmodel.dev/catalogs/opm/traits/backup@v1alpha1`
- **AND** the reconciler records `accepted: true` and `Ready=True` with reason `Accepted`

### Requirement: A registry-backed integration spec runs the reconciler once with each kind of catalog
As in requirement "A claim gets the same verdict on both paths of the provider-set derivation": A registry-backed integration spec SHALL run the `TransformerRegistrationReconciler` on the claim `backup_provider` renders, against a generated platform, once with each kind of catalog. Each time it SHALL assert that the claim is accepted and that the catalog's derived provider set is exactly opm's backup trait.

#### Scenario: The integration spec runs the reconciler on the backup claim for each kind of catalog

- **WHEN** the integration spec runs for one kind of catalog
- **THEN** it runs the `TransformerRegistrationReconciler` on the claim `backup_provider` renders, against a generated platform, and asserts that the claim is accepted and that the catalog's derived provider set is exactly opm's backup trait

### Requirement: The old catalog is the published backup fixture and the new catalog moves its core pin
As in requirement "A claim gets the same verdict on both paths of the provider-set derivation": The old catalog SHALL be the published `backup` catalog fixture, acquired from the registry. The new catalog SHALL be the same fixture tree with its core pin moved to exactly the library's `ProvidesSince`, the first core that derives `provides`.

#### Scenario: The two integration catalogs come from one fixture

- **WHEN** the old-catalog spec and the new-catalog spec acquire their catalogs
- **THEN** the old catalog is the published `backup` catalog fixture acquired from the registry, and the new catalog is the same fixture tree with its core pin moved to exactly the library's `ProvidesSince`

### Requirement: Each integration spec asserts the facts that select its path
As in requirement "A registry-backed integration spec runs the reconciler once with each kind of catalog": Each spec SHALL also assert the facts that select its path: the old catalog pins a core older than `ProvidesSince` and carries no `provides` field, and the new catalog pins exactly `ProvidesSince` and carries a `provides` field that equals the derived set. The specs cannot tell the fold from the field on the new catalog, because both give the same set there.

#### Scenario: The old-catalog spec asserts the facts of the fold path

- **WHEN** the old-catalog spec runs
- **THEN** it asserts that the old catalog pins a core older than `ProvidesSince` and carries no `provides` field

### Requirement: A broken premise fails the spec that relies on it
As in requirement "Each integration spec asserts the facts that select its path": The old catalog SHALL be the published `backup` catalog fixture, acquired from the registry. A fixture or library move that breaks a premise then fails the spec that relies on it, and a library that stops folding a catalog without the field fails the old spec; one that picks its path by the field's presence, ignoring the pin, passes it, since the old catalog has no field.

#### Scenario: The old-catalog case is not lost silently

- **WHEN** the `backup` catalog fixture's core pin is moved to `ProvidesSince` or later, or the library's `ProvidesSince` names a core that does not derive `provides`
- **THEN** the spec whose path premise no longer holds fails, naming that premise, instead of passing on the other path

### Requirement: A transient registry failure leaves an accepted claim's verdict as it was
When the acquisition of the claimed catalog fails with a transient registry failure, a claim whose `status.accepted` is true and whose `status.observedGeneration` equals its `metadata.generation` SHALL keep its verdict: `status.accepted`, `status.active`, the `Active` condition and a `Ready` condition that reports the acceptance SHALL stay as they were.

#### Scenario: An accepted, active claim survives an unreachable registry

- **WHEN** an accepted and active claim is reconciled while the registry cannot be reached and the catalog is not in the module cache
- **THEN** `status.accepted` and `status.active` stay true, `Ready` stays `True` with reason `Accepted`, and `Active` stays `True`
- **AND** the claim reports `Reconciling=True` with reason `CatalogUnresolved`, a Warning event carries the registry error, and the claim is requeued after about 5 seconds

### Requirement: A Ready condition that reports no verdict becomes Unknown with reason CatalogUnresolved
As in requirement "A transient registry failure leaves an accepted claim's verdict as it was": A `Ready` condition that is not `True` reports no verdict (an earlier reconcile deferred it, as the first reconcile after a controller restart does); it SHALL become `Unknown` with reason `CatalogUnresolved`, so it names the current cause.

#### Scenario: A claim deferred after a controller restart stays accepted

- **WHEN** an accepted and active claim reports `Ready=Unknown` because an earlier reconcile found no generated platform, and its next reconcile cannot reach the registry
- **THEN** `status.accepted` and `status.active` stay true and `Active` stays `True`
- **AND** `Ready` is `Unknown` with reason `CatalogUnresolved`

### Requirement: A transient registry failure is one the library classifies as transient or a 429 answer
A transient registry failure SHALL be the one the library classifies as transient (the registry gave no HTTP response, or it answered with a 5xx status), or a 429 answer, which the controller reads from the status of the library's typed fetch failure. The controller SHALL NOT classify by message text.

#### Scenario: A 5xx answer is transient

- **WHEN** an accepted claim is reconciled and the registry answers the catalog fetch with a 5xx status
- **THEN** the claim keeps its verdict and is retried

### Requirement: The controller reports a transient failure without changing the verdict
As in requirement "A transient registry failure leaves an accepted claim's verdict as it was": The controller SHALL report the failure without changing the verdict: the claim SHALL carry `Reconciling=True` with reason `CatalogUnresolved`, and the controller SHALL emit one Warning event that carries the registry error when the claim enters this state. The condition message SHALL NOT change between attempts, so repeated attempts write no status and emit no further event.

#### Scenario: Repeated attempts do not repeat the event

- **WHEN** an accepted claim fails a second attempt in a row with a transient registry failure
- **THEN** no second Warning event is emitted

### Requirement: The retry of the acquisition backs off from the time the claim entered the state
As in requirement "A transient registry failure leaves an accepted claim's verdict as it was": The controller SHALL retry the acquisition with an exponential backoff that starts at 5 seconds, doubles, and is capped at 5 minutes, with jitter. The backoff SHALL be measured from the time the claim entered this state, so a controller restart does not reset it.

#### Scenario: The retry backs off

- **WHEN** an accepted claim has been in this state for longer than 5 seconds
- **THEN** the next retry waits about as long as the claim has been in the state, and never longer than about 5 minutes

### Requirement: A claim is judged again when an attempt gets an answer
As in requirement "A transient registry failure leaves an accepted claim's verdict as it was": When an attempt gets an answer, the claim SHALL be judged on that answer and SHALL no longer report `Reconciling` with reason `CatalogUnresolved`: an acceptance or a refusal removes the condition, and a deferred verdict replaces its reason.

#### Scenario: The claim is judged again when the registry answers

- **WHEN** a claim in this state is reconciled and the catalog is acquired
- **THEN** the claim is judged on the acquired catalog and no longer reports `Reconciling` with reason `CatalogUnresolved`

### Requirement: Every other acquisition failure refuses the claim
As in requirement "A transient registry failure leaves an accepted claim's verdict as it was": Every other acquisition failure SHALL refuse the claim, accepted or not: a catalog or version the registry does not hold, refused credentials, a wrong-kind artifact, and a failure the library did not classify. A claim that is not accepted SHALL be refused on a transient registry failure too, because it has no verdict to keep.

#### Scenario: A catalog the registry does not hold still un-accepts

- **WHEN** an accepted claim is reconciled and the registry answers that it does not hold the catalog at that version
- **THEN** the claim is refused with reason `CatalogUnresolved` and `status.accepted` is false

#### Scenario: An unclassified failure still un-accepts

- **WHEN** an accepted claim is reconciled and the acquisition fails with an error the library did not classify as a registry fetch failure
- **THEN** the claim is refused with reason `CatalogUnresolved` and `status.accepted` is false

#### Scenario: A claim that is not accepted is refused as before

- **WHEN** a claim whose `status.accepted` is false is reconciled while the registry cannot be reached
- **THEN** the claim is refused with reason `CatalogUnresolved`

### Requirement: An edited claim is refused, not held
As in requirement "A transient registry failure leaves an accepted claim's verdict as it was": A claim whose spec was edited after its verdict (`status.observedGeneration` is behind `metadata.generation`) SHALL be refused too: its acceptance is for another spec, and the Platform reads the catalog and the version from the spec, so keeping it would keep a coordinate that was never judged.

#### Scenario: An edited claim is refused, not held

- **WHEN** an accepted and active claim's `spec.version` is changed and its next reconcile cannot reach the registry
- **THEN** the claim is refused with reason `CatalogUnresolved`, naming the new version, and `status.accepted` is false

### Requirement: A token endpoint answer counts as an answer of the registry
A registry that uses token authentication has a token endpoint, and an answer from that endpoint SHALL count as an answer from the registry. When the acquisition of the claimed catalog fails because the token endpoint answered with an error status, the controller SHALL judge the claim by that status, as the library classifies it, and SHALL NOT treat the failure as an unreachable registry.

#### Scenario: A token endpoint error status is judged as the library classifies it

- **WHEN** the acquisition of the claimed catalog fails because the token endpoint answered with an error status
- **THEN** the controller judges the claim by that status, as the library classifies it, and does not treat the failure as an unreachable registry

### Requirement: A token endpoint answer that reaches the controller through a dependency load is judged by its status
As in requirement "A token endpoint answer counts as an answer of the registry": This holds where the answer reaches the controller through a failed dependency load, which carries no typed status. The controller SHALL NOT classify by message text itself.

#### Scenario: A refused token during a dependency load un-accepts

- **WHEN** an accepted claim is reconciled and the catalog acquisition fails while loading a dependency, because the registry's token endpoint answers 401
- **THEN** the failure is a typed fetch failure of kind unauthorized with status 401, and it is not transient
- **AND** the claim is refused with reason `CatalogUnresolved` and `status.accepted` is false

### Requirement: A refusal from the token endpoint refuses the claim
As in requirement "A token endpoint answer counts as an answer of the registry": A refusal (401) SHALL refuse the claim, accepted or not, with reason `CatalogUnresolved`.

#### Scenario: A claim that is not accepted is refused on a 401 from the token endpoint

- **WHEN** a claim that is not accepted is reconciled and the token endpoint answers 401
- **THEN** the claim is refused with reason `CatalogUnresolved`

### Requirement: A 5xx answer from the token endpoint is a transient registry failure
A 5xx answer from the token endpoint is a transient registry failure: an accepted claim SHALL keep its verdict through it.

#### Scenario: A token endpoint that is down holds the claim

- **WHEN** the same acquisition fails because the token endpoint answers 503
- **THEN** the failure is transient, and an accepted claim keeps `status.accepted`, `status.active` and its `Ready` condition

### Requirement: A 429 answer holds an accepted claim's verdict
As in requirement "A 5xx answer from the token endpoint is a transient registry failure": A 429 answer (a rate limit) SHALL keep an accepted claim's verdict in the same way, with the same report and the same backoff: a rate limit passes with nothing changed and says nothing about the claim. The controller SHALL tell a 429 by the status of the library's typed fetch failure. That status does not say whether the registry or its token endpoint answered, so a 429 from either holds the claim.

#### Scenario: A rate-limited token endpoint holds the claim

- **WHEN** the same acquisition fails because the token endpoint answers 429
- **THEN** the failure is a typed fetch failure with status 429, which the library does not class as transient
- **AND** an accepted claim keeps `status.accepted`, `status.active` and its `Ready` condition, reports `Reconciling=True` with reason `CatalogUnresolved`, and is retried on the backoff
