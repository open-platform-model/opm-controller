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

- `openspec validate --all --strict --no-interactive` under 1.14.1 exits 0 on a tree where the two
  open changes are archived, whichever order they and this change archive in.
- Every rewrite is a split or a move into scenarios. No rule, value, default, exit code,
  citation or scenario is lost or changed.
- The gate runs in `task` and in CI, pinned.

Non-goals: any code change; any change of meaning of a requirement or scenario, in a main spec or
in an open change; the proposal intent, the tasks or the code of an open change; the Purpose that
archive writes into the new spec `cross-namespace-source-grants` (see "Risks").

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

**Context**: a rewrite of a requirement that an open change MODIFIES or REMOVES breaks that
change at its archive, and the open change replaces the main text anyway. So a long requirement
that an open change touches is fixed in the open change's delta, not in the main spec. The list
below is fresh: it compares the 84 long main requirements with every requirement the two deltas
name, and measures the delta blocks (`validate <change> --strict --json`, and a length count of
every delta block, MODIFIED included, because validate skips those).
**Decision**: edit only `specs/` of the open changes. The post-archive text of a MODIFIED block
is the delta's, so a delta block that would be over 500 after archive is split now, by the same
verbatim rule. Each edit is made in the implement unit "Open-change fixes" and recorded below.

| Open change | Requirement (header) | Delta op | Main length | Delta length | Where it is fixed |
| --- | --- | --- | --- | --- | --- |
| stop-operator-install-manifest | The controller's pod and Service stay those of the kustomize tree (operator-module) | MODIFIED | 581 | 670 | Delta: restore scenario "An object added to the kustomize tree only fails the test" byte equal to `origin/main`; split the body (form a or b) |
| stop-operator-install-manifest | Release published once after every release job (release-automation) | MODIFIED | 837 | 735 | Delta: split the body, keep every live scenario |
| stop-operator-install-manifest | Beta prerelease line (release-automation) | MODIFIED | 1475 | 1475 | Delta: split the body, keep every live scenario |
| stop-operator-install-manifest | Release install manifest with digest-pinned image (container-image-publish) | REMOVED | 685 | none | Nowhere: archive removes it. Not rewritten in the main spec |
| stop-operator-install-manifest | Release assets upload only to a draft release (release-automation) | MODIFIED | 488 | 437 | No change: 500 or less. Listed so the list is complete |
| add-cross-namespace-source-grants | OCIRepository lookup (source-resolution) | MODIFIED | 280 | 633 | Delta: the main text is short, the delta text would be long after archive; split it (form a or b) |
| add-cross-namespace-source-grants | SourceRefGrant CRD (cross-namespace-source-grants) | ADDED | none | 571 | Delta: split by the same rule |
| add-cross-namespace-source-grants | Grant-matching policy permits a reference (cross-namespace-source-grants) | ADDED | none | 547 | Delta: split by the same rule |

The first four rows are the long main requirements that an open change touches and that stay long
on `main`. `Release install manifest ...` is not fixed anywhere because archive removes it. The brief of the human named one long ADDED block; the measure shows two
(`validate` lists both), and the MODIFIED `OCIRepository lookup` would fail after archive, so
three blocks of `add-cross-namespace-source-grants` are split. The last three rows do not appear in the 84.

The three open-change blocks of `release-automation` and `operator-module` are not split in the
main spec, so the main specs keep 4 long requirements (`operator-module` 1, `release-automation`
2, `container-image-publish` 1) until the open change archives. This change rewrites 84 - 4 = 80.

Edit rules for a delta, besides the verbatim rule:

- A MODIFIED block keeps its header and every live scenario name and text from `origin/main`,
  byte for byte. New ADDED requirements under new names carry the split-off behaviors, with their
  own scenarios. A live scenario never moves from the MODIFIED parent to an ADDED block.
- An ADDED block of the open change that was long is split into ADDED blocks of the same delta
  (form c2 inside one delta: the long header disappears, the new headers appear). No name in the
  delta is in two operations. A new name is checked against the names this change adds to the
  same spec.
- A scenario that the open change wrote is kept byte for byte, like a scenario of a main spec.
- Each edit is checked by `openspec validate <open-change> --strict` (exit 0), by the byte-equality
  command (against `origin/main` for a MODIFIED block, against the open change's delta on
  `origin/main` for an ADDED block), and by the archive simulation of "The gate".

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
**Decision**: Purpose, open-change fixes, six rewrite units, gate.

| Group | Content | Requirements | N after it |
| --- | --- | --- | --- |
| 1 Purpose | `source-resolution` Purpose, direct edit | 0 | 0 |
| 2 Open-change fixes | Deltas of the two open changes | the 2 long ADDED blocks | 2 |
| 3 Unit 1 | finalizer-and-deletion 11, instance-health 4 | 15 | 17 |
| 4 Unit 2 | release-automation 11, container-image-publish 2 | 13 | 30 |
| 5 Unit 3 | reconcile-loop-assembly 9, operator-module 4 | 13 | 43 |
| 6 Unit 4 | registration-acceptance 9, cascade-receiver 5 | 14 | 57 |
| 7 Unit 5 | example-test-modules 7, prune-stale-resources 6 | 13 | 70 |
| 8 Unit 6 | drift-detection 5, modulepackage-kernel-rendering 5, example-module-publishing 2 | 12 | 82 |
| 9 Gate | Taskfile tasks, CI step, final archive simulation | 0 | 82 |

`N` is the number of WARNING lines the whole change on the branch has removed from
`validate --all --strict`, counted against `origin/main` (86 WARNING lines: the 84 long main
requirements and the 2 long ADDED blocks of the open change). Group 2 removes 2 lines. After
group 8, 4 lines remain: the four long main requirements of the overlap list. The count per spec
is from `openspec validate --specs --strict --json` on `origin/main` (13 release-automation,
11 finalizer-and-deletion, 9 reconcile-loop-assembly, 9 registration-acceptance, 7
example-test-modules, 6 prune-stale-resources, 5 each operator-module,
modulepackage-kernel-rendering, drift-detection and cascade-receiver, 4 instance-health, 3
container-image-publish, 2 example-module-publishing). The units that touch a spec of an open
change (2 and 3) leave the overlap requirement of that spec alone.
**Alternative rejected**: rewrite all 84 in the main specs and move the open changes onto the
new text. It breaks their MODIFIED and REMOVED blocks (a REMOVED header that was split no longer
matches, and a MODIFIED block replaces the shortened text with its own), and the open-change rule of the brief rules
it out.

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
contains both open changes, as edited by group 2 once it is done. Pass rules:

- `validate <change> --strict` alone is not enough: it skips long MODIFIED text.
- ERRORs are compared as sets of `type id path`; the set may shrink and never gains an item.
  `origin/main` has the `source-resolution` item and the `stop-operator-install-manifest` item.
- `validate --all` exits 1 until the end, because the four overlap requirements stay long. A
  unit passes when no ERROR item is new and the WARNING count fell by `N`. For the unit's specs
  also run `openspec validate <spec> --type spec --strict --json 2>&1`: no WARNING for a
  rewritten requirement.
- The 14 review items of the guide run per unit, including the byte-equality command for each
  rewritten requirement.

### Gate for the open-change edits

After group 2 and again in group 9, in a scratch copy:

1. `openspec validate stop-operator-install-manifest --strict` and
   `openspec validate add-cross-namespace-source-grants --strict` exit 0.
2. Archive simulation, order A: archive `add-cross-namespace-source-grants`, then
   `stop-operator-install-manifest`, then this change. Order B: this change first, then the two.
   Each `openspec archive <change> --yes` exits 0 and prints no refusal.
3. After the simulated archives, `openspec validate --all --strict --no-interactive` exits 0,
   after one simulated direct edit: the placeholder Purpose that archive writes into the new
   spec `cross-namespace-source-grants` is replaced by real text (a WARNING otherwise; measured
   in a scratch copy). The real fix of that Purpose belongs to the archive step of the open
   change.
4. The byte-equality command of the guide, run against the delta of each edited block.

This is the check of the wanted behavior "keeps passing after the two open changes are
archived".

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

**Context**: until both open changes archive, four long requirements stay in the main specs, so
`task openspec:check` and the new CI step exit 1 on `main` and on this branch. And the open
changes cannot archive cleanly before the edits of group 2 land.
**Decision**: this is a delivery question for the human, not a spec question. The recommended
order is two PRs from this branch: PR 1 holds groups 1 to 8 (specs and Purpose; no gate step),
PR 2 holds group 9 and lands after both open changes have archived, so the step is green on
`main` when it arrives. A single PR would put a red required check on `main`. Group 9's Taskfile
tasks do not fail anything and could ride PR 1; the `lint.yml` step is what must wait.
**Rationale**: it is the order of the guide (gate PRs after the open changes land) and keeps
Principle VIII's rule that every merge leaves `main` releasable.

## Risks / Trade-offs

- A rewrite that loses a rule while the validator is green. Counter: ledger, byte-equality
  command, clause table and strength table, and a review per unit.
- The Purpose of the new spec `cross-namespace-source-grants` is a placeholder after the archive of
  its change, a WARNING that fails the gate. No delta can carry it. The archive step of the open
  change must replace it by direct edit; this change cannot do it before the spec exists. Its
  proposal intent and tasks are out of scope here, so this is reported to the human.
- The gate cannot pass on `main` before both open changes archive. See "Delivery order".
- An open change edited by this change can drift if its author edits the same deltas. Counter:
  the edits are in the open change's folder on this branch; the author merges this branch's
  delta files, and a conflict is a stop.
- Section count: nine groups, more than about five; each ends green and the change runs as units.
- No new trust boundary, secret or personal data. The new CI step installs an npm package pinned
  by version with scripts off, in a read-only job; the lockfile gap is the same accepted residual
  as in core and cli, owner: this repo's maintainers.

## Open Questions

- Delivery order (see above): two PRs, with the CI step after the open changes archive. The human
  decides.
- The Purpose of `cross-namespace-source-grants` at archive time: who writes it.
