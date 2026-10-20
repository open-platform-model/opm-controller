## REMOVED Requirements

### Requirement: Pull-request image build trigger and tags
**Reason**: Split into three requirements that each state one behavior.
**Migration**: See "The PR image workflow runs on pull request events and skips bot heads", "The PR image build pushes a sha tag and a PR tag" and "The PR image build does not use the GitHub Actions cache".

### Requirement: Release image tags
**Reason**: Split into five requirements that each state one behavior.
**Migration**: See "A gated release run pushes three tags", "The release image carries the release commit and the release tag", "The version tag always resolves to the manifest list first pushed under it", "A version tag built from the release commit is reused" and "A version tag held by another revision fails the job".

## ADDED Requirements

### Requirement: The PR image workflow runs on pull request events and skips bot heads
The PR image workflow SHALL run on `pull_request` events (types: `opened`, `synchronize`, `reopened`) whatever paths they touch (the workflow has no `paths` filter), except for a head named `deps/cascade` or starting with `release-please--`, `module/` (the operator module's bot heads) or `dependabot/`, where the job SHALL be skipped so that unreviewed bot-head code never builds under a `packages: write` and `id-token: write` token.

#### Scenario: Non-applicable event
- **WHEN** a comment-only or label-only event fires on a pull request (without code changes)
- **THEN** the workflow SHALL NOT run an image build

#### Scenario: Bot head skipped
- **WHEN** a pull request from `deps/cascade`, a `release-please--*` head, a `module/*` head or a `dependabot/*` head opens or updates
- **THEN** the image job SHALL be skipped and SHALL push nothing

### Requirement: The PR image build pushes a sha tag and a PR tag
As in requirement "The PR image workflow runs on pull request events and skips bot heads": On every run it SHALL build and push exactly two tags: `:sha-<short>` where `<short>` is the 7-character short commit SHA of the PR head, and `:pr-<PR_ID>` where `<PR_ID>` is the GitHub pull request number. The `:pr-<PR_ID>` tag MAY be overwritten by subsequent pushes to the same PR; `:sha-<short>` is effectively immutable because the SHA changes when the commit changes.

#### Scenario: New PR opened
- **WHEN** a contributor opens a pull request whose head commit is `abcd123...`
- **THEN** the workflow SHALL push `ghcr.io/open-platform-model/opm-operator:sha-abcd123` and `ghcr.io/open-platform-model/opm-operator:pr-<PR_ID>`

#### Scenario: Force push updates an open PR
- **WHEN** a contributor force-pushes a new head commit `ef01234...` to an existing PR with `PR_ID=42`
- **THEN** the workflow SHALL push `:sha-ef01234` as a new tag AND overwrite `:pr-42` to point at the new image

### Requirement: The PR image build does not use the GitHub Actions cache
As in requirement "The PR image workflow runs on pull request events and skips bot heads": The build SHALL NOT use the GitHub Actions cache.

#### Scenario: The PR image build uses no Actions cache
- **WHEN** the PR image workflow builds an image
- **THEN** the build does not use the GitHub Actions cache

### Requirement: A gated release run pushes three tags
On a gated release run, the image-release job SHALL push, for a release version `v<MAJOR>.<MINOR>.<PATCH>`, `:sha-<short>` (7-character short SHA of the release commit), `:v<MAJOR>.<MINOR>.<PATCH>` (exact release version including leading `v`) and `:latest`. `:latest` and `:sha-<short>` are mutable by design. (Requirement: The release image carries the release commit and the release tag)

#### Scenario: First release v0.1.0
- **WHEN** release-please cuts `v0.1.0` at commit `abcd123...`
- **THEN** the job SHALL push `:sha-abcd123`, `:v0.1.0`, and `:latest` all pointing at the same manifest list

#### Scenario: Subsequent release v0.2.0
- **WHEN** release-please cuts `v0.2.0` at commit `ef01234...`
- **THEN** the job SHALL push `:sha-ef01234`, `:v0.2.0`, and update `:latest` to point at the new manifest list

### Requirement: The release image carries the release commit and the release tag
The release commit is the commit the release tag names, read from the checked-out tag, not the pushed commit; the image's `org.opencontainers.image.revision` label SHALL carry it and its `org.opencontainers.image.version` label SHALL be the release tag.

#### Scenario: Release cut on a later push
- **WHEN** release-please creates the release for a Release PR merged at `abcd123...` during a workflow run triggered by a later push at `ef01234...`
- **THEN** the image SHALL be built from the tag, tagged `:sha-abcd123`, and labelled with revision `abcd123...`

### Requirement: The version tag always resolves to the manifest list first pushed under it
The version tag SHALL always resolve to the manifest list first pushed under it (enhancement 0021, "Release tags are immutable"): before building, the job SHALL probe the version tag. When it is absent the job builds and pushes, then SHALL confirm the version tag resolves to the pushed digest.

#### Scenario: An absent version tag is built, pushed and confirmed
- **WHEN** the job probes the version tag before building and it is absent
- **THEN** the job builds and pushes, then confirms the version tag resolves to the pushed digest

### Requirement: A version tag built from the release commit is reused
As in requirement "The version tag always resolves to the manifest list first pushed under it": When it exists with a revision label equal to the release commit the job SHALL push nothing under the version tag, SHALL reuse the existing digest for every later step, and MAY re-point the mutable tags to it. (Requirement: The release image carries the release commit and the release tag)

#### Scenario: Re-run after the version tag was pushed
- **WHEN** the image-release job is re-run for `v1.0.0-beta.3` after a previous attempt pushed `:v1.0.0-beta.3` built from the release commit
- **THEN** the job SHALL NOT build or push `:v1.0.0-beta.3`, its digest stays unchanged, and the job SHALL sign, attest and render `install.yaml` against that digest

### Requirement: A version tag held by another revision fails the job
As in requirement "The version tag always resolves to the manifest list first pushed under it": When it exists with any other revision, or the probe fails for any reason other than the tag being absent, the job SHALL fail without pushing.

#### Scenario: Version tag held by a different image
- **WHEN** `:v1.0.0-beta.3` already exists and its revision label names a commit other than the release commit
- **THEN** the job SHALL fail without pushing any tag, and the release stays a draft
