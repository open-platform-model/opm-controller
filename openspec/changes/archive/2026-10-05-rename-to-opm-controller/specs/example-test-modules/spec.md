## MODIFIED Requirements

### Requirement: Public example module path

All example test modules SHALL be authored under the CUE module path `testing.opmodel.dev/modules/controller/<module>@v0`, where `<module>` is the module's short name. Each module's `cue.mod/module.cue` `module:` field and its `#Module.metadata.modulePath` SHALL agree with this path, and SHALL be published to `ghcr.io/open-platform-model` so the module resolves under the canonical registry mapping, which routes both `opmodel.dev` (core, catalogs) and `testing.opmodel.dev` (these fixtures) to GHCR.

Fixtures SHALL NOT be authored under `opmodel.dev/*`. CUE resolves modules by longest-prefix match on the module path, so a fixture in the production namespace forces that entire prefix — core and the catalogs included — onto whatever registry serves the fixture. The publish gates enforce this independently: a nested path under `opmodel.dev` is refused.

#### Scenario: New module declares public path

- **WHEN** a new example module is added
- **THEN** its `cue.mod/module.cue` declares `module: "testing.opmodel.dev/modules/controller/podinfo@v0"` and `metadata.modulePath` resolves to the same value

#### Scenario: Consumer resolves without extra config

- **WHEN** a user with the canonical `CUE_REGISTRY` mapping resolves `testing.opmodel.dev/modules/controller/podinfo@v0`
- **THEN** it resolves from `ghcr.io/open-platform-model` with no additional registry configuration

#### Scenario: Production namespace is refused

- **WHEN** a fixture declares a path under `opmodel.dev/modules/test/`
- **THEN** `opm module publish` refuses it as not fitting the namespace the domain publishes

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

### Requirement: ModulePackage fixture parity for example modules

Each example module of the workload fleet (`hello`, `hello_web`, `podinfo`, `redis`) SHALL have a sibling modulepackage fixture declaring `testing.opmodel.dev/releases/controller/<module>@v0`, an `instance.cue` that imports the published module, and a `cue.mod/module.cue` that pins the same `testing.opmodel.dev/modules/controller/<module>@v0` version the module declares. The backup fixture set (`backup_provider`, `backup_consumer`) has no modulepackage fixture: it exists to exercise a registration, and `backup_consumer` does not render on a platform without an active provider.

#### Scenario: Each module has a modulepackage fixture

- **WHEN** the modulepackage fixtures are enumerated
- **THEN** there SHALL be one per workload fleet module, declaring `testing.opmodel.dev/releases/controller/<module>@v0`

#### Scenario: instance.cue imports and embeds the published module

- **WHEN** a modulepackage fixture is loaded
- **THEN** it embeds `core.#ModuleInstance`, imports `testing.opmodel.dev/modules/controller/<module>@v0`, and sets `#module` to the imported module

#### Scenario: ModulePackage CR references its OCIRepository

- **WHEN** a modulepackage fixture's `ModulePackage` is inspected
- **THEN** its `spec.sourceRef` names the `OCIRepository` declared in the sibling `ocirepository.yaml`, whose `url` ends in `testing.opmodel.dev/releases/controller/<module>`

### Requirement: hello_web ready-to-apply ModuleInstance

The `hello_web` example SHALL ship a `ModuleInstance` manifest referencing the `testing.opmodel.dev/modules/controller/hello_web@v0` path and a concrete version, so all four fleet members are applyable.

#### Scenario: hello_web manifest references public module

- **WHEN** a user inspects `hello_web`'s `moduleinstance.yaml`
- **THEN** its `spec.module.path` is `testing.opmodel.dev/modules/controller/hello_web@v0` with an explicit `spec.module.version`

### Requirement: A test catalog is published beside the module fleet

The repo SHALL carry a catalog fixture under `test/fixtures/catalogs/provider`, declaring the CUE module path `testing.opmodel.dev/catalogs/controller/provider@v0`. It SHALL be the smallest valid core `#Catalog`: one resource and one transformer that requires it, each with an FQN under the catalog's own registry path, so a platform can subscribe it beside `opmodel.dev/catalogs/opm@v4` without a key collision. It SHALL carry an `identity/identity.cue` package as the single source of its module path and version, with `Version` a plain SemVer literal.

The fixture SHALL be published through `hack/fixtures.sh`, the same flow as the module fleet, with `opm catalog publish` in place of `opm module publish`: `check` runs every catalog publish gate as a dry run and enforces changed-implies-bumped, `seed` publishes it into the job-local registry before the registry-backed specs run in `test.yml`, and `publish` puts it on GHCR when `publish-fixtures.yml` runs after a merge that touches `test/fixtures/catalogs/`. Registry-backed specs that need a resolvable catalog other than `opmodel.dev/catalogs/opm@v4` SHALL use this fixture and SHALL read its coordinate from its identity package (`fixtures.MustCatalog`), never a literal.

#### Scenario: The catalog passes the catalog publish gates

- **WHEN** `opm catalog publish --dry-run test/fixtures/catalogs/provider` runs
- **THEN** the plan resolves with no refusals, checking two members through the member FQN gate, and derives repository `testing.opmodel.dev/catalogs/controller/provider` and the tag from the identity package

#### Scenario: PR CI seeds the catalog from the tree

- **WHEN** `test.yml` runs on a pull request
- **THEN** the seed step publishes `testing.opmodel.dev/catalogs/controller/provider` at its declared version into the job-local registry before any registry-backed spec runs

#### Scenario: A changed catalog must carry a new version

- **WHEN** a pull request changes `test/fixtures/catalogs/provider` but keeps a version GHCR already holds
- **THEN** `task examples:check` fails, naming the fixture and `opm catalog version set` as the fix

#### Scenario: The claim specs build against the test catalog

- **WHEN** the claim-driven regeneration specs activate a claim for a catalog the Platform does not subscribe
- **THEN** that catalog is `testing.opmodel.dev/catalogs/controller/provider@v0` at the version its identity package declares, and the regenerated platform pins, imports and builds it

### Requirement: A provider catalog fixture implements a provider-fulfilled contract

The repo SHALL carry a catalog fixture under `test/fixtures/catalogs/backup`, declaring `testing.opmodel.dev/catalogs/controller/backup@v0`, whose transformers implement opm's provider-fulfilled backup trait (`opmodel.dev/catalogs/opm/traits/backup@v1alpha1`) and nothing else, so the provider contracts derived from it are exactly that trait. It SHALL render the trait without any CRD, publish through `hack/fixtures.sh` like the provider catalog fixture, and pin core and the opm catalog no newer than the build the registry-backed specs subscribe. `task deps:cascade` SHALL NOT be required to move it. While the library keeps its deprecated provider-set fold, its core pin SHALL stay older than the library's `ProvidesSince`, since the controller's old-catalog acceptance spec relies on it.

#### Scenario: The catalog passes the catalog publish gates

- **WHEN** `opm catalog publish --dry-run test/fixtures/catalogs/backup` runs against GHCR
- **THEN** the plan resolves with no refusals, checking one member through the member FQN gate

#### Scenario: The catalog provides exactly the backup trait

- **WHEN** the backup catalog fixture is acquired from the registry at its declared build and its provider contracts are derived
- **THEN** they are exactly `opmodel.dev/catalogs/opm/traits/backup@v1alpha1`

### Requirement: The backup fixture set exercises a registration end to end

`backup_provider` SHALL render exactly one `TransformerRegistration`, through opm's `transformer-registration` contract, that names the backup catalog fixture at the version its identity package declares and lists exactly the provider contracts that catalog implements, so a cluster running a controller that resolves the catalog accepts it (enhancement 0015, "A transformer registration is a cluster-scoped CR, gated by the RBAC the operator already enforces" and "Every field of the claim is derived or stamped; the only authored fact is the module's catalog dependency version"). The claim's `catalog`, `version` and `provides` are authored literals, a deliberate deviation from enhancement 0015 ("Every field of the claim is derived or stamped; the only authored fact is the module's catalog dependency version"): a module fixture cannot import a catalog fixture version that is new in the same pull request, so a derived claim could not land with a `backup` catalog bump. Its `moduleinstance.yaml` SHALL apply it under a ServiceAccount bound to a ClusterRole that may write `transformerregistrations`, since only a platform-team identity may register a provider (enhancement 0015, "A transformer registration is a cluster-scoped CR, gated by the RBAC the operator already enforces"). `backup_consumer` SHALL attach opm's backup trait, and its render SHALL be refused naming that contract while no provider of it is in the platform's registry (enhancement 0010, "An unresolved primitive demand is an error; a trait's optionality is stated by its catalog and overridden by the attachment"), and SHALL succeed with the backup catalog in the registry, rendering through the backup catalog's transformer. A registry-backed integration spec SHALL check all three.

#### Scenario: The rendered claim fits the backup catalog

- **WHEN** `backup_provider` is rendered as instance `backup-provider` in namespace `default`
- **THEN** the output is one `TransformerRegistration` named `default.backup-provider`
- **AND** its `spec.catalog` and `spec.version` equal the backup catalog fixture's `ModulePath` and `Version`
- **AND** the catalog resolves from the registry with the claim's own `spec.catalog` and bare `spec.version`, as acceptance resolves it
- **AND** its `spec.provides` equals the provider contracts derived from that catalog build
- **AND** its `spec.providerRef` names `default/backup-provider`

#### Scenario: The consumer is refused without a provider

- **WHEN** `backup_consumer` is rendered against a platform that subscribes only `opmodel.dev/catalogs/opm@v4`
- **THEN** the render is refused, naming `opmodel.dev/catalogs/opm/traits/backup@v1alpha1`

#### Scenario: The consumer renders through the provider

- **WHEN** `backup_consumer` is rendered against a platform whose registry also holds the backup catalog fixture
- **THEN** it renders one ConfigMap, produced by the backup catalog's transformer
- **AND** the instance's required contracts include `opmodel.dev/catalogs/opm/traits/backup@v1alpha1`

#### Scenario: A cluster accepts and activates the claim

- **WHEN** both `moduleinstance.yaml` files are applied to a cluster whose controller resolves the published fixtures and accepts a bare SemVer in `spec.version`
- **THEN** the claim `default.backup-provider` reports `accepted: true` and, once `backup-provider` is Ready, `active: true`
- **AND** `backup-consumer` reaches Ready
