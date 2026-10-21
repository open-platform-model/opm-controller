## Purpose

The Platform reconciler's read of the built platform's contract inventory: which of core's reports refuse platform-package generation, which only report, and how each surfaces on `Platform.status` (enhancement 0015, "Comparable predicates within one catalog-fulfilled bucket are refused, not arbitrated" and "An unfulfilled contract is reported, never refused; only over-subscription refuses platform-package generation").

## Requirements

### Requirement: The gate is re-evaluated from current state

A refusal SHALL be level-computed from the current tuple: a Platform edit or a claim change wakes the reconciler through its existing watches, and a tuple whose built inventory is routable and discriminated SHALL be recorded and reported `Ready=True` with reason `Generated`, with nothing to clear by hand.

#### Scenario: Removing the competing catalog recovers

- **WHEN** a refused tuple's competing catalog is disabled in `spec.registry` or its claim is removed
- **THEN** the next reconcile records the package and the Platform reports `Ready=True` with reason `Generated`

#### Scenario: Disabling all but one colliding entry recovers

- **WHEN** a tuple refused with `ContractCollisions` is edited so that only one of the registry entries defining each colliding key stays enabled, and the rebuilt inventory is routable and discriminated
- **THEN** the next reconcile records the package and the Platform reports `Ready=True` with reason `Generated`

### Requirement: An unreadable inventory is a build failure

When the built platform carries no readable contract inventory, or lacks a report field the library expects, the reconciler SHALL report `Ready=False` with reason `BuildFailed` naming the missing field, and SHALL NOT record the package.

#### Scenario: A platform without the report is not a pass

- **WHEN** the built platform's inventory cannot be decoded
- **THEN** the Platform reports `Ready=False` with reason `BuildFailed` and the message names the field the read failed on

### Requirement: The reconciler reads the contract inventory of the built platform before recording the package
After a generated platform module builds, the reconciler SHALL read the built platform's contract inventory before recording the package.

#### Scenario: The inventory is read before the package is recorded
- **WHEN** a generated platform module builds
- **THEN** the reconciler reads the built platform's contract inventory before recording the package

### Requirement: An over-subscribed platform is refused with the reason OverSubscribedContracts
When any provider-fulfilled contract is provided by more than one enabled registry entry (the inventory is not routable; enhancement 0010 ("A contract declares where its fulfilment comes from, and a provider-fulfilled contract admits exactly one provider")), the reconciler SHALL refuse the package: `Ready=False` with reason `OverSubscribedContracts`.

#### Scenario: Two catalogs providing one contract are refused

- **WHEN** the built platform's inventory lists a provider-fulfilled contract provided by each of two enabled registry entries
- **THEN** the Platform reports `Ready=False` with reason `OverSubscribedContracts`, the message names the contract, its defining catalog and both providing registry entries, and no package is recorded for that tuple

### Requirement: A contract provided by more than one enabled registry entry refuses the package
When any provider-fulfilled contract is provided by more than one enabled registry entry (enhancement 0015 ("One provider per contract stands; provider classes are rejected in this entry" and "An unfulfilled contract is reported, never refused; only over-subscription refuses platform-package generation")), the reconciler SHALL refuse the package.

#### Scenario: More than one provider refuses the package
- **WHEN** a provider-fulfilled contract is provided by more than one enabled registry entry
- **THEN** the reconciler refuses the package

### Requirement: The over-subscription message names the contracts, the catalogs and the providers
When any provider-fulfilled contract is provided by more than one enabled registry entry, the reconciler SHALL refuse the package: a message naming, in a deterministic order, each over-subscribed contract, the catalog that defines it when an enabled catalog does, and every registry entry that provides it.

#### Scenario: The message names the over-subscribed contract, its defining catalog and its providers
- **WHEN** a provider-fulfilled contract is provided by more than one enabled registry entry
- **THEN** the message names, in a deterministic order, each over-subscribed contract, the catalog that defines it when an enabled catalog does, and every registry entry that provides it

### Requirement: A collision changes the reason of an over-subscribed platform
When any provider-fulfilled contract is provided by more than one enabled registry entry, the reconciler SHALL refuse the package: `Ready=False` with reason `OverSubscribedContracts`. When the inventory also reports a colliding contract key, the reason SHALL be `ContractCollisions` instead and the message SHALL carry the over-subscription finding after the collision finding.

#### Scenario: A colliding key changes the reason
- **WHEN** the inventory also reports a colliding contract key
- **THEN** the reason is `ContractCollisions` and the message carries the over-subscription finding after the collision finding

### Requirement: The over-subscription count is the one core computes
When any provider-fulfilled contract is provided by more than one enabled registry entry, the reconciler SHALL refuse the package. The count SHALL be the platform's own, as core computes it and the render build refuses on: every transformer of every enabled registry entry counts, a provider is identified by its registry entry (the catalog module path with its major), and a contract counts whether or not an enabled catalog defines it.

#### Scenario: Two majors of one provider catalog are refused

- **WHEN** two majors of one provider catalog are enabled and each carries a transformer requiring the same provider-fulfilled contract
- **THEN** the Platform reports `Ready=False` with reason `OverSubscribedContracts`, the message names both registry entries (the same path with two majors), and no package is recorded for that tuple

#### Scenario: Providers of a contract whose defining catalog is not enabled are refused

- **WHEN** the catalog defining a provider-fulfilled contract is disabled or absent and two enabled registry entries provide that contract
- **THEN** the Platform reports `Ready=False` with reason `OverSubscribedContracts`, the message names the contract and both providing registry entries with no defining catalog, and no package is recorded for that tuple

### Requirement: The controller does not recount the over-subscription
When any provider-fulfilled contract is provided by more than one enabled registry entry, the reconciler SHALL refuse the package. As listed in requirement "The over-subscription count is the one core computes": The controller SHALL NOT recount it, so the generation gate and the render refusal fire on the same platforms.

#### Scenario: The controller does not recount
- **WHEN** a provider-fulfilled contract is provided by more than one enabled registry entry
- **THEN** the controller does not recount it, so the generation gate and the render refusal fire on the same platforms

### Requirement: The over-subscription finding is printed only for an over-subscribed contract
The over-subscription finding SHALL be printed only when the inventory lists an over-subscribed contract. An inventory that is not routable while listing neither an over-subscribed contract nor a colliding key SHALL still be refused (fail closed), with reason `OverSubscribedContracts` and a message saying the platform is not routable and the inventory names no over-subscribed or colliding contract.

#### Scenario: An unroutable inventory that names nothing still refuses

- **WHEN** the built platform's inventory reads `routable` false and lists no over-subscribed contract and no colliding key
- **THEN** the Platform reports `Ready=False` with reason `OverSubscribedContracts`, the message says the inventory names no over-subscribed or colliding contract and never reports a count of zero over-subscribed contracts, and no package is recorded for that tuple

#### Scenario: No over-subscribed contract, no finding
- **WHEN** the inventory lists no over-subscribed contract
- **THEN** the over-subscription finding is not printed

### Requirement: A platform with comparable transformer predicates is refused with the reason ComparablePredicates
When the built platform's inventory reports a pair of enabled transformers whose match predicates are comparable over a shared catalog-fulfilled contract (the inventory is not discriminated; enhancement 0015 ("Comparable predicates within one catalog-fulfilled bucket are refused, not arbitrated")), the reconciler SHALL refuse the package: `Ready=False` with reason `ComparablePredicates`.

#### Scenario: Discriminated transformers coexist

- **WHEN** the built platform's inventory reports no comparable pair, however many transformers share a catalog-fulfilled contract
- **THEN** the gate refuses nothing on discrimination and the package is recorded

### Requirement: The comparable-predicates message names the pair and the shared contracts
When the built platform's inventory reports a pair of enabled transformers whose match predicates are comparable over a shared catalog-fulfilled contract, the reconciler SHALL refuse the package: a message naming, for each reported pair, the broader transformer, the narrower transformer and the contracts they share, in a deterministic order.

#### Scenario: A comparable pair is refused naming both transformers

- **WHEN** the built platform's inventory reports one comparable pair over one catalog-fulfilled contract
- **THEN** the Platform reports `Ready=False` with reason `ComparablePredicates`, the message names the broader FQN, the narrower FQN and the shared contract, and no package is recorded for that tuple

### Requirement: No arbitration orders comparable transformers
When the built platform's inventory reports a pair of enabled transformers whose match predicates are comparable over a shared catalog-fulfilled contract, the reconciler SHALL refuse the package. No arbitration, ordering or most-specific-wins SHALL be applied.

#### Scenario: No arbitration of a comparable pair
- **WHEN** the built platform's inventory reports a comparable pair
- **THEN** no arbitration, ordering or most-specific-wins is applied

### Requirement: Over-subscription and comparable predicates report together
When a platform is both over-subscribed and undiscriminated and no contract key collides, the reason SHALL be `OverSubscribedContracts` and the message SHALL carry both findings. When a contract key also collides, the reason SHALL be `ContractCollisions` and the message SHALL carry every finding.

#### Scenario: Both findings report under the routing reason

- **WHEN** the built platform's inventory is over-subscribed and not discriminated, and reports no colliding contract key
- **THEN** the Platform reports `Ready=False` with reason `OverSubscribedContracts` and the message names the over-subscribed contracts and the comparable pairs

### Requirement: A refused generation behaves as a failed build outside the Ready condition
A refused generation SHALL behave as a failed build toward everything outside the Ready condition: the process-local store SHALL keep the last good package, `status.packageIdentity` and `status.registry` SHALL keep describing it, `status.observedGeneration` SHALL be set, a warning event SHALL be emitted on transition only, and the reconcile SHALL requeue on the stalled recheck interval. Renders SHALL keep consuming the last good package.

#### Scenario: The last good package survives a refusal

- **WHEN** a package is held for one tuple and the next tuple is refused
- **THEN** the store still returns the held package and its identity, `status.packageIdentity` and `status.registry` are unchanged, and the Ready condition alone carries the refusal

### Requirement: A refusal with no last good package leaves the store empty
A refusal with no last good package SHALL leave the store empty, so instances wait at `PlatformNotReady` while the cause is on the Platform. The refused module directory MAY remain on disk until the next successful generation prunes it.

#### Scenario: A first-boot refusal holds nothing

- **WHEN** no package has been recorded and the first tuple is refused
- **THEN** the store holds no package, renders report `PlatformNotReady`, and the Platform's Ready condition names the refusal

#### Scenario: The refused module directory stays until the next generation
- **WHEN** a generation is refused
- **THEN** the refused module directory may remain on disk until the next successful generation prunes it

### Requirement: The ContractsFulfilled condition is written on both success paths and never moves Ready
The reconciler SHALL write a `ContractsFulfilled` condition from the effective package's inventory on both success paths (a fresh build and the current-package skip) and SHALL declare it among the conditions it owns. The condition SHALL NOT move `Ready` and SHALL NOT be rewritten by a refused generation, which describes only the package renders consume (enhancement 0015, "An unfulfilled contract is reported, never refused; only over-subscription refuses platform-package generation").

#### Scenario: The skip path keeps the condition current

- **WHEN** a reconcile finds the held package already current for the tuple
- **THEN** `ContractsFulfilled` is written from the held package's inventory

### Requirement: The ContractsFulfilled reasons report unfulfilled, fulfilled and undefined contracts
`ContractsFulfilled=False` with reason `UnfulfilledContracts` SHALL name each provider-fulfilled contract nothing on the platform implements and the catalog that defines it; `ContractsFulfilled=True` with reason `ContractsFulfilled` SHALL report that every defined provider-fulfilled contract has a provider; `ContractsFulfilled=True` with reason `NoContractsDefined` SHALL report that the enabled catalogs list no contract, so nothing was verified.

#### Scenario: An unfulfilled contract is a report beside a generated platform

- **WHEN** the effective package's inventory lists a provider-fulfilled contract that no enabled transformer requires
- **THEN** the Platform reports `Ready=True` with reason `Generated` and `ContractsFulfilled=False` with reason `UnfulfilledContracts`, the message naming the contract and its defining catalog

#### Scenario: Every contract fulfilled

- **WHEN** the effective package's inventory defines contracts and lists none unfulfilled
- **THEN** the Platform reports `ContractsFulfilled=True` with reason `ContractsFulfilled`

#### Scenario: No contract defined is visible, not a pass

- **WHEN** the effective package's inventory defines no contract
- **THEN** the Platform reports `ContractsFulfilled=True` with reason `NoContractsDefined` and a message saying nothing was verified

### Requirement: A platform whose enabled entries share contract keys is refused with the reason ContractCollisions
After a generated platform module builds, when the built platform's contract inventory reports any colliding contract key (a key that the catalogs of more than one enabled registry entry list, as two majors of one catalog do), the reconciler SHALL refuse the package: `Ready=False` with reason `ContractCollisions`.

#### Scenario: Two majors sharing contract keys are refused naming the keys and the entries

- **WHEN** two majors of one catalog are enabled and the built platform's inventory reports two contract keys each defined by both registry entries, and no over-subscribed contract
- **THEN** the Platform reports `Ready=False` with reason `ContractCollisions`, the message names both keys in sorted order and, for each, both registry entries in sorted order, it does not mention over-subscribed contracts, and no package is recorded for that tuple

### Requirement: The collision message names each colliding key and its defining entries
After a generated platform module builds, when the built platform's contract inventory reports any colliding contract key (a key that the catalogs of more than one enabled registry entry list, as two majors of one catalog do), the reconciler SHALL refuse the package: a message naming, in a deterministic order, each colliding key and every enabled registry entry (the catalog module path with its major) defining it.

#### Scenario: The collision message does not depend on report order

- **WHEN** two builds report the same colliding keys and defining entries in different orders
- **THEN** both produce the same message, so the warning event fires once

### Requirement: The collision message says a package cannot be generated until all but one entry is disabled
After a generated platform module builds, when the built platform's contract inventory reports any colliding contract key (a key that the catalogs of more than one enabled registry entry list, as two majors of one catalog do), the reconciler SHALL refuse the package: a message saying that a package cannot be generated until all but one of those entries is disabled.

#### Scenario: The collision message names the way out
- **WHEN** a package cannot be generated because enabled entries share contract keys
- **THEN** the message says that a package cannot be generated until all but one of those entries is disabled

### Requirement: The controller reads colliding keys from the inventory
The colliding keys and their defining entries SHALL be read from the inventory as core reports them; the controller SHALL NOT count definers.

#### Scenario: Definers are read, not counted
- **WHEN** core reports a colliding key and its defining entries
- **THEN** they are read from the inventory as core reports them, and the controller does not count definers

### Requirement: The collision refusal holds whatever the other verdicts read
The refusal SHALL hold whatever the inventory's other verdicts read: a colliding key is absent from the defined, required and comparable reports, so `fulfilled` and `discriminated` can read true while a collision exists, and neither SHALL be taken as clearing it.

#### Scenario: A collision refuses although the fulfilment and discrimination reports read true

- **WHEN** the built platform's inventory reports a colliding contract key while reading `fulfilled` and `discriminated` true
- **THEN** the Platform reports `Ready=False` with reason `ContractCollisions` and no package is recorded

### Requirement: ContractCollisions takes precedence over the other refusals
`ContractCollisions` SHALL take precedence over `OverSubscribedContracts` and `ComparablePredicates`. When the inventory also reports an over-subscribed contract or a comparable pair, the reason SHALL be `ContractCollisions` and the message SHALL carry every finding, the collision first.

#### Scenario: A collision is reported ahead of an over-subscription

- **WHEN** the built platform's inventory reports a colliding contract key and an over-subscribed provider-fulfilled contract
- **THEN** the Platform reports `Ready=False` with reason `ContractCollisions` and the message names the colliding key with its defining entries first, then the over-subscribed contract with its providing registry entries

#### Scenario: A collision is reported ahead of a comparable pair

- **WHEN** the built platform's inventory reports a colliding contract key and a comparable transformer pair over another key
- **THEN** the Platform reports `Ready=False` with reason `ContractCollisions` and the message carries the collision finding and the comparable-pair finding

### Requirement: A platform whose enabled entries define every contract key once is not refused for collisions
After a generated platform module builds, when the built platform's contract inventory reports any colliding contract key (a key that the catalogs of more than one enabled registry entry list, as two majors of one catalog do), the reconciler SHALL refuse the package: `Ready=False` with reason `ContractCollisions`. A platform whose enabled entries define every contract key once SHALL NOT be refused on this ground.

#### Scenario: Every contract key defined once
- **WHEN** the enabled entries define every contract key once
- **THEN** the platform is not refused on this ground
