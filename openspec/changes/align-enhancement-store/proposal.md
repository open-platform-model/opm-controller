## Why

The enhancements now live in an OpenSpec store (`opm-enhancements`). The operator still cites them in the retired form: a decision number, a lone id, or an `enhancement.yaml` file per change. A decision number is not a stable address and a lone id gives the reader no reason. The enhancements design (rule: an id only beside a reason in words, never alone, never a decision number) replaces this form in every repo.

## What Changes

- Every enhancement citation in comments, live spec text, archived changes, docs and test messages becomes an id beside the reason in words. A tool does the mechanical part (`enh cite-rewrite`); the rows it cannot rewrite are rewritten by hand.
- Citations of a retired id lose the id and keep the reason text.
- The 21 archived `enhancement.yaml` files are deleted. `enhancements/000-TEMPLATE/` (a template of the retired layout) is deleted.
- `openspec/config.yaml` gains `references:` to the store, the `Delivers` rule, a `context` line to read an entry, and the new comment rule. The `enhancement.yaml` and `task enhancements:delivery:*` rules go. `AGENTS.md` says the same.
- Generated copies (`config/crd/bases/*`, `dist/install.yaml`, `modules/opm_operator/zz_generated_crds.cue`) are regenerated from the changed Go type comments, never edited by hand.
- `docs-kit.cue` drops `citations: "link"`; `openspec/specs/docs-bundle/spec.md` is updated in wording.
- Excluded on purpose: published fixtures under `test/fixtures/modules/` (an edit forces a fixture version bump) and the active changes `add-cross-namespace-source-grants` and `stop-operator-install-manifest` (they belong to their authors).

No behaviour changes. The only code touched is comments, string text and test messages; one production string (`internal/controller/platform_inventory.go`) reaches the status message and loses its decision number.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

None. Live specs (`viewer-roles`, `docs-bundle` and others) change in wording only: no requirement is added, removed or changed, so there is no spec delta. The change sets `skip_specs: true` in `.openspec.yaml`.

## Impact

- SemVer: none (documentation and comments only; commits use `docs`, `test`, `chore`).
- Generated CRD descriptions change in text; the CRD schema does not.
- Owner action after merge: `.claude/skills/openspec-new-change` and `openspec-propose` carry a `REPO-LOCAL PATCH` ("Declare the enhancement now", `enhancement.yaml`). The owner removes both by running `openspec update`, because agents do not write under `.claude/`. `config.yaml` now carries what the patch said.
