# comment-rule-check Specification

## Purpose
Run the shared comment rule of the enhancements store on Go and CUE files, locally with `task cite-lint` and in CI through the reusable workflow of `.github`, so a bad citation fails the pull request that adds it.

## Requirements

### Requirement: The local task runs the comment rule at the pinned tool commit
`task cite-lint` SHALL run `enh cite-lint` from `github.com/open-platform-model/enhancements/tools/enh` at commit `e0335b2ec44e659c0191bc3523d79d17f26f8e7d`, the `enh-ref` default of the reusable workflow at the repo's pinned `.github` SHA. Its excludes SHALL be exactly `test/fixtures/modules/` and `test/fixtures/catalogs/`, the published fixtures that keep old citations until their next republish.

#### Scenario: A decision number in a comment
- **WHEN** a pull request adds a Go comment that cites a decision by number outside the excluded paths
- **THEN** `task cite-lint` fails with the rule `decision-number` and the file and line

### Requirement: CI calls the shared workflow with the same excludes
`.github/workflows/cite-lint.yml` SHALL have one job that calls `open-platform-model/.github/.github/workflows/cite-lint.yml` at the repo's pinned `.github` SHA, on `pull_request` and on `push` to `main`. The job SHALL have only `name`, `permissions` (`contents: read`), `uses` and `with` (`paths` and `exclude`). The excludes SHALL equal those of the task, and a comment in each place SHALL name the other.

#### Scenario: A decision number in CI
- **WHEN** a pull request adds a comment with a decision number outside the excluded paths
- **THEN** the `Cite lint` job fails and its log lists `path:line: decision-number`

### Requirement: The wiring check holds the call to the plain shape
`.tasks/cascade/wiring-check.yaml` SHALL declare `cite-lint.yml` as an `extra-references` item of kind `cite-lint`. The file SHALL have no `pull_request_target` or `workflow_run` trigger.

#### Scenario: A call that is not plain
- **WHEN** a pull request adds `enh-ref` to the `with` of the call, or a `secrets` key to the job
- **THEN** the "Verify the cascade wiring" step of `Lint` fails naming `cite-lint.yml`

#### Scenario: The call keeps its own pin
- **WHEN** a pull request moves the call to another `.github` SHA than the cascade references
- **THEN** the wiring check fails on `one .github SHA`
