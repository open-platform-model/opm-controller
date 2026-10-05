## MODIFIED Requirements

### Requirement: Code owners on the release machinery
`.github/CODEOWNERS` SHALL name the owner (`@emil-jacero`) as the only code owner for `/.github/`, `/.tasks/`, `/Taskfile*.yml`, `/release-please-config.json`, `/.release-please-manifest.json`, `/.cascade-frozen`, `/hack/` (scripts that run in jobs holding write tokens or the release App key, `hack/controller-module/` included), `/.opm-cli-version` (the opm binary the publishing jobs install), and `/.opm-docs-version` and `/docs-kit.cue` (the signed docs publish).

#### Scenario: Workflow edit needs a code owner
- **WHEN** a pull request edits a file under `.github/workflows/`
- **THEN** GitHub requests review from the code owners listed for `/.github/`

#### Scenario: One code owner per path
- **WHEN** `.github/CODEOWNERS` is read
- **THEN** every path line names `@emil-jacero` and no other owner

## REMOVED Requirements

### Requirement: No Actions cache in a publishing job

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-019); this requirement is restated under the controller's name.

**Migration**: See the requirement "A publishing job uses no Actions cache" in this capability.

### Requirement: Bot-head code never holds a write token or OIDC

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-019); this requirement is restated under the controller's name.

**Migration**: See the requirement "Code on a bot head never holds a write token or OIDC" in this capability.

### Requirement: Operator module bot PRs commit plain files only

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-019); this requirement is restated under the controller's name.

**Migration**: See the requirement "Controller module bot PRs commit plain files only" in this capability.

## ADDED Requirements

### Requirement: A publishing job uses no Actions cache
A job that publishes, signs or attests anything (a container image, a CUE module, a release asset) SHALL NOT restore or save the GitHub Actions cache: no `type=gha` cache in a Docker build, `cache: false` on `actions/setup-go`, and no `actions/cache`. Jobs that publish nothing MAY keep their caches, except in a workflow that also has a publishing job: `test-e2e.yml` and `image-pr.yml` use no cache at all, so the cascade wiring check's per-workflow no-cache rule can cover them.

#### Scenario: Release image built without the Actions cache
- **WHEN** `release.yml`'s `image-release` job builds the release image
- **THEN** its build step has no `cache-from` or `cache-to`

#### Scenario: E2E workflow holds no cache
- **WHEN** either job of `test-e2e.yml` sets up Go
- **THEN** `actions/setup-go` runs with `cache: false`

#### Scenario: Module publishing jobs skip the Go cache
- **WHEN** `release.yml`'s `publish-examples`, `publish-fixtures.yml` or `test-e2e.yml`'s `publish-fixtures` job sets up Go
- **THEN** `actions/setup-go` runs with `cache: false`

#### Scenario: Controller module release jobs skip the Go cache
- **WHEN** `release.yml`'s `module-identity-advance` (which later holds the release App key), `module-publish` or `module-manifest` job sets up Go
- **THEN** `actions/setup-go` runs with `cache: false`, and `module-image.yml` and `module-deps.yml` use no cache in any job

### Requirement: Code on a bot head never holds a write token or OIDC
A job triggered by `pull_request` that holds any write scope or `id-token: write` SHALL NOT run for a head named `deps/cascade` or starting with `release-please--`, `module/` or `dependabot/`. The code on those heads comes from the release cascade, release-please, the controller module's bot (`module/deps`, `module/controller-image`, pushed with the release App's token) or Dependabot, and no human has reviewed it yet; a Dependabot `github-actions` bump runs a new upstream action release in the job itself. A job that runs such code SHALL hold only read scopes.

#### Scenario: Cascade PR image build skipped
- **WHEN** a pull request from `deps/cascade` opens or updates
- **THEN** `image-pr.yml`'s job is skipped

#### Scenario: Cascade PR e2e runs read-only
- **WHEN** a pull request from `deps/cascade` or `release-please--*` opens or updates
- **THEN** `test-e2e.yml`'s `publish-fixtures` job is skipped, and the `test-e2e` job runs the suite with `contents: read` and `packages: read`, without the pre-release pins or a GHCR credential

#### Scenario: E2E without a publish still runs the podinfo spec
- **WHEN** `test-e2e.yml`'s `publish-fixtures` job is skipped (a fork, or a `deps/cascade`, `release-please--*`, `module/*` or `dependabot/*` head)
- **THEN** the `test-e2e` job seeds a job-local registry from the tree, connects it to the kind network, sets `LOCAL_REGISTRY` so the controller resolves the fixtures from it, and the podinfo and redis specs run instead of skipping

#### Scenario: Dependabot and controller module bot PRs hold no write token
- **WHEN** a pull request from a `dependabot/*` or `module/*` head opens or updates
- **THEN** `image-pr.yml`'s job and `test-e2e.yml`'s `publish-fixtures` job are skipped

#### Scenario: Human PR e2e publishes in its own job
- **WHEN** a same-repo pull request from any other head opens or updates
- **THEN** `publish-fixtures` publishes the pre-release fixtures holding `packages: write`, and the `test-e2e` job pins them and runs the suite holding only read scopes, with a docker config built from its own read-only token

### Requirement: Controller module bot PRs commit plain files only
The `publish` jobs of `module-image.yml` and `module-deps.yml` SHALL refuse a handed-over change tarball holding any member that is not a plain file or a directory, before extracting it. `hack/controller-module/bot-pr.sh` SHALL refuse a changed path that is a symlink or a hard link, and a staged change whose mode is anything but `100644` (an executable bit, a symlink or a gitlink), before it commits or pushes.

#### Scenario: A link member in the tarball is refused
- **WHEN** the module change tarball holds a symlink or hard link member under `modules/opm_controller/`
- **THEN** the publish job fails before extracting it, and no token-bearing step reads the link's target

#### Scenario: bot-pr.sh refuses a symlink, a hard link or a mode change
- **WHEN** the working tree under `modules/opm_controller/` holds a symlink, a hard link, or a file whose mode changed
- **THEN** `bot-pr.sh` exits non-zero naming the path or the mode change, and pushes nothing
