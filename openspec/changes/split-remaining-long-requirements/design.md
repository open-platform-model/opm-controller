## Context

The validator measures the body of a requirement: the lines between `### Requirement:` and the
first `#` line, with blank lines, fenced lines and `**Key**:` lines dropped (the last only when
other text exists), joined by newlines. Scenarios do not count. More than 500 UTF-16 units is a
WARNING; `--strict` fails on it. The check runs on main specs and on ADDED blocks of a delta, not
on MODIFIED blocks. A MODIFIED block that stays over 500 therefore passes
`validate <change> --strict`, archives, and then fails `validate --all --strict` on the spec.

Current state, measured on `origin/main` 693b6451 with OpenSpec 1.14.1 in a scratch copy
(`validate --specs --strict --json`, 2026-10-20): 59 main specs, 45 pass, 14 fail. 84 requirements
are over 500 characters in 13 specs; the 14th failing spec is `source-resolution`, one ERROR (no
`## Purpose`) and no long requirement. The count per spec is in "The units". The archived change
`split-long-requirements` rewrote the other 112 long requirements; its conventions (forms,
ledger, strength table) carry over, and this file does not restate its text.

Two open changes exist: `add-cross-namespace-source-grants` (specs `cross-namespace-source-grants`
and `source-resolution`) and `stop-operator-install-manifest` (specs `container-image-publish`,
`operator-module`, `release-automation`).

The rules of the rewrite are: copy whole sentences of the old body, state one behavior in each
requirement, keep every scenario byte for byte, and add no word beyond a subject or condition
copied from the old requirement's own text and the RFC 2119 keyword where a sentence has none.
Every sentence that gains a keyword goes into the strength table below.

## Goals / Non-Goals

Goals:

- `openspec validate --all --strict --no-interactive` under 1.14.1 exits 0 on the branch, and
  exits 0 after the two open changes are archived in the order: this change, then
  `stop-operator-install-manifest`, then `add-cross-namespace-source-grants` (with the Purpose of
  the new spec written, see "Risks").
- Every rewrite is a split or a move into scenarios. No rule, value, default, exit code,
  citation or scenario is lost or changed.
- The gate runs in `task` and in CI, pinned.

Non-goals: any code change; any change of meaning of a requirement or scenario, in a main spec or
in an open change; the proposal intent, the design or the code of an open change, and its tasks
except the one Purpose line named below; the Purpose that archive writes into the new spec
`cross-namespace-source-grants` (the archiver writes it, see "Risks").

## Research & Decisions

### Which form per requirement

**Context**: A long requirement can shrink four ways (forms a, b, c2, d of the approved guide).
**Decision**: choose per requirement by reading the text, never from an inventory label:

| The excess is | Form |
| --- | --- |
| Examples, lists, formats, rationale around one behavior | a: MODIFIED, shorter body, more scenarios |
| Several behaviors, and the old scenarios fit the parent that stays | b: MODIFIED parent plus ADDED splits |
| Several behaviors, and the old scenarios belong to different behaviors | c2: REMOVED plus ADDED under new names; the ledger is mandatory |
| A name that must change, with no change of rule | d: RENAMED plus MODIFIED under the new name, alone |

**Rationale**: b and c2 give one testable behavior per requirement, which is the reason for the
warning. Form a alone on a requirement with several behaviors would only hide them in scenarios.
Not allowed: text moved into `**Key**:` lines, under a `#####` heading or into a fenced block;
a live scenario moved from a MODIFIED parent into an ADDED block; two operations on one name.

### What this repo adds to the guide

**Context**: the `specs` rules of `openspec/config.yaml`.
**Decision**: RFC 2119 keywords stay in every body with the same strength. A citation (an
enhancement id beside a reason in words, never alone, never a decision number) stays in the body of
the block that carries the rule it supports, never in the header. A scenario never adds a case
that is not in the old text; a new case is a question for the human. Before a split or a c2,
`grep -rn --exclude-dir=archive "<old header>" .` runs and stale references are fixed in the same
unit. An em-dash of an old body is written as `:` or `,` in a new body (workspace style).

### The open changes: the overlap list

**Context**: this change splits all 84 long requirements on `main`, so a split requirement that an
open change MODIFIES or REMOVES no longer has the header that the open change names, and the
archive of that change would refuse. The list below is fresh: it compares the 84 long main
requirements with every requirement the two deltas name, and measures the delta blocks
(`validate <change> --strict --json`, and a length count of every delta block, MODIFIED included,
because validate skips those).
**Decision** (the human's gate, 2026-10-20 15:40): the four long main requirements in the list are
split on `main` in this change, by the same verbatim rule, and the delta of
`stop-operator-install-manifest` is retargeted in the same change. The rejected alternative of the
first draft ("leave the four long on main and fix them in the delta") is replaced by this.

| Open change | Requirement (header) | Delta op | Main length | Delta length | What this change does |
| --- | --- | --- | --- | --- | --- |
| stop-operator-install-manifest | The controller's pod and Service stay those of the kustomize tree (operator-module) | MODIFIED | 581 | 670 | Split on main; retarget the MODIFIED block to one MODIFIED block per new header; restore the dropped scenario |
| stop-operator-install-manifest | Release published once after every release job (release-automation) | MODIFIED | 837 | 735 | Split on main; retarget |
| stop-operator-install-manifest | Beta prerelease line (release-automation) | MODIFIED | 1475 | 1475 | Split on main; retarget |
| stop-operator-install-manifest | Release install manifest with digest-pinned image (container-image-publish) | REMOVED | 685 | none | Split on main; the REMOVED block becomes one REMOVED block per new header |
| stop-operator-install-manifest | Release assets upload only to a draft release (release-automation) | MODIFIED | 488 | 437 | No change: 500 or less, header kept. Listed so the list is complete |
| add-cross-namespace-source-grants | OCIRepository lookup (source-resolution) | MODIFIED | 280 | 633 | Not long on main; the delta text would be long after archive: split the delta (form a or b) |
| add-cross-namespace-source-grants | SourceRefGrant CRD (cross-namespace-source-grants) | ADDED | none | 571 | Split the delta block by the same rule |
| add-cross-namespace-source-grants | Grant-matching policy permits a reference (cross-namespace-source-grants) | ADDED | none | 547 | Split the delta block by the same rule |

The first four rows are the long main requirements that an open change touches. The last three
rows do not appear in the 84. The brief of the human named one long ADDED block of
`add-cross-namespace-source-grants`; the measure shows two (`validate` lists both), and the MODIFIED
`OCIRepository lookup` would fail after archive, so three blocks of that delta are split.

Rules to retarget a delta of `stop-operator-install-manifest`, besides the verbatim rule:

- **MODIFIED of a split requirement.** The main split gives the old header a set of new headers
  (form b: the parent plus ADDED splits; form c2: new headers only). The delta gets one MODIFIED
  block per new header. Its text is the delta text cut at the same kind of boundary as the main
  split (whole sentences, verbatim), each part 500 or less. Each live scenario of the delta block
  goes under the new header that holds that scenario on main (the ledger of the main split names
  it), byte for byte.
- **A part with no counterpart.** When a new header on main has no sentence in the delta text,
  the delta did not change those sentences. Then either the block is left out (the main text
  stays: say so in the ledger, with the reason "the delta did not change these sentences") or it is
  copied verbatim from the new main text as a MODIFIED block (only when a scenario of the delta
  belongs under that header). Each case is written in the ledger.
- **REMOVED of a split requirement.** One REMOVED block per new header, with the header (and the
  Reason and Migration lines the delta has) as the old block had. All of them together remove
  what the old block removed. A new header that the old removal did not cover by its text is a
  question for the human.
- **ERROR repair.** The dropped scenario "An object added to the kustomize tree only fails the
  test" of the `operator-module` block is restored byte equal to `origin/main`, under the new
  header that holds it.
- No name in the delta is in two operations. A scenario that the open change wrote is kept byte
  for byte. The meaning of the delta text never changes: a part that the cut would change is a
  question for the human.
- Each edit is checked by `openspec validate stop-operator-install-manifest --strict` (exit 0),
  the byte-equality command (against `origin/main` for a MODIFIED part, against the delta on
  `origin/main` for text the delta wrote), and the archive simulation of "The gate".

Rules for `add-cross-namespace-source-grants` (deltas only, besides one tasks line):

- An ADDED block that was long is split into ADDED blocks of the same delta (form c2 inside one
  delta: the long header disappears, the new headers appear). The MODIFIED `OCIRepository lookup`
  keeps its header and every live scenario byte for byte; the split-off behaviors become ADDED
  requirements under new names with their own scenarios. A live scenario never moves from a
  MODIFIED parent to an ADDED block. A new name is checked against the names this change adds to
  the same spec.
- `tasks.md` of that change gets one line, at the end of its last section: at archive, write the
  real Purpose of the new spec `cross-namespace-source-grants`, because the strict gate fails on
  the placeholder. No other task is added, changed or ticked.

### The Purpose of source-resolution

**Context**: `openspec archive add-cross-namespace-source-grants` aborts today: "Validation
errors in rebuilt spec for source-resolution: Spec must have a Purpose section". A Purpose must
exist before that archive can run. Archive ignores a delta Purpose for an existing spec.
**Decision**: add `## Purpose` above `## Requirements` of `openspec/specs/source-resolution/spec.md`
by direct edit, in its own commit, with text read from the spec's requirements. The text is at
least 50 characters (a shorter one is a WARNING, "Purpose section is too brief"; measured in a
scratch copy).
**Rationale**: it is the only way to clear the last main-spec ERROR and to let the archive of the
open change run.

### Order of the units

**Context**: the human asked for units of about 14 requirements, by spec, with the open-change
fixes and the Purpose as their own groups and the gate last.
**Decision**: Purpose, fixes of `add-cross-namespace-source-grants`, six rewrite units, one group for
the four overlap requirements with the retarget of `stop-operator-install-manifest`, gate. Units
are about 14 requirements. The four overlap requirements (operator-module 1, release-automation 2,
container-image-publish 1) are planned in their own group, not in the units of their specs (2 and
3), because their split and the retarget must land in one commit; the units of their specs then
hold 13 requirements each, and the overlap group holds 4.

| Group | Content | Requirements | N after it |
| --- | --- | --- | --- |
| 1 Purpose | `source-resolution` Purpose, direct edit | 0 | 0 |
| 2 Open-change fixes | Deltas of `add-cross-namespace-source-grants`; the Purpose line in its `tasks.md` | the 2 long ADDED blocks | 2 |
| 3 Unit 1 | finalizer-and-deletion 11, instance-health 4 | 15 | 17 |
| 4 Unit 2 | release-automation 11, container-image-publish 2 | 13 | 30 |
| 5 Unit 3 | reconcile-loop-assembly 9, operator-module 4 | 13 | 43 |
| 6 Unit 4 | registration-acceptance 9, cascade-receiver 5 | 14 | 57 |
| 7 Unit 5 | example-test-modules 7, prune-stale-resources 6 | 13 | 70 |
| 8 Unit 6 | drift-detection 5, modulepackage-kernel-rendering 5, example-module-publishing 2 | 12 | 82 |
| 9 Overlap | operator-module 1, release-automation 2, container-image-publish 1, with the retarget of `stop-operator-install-manifest` | 4 | 86 |
| 10 Gate | Taskfile tasks, CI step, final archive simulation, archive of this change | 0 | 86 |

`N` is the number of WARNING lines the whole change on the branch has removed from
`validate --all --strict`, counted against `origin/main` (86 WARNING lines: the 84 long main
requirements and the 2 long ADDED blocks of the open change). Group 2 removes 2 lines. After
group 9, no line remains. The count per spec is from `openspec validate --specs --strict --json`
on `origin/main` (13 release-automation, 11 finalizer-and-deletion, 9 reconcile-loop-assembly, 9
registration-acceptance, 7 example-test-modules, 6 prune-stale-resources, 5 each operator-module,
modulepackage-kernel-rendering, drift-detection and cascade-receiver, 4 instance-health, 3
container-image-publish, 2 example-module-publishing). Units 2 and 3 leave the overlap
requirement of their specs to group 9.
**Alternative rejected**: leave the four overlap requirements long on `main` and fix the open
change's delta text only (the first draft). It leaves `validate --all --strict` red on the branch
and forces a second PR. The human rejected it at the gate (2026-10-20 15:40).

### Reversibility

high (hours). Spec text is reverted by a PR; the gate is one task and one step. No contract
changes and no accepted decision record is edited.

## The scenario ledger

Each unit writes its ledger in this file before it writes the delta of a spec, under the heading
`### Ledger: <spec>`, from the requirement text read when the unit starts. The open-change
fixes write theirs under `### Ledger: <open change>, <spec>`. The first table has the form of the
guide:

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |

A second table lists the rule clauses of the old body (one per SHALL, MUST, number, list item,
format) and the clause or scenario THEN line that carries each. Rules of the ledger:

- Forms a, b, d: every old scenario name stays under the same requirement, byte for byte.
- Form c2: every old scenario name appears in exactly one new requirement, WHEN and THEN byte
  for byte. OpenSpec does not check this: an old scenario in no new requirement disappears with
  no error. The reviewer checks it by name.
- A clause that moves from a body SHALL into a scenario THEN keeps its SHALL there, and only a
  default or an exact value may move this way. A clause that adds or changes a condition stays in
  the body. A separate behavior goes to b or c2.

No ledger is written yet. This section fills in as the units run.

## Strength table

Every sentence that gains SHALL, MUST or MAY, collected from the clause tables of all units, with
the old sentence. Empty until the units run.

| Requirement (new) | Sentence | Old sentence, no keyword |
| --- | --- | --- |

## The gate

### Per-unit gate (guide section 6)

Run in a scratch copy of the worktree's `openspec/` folder, never in the worktree, with the
pinned CLI (1.14.1) and `jq`. Values:

```bash
WT=<worktree>
CH=split-remaining-long-requirements
S=<scratch folder under the session scratchpad>
N=<WARNING lines removed so far, counted against origin/main, earlier groups included>
# LAST is unset in every unit; the final check is the next block
```

The block of the guide runs unchanged, with `origin/main` in place of `main` in its
`git archive` line (the local `main` of the repo can be behind). It validates the change,
archives it in the copy, runs `validate --all --strict`, fails on an "Archive would refuse" line
and on a new ERROR item, and requires WARNING lines to fall by `N` (baseline 86). The copy
contains both open changes, as edited by groups 2 and 9 once they are done. Pass rules:

- `validate <change> --strict` alone is not enough: it skips long MODIFIED text.
- ERRORs are compared as sets of `type id path`; the set may shrink and never gains an item.
  `origin/main` has the `source-resolution` item and the `stop-operator-install-manifest` item.
- `validate --all` exits 1 until group 9, because the overlap requirements are long until then. A
  unit passes when no ERROR item is new and the WARNING count fell by `N`. For the unit's specs
  also run `openspec validate <spec> --type spec --strict --json 2>&1`: no WARNING for a
  rewritten requirement.
- The 14 review items of the guide run per unit, including the byte-equality command for each
  rewritten requirement.

### Gate for the open-change edits

After group 2, after group 9, and again in group 10, in a scratch copy of the branch's `openspec/`
folder, with the pinned CLI (1.14.1):

1. `openspec validate stop-operator-install-manifest --strict` and
   `openspec validate add-cross-namespace-source-grants --strict` exit 0 (after group 9; after
   group 2 only the second).
2. **The archive check.** Archive in the scratch copy, in this order: this change, then
   `stop-operator-install-manifest`, then `add-cross-namespace-source-grants`. Each
   `openspec archive <change> --yes` exits 0 and prints no "would refuse". The order is the
   only one required: the reverse cannot work on `main`, because the open changes name headers
   that this change creates.
3. After the three archives, `openspec validate --specs --strict --no-interactive` shows no
   WARNING except the placeholder Purpose of the new spec `cross-namespace-source-grants`. That
   one is written by the archiver of `add-cross-namespace-source-grants` (task line added in
   group 2). After one simulated direct edit of that Purpose,
   `openspec validate --all --strict --no-interactive` exits 0.
4. The byte-equality command of the guide, run against the delta of each edited block.

This is the check of the wanted behavior "keeps passing after the two open changes are
archived". It is part of the final gate (group 10) and of the gate of group 9.

### The pinned gate in `task` and in CI

**Context**: core's Taskfile pins `OPENSPEC_VERSION` once and has two tasks: `openspec:install`
(`npm install -g --ignore-scripts @fission-ai/openspec@{{.OPENSPEC_VERSION}}`) and
`openspec:check` (preconditions: the binary exists and its version equals the pin;
`OPENSPEC_TELEMETRY=0`; `openspec validate --all --strict --no-interactive`). Its CI job runs
`setup-node`, then `task openspec:install`, then `task openspec:check`, as steps in a job whose token
only reads. Its branch-publish job, which holds a package-write token, calls a task that
leaves the OpenSpec check out, because npm code that no lockfile pins must not run there (core
PR 139 lesson).
**Decision**: copy the form. Add `OPENSPEC_VERSION: 1.14.1` and the two tasks to the root
`Taskfile.yml` of this repo. Add three steps (setup Node, install, check) to the `Lint` job of
`.github/workflows/lint.yml`: that job has `permissions: contents: read` and no job in the file
writes. Add no step to a workflow with a write permission (`image-pr.yml`, `release.yml`,
`publish-fixtures.yml`, `docs.yml`, `deps-cascade.yml`, `test-e2e.yml`, `cascade-gates.yml`).
Pin `actions/setup-node` by commit SHA like the other actions of the file. If an aggregate task
runs the repo's other gates, the OpenSpec check stays out of the one a write-token job calls.
**Rationale**: the pin lives once; the install runs no script; no write token reaches the CLI.
**Alternative rejected**: a `Makefile` target (the repo says new automation goes in `.tasks/`
or `Taskfile.yml`).

### Delivery order

**Context**: the first draft planned two PRs, because four long requirements stayed on `main` and
the CI step would be red. The human's gate removes the cause: all 84 are split in this change.
**Decision**: one PR, from this branch. The gate group (10) is last. The branch is green at the
end: `openspec validate --all --strict --no-interactive` exits 0 on it, and so does
`task openspec:check`. Groups 1 to 9 are separate commits, each green under the per-unit gate
of this file. The archive of this change rides the PR (group 10). The two open changes archive
after the PR merges; the archive check of "The gate" shows that they can, in the order: this
change, `stop-operator-install-manifest`, `add-cross-namespace-source-grants`.
**Rationale**: it keeps Principle VIII's rule that every merge leaves `main` releasable and puts
no red required check on `main`. The two-PR plan is dropped.

## Risks / Trade-offs

- A rewrite that loses a rule while the validator is green. Counter: ledger, byte-equality
  command, clause table and strength table, and a review per unit.
- The Purpose of the new spec `cross-namespace-source-grants` is a placeholder after the archive of
  its change, a WARNING that fails the strict gate. No delta can carry it. The archiver of
  `add-cross-namespace-source-grants` writes the real text; the one task line added to that
  change in group 2 says so (the human's gate, 2026-10-20 15:40).
- An open change edited by this change can drift if its author edits the same deltas. Counter:
  the edits are in the open change's folder on this branch; the author merges this branch's
  delta files, and a conflict is a stop.
- A retargeted delta that does not match the new main headers makes the archive of
  `stop-operator-install-manifest` refuse. Counter: the archive check in "The gate", run after
  group 9 and in group 10.
- Section count: ten groups, more than about five; each ends green and the change runs as units.
- No new trust boundary, secret or personal data. The new CI step installs an npm package pinned
  by version with scripts off, in a read-only job; the lockfile gap is the same accepted residual
  as in core and cli, owner: this repo's maintainers.

## Open Questions

None. The delivery order, the four overlap requirements and the Purpose of
`cross-namespace-source-grants` are decided by the human's gate (2026-10-20 15:40).
