## REMOVED Requirements

### Requirement: CUE_REGISTRY configuration for tests
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "Tests resolve fixtures and their dependencies through the configured mapping", "A local-iteration mapping redirects only the fixture prefix" and "Publish tooling exports OPM_REGISTRY as well as CUE_REGISTRY".

### Requirement: End-to-end integration tests
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "An integration test exercises the real renderer against the local registry", "The integration test resolves the catalog from the generated platform" and "The render result carries no inventory entries and full apply is covered by the Kind-backed e2e suite".

## ADDED Requirements

### Requirement: Tests resolve fixtures and their dependencies through the configured mapping
Tests SHALL resolve the fixture modules from whatever registry the configured mapping routes `testing.opmodel.dev` to, and their core and catalog dependencies from wherever it routes `opmodel.dev`. The canonical mapping routes both to `ghcr.io/open-platform-model`.

#### Scenario: Catalog present enables the platform build

- **WHEN** the configured mapping resolves the subscribed catalog build
- **THEN** the platform module builds

### Requirement: A local-iteration mapping redirects only the fixture prefix
Because the fixtures and their dependencies now occupy SEPARATE domains, a local-iteration mapping SHALL redirect only the fixture prefix, `testing.opmodel.dev=localhost:5000+insecure`, while core and the catalogs continue to resolve from GHCR. Mapping all of `opmodel.dev` to a local registry SHALL NOT be required, and mirroring core or the catalogs locally SHALL NOT be required.

#### Scenario: Fixture prefix redirected alone

- **WHEN** `CUE_REGISTRY` maps `testing.opmodel.dev` to the local registry and `opmodel.dev` to GHCR
- **THEN** the fixture resolves locally and core and the catalogs resolve from GHCR

### Requirement: Publish tooling exports OPM_REGISTRY as well as CUE_REGISTRY
Tooling that invokes the publish CLI SHALL export `OPM_REGISTRY` as well as `CUE_REGISTRY`: the CLI resolves `--registry` > `OPM_REGISTRY` > its config file and never reads `CUE_REGISTRY`, so exporting only the latter silently leaves it on the caller's personal configuration.

#### Scenario: Exporting only CUE_REGISTRY leaves the CLI on the caller's personal configuration
- **WHEN** tooling that invokes the publish CLI exports only `CUE_REGISTRY`
- **THEN** the CLI never reads `CUE_REGISTRY`, so it stays on the caller's personal configuration

### Requirement: An integration test exercises the real renderer against the local registry
At least one integration test MUST exercise the real renderer (`render.KernelModuleRenderer`) against the local OCI registry, generating a platform from the real catalog, to validate the registry-backed render pipeline: module acquisition from the registry (`moduleacquire.Acquire`) → kernel `SynthesizeInstance` → kernel `Render` → rendered resources with provenance and runtime-identity labels.

#### Scenario: Real-renderer pipeline validated against the registry
- **WHEN** the integration test runs with the local registry available
- **THEN** it constructs `render.KernelModuleRenderer` with a generated platform (built via `platformmodule` generation and `AcquirePlatformFromDir`, recorded with `SetGenerated`), renders a ModuleInstance, and the rendered resources carry the runtime-identity labels (`managed-by = opm-controller`, non-empty module-instance uuid label)

### Requirement: The integration test resolves the catalog from the generated platform
As in requirement "An integration test exercises the real renderer against the local registry": The test MUST resolve the catalog from the generated platform (the same path the `PlatformReconciler` uses) rather than copying catalog sources into `test/fixtures/`, so it tracks production composition automatically.

#### Scenario: Catalog resolved from the generated platform
- **WHEN** the integration test generates the platform through `platformmodule` generation and `AcquirePlatformFromDir` and records it with `SetGenerated`
- **THEN** the catalog is resolved from the registry via the kernel rather than a copy under `test/fixtures/`, so the test automatically tracks production composition

### Requirement: The render result carries no inventory entries and full apply is covered by the Kind-backed e2e suite
As in requirement "An integration test exercises the real renderer against the local registry": The render result carries no inventory entries; the reconciler builds them from its one export of the resources. Full apply → `Ready=True` on a live cluster SHALL be covered by the Kind-backed `test/e2e` suite, not this integration-tier test.

#### Scenario: The render result carries no inventory entries
- **WHEN** the integration test renders a ModuleInstance
- **THEN** the render result carries no inventory entries; the reconciler builds them from its one export of the resources

#### Scenario: Full apply to Ready is left to the Kind-backed suite
- **WHEN** the integration test runs
- **THEN** full apply to `Ready=True` on a live cluster is covered by the Kind-backed `test/e2e` suite, not this integration-tier test
