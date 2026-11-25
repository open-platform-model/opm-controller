## MODIFIED Requirements

### Requirement: The registry maps catalog paths to subscriptions of one published version
`PlatformSpec.Registry` SHALL map major-suffixed catalog module paths (`opmodel.dev/catalogs/opm@v2`) to subscriptions of the form `{enable?, version}`. `version` SHALL name exactly one published catalog build as a SemVer string (`MinLength=1`), bare (`2.0.0-alpha.3`) or with a leading `v` (`v2.0.0-alpha.3`); both spellings SHALL name the same build. The filter vocabulary (`filter.range`, `filter.allow`, `filter.deny`) SHALL NOT exist in the schema.

#### Scenario: Scalar subscription accepted

- **WHEN** a Platform is applied with `registry: {"opmodel.dev/catalogs/opm@v2": {version: "2.0.0-alpha.3"}}`
- **THEN** the API server accepts it and the reconciler generates a platform module that pins exactly that build

#### Scenario: Filter rejected

- **WHEN** a Platform is applied carrying `filter` under a subscription
- **THEN** the structural schema prunes or rejects the field; it never reaches the reconciler
