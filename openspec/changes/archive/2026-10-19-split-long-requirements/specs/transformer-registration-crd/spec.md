## REMOVED Requirements

### Requirement: Every claim field is required at admission
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "TransformerRegistrationSpec requires catalog, version, provides and providerRef" and "Claim validation is not delegated to the authoring catalog".

## ADDED Requirements

### Requirement: TransformerRegistrationSpec requires catalog, version, provides and providerRef
`TransformerRegistrationSpec` SHALL require `catalog` (the provider catalog's module path), `version` (its build), `provides` (the contract keys it claims to implement) and `providerRef` (the namespace and name of the claiming instance). The API server SHALL reject a claim missing any of them. `provides` is required but MAY be empty: the CRD sets no minimum on its length, and acceptance (not admission) compares the list with the contracts re-derived from the catalog, for exact equality.

#### Scenario: A claim missing its catalog is rejected

- **WHEN** a `TransformerRegistration` is applied with `spec.version`, `spec.provides` and `spec.providerRef` but no `spec.catalog`
- **THEN** the API server rejects it, naming the missing field

#### Scenario: A claim with an empty provides list is accepted

- **WHEN** a `TransformerRegistration` is applied whose `spec.provides` is an empty list
- **THEN** the API server accepts it, because the CRD sets no minimum on `provides`
- **AND** acceptance treats the empty list like any other list: it passes when the catalog implements no provider-fulfilled contract and is refused with reason `ProvidesMismatch` when the catalog implements any

### Requirement: Claim validation is not delegated to the authoring catalog
As in requirement "TransformerRegistrationSpec requires catalog, version, provides and providerRef": This validation SHALL NOT be delegated to the authoring catalog: a missing required field in CUE is an incomplete value rather than an error, so a claim can reach the cluster with a field absent even though the contract declares it required.

#### Scenario: A claim can reach the cluster with a field absent
- **WHEN** a required field is missing in CUE, which is an incomplete value rather than an error
- **THEN** a claim can reach the cluster with a field absent even though the contract declares it required, so this validation is not delegated to the authoring catalog
