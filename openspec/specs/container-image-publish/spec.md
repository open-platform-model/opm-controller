## Purpose

Publish the controller manager container image to GHCR from CI, with distinct PR-validation and release-gated workflows, cosign keyless signatures, SBOM and SLSA provenance on releases, with the install manifest left to the controller module's releases.

## Requirements

### Requirement: Image registry and repository

The CI pipeline SHALL publish the controller manager container image to `ghcr.io/open-platform-model/opm-controller`. No other registry or repository path is used by this capability.

#### Scenario: PR build push target
- **WHEN** the PR image workflow builds and pushes an image
- **THEN** the image reference SHALL begin with `ghcr.io/open-platform-model/opm-controller:`

#### Scenario: Release build push target
- **WHEN** the release image job builds and pushes an image
- **THEN** the image reference SHALL begin with `ghcr.io/open-platform-model/opm-controller:`

### Requirement: Release image build trigger gated on release-please

The release image job SHALL only execute when the `release-please` job in the same workflow run reports that the controller package `"."` created a release (`outputs.release_created == 'true'`), and SHALL take the release tag from that package's `outputs.tag_name`. A release created only by the controller module's package SHALL NOT run it. When a Release PR is opened or updated without a release being cut, the image job SHALL be skipped.

#### Scenario: Release cut after merge of Release PR
- **WHEN** a push to `main` causes release-please to tag version `v1.2.3` and create a GitHub release
- **THEN** the image-release job SHALL run and publish the image

#### Scenario: Push to main opens or updates a Release PR only
- **WHEN** a push to `main` causes release-please to open or update a Release PR (but not cut a release)
- **THEN** the image-release job SHALL be skipped and its status SHALL be reported as skipped (not failed) in the GitHub Actions UI

#### Scenario: Push to main with no releasable commits
- **WHEN** a push to `main` contains only `chore`, `docs`, `test`, `ci`, or `build` commits and release-please takes no action
- **THEN** the image-release job SHALL be skipped

#### Scenario: Module release only
- **WHEN** a push to `main` merges the module's release PR and release-please creates only `opm_controller-v0.2.0`
- **THEN** the image-release job SHALL be skipped, and no image tag is pushed

### Requirement: Multi-architecture support

PR image builds SHALL produce a multi-architecture manifest list covering `linux/amd64` and `linux/arm64`, the same set as a release. Release image builds SHALL produce a multi-architecture manifest list covering `linux/amd64` and `linux/arm64`.

#### Scenario: PR build architecture set
- **WHEN** the PR image workflow runs
- **THEN** the resulting pushed manifest list SHALL expose exactly two platform descriptors: `linux/amd64` and `linux/arm64`

#### Scenario: Release build architecture set
- **WHEN** the release image-release job runs successfully
- **THEN** the resulting pushed manifest list SHALL expose exactly two platform descriptors: `linux/amd64` and `linux/arm64`

### Requirement: Cosign keyless signatures

Every image manifest pushed by either the PR or release workflow SHALL be signed with cosign using Sigstore's keyless OIDC flow (Fulcio-issued certificate, Rekor transparency log entry). No long-lived cosign private key is stored in the repository or in GitHub secrets.

#### Scenario: PR image signature
- **WHEN** the PR image workflow successfully pushes an image
- **THEN** a cosign signature keyed by the image's SHA256 digest SHALL be published alongside the image in GHCR AND a Rekor transparency log entry SHALL exist referencing the GitHub Actions OIDC token

#### Scenario: Release image signature
- **WHEN** the release image-release job successfully pushes a manifest list
- **THEN** a cosign signature SHALL be published for the manifest list digest

#### Scenario: Signature verification recipe
- **WHEN** a consumer runs `cosign verify` against a release image with the repository's documented `--certificate-identity-regexp` and `--certificate-oidc-issuer=https://token.actions.githubusercontent.com`
- **THEN** verification SHALL succeed for images produced by the release workflow

### Requirement: SBOM and SLSA provenance on release builds

Release image-release job SHALL attach, in addition to the cosign signature, an SPDX-JSON software bill of materials and a SLSA build provenance attestation for the manifest list digest. PR image builds SHALL NOT produce SBOM or provenance attestations (signature only).

#### Scenario: Release SBOM attached
- **WHEN** the release image-release job completes successfully
- **THEN** an SPDX-JSON SBOM attestation SHALL be discoverable via `cosign download attestation` against the manifest-list digest

#### Scenario: Release provenance attached
- **WHEN** the release image-release job completes successfully
- **THEN** a SLSA provenance predicate SHALL be discoverable via the repository's attestations UI AND via `cosign download attestation`

#### Scenario: PR builds omit heavy attestations
- **WHEN** the PR image workflow completes successfully
- **THEN** no SBOM or SLSA provenance attestation SHALL be generated for the PR image

### Requirement: Workflow permissions and action pinning

The PR image workflow and the release image-release job SHALL declare the minimal permissions required: `contents: read`, `packages: write`, `id-token: write`, and `attestations: write` (release job only). All third-party GitHub Actions SHALL be pinned by full commit SHA, not by floating tag.

#### Scenario: PR workflow permissions
- **WHEN** the PR image workflow is defined
- **THEN** it SHALL declare `contents: read`, `packages: write`, and `id-token: write` (no `attestations: write`, no `contents: write`)

#### Scenario: Release job permissions
- **WHEN** the image-release job is defined
- **THEN** it SHALL declare `contents: read`, `packages: write`, `id-token: write`, and `attestations: write`, and no `contents: write`

#### Scenario: Action pinning
- **WHEN** any third-party action is referenced in the PR workflow or the image-release job
- **THEN** it SHALL be referenced by a 40-character commit SHA, never by `@v1`, `@v5`, `@main`, or any other floating ref

### Requirement: The PR image workflow runs on pull request events and skips bot heads
The PR image workflow SHALL run on `pull_request` events (types: `opened`, `synchronize`, `reopened`) whatever paths they touch (the workflow has no `paths` filter), except for a head named `deps/cascade` or starting with `release-please--`, `module/` (the controller module's bot heads) or `dependabot/`, where the job SHALL be skipped so that unreviewed bot-head code never builds under a `packages: write` and `id-token: write` token.

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
- **THEN** the workflow SHALL push `ghcr.io/open-platform-model/opm-controller:sha-abcd123` and `ghcr.io/open-platform-model/opm-controller:pr-<PR_ID>`

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
The version tag SHALL always resolve to the manifest list first pushed under it (enhancement 0021, "Release tags are immutable"): before building, the job SHALL probe the version tag. When it is absent the job builds and pushes, then SHALL confirm the version tag resolves to the pushed digest. (Requirement: A gated release run pushes three tags)

#### Scenario: An absent version tag is built, pushed and confirmed
- **WHEN** the job probes the version tag before building and it is absent
- **THEN** the job builds and pushes, then confirms the version tag resolves to the pushed digest

### Requirement: A version tag built from the release commit is reused
As in requirement "The version tag always resolves to the manifest list first pushed under it": When it exists with a revision label equal to the release commit the job SHALL push nothing under the version tag, SHALL reuse the existing digest for every later step, and MAY re-point the mutable tags to it. (Requirement: The release image carries the release commit and the release tag) (Requirement: A gated release run pushes three tags)

#### Scenario: Re-run after the version tag was pushed
- **WHEN** the image-release job is re-run for `v1.0.0-beta.3` after a previous attempt pushed `:v1.0.0-beta.3` built from the release commit
- **THEN** the job SHALL NOT build or push `:v1.0.0-beta.3`, its digest stays unchanged, and the job SHALL sign and attest that digest

### Requirement: A version tag held by another revision fails the job
As in requirement "The version tag always resolves to the manifest list first pushed under it": When it exists with any other revision, or the probe fails for any reason other than the tag being absent, the job SHALL fail without pushing. (Requirement: The release image carries the release commit and the release tag)

#### Scenario: Version tag held by a different image
- **WHEN** `:v1.0.0-beta.3` already exists and its revision label names a commit other than the release commit
- **THEN** the job SHALL fail without pushing any tag, and the release stays a draft

### Requirement: The installer target stays unchanged in behavior
The existing `task controller:installer` target SHALL remain unchanged in behavior and its default invocation (no `IMG` override) SHALL still produce a manifest using `controller:latest`.

#### Scenario: Default build-installer unchanged
- **WHEN** a developer runs `task controller:installer` locally without an `IMG` override
- **THEN** the rendered `dist/install.yaml` SHALL continue to reference `controller:latest` as before this change
