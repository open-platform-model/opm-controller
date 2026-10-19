## Why

OpenSpec 1.14.1 raises "Requirement text is very long (>500 characters)" from INFO to WARNING.
`openspec validate --strict` fails on any warning, and no config key or flag turns it off: the
limit is a constant in the validator. The warning is right on the merits: a long requirement
usually holds several behaviors, and one behavior per requirement makes each one testable by its
own scenarios. This change rewrites the long requirements of this repo so that `--strict` can
pass again with no loss of rules.

Today `openspec validate --all --strict` under 1.14.1 reports 196 requirements over 500
characters in 44 main specs, and 8 ERRORs. Seven ERRORs are specs without a `## Purpose`. A
placeholder Purpose ("TBD - created by archiving change ...") is a second warning kind that also
fails `--strict`. The eighth ERROR sits in the open change `stop-operator-install-manifest`; its
author fixes it (see Scope).

The rewrite follows the swarm guide `design/T2.1.md` (approved 2026-10-19 13:40). Its rule is
that a rewrite is a split or a move into scenarios: the body itself gets shorter, and no rule is
lost or changed.

## What Changes

No code, no API type, no controller and no CRD changes. Only `openspec/specs/` text, a Taskfile
task and a CI step change. Spec deltas are not written in this step: later implement units write
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
- **opm-operator-B1 (28), B3 (28), B5 (28) run last.** They hold specs that the open changes
  `stop-operator-install-manifest` and `add-cross-namespace-source-grants` touch or sit near
  (B1: `release-automation`, `finalizer-and-deletion`, `instance-health`; B3:
  `reconcile-loop-assembly`, `registration-acceptance`, `example-test-modules`,
  `container-image-publish`; B5: `prune-stale-resources`, `cascade-receiver`, `drift-detection`,
  `modulepackage-kernel-rendering`, `operator-module`, `example-module-publishing`). They start
  after the human lands both open changes, and each starts with a fresh inventory (guide
  section 8), because the open changes alter lengths and may add long requirements.
- **Purpose, after the open changes land.** `source-resolution` gets its Purpose in the last
  group, after `add-cross-namespace-source-grants` lands. `example-module-publishing` gets real
  Purpose text (an archive placeholder today) in the B5 unit, which already rewrites that spec.
- **Gate.** A pinned OpenSpec 1.14.1 `openspec validate --all --strict` gate in `Taskfile.yml`
  (via `.tasks/`) and in CI, as the last unit, delivered after both open changes land.

## Scope

In scope: the specs named above, the Purpose text, the gate.

The 8th ERROR (the MODIFIED block of `operator-module` in `stop-operator-install-manifest`
omits a live scenario) is that change's to fix before it lands (DECISIONS.md 2026-10-19 12:45).
This change does not edit that change or `add-cross-namespace-source-grants`. Their ADDED blocks
over 500 characters belong to their authors.

## Impact

- Affected API types and controllers: none.
- SemVer class: none (specification text and CI only). The commits are `docs` and `ci`, which
  release nothing, so no `-beta.N` is cut.
- Complexity: no new component. The one added moving part is a CI step; it keeps the repo
  clean once the rewrite is done (Principle VII).
- Sections: more than about five (see `tasks.md`). Principle VIII allows one PR per change with
  one commit per section; a section that edits specs only ends green under the guide gate
  and leaves `main` releasable. The change runs as several units on one branch (DECISIONS.md
  2026-10-19 09:50). The delivery order is a question for the gate (see `design.md`).
- Risk: a rewrite that drops or changes a rule while the validator stays green. The scenario
  ledger and the byte-equality check of the guide (section 7, item 3) are the counter.

## Enhancement

None.
