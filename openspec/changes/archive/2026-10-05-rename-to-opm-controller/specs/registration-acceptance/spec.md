## REMOVED Requirements

### Requirement: An accepted claim keeps its verdict through a transient registry failure

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-021); this requirement is restated under the controller's name.

**Migration**: See the requirement "An accepted claim keeps its verdict while the registry fails transiently" in this capability.

## MODIFIED Requirements

### Requirement: Every claim receives a verdict

The controller SHALL reconcile every `TransformerRegistration`, recording the outcome on its status as `conditions`, `accepted` and `observedGeneration`. A refusal SHALL name what failed and the value that failed it, so the claimant can act on it without reading controller logs. A claim SHALL NOT be partially honoured: acceptance is whole or it is a refusal.

#### Scenario: An accepted claim reports acceptance

- **WHEN** a claim passes every check
- **THEN** `status.accepted` is true and a condition records the acceptance

#### Scenario: A refused claim names its reason

- **WHEN** a claim fails any check
- **THEN** `status.accepted` is false and a condition names the failing check and the offending value

### Requirement: The claimed artifact must be a catalog

Acceptance SHALL acquire the artifact named by `spec.catalog` at `spec.version` and SHALL refuse the claim unless that artifact is a `#Catalog` (enhancement 0015 D10). The refusal SHALL be structural (the kind gate the acquisition applies) rather than a rule this repo maintains. An artifact that cannot be resolved at all SHALL be refused distinguishably from one that resolves and is the wrong kind, since the first is a registry or coordinate problem and the second is an authoring one. One case is not a refusal: an accepted claim whose acquisition fails with a transient registry failure keeps its verdict, as the requirement "An accepted claim keeps its verdict while the registry fails transiently" states.

#### Scenario: A module artifact is refused

- **WHEN** `spec.catalog` names an artifact whose kind is `Module`
- **THEN** the claim is refused, naming the kind found

#### Scenario: An unresolvable coordinate is refused distinguishably

- **WHEN** `spec.catalog` at `spec.version` resolves to nothing in the registry
- **THEN** the claim is refused with a reason distinct from the wrong-kind refusal, naming the coordinate

### Requirement: A contract has one provider

Acceptance SHALL refuse a claim whose `provides` names a contract that an enabled entry of the built platform's registry or another **active** claim already provides, naming the holder and the contract (enhancement 0015 D2, keeping 0010 D37's exactly-one-provider rule). This refusal is distinct from the duplicate-claim refusal: that one is about two claims naming the same catalog, this one is about two providers of the same contract, which can arrive from different catalogs entirely.

Which enabled entries provide a contract SHALL be read from the built platform's own provider count (the contract inventory's providers, which core computes once and the render build refuses on), never recounted by the controller. A provider SHALL be identified by its registry entry: the catalog module path with its major, the same string a claim's `spec.catalog` and a `spec.registry` key carry. Two majors of one catalog are two providers. When the built platform's provider count cannot be read, acceptance SHALL defer the verdict (not refuse), as it does when no platform has been built.

An **inactive** accepted claim SHALL NOT hold a contract against a competitor, because a claim that has never served has no dependents to protect.

A claim whose deletion is **blocked** SHALL keep holding its contracts. Its catalog is still supplying them to the generated platform, so admitting a second provider would not replace it: it would over-subscribe the contract and refuse platform generation for every instance in the cluster. Refusing the newcomer names both ends of the problem instead.

A claim's **own** catalog SHALL NOT hold a contract against it. Once a claim is active its catalog is an enabled entry of the very platform the next reconcile judges it against (enhancement 0015 D13), so counting the claim's own entry would refuse the claim for providing what it exists to provide, deactivate it, drop its catalog from the next generated package and accept it again, an oscillation rather than a verdict. The claim's own entry is the registry entry equal to its `spec.catalog`; another major of the same catalog is not its own entry. Only ANOTHER provider of the contract is a conflict, and the claim's own catalog being excused SHALL NOT excuse a different catalog providing the same contract.

#### Scenario: A claim is refused against an enabled subscription

- **WHEN** a claim provides a contract an enabled registry subscription's catalog already provides
- **THEN** it is refused, naming the contract and the subscribed catalog by its registry entry (path with major)

#### Scenario: A claim is refused against an active claim

- **WHEN** a claim provides a contract another active claim already provides
- **THEN** it is refused, naming the contract and the holding claim

#### Scenario: An inactive accepted claim does not block a competitor

- **WHEN** the holder of a contract is accepted but has never activated
- **THEN** a competing claim for that contract is not refused against it

#### Scenario: An active claim is not refused against its own catalog

- **WHEN** an active claim is re-judged against a platform whose registry now carries its own catalog, with transformers stamped the way core stamps them (a major-free module path)
- **THEN** it is not refused for providing its own contracts

#### Scenario: Another provider still refuses a claim whose catalog is in the registry

- **WHEN** a different catalog provides the same contract as an active claim whose own catalog is in the registry
- **THEN** the competing claim is still refused

#### Scenario: A blocked claim still holds its contracts against another catalog

- **WHEN** a claim from a different catalog providing the same contract is judged while the holder's deletion is blocked
- **THEN** it is refused, naming the blocked claim

#### Scenario: Another major of the claim's own catalog is another provider

- **WHEN** a claim for one major of a catalog provides a contract that an enabled registry entry for a different major of the same catalog already provides
- **THEN** it is refused with reason `ContractSubscribed`, naming the contract and the other major's registry entry

#### Scenario: An unreadable provider count defers the verdict

- **WHEN** the built platform's contract inventory cannot be read
- **THEN** the claim is not refused: it reports `Ready=Unknown` with reason `PlatformNotReady` and is requeued

### Requirement: A token endpoint answer is judged by its status

A registry that uses token authentication has a token endpoint, and an answer from that endpoint SHALL count as an answer from the registry. When the acquisition of the claimed catalog fails because the token endpoint answered with an error status, the controller SHALL judge the claim by that status, as the library classifies it, and SHALL NOT treat the failure as an unreachable registry. This holds where the answer reaches the controller through a failed dependency load, which carries no typed status. The controller SHALL NOT classify by message text itself.

A refusal (401) SHALL refuse the claim, accepted or not, with reason `CatalogUnresolved`. A 5xx answer from the token endpoint is a transient registry failure: an accepted claim SHALL keep its verdict through it. A 429 answer (a rate limit) SHALL keep an accepted claim's verdict in the same way, with the same report and the same backoff: a rate limit passes with nothing changed and says nothing about the claim. The controller SHALL tell a 429 by the status of the library's typed fetch failure. That status does not say whether the registry or its token endpoint answered, so a 429 from either holds the claim.

#### Scenario: A refused token during a dependency load un-accepts

- **WHEN** an accepted claim is reconciled and the catalog acquisition fails while loading a dependency, because the registry's token endpoint answers 401
- **THEN** the failure is a typed fetch failure of kind unauthorized with status 401, and it is not transient
- **AND** the claim is refused with reason `CatalogUnresolved` and `status.accepted` is false

#### Scenario: A rate-limited token endpoint holds the claim

- **WHEN** the same acquisition fails because the token endpoint answers 429
- **THEN** the failure is a typed fetch failure with status 429, which the library does not class as transient
- **AND** an accepted claim keeps `status.accepted`, `status.active` and its `Ready` condition, reports `Reconciling=True` with reason `CatalogUnresolved`, and is retried on the backoff

#### Scenario: A token endpoint that is down holds the claim

- **WHEN** the same acquisition fails because the token endpoint answers 503
- **THEN** the failure is transient, and an accepted claim keeps `status.accepted`, `status.active` and its `Ready` condition

## ADDED Requirements

### Requirement: An accepted claim keeps its verdict while the registry fails transiently

When the acquisition of the claimed catalog fails with a transient registry failure, a claim whose `status.accepted` is true and whose `status.observedGeneration` equals its `metadata.generation` SHALL keep its verdict: `status.accepted`, `status.active`, the `Active` condition and a `Ready` condition that reports the acceptance SHALL stay as they were. A `Ready` condition that is not `True` reports no verdict (an earlier reconcile deferred it, as the first reconcile after a controller restart does); it SHALL become `Unknown` with reason `CatalogUnresolved`, so it names the current cause. A transient registry failure is the one the library classifies as transient (the registry gave no HTTP response, or it answered with a 5xx status), or a 429 answer, which the controller reads from the status of the library's typed fetch failure. The controller SHALL NOT classify by message text.

The controller SHALL report the failure without changing the verdict: the claim SHALL carry `Reconciling=True` with reason `CatalogUnresolved`, and the controller SHALL emit one Warning event that carries the registry error when the claim enters this state. The condition message SHALL NOT change between attempts, so repeated attempts write no status and emit no further event.

The controller SHALL retry the acquisition with an exponential backoff that starts at 5 seconds, doubles, and is capped at 5 minutes, with jitter. The backoff SHALL be measured from the time the claim entered this state, so a controller restart does not reset it. When an attempt gets an answer, the claim SHALL be judged on that answer and SHALL no longer report `Reconciling` with reason `CatalogUnresolved`: an acceptance or a refusal removes the condition, and a deferred verdict replaces its reason.

Every other acquisition failure SHALL refuse the claim, accepted or not: a catalog or version the registry does not hold, refused credentials, a wrong-kind artifact, and a failure the library did not classify. A claim that is not accepted SHALL be refused on a transient registry failure too, because it has no verdict to keep. A claim whose spec was edited after its verdict (`status.observedGeneration` is behind `metadata.generation`) SHALL be refused too: its acceptance is for another spec, and the Platform reads the catalog and the version from the spec, so keeping it would keep a coordinate that was never judged.

#### Scenario: An accepted, active claim survives an unreachable registry

- **WHEN** an accepted and active claim is reconciled while the registry cannot be reached and the catalog is not in the module cache
- **THEN** `status.accepted` and `status.active` stay true, `Ready` stays `True` with reason `Accepted`, and `Active` stays `True`
- **AND** the claim reports `Reconciling=True` with reason `CatalogUnresolved`, a Warning event carries the registry error, and the claim is requeued after about 5 seconds

#### Scenario: A claim deferred after a controller restart stays accepted

- **WHEN** an accepted and active claim reports `Ready=Unknown` because an earlier reconcile found no generated platform, and its next reconcile cannot reach the registry
- **THEN** `status.accepted` and `status.active` stay true and `Active` stays `True`
- **AND** `Ready` is `Unknown` with reason `CatalogUnresolved`

#### Scenario: A 5xx answer is transient

- **WHEN** an accepted claim is reconciled and the registry answers the catalog fetch with a 5xx status
- **THEN** the claim keeps its verdict and is retried

#### Scenario: The retry backs off

- **WHEN** an accepted claim has been in this state for longer than 5 seconds
- **THEN** the next retry waits about as long as the claim has been in the state, and never longer than about 5 minutes

#### Scenario: Repeated attempts do not repeat the event

- **WHEN** an accepted claim fails a second attempt in a row with a transient registry failure
- **THEN** no second Warning event is emitted

#### Scenario: The claim is judged again when the registry answers

- **WHEN** a claim in this state is reconciled and the catalog is acquired
- **THEN** the claim is judged on the acquired catalog and no longer reports `Reconciling` with reason `CatalogUnresolved`

#### Scenario: A catalog the registry does not hold still un-accepts

- **WHEN** an accepted claim is reconciled and the registry answers that it does not hold the catalog at that version
- **THEN** the claim is refused with reason `CatalogUnresolved` and `status.accepted` is false

#### Scenario: An unclassified failure still un-accepts

- **WHEN** an accepted claim is reconciled and the acquisition fails with an error the library did not classify as a registry fetch failure
- **THEN** the claim is refused with reason `CatalogUnresolved` and `status.accepted` is false

#### Scenario: A claim that is not accepted is refused as before

- **WHEN** a claim whose `status.accepted` is false is reconciled while the registry cannot be reached
- **THEN** the claim is refused with reason `CatalogUnresolved`

#### Scenario: An edited claim is refused, not held

- **WHEN** an accepted and active claim's `spec.version` is changed and its next reconcile cannot reach the registry
- **THEN** the claim is refused with reason `CatalogUnresolved`, naming the new version, and `status.accepted` is false
