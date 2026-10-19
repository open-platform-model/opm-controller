## REMOVED Requirements

### Requirement: The generated package is a function of the platform spec and the active-claim set
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "The platform package is generated from the Platform spec and the accepted and active claims", "Generation computes from the current state, never from the waking event", "A claim whose deletion is blocked remains in the active set".

### Requirement: The effective registry is readable from the Platform
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "The Platform status carries the union of the authored subscriptions and the active claims' catalogs with the source of each entry", "The Platform status shows the pinned version and whether the transformers of an entry register".


## ADDED Requirements

### Requirement: The platform package is generated from the Platform spec and the accepted and active claims
The operator SHALL generate the platform package from exactly two inputs: the `Platform` CR's spec, and the set of claims that are both accepted and active. A platform with no active claims SHALL generate the package it would have generated from its spec alone.

#### Scenario: An active claim's catalog reaches the generated package

- **WHEN** a claim becomes active
- **THEN** the regenerated package carries that provider catalog as a registry entry at the claim's version

#### Scenario: No active claims leaves generation unchanged

- **WHEN** no claim is active
- **THEN** the generated package is a function of the `Platform` spec alone

### Requirement: Generation computes from the current state, never from the waking event
The operator SHALL generate the platform package from exactly two inputs: the `Platform` CR's spec, and the set of claims that are both accepted and active. Generation SHALL compute from the current state of both, never from the content of the event that woke it, so that repeated or reordered events yield the same package.

#### Scenario: Generation ignores the waking event's content

- **WHEN** regeneration is woken by a stale or duplicated event
- **THEN** the package computed matches the one the current state implies

### Requirement: A claim whose deletion is blocked remains in the active set
A claim whose deletion is blocked SHALL remain in the active set until the block releases. A deletion timestamp alone SHALL NOT remove a claim's catalog from the generated package: while the block holds, the claim's dependents are still rendering against that catalog, and dropping it would abandon them through the very door the block exists to close. Once the block releases the claim leaves the active set as any deleted claim does.

#### Scenario: A blocked claim's catalog survives regeneration

- **WHEN** the package is regenerated while an active claim is terminating and its deletion is blocked
- **THEN** the regenerated package still carries that claim's catalog, and its dependents keep rendering

#### Scenario: A released claim leaves the active set

- **WHEN** a terminating claim's block releases
- **THEN** the next regenerated package no longer carries that claim's catalog

### Requirement: The Platform status carries the union of the authored subscriptions and the active claims' catalogs with the source of each entry
`Platform` status SHALL carry the resolved union of the registry the platform is running: the subscriptions its spec authored and the catalogs its active claims contributed. An entry SHALL indicate which of the two it came from, so an operator can tell an authored subscription from one a provider registered.

#### Scenario: The union shows both sources

- **WHEN** a platform has an authored subscription and an active claim
- **THEN** the status union lists both, distinguishing which is which

#### Scenario: The union follows the active set

- **WHEN** a claim stops being active
- **THEN** the union no longer lists its catalog

### Requirement: The Platform status shows the pinned version and whether the transformers of an entry register
`Platform` status SHALL carry the resolved union of the registry the platform is running: the subscriptions its spec authored and the catalogs its active claims contributed. An entry SHALL also carry the version the package pins it at and whether its transformers register, because a disabled subscription belongs in the union, it is still pinned and imported, but contributes no transformer, and a union read without that distinction overstates what the platform runs.

#### Scenario: A disabled subscription is listed as not registering

- **WHEN** a platform authors a subscription with transformer registration disabled
- **THEN** the status union lists its catalog and marks it as contributing no transformer
