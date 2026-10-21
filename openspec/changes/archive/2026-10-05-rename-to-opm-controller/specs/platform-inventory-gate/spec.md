## MODIFIED Requirements

### Requirement: An over-subscribed platform is refused at generation

After a generated platform module builds, the reconciler SHALL read the built platform's contract inventory before recording the package. When any provider-fulfilled contract is provided by more than one enabled registry entry (the inventory is not routable; enhancements 0010 and 0015), the reconciler SHALL refuse the package: `Ready=False` with reason `OverSubscribedContracts` and a message naming, in a deterministic order, each over-subscribed contract, the catalog that defines it when an enabled catalog does, and every registry entry that provides it. When the inventory also reports a colliding contract key, the reason SHALL be `ContractCollisions` instead and the message SHALL carry the over-subscription finding after the collision finding.

The count SHALL be the platform's own, as core computes it and the render build refuses on: every transformer of every enabled registry entry counts, a provider is identified by its registry entry (the catalog module path with its major), and a contract counts whether or not an enabled catalog defines it. The controller SHALL NOT recount it, so the generation gate and the render refusal fire on the same platforms.

The over-subscription finding SHALL be printed only when the inventory lists an over-subscribed contract. An inventory that is not routable while listing neither an over-subscribed contract nor a colliding key SHALL still be refused (fail closed), with reason `OverSubscribedContracts` and a message saying the platform is not routable and the inventory names no over-subscribed or colliding contract.

#### Scenario: Two catalogs providing one contract are refused

- **WHEN** the built platform's inventory lists a provider-fulfilled contract provided by each of two enabled registry entries
- **THEN** the Platform reports `Ready=False` with reason `OverSubscribedContracts`, the message names the contract, its defining catalog and both providing registry entries, and no package is recorded for that tuple

#### Scenario: Two majors of one provider catalog are refused

- **WHEN** two majors of one provider catalog are enabled and each carries a transformer requiring the same provider-fulfilled contract
- **THEN** the Platform reports `Ready=False` with reason `OverSubscribedContracts`, the message names both registry entries (the same path with two majors), and no package is recorded for that tuple

#### Scenario: Providers of a contract whose defining catalog is not enabled are refused

- **WHEN** the catalog defining a provider-fulfilled contract is disabled or absent and two enabled registry entries provide that contract
- **THEN** the Platform reports `Ready=False` with reason `OverSubscribedContracts`, the message names the contract and both providing registry entries with no defining catalog, and no package is recorded for that tuple

#### Scenario: An unroutable inventory that names nothing still refuses

- **WHEN** the built platform's inventory reads `routable` false and lists no over-subscribed contract and no colliding key
- **THEN** the Platform reports `Ready=False` with reason `OverSubscribedContracts`, the message says the inventory names no over-subscribed or colliding contract and never reports a count of zero over-subscribed contracts, and no package is recorded for that tuple

### Requirement: A platform whose enabled entries share contract keys is refused at generation

After a generated platform module builds, when the built platform's contract inventory reports any colliding contract key (a key that the catalogs of more than one enabled registry entry list, as two majors of one catalog do), the reconciler SHALL refuse the package: `Ready=False` with reason `ContractCollisions` and a message naming, in a deterministic order, each colliding key and every enabled registry entry (the catalog module path with its major) defining it, and saying that a package cannot be generated until all but one of those entries is disabled.

The colliding keys and their defining entries SHALL be read from the inventory as core reports them; the controller SHALL NOT count definers. The refusal SHALL hold whatever the inventory's other verdicts read: a colliding key is absent from the defined, required and comparable reports, so `fulfilled` and `discriminated` can read true while a collision exists, and neither SHALL be taken as clearing it.

`ContractCollisions` SHALL take precedence over `OverSubscribedContracts` and `ComparablePredicates`. When the inventory also reports an over-subscribed contract or a comparable pair, the reason SHALL be `ContractCollisions` and the message SHALL carry every finding, the collision first.

A platform whose enabled entries define every contract key once SHALL NOT be refused on this ground.

#### Scenario: Two majors sharing contract keys are refused naming the keys and the entries

- **WHEN** two majors of one catalog are enabled and the built platform's inventory reports two contract keys each defined by both registry entries, and no over-subscribed contract
- **THEN** the Platform reports `Ready=False` with reason `ContractCollisions`, the message names both keys in sorted order and, for each, both registry entries in sorted order, it does not mention over-subscribed contracts, and no package is recorded for that tuple

#### Scenario: A collision is reported ahead of an over-subscription

- **WHEN** the built platform's inventory reports a colliding contract key and an over-subscribed provider-fulfilled contract
- **THEN** the Platform reports `Ready=False` with reason `ContractCollisions` and the message names the colliding key with its defining entries first, then the over-subscribed contract with its providing registry entries

#### Scenario: A collision is reported ahead of a comparable pair

- **WHEN** the built platform's inventory reports a colliding contract key and a comparable transformer pair over another key
- **THEN** the Platform reports `Ready=False` with reason `ContractCollisions` and the message carries the collision finding and the comparable-pair finding

#### Scenario: A collision refuses although the fulfilment and discrimination reports read true

- **WHEN** the built platform's inventory reports a colliding contract key while reading `fulfilled` and `discriminated` true
- **THEN** the Platform reports `Ready=False` with reason `ContractCollisions` and no package is recorded

#### Scenario: The collision message does not depend on report order

- **WHEN** two builds report the same colliding keys and defining entries in different orders
- **THEN** both produce the same message, so the warning event fires once
