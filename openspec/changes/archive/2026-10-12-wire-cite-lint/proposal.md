## Why

`.github` PR 22 merged as `2b66da51fede4332969c324b069012632d8edea9`. It adds the reusable `cite-lint.yml` workflow and a `cite-lint` reference kind in `wiring-check.sh`. The comment rule (no decision number, no lone enhancement id in a Go or CUE comment) is the same in every repo, so the operator runs it locally and in CI, and a later bad citation fails the pull request that adds it. The operator is aligned with the rule since `align-enhancement-store`, so the first run has no finding outside the published fixtures.

## What Changes

- All `.github` references in `.github/workflows` move from `6df4000f` to `2b66da51fede4332969c324b069012632d8edea9 # .github main`. `.tasks/cascade/wiring-check.sh` becomes the file at that SHA (byte for byte). Between the two SHAs the script only gains the `cite-lint` kind and its plain-call rule; no other cascade code or caller shape changed, so the cascade jobs behave as before.
- New workflow `.github/workflows/cite-lint.yml`: one job that calls the reusable workflow on `pull_request` and on `push` to `main`, excluding the published fixtures `test/fixtures/modules/` and `test/fixtures/catalogs/`.
- New task `task cite-lint` with the same tool commit and excludes. The excludes sit in two places (task and workflow); a comment in each names the other.
- `.tasks/cascade/wiring-check.yaml` declares the call (kind `cite-lint`).
- `AGENTS.md` lists the task.

## Capabilities

### New Capabilities
- `comment-rule-check`: the local task and the CI call of the comment rule.

### Modified Capabilities
- `cascade-receiver`: the pin requirement counts the `cite-lint.yml` call and cites the `.github` README at the new SHA.

## Impact

CI and task files only (`ci`, `ci(deps)`, `build`, `docs`, `chore(openspec)`); no operator or module release. The check is not required by the ruleset; making it required is an owner step. The OpenSpec CLI pin: this repo pins none, and 53 specs fail `--strict` under 1.14.1 at base, so no pin is added.
