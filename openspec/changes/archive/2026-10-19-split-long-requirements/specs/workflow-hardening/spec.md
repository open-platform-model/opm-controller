## REMOVED Requirements

### Requirement: Release App key read only in the release Environment
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "Every job that reads the release App key runs in the release Environment" and "A reusable-workflow call does not pass the release App key".

### Requirement: Bot-head code never holds a write token or OIDC
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "A write-scope pull_request job does not run for a bot head" and "A job that runs bot-head code holds only read scopes".

## ADDED Requirements

### Requirement: Every job that reads the release App key runs in the release Environment
Every job that reads `secrets.RELEASE_APP_PRIVATE_KEY` SHALL declare `environment: release`, and no other job SHALL declare that Environment. Today those jobs are `release.yml`'s `release-please` and `module-identity-advance` (on a push to `main`), and the `publish` jobs of `module-image.yml` and `module-deps.yml` (which run only on `refs/heads/main`).

#### Scenario: Every key reader runs in the release Environment
- **WHEN** `.github/workflows/` is searched for `RELEASE_APP_PRIVATE_KEY`
- **THEN** every job that names it declares `environment: release`, and no other job declares `release`

#### Scenario: Key not yet moved
- **WHEN** Environment `release` has no `RELEASE_APP_PRIVATE_KEY` and the org secret still exists
- **THEN** the jobs of a top-level run (`release-please`, `module-identity-advance`, and `module-image.yml` or `module-deps.yml` run by dispatch) read the org secret and run as before, while `module-image.yml`'s `publish` job called from `module-image-pr` gets no key until the owner moves it into the Environment (a called workflow sees no org secret that is not passed), and a dispatch of `module-image.yml` recovers that run

### Requirement: A reusable-workflow call does not pass the release App key
As in requirement "Every job that reads the release App key runs in the release Environment": A reusable-workflow call SHALL NOT pass the key through `secrets:` (nor `secrets: inherit`): a calling job cannot declare an Environment, so the called job that declares `environment: release` reads the Environment secret itself, and the called workflow declares no `workflow_call` secret for it.

#### Scenario: Module image PR call passes no key
- **WHEN** `release.yml`'s `module-image-pr` job calls `module-image.yml`
- **THEN** it has no `secrets:` key, `module-image.yml` declares no `workflow_call` secret, and its `publish` job declares `environment: release`

### Requirement: A write-scope pull_request job does not run for a bot head
A job triggered by `pull_request` that holds any write scope or `id-token: write` SHALL NOT run for a head named `deps/cascade` or starting with `release-please--`, `module/` or `dependabot/`.

#### Scenario: Cascade PR image build skipped
- **WHEN** a pull request from `deps/cascade` opens or updates
- **THEN** `image-pr.yml`'s job is skipped

#### Scenario: Dependabot and operator module bot PRs hold no write token
- **WHEN** a pull request from a `dependabot/*` or `module/*` head opens or updates
- **THEN** `image-pr.yml`'s job and `test-e2e.yml`'s `publish-fixtures` job are skipped

#### Scenario: Human PR e2e publishes in its own job
- **WHEN** a same-repo pull request from any other head opens or updates
- **THEN** `publish-fixtures` publishes the pre-release fixtures holding `packages: write`, and the `test-e2e` job pins them and runs the suite holding only read scopes, with a docker config built from its own read-only token

### Requirement: A job that runs bot-head code holds only read scopes
As in requirement "A write-scope pull_request job does not run for a bot head": The code on those heads comes from the release cascade, release-please, the operator module's bot (`module/deps`, `module/operator-image`, pushed with the release App's token) or Dependabot, and no human has reviewed it yet; a Dependabot `github-actions` bump runs a new upstream action release in the job itself. A job that runs such code SHALL hold only read scopes.

#### Scenario: Cascade PR e2e runs read-only
- **WHEN** a pull request from `deps/cascade` or `release-please--*` opens or updates
- **THEN** `test-e2e.yml`'s `publish-fixtures` job is skipped, and the `test-e2e` job runs the suite with `contents: read` and `packages: read`, without the pre-release pins or a GHCR credential

#### Scenario: E2E without a publish still runs the podinfo spec
- **WHEN** `test-e2e.yml`'s `publish-fixtures` job is skipped (a fork, or a `deps/cascade`, `release-please--*`, `module/*` or `dependabot/*` head)
- **THEN** the `test-e2e` job seeds a job-local registry from the tree, connects it to the kind network, sets `LOCAL_REGISTRY` so the controller resolves the fixtures from it, and the podinfo and redis specs run instead of skipping
