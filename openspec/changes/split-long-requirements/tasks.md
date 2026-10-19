# Tasks: split-long-requirements

Branch `docs/split-long-requirements`, created from `origin/main`. Run every command inside its worktree. The gate of `design.md` runs in a scratch copy of `openspec/`, never in the worktree.

Every unit of sections 1 to 10 edits spec text only, so no Go file, CRD or generated file changes and `task dev:fmt dev:vet dev:lint dev:test` has nothing to check there; the unit's gate is the gate of `design.md`, and its review is the 14 items of guide section 7. `N` in a gate is the cumulative count of requirements rewritten on the branch.

Each unit starts by writing its scenario ledger into `design.md` (heading `### Ledger: <spec>`), reading `git show main:openspec/specs/<spec>/spec.md` and running `grep -rn --exclude-dir=archive "<old header>" .` before any split.

## 1. Purpose for the six specs with none

- [ ] 1.1 Add a `## Purpose` section of one or two sentences, above `## Requirements`, to `artifact-fetch`, `catalog-registry-resolution`, `failure-counters`, `history-tracking`, `metrics` and `suspend-resume` (6 specs, 0 long requirements; text is read from each spec, never invented from its name). The edit is a delta per spec, ADDED or MODIFIED as `openspec validate` accepts; if the validator cannot take a Purpose-only delta, the unit reports it and asks.
- [ ] 1.2 Gate (`N=0`): no new ERROR item, the ERROR set shrinks by those 6 specs, WARNING lines unchanged.
- [ ] 1.3 Commit `docs(openspec): add a Purpose to six specs`

## 2. opm-operator-B2: 28 requirements

Specs: deps-cascade 10, operator-module-release 10, events-emission 8.

- [ ] 2.1 For each long requirement pick the form (design.md), write the ledger, write the delta under `specs/<spec>/spec.md` of this change. One delta file per spec; no spec of another batch.
- [ ] 2.2 Byte-equality command of guide section 7 item 3 for every rewritten requirement: no `DIFF`, no `LOST`.
- [ ] 2.3 Gate (`N=28`): `validate split-long-requirements --strict` exit 0, archive in the copy exit 0, no "Archive would refuse", no new ERROR item, WARNING lines fell by `N`; `validate <spec> --type spec --strict` shows no WARNING for a rewritten requirement.
- [ ] 2.4 Commit `docs(openspec): split the long requirements of opm-operator-B2`

## 3. opm-operator-B4: 28 requirements

Specs: status-conditions 7, kubernetes-tier-adoption 6, module-instance-ownership 6, platform-reconciler 6, kernel-module-renderer 3.

- [ ] 3.1 For each long requirement pick the form (design.md), write the ledger, write the delta under `specs/<spec>/spec.md` of this change. One delta file per spec; no spec of another batch.
- [ ] 3.2 Byte-equality command of guide section 7 item 3 for every rewritten requirement: no `DIFF`, no `LOST`.
- [ ] 3.3 Gate (`N=56`): `validate split-long-requirements --strict` exit 0, archive in the copy exit 0, no "Archive would refuse", no new ERROR item, WARNING lines fell by `N`; `validate <spec> --type spec --strict` shows no WARNING for a rewritten requirement.
- [ ] 3.4 Commit `docs(openspec): split the long requirements of opm-operator-B4`

## 4. opm-operator-B6: 28 requirements

Specs: platform-inventory-gate 5, render-input-key 5, ssa-apply 5, library-kernel-runtime 4, module-instance-synthesis 4, modulepackage-reconcile-loop 4, digest-computation 1.

- [ ] 4.1 For each long requirement pick the form (design.md), write the ledger, write the delta under `specs/<spec>/spec.md` of this change. One delta file per spec; no spec of another batch.
- [ ] 4.2 Byte-equality command of guide section 7 item 3 for every rewritten requirement: no `DIFF`, no `LOST`.
- [ ] 4.3 Gate (`N=84`): `validate split-long-requirements --strict` exit 0, archive in the copy exit 0, no "Archive would refuse", no new ERROR item, WARNING lines fell by `N`; `validate <spec> --type spec --strict` shows no WARNING for a rewritten requirement.
- [ ] 4.4 Commit `docs(openspec): split the long requirements of opm-operator-B6`

## 5. opm-operator-B7: 28 requirements

Specs: serviceaccount-impersonation 4, platform-crd 3, platform-module-generation 3, reconcile-backoff 2, registration-driven-regeneration 2, registration-removal-guard 2, test-registry-lifecycle 2, workflow-hardening 2, docs-bundle 1, inventory-bridge 1, live-flux-e2e 1, module-acquisition 1, operator-resource-reference 1, platform-gated-rendering 1, registration-activation 1, transformer-registration-crd 1.

- [ ] 5.1 For each long requirement pick the form (design.md), write the ledger, write the delta under `specs/<spec>/spec.md` of this change. One delta file per spec; no spec of another batch.
- [ ] 5.2 Byte-equality command of guide section 7 item 3 for every rewritten requirement: no `DIFF`, no `LOST`.
- [ ] 5.3 Gate (`N=112`): `validate split-long-requirements --strict` exit 0, archive in the copy exit 0, no "Archive would refuse", no new ERROR item, WARNING lines fell by `N`; `validate <spec> --type spec --strict` shows no WARNING for a rewritten requirement.
- [ ] 5.4 Commit `docs(openspec): split the long requirements of opm-operator-B7`

## 6. Pause for the open changes (the human's step)

- [ ] 6.1 The human lands `stop-operator-install-manifest` (with its 8th ERROR fixed by its author) and `add-cross-namespace-source-grants`. A worker never ticks this and never edits those changes. After it, merge `origin/main` into the branch and go on with section 7.

## 7. opm-operator-B1: 28 requirements

Specs: release-automation 13, finalizer-and-deletion 11, instance-health 4.

- [ ] 7.0 Precondition: the human has landed `stop-operator-install-manifest` and `add-cross-namespace-source-grants`. Merge the base into the branch, then re-run the inventory (`validate --all --strict --json`, count `Requirement text is very long` per spec) and compare with this group. Differences go to the report; `N` becomes the fresh cumulative count.
- [ ] 7.1 For each long requirement pick the form (design.md), write the ledger, write the delta under `specs/<spec>/spec.md` of this change. One delta file per spec; no spec of another batch.
- [ ] 7.2 Byte-equality command of guide section 7 item 3 for every rewritten requirement: no `DIFF`, no `LOST`.
- [ ] 7.3 Gate (`N=140`): `validate split-long-requirements --strict` exit 0, archive in the copy exit 0, no "Archive would refuse", no new ERROR item, WARNING lines fell by `N`; `validate <spec> --type spec --strict` shows no WARNING for a rewritten requirement.
- [ ] 7.4 Commit `docs(openspec): split the long requirements of opm-operator-B1`

## 8. opm-operator-B3: 28 requirements

Specs: reconcile-loop-assembly 9, registration-acceptance 9, example-test-modules 7, container-image-publish 3.

- [ ] 8.0 Precondition: the human has landed `stop-operator-install-manifest` and `add-cross-namespace-source-grants`. Merge the base into the branch, then re-run the inventory (`validate --all --strict --json`, count `Requirement text is very long` per spec) and compare with this group. Differences go to the report; `N` becomes the fresh cumulative count.
- [ ] 8.1 For each long requirement pick the form (design.md), write the ledger, write the delta under `specs/<spec>/spec.md` of this change. One delta file per spec; no spec of another batch.
- [ ] 8.2 Byte-equality command of guide section 7 item 3 for every rewritten requirement: no `DIFF`, no `LOST`.
- [ ] 8.3 Gate (`N=168`): `validate split-long-requirements --strict` exit 0, archive in the copy exit 0, no "Archive would refuse", no new ERROR item, WARNING lines fell by `N`; `validate <spec> --type spec --strict` shows no WARNING for a rewritten requirement.
- [ ] 8.4 Commit `docs(openspec): split the long requirements of opm-operator-B3`

## 9. opm-operator-B5: 28 requirements

Specs: prune-stale-resources 6, cascade-receiver 5, drift-detection 5, modulepackage-kernel-rendering 5, operator-module 5, example-module-publishing 2.

- [ ] 9.0 Precondition: the human has landed `stop-operator-install-manifest` and `add-cross-namespace-source-grants`. Merge the base into the branch, then re-run the inventory (`validate --all --strict --json`, count `Requirement text is very long` per spec) and compare with this group. Differences go to the report; `N` becomes the fresh cumulative count.
- [ ] 9.1 For each long requirement pick the form (design.md), write the ledger, write the delta under `specs/<spec>/spec.md` of this change. One delta file per spec; no spec of another batch.
- [ ] 9.2 Byte-equality command of guide section 7 item 3 for every rewritten requirement: no `DIFF`, no `LOST`.
- [ ] 9.3 Gate (`N=196`): `validate split-long-requirements --strict` exit 0, archive in the copy exit 0, no "Archive would refuse", no new ERROR item, WARNING lines fell by `N`; `validate <spec> --type spec --strict` shows no WARNING for a rewritten requirement.
- [ ] 9.4 Commit `docs(openspec): split the long requirements of opm-operator-B5`

## 10. Purpose of the remaining specs

- [ ] 10.1 `example-module-publishing`: replace the archive placeholder Purpose (a WARNING under `--strict`) with real text read from the spec's requirements. The spec is also in section 9; the Purpose text is recorded in its ledger.
- [ ] 10.2 `source-resolution`: add its `## Purpose` (now that `add-cross-namespace-source-grants` is landed, the spec's live text is read first). If that change already wrote a Purpose, do nothing and say so.
- [ ] 10.3 Gate (`N=196` or the fresh total): no new ERROR item; the Purpose ERRORs are gone.
- [ ] 10.4 Commit `docs(openspec): add Purpose text to source-resolution and example-module-publishing`

## 11. Final gate and archive

The pinned OpenSpec 1.14.1 gate (Taskfile and CI) is not part of this change. It moves to a separate change and PR after `stop-operator-install-manifest` and `add-cross-namespace-source-grants` land (DECISIONS.md 2026-10-19 14:05).

- [ ] 11.1 Final gate (`LAST=1`): in a scratch copy, `validate split-long-requirements --strict` exit 0; archive exit 0; `validate --specs --strict` exit 0; no "Archive would refuse"; no new ERROR item. Report `validate --all --strict` as well (expected exit 0 on main after both open changes).
- [ ] 11.2 Archive the change on this branch (`openspec-archive-change`), so the archive commit rides the PR. Commit `docs(openspec): archive split-long-requirements`
