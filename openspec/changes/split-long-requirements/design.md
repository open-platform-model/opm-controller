## Context

The validator measures the body of a requirement: the lines between `### Requirement:` and the
first `#` line, with blank lines, fenced lines and `**Key**:` lines dropped (the last only when
other text exists), joined by newlines. Scenarios do not count. More than 500 UTF-16 units is a
WARNING; `--strict` fails on it. The check runs on main specs and on ADDED blocks of a delta, not
on MODIFIED blocks. A MODIFIED block that stays over 500 therefore passes
`validate <change> --strict`, archives, and then fails `validate --all --strict` on the spec.

Current state, measured on `origin/main` 7713cc86 with OpenSpec 1.14.1 in a scratch copy
(`validate --all --strict --json`, 2026-10-19): 196 requirements over 500 in 44 specs, equal to
`explore/T1.3.md` spec by spec; 8 ERRORs. Seven are specs without `## Purpose` (the eighth is in
the open change `stop-operator-install-manifest`). The `example-module-publishing` spec has a
placeholder Purpose, a WARNING. The repo has 59 main specs. Two open changes exist and this
change does not touch them: `add-cross-namespace-source-grants` (specs
`cross-namespace-source-grants` and `source-resolution`) and `stop-operator-install-manifest`
(specs `container-image-publish`, `operator-module`, `release-automation`).

The swarm guide `design/T2.1.md` (approved 2026-10-19 13:40) holds the rules. This file applies
them to this repo and does not restate them in full.

## Goals / Non-Goals

Goals:

- Every main spec passes `openspec validate --specs --strict` under 1.14.1.
- Every rewrite is a split or a move into scenarios. No rule, value, default, exit code,
  citation or scenario is lost or changed.
- The extras of `proposal.md` are done.

Non-goals: any code change; any change of meaning; editing the two open changes; fixing their
ADDED blocks over 500 or their placeholder Purpose (the authors' work at archive).

## Research & Decisions

### Which form per requirement

**Context**: A long requirement can shrink four ways. The guide allows exactly these forms.
**Explored**: `explore/T1.1.md`, the runs in `reviews/T1.4.md`, guide section 3. Counts and kind
labels: `explore/T1.3.md`.
**Decision**: choose per requirement by reading the text, never from the T1.3 kind label (a
heuristic):

| The excess is | Form |
| --- | --- |
| Examples, lists, formats, rationale around one behavior | a: MODIFIED, shorter body, more scenarios |
| Several behaviors, and the old scenarios fit the parent that stays | b: MODIFIED parent plus ADDED splits |
| Several behaviors, and the old scenarios belong to different behaviors | c2: REMOVED plus ADDED under new names; the ledger is mandatory |
| A name that must change, with no change of rule | d: RENAMED plus MODIFIED under the new name, alone |

Almost every long requirement here states several rules (T1.3: "several behaviors" is the
dominant label), so b and c2 lead, with detail moved to scenarios by a.
**Rationale**: b and c2 give one testable behavior per requirement, which is the reason for the
warning. Form a alone on a requirement with several behaviors would only hide them in scenarios.

Not allowed: moving text into `**Key**:` lines, under a `#####` heading, or into a fenced block.
They pass the length check and remove nothing. Moving a live scenario out of a MODIFIED parent
into an ADDED block is an ERROR ("b-move"). The operation pairs that are errors on one name are
in guide section 3.

### What this repo adds to the guide

**Context**: the `specs` rules of `openspec/config.yaml`.
**Decision**: these hold for every body and every new ADDED requirement:

- RFC 2119 keywords (SHALL, MUST) stay in every body, with the same strength as before.
- A citation (an enhancement id beside a reason in words, never alone, never a decision number)
  stays in the body of the block that carries the rule it supports, never in the header.
- Scenarios describe observable controller behavior (status changes, created resources, logs).
  Where the old text names transient and permanent outcomes, the split keeps both as scenarios
  (the config's "Transient Error vs Stalled" example). The guide's rule comes first: a scenario
  never adds a case that is not in the old text; a new case is a question to the human.
- Before a split or a c2, grep the repo for the old header text
  (`grep -rn --exclude-dir=archive "<old header>" .`) and fix stale references in the same unit.

**Rationale**: these are the repo's rules for specs; the guide asks the rewriter to read them
before the first split.

### Order of the units

**Context**: Rewriting a main spec that an open change also edits makes the open change's
MODIFIED block replace the shortened text with the old long one and re-add dropped scenarios (a
c2 leaves it with only the INFO "Archive would refuse"). The human decided on 2026-10-19 10:50
that the colliding batches run last.
**Decision**: Purpose-only specs, then B2, B4, B6, B7 (no spec of an open change), then a stop
for the human to land both open changes, then a fresh inventory, then B1, B3, B5, then the
`source-resolution` Purpose, then the gate.

| Unit | Why here | N after it |
| --- | --- | --- |
| Purpose for 6 specs | No collision, no length work; removes 6 of the 8 ERRORs on main | 0 |
| B2, B4, B6, B7 | No spec of these batches is edited by an open change | 28, 56, 84, 112 |
| Pause: the human lands both open changes | They edit `release-automation` (B1), `container-image-publish` (B3), `operator-module` (B5) and `source-resolution` | |
| B1, B3, B5 | After the pause, with an inventory re-run first | 140, 168, 196, or the fresh count |
| `source-resolution` Purpose | Its change must have landed | unchanged |
| Gate (Taskfile, CI) | A step that runs `--all --strict` would fail its own PR while the open changes carry warnings and errors | unchanged |

**Rationale**: each unit stays independent of the open changes until they land, as the human
ruled. The last column assumes an unchanged inventory; B1, B3 and B5 re-measure first, and `N`
is the fresh cumulative count.
**Alternative rejected**: all units first, with the open changes moved onto the new text
afterwards. It breaks their MODIFIED blocks (guide section 3, "unsafe because silent"), and
editing them is out of scope.

### Reversibility

high (hours). Spec text and a CI step are reverted by a PR. No contract changes and no accepted
decision record is edited.

## The scenario ledger

Each unit writes its ledger in this file before it writes the delta of a spec, under the heading
`### Ledger: <spec>`. It is not written now, because the inventory moves for B1, B3 and B5 and
the requirement text is read when the unit starts. The first table has the form of guide
section 4:

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |

A second table lists the rule clauses of the old body (one per SHALL, MUST, number, list item,
format) and the clause or scenario THEN line that carries each. Rules of the ledger:

- Forms a, b, d: every old scenario name stays under the same requirement, byte for byte.
- Form c2: every old scenario name appears in exactly one new requirement, WHEN and THEN byte
  for byte. OpenSpec does not check this: an old scenario that is in no new requirement
  disappears with no error.
- A clause that moves from a body SHALL into a scenario THEN keeps its SHALL there, and only a
  default or an exact value may move this way. A clause that adds or changes a condition stays
  in the body. A separate behavior goes to b or c2.

At this step no ledger row exists.

## The gate

The gate of guide section 6 runs in a scratch copy of the worktree's `openspec/` folder, never in
the worktree. It uses the pinned CLI (1.14.1) and `jq`. Values for this repo:

```bash
WT=/var/home/emil/dev/open-platform-model/agent-swarm/openspec-requirement-split/worktrees/opm-operator-split-long-requirements
CH=split-long-requirements
S=<scratch folder under the session scratchpad>
N=<requirements rewritten on the branch so far, earlier units included>
# LAST=1 only in the final unit
```

Then the block of guide section 6 runs unchanged: it archives the change in the copy, runs
`validate --all --strict`, fails on any "Archive would refuse" line and on any new ERROR item,
and either checks that WARNING lines fell by `N` or, with `LAST=1`, requires
`openspec validate --specs --strict` to exit 0. The block reads the baseline with
`git -C "$WT" archive main openspec`. The worktree branch is `docs/split-long-requirements`, so
`main` there is the local `main` ref of the repo; if it is behind `origin/main`, the unit uses
`origin/main` in that line and records it.

Pass rules:

- `validate split-long-requirements --strict` alone is not enough: it skips long MODIFIED text.
- ERRORs are compared as sets of `type id path`. `main` has 8 ERROR items until the Purpose
  units land. The set may shrink and never gains an item.
- Earlier units: other batches still hold warnings, so `validate --all` exits 1. The unit passes
  when no ERROR item is new and WARNING lines fell by `N`. For the batch's specs also run
  `openspec validate <spec> --type spec --strict --json 2>&1` and read it.
- Last unit (`LAST=1`): passes when `validate --specs --strict` exits 0 after the archive in the
  copy and the change passes `validate --strict`. Both open changes have landed by then, so
  `validate --all --strict` is expected to exit 0 as well; the unit reports its exit code and
  names the cause if it is not 0.
- The 14 review items of guide section 7 run per unit, including the byte-equality command for
  each rewritten requirement.

## The pinned gate (Taskfile and CI)

**Context**: this repo has no OpenSpec task and no workflow that runs the CLI
(`grep -rn openspec Taskfile.yml .tasks .github/workflows` finds nothing). The `cli` repo's
`openspec:install` and `openspec:check` tasks are the model.
**Decision**: add two tasks in `.tasks/`: `openspec:install` (npm install of
`@fission-ai/openspec` at exactly 1.14.1) and `openspec:check`
(`openspec validate --all --strict --no-interactive`, with a clear message when the CLI is
missing), included from `Taskfile.yml` like the other `.tasks/*.yaml` files. Add one step in the
required `Lint` job of `.github/workflows/lint.yml` that installs the pinned version through
the task and runs the check. The pin is written once, in the task, and the workflow carries no
version literal (the repo's own rule for the opm CLI pin).
**Rationale**: one pin in one place; `Lint` is the required check, so a regression of a long
requirement fails the PR.
**Alternative rejected**: a separate workflow file. It adds a trigger and a permission set and
is not a required check.
**Order**: delivered after both open changes land, as the human chose on 2026-10-19 10:50.
Before that, the step would fail its own PR.
**Confirmation**: the unit shows the task failing on a scratch copy with one long requirement
restored.
**Security**: the step runs a pinned npm package. The unit loads the supply-chain rules
(exact version, cooldown, install scripts) and the workflow-security rules before it edits
`lint.yml`.

## Risks / Trade-offs

- A rewrite that loses a rule while the validator is green. Counter: ledger, byte-equality
  command, clause table, and a review per unit (guide section 7).
- The two open changes move text before B1, B3 and B5. Counter: the pause and the fresh
  inventory. An "Archive would refuse" line is reported, never repaired by editing the other
  change.
- Section count: 11 task groups (9 units plus a pause and the gate), against the config's "at most about seven sections" rule. The
  groups are mechanical units of one kind, each ends green, and DECISIONS.md 2026-10-19 09:50
  sets one change per repo run as several units. The proposal gate may ask instead for two
  changes (before and after the open changes land).
- No new trust boundary, secret or personal data. The one new CI step is covered above.

## Open Questions

- Does the proposal gate accept 11 groups in one change, or want two changes?
- The counts of B1, B3 and B5 are known only after the fresh inventory.
