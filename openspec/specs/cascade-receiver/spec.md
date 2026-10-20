# cascade-receiver Specification

## Purpose
The opm-operator side of the release cascade wiring: the `deps-cascade.yml` receiver that runs `task -x deps:cascade` through the shared receive workflow on every upstream release, a daily sweep and by hand, and publishes through a caller-owned job; the `cascade-gates.yml` caller that keeps the G2 and G3 commit statuses on every pull request; and the one-SHA pin of the cascade code, checked on every pull request.

## Requirements

### Requirement: Gates-only and branch runs never displace a pending real run
The receiver SHALL declare `cancel-in-progress: false` and the concurrency group `deps-cascade` for real runs on `main` (dispatch, schedule, manual), `deps-cascade-gates` for gates-only runs on `main`, and `deps-cascade-<ref>` for runs from any other ref (Phase 3 wiring contract (version 3.1) §5; workspace RELEASING.md, "Concurrency").

#### Scenario: Burst of releases
- **WHEN** three dispatches arrive while a real run is active
- **THEN** one run stays active and one is pending in `deps-cascade`, and the pending run resolves the newest versions itself

#### Scenario: Gates-only run during a pending real run
- **WHEN** a release PR head moves while a real run is pending
- **THEN** the gates-only run queues in `deps-cascade-gates`, and the pending real run stays pending

### Requirement: The Deps cascade workflow
The repo SHALL have `.github/workflows/deps-cascade.yml`, named `Deps cascade`, triggered by `repository_dispatch` of type `upstream-released`, by the schedule `47 5 * * *`, and by `workflow_dispatch` with the boolean inputs `dry_run` and `gates_only` (both default `false`). Its top-level keys SHALL be exactly `name`, `on`, `permissions`, `concurrency` and `jobs`, with `permissions: {}`.

#### Scenario: A library release reaches the operator
- **WHEN** library publishes `v1.0.0-beta.5`, its notify job dispatches `upstream-released` to opm-operator, and `CASCADE_DRY_RUN` is `false`
- **THEN** a `Deps cascade` run starts on `main`, runs `task -x deps:cascade`, and the `deps/cascade` PR names library `v1.0.0-beta.5`

#### Scenario: Daily sweep
- **WHEN** no dispatch arrived and the schedule fires at 05:47 UTC
- **THEN** a `Deps cascade` run starts with no payload and behaves as a sweep

### Requirement: The cascade job's call and grants
As in requirement "The Deps cascade workflow": Its job `cascade` SHALL call `open-platform-model/.github/.github/workflows/cascade-receive.yml` at the repo's pinned `.github` SHA, granting `contents: read`, `pull-requests: read` and `statuses: write` and nothing else, passing no `secrets`. Payload validation, the task run and the plan belong to the called workflow, whose jobs hold no secret (Phase 3 wiring contract (version 3.1) §5, §5.1, §5.2, §6).
(Requirement: The cascade job's inputs)

#### Scenario: Workflow grants nothing at the top
- **WHEN** `deps-cascade.yml` is read
- **THEN** the top-level `permissions` is empty, the `cascade` job grants only `contents: read`, `pull-requests: read` and `statuses: write`, and no job that calls into `open-platform-model/.github` has a `secrets` key

### Requirement: The cascade job's inputs
As in requirement "The Deps cascade workflow": Its job `cascade` SHALL call `open-platform-model/.github/.github/workflows/cascade-receive.yml` at the repo's pinned `.github` SHA, passing `setup-go: true`, the dry-run expression `${{ inputs.dry_run == true || vars.CASCADE_DRY_RUN != 'false' }}`, `gates-only`, and `g2-mode` and `g3-mode` from the repo variables `CASCADE_G2_MODE` and `CASCADE_G3_MODE`, each defaulting to `warn`.
(Requirement: The cascade job's call and grants)

#### Scenario: The cascade job passes the dry-run expression and the modes
- **WHEN** the `cascade` job calls the shared receive workflow
- **THEN** it passes `setup-go: true`, the dry-run expression, `gates-only`, and `g2-mode` and `g3-mode` from the repo variables `CASCADE_G2_MODE` and `CASCADE_G3_MODE`, each defaulting to `warn`

### Requirement: The publish job's shape
`deps-cascade.yml` SHALL have a job `publish` that needs `cascade`, runs on `ubuntu-latest`, declares `environment: cascade`, grants `contents: read` and `pull-requests: read` and nothing else.
(Requirement: The publish job has exactly one step, the pinned cascade-publish action)

#### Scenario: The plan becomes a push
- **WHEN** `CASCADE_DRY_RUN` is `false`, the run is on `main`, `compute` succeeded and planned `push`
- **THEN** `publish` runs in the `cascade` Environment, and the `cascade-publish` action verifies the plan, mints the App token and pushes `deps/cascade`

### Requirement: The publish job has exactly one step, the pinned cascade-publish action
`deps-cascade.yml` SHALL have a job `publish` that has exactly one step: the action `open-platform-model/.github/.github/actions/cascade-publish` at the repo's pinned `.github` SHA, with the inputs `dry-run` (the same expression as the `cascade` job), `gates-only: ${{ inputs.gates_only == true }}`, `labels-managed: false`, `client-id: ${{ vars.CASCADE_APP_CLIENT_ID }}` and `private-key: ${{ secrets.CASCADE_APP_PRIVATE_KEY }}`, and no other input.
(Requirement: The publish job's shape)

#### Scenario: A gates-only run reaches the action
- **WHEN** the `publish` job's `if:` were bypassed on a run dispatched with `gates_only=true`
- **THEN** `cascade-publish` receives `gates-only: true` and fails before it mints the App token

### Requirement: Only a caller-owned job reads the cascade App key
As in requirement "The publish job has exactly one step, the pinned cascade-publish action": The key SHALL be read only in the caller-owned `notify-downstream` or `publish` job, which declares `environment: cascade` and passes `secrets.CASCADE_APP_PRIVATE_KEY` only as the `private-key` input of the SHA-pinned cascade action.
(Requirement: The key-reading job has no checkout, run, env, container or services)

#### Scenario: Another job reads the key
- **WHEN** a pull request adds `secrets.CASCADE_APP_PRIVATE_KEY` to any other job or declares `environment: cascade` on any other job
- **THEN** the wiring check fails naming the unexpected reader or Environment job

### Requirement: The key-reading job has no checkout, run, env, container or services
As in requirement "Only a caller-owned job reads the cascade App key": that job SHALL have no checkout or `run:` of its own, and no `env:`, `container:` or `services:` (the action checks out the repo but never runs it); no reusable call SHALL pass `secrets:` or `secrets: inherit` (Phase 3 wiring contract (version 3.1) §2.2, §5.2, §10.1 item 9; `.github` README at `2b66da5`, "Receiver caller").

#### Scenario: An env on the publish job is refused
- **WHEN** a pull request adds `env: {BASH_ENV: repo/x.sh}` to the `publish` job
- **THEN** the `Lint` job's "Verify the cascade wiring" step fails naming `deps-cascade.yml:publish keys`

### Requirement: The receiver plans a dry run unless CASCADE_DRY_RUN is exactly false
The receiver SHALL plan a dry run, and SHALL NOT start the `publish` job, unless the repo variable `CASCADE_DRY_RUN` is exactly `false`, the `dry_run` input is not `true`, the `gates_only` input is not `true`, and the run is on `main`; the `publish` job's `if:` SHALL carry the clause `inputs.gates_only != true`. Unset, deleted or any other value SHALL mean a dry run.

#### Scenario: Variable unset
- **WHEN** `CASCADE_DRY_RUN` does not exist and a dispatch arrives
- **THEN** the run computes the diff, writes it to the job summary with "DRY RUN: nothing was pushed", `Publish` is skipped, and `deps/cascade` is not created or moved

#### Scenario: Live repo
- **WHEN** `CASCADE_DRY_RUN` is `false` and a dispatch arrives that moves a pin
- **THEN** the run pushes `deps/cascade` and creates or updates the cascade PR

#### Scenario: A gates-only run on a live repo
- **WHEN** `CASCADE_DRY_RUN` is `false` and `cascade-gates.yml` dispatches `deps-cascade.yml` with `gates_only=true`
- **THEN** the run posts the gate statuses, `Publish` is skipped, and nothing is pushed

### Requirement: The cascade-publish action publishes nothing unless its dry-run input is exactly false
As in requirement "The receiver plans a dry run unless CASCADE_DRY_RUN is exactly false": The `cascade-publish` action's required `dry-run` input SHALL receive the same expression as the `cascade` job, so the action itself publishes nothing unless the value is exactly `false`, whatever the `publish` job's `if:` says.

#### Scenario: A mistyped publish condition
- **WHEN** a pull request changes the `publish` job's `if:` or replaces either `dry-run` value with a literal
- **THEN** the wiring check fails naming the changed value, and even if it merged, the action would publish nothing while `CASCADE_DRY_RUN` is not `false`

### Requirement: A dry run pushes nothing and writes no PR, label or comment
As in requirement "The receiver plans a dry run unless CASCADE_DRY_RUN is exactly false": In a dry run nothing SHALL be pushed, and no PR, label or comment SHALL be written; the gate statuses SHALL still be posted (Phase 3 wiring contract (version 3.1) §5, §9.1; workspace RELEASING.md, "Stop switches").

#### Scenario: Live repo, manual dry run
- **WHEN** `CASCADE_DRY_RUN` is `false` and someone runs `gh workflow run deps-cascade.yml -f dry_run=true`
- **THEN** that run pushes nothing

### Requirement: The Cascade gates workflow runs on pull_request_target and calls the shared gates workflow
The repo SHALL have `.github/workflows/cascade-gates.yml`, named `Cascade gates`, triggered by `pull_request_target` with types `opened`, `reopened` and `synchronize`. It SHALL declare `permissions: {}` at the top, use the concurrency group `cascade-gates-<PR number>` with `cancel-in-progress: true`, and have a single job, `gates`, that calls `open-platform-model/.github/.github/workflows/cascade-gates.yml` at the repo's pinned `.github` SHA.

#### Scenario: Ordinary pull request
- **WHEN** a PR from a branch other than `release-please--*` is opened
- **THEN** `cascade/freshness` and `cascade/settled` appear on its head as `success` with `n/a: not a release PR`

#### Scenario: Release PR head moves
- **WHEN** release-please pushes a new head to the open release PR
- **THEN** the caller dispatches `deps-cascade.yml` with `gates_only=true`, and that run posts both statuses on the new head without running the cascade on `main`

### Requirement: The gates job grants two permissions, passes the modes and checks out no pull request code
As in requirement "The Cascade gates workflow runs on pull_request_target and calls the shared gates workflow": The job SHALL grant `statuses: write` and `actions: write` and nothing else, and SHALL pass `g2-mode` and `g3-mode` from `CASCADE_G2_MODE` and `CASCADE_G3_MODE` (default `warn`). The workflow SHALL contain no step of its own and SHALL check out no pull request code.

#### Scenario: Behind shipped pin in warn mode
- **WHEN** the release PR's `go.mod` pins library `v1.0.0-beta.4` while `v1.0.0-beta.5` is published, and `CASCADE_G2_MODE` is unset
- **THEN** `cascade/freshness` is `success` with a description starting `WARN:` that names library, and the release PR can still merge

### Requirement: Neither gate status is a required check before Phase 5
Neither `cascade/freshness` nor `cascade/settled` SHALL be a required check before Phase 5 (Phase 3 wiring contract (version 3.1) §8.3, §8.4; workspace RELEASING.md, "Gates").

#### Scenario: The gate statuses are not required checks before Phase 5
- **WHEN** the checks required on the repo are read before Phase 5
- **THEN** neither `cascade/freshness` nor `cascade/settled` is a required check

### Requirement: The first cascade references name one .github main SHA
Every reference to the cascade code SHALL name one full 40-character SHA of a commit on `open-platform-model/.github` `main`, never a branch or a tag, followed by the comment `# .github main`: the `cascade-notify` step in `release.yml`, the `cascade-receive.yml` call and the `cascade-publish` step in `deps-cascade.yml`, the `cascade-gates.yml` call.
(Requirement: The other cascade references name one .github main SHA)

#### Scenario: A branch ref
- **WHEN** a pull request names `cascade-notify@main`
- **THEN** the wiring check fails on the notify `uses` and on `one .github SHA`

### Requirement: The other cascade references name one .github main SHA
Every reference to the cascade code SHALL name one full 40-character SHA of a commit on `open-platform-model/.github` `main`, never a branch or a tag, followed by the comment `# .github main`: the `ref:` of the `open-platform-model/.github` checkout in `cascade-task.yml`, and the `ref:` of the resolver checkout in `module-deps.yml` (the operator module's own cascade PR, see `deps-cascade`).
(Requirement: The cite-lint call names one .github main SHA)

#### Scenario: The resolver checkouts name one .github SHA
- **WHEN** the `ref:` of the `open-platform-model/.github` checkout in `cascade-task.yml` and the `ref:` of the resolver checkout in `module-deps.yml` are read
- **THEN** each names one full 40-character SHA of a commit on `open-platform-model/.github` `main`, followed by the comment `# .github main`

### Requirement: The cite-lint call names one .github main SHA
Every reference to the cascade code SHALL name one full 40-character SHA of a commit on `open-platform-model/.github` `main`, never a branch or a tag, followed by the comment `# .github main`: the `uses` of the `cite-lint.yml` call in `cite-lint.yml` (see `comment-rule-check`).
(Requirement: The first cascade references name one .github main SHA)

#### Scenario: The cite-lint call names one .github SHA
- **WHEN** the `uses` of the `cite-lint.yml` call in `cite-lint.yml` is read
- **THEN** it names one full 40-character SHA of a commit on `open-platform-model/.github` `main`, followed by the comment `# .github main`

### Requirement: All seven cascade references carry the same SHA
All seven SHALL carry the same SHA, which moves only through a `ci(deps): pin the cascade to .github <sha7>` pull request.
(Requirement: The first cascade references name one .github main SHA)
(Requirement: The other cascade references name one .github main SHA)
(Requirement: The cite-lint call names one .github main SHA)

#### Scenario: One reference moves alone
- **WHEN** a pull request changes only the `cascade-gates.yml` call to another SHA
- **THEN** the "Verify the cascade wiring" step fails with `one .github SHA` naming both SHAs, and the `Lint` check is red

#### Scenario: The module's resolver checkout follows the pin
- **WHEN** a pull request moves the other `.github` references to a new SHA and leaves `module-deps.yml`'s resolver checkout behind
- **THEN** the wiring check fails on `one .github SHA`

### Requirement: The cascade-task workflow fails when the resolver is missing at the pinned commit
`cascade-task.yml` SHALL fail, not skip, when the resolver is missing at the pinned commit.

#### Scenario: The resolver is missing at the pin
- **WHEN** `cascade-task.yml` runs and the pinned `.github` commit has no `cascade-resolve.sh`
- **THEN** the run fails with `no cascade resolver at the pinned .github commit` instead of skipping S5

### Requirement: The wiring check script is a copy of the .github script
`.tasks/cascade/wiring-check.sh` SHALL be byte-identical to `.github/scripts/cascade/wiring-check.sh` at the pinned SHA and SHALL move only with the pin.
(Requirement: The wiring-check.yaml values hold the receiver, env, CI and notify settings)

#### Scenario: A copy left behind at a pin bump
- **WHEN** a pull request moves every reference to a new SHA but keeps the old `wiring-check.sh`
- **THEN** the "Verify the cascade wiring" step fails with `differs from` naming the file at the new SHA

### Requirement: The wiring-check.yaml values hold the receiver, env, CI and notify settings
As in requirement "The wiring check script is a copy of the .github script": this repo's values SHALL live in `.tasks/cascade/wiring-check.yaml` (`receiver: true`, `env-allow` `REGISTRY`, `IMAGE_NAME` and `CUE_VERSION`, `ci` `lint.yml` job `lint`, the notify `needs`, `if` and `tag` below, `labels-managed: false`).
(Requirement: The wiring check fails on notify job values)
(Requirement: The wiring-check.yaml values name the extra references and the publishing workflows)

#### Scenario: The values file holds the receiver, env, CI and notify settings
- **WHEN** `.tasks/cascade/wiring-check.yaml` is read
- **THEN** it holds `receiver: true`, `env-allow` `REGISTRY`, `IMAGE_NAME` and `CUE_VERSION`, `ci` `lint.yml` job `lint`, the notify `needs`, `if` and `tag`, and `labels-managed: false`

### Requirement: The wiring-check.yaml values name the extra references and the publishing workflows
As in requirement "The wiring-check.yaml values hold the receiver, env, CI and notify settings": this repo's values SHALL live in `.tasks/cascade/wiring-check.yaml` (`extra-references` naming the `module-deps.yml` resolver and the `cite-lint.yml` call (kind `cite-lint`), and `publish-workflows` naming every workflow that publishes: `release.yml`, `publish-fixtures.yml`, `docs.yml`, `image-pr.yml`, `test-e2e.yml`, `module-image.yml` and `module-deps.yml`).

#### Scenario: The values file names the extra references and the publishing workflows
- **WHEN** `.tasks/cascade/wiring-check.yaml` is read
- **THEN** it holds `extra-references` naming the `module-deps.yml` resolver and the `cite-lint.yml` call (kind `cite-lint`), and `publish-workflows` naming every workflow that publishes

### Requirement: The Lint job runs the wiring check on every pull request
As in requirement "The wiring check script is a copy of the .github script": The required `Lint` job SHALL run it on every pull request as the step "Verify the cascade wiring", exactly `bash .tasks/cascade/wiring-check.sh --pin-on-main` with only `GH_TOKEN: ${{ github.token }}` in its `env`, preceded only by SHA-pinned actions from other repos; `task cascade:wiring:check` runs it offline.

#### Scenario: The wired tree passes
- **WHEN** the `Lint` job runs on a pull request that leaves the cascade wiring as the contract shapes it
- **THEN** the step prints `cascade wiring: ok, .github <SHA> (.github main)`

#### Scenario: A run step before the wiring step
- **WHEN** a pull request adds a `run:` step to the `Lint` job above "Verify the cascade wiring"
- **THEN** the wiring check fails naming the step that is not a pinned action

### Requirement: The wiring check fails on a key job's shape
`.tasks/cascade/wiring-check.sh` SHALL fail when: a key-holding job (`notify-downstream`, `publish`) differs from its contract shape in its keys, step keys, `with` keys, Environment, permissions, step count, action reference, key inputs or a `runs-on` other than `ubuntu-latest`.
(Requirement: The wiring check fails on notify job values)

#### Scenario: A self-hosted runner for a key job
- **WHEN** a pull request sets `runs-on: self-hosted` on `notify-downstream`
- **THEN** the wiring check fails naming `release.yml:notify-downstream runs-on`

### Requirement: The wiring check fails on notify job values
`.tasks/cascade/wiring-check.sh` SHALL fail when: the `notify-downstream` job's `needs` differs from `[release-please, publish-release]`, or its `if` from `needs.release-please.outputs.release_created == 'true' && vars.CASCADE_NOTIFY != 'off'` (contract §4.6 prints `releases_created`, which a module-only release also sets), or its `tag` input or `timeout-minutes` (20) from contract §4.6.
(Requirement: The wiring check fails on publish job values and the release env)

#### Scenario: The wiring check fails on a notify job value off the contract
- **WHEN** the `needs` of the `notify-downstream` job differs from `[release-please, publish-release]`, or its `if` from `needs.release-please.outputs.release_created == 'true' && vars.CASCADE_NOTIFY != 'off'`, or its `tag` input or `timeout-minutes` (20) from contract §4.6
- **THEN** the step "Verify the cascade wiring" fails

### Requirement: The wiring check fails on publish job values and the release env
`.tasks/cascade/wiring-check.sh` SHALL fail when: the `publish` job's `needs` differs from `cascade`, its `if:` or `gates-only` input from the receiver caller shape, or its `timeout-minutes` from 15; `release.yml`'s workflow-level `env` is not a map or has a key other than `REGISTRY`, `IMAGE_NAME` and `CUE_VERSION`.
(Requirement: The wiring check fails on receiver values and key readers)

#### Scenario: A startup-code variable in release.yml
- **WHEN** a pull request adds `BASH_ENV` or any key other than `REGISTRY`, `IMAGE_NAME` and `CUE_VERSION` to `release.yml`'s workflow-level `env`
- **THEN** the wiring check fails naming the key

### Requirement: The wiring check fails on receiver values and key readers
`.tasks/cascade/wiring-check.sh` SHALL fail when: `deps-cascade.yml`'s top-level keys, concurrency group, `publish` condition or either `dry-run` value differ from the contract; any other workflow value reads the key or any other job declares `environment: cascade`; a job reads `RELEASE_APP_PRIVATE_KEY` without declaring exactly `environment: release`, or a job declares `release` without reading it.
(Requirement: The wiring check fails on caches, checkouts and the Lint env)

#### Scenario: The wiring check fails on a second reader of the key
- **WHEN** any other workflow value reads the key or any other job declares `environment: cascade`
- **THEN** the step "Verify the cascade wiring" fails

### Requirement: The wiring check fails on caches, checkouts and the Lint env
`.tasks/cascade/wiring-check.sh` SHALL fail when: a workflow listed in `publish-workflows` uses an Actions cache; a `.github` checkout is not `actions/checkout@<SHA>` with exactly `repository`, `ref`, `path` and `persist-credentials: false`; the `Lint` workflow or job `env` names anything outside `CUE_*`, `OPM_*`, `REGISTRY` and `IMAGE_NAME`, or the job has a `container` or `services`.
(Requirement: The wiring check fails on a secrets key or a reference set)

#### Scenario: A cache in a publishing workflow
- **WHEN** a pull request drops `cache: false` from a `setup-go` step in `test-e2e.yml`
- **THEN** the wiring check fails naming `test-e2e.yml cache use`

### Requirement: The wiring check fails on a secrets key or a reference set
`.tasks/cascade/wiring-check.sh` SHALL fail when: a call into `.github` has a `secrets` key; or the set of `.github` references, their SHA or their comment differs.
(Requirement: The wiring check fails on a key job's shape)

#### Scenario: The wiring check fails when the set of .github references differs
- **WHEN** the set of `.github` references, their SHA or their comment differs
- **THEN** the step "Verify the cascade wiring" fails

### Requirement: The wiring check also fails on a SHA that is not on .github main
As in requirement "The Lint job runs the wiring check on every pull request": With `--pin-on-main` it SHALL also fail when the SHA is not on `.github` `main` or the running copy differs from the file at the SHA.
(Requirement: The wiring check fails on a key job's shape)

#### Scenario: A SHA that is not on .github main fails the check
- **WHEN** the pinned SHA is not on `.github` `main` and the check runs with `--pin-on-main`
- **THEN** the check fails

### Requirement: The wiring check guards against mistakes, not a deliberate edit
`.tasks/cascade/wiring-check.sh` SHALL be byte-identical to `.github/scripts/cascade/wiring-check.sh` at the pinned SHA and SHALL move only with the pin.
It guards against mistakes, not a deliberate edit: code-owner review is the guard against that, and the `main` ruleset does not require review while OPM is in beta (no approvals, workspace PR 29) (Phase 3 wiring contract (version 3.1) §2.4, §10.1 items 2, 3 and 6; the supervisor's addendum; `.github` README at `2b66da5`, "The wiring check").

#### Scenario: A deliberate edit is guarded by code-owner review
- **WHEN** a deliberate edit changes the cascade wiring
- **THEN** code-owner review is the guard against it, and the `main` ruleset does not require review while OPM is in beta
