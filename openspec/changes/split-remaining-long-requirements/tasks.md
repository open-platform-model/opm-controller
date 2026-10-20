# Tasks: split-remaining-long-requirements

Branch `docs/split-remaining-long-requirements`, created from `origin/main`. Run every command inside its worktree. The gate of `design.md` runs in a scratch copy of `openspec/`, never in the worktree. Every command that runs `openspec` uses version 1.14.1 and `OPENSPEC_TELEMETRY=0`.

Sections 1 to 8 edit spec text only, so no Go file, CRD or generated file changes and `task dev:fmt dev:vet dev:lint dev:test` has nothing to check there; the section's gate is the gate of `design.md`, and its review is the 14 items of the guide. `N` in a gate is the number of WARNING lines the branch has removed from `validate --all --strict`, counted against `origin/main` (86 WARNING lines at the start).

Each rewrite unit starts by writing its scenario ledger into `design.md` (heading `### Ledger: <spec>`), reading `git show origin/main:openspec/specs/<spec>/spec.md` and running `grep -rn --exclude-dir=archive "<old header>" .` before any split. It adds every keyword it gains to the strength table of `design.md`.

## 1. Purpose of source-resolution

- [ ] 1.1 Add `## Purpose` (at least 50 characters, text read from the requirements of the spec) above `## Requirements` of `openspec/specs/source-resolution/spec.md`, by direct edit. Record it in `design.md`.
- [ ] 1.2 Gate (`N=0`): no new ERROR item, the ERROR set loses the `spec source-resolution file` item, WARNING lines unchanged.
- [ ] 1.3 Commit `docs(openspec): add a Purpose to source-resolution`

## 2. Open-change fixes (deltas only)

Edit only files under `openspec/changes/<open change>/specs/`. Never the proposal, design, tasks or code of an open change. Write the ledger `### Ledger: <open change>, <spec>` before each edit.

- [ ] 2.1 `stop-operator-install-manifest`, `operator-module`: restore scenario "An object added to the kustomize tree only fails the test" in the MODIFIED block of "The controller's pod and Service stay those of the kustomize tree", byte equal to `origin/main`; split the body so that every block is 500 characters or less.
- [ ] 2.2 `stop-operator-install-manifest`, `release-automation`: split the MODIFIED blocks "Release published once after every release job" and "Beta prerelease line"; keep every live scenario name and text.
- [ ] 2.3 `add-cross-namespace-source-grants`: split the ADDED blocks "SourceRefGrant CRD" and "Grant-matching policy permits a reference", and the MODIFIED block "OCIRepository lookup" of `source-resolution` (the delta text is over 500).
- [ ] 2.4 Byte-equality command for every edited block (MODIFIED against `origin/main`, ADDED against the delta on `origin/main`): no `DIFF`, no `LOST`.
- [ ] 2.5 Gate (`N=2`): both open changes pass `validate <change> --strict`; the archive simulation of design.md ("Gate for the open-change edits"), orders A and B, exits 0; `validate --all --strict` after the simulated archives and the simulated Purpose edit exits 0; no new ERROR item.
- [ ] 2.6 Commit `docs(openspec): fix the spec deltas of the two open changes`

## 3. Unit 1: finalizer-and-deletion, instance-health: 15 requirements

Specs: finalizer-and-deletion 11, instance-health 4.

- [ ] 3.1 For each long requirement pick the form (design.md), write the ledger, write the delta under `specs/<spec>/spec.md` of this change. One delta file per spec; no spec of another unit.
- [ ] 3.2 Byte-equality command of the guide (review item 3) for every rewritten requirement: no `DIFF`, no `LOST`.
- [ ] 3.3 Gate (`N=17`): `validate split-remaining-long-requirements --strict` exit 0, archive in the copy exit 0, no "Archive would refuse", no new ERROR item, WARNING lines fell by `N`; `validate <spec> --type spec --strict` shows no WARNING for a rewritten requirement.
- [ ] 3.4 Commit `docs(openspec): split the long requirements of finalizer-and-deletion and instance-health`

## 4. Unit 2: release-automation, container-image-publish: 13 requirements

Specs: release-automation 11, container-image-publish 2. The open change `stop-operator-install-manifest` MODIFIES two and REMOVES one long requirement of these specs; those three are not touched here (design.md, overlap list)..

- [ ] 4.1 For each long requirement pick the form (design.md), write the ledger, write the delta under `specs/<spec>/spec.md` of this change. One delta file per spec; no spec of another unit.
- [ ] 4.2 Byte-equality command of the guide (review item 3) for every rewritten requirement: no `DIFF`, no `LOST`.
- [ ] 4.3 Gate (`N=30`): `validate split-remaining-long-requirements --strict` exit 0, archive in the copy exit 0, no "Archive would refuse", no new ERROR item, WARNING lines fell by `N`; `validate <spec> --type spec --strict` shows no WARNING for a rewritten requirement.
- [ ] 4.4 Commit `docs(openspec): split the long requirements of release-automation and container-image-publish`

## 5. Unit 3: reconcile-loop-assembly, operator-module: 13 requirements

Specs: reconcile-loop-assembly 9, operator-module 4. The open change MODIFIES one long requirement of `operator-module` ("The controller's pod and Service stay those of the kustomize tree"); it is not touched here..

- [ ] 5.1 For each long requirement pick the form (design.md), write the ledger, write the delta under `specs/<spec>/spec.md` of this change. One delta file per spec; no spec of another unit.
- [ ] 5.2 Byte-equality command of the guide (review item 3) for every rewritten requirement: no `DIFF`, no `LOST`.
- [ ] 5.3 Gate (`N=43`): `validate split-remaining-long-requirements --strict` exit 0, archive in the copy exit 0, no "Archive would refuse", no new ERROR item, WARNING lines fell by `N`; `validate <spec> --type spec --strict` shows no WARNING for a rewritten requirement.
- [ ] 5.4 Commit `docs(openspec): split the long requirements of reconcile-loop-assembly and operator-module`

## 6. Unit 4: registration-acceptance, cascade-receiver: 14 requirements

Specs: registration-acceptance 9, cascade-receiver 5.

- [ ] 6.1 For each long requirement pick the form (design.md), write the ledger, write the delta under `specs/<spec>/spec.md` of this change. One delta file per spec; no spec of another unit.
- [ ] 6.2 Byte-equality command of the guide (review item 3) for every rewritten requirement: no `DIFF`, no `LOST`.
- [ ] 6.3 Gate (`N=57`): `validate split-remaining-long-requirements --strict` exit 0, archive in the copy exit 0, no "Archive would refuse", no new ERROR item, WARNING lines fell by `N`; `validate <spec> --type spec --strict` shows no WARNING for a rewritten requirement.
- [ ] 6.4 Commit `docs(openspec): split the long requirements of registration-acceptance and cascade-receiver`

## 7. Unit 5: example-test-modules, prune-stale-resources: 13 requirements

Specs: example-test-modules 7, prune-stale-resources 6.

- [ ] 7.1 For each long requirement pick the form (design.md), write the ledger, write the delta under `specs/<spec>/spec.md` of this change. One delta file per spec; no spec of another unit.
- [ ] 7.2 Byte-equality command of the guide (review item 3) for every rewritten requirement: no `DIFF`, no `LOST`.
- [ ] 7.3 Gate (`N=70`): `validate split-remaining-long-requirements --strict` exit 0, archive in the copy exit 0, no "Archive would refuse", no new ERROR item, WARNING lines fell by `N`; `validate <spec> --type spec --strict` shows no WARNING for a rewritten requirement.
- [ ] 7.4 Commit `docs(openspec): split the long requirements of example-test-modules and prune-stale-resources`

## 8. Unit 6: drift-detection, modulepackage-kernel-rendering, example-module-publishing: 12 requirements

Specs: drift-detection 5, modulepackage-kernel-rendering 5, example-module-publishing 2.

- [ ] 8.1 For each long requirement pick the form (design.md), write the ledger, write the delta under `specs/<spec>/spec.md` of this change. One delta file per spec; no spec of another unit.
- [ ] 8.2 Byte-equality command of the guide (review item 3) for every rewritten requirement: no `DIFF`, no `LOST`.
- [ ] 8.3 Gate (`N=82`): `validate split-remaining-long-requirements --strict` exit 0, archive in the copy exit 0, no "Archive would refuse", no new ERROR item, WARNING lines fell by `N`; `validate <spec> --type spec --strict` shows no WARNING for a rewritten requirement.
- [ ] 8.4 Commit `docs(openspec): split the long requirements of drift-detection, modulepackage-kernel-rendering and example-module-publishing`

## 9. The pinned strict gate (last group)

- [ ] 9.1 `Taskfile.yml`: add the variable `OPENSPEC_VERSION: 1.14.1` with a comment that it is the one pin, and the tasks `openspec:install` (`npm install -g --ignore-scripts @fission-ai/openspec@{{.OPENSPEC_VERSION}}`) and `openspec:check` (preconditions: `openspec` on PATH and its version equals the pin; env `OPENSPEC_TELEMETRY: "0"`; command `openspec validate --all --strict --no-interactive`), in the form of core's `Taskfile.yml` on `origin/main`.
- [ ] 9.2 `.github/workflows/lint.yml`: add to the `Lint` job a setup-node step (action pinned by commit SHA, as the other actions of the file), `task openspec:install` and `task openspec:check`. Add nothing to a workflow that holds a write permission. Run `task cascade:wiring:check` and `zizmor` on the file if the repo's tools list them.
- [ ] 9.3 Final gate in a scratch copy of the branch: `validate split-remaining-long-requirements --strict` exit 0; the archive simulation of design.md, orders A and B, ends with `validate --all --strict` exit 0 after the simulated Purpose edit of `cross-namespace-source-grants`; `task openspec:check` on that copy exits 0; on the branch itself it exits 1 only on the four overlap requirements of design.md. Report the WARNING count (86 at the start, 4 at the end).
- [ ] 9.4 Archive this change on the branch (`openspec-archive-change`), so the archive commit rides the PR. Commit `ci(openspec): gate validate --all --strict with a pinned OpenSpec 1.14.1`. Delivery: if the human chooses two PRs (design.md, "Delivery order"), this group goes in the second PR, after both open changes have archived, and the archive commit goes with it.
