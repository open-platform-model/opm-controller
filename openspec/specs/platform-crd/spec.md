## Purpose

Define the cluster-scoped singleton `Platform` custom resource: its CRD shape
(spec projecting core `#Platform`, status carrying conditions and
`observedGeneration`), its singleton enforcement, and its registration in the
runtime scheme, where a controller reconciles it (capability `platform-reconciler`).

## Requirements

### Requirement: Cluster-scoped singleton Platform resource

The operator SHALL define a `Platform` custom resource that is cluster-scoped and
constrained to a single instance. The only permitted `metadata.name` SHALL be
`cluster`, enforced declaratively by a CEL validation rule on the resource root
(no admission webhook). Because the resource is cluster-scoped, name uniqueness
guarantees at most one `Platform` can exist.

#### Scenario: Platform named cluster is accepted

- **WHEN** a `Platform` with `metadata.name: cluster` is applied
- **THEN** the API server accepts it

#### Scenario: Platform with any other name is rejected

- **WHEN** a `Platform` with `metadata.name` other than `cluster` is applied
- **THEN** the API server rejects it with the CEL validation message identifying the singleton constraint

#### Scenario: Platform is cluster-scoped

- **WHEN** the CRD is installed
- **THEN** `Platform` is registered with `scope: Cluster`
- **AND** `Platform` objects carry no namespace

### Requirement: Platform is registered in the scheme and reconciled

The operator SHALL register `Platform` and `PlatformList` with the runtime scheme so the types are serializable and installable. The manager SHALL also run a controller for `Platform` that reconciles the singleton named `cluster` (capability `platform-reconciler`), and it SHALL be woken by a change to the `Platform` spec and by a change to the set of accepted and active `TransformerRegistration` claims.

#### Scenario: Types are registered in the scheme

- **WHEN** the manager starts
- **THEN** `Platform` and `PlatformList` are registered in the scheme

#### Scenario: A Platform named cluster is reconciled

- **WHEN** the manager runs and a `Platform` named `cluster` is applied
- **THEN** the Platform controller reconciles it and records the outcome on `status.conditions`

### Requirement: PlatformSpec carries type, registry and skewPolicy and generates the platform module from type and registry
`PlatformSpec` SHALL carry `type` (required, informational discriminator), `registry` (path-keyed catalog subscriptions) and `skewPolicy` (optional enum `Warn` | `Refuse`; unset means `Warn`). The operator generates a platform CUE module from `type` and `registry`; `skewPolicy` is not part of the module.

#### Scenario: Skew policy defaults to Warn

- **WHEN** a Platform is applied without `spec.skewPolicy`
- **THEN** it is admitted and reconciles as `Warn`

#### Scenario: Minimal valid platform spec

- **WHEN** a `Platform` is applied with `spec.type` set and a `spec.registry` entry keyed by a catalog module path carrying only a `version`
- **THEN** the API server accepts it, the omitted `enable` is understood as the schema default (true) and the omitted `skewPolicy` as `Warn`

#### Scenario: type is required

- **WHEN** a `Platform` is applied without `spec.type`
- **THEN** the API server rejects it as a missing required field

#### Scenario: The skewPolicy is not part of the generated module
- **WHEN** the operator generates a platform CUE module
- **THEN** it is generated from `type` and `registry`; `skewPolicy` is not part of the module

### Requirement: skewPolicy answers a catalog requirement newer than the pinned build and is validated at admission
`skewPolicy` is the operator's response when a module's declared catalog requirement is newer than the build the platform pins: `Warn` renders and reports the skew; `Refuse` refuses the render. The field SHALL be validated at admission by an enum constraint.

#### Scenario: Invalid skew policy is rejected at admission

- **WHEN** a Platform sets `spec.skewPolicy: Strict`
- **THEN** the API server rejects it

### Requirement: The registry maps catalog paths to subscriptions of one published version
`PlatformSpec.Registry` SHALL map major-suffixed catalog module paths (`opmodel.dev/catalogs/opm@v2`) to subscriptions of the form `{enable?, version}`. `version` SHALL name exactly one published catalog build as a bare SemVer string (`MinLength=1`); the filter vocabulary (`filter.range`, `filter.allow`, `filter.deny`) SHALL NOT exist in the schema.

#### Scenario: Scalar subscription accepted

- **WHEN** a Platform is applied with `registry: {"opmodel.dev/catalogs/opm@v2": {version: "2.0.0-alpha.3"}}`
- **THEN** the API server accepts it and the reconciler generates a platform module that pins exactly that build

#### Scenario: Filter rejected

- **WHEN** a Platform is applied carrying `filter` under a subscription
- **THEN** the structural schema prunes or rejects the field; it never reaches the reconciler

### Requirement: A subscription without a version never reaches a platform build
A subscription without a version SHALL never reach a platform build. Whether `version` is schema-`required` or schema-optional follows the recorded ratcheting measurement; the posture changes only the rejecting actor: admission in the required posture, the reconciler before any registry I/O (surfaced as `BuildFailed` naming the subscription path) in the optional posture.

#### Scenario: Stored legacy singleton remains serviceable

- **WHEN** the stored `cluster` Platform predates the schema change (carries `filter`, lacks `version`)
- **THEN** the operator's status patches against it continue to succeed
- **AND** its subscriptions fail with `BuildFailed` naming the registry entry's missing version until the spec is re-applied

### Requirement: PlatformStatus carries conditions, observedGeneration and operatorVersion
`PlatformStatus` SHALL carry `conditions` (a `metav1.Condition` list keyed by type), `observedGeneration` and `operatorVersion`.

#### Scenario: Status subresource present

- **WHEN** the CRD is installed
- **THEN** `Platform` exposes a `/status` subresource, and `status.conditions`, `status.observedGeneration` and `status.operatorVersion` are part of the schema

#### Scenario: Operator printcolumn

- **WHEN** `kubectl get platform` runs against a reconciled Platform
- **THEN** the output includes an `Operator` column showing `.status.operatorVersion`

### Requirement: The Ready and ContractsFulfilled conditions carry their reasons
The Ready condition SHALL summarise module generation: `Ready=True` reason `Generated`; `Ready=False` reason `BuildFailed`, `GenerateFailed`, `ContractCollisions`, `OverSubscribedContracts` or `ComparablePredicates`. A `ContractsFulfilled` condition SHALL report the effective package's unfulfilled provider-fulfilled contracts without affecting Ready (reasons `UnfulfilledContracts`, `ContractsFulfilled`, `NoContractsDefined`).

#### Scenario: Status reflects the generate-and-build outcome

- **WHEN** the Platform reconciles
- **THEN** `status.conditions` carries a Ready condition with one of the reasons `Generated`, `BuildFailed`, `GenerateFailed`, `ContractCollisions`, `OverSubscribedContracts` or `ComparablePredicates`

#### Scenario: Status carries the contract report

- **WHEN** the Platform reconciles to `Ready=True`
- **THEN** `status.conditions` carries a `ContractsFulfilled` condition with one of the reasons `UnfulfilledContracts`, `ContractsFulfilled` or `NoContractsDefined`

### Requirement: The documentation of the conditions field names every reason
As in requirement "PlatformStatus carries conditions, observedGeneration and operatorVersion": The field's documentation SHALL name every reason.

#### Scenario: The field documentation names every reason
- **WHEN** the `conditions` field of `PlatformStatus` is documented
- **THEN** its documentation names every reason
