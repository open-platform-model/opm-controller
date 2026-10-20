# example-test-modules Specification

## Purpose
The controller's example test module fleet and its modulepackage fixtures: how they are versioned, pinned, checked and published.

## Requirements

### Requirement: podinfo web example module

The repo SHALL provide a podinfo example module modelling a stateless web workload. It SHALL render a Deployment and a Service exposing the podinfo HTTP port (9898), and SHALL declare an HTTP `livenessProbe` against `/healthz` and an HTTP `readinessProbe` against `/readyz` on that port.

#### Scenario: Renders deployment with probes
- **WHEN** the podinfo module is rendered
- **THEN** the output includes a Deployment whose container declares a `livenessProbe.httpGet` path `/healthz` and a `readinessProbe.httpGet` path `/readyz`, both on port 9898
- **AND** the output includes a Service targeting port 9898

#### Scenario: Configurable replicas and image
- **WHEN** a ModuleRelease overrides the podinfo replica count or image tag
- **THEN** the rendered Deployment reflects the overridden values

### Requirement: redis stateful example module

The repo SHALL provide a redis example module modelling a stateful workload. It SHALL render a StatefulSet, a headless Service, and a PersistentVolumeClaim (or volumeClaimTemplate), and SHALL declare an exec readiness probe running `redis-cli ping`.

#### Scenario: Renders statefulset with persistence and probe
- **WHEN** the redis module is rendered
- **THEN** the output includes a StatefulSet with a volume claim and a headless Service
- **AND** the container declares an exec readiness probe invoking `redis-cli ping`

#### Scenario: Persistence default is documented and overridable
- **WHEN** the redis module is authored
- **THEN** its persistence behavior (ephemeral vs PVC) has an explicit documented default and is overridable via module config

### Requirement: Example modules ship ready-to-apply manifests

Each example module SHALL include a `ModuleInstance` manifest (and, where applicable, `ModulePackage`/`OCIRepository` manifests) that a user can apply against a running controller to deploy the example. The manifests SHALL reference the `testing.opmodel.dev/modules/controller/<m>@v0` path and a concrete version.

#### Scenario: Manifest references public module

- **WHEN** a user inspects an example module's `moduleinstance.yaml`
- **THEN** its `spec.module.path` is `testing.opmodel.dev/modules/controller/podinfo@v0` with an explicit `spec.module.version`

### Requirement: podinfo liveness/readiness e2e validation

An e2e test SHALL deploy the podinfo `ModuleRelease` against a Kind-backed controller and assert that the rendered probes function — the Deployment's pods become Ready, which requires both the liveness and readiness probes to pass.

#### Scenario: podinfo pods become ready
- **WHEN** the e2e suite applies the podinfo ModuleRelease and the controller reconciles it
- **THEN** the resulting Deployment's pods SHALL reach Ready within the test timeout, demonstrating the modelled probes succeed against the running container

#### Scenario: Rendered probe contract matches container
- **WHEN** the e2e suite inspects the deployed podinfo container
- **THEN** the container's liveness/readiness probe HTTP paths and port match the values declared in the module (`/healthz`, `/readyz`, 9898)

### Requirement: hello_web ready-to-apply ModuleInstance

The `hello_web` example SHALL ship a `ModuleInstance` manifest referencing the `testing.opmodel.dev/modules/controller/hello_web@v0` path and a concrete version, so all four fleet members are applyable.

#### Scenario: hello_web manifest references public module

- **WHEN** a user inspects `hello_web`'s `moduleinstance.yaml`
- **THEN** its `spec.module.path` is `testing.opmodel.dev/modules/controller/hello_web@v0` with an explicit `spec.module.version`

### Requirement: ModulePackage render integration coverage

The `KernelPackageRenderer` integration test SHALL exercise every modulepackage fixture
(`hello`, `hello_web`, `podinfo`, `redis`), not only `hello`, asserting that each authored
package loads its imported `#Module` and renders at least one resource carrying the
controller's ownership labels. The coverage SHALL remain gated on the local test registry
(skipping when it is unavailable).

#### Scenario: Every modulepackage renders under a generated platform
- **WHEN** the integration suite runs with a reachable test registry and a generated platform
- **THEN** each of the `hello`, `hello_web`, `podinfo`, `redis` modulepackage fixtures loads without a "field not allowed" error and renders at least one resource
- **AND** each rendered resource carries the `managed-by` controller label and a non-empty module-instance UUID label

#### Scenario: Coverage skips without the test registry
- **WHEN** the integration suite runs without the local test registry configured
- **THEN** the modulepackage render cases skip rather than fail

### Requirement: Example test modules are authored under the public module path
All example test modules SHALL be authored under the CUE module path `testing.opmodel.dev/modules/controller/<module>@v0`, where `<module>` is the module's short name.

#### Scenario: New module declares public path

- **WHEN** a new example module is added
- **THEN** its `cue.mod/module.cue` declares `module: "testing.opmodel.dev/modules/controller/podinfo@v0"` and `metadata.modulePath` resolves to the same value

### Requirement: A module declares and publishes the public module path
As in requirement "Example test modules are authored under the public module path": Each module's `cue.mod/module.cue` `module:` field and its `#Module.metadata.modulePath` SHALL agree with this path, and SHALL be published to `ghcr.io/open-platform-model` so the module resolves under the canonical registry mapping, which routes both `opmodel.dev` (core, catalogs) and `testing.opmodel.dev` (these fixtures) to GHCR.

#### Scenario: Consumer resolves without extra config

- **WHEN** a user with the canonical `CUE_REGISTRY` mapping resolves `testing.opmodel.dev/modules/controller/podinfo@v0`
- **THEN** it resolves from `ghcr.io/open-platform-model` with no additional registry configuration

### Requirement: Fixtures are not authored under opmodel.dev
Fixtures SHALL NOT be authored under `opmodel.dev/*`. CUE resolves modules by longest-prefix match on the module path, so a fixture in the production namespace forces that entire prefix, core and the catalogs included, onto whatever registry serves the fixture. The publish gates enforce this independently: a nested path under `opmodel.dev` is refused.

#### Scenario: Production namespace is refused

- **WHEN** a fixture declares a path under `opmodel.dev/modules/test/`
- **THEN** `opm module publish` refuses it as not fitting the namespace the domain publishes

### Requirement: Each workload fleet module has a sibling modulepackage fixture
Each example module of the workload fleet (`hello`, `hello_web`, `podinfo`, `redis`) SHALL have a sibling modulepackage fixture declaring `testing.opmodel.dev/releases/controller/<module>@v0`, an `instance.cue` that imports the published module, and a `cue.mod/module.cue` that pins the same `testing.opmodel.dev/modules/controller/<module>@v0` version the module declares.

#### Scenario: Each module has a modulepackage fixture

- **WHEN** the modulepackage fixtures are enumerated
- **THEN** there SHALL be one per workload fleet module, declaring `testing.opmodel.dev/releases/controller/<module>@v0`

#### Scenario: instance.cue imports and embeds the published module

- **WHEN** a modulepackage fixture is loaded
- **THEN** it embeds `core.#ModuleInstance`, imports `testing.opmodel.dev/modules/controller/<module>@v0`, and sets `#module` to the imported module

#### Scenario: ModulePackage CR references its OCIRepository

- **WHEN** a modulepackage fixture's `ModulePackage` is inspected
- **THEN** its `spec.sourceRef` names the `OCIRepository` declared in the sibling `ocirepository.yaml`, whose `url` ends in `testing.opmodel.dev/releases/controller/<module>`

### Requirement: The backup fixture set has no modulepackage fixture
As in requirement "Each workload fleet module has a sibling modulepackage fixture": The backup fixture set (`backup_provider`, `backup_consumer`) SHALL have no modulepackage fixture: it exists to exercise a registration, and `backup_consumer` does not render on a platform without an active provider.

#### Scenario: The backup fixture set has no modulepackage fixture

- **WHEN** the modulepackage fixtures are enumerated
- **THEN** none is declared for `backup_provider` or `backup_consumer`

### Requirement: The example test module fleet is the workload fleet and the backup fixture set
The example test module fleet SHALL be the workload fleet `hello`, `hello_web`, `podinfo` and `redis`, plus the backup fixture set `backup_provider` and `backup_consumer`. Catalog imports SHALL stay on the versioned packages of `opmodel.dev/catalogs/opm@v4`. Each module's `moduleinstance.yaml` SHALL pin the module's current published `v`-prefixed version.

#### Scenario: Fleet renders on the v2 line

- **WHEN** each workload fleet member is loaded and rendered
- **THEN** it renders against `opmodel.dev/core@v2` and the versioned catalog packages without error

### Requirement: Each fleet module carries an identity package that its metadata derives from
As in requirement "The example test module fleet is the workload fleet and the backup fixture set": Each SHALL carry an `identity/identity.cue` package as the single source of its module path and version (core `#IdentityPackage`), and its `#Module.metadata` SHALL DERIVE from that package rather than restate it: `metadata.modulePath` is the identity package's `ModulePath`, `metadata.version` its `Version`, and `metadata.name` the path's leaf in snake case.

#### Scenario: Metadata derives from identity

- **WHEN** a fixture's `identity/identity.cue` declares a path and version
- **THEN** its `metadata.modulePath`, `metadata.version`, and `metadata.name` evaluate to values derived from that package
- **AND** `opm module publish` reports no derivation disagreement

### Requirement: A version bump is an edit to the identity package
As in requirement "Each fleet module carries an identity package that its metadata derives from": A version bump SHALL therefore be an edit to `identity/identity.cue` (directly or via `opm module version set`), never to the metadata block.

#### Scenario: A version bump is made in the identity package

- **WHEN** a fixture's version is bumped
- **THEN** the edit is to `identity/identity.cue` (directly or via `opm module version set`), never to the metadata block

### Requirement: The hyphenated hello-web name is retired
The hyphenated `hello-web` name SHALL be retired at the source: a core-v2 module cannot carry a hyphen, so the fixture publishes as `hello_web`. Artifacts previously published under `opmodel.dev/modules/test/*`, including the hyphenated `hello-web`, are unmodified by this change; their deletion is owned by enhancement 0011's `registry-cleanup`.

#### Scenario: Hyphenated name absent from the fleet

- **WHEN** the fleet's module paths are enumerated
- **THEN** none contains a hyphen in its leaf segment

### Requirement: A modulepackage pins what CUE resolves for its module version
Each modulepackage fixture under `test/fixtures/modulepackages/<m>` SHALL pin every dependency it shares with its module (`test/fixtures/modules/<m>`) exactly as CUE resolves it for the module version it names: what `cue mod get <module>@<version>` followed by `cue mod tidy` writes. CUE keeps a dependency the modulepackage already lists at its listed version, so a version re-pin alone renders the new module against an older core without any error.

#### Scenario: A modulepackage pins what CUE resolves for its module version

- **WHEN** a modulepackage is checked against its module
- **THEN** every dependency they share is pinned exactly as CUE resolves it for the module version the modulepackage names

### Requirement: task examples:pin sets every shared dependency to the module's pin
`task examples:pin` SHALL, after re-pinning a modulepackage to its module's version, set every other dependency the modulepackage shares with the module to the module's pin.

#### Scenario: Manual bump moves core with the module

- **WHEN** a module's `cue.mod` pins a newer core than its modulepackage and `opm module version set` then `task examples:pin` run
- **THEN** the modulepackage pins the module's new version and the module's core and catalog versions

### Requirement: task examples:consumers checks every modulepackage and the test workflow runs it
`task examples:consumers` SHALL run `hack/fixtures.sh consumers` over every modulepackage, resolving through `CONSUMERS_CUE_REGISTRY` (GHCR by default), and the `test.yml` workflow SHALL run it after seeding the job-local registry, with the mixed mapping.

#### Scenario: Modulepackages that follow their module pass

- **WHEN** every modulepackage pins what CUE resolves for its module version
- **THEN** `task examples:consumers` prints `ok` for each and exits zero

### Requirement: The consumers check fails naming each modulepackage that differs or cannot be resolved
As in requirement "task examples:consumers checks every modulepackage and the test workflow runs it": It SHALL fail naming each modulepackage whose pins differ from what CUE resolves or whose module version cannot be resolved, after checking all of them.

#### Scenario: Stale modulepackage fails the test workflow

- **WHEN** a pull request re-pins a modulepackage to a new module version but leaves its core pin lower than the module's
- **THEN** the `test.yml` consumers step fails, printing the core pin diff and a `FAIL` line naming that modulepackage

#### Scenario: Unresolvable module version is a named failure

- **WHEN** a modulepackage pins a module version the configured registry does not hold
- **THEN** the step prints a `FAIL` line naming the modulepackage and the version, still checks the others, and exits non-zero

### Requirement: The repo carries a provider catalog fixture
The repo SHALL carry a catalog fixture under `test/fixtures/catalogs/provider`, declaring the CUE module path `testing.opmodel.dev/catalogs/controller/provider@v0`. It SHALL be the smallest valid core `#Catalog`: one resource and one transformer that requires it, each with an FQN under the catalog's own registry path, so a platform can subscribe it beside `opmodel.dev/catalogs/opm@v4` without a key collision.

#### Scenario: The provider catalog fixture declares its module path

- **WHEN** the catalog fixtures of the repo are enumerated
- **THEN** `test/fixtures/catalogs/provider` declares `testing.opmodel.dev/catalogs/controller/provider@v0`

### Requirement: The provider catalog fixture carries an identity package
As in requirement "The repo carries a provider catalog fixture": It SHALL carry an `identity/identity.cue` package as the single source of its module path and version, with `Version` a plain SemVer literal.

#### Scenario: The identity package carries a plain SemVer Version

- **WHEN** the identity package of the provider catalog fixture is read
- **THEN** it is the single source of the module path and version, and `Version` is a plain SemVer literal

### Requirement: The provider catalog fixture is checked by hack/fixtures.sh
As in requirement "The repo carries a provider catalog fixture": The fixture SHALL be published through `hack/fixtures.sh`, the same flow as the module fleet, with `opm catalog publish` in place of `opm module publish`: `check` runs every catalog publish gate as a dry run and enforces changed-implies-bumped.
(Requirement: The provider catalog fixture is seeded and published by hack/fixtures.sh)

#### Scenario: The catalog passes the catalog publish gates

- **WHEN** `opm catalog publish --dry-run test/fixtures/catalogs/provider` runs
- **THEN** the plan resolves with no refusals, checking two members through the member FQN gate, and derives repository `testing.opmodel.dev/catalogs/controller/provider` and the tag from the identity package

#### Scenario: A changed catalog must carry a new version

- **WHEN** a pull request changes `test/fixtures/catalogs/provider` but keeps a version GHCR already holds
- **THEN** `task examples:check` fails, naming the fixture and `opm catalog version set` as the fix

### Requirement: The provider catalog fixture is seeded and published by hack/fixtures.sh
As in requirement "The repo carries a provider catalog fixture": The fixture SHALL be published through `hack/fixtures.sh`, the same flow as the module fleet, with `opm catalog publish` in place of `opm module publish`: `seed` publishes it into the job-local registry before the registry-backed specs run in `test.yml`, and `publish` puts it on GHCR when `publish-fixtures.yml` runs after a merge that touches `test/fixtures/catalogs/`.

#### Scenario: PR CI seeds the catalog from the tree

- **WHEN** `test.yml` runs on a pull request
- **THEN** the seed step publishes `testing.opmodel.dev/catalogs/controller/provider` at its declared version into the job-local registry before any registry-backed spec runs

### Requirement: Registry-backed specs use the provider catalog fixture
As in requirement "The repo carries a provider catalog fixture": Registry-backed specs that need a resolvable catalog other than `opmodel.dev/catalogs/opm@v4` SHALL use this fixture and SHALL read its coordinate from its identity package (`fixtures.MustCatalog`), never a literal.

#### Scenario: The claim specs build against the test catalog

- **WHEN** the claim-driven regeneration specs activate a claim for a catalog the Platform does not subscribe
- **THEN** that catalog is `testing.opmodel.dev/catalogs/controller/provider@v0` at the version its identity package declares, and the regenerated platform pins, imports and builds it

### Requirement: The repo carries a backup catalog fixture that implements only the backup trait
The repo SHALL carry a catalog fixture under `test/fixtures/catalogs/backup`, declaring `testing.opmodel.dev/catalogs/controller/backup@v0`, whose transformers implement opm's provider-fulfilled backup trait (`opmodel.dev/catalogs/opm/traits/backup@v1alpha1`) and nothing else, so the provider contracts derived from it are exactly that trait.

#### Scenario: The catalog provides exactly the backup trait

- **WHEN** the backup catalog fixture is acquired from the registry at its declared build and its provider contracts are derived
- **THEN** they are exactly `opmodel.dev/catalogs/opm/traits/backup@v1alpha1`

### Requirement: The backup catalog fixture renders without a CRD, publishes through hack/fixtures.sh and is not moved by the cascade
As in requirement "The repo carries a backup catalog fixture that implements only the backup trait": It SHALL render the trait without any CRD, publish through `hack/fixtures.sh` like the provider catalog fixture, and pin core and the opm catalog no newer than the build the registry-backed specs subscribe. `task deps:cascade` SHALL NOT be required to move it.

#### Scenario: The catalog passes the catalog publish gates

- **WHEN** `opm catalog publish --dry-run test/fixtures/catalogs/backup` runs against GHCR
- **THEN** the plan resolves with no refusals, checking one member through the member FQN gate

### Requirement: The backup catalog fixture keeps a core pin older than ProvidesSince while the fold exists
As in requirement "The repo carries a backup catalog fixture that implements only the backup trait": While the library keeps its deprecated provider-set fold, its core pin SHALL stay older than the library's `ProvidesSince`, since the controller's old-catalog acceptance spec relies on it.

#### Scenario: The core pin stays older than ProvidesSince while the fold exists

- **WHEN** the library keeps its deprecated provider-set fold
- **THEN** the core pin of `test/fixtures/catalogs/backup` stays older than the library's `ProvidesSince`

### Requirement: backup_provider renders one TransformerRegistration
`backup_provider` SHALL render exactly one `TransformerRegistration`, through opm's `transformer-registration` contract, that names the backup catalog fixture at the version its identity package declares and lists exactly the provider contracts that catalog implements.

#### Scenario: The rendered claim fits the backup catalog

- **WHEN** `backup_provider` is rendered as instance `backup-provider` in namespace `default`
- **THEN** the output is one `TransformerRegistration` named `default.backup-provider`
- **AND** its `spec.catalog` and `spec.version` equal the backup catalog fixture's `ModulePath` and `Version`
- **AND** the catalog resolves from the registry with the claim's own `spec.catalog` and bare `spec.version`, as acceptance resolves it
- **AND** its `spec.provides` equals the provider contracts derived from that catalog build
- **AND** its `spec.providerRef` names `default/backup-provider`

### Requirement: The claim's catalog, version and provides are authored literals
As in requirement "backup_provider renders one TransformerRegistration": The claim's `catalog`, `version` and `provides` SHALL be authored literals, a deliberate deviation from enhancement 0015 ("Every field of the claim is derived or stamped; the only authored fact is the module's catalog dependency version"): a module fixture cannot import a catalog fixture version that is new in the same pull request, so a derived claim could not land with a `backup` catalog bump.

#### Scenario: The claim fields are authored literals

- **WHEN** `backup_provider` is rendered
- **THEN** `spec.catalog`, `spec.version` and `spec.provides` of the claim are authored literals

### Requirement: backup_provider is applied under a ServiceAccount that may write transformerregistrations
As in requirement "backup_provider renders one TransformerRegistration": Its `moduleinstance.yaml` SHALL apply it under a ServiceAccount bound to a ClusterRole that may write `transformerregistrations`, since only a platform-team identity may register a provider (enhancement 0015, "A transformer registration is a cluster-scoped CR, gated by the RBAC the controller already enforces").

#### Scenario: The provider instance is applied under a ServiceAccount that may write transformerregistrations

- **WHEN** the `moduleinstance.yaml` of `backup_provider` is read
- **THEN** it applies the instance under a ServiceAccount bound to a ClusterRole that may write `transformerregistrations`

### Requirement: backup_consumer attaches the backup trait
`backup_consumer` SHALL attach opm's backup trait, and its render SHALL be refused naming that contract while no provider of it is in the platform's registry (enhancement 0010, "An unresolved primitive demand is an error; a trait's optionality is stated by its catalog and overridden by the attachment"), and SHALL succeed with the backup catalog in the registry, rendering through the backup catalog's transformer.

#### Scenario: The consumer is refused without a provider

- **WHEN** `backup_consumer` is rendered against a platform that subscribes only `opmodel.dev/catalogs/opm@v4`
- **THEN** the render is refused, naming `opmodel.dev/catalogs/opm/traits/backup@v1alpha1`

#### Scenario: The consumer renders through the provider

- **WHEN** `backup_consumer` is rendered against a platform whose registry also holds the backup catalog fixture
- **THEN** it renders one ConfigMap, produced by the backup catalog's transformer
- **AND** the instance's required contracts include `opmodel.dev/catalogs/opm/traits/backup@v1alpha1`

### Requirement: A registry-backed integration spec checks the registration, the refusal and the render
As in requirements "backup_provider renders one TransformerRegistration" and "backup_consumer attaches the backup trait": So a cluster running a controller that resolves the catalog accepts it (enhancement 0015, "A transformer registration is a cluster-scoped CR, gated by the RBAC the controller already enforces" and "Every field of the claim is derived or stamped; the only authored fact is the module's catalog dependency version"). A registry-backed integration spec SHALL check all three.

#### Scenario: A cluster accepts and activates the claim

- **WHEN** both `moduleinstance.yaml` files are applied to a cluster whose controller resolves the published fixtures and accepts a bare SemVer in `spec.version`
- **THEN** the claim `default.backup-provider` reports `accepted: true` and, once `backup-provider` is Ready, `active: true`
- **AND** `backup-consumer` reaches Ready
