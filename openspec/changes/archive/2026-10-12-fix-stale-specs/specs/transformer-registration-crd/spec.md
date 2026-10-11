## MODIFIED Requirements

### Requirement: Every claim field is required at admission

`TransformerRegistrationSpec` SHALL require `catalog` (the provider catalog's module path), `version` (its build), `provides` (the contract keys it claims to implement) and `providerRef` (the namespace and name of the claiming instance). The API server SHALL reject a claim missing any of them. `provides` is required but MAY be empty: the CRD sets no minimum on its length, and acceptance (not admission) compares the list with the contracts re-derived from the catalog, for exact equality. This validation SHALL NOT be delegated to the authoring catalog: a missing required field in CUE is an incomplete value rather than an error, so a claim can reach the cluster with a field absent even though the contract declares it required.

#### Scenario: A claim missing its catalog is rejected

- **WHEN** a `TransformerRegistration` is applied with `spec.version`, `spec.provides` and `spec.providerRef` but no `spec.catalog`
- **THEN** the API server rejects it, naming the missing field

#### Scenario: A claim with an empty provides list is accepted

- **WHEN** a `TransformerRegistration` is applied whose `spec.provides` is an empty list
- **THEN** the API server accepts it, because the CRD sets no minimum on `provides`
- **AND** acceptance treats the empty list like any other list: it passes when the catalog implements no provider-fulfilled contract and is refused with reason `ProvidesMismatch` when the catalog implements any

## REMOVED Requirements

### Requirement: Claim status separates acceptance from activation

**Reason**: The scenario `A stored claim is inert` is false. An active claim is an input to the generated platform module, so it does change what the operator renders. OpenSpec refuses to drop or rename a scenario inside a MODIFIED requirement, so the requirement is retired and added again under a new name with the corrected scenario.

**Migration**: See the requirement `Claim status reports acceptance and activation as separate states`, added below. Its first scenario is unchanged.

## ADDED Requirements

### Requirement: Claim status reports acceptance and activation as separate states

`TransformerRegistrationStatus` SHALL carry `conditions`, `accepted` and `active`, so that a claim reports one of three states: not yet judged, accepted but inactive, or active. A controller SHALL watch the kind and record a verdict on every claim, setting `conditions`, `accepted` and `observedGeneration`, and SHALL set `active` when an accepted claim's provider is ready. Only a claim that is both accepted and active is an input to the generated platform module.

#### Scenario: A stored claim receives a verdict

- **WHEN** a valid `TransformerRegistration` is applied and the manager is running
- **THEN** a reconcile is triggered for it and its status reports whether it was accepted
- **AND** `observedGeneration` matches the claim's generation

#### Scenario: Only an accepted and active claim feeds the generated platform

- **WHEN** a `TransformerRegistration` is stored and reconciled to refused, accepted but inactive, or active
- **THEN** a refused claim and an accepted but inactive claim are not inputs to the generated platform module
- **AND** an active claim is an input, and its catalog is imported unless a subscription in `Platform.spec.registry` already names that catalog
