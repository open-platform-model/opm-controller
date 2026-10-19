## Why

OpenSpec 1.14.1 raises "Requirement text is very long (>500 characters)" from INFO to WARNING.
`openspec validate --strict` fails on any warning, and no config key or flag turns it off: the
limit is a constant in the validator. The warning is right on the merits: a long requirement
usually holds several behaviors, and one behavior per requirement makes each one testable by its
own scenarios. This change rewrites the long requirements of this repo so that `--strict` can
pass again with no loss of rules.

Today `openspec validate --all --strict` under 1.14.1 reports 196 requirements over 500
characters in 44 main specs (this change rewrites 112 of them), and 8 ERRORs. Seven ERRORs are specs without a `## Purpose`. A
placeholder Purpose ("TBD - created by archiving change ...") is a second warning kind that also
fails `--strict`. The eighth ERROR sits in the open change `stop-operator-install-manifest`; its
author fixes it (see Scope).

The rewrite follows the swarm guide `design/T2.1.md`. Its rule is
that a rewrite is a split or a move into scenarios: the body itself gets shorter, and no rule is
lost or changed.

## What Changes

No code, no API type, no controller and no CRD changes. Only `openspec/specs/` text changes. Spec deltas are not written in this step: later implement units write
them batch by batch, each under the human's yes at this proposal gate.

Specs by batch (counts are long requirements, from `explore/T1.3.md` and re-measured on
`origin/main` 7713cc86):

- **Purpose, first.** Add `## Purpose` to `artifact-fetch`, `catalog-registry-resolution`,
  `failure-counters`, `history-tracking`, `metrics` and `suspend-resume`. They hold no long
  requirement.
- **opm-operator-B2 (28):** `deps-cascade`, `operator-module-release`, `events-emission`.
- **opm-operator-B4 (28):** `status-conditions`, `kubernetes-tier-adoption`,
  `module-instance-ownership`, `platform-reconciler`, `kernel-module-renderer`.
- **opm-operator-B6 (28):** `platform-inventory-gate`, `render-input-key`, `ssa-apply`,
  `library-kernel-runtime`, `module-instance-synthesis`, `modulepackage-reconcile-loop`,
  `digest-computation`.
- **opm-operator-B7 (28):** `serviceaccount-impersonation`, `platform-crd`,
  `platform-module-generation`, `reconcile-backoff`, `registration-driven-regeneration`,
  `registration-removal-guard`, `test-registry-lifecycle`, `workflow-hardening`, `docs-bundle`,
  `inventory-bridge`, `live-flux-e2e`, `module-acquisition`, `operator-resource-reference`,
  `platform-gated-rendering`, `registration-activation`, `transformer-registration-crd`.
- **Purpose of `example-module-publishing`.** Real Purpose text replaces its archive placeholder, by direct edit.
- **Workflow comments.** Seven workflow comments that cite removed requirement names are updated (comment text only).

## Not in this change

The human cut the scope: the two open changes `stop-operator-install-manifest` and
`add-cross-namespace-source-grants` are not waited for. This change rewrites 112 of the 196 long
requirements. The other 84 stay long and stay for a later change:

- B1 (`release-automation`, `finalizer-and-deletion`, `instance-health`), B3
  (`reconcile-loop-assembly`, `registration-acceptance`, `example-test-modules`,
  `container-image-publish`) and B5 (`prune-stale-resources`, `cascade-receiver`,
  `drift-detection`, `modulepackage-kernel-rendering`, `operator-module`,
  `example-module-publishing`). These specs are touched or neighbored by the open changes, whose
  MODIFIED blocks a rewrite would break.
- The `## Purpose` of `source-resolution`, which the open change `add-cross-namespace-source-grants` owns.
- The pinned OpenSpec 1.14.1 `openspec validate --all --strict` gate (Taskfile and CI). It becomes a separate change.

Until those land, `openspec validate --all --strict` still fails on the 84 requirements and on the open changes.

## Scope

In scope: the specs named above, the Purpose text.

The 8th ERROR (the MODIFIED block of `operator-module` in `stop-operator-install-manifest`
omits a live scenario) is that change's to fix before it lands.
This change does not edit that change or `add-cross-namespace-source-grants`. Their ADDED blocks
over 500 characters belong to their authors.

## Impact

- Affected API types and controllers: none.
- SemVer class: none (specification text only). The commits are `docs`, which
  release nothing, so no `-beta.N` is cut.
- Complexity: no new component.
- Sections: more than about five (see `tasks.md`). Principle VIII allows one PR per change with
  one commit per section; a section that edits specs only ends green under the guide gate
  and leaves `main` releasable. The change runs as several units on one branch.
- Risk: a rewrite that drops or changes a rule while the validator stays green. The scenario
  ledger and the byte-equality check of the guide (section 7, item 3) are the counter.

## Enhancement

None.
