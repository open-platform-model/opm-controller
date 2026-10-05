## Context

The owner's rename plan (decisions D1 to D16, 2026-10-05) renames opm-operator to opm-controller across the workspace. This change is the controller repository's part. It is built in a worktree now and merged in one cutover sitting with the other repositories' changes, after the GitHub repository is renamed.

## Goals / Non-Goals

**Goals:**

- One name in every identifier that is not immutable history, with no alias, migration or warning (D2).
- The controller module starts clean on `opmodel.dev/modules/opm_controller@v0` at 0.1.0.
- A controller release stops publishing an install manifest of its own (D14).

**Non-Goals:**

- Rewriting history: CHANGELOG entries, `openspec/changes/archive/**` and git tags stay as written (D11).
- The cli, `.github`, the site and every other consumer (their own changes).
- The annotation `opm.dev/force-delete-orphan`.

## Decisions

### Reconcile phase impact

None beyond names. Source, Render, Apply and Prune behave as before. Status writes `controllerVersion` instead of `operatorVersion` and the render input key carries `ControllerVersion` in the same position, so no recorded digest changes. The self-skip gate compares against the new coordinates:

```go
const (
	ownInstanceName       = "opm-controller"
	ownInstanceNamespace  = "opm-controller-system"
	ownInstanceModulePath = "opmodel.dev/modules/opm_controller"
)
```

The API enum is a hard rename: `+kubebuilder:validation:Enum=cli;controller`. An absent owner still means controller-managed.

### Base name `manager` (D4)

`config/manager/manager.yaml` and `config/rbac/service_account.yaml` name the Deployment and ServiceAccount `manager`, so the prefix `opm-controller-` gives `opm-controller-manager`; the module derives `"\(P)-manager"`. The pod label `control-plane: controller-manager` and the module's component key stay, so the e2e selectors and the module's fixed Deployment selector keep their shape.

### Release-please

The module key leaves `.release-please-manifest.json`, so release-please treats `modules/opm_controller` as a new package and proposes `initial-version` 0.1.0. The root keeps `bootstrap-sha` and gets no `last-release-sha` (a global setting that would truncate root commits). The root excludes both `modules/opm_controller` and `modules/opm_operator` for one release, so commits that touched only the old path stay out of the next controller changelog.

### Image pin of the module

`controller/controller.cue` keeps the version and digest of the last release published under the old image name, marked `TODO(cutover)`, because no `opm-controller` image exists until the first release after the merge. The module image PR of that release replaces both. `min-controller-version` moves to that release in a separate hidden-type `ci(module):` commit after the bot PR merges, because the bot branch's tar allow-path accepts only `modules/opm_controller/`; the same follow-up drops the `TODO(cutover)` note.

## Risks / Trade-offs

- [The `controller-module.yml` install job pulls `ghcr.io/open-platform-model/opm-controller:v1.0.0-beta.8`, which does not exist] → The job is red on this PR and green after the first controller release and its module image PR; it is not a required check.
- [Specs retitled by REMOVED + ADDED move to the end of their files] → The archive commit restores the original order by hand, with identical requirement text.
- [Intermediate commits are not each green] → The PR is squash-merged as one `refactor!:` commit; the branch's commits are history for review.

## Research & Decisions

### Capability rename mechanics
**Context**: OpenSpec 1.12 has no capability rename; a delta to a missing spec accepts only ADDED, and MODIFIED refuses to drop a scenario.
**Explored**: INVENTORY "OpenSpec mechanics" and gap report `gap-r2-openspec-capability-rename` (rename plan directory).
**Decision**: Recreate each `operator-*` capability as `controller-*` (Purpose + ADDED) with the old one REMOVED and `retire_capabilities: true`; a requirement whose title or scenario names the operator is REMOVED + ADDED under a new title; a body-only change is MODIFIED.
**Rationale**: The canonical recipe keeps an archive trail and is what the archive engine accepts.

### Folding stop-operator-install-manifest
**Context**: The blocked change waited on the cli, the resolver and the module release.
**Explored**: cli `origin/main` installs from the module and resolves only the module pin; the `.github` resolver checks an asset only when a caller passes `--asset`, and no caller does for the controller.
**Decision**: Fold its asset removal and its `controller-module` rewording; drop its install-page task (legacy text, D2) and its notify change (no module notification exists to replace it).
**Rationale**: The gates hold, and the rename is the cheapest moment to drop an asset no consumer reads.
