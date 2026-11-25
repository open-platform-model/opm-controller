## ADDED Requirements

### Requirement: A catalog version is read in its bare spelling

The controller SHALL read every catalog version that a subscription or a claim names in its bare SemVer spelling, the spelling the generated platform module stamps: `0.1.0` and `v0.1.0` SHALL be the same version. The package identity SHALL use that spelling for each claim coordinate, and the status union SHALL record it for each entry.

#### Scenario: A spelling-only version edit keeps the identity

- **WHEN** an active claim's `spec.version` changes from `0.1.0` to `v0.1.0`
- **THEN** `status.packageIdentity` keeps its value, and the platform package is not regenerated

#### Scenario: A v-prefixed version is listed bare

- **WHEN** a subscription names `v2.0.0-alpha.3` and an active claim names `v0.1.0`
- **THEN** the status union lists their versions as `2.0.0-alpha.3` and `0.1.0`
