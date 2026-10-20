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

### Record: Purpose of source-resolution

Added by direct edit above `## Requirements` of `openspec/specs/source-resolution/spec.md`
(226 characters). Each part is read from a requirement of the spec: the lookup of OCIRepository,
GitRepository and Bucket sources and the artifact URL, revision and digest ("OCIRepository
lookup", "ArtifactRef carries full metadata"), the sentinel errors ("Typed source errors") and the
watch that enqueues reconciliation ("Source watch triggers reconciliation"). No requirement text
changed.

### Ledger: add-cross-namespace-source-grants, cross-namespace-source-grants

Form c2 inside the one delta (ADDED blocks only). Old lengths 571 and 547; new lengths in the
clause table.

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| Grant is namespaced and lives with the source | SourceRefGrant CRD | no |
| Optional referent name | SourceRefGrant spec lists referrers and referents | no |
| Matching grant permits the reference | Grant-matching policy permits a reference | no |
| No matching grant denies the reference | Grant-matching policy permits a reference | no |
| Revocation takes effect on next reconcile | Grant evaluation runs on every reconcile | no |

| Old clause | New requirement (length) | Note |
| --- | --- | --- |
| The operator MUST define a namespaced `SourceRefGrant` custom resource ... | SourceRefGrant CRD (403) | verbatim |
| A `SourceRefGrant` lives in the namespace that owns the source ... | SourceRefGrant CRD | verbatim |
| Its spec MUST contain a `from` list ... and a `to` list ... | SourceRefGrant spec lists referrers and referents (188) | subject "Its spec" written "The spec of a `SourceRefGrant`" (header name) |
| The CRD MUST be additive and permissive only, it grants access and never denies it; absence of a matching grant is a denial | SourceRefGrant CRD | em-dash written as a colon |
| When the master switch is enabled, a cross-namespace reference ... MUST be permitted if and only if ... | Grant-matching policy permits a reference (441) | verbatim |
| Evaluation MUST be performed on every reconcile so that deleting a grant revokes access on the next pass | Grant evaluation runs on every reconcile (169) | subject and condition "When the master switch is enabled, evaluation of the grant-matching policy" copied from the old sentence and header |

### Ledger: add-cross-namespace-source-grants, source-resolution

Form b: the MODIFIED block "OCIRepository lookup" keeps its header and its seven live scenarios
byte for byte. Delta text 633, now 496. The sentence about the empty or equal namespace becomes
an ADDED requirement with one new scenario that restates the case the sentence names.

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| Source exists and is ready | OCIRepository lookup | no |
| Source not found | OCIRepository lookup | no |
| Source exists but not ready | OCIRepository lookup | no |
| Source ready but no artifact | OCIRepository lookup | no |
| Unsupported source kind | OCIRepository lookup | no |
| Cross-namespace reference denied by default | OCIRepository lookup | no |
| Cross-namespace reference permitted by policy | OCIRepository lookup | no |
| (new) Empty namespace resolves in the release namespace | Same-namespace lookup uses the release namespace | new: restates "When `SourceReference.Namespace` is empty" |

| Old clause | New requirement (length) | Note |
| --- | --- | --- |
| The `internal/source` package MUST provide a `Resolve` function ... | OCIRepository lookup (496) | verbatim |
| The function MUST support OCIRepository, GitRepository, and Bucket source kinds | OCIRepository lookup | verbatim |
| When `SourceReference.Namespace` is empty or equal to the release namespace, `Resolve` MUST look up the source in the release namespace | Same-namespace lookup uses the release namespace (136) | verbatim |
| When `SourceReference.Namespace` names a different namespace, `Resolve` MUST consult the `CrossNamespacePolicy` and MUST NOT read the foreign source unless the policy permits the reference | OCIRepository lookup | verbatim |

The two ADDED blocks of this delta that were not long are unchanged. No keyword was added, so
the strength table has no row for group 2.

### Ledger: finalizer-and-deletion

Forms and lengths (old body, then each new body):

| Old requirement | Form | Old length | New blocks (length) |
| --- | --- | --- | --- |
| Finalizer registration | c2 (REMOVED plus ADDED under new names) | 695 | The cleanup finalizer is registered in Phase 0 (410); The finalizer-add reconcile requests an immediate requeue (362) |
| Deletion cleanup with prune enabled | b (MODIFIED parent kept, ADDED splits) | 667 | Deletion cleanup with prune enabled (426); The deletion completes once every deleted object is gone (301) |
| Finalizer retained on DeletionSAMissing stall | b (MODIFIED parent kept, ADDED splits) | 887 | Finalizer retained on DeletionSAMissing stall (352); A DeletionSAMissing stall ends by a restored ServiceAccount or the orphan annotation (497); A DeletionSAMissing stall ends when spec.prune is set to false (422) |
| Recovery signals of a stalled deletion trigger a reconcile | b (MODIFIED parent kept, ADDED splits) | 1157 | Recovery signals of a stalled deletion trigger a reconcile (461); A recovery signal enqueues the instance and the existing deletion rules apply (307); A ServiceAccount created before its RBAC stalls the deletion with ImpersonationFailed (328); Other annotation changes and status-only writes trigger no reconcile (389) |
| A kept claim never holds the finalizer | c2 (REMOVED plus ADDED under new names) | 680 | A kept claim does not hold the finalizer (477); An inventory of only kept claims releases the finalizer without a delete (264) |
| Deletion cleanup judges every object with the delete verdict | b (MODIFIED parent kept, ADDED splits) | 1578 | Deletion cleanup judges every object with the delete verdict (375); The deletion cleanup judges with the identities of the instance (258); An empty instance identity is judged with no identity (363); A skipped object does not hold the finalizer and a failed step does (444) |
| Deletion cleanup runs the library's deletion plan | c2 (REMOVED plus ADDED under new names) | 1047 | The deletion cleanup deletes an object only as the library's deletion plan names it (454); The deletion cleanup runs one plan per identity (398); A kept claim is not an entry of any plan (110) |
| The cleanup finalizer follows the library's hold verdict | b (MODIFIED parent kept, ADDED splits) | 1823 | The cleanup finalizer follows the library's hold verdict (468); The inputs of the hold verdict (353); Each hold verdict keeps its status (428); Only kept claims release the finalizer before the identity is looked at (458); A release without a read emits DeletionUnconfirmed and no ClaimsKept (388) |
| The finalizer is kept until the deleted objects are gone | b (MODIFIED parent kept, ADDED splits) | 1001 | The finalizer is kept until the deleted objects are gone (462); A deleted object that still exists holds the finalizer with DeletionInProgress (389); Each reconcile of a waiting deletion judges the inventory again from the cluster (209) |
| A wait survives the loss of the ServiceAccount | b (MODIFIED parent kept, ADDED splits) | 3793 | A wait survives the loss of the ServiceAccount (456); The record is read only from the Ready condition at the start of the reconcile (436); The ServiceAccount is gone in exactly one case (450); A NotFound that names another object does not make the ServiceAccount gone (382); A gone ServiceAccount releases the finalizer (464); Other answers hold the finalizer (470); A held finalizer is retried with a Ready message (481); The Ready message of a failed recheck is bounded (311); A failed recheck is reported as DeletionBlocked (481); A ServiceAccount that loses its rights does not release the deletion by itself (411); The wait reasons are set only during a deletion (347); Without the record the stall reasons hold the finalizer (276) |
| A deletion that does not finish is reported as blocked | b (MODIFIED parent kept, ADDED splits) | 847 | A deletion that does not finish is reported as blocked (387); The DeletionBlocked message names the remaining objects and the ways out (197); A blocked deletion keeps the finalizer and keeps rechecking (441) |

`grep -rn --exclude-dir=archive "<old header>" .` found no reference to a header that this unit renames or removes outside this spec; the in-spec references name headers that stay (form b parents) or are updated here.

Scenarios:

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| First reconcile adds finalizer | The cleanup finalizer is registered in Phase 0, First reconcile adds finalizer | no (byte equal, checked by the item 3 command) |
| Subsequent reconciles preserve finalizer | The cleanup finalizer is registered in Phase 0, Subsequent reconciles preserve finalizer | no (byte equal, checked by the item 3 command) |
| Finalizer-add reconcile requests immediate requeue | The finalizer-add reconcile requests an immediate requeue, Finalizer-add reconcile requests immediate requeue | no (byte equal, checked by the item 3 command) |
| Status conditions appear after creation without manual triggers | The cleanup finalizer is registered in Phase 0, Status conditions appear after creation without manual triggers | no (byte equal, checked by the item 3 command) |
| Delete all owned resources on CR deletion | Deletion cleanup with prune enabled, Delete all owned resources on CR deletion | no (byte equal, checked by the item 3 command) |
| Safety exclusions during deletion | Deletion cleanup with prune enabled, Safety exclusions during deletion | no (byte equal, checked by the item 3 command) |
| Claims are kept on deletion by default | Deletion cleanup with prune enabled, Claims are kept on deletion by default | no (byte equal, checked by the item 3 command) |
| Claims are deleted on deletion when the instance opts out | Deletion cleanup with prune enabled, Claims are deleted on deletion when the instance opts out | no (byte equal, checked by the item 3 command) |
| Release not garbage-collected while stalled on DeletionSAMissing | Finalizer retained on DeletionSAMissing stall, Release not garbage-collected while stalled on DeletionSAMissing | no (byte equal, checked by the item 3 command) |
| Setting the orphan annotation releases a stalled deletion without waiting | Recovery signals of a stalled deletion trigger a reconcile, Setting the orphan annotation releases a stalled deletion without waiting | no (byte equal, checked by the item 3 command) |
| The return of the ServiceAccount completes a stalled deletion without waiting | Recovery signals of a stalled deletion trigger a reconcile, The return of the ServiceAccount completes a stalled deletion without waiting | no (byte equal, checked by the item 3 command) |
| An annotation change on a live instance is not a trigger | Recovery signals of a stalled deletion trigger a reconcile, An annotation change on a live instance is not a trigger | no (byte equal, checked by the item 3 command) |
| A ServiceAccount that returns before its RBAC does not complete the deletion | Recovery signals of a stalled deletion trigger a reconcile, A ServiceAccount that returns before its RBAC does not complete the deletion | no (byte equal, checked by the item 3 command) |
| Deletion with only claims left completes | An inventory of only kept claims releases the finalizer without a delete, Deletion with only claims left completes | no (byte equal, checked by the item 3 command) |
| Another failed delete still holds the finalizer | A kept claim does not hold the finalizer, Another failed delete still holds the finalizer | no (byte equal, checked by the item 3 command) |
| Object adopted by another instance survives the deletion | Deletion cleanup judges every object with the delete verdict, Object adopted by another instance survives the deletion | no (byte equal, checked by the item 3 command) |
| Object of another instance survives the deletion | Deletion cleanup judges every object with the delete verdict, Object of another instance survives the deletion | no (byte equal, checked by the item 3 command) |
| Deletion before an identity change is settled removes the live workload | Deletion cleanup judges every object with the delete verdict, Deletion before an identity change is settled removes the live workload | no (byte equal, checked by the item 3 command) |
| Deletion with no recorded identity leaves an annotated object | Deletion cleanup judges every object with the delete verdict, Deletion with no recorded identity leaves an annotated object | no (byte equal, checked by the item 3 command) |
| A replaced object holds the finalizer for one more attempt | Deletion cleanup judges every object with the delete verdict, A replaced object holds the finalizer for one more attempt | no (byte equal, checked by the item 3 command) |
| Higher kind weight is deleted first | The deletion cleanup deletes an object only as the library's deletion plan names it, Higher kind weight is deleted first | no (byte equal, checked by the item 3 command) |
| Deletes use Foreground propagation | The deletion cleanup deletes an object only as the library's deletion plan names it, Deletes use Foreground propagation | no (byte equal, checked by the item 3 command) |
| A kept claim is never sent to the plan | A kept claim is not an entry of any plan, A kept claim is never sent to the plan | no (byte equal, checked by the item 3 command) |
| An object of the previous identity is deleted by the second plan | The deletion cleanup runs one plan per identity, An object of the previous identity is deleted by the second plan | no (byte equal, checked by the item 3 command) |
| Prune disabled releases without a read | The cleanup finalizer follows the library's hold verdict, Prune disabled releases without a read | no (byte equal, checked by the item 3 command) |
| A failed step holds the finalizer | The cleanup finalizer follows the library's hold verdict, A failed step holds the finalizer | no (byte equal, checked by the item 3 command) |
| A missing ServiceAccount holds without a plan | The cleanup finalizer follows the library's hold verdict, A missing ServiceAccount holds without a plan | no (byte equal, checked by the item 3 command) |
| Force-orphan lifts only a missing ServiceAccount | The cleanup finalizer follows the library's hold verdict, Force-orphan lifts only a missing ServiceAccount | no (byte equal, checked by the item 3 command) |
| Only kept claims and a missing ServiceAccount | The cleanup finalizer follows the library's hold verdict, Only kept claims and a missing ServiceAccount | no (byte equal, checked by the item 3 command) |
| A claim the instance may delete still needs the ServiceAccount | The cleanup finalizer follows the library's hold verdict, A claim the instance may delete still needs the ServiceAccount | no (byte equal, checked by the item 3 command) |
| A terminating Deployment holds the finalizer | The finalizer is kept until the deleted objects are gone, A terminating Deployment holds the finalizer | no (byte equal, checked by the item 3 command) |
| The finalizer goes when the objects are gone | The finalizer is kept until the deleted objects are gone, The finalizer goes when the objects are gone | no (byte equal, checked by the item 3 command) |
| Nothing left at the check releases in the same reconcile | The finalizer is kept until the deleted objects are gone, Nothing left at the check releases in the same reconcile | no (byte equal, checked by the item 3 command) |
| An object the garbage collector has not removed yet is waited for | The finalizer is kept until the deleted objects are gone, An object the garbage collector has not removed yet is waited for | no (byte equal, checked by the item 3 command) |
| A name taken by a new object does not hold the finalizer | The finalizer is kept until the deleted objects are gone, A name taken by a new object does not hold the finalizer | no (byte equal, checked by the item 3 command) |
| Turning prune off releases a waiting deletion | The finalizer is kept until the deleted objects are gone, Turning prune off releases a waiting deletion | no (byte equal, checked by the item 3 command) |
| ServiceAccount deleted during the wait | A wait survives the loss of the ServiceAccount, ServiceAccount deleted during the wait | no (byte equal, checked by the item 3 command) |
| A blocked deletion is released the same way | A wait survives the loss of the ServiceAccount, A blocked deletion is released the same way | no (byte equal, checked by the item 3 command) |
| RBAC removed during the wait holds | A wait survives the loss of the ServiceAccount, RBAC removed during the wait holds | no (byte equal, checked by the item 3 command) |
| One refused object next to readable ones holds | A wait survives the loss of the ServiceAccount, One refused object next to readable ones holds | no (byte equal, checked by the item 3 command) |
| A NotFound that is not the ServiceAccount's does not release | A wait survives the loss of the ServiceAccount, A NotFound that is not the ServiceAccount's does not release | no (byte equal, checked by the item 3 command) |
| A recreated ServiceAccount does not release | A wait survives the loss of the ServiceAccount, A recreated ServiceAccount does not release | no (byte equal, checked by the item 3 command) |
| A failed lookup of the ServiceAccount during the wait keeps the finalizer | A wait survives the loss of the ServiceAccount, A failed lookup of the ServiceAccount during the wait keeps the finalizer | no (byte equal, checked by the item 3 command) |
| An Unauthorized answer does not release | A wait survives the loss of the ServiceAccount, An Unauthorized answer does not release | no (byte equal, checked by the item 3 command) |
| A wait whose rechecks keep failing is reported as blocked | A wait survives the loss of the ServiceAccount, A wait whose rechecks keep failing is reported as blocked | no (byte equal, checked by the item 3 command) |
| The message of a failed recheck is bounded | A wait survives the loss of the ServiceAccount, The message of a failed recheck is bounded | no (byte equal, checked by the item 3 command) |
| Identity lost before every delete was sent | A wait survives the loss of the ServiceAccount, Identity lost before every delete was sent | no (byte equal, checked by the item 3 command) |
| A readable terminating object still holds | A wait survives the loss of the ServiceAccount, A readable terminating object still holds | no (byte equal, checked by the item 3 command) |
| The record is ignored on a live object | A wait survives the loss of the ServiceAccount, The record is ignored on a live object | no (byte equal, checked by the item 3 command) |
| The shipped roles cannot write the record | A wait survives the loss of the ServiceAccount, The shipped roles cannot write the record | no (byte equal, checked by the item 3 command) |
| A Pod that cannot terminate blocks the deletion visibly | A deletion that does not finish is reported as blocked, A Pod that cannot terminate blocks the deletion visibly | no (byte equal, checked by the item 3 command) |
| A blocked deletion completes when the object goes | A deletion that does not finish is reported as blocked, A blocked deletion completes when the object goes | no (byte equal, checked by the item 3 command) |
| The threshold is reached without a sleep | A deletion that does not finish is reported as blocked, The threshold is reached without a sleep | no (byte equal, checked by the item 3 command) |

New scenarios (each restates a case the old text names):

| New place | New scenario | Old text that names the case |
| --- | --- | --- |
| The deletion completes once every deleted object is gone | The deletion completes in a later reconcile | the sentence of this block |
| A DeletionSAMissing stall ends by a restored ServiceAccount or the orphan annotation | A restored ServiceAccount prunes the inventory | the sentence of this block |
| A DeletionSAMissing stall ends by a restored ServiceAccount or the orphan annotation | The orphan annotation removes the finalizer | the sentence of this block |
| A DeletionSAMissing stall ends when spec.prune is set to false | Disabling prune ends the stall | the sentence of this block |
| A recovery signal enqueues the instance and the existing deletion rules apply | The orphan annotation takes the orphan-exit path | the sentence of this block |
| A recovery signal enqueues the instance and the existing deletion rules apply | The return of the ServiceAccount prunes as that ServiceAccount | the sentence of this block |
| A ServiceAccount created before its RBAC stalls the deletion with ImpersonationFailed | The reconcile of an early ServiceAccount is forbidden | the sentence of this block |
| Other annotation changes and status-only writes trigger no reconcile | A status-only write is not a trigger | the sentence of this block |
| Other annotation changes and status-only writes trigger no reconcile | A ModulePackage keeps waiting for its stalled recheck | the sentence of this block |
| The deletion cleanup judges with the identities of the instance | An object that carries the previous identity is the instance's own | the sentence of this block |
| An empty instance identity is judged with no identity | An object annotated for this instance stays when no identity is recorded | the sentence of this block |
| A skipped object does not hold the finalizer and a failed step does | A skipped object stays and does not hold the finalizer | the sentence of this block |
| A skipped object does not hold the finalizer and a failed step does | A read that fails other than NotFound holds the finalizer | the sentence of this block |
| The inputs of the hold verdict | Force-orphan is true only for the literal true | the sentence of this block |
| Each hold verdict keeps its status | A failed identity stalls with ImpersonationFailed | the sentence of this block |
| Each hold verdict keeps its status | Force-orphan emits OrphanedOnDeletion and clears the inventory | the sentence of this block |
| Only kept claims release the finalizer before the identity is looked at | An empty plan releases before the identity is looked at | the sentence of this block |
| A release without a read emits DeletionUnconfirmed and no ClaimsKept | One DeletionUnconfirmed event states the number of claims | the sentence of this block |
| A deleted object that still exists holds the finalizer with DeletionInProgress | A deleted object that still exists ends the reconcile with DeletionInProgress | the sentence of this block |
| Each reconcile of a waiting deletion judges the inventory again from the cluster | Each reconcile of a waiting deletion judges the inventory again | the sentence of this block |
| The record is read only from the Ready condition at the start of the reconcile | The record is read from the Ready condition at the start of the reconcile | the sentence of this block |
| The ServiceAccount is gone in exactly one case | The ServiceAccount is gone when its own read is answered NotFound | the sentence of this block |
| A NotFound that names another object does not make the ServiceAccount gone | A NotFound that names another object does not release | the sentence of this block |
| A gone ServiceAccount releases the finalizer | A gone ServiceAccount releases the finalizer without a read | the sentence of this block |
| Other answers hold the finalizer | A server error holds the finalizer | the sentence of this block |
| A held finalizer is retried with a Ready message | A held finalizer is retried and the wait reason stays | the sentence of this block |
| The Ready message of a failed recheck is bounded | A status that cannot be written fails the reconcile | the sentence of this block |
| A failed recheck is reported as DeletionBlocked | Ten minutes of failed rechecks report DeletionBlocked | the sentence of this block |
| A ServiceAccount that loses its rights does not release the deletion by itself | Disabling prune releases a deletion whose ServiceAccount lost its rights | the sentence of this block |
| The wait reasons are set only during a deletion | A wait reason is not set on a live object | the sentence of this block |
| Without the record the stall reasons hold the finalizer | Without the record a missing identity holds with DeletionSAMissing | the sentence of this block |
| The DeletionBlocked message names the remaining objects and the ways out | The message names each remaining object and the ways out | the sentence of this block |
| A blocked deletion keeps the finalizer and keeps rechecking | A blocked deletion is not given up on a timeout | the sentence of this block |
| A blocked deletion keeps the finalizer and keeps rechecking | An object terminating before the deletion is reported by the first reconcile | the sentence of this block |

Clauses (second table):

| Old clause (requirement, sentence) | New requirement | Edit |
| --- | --- | --- |
| Finalizer registration #0: The controller MUST add the finalizer `releases.opmodel.dev/cleanup... | The cleanup finalizer is registered in Phase 0 | verbatim |
| Finalizer registration #1: The reconcile call that adds the finalizer MUST request an immediat... | The finalizer-add reconcile requests an immediate requeue | verbatim |
| Finalizer registration #2: Without an explicit requeue, the next reconcile would be deferred t... | The cleanup finalizer is registered in Phase 0 | verbatim |
| Deletion cleanup with prune enabled #0: When a ModuleRelease with `spec.prune=true` is deleted, the control... | Deletion cleanup with prune enabled | verbatim |
| Deletion cleanup with prune enabled #1: The finalizer is removed, and the deletion completes, once every ob... | The deletion completes once every deleted object is gone | "The finalizer is removed" becomes "the finalizer MUST be removed" |
| Finalizer retained on DeletionSAMissing stall #0: While a release is stalled with reason `DeletionSAMissing`, the fin... | Finalizer retained on DeletionSAMissing stall | verbatim |
| Finalizer retained on DeletionSAMissing stall #1: A deletion never enters this stall when its inventory holds nothing... | Finalizer retained on DeletionSAMissing stall | verbatim |
| Finalizer retained on DeletionSAMissing stall #2: The release remains blocked from garbage collection by the apiserve... | A DeletionSAMissing stall ends by a restored ServiceAccount or the orphan annotation | "The release remains" becomes "the release MUST remain" |
| Finalizer retained on DeletionSAMissing stall #2: The release remains blocked from garbage collection by the apiserve... | A DeletionSAMissing stall ends when spec.prune is set to false | "The release remains" becomes "the release MUST remain" |
| Finalizer retained on DeletionSAMissing stall #3: - The ServiceAccount is restored in the release's namespace and the... | A DeletionSAMissing stall ends by a restored ServiceAccount or the orphan annotation | verbatim |
| Finalizer retained on DeletionSAMissing stall #4: - The operator sets annotation `opm.dev/force-delete-orphan=true` o... | A DeletionSAMissing stall ends by a restored ServiceAccount or the orphan annotation | ", OR" becomes "." |
| Finalizer retained on DeletionSAMissing stall #5: - The operator sets `spec.prune=false` on the release and the next ... | A DeletionSAMissing stall ends when spec.prune is set to false | verbatim |
| Recovery signals of a stalled deletion trigger a reconcile #0: A `ModuleInstance` whose deletion is stalled with reason `DeletionS... | Recovery signals of a stalled deletion trigger a reconcile | verbatim |
| Recovery signals of a stalled deletion trigger a reconcile #1: - its annotation `opm.dev/force-delete-orphan` becomes the literal ... | Recovery signals of a stalled deletion trigger a reconcile | verbatim |
| Recovery signals of a stalled deletion trigger a reconcile #2: - a ServiceAccount with the name the instance impersonates is creat... | Recovery signals of a stalled deletion trigger a reconcile | verbatim |
| Recovery signals of a stalled deletion trigger a reconcile #3: The creation of the ServiceAccount is the only trigger of the secon... | Recovery signals of a stalled deletion trigger a reconcile | verbatim |
| Recovery signals of a stalled deletion trigger a reconcile #4: When the ServiceAccount is created before the RBAC that lets it del... | A ServiceAccount created before its RBAC stalls the deletion with ImpersonationFailed | "the deletion stalls with" becomes "the deletion MUST stall with" |
| Recovery signals of a stalled deletion trigger a reconcile #5: "Promptly" means that the event itself enqueues the instance. | A recovery signal enqueues the instance and the existing deletion rules apply | verbatim |
| Recovery signals of a stalled deletion trigger a reconcile #6: The reconcile that follows applies the existing deletion rules unch... | A recovery signal enqueues the instance and the existing deletion rules apply | "The reconcile that follows applies" becomes "The reconcile that follows MUST apply" |
| Recovery signals of a stalled deletion trigger a reconcile #7: An annotation change on an instance that is not being deleted, and ... | Other annotation changes and status-only writes trigger no reconcile | verbatim |
| Recovery signals of a stalled deletion trigger a reconcile #8: A status-only write MUST NOT trigger a reconcile, as before. | Other annotation changes and status-only writes trigger no reconcile | verbatim |
| Recovery signals of a stalled deletion trigger a reconcile #9: This requirement covers `ModuleInstance`. | Recovery signals of a stalled deletion trigger a reconcile | verbatim (repeated as context) |
| Recovery signals of a stalled deletion trigger a reconcile #9: This requirement covers `ModuleInstance`. | A ServiceAccount created before its RBAC stalls the deletion with ImpersonationFailed | verbatim (repeated as context) |
| Recovery signals of a stalled deletion trigger a reconcile #9: This requirement covers `ModuleInstance`. | Other annotation changes and status-only writes trigger no reconcile | verbatim (repeated as context) |
| Recovery signals of a stalled deletion trigger a reconcile #10: A `ModulePackage` stalled the same way still waits for its stalled ... | Other annotation changes and status-only writes trigger no reconcile | verbatim |
| A kept claim never holds the finalizer #0: A PersistentVolumeClaim that the deletion cleanup keeps MUST NOT ke... | A kept claim does not hold the finalizer | verbatim |
| A kept claim never holds the finalizer #1: The finalizer MUST be removed in the reconcile in which every other... | A kept claim does not hold the finalizer | verbatim |
| A kept claim never holds the finalizer #2: An inventory that holds only kept claims MUST release the finalizer... | An inventory of only kept claims releases the finalizer without a delete | verbatim |
| A kept claim never holds the finalizer #3: A failure to delete another entry MUST still keep the finalizer, as... | A kept claim does not hold the finalizer | verbatim |
| Deletion cleanup judges every object with the delete verdict #0: With `spec.prune` true, the deletion cleanup of a ModuleInstance or... | Deletion cleanup judges every object with the delete verdict | verbatim |
| Deletion cleanup judges every object with the delete verdict #1: It MUST judge with `status.instanceUUID`, and with `status.previous... | The deletion cleanup judges with the identities of the instance | "It MUST judge" becomes "The deletion cleanup MUST judge" |
| Deletion cleanup judges every object with the delete verdict #2: Source: enhancement 0012 ("The library owns the whole deletion sequ... | Deletion cleanup judges every object with the delete verdict | moved to a `**Source**:` line, full old Source line |
| Deletion cleanup judges every object with the delete verdict #2: Source: enhancement 0012 ("The library owns the whole deletion sequ... | The deletion cleanup judges with the identities of the instance | moved to a `**Source**:` line, full old Source line |
| Deletion cleanup judges every object with the delete verdict #2: Source: enhancement 0012 ("The library owns the whole deletion sequ... | An empty instance identity is judged with no identity | moved to a `**Source**:` line, full old Source line |
| Deletion cleanup judges every object with the delete verdict #2: Source: enhancement 0012 ("The library owns the whole deletion sequ... | A skipped object does not hold the finalizer and a failed step does | moved to a `**Source**:` line, full old Source line |
| Deletion cleanup judges every object with the delete verdict #3: With an empty `status.instanceUUID` the cleanup MUST ask the verdic... | An empty instance identity is judged with no identity | verbatim |
| Deletion cleanup judges every object with the delete verdict #4: The verdict then compares no UUID label, and it leaves in place eve... | An empty instance identity is judged with no identity | verbatim |
| Deletion cleanup judges every object with the delete verdict #5: An object the verdict skips MUST be left in the cluster and MUST NO... | A skipped object does not hold the finalizer and a failed step does | verbatim |
| Deletion cleanup judges every object with the delete verdict #6: In the scenarios below "the finalizer is removed" means: once every... | Deletion cleanup judges every object with the delete verdict | verbatim |
| Deletion cleanup judges every object with the delete verdict #7: A read that fails with an error other than NotFound, and a DELETE t... | A skipped object does not hold the finalizer and a failed step does | verbatim |
| Deletion cleanup judges every object with the delete verdict #8: The rules for PersistentVolumeClaims and for a missing ServiceAccou... | A skipped object does not hold the finalizer and a failed step does | verbatim |
| Deletion cleanup runs the library's deletion plan #0: With `spec.prune` true, the deletion cleanup of a ModuleInstance or... | The deletion cleanup deletes an object only as the library's deletion plan names it | verbatim |
| Deletion cleanup runs the library's deletion plan #1: Source: enhancement 0012 ("The library owns the whole deletion sequ... | The deletion cleanup deletes an object only as the library's deletion plan names it | moved to a `**Source**:` line, full old Source line |
| Deletion cleanup runs the library's deletion plan #1: Source: enhancement 0012 ("The library owns the whole deletion sequ... | The deletion cleanup runs one plan per identity | moved to a `**Source**:` line, full old Source line |
| Deletion cleanup runs the library's deletion plan #1: Source: enhancement 0012 ("The library owns the whole deletion sequ... | A kept claim is not an entry of any plan | moved to a `**Source**:` line, full old Source line |
| Deletion cleanup runs the library's deletion plan #2: The cleanup MUST send each delete in the order the plan gives, desc... | The deletion cleanup deletes an object only as the library's deletion plan names it | verbatim |
| Deletion cleanup runs the library's deletion plan #3: Each delete MUST carry the propagation policy the action names, For... | The deletion cleanup deletes an object only as the library's deletion plan names it | verbatim |
| Deletion cleanup runs the library's deletion plan #4: When the object records two identities, or none on a ModulePackage,... | The deletion cleanup runs one plan per identity | "When the" becomes "when the" |
| Deletion cleanup runs the library's deletion plan #5: An entry that no plan deletes MUST be reported with the reason of t... | The deletion cleanup runs one plan per identity | verbatim |
| Deletion cleanup runs the library's deletion plan #6: A PersistentVolumeClaim that `spec.dataPolicy` keeps MUST NOT be an... | A kept claim is not an entry of any plan | "A PersistentVolumeClaim" becomes "a PersistentVolumeClaim" |
| The cleanup finalizer follows the library's hold verdict #0: The controller MUST remove the cleanup finalizer of a deleting Modu... | The cleanup finalizer follows the library's hold verdict | verbatim |
| The cleanup finalizer follows the library's hold verdict #1: Source: enhancement 0012 ("The library owns the whole deletion sequ... | The cleanup finalizer follows the library's hold verdict | moved to a `**Source**:` line, full old Source line |
| The cleanup finalizer follows the library's hold verdict #1: Source: enhancement 0012 ("The library owns the whole deletion sequ... | The inputs of the hold verdict | moved to a `**Source**:` line, full old Source line |
| The cleanup finalizer follows the library's hold verdict #1: Source: enhancement 0012 ("The library owns the whole deletion sequ... | Each hold verdict keeps its status | moved to a `**Source**:` line, full old Source line |
| The cleanup finalizer follows the library's hold verdict #1: Source: enhancement 0012 ("The library owns the whole deletion sequ... | Only kept claims release the finalizer before the identity is looked at | moved to a `**Source**:` line, full old Source line |
| The cleanup finalizer follows the library's hold verdict #1: Source: enhancement 0012 ("The library owns the whole deletion sequ... | A release without a read emits DeletionUnconfirmed and no ClaimsKept | moved to a `**Source**:` line, full old Source line |
| The cleanup finalizer follows the library's hold verdict #2: The verdict's inputs MUST be: prune from `spec.prune`; force-orphan... | The inputs of the hold verdict | verbatim |
| The cleanup finalizer follows the library's hold verdict #3: The verdict MUST be asked with the plans as they are built, that is... | The cleanup finalizer follows the library's hold verdict | verbatim |
| The cleanup finalizer follows the library's hold verdict #4: Each verdict MUST keep the status it has today: a missing identity ... | Each hold verdict keeps its status | verbatim |
| The cleanup finalizer follows the library's hold verdict #5: One outcome differs from the behaviour before this requirement. | Only kept claims release the finalizer before the identity is looked at | verbatim |
| The cleanup finalizer follows the library's hold verdict #6: When the inventory holds only PersistentVolumeClaims that `spec.dat... | Only kept claims release the finalizer before the identity is looked at | verbatim |
| The cleanup finalizer follows the library's hold verdict #7: The controller MUST then remove the finalizer without a read, also ... | Only kept claims release the finalizer before the identity is looked at | verbatim |
| The cleanup finalizer follows the library's hold verdict #8: It MUST NOT emit a `ClaimsKept` event for claims it did not read; i... | A release without a read emits DeletionUnconfirmed and no ClaimsKept | "It MUST NOT" becomes "the controller MUST NOT" |
| The finalizer is kept until the deleted objects are gone #0: After a release verdict, the controller MUST remove the cleanup fin... | The finalizer is kept until the deleted objects are gone | verbatim |
| The finalizer is kept until the deleted objects are gone #1: A kept PersistentVolumeClaim, an object left behind and an object t... | The finalizer is kept until the deleted objects are gone | verbatim |
| The finalizer is kept until the deleted objects are gone #2: Source: owner decision of 2026-10-09 (the finalizer waits until the... | The finalizer is kept until the deleted objects are gone | moved to a `**Source**:` line, full old Source line |
| The finalizer is kept until the deleted objects are gone #2: Source: owner decision of 2026-10-09 (the finalizer waits until the... | A deleted object that still exists holds the finalizer with DeletionInProgress | moved to a `**Source**:` line, full old Source line |
| The finalizer is kept until the deleted objects are gone #2: Source: owner decision of 2026-10-09 (the finalizer waits until the... | Each reconcile of a waiting deletion judges the inventory again from the cluster | moved to a `**Source**:` line, full old Source line |
| The finalizer is kept until the deleted objects are gone #3: While at least one deleted object still exists, the reconcile MUST ... | A deleted object that still exists holds the finalizer with DeletionInProgress | verbatim |
| The finalizer is kept until the deleted objects are gone #4: It MUST NOT wait inside the reconcile. | A deleted object that still exists holds the finalizer with DeletionInProgress | verbatim |
| The finalizer is kept until the deleted objects are gone #5: Each such reconcile MUST judge every inventory entry again from the... | Each reconcile of a waiting deletion judges the inventory again from the cluster | verbatim |
| The finalizer is kept until the deleted objects are gone #6: Setting `spec.prune` to false on the deleting object MUST release t... | The finalizer is kept until the deleted objects are gone | verbatim |
| A wait survives the loss of the ServiceAccount #0: A deletion cleanup has sent every delete when one reconcile of it e... | A wait survives the loss of the ServiceAccount | verbatim |
| A wait survives the loss of the ServiceAccount #1: That reconcile records the fact by setting the `Ready` reason `Dele... | A wait survives the loss of the ServiceAccount | verbatim |
| A wait survives the loss of the ServiceAccount #2: The reason `DeletionBlocked` carries the same fact. | The record is read only from the Ready condition at the start of the reconcile | verbatim |
| A wait survives the loss of the ServiceAccount #3: No other field records it. | The record is read only from the Ready condition at the start of the reconcile | verbatim |
| A wait survives the loss of the ServiceAccount #4: Source: owner decisions of 2026-10-09 (the `Ready` reason is the re... | A wait survives the loss of the ServiceAccount | moved to a `**Source**:` line, full old Source line |
| A wait survives the loss of the ServiceAccount #4: Source: owner decisions of 2026-10-09 (the `Ready` reason is the re... | The record is read only from the Ready condition at the start of the reconcile | moved to a `**Source**:` line, full old Source line |
| A wait survives the loss of the ServiceAccount #4: Source: owner decisions of 2026-10-09 (the `Ready` reason is the re... | The ServiceAccount is gone in exactly one case | moved to a `**Source**:` line, full old Source line |
| A wait survives the loss of the ServiceAccount #4: Source: owner decisions of 2026-10-09 (the `Ready` reason is the re... | A NotFound that names another object does not make the ServiceAccount gone | moved to a `**Source**:` line, full old Source line |
| A wait survives the loss of the ServiceAccount #4: Source: owner decisions of 2026-10-09 (the `Ready` reason is the re... | A gone ServiceAccount releases the finalizer | moved to a `**Source**:` line, full old Source line |
| A wait survives the loss of the ServiceAccount #4: Source: owner decisions of 2026-10-09 (the `Ready` reason is the re... | Other answers hold the finalizer | moved to a `**Source**:` line, full old Source line |
| A wait survives the loss of the ServiceAccount #4: Source: owner decisions of 2026-10-09 (the `Ready` reason is the re... | A held finalizer is retried with a Ready message | moved to a `**Source**:` line, full old Source line |
| A wait survives the loss of the ServiceAccount #4: Source: owner decisions of 2026-10-09 (the `Ready` reason is the re... | The Ready message of a failed recheck is bounded | moved to a `**Source**:` line, full old Source line |
| A wait survives the loss of the ServiceAccount #4: Source: owner decisions of 2026-10-09 (the `Ready` reason is the re... | A failed recheck is reported as DeletionBlocked | moved to a `**Source**:` line, full old Source line |
| A wait survives the loss of the ServiceAccount #4: Source: owner decisions of 2026-10-09 (the `Ready` reason is the re... | A ServiceAccount that loses its rights does not release the deletion by itself | moved to a `**Source**:` line, full old Source line |
| A wait survives the loss of the ServiceAccount #4: Source: owner decisions of 2026-10-09 (the `Ready` reason is the re... | The wait reasons are set only during a deletion | moved to a `**Source**:` line, full old Source line |
| A wait survives the loss of the ServiceAccount #4: Source: owner decisions of 2026-10-09 (the `Ready` reason is the re... | Without the record the stall reasons hold the finalizer | moved to a `**Source**:` line, full old Source line |
| A wait survives the loss of the ServiceAccount #5: The controller MUST read the record only from the `Ready` condition... | The record is read only from the Ready condition at the start of the reconcile | verbatim |
| A wait survives the loss of the ServiceAccount #6: It MUST use the record for one decision: to remove the finalizer wh... | The record is read only from the Ready condition at the start of the reconcile | verbatim |
| A wait survives the loss of the ServiceAccount #7: While the ServiceAccount exists and every request succeeds, the rec... | A wait survives the loss of the ServiceAccount | verbatim |
| A wait survives the loss of the ServiceAccount #8: The ServiceAccount is gone in exactly one case: the controller's ow... | The ServiceAccount is gone in exactly one case | verbatim |
| A wait survives the loss of the ServiceAccount #9: The controller MUST decide this from the typed API status (a NotFou... | The ServiceAccount is gone in exactly one case | verbatim (repeated as context) |
| A wait survives the loss of the ServiceAccount #9: The controller MUST decide this from the typed API status (a NotFou... | A NotFound that names another object does not make the ServiceAccount gone | verbatim (repeated as context) |
| A wait survives the loss of the ServiceAccount #10: A NotFound that names another object or none is not this case. | A NotFound that names another object does not make the ServiceAccount gone | verbatim |
| A wait survives the loss of the ServiceAccount #11: A ServiceAccount that was deleted and created again under the same ... | A NotFound that names another object does not make the ServiceAccount gone | verbatim |
| A wait survives the loss of the ServiceAccount #12: With the record present: | A gone ServiceAccount releases the finalizer | verbatim (repeated as context) |
| A wait survives the loss of the ServiceAccount #12: With the record present: | Other answers hold the finalizer | verbatim (repeated as context) |
| A wait survives the loss of the ServiceAccount #12: With the record present: | A failed recheck is reported as DeletionBlocked | verbatim (repeated as context) |
| A wait survives the loss of the ServiceAccount #13: - When the ServiceAccount is gone, the controller MUST remove the f... | A gone ServiceAccount releases the finalizer | verbatim |
| A wait survives the loss of the ServiceAccount #14: - Every other answer MUST hold the finalizer, on any read or delete... | Other answers hold the finalizer | verbatim |
| A wait survives the loss of the ServiceAccount #15: This includes a Forbidden answer for one object while other objects... | A held finalizer is retried with a Ready message | verbatim |
| A wait survives the loss of the ServiceAccount #16: The controller MUST retry with its normal backoff, MUST NOT replace... | A held finalizer is retried with a Ready message | verbatim |
| A wait survives the loss of the ServiceAccount #17: The message MUST name at most three failed requests, MUST count the... | The Ready message of a failed recheck is bounded | verbatim |
| A wait survives the loss of the ServiceAccount #18: When the status cannot be written, the reconcile MUST fail so that ... | The Ready message of a failed recheck is bounded | verbatim |
| A wait survives the loss of the ServiceAccount #19: - Once the object has been deleting for 10 minutes by the controlle... | A failed recheck is reported as DeletionBlocked | verbatim |
| A wait survives the loss of the ServiceAccount #20: Known limit, by the owner's decision: when the rights of the Servic... | A ServiceAccount that loses its rights does not release the deletion by itself | "the deletion does not release by itself" becomes "the deletion MUST NOT release by itself" |
| A wait survives the loss of the ServiceAccount #21: It reports `DeletionBlocked` after 10 minutes and names the cause; | A ServiceAccount that loses its rights does not release the deletion by itself | verbatim |
| A wait survives the loss of the ServiceAccount #22: `spec.prune=false` releases it, and so does deleting the ServiceAcc... | A ServiceAccount that loses its rights does not release the deletion by itself | verbatim |
| A wait survives the loss of the ServiceAccount #23: Without the record, a missing or failed identity and a Forbidden st... | Without the record the stall reasons hold the finalizer | verbatim |
| A wait survives the loss of the ServiceAccount #24: The controller MUST NOT set the reasons `DeletionInProgress` and `D... | The wait reasons are set only during a deletion | verbatim |
| A wait survives the loss of the ServiceAccount #25: The roles the operator ships MUST NOT grant a user a write verb on ... | The wait reasons are set only during a deletion | verbatim |
| A deletion that does not finish is reported as blocked #0: When a deleted object still exists 10 minutes after its `deletionTi... | A deletion that does not finish is reported as blocked | verbatim |
| A deletion that does not finish is reported as blocked #1: The message MUST name each remaining object by kind, namespace and ... | The DeletionBlocked message names the remaining objects and the ways out | verbatim |
| A deletion that does not finish is reported as blocked #2: The controller MUST keep the finalizer, MUST keep rechecking at an ... | A blocked deletion keeps the finalizer and keeps rechecking | "The controller MUST keep" becomes "the controller MUST keep" |
| A deletion that does not finish is reported as blocked #3: An object that had been terminating for more than 10 minutes before... | A blocked deletion keeps the finalizer and keeps rechecking | verbatim |
| A deletion that does not finish is reported as blocked #4: The threshold and the clock MUST be values a test can set, so that ... | A deletion that does not finish is reported as blocked | verbatim |
| A deletion that does not finish is reported as blocked #5: They MUST NOT be command-line flags. | A deletion that does not finish is reported as blocked | verbatim |

The em-dash of sentence 1 of "Finalizer registration" is written as a colon (workspace style). Every `**Source**:` line holds the full old Source line and nothing else.

Added words (the only ones): subjects and conditions copied from the same old requirement, the exact header of a sibling ("As in requirement" prefix or "(Requirement: ...)" suffix), and the keywords in the Strength table.


### Ledger: instance-health

Forms and lengths (old body, then each new body):

| Old requirement | Form | Old length | New blocks (length) |
| --- | --- | --- | --- |
| A Healthy condition reports whether the applied objects have rolled out | c2 (REMOVED plus ADDED under new names) | 2086 | Healthy reports whether the inventory rolled out (375); Healthy is RolledOut or ProgressDeadlineExceeded (362); Healthy is NotRolledOut or HealthUnknown (394); The Healthy message begins with the ready count and names the objects that are not healthy (392); The Healthy message orders the named objects and is deterministic (355); The Healthy message names the error of an unreadable object (351); Ready keeps its meaning beside Healthy (243) |
| Objects are read uncached, after the apply, through the applying identity | c2 (REMOVED plus ADDED under new names) | 916 | Each inventory object is read uncached through the applying identity (290); Health is read only after the apply (233); Health reads are bounded (202); An unbuildable impersonated reader makes Healthy Unknown (103) |
| A reconcile requeues until the instance has rolled out | c2 (REMOVED plus ADDED under new names) | 895 | A reconcile requeues after half the time since the last apply (387); A stalled rollout requeues at the stalled recheck interval (168); A rolled-out instance is not requeued for health (247); A package requeues after the shorter of the health requeue and its interval (173) |
| An expired Job does not make a ModuleInstance unhealthy | c2 (REMOVED plus ADDED under new names) | 1099 | An expired Job leaves the inventory and is not judged (475); A render-skipping reconcile does not report an absent Job as Missing (226); A failed Job with a TTL counts while it exists (354); A removal that empties the inventory leaves Healthy Unknown (197) |

`grep -rn --exclude-dir=archive "<old header>" .` found no reference to a header that this unit renames or removes outside this spec; the in-spec references name headers that stay (form b parents) or are updated here.

Scenarios:

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| Every object healthy | Healthy is RolledOut or ProgressDeadlineExceeded, Every object healthy | no (byte equal, checked by the item 3 command) |
| A rollout in progress | The Healthy message begins with the ready count and names the objects that are not healthy, A rollout in progress | no (byte equal, checked by the item 3 command) |
| A stalled Deployment | The Healthy message orders the named objects and is deterministic, A stalled Deployment | no (byte equal, checked by the item 3 command) |
| A missing object | The Healthy message begins with the ready count and names the objects that are not healthy, A missing object | no (byte equal, checked by the item 3 command) |
| A kind that is no longer served | The Healthy message names the error of an unreadable object, A kind that is no longer served | no (byte equal, checked by the item 3 command) |
| An unreadable object | The Healthy message names the error of an unreadable object, An unreadable object | no (byte equal, checked by the item 3 command) |
| An empty inventory | Healthy is NotRolledOut or HealthUnknown, An empty inventory | no (byte equal, checked by the item 3 command) |
| Ready and its readers are unchanged | Ready keeps its meaning beside Healthy, Ready and its readers are unchanged | no (byte equal, checked by the item 3 command) |
| An impersonated instance reads as its ServiceAccount | Each inventory object is read uncached through the applying identity, An impersonated instance reads as its ServiceAccount | no (byte equal, checked by the item 3 command) |
| An instance without a ServiceAccount reads uncached | Each inventory object is read uncached through the applying identity, An instance without a ServiceAccount reads uncached | no (byte equal, checked by the item 3 command) |
| A just-updated Deployment is not reported rolled out | Health is read only after the apply, A just-updated Deployment is not reported rolled out | no (byte equal, checked by the item 3 command) |
| A missing ServiceAccount on a skip | An unbuildable impersonated reader makes Healthy Unknown, A missing ServiceAccount on a skip | no (byte equal, checked by the item 3 command) |
| A fresh rollout is checked quickly | A reconcile requeues after half the time since the last apply, A fresh rollout is checked quickly | no (byte equal, checked by the item 3 command) |
| A long rollout settles at the ceiling | A reconcile requeues after half the time since the last apply, A long rollout settles at the ceiling | no (byte equal, checked by the item 3 command) |
| A stalled rollout stops the fast requeue | A stalled rollout requeues at the stalled recheck interval, A stalled rollout stops the fast requeue | no (byte equal, checked by the item 3 command) |
| A rolled-out instance is not requeued for health | A rolled-out instance is not requeued for health, A rolled-out instance is not requeued for health | no (byte equal, checked by the item 3 command) |
| A package keeps its interval as an upper bound | A package requeues after the shorter of the health requeue and its interval, A package keeps its interval as an upper bound | no (byte equal, checked by the item 3 command) |
| An instance with an expired Job stays healthy | An expired Job leaves the inventory and is not judged, An instance with an expired Job stays healthy | no (byte equal, checked by the item 3 command) |
| The only inventory object was an expired Job | A removal that empties the inventory leaves Healthy Unknown, The only inventory object was an expired Job | no (byte equal, checked by the item 3 command) |
| An absent Job after a failed drift detection | A render-skipping reconcile does not report an absent Job as Missing, An absent Job after a failed drift detection | no (byte equal, checked by the item 3 command) |

New scenarios (each restates a case the old text names):

| New place | New scenario | Old text that names the case |
| --- | --- | --- |
| Healthy reports whether the inventory rolled out | More than one reason applies | the sentence of this block |
| Health reads are bounded | A read that does not finish in time counts as unreadable | the sentence of this block |
| A failed Job with a TTL counts while it exists | A failed Job with a TTL counts while it exists | the sentence of this block |

Clauses (second table):

| Old clause (requirement, sentence) | New requirement | Edit |
| --- | --- | --- |
| A Healthy condition reports whether the applied objects have rolled out #0: The operator SHALL write a `Healthy` condition on a ModuleInstance ... | Healthy reports whether the inventory rolled out | verbatim |
| A Healthy condition reports whether the applied objects have rolled out #1: It reports whether the objects in `status.inventory` have rolled ou... | Healthy reports whether the inventory rolled out | verbatim |
| A Healthy condition reports whether the applied objects have rolled out #2: It SHALL be one of: | Healthy is RolledOut or ProgressDeadlineExceeded | "It SHALL be one of:" becomes "The `Healthy` condition SHALL be one of:" |
| A Healthy condition reports whether the applied objects have rolled out #2: It SHALL be one of: | Healthy is NotRolledOut or HealthUnknown | "It SHALL be one of:" becomes "The `Healthy` condition SHALL be one of:" |
| A Healthy condition reports whether the applied objects have rolled out #3: - `True` with reason `RolledOut`: every inventory object was read a... | Healthy is RolledOut or ProgressDeadlineExceeded | verbatim |
| A Healthy condition reports whether the applied objects have rolled out #4: - `False` with reason `ProgressDeadlineExceeded`: at least one obje... | Healthy is RolledOut or ProgressDeadlineExceeded | "holds;" becomes "holds." |
| A Healthy condition reports whether the applied objects have rolled out #5: - `False` with reason `NotRolledOut`: no Deployment is stalled, and... | Healthy is NotRolledOut or HealthUnknown | verbatim |
| A Healthy condition reports whether the applied objects have rolled out #6: - `Unknown` with reason `HealthUnknown`: no object is known to be u... | Healthy is NotRolledOut or HealthUnknown | verbatim |
| A Healthy condition reports whether the applied objects have rolled out #7: When more than one applies, `ProgressDeadlineExceeded` wins over `N... | Healthy reports whether the inventory rolled out | verbatim |
| A Healthy condition reports whether the applied objects have rolled out #8: `RolledOut` holds only when none of the others applies. | Healthy reports whether the inventory rolled out | verbatim |
| A Healthy condition reports whether the applied objects have rolled out #9: The message SHALL begin with the ready count and the total `health.... | The Healthy message begins with the ready count and names the objects that are not healthy | verbatim |
| A Healthy condition reports whether the applied objects have rolled out #10: For `ProgressDeadlineExceeded` and `NotRolledOut` it SHALL name up ... | The Healthy message begins with the ready count and names the objects that are not healthy | verbatim |
| A Healthy condition reports whether the applied objects have rolled out #11: For `ProgressDeadlineExceeded` the stalled Deployments come first, ... | The Healthy message orders the named objects and is deterministic | verbatim |
| A Healthy condition reports whether the applied objects have rolled out #12: For a read failure it SHALL name the error of the first unreadable ... | The Healthy message names the error of an unreadable object | verbatim |
| A Healthy condition reports whether the applied objects have rolled out #13: A read refused because the object's kind is no longer served (the A... | The Healthy message names the error of an unreadable object | verbatim |
| A Healthy condition reports whether the applied objects have rolled out #14: Two judgements of the same cluster state SHALL produce the same mes... | The Healthy message orders the named objects and is deterministic | verbatim |
| A Healthy condition reports whether the applied objects have rolled out #15: `Ready` SHALL keep its meaning: the render was applied. | Ready keeps its meaning beside Healthy | verbatim |
| A Healthy condition reports whether the applied objects have rolled out #16: No reason, status or timing of `Ready` depends on `Healthy`. | Ready keeps its meaning beside Healthy | verbatim |
| A Healthy condition reports whether the applied objects have rolled out #17: `spec.dependsOn` SHALL keep waiting on `Ready`, and a TransformerRe... | Ready keeps its meaning beside Healthy | verbatim |
| Objects are read uncached, after the apply, through the applying identity #0: The reconciler SHALL read each inventory object uncached. | Each inventory object is read uncached through the applying identity | verbatim (repeated as context) |
| Objects are read uncached, after the apply, through the applying identity #0: The reconciler SHALL read each inventory object uncached. | Health is read only after the apply | verbatim (repeated as context) |
| Objects are read uncached, after the apply, through the applying identity #0: The reconciler SHALL read each inventory object uncached. | Health reads are bounded | verbatim (repeated as context) |
| Objects are read uncached, after the apply, through the applying identity #1: When an effective ServiceAccount applies the instance, it SHALL rea... | Each inventory object is read uncached through the applying identity | verbatim |
| Objects are read uncached, after the apply, through the applying identity #2: Otherwise it SHALL read through the manager's uncached API reader, ... | Each inventory object is read uncached through the applying identity | verbatim |
| Objects are read uncached, after the apply, through the applying identity #3: On a successful apply it SHALL read only after the apply and prune ... | Health is read only after the apply | verbatim |
| Objects are read uncached, after the apply, through the applying identity #4: The reconciler SHALL NOT set `Healthy=True` from a read taken befor... | Health is read only after the apply | verbatim |
| Objects are read uncached, after the apply, through the applying identity #5: Reads SHALL run with at most 8 in flight and SHALL share one deadli... | Health reads are bounded | verbatim |
| Objects are read uncached, after the apply, through the applying identity #6: A read that does not finish within it counts as unreadable. | Health reads are bounded | verbatim |
| Objects are read uncached, after the apply, through the applying identity #7: When the impersonated reader cannot be built, `Healthy` SHALL be `U... | An unbuildable impersonated reader makes Healthy Unknown | verbatim |
| Objects are read uncached, after the apply, through the applying identity #8: Source: enhancement 0012 ("The Kubernetes runtime decisions live in... | Each inventory object is read uncached through the applying identity | moved to a `**Source**:` line, full old Source line |
| Objects are read uncached, after the apply, through the applying identity #8: Source: enhancement 0012 ("The Kubernetes runtime decisions live in... | Health is read only after the apply | moved to a `**Source**:` line, full old Source line |
| Objects are read uncached, after the apply, through the applying identity #8: Source: enhancement 0012 ("The Kubernetes runtime decisions live in... | Health reads are bounded | moved to a `**Source**:` line, full old Source line |
| Objects are read uncached, after the apply, through the applying identity #8: Source: enhancement 0012 ("The Kubernetes runtime decisions live in... | An unbuildable impersonated reader makes Healthy Unknown | moved to a `**Source**:` line, full old Source line |
| A reconcile requeues until the instance has rolled out #0: When `Healthy` is `False` with reason `NotRolledOut`, or `Unknown` ... | A reconcile requeues after half the time since the last apply | verbatim |
| A reconcile requeues until the instance has rolled out #1: When `Healthy` is `False` with reason `ProgressDeadlineExceeded`, t... | A stalled rollout requeues at the stalled recheck interval | verbatim |
| A reconcile requeues until the instance has rolled out #2: When `Healthy` is `True`, or `Unknown` because the inventory is emp... | A rolled-out instance is not requeued for health | verbatim |
| A reconcile requeues until the instance has rolled out #3: A ModulePackage SHALL requeue after the shorter of the health reque... | A package requeues after the shorter of the health requeue and its interval | verbatim |
| A reconcile requeues until the instance has rolled out #4: A health requeue SHALL NOT set `status.nextRetryAt` and SHALL NOT c... | A reconcile requeues after half the time since the last apply | verbatim |
| An expired Job does not make a ModuleInstance unhealthy #0: A `batch/v1` Job that a ModuleInstance renders with `ttlSecondsAfte... | An expired Job leaves the inventory and is not judged | verbatim |
| An expired Job does not make a ModuleInstance unhealthy #1: The reconcile that finds it absent with unchanged digests removes i... | An expired Job leaves the inventory and is not judged | verbatim |
| An expired Job does not make a ModuleInstance unhealthy #2: The ready count and the total in the message SHALL count those entr... | An expired Job leaves the inventory and is not judged | verbatim |
| An expired Job does not make a ModuleInstance unhealthy #3: A reconcile that would skip its render SHALL NOT report an absent J... | A render-skipping reconcile does not report an absent Job as Missing | verbatim |
| An expired Job does not make a ModuleInstance unhealthy #4: The operator does not know how the Job ended. | A failed Job with a TTL counts while it exists | verbatim |
| An expired Job does not make a ModuleInstance unhealthy #5: A Job with a TTL that failed makes `Healthy` `False` while it exist... | A failed Job with a TTL counts while it exists | verbatim |
| An expired Job does not make a ModuleInstance unhealthy #6: When the removal leaves `status.inventory` with no entries, `Health... | A removal that empties the inventory leaves Healthy Unknown | verbatim |
| An expired Job does not make a ModuleInstance unhealthy #7: A Job without a TTL that does not exist, and every other missing ob... | A failed Job with a TTL counts while it exists | verbatim |
| An expired Job does not make a ModuleInstance unhealthy #8: This requirement does not apply to a ModulePackage. | An expired Job leaves the inventory and is not judged | verbatim (repeated as context) |
| An expired Job does not make a ModuleInstance unhealthy #8: This requirement does not apply to a ModulePackage. | A render-skipping reconcile does not report an absent Job as Missing | verbatim (repeated as context) |
| An expired Job does not make a ModuleInstance unhealthy #8: This requirement does not apply to a ModulePackage. | A failed Job with a TTL counts while it exists | verbatim (repeated as context) |
| An expired Job does not make a ModuleInstance unhealthy #8: This requirement does not apply to a ModulePackage. | A removal that empties the inventory leaves Healthy Unknown | verbatim (repeated as context) |

Added words (the only ones): subjects and conditions copied from the same old requirement, the exact header of a sibling ("As in requirement" prefix or "(Requirement: ...)" suffix), and the keywords in the Strength table.

This section fills in with the ledgers as the groups run.

## Strength table

Every sentence that gains SHALL, MUST or MAY, collected from the clause tables of all units, with
the old sentence.

| Requirement (new) | Sentence | Old sentence, no keyword |
| --- | --- | --- |
| The deletion completes once every deleted object is gone | the finalizer MUST be removed, and the deletion completes, once every object the cleanup deleted is gone (requirement "The finalizer is kept until the deleted objects are gone"); that can be a late... | The finalizer is removed, and the deletion completes, once every object the cleanup deleted is gone (requirement "The finalizer is kept until the deleted objects are gone"); that can be a later rec... |
| A DeletionSAMissing stall ends by a restored ServiceAccount or the orphan annotation | the release MUST remain blocked from garbage collection by the apiserver until either: | The release remains blocked from garbage collection by the apiserver until either: |
| A DeletionSAMissing stall ends when spec.prune is set to false | the release MUST remain blocked from garbage collection by the apiserver until either: | The release remains blocked from garbage collection by the apiserver until either: |
| A recovery signal enqueues the instance and the existing deletion rules apply | The reconcile that follows MUST apply the existing deletion rules unchanged: the orphan-exit path for the annotation, the prune as the ServiceAccount for its return. | The reconcile that follows applies the existing deletion rules unchanged: the orphan-exit path for the annotation, the prune as the ServiceAccount for its return. |
| A ServiceAccount created before its RBAC stalls the deletion with ImpersonationFailed | When the ServiceAccount is created before the RBAC that lets it delete the inventory, the reconcile it triggers is forbidden and the deletion MUST stall with `ImpersonationFailed` until the stalled... | When the ServiceAccount is created before the RBAC that lets it delete the inventory, the reconcile it triggers is forbidden and the deletion stalls with `ImpersonationFailed` until the stalled rec... |
| A ServiceAccount that loses its rights does not release the deletion by itself | Known limit, by the owner's decision: when the rights of the ServiceAccount are removed during the wait and the ServiceAccount stays (its RoleBinding is deleted), the deletion MUST NOT release by i... | Known limit, by the owner's decision: when the rights of the ServiceAccount are removed during the wait and the ServiceAccount stays (its RoleBinding is deleted), the deletion does not release by i... |

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
