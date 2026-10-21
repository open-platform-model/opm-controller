## MODIFIED Requirements

### Requirement: Example modules published to GHCR on release

On a controller release, CI SHALL publish each example module to `ghcr.io/open-platform-model` under its `testing.opmodel.dev/modules/controller/<m>` path, using a registry mapping that routes BOTH `testing.opmodel.dev` (the fixtures) and `opmodel.dev` (their core and catalog dependencies) to GHCR. A mapping missing either domain falls through to a public mirror and fails with 401.

Publishing SHALL go through `opm module publish`, installed in CI from a pinned cli release, so every publish gate runs over each fixture and a fixture that violates a gate fails the release instead of shipping. The declared version SHALL be read from the module's own `identity/` package; it SHALL NOT be passed in by the caller, and a module without an identity package SHALL be a hard error rather than a silent skip.

#### Scenario: Release triggers module publish

- **WHEN** a release is created
- **THEN** the release job publishes each example module via `opm module publish` to `ghcr.io/open-platform-model` under its `testing.opmodel.dev/modules/controller/<m>` path

#### Scenario: No publish without a release

- **WHEN** a commit lands on the default branch without creating a release
- **THEN** no example module is published

#### Scenario: Fixture missing an identity package

- **WHEN** a directory under the fixtures root carries no `identity/` package
- **THEN** the publish step SHALL fail, naming the directory
- **AND** SHALL NOT report success

### Requirement: Example manifests uploaded as release artifacts

The release SHALL attach the example `moduleinstance.yaml` manifests, each referencing a version the same release published.

#### Scenario: Attached manifests reference published modules

- **WHEN** a user downloads an attached example manifest
- **THEN** its `spec.module.path` references a `testing.opmodel.dev/modules/controller/<m>@v0` version that the same release published

## REMOVED Requirements

### Requirement: Independent per-module versioning

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-021); this requirement is restated under the controller's name.

**Migration**: See the requirement "Each example module versions on its own" in this capability.

## ADDED Requirements

### Requirement: Each example module versions on its own

Each example module SHALL declare its own semver in its `identity/` package, independent of the controller's release version. The publish step SHALL publish that version and no other.

#### Scenario: Module version differs from controller version

- **WHEN** the controller releases `v0.5.0` and `podinfo`'s identity package declares `0.1.0`
- **THEN** the publish step publishes `testing.opmodel.dev/modules/controller/podinfo` at `v0.1.0`, not `v0.5.0`
