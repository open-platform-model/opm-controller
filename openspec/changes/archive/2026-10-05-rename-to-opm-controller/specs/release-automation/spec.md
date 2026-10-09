## MODIFIED Requirements

### Requirement: Git tag and GitHub Release on merge
When a Release PR is merged to `main`, release-please SHALL create the git tag eagerly at the release commit and a GitHub Release in draft state with the changelog section as release notes. The controller package `"."` tags `vX.Y.Z` (e.g. `v0.2.0`, or `v1.0.0-beta.6` on the beta line); the controller module's package tags `opm_controller-vX.Y.Z` (see `controller-module-release`). Each release-please package SHALL be configured with `draft: true` and `force-tag-creation: true`. A controller release's draft SHALL become public only through the publish job `publish-release` defined in "Release published once after every release job"; a module release's draft only through the module's final publish job `module-publish-release`.

#### Scenario: Release PR merged
- **WHEN** the Release PR is merged to `main`
- **THEN** release-please creates a git tag matching the version (prefixed with `v` for the controller, with `opm_controller-v` for the module) at the release commit and a draft GitHub Release with the changelog for that version as the body, and the release is not visible as published until every release job has succeeded

#### Scenario: Release PR closed without merge
- **WHEN** the Release PR is closed without merging
- **THEN** no tag or release SHALL be created; the next push to `main` re-opens or creates a new Release PR

#### Scenario: Config enables draft-first
- **WHEN** `release-please-config.json` is inspected on `main`
- **THEN** package `"."` and package `"modules/opm_controller"` SHALL each set `"draft": true` and `"force-tag-creation": true`

### Requirement: Release assets upload only to a draft release
Every step that uploads an asset to the GitHub Release (`install.yaml` of a module release, `opm-examples.tar.gz`, the example manifests) SHALL first confirm that the single release for the tag is a draft and SHALL fail without uploading when it is published. Replacing an existing asset (`--clobber`) SHALL be used only on a draft. Every GitHub CLI call in the release workflow SHALL name the repository explicitly. Source: 0021:D10:R8.

#### Scenario: Upload to the draft
- **WHEN** the example publishing job uploads `opm-examples.tar.gz` for a release that is still a draft
- **THEN** the upload succeeds, and re-running the job replaces the asset on the draft

#### Scenario: Module manifest upload to the draft
- **WHEN** the module publish job uploads `install.yaml` for `opm_controller-v0.2.0` while its release is a draft
- **THEN** the upload succeeds, and re-running the job replaces the asset on the draft

#### Scenario: Upload to a published release refused
- **WHEN** an upload step runs for a tag whose release is already published
- **THEN** the step fails before uploading and tells the operator to release the next version instead

### Requirement: Release published once after every release job
The release workflow SHALL contain a publish job (`publish-release`) that depends on every job producing a release artifact (the image job and the example publishing job), runs only when the controller package `"."` created a release (a module-only release does not run it), holds `contents: write` and no other write permission, checks out only the files it runs, confirms the required asset `opm-examples.tar.gz` is attached to the draft, and then publishes the draft. It SHALL leave the Pre-release flag set by release-please unchanged. When the release is already published it SHALL succeed without changing anything. The only jobs that SHALL depend on it are the cascade notify job and the module's image-bump PR job (`module-image-pr`, see `controller-module-release`), neither of which produces a release artifact. A controller release SHALL NOT carry `install.yaml`: the install manifest is an asset of the module's release (`controller-module-release`).

#### Scenario: All release jobs succeed
- **WHEN** the image and example jobs finish successfully for `v1.0.0-beta.3`
- **THEN** `publish-release` publishes the draft, and the GitHub Release is public, flagged Pre-release and carries every asset

#### Scenario: A release job fails
- **WHEN** the example publishing job fails for `v1.0.0-beta.3`
- **THEN** `publish-release` does not run, the release stays a draft, and "Re-run failed jobs" on the same workflow run completes and publishes it

#### Scenario: Required asset missing
- **WHEN** `publish-release` runs and the draft lacks `opm-examples.tar.gz`
- **THEN** the job fails and the release stays a draft

#### Scenario: Module-only release
- **WHEN** a push to `main` creates only a module release
- **THEN** `publish-release` does not run

#### Scenario: No controller manifest
- **WHEN** a controller release `v1.0.0-beta.9` is published
- **THEN** its GitHub Release carries `opm-examples.tar.gz` and no `install.yaml`

### Requirement: Release PR opens for releasable commits on push to main
The release-please workflow SHALL run on every push to the `main` branch. It SHALL open a Release PR if releasable commits exist since the last release tag. Releasable commits are those whose type has a visible changelog section in `release-please-config.json` (`feat`, `fix`, `perf`, `revert`, `deps`, `refactor`). `docs`, `chore`, `test`, `ci` and `build` are hidden and SHALL NOT by themselves open a Release PR, and a `release-as` value in `release-please-config.json` does not make them releasable. If a Release PR already exists, it SHALL update the PR with the latest accumulated changes. Hiding `docs` keeps a doc-only commit from cutting a controller release that would cascade a pointless bump into the cli (owner decision 2026-10-01 (RELEASING.md, Pin classes)).

#### Scenario: First feat commit after a release
- **WHEN** a `feat(scope): description` commit is pushed to `main` and no open Release PR exists
- **THEN** release-please opens a new Release PR proposing the next version per "Version bump determination per release line", with an updated CHANGELOG.md

#### Scenario: Subsequent fix commit with open Release PR
- **WHEN** a `fix(scope): description` commit is pushed to `main` and a Release PR already exists
- **THEN** release-please updates the existing Release PR to include the new fix in the changelog, adjusting the proposed version if needed

#### Scenario: Only non-releasable commits
- **WHEN** only `docs`, `chore`, `test`, `ci` or `build` commits are pushed since the last release, whether or not `release-please-config.json` carries a `release-as` value
- **THEN** release-please SHALL NOT open a Release PR (no version bump warranted)

#### Scenario: Docs-only commits do not cut a release
- **WHEN** only `docs` commits are pushed since the last release
- **THEN** release-please SHALL NOT open a Release PR, and `release-please-config.json` SHALL declare the `docs` section with `"hidden": true`

#### Scenario: Refactor commits still cut a release
- **WHEN** only `refactor` commits are pushed since the last release
- **THEN** release-please opens a Release PR listing them under Code Refactoring

### Requirement: Release PR passes the release-pin gate
On every Release PR, the release-pin gate (G1, workspace RELEASING.md, section "Gates") SHALL run as a step inside an existing CI job that runs on every pull request, never as a job of its own, so a skipped job can never report a pass. A PR is a Release PR when `${{ github.head_ref || github.ref_name }}` starts with `release-please--`, except the controller module's release PR (`release-please--branches--main--components--opm_controller`): it ships none of the files this gate checks, and the module release gate (`controller-module-release`) checks its pins. The gate SHALL fail the job when any of the following holds on the PR head:
- `go.mod` carries a `replace` directive;
- a required module under `github.com/open-platform-model/` is pinned to a Go pseudo-version, or to a version that is not an existing tag of that module's repository;
- a published fixture's `cue.mod/module.cue` (`test/fixtures/modules/*`, `test/fixtures/modulepackages/*`) pins a dependency version containing `-0.dev.`;
- any `cue.mod/local-module.cue` is tracked by git.

Each failure SHALL name the offending file and pin. On pull requests that are not Release PRs the gate step SHALL be skipped. The same check SHALL be runnable locally as one task. The job carrying the gate SHALL report a check name no other workflow job in the repo uses (`Lint`), so a ruleset can require exactly that check.

#### Scenario: Replace directive blocks the release
- **WHEN** a Release PR head has `replace github.com/open-platform-model/library => ../library` in `go.mod`
- **THEN** the gate step fails and its output names the replace directive

#### Scenario: Pseudo-version blocks the release
- **WHEN** a Release PR head pins `github.com/open-platform-model/library v1.0.0-beta.2.0.20261001120000-abcdef123456`
- **THEN** the gate step fails and names the module and version

#### Scenario: Untagged version blocks the release
- **WHEN** a Release PR head pins `github.com/open-platform-model/library v1.0.0-beta.9` and the library repository has no tag `v1.0.0-beta.9`
- **THEN** the gate step fails and names the module and version

#### Scenario: Dev CUE pin in a published fixture blocks the release
- **WHEN** a Release PR head has `test/fixtures/modules/hello/cue.mod/module.cue` pinning `v: "v2.0.0-0.dev.20261001"`
- **THEN** the gate step fails and names the file

#### Scenario: Tracked local-module.cue blocks the release
- **WHEN** a Release PR head tracks `test/fixtures/modules/hello/cue.mod/local-module.cue`
- **THEN** the gate step fails and names the file

#### Scenario: Clean release PR passes
- **WHEN** a Release PR head has no replace directive, pins `github.com/open-platform-model/library` to an existing tag, and no published fixture carries a dev pin or a tracked `local-module.cue`
- **THEN** the gate step passes

#### Scenario: The gate's check name is unique
- **WHEN** the `name:` of every job under `.github/workflows/` is listed
- **THEN** the job carrying the release-pin gate reports `Lint`, and no other job reports that name

#### Scenario: Ordinary PRs skip the gate
- **WHEN** a pull request's head branch does not start with `release-please--`
- **THEN** the gate step is skipped and the job's result depends only on its other steps

#### Scenario: The module's release PR skips the gate
- **WHEN** the head branch is `release-please--branches--main--components--opm_controller`
- **THEN** the release-pin gate step is skipped and the module release gate step runs

### Requirement: Dependabot leaves OPM Go modules to the release cascade
The `gomod` entry of `.github/dependabot.yml` SHALL ignore every dependency matching `github.com/open-platform-model/*`. Those pins move only through the release cascade, which titles a shipped bump `fix(deps)` so it releases (workspace RELEASING.md, section "Pin classes"). The `github-actions` entry SHALL ignore `open-platform-model/.github*`: the cascade references move together, one `.github` SHA for the repo, only through a `ci(deps): pin the cascade to .github <sha7>` pull request (Phase 3 wiring contract (version 3.1) §2.4, §10.1 item 7). Third-party Go modules, except `cuelang.org/go` and `cuelabs.dev/go/oci/ociregistry` (see "Dependabot leaves cuelang.org/go to library releases"), and other GitHub Actions SHALL keep their Dependabot updates.

#### Scenario: Library release opens no Dependabot PR
- **WHEN** library publishes a new tag
- **THEN** Dependabot opens no PR bumping `github.com/open-platform-model/library` in opm-controller

#### Scenario: Third-party bumps continue
- **WHEN** a new `k8s.io/api` release exists
- **THEN** Dependabot still proposes the grouped Kubernetes bump

#### Scenario: A .github main commit opens no Dependabot PR
- **WHEN** a commit lands on `open-platform-model/.github` `main` after the repo's pinned SHA
- **THEN** Dependabot opens no PR moving any `open-platform-model/.github` reference, and the five references keep one SHA

### Requirement: Release notifies downstream after it is published
The release workflow SHALL contain a job `notify-downstream`, owned by this repo, that needs `release-please` and `publish-release`, runs on `ubuntu-latest`, declares `environment: cascade`, has a 20-minute timeout, grants `contents: read` and no other permission, and has exactly one step: the action `open-platform-model/.github/.github/actions/cascade-notify` at the repo's pinned `.github` `main` SHA, with the inputs `tag: ${{ needs.release-please.outputs.tag_name }}`, `client-id: ${{ vars.CASCADE_APP_CLIENT_ID }}` and `private-key: ${{ secrets.CASCADE_APP_PRIVATE_KEY }}` and no other input. Its `if` SHALL be `needs.release-please.outputs.release_created == 'true' && vars.CASCADE_NOTIFY != 'off'`: it SHALL run only when the controller package `"."` created a release, never for a release created only by the controller module's package, `publish-release` succeeded, and the repo variable `CASCADE_NOTIFY` is not `off`. It SHALL NOT wait for, or be skipped by, `publish-docs`. The key SHALL be read only in the caller-owned `notify-downstream` or `publish` job, which declares `environment: cascade` and passes `secrets.CASCADE_APP_PRIVATE_KEY` only as the `private-key` input of the SHA-pinned cascade action; that job SHALL have no checkout or `run:` of its own, and no `env:`, `container:` or `services:`; no reusable call SHALL pass `secrets:` or `secrets: inherit`. The action mints the App token and SHALL dispatch `upstream-released` only to `cli` (Phase 3 wiring contract (version 3.1) §3.1, §4.3, §4.6, §10.1 item 9; workspace RELEASING.md, "Notify after publish").

#### Scenario: Published release notifies the cli
- **WHEN** the release run for `v1.0.0-beta.9` publishes the draft with `opm-examples.tar.gz` attached
- **THEN** `notify-downstream` runs after `publish-release` in the `cascade` Environment, and the cli receives a `repository_dispatch` of type `upstream-released` with source `opm-controller` and tag `v1.0.0-beta.9`

#### Scenario: Release left a draft does not notify
- **WHEN** `publish-examples` fails, so `publish-release` does not run and the release stays a draft
- **THEN** `notify-downstream` does not run, and "Re-run failed jobs" that completes and publishes the release also runs `notify-downstream`

#### Scenario: Failed notify is recovered without republishing
- **WHEN** `publish-release` succeeded and `notify-downstream` failed, for example because the App installation was suspended
- **THEN** "Re-run failed jobs" re-runs only `notify-downstream`, and the release stays as it is

#### Scenario: Notify switched off
- **WHEN** the repo variable `CASCADE_NOTIFY` is `off` and a release is published
- **THEN** `notify-downstream` is skipped, the release is still published, and the cli's daily sweep picks the release up

#### Scenario: Docs failure does not hold notify
- **WHEN** `publish-docs` fails and `publish-release` succeeds
- **THEN** `notify-downstream` still runs

#### Scenario: Run without a release
- **WHEN** a push to `main` cuts no release
- **THEN** `notify-downstream` is skipped

#### Scenario: A step added beside the action is refused
- **WHEN** a pull request adds a checkout or a `run:` step to `notify-downstream`
- **THEN** the `Lint` job's "Verify the cascade wiring" step fails naming `release.yml:notify-downstream step count`

#### Scenario: Module-only release does not notify
- **WHEN** a push to `main` creates only `opm_controller-v0.2.0`
- **THEN** `notify-downstream` is skipped, and no dispatch carries an empty tag

### Requirement: Dependabot leaves cuelang.org/go to library releases
The `gomod` entry of `.github/dependabot.yml` SHALL ignore the dependencies `cuelang.org/go` and `cuelabs.dev/go/oci/ociregistry` (the CUE team's OCI module, whose version the library's `go.mod` sets, `cuelang.org/go` setting only a floor), so the controller's CUE version moves only through a library release: the library bumps CUE in its own pull request, where the closedness canary runs, and the controller picks the new version up when the release cascade (or a hand-made `fix(deps)` pull request) moves the library pin, whose `go get` of the library raises `cuelang.org/go` by minimal version selection (kept by `go mod tidy`). Other third-party Go modules SHALL keep their Dependabot updates. Source: owner decision j4 of the kernel beta.1 plan walkthrough (2026-10-03); workspace RELEASING.md, section "Pin classes".

#### Scenario: New CUE release opens no Dependabot PR
- **WHEN** a new `cuelang.org/go` version is published and Dependabot's weekly `gomod` run executes
- **THEN** Dependabot opens no pull request for `cuelang.org/go` or `cuelabs.dev/go/oci/ociregistry` in opm-controller

#### Scenario: CUE arrives with a library release
- **WHEN** a library release that requires a newer `cuelang.org/go` is pinned by the cascade's library bump
- **THEN** that bump pull request moves `cuelang.org/go` in `go.mod` and `go.sum` to at least the version the library requires, and the cascade reports the move as a warning (`.tasks/cascade/cascade.sh`, Phase C, `warn_deps`)

#### Scenario: Other third-party modules still update
- **WHEN** a newer `github.com/fluxcd/pkg/runtime` is published
- **THEN** Dependabot still proposes the bump

## REMOVED Requirements

### Requirement: Changelog generation

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-021); this requirement is restated under the controller's name.

**Migration**: See the requirement "The changelog is generated per package" in this capability.

### Requirement: Beta prerelease line

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-021); this requirement is restated under the controller's name.

**Migration**: See the requirement "The controller's beta prerelease line" in this capability.

### Requirement: The Release PR rewrites the install page's operator version

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-021); this requirement is restated under the controller's name.

**Migration**: See the requirement "The Release PR rewrites the install page's controller version" in this capability.

## ADDED Requirements

### Requirement: The changelog is generated per package
The workflow SHALL generate and maintain a `CHANGELOG.md` file at the repository root for the controller package. Entries SHALL be grouped under the visible sections that `release-please-config.json` declares (Features, Bug Fixes, Performance Improvements, Reverts, Dependencies, Code Refactoring). Commits of a hidden type (`docs`, `chore`, `test`, `ci`, `build`) SHALL NOT appear. A commit that changes only files under the controller module's directory `modules/opm_controller/` belongs to the module's own changelog (see `controller-module-release`) and SHALL NOT appear in the root `CHANGELOG.md`. Documentation entries already in CHANGELOG.md from earlier releases SHALL stay as they are.

#### Scenario: Changelog includes all commit types
- **WHEN** the Release PR is created or updated
- **THEN** `CHANGELOG.md` SHALL list every commit since the last release whose type has a visible section in `release-please-config.json` and that changes a file outside `modules/opm_controller/`, grouped under that section, with commit messages as entries, and SHALL list no commit of a hidden type

#### Scenario: Changelog preserves history
- **WHEN** a new release is cut
- **THEN** the new changelog section SHALL be prepended to existing content, preserving prior release entries

#### Scenario: Module-only commit stays out of the controller changelog
- **WHEN** `fix(deps): deploy controller v1.0.0-beta.6 from the controller module`, which changes only `modules/opm_controller/`, is merged
- **THEN** the controller's Release PR and root `CHANGELOG.md` do not list it

### Requirement: The controller's beta prerelease line
The release-please package `"."` SHALL be configured with `versioning: prerelease`, `prerelease: true` and `prerelease-type: beta` while the controller is on its beta line. From its first beta, a prerelease line (opmodel.dev/core@v2, library, cli, opm-controller) is on the path to GA. A breaking change is still allowed during beta, but only as a `!` in the PR title (`feat!:`), which becomes the CHANGELOG entry; its migration note goes in the PR body, which the entry links and which never reaches `main`. It advances the `-beta.N` counter and never moves the module path to a new major. Stable lines (opmodel.dev/catalogs/opm@v4 and the module fleets) keep the normal SemVer rule: a break is a new major. The controller module's own package follows its 0.x rule instead (`controller-module-release`). A core beta break that would force a catalogs/opm major needs owner sign-off. A controller `feat!` that the released cli cannot drive SHALL merge only after the cli release that can drive it, and no `release-as` value SHALL hop the controller to a new minor or major (such as `1.1.0-beta.1`) during beta. Every beta GitHub Release SHALL be flagged Pre-release. Changing `prerelease-type` alone SHALL NOT be relied on to change the label of a version that already carries a suffix; a label change needs a `release-as` value (see "Forced version via the release-as config key"). GA drops the suffix: `prerelease: false` plus a visible carrier commit per package, in dependency order.

#### Scenario: Beta releases are flagged Pre-release
- **WHEN** the Release PR for `1.0.0-beta.N` is merged
- **THEN** release-please creates tag `v1.0.0-beta.N` and a GitHub Release marked Pre-release, and the image and the example assets are published under that tag

#### Scenario: Breaking change during beta carries its migration note
- **WHEN** a pull request titled `feat(api)!: ...`, whose body describes the migration, is squash-merged on the beta line
- **THEN** the next Release PR proposes the next `-beta.N`, its CHANGELOG section lists the PR title as a breaking-change entry linking the pull request, and the migration note is read from that pull request

#### Scenario: Controller break waits for a cli that can drive it
- **WHEN** a controller `feat!` changes behavior the newest released cli cannot drive
- **THEN** it merges only after a cli release that can drive it, ships as the next `1.0.0-beta.N`, and no `release-as` value moves the controller to `1.1.0-beta.1` or any other minor or major

#### Scenario: Label flip alone keeps the old label
- **WHEN** only `prerelease-type` changes from `alpha` to `beta` while the manifest reads `1.0.0-alpha.22` and `release-please-config.json` carries no `release-as`
- **THEN** the next Release PR proposes `1.0.0-alpha.23`, which is why the crossing sets `"release-as": "1.0.0-beta.1"` in `release-please-config.json`

#### Scenario: GA drops the suffix
- **WHEN** the config sets `prerelease: false` and a visible carrier commit lands on `main` while the manifest reads `1.0.0-beta.N`
- **THEN** the Release PR proposes `1.0.0` and the GitHub Release is not marked Pre-release

### Requirement: The Release PR rewrites the install page's controller version

`release-please-config.json` SHALL list `docs/site/start/install-the-controller.md` under the root package's `extra-files`, and every passage of that page that names the controller's version SHALL sit inside an `x-release-please-start-version` ... `x-release-please-end` block of HTML comment markers, each on its own line outside code fences and at the content indentation of its enclosing block (three spaces inside an ordered-list item), so the page renders as it would without them. A block SHALL name no version other than the controller's. A test in `internal/version` SHALL fail when the page has no block, when a block is not closed, or when a version inside a block differs from the `Version` constant.

#### Scenario: A Release PR moves the page with the constant

- **WHEN** release-please opens the Release PR for `1.0.0-beta.6`
- **THEN** the PR rewrites `Version` in `internal/version/version.go` and every version inside the install page's blocks to `1.0.0-beta.6`, and the release's docs bundle names `v1.0.0-beta.6` in its install commands

#### Scenario: Another project's version inside a block is refused

- **WHEN** a commit puts a core or catalog version inside one of the install page's blocks
- **THEN** `task dev:test` fails, naming the version and the expected controller version

#### Scenario: A detached marker is refused

- **WHEN** an edit removes every block from the install page, or leaves a start marker without its end
- **THEN** `task dev:test` fails, naming the page
