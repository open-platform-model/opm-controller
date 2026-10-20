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
- **Open-change fixes.** The `specs/` deltas of the two open changes are edited, and nothing else
  of them: the dropped scenario of `stop-operator-install-manifest` is restored, and every
  block of the two deltas that is, or after archive would become, over 500 is split. The proposal
  intent, the tasks and the code of those changes stay as they are.
- **Rewrite of the main specs, in six units of about 14 requirements, by spec.** 80 requirements
  in 13 specs: the first unit holds `finalizer-and-deletion` and `instance-health`; the second
  `release-automation` and `container-image-publish`; the third `reconcile-loop-assembly` and
  `operator-module`; the fourth `registration-acceptance` and `cascade-receiver`; the fifth
  `example-test-modules` and `prune-stale-resources`; the sixth `drift-detection`,
  `modulepackage-kernel-rendering` and `example-module-publishing`.
- **Four long main requirements stay as they are.** The open change
  `stop-operator-install-manifest` MODIFIES three of them and REMOVES one, so a rewrite in the
  main spec would break its archive. The fix of the three is in its delta (`design.md`, "the
  overlap list"); the removed one needs no fix. They leave or change in the main specs when that
  change archives.
- **The gate.** A Taskfile task installs OpenSpec 1.14.1 with `--ignore-scripts` and a second
  task runs `openspec validate --all --strict --no-interactive`. A step in the read-only `Lint`
  workflow job runs both.

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
  ledger and the byte-equality check are the counter. A second risk is delivery order: the CI step
  fails on `main` until both open changes have archived (see `design.md`, "Delivery order").

## Enhancement

None.
