## Why

OpenSpec 1.14.1 raises "Requirement text is very long (>500 characters)" from INFO to WARNING.
`openspec validate --strict` fails on any warning, and no config key or flag turns it off: the
limit is a constant in the validator. The warning is right on the merits: a long requirement
usually holds several behaviors, and one behavior per requirement makes each one testable by its
own scenarios. This change rewrites the long requirements that the archived change
`split-long-requirements` left, so that `openspec validate --all --strict` can pass again with no
loss of rules, and it adds a pinned gate so that the limit holds.

Today `openspec validate --all --strict` under 1.14.1 fails on `origin/main` of this repo. The
main specs hold 84 requirements over 500 characters in 13 specs, and `source-resolution` has no
`## Purpose` (an ERROR, and the last ERROR in the main specs). The two open changes fail too:
`stop-operator-install-manifest` has an ERROR (its MODIFIED block of "The controller's pod and
Service stay those of the kustomize tree" drops a live scenario), and
`add-cross-namespace-source-grants` has two ADDED blocks over 500. The repo has no OpenSpec gate,
so the limit can be crossed again without notice.

The rewrite follows one rule: a rewrite is a split or a move into scenarios. The body itself
gets shorter, a new body is built from whole sentences of the old body, and no rule or scenario
is lost or changed.

## What Changes

No code, no API type, no controller and no CRD changes. Spec text changes, and the repo gains a
validation task and one CI step.

- **Purpose.** `source-resolution` gets a real `## Purpose`, by direct edit of the main spec
  (archive ignores a delta Purpose for an existing spec). The open change cannot archive without
  it.
- **Open-change fixes.** The `specs/` deltas of `add-cross-namespace-source-grants` are edited:
  every block that is, or after archive would become, over 500 is split. The tasks file of that
  change gets one line: at archive, write the real Purpose of the new spec
  `cross-namespace-source-grants`, because the strict gate fails on the placeholder. This is the
  only `tasks.md` edit of an open change. Nothing else of that change moves.
- **Rewrite of the main specs, in about 14 requirements per unit, by spec.** All 84
  requirements in 13 specs. Six units, then one group for the overlap: the first holds `finalizer-and-deletion` and
  `instance-health`; the second `release-automation` and `container-image-publish`; the third
  `reconcile-loop-assembly` and `operator-module`; the fourth `registration-acceptance` and
  `cascade-receiver`; the fifth `example-test-modules` and `prune-stale-resources`; the sixth
  `drift-detection`, `modulepackage-kernel-rendering` and `example-module-publishing`. The four
  requirements that `stop-operator-install-manifest` MODIFIES or REMOVES are split too, but in
  their own group, together with the retarget of that change (next bullet), so the main spec and
  the delta move in one commit.
- **Retarget of `stop-operator-install-manifest`.** Its delta is edited in this change. A MODIFIED
  block of a split requirement becomes one MODIFIED block per new header, its text cut at the same
  boundary as the main split, verbatim, each part 500 or less. A part with no counterpart in the
  delta text keeps its main text. A REMOVED block of a split requirement becomes one REMOVED block
  per new header. The dropped scenario of its `operator-module` block is restored. Its proposal,
  design, tasks and code stay as they are.
- **The gate.** A Taskfile task installs OpenSpec 1.14.1 with `--ignore-scripts` and a second
  task runs `openspec validate --all --strict --no-interactive`. A job of its own in `lint.yml`, with a read-only token and a
  checkout that keeps no credentials, runs both. It is the last group, in the same PR, and it is green on the branch.

## Impact

- Affected API types and controllers: none.
- SemVer class: none (specification text, a Taskfile task and a CI step). The commits are `docs`
  and `ci`, which release nothing, so no `-beta.N` is cut.
- Complexity: no new component. One pinned npm CLI in CI, installed without install scripts, in a
  job whose token only reads.
- Sections: more than about five (see `tasks.md`). Principle VIII allows one PR per change with
  one commit per section. Each spec section ends green under the gate of `design.md` and leaves
  `main` releasable. The change runs as several units on one branch.
- Risk: a rewrite that drops or changes a rule while the validator stays green. The scenario
  ledger and the byte-equality check are the counter. A second risk is the retarget: a MODIFIED or
  REMOVED block that names a header that the archive of this change did not create makes the
  archive of `stop-operator-install-manifest` refuse. The archive simulation in a scratch copy
  (this change, then `stop-operator-install-manifest`, then `add-cross-namespace-source-grants`)
  is the counter. The change is one PR (see `design.md`, "Delivery order").

## Enhancement

None.
