## MODIFIED Requirements

### Requirement: Test module publication

A fixture module MUST be publishable to a registry as `testing.opmodel.dev/modules/controller/<name>@v0` at the version its own identity package declares. The publish task SHALL force its registry mapping in-script so an ambient environment cannot redirect it, and SHALL derive the version from the identity package rather than accepting it as a parameter.

#### Scenario: Publish task uses task-local registry mapping

- **WHEN** the publish task runs with an unrelated mapping exported in the shell
- **THEN** it publishes against its own forced mapping

#### Scenario: Fixture published at its declared coordinates

- **WHEN** the publish task completes for `hello`
- **THEN** `testing.opmodel.dev/modules/controller/hello@v0` is available at the version its identity package declares

#### Scenario: Stale old-era artifact is inert

- **WHEN** an artifact published under a retired path remains in a registry
- **THEN** no test resolves it

### Requirement: Test fixture module path

The fixture modules MUST use module paths `testing.opmodel.dev/modules/controller/<name>@v0`, and each `#Module.metadata.modulePath` MUST equal that full major-suffixed address, derived from the module's identity package.

#### Scenario: Fixture published under the fixture namespace

- **WHEN** a fixture is resolved by a consumer
- **THEN** it resolves under the `testing.opmodel.dev` domain, which is reserved for fixtures and experiments and is routed to GHCR by the canonical mapping

#### Scenario: Fixture is concrete-valid standalone

- **WHEN** a fixture module is evaluated on its own
- **THEN** it evaluates without error and its identity package yields a concrete version
