# Tasks: split-remaining-long-requirements

Branch `docs/split-remaining-long-requirements`, created from `origin/main`. Run every command inside its worktree. The gate of `design.md` runs in a scratch copy of `openspec/`, never in the worktree. Every command that runs `openspec` uses version 1.14.1 and `OPENSPEC_TELEMETRY=0`.

Sections 1 to 9 edit spec text only, so no Go file, CRD or generated file changes and `task dev:fmt dev:vet dev:lint dev:test` has nothing to check there; the section's gate is the gate of `design.md`, and its review is the 14 items of the guide. `N` in a gate is the number of WARNING lines the branch has removed from `validate --all --strict`, counted against `origin/main` (86 WARNING lines at the start).

Each rewrite unit starts by writing its scenario ledger into `design.md` (heading `### Ledger: <spec>`), reading `git show origin/main:openspec/specs/<spec>/spec.md` and running `grep -rn --exclude-dir=archive "<old header>" .` before any split. It adds every keyword it gains to the strength table of `design.md`.

## 1. Purpose of source-resolution

- [x] 1.1 Add `## Purpose` (at least 50 characters, text read from the requirements of the spec) above `## Requirements` of `openspec/specs/source-resolution/spec.md`, by direct edit. Record it in `design.md`.
- [x] 1.2 Gate (`N=0`): no new ERROR item, the ERROR set loses the `spec source-resolution file` item, WARNING lines unchanged.
- [x] 1.3 Commit `docs(openspec): add a Purpose to source-resolution`

## 2. Fixes of add-cross-namespace-source-grants (deltas and one task line)

Edit only files under `openspec/changes/add-cross-namespace-source-grants/specs/`, and add the one task line of 2.4 to its `tasks.md`. Never its proposal, design or code. Write the ledger `### Ledger: add-cross-namespace-source-grants, <spec>` before each edit.

- [x] 2.1 `cross-namespace-source-grants`: split the ADDED blocks "SourceRefGrant CRD" and "Grant-matching policy permits a reference" into ADDED blocks of 500 characters or less.
- [x] 2.2 `source-resolution`: split the MODIFIED block "OCIRepository lookup" (the delta text is over 500): keep the header and every live scenario byte for byte; the split-off behaviors become ADDED requirements.
- [x] 2.3 Byte-equality command for every edited block (against the delta on `origin/main`, MODIFIED against `origin/main`): no `DIFF`, no `LOST`.
- [x] 2.4 Add one task line to `tasks.md` of `add-cross-namespace-source-grants`, at the end of its last section: at archive, write the real Purpose of the new spec `cross-namespace-source-grants`, because the strict gate fails on the placeholder Purpose that archive writes. This is the only `tasks.md` edit of an open change.
- [x] 2.5 Gate (`N=2`): `validate add-cross-namespace-source-grants --strict` exits 0; `validate --all --strict` has no new ERROR item and the WARNING lines fell by 2.
- [x] 2.6 Commit `docs(openspec): fix the spec deltas of add-cross-namespace-source-grants`

## 3. Unit 1: finalizer-and-deletion, instance-health: 15 requirements

Specs: finalizer-and-deletion 11, instance-health 4.

- [x] 3.1 For each long requirement pick the form (design.md), write the ledger, write the delta under `specs/<spec>/spec.md` of this change. One delta file per spec; no spec of another unit.
- [x] 3.2 Byte-equality command of the guide (review item 3) for every rewritten requirement: no `DIFF`, no `LOST`.
- [x] 3.3 Gate (`N=17`): `validate split-remaining-long-requirements --strict` exit 0, archive in the copy exit 0, no "Archive would refuse", no new ERROR item, WARNING lines fell by `N`; `validate <spec> --type spec --strict` shows no WARNING for a rewritten requirement.
- [x] 3.4 Commit `docs(openspec): split the long requirements of finalizer-and-deletion and instance-health`

## 4. Unit 2: release-automation, container-image-publish: 13 requirements

Specs: release-automation 11, container-image-publish 2. The three long requirements that `stop-operator-install-manifest` MODIFIES or REMOVES in these specs are not in this unit: group 9 splits them (design.md, overlap list).

- [x] 4.1 For each long requirement pick the form (design.md), write the ledger, write the delta under `specs/<spec>/spec.md` of this change. One delta file per spec; no spec of another unit.
- [x] 4.2 Byte-equality command of the guide (review item 3) for every rewritten requirement: no `DIFF`, no `LOST`.
- [x] 4.3 Gate (`N=30`): `validate split-remaining-long-requirements --strict` exit 0, archive in the copy exit 0, no "Archive would refuse", no new ERROR item, WARNING lines fell by `N`; `validate <spec> --type spec --strict` shows no WARNING for a rewritten requirement.
- [x] 4.4 Commit `docs(openspec): split the long requirements of release-automation and container-image-publish`

## 5. Unit 3: reconcile-loop-assembly, operator-module: 13 requirements

Specs: reconcile-loop-assembly 9, operator-module 4. The long requirement that `stop-operator-install-manifest` MODIFIES in `operator-module` ("The controller's pod and Service stay those of the kustomize tree") is not in this unit: group 9 splits it (design.md, overlap list).

- [x] 5.1 For each long requirement pick the form (design.md), write the ledger, write the delta under `specs/<spec>/spec.md` of this change. One delta file per spec; no spec of another unit.
- [x] 5.2 Byte-equality command of the guide (review item 3) for every rewritten requirement: no `DIFF`, no `LOST`.
- [x] 5.3 Gate (`N=43`): `validate split-remaining-long-requirements --strict` exit 0, archive in the copy exit 0, no "Archive would refuse", no new ERROR item, WARNING lines fell by `N`; `validate <spec> --type spec --strict` shows no WARNING for a rewritten requirement.
- [x] 5.4 Commit `docs(openspec): split the long requirements of reconcile-loop-assembly and operator-module`

## 6. Unit 4: registration-acceptance, cascade-receiver: 14 requirements

Specs: registration-acceptance 9, cascade-receiver 5.

- [x] 6.1 For each long requirement pick the form (design.md), write the ledger, write the delta under `specs/<spec>/spec.md` of this change. One delta file per spec; no spec of another unit.
- [x] 6.2 Byte-equality command of the guide (review item 3) for every rewritten requirement: no `DIFF`, no `LOST`.
- [x] 6.3 Gate (`N=57`): `validate split-remaining-long-requirements --strict` exit 0, archive in the copy exit 0, no "Archive would refuse", no new ERROR item, WARNING lines fell by `N`; `validate <spec> --type spec --strict` shows no WARNING for a rewritten requirement.
- [x] 6.4 Commit `docs(openspec): split the long requirements of registration-acceptance and cascade-receiver`

## 7. Unit 5: example-test-modules, prune-stale-resources: 13 requirements

Specs: example-test-modules 7, prune-stale-resources 6.

- [x] 7.1 For each long requirement pick the form (design.md), write the ledger, write the delta under `specs/<spec>/spec.md` of this change. One delta file per spec; no spec of another unit.
- [x] 7.2 Byte-equality command of the guide (review item 3) for every rewritten requirement: no `DIFF`, no `LOST`.
- [x] 7.3 Gate (`N=70`): `validate split-remaining-long-requirements --strict` exit 0, archive in the copy exit 0, no "Archive would refuse", no new ERROR item, WARNING lines fell by `N`; `validate <spec> --type spec --strict` shows no WARNING for a rewritten requirement.
- [x] 7.4 Commit `docs(openspec): split the long requirements of example-test-modules and prune-stale-resources`

## 8. Unit 6: drift-detection, modulepackage-kernel-rendering, example-module-publishing: 12 requirements

Specs: drift-detection 5, modulepackage-kernel-rendering 5, example-module-publishing 2.

- [x] 8.1 For each long requirement pick the form (design.md), write the ledger, write the delta under `specs/<spec>/spec.md` of this change. One delta file per spec; no spec of another unit.
- [x] 8.2 Byte-equality command of the guide (review item 3) for every rewritten requirement: no `DIFF`, no `LOST`.
- [x] 8.3 Gate (`N=82`): `validate split-remaining-long-requirements --strict` exit 0, archive in the copy exit 0, no "Archive would refuse", no new ERROR item, WARNING lines fell by `N`; `validate <spec> --type spec --strict` shows no WARNING for a rewritten requirement.
- [x] 8.4 Commit `docs(openspec): split the long requirements of drift-detection, modulepackage-kernel-rendering and example-module-publishing`

## 9. The four overlap requirements and the retarget of stop-operator-install-manifest: 4 requirements

Requirements: operator-module "The controller's pod and Service stay those of the kustomize tree" (MODIFIED by the open change), release-automation "Release published once after every release job" and "Beta prerelease line" (MODIFIED), container-image-publish "Release install manifest with digest-pinned image" (REMOVED). "Release assets upload only to a draft release" is 488 and stays. Main specs and the delta move in the same commit.

- [x] 9.1 For each of the four, pick the form (design.md) and write the ledger, then write the delta under `specs/<spec>/spec.md` of this change, as in units 1 to 6. One delta file per spec: add these requirements to the delta file the unit of the spec wrote, or create it.
- [x] 9.2 Retarget `openspec/changes/stop-operator-install-manifest/specs/` (design.md, "the overlap list"): each MODIFIED block of a split requirement becomes one MODIFIED block per new header, its delta text cut at the same boundary as the main split, verbatim, each part 500 or less; each live scenario goes under the header that holds it on main, byte for byte. A part with no counterpart in the delta text keeps its main text: leave it unmodified or copy it verbatim as a MODIFIED block, and record which and why in the ledger. The REMOVED block of `container-image-publish` becomes one REMOVED block per new header. Restore the dropped scenario "An object added to the kustomize tree only fails the test", byte equal to `origin/main`. Edit nothing else of that change.
- [x] 9.3 Byte-equality command (review item 3) for every rewritten requirement and every edited delta part: no `DIFF`, no `LOST`.
- [x] 9.4 Gate (`N=86`): `validate split-remaining-long-requirements --strict` and `validate stop-operator-install-manifest --strict` exit 0; `validate --all --strict` exits 0 on the branch copy (no WARNING, no ERROR); the archive check of design.md ("Gate for the open-change edits"): archive this change, then `stop-operator-install-manifest`, then `add-cross-namespace-source-grants`, each exit 0 with no "would refuse", and the final `validate --specs --strict` shows no WARNING except the placeholder Purpose of `cross-namespace-source-grants`.
- [x] 9.5 Commit `docs(openspec): split the requirements that stop-operator-install-manifest changes and retarget its delta`

## 10. The pinned strict gate (last group)

- [x] 10.1 `Taskfile.yml`: add the variable `OPENSPEC_VERSION: 1.14.1` with a comment that it is the one pin, and the tasks `openspec:install` (`npm install -g --ignore-scripts @fission-ai/openspec@{{.OPENSPEC_VERSION}}`) and `openspec:check` (preconditions: `openspec` on PATH and its version equals the pin; env `OPENSPEC_TELEMETRY: "0"`; command `openspec validate --all --strict --no-interactive`), in the form of core's `Taskfile.yml` on `origin/main`.
- [x] 10.2 `.github/workflows/lint.yml`: add a job `openspec` (its own job, not a step of the `Lint` job: a timeout, a checkout with `persist-credentials: false`, the Task setup, a setup-node step with the action pinned by commit SHA, `task openspec:install`, `task openspec:check`). Add nothing to a workflow that holds a write permission. Run `task cascade:wiring:check` and `zizmor` on the file if the repo's tools list them.
- [ ] 10.3 Final gate in a scratch copy of the branch: `validate split-remaining-long-requirements --strict` exit 0; the archive check of design.md (this change, `stop-operator-install-manifest`, `add-cross-namespace-source-grants`, each exit 0, no "would refuse"; the final `validate --specs --strict` shows no WARNING except the placeholder Purpose of `cross-namespace-source-grants`, and exits 0 after one simulated edit of it); `task openspec:check` on the branch itself exits 0. Report the WARNING count (86 at the start, 0 at the end).
- [ ] 10.4 Archive this change on the branch (`openspec-archive-change`), so the archive commit rides the PR. Commit `ci(openspec): gate validate --all --strict with a pinned OpenSpec 1.14.1`. The gate group stays last, and `openspec validate --all --strict` exits 0 on the branch at the end. One PR; no second PR.
