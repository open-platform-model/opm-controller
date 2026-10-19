# Tasks: split-long-requirements

Branch `docs/split-long-requirements`, created from `origin/main`. Run every command inside its worktree. The gate of `design.md` runs in a scratch copy of `openspec/`, never in the worktree.

Every unit of sections 1 to 10 edits spec text only, so no Go file, CRD or generated file changes and `task dev:fmt dev:vet dev:lint dev:test` has nothing to check there; the unit's gate is the gate of `design.md`, and its review is the 14 items of guide section 7. `N` in a gate is the cumulative count of requirements rewritten on the branch.

Each unit starts by writing its scenario ledger into `design.md` (heading `### Ledger: <spec>`), reading `git show main:openspec/specs/<spec>/spec.md` and running `grep -rn --exclude-dir=archive "<old header>" .` before any split.

## 1. Purpose for the six specs with none

- [x] 1.1 Add a `## Purpose` section of one or two sentences, above `## Requirements`, to `artifact-fetch`, `catalog-registry-resolution`, `failure-counters`, `history-tracking`, `metrics` and `suspend-resume` (6 specs, 0 long requirements; text is read from each spec, never invented from its name). The edit is a direct edit of `## Purpose` in each main spec, because archive ignores a delta Purpose for an existing spec; `design.md` records it.
- [x] 1.2 Gate (`N=0`): no new ERROR item, the ERROR set shrinks by those 6 specs, WARNING lines unchanged.
- [x] 1.3 Commit `docs(openspec): add a Purpose to six specs`

## 2. opm-operator-B2: 28 requirements

Specs: deps-cascade 10, operator-module-release 10, events-emission 8.

- [x] 2.1 For each long requirement pick the form (design.md), write the ledger, write the delta under `specs/<spec>/spec.md` of this change. One delta file per spec; no spec of another batch.
- [x] 2.2 Byte-equality command of guide section 7 item 3 for every rewritten requirement: no `DIFF`, no `LOST`.
- [x] 2.3 Gate (`N=28`): `validate split-long-requirements --strict` exit 0, archive in the copy exit 0, no "Archive would refuse", no new ERROR item, WARNING lines fell by `N`; `validate <spec> --type spec --strict` shows no WARNING for a rewritten requirement.
- [x] 2.4 Commit `docs(openspec): split the long requirements of opm-operator-B2`

## 3. opm-operator-B4: 28 requirements

Specs: status-conditions 7, kubernetes-tier-adoption 6, module-instance-ownership 6, platform-reconciler 6, kernel-module-renderer 3.

- [x] 3.1 For each long requirement pick the form (design.md), write the ledger, write the delta under `specs/<spec>/spec.md` of this change. One delta file per spec; no spec of another batch.
- [x] 3.2 Byte-equality command of guide section 7 item 3 for every rewritten requirement: no `DIFF`, no `LOST`.
- [x] 3.3 Gate (`N=56`): `validate split-long-requirements --strict` exit 0, archive in the copy exit 0, no "Archive would refuse", no new ERROR item, WARNING lines fell by `N`; `validate <spec> --type spec --strict` shows no WARNING for a rewritten requirement.
- [x] 3.4 Commit `docs(openspec): split the long requirements of opm-operator-B4`

## 4. opm-operator-B6: 28 requirements

Specs: platform-inventory-gate 5, render-input-key 5, ssa-apply 5, library-kernel-runtime 4, module-instance-synthesis 4, modulepackage-reconcile-loop 4, digest-computation 1.

- [x] 4.1 For each long requirement pick the form (design.md), write the ledger, write the delta under `specs/<spec>/spec.md` of this change. One delta file per spec; no spec of another batch.
- [x] 4.2 Byte-equality command of guide section 7 item 3 for every rewritten requirement: no `DIFF`, no `LOST`. The command prints one false `DIFF` for "Content sensitivity" of `digest-computation`, because the old spec holds two scenarios of that name; a per-requirement byte comparison of all 144 kept scenarios shows none changed.
- [x] 4.3 Gate (`N=84`): `validate split-long-requirements --strict` exit 0, archive in the copy exit 0, no "Archive would refuse", no new ERROR item, WARNING lines fell by `N`; `validate <spec> --type spec --strict` shows no WARNING for a rewritten requirement.
- [x] 4.4 Commit `docs(openspec): split the long requirements of opm-operator-B6`

## 5. opm-operator-B7: 28 requirements

Specs: serviceaccount-impersonation 4, platform-crd 3, platform-module-generation 3, reconcile-backoff 2, registration-driven-regeneration 2, registration-removal-guard 2, test-registry-lifecycle 2, workflow-hardening 2, docs-bundle 1, inventory-bridge 1, live-flux-e2e 1, module-acquisition 1, operator-resource-reference 1, platform-gated-rendering 1, registration-activation 1, transformer-registration-crd 1.

- [x] 5.1 For each long requirement pick the form (design.md), write the ledger, write the delta under `specs/<spec>/spec.md` of this change. One delta file per spec; no spec of another batch. (Unit 5a: `serviceaccount-impersonation`, `platform-crd`, `platform-module-generation`, `reconcile-backoff`, `registration-driven-regeneration`, 14 requirements. Unit 5b: `registration-removal-guard`, `test-registry-lifecycle`, `workflow-hardening`, `docs-bundle`, `inventory-bridge`, `live-flux-e2e`, `module-acquisition`, `operator-resource-reference`, `platform-gated-rendering`, `registration-activation`, `transformer-registration-crd`, 14 requirements.)
- [x] 5.2 Byte-equality command of guide section 7 item 3 for every rewritten requirement: no `DIFF`, no `LOST`. (Units 5a and 5b: the 28 requirements above; no `DIFF`, no `LOST`.)
- [x] 5.3 Gate (`N=112`): `validate split-long-requirements --strict` exit 0, archive in the copy exit 0, no "Archive would refuse", no new ERROR item, WARNING lines fell by `N`; `validate <spec> --type spec --strict` shows no WARNING for a rewritten requirement. (Unit 5a: gate run with `N=98`. Unit 5b: gate run with `N=112`, `GATE PASS`; the 11 specs of unit 5b show no WARNING.)
- [x] 5.4 Commit `docs(openspec): split the long requirements of opm-operator-B7`

## 10. Purpose of example-module-publishing

- [x] 10.1 `example-module-publishing`: replace the archive placeholder Purpose (a WARNING under `--strict`) with real text read from the spec's requirements, by direct edit of `## Purpose`. The spec's two long requirements are not rewritten in this change (see "Not in this change" in `proposal.md`).

## 11. Final gate and archive

The pinned OpenSpec 1.14.1 gate (Taskfile and CI) is not part of this change. It moves to a separate change and PR.

- [x] 11.0 Update the seven workflow comments that cite requirement names this change removed (comment text only), in its own `ci` commit. Remove the decision times from `design.md`.
- [x] 11.1 Final gate (`N=113`: 112 requirements and the placeholder Purpose warning of `example-module-publishing`, `LAST` unset, base `origin/main`): in a scratch copy, `validate split-long-requirements --strict` exit 0; archive exit 0; WARNING lines fall by 113; no "Archive would refuse"; no new ERROR item. Report `validate --all --strict` as well (it still fails on the 84 long requirements of B1, B3 and B5 and on the open changes).
- [x] 11.2 Archive the change on this branch (`openspec-archive-change`), so the archive commit rides the PR. Commit `docs(openspec): archive split-long-requirements`

## Not in this change

Groups B1, B3 and B5 (84 requirements) and the `source-resolution` Purpose stay out. They wait for the open changes `stop-operator-install-manifest` and `add-cross-namespace-source-grants`. The pinned CI gate is also out. See `proposal.md`.
