## Context

The rewrite rule is design T2.3 sections 6 and 9 of the refactoring-enhancements swarm. The tool `enh cite-rewrite` (built at `fdb2ec9e` on the `enhancements` branch `p5-rewrite`) rewrites by form and writes a worklist for rows it cannot rewrite safely. `enh cite-lint` is the closing check. Neither is wired into this repo yet; that is a later task (a local task and a CI call, pinned to merged commits).

## Decisions

- **Tool first, hand second.** Run `cite-rewrite --write` once with a worklist, commit the result per place (code, live specs, archived changes, docs and rest). Then clear every worklist row by hand, or leave it with a reason (excluded path, generated file).
- **Generated files are regenerated.** The Go type comments in `api/v1alpha1` change; `task dev:manifests dev:generate`, `task operator:installer` regenerate the copies, and `task operator-module:drift` proves the CUE copy matches.
- **`0030:`** loses the id and keeps the reason text, in the live spec, the RBAC kustomization comment and the test.
- **Excluded paths** stay as they are and are named in the PR body.
- **Rules in `config.yaml`, not in skills.** `rules.proposal` carries the `Delivers` line; `context` carries the line to read an entry. `references:` uses `- id: opm-enhancements` with the remote of the store.
- **No new behaviour, no new check in this change.** The lint task and CI call come later.

## Risks

- A wrong tool edit in a comment reads oddly but changes nothing; the sample read and the review catch it.
- A regenerated CRD that differs in schema (not only text) would be a defect; `operator-module:drift` and a diff of the CRD without `description` fields guard it.
