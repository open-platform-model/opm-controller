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
| A wait survives the loss of the ServiceAccount #20: Known limit, by the owner's decision: when the rights of the Servic... | A ServiceAccount that loses its rights does not release the deletion by itself | verbatim: "the deletion does not release by itself" |
| A wait survives the loss of the ServiceAccount #21: It reports `DeletionBlocked` after 10 minutes and names the cause; | A ServiceAccount that loses its rights does not release the deletion by itself | "It reports ... and names the cause" becomes "It SHALL report ... and SHALL name the cause" (Strength table) |
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

### Ledger: release-automation

Forms and lengths (old body, then each new body; a MODIFIED parent is measured on its body only):

| Old requirement | Form | Old length | New blocks (length) |
| --- | --- | --- | --- |
| Changelog generation | c2 (REMOVED plus ADDED under new names) | 686 | The changelog groups entries under the visible sections (380); Earlier changelog entries stay as they are (91); A module-only commit stays out of the root changelog (213) |
| Git tag and GitHub Release on merge | c2 (REMOVED plus ADDED under new names) | 724 | A merged Release PR creates the tag and a draft GitHub Release (374); Releases are drafts that a publish job makes public (349) |
| Version bump determination per release line | b (MODIFIED parent kept, ADDED splits) | 865 | Version bump determination per release line (parent, MODIFIED) (496); Commits set the proposed version (475) |
| Release PR opens for releasable commits on push to main | c2 (REMOVED plus ADDED under new names) | 772 | The release-please workflow opens or updates a Release PR on every push to main (407); Hidden commit types do not open a Release PR (465) |
| Release PR passes the release-pin gate | c2 (REMOVED plus ADDED under new names) | 1375 | The release-pin gate is a step of an existing job (435); A Release PR is told by its head branch (466); The gate fails on a replace directive or an unreleased pin (481); The gate fails on a dev pin or a tracked local-module file (479) |
| Dependabot leaves OPM Go modules to the release cascade | c2 (REMOVED plus ADDED under new names) | 759 | Dependabot ignores the OPM Go modules (274); Dependabot ignores the cascade references of the .github repository (334); Third-party updates continue under Dependabot (344) |
| The opm CLI version is pinned in one file | b (MODIFIED parent kept, ADDED splits) | 573 | The opm CLI version is pinned in one file (parent, MODIFIED) (440); Moving the opm CLI pin cuts no release (206) |
| Forced version via the release-as config key | b (MODIFIED parent kept, ADDED splits) | 1124 | Forced version via the release-as config key (parent, MODIFIED) (471); The release-as key repeats its version and crosses a boundary (457); A Release-As commit footer is not relied on (331) |
| The Release PR rewrites the install page's operator version | c2 (REMOVED plus ADDED under new names) | 705 | The install page's operator version sits in release-please blocks (482); The install page's version blocks are tested (309) |
| Release notifies downstream after it is published | c2 (REMOVED plus ADDED under new names) | 1580 | The notify-downstream job needs release-please and publish-release (403); The notify-downstream job has exactly one step, the cascade action (499); The notify-downstream job runs only after a published operator release (413); The cascade action dispatches upstream-released only to the cli (291); The cascade App key is read only by the notify job (481) |
| Dependabot leaves cuelang.org/go to library releases | b (MODIFIED parent kept, ADDED splits) | 820 | Dependabot leaves cuelang.org/go to library releases (parent, MODIFIED) (374); The operator's CUE version moves only through a library release (390) |

`grep -rn --exclude-dir=archive` for each old header found references to three headers outside their own requirement: "Version bump determination per release line" (a scenario of "Release PR opens for releasable commits on push to main" and a scenario of "Forced version via the release-as config key"), "Forced version via the release-as config key" (the body of "Beta prerelease line") and "Dependabot leaves cuelang.org/go to library releases" (the body of "Dependabot leaves OPM Go modules to the release cascade"). Those three headers stay as form b parents, so no reference changes. One comment of `.github/workflows/lint.yml` names the gate by G1, not by a header. The header "Release published once after every release job", cited in a new body, stays until group 9.

Scenarios:

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| Changelog includes all commit types | The changelog groups entries under the visible sections, Changelog includes all commit types | no (byte equal, checked by the item 3 command) |
| Changelog preserves history | Earlier changelog entries stay as they are, Changelog preserves history | no (byte equal, checked by the item 3 command) |
| Module-only commit stays out of the operator changelog | A module-only commit stays out of the root changelog, Module-only commit stays out of the operator changelog | no (byte equal, checked by the item 3 command) |
| Release PR merged | A merged Release PR creates the tag and a draft GitHub Release, Release PR merged | no (byte equal, checked by the item 3 command) |
| Release PR closed without merge | A merged Release PR creates the tag and a draft GitHub Release, Release PR closed without merge | no (byte equal, checked by the item 3 command) |
| Config enables draft-first | Releases are drafts that a publish job makes public, Config enables draft-first | no (byte equal, checked by the item 3 command) |
| fix commit on the beta line advances the counter | Version bump determination per release line (parent, MODIFIED), fix commit on the beta line advances the counter | no (byte equal, checked by the item 3 command) |
| feat commit on the beta line advances the counter | Version bump determination per release line (parent, MODIFIED), feat commit on the beta line advances the counter | no (byte equal, checked by the item 3 command) |
| Breaking-change commit on the beta line advances the counter | Version bump determination per release line (parent, MODIFIED), Breaking-change commit on the beta line advances the counter | no (byte equal, checked by the item 3 command) |
| feat commit on a stable line triggers MINOR bump | Version bump determination per release line (parent, MODIFIED), feat commit on a stable line triggers MINOR bump | no (byte equal, checked by the item 3 command) |
| fix-only commits on a stable line trigger PATCH bump | Version bump determination per release line (parent, MODIFIED), fix-only commits on a stable line trigger PATCH bump | no (byte equal, checked by the item 3 command) |
| Breaking-change commit on a stable line triggers MAJOR bump | Version bump determination per release line (parent, MODIFIED), Breaking-change commit on a stable line triggers MAJOR bump | no (byte equal, checked by the item 3 command) |
| (new) A breaking change is a ! in the PR title | Commits set the proposed version, A breaking change is a ! in the PR title | new: restates a case the old text names |
| (new) A docs commit proposes no bump | Commits set the proposed version, A docs commit proposes no bump | new: restates a case the old text names |
| First feat commit after a release | The release-please workflow opens or updates a Release PR on every push to main, First feat commit after a release | no (byte equal, checked by the item 3 command) |
| Subsequent fix commit with open Release PR | The release-please workflow opens or updates a Release PR on every push to main, Subsequent fix commit with open Release PR | no (byte equal, checked by the item 3 command) |
| Only non-releasable commits | Hidden commit types do not open a Release PR, Only non-releasable commits | no (byte equal, checked by the item 3 command) |
| Docs-only commits do not cut a release | Hidden commit types do not open a Release PR, Docs-only commits do not cut a release | no (byte equal, checked by the item 3 command) |
| Refactor commits still cut a release | The release-please workflow opens or updates a Release PR on every push to main, Refactor commits still cut a release | no (byte equal, checked by the item 3 command) |
| Replace directive blocks the release | The gate fails on a replace directive or an unreleased pin, Replace directive blocks the release | no (byte equal, checked by the item 3 command) |
| Pseudo-version blocks the release | The gate fails on a replace directive or an unreleased pin, Pseudo-version blocks the release | no (byte equal, checked by the item 3 command) |
| Untagged version blocks the release | The gate fails on a replace directive or an unreleased pin, Untagged version blocks the release | no (byte equal, checked by the item 3 command) |
| Dev CUE pin in a published fixture blocks the release | The gate fails on a dev pin or a tracked local-module file, Dev CUE pin in a published fixture blocks the release | no (byte equal, checked by the item 3 command) |
| Tracked local-module.cue blocks the release | The gate fails on a dev pin or a tracked local-module file, Tracked local-module.cue blocks the release | no (byte equal, checked by the item 3 command) |
| Clean release PR passes | The gate fails on a dev pin or a tracked local-module file, Clean release PR passes | no (byte equal, checked by the item 3 command) |
| The gate's check name is unique | The release-pin gate is a step of an existing job, The gate's check name is unique | no (byte equal, checked by the item 3 command) |
| Ordinary PRs skip the gate | A Release PR is told by its head branch, Ordinary PRs skip the gate | no (byte equal, checked by the item 3 command) |
| The module's release PR skips the gate | A Release PR is told by its head branch, The module's release PR skips the gate | no (byte equal, checked by the item 3 command) |
| Library release opens no Dependabot PR | Dependabot ignores the OPM Go modules, Library release opens no Dependabot PR | no (byte equal, checked by the item 3 command) |
| Third-party bumps continue | Third-party updates continue under Dependabot, Third-party bumps continue | no (byte equal, checked by the item 3 command) |
| A .github main commit opens no Dependabot PR | Dependabot ignores the cascade references of the .github repository, A .github main commit opens no Dependabot PR | no (byte equal, checked by the item 3 command) |
| Every install reads the file | The opm CLI version is pinned in one file (parent, MODIFIED), Every install reads the file | no (byte equal, checked by the item 3 command) |
| Malformed pin fails fast | The opm CLI version is pinned in one file (parent, MODIFIED), Malformed pin fails fast | no (byte equal, checked by the item 3 command) |
| Release job installs the pin of the release tag | The opm CLI version is pinned in one file (parent, MODIFIED), Release job installs the pin of the release tag | no (byte equal, checked by the item 3 command) |
| (new) Moving the pin is a ci(deps) commit | Moving the opm CLI pin cuts no release, Moving the pin is a ci(deps) commit | new: restates a case the old text names |
| Forced version set by a normal PR | Forced version via the release-as config key (parent, MODIFIED), Forced version set by a normal PR | no (byte equal, checked by the item 3 command) |
| The next PR removes the value once the release is cut | Forced version via the release-as config key (parent, MODIFIED), The next PR removes the value once the release is cut | no (byte equal, checked by the item 3 command) |
| Config carries release-as only while a forced release is pending | Forced version via the release-as config key (parent, MODIFIED), Config carries release-as only while a forced release is pending | no (byte equal, checked by the item 3 command) |
| A footer forces nothing | Forced version via the release-as config key (parent, MODIFIED), A footer forces nothing | no (byte equal, checked by the item 3 command) |
| A hidden-type carrier waits for a visible commit | Forced version via the release-as config key (parent, MODIFIED), A hidden-type carrier waits for a visible commit | no (byte equal, checked by the item 3 command) |
| (new) A value left in place re-proposes its version | The release-as key repeats its version and crosses a boundary, A value left in place re-proposes its version | new: restates a case the old text names |
| (new) A footer in the PR body never reaches main | A Release-As commit footer is not relied on, A footer in the PR body never reaches main | new: restates a case the old text names |
| A Release PR moves the page with the constant | The install page's operator version sits in release-please blocks, A Release PR moves the page with the constant | no (byte equal, checked by the item 3 command) |
| Another project's version inside a block is refused | The install page's version blocks are tested, Another project's version inside a block is refused | no (byte equal, checked by the item 3 command) |
| A detached marker is refused | The install page's version blocks are tested, A detached marker is refused | no (byte equal, checked by the item 3 command) |
| Published release notifies the cli | The cascade action dispatches upstream-released only to the cli, Published release notifies the cli | no (byte equal, checked by the item 3 command) |
| Release left a draft does not notify | The notify-downstream job runs only after a published operator release, Release left a draft does not notify | no (byte equal, checked by the item 3 command) |
| Failed notify is recovered without republishing | The notify-downstream job runs only after a published operator release, Failed notify is recovered without republishing | no (byte equal, checked by the item 3 command) |
| Notify switched off | The notify-downstream job runs only after a published operator release, Notify switched off | no (byte equal, checked by the item 3 command) |
| Docs failure does not hold notify | The notify-downstream job needs release-please and publish-release, Docs failure does not hold notify | no (byte equal, checked by the item 3 command) |
| Run without a release | The notify-downstream job runs only after a published operator release, Run without a release | no (byte equal, checked by the item 3 command) |
| A step added beside the action is refused | The notify-downstream job has exactly one step, the cascade action, A step added beside the action is refused | no (byte equal, checked by the item 3 command) |
| Module-only release does not notify | The notify-downstream job runs only after a published operator release, Module-only release does not notify | no (byte equal, checked by the item 3 command) |
| (new) The private key is read only by the notify job | The cascade App key is read only by the notify job, The private key is read only by the notify job | new: restates a case the old text names |
| New CUE release opens no Dependabot PR | Dependabot leaves cuelang.org/go to library releases (parent, MODIFIED), New CUE release opens no Dependabot PR | no (byte equal, checked by the item 3 command) |
| CUE arrives with a library release | Dependabot leaves cuelang.org/go to library releases (parent, MODIFIED), CUE arrives with a library release | no (byte equal, checked by the item 3 command) |
| Other third-party modules still update | Dependabot leaves cuelang.org/go to library releases (parent, MODIFIED), Other third-party modules still update | no (byte equal, checked by the item 3 command) |
| (new) The operator picks up CUE with the library pin | The operator's CUE version moves only through a library release, The operator picks up CUE with the library pin | new: restates a case the old text names |

Clauses (the rule sentences of each old body and the block that carries each):

| Old clause | New requirement | Note |
| --- | --- | --- |
| Changelog generation #1: The workflow SHALL generate and maintain a `CHANGELOG.md` file at the repository root for the operator pack... | The changelog groups entries under the visible sections | verbatim |
| Changelog generation #2: Entries SHALL be grouped under the visible sections that `release-please-config.json` declares (Features, B... | The changelog groups entries under the visible sections | verbatim |
| Changelog generation #3: Commits of a hidden type (`docs`, `chore`, `test`, `ci`, `build`) SHALL NOT appear. | The changelog groups entries under the visible sections | verbatim |
| Changelog generation #4: A commit that changes only files under the operator module's directory `modules/opm_operator/` belongs to t... | A module-only commit stays out of the root changelog | verbatim |
| Changelog generation #5: Documentation entries already in CHANGELOG.md from earlier releases SHALL stay as they are. | Earlier changelog entries stay as they are | verbatim |
| Git tag and GitHub Release on merge #1: When a Release PR is merged to `main`, release-please SHALL create the git tag eagerly at the release commi... | A merged Release PR creates the tag and a draft GitHub Release | verbatim |
| Git tag and GitHub Release on merge #2: The operator package `"."` tags `vX.Y.Z` (e.g. `v0.2.0`, or `v1.0.0-beta.6` on the beta line); the operator... | A merged Release PR creates the tag and a draft GitHub Release | verbatim |
| Git tag and GitHub Release on merge #3: Each release-please package SHALL be configured with `draft: true` and `force-tag-creation: true`. | Releases are drafts that a publish job makes public | verbatim |
| Git tag and GitHub Release on merge #4: An operator release's draft SHALL become public only through the publish job `publish-release` defined in "... | Releases are drafts that a publish job makes public | verbatim |
| Version bump determination per release line #1: The workflow SHALL determine the proposed version from the commits since the last release tag using Convent... | Commits set the proposed version | verbatim |
| Version bump determination per release line #2: On the prerelease line (`versioning: prerelease`, `prerelease: true`, current version `X.0.0-beta.N`), ever... | Version bump determination per release line (parent, MODIFIED) | verbatim |
| Version bump determination per release line #3: On a stable line (after GA, no prerelease suffix), `feat` SHALL propose a MINOR bump, a breaking change a M... | Version bump determination per release line (parent, MODIFIED) | verbatim |
| Version bump determination per release line #4: A breaking change is a `!` after the type or scope in the commit subject, which for a squash-merged PR is t... | Commits set the proposed version | verbatim |
| Version bump determination per release line #5: A `docs` commit is not releasable and proposes no bump. | Commits set the proposed version | verbatim |
| Release PR opens for releasable commits on push to main #1: The release-please workflow SHALL run on every push to the `main` branch. | The release-please workflow opens or updates a Release PR on every push to main | verbatim |
| Release PR opens for releasable commits on push to main #2: It SHALL open a Release PR if releasable commits exist since the last release tag. | The release-please workflow opens or updates a Release PR on every push to main | verbatim |
| Release PR opens for releasable commits on push to main #3: Releasable commits are those whose type has a visible changelog section in `release-please-config.json` (`f... | The release-please workflow opens or updates a Release PR on every push to main | verbatim |
| Release PR opens for releasable commits on push to main #4: `docs`, `chore`, `test`, `ci` and `build` are hidden and SHALL NOT by themselves open a Release PR, and a `... | Hidden commit types do not open a Release PR | verbatim |
| Release PR opens for releasable commits on push to main #5: If a Release PR already exists, it SHALL update the PR with the latest accumulated changes. | The release-please workflow opens or updates a Release PR on every push to main | verbatim |
| Release PR opens for releasable commits on push to main #6: Hiding `docs` keeps a doc-only commit from cutting an operator release that would cascade a pointless bump ... | Hidden commit types do not open a Release PR | verbatim |
| Release PR passes the release-pin gate #1: On every Release PR, the release-pin gate (G1, workspace RELEASING.md, section "Gates") SHALL run as a step... | The release-pin gate is a step of an existing job | verbatim |
| Release PR passes the release-pin gate #2: A PR is a Release PR when `${{ github.head_ref \|\| github.ref_name }}` starts with `release-please--`, exc... | A Release PR is told by its head branch | verbatim |
| Release PR passes the release-pin gate #3: The gate SHALL fail the job when any of the following holds on the PR head: | The gate fails on a replace directive or an unreleased pin; The gate fails on a dev pin or a tracked local-module file | list cut: the sentence with its four items is 491 on its own and does not fit with its reference; both parts verbatim, subject and the last sentence repeated, each part names the block that holds the rest |
| Release PR passes the release-pin gate #4: - `go.mod` carries a `replace` directive; | The gate fails on a replace directive or an unreleased pin | verbatim |
| Release PR passes the release-pin gate #5: - a required module under `github.com/open-platform-model/` is pinned to a Go pseudo-version, or to a versi... | The gate fails on a replace directive or an unreleased pin | verbatim |
| Release PR passes the release-pin gate #6: - a published fixture's `cue.mod/module.cue` (`test/fixtures/modules/*`, `test/fixtures/modulepackages/*`) ... | The gate fails on a dev pin or a tracked local-module file | verbatim |
| Release PR passes the release-pin gate #7: - any `cue.mod/local-module.cue` is tracked by git. | The gate fails on a dev pin or a tracked local-module file | verbatim |
| Release PR passes the release-pin gate #8: Each failure SHALL name the offending file and pin. | The gate fails on a replace directive or an unreleased pin; The gate fails on a dev pin or a tracked local-module file | verbatim |
| Release PR passes the release-pin gate #9: On pull requests that are not Release PRs the gate step SHALL be skipped. | A Release PR is told by its head branch | verbatim |
| Release PR passes the release-pin gate #10: The same check SHALL be runnable locally as one task. | The release-pin gate is a step of an existing job | verbatim |
| Release PR passes the release-pin gate #11: The job carrying the gate SHALL report a check name no other workflow job in the repo uses (`Lint`), so a r... | The release-pin gate is a step of an existing job | verbatim |
| Dependabot leaves OPM Go modules to the release cascade #1: The `gomod` entry of `.github/dependabot.yml` SHALL ignore every dependency matching `github.com/open-platf... | Dependabot ignores the OPM Go modules | verbatim |
| Dependabot leaves OPM Go modules to the release cascade #2: Those pins move only through the release cascade, which titles a shipped bump `fix(deps)` so it releases (w... | Dependabot ignores the OPM Go modules | verbatim |
| Dependabot leaves OPM Go modules to the release cascade #3: The `github-actions` entry SHALL ignore `open-platform-model/.github*`: the cascade references move togethe... | Dependabot ignores the cascade references of the .github repository | verbatim |
| Dependabot leaves OPM Go modules to the release cascade #4: Third-party Go modules, except `cuelang.org/go` and `cuelabs.dev/go/oci/ociregistry` (see "Dependabot leave... | Third-party updates continue under Dependabot | verbatim |
| The opm CLI version is pinned in one file #1: The opm CLI release that CI and the release workflow install SHALL be declared once, in a repo-root file `.... | The opm CLI version is pinned in one file (parent, MODIFIED) | verbatim |
| The opm CLI version is pinned in one file #2: Every workflow step that installs the opm CLI SHALL read the version from that file and SHALL fail before i... | The opm CLI version is pinned in one file (parent, MODIFIED) | verbatim |
| The opm CLI version is pinned in one file #3: No workflow file SHALL carry an opm CLI version literal. | The opm CLI version is pinned in one file (parent, MODIFIED) | verbatim |
| The opm CLI version is pinned in one file #4: The pin is a release-tool pin: moving it is a `ci(deps)` commit and cuts no release (workspace RELEASING.md... | Moving the opm CLI pin cuts no release | keyword added (Strength table) |
| Forced version via the release-as config key #1: A forced version SHALL be a `release-as` value on package `"."` of `release-please-config.json`, set by a n... | Forced version via the release-as config key (parent, MODIFIED) | verbatim |
| Forced version via the release-as config key #2: While the value stays, release-please proposes exactly that version on every Release PR, so a value left in... | The release-as key repeats its version and crosses a boundary | keyword added (Strength table) |
| Forced version via the release-as config key #3: The key is the required mechanism to cross a boundary the automated bump cannot reach (changing the prerele... | The release-as key repeats its version and crosses a boundary | verbatim |
| Forced version via the release-as config key #4: A `Release-As:` commit footer SHALL NOT be relied on: the squash commit on `main` carries only the PR title... | A Release-As commit footer is not relied on | verbatim |
| Forced version via the release-as config key #5: GA needs no forced version: `prerelease: false` plus a visible carrier commit is enough. | Forced version via the release-as config key (parent, MODIFIED) | verbatim |
| Forced version via the release-as config key #6: A `release-as` value sets the proposed version but does not by itself open a Release PR; one opens only onc... | Forced version via the release-as config key (parent, MODIFIED) | verbatim |
| The Release PR rewrites the install page's operator version #1: `release-please-config.json` SHALL list `docs/site/start/install-the-operator.md` under the root package's ... | The install page's operator version sits in release-please blocks | verbatim |
| The Release PR rewrites the install page's operator version #2: A block SHALL name no version other than the operator's. | The install page's version blocks are tested | verbatim |
| The Release PR rewrites the install page's operator version #3: A test in `internal/version` SHALL fail when the page has no block, when a block is not closed, or when a v... | The install page's version blocks are tested | verbatim |
| Release notifies downstream after it is published #1: The release workflow SHALL contain a job `notify-downstream`, owned by this repo, that | The notify-downstream job needs release-please and publish-release; The notify-downstream job has exactly one step, the cascade action | subject repeated in both parts of the cut sentence |
| Release notifies downstream after it is published #2: needs `release-please` and `publish-release`, runs on `ubuntu-latest`, declares `environment: cascade`, has... | The notify-downstream job needs release-please and publish-release | list cut; the other item of the list is in the next block |
| Release notifies downstream after it is published #3: has exactly one step: the action `open-platform-model/.github/.github/actions/cascade-notify` at the repo's... | The notify-downstream job has exactly one step, the cascade action | list cut; joined to the subject and the job name of the old sentence |
| Release notifies downstream after it is published #4: Its `if` SHALL be `needs.release-please.outputs.release_created == 'true' && vars.CASCADE_NOTIFY != 'off'`:... | The notify-downstream job runs only after a published operator release | verbatim |
| Release notifies downstream after it is published #5: It SHALL NOT wait for, or be skipped by, `publish-docs`. | The notify-downstream job needs release-please and publish-release | verbatim |
| Release notifies downstream after it is published #6: The key SHALL be read only in the caller-owned `notify-downstream` or `publish` job, which declares `enviro... | The cascade App key is read only by the notify job | verbatim |
| Release notifies downstream after it is published #7: The action mints the App token and SHALL dispatch `upstream-released` only to `cli` (Phase 3 wiring contrac... | The cascade action dispatches upstream-released only to the cli | verbatim |
| Dependabot leaves cuelang.org/go to library releases #1: The `gomod` entry of `.github/dependabot.yml` SHALL ignore the dependencies `cuelang.org/go` and `cuelabs.d... | Dependabot leaves cuelang.org/go to library releases (parent, MODIFIED) | cut at the clause boundary before the colon; the rest is in the next block |
| Dependabot leaves cuelang.org/go to library releases #2: the library bumps CUE in its own pull request, where the closedness canary runs, and the operator picks the... | The operator's CUE version moves only through a library release | subject "The operator's CUE version SHALL move only through a library release" copied from the clause before the colon; keyword added (Strength table) |
| Dependabot leaves cuelang.org/go to library releases #3: Other third-party Go modules SHALL keep their Dependabot updates. | Dependabot leaves cuelang.org/go to library releases (parent, MODIFIED) | verbatim |
| Dependabot leaves cuelang.org/go to library releases #4: Source: owner decision j4 of the kernel beta.1 plan walkthrough (2026-10-03); workspace RELEASING.md, secti... | Dependabot leaves cuelang.org/go to library releases (parent, MODIFIED); The operator's CUE version moves only through a library release | carried as the **Source** line of the parent and of the new block (full old Source line) |

### Ledger: container-image-publish

Forms and lengths (old body, then each new body; a MODIFIED parent is measured on its body only):

| Old requirement | Form | Old length | New blocks (length) |
| --- | --- | --- | --- |
| Pull-request image build trigger and tags | c2 (REMOVED plus ADDED under new names) | 859 | The PR image workflow runs on pull request events and skips bot heads (435); The PR image build pushes a sha tag and a PR tag (464); The PR image build does not use the GitHub Actions cache (140) |
| Release image tags | c2 (REMOVED plus ADDED under new names) | 1212 | A gated release run pushes three tags (390); The release image carries the release commit and the release tag (253); The version tag always resolves to the manifest list first pushed under it (347); A version tag built from the release commit is reused (438); A version tag held by another revision fails the job (318) |

`grep -rn --exclude-dir=archive` for the old headers found no reference outside this spec.

Scenarios:

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| New PR opened | The PR image build pushes a sha tag and a PR tag, New PR opened | no (byte equal, checked by the item 3 command) |
| Force push updates an open PR | The PR image build pushes a sha tag and a PR tag, Force push updates an open PR | no (byte equal, checked by the item 3 command) |
| Non-applicable event | The PR image workflow runs on pull request events and skips bot heads, Non-applicable event | no (byte equal, checked by the item 3 command) |
| Bot head skipped | The PR image workflow runs on pull request events and skips bot heads, Bot head skipped | no (byte equal, checked by the item 3 command) |
| (new) The PR image build uses no Actions cache | The PR image build does not use the GitHub Actions cache, The PR image build uses no Actions cache | new: restates a case the old text names |
| First release v0.1.0 | A gated release run pushes three tags, First release v0.1.0 | no (byte equal, checked by the item 3 command) |
| Subsequent release v0.2.0 | A gated release run pushes three tags, Subsequent release v0.2.0 | no (byte equal, checked by the item 3 command) |
| Re-run after the version tag was pushed | A version tag built from the release commit is reused, Re-run after the version tag was pushed | no (byte equal, checked by the item 3 command) |
| Version tag held by a different image | A version tag held by another revision fails the job, Version tag held by a different image | no (byte equal, checked by the item 3 command) |
| Release cut on a later push | The release image carries the release commit and the release tag, Release cut on a later push | no (byte equal, checked by the item 3 command) |
| (new) An absent version tag is built, pushed and confirmed | The version tag always resolves to the manifest list first pushed under it, An absent version tag is built, pushed and confirmed | new: restates a case the old text names |

Clauses (the rule sentences of each old body and the block that carries each):

| Old clause | New requirement | Note |
| --- | --- | --- |
| Pull-request image build trigger and tags #1: The PR image workflow SHALL run on `pull_request` events (types: `opened`, `synchronize`, `reopened`) whate... | The PR image workflow runs on pull request events and skips bot heads | verbatim |
| Pull-request image build trigger and tags #2: On every run it SHALL build and push exactly two tags: `:sha-<short>` where `<short>` is the 7-character sh... | The PR image build pushes a sha tag and a PR tag | verbatim |
| Pull-request image build trigger and tags #3: The `:pr-<PR_ID>` tag MAY be overwritten by subsequent pushes to the same PR; `:sha-<short>` is effectively... | The PR image build pushes a sha tag and a PR tag | verbatim |
| Pull-request image build trigger and tags #4: The build SHALL NOT use the GitHub Actions cache. | The PR image build does not use the GitHub Actions cache | verbatim |
| Release image tags #1: On a gated release run, the image-release job SHALL push, for a release version `v<MAJOR>.<MINOR>.<PATCH>`,... | A gated release run pushes three tags | verbatim |
| Release image tags #2: The release commit is the commit the release tag names, read from the checked-out tag, not the pushed commi... | The release image carries the release commit and the release tag | verbatim |
| Release image tags #3: The version tag SHALL always resolve to the manifest list first pushed under it (enhancement 0021, "Release... | The version tag always resolves to the manifest list first pushed under it | verbatim |
| Release image tags #4: When it is absent the job builds and pushes, then SHALL confirm the version tag resolves to the pushed digest. | The version tag always resolves to the manifest list first pushed under it | verbatim |
| Release image tags #5: When it exists with a revision label equal to the release commit the job SHALL push nothing under the versi... | A version tag built from the release commit is reused | verbatim |
| Release image tags #6: When it exists with any other revision, or the probe fails for any reason other than the tag being absent, ... | A version tag held by another revision fails the job | verbatim |
| Release image tags #7: `:latest` and `:sha-<short>` are mutable by design. | A gated release run pushes three tags | verbatim |

Added words (the only ones): subjects and conditions copied from the same old requirement, the exact header of a sibling ("As in requirement" prefix or "(Requirement: ...)" suffix), and the keywords in the Strength table.

### Ledger: reconcile-loop-assembly

Forms and lengths (old body, then each new body; a MODIFIED parent is measured on its body only):

| Old requirement | Form | Old length | New blocks (length) |
| --- | --- | --- | --- |
| Status always patched | c2 (REMOVED plus ADDED under new names) | 2906 | Every reconcile attempt ends with a status patch (351); A successful outcome patches the Healthy condition it judged (236); A NoOp patch is bounded to a fixed set of fields (469); A NoOp patch leaves lastAttempted, history and inventory alone (232); A reconcile that skips its render patches only a changed Healthy condition (351); requiredContracts is in the NoOp patch set (125); A NoOp patch carries requiredContracts because the demand can move without a digest (450); lastAppliedVersion is in the NoOp patch set (431); lastAppliedInputs is in the NoOp patch set while the drift render interval is positive (420); Storm safety comes from the generation predicate (225) |
| Inventory updated only on full success | c2 (REMOVED plus ADDED under new names) | 947 | The inventory is replaced only after a full success (282); An expired Job entry leaves the inventory (209); An entry the apply verdict refuses as adopted-elsewhere leaves the inventory (366); The inventory digest stays the digest of the rendered set (284); A refused or failed reconcile keeps the inventory (100) |
| The instance's contract demand is recorded on its status | b (MODIFIED parent plus ADDED split) | 678 | The instance's contract demand is recorded on its status (368); A reconcile that does not render leaves requiredContracts at its previous value (387) |
| A rendered resource may be withheld from apply | c2 (REMOVED plus ADDED under new names) | 749 | The reconciler evaluates rendered resources and may withhold one (340); Withholding does not change the inventory (337); A reconcile that withheld a resource is not ready and not a no-op (156) |
| A recovered panic is recorded as a failed attempt | c2 (REMOVED plus ADDED under new names) | 1762 | A recovered panic is never reported as a success (159); A panic is recorded in Ready and history (495); A panic counts one failure (497); A panic keeps the applied state (411); A panicking attempt counts as FailedTransient for counters, history and metrics only (495); A recovered panic is raised again with its original value (342) |
| The applied module version is recorded in plain text | c2 (REMOVED plus ADDED under new names) | 1820 | The applied module version is recorded on status.lastAppliedVersion (438); The applied version is written with the lastApplied fields on a success (373); A NoOp that rendered sets the applied version (375); A NoOp fills an empty applied version and a render-skipping reconcile leaves it (364); A version that cannot be read does not fail the reconcile (276); The applied version is not an input to no-op detection (179) |
| A healthy ModuleInstance is reconciled periodically | c2 (REMOVED plus ADDED under new names) | 1167 | The instance reconcile interval comes from a flag (410); A reconciled instance requeues after the interval plus a random addition (360); The periodic requeue sets no nextRetryAt and is not a failure (265); A suspended, CLI-owned or operator instance is not requeued by the interval (260); An unchanged periodic reconcile ends as a skipped render or a NoOp (192) |
| The instance identities are stored before the first write of an apply | c2 (REMOVED plus ADDED under new names) | 1779 | An instance carries its identity and the earlier one until a change is settled (284); A first identity is stored (424); A changed identity keeps the earlier one (479); A return to the earlier identity swaps the two (411); previousInstanceUUID is cleared only by a full success (292); A reconcile that applies nothing writes an empty instanceUUID and leaves both fields alone otherwise (226); An unsettled identity change is applied and pruned even when the render equals the last apply (302) |
| A reconcile the apply verdict refuses is retried and not stalled | b (MODIFIED parent plus ADDED split) | 886 | A reconcile the apply verdict refuses is retried and not stalled (418); A reconcile with unchanged digests is not refused for an object being deleted (318); A ModulePackage reconcile on matching digests is refused until the object is gone (327) |

`grep -rn --exclude-dir=archive` for the old headers found `The instance's contract demand is recorded on its status` cited by `kernel-module-renderer` and `A reconcile the apply verdict refuses is retried and not stalled` cited by `modulepackage-reconcile-loop`. Both keep their header (form b: MODIFIED parent with its scenarios, split-off behavior ADDED), so no reference changes. No other old header of this unit is cited outside this spec.

Scenarios:

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| Status updated on failure | Every reconcile attempt ends with a status patch, Status updated on failure | no (byte equal, checked by the item 3 command) |
| NoOp removes an expired Job entry | A NoOp patch is bounded to a fixed set of fields, NoOp removes an expired Job entry | no (byte equal, checked by the item 3 command) |
| NoOp patches drift and counters | A NoOp patch leaves lastAttempted, history and inventory alone, NoOp patches drift and counters | no (byte equal, checked by the item 3 command) |
| A skipped render patches nothing | A reconcile that skips its render patches only a changed Healthy condition, A skipped render patches nothing | no (byte equal, checked by the item 3 command) |
| A skipped render records a changed health judgement | A reconcile that skips its render patches only a changed Healthy condition, A skipped render records a changed health judgement | no (byte equal, checked by the item 3 command) |
| NoOp refreshes the demand the render just reported | requiredContracts is in the NoOp patch set, NoOp refreshes the demand the render just reported | no (byte equal, checked by the item 3 command) |
| NoOp patch does not trigger reconcile loop | Storm safety comes from the generation predicate, NoOp patch does not trigger reconcile loop | no (byte equal, checked by the item 3 command) |
| Partial failure preserves inventory | The inventory is replaced only after a full success, Partial failure preserves inventory | no (byte equal, checked by the item 3 command) |
| An expired Job leaves the inventory without an apply | An expired Job entry leaves the inventory, An expired Job leaves the inventory without an apply | no (byte equal, checked by the item 3 command) |
| An adopted object leaves the inventory without an apply | An entry the apply verdict refuses as adopted-elsewhere leaves the inventory, An adopted object leaves the inventory without an apply | no (byte equal, checked by the item 3 command) |
| A refused reconcile keeps the inventory | A refused or failed reconcile keeps the inventory, A refused reconcile keeps the inventory | no (byte equal, checked by the item 3 command) |
| A successful render records the demand | The instance's contract demand is recorded on its status, A successful render records the demand | no (byte equal, checked by the item 3 command) |
| A failed render leaves the previous value | The instance's contract demand is recorded on its status, A failed render leaves the previous value | no (byte equal, checked by the item 3 command) |
| The demand follows the spec | The instance's contract demand is recorded on its status, The demand follows the spec | no (byte equal, checked by the item 3 command) |
| The rest of the module still applies | The reconciler evaluates rendered resources and may withhold one, The rest of the module still applies | no (byte equal, checked by the item 3 command) |
| Nothing withheld behaves as before | The reconciler evaluates rendered resources and may withhold one, Nothing withheld behaves as before | no (byte equal, checked by the item 3 command) |
| The withheld resource stays owned | Withholding does not change the inventory, The withheld resource stays owned | no (byte equal, checked by the item 3 command) |
| A withheld resource keeps the instance unconverged | A reconcile that withheld a resource is not ready and not a no-op, A withheld resource keeps the instance unconverged | no (byte equal, checked by the item 3 command) |
| A panicking render does not mark a ready instance Ready | A recovered panic is never reported as a success, A panicking render does not mark a ready instance Ready | no (byte equal, checked by the item 3 command) |
| A panicking render does not mark a ready package Ready | A recovered panic is never reported as a success, A panicking render does not mark a ready package Ready | no (byte equal, checked by the item 3 command) |
| A panic inside a phase does not reset that phase's counter | A panic counts one failure, A panic inside a phase does not reset that phase's counter | no (byte equal, checked by the item 3 command) |
| A panic during conversion is recorded and frees the slot | A recovered panic is raised again with its original value, A panic during conversion is recorded and frees the slot | no (byte equal, checked by the item 3 command) |
| A successful apply records the module version | The applied module version is recorded on status.lastAppliedVersion, A successful apply records the module version | no (byte equal, checked by the item 3 command) |
| A ModulePackage records the version of the module its instance renders | The applied module version is recorded on status.lastAppliedVersion, A ModulePackage records the version of the module its instance renders | no (byte equal, checked by the item 3 command) |
| A failed apply keeps the previous version | The applied version is written with the lastApplied fields on a success, A failed apply keeps the previous version | no (byte equal, checked by the item 3 command) |
| A NoOp fills an empty field | A NoOp fills an empty applied version and a render-skipping reconcile leaves it, A NoOp fills an empty field | no (byte equal, checked by the item 3 command) |
| A NoOp after an ownership handback corrects a stale version | A NoOp fills an empty applied version and a render-skipping reconcile leaves it, A NoOp after an ownership handback corrects a stale version | no (byte equal, checked by the item 3 command) |
| An apply whose render reports no version clears the field | A version that cannot be read does not fail the reconcile, An apply whose render reports no version clears the field | no (byte equal, checked by the item 3 command) |
| The interval is disabled | The instance reconcile interval comes from a flag, The interval is disabled | no (byte equal, checked by the item 3 command) |
| A negative interval stops the operator | The instance reconcile interval comes from a flag, A negative interval stops the operator | no (byte equal, checked by the item 3 command) |
| A rolled-out instance requeues on the interval | A reconciled instance requeues after the interval plus a random addition, A rolled-out instance requeues on the interval | no (byte equal, checked by the item 3 command) |
| A NoOp and a skipped render requeue on the interval | A reconciled instance requeues after the interval plus a random addition, A NoOp and a skipped render requeue on the interval | no (byte equal, checked by the item 3 command) |
| A suspended instance is not requeued | A suspended, CLI-owned or operator instance is not requeued by the interval, A suspended instance is not requeued | no (byte equal, checked by the item 3 command) |
| A CLI-owned instance is not requeued | A suspended, CLI-owned or operator instance is not requeued by the interval, A CLI-owned instance is not requeued | no (byte equal, checked by the item 3 command) |
| A skipped periodic reconcile writes nothing | An unchanged periodic reconcile ends as a skipped render or a NoOp, A skipped periodic reconcile writes nothing | no (byte equal, checked by the item 3 command) |
| First reconcile records the identity | A first identity is stored, First reconcile records the identity | no (byte equal, checked by the item 3 command) |
| The identities cannot be stored | A first identity is stored, The identities cannot be stored | no (byte equal, checked by the item 3 command) |
| Both identities are stored before the apply | A changed identity keeps the earlier one, Both identities are stored before the apply | no (byte equal, checked by the item 3 command) |
| Going back before the change is settled | A return to the earlier identity swaps the two, Going back before the change is settled | no (byte equal, checked by the item 3 command) |
| Full success settles the change | previousInstanceUUID is cleared only by a full success, Full success settles the change | no (byte equal, checked by the item 3 command) |
| A refused apply keeps both identities | previousInstanceUUID is cleared only by a full success, A refused apply keeps both identities | no (byte equal, checked by the item 3 command) |
| A panic keeps both identities | previousInstanceUUID is cleared only by a full success, A panic keeps both identities | no (byte equal, checked by the item 3 command) |
| An object without the field gains it without an apply | A reconcile that applies nothing writes an empty instanceUUID and leaves both fields alone otherwise, An object without the field gains it without an apply | no (byte equal, checked by the item 3 command) |
| Going back renders what was last applied | An unsettled identity change is applied and pruned even when the render equals the last apply, Going back renders what was last applied | no (byte equal, checked by the item 3 command) |
| A refusal is transient | A reconcile the apply verdict refuses is retried and not stalled, A refusal is transient | no (byte equal, checked by the item 3 command) |
| The refusal ends when the object is adopted | A reconcile the apply verdict refuses is retried and not stalled, The refusal ends when the object is adopted | no (byte equal, checked by the item 3 command) |
| A refused identity change stores nothing | A reconcile the apply verdict refuses is retried and not stalled, A refused identity change stores nothing | no (byte equal, checked by the item 3 command) |
| An inventoried object being deleted on unchanged digests | A reconcile the apply verdict refuses is retried and not stalled, An inventoried object being deleted on unchanged digests | no (byte equal, checked by the item 3 command) |
| A rendered object held by another instance on unchanged digests | A reconcile the apply verdict refuses is retried and not stalled, A rendered object held by another instance on unchanged digests | no (byte equal, checked by the item 3 command) |

New scenarios (each restates a case the old text names):

| New place | New scenario | Old text that names the case |
| --- | --- | --- |
| A successful outcome patches the Healthy condition it judged | A failed attempt leaves Healthy as it was | the sentence(s) of this block |
| A NoOp patch carries requiredContracts because the demand can move without a digest | A regenerated platform moves the demand without a digest change | the sentence(s) of this block |
| lastAppliedVersion is in the NoOp patch set | A NoOp that rendered writes the version it reported | the sentence(s) of this block |
| lastAppliedInputs is in the NoOp patch set while the drift render interval is positive | A NoOp does not write lastAppliedInputs while the drift render interval is zero | the sentence(s) of this block |
| The inventory digest stays the digest of the rendered set | A removal leaves the digest of the rendered set | the sentence(s) of this block |
| A reconcile that does not render leaves requiredContracts at its previous value | A suspended instance leaves the previous value | the sentence(s) of this block |
| A panic is recorded in Ready and history | A panic is recorded in Ready and history | the sentence(s) of this block |
| A panic keeps the applied state | A panic clears nextRetryAt and keeps the applied state | the sentence(s) of this block |
| A panicking attempt counts as FailedTransient for counters, history and metrics only | A panicking attempt returns no ctrl.Result | the sentence(s) of this block |
| A NoOp that rendered sets the applied version | A NoOp that rendered sets the version the render reported | the sentence(s) of this block |
| A NoOp fills an empty applied version and a render-skipping reconcile leaves it | A reconcile that skips its render leaves the applied version | the sentence(s) of this block |
| The applied version is not an input to no-op detection | The applied version is not compared by no-op detection | the sentence(s) of this block |
| The periodic requeue sets no nextRetryAt and is not a failure | The periodic requeue is not a failure | the sentence(s) of this block |
| An instance carries its identity and the earlier one until a change is settled | An unsettled change of identity carries both identities | the sentence(s) of this block |
| A reconcile with unchanged digests is not refused for an object being deleted | A no-op over an inventoried object being deleted is not refused | the sentence(s) of this block |
| A ModulePackage reconcile on matching digests is refused until the object is gone | A ModulePackage reconcile on matching digests is refused until the object is gone | the sentence(s) of this block |

Clauses (the rule sentences of each old body and the block that carries each):

| Old clause | New requirement | Note |
| --- | --- | --- |
| Status always patched #1: The reconciler MUST patch `ModuleRelease.status` at the end of every reconcile attempt, including `NoOp`. | Every reconcile attempt ends with a status patch | verbatim |
| Status always patched #2: On meaningful outcomes (`Applied`, `AppliedAndPruned`, `FailedTransient`, `FailedStalled`), the patch updates ... | Every reconcile attempt ends with a status patch | keyword added (strength table) |
| Status always patched #3: On `NoOp`, the patch is bounded to: drift condition (`Drifted`), the `Healthy` condition (`instance-health`), ... | A NoOp patch is bounded to a fixed set of fields | keyword added (strength table) |
| Status always patched #4: `lastAttempted*` and history MUST NOT be modified on `NoOp`, and `inventory` MUST NOT be modified in any other... | A NoOp patch leaves lastAttempted, history and inventory alone | verbatim |
| Status always patched #5: A reconcile that skips its render because its inputs are unchanged (`render-input-key`) is not an attempt. | A reconcile that skips its render patches only a changed Healthy condition | keyword added (strength table) |
| Status always patched #6: It judges health (`instance-health`) and MUST NOT patch anything but the `Healthy` condition, which it patches... | A reconcile that skips its render patches only a changed Healthy condition | verbatim |
| Status always patched #7: On a successful outcome the patch also carries the `Healthy` condition the reconcile judged (`instance-health`... | A successful outcome patches the Healthy condition it judged | keyword added (strength table) |
| Status always patched #8: A failed attempt leaves `Healthy` as it was. | A successful outcome patches the Healthy condition it judged | keyword added (strength table) |
| Status always patched #9: `requiredContracts` is in the `NoOp` set deliberately, and it is the one field there that does not describe an... | requiredContracts is in the NoOp patch set | keyword added (strength table) |
| Status always patched #10: It describes what the instance demands, which a reconcile can observe as changed without any digest changing: ... | A NoOp patch carries requiredContracts because the demand can move without a digest | keyword added (strength table) |
| Status always patched #11: Bounding it out of the `NoOp` patch would leave the field describing a render the operator has already superse... | A NoOp patch carries requiredContracts because the demand can move without a digest | verbatim |
| Status always patched #12: `lastAppliedVersion` is in the `NoOp` set because a `NoOp` re-proves it. | lastAppliedVersion is in the NoOp patch set | keyword added (strength table) |
| Status always patched #13: A `NoOp` means the source digest equals the last apply's, and the source digest pins the module version, so th... | lastAppliedVersion is in the NoOp patch set | verbatim |
| Status always patched #14: Writing it fills the field on an object applied before it existed and corrects it after an ownership handback. | lastAppliedVersion is in the NoOp patch set | verbatim |
| Status always patched #15: A reconcile that did not render leaves it as it was. | lastAppliedVersion is in the NoOp patch set | keyword added (strength table) |
| Status always patched #16: `lastAppliedInputs` is in the `NoOp` set for the same reason: a `NoOp` that rendered proves the cluster holds ... | lastAppliedInputs is in the NoOp patch set while the drift render interval is positive | keyword added (strength table) |
| Status always patched #17: It is written on a `NoOp` only while the drift render interval is greater than zero; with the skip disabled no... | lastAppliedInputs is in the NoOp patch set while the drift render interval is positive | keyword added (strength table) |
| Status always patched #18: Storm safety is provided by `WithEventFilter(predicate.GenerationChangedPredicate{})` on the controller. | Storm safety comes from the generation predicate | keyword added (strength table) |
| Status always patched #19: Status subresource patches do not bump `metadata.generation`, so the predicate filters them at the watch bound... | Storm safety comes from the generation predicate | verbatim |
| Inventory updated only on full success #1: The `status.inventory` MUST only be replaced after a fully successful apply (and prune, if enabled). | The inventory is replaced only after a full success | verbatim |
| Inventory updated only on full success #2: Two narrower writes are allowed. | The inventory is replaced only after a full success | verbatim |
| Inventory updated only on full success #3: A ModuleInstance reconcile whose digests are unchanged MUST remove the entries of expired Jobs (`drift-detecti... | An expired Job entry leaves the inventory | verbatim |
| Inventory updated only on full success #4: A reconcile that renders MUST leave out of the inventory it records, on a `NoOp` as well, the entry of every o... | An entry the apply verdict refuses as adopted-elsewhere leaves the inventory | verbatim |
| Inventory updated only on full success #5: `status.inventory.digest` MUST stay the digest of the rendered set, so after such a removal it is not the dige... | The inventory digest stays the digest of the rendered set | verbatim |
| Inventory updated only on full success #6: A reconcile that is refused by the apply verdict, or that fails, MUST NOT change `status.inventory`. | A refused or failed reconcile keeps the inventory | verbatim |
| The instance's contract demand is recorded on its status #1: The reconciler SHALL record on `ModuleInstance.status.requiredContracts` the set of contract FQNs the instance... | The instance's contract demand is recorded on its status | verbatim |
| The instance's contract demand is recorded on its status #2: The field SHALL be derived from the rendered instance and SHALL NOT be readable from, or influenced by, anythi... | The instance's contract demand is recorded on its status | verbatim |
| The instance's contract demand is recorded on its status #3: A reconcile that does not render SHALL leave the field at its previous value. | A reconcile that does not render leaves requiredContracts at its previous value | verbatim |
| The instance's contract demand is recorded on its status #4: This includes a failed render, a suspended instance and a CLI-owned instance: a stale value over-reports deman... | A reconcile that does not render leaves requiredContracts at its previous value | verbatim |
| A rendered resource may be withheld from apply #1: The reconciler SHALL evaluate the rendered resources before applying them and MAY withhold an individual resou... | The reconciler evaluates rendered resources and may withhold one | verbatim |
| A rendered resource may be withheld from apply #2: A withheld resource SHALL NOT be applied, and every other rendered resource SHALL be applied as usual, so one ... | The reconciler evaluates rendered resources and may withhold one | verbatim |
| A rendered resource may be withheld from apply #3: Withholding SHALL NOT change the inventory. | Withholding does not change the inventory | verbatim |
| A rendered resource may be withheld from apply #4: The inventory records which resources the instance owns, and the operator still owns a resource it declined to... | Withholding does not change the inventory | verbatim |
| A rendered resource may be withheld from apply #5: A reconcile that withheld a resource SHALL NOT report ready and SHALL NOT be treated as a no-op, because the c... | A reconcile that withheld a resource is not ready and not a no-op | verbatim |
| A recovered panic is recorded as a failed attempt #1: When a ModuleInstance or ModulePackage reconcile panics after its deferred status commit is installed, the rec... | A recovered panic is never reported as a success | verbatim |
| A recovered panic is recorded as a failed attempt #2: Before the panic leaves the reconcile function, the deferred commit MUST record the attempt as a transient fai... | A panic is recorded in Ready and history | verbatim |
| A recovered panic is recorded as a failed attempt #3: - `Ready=False` with reason `ReconcilePanic` and a message that starts with `reconcile panicked: | A panic is recorded in Ready and history | verbatim |
| A recovered panic is recorded as a failed attempt #4: ` followed by the panic value. | A panic is recorded in Ready and history | verbatim |
| A recovered panic is recorded as a failed attempt #5: - `Reconciling=True`, and no `Stalled` condition. | A panic is recorded in Ready and history | verbatim |
| A recovered panic is recorded as a failed attempt #6: - `lastAttempted*` set for this attempt, and a failure history entry carrying the same message. | A panic is recorded in Ready and history | verbatim |
| A recovered panic is recorded as a failed attempt #7: - `failureCounters.reconcile` incremented, as for any `FailedTransient` outcome. | A panic counts one failure | verbatim |
| A recovered panic is recorded as a failed attempt #8: - The drift, apply and prune failure counters left as they were. | A panic counts one failure | verbatim |
| A recovered panic is recorded as a failed attempt #9: A phase that was in flight when the panic hit did not succeed, so a panicking attempt never resets a phase cou... | A panic counts one failure | verbatim |
| A recovered panic is recorded as a failed attempt #10: - `nextRetryAt` cleared, because the retry is scheduled by the controller runtime's rate limiter, not by the o... | A panic keeps the applied state | verbatim |
| A recovered panic is recorded as a failed attempt #11: - `lastApplied*` and `inventory` left as they were. | A panic keeps the applied state | verbatim |
| A recovered panic is recorded as a failed attempt #12: Fields the attempt already wrote on the object before the panic (`requiredContracts`, `instanceUUID`, the `Dri... | A panicking attempt counts as FailedTransient for counters, history and metrics only | keyword added (strength table) |
| A recovered panic is recorded as a failed attempt #13: A panicking attempt is counted as `FailedTransient` for failure counters, history and reconcile metrics only. | A panicking attempt counts as FailedTransient for counters, history and metrics only | keyword added (strength table) |
| A recovered panic is recorded as a failed attempt #14: It returns no `ctrl.Result`: the requeue scenarios of "Outcome classification" and the `nextRetryAt` rule of t... | A panicking attempt counts as FailedTransient for counters, history and metrics only | keyword added (strength table) |
| A recovered panic is recorded as a failed attempt #15: The reconciler MUST then re-panic with the original panic value, so the controller runtime still recovers, log... | A recovered panic is raised again with its original value | verbatim |
| A recovered panic is recorded as a failed attempt #16: The operator SHOULD log the panic value with the object's name and namespace before it re-panics. | A recovered panic is raised again with its original value | verbatim |
| The applied module version is recorded in plain text #1: The ModuleInstance and ModulePackage reconcilers SHALL record on `status.lastAppliedVersion` the version of th... | The applied module version is recorded on status.lastAppliedVersion | verbatim |
| The applied module version is recorded in plain text #2: The value SHALL be the version the rendered instance's source module declares in its `metadata.version` (bare ... | The applied module version is recorded on status.lastAppliedVersion | verbatim |
| The applied module version is recorded in plain text #3: The field SHALL be written together with the other `lastApplied*` fields, on an attempt whose apply (and prune... | The applied version is written with the lastApplied fields on a success | verbatim |
| The applied module version is recorded in plain text #4: An attempt that fails, panics, is suspended, belongs to a CLI-owned instance, or is the operator's own instanc... | The applied version is written with the lastApplied fields on a success | verbatim |
| The applied module version is recorded in plain text #5: When the outcome is `NoOp` and the attempt rendered, the reconciler SHALL set it to the version that render re... | A NoOp that rendered sets the applied version | verbatim |
| The applied module version is recorded in plain text #6: This fills the field on an object applied before it existed and corrects it after an ownership handback, where... | A NoOp fills an empty applied version and a render-skipping reconcile leaves it | verbatim |
| The applied module version is recorded in plain text #7: A reconcile that did not render (for example one that skips the render because its inputs are unchanged) SHALL... | A NoOp fills an empty applied version and a render-skipping reconcile leaves it | verbatim |
| The applied module version is recorded in plain text #8: A version that cannot be read as a concrete string SHALL NOT fail the reconcile. | A version that cannot be read does not fail the reconcile | verbatim |
| The applied module version is recorded in plain text #9: A successful apply then writes the field empty, so it never names a version the last apply did not render. | A version that cannot be read does not fail the reconcile | verbatim |
| The applied module version is recorded in plain text #10: The version SHALL NOT be an input to no-op detection; the source digest already covers it. | The applied version is not an input to no-op detection | verbatim |
| The applied module version is recorded in plain text #11: Source: owner decision j5 of the kernel-plan walkthrough (2026-10-03): "Operator adds an additive status.lastA... | The applied module version is recorded on status.lastAppliedVersion | the citation moves to a `**Source**:` metadata line of this block (full old Source line) |
| A healthy ModuleInstance is reconciled periodically #1: The operator SHALL take an instance reconcile interval from the flag `--instance-reconcile-interval` (default ... | The instance reconcile interval comes from a flag | verbatim |
| A healthy ModuleInstance is reconciled periodically #2: When a ModuleInstance reconcile ends with outcome `NoOp`, `Applied` or `AppliedAndPruned`, or skips its render... | A reconciled instance requeues after the interval plus a random addition | verbatim |
| A healthy ModuleInstance is reconciled periodically #3: When the health judgement asks for a requeue, that requeue SHALL apply unchanged. | The periodic requeue sets no nextRetryAt and is not a failure | verbatim |
| A healthy ModuleInstance is reconciled periodically #4: The periodic requeue SHALL NOT set `status.nextRetryAt` and SHALL NOT count as a failure. | The periodic requeue sets no nextRetryAt and is not a failure | verbatim |
| A healthy ModuleInstance is reconciled periodically #5: An interval of `0` SHALL disable the periodic requeue. | The instance reconcile interval comes from a flag | verbatim |
| A healthy ModuleInstance is reconciled periodically #6: The operator MUST refuse to start with a negative interval. | The instance reconcile interval comes from a flag | verbatim |
| A healthy ModuleInstance is reconciled periodically #7: The interval is operator-wide: the ModuleInstance API SHALL NOT gain a field for it. | The instance reconcile interval comes from a flag | verbatim |
| A healthy ModuleInstance is reconciled periodically #8: A suspended ModuleInstance, a CLI-owned ModuleInstance and the operator's own instance SHALL NOT be requeued b... | A suspended, CLI-owned or operator instance is not requeued by the interval | verbatim |
| A healthy ModuleInstance is reconciled periodically #9: A failed reconcile SHALL keep the requeue of its failure class. | A suspended, CLI-owned or operator instance is not requeued by the interval | verbatim |
| A healthy ModuleInstance is reconciled periodically #10: The periodic reconcile of an instance whose inputs and objects are unchanged SHALL end as a skipped render or ... | An unchanged periodic reconcile ends as a skipped render or a NoOp | verbatim |
| The instance identities are stored before the first write of an apply #1: A ModuleInstance and a ModulePackage MUST carry in `status.instanceUUID` the identity of the instance the oper... | An instance carries its identity and the earlier one until a change is settled | verbatim |
| The instance identities are stored before the first write of an apply #2: Source: owner decision of 2026-10-08 (the status holds the earlier and the new identity until the prune of tha... | An instance carries its identity and the earlier one until a change is settled | the citation moves to a `**Source**:` metadata line of this block (full old Source line) |
| The instance identities are stored before the first write of an apply #3: When a render's identity differs from `status.instanceUUID`, the reconciler MUST store the new state in the ob... | A first identity is stored | verbatim |
| The instance identities are stored before the first write of an apply #4: - `status.instanceUUID` empty: it becomes the render's identity and `status.previousInstanceUUID` stays empty. | A first identity is stored | verbatim |
| The instance identities are stored before the first write of an apply #5: - `status.instanceUUID` is `A`, `status.previousInstanceUUID` is empty and the render carries `B`: | A changed identity keeps the earlier one | verbatim |
| The instance identities are stored before the first write of an apply #6: `status.instanceUUID` becomes `B` and `status.previousInstanceUUID` becomes `A`. | A changed identity keeps the earlier one | verbatim |
| The instance identities are stored before the first write of an apply #7: - `status.instanceUUID` is `B`, `status.previousInstanceUUID` is `A` and the render carries `A`: the two value... | A return to the earlier identity swaps the two | verbatim |
| The instance identities are stored before the first write of an apply #8: `status.previousInstanceUUID` MUST be cleared in the status commit of a reconcile whose apply and prune succee... | previousInstanceUUID is cleared only by a full success | verbatim |
| The instance identities are stored before the first write of an apply #9: A reconcile that applies nothing MUST write `status.instanceUUID` when it is empty and MUST leave both fields ... | A reconcile that applies nothing writes an empty instanceUUID and leaves both fields alone otherwise | verbatim |
| The instance identities are stored before the first write of an apply #10: While `status.previousInstanceUUID` is set, or a render is about to set it, a reconcile MUST NOT end as one th... | An unsettled identity change is applied and pruned even when the render equals the last apply | verbatim |
| A reconcile the apply verdict refuses is retried and not stalled #1: When the apply verdict refuses a reconcile (`ssa-apply`, "An apply is judged by the ownership verdict before i... | A reconcile the apply verdict refuses is retried and not stalled | verbatim |
| A reconcile the apply verdict refuses is retried and not stalled #2: `status.inventory`, the applied digests, `status.instanceUUID` and `status.previousInstanceUUID` MUST keep the... | A reconcile the apply verdict refuses is retried and not stalled | verbatim |
| A reconcile the apply verdict refuses is retried and not stalled #3: A reconcile whose digests are unchanged MUST NOT refuse because an inventoried object is being deleted, when i... | A reconcile with unchanged digests is not refused for an object being deleted | verbatim |
| A reconcile the apply verdict refuses is retried and not stalled #4: A ModulePackage reconcile that applies on matching digests (`modulepackage-reconcile-loop`, "No-op detection")... | A ModulePackage reconcile on matching digests is refused until the object is gone | verbatim |

Added words (the only ones): subjects and conditions copied from the same old requirement, the exact header of a sibling ("As in requirement" prefix or "(Requirement: ...)" suffix), and the keywords in the Strength table. Sentences left without a keyword because they explain a rule and state none: "Two narrower writes are allowed.", "A `NoOp` means the source digest equals the last apply's, ...", "Writing it fills the field ...", "Bounding it out of the `NoOp` patch would leave ...", "Status subresource patches do not bump `metadata.generation`, ...", "This fills the field ...". A cut list repeats the lead-in sentence of the old body in each part and names the other parts; the first block of a cut carries the lead-in with `the panic` as in the old body. In "A recovered panic is recorded as a failed attempt" the sentences of the closing paragraphs are regrouped (counted as `FailedTransient`, no `ctrl.Result`, fields already written) in one block; their order within it follows the sentence that names the subject first.

### Ledger: operator-module

Forms and lengths (old body, then each new body; a MODIFIED parent is measured on its body only):

| Old requirement | Form | Old length | New blocks (length) |
| --- | --- | --- | --- |
| The module names the one operator release it deploys | c2 (REMOVED plus ADDED under new names) | 569 | The module carries its own version and names one operator image (443); Every rendered controller container runs the image the module names (210) |
| Every object renders through a catalog resource | c2 (REMOVED plus ADDED under new names) | 770 | Operator objects render through the catalog resource made for their kind (486); Administrator ClusterRoles render unbound through the raw-objects resource (373) |
| Names stay those of the earlier manifest | c2 (REMOVED plus ADDED under new names) | 970 | Fixed names: namespaced objects (458); Fixed names: first ClusterRoles (447); Fixed names: other ClusterRoles (483); Fixed names: CRDs (360) |
| The operator's tuning is typed instance values | c2 (REMOVED plus ADDED under new names) | 956 | The operator's tuning values are typed in #config and rendered into the Deployment (495); Defaults reproduce the earlier manifest's Deployment and each resource quantity defaults on its own (381); The memory limit always resolves and sets GOMEMLIMIT (286) |

`grep -rn --exclude-dir=archive` for the old headers found no reference outside this spec. The overlap requirement "The controller's pod and Service stay those of the kustomize tree" is not in this unit (group 9).

Scenarios:

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| The operator version is read without a render | The module carries its own version and names one operator image, The operator version is read without a render | no (byte equal, checked by the item 3 command) |
| The rendered image matches the named release | Every rendered controller container runs the image the module names, The rendered image matches the named release | no (byte equal, checked by the item 3 command) |
| The render contains the whole install shape | Operator objects render through the catalog resource made for their kind, The render contains the whole install shape | no (byte equal, checked by the item 3 command) |
| Raw objects only for the administrator roles | Administrator ClusterRoles render unbound through the raw-objects resource, Raw objects only for the administrator roles | no (byte equal, checked by the item 3 command) |
| Administrator roles render unbound | Administrator ClusterRoles render unbound through the raw-objects resource, Administrator roles render unbound | no (byte equal, checked by the item 3 command) |
| The fixed names render | Fixed names: namespaced objects, The fixed names render | no (byte equal, checked by the item 3 command) |
| Bindings take catalog names | Fixed names: other ClusterRoles, Bindings take catalog names | no (byte equal, checked by the item 3 command) |
| Every value reaches the Deployment | The operator's tuning values are typed in #config and rendered into the Deployment, Every value reaches the Deployment | no (byte equal, checked by the item 3 command) |
| A CPU string the catalog cannot render is refused at the value | The operator's tuning values are typed in #config and rendered into the Deployment, A CPU string the catalog cannot render is refused at the value | no (byte equal, checked by the item 3 command) |
| One resources field keeps the other defaults | Defaults reproduce the earlier manifest's Deployment and each resource quantity defaults on its own, One resources field keeps the other defaults | no (byte equal, checked by the item 3 command) |
| The Go memory limit follows the memory limit | The memory limit always resolves and sets GOMEMLIMIT, The Go memory limit follows the memory limit | no (byte equal, checked by the item 3 command) |
| An unsupported memory unit is refused | The memory limit always resolves and sets GOMEMLIMIT, An unsupported memory unit is refused | no (byte equal, checked by the item 3 command) |

New scenarios (each restates a case the old text names):

| New place | New scenario | Old text that names the case |
| --- | --- | --- |
| Fixed names: first ClusterRoles | The first ClusterRole names render | the sentence(s) of this block |
| Fixed names: other ClusterRoles | The other ClusterRole names render | the sentence(s) of this block |
| Fixed names: CRDs | The CRD names render | the sentence(s) of this block |

Clauses (the rule sentences of each old body and the block that carries each):

| Old clause | New requirement | Note |
| --- | --- | --- |
| The module names the one operator release it deploys #1: The repository SHALL carry the operator's OPM module at the path `opmodel.dev/modules/opm_operator` on the v0 ... | The module carries its own version and names one operator image | verbatim |
| The module names the one operator release it deploys #2: The module's source SHALL name the operator image it deploys by a published operator release's version tag and... | The module carries its own version and names one operator image | verbatim |
| The module names the one operator release it deploys #3: Every rendered controller container SHALL run exactly that image, so one module version deploys exactly one op... | Every rendered controller container runs the image the module names | verbatim |
| Every object renders through a catalog resource #1: The module SHALL render the operator's CRDs, its Namespace, the controller Deployment, its ServiceAccount, the... | Operator objects render through the catalog resource made for their kind | verbatim |
| Every object renders through a catalog resource #2: The administrator ClusterRoles the operator ships for users to bind, its viewer roles among them, SHALL render... | Administrator ClusterRoles render unbound through the raw-objects resource | verbatim |
| Every object renders through a catalog resource #3: Those administrator ClusterRoles SHALL render through the catalog's raw-objects resource, in one component tha... | Administrator ClusterRoles render unbound through the raw-objects resource | verbatim |
| Every object renders through a catalog resource #4: The module SHALL render the operator's Namespace itself, so an install records it as an object of the instance... | Operator objects render through the catalog resource made for their kind | verbatim |
| Names stay those of the earlier manifest #1: The module SHALL render the Namespace `opm-operator-system`, the Deployment `opm-operator-controller-manager`,... | Fixed names: namespaced objects; Fixed names: first ClusterRoles; Fixed names: other ClusterRoles; Fixed names: CRDs | the one long sentence is cut at list boundaries into four parts (verbatim fragments, `The module SHALL render` and the closing clause repeated in each) |
| Names stay those of the earlier manifest #2: Only a role binding MAY take the name the catalog derives from its role. | Fixed names: other ClusterRoles | verbatim |
| The operator's tuning is typed instance values #1: The module's `#config` SHALL type these values and render each into the controller Deployment: the operator im... | The operator's tuning values are typed in #config and rendered into the Deployment | verbatim |
| The operator's tuning is typed instance values #2: Each unset optional value SHALL leave its argument out. | The operator's tuning values are typed in #config and rendered into the Deployment | verbatim |
| The operator's tuning is typed instance values #3: Defaults SHALL reproduce the earlier manifest's Deployment: one replica, requests `100m` CPU and `256Mi` memor... | Defaults reproduce the earlier manifest's Deployment and each resource quantity defaults on its own | verbatim |
| The operator's tuning is typed instance values #4: Each of those four resource quantities SHALL default on its own, so a value that sets one of them keeps the de... | Defaults reproduce the earlier manifest's Deployment and each resource quantity defaults on its own | verbatim |
| The operator's tuning is typed instance values #5: The memory limit SHALL always resolve to a value, default or given, so the container's `GOMEMLIMIT` is always ... | The memory limit always resolves and sets GOMEMLIMIT | verbatim |

Added words (the only ones): subjects and conditions copied from the same old requirement, the exact header of a sibling ("As in requirement" prefix or "(Requirement: ...)" suffix), and the keywords in the Strength table. The one long sentence of "Names stay those of the earlier manifest" is cut at list boundaries into four parts; each part repeats `The module SHALL render` and the closing clause "the names the operator's own install manifest gives them" and names the other three parts. "Only a role binding MAY take the name the catalog derives from its role." keeps its MAY with no added keyword; it sits in the part that holds the other ClusterRoles, with the scenario "Bindings take catalog names".

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
| A ServiceAccount that loses its rights does not release the deletion by itself | It SHALL report `DeletionBlocked` after 10 minutes and SHALL name the cause; `spec.prune=false` releases it, and so does deleting the ServiceAccount. | It reports `DeletionBlocked` after 10 minutes and names the cause; `spec.prune=false` releases it, and so does deleting the ServiceAccount. |
| Moving the opm CLI pin cuts no release | moving it SHALL be a `ci(deps)` commit and SHALL cut no release (workspace RELEASING.md, section "Pin classes") | moving it is a `ci(deps)` commit and cuts no release (workspace RELEASING.md, section "Pin classes") |
| The release-as key repeats its version and crosses a boundary | While the value stays, release-please SHALL propose exactly that version on every Release PR, so a value left in place after its release re-proposes an already-published version. | While the value stays, release-please proposes exactly that version on every Release PR, so a value left in place after its release re-proposes an already-published version. |
| The operator's CUE version moves only through a library release | The operator's CUE version SHALL move only through a library release: the library bumps CUE in its own pull request, ... | ... so the operator's CUE version moves only through a library release: the library bumps CUE in its own pull request, ... (the clause had no keyword of its own; the sentence's SHALL belongs to "SHALL ignore" in the other block) |
| Every reconcile attempt ends with a status patch | On meaningful outcomes (`Applied`, `AppliedAndPruned`, `FailedTransient`, `FailedStalled`), the patch SHALL update conditions, `lastAttempted*`, `lastApplied*` (on success), `inven... | On meaningful outcomes (`Applied`, `AppliedAndPruned`, `FailedTransient`, `FailedStalled`), the patch updates conditions, `lastAttempted*`, `lastApplied*` (on success), `inventory`... |
| A NoOp patch is bounded to a fixed set of fields | On `NoOp`, the patch SHALL be bounded to: drift condition (`Drifted`), the `Healthy` condition (`instance-health`), failure counter deltas (incl. drift counter), `nextRetryAt` clea... | On `NoOp`, the patch is bounded to: drift condition (`Drifted`), the `Healthy` condition (`instance-health`), failure counter deltas (incl. drift counter), `nextRetryAt` clearing, ... |
| A reconcile that skips its render patches only a changed Healthy condition | A reconcile that skips its render because its inputs are unchanged (`render-input-key`) SHALL NOT be an attempt. | A reconcile that skips its render because its inputs are unchanged (`render-input-key`) is not an attempt. |
| A successful outcome patches the Healthy condition it judged | On a successful outcome the patch SHALL also carry the `Healthy` condition the reconcile judged (`instance-health`). | On a successful outcome the patch also carries the `Healthy` condition the reconcile judged (`instance-health`). |
| A successful outcome patches the Healthy condition it judged | A failed attempt SHALL leave `Healthy` as it was. | A failed attempt leaves `Healthy` as it was. |
| requiredContracts is in the NoOp patch set | `requiredContracts` SHALL be in the `NoOp` set deliberately, and it is the one field there that does not describe an outcome. | `requiredContracts` is in the `NoOp` set deliberately, and it is the one field there that does not describe an outcome. |
| A NoOp patch carries requiredContracts because the demand can move without a digest | It SHALL describe what the instance demands, which a reconcile can observe as changed without any digest changing: a regenerated platform re-enqueues every instance, and the render... | It describes what the instance demands, which a reconcile can observe as changed without any digest changing: a regenerated platform re-enqueues every instance, and the render that... |
| lastAppliedVersion is in the NoOp patch set | `lastAppliedVersion` SHALL be in the `NoOp` set because a `NoOp` re-proves it. | `lastAppliedVersion` is in the `NoOp` set because a `NoOp` re-proves it. |
| lastAppliedVersion is in the NoOp patch set | A reconcile that did not render SHALL leave it as it was. | A reconcile that did not render leaves it as it was. |
| lastAppliedInputs is in the NoOp patch set while the drift render interval is positive | `lastAppliedInputs` SHALL be in the `NoOp` set for the same reason: a `NoOp` that rendered proves the cluster holds what those inputs produce, so its key and render time are record... | `lastAppliedInputs` is in the `NoOp` set for the same reason: a `NoOp` that rendered proves the cluster holds what those inputs produce, so its key and render time are recorded, an... |
| lastAppliedInputs is in the NoOp patch set while the drift render interval is positive | It SHALL be written on a `NoOp` only while the drift render interval is greater than zero; with the skip disabled nothing reads it. | It is written on a `NoOp` only while the drift render interval is greater than zero; with the skip disabled nothing reads it. |
| Storm safety comes from the generation predicate | Storm safety SHALL be provided by `WithEventFilter(predicate.GenerationChangedPredicate{})` on the controller. | Storm safety is provided by `WithEventFilter(predicate.GenerationChangedPredicate{})` on the controller. |
| A panicking attempt counts as FailedTransient for counters, history and metrics only | Fields the attempt already wrote on the object before the panic (`requiredContracts`, `instanceUUID`, the `Drifted` condition, a ModulePackage's `status.source`) SHALL be committed... | Fields the attempt already wrote on the object before the panic (`requiredContracts`, `instanceUUID`, the `Drifted` condition, a ModulePackage's `status.source`) are committed as o... |
| A panicking attempt counts as FailedTransient for counters, history and metrics only | A panicking attempt SHALL be counted as `FailedTransient` for failure counters, history and reconcile metrics only. | A panicking attempt is counted as `FailedTransient` for failure counters, history and reconcile metrics only. |
| A panicking attempt counts as FailedTransient for counters, history and metrics only | It SHALL return no `ctrl.Result`: the requeue scenarios of "Outcome classification" and the `nextRetryAt` rule of the reconcile-backoff capability do not apply to it. | It returns no `ctrl.Result`: the requeue scenarios of "Outcome classification" and the `nextRetryAt` rule of the reconcile-backoff capability do not apply to it. |

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
