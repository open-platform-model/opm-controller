## REMOVED Requirements

### Requirement: The module is a release unit of its own
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "The module is a second release-please package", "Module release tags cannot collide with operator tags", "A commit releases only the package whose files it changes".

### Requirement: The module version follows the 0.x bump rule on the v0 path
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "Module versions on the v0 path are 0.y.z", "The first module release is 0.1.0", "The module moves to the v1 path only after operator GA".

### Requirement: The module's identity version has one writer
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "The identity-advance step is the only writer of the module's identity version", "No release-please updater writes the module's identity package", "The identity-advance step is idempotent", "The identity-advance step runs in a job of its own".

### Requirement: Module release gate
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "The module release check runs in Lint on the module's release PR", "The module release check refuses a development pin or a local module", "The module release check verifies the operator image reference", "The module release check verifies the module version", "The module release check refuses an operator below the minimum version", "The module release check compares generated CRDs and RBAC with the deployed operator", "Module releases are frozen while config differs from the deployed tag", "The module release check runs in the publish job before any push", "An offline test covers the module release checks that need no network".

### Requirement: The module publishes its tagged source once
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "The module publish job publishes the tagged source", "The module publish job runs only for a module release", "The module publish job is idempotent for the same release", "The module publish job reads registry references from cue mod resolve", "An offline test asserts the registry references of the publish job".

### Requirement: Every module release publishes the install manifest rendered from it
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "The install manifest is the module's render at its defaults", "install.yaml is attached to the draft release with the Namespace and CRDs first", "The render and the upload run apart from the publish", "The install manifest names the operator image the module names", "The install manifest yields a running operator".

### Requirement: Module releases are drafts until every asset is attached
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "A module release stays a draft with one tag until it is recovered or published", "The final module publish job publishes the draft only with install.yaml attached", "A module release is never the latest release".

### Requirement: A stranded module draft is recovered by dispatch
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "The release workflow accepts a module recovery dispatch", "The recovery dispatch refuses before any registry login", "The recovery dispatch runs main's scripts against the tag's tree", "Each recovery module job gates on its predecessor explicitly", "The recovery dispatch publishes, attaches and releases in order", "The recovery dispatch runs in a concurrency group per tag", "A push-triggered release runs the tag's own scripts", "Recovery never moves a tag or edits a release by hand".

### Requirement: Only this repository's module release publishes under the module path
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "Only the module publish job publishes under the module path", "The module publish job refuses outside the module release", "Local and test flows publish the module away from GHCR".

### Requirement: An operator release opens the PR that moves the module's image
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "An operator release opens one pull request on module/operator-image", "The module image PR sets the deployed operator release and nothing else", "The module image run fails when config differs from the operator tag", "The module image run fails for a tag that does not descend from the deployed one", "The module image PR title carries the change type", "The module image workflow accepts a dispatch for recovery", "The module/operator-image branch is rebuilt only while it holds bot commits".

## ADDED Requirements

### Requirement: The module is a second release-please package
`release-please-config.json` SHALL declare the module directory `modules/opm_operator` as a second package beside the operator package `"."`. The module package SHALL have its own version in `.release-please-manifest.json`, its own changelog, its own release PR (`separate-pull-requests: true`) and its own release tags of the form `opm_operator-vX.Y.Z`. The module's versions SHALL NOT follow the operator's versions or the cli's versions.

#### Scenario: Module tag shape

- **WHEN** the module's release PR proposing `0.2.0` is merged
- **THEN** release-please creates the tag `opm_operator-v0.2.0` and a draft GitHub Release for it, and no tag starting with `v` is created by that merge

### Requirement: Module release tags cannot collide with operator tags
No operator release tag (`vX.Y.Z`) and no operator release branch (`release/vX.Y`) can take the form `opm_operator-vX.Y.Z` of the module's release tags. Only release-please SHALL create the module's release tags, and no workflow SHALL move one.

#### Scenario: The operator's tag listing ignores module tags

- **WHEN** the tags `v1.0.0-beta.6` and `opm_operator-v0.2.0` both exist
- **THEN** a listing of `refs/tags/v*`, which the cascade resolver uses for the operator, contains only `v1.0.0-beta.6`

### Requirement: A commit releases only the package whose files it changes
The operator package SHALL exclude the module directory, so a commit that changes only the module directory proposes no operator release. A commit that changes no file under the module directory SHALL propose no module release.

#### Scenario: A module-only change releases only the module

- **WHEN** a `feat(module): ...` PR that changes only files under `modules/opm_operator/` is squash-merged to `main`
- **THEN** release-please opens or updates the module's release PR and leaves the operator's release PR unchanged

#### Scenario: An operator-only change releases only the operator

- **WHEN** a `fix(controller): ...` PR that changes no file under `modules/opm_operator/` is squash-merged
- **THEN** release-please opens or updates the operator's release PR and leaves the module's release PR unchanged

### Requirement: Module versions on the v0 path are 0.y.z
While the module path is `opmodel.dev/modules/opm_operator@v0`, every module version SHALL be `0.y.z`: a breaking release SHALL raise `y` and reset `z`, and every other release SHALL raise `z`, `feat` included. A module release is breaking when its `#config` change is breaking, when the operator release it deploys is a declared breaking change against the one the previous module release deployed, or when its CRDs stop serving a version the previous release's CRDs served.

#### Scenario: A breaking change raises the minor

- **WHEN** the module's last release is `0.3.2` and an unreleased squash commit is titled `fix(deps)!: deploy operator v1.0.0-beta.9 from the operator module`
- **THEN** the module's release PR proposes `0.4.0`

#### Scenario: A feature raises the patch

- **WHEN** the module's last release is `0.3.2` and the unreleased commits include `feat(module): add a pod annotations value` and no breaking change
- **THEN** the module's release PR proposes `0.3.3`

### Requirement: The first module release is 0.1.0
The first module release SHALL be `0.1.0`.

#### Scenario: The first release

- **WHEN** the module package has never released and a squash commit titled `feat(module): ...` that changes only `modules/opm_operator/` lands on `main`
- **THEN** the module's release PR proposes `0.1.0`, and the operator's release PR is unchanged

#### Scenario: Hidden module commits open no release

- **WHEN** only `build`, `test`, `docs`, `ci` or `chore` commits have changed `modules/opm_operator/` since the module's last release, or since it was added
- **THEN** release-please opens no module release PR

### Requirement: The module moves to the v1 path only after operator GA
The module SHALL move to the `@v1` path only with a `1.0.0` release made no earlier than the operator's first GA release, and the module release gate SHALL refuse `1.0.0` or higher while the operator release the module deploys is a prerelease.

#### Scenario: 1.0.0 before operator GA is refused

- **WHEN** the module's release PR proposes `1.0.0` and the module names `opm-operator:v1.0.0-beta.9`
- **THEN** the module release gate fails and names the operator prerelease

### Requirement: The identity-advance step is the only writer of the module's identity version
The module's `identity.Version` SHALL be written only by the release workflow's identity-advance step, so the module the publish job pushes is exactly its tagged source with nothing added at publish time. That step runs `opm module version set <version>` on the module's open release PR branch, commits `chore: advance opm_operator identity.Version to <version>`, and pushes to that branch.

#### Scenario: Release PR gains the identity commit

- **WHEN** release-please opens the module's release PR proposing `0.2.0` and the module's `identity.Version` reads `0.1.0`
- **THEN** the same workflow run pushes one commit to that PR's branch setting `identity.Version` to `0.2.0`, and the merged release commit declares `0.2.0`

### Requirement: No release-please updater writes the module's identity package
No `x-release-please-version` annotation or other release-please updater SHALL write the module's identity package, and `opm module publish --version` SHALL only assert the version.

#### Scenario: Publish asserts the version
- **WHEN** `opm module publish --version` runs for the module
- **THEN** it only asserts the version, and no release-please updater has written the module's identity package

### Requirement: The identity-advance step is idempotent
The identity-advance step SHALL be idempotent: when the branch already declares the version it SHALL change nothing.

#### Scenario: Re-run converges

- **WHEN** the step runs again for a branch whose `identity.Version` already reads `0.2.0`
- **THEN** it pushes nothing

### Requirement: The identity-advance step runs in a job of its own
The identity-advance step SHALL run in a job of its own after the release-please job, so its failure never fails the job whose outputs start the release jobs, and it SHALL hold the release App token only for its push.

#### Scenario: A failed identity advance strands no release

- **WHEN** the identity advance fails in a run where release-please created an operator release
- **THEN** the operator's release jobs still run, and "Re-run failed jobs" re-runs only the identity advance

### Requirement: The module release check runs in Lint on the module's release PR
On the module's release PR, the required `Lint` job SHALL run a module release check that fails, naming every failure before it exits.

#### Scenario: Not a release PR

- **WHEN** the `Lint` job runs on any other pull request
- **THEN** the module release check does not run

### Requirement: The module release check refuses a development pin or a local module
The module release check SHALL fail when the module's `cue.mod/module.cue` pins a `-0.dev.` version or the module tree holds a `cue.mod/local-module.cue`.

#### Scenario: Development pin

- **WHEN** the module's `cue.mod/module.cue` pins catalog `v4.6.0-0.dev.3`
- **THEN** the check fails naming the pin

### Requirement: The module release check verifies the operator image reference
The module release check SHALL fail when the module's operator image reference is not `ghcr.io/open-platform-model/opm-operator:<tag>@<digest>`, where `<tag>` is a published, non-draft operator release and `<digest>` is the manifest-list digest GHCR serves under that tag.

#### Scenario: Image digest does not match GHCR

- **WHEN** the module names `opm-operator:v1.0.0-beta.6@sha256:aaa…` and GHCR serves `sha256:bbb…` under `v1.0.0-beta.6`
- **THEN** the check fails naming both digests, and the publish job, if reached, pushes nothing

#### Scenario: Image tag is a draft release

- **WHEN** the module names an operator tag whose GitHub Release is still a draft
- **THEN** the check fails naming the tag

### Requirement: The module release check verifies the module version
The module release check SHALL fail when `identity.Version` differs from the version the release PR proposes, when its major disagrees with the module path's major, or when the proposed version breaks the v0 rules of "Module versions on the v0 path are 0.y.z", "The first module release is 0.1.0" or "The module moves to the v1 path only after operator GA".

#### Scenario: Identity version differs from the proposed version
- **WHEN** `identity.Version` differs from the version the release PR proposes
- **THEN** the module release check fails

### Requirement: The module release check refuses an operator below the minimum version
The module release check SHALL fail when the `<tag>` of the module's operator image reference is older than the repository's minimum operator version, the first operator release that refuses to reconcile the operator's own instance.

#### Scenario: Operator below the minimum

- **WHEN** the module names `opm-operator:v1.0.0-beta.5` and the minimum operator version is `v1.0.0-beta.6`
- **THEN** the check fails naming both tags

### Requirement: The module release check compares generated CRDs and RBAC with the deployed operator
The module release check SHALL fail when the module's generated CRDs or controller RBAC differ from what the operator source at the `<tag>` of the module's operator image generates. The last check (The module release check compares generated CRDs and RBAC with the deployed operator) exists because a module release must render exactly the CRDs the controller it deploys serves and exactly the cluster permissions that controller declares, and `main` can move ahead of the deployed operator release.

#### Scenario: CRDs ahead of the deployed operator

- **WHEN** `main` changed a CRD after `v1.0.0-beta.6` and the module still names `v1.0.0-beta.6` but carries the new CRD
- **THEN** the check fails naming the CRD, and the module waits for the image-bump PR of the next operator release

### Requirement: Module releases are frozen while config differs from the deployed tag
Because the module's generated data follow `main`'s `config/` on every pull request, module releases SHALL be frozen while `main`'s `config/` differs from the deployed tag's, until the image-bump PR of an operator release carrying that change merges.

#### Scenario: Release frozen while config differs
- **WHEN** `main`'s `config/` differs from the deployed tag's
- **THEN** module releases are frozen until the image-bump PR of an operator release carrying that change merges

### Requirement: The module release check runs in the publish job before any push
The module release check SHALL also run in the module publish job before anything is pushed, so a release merged past a red check still publishes nothing.

#### Scenario: Release merged past a red check
- **WHEN** a module release is merged while the module release check is red
- **THEN** the module publish job runs the check before anything is pushed and publishes nothing

### Requirement: An offline test covers the module release checks that need no network
The checks that need no network (development pin, `local-module.cue`, version and major, the `1.0.0` rule, the minimum operator version, reporting every failure) SHALL be covered by an offline test that the `Lint` job runs on every pull request.

#### Scenario: Several failures named at once

- **WHEN** the offline test runs the check over a fixture tree with a development pin and an `identity.Version` that differs from the proposed version
- **THEN** the check exits non-zero and its output names both failures

### Requirement: The module publish job publishes the tagged source
When release-please creates a module release, the module publish job SHALL check out the module's release tag, install the opm CLI named by `.opm-cli-version`, run the module release check, and publish `modules/opm_operator` with `opm module publish <dir> --version <version>`.

#### Scenario: First publish

- **WHEN** the module release `opm_operator-v0.2.0` is created and GHCR holds no `v0.2.0` of `opmodel.dev/modules/opm_operator`
- **THEN** the job publishes `v0.2.0` from the tagged tree and records its manifest digest

### Requirement: The module publish job runs only for a module release
The module publish job SHALL run only when the module package created a release. It SHALL never run for an operator-only release.

#### Scenario: Operator-only release
- **WHEN** release-please creates an operator-only release
- **THEN** the module publish job does not run

### Requirement: The module publish job is idempotent for the same release
Re-running the module publish job for the same release SHALL succeed without pushing when the registry already holds that version with the content this tag produces. It SHALL fail when the registry holds that version with other content. The job SHALL NOT rely on the registry to refuse an overwrite.

#### Scenario: Re-run after a later step failed

- **WHEN** the publish job is re-run after its publish succeeded and a later step of that job failed
- **THEN** the publish step reports the version as already present with the same content and pushes nothing

#### Scenario: Version held by other content

- **WHEN** GHCR already holds `v0.2.0` with a digest that differs from what the tagged tree publishes
- **THEN** the job fails, the release stays a draft, and the fix is the next module version

### Requirement: The module publish job reads registry references from cue mod resolve
Every registry reference the module publish job probes or reads back SHALL be the one `cue mod resolve` gives the module path in `cue.mod/module.cue` under the registry mapping the publish uses, never a hand-spelled repository: CUE keeps the full module path under the mapped prefix, so `opmodel.dev/modules/opm_operator` under `opmodel.dev=ghcr.io/open-platform-model` lands at `ghcr.io/open-platform-model/opmodel.dev/modules/opm_operator`.

#### Scenario: Read back where CUE published

- **WHEN** the job publishes `v0.2.0` under the mapping `opmodel.dev=ghcr.io/open-platform-model`
- **THEN** it reads the digest back from `ghcr.io/open-platform-model/opmodel.dev/modules/opm_operator:v0.2.0`, and the reuse comparison reads the job-local copy from `localhost:5000/opmodel.dev/modules/opm_operator:v0.2.0`

### Requirement: An offline test asserts the registry references of the publish job
An offline test SHALL assert the registry references the module publish job probes or reads back.

#### Scenario: Offline test asserts the references
- **WHEN** the offline test runs
- **THEN** it asserts the registry references the module publish job probes or reads back

### Requirement: The install manifest is the module's render at its defaults
The module publish job SHALL render the published module version from the registry at its default values, as the instance `opm-operator` in the namespace `opm-operator-system`, against the platform generated from the module's own pins and with no cluster.

#### Scenario: Manifest equals the module's render

- **WHEN** the module release `opm_operator-v0.2.0` completes
- **THEN** its GitHub Release carries `install.yaml`, and rendering `opmodel.dev/modules/opm_operator` version `0.2.0` with no values as `opm-operator` in `opm-operator-system` yields the same objects

#### Scenario: Module debug values are not used

- **WHEN** the module declares non-empty `debugValues`
- **THEN** the manifest is still the render at the `#config` defaults

### Requirement: install.yaml is attached to the draft release with the Namespace and CRDs first
The module publish job SHALL write the render to `install.yaml`, with the Namespace and the CRDs before every namespaced object, and attach it to the module's draft release.

#### Scenario: Namespace first

- **WHEN** `install.yaml` is read top to bottom
- **THEN** the Namespace `opm-operator-system` and the four CRDs come before every namespaced object

### Requirement: The render and the upload run apart from the publish
The render and the upload of `install.yaml` SHALL run in a job apart from the publish, so re-running them never re-enters the publish.

#### Scenario: Upload fails after the publish

- **WHEN** the install-manifest upload fails after the module version was published
- **THEN** "Re-run failed jobs" re-runs the render and the upload, never the publish

### Requirement: The install manifest names the operator image the module names
The install manifest SHALL name the operator image by the tag and digest the module names, and the job apart from the publish SHALL check that before the upload.

#### Scenario: Manifest names another image

- **WHEN** the rendered manifest runs an operator image other than the module's tag and digest
- **THEN** the job fails before the upload and the release stays a draft

### Requirement: The install manifest yields a running operator
Applying `install.yaml` with `kubectl apply --server-side` on a cluster with no operator SHALL yield a running operator. A pull request that changes `modules/opm_operator/` or the scripts that render the manifest SHALL prove this on a kind cluster.

#### Scenario: Fresh cluster

- **WHEN** CI applies the manifest rendered from a pull request's module to a kind cluster with no OPM CRDs
- **THEN** the operator's Deployment completes its rollout

### Requirement: A module release stays a draft with one tag until it is recovered or published
The module package SHALL set `draft: true` and `force-tag-creation: true`. Every upload to a module release SHALL first confirm that exactly one release carries the tag and that it is a draft. Recovery before publication is "Re-run failed jobs", or the recovery dispatch when the failure lies in a script at the tag. After publication it is the next module version; a published tag is never moved.

#### Scenario: Upload confirms a single draft
- **WHEN** an upload to a module release starts
- **THEN** it first confirms that exactly one release carries the tag and that it is a draft

### Requirement: The final module publish job publishes the draft only with install.yaml attached
A final module publish job SHALL depend on the publish job and SHALL confirm `install.yaml` is attached before it publishes the draft.

#### Scenario: Asset missing

- **WHEN** the final module publish job finds the draft for `opm_operator-v0.2.0` without `install.yaml`
- **THEN** it fails and the release stays a draft

#### Scenario: Already published

- **WHEN** the final job runs for a module release that is already published
- **THEN** it succeeds without changing anything

### Requirement: A module release is never the latest release
A module release is not a prerelease, and it SHALL be published with GitHub's "latest" mark withheld (`make_latest=false`), so it never becomes the repository's latest release (enhancement 0021, "The operator's install artifact is a registry module on its own version train").

#### Scenario: Never the latest release

- **WHEN** the final module publish job publishes `opm_operator-v0.2.0` and the repository's other releases are operator prereleases
- **THEN** the release is published with `make_latest=false`, and GitHub's latest release of the repository is not `opm_operator-v0.2.0`

### Requirement: The release workflow accepts a module recovery dispatch
The release workflow SHALL accept a `workflow_dispatch` with the input `module_tag`, for a module release whose run left its draft unpublished. It SHALL run only the module jobs, never release-please, the identity advance or an operator job.

#### Scenario: Dispatch runs the module jobs only
- **WHEN** a `workflow_dispatch` with `module_tag` runs
- **THEN** it runs only the module jobs and not release-please, the identity advance or an operator job

### Requirement: The recovery dispatch refuses before any registry login
The recovery dispatch SHALL refuse, before any registry login, unless it runs in `open-platform-model/opm-operator` on `main`, `module_tag` matches `^opm_operator-v[0-9]+\.[0-9]+\.[0-9]+$`, the tag exists and is reachable from `main`, and exactly one release carries it as a draft.

#### Scenario: Release already published

- **WHEN** the dispatch names a module tag whose release is published
- **THEN** it fails before any module job runs

#### Scenario: Not a module tag, or not main

- **WHEN** the dispatch names `v1.0.0-beta.7`, a module tag not reachable from `main`, or runs on a branch other than `main`
- **THEN** it fails before any module job runs

### Requirement: The recovery dispatch runs main's scripts against the tag's tree
In a recovery dispatch, the module jobs SHALL check the tag's tree out and run `main`'s release scripts (`hack/operator-module/` and `.github/scripts/`) against it, so a script fixed after the tag is the one that runs; the release gate and the opm CLI stay the tag's.

#### Scenario: Fixed script runs against the tag
- **WHEN** a script under `hack/operator-module/` or `.github/scripts/` was fixed on `main` after the module tag
- **THEN** the recovery dispatch runs the fixed script against the tag's tree, with the tag's release gate and opm CLI

### Requirement: Each recovery module job gates on its predecessor explicitly
Each module job of a recovery dispatch SHALL gate on its predecessor's result explicitly, never on the implicit success check, which would skip it behind the skipped release-please.

#### Scenario: Module job gates explicitly
- **WHEN** a recovery dispatch runs the module jobs while release-please is skipped
- **THEN** each module job gates on its predecessor's result explicitly and is not skipped by the implicit success check

### Requirement: The recovery dispatch publishes, attaches and releases in order
A recovery dispatch SHALL run, in order: the publish, which reuses a version GHCR already holds only when it equals what the tag's tree publishes, and publishes the tag's tree when GHCR does not hold the version yet (a draft stranded before its first publish); the install manifest rendered from the published version and attached to the draft; then the draft's publication with `make_latest=false`.

#### Scenario: Recover after a broken read-back

- **WHEN** `opm_operator-v0.1.0` was published to GHCR but the run failed reading it back, and the release is a draft without `install.yaml`
- **THEN** a dispatch with `module_tag=opm_operator-v0.1.0` on `main` reuses the held `v0.1.0` after comparing it with the tag's tree, attaches `install.yaml`, and publishes the release with `make_latest=false`

#### Scenario: Stranded before its first publish

- **WHEN** a module release's run failed before it published, and GHCR does not hold the version
- **THEN** the dispatch publishes the tag's tree, attaches `install.yaml`, and publishes the release with `make_latest=false`

#### Scenario: Held version differs from the tag

- **WHEN** the dispatch finds GHCR holding the version with a digest the tag's tree does not publish
- **THEN** the publish fails, the release stays a draft, and the fix is the next module version

### Requirement: The recovery dispatch runs in a concurrency group per tag
The recovery dispatch SHALL run in a concurrency group of its own per tag, so it never replaces or is replaced by a pending push run.

#### Scenario: Dispatch has its own concurrency group
- **WHEN** a recovery dispatch for a module tag runs while a push run is pending
- **THEN** the dispatch runs in a concurrency group of its own for that tag and neither replaces nor is replaced by the push run

### Requirement: A push-triggered release runs the tag's own scripts
A push-triggered release SHALL run the tag's own scripts and start the module jobs only when release-please succeeded and created a module release, as before.

#### Scenario: Push run starts module jobs only after a module release
- **WHEN** a push-triggered release runs
- **THEN** it runs the tag's own scripts and starts the module jobs only when release-please succeeded and created a module release

### Requirement: Recovery never moves a tag or edits a release by hand
Recovery SHALL never move, delete or re-create a tag and SHALL never edit a release by hand.

#### Scenario: Recovery leaves tags and releases alone
- **WHEN** a recovery dispatch runs
- **THEN** it moves, deletes or re-creates no tag and edits no release by hand

### Requirement: Only the module publish job publishes under the module path
No workflow, task or script in this repository other than the module publish job (on a release or its recovery dispatch) SHALL publish to `opmodel.dev/modules/opm_operator` or to its GHCR repository `ghcr.io/open-platform-model/opmodel.dev/modules/opm_operator`.

#### Scenario: Search for publishers

- **WHEN** `.github/workflows/`, `.github/scripts/`, `Taskfile.yml`, `.tasks/` and `hack/` are searched for a publish of the module directory or of `ghcr.io/open-platform-model/opmodel.dev/modules/opm_operator`
- **THEN** the only match is the module publish job and the task it calls

### Requirement: The module publish job refuses outside the module release
The module publish job SHALL refuse to run unless the repository is `open-platform-model/opm-operator` and the checked-out ref is a module release tag.

#### Scenario: Fork or dispatch on another ref

- **WHEN** the module publish job is started in a fork or on a ref that is not a module release tag
- **THEN** it fails before logging in to the registry

### Requirement: Local and test flows publish the module away from GHCR
Local and test flows SHALL publish the module only to a registry that maps `opmodel.dev` away from GHCR.

#### Scenario: Local flow publishes away from GHCR
- **WHEN** a local or test flow publishes the module
- **THEN** it publishes only to a registry that maps `opmodel.dev` away from GHCR

### Requirement: An operator release opens one pull request on module/operator-image
Once an operator release has been published, the release workflow SHALL open or update one pull request on the branch `module/operator-image`. The PR SHALL be created with the release App's token so its CI runs.

#### Scenario: Image PR created with the release App's token
- **WHEN** an operator release has been published
- **THEN** the release workflow opens or updates one pull request on `module/operator-image`, created with the release App's token so its CI runs

### Requirement: The module image PR sets the deployed operator release and nothing else
The module image PR SHALL set the operator release the module deploys to that release's tag and the manifest-list digest GHCR serves under it, and SHALL change nothing else: no file outside `modules/opm_operator/`, and not the module's generated CRD and RBAC files, which follow `main`'s `config/`.

#### Scenario: Operator release published

- **WHEN** `v1.0.0-beta.6` is published with manifest-list digest `sha256:abc…`
- **THEN** a PR from `module/operator-image` sets the module's image to `ghcr.io/open-platform-model/opm-operator:v1.0.0-beta.6@sha256:abc…`, changes only `modules/opm_operator/operator/operator.cue`, and is titled `fix(deps): deploy operator v1.0.0-beta.6 from the operator module`

#### Scenario: Image PR never releases the operator

- **WHEN** the image PR is merged
- **THEN** release-please proposes a module release and no operator release

### Requirement: The module image run fails when config differs from the operator tag
When the CRDs or RBAC roles under `config/` at the operator release tag differ from `main`'s, the module image run SHALL fail naming the files and open no PR, because the module release gate would refuse that module.

#### Scenario: Tag behind main

- **WHEN** the dispatch names `v1.0.0-beta.6` and `main` changed a CRD after that tag
- **THEN** the run fails naming the CRD file and opens no PR

### Requirement: The module image run fails for a tag that does not descend from the deployed one
The module image run SHALL also fail and open no PR when the operator release tag does not descend from the operator tag the module deploys, so a re-run or a dispatch never moves the module back.

#### Scenario: Older operator tag

- **WHEN** the module deploys `v1.0.0-beta.7` and a re-run or dispatch names `v1.0.0-beta.6`
- **THEN** the run fails and opens no PR

### Requirement: The module image PR title carries the change type
The module image PR's title SHALL be `fix(deps): deploy operator <tag> from the operator module`, where `<tag>` is the published operator release. The title SHALL be `fix(deps)!: ...` instead when the operator CHANGELOG sections between the module's previous operator tag and `<tag>` carry a breaking-change entry, or when a CRD under `config/crd/bases` at `<tag>` stops serving a version it served at the module's previous operator tag.

#### Scenario: Breaking operator release

- **WHEN** the operator CHANGELOG section for `v1.0.0-beta.7` lists a breaking change
- **THEN** the PR title carries `!`, and merging it makes the module's next release raise the minor

### Requirement: The module image workflow accepts a dispatch for recovery
As in requirement "An operator release opens one pull request on module/operator-image": A `workflow_dispatch` with a `tag` input SHALL do the same for a published operator release, for recovery.

#### Scenario: Draft operator release

- **WHEN** the dispatch names an operator tag whose release is still a draft
- **THEN** the run fails and opens no PR

### Requirement: The module/operator-image branch is rebuilt only while it holds bot commits
While the branch `module/operator-image` carries only bot commits, a later run SHALL rebuild it from `main`. Once a human commit is on it, the run SHALL add its own commit on top and never rewrite the branch.

#### Scenario: Bot-only image branch is rebuilt
- **WHEN** `module/operator-image` carries only bot commits and the workflow runs again
- **THEN** the run rebuilds the branch from `main`

#### Scenario: Human commit on the image branch
- **WHEN** a human commit is on `module/operator-image` and the workflow runs again
- **THEN** the run adds its own commit on top and does not rewrite the branch
