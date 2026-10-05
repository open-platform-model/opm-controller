## MODIFIED Requirements

### Requirement: The cascade references are pinned to one .github main SHA and checked in the Lint job
Every reference to the cascade code SHALL name one full 40-character SHA of a commit on `open-platform-model/.github` `main`, never a branch or a tag, followed by the comment `# .github main`: the `cascade-notify` step in `release.yml`, the `cascade-receive.yml` call and the `cascade-publish` step in `deps-cascade.yml`, the `cascade-gates.yml` call, the `ref:` of the `open-platform-model/.github` checkout in `cascade-task.yml`, and the `ref:` of the resolver checkout in `module-deps.yml` (the controller module's own cascade PR, see `deps-cascade`). All six SHALL carry the same SHA, which moves only through a `ci(deps): pin the cascade to .github <sha7>` pull request. `cascade-task.yml` SHALL fail, not skip, when the resolver is missing at the pinned commit. `.tasks/cascade/wiring-check.sh` SHALL be byte-identical to `.github/scripts/cascade/wiring-check.sh` at the pinned SHA and SHALL move only with the pin; this repo's values SHALL live in `.tasks/cascade/wiring-check.yaml` (`receiver: true`, `env-allow` `REGISTRY`, `IMAGE_NAME` and `CUE_VERSION`, `ci` `lint.yml` job `lint`, the notify `needs`, `if` and `tag` below, `labels-managed: false`, `extra-references` naming the `module-deps.yml` resolver, and `publish-workflows` naming every workflow that publishes: `release.yml`, `publish-fixtures.yml`, `docs.yml`, `image-pr.yml`, `test-e2e.yml`, `module-image.yml` and `module-deps.yml`). The required `Lint` job SHALL run it on every pull request as the step "Verify the cascade wiring", exactly `bash .tasks/cascade/wiring-check.sh --pin-on-main` with only `GH_TOKEN: ${{ github.token }}` in its `env`, preceded only by SHA-pinned actions from other repos; `task cascade:wiring:check` runs it offline. It SHALL fail when: a key-holding job (`notify-downstream`, `publish`) differs from its contract shape in its keys, step keys, `with` keys, Environment, permissions, step count, action reference, key inputs or a `runs-on` other than `ubuntu-latest`; the `notify-downstream` job's `needs` differs from `[release-please, publish-release]`, or its `if` from `needs.release-please.outputs.release_created == 'true' && vars.CASCADE_NOTIFY != 'off'` (contract §4.6 prints `releases_created`, which a module-only release also sets), or its `tag` input or `timeout-minutes` (20) from contract §4.6; the `publish` job's `needs` differs from `cascade`, its `if:` or `gates-only` input from the receiver caller shape, or its `timeout-minutes` from 15; `release.yml`'s workflow-level `env` is not a map or has a key other than `REGISTRY`, `IMAGE_NAME` and `CUE_VERSION`; `deps-cascade.yml`'s top-level keys, concurrency group, `publish` condition or either `dry-run` value differ from the contract; any other workflow value reads the key or any other job declares `environment: cascade`; a job reads `RELEASE_APP_PRIVATE_KEY` without declaring exactly `environment: release`, or a job declares `release` without reading it; a workflow listed in `publish-workflows` uses an Actions cache; a `.github` checkout is not `actions/checkout@<SHA>` with exactly `repository`, `ref`, `path` and `persist-credentials: false`; the `Lint` workflow or job `env` names anything outside `CUE_*`, `OPM_*`, `REGISTRY` and `IMAGE_NAME`, or the job has a `container` or `services`; a call into `.github` has a `secrets` key; or the set of `.github` references, their SHA or their comment differs. With `--pin-on-main` it SHALL also fail when the SHA is not on `.github` `main` or the running copy differs from the file at the SHA. It guards against mistakes, not a deliberate edit: code-owner review is the guard against that, and the `main` ruleset does not require review while OPM is in beta (no approvals, workspace PR 29) (Phase 3 wiring contract (version 3.1) §2.4, §10.1 items 2, 3 and 6; the supervisor's addendum; `.github` README at `0f9c6ac`, "The wiring check").

#### Scenario: One reference moves alone
- **WHEN** a pull request changes only the `cascade-gates.yml` call to another SHA
- **THEN** the "Verify the cascade wiring" step fails with `one .github SHA` naming both SHAs, and the `Lint` check is red

#### Scenario: A branch ref
- **WHEN** a pull request names `cascade-notify@main`
- **THEN** the wiring check fails on the notify `uses` and on `one .github SHA`

#### Scenario: A startup-code variable in release.yml
- **WHEN** a pull request adds `BASH_ENV` or any key other than `REGISTRY`, `IMAGE_NAME` and `CUE_VERSION` to `release.yml`'s workflow-level `env`
- **THEN** the wiring check fails naming the key

#### Scenario: A self-hosted runner for a key job
- **WHEN** a pull request sets `runs-on: self-hosted` on `notify-downstream`
- **THEN** the wiring check fails naming `release.yml:notify-downstream runs-on`

#### Scenario: The resolver is missing at the pin
- **WHEN** `cascade-task.yml` runs and the pinned `.github` commit has no `cascade-resolve.sh`
- **THEN** the run fails with `no cascade resolver at the pinned .github commit` instead of skipping S5

#### Scenario: The wired tree passes
- **WHEN** the `Lint` job runs on a pull request that leaves the cascade wiring as the contract shapes it
- **THEN** the step prints `cascade wiring: ok, .github <SHA> (.github main)`

#### Scenario: The module's resolver checkout follows the pin
- **WHEN** a pull request moves the five other `.github` references to a new SHA and leaves `module-deps.yml`'s resolver checkout behind
- **THEN** the wiring check fails on `one .github SHA`

#### Scenario: A copy left behind at a pin bump
- **WHEN** a pull request moves every reference to a new SHA but keeps the old `wiring-check.sh`
- **THEN** the "Verify the cascade wiring" step fails with `differs from` naming the file at the new SHA

#### Scenario: A run step before the wiring step
- **WHEN** a pull request adds a `run:` step to the `Lint` job above "Verify the cascade wiring"
- **THEN** the wiring check fails naming the step that is not a pinned action

#### Scenario: A cache in a publishing workflow
- **WHEN** a pull request drops `cache: false` from a `setup-go` step in `test-e2e.yml`
- **THEN** the wiring check fails naming `test-e2e.yml cache use`

## REMOVED Requirements

### Requirement: The receiver runs the shared receive workflow on dispatch, sweep and demand

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-019); this requirement is restated under the controller's name.

**Migration**: See the requirement "The receiver runs the shared receive workflow on a dispatch, the sweep and demand" in this capability.

## ADDED Requirements

### Requirement: The receiver runs the shared receive workflow on a dispatch, the sweep and demand
The repo SHALL have `.github/workflows/deps-cascade.yml`, named `Deps cascade`, triggered by `repository_dispatch` of type `upstream-released`, by the schedule `47 5 * * *`, and by `workflow_dispatch` with the boolean inputs `dry_run` and `gates_only` (both default `false`). Its top-level keys SHALL be exactly `name`, `on`, `permissions`, `concurrency` and `jobs`, with `permissions: {}`. Its job `cascade` SHALL call `open-platform-model/.github/.github/workflows/cascade-receive.yml` at the repo's pinned `.github` SHA, granting `contents: read`, `pull-requests: read` and `statuses: write` and nothing else, passing no `secrets`, and passing `setup-go: true`, the dry-run expression `${{ inputs.dry_run == true || vars.CASCADE_DRY_RUN != 'false' }}`, `gates-only`, and `g2-mode` and `g3-mode` from the repo variables `CASCADE_G2_MODE` and `CASCADE_G3_MODE`, each defaulting to `warn`. Payload validation, the task run and the plan belong to the called workflow, whose jobs hold no secret (Phase 3 wiring contract (version 3.1) §5, §5.1, §5.2, §6).

#### Scenario: A library release reaches the controller
- **WHEN** library publishes `v1.0.0-beta.5`, its notify job dispatches `upstream-released` to opm-controller, and `CASCADE_DRY_RUN` is `false`
- **THEN** a `Deps cascade` run starts on `main`, runs `task -x deps:cascade`, and the `deps/cascade` PR names library `v1.0.0-beta.5`

#### Scenario: Daily sweep
- **WHEN** no dispatch arrived and the schedule fires at 05:47 UTC
- **THEN** a `Deps cascade` run starts with no payload and behaves as a sweep

#### Scenario: Workflow grants nothing at the top
- **WHEN** `deps-cascade.yml` is read
- **THEN** the top-level `permissions` is empty, the `cascade` job grants only `contents: read`, `pull-requests: read` and `statuses: write`, and no job that calls into `open-platform-model/.github` has a `secrets` key
