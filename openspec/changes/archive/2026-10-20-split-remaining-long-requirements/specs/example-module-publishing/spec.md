## REMOVED Requirements

### Requirement: Example modules published to GHCR on release
**Reason**: Split into three requirements that each state one behavior.
**Migration**: See "Example modules are published to GHCR on a release through a registry mapping that routes both domains", "Publishing goes through opm module publish so every gate runs over each fixture" and "The declared version is read from the identity package of the module".

### Requirement: Pre-release publishes move the declared version with the tag
**Reason**: Split into two requirements that each state one behavior.
**Migration**: See "A pre-release publish tags the artifact and declares the same version" and "A pre-release publish rewrites the declared version on a staged copy without the version flag".

## ADDED Requirements

### Requirement: Example modules are published to GHCR on a release through a registry mapping that routes both domains
On an operator release, CI SHALL publish each example module to `ghcr.io/open-platform-model` under its `testing.opmodel.dev/modules/operator/<m>` path, using a registry mapping that routes BOTH `testing.opmodel.dev` (the fixtures) and `opmodel.dev` (their core and catalog dependencies) to GHCR. A mapping missing either domain falls through to a public mirror and fails with 401.

#### Scenario: Release triggers module publish

- **WHEN** a release is created
- **THEN** the release job publishes each example module via `opm module publish` to `ghcr.io/open-platform-model` under its `testing.opmodel.dev/modules/operator/<m>` path

#### Scenario: No publish without a release

- **WHEN** a commit lands on the default branch without creating a release
- **THEN** no example module is published

### Requirement: Publishing goes through opm module publish so every gate runs over each fixture
As in requirement "Example modules are published to GHCR on a release through a registry mapping that routes both domains": Publishing SHALL go through `opm module publish`, installed in CI from a pinned cli release, so every publish gate runs over each fixture and a fixture that violates a gate fails the release instead of shipping.

#### Scenario: A fixture that violates a publish gate fails the release
- **WHEN** a fixture violates a gate of `opm module publish` during a release
- **THEN** the release fails instead of shipping the fixture

### Requirement: The declared version is read from the identity package of the module
As in requirement "Example modules are published to GHCR on a release through a registry mapping that routes both domains": The declared version SHALL be read from the module's own `identity/` package; it SHALL NOT be passed in by the caller, and a module without an identity package SHALL be a hard error rather than a silent skip.

#### Scenario: Fixture missing an identity package

- **WHEN** a directory under the fixtures root carries no `identity/` package
- **THEN** the publish step SHALL fail, naming the directory
- **AND** SHALL NOT report success

### Requirement: A pre-release publish tags the artifact and declares the same version
When publishing with a pre-release identifier, the publish step SHALL tag the artifact `v<version>-<prerelease>` AND SHALL make the artifact's own declared version equal that same value, because the acquire-time identity check refuses a module whose declared version differs from the fetched tag.

#### Scenario: Pre-release tag and declared version agree

- **WHEN** the fixtures are published with pre-release identifier `e2e.gabc1234`
- **THEN** each is tagged `v<version>-e2e.gabc1234`
- **AND** the published artifact's declared version is `<version>-e2e.gabc1234`

### Requirement: A pre-release publish rewrites the declared version on a staged copy without the version flag
As in requirement "A pre-release publish tags the artifact and declares the same version": The declared version SHALL be rewritten on a staged copy of the tree using the identity-package writer, leaving the working tree unmodified. The publish command's `--version` flag SHALL NOT be used for this: it asserts an already-declared version rather than filling one, and these fixtures declare a defaulted version.

#### Scenario: Working tree unmodified by a pre-release publish

- **WHEN** a pre-release publish completes
- **THEN** no file under the fixture trees has changed on disk
