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

The swarm guide `design/T2.1.md` holds the rules. This file applies
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
c2 leaves it with only the INFO "Archive would refuse"). The human decided
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
| Gate (Taskfile, CI) | Moved to a separate change after the open changes land | unchanged |

**Rationale**: each unit stays independent of the open changes until they land, as the human
ruled. The last column assumes an unchanged inventory; B1, B3 and B5 re-measure first, and `N`
is the fresh cumulative count.
**Alternative rejected**: all units first, with the open changes moved onto the new text
afterwards. It breaks their MODIFIED blocks (guide section 3, "unsafe because silent"), and
editing them is out of scope.

### Reversibility

high (hours). Spec text is reverted by a PR. No contract changes and no accepted
decision record is edited.

## Purpose fixes by direct edit

`openspec archive` ignores a delta Purpose for an existing spec (verified in the CLI's `specs-apply`). A Purpose is therefore a direct edit of `## Purpose` in
`openspec/specs/<spec>/spec.md`, in its own commit, and is recorded here. Unit 1 added a
`## Purpose` above `## Requirements` to `artifact-fetch`, `catalog-registry-resolution`,
`failure-counters`, `history-tracking`, `metrics` and `suspend-resume`. Each text was read from the
requirements of the spec. Gate on `origin/main` 7713cc86 in a scratch copy: the ERROR set lost
those six specs (7 spec items to 1, plus the open change's one); WARNING lines stayed
equal. The remaining spec ERROR, `source-resolution`, waits for its open change.

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

Ledgers written so far: unit 1 (opm-operator-B2), unit 2 (opm-operator-B4), unit 3 (opm-operator-B6) and unit 4a (the first five specs of opm-operator-B7: `serviceaccount-impersonation`, `platform-crd`, `platform-module-generation`, `reconcile-backoff`, `registration-driven-regeneration`), below.

Unit 2 (opm-operator-B4) rules of text, for the five ledgers below:

- Em-dash characters of the old bodies (four places: `module-instance-ownership` "The operator acknowledges a CLI-owned instance with a single condition", `platform-reconciler` "Generate and build the singleton Platform on reconcile", `kernel-module-renderer` "Adapt compiled output to operator resources", twice) are written as `:` or `,` in the new bodies, because the workspace bans the character. The word and the order stay.
- A sentence whose subject was "It" gets the subject of the sentence before it, copied (the "Edit" column shows it).
- Where a body relies on a sibling for a condition, the condition is copied from the old text in front of its sentence (the "Edit" column shows it).
- A sentence over 500 on its own (`Reason constants`, the outcome list of `Surface build outcome on status`, the identity sentence of `Single-slot identity-keyed store for concurrent read`) is cut at a list boundary; the parts are verbatim and the head is repeated.
- `Source:` lines: the whole line stays where it fits with the rule it supports, and bodies are split further until it fits (the enhancement id is never alone, `openspec/config.yaml`). The old `Source:` line covers the whole old requirement, so every split block of that requirement carries it (the library-owns-the-deletion-sequence title for the prune and cleanup deletes; the apply guard title for the Ready message). A block that cannot hold its whole line in 500 characters carries it as a `**Source**:` metadata line, which the length check skips: the Ready message body of `status-conditions` ("The Ready message counts the objects adopted by another instance"), the first AdoptedElsewhere body of `events-emission`, and the definition body of the verdict identity in `ssa-apply`; the same form holds the line of two order-check bodies of `ssa-apply`. Only the citation is in that line.

Unit 3 (opm-operator-B6) rules of text, for the seven ledgers of unit 3:

- A sentence that does not fit with its `Source:` line, or that is over 500 on its own, is cut at a list or clause boundary: every part is a substring of the old sentence, the head is repeated, and the Edit column lists the added punctuation. The cut places are the list of reasons, the list of inputs, the list of citations and the list of outcomes that the old sentence itself held.
- The old `Source:` line covers the whole old requirement, so every block cut from a requirement that has one carries it whole: `render-input-key` (two requirements), `ssa-apply` (three), `modulepackage-reconcile-loop` (two) and `digest-computation` (one). No block of unit 3 uses the short form.
- A list of rules that cannot fit moves into a scenario: the six parts of the render input key and the seven skip conditions. The body keeps the SHALL, repeats the old subject and names the scenario.
- Three bodies gained one RFC 2119 keyword where the old sentence had none; the Strength table lists them. Where a sibling sentence carries the keyword, the sibling is copied as context instead.
- A "Subject" of the old text ("It", "Both fields") is replaced by a subject copied from the same old requirement (the Edit column shows it).

Unit 4a (opm-operator-B7, first five specs; 14 requirements) rules of text, for the five ledgers of unit 4a:

- The header of the `ssa-apply` requirement "A claim whose update the API server refuses fails the whole apply with a named error" is renamed "A claim whose update the API server refuses returns a named error" in its own commit before the 14 requirements (DECISIONS.md 2026-10-20 01:35). Its body says that the apply returns an error that names the claim, the fields and the API server's message; the rule that no object is written is in the requirement "A refused claim update applies, creates and deletes no object of the set". Every exact-header reference in this change (that delta, its migration line, its ledger and clause rows) uses the new header. The six-header identity list of "The identity of the apply verdict is the same for a ModuleInstance and a ModulePackage" does not hold this header.
- A body that relies on a sibling for its subject or condition names the sibling by its exact header, in the form `As in requirement "<header>": <old sentence>` (DECISIONS.md 2026-10-19 16:55). The prefix is the only added text of such a body. It is used where the old text says "this stall", "both", "other", "the annotation" or "the message", and where a sentence has no subject of its own.
- A list whose lead sentence is a head ("the controller MUST:" with bullets, "The Ready condition SHALL reflect the outcome:" with items, "The registry SHALL be resolved from BOTH ...:" with two entries) is cut at an item boundary: the head is repeated and the bullet lines or items are verbatim. The first sentence of `platform-module-generation` "The reconciler generates a platform module from the CR" is over 500 on its own; it is cut at the boundaries of its own list (the `cue.mod/module.cue` path, the core pin with the catalogs, the `platform.cue` entry), each part with the head repeated. The dependency sentence of that requirement is cut at its ";" (the closure, then the module-file source); its roots are named by the exact headers of the two requirements that hold them. A closing "." replaces the "," or ";" that ended the cut item.
- The two em-dashes of `registration-driven-regeneration` and the one of `platform-crd` are written as "," and ":" in the bodies, because the workspace bans the character.
- `reconcile-backoff` keeps both headers (form b), because other specs cite them by exact text (see its ledger). The parent bodies are the sentences that carry the header's rule; every other rule is an ADDED requirement with a restatement scenario in words of the old body.
- Five sentences gained a keyword (Strength table). Where a body already holds a SHALL, a sentence without a keyword stays as it was (`platform-crd`: two bodies).

### Ledger: deps-cascade

Form c2 for every requirement below: REMOVED under the old name, ADDED under new names, new blocks at the end of the spec. `grep -rn --exclude-dir=archive` for each old header found no reference to fix, except the self-reference in the body of "Module release gate", which the new "The module release check verifies the module version" updates, and the Reason text of the open change `stop-operator-install-manifest` (not edited).

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| Up-to-date tree | deps:cascade moves pins in the working tree and reports by exit code, Up-to-date tree | no (byte equal, checked by the item 3 command) |
| A pin is behind | deps:cascade moves pins in the working tree and reports by exit code, A pin is behind | no (byte equal, checked by the item 3 command) |
| Resolver not found | deps:cascade takes every version decision from the shared cascade resolver, Resolver not found | no (byte equal, checked by the item 3 command) |
| A predicate's "no" is not "nothing to do" | deps:cascade moves pins in the working tree and reports by exit code, A predicate's "no" is not "nothing to do" | no (byte equal, checked by the item 3 command) |
| Relative resolver path | deps:cascade takes every version decision from the shared cascade resolver, Relative resolver path | no (byte equal, checked by the item 3 command) |
| Resolver error leaves the tree untouched | deps:cascade resolves every target before it edits, Resolver error leaves the tree untouched | no (byte equal, checked by the item 3 command) |
| Dirty tree refused | deps:cascade refuses a dirty working tree, Dirty tree refused | no (byte equal, checked by the item 3 command) |
| Catalog and core move as a consistent set | deps:cascade moves core to the core its catalog pins, Catalog and core move as a consistent set | no (byte equal, checked by the item 3 command) |
| Core already ahead of the catalog's core | deps:cascade moves core to the core its catalog pins, Core already ahead of the catalog's core | no (byte equal, checked by the item 3 command) |
| Newer core published but not pinned by the catalog | deps:cascade moves core to the core its catalog pins, Newer core published but not pinned by the catalog | no (byte equal, checked by the item 3 command) |
| Sample Platform moved | deps:cascade leaves the generated operator reference alone, Sample Platform moved | no (byte equal, checked by the item 3 command) |
| Tidy adds a dependency | deps:cascade warns about changes beyond the pins, Tidy adds a dependency | no (byte equal, checked by the item 3 command) |
| Third-party pins untouched | deps:cascade runs cue mod get on OPM modules only, Third-party pins untouched | no (byte equal, checked by the item 3 command) |
| Frozen fixture left alone | deps:cascade skips frozen pins, Frozen fixture left alone | no (byte equal, checked by the item 3 command) |
| Frozen catalog keeps its core | deps:cascade skips frozen pins, Frozen catalog keeps its core | no (byte equal, checked by the item 3 command) |
| MVS raises a frozen key | deps:cascade skips frozen pins, MVS raises a frozen key | no (byte equal, checked by the item 3 command) |
| Core hold holds the catalog | deps:cascade honours a hold on core, Core hold holds the catalog | no (byte equal, checked by the item 3 command) |
| Second run does not advance again | Fixture versions advance once per PR, Second run does not advance again | no (byte equal, checked by the item 3 command) |
| Consumers name the advanced version | Fixture consumers follow the advanced version, Consumers name the advanced version | no (byte equal, checked by the item 3 command) |
| Pending version not bumped twice | Fixture versions advance once per PR, Pending version not bumped twice | no (byte equal, checked by the item 3 command) |
| Release files untouched | deps:cascade never modifies release, workflow or shared files, Release files untouched | no (byte equal, checked by the item 3 command) |
| The operator module is left to its own task | deps:cascade leaves the operator module to its own task, The operator module is left to its own task | no (byte equal, checked by the item 3 command) |
| Install page example output drifted | deps:cascade warns about install page drift and leaves the page unchanged, Install page example output drifted | no (byte equal, checked by the item 3 command) |
| Newer CUE language warned | deps:cascade warns when an upstream CUE language is newer, Newer CUE language warned | no (byte equal, checked by the item 3 command) |
| Offline set in PR CI | The offline cascade test set needs no network and runs in the Lint job, Offline set in PR CI | no (byte equal, checked by the item 3 command) |
| Older pins produce the expected diff | The full cascade test set runs apart from required CI, Older pins produce the expected diff | no (byte equal, checked by the item 3 command) |
| Stub drift caught | deps:cascade:test reports each scenario against a stub resolver, Stub drift caught | no (byte equal, checked by the item 3 command) |
| Library and catalog moved | Title and body come from the shared resolver, Library and catalog moved | no (byte equal, checked by the item 3 command) |
| Only fixtures moved | The classes file holds the contract's path classes, Only fixtures moved | no (byte equal, checked by the item 3 command) |
| Only the sample moved, no page regenerated | The classes file holds the contract's path classes, Only the sample moved, no page regenerated | no (byte equal, checked by the item 3 command) |
| Only the opm CLI moved | The classes file holds the contract's path classes, Only the opm CLI moved | no (byte equal, checked by the item 3 command) |
| Pin report agrees between tree and HEAD | pins.sh prints one row per logical pin, Pin report agrees between tree and HEAD | no (byte equal, checked by the item 3 command) |
| Module pins move, image and version stay | deps:cascade:module moves the module's catalog and core, Module pins move, image and version stay | no (byte equal, checked by the item 3 command) |
| Catalog hold holds the module | deps:cascade:module honours holds and frozen paths, Catalog hold holds the module | no (byte equal, checked by the item 3 command) |
| Nothing to move | deps:cascade:module reports by the exit codes of deps:cascade, Nothing to move | no (byte equal, checked by the item 3 command) |
| Catalog published | module-deps.yml opens one pull request for the module's pin move, Catalog published | no (byte equal, checked by the item 3 command) |
| Dry run | module-deps.yml writes only a summary in a dry run, Dry run | no (byte equal, checked by the item 3 command) |
| Pin report keys are unique | The module pull request title and body come from the resolver, Pin report keys are unique | no (byte equal, checked by the item 3 command) |

New scenarios (restate a clause of the old body, with words of the old text only):

| New scenario | Requirement | Old clause restated |
| --- | --- | --- |
| Own registry mapping | deps:cascade sets its own registry environment | the clause of that requirement body |
| Library moved | deps:cascade moves the library pin | the clause of that requirement body |
| Library not moved | deps:cascade moves the library pin | the clause of that requirement body |
| File below the target catalog moves | deps:cascade moves the opm catalog pin to one target | the clause of that requirement body |
| File above the target catalog is not lowered | deps:cascade moves the opm catalog pin to one target | the clause of that requirement body |
| CLI pin written last | deps:cascade writes the opm CLI pin last | the clause of that requirement body |
| Version written only when it differs | Fixture versions are written with the opm CLI | the clause of that requirement body |
| Registries are only read | deps:cascade only reads from registries | the clause of that requirement body |
| Files the module task never edits | deps:cascade:module edits only the module's pins | the clause of that requirement body |
| Offline set covers the module task | The offline test set covers deps:cascade:module | the clause of that requirement body |
| Bot-only branch is rebuilt | The module/deps branch is rebuilt only while it holds bot commits | the clause of that requirement body |
| Human commit on the branch | The module/deps branch is rebuilt only while it holds bot commits | the clause of that requirement body |
| Module pull request touches only the module | The module pull request changes only the module directory | the clause of that requirement body |

Clause table (old clause, new requirement, sentence of its body):

| Old requirement | Old clause | New requirement | Sentence |
| --- | --- | --- | --- |
| deps:cascade moves upstream pins and reports by exit code | provides `task deps:cascade`; moves the upstream pins to the newest published versions in the working tree only | deps:cascade moves pins in the working tree and reports by exit code | 1 |
| deps:cascade moves upstream pins and reports by exit code | SHALL NOT commit, push or create a branch | deps:cascade moves pins in the working tree and reports by exit code | 2 |
| deps:cascade moves upstream pins and reports by exit code | exit 0 changed, 3 nothing to do, other code on error; never turn a failure into 0 or 3 | deps:cascade moves pins in the working tree and reports by exit code | 3 |
| deps:cascade moves upstream pins and reports by exit code | SHALL move exactly these pins (library, opm catalog, core, opm CLI): the closed set | deps:cascade moves pins in the working tree and reports by exit code | 4 |
| deps:cascade moves upstream pins and reports by exit code | every version decision from the shared resolver (RELEASING.md citation) | deps:cascade takes every version decision from the shared cascade resolver | 1 |
| deps:cascade moves upstream pins and reports by exit code | resolver found through absolute `CASCADE_RESOLVER`, else `../.github/.github/scripts/cascade/cascade-resolve.sh` | deps:cascade takes every version decision from the shared cascade resolver | 2 |
| deps:cascade moves upstream pins and reports by exit code | resolver missing: SHALL fail with a message naming both options | deps:cascade takes every version decision from the shared cascade resolver | 3 |
| deps:cascade resolves everything before it edits | SHALL refuse (exit 1) on a changed or untracked tree unless `CASCADE_ALLOW_DIRTY=1` | deps:cascade refuses a dirty working tree | 1 |
| deps:cascade resolves everything before it edits | then judge the result by a before and after snapshot | deps:cascade refuses a dirty working tree | 1 |
| deps:cascade resolves everything before it edits | validate `.cascade-frozen` and `.cascade-hold` through the resolver first | deps:cascade resolves every target before it edits | 1 |
| deps:cascade resolves everything before it edits | resolve every target before editing any file | deps:cascade resolves every target before it edits | 2 |
| deps:cascade resolves everything before it edits | own `CUE_REGISTRY` and `OPM_REGISTRY` value; never inherit the Taskfile's `localhost:5000` | deps:cascade sets its own registry environment | 1 |
| deps:cascade moves library, the catalog and core, and the opm CLI | `go get library@<exact version>` then `go mod tidy`, only when library moved | deps:cascade moves the library pin | 1 |
| deps:cascade moves library, the catalog and core, and the opm CLI | catalog to one target `K` in every file below `K`: sample Platform `version:`, `CatalogVersion()`, four fixture `module.cue` (bare / bare) | deps:cascade moves the opm catalog pin to one target | 1 |
| deps:cascade moves library, the catalog and core, and the opm CLI | a file whose catalog is above `K` SHALL NOT be lowered | deps:cascade moves the opm catalog pin to one target | 2 |
| deps:cascade moves library, the catalog and core, and the opm CLI | core in those modules and the provider fixture to the core the file's catalog pins after the move (higher of file catalog and `K`; `K` for provider), only when greater than the file's own core | deps:cascade moves core to the core its catalog pins | 1 |
| deps:cascade moves library, the catalog and core, and the opm CLI | core SHALL never come from the newest published core directly | deps:cascade moves core to the core its catalog pins | 2 |
| deps:cascade moves library, the catalog and core, and the opm CLI | opm CLI: `.opm-cli-version` written last | deps:cascade writes the opm CLI pin last | 1 |
| deps:cascade moves library, the catalog and core, and the opm CLI | SHALL NOT edit `docs/site/reference/operator/_index.md`; docs bundle generates it from `config/samples` | deps:cascade leaves the generated operator reference alone | 1 |
| deps:cascade moves library, the catalog and core, and the opm CLI | `cue mod get` names only `opmodel.dev/*` and `testing.opmodel.dev/*`, exact version; runs with one `cue mod tidy` only where a pin moved | deps:cascade runs cue mod get on OPM modules only | 1 |
| deps:cascade moves library, the catalog and core, and the opm CLI | third-party pins SHALL never be named | deps:cascade runs cue mod get on OPM modules only | 2 |
| deps:cascade moves library, the catalog and core, and the opm CLI | never edit an import path or a `@vN` key; new major only as the resolver's warning | deps:cascade runs cue mod get on OPM modules only | 3 |
| deps:cascade moves library, the catalog and core, and the opm CLI | warnings: third-party pin changed by `tidy`; `go`/`toolchain` directive; dep a fixture gains that its modulepackage lacks | deps:cascade warns about changes beyond the pins | 1 |
| deps:cascade honours frozen and held pins | ask the resolver whether the exact file is frozen for the pin key before editing a literal or naming a `cue mod get` key | deps:cascade skips frozen pins | 1 |
| deps:cascade honours frozen and held pins | skip that file for that key when frozen | deps:cascade skips frozen pins | 2 |
| deps:cascade honours frozen and held pins | after `cue mod tidy` every frozen key byte-unchanged, else exit 1 naming file and key | deps:cascade skips frozen pins | 3 |
| deps:cascade honours frozen and held pins | hold on `opmodel.dev/core@v2` with `max` below the target catalog's core keeps catalog and core, with a warning | deps:cascade honours a hold on core | 1 |
| Fixture versions advance once per PR and consumers follow | per fixture module and the provider fixture, set the declared version in `identity/identity.cue`; `B` = declared version at the merge-base of `CASCADE_BASE` (default `origin/main`) and `HEAD` | Fixture versions advance once per PR | 1 |
| Fixture versions advance once per PR and consumers follow | `B` when unchanged; next patch of `B` when changed and `B` published; `B` when changed and not yet published | Fixture versions advance once per PR | 1 |
| Fixture versions advance once per PR and consumers follow | written with `opm module version set` or `opm catalog version set`, CLI installed from the unmodified `.opm-cli-version` under the git directory, only when it differs | Fixture versions are written with the opm CLI | 1 |
| Fixture versions advance once per PR and consumers follow | consumers follow in the same run: modulepackage `module.cue`, `moduleinstance.yaml`, the sample ModuleInstance | Fixture consumers follow the advanced version | 1 |
| deps:cascade never touches release, workflow or shared files | SHALL NOT modify: `.cascade-*`; release-please files and `CHANGELOG.md`; `.github/`; docs-kit files; `hack/fixtures.sh`, `.tasks/examples.yaml`; CUE `language.version`; Jellyfin sample, `ocirepository.yaml` | deps:cascade never modifies release, workflow or shared files | 1 |
| deps:cascade never touches release, workflow or shared files | SHALL NOT modify anything under `modules/opm_operator/`, moved by `deps:cascade:module` in its own PR | deps:cascade leaves the operator module to its own task | 1 |
| deps:cascade never touches release, workflow or shared files | only read from registries; never publish or seed | deps:cascade only reads from registries | 1 |
| deps:cascade never touches release, workflow or shared files | warn when an upstream `language.version` is newer than `CUE_VERSION` in `test.yml` | deps:cascade warns when an upstream CUE language is newer | 1 |
| deps:cascade never touches release, workflow or shared files | warn, and not edit, when the install page example output names another catalog or core | deps:cascade warns about install page drift and leaves the page unchanged | 1 |
| deps:cascade is tested offline in required CI and online on demand | `deps:cascade:test` runs in throwaway copies against the stub resolver; `PASS`/`FAIL` per scenario | deps:cascade:test reports each scenario against a stub resolver | 1 |
| deps:cascade is tested offline in required CI and online on demand | exit 0 only when every scenario passes | deps:cascade:test reports each scenario against a stub resolver | 2 |
| deps:cascade is tested offline in required CI and online on demand | fails when the stub's `sha256sum` differs from the contract's | deps:cascade:test reports each scenario against a stub resolver | 3 |
| deps:cascade is tested offline in required CI and online on demand | offline set: no GHCR or module proxy beyond a warm Go cache; runs in `Lint` (the required check) | The offline cascade test set needs no network and runs in the Lint job | 1 |
| deps:cascade is tested offline in required CI and online on demand | SHALL NOT need the real resolver or a `.github` checkout, except for the title and body check | The offline cascade test set needs no network and runs in the Lint job | 2 |
| deps:cascade is tested offline in required CI and online on demand | full set in a separate, non-required workflow on cascade-file changes, by hand, weekly | The full cascade test set runs apart from required CI | 1 |
| deps:cascade is tested offline in required CI and online on demand | full set includes older pins with the version-advance diff, frozen pins, title and body against the real resolver | The full cascade test set runs apart from required CI | 2 |
| Title and body come from the shared resolver and the contract's path classes | `deps:cascade:title` and `deps:cascade:body` call the resolver's `title` and `body` with `classes` and `pins.sh` | Title and body come from the shared resolver | 1 |
| Title and body come from the shared resolver and the contract's path classes | `classes`: release-tool, test classes, the contract's opm-operator block and nothing more; other paths shipped | The classes file holds the contract's path classes | 1 |
| Title and body come from the shared resolver and the contract's path classes | `classes`: any other path, including `go.mod` and `go.sum`, SHALL be shipped | The classes file holds the contract's path classes | 2 |
| Title and body come from the shared resolver and the contract's path classes | `pins.sh <ref>`: one row per logical pin, five tab-separated fields; the four rows | pins.sh prints one row per logical pin | 1 |
| deps:cascade:module moves only the operator module's catalog and core | moves the module's pins in `modules/opm_operator/cue.mod/module.cue` and nothing else, with the fixture rules: catalog to `K` when below, core to the catalog's core when greater | deps:cascade:module moves the module's catalog and core | 1 |
| deps:cascade:module moves only the operator module's catalog and core | the module's pins are shipped pins | deps:cascade:module moves the module's catalog and core | 2 |
| deps:cascade:module moves only the operator module's catalog and core | holds and frozen paths under the same pin keys; a hold on either holds the module; a `.cascade-frozen` entry for the module freezes only the module | deps:cascade:module honours holds and frozen paths | 1 |
| deps:cascade:module moves only the operator module's catalog and core | same exit codes as `deps:cascade` | deps:cascade:module reports by the exit codes of deps:cascade | 1 |
| deps:cascade:module moves only the operator module's catalog and core | SHALL NOT edit outside `modules/opm_operator/`, nor `operator/operator.cue`, `identity.Version`, generated CRD/RBAC, `CHANGELOG.md`, `RELEASE`; the four reasons | deps:cascade:module edits only the module's pins | 1 |
| deps:cascade:module moves only the operator module's catalog and core | the offline test set SHALL cover it | The offline test set covers deps:cascade:module | 1 |
| The operator module's pin move opens its own pull request | `module-deps.yml` runs `task -x deps:cascade:module` on `repository_dispatch` `upstream-released` and `workflow_dispatch` | module-deps.yml opens one pull request for the module's pin move | 1 |
| The operator module's pin move opens its own pull request | when the tree changed: push `module/deps` with the release App's token; open or update one PR so its CI runs; the dry-run condition of the old last sentence ("the repository variable `CASCADE_DRY_RUN` is not exactly `false`") copied in front as "Unless ... " so the body keeps its exception | module-deps.yml opens one pull request for the module's pin move | 2 |
| The operator module's pin move opens its own pull request | bot-only branch rebuilt from `main`; human commit: add own commit, never rewrite | The module/deps branch is rebuilt only while it holds bot commits | 1 |
| The operator module's pin move opens its own pull request | title and body from the resolver with `module-classes` (all shipped) and `module-pins.sh` (two keyed rows); each key appears once | The module pull request title and body come from the resolver | 1 |
| The operator module's pin move opens its own pull request | PR changes no file outside `modules/opm_operator/`, so its squash commit proposes a module release and no operator release | The module pull request changes only the module directory | 1 |
| The operator module's pin move opens its own pull request | while `CASCADE_DRY_RUN` is not exactly `false`: write the diff to the job summary, push nothing | module-deps.yml writes only a summary in a dry run | 1 |

Lengths (old body, new bodies; characters by the guide's rule):

| Old requirement | Old | New |
| --- | --- | --- |
| deps:cascade moves upstream pins and reports by exit code | 779 | 446, 426 |
| deps:cascade resolves everything before it edits | 606 | 239, 155, 255 |
| deps:cascade moves library, the catalog and core, and the opm CLI | 1680 | 170, 352, 429, 92, 165, 370, 228 |
| deps:cascade honours frozen and held pins | 538 | 383, 189 |
| Fixture versions advance once per PR and consumers follow | 978 | 483, 246, 378 |
| deps:cascade never touches release, workflow or shared files | 898 | 373, 167, 84, 134, 233 |
| deps:cascade is tested offline in required CI and online on demand | 910 | 296, 345, 318 |
| Title and body come from the shared resolver and the contract's path classes | 801 | 171, 276, 370 |
| deps:cascade:module moves only the operator module's catalog and core | 1238 | 482, 267, 110, 414, 88 |
| The operator module's pin move opens its own pull request | 1182 | 395, 210, 448, 163, 156 |

### Ledger: operator-module-release

Form c2 for every requirement below: REMOVED under the old name, ADDED under new names, new blocks at the end of the spec. `grep -rn --exclude-dir=archive` for each old header found no reference to fix, except the self-reference in the body of "Module release gate", which the new "The module release check verifies the module version" updates, and the Reason text of the open change `stop-operator-install-manifest` (not edited).

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| A module-only change releases only the module | A commit releases only the package whose files it changes, A module-only change releases only the module | no (byte equal, checked by the item 3 command) |
| An operator-only change releases only the operator | A commit releases only the package whose files it changes, An operator-only change releases only the operator | no (byte equal, checked by the item 3 command) |
| Module tag shape | The module is a second release-please package, Module tag shape | no (byte equal, checked by the item 3 command) |
| The operator's tag listing ignores module tags | Module release tags cannot collide with operator tags, The operator's tag listing ignores module tags | no (byte equal, checked by the item 3 command) |
| A breaking change raises the minor | Module versions on the v0 path are 0.y.z, A breaking change raises the minor | no (byte equal, checked by the item 3 command) |
| A feature raises the patch | Module versions on the v0 path are 0.y.z, A feature raises the patch | no (byte equal, checked by the item 3 command) |
| The first release | The first module release is 0.1.0, The first release | no (byte equal, checked by the item 3 command) |
| Hidden module commits open no release | The first module release is 0.1.0, Hidden module commits open no release | no (byte equal, checked by the item 3 command) |
| 1.0.0 before operator GA is refused | The module moves to the v1 path only after operator GA, 1.0.0 before operator GA is refused | no (byte equal, checked by the item 3 command) |
| Release PR gains the identity commit | The identity-advance step is the only writer of the module's identity version, Release PR gains the identity commit | no (byte equal, checked by the item 3 command) |
| Re-run converges | The identity-advance step is idempotent, Re-run converges | no (byte equal, checked by the item 3 command) |
| A failed identity advance strands no release | The identity-advance step runs in a job of its own, A failed identity advance strands no release | no (byte equal, checked by the item 3 command) |
| Image digest does not match GHCR | The module release check verifies the operator image reference, Image digest does not match GHCR | no (byte equal, checked by the item 3 command) |
| Image tag is a draft release | The module release check verifies the operator image reference, Image tag is a draft release | no (byte equal, checked by the item 3 command) |
| CRDs ahead of the deployed operator | The module release check compares generated CRDs and RBAC with the deployed operator, CRDs ahead of the deployed operator | no (byte equal, checked by the item 3 command) |
| Operator below the minimum | The module release check refuses an operator below the minimum version, Operator below the minimum | no (byte equal, checked by the item 3 command) |
| Development pin | The module release check refuses a development pin or a local module, Development pin | no (byte equal, checked by the item 3 command) |
| Several failures named at once | An offline test covers the module release checks that need no network, Several failures named at once | no (byte equal, checked by the item 3 command) |
| Not a release PR | The module release check runs in Lint on the module's release PR, Not a release PR | no (byte equal, checked by the item 3 command) |
| First publish | The module publish job publishes the tagged source, First publish | no (byte equal, checked by the item 3 command) |
| Re-run after a later step failed | The module publish job is idempotent for the same release, Re-run after a later step failed | no (byte equal, checked by the item 3 command) |
| Read back where CUE published | The module publish job reads registry references from cue mod resolve, Read back where CUE published | no (byte equal, checked by the item 3 command) |
| Version held by other content | The module publish job is idempotent for the same release, Version held by other content | no (byte equal, checked by the item 3 command) |
| Manifest equals the module's render | The install manifest is the module's render at its defaults, Manifest equals the module's render | no (byte equal, checked by the item 3 command) |
| Module debug values are not used | The install manifest is the module's render at its defaults, Module debug values are not used | no (byte equal, checked by the item 3 command) |
| Namespace first | install.yaml is attached to the draft release with the Namespace and CRDs first, Namespace first | no (byte equal, checked by the item 3 command) |
| Upload fails after the publish | The render and the upload run apart from the publish, Upload fails after the publish | no (byte equal, checked by the item 3 command) |
| Manifest names another image | The install manifest names the operator image the module names, Manifest names another image | no (byte equal, checked by the item 3 command) |
| Fresh cluster | The install manifest yields a running operator, Fresh cluster | no (byte equal, checked by the item 3 command) |
| Asset missing | The final module publish job publishes the draft only with install.yaml attached, Asset missing | no (byte equal, checked by the item 3 command) |
| Never the latest release | A module release is never the latest release, Never the latest release | no (byte equal, checked by the item 3 command) |
| Already published | The final module publish job publishes the draft only with install.yaml attached, Already published | no (byte equal, checked by the item 3 command) |
| Recover after a broken read-back | The recovery dispatch publishes, attaches and releases in order, Recover after a broken read-back | no (byte equal, checked by the item 3 command) |
| Stranded before its first publish | The recovery dispatch publishes, attaches and releases in order, Stranded before its first publish | no (byte equal, checked by the item 3 command) |
| Held version differs from the tag | The recovery dispatch publishes, attaches and releases in order, Held version differs from the tag | no (byte equal, checked by the item 3 command) |
| Release already published | The recovery dispatch refuses before any registry login, Release already published | no (byte equal, checked by the item 3 command) |
| Not a module tag, or not main | The recovery dispatch refuses before any registry login, Not a module tag, or not main | no (byte equal, checked by the item 3 command) |
| Search for publishers | Only the module publish job publishes under the module path, Search for publishers | no (byte equal, checked by the item 3 command) |
| Fork or dispatch on another ref | The module publish job refuses outside the module release, Fork or dispatch on another ref | no (byte equal, checked by the item 3 command) |
| Operator release published | The module image PR sets the deployed operator release and nothing else, Operator release published | no (byte equal, checked by the item 3 command) |
| Breaking operator release | The module image PR title carries the change type, Breaking operator release | no (byte equal, checked by the item 3 command) |
| Draft operator release | The module image workflow accepts a dispatch for recovery, Draft operator release | no (byte equal, checked by the item 3 command) |
| Tag behind main | The module image run fails when config differs from the operator tag, Tag behind main | no (byte equal, checked by the item 3 command) |
| Older operator tag | The module image run fails for a tag that does not descend from the deployed one, Older operator tag | no (byte equal, checked by the item 3 command) |
| Image PR never releases the operator | The module image PR sets the deployed operator release and nothing else, Image PR never releases the operator | no (byte equal, checked by the item 3 command) |

New scenarios (restate a clause of the old body, with words of the old text only):

| New scenario | Requirement | Old clause restated |
| --- | --- | --- |
| Publish asserts the version | No release-please updater writes the module's identity package | the clause of that requirement body |
| Identity version differs from the proposed version | The module release check verifies the module version | the clause of that requirement body |
| Release frozen while config differs | Module releases are frozen while config differs from the deployed tag | the clause of that requirement body |
| Release merged past a red check | The module release check runs in the publish job before any push | the clause of that requirement body |
| Operator-only release | The module publish job runs only for a module release | the clause of that requirement body |
| Offline test asserts the references | An offline test asserts the registry references of the publish job | the clause of that requirement body |
| Upload confirms a single draft | A module release stays a draft with one tag until it is recovered or published | the clause of that requirement body |
| Dispatch runs the module jobs only | The release workflow accepts a module recovery dispatch | the clause of that requirement body |
| Fixed script runs against the tag | The recovery dispatch runs main's scripts against the tag's tree | the clause of that requirement body |
| Module job gates explicitly | Each recovery module job gates on its predecessor explicitly | the clause of that requirement body |
| Dispatch has its own concurrency group | The recovery dispatch runs in a concurrency group per tag | the clause of that requirement body |
| Push run starts module jobs only after a module release | A push-triggered release runs the tag's own scripts | the clause of that requirement body |
| Recovery leaves tags and releases alone | Recovery never moves a tag or edits a release by hand | the clause of that requirement body |
| Local flow publishes away from GHCR | Local and test flows publish the module away from GHCR | the clause of that requirement body |
| Image PR created with the release App's token | An operator release opens one pull request on module/operator-image | the clause of that requirement body |
| Bot-only image branch is rebuilt | The module/operator-image branch is rebuilt only while it holds bot commits | the clause of that requirement body |
| Human commit on the image branch | The module/operator-image branch is rebuilt only while it holds bot commits | the clause of that requirement body |

Clause table (old clause, new requirement, sentence of its body):

| Old requirement | Old clause | New requirement | Sentence |
| --- | --- | --- | --- |
| The module is a release unit of its own | `release-please-config.json` declares `modules/opm_operator` as a second package beside `"."` | The module is a second release-please package | 1 |
| The module is a release unit of its own | own version, changelog, release PR, tags `opm_operator-vX.Y.Z` | The module is a second release-please package | 2 |
| The module is a release unit of its own | no operator tag or `release/vX.Y` branch can take that form | Module release tags cannot collide with operator tags | 1 |
| The module is a release unit of its own | only release-please creates these tags; no workflow moves one | Module release tags cannot collide with operator tags | 2 |
| The module is a release unit of its own | operator package excludes the module directory; module-only commit proposes no operator release | A commit releases only the package whose files it changes | 1 |
| The module is a release unit of its own | a commit with no file under the module directory proposes no module release | A commit releases only the package whose files it changes | 2 |
| The module is a release unit of its own | module versions do not follow the operator's or the cli's | The module is a second release-please package | 3 |
| The module version follows the 0.x bump rule on the v0 path | definition of a breaking module release (three conditions) | Module versions on the v0 path are 0.y.z | 2 |
| The module version follows the 0.x bump rule on the v0 path | on `@v0` every version `0.y.z`; breaking raises `y` and resets `z`; other releases raise `z`, `feat` included | Module versions on the v0 path are 0.y.z | 1 |
| The module version follows the 0.x bump rule on the v0 path | first module release `0.1.0` | The first module release is 0.1.0 | 1 |
| The module version follows the 0.x bump rule on the v0 path | move to `@v1` only with `1.0.0` no earlier than the operator's first GA | The module moves to the v1 path only after operator GA | 1 |
| The module version follows the 0.x bump rule on the v0 path | gate refuses `1.0.0` or higher while the deployed operator release is a prerelease | The module moves to the v1 path only after operator GA | 1 |
| The module's identity version has one writer | `identity.Version` written only by the identity-advance step; module published equals tagged source | The identity-advance step is the only writer of the module's identity version | 1 |
| The module's identity version has one writer | step runs `opm module version set`, commits the named message, pushes to the release PR branch | The identity-advance step is the only writer of the module's identity version | 2 |
| The module's identity version has one writer | no `x-release-please-version` or other updater writes the identity package; `publish --version` only asserts | No release-please updater writes the module's identity package | 1 |
| The module's identity version has one writer | step idempotent, no change when the branch already declares the version | The identity-advance step is idempotent | 1 |
| The module's identity version has one writer | own job after release-please; failure never fails the job whose outputs start the release jobs | The identity-advance step runs in a job of its own | 1 |
| The module's identity version has one writer | holds the release App token only for its push | The identity-advance step runs in a job of its own | 1 |
| Module release gate | required `Lint` job runs a release check that fails, naming every failure before it exits | The module release check runs in Lint on the module's release PR | 1 |
| Module release gate | fails on `-0.dev.` pin or `cue.mod/local-module.cue` | The module release check refuses a development pin or a local module | 1 |
| Module release gate | fails when the image reference is not `...opm-operator:<tag>@<digest>` with a published non-draft tag and the GHCR manifest-list digest | The module release check verifies the operator image reference | 1 |
| Module release gate | fails when `identity.Version` differs from the proposed version, or its major disagrees with the path's | The module release check verifies the module version | 1 |
| Module release gate | fails when the proposed version breaks the v0 rules (reference now names the three requirements) | The module release check verifies the module version | 1 |
| Module release gate | fails when `<tag>` is older than the minimum operator version (first release that refuses to reconcile the operator's own instance) | The module release check refuses an operator below the minimum version | 1 |
| Module release gate | fails when generated CRDs or RBAC differ from the operator source at `<tag>` | The module release check compares generated CRDs and RBAC with the deployed operator | 1 |
| Module release gate | rationale: module must render exactly the CRDs and permissions of the controller it deploys; `main` can move ahead | The module release check compares generated CRDs and RBAC with the deployed operator | 2 |
| Module release gate | releases are frozen while the generated CRDs or RBAC under `main`'s `config/` differ from the deployed tag's, until the image-bump PR merges (stated as the effect of the CRD and RBAC check; the check SHALL fail in that state) | Module releases are frozen while config differs from the deployed tag | 1 |
| Module release gate | check also runs in the publish job before anything is pushed | The module release check runs in the publish job before any push | 1 |
| Module release gate | no-network checks covered by an offline test run by `Lint` on every PR | An offline test covers the module release checks that need no network | 1 |
| The module publishes its tagged source once | publish job checks out the tag, installs the opm CLI, runs the release check, runs `opm module publish` | The module publish job publishes the tagged source | 1 |
| The module publishes its tagged source once | job runs only when the module package created a release; never for an operator-only release | The module publish job runs only for a module release | 1 |
| The module publishes its tagged source once | re-run succeeds without pushing when the registry holds the same content | The module publish job is idempotent for the same release | 1 |
| The module publishes its tagged source once | fails when the registry holds the version with other content | The module publish job is idempotent for the same release | 2 |
| The module publishes its tagged source once | does not rely on the registry to refuse an overwrite | The module publish job is idempotent for the same release | 3 |
| The module publishes its tagged source once | registry references from `cue mod resolve`, never hand-spelled (worked example) | The module publish job reads registry references from cue mod resolve | 1 |
| The module publishes its tagged source once | an offline test asserts those references | An offline test asserts the registry references of the publish job | 1 |
| Every module release publishes the install manifest rendered from it | render at default values, instance `opm-operator`, namespace `opm-operator-system`, platform from the module's pins, no cluster | The install manifest is the module's render at its defaults | 1 |
| Every module release publishes the install manifest rendered from it | write `install.yaml`, Namespace and CRDs first; attach to the draft release | install.yaml is attached to the draft release with the Namespace and CRDs first | 1 |
| Every module release publishes the install manifest rendered from it | render and upload in a job apart from the publish | The render and the upload run apart from the publish | 1 |
| Every module release publishes the install manifest rendered from it | manifest names the module's tag and digest; the job (old text: "a job apart from the publish") checks before upload; the body names that job in the words of the old text | The install manifest names the operator image the module names | 1 |
| Every module release publishes the install manifest rendered from it | `kubectl apply --server-side` on a cluster with no operator yields a running operator | The install manifest yields a running operator | 1 |
| Every module release publishes the install manifest rendered from it | a PR changing the module or the render scripts proves this on kind | The install manifest yields a running operator | 2 |
| Module releases are drafts until every asset is attached | module package sets `draft: true` and `force-tag-creation: true` | A module release stays a draft with one tag until it is recovered or published | 1 |
| Module releases are drafts until every asset is attached | every upload first confirms exactly one release carries the tag and it is a draft | A module release stays a draft with one tag until it is recovered or published | 2 |
| Module releases are drafts until every asset is attached | recovery before publication; after publication next version; published tag never moved | A module release stays a draft with one tag until it is recovered or published | 3 |
| Module releases are drafts until every asset is attached | final job depends on the publish job and confirms `install.yaml` before publishing | The final module publish job publishes the draft only with install.yaml attached | 1 |
| Module releases are drafts until every asset is attached | not a prerelease; published with `make_latest=false`; enhancement 0021 citation | A module release is never the latest release | 1 |
| A stranded module draft is recovered by dispatch | accepts `workflow_dispatch` with `module_tag` for a module release whose draft is unpublished | The release workflow accepts a module recovery dispatch | 1 |
| A stranded module draft is recovered by dispatch | runs only the module jobs, never release-please, identity advance or an operator job | The release workflow accepts a module recovery dispatch | 2 |
| A stranded module draft is recovered by dispatch | refuses before any registry login unless: repo and `main`; tag regex; tag exists and reachable from `main`; exactly one draft release | The recovery dispatch refuses before any registry login | 1 |
| A stranded module draft is recovered by dispatch | checks the tag's tree out; runs `main`'s release scripts; gate and opm CLI stay the tag's | The recovery dispatch runs main's scripts against the tag's tree | 1 |
| A stranded module draft is recovered by dispatch | each module job gates explicitly on its predecessor | Each recovery module job gates on its predecessor explicitly | 1 |
| A stranded module draft is recovered by dispatch | in order: publish (reuse only when equal; publish when not held), install manifest, publication with `make_latest=false` | The recovery dispatch publishes, attaches and releases in order | 1 |
| A stranded module draft is recovered by dispatch | concurrency group per tag | The recovery dispatch runs in a concurrency group per tag | 1 |
| A stranded module draft is recovered by dispatch | push-triggered release runs the tag's own scripts, starts module jobs only after a module release | A push-triggered release runs the tag's own scripts | 1 |
| A stranded module draft is recovered by dispatch | recovery never moves, deletes or re-creates a tag, never edits a release by hand | Recovery never moves a tag or edits a release by hand | 1 |
| Only this repository's module release publishes under the module path | no other workflow, task or script publishes to the module path or its GHCR repository | Only the module publish job publishes under the module path | 1 |
| Only this repository's module release publishes under the module path | publish job refuses unless the repository and a module release tag | The module publish job refuses outside the module release | 1 |
| Only this repository's module release publishes under the module path | local and test flows publish only to a registry that maps `opmodel.dev` away from GHCR | Local and test flows publish the module away from GHCR | 1 |
| An operator release opens the PR that moves the module's image | open or update one PR on `module/operator-image` after a published operator release | An operator release opens one pull request on module/operator-image | 1 |
| An operator release opens the PR that moves the module's image | created with the release App's token so CI runs | An operator release opens one pull request on module/operator-image | 2 |
| An operator release opens the PR that moves the module's image | PR sets tag and manifest-list digest; changes nothing else (no file outside the module, no generated CRD/RBAC) | The module image PR sets the deployed operator release and nothing else | 1 |
| An operator release opens the PR that moves the module's image | CRDs or RBAC under `config/` at the tag differ from `main`: fail naming files, no PR, because the gate would refuse | The module image run fails when config differs from the operator tag | 1 |
| An operator release opens the PR that moves the module's image | fail, no PR, when `<tag>` does not descend from the deployed tag | The module image run fails for a tag that does not descend from the deployed one | 1 |
| An operator release opens the PR that moves the module's image | title `fix(deps): deploy operator <tag> ...` | The module image PR title carries the change type | 1 |
| An operator release opens the PR that moves the module's image | title `fix(deps)!: ...` on a breaking CHANGELOG entry or a CRD that stops serving a version | The module image PR title carries the change type | 2 |
| An operator release opens the PR that moves the module's image | `workflow_dispatch` with `tag` does the same for a published release, for recovery | The module image workflow accepts a dispatch for recovery | 1 |
| An operator release opens the PR that moves the module's image | bot-only branch rebuilt from `main`; human commit: add on top, never rewrite | The module/operator-image branch is rebuilt only while it holds bot commits | 1 |

Lengths (old body, new bodies; characters by the guide's rule):

| Old requirement | Old | New |
| --- | --- | --- |
| The module is a release unit of its own | 848 | 440, 244, 227 |
| The module version follows the 0.x bump rule on the v0 path | 762 | 475, 42, 243 |
| The module's identity version has one writer | 867 | 390, 181, 116, 217 |
| Module release gate | 1843 | 134, 154, 272, 357, 232, 392, 339, 154, 245 |
| The module publishes its tagged source once | 1152 | 277, 129, 299, 442, 97 |
| Every module release publishes the install manifest rendered from it | 894 | 255, 173, 134, 161, 248 |
| Module releases are drafts until every asset is attached | 809 | 397, 134, 276 |
| A stranded module draft is recovered by dispatch | 1644 | 241, 281, 267, 181, 399, 133, 158, 92 |
| Only this repository's module release publishes under the module path | 519 | 262, 151, 104 |
| An operator release opens the PR that moves the module's image | 1519 | 211, 298, 216, 195, 437, 163, 208 |

### Ledger: events-emission

Form c2 for every requirement below: REMOVED under the old name, ADDED under new names, new blocks at the end of the spec. `grep -rn --exclude-dir=archive` for each old header found no reference to fix, except the self-reference in the body of "Module release gate", which the new "The module release check verifies the module version" updates, and the Reason text of the open change `stop-operator-install-manifest` (not edited).

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| Skew under Warn is reported once | A render with advisory findings emits one RenderWarning event per finding, Skew under Warn is reported once | no (byte equal, checked by the item 3 command) |
| Unchanged warnings do not repeat | RenderWarning events are emitted only when the findings change, Unchanged warnings do not repeat | no (byte equal, checked by the item 3 command) |
| Rewording does not re-emit | RenderWarning text is authored by the operator and keyed on facts, Rewording does not re-emit | no (byte equal, checked by the item 3 command) |
| Kept claims on prune | A ClaimsKept event is emitted when claims are kept, Kept claims on prune | no (byte equal, checked by the item 3 command) |
| Kept claims on deletion | A ClaimsKept event is emitted when claims are kept, Kept claims on deletion | no (byte equal, checked by the item 3 command) |
| Nothing kept, no event | No ClaimsKept event is emitted when nothing was kept, Nothing kept, no event | no (byte equal, checked by the item 3 command) |
| More than ten kept claims | The ClaimsKept message states the count, the claims and how to delete one, More than ten kept claims | no (byte equal, checked by the item 3 command) |
| Refused claim recreate | A ClaimConflict event is emitted when a claim recreate is refused, Refused claim recreate | no (byte equal, checked by the item 3 command) |
| Claim kept by the delete guard | The ClaimConflict event names the claim and carries the Ready text, Claim kept by the delete guard | no (byte equal, checked by the item 3 command) |
| An object of another instance left on prune | A LeftBehind event is emitted when objects are left in the cluster, An object of another instance left on prune | no (byte equal, checked by the item 3 command) |
| Only a Namespace left on deletion | The LeftBehind event type depends on why objects were left, Only a Namespace left on deletion | no (byte equal, checked by the item 3 command) |
| A Namespace and an adopted object left on deletion | The LeftBehind message states the count and the library messages, A Namespace and an adopted object left on deletion | no (byte equal, checked by the item 3 command) |
| A failed prune reports nothing | A failed prune or cleanup emits no LeftBehind event, A failed prune reports nothing | no (byte equal, checked by the item 3 command) |
| Nothing left, no event | A LeftBehind event is emitted when objects are left in the cluster, Nothing left, no event | no (byte equal, checked by the item 3 command) |
| A waiting deletion reports what it left once | A LeftBehind event is emitted only by a committed run, A waiting deletion reports what it left once | no (byte equal, checked by the item 3 command) |
| An object is let go | An AdoptedElsewhere event is emitted when the count of adopted objects changes, whether the reconcile applies, restores or does nothing, An object is let go | no (byte equal, checked by the item 3 command) |
| No event while nothing changes | No AdoptedElsewhere event while the count is unchanged, No event while nothing changes | no (byte equal, checked by the item 3 command) |
| A second object is let go | The AdoptedElsewhere event states the count and the library messages, A second object is let go | no (byte equal, checked by the item 3 command) |
| Nothing adopted elsewhere, no event | An AdoptedElsewhere event is emitted when the count of adopted objects changes, whether the reconcile applies, restores or does nothing, Nothing adopted elsewhere, no event | no (byte equal, checked by the item 3 command) |
| Twenty unset required values | A stalled render event message that does not fit keeps whole findings, Twenty unset required values | no (byte equal, checked by the item 3 command) |
| A message that fits | A stalled render failure event carries at most 1024 characters, A message that fits | no (byte equal, checked by the item 3 command) |
| A long message without findings | A stalled render event message without findings is cut and counted, A long message without findings | no (byte equal, checked by the item 3 command) |
| One event when the wait starts | A DeletionInProgress event is emitted when a deletion wait starts, One event when the wait starts | no (byte equal, checked by the item 3 command) |
| One warning when the wait is blocked | A DeletionBlocked event is emitted when a deletion wait becomes blocked, One warning when the wait is blocked | no (byte equal, checked by the item 3 command) |
| Released after the ServiceAccount went | A DeletionUnconfirmed event is emitted when a deletion is released unread, Released after the ServiceAccount went | no (byte equal, checked by the item 3 command) |
| A confirmed deletion emits no such event | A DeletionUnconfirmed event is emitted when a deletion is released unread, A confirmed deletion emits no such event | no (byte equal, checked by the item 3 command) |

New scenarios (restate a clause of the old body, with words of the old text only):

| New scenario | Requirement | Old clause restated |
| --- | --- | --- |
| Absent object and kept claim | No LeftBehind event for an absent object or a kept claim | the clause of that requirement body |
| AdoptedElsewhere event on a no-op reconcile | An AdoptedElsewhere event is emitted when the count of adopted objects changes, whether the reconcile applies, restores or does nothing | the clause of that requirement body |
| Message of a released deletion | The DeletionUnconfirmed message states the count, the reason and the doubt | the clause of that requirement body |

Clause table (old clause, new requirement, sentence of its body):

| Old requirement | Old clause | New requirement | Sentence |
| --- | --- | --- | --- |
| Render warnings are emitted as events on transition | render with advisory findings: one Warning event per distinct finding, reason `RenderWarning`, action `Render` | A render with advisory findings emits one RenderWarning event per finding | 1 |
| Render warnings are emitted as events on transition | a render with no findings emits none | A render with advisory findings emits one RenderWarning event per finding | 2 |
| Render warnings are emitted as events on transition | emit only when the set of findings changes between reconciles | RenderWarning events are emitted only when the findings change | 1 |
| Render warnings are emitted as events on transition | event text authored by the operator from the advisory rows | RenderWarning text is authored by the operator and keyed on facts | 1 |
| Render warnings are emitted as events on transition | transition check keyed on the facts the rows carry (em-dash pair turned into parentheses) | RenderWarning text is authored by the operator and keyed on facts | 2 |
| Events emitted when claims are kept | one Normal `ClaimsKept` event when a prune or a deletion cleanup keeps PVCs | A ClaimsKept event is emitted when claims are kept | 1 |
| Events emitted when claims are kept | states the count; names claims `<namespace>/<name>`, at most ten, fewer past 1024 chars, then the rest; says no longer tracked; says how to delete | The ClaimsKept message states the count, the claims and how to delete one | 1 |
| Events emitted when claims are kept | `action` `Prune` on the stale-prune path, `Delete` on the deletion path | A ClaimsKept event is emitted when claims are kept | 2 |
| Events emitted when claims are kept | on the deletion path emitted before the finalizer is removed | A ClaimsKept event is emitted when claims are kept | 3 |
| Events emitted when claims are kept | text MUST NOT carry an enhancement reference | The ClaimsKept message states the count, the claims and how to delete one | 2 |
| Events emitted when claims are kept | no event when nothing was kept: prune not true; `dataPolicy` `Delete`; no live owned claim | No ClaimsKept event is emitted when nothing was kept | 1 |
| Events emitted when a claim recreate is refused | one Warning `ClaimConflict`, action `Apply`, when an apply refused the forced recreate of a PVC | A ClaimConflict event is emitted when a claim recreate is refused | 1 |
| Events emitted when a claim recreate is refused | names the claim `<namespace>/<name>`; same text as `Ready`; field named or not | The ClaimConflict event names the claim and carries the Ready text | 1 |
| Events emitted when a claim recreate is refused | MUST NOT also emit `ApplyFailed` | A ClaimConflict event is emitted when a claim recreate is refused | 2 |
| Events emitted when a claim recreate is refused | text MUST NOT carry an enhancement reference | The ClaimConflict event names the claim and carries the Ready text | 2 |
| Events emitted when objects are left behind | one event `LeftBehind` when a prune or a deletion cleanup leaves objects the delete verdict skipped | A LeftBehind event is emitted when objects are left in the cluster | 1 |
| Events emitted when objects are left behind | type `Normal` for Namespaces and CRDs OPM never deletes; `Warning` for unmanaged, other-instance, adopted-elsewhere | The LeftBehind event type depends on why objects were left | 1 |
| Events emitted when objects are left behind | states the count; library skip message unchanged per object; at most ten, 1024 limit, rest count; ownership reasons first | The LeftBehind message states the count and the library messages | 1 |
| Events emitted when objects are left behind | `action` `Prune` / `Delete` | A LeftBehind event is emitted when objects are left in the cluster | 2 |
| Events emitted when objects are left behind | text MUST NOT carry an enhancement reference | The LeftBehind message states the count and the library messages | 2 |
| Events emitted when objects are left behind | emitted only by a committed run: prune returned no error; first reconcile reaching a release verdict | A LeftBehind event is emitted only by a committed run | 1 |
| Events emitted when objects are left behind | a later reconcile of the same wait MUST NOT emit again | A LeftBehind event is emitted only by a committed run | 2 |
| Events emitted when objects are left behind | a prune or cleanup that failed for another entry MUST NOT emit it; the retry judges again | A failed prune or cleanup emits no LeftBehind event | 1 |
| Events emitted when objects are left behind | none for an absent object or a kept PVC, which has its own event | No LeftBehind event for an absent object or a kept claim | 1 |
| Events emitted when an object is adopted by another instance | sentence 1, whole and unchanged: the `Warning` event with reason `AdoptedElsewhere` and action `Apply`, the whole old condition and the citation of `status-conditions` | An AdoptedElsewhere event is emitted when the count of adopted objects changes, whether the reconcile applies, restores or does nothing | 1 | none |
| Events emitted when an object is adopted by another instance | the citation `status-conditions`, "The Ready message counts the objects adopted by another instance" (copied from sentence 1 of the old body into the no-event requirement, where the Ready message carries the state) | No AdoptedElsewhere event while the count is unchanged | 1 | none |
| Events emitted when an object is adopted by another instance | "This holds whether the reconcile applies, restores or is a no-op." (sentence 2, copied unchanged, after sentence 1 in the same body) | An AdoptedElsewhere event is emitted when the count of adopted objects changes, whether the reconcile applies, restores or does nothing | 2 | none |
| Events emitted when an object is adopted by another instance | states the count; the library's message for each object unchanged; within the `ApplyRefused` limits (limits restated) | The AdoptedElsewhere event states the count and the library messages | 1 |
| Events emitted when an object is adopted by another instance | same number as before, or a skipped render, emits none; `Ready` carries the state between events | No AdoptedElsewhere event while the count is unchanged | 1 |
| Events emitted when an object is adopted by another instance | `Source:` line (enhancement 0012): whole, as a `**Source**:` line on the first requirement, and inline on the other two | An AdoptedElsewhere event is emitted when the count of adopted objects changes, whether the reconcile applies, restores or does nothing | 3 | none |
| The event of a stalled render failure stays inside the event note limit | stalled render or package load Warning event: message at most 1024 chars | A stalled render failure event carries at most 1024 characters | 1 |
| The event of a stalled render failure stays inside the event note limit | when the failure's message fits: the `Ready` message | A stalled render failure event carries at most 1024 characters | 2 |
| The event of a stalled render failure stays inside the event note limit | findings overflow: text in front plus whole findings, end `; and <N> more findings` | A stalled render event message that does not fit keeps whole findings | 1 |
| The event of a stalled render failure stays inside the event note limit | any other message that does not fit: cut between two characters, end ` ... (<N> more characters)` | A stalled render event message without findings is cut and counted | 1 |
| The event of a stalled render failure stays inside the event note limit | `Ready` keeps the whole message | A stalled render failure event carries at most 1024 characters | 3 |
| Events emitted while a deletion waits | one Normal `DeletionInProgress`, action `Delete`, when the cleanup starts to wait | A DeletionInProgress event is emitted when a deletion wait starts | 1 |
| Events emitted while a deletion waits | one Warning `DeletionBlocked`, action `Delete`, when the wait becomes blocked | A DeletionBlocked event is emitted when a deletion wait becomes blocked | 1 |
| Events emitted while a deletion waits | each event carries the message of the condition | A DeletionInProgress event is emitted when a deletion wait starts | 2 |
| Events emitted while a deletion waits | MUST NOT repeat either while `Ready` already carries the reason | A DeletionInProgress event is emitted when a deletion wait starts | 3 |
| Events emitted while a deletion waits | a deletion blocked at its first reconcile emits `DeletionBlocked` only (old sentence, copied; it is in both bodies, so neither states its rule without the exception) | A DeletionBlocked event is emitted when a deletion wait becomes blocked; A DeletionInProgress event is emitted when a deletion wait starts | 4 |
| Events emitted when a deletion is released without confirmation | one Warning `DeletionUnconfirmed`, action `Delete`, in the reconcile that removes the finalizer without reading every object (two cases) | A DeletionUnconfirmed event is emitted when a deletion is released unread | 1 |
| Events emitted when a deletion is released without confirmation | message states the number of unread objects, says why, says they may still exist | The DeletionUnconfirmed message states the count, the reason and the doubt | 1 |
| Events emitted when a deletion is released without confirmation | emitted before the finalizer is removed | A DeletionUnconfirmed event is emitted when a deletion is released unread | 2 |
| Events emitted when a deletion is released without confirmation | text MUST NOT carry an enhancement reference | The DeletionUnconfirmed message states the count, the reason and the doubt | 2 |

Lengths (old body, new bodies; characters by the guide's rule):

| Old requirement | Old | New |
| --- | --- | --- |
| Render warnings are emitted as events on transition | 768 | 269, 155, 429 |
| Events emitted when claims are kept | 845 | 319, 221, 333 |
| Events emitted when a claim recreate is refused | 545 | 234, 326 |
| Events emitted when objects are left behind | 1509 | 271, 281, 362, 380, 146, 149 |
| Events emitted when an object is adopted by another instance | 969 | 475, 492, 467 |
| The event of a stalled render failure stays inside the event note limit | 714 | 358, 278, 154 |
| Events emitted while a deletion waits | 507 | 393, 382 |
| Events emitted when a deletion is released without confirmation | 606 | 366, 261 |

### Ledger: status-conditions

Form c2 for every requirement below except where marked b: REMOVED under the old name, ADDED under new names, new blocks at the end of the spec. Form b for "The Ready message counts the objects adopted by another instance": the parent stays (MODIFIED, shorter body, both of its scenarios in place) and the split-off sentence is ADDED under a new name, because `events-emission` quotes the parent header. `grep -rn --exclude-dir=archive` for each old header found no reference to fix, except the quotation in `events-emission` (kept valid by form b). Old scenario text is copied by script from `origin/main`; item 3 of the guide checklist shows `same` for each.

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| Reason constants available | The existing reason constants stay defined, Reason constants available | no (byte equal, checked by the item 3 command) |
| Skew refusal reason is available | The SkewRefused reason constant is defined, Skew refusal reason is available | no (byte equal, checked by the item 3 command) |
| Duplicate-identity reason is available | The DuplicateIdentities reason constant is defined, Duplicate-identity reason is available | no (byte equal, checked by the item 3 command) |
| Reconcile-panic reason is available | The ReconcilePanic reason constant is defined, Reconcile-panic reason is available | no (byte equal, checked by the item 3 command) |
| Render-timeout reason is available | The RenderTimedOut reason constant is defined, Render-timeout reason is available | no (byte equal, checked by the item 3 command) |
| The constants are available | The Healthy condition type and its reasons are defined, The constants are available | no (byte equal, checked by the item 3 command) |
| Marking Ready leaves Healthy alone | Healthy is independent of Ready, Marking Ready leaves Healthy alone | no (byte equal, checked by the item 3 command) |
| The helper sets only Healthy | One helper sets Healthy and no other condition, The helper sets only Healthy | no (byte equal, checked by the item 3 command) |
| A refused claim on a ModuleInstance | A ClaimConflict reconcile does not stall, prune or record, A refused claim on a ModuleInstance | no (byte equal, checked by the item 3 command) |
| The conflict clears when the render fits the claim again | A ClaimConflict reconcile does not stall, prune or record, The conflict clears when the render fits the claim again | no (byte equal, checked by the item 3 command) |
| Delete recreates the claim | A refused claim recreate sets the ClaimConflict reason, Delete recreates the claim | no (byte equal, checked by the item 3 command) |
| A claim kept by the delete guard | The ClaimConflict message after the delete guard kept the claim, A claim kept by the delete guard | no (byte equal, checked by the item 3 command) |
| The reason on a refused second change | The IdentityChangeUnsettled message names the way out, The reason on a refused second change | no (byte equal, checked by the item 3 command) |
| The reason on a refused apply | The ApplyRefused message carries the library messages, The reason on a refused apply | no (byte equal, checked by the item 3 command) |
| A let-go object does not set the reason | A refused apply sets the ApplyRefused reason, A let-go object does not set the reason | no (byte equal, checked by the item 3 command) |
| Waiting is not stalled | DeletionInProgress is set while the cleanup waits, Waiting is not stalled | no (byte equal, checked by the item 3 command) |
| Blocked is stalled and names the way out | DeletionBlocked is set when an object has not gone for 10 minutes, Blocked is stalled and names the way out | no (byte equal, checked by the item 3 command) |
| A long list is cut | A deletion wait message lists at most ten objects, A long list is cut | no (byte equal, checked by the item 3 command) |
| The finalizer that holds the object is named first | A deletion wait message names the finalizer that holds the object first, The finalizer that holds the object is named first | no (byte equal, checked by the item 3 command) |
| Only the collector's finalizer | A deletion wait message names the finalizer that holds the object first, Only the collector's finalizer | no (byte equal, checked by the item 3 command) |
| The count after a let-go | The Ready message counts the objects adopted by another instance, The count after a let-go | no (byte equal, checked by the item 3 command) |
| The count goes away | The Ready message counts the objects adopted by another instance, The count goes away | no (byte equal, checked by the item 3 command) |

New scenarios (restate a clause of the old body, with words of the old text only):

| New scenario | Requirement | Old clause restated |
| --- | --- | --- |
| A claim refused by the check before the apply | The ClaimConflict message names the refused field when the check refused the claim | the clause of that requirement body |
| A ClaimConflict reconcile prunes and records nothing | A ClaimConflict reconcile does not stall, prune or record | the clause of that requirement body |
| The IdentityChangeUnsettled message carries no enhancement reference | A refused second identity change sets the IdentityChangeUnsettled reason | the clause of that requirement body |
| A skipped render keeps the Ready message | A reconcile that skips its render leaves the Ready message as written | the clause of that requirement body |
| The message states the count and that nothing was applied | The ApplyRefused message states the count | the clause of that requirement body |
| The deletion reasons are defined | The deletion wait reasons are defined | the clause of that requirement body |
| An object that is not being deleted carries neither reason | The deletion wait reasons are defined | the clause of that requirement body |

Clause table (old clause, new requirement, sentence of its body; "Edit" lists the words that differ from the old sentence):

| Old requirement | Old clause | New requirement | Sentence | Edit |
| --- | --- | --- | --- | --- |
| Reason constants | list item `SkewRefused` (Ready=False: the platform's skew policy is ...) | The SkewRefused reason constant is defined | 1 | list cut at an item boundary; the sentence head "The status package SHALL define the reason constants used on the Ready condition, including" repeated; none else |
| Reason constants | list item `DuplicateIdentities` (... enhancement 0015 ("One registration per module")) | The DuplicateIdentities reason constant is defined | 1 | list cut at an item boundary; head repeated; the citation stays in the item; none else |
| Reason constants | list item `ReconcilePanic` (Ready=False, Reconciling=True: ...) | The ReconcilePanic reason constant is defined | 1 | list cut at an item boundary; head repeated; none else |
| Reason constants | list item `RenderTimedOut` (Ready=False, Reconciling=True, not stalled: ...) | The RenderTimedOut reason constant is defined | 1 | list cut at an item boundary; head repeated; none else |
| Reason constants | "beside the existing `ResolutionFailed`, `RenderFailed`, `PlatformNotReady` and the Platform reasons `Generated`, `BuildFailed`, `GenerateFailed`" | The existing reason constants stay defined | 1 | head repeated, the including-list left out; none else |
| The Healthy condition vocabulary | The `internal/status` package SHALL define the condition type `Healthy... | The Healthy condition type and its reasons are defined | 1 | none |
| The Healthy condition vocabulary | It SHALL provide one helper that sets `Healthy` from a status, a reaso... | One helper sets Healthy and no other condition | 1 | -It +The `internal/status` package |
| The Healthy condition vocabulary | `Healthy` is independent of `Ready`: no helper that sets `Ready`, `Rec... | Healthy is independent of Ready | 1 | -is +SHALL be |
| The ClaimConflict reason | The status package SHALL define the reason constant `ClaimConflict`. | A refused claim recreate sets the ClaimConflict reason | 1 | none |
| The ClaimConflict reason | A ModuleInstance or ModulePackage reconcile whose apply refused the fo... | A refused claim recreate sets the ClaimConflict reason | 2 | none |
| The ClaimConflict reason | A ModuleInstance or ModulePackage reconcile whose apply refused the fo... | The ClaimConflict message names the refused field when the c... | 1 | none |
| The ClaimConflict reason | A ModuleInstance or ModulePackage reconcile whose apply refused the fo... | The ClaimConflict message after the delete guard kept the cl... | 1 | none |
| The ClaimConflict reason | A ModuleInstance or ModulePackage reconcile whose apply refused the fo... | A ClaimConflict reconcile does not stall, prune or record | 1 | none |
| The ClaimConflict reason | The message MUST name the claim as `<namespace>/<name>`, MUST say that... | A refused claim recreate sets the ClaimConflict reason | 3 | none |
| The ClaimConflict reason | When the check before the apply refused the claim, the message MUST al... | The ClaimConflict message names the refused field when the c... | 2 | none |
| The ClaimConflict reason | When the resource manager's delete guard kept the claim (the claim cha... | The ClaimConflict message after the delete guard kept the cl... | 2 | none |
| The ClaimConflict reason | The reconcile MUST NOT set `Stalled`: it MUST retry on the bounded bac... | A ClaimConflict reconcile does not stall, prune or record | 2 | none |
| The ClaimConflict reason | The reconcile MUST NOT prune, MUST NOT record a new inventory and MUST... | A ClaimConflict reconcile does not stall, prune or record | 3 | none |
| The IdentityChangeUnsettled reason | The status package SHALL define the reason constant `IdentityChangeUns... | A refused second identity change sets the IdentityChangeUnse... | 1 | none |
| The IdentityChangeUnsettled reason | The status package SHALL define the reason constant `IdentityChangeUns... | The IdentityChangeUnsettled message names the way out | 1 | none |
| The IdentityChangeUnsettled reason | A ModuleInstance or ModulePackage reconcile that refuses a second chan... | A refused second identity change sets the IdentityChangeUnse... | 2 | none |
| The IdentityChangeUnsettled reason | A ModuleInstance or ModulePackage reconcile that refuses a second chan... | The IdentityChangeUnsettled message names the way out | 2 | none |
| The IdentityChangeUnsettled reason | The message MUST say that an earlier change of the instance identity i... | The IdentityChangeUnsettled message names the way out | 3 | none |
| The IdentityChangeUnsettled reason | The message MUST NOT carry an enhancement reference. | A refused second identity change sets the IdentityChangeUnse... | 3 | none |
| The ApplyRefused reason | The status package SHALL define the reason constant `ApplyRefused`. | A refused apply sets the ApplyRefused reason | 1 | none |
| The ApplyRefused reason | A ModuleInstance or ModulePackage reconcile that the apply verdict ref... | A refused apply sets the ApplyRefused reason | 2 | none |
| The ApplyRefused reason | The status package SHALL define the reason constant `ApplyRefused`. | The ApplyRefused message states the count | 1 | none (context copied; "that reason" of sentence 2 keeps its referent) |
| The ApplyRefused reason | A ModuleInstance or ModulePackage reconcile that the apply verdict ref... | The ApplyRefused message states the count | 2 | none |
| The ApplyRefused reason | The message MUST state how many objects were refused and that nothing was applied, | The ApplyRefused message states the count | 3 | "," at the end becomes "." (list cut at ", and MUST carry") |
| The ApplyRefused reason | The status package SHALL define the reason constant `ApplyRefused`. | The ApplyRefused message carries the library messages | 1 | none (context copied) |
| The ApplyRefused reason | A ModuleInstance or ModulePackage reconcile that the apply verdict ref... | The ApplyRefused message carries the library messages | 2 | none |
| The ApplyRefused reason | and MUST carry, for each refused object, the message the library words... | The ApplyRefused message carries the library messages | 3 | -and +The message (subject repeated; list cut) |
| The ApplyRefused reason | The message MUST NOT carry an enhancement reference. | A refused apply sets the ApplyRefused reason | 3 | none |
| The Ready message counts the objects ado... | When a ModuleInstance or ModulePackage reconcile renders, ends `Ready=... | The Ready message counts the objects adopted by another inst... | 1 | none |
| The Ready message counts the objects ado... | When it found none, the message MUST NOT mention such objects. | The Ready message counts the objects adopted by another inst... | 2 | none |
| The Ready message counts the objects ado... | A reconcile that skips its render MUST leave the message as the last r... | A reconcile that skips its render leaves the Ready message a... | 1 |  +of the `Ready` condition |
| The Ready message counts the objects ado... | The message MUST NOT carry an enhancement reference. | The Ready message counts the objects adopted by another inst... | 3 | none |
| The Ready message counts the objects ado... | Source: enhancement 0012 ("The apply guard runs on every apply, refusi... | A reconcile that skips its render leaves the Ready message a... | 2 | none |
| The Ready message counts the objects ado... | Source: enhancement 0012 ("The apply guard runs on every apply, refusi... | The Ready message counts the objects adopted by another inst... | 4 | the whole line as a `**Source**:` line, with the old title (the id was alone before) |
| The deletion wait reasons | The `internal/status` package MUST define the reasons `DeletionInProgr... | The deletion wait reasons are defined | 1 | none |
| The deletion wait reasons | `DeletionInProgress` MUST be set with `Ready=False` and `Reconciling=T... | DeletionInProgress is set while the cleanup waits | 1 | none |
| The deletion wait reasons | `DeletionInProgress` MUST be set with `Ready=False` and `Reconciling=T... | DeletionBlocked is set when an object has not gone for 10 mi... | 1 | none |
| The deletion wait reasons | `DeletionInProgress` MUST be set with `Ready=False` and `Reconciling=T... | A deletion wait message lists at most ten objects | 1 | none |
| The deletion wait reasons | `DeletionBlocked` MUST be set with `Ready=False` and `Stalled=True`: a... | DeletionBlocked is set when an object has not gone for 10 mi... | 2 | none |
| The deletion wait reasons | `DeletionBlocked` MUST be set with `Ready=False` and `Stalled=True`: a... | A deletion wait message lists at most ten objects | 2 | none |
| The deletion wait reasons | Both messages MUST name the remaining objects and their finalizers; th... | DeletionBlocked is set when an object has not gone for 10 mi... | 3 | none |
| The deletion wait reasons | Both messages MUST name the remaining objects and their finalizers; th... | A deletion wait message lists at most ten objects | 3 | none |
| The deletion wait reasons | A message MUST list at most ten objects and say how many more remain. | A deletion wait message lists at most ten objects | 4 | none |
| The deletion wait reasons | A message MUST list at most ten objects and say how many more remain. | A deletion wait message names the finalizer that holds the o... | 1 | none |
| The deletion wait reasons | For each object it MUST name at most three finalizers and count the re... | A deletion wait message names the finalizer that holds the o... | 2 | none |
| The deletion wait reasons | When `foregroundDeletion` is the only finalizer, the message MUST say... | A deletion wait message names the finalizer that holds the o... | 3 | none |
| The deletion wait reasons | Neither reason MUST appear on an object that is not being deleted. | The deletion wait reasons are defined | 2 | none |


Lengths (old body, new bodies; characters by the guide's rule):

| Old requirement | Old | New |
| --- | --- | --- |
| Reason constants | 962 | 229, 316, 265, 277, 226 |
| The Healthy condition vocabulary | 598 | 295, 145, 189 |
| The ClaimConflict reason | 1081 | 388, 326, 401, 452 |
| The IdentityChangeUnsettled reason | 520 | 320, 467 |
| The ApplyRefused reason | 566 | 263, 293, 438 |
| The deletion wait reasons | 1155 | 295, 170, 405, 475, 453 |
| The Ready message counts the objects adopted by another instance | 746 | 451, 319 |

### Ledger: kubernetes-tier-adoption

Form c2 for every requirement below: REMOVED under the old name, ADDED under new names, new blocks at the end of the spec. `grep -rn --exclude-dir=archive` for each old header found no reference to fix. Old scenario text is copied by script from `origin/main`; item 3 of the guide checklist shows `same` for each.

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| No local copy | The operator takes the Kubernetes object code from the library and declares no copy, No local copy | no (byte equal, checked by the item 3 command) |
| Lint refuses the old path | Lint refuses an import of the operator's earlier object copy, Lint refuses the old path | no (byte equal, checked by the item 3 command) |
| Digest and apply objects from one export | The reconciler exports a render once and takes the digest, the objects and the entries from it, Digest and apply objects from one export | no (byte equal, checked by the item 3 command) |
| Export failure keeps its reason | A failure of the export is reported by the conversion, Export failure keeps its reason | no (byte equal, checked by the item 3 command) |
| Decode failure keeps its reason | A failure of the export is reported by the conversion, Decode failure keeps its reason | no (byte equal, checked by the item 3 command) |
| Inventory entries from the same export | The reconciler exports a render once and takes the digest, the objects and the entries from it, Inventory entries from the same export | no (byte equal, checked by the item 3 command) |
| No local copy | The operator takes the Kubernetes object code from the library and declares no copy, No local copy | no (byte equal, checked by the item 3 command) |
| A copy in the conversion package fails the tests | A test keeps the inventory conversion package closed, A copy in the conversion package fails the tests | no (byte equal, checked by the item 3 command) |
| No local evaluator | The Healthy condition code declares no readiness rule of its own, No local evaluator | no (byte equal, checked by the item 3 command) |
| The prune's outcome is the verdict's | A prune or cleanup delete follows a delete verdict of the library's ownership package, The prune's outcome is the verdict's | no (byte equal, checked by the item 3 command) |
| A new delete outside the listed places fails the tests | A test keeps the places that may delete a cluster object closed, A new delete outside the listed places fails the tests | no (byte equal, checked by the item 3 command) |
| A method of the same name on another type is not a delete | A test finds a delete by the type of the receiver, A method of the same name on another type is not a delete | no (byte equal, checked by the item 3 command) |
| The plan's order and propagation reach the API server unchanged | The operator builds each deletion plan with the library's constructor, The plan's order and propagation reach the API server unchanged | no (byte equal, checked by the item 3 command) |
| The forced recreate is not a plan step | The forced recreate delete carries the UID precondition, The forced recreate is not a plan step | no (byte equal, checked by the item 3 command) |
| The apply's outcome is the verdict's | A write follows an apply verdict of the library's ownership package, The apply's outcome is the verdict's | no (byte equal, checked by the item 3 command) |
| A new write outside the listed places fails the tests | A test keeps the places that may write a cluster object closed, A new write outside the listed places fails the tests | no (byte equal, checked by the item 3 command) |
| No write of the adopt annotation | The operator never writes the adopt annotation, No write of the adopt annotation | no (byte equal, checked by the item 3 command) |

New scenarios (restate a clause of the old body, with words of the old text only):

| New scenario | Requirement | Old clause restated |
| --- | --- | --- |
| The label code comes from the library | The operator takes the OPM label code from the library and declares no copy | the clause of that requirement body |
| The stale set and the digests come from the library | The operator computes the inventory stale set and the digests only with the library | the clause of that requirement body |
| The operator passes the fetched object to the library | The operator judges the readiness of an inventory object only with the library's health package | the clause of that requirement body |
| The operator's own resources keep their own condition checks | The health package covers applied inventory objects only | the clause of that requirement body |
| The install admission input is not set | The operator sets no install admission input in a delete verdict | the clause of that requirement body |
| A delete is an action of the lifecycle transition | Prune and cleanup deletes are lifecycle actions | the clause of that requirement body |
| The forced recreate delete asks no verdict | The forced recreate delete asks no verdict | the clause of that requirement body |
| One function performs the plan's reads and deletes | One function of the operator performs the deletion plan | the clause of that requirement body |
| No ownership rule of the operator's own | The operator declares no ownership rule of its own for an apply | the clause of that requirement body |
| A write is found by the type of the receiver | A test finds a write by the type of the receiver | the clause of that requirement body |

Clause table (old clause, new requirement, sentence of its body; "Edit" lists the words that differ from the old sentence):

| Old requirement | Old clause | New requirement | Sentence | Edit |
| --- | --- | --- | --- | --- |
| The operator uses the library's object a... | The operator SHALL take the Kubernetes resource wrapper over the kerne... | The operator takes the OPM label code from the library and d... | 1 | -the Kubernetes resource wrapper over the kernel's compiled output, its JSON and unstructured conversion, the duplicate-identity check, |
| The operator uses the library's object a... | It SHALL declare no type, constant or function of its own for any of t... | The operator takes the Kubernetes object code from the libra... | 2 | -It +The operator |
| The operator uses the library's object a... | It SHALL declare no type, constant or function of its own for any of t... | The operator takes the OPM label code from the library and d... | 2 | -It +The operator |
| The operator uses the library's object a... | Its lint SHALL refuse an import of the package path that held its earl... | Lint refuses an import of the operator's earlier object copy | 1 | -Its +The operator's |
| The operator uses the library's object a... | Source: enhancement 0012 ("The Kubernetes runtime decisions live in a... | The operator takes the Kubernetes object code from the libra... | 3 | none |
| The operator uses the library's object a... | Source: enhancement 0012 ("The Kubernetes runtime decisions live in a... | The operator takes the OPM label code from the library and d... | 3 | none |
| The operator uses the library's object a... | Source: enhancement 0012 ("The Kubernetes runtime decisions live in a... | Lint refuses an import of the operator's earlier object copy | 2 | none |
| The reconciler converts a render with on... | The reconciler SHALL convert a render's resources with one call to the... | The reconciler exports a render once and takes the digest, t... | 1 | none |
| The reconciler converts a render with on... | It SHALL take the render digest, the unstructured objects it applies a... | The reconciler exports a render once and takes the digest, t... | 2 | -It +The reconciler |
| The reconciler converts a render with on... | It SHALL drop the rendered CUE values as soon as the export returns, s... | The reconciler exports a render once and takes the digest, t... | 3 | -It +The reconciler |
| The reconciler converts a render with on... | No other code path SHALL export a rendered resource. | The reconciler exports a render once and takes the digest, t... | 4 | none |
| The reconciler converts a render with on... | A failure of the export SHALL be reported by the conversion: a value t... | A failure of the export is reported by the conversion | 1 | none |
| The reconciler converts a render with on... | With the renderer's export gone, the conversion is the first place suc... | A failure of the export is reported by the conversion | 2 | none |
| The operator uses the library's inventor... | The operator SHALL compute the inventory stale set, the inventory dige... | The operator computes the inventory stale set and the digest... | 1 | none |
| The operator uses the library's inventor... | It SHALL build inventory entries with that package's `NewEntry`. | The operator computes the inventory stale set and the digest... | 2 | -It +The operator |
| The operator uses the library's inventor... | It SHALL declare no entry constructor, identity relation, stale-set fu... | The operator declares no inventory code beyond the conversio... | 1 | -It +The operator |
| The operator uses the library's inventor... | Its only inventory code is the conversion between the API's `Inventory...` | The operator declares no inventory code beyond the conversion | 2 | none (cut at ", and a test") |
| The operator uses the library's inventor... | a test of `internal/inventory` SHALL refuse any further export ... unnoticed | A test keeps the inventory conversion package closed | 1 | cut at ", and a test"; "a test" to "A test" |
| The operator uses the library's inventor... | Source: enhancement 0012 ("The Kubernetes runtime decisions live in a... | The operator computes the inventory stale set and the digest... | 3 | none |
| The operator uses the library's inventor... | Source: enhancement 0012 ("The Kubernetes runtime decisions live in a... | The operator declares no inventory code beyond the conversio... | 3 | none |
| The operator uses the library's inventor... | Source: enhancement 0012 ("The Kubernetes runtime decisions live in a... | A test keeps the inventory conversion package closed | 2 | none |
| The operator judges the health of invent... | The operator SHALL judge the readiness of an object in `status.invento... | The operator judges the readiness of an inventory object onl... | 1 | none |
| The operator judges the health of invent... | The code that computes the `Healthy` condition SHALL read no status fi... | The Healthy condition code declares no readiness rule of its... | 1 | none |
| The operator judges the health of invent... | The operator fetches each object with its own client and passes it in,... | The operator judges the readiness of an inventory object onl... | 2 | none |
| The operator judges the health of invent... | This requirement covers applied inventory objects only; the operator's... | The health package covers applied inventory objects only | 1 |  +SHALL |
| The operator judges the health of invent... | Source: enhancement 0012 ("The Kubernetes runtime decisions live in a... | The operator judges the readiness of an inventory object onl... | 3 | none |
| The operator judges the health of invent... | Source: enhancement 0012 ("The Kubernetes runtime decisions live in a... | The Healthy condition code declares no readiness rule of its... | 2 | none |
| The operator judges the health of invent... | Source: enhancement 0012 ("The Kubernetes runtime decisions live in a... | The health package covers applied inventory objects only | 2 | none |
| The operator decides prune and cleanup d... | Every delete of a cluster object that the operator makes in a prune of... | A prune or cleanup delete follows a delete verdict of the li... | 1 | none |
| The operator decides prune and cleanup d... | Every delete of a cluster object that the operator makes in a prune of... | The operator sets no install admission input in a delete ver... | 1 | none |
| The operator decides prune and cleanup d... | Every delete of a cluster object that the operator makes in a prune of... | The forced recreate delete asks no verdict | 1 | none (context copied) |
| The operator decides prune and cleanup d... | The operator MUST take the kinds that are never deleted from the same... | A prune or cleanup delete follows a delete verdict of the li... | 2 | none |
| The operator decides prune and cleanup d... | The operator MUST NOT set the verdict's install admission input. | The operator sets no install admission input in a delete ver... | 2 | none |
| The operator decides prune and cleanup d... | Source: enhancement 0012 ("The library owns the whole deletion sequenc... | A prune or cleanup delete follows a delete verdict of the li... | 3 | none |
| The operator decides prune and cleanup d... | Source: enhancement 0012 ("The library owns the whole deletion sequenc... | The operator sets no install admission input in a delete ver... | 3 | none |
| The operator decides prune and cleanup d... | Source: enhancement 0012 ("The library owns the whole deletion sequenc... | A test keeps the places that may delete a cluster object clo... | 2 | none |
| The operator decides prune and cleanup d... | Source: enhancement 0012 ("The library owns the whole deletion sequenc... | A test finds a delete by the type of the receiver | 2 | none |
| The operator decides prune and cleanup d... | Each such delete MUST be an action the library's `opm/k8s/lifecycle` transition named, ... | Prune and cleanup deletes are lifecycle actions | 1 | "Each such delete" replaced by "Every delete of a cluster object that the operator makes in a prune of stale resources or in a deletion cleanup, for a ModuleInstance or a ModulePackage," copied from the first sentence |
| The operator decides prune and cleanup d... | The operator MUST build each plan with the library's constructor and M... | The operator builds each deletion plan with the library's constructor | 1 | none |
| The operator decides prune and cleanup d... | The operator MUST build each plan with the library's constructor and M... | One function of the operator performs the deletion plan | 1 | none |
| The operator decides prune and cleanup d... | One function of the operator MUST perform the plan's reads and deletes... | One function of the operator performs the deletion plan | 2 | none |
| The operator decides prune and cleanup d... | Source: enhancement 0012 ("The library owns the whole deletion sequenc... | Prune and cleanup deletes are lifecycle actions | 2 | none |
| The operator decides prune and cleanup d... | Source: enhancement 0012 ("The library owns the whole deletion sequenc... | The operator builds each deletion plan with the library's constructor | 2 | none |
| The operator decides prune and cleanup d... | Source: enhancement 0012 ("The library owns the whole deletion sequenc... | One function of the operator performs the deletion plan | 3 | none |
| The operator decides prune and cleanup d... | Source: enhancement 0012 ("The library owns the whole deletion sequenc... | The forced recreate delete asks no verdict | 3 | none |
| The operator decides prune and cleanup d... | Source: enhancement 0012 ("The library owns the whole deletion sequenc... | The forced recreate delete carries the UID precondition | 2 | none |
| The operator decides prune and cleanup d... | The one other delete the operator makes, the delete inside a forced recreate, asks no verdict ... | The forced recreate delete asks no verdict | 2 | "and MUST carry the UID precondition (capability `ssa-apply`)" left out of this block (list cut at the "and") |
| The operator decides prune and cleanup d... | The one other delete the operator makes, the delete inside a forced recreate, ... MUST carry the UID precondition (capability `ssa-apply`). | The forced recreate delete carries the UID precondition | 1 | "asks no verdict and" left out (list cut); the subject is repeated |
| The operator decides prune and cleanup d... | A test MUST keep the list of places that may delete a cluster object c... | A test keeps the places that may delete a cluster object clo... | 1 | none |
| The operator decides prune and cleanup d... | It MUST find a delete by the type of the receiver, a Kubernetes client... | A test finds a delete by the type of the receiver | 1 | -It +A test |
| The operator decides every apply only th... | Every write of a cluster object that the operator makes on behalf of a... | A write follows an apply verdict of the library's ownership... | 1 | none |
| The operator decides every apply only th... | The operator MUST declare no ownership rule of its own for an apply: i... | The operator declares no ownership rule of its own for an ap... | 1 | none |
| The operator decides every apply only th... | The operator MUST NOT set, change or remove the `opmodel.dev/adopt` an... | The operator never writes the adopt annotation | 1 | none |
| The operator decides every apply only th... | Source: enhancement 0012 ("The library owns the whole deletion sequenc... | A write follows an apply verdict of the library's ownership... | 2 | -library owns the whole deletion sequence but only the per-object verdict and the order for apply" and "The |
| The operator decides every apply only th... | Source: enhancement 0012 ("The library owns the whole deletion sequenc... | The operator declares no ownership rule of its own for an ap... | 2 | -library owns the whole deletion sequence but only the per-object verdict and the order for apply" and "The |
| The operator decides every apply only th... | Source: enhancement 0012 ("The library owns the whole deletion sequenc... | The operator never writes the adopt annotation | 2 | none |
| The operator decides every apply only th... | Source: enhancement 0012 ("The library owns the whole deletion sequenc... | A test keeps the places that may write a cluster object clos... | 2 | -library owns the whole deletion sequence but only the per-object verdict and the order for apply" and "The |
| The operator decides every apply only th... | Source: enhancement 0012 ("The library owns the whole deletion sequenc... | A test finds a write by the type of the receiver | 2 | -library owns the whole deletion sequence but only the per-object verdict and the order for apply" and "The |
| The operator decides every apply only th... | A test MUST keep the list of places that may write a cluster object cl... | A test keeps the places that may write a cluster object clos... | 1 | none |
| The operator decides every apply only th... | It MUST find a write by the type of the receiver, as the test of the d... | A test finds a write by the type of the receiver | 1 | -It +A test |


Lengths (old body, new bodies; characters by the guide's rule):

| Old requirement | Old | New |
| --- | --- | --- |
| The operator uses the library's object and label packages and keeps no copy | 692 | 487, 431, 291 |
| The reconciler converts a render with one export | 816 | 482, 357 |
| The operator uses the library's inventory package and keeps no copy | 742 | 398, 378, 270 |
| The operator judges the health of inventory objects only through the library's health package | 925 | 436, 423, 360 |
| The operator decides prune and cleanup deletes only through the library's ownership package | 1472 | 468, 454, 429, 276, 348, 484, 266, 269, 398 |
| The operator decides every apply only through the library's ownership package | 1191 | 418, 356, 409, 424, 420 |

### Ledger: module-instance-ownership

Form c2 for every requirement below: REMOVED under the old name, ADDED under new names, new blocks at the end of the spec. `grep -rn --exclude-dir=archive` for each old header found no reference to fix. Old scenario text is copied by script from `origin/main`; item 3 of the guide checklist shows `same` for each.

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| Field accepts the two enum values | ModuleInstance carries an optional owner field with two values, Field accepts the two enum values | no (byte equal, checked by the item 3 command) |
| Field is optional with no default | ModuleInstance carries an optional owner field with two values, Field is optional with no default | no (byte equal, checked by the item 3 command) |
| Empty owner is reconciled as operator-managed | An absent or empty owner is operator-managed, Empty owner is reconciled as operator-managed | no (byte equal, checked by the item 3 command) |
| Unknown owner values cannot reach the reconciler | ModuleInstance carries an optional owner field with two values, Unknown owner values cannot reach the reconciler | no (byte equal, checked by the item 3 command) |
| CLI-owned instance is not reconciled and gets no finalizer | A CLI-owned instance is not reconciled and gets no finalizer, CLI-owned instance is not reconciled and gets no finalizer | no (byte equal, checked by the item 3 command) |
| Deleting a CLI-owned instance is a no-op for the operator | A CLI-owned instance is not reconciled and gets no finalizer, Deleting a CLI-owned instance is a no-op for the operator | no (byte equal, checked by the item 3 command) |
| Deleting a CLI-owned instance that carries the finalizer releases it | A CLI-owned instance that carries the finalizer releases it when it is deleted, Deleting a CLI-owned instance that carries the finalizer releases it | no (byte equal, checked by the item 3 command) |
| A live CLI-owned instance keeps a leftover finalizer | A live CLI-owned instance keeps its finalizers, A live CLI-owned instance keeps a leftover finalizer | no (byte equal, checked by the item 3 command) |
| Operator-managed instances are unaffected | A CLI-owned instance is not reconciled and gets no finalizer, Operator-managed instances are unaffected | no (byte equal, checked by the item 3 command) |
| ManagedExternally condition is set | The status package exposes the ManagedExternally reason and its helper, ManagedExternally condition is set | no (byte equal, checked by the item 3 command) |
| No observedGeneration and no CLI-written status is touched | A CLI-owned instance is acknowledged with Ready Unknown and no other status, No observedGeneration and no CLI-written status is touched | no (byte equal, checked by the item 3 command) |
| Re-acknowledgement is a no-op | A CLI-owned instance is acknowledged with Ready Unknown and no other status, Re-acknowledgement is a no-op | no (byte equal, checked by the item 3 command) |
| Fixed coordinates identify the instance | The fixed coordinates identify the operator's own instance, Fixed coordinates identify the instance | no (byte equal, checked by the item 3 command) |
| The operator module identifies the instance in any namespace and any major | The operator module path identifies the operator's own instance, The operator module identifies the instance in any namespace and any major | no (byte equal, checked by the item 3 command) |
| A recorded operator CRD identifies the instance | A recorded operator CRD identifies the operator's own instance, A recorded operator CRD identifies the instance | no (byte equal, checked by the item 3 command) |
| Similar names are not the operator's own instance | A recorded operator CRD identifies the operator's own instance, Similar names are not the operator's own instance | no (byte equal, checked by the item 3 command) |
| An own instance flipped to operator is refused without a finalizer | The controller never applies, prunes or adds a finalizer on its own instance, An own instance flipped to operator is refused without a finalizer | no (byte equal, checked by the item 3 command) |
| A suspended own instance is refused, not suspended | An own instance with no owner or owner operator is refused before any render, A suspended own instance is refused, not suspended | no (byte equal, checked by the item 3 command) |
| Refusal clears conditions left by an earlier adoption | A refusal sets Ready and Stalled with SelfManagementRefused, Refusal clears conditions left by an earlier adoption | no (byte equal, checked by the item 3 command) |
| An own instance with no owner is refused | An own instance with no owner or owner operator is refused before any render, An own instance with no owner is refused | no (byte equal, checked by the item 3 command) |
| Refusal leaves CLI-written status alone | A refusal leaves the CLI-written status alone and does not requeue, Refusal leaves CLI-written status alone | no (byte equal, checked by the item 3 command) |
| Re-refusal is a no-op | Re-refusing an already-refused instance changes no status, Re-refusal is a no-op | no (byte equal, checked by the item 3 command) |
| Handing the instance back to the CLI clears the refusal | An own instance with no owner or owner operator is refused before any render, Handing the instance back to the CLI clears the refusal | no (byte equal, checked by the item 3 command) |
| Deleting an own instance with a leftover finalizer prunes nothing | The operator releases a leftover finalizer from its own instance without pruning, Deleting an own instance with a leftover finalizer prunes nothing | no (byte equal, checked by the item 3 command) |
| A live own instance loses a leftover finalizer | The operator releases a leftover finalizer from its own instance without pruning, A live own instance loses a leftover finalizer | no (byte equal, checked by the item 3 command) |
| A CLI-owned own instance is left alone | The operator writes no status on its own instance that is being deleted, A CLI-owned own instance is left alone | no (byte equal, checked by the item 3 command) |

New scenarios (restate a clause of the old body, with words of the old text only):

| New scenario | Requirement | Old clause restated |
| --- | --- | --- |
| Refusal records the generation | A refusal records the generation and emits one Warning event | the clause of that requirement body |
| Refusal emits a Warning event once | A refusal records the generation and emits one Warning event | the clause of that requirement body |
| The reason constant is exposed | The status package exposes the SelfManagementRefused reason | the clause of that requirement body |

Clause table (old clause, new requirement, sentence of its body; "Edit" lists the words that differ from the old sentence):

| Old requirement | Old clause | New requirement | Sentence | Edit |
| --- | --- | --- | --- | --- |
| ModuleInstance carries an ownership mark... | `ModuleInstanceSpec` SHALL provide an optional `owner` field of a type... | ModuleInstance carries an optional owner field with two valu... | 1 | none |
| ModuleInstance carries an ownership mark... | The field SHALL NOT define a CRD-level default. | ModuleInstance carries an optional owner field with two valu... | 2 | none |
| ModuleInstance carries an ownership mark... | The API SHALL define exported constants for both values. | ModuleInstance carries an optional owner field with two valu... | 3 | none |
| ModuleInstance carries an ownership mark... | An absent or empty `owner` SHALL be treated by the controller as opera... | An absent or empty owner is operator-managed | 1 | none |
| The operator skips CLI-owned instances b... | When `spec.owner == cli`, the controller MUST return from `Reconcile`... | A CLI-owned instance is not reconciled and gets no finalizer | 1 | none |
| The operator skips CLI-owned instances b... | The owner check MUST occur before finalizer registration so that an in... | A CLI-owned instance is not reconciled and gets no finalizer | 2 | none |
| The operator skips CLI-owned instances b... | An instance that was operator-owned before its `spec.owner` became `cl... | A CLI-owned instance that carries the finalizer releases it... | 1 | none |
| The operator skips CLI-owned instances b... | An instance that was operator-owned before its `spec.owner` became `cl... | A live CLI-owned instance keeps its finalizers | 1 | none |
| The operator skips CLI-owned instances b... | When such an instance is being deleted, the controller MUST remove the... | A CLI-owned instance that carries the finalizer releases it... | 2 | none |
| The operator skips CLI-owned instances b... | On a CLI-owned instance that is not being deleted the controller MUST... | A live CLI-owned instance keeps its finalizers | 2 | none |
| The operator skips CLI-owned instances b... | Rationale: only the operator removes its finalizer, and the owner-skip... | A CLI-owned instance that carries the finalizer releases it... | 3 | none |
| The operator skips CLI-owned instances b... | Without the release the delete never completes. | A CLI-owned instance that carries the finalizer releases it... | 4 | none |
| The operator acknowledges a CLI-owned in... | On a non-deleting `ModuleInstance` with `spec.owner == cli`, the contr... | A CLI-owned instance is acknowledged with Ready Unknown and... | 1 | -field (em-dash) +field: |
| The operator acknowledges a CLI-owned in... | The `internal/status` package MUST expose a `ManagedExternally` reason... | The status package exposes the ManagedExternally reason and... | 1 | none |
| The operator acknowledges a CLI-owned in... | The acknowledgement MUST be idempotent: reconciling an already-acknowl... | A CLI-owned instance is acknowledged with Ready Unknown and... | 2 | none |
| The operator recognises the instance tha... | The controller SHALL treat a `ModuleInstance` as the operator's own in... | The fixed coordinates identify the operator's own instance | 1 | none |
| The operator recognises the instance tha... | The controller SHALL treat a `ModuleInstance` as the operator's own in... | The operator module path identifies the operator's own insta... | 1 | none |
| The operator recognises the instance tha... | The controller SHALL treat a `ModuleInstance` as the operator's own in... | A recorded operator CRD identifies the operator's own instan... | 1 | none |
| The operator recognises the instance tha... | - its `metadata.name` is `opm-operator` and its `metadata.namespace` i... | The fixed coordinates identify the operator's own instance | 2 | none |
| The operator recognises the instance tha... | - its `spec.module.path`, with any `@<major>` suffix removed, is `opmo... | The operator module path identifies the operator's own insta... | 2 | none |
| The operator recognises the instance tha... | - its `status.inventory.entries` holds an entry of group `apiextension... | A recorded operator CRD identifies the operator's own instan... | 2 | none |
| The operator recognises the instance tha... | The coordinates are the ones every CLI install records the operator's... | The fixed coordinates identify the operator's own instance | 3 | none |
| The operator recognises the instance tha... | The coordinates are the ones every CLI install records the operator's... | The operator module path identifies the operator's own insta... | 3 | none |
| The operator recognises the instance tha... | A refusal must be decided before the finalizer is registered, which is... | A recorded operator CRD identifies the operator's own instan... | 3 | none |
| The operator never reconciles its own in... | On the operator's own instance, whatever `spec.owner` says, the contro... | The controller never applies, prunes or adds a finalizer on... | 1 | none |
| The operator never reconciles its own in... | On the operator's own instance, whatever `spec.owner` says, the contro... | An own instance with no owner or owner operator is refused b... | 1 | none |
| The operator never reconciles its own in... | An own instance with `spec.owner: cli` is handled by the owner-skip ga... | An own instance with no owner or owner operator is refused b... | 2 | none |
| The operator never reconciles its own in... | On refusal the controller MUST: | A refusal sets Ready and Stalled with SelfManagementRefused; A refusal records the generation and emits one Warning event; A refusal leaves the CLI-written status alone and does not requeue | 1 | "of the operator's own instance" added after "On refusal" (condition copied from the first sentence of the old body), because no sibling is read for it |
| The operator never reconciles its own in... | - set `Ready=False` and `Stalled=True`, both with reason `SelfManageme... | A refusal sets Ready and Stalled with SelfManagementRefused | 2 | none |
| The operator never reconciles its own in... | - set `status.observedGeneration` to the instance's generation, so a c... | A refusal records the generation and emits one Warning event | 2 | none |
| The operator never reconciles its own in... | - emit a Warning event with reason `SelfManagementRefused` when the in... | A refusal records the generation and emits one Warning event | 3 | none |
| The operator never reconciles its own in... | - leave `status.inventory`, the `lastApplied*` digests, `instanceUUID`... | A refusal leaves the CLI-written status alone and does not r... | 2 | none |
| The operator never reconciles its own in... | - not requeue: the refusal is re-evaluated when the instance changes. | A refusal leaves the CLI-written status alone and does not r... | 3 | none |
| The operator never reconciles its own in... | The `internal/status` package MUST expose a `SelfManagementRefused` re... | The status package exposes the SelfManagementRefused reason | 1 | none |
| The operator never reconciles its own in... | Re-refusing an already-refused instance at the same generation MUST pr... | Re-refusing an already-refused instance changes no status | 1 | none |
| The operator never reconciles its own in... | Rationale: the operator is the one workload whose failure stops every... | The controller never applies, prunes or adds a finalizer on... | 2 | none |
| The operator never reconciles its own in... | Nothing the operator does to its own instance may stand between a re-r... | The controller never applies, prunes or adds a finalizer on... | 3 | none |
| The operator releases its finalizer from... | When the operator's own instance with `spec.owner` absent or `operator... | The operator releases a leftover finalizer from its own inst... | 1 | none |
| The operator releases its finalizer from... | On an own instance being deleted, the controller MUST write no status. | The operator writes no status on its own instance that is be... | 1 | none |
| The operator releases its finalizer from... | An own instance with `spec.owner: cli` is left to the owner-skip gate,... | The operator writes no status on its own instance that is be... | 2 | none |
| The operator releases its finalizer from... | Rationale: a finalizer only a running operator can clear would block t... | The operator writes no status on its own instance that is be... | 3 | none |


Lengths (old body, new bodies; characters by the guide's rule):

| Old requirement | Old | New |
| --- | --- | --- |
| ModuleInstance carries an ownership marker | 536 | 274, 261 |
| The operator skips CLI-owned instances before registering a finalizer | 926 | 338, 490, 198 |
| The operator acknowledges a CLI-owned instance with a single condition | 804 | 463, 339 |
| The operator recognises the instance that deploys it | 845 | 434, 440, 482 |
| The operator never reconciles its own instance | 1619 | 446, 422, 302, 312, 262, 84, 98 |
| The operator releases its finalizer from its own instance without pruning | 767 | 341, 425 |

### Ledger: platform-reconciler

Form c2 for every requirement below: REMOVED under the old name, ADDED under new names, new blocks at the end of the spec. `grep -rn --exclude-dir=archive` for each old header found no reference to fix. Old scenario text is copied by script from `origin/main`; item 3 of the guide checklist shows `same` for each.

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| Valid platform is generated | The Platform named cluster is generated, built and recorded in the store, Valid platform is generated | no (byte equal, checked by the item 3 command) |
| Non-cluster object ignored | Only the Platform named cluster is reconciled, Non-cluster object ignored | no (byte equal, checked by the item 3 command) |
| A claim change wakes the reconciler | The generated package is a function of the spec and the active claims, A claim change wakes the reconciler | no (byte equal, checked by the item 3 command) |
| An unchanged tuple regenerates nothing | A reconcile whose inputs resolve to the held identity skips regeneration, An unchanged tuple regenerates nothing | no (byte equal, checked by the item 3 command) |
| Success sets Ready and observedGeneration | A successful build sets Ready True with reason Generated, Success sets Ready and observedGeneration | no (byte equal, checked by the item 3 command) |
| Build failure surfaces structured error | A build failure sets Ready False with reason BuildFailed, Build failure surfaces structured error | no (byte equal, checked by the item 3 command) |
| A refusal surfaces a structured verdict | A contract collision sets Ready False with reason ContractCollisions, A refusal surfaces a structured verdict | no (byte equal, checked by the item 3 command) |
| Failure preserves last-good generated platform | Every outcome sets observedGeneration and a failure keeps the last good module, Failure preserves last-good generated platform | no (byte equal, checked by the item 3 command) |
| A spec change replaces the slot | The package identity is derived from the generation and the active claims, A spec change replaces the slot | no (byte equal, checked by the item 3 command) |
| A claim change replaces the slot without a spec change | The generated-platform record is held per package identity, A claim change replaces the slot without a spec change | no (byte equal, checked by the item 3 command) |
| A leased package survives the swap | Readers lease the record and superseded directories are removed only when unleased, A leased package survives the swap | no (byte equal, checked by the item 3 command) |
| Concurrent reads are safe | The store holds at most one generated-module record, Concurrent reads are safe | no (byte equal, checked by the item 3 command) |
| Failure requeues instead of stalling indefinitely | A build failure requeues after a bounded interval and keeps the last good module, Failure requeues instead of stalling indefinitely | no (byte equal, checked by the item 3 command) |
| Recovery without a spec change | A build failure requeues after a bounded interval and keeps the last good module, Recovery without a spec change | no (byte equal, checked by the item 3 command) |
| A transient registry failure does not outlive its reconcile | Each reconcile resolves module files through a source built for it, A transient registry failure does not outlive its reconcile | no (byte equal, checked by the item 3 command) |
| Transient cause retries quickly | A failure is transient when its error chain holds a registry fetch failure, Transient cause retries quickly | no (byte equal, checked by the item 3 command) |
| Semantic or unknown cause retries slowly | Transient and other failures are requeued on different intervals, Semantic or unknown cause retries slowly | no (byte equal, checked by the item 3 command) |
| Missing version surfaces as BuildFailed | A subscription that cannot be built surfaces as BuildFailed, Missing version surfaces as BuildFailed | no (byte equal, checked by the item 3 command) |
| Named build not published | A subscription that cannot be built surfaces as BuildFailed, Named build not published | no (byte equal, checked by the item 3 command) |
| Identity mismatch surfaces with both values | A subscription that cannot be built surfaces as BuildFailed, Identity mismatch surfaces with both values | no (byte equal, checked by the item 3 command) |

New scenarios (restate a clause of the old body, with words of the old text only):

| New scenario | Requirement | Old clause restated |
| --- | --- | --- |
| A stale event yields the package the current state implies | The waking event's content is not an input to the generated package | the clause of that requirement body |
| A module that could not be written surfaces as GenerateFailed | A write failure sets Ready False with reason GenerateFailed | the clause of that requirement body |
| A platform that is not routable surfaces as OverSubscribedContracts | An over-subscribed contract sets Ready False with reason OverSubscribedContracts | the clause of that requirement body |
| A platform that is not discriminated surfaces as ComparablePredicates | Comparable predicates set Ready False with reason ComparablePredicates | the clause of that requirement body |
| Several refusals hold at once | When several refusals hold the first reason applies | the clause of that requirement body |
| The failure class does not change the condition | The failure class picks the requeue interval only | the clause of that requirement body |
| A subscription version is used as written | A subscription version is used verbatim as the pin and the stamp | the clause of that requirement body |

Clause table (old clause, new requirement, sentence of its body; "Edit" lists the words that differ from the old sentence):

| Old requirement | Old clause | New requirement | Sentence | Edit |
| --- | --- | --- | --- | --- |
| Generate and build the singleton Platfor... | The operator SHALL reconcile the `Platform` named `cluster` by generat... | The Platform named cluster is generated, built and recorded... | 1 | none |
| Generate and build the singleton Platfor... | The reconciler SHALL reconcile only the object named `cluster`; any ot... | Only the Platform named cluster is reconciled | 1 | none |
| Generate and build the singleton Platfor... | No second, derived copy of the platform exists. | The Platform named cluster is generated, built and recorded... | 2 | none |
| Generate and build the singleton Platfor... | The package the reconcile generates SHALL be a function of exactly two... | The generated package is a function of the spec and the acti... | 1 | -inputs (em-dash) +inputs:; -claims (em-dash) +claims, |
| Generate and build the singleton Platfor... | The waking event's content SHALL NOT be an input: the reconcile comput... | The waking event's content is not an input to the generated... | 2 | none |
| Generate and build the singleton Platfor... | A reconcile whose inputs resolve to the identity the store already hol... | A reconcile whose inputs resolve to the held identity skips... | 1 | none |
| Surface build outcome on status | prefix "The reconciler SHALL record the outcome on the `Platform` status:" | every one of the six outcome requirements | 1 | list cut at item boundaries; prefix repeated in each |
| Surface build outcome on status | item `Ready=True` with reason `Generated` on success | A successful build sets Ready True with reason Generated | 1 | none |
| Surface build outcome on status | item `Ready=False` with reason `BuildFailed` (...) | A build failure sets Ready False with reason BuildFailed | 1 | none |
| Surface build outcome on status | item `GenerateFailed` (the module could not be written) | A write failure sets Ready False with reason GenerateFailed | 1 | "`Ready=False` with reason" copied in front from the `BuildFailed` item, whose reason it shares |
| Surface build outcome on status | item `ContractCollisions` (...) | A contract collision sets Ready False with reason ContractCollisions | 1 | "`Ready=False` with reason" copied in front, as above |
| Surface build outcome on status | item `OverSubscribedContracts` (...) | An over-subscribed contract sets Ready False with reason OverSubscribedContracts | 1 | "`Ready=False` with reason" copied in front, as above |
| Surface build outcome on status | item `ComparablePredicates` (...) | Comparable predicates set Ready False with reason ComparablePredicates | 1 | "`Ready=False` with reason" copied in front, as above |
| Surface build outcome on status | When several refusals hold, the reason SHALL be the first of `Contract... | When several refusals hold the first reason applies | 1 | none |
| Surface build outcome on status | On success the reconciler SHALL also write the non-gating `ContractsFu... | A successful build sets Ready True with reason Generated | 2 | none |
| Surface build outcome on status | `status.observedGeneration` SHALL be set on every outcome. | Every outcome sets observedGeneration and a failure keeps th... | 1 | none |
| Surface build outcome on status | A failure or a refusal SHALL NOT overwrite a previously recorded good... | Every outcome sets observedGeneration and a failure keeps th... | 2 | none |
| Single-slot identity-keyed store for con... | The store SHALL hold at most one current generated-module record, keye... | The store holds at most one generated-module record | 1 | none |
| Single-slot identity-keyed store for con... | Readers SHALL obtain the record through a lease that they release when... | Readers lease the record and superseded directories are remo... | 1 | none |
| Single-slot identity-keyed store for con... | Superseded module directories SHALL be removed only when no lease hold... | Readers lease the record and superseded directories are remo... | 2 | none |
| Single-slot identity-keyed store for con... | The identity SHALL be derived from the `Platform` CR's generation and... | The package identity is derived from the generation and the... | 1 | -accepted" and "The generated-platform record survives the render-parity enhancement and is held per package identity, not per Platform CR generation"). +accept |
| Single-slot identity-keyed store for con... | The identity SHALL be derived from the `Platform` CR's generation and... | The generated-platform record is held per package identity | 1 | -platform package regenerates edge-triggered and level-computed, keyed by generation plus the active-claim set; blast radius accepted" and "The |
| Build failures requeue on a bounded inte... | When closure derivation, generation or the build fails, the `PlatformR... | A build failure requeues after a bounded interval and keeps... | 1 | none |
| Build failures requeue on a bounded inte... | The reconciler SHALL set the failure reason (`BuildFailed` or `Generat... | A build failure requeues after a bounded interval and keeps... | 2 | none |
| Build failures requeue on a bounded inte... | Each reconcile SHALL resolve module files for closure derivation throu... | Each reconcile resolves module files through a source built... | 1 | none |
| Build failures requeue on a bounded inte... | A module-file source injected into the reconciler (tests) SHALL be use... | Each reconcile resolves module files through a source built... | 2 | none |
| Transient failures retry faster than sem... | The reconciler SHALL requeue a transient failure on a short interval a... | Transient and other failures are requeued on different inter... | 1 | none |
| Transient failures retry faster than sem... | The reconciler SHALL requeue a transient failure on a short interval a... | The failure class picks the requeue interval only | 1 | none |
| Transient failures retry faster than sem... | A failure is transient when its error chain holds the library's typed... | A failure is transient when its error chain holds a registry... | 1 | -is +SHALL be |
| Transient failures retry faster than sem... | The shared classification picks the Platform's requeue interval only:... | The failure class picks the requeue interval only | 2 | none |
| Transient failures retry faster than sem... | The reconciler SHALL classify by type, never by message text; an unrec... | Transient and other failures are requeued on different inter... | 2 | none |
| Transient failures retry faster than sem... | The condition the failure sets is unchanged by its class. | The failure class picks the requeue interval only | 3 | none |
| Subscription mapping and build failures | Each subscription's `version` SHALL be used verbatim as the generated... | A subscription version is used verbatim as the pin and the s... | 1 | none |
| Subscription mapping and build failures | A subscription with an empty version (a stored object predating the re... | A subscription that cannot be built surfaces as BuildFailed | 1 | none |
| Subscription mapping and build failures | A version absent from the registry SHALL surface as `BuildFailed` nami... | A subscription that cannot be built surfaces as BuildFailed | 2 | none |
| Subscription mapping and build failures | A pin whose bytes disagree with the stamp SHALL surface as `BuildFaile... | A subscription that cannot be built surfaces as BuildFailed | 3 | none |
| Subscription mapping and build failures | All retain the stalled recheck interval. | A subscription that cannot be built surfaces as BuildFailed | 4 | none |


Lengths (old body, new bodies; characters by the guide's rule):

| Old requirement | Old | New |
| --- | --- | --- |
| Generate and build the singleton Platform on reconcile | 1379 | 426, 110, 377, 354, 268 |
| Surface build outcome on status | 1387 | 243, 275, 143, 293, 296, 237, 189, 136 |
| Single-slot identity-keyed store for concurrent read | 936 | 156, 234, 404, 401 |
| Build failures requeue on a bounded interval | 675 | 338, 336 |
| Transient failures retry faster than semantic failures | 844 | 250, 382, 347 |
| Subscription mapping and build failures | 519 | 127, 391 |

### Ledger: kernel-module-renderer

Form c2 for every requirement below: REMOVED under the old name, ADDED under new names, new blocks at the end of the spec. `grep -rn --exclude-dir=archive` for each old header found no reference to fix. Old scenario text is copied by script from `origin/main`; item 3 of the guide checklist shows `same` for each.

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| Renders resources from a generated platform | KernelModuleRenderer renders an instance through the kernel and shares nothing between kernel calls, Renders resources from a generated platform | no (byte equal, checked by the item 3 command) |
| Values are applied when supplied | Values reach synthesis as a stack of values sources with an origin, Values are applied when supplied | no (byte equal, checked by the item 3 command) |
| A values error names its origin | Values reach synthesis as a stack of values sources with an origin, A values error names its origin | no (byte equal, checked by the item 3 command) |
| Compiled item maps to a resource | The renderer adapts compiled objects to the library's object resources, Compiled item maps to a resource | no (byte equal, checked by the item 3 command) |
| One compiled-output type | The renderer adapts compiled objects to the library's object resources, One compiled-output type | no (byte equal, checked by the item 3 command) |
| Warnings reach the reconciler | The renderer composes the result's warnings from the advisory rows, Warnings reach the reconciler | no (byte equal, checked by the item 3 command) |
| Two registrations in one module are refused before apply | The refusal of duplicate identities builds no resource, Two registrations in one module are refused before apply | no (byte equal, checked by the item 3 command) |
| Distinct identities adapt as before | The renderer adapts compiled objects to the library's object resources, Distinct identities adapt as before | no (byte equal, checked by the item 3 command) |
| A rendered module reports its components' contracts | Both kernel renderers report the contract demand of the kernel's render, A rendered module reports its components' contracts | no (byte equal, checked by the item 3 command) |
| Both renderers fill the demand the same way | Both kernel renderers report the contract demand of the kernel's render, Both renderers fill the demand the same way | no (byte equal, checked by the item 3 command) |
| An empty demand is an empty list | An empty demand is reported as an empty, non-nil list, An empty demand is an empty list | no (byte equal, checked by the item 3 command) |
| No operator walk of the components | The operator does not compute the contract demand itself, No operator walk of the components | no (byte equal, checked by the item 3 command) |

New scenarios (restate a clause of the old body, with words of the old text only):

| New scenario | Requirement | Old clause restated |
| --- | --- | --- |
| A call passes no per-call registry argument | No acquisition or synthesis call passes a per-call registry or load-options argument | the clause of that requirement body |
| A refusal carries the library's error unchanged | The refusal of duplicate identities uses the library check and carries its error | the clause of that requirement body |
| A failed render yields no demand | The demand is read only from a successful render | the clause of that requirement body |

Clause table (old clause, new requirement, sentence of its body; "Edit" lists the words that differ from the old sentence):

| Old requirement | Old clause | New requirement | Sentence | Edit |
| --- | --- | --- | --- | --- |
| Render a ModuleRelease through the kerne... | `KernelModuleRenderer` SHALL acquire the module from the registry, syn... | KernelModuleRenderer renders an instance through the kernel... | 1 | none |
| Render a ModuleRelease through the kerne... | Every kernel call SHALL share nothing: acquisition, synthesis and the... | KernelModuleRenderer renders an instance through the kernel... | 2 | none |
| Render a ModuleRelease through the kerne... | Values SHALL reach synthesis as a stack of values sources, each carryi... | Values reach synthesis as a stack of values sources with an... | 1 | none |
| Render a ModuleRelease through the kerne... | No acquisition or synthesis call SHALL pass a per-call registry or loa... | No acquisition or synthesis call passes a per-call registry... | 1 | none |
| Adapt compiled output to operator resour... | The renderer SHALL adapt the single-build render's compiled objects (em-dash)... | The renderer adapts compiled objects to the library's object... | 1 | -objects (em-dash) +objects,; -verb (em-dash) +verb, |
| Adapt compiled output to operator resour... | The refusal SHALL use the library's duplicate-identity check and SHALL... | The refusal of duplicate identities builds no resource | 1 | -SHALL use the library's duplicate-identity check and SHALL carry its error unchanged, naming each shared identity once and every component and transformer that |
| Adapt compiled output to operator resour... | The render SHALL report no message strings. | The renderer composes the result's warnings from the advisor... | 1 | none |
| Adapt compiled output to operator resour... | The renderer SHALL compose the result's warnings itself from the rende... | The renderer composes the result's warnings from the advisor... | 2 | none |
| Adapt compiled output to operator resour... | Each composed warning SHALL name the same facts as before (em-dash) for skew,... | The renderer composes the result's warnings from the advisor... | 3 | -before (em-dash) +before: |
| The render reports the instance's contra... | Both kernel renderers (`KernelModuleRenderer` for a ModuleInstance, `K... | Both kernel renderers report the contract demand of the kern... | 1 | none |
| The render reports the instance's contra... | Both kernel renderers (`KernelModuleRenderer` for a ModuleInstance, `K... | The operator does not compute the contract demand itself | 1 | -diagnostics: every `#resources` and `#traits` key of every component of the instance, sorted and deduplicated. +diagnostics. |
| The render reports the instance's contra... | The operator SHALL NOT compute the demand itself by reading the instan... | The operator does not compute the contract demand itself | 2 | none |
| The render reports the instance's contra... | An empty demand SHALL be reported as an empty, non-nil list. | An empty demand is reported as an empty, non-nil list | 1 | none |
| The render reports the instance's contra... | Source: enhancement 0013 ("The kernel reports every contract a render... | Both kernel renderers report the contract demand of the kern... | 2 | none |
| The render reports the instance's contra... | Source: enhancement 0013 ("The kernel reports every contract a render... | The operator does not compute the contract demand itself | 3 | none |
| The render reports the instance's contra... | Source: enhancement 0013 ("The kernel reports every contract a render... | An empty demand is reported as an empty, non-nil list | 2 | none |
| The render reports the instance's contra... | Source: enhancement 0013 ("The kernel reports every contract a render... | The demand is read only from a successful render | 2 | none |
| The render reports the instance's contra... | The demand is read only from a successful render; this requirement doe... | The demand is read only from a successful render | 1 | -is +SHALL be |


Lengths (old body, new bodies; characters by the guide's rule):

| Old requirement | Old | New |
| --- | --- | --- |
| Render a ModuleRelease through the kernel | 778 | 413, 202, 161 |
| Adapt compiled output to operator resources | 1352 | 352, 255, 463, 431 |
| The render reports the instance's contract demand from the kernel | 796 | 442, 431, 175, 326 |

### Ledger: platform-inventory-gate

Form c2 for every requirement below: REMOVED under the old name, ADDED under new names, new blocks at the end of the spec. `grep -rn --exclude-dir=archive` for each old header found no reference to fix. Old scenario text is copied by script from `origin/main`; item 3 of the guide checklist shows `same` for each.

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| Two catalogs providing one contract are refused | An over-subscribed platform is refused with the reason OverSubscribedContracts, Two catalogs providing one contract are refused | no (byte equal, checked by the item 3 command) |
| Two majors of one provider catalog are refused | The over-subscription count is the one core computes, Two majors of one provider catalog are refused | no (byte equal, checked by the item 3 command) |
| Providers of a contract whose defining catalog is not enabled are refused | The over-subscription count is the one core computes, Providers of a contract whose defining catalog is not enabled are refused | no (byte equal, checked by the item 3 command) |
| An unroutable inventory that names nothing still refuses | The over-subscription finding is printed only for an over-subscribed contract, An unroutable inventory that names nothing still refuses | no (byte equal, checked by the item 3 command) |
| Discriminated transformers coexist | A platform with comparable transformer predicates is refused with the reason ComparablePredicates, Discriminated transformers coexist | no (byte equal, checked by the item 3 command) |
| A comparable pair is refused naming both transformers | The comparable-predicates message names the pair and the shared contracts, A comparable pair is refused naming both transformers | no (byte equal, checked by the item 3 command) |
| Both findings report under the routing reason | Over-subscription and comparable predicates report together, Both findings report under the routing reason | no (byte equal, checked by the item 3 command) |
| The last good package survives a refusal | A refused generation behaves as a failed build outside the Ready condition, The last good package survives a refusal | no (byte equal, checked by the item 3 command) |
| A first-boot refusal holds nothing | A refusal with no last good package leaves the store empty, A first-boot refusal holds nothing | no (byte equal, checked by the item 3 command) |
| The skip path keeps the condition current | The ContractsFulfilled condition is written on both success paths and never moves Ready, The skip path keeps the condition current | no (byte equal, checked by the item 3 command) |
| An unfulfilled contract is a report beside a generated platform | The ContractsFulfilled reasons report unfulfilled, fulfilled and undefined contracts, An unfulfilled contract is a report beside a generated platform | no (byte equal, checked by the item 3 command) |
| Every contract fulfilled | The ContractsFulfilled reasons report unfulfilled, fulfilled and undefined contracts, Every contract fulfilled | no (byte equal, checked by the item 3 command) |
| No contract defined is visible, not a pass | The ContractsFulfilled reasons report unfulfilled, fulfilled and undefined contracts, No contract defined is visible, not a pass | no (byte equal, checked by the item 3 command) |
| Two majors sharing contract keys are refused naming the keys and the entries | A platform whose enabled entries share contract keys is refused with the reason ContractCollisions, Two majors sharing contract keys are refused naming the keys and the entries | no (byte equal, checked by the item 3 command) |
| The collision message does not depend on report order | The collision message names each colliding key and its defining entries, The collision message does not depend on report order | no (byte equal, checked by the item 3 command) |
| A collision refuses although the fulfilment and discrimination reports read true | The collision refusal holds whatever the other verdicts read, A collision refuses although the fulfilment and discrimination reports read true | no (byte equal, checked by the item 3 command) |
| A collision is reported ahead of an over-subscription | ContractCollisions takes precedence over the other refusals, A collision is reported ahead of an over-subscription | no (byte equal, checked by the item 3 command) |
| A collision is reported ahead of a comparable pair | ContractCollisions takes precedence over the other refusals, A collision is reported ahead of a comparable pair | no (byte equal, checked by the item 3 command) |

New scenarios (restate a clause of the old body, with words of the old text only):

| New scenario | Requirement | Old clause restated |
| --- | --- | --- |
| The inventory is read before the package is recorded | The reconciler reads the contract inventory of the built platform before recording the package | the clause of that requirement body |
| More than one provider refuses the package | A contract provided by more than one enabled registry entry refuses the package | the clause of that requirement body |
| The message names the over-subscribed contract, its defining catalog and its providers | The over-subscription message names the contracts, the catalogs and the providers | the clause of that requirement body |
| A colliding key changes the reason | A collision changes the reason of an over-subscribed platform | the clause of that requirement body |
| The operator does not recount | The operator does not recount the over-subscription | the clause of that requirement body |
| No over-subscribed contract, no finding | The over-subscription finding is printed only for an over-subscribed contract | the clause of that requirement body |
| No arbitration of a comparable pair | No arbitration orders comparable transformers | the clause of that requirement body |
| The refused module directory stays until the next generation | A refusal with no last good package leaves the store empty | the clause of that requirement body |
| The collision message names the way out | The collision message says a package cannot be generated until all but one entry is disabled | the clause of that requirement body |
| Definers are read, not counted | The operator reads colliding keys from the inventory | the clause of that requirement body |
| Every contract key defined once | A platform whose enabled entries define every contract key once is not refused for collisions | the clause of that requirement body |

Clause table (old clause, new requirement, sentence of its body; "Edit" lists the words that differ from the old sentence):

| Old requirement | Old clause | New requirement | Sentence | Edit |
| --- | --- | --- | --- | --- |
| An over-subscribed platform is refuse... | After a generated platform module builds, the reconciler SHALL read... | The reconciler reads the contract inventory of the built ... | 1 | none |
| An over-subscribed platform is refuse... | When any provider-fulfilled contract is provided by more than one e... | An over-subscribed platform is refused with the reason Ov... | 1 | list cut inside the parenthesis at the citation boundary; the head is kept; added: ")" / "." |
| An over-subscribed platform is refuse... | When any provider-fulfilled contract is provided by more than one e... | A contract provided by more than one enabled registry ent... | 1 | list cut inside the parenthesis at the citation boundary; the head and the refusal are kept; added: " (" / ")" / "." |
| An over-subscribed platform is refuse... | When any provider-fulfilled contract is provided by more than one e... | The over-subscription message names the contracts, the ca... | 1 | list cut at "and a message naming"; the head and the refusal are kept; added: ": " / "." |
| An over-subscribed platform is refuse... | When any provider-fulfilled contract is provided by more than one e... | A collision changes the reason of an over-subscribed plat... | 1 | the condition and the refusal of the first sentence are copied as context; added: "." |
| An over-subscribed platform is refuse... | When the inventory also reports a colliding contract key, the reaso... | A collision changes the reason of an over-subscribed plat... | 2 | none |
| An over-subscribed platform is refuse... | When any provider-fulfilled contract is provided by more than one e... | The over-subscription count is the one core computes | 1 | the condition and the refusal of sentence 1 are copied as context; added: "." |
| An over-subscribed platform is refuse... | The count SHALL be the platform's own, as core computes it and the ... | The over-subscription count is the one core computes | 2 | none |
| An over-subscribed platform is refuse... | When any provider-fulfilled contract is provided by more than one e... | The operator does not recount the over-subscription | 1 | the condition and the refusal of sentence 1 are copied as context; added: "." |
| An over-subscribed platform is refuse... | As listed in requirement "The over-subscription count is the one co... | The operator does not recount the over-subscription | 2 | requirement named by its exact header: the old sentence says "it" for the count of the sentence before |
| An over-subscribed platform is refuse... | The operator SHALL NOT recount it, so the generation gate and the r... | The operator does not recount the over-subscription | 3 | none |
| An over-subscribed platform is refuse... | The over-subscription finding SHALL be printed only when the invent... | The over-subscription finding is printed only for an over... | 1 | none |
| An over-subscribed platform is refuse... | An inventory that is not routable while listing neither an over-sub... | The over-subscription finding is printed only for an over... | 2 | none |
| An undiscriminated platform is refuse... | When the built platform's inventory reports a pair of enabled trans... | A platform with comparable transformer predicates is refu... | 1 | list cut at "and a message naming"; added: "." |
| An undiscriminated platform is refuse... | When the built platform's inventory reports a pair of enabled trans... | The comparable-predicates message names the pair and the ... | 1 | list cut at "and a message naming"; the head and the refusal are kept; added: "." |
| An undiscriminated platform is refuse... | When the built platform's inventory reports a pair of enabled trans... | No arbitration orders comparable transformers | 1 | the condition and the refusal of sentence 0 are copied as context; added: "." |
| An undiscriminated platform is refuse... | No arbitration, ordering or most-specific-wins SHALL be applied. | No arbitration orders comparable transformers | 2 | none |
| An undiscriminated platform is refuse... | When a platform is both over-subscribed and undiscriminated and no ... | Over-subscription and comparable predicates report together | 1 | none |
| An undiscriminated platform is refuse... | When a contract key also collides, the reason SHALL be `ContractCol... | Over-subscription and comparable predicates report together | 2 | none |
| A refusal preserves the effective pac... | A refused generation SHALL behave as a failed build toward everythi... | A refused generation behaves as a failed build outside th... | 1 | none |
| A refusal preserves the effective pac... | Renders SHALL keep consuming the last good package. | A refused generation behaves as a failed build outside th... | 2 | none |
| A refusal preserves the effective pac... | A refusal with no last good package SHALL leave the store empty, so... | A refusal with no last good package leaves the store empty | 1 | none |
| A refusal preserves the effective pac... | The refused module directory MAY remain on disk until the next succ... | A refusal with no last good package leaves the store empty | 2 | none |
| Unfulfilled contracts are reported, n... | The reconciler SHALL write a `ContractsFulfilled` condition from th... | The ContractsFulfilled condition is written on both succe... | 1 | none |
| Unfulfilled contracts are reported, n... | The condition SHALL NOT move `Ready` and SHALL NOT be rewritten by ... | The ContractsFulfilled condition is written on both succe... | 2 | none |
| Unfulfilled contracts are reported, n... | `ContractsFulfilled=False` with reason `UnfulfilledContracts` SHALL... | The ContractsFulfilled reasons report unfulfilled, fulfil... | 1 | none |
| A platform whose enabled entries shar... | After a generated platform module builds, when the built platform's... | A platform whose enabled entries share contract keys is r... | 1 | list cut at "and a message naming"; added: "." |
| A platform whose enabled entries shar... | After a generated platform module builds, when the built platform's... | The collision message names each colliding key and its de... | 1 | list cut at "and saying that"; the head is repeated; added: "." |
| A platform whose enabled entries shar... | After a generated platform module builds, when the built platform's... | The collision message says a package cannot be generated ... | 1 | list cut at "and saying that"; the head is repeated; added: "." |
| A platform whose enabled entries shar... | The colliding keys and their defining entries SHALL be read from th... | The operator reads colliding keys from the inventory | 1 | none |
| A platform whose enabled entries shar... | The refusal SHALL hold whatever the inventory's other verdicts read... | The collision refusal holds whatever the other verdicts read | 1 | none |
| A platform whose enabled entries shar... | `ContractCollisions` SHALL take precedence over `OverSubscribedCont... | ContractCollisions takes precedence over the other refusals | 1 | none |
| A platform whose enabled entries shar... | When the inventory also reports an over-subscribed contract or a co... | ContractCollisions takes precedence over the other refusals | 2 | none |
| A platform whose enabled entries shar... | After a generated platform module builds, when the built platform's... | A platform whose enabled entries define every contract ke... | 1 | the condition and the refusal of sentence 0 are copied as context; added: "." |
| A platform whose enabled entries shar... | A platform whose enabled entries define every contract key once SHA... | A platform whose enabled entries define every contract ke... | 2 | none |

Lengths (old body, new bodies; characters by the guide's rule):

| Old requirement | Old | New |
| --- | --- | --- |
| An over-subscribed platform is refused at generation | 1901 | 137, 352, 348, 306, 377, 437, 320, 401 |
| An undiscriminated platform is refused at generation | 904 | 390, 343, 260, 297 |
| A refusal preserves the effective package and its status | 678 | 441, 236 |
| Unfulfilled contracts are reported, never refused | 938 | 488, 449 |
| A platform whose enabled entries share contract keys is refused at generation | 1356 | 315, 417, 367, 139, 263, 288, 416 |

### Ledger: render-input-key

Form c2 for every requirement below: REMOVED under the old name, ADDED under new names, new blocks at the end of the spec. `grep -rn --exclude-dir=archive` for each old header found no reference to fix. Old scenario text is copied by script from `origin/main`; item 3 of the guide checklist shows `same` for each. The two list bodies of "The render input key names what a render is a function of" and "A reconcile skips its render when its inputs are unchanged" move into the scenarios "The six parts of the render input key" and "The skip conditions"; their bodies keep the SHALL and name the scenario.

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| Any part changes the digest | The key digest is a fixed encoding of its parts, Any part changes the digest | no (byte equal, checked by the item 3 command) |
| The same parts give the same digest | The key digest is a fixed encoding of its parts, The same parts give the same digest | no (byte equal, checked by the item 3 command) |
| A missing part makes the key incomplete | A key with an empty part is incomplete, A missing part makes the key incomplete | no (byte equal, checked by the item 3 command) |
| A successful apply records the key | The last confirming render records its key on status, A successful apply records the key | no (byte equal, checked by the item 3 command) |
| A NoOp re-proves the key | The last confirming render records its key on status, A NoOp re-proves the key | no (byte equal, checked by the item 3 command) |
| The recorded identity is the one the render leased | The recorded key is built from what the render used, The recorded identity is the one the render leased | no (byte equal, checked by the item 3 command) |
| A failed apply keeps the previous key | An attempt that fails or skips its render leaves the recorded key as it was, A failed apply keeps the previous key | no (byte equal, checked by the item 3 command) |
| Unchanged inputs skip the render | A reconcile skips its render when, and only when, all of the listed conditions hold, Unchanged inputs skip the render | no (byte equal, checked by the item 3 command) |
| A changed input renders | A reconcile skips its render when, and only when, all of the listed conditions hold, A changed input renders | no (byte equal, checked by the item 3 command) |
| A failed attempt renders on the next reconcile | A reconcile skips its render when, and only when, all of the listed conditions hold, A failed attempt renders on the next reconcile | no (byte equal, checked by the item 3 command) |
| A spec edit outside the key renders | A reconcile skips its render when, and only when, all of the listed conditions hold, A spec edit outside the key renders | no (byte equal, checked by the item 3 command) |
| A new revision with the same digest renders once | A reconcile skips its render when, and only when, all of the listed conditions hold, A new revision with the same digest renders once | no (byte equal, checked by the item 3 command) |
| An absent inventory Job forces a render | A failed drift detection skips the render and reports the Job as Missing, An absent inventory Job forces a render | no (byte equal, checked by the item 3 command) |
| An absent inventory Job after a failed drift detection | A failed drift detection skips the render and reports the Job as Missing, An absent inventory Job after a failed drift detection | no (byte equal, checked by the item 3 command) |
| A skipped instance that has not rolled out requeues | A skipped reconcile judges health and patches only the Healthy condition, A skipped instance that has not rolled out requeues | no (byte equal, checked by the item 3 command) |
| A skipped package keeps its interval | A skipped reconcile requeues by its kind, A skipped package keeps its interval | no (byte equal, checked by the item 3 command) |
| An old render is repeated | The drift render interval is the longest an unchanged reconcile skips its render, An old render is repeated | no (byte equal, checked by the item 3 command) |
| Zero disables the skip | Zero disables the skip, Zero disables the skip | no (byte equal, checked by the item 3 command) |
| A negative interval is refused | A negative drift render interval is refused at startup, A negative interval is refused | no (byte equal, checked by the item 3 command) |
| A newer operator renders and applies once | The first reconcile after an operator or library upgrade renders, A newer operator renders and applies once | no (byte equal, checked by the item 3 command) |
| The library inventory digests apply once and then converge | A change to a stored digest ships with a version change, The library inventory digests apply once and then converge | no (byte equal, checked by the item 3 command) |
| A key recorded by the same version keeps the old digests until the interval | A change to a stored digest ships with a version change, A key recorded by the same version keeps the old digests until the interval | no (byte equal, checked by the item 3 command) |
| A restart does not render unchanged objects | A restart does not render unchanged objects, A restart does not render unchanged objects | no (byte equal, checked by the item 3 command) |

New scenarios (restate a clause of the old body, with words of the old text only):

| New scenario | Requirement | Old clause restated |
| --- | --- | --- |
| The six parts of the render input key | The render input key has six parts | the six list items of the old body, copied unchanged |
| The key is computed before the render | A reconcile computes the key before it renders | the clause of that requirement body |
| The skip conditions | A reconcile skips its render when, and only when, all of the listed conditions hold | the seven list items of the old body, copied unchanged |
| An absent Job is judged by a render | An absent batch Job forces a render because only a render tells whether it is expired | the clause of that requirement body |
| A render for an absent Job judges health again | A reconcile that renders for an absent batch Job judges health twice | the clause of that requirement body |
| A skipped reconcile does nothing beyond health | A skipped reconcile takes no slot and moves nothing | the clause of that requirement body |
| The interval schedules nothing | The drift render interval schedules no reconcile | the clause of that requirement body |

Clause table (old clause, new requirement, sentence of its body; "Edit" lists the words that differ from the old sentence):

| Old requirement | Old clause | New requirement | Sentence | Edit |
| --- | --- | --- | --- | --- |
| The render input key names what a render is... | the six list items (source digest, config digest, platform package identity, skew policy, operator version, library version) | scenario "The six parts of the render input key" of The render input key has six parts | THEN and AND lines | each item copied unchanged, the final ";" and "." dropped |
| The render input key names what a ren... | The operator SHALL compute a render input key for a ModuleInstance ... | The render input key has six parts | 1 | the list of six parts moves to the scenario "The six parts of the render input key"; the body keeps the SHALL and names the scenario; added: ", as the scenario "The six parts of the render input key" lists them" / "." |
| The render input key names what a ren... | Source: owner decision g3 of the kernel-plan walkthrough (2026-10-0... | The render input key has six parts | 2 | none |
| The render input key names what a ren... | The key's digest SHALL be `sha256:<hex>` over a fixed, versioned en... | The key digest is a fixed encoding of its parts | 1 | none |
| The render input key names what a ren... | Source: owner decision g3 of the kernel-plan walkthrough (2026-10-0... | The key digest is a fixed encoding of its parts | 2 | none |
| The render input key names what a ren... | A key with an empty part is incomplete: it SHALL never be recorded ... | A key with an empty part is incomplete | 1 | none |
| The render input key names what a ren... | Source: owner decision g3 of the kernel-plan walkthrough (2026-10-0... | A key with an empty part is incomplete | 2 | none |
| The key of the last confirming render... | ModuleInstance and ModulePackage SHALL carry an optional `status.la... | The last confirming render records its key on status | 1 | none |
| The key of the last confirming render... | The reconcilers SHALL write it on an attempt whose apply (and prune... | The last confirming render records its key on status | 2 | none |
| The key of the last confirming render... | ModuleInstance and ModulePackage SHALL carry an optional `status.la... | The recorded key is built from what the render used | 1 | none |
| The key of the last confirming render... | The recorded key SHALL be built from the platform package identity ... | The recorded key is built from what the render used | 2 | none |
| The key of the last confirming render... | When that key is incomplete, the field SHALL be cleared. | The recorded key is built from what the render used | 3 | none |
| The key of the last confirming render... | ModuleInstance and ModulePackage SHALL carry an optional `status.la... | An attempt that fails or skips its render leaves the reco... | 1 | none |
| The key of the last confirming render... | An attempt that fails, is refused, panics, or skips its render SHAL... | An attempt that fails or skips its render leaves the reco... | 2 | none |
| A reconcile skips its render when its... | Before rendering, a ModuleInstance or ModulePackage reconcile SHALL... | A reconcile computes the key before it renders | 1 | none |
| A reconcile skips its render when its... | A ModuleInstance or ModulePackage reconcile SHALL skip the render w... | A reconcile skips its render when, and only when, all of ... | 1 | the subject "It" is replaced by the subject of sentence 0; the list of seven conditions moves to the scenario "The skip conditions"; the body keeps the SHALL and names the scenario; added: "A " / ", as the scenario "The skip conditions" lists them" / "." |
| A reconcile skips its render when its... | Rationale: a render is the most expensive step of a reconcile, and ... | A reconcile skips its render when, and only when, all of ... | 2 | none |
| A reconcile skips its render when its... | A ModuleInstance or ModulePackage reconcile SHALL skip the render w... | An absent batch Job forces a render because only a render... | 1 | the subject "It" is replaced by the subject of sentence 0; the list of seven conditions is in the scenario "The skip conditions"; added: "A " / ", as the scenario "The skip conditions" lists them" / "." |
| A reconcile skips its render when its... | The last condition exists because only a render tells whether an ab... | An absent batch Job forces a render because only a render... | 2 | none |
| A reconcile skips its render when its... | A ModuleInstance or ModulePackage reconcile SHALL skip the render w... | A reconcile that renders for an absent batch Job judges h... | 1 | the subject "It" is replaced by the subject of sentence 0; the list of seven conditions is in the scenario "The skip conditions"; added: "A " / ", as the scenario "The skip conditions" lists them" / "." |
| A reconcile skips its render when its... | The last condition exists because only a render tells whether an ab... | A reconcile that renders for an absent batch Job judges h... | 2 | none |
| A reconcile skips its render when its... | A reconcile that renders for it has judged health once without a pa... | A reconcile that renders for an absent batch Job judges h... | 3 | none |
| A reconcile skips its render when its... | A ModuleInstance or ModulePackage reconcile SHALL skip the render w... | A failed drift detection skips the render and reports the... | 1 | the subject "It" is replaced by the subject of sentence 0; the list of seven conditions is in the scenario "The skip conditions"; added: "A " / ", as the scenario "The skip conditions" lists them" / "." |
| A reconcile skips its render when its... | While the last drift detection failed, the missing set is unknown a... | A failed drift detection skips the render and reports the... | 2 | none |
| A reconcile skips its render when its... | A skipped reconcile SHALL NOT take a render slot, lease the platfor... | A skipped reconcile takes no slot and moves nothing | 1 | none |
| A reconcile skips its render when its... | It is not a reconcile attempt: it records no outcome, no history an... | A skipped reconcile takes no slot and moves nothing | 2 | none |
| A reconcile skips its render when its... | A skipped reconcile SHALL judge health over `status.inventory` (`in... | A skipped reconcile judges health and patches only the He... | 1 | the subject "It" is replaced by the subject of sentence 12; added: "A skipped reconcile" |
| A reconcile skips its render when its... | The `Healthy` condition is the only status a skip writes; it reads ... | A skipped reconcile judges health and patches only the He... | 2 | none |
| A reconcile skips its render when its... | A skipped ModuleInstance reconcile SHALL return without a requeue u... | A skipped reconcile requeues by its kind | 1 | none |
| The drift render interval bounds how ... | The manager SHALL take `--drift-render-interval` (a duration, defau... | The drift render interval is the longest an unchanged rec... | 1 | none |
| The drift render interval bounds how ... | A reconcile triggered after that renders, so drift detection runs a... | The drift render interval is the longest an unchanged rec... | 2 | none |
| The drift render interval bounds how ... | Source: owner decision g3 of the kernel-plan walkthrough (2026-10-0... | The drift render interval is the longest an unchanged rec... | 3 | none |
| The drift render interval bounds how ... | The manager SHALL take `--drift-render-interval` (a duration, defau... | The drift render interval schedules no reconcile | 1 | none |
| The drift render interval bounds how ... | The interval SHALL NOT schedule a reconcile of its own. | The drift render interval schedules no reconcile | 2 | none |
| The drift render interval bounds how ... | Source: owner decision g3 of the kernel-plan walkthrough (2026-10-0... | The drift render interval schedules no reconcile | 3 | none |
| The drift render interval bounds how ... | The manager SHALL take `--drift-render-interval` (a duration, defau... | Zero disables the skip | 1 | none |
| The drift render interval bounds how ... | `0` SHALL disable the skip, so every reconcile renders, and a `NoOp... | Zero disables the skip | 2 | none |
| The drift render interval bounds how ... | Source: owner decision g3 of the kernel-plan walkthrough (2026-10-0... | Zero disables the skip | 3 | none |
| The drift render interval bounds how ... | The manager SHALL take `--drift-render-interval` (a duration, defau... | A negative drift render interval is refused at startup | 1 | none |
| The drift render interval bounds how ... | A negative value SHALL make the manager exit at startup with an err... | A negative drift render interval is refused at startup | 2 | none |
| The drift render interval bounds how ... | Source: owner decision g3 of the kernel-plan walkthrough (2026-10-0... | A negative drift render interval is refused at startup | 3 | none |
| An operator or library upgrade render... | Because the operator and library versions are key parts, the first ... | The first reconcile after an operator or library upgrade ... | 1 | none |
| An operator or library upgrade render... | The rule rests on every operator release changing `version.Version`... | The first reconcile after an operator or library upgrade ... | 2 | none |
| An operator or library upgrade render... | A change to how the operator computes a stored digest (`lastApplied... | A change to a stored digest ships with a version change | 1 | none |
| An operator or library upgrade render... | Starting the same operator and library again (a restart) SHALL NOT ... | A restart does not render unchanged objects | 1 | none |

Lengths (old body, new bodies; characters by the guide's rule):

| Old requirement | Old | New |
| --- | --- | --- |
| The render input key names what a render is a function of | 1205 | 359, 388, 273 |
| The key of the last confirming render is recorded on status | 726 | 392, 403, 261 |
| A reconcile skips its render when its inputs are unchanged | 2537 | 299, 391, 315, 417, 316, 282, 373, 217 |
| The drift render interval bounds how long a render is skipped | 717 | 450, 371, 439, 402 |
| An operator or library upgrade renders every object once | 669 | 267, 247, 153 |

### Ledger: ssa-apply

Form c2 for "A forced recreate keeps PersistentVolumeClaims unless data deletion is allowed" and "The staged apply order never contradicts the library's kind weights". Form b for the other three, because their header is cited: "Custom resources wait for discovery of a CRD in the same set" (quoted by "Staged apply ordering"), "A forced recreate deletes exactly the object that was read" (quoted by the PersistentVolumeClaim requirement) and "An apply is judged by the ownership verdict before its first write" (cited by `modulepackage-reconcile-loop` and `reconcile-loop-assembly`). A form b parent keeps every scenario of its old requirement, so the ADDED siblings carry new scenarios that restate their clause. The `Source:` line of each of the three old requirements that have one stays whole in every block cut from it. Old scenario text is copied by script from `origin/main`; item 3 of the guide checklist shows `same` for each.

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| Discovery serves the new kind late | Custom resources wait for discovery of a CRD in the same set, Discovery serves the new kind late | no (byte equal, checked by the item 3 command) |
| Custom resource without its CRD in the set fails at once | Custom resources wait for discovery of a CRD in the same set, Custom resource without its CRD in the set fails at once | no (byte equal, checked by the item 3 command) |
| Discovery never serves the kind | Custom resources wait for discovery of a CRD in the same set, Discovery never serves the kind | no (byte equal, checked by the item 3 command) |
| Other apply errors are not retried | Custom resources wait for discovery of a CRD in the same set, Other apply errors are not retried | no (byte equal, checked by the item 3 command) |
| An attempt is not bounded by the retry window | Custom resources wait for discovery of a CRD in the same set, An attempt is not bounded by the retry window | no (byte equal, checked by the item 3 command) |
| A claim with a changed immutable field is kept | A claim whose update the API server refuses returns a named error, A claim with a changed immutable field is kept | no (byte equal, checked by the item 3 command) |
| A claim that needs no recreate is applied | A claim whose update the API server refuses returns a named error, A claim that needs no recreate is applied | no (byte equal, checked by the item 3 command) |
| The claim is recreated when data deletion is allowed | A claim is recreated when data deletion is allowed, The claim is recreated when data deletion is allowed | no (byte equal, checked by the item 3 command) |
| Another kind is recreated as before | The forced recreate of every other kind is unchanged, Another kind is recreated as before | no (byte equal, checked by the item 3 command) |
| The resource manager refuses a claim delete | The resource manager refuses to delete a core PersistentVolumeClaim outside a data-deleting apply, The resource manager refuses a claim delete | no (byte equal, checked by the item 3 command) |
| The recreated object is the one that was read | A forced recreate deletes exactly the object that was read, The recreated object is the one that was read | no (byte equal, checked by the item 3 command) |
| An object replaced since the read is not deleted | A forced recreate deletes exactly the object that was read, An object replaced since the read is not deleted | no (byte equal, checked by the item 3 command) |
| A delete of a collection is refused | A forced recreate deletes exactly the object that was read, A delete of a collection is refused | no (byte equal, checked by the item 3 command) |
| A real staged apply follows the weights | The staged apply order is never the opposite of the library kind weights, A real staged apply follows the weights | no (byte equal, checked by the item 3 command) |
| A Flux bump moves a kind | A kind Flux applies first never weighs more than the kind it applies second, A Flux bump moves a kind | no (byte equal, checked by the item 3 command) |
| A Flux bump adds a kind to its order | A kind Flux applies first never weighs more than the kind it applies second, A Flux bump adds a kind to its order | no (byte equal, checked by the item 3 command) |
| A library bump moves a weight | A kind Flux applies first never weighs more than the kind it applies second, A library bump moves a weight | no (byte equal, checked by the item 3 command) |
| The pinned versions agree | The test suite fails when the pinned versions break the order rule, The pinned versions agree | no (byte equal, checked by the item 3 command) |
| Flux refines a tie | The check reports two kinds of equal weight that Flux orders, Flux refines a tie | no (byte equal, checked by the item 3 command) |
| A kind that only one order names | The check reports a kind that only one of the two orders names, A kind that only one order names | no (byte equal, checked by the item 3 command) |
| The failure says what to do | A failure tells the maintainer what changes, The failure says what to do | no (byte equal, checked by the item 3 command) |
| A new object is created | An apply is judged by the ownership verdict before its first write, A new object is created | no (byte equal, checked by the item 3 command) |
| A foreign object refuses the whole apply | An apply is judged by the ownership verdict before its first write, A foreign object refuses the whole apply | no (byte equal, checked by the item 3 command) |
| Another instance's object refuses the apply | An apply is judged by the ownership verdict before its first write, Another instance's object refuses the apply | no (byte equal, checked by the item 3 command) |
| The adopt annotation lets the instance take the object | An apply is judged by the ownership verdict before its first write, The adopt annotation lets the instance take the object | no (byte equal, checked by the item 3 command) |
| An object being deleted refuses a changed render | An apply is judged by the ownership verdict before its first write, An object being deleted refuses a changed render | no (byte equal, checked by the item 3 command) |
| An object adopted by another instance is not applied | An apply is judged by the ownership verdict before its first write, An object adopted by another instance is not applied | no (byte equal, checked by the item 3 command) |
| The earlier identity is not asked with while a change is not settled | An apply is judged by the ownership verdict before its first write, The earlier identity is not asked with while a change is not settled | no (byte equal, checked by the item 3 command) |
| An object adopted under the earlier identity is let go at once | An apply is judged by the ownership verdict before its first write, An object adopted under the earlier identity is let go at once | no (byte equal, checked by the item 3 command) |
| An unreadable object stops the apply | An apply is judged by the ownership verdict before its first write, An unreadable object stops the apply | no (byte equal, checked by the item 3 command) |
| A custom resource of a CRD in the same render | An apply is judged by the ownership verdict before its first write, A custom resource of a CRD in the same render | no (byte equal, checked by the item 3 command) |
| A package with no recorded identity lets an adopted object go | An apply is judged by the ownership verdict before its first write, A package with no recorded identity lets an adopted object go | no (byte equal, checked by the item 3 command) |
| An empty render with no identity is not failed | An apply is judged by the ownership verdict before its first write, An empty render with no identity is not failed | no (byte equal, checked by the item 3 command) |
| An inventoried object with another instance's UUID label is applied | An apply is judged by the ownership verdict before its first write, An inventoried object with another instance's UUID label is applied | no (byte equal, checked by the item 3 command) |
| An inventoried object without OPM labels is applied | An apply is judged by the ownership verdict before its first write, An inventoried object without OPM labels is applied | no (byte equal, checked by the item 3 command) |
| A taken-in object is not recreated | An apply is judged by the ownership verdict before its first write, A taken-in object is not recreated | no (byte equal, checked by the item 3 command) |

New scenarios (restate a clause of the old body, with words of the old text only):

| New scenario | Requirement | Old clause restated |
| --- | --- | --- |
| The retry interval is fixed and bounded | The discovery retry runs at a fixed interval within a bound | the clause of that requirement body |
| No retry of another error | The discovery retry covers no other error | the clause of that requirement body |
| The retry ends on a no-match error | Apply returns the no-match error when the retry ends | the clause of that requirement body |
| The data deletion option is off at its zero value | The data deletion option is off at its zero value and has no effect without force | the clause of that requirement body |
| The option has no effect without force | The data deletion option is off at its zero value and has no effect without force | the clause of that requirement body |
| A forced apply with the option off keeps the claim | A forced apply with the option off never deletes a PersistentVolumeClaim | the clause of that requirement body |
| Every claim is checked first | The apply checks every claim of the set before it applies or deletes any object | the clause of that requirement body |
| No object of the set is written after a refused claim update | A refused claim update applies, creates and deletes no object of the set | the clause of that requirement body |
| The claim rule comes first | The PersistentVolumeClaim rule is checked before the UID precondition | the clause of that requirement body |
| The client refuses a collection delete | The resource manager client refuses a delete of a collection | the clause of that requirement body |
| Flux orders two kinds of equal weight | Flux may refine the library order | the clause of that requirement body |
| The check reads the pinned module | The check reads the order from the pinned Flux module | the clause of that requirement body |
| The check covers a custom kind whose name ends in Class | The check covers the kinds of both orders and the kinds of another owner | the clause of that requirement body |
| Two objects of one kind | The check reports two objects of the same group and kind | the clause of that requirement body |
| The failure names the pairs and the weights | A failure names every pair of kinds in opposite order | the clause of that requirement body |
| The inventory decides the verdict input | The controller asks the apply verdict whether status.inventory lists the object | the clause of that requirement body |
| The verdict is asked on every apply | The apply verdict is asked for a ModuleInstance and a ModulePackage, for a changed render and for a restore | the clause of that requirement body |
| The read uses the applying client | The read is made by the client that applies | the clause of that requirement body |
| The install admission input is not set | The controller sets no install admission input in an apply verdict | the clause of that requirement body |
| One identity for both kinds | The identity of the apply verdict is the same for a ModuleInstance and a ModulePackage | the clause of that requirement body |
| The verdict identity is the render's or the recorded one | The verdict is asked with the render's or recorded identity | the clause of that requirement body |
| Verdict identity when none is recorded | The verdict identity when none is recorded or a change is unsettled | the clause of that requirement body |
| No earlier identity while a change is unsettled | The earlier identity of an unsettled change is not asked with | the clause of that requirement body |
| No empty identity and no prune list | The verdict is not asked with an empty identity or prune list | the clause of that requirement body |
| Objects to apply and no identity | A reconcile with objects and no identity fails | the clause of that requirement body |
| Empty apply list | A reconcile with no objects asks no verdict | the clause of that requirement body |
| A terminating object stops all writes | An object refused as terminating stops all writes | the clause of that requirement body |
| A foreign object stops all writes | An object refused as foreign-object stops all writes | the clause of that requirement body |
| Another instance's object stops all writes | An object refused as other-instance stops all writes | the clause of that requirement body |
| The refused object is untouched | A refused object is left as it is | the clause of that requirement body |
| An adopted-elsewhere object is skipped | An object adopted elsewhere is not applied and the others are | the clause of that requirement body |
| A taken-in object keeps its UID | A taken-in object is not deleted and created again | the clause of that requirement body |
| A taken-in object fails before the first write | A taken-in object whose update is refused as immutable fails before the first write | the clause of that requirement body |
| The message names the object and the fields | The failure of a taken-in object names the object and the refused fields | the clause of that requirement body |
| An unreadable object writes nothing | An unreadable object stops the apply | the clause of that requirement body |
| A kind a CRD in the list defines counts as absent | An object whose kind a CRD of the same apply list defines counts as an object that does not exist | the clause of that requirement body |

Clause table (old clause, new requirement, sentence of its body; "Edit" lists the words that differ from the old sentence):

| Old requirement | Old clause | New requirement | Sentence | Edit |
| --- | --- | --- | --- | --- |
| Custom resources wait for discovery o... | When a staged apply fails only because the API server does not serv... | Custom resources wait for discovery of a CRD in the same set | 1 | none |
| Custom resources wait for discovery o... | A CRD can report `Established` before API discovery serves its kind. | Custom resources wait for discovery of a CRD in the same set | 2 | none |
| Custom resources wait for discovery o... | In the retry that the requirement "Custom resources wait for discov... | The discovery retry runs at a fixed interval within a bound | 1 | the subject "It" is replaced by the subject of sentence 0; the retry is named by the exact header of the requirement that holds its condition; added: "In the retry that the requirement "Custom resources wait for discovery of a CRD in the same set" states, the " |
| Custom resources wait for discovery o... | When a staged apply fails only because the API server does not serv... | The discovery retry covers no other error | 1 | none |
| Custom resources wait for discovery o... | It MUST NOT retry any other error. | The discovery retry covers no other error | 2 | none |
| Custom resources wait for discovery o... | In the retry that the requirement "Custom resources wait for discov... | Apply returns the no-match error when the retry ends | 1 | the subject "It" is replaced by the subject of sentence 0; the retry is named by the exact header of the requirement that holds its condition; added: "In the retry that the requirement "Custom resources wait for discovery of a CRD in the same set" states, the " |
| Custom resources wait for discovery o... | When the bound or the caller's context ends while the error is stil... | Apply returns the no-match error when the retry ends | 2 | none |
| A forced recreate keeps PersistentVol... | The apply MUST take, beside `force`, an option that allows the dele... | The data deletion option is off at its zero value and has... | 1 | none |
| A forced recreate keeps PersistentVol... | When `force` is false, the option has no effect. | The data deletion option is off at its zero value and has... | 2 | none |
| A forced recreate keeps PersistentVol... | The apply MUST take, beside `force`, an option that allows the dele... | A forced apply with the option off never deletes a Persis... | 1 | none |
| A forced recreate keeps PersistentVol... | When `force` is true and the option is off, the apply MUST NOT dele... | A forced apply with the option off never deletes a Persis... | 2 | none |
| A forced recreate keeps PersistentVol... | The apply MUST take, beside `force`, an option that allows the dele... | The apply checks every claim of the set before it applies... | 1 | none |
| A forced recreate keeps PersistentVol... | When `force` is true and the option is off, the apply MUST NOT dele... | The apply checks every claim of the set before it applies... | 2 | none |
| A forced recreate keeps PersistentVol... | Before it applies or deletes any object, it MUST check every such c... | The apply checks every claim of the set before it applies... | 3 | none |
| A forced recreate keeps PersistentVol... | The apply MUST take, beside `force`, an option that allows the dele... | A claim whose update the API server refuses returns a nam... | 1 | sentence 0 is cut at ", and the option MUST be off" so that "the option" has its referent; added: "." |
| A forced recreate keeps PersistentVol... | When `force` is true and the option is off, the apply MUST NOT dele... | A claim whose update the API server refuses returns a nam... | 2 | none |
| A forced recreate keeps PersistentVol... | When the API server refuses the update of one (the error the forced... | A claim whose update the API server refuses returns a nam... | 3 | none |
| A forced recreate keeps PersistentVol... | In that case no object of the resource set is applied, created or d... | A refused claim update applies, creates and deletes no ob... | 1 | +SHALL (Strength table); the case is named by the exact header of the requirement that holds sentences 2 and 3 in full |
| A forced recreate keeps PersistentVol... | The apply MUST take, beside `force`, an option that allows the dele... | A claim is recreated when data deletion is allowed | 1 | none |
| A forced recreate keeps PersistentVol... | When the option is on, a forced recreate MUST delete and recreate a... | A claim is recreated when data deletion is allowed | 2 | none |
| A forced recreate keeps PersistentVol... | The apply MUST take, beside `force`, an option that allows the dele... | The forced recreate of every other kind is unchanged | 1 | none |
| A forced recreate keeps PersistentVol... | When `force` is true and the option is off, the apply MUST NOT dele... | The forced recreate of every other kind is unchanged | 2 | none |
| A forced recreate keeps PersistentVol... | The forced recreate of every other kind MUST delete and recreate as... | The forced recreate of every other kind is unchanged | 3 | none |
| A forced recreate keeps PersistentVol... | The resource manager MUST refuse on its own to delete a core Persis... | The resource manager refuses to delete a core PersistentV... | 1 | none |
| A forced recreate deletes exactly the... | When an apply with `force` deletes an object to create it again, th... | A forced recreate deletes exactly the object that was read | 1 | none |
| A forced recreate deletes exactly the... | The apply MUST NOT ask an ownership verdict for that delete. | A forced recreate deletes exactly the object that was read | 2 | none |
| A forced recreate deletes exactly the... | A DELETE the API server refuses on the precondition MUST fail the a... | A forced recreate deletes exactly the object that was read | 3 | none |
| A forced recreate deletes exactly the... | Source: owner decision of 2026-10-08 (the forced recreate carries t... | A forced recreate deletes exactly the object that was read | 4 | none |
| A forced recreate deletes exactly the... | When an apply with `force` deletes an object to create it again, th... | The PersistentVolumeClaim rule is checked before the UID ... | 1 | none |
| A forced recreate deletes exactly the... | The rule for PersistentVolumeClaims is checked first and is not cha... | The PersistentVolumeClaim rule is checked before the UID ... | 2 | none |
| A forced recreate deletes exactly the... | Source: owner decision of 2026-10-08 (the forced recreate carries t... | The PersistentVolumeClaim rule is checked before the UID ... | 3 | none |
| A forced recreate deletes exactly the... | The resource manager's client MUST refuse a delete of a collection ... | The resource manager client refuses a delete of a collection | 1 | none |
| A forced recreate deletes exactly the... | Source: owner decision of 2026-10-08 (the forced recreate carries t... | The resource manager client refuses a delete of a collection | 2 | none |
| The staged apply order never contradi... | The order in which the staged apply writes two objects of different... | The staged apply order is never the opposite of the libra... | 1 | none |
| The staged apply order never contradi... | Source: enhancement 0012 ("Ordering is kind-class order in the Kube... | The staged apply order is never the opposite of the libra... | 2 | none |
| The staged apply order never contradi... | When Flux's staged apply, at the version the operator pins, applies... | A kind Flux applies first never weighs more than the kind... | 1 | none |
| The staged apply order never contradi... | Source: enhancement 0012 ("Ordering is kind-class order in the Kube... | A kind Flux applies first never weighs more than the kind... | 2 | none |
| The staged apply order never contradi... | The order in which the staged apply writes two objects of different... | Flux may refine the library order | 1 | none |
| The staged apply order never contradi... | Flux's order MAY refine the library's order: it MAY put two kinds t... | Flux may refine the library order | 2 | none |
| The staged apply order never contradi... | Source: enhancement 0012 ("Ordering is kind-class order in the Kube... | Flux may refine the library order | 3 | none |
| The staged apply order never contradi... | The order in which the staged apply writes two objects of different... | The test suite fails when the pinned versions break the o... | 1 | none |
| The staged apply order never contradi... | The operator's test suite MUST fail when the two pinned versions br... | The test suite fails when the pinned versions break the o... | 2 | none |
| The staged apply order never contradi... | Source: enhancement 0012 ("Ordering is kind-class order in the Kube... | The test suite fails when the pinned versions break the o... | 3 | none |
| The staged apply order never contradi... | The order in which the staged apply writes two objects of different... | The check reads the order from the pinned Flux module | 1 | none (copied in front so that "this rule" has its referent) |
| The staged apply order never contradi... | The operator's test suite MUST fail when the two pinned versions br... | The check reads the order from the pinned Flux module | 2 | none |
| The staged apply order never contradi... | The check MUST read Flux's order from the pinned Flux module itself... | The check reads the order from the pinned Flux module | 3 | none |
| The staged apply order never contradi... | Source: enhancement 0012 ("Ordering is kind-class order in the Kube... | The check reads the order from the pinned Flux module | 4 | none (a `**Source**:` line) |
| The staged apply order never contradi... | The order in which the staged apply writes two objects of different... | The check covers the kinds of both orders and the kinds o... | 1 | none (copied in front so that "this rule" has its referent) |
| The staged apply order never contradi... | The operator's test suite MUST fail when the two pinned versions br... | The check covers the kinds of both orders and the kinds o... | 2 | none |
| The staged apply order never contradi... | The check MUST cover every kind Flux's order names, every kind the ... | The check covers the kinds of both orders and the kinds o... | 3 | none |
| The staged apply order never contradi... | Source: enhancement 0012 ("Ordering is kind-class order in the Kube... | The check covers the kinds of both orders and the kinds o... | 4 | none (a `**Source**:` line) |
| The staged apply order never contradi... | These cases are not a failure by themselves, and the check MUST rep... | The check reports two kinds of equal weight that Flux orders | 1 | list cut at the three cases; the head is repeated; added: "." |
| The staged apply order never contradi... | Source: enhancement 0012 ("Ordering is kind-class order in the Kube... | The check reports two kinds of equal weight that Flux orders | 2 | none |
| The staged apply order never contradi... | These cases are not a failure by themselves, and the check MUST rep... | The check reports a kind that only one of the two orders ... | 1 | list cut at the three cases; the head is repeated; added: "." |
| The staged apply order never contradi... | Source: enhancement 0012 ("Ordering is kind-class order in the Kube... | The check reports a kind that only one of the two orders ... | 2 | none |
| The staged apply order never contradi... | These cases are not a failure by themselves, and the check MUST rep... | The check reports two objects of the same group and kind | 1 | list cut at the three cases; the head is repeated |
| The staged apply order never contradi... | Source: enhancement 0012 ("Ordering is kind-class order in the Kube... | The check reports two objects of the same group and kind | 2 | none |
| The staged apply order never contradi... | A failure MUST name every pair of kinds that the two orders put the... | A failure names every pair of kinds in opposite order | 1 | none |
| The staged apply order never contradi... | Source: enhancement 0012 ("Ordering is kind-class order in the Kube... | A failure names every pair of kinds in opposite order | 2 | none |
| The staged apply order never contradi... | A failure MUST tell the maintainer that the library's weight table ... | A failure tells the maintainer what changes | 1 | the subject "It" is replaced by the subject of the sentence before; added: "A failure" |
| The staged apply order never contradi... | Source: enhancement 0012 ("Ordering is kind-class order in the Kube... | A failure tells the maintainer what changes | 2 | none |
| An apply is judged by the ownership v... | Before the first write of a reconcile, the controller MUST read eve... | An apply is judged by the ownership verdict before its fi... | 1 | list cut at "and with whether"; the head is kept; added: "." |
| An apply is judged by the ownership v... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | An apply is judged by the ownership verdict before its fi... | 2 | none |
| An apply is judged by the ownership v... | Before the first write of a reconcile, the controller MUST read eve... | The controller asks the apply verdict whether status.inve... | 1 | list cut at "and with whether"; the head is repeated; added: ", " / "." |
| An apply is judged by the ownership v... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | The controller asks the apply verdict whether status.inve... | 2 | none |
| An apply is judged by the ownership v... | Before the first write of a reconcile, the controller MUST read eve... | The apply verdict is asked for a ModuleInstance and a Mod... | 1 | the head of sentence 0 is copied as the subject of "This holds"; added: "." |
| An apply is judged by the ownership v... | This holds for a ModuleInstance and a ModulePackage, for the apply ... | The apply verdict is asked for a ModuleInstance and a Mod... | 2 | none |
| An apply is judged by the ownership v... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | The apply verdict is asked for a ModuleInstance and a Mod... | 3 | none |
| An apply is judged by the ownership v... | The read MUST be made by the client that applies: the impersonated ... | The read is made by the client that applies | 1 | none |
| An apply is judged by the ownership v... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | The read is made by the client that applies | 2 | none |
| An apply is judged by the ownership v... | The controller MUST NOT set the verdict's install admission input. | The controller sets no install admission input in an appl... | 1 | none |
| An apply is judged by the ownership v... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | The controller sets no install admission input in an appl... | 2 | none |
| An apply is judged by the ownership v... | The identity the verdict is asked with is the same for a ModuleInst... | The identity of the apply verdict is the same for a Modul... | 1 | "this paragraph" becomes "requirements" and the exact headers of the six requirements that hold the rest of the old paragraph (the six headers are shortened so that the list fits 500); "is the same" becomes "MUST be the same" (Strength table) |
| An apply is judged by the ownership v... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | The identity of the apply verdict is the same for a Modul... | 2 | none (a `**Source**:` line) |
| An apply is judged by the ownership v... | The identity the verdict is asked with is the same for a ModuleInst... | The verdict is asked with the render's or recorded identity | 1 | sentence 5 is cut at ", and this paragraph" so that "the verdict" has its referent; added: "." |
| An apply is judged by the ownership v... | The controller MUST ask once per object, with the instance's identi... | The verdict is asked with the render's or recorded identity | 2 | none |
| An apply is judged by the ownership v... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | The verdict is asked with the render's or recorded identity | 3 | none |
| An apply is judged by the ownership v... | The identity the verdict is asked with is the same for a ModuleInst... | The verdict identity when none is recorded or a change is unsettled | 1 | sentence 5 is cut at ", and this paragraph" so that "the verdict" has its referent; added: "." |
| An apply is judged by the ownership v... | The controller MUST ask once per object, with the instance's identity. | The verdict identity when none is recorded or a change is unsettled | 2 | none |
| An apply is judged by the ownership v... | This holds when no identity is recorded yet and while an identity c... | The verdict identity when none is recorded or a change is unsettled | 3 | none |
| An apply is judged by the ownership v... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | The verdict identity when none is recorded or a change is unsettled | 4 | none |
| An apply is judged by the ownership v... | The controller MUST NOT ask with the earlier identity of an unsettl... | The earlier identity of an unsettled change is not asked ... | 1 | none |
| An apply is judged by the ownership v... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | The earlier identity of an unsettled change is not asked ... | 2 | none |
| An apply is judged by the ownership v... | The controller MUST NOT ask the apply verdict with an empty identit... | The verdict is not asked with an empty identity or prune list | 1 | none |
| An apply is judged by the ownership v... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | The verdict is not asked with an empty identity or prune list | 2 | none |
| An apply is judged by the ownership v... | A reconcile that has one or more objects to apply and no identity a... | A reconcile with objects and no identity fails | 1 | none |
| An apply is judged by the ownership v... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | A reconcile with objects and no identity fails | 2 | none |
| An apply is judged by the ownership v... | A reconcile whose apply list is empty has no object to judge: it as... | A reconcile with no objects asks no verdict | 1 | none |
| An apply is judged by the ownership v... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | A reconcile with no objects asks no verdict | 2 | none |
| An apply is judged by the ownership v... | When the verdict refuses an object as `terminating`, and the reconc... | An object refused as terminating stops all writes | 1 | list cut at the three reasons; the head is repeated |
| An apply is judged by the ownership v... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | An object refused as terminating stops all writes | 2 | none |
| An apply is judged by the ownership v... | When the verdict refuses an object as `foreign-object`, and the rec... | An object refused as foreign-object stops all writes | 1 | list cut at the three reasons; the head is repeated |
| An apply is judged by the ownership v... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | An object refused as foreign-object stops all writes | 2 | none |
| An apply is judged by the ownership v... | When the verdict refuses an object as `other-instance`, and the rec... | An object refused as other-instance stops all writes | 1 | list cut at the three reasons; the head is repeated |
| An apply is judged by the ownership v... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | An object refused as other-instance stops all writes | 2 | none |
| An apply is judged by the ownership v... | When the verdict refuses an object as `terminating`, `foreign-objec... | A refused object is left as it is | 1 | the head of sentence 12 is copied as context; the subject is kept; added: ", the" |
| An apply is judged by the ownership v... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | A refused object is left as it is | 2 | none |
| An apply is judged by the ownership v... | When the verdict refuses an object as `adopted-elsewhere`, the cont... | An object adopted elsewhere is not applied and the others... | 1 | none |
| An apply is judged by the ownership v... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | An object adopted elsewhere is not applied and the others... | 2 | none |
| An apply is judged by the ownership v... | An object that the verdict allows, that exists and that is not in `... | A taken-in object is not deleted and created again | 1 | none |
| An apply is judged by the ownership v... | The controller MUST NOT delete and create such an object again in t... | A taken-in object is not deleted and created again | 2 | sentence 16 is cut at its colon; added: "." |
| An apply is judged by the ownership v... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | A taken-in object is not deleted and created again | 3 | none |
| An apply is judged by the ownership v... | An object that the verdict allows, that exists and that is not in `... | A taken-in object whose update is refused as immutable fa... | 1 | none |
| An apply is judged by the ownership v... | With `spec.rollout.forceConflicts`, when the API server refuses its... | A taken-in object whose update is refused as immutable fa... | 2 | sentence 16 is cut at its colon and at ", and the message"; added: "W" / "." |
| An apply is judged by the ownership v... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | A taken-in object whose update is refused as immutable fa... | 3 | none |
| An apply is judged by the ownership v... | An object that the verdict allows, that exists and that is not in `... | The failure of a taken-in object names the object and the... | 1 | none |
| An apply is judged by the ownership v... | With `spec.rollout.forceConflicts`, when the API server refuses its... | The failure of a taken-in object names the object and the... | 2 | sentence 16 is cut at its colon and at ", and the message"; added: "W" / ", " / "." |
| An apply is judged by the ownership v... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | The failure of a taken-in object names the object and the... | 3 | none |
| An apply is judged by the ownership v... | When the read of an object fails for a reason other than that the o... | An unreadable object stops the apply | 1 | none |
| An apply is judged by the ownership v... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | An unreadable object stops the apply | 2 | none |
| An apply is judged by the ownership v... | An object whose kind the API server does not serve yet, while a Cus... | An object whose kind a CRD of the same apply list defines... | 1 | keyword added where the sentence has none; added: "SHALL count" |
| An apply is judged by the ownership v... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | An object whose kind a CRD of the same apply list defines... | 2 | none |

Lengths (old body, new bodies; characters by the guide's rule):

| Old requirement | Old | New |
| --- | --- | --- |
| Custom resources wait for discovery of a CRD in the same set | 766 | 367, 333, 333, 494 |
| A forced recreate keeps PersistentVolumeClaims unless data deletion is allowed | 1256 | 175, 243, 363, 467, 378, 229, 422, 216 |
| A forced recreate deletes exactly the object that was read | 680 | 449, 314, 251 |
| The staged apply order never contradicts the library's kind weights | 1917 | 322, 390, 444, 404, 486, 465, 324, 435, 351, 293, 394 |
| An apply is judged by the ownership verdict before its first write | 3220 | 417, 482, 493, 354, 275, 478, 436, 461, 428, 437, 331, 331, 477, 480, 480, 341, 343, 466, 488, 452, 383, 385 |

### Ledger: library-kernel-runtime

Form c2 for every requirement below: REMOVED under the old name, ADDED under new names, new blocks at the end of the spec. `grep -rn --exclude-dir=archive` for each old header found no reference to fix. Old scenario text is copied by script from `origin/main`; item 3 of the guide checklist shows `same` for each.

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| Kernel constructed once at startup | The manager constructs one library Kernel for the process lifetime, Kernel constructed once at startup | no (byte equal, checked by the item 3 command) |
| Kernel survives across reconciles | The manager constructs one library Kernel for the process lifetime, Kernel survives across reconciles | no (byte equal, checked by the item 3 command) |
| Two renders overlap | The Kernel is safe for concurrent use and no kernel call is serialised, Two renders overlap | no (byte equal, checked by the item 3 command) |
| No kernel gate | No lock or ordering gate is held across a kernel call, No kernel gate | no (byte equal, checked by the item 3 command) |
| Default keeps reconciles serial | The manager takes --max-concurrent-renders and builds one pool of render slots, Default keeps reconciles serial | no (byte equal, checked by the item 3 command) |
| Raising the bound allows overlap | The manager takes --max-concurrent-renders and builds one pool of render slots, Raising the bound allows overlap | no (byte equal, checked by the item 3 command) |
| Both kinds share the slots | The manager takes --max-concurrent-renders and builds one pool of render slots, Both kinds share the slots | no (byte equal, checked by the item 3 command) |
| Conversion runs while the slot is held | A reconcile holds a render slot from the renderer call to the dropped values, Conversion runs while the slot is held | no (byte equal, checked by the item 3 command) |
| A panicking render frees its slot | A reconcile holds a render slot from the renderer call to the dropped values, A panicking render frees its slot | no (byte equal, checked by the item 3 command) |
| Shutdown while waiting for a slot | A reconcile whose wait for a slot is cancelled returns the context error, Shutdown while waiting for a slot | no (byte equal, checked by the item 3 command) |
| A ModuleInstance reconcile drops the rendered resources | The reconcilers drop the rendered CUE-backed resources after the export, A ModuleInstance reconcile drops the rendered resources | no (byte equal, checked by the item 3 command) |
| A ModulePackage reconcile drops the rendered resources | The reconcilers drop the rendered CUE-backed resources after the export, A ModulePackage reconcile drops the rendered resources | no (byte equal, checked by the item 3 command) |
| Default timeout | The manager takes --render-timeout and refuses a negative value at startup, Default timeout | no (byte equal, checked by the item 3 command) |
| Negative timeout is refused | The manager takes --render-timeout and refuses a negative value at startup, Negative timeout is refused | no (byte equal, checked by the item 3 command) |
| Queued renders are not timed out | A render runs under the render deadline while it holds its slot, Queued renders are not timed out | no (byte equal, checked by the item 3 command) |
| A render that finishes in time is unaffected | A render runs under the render deadline while it holds its slot, A render that finishes in time is unaffected | no (byte equal, checked by the item 3 command) |
| Zero disables the deadline | A render timeout of zero disables the deadline, Zero disables the deadline | no (byte equal, checked by the item 3 command) |
| The status patch outlives the render deadline | A reconcile stops waiting for a render that passes its deadline, The status patch outlives the render deadline | no (byte equal, checked by the item 3 command) |
| A timed-out render keeps its slot until it returns | An abandoned render keeps its slot until it returns, A timed-out render keeps its slot until it returns | no (byte equal, checked by the item 3 command) |
| A ModulePackage render past its deadline keeps its files | A reconcile does not read what an abandoned render writes and the render keeps what it reads, A ModulePackage render past its deadline keeps its files | no (byte equal, checked by the item 3 command) |
| A hung object holds at most one slot | A timed-out render that is still running blocks a second render of the object, A hung object holds at most one slot | no (byte equal, checked by the item 3 command) |
| A panic before the deadline keeps its stack | A panic in the render before its deadline reaches the reconcile, A panic before the deadline keeps its stack | no (byte equal, checked by the item 3 command) |
| A panic after the deadline does not stop the operator | A panic after the reconcile stopped waiting does not stop the process, A panic after the deadline does not stop the operator | no (byte equal, checked by the item 3 command) |

New scenarios (restate a clause of the old body, with words of the old text only):

| New scenario | Requirement | Old clause restated |
| --- | --- | --- |
| No rendered value without a slot | No reconcile holds a rendered CUE value without a slot | the clause of that requirement body |
| Reconciles outside the render are not serialised | Each controller uses the render flag as its maximum concurrent reconciles | the clause of that requirement body |
| The Platform controller takes no slot | The Platform controller stays serial and takes no slot | the clause of that requirement body |
| The flag help states the sizing rule | The flag help states the memory sizing rule | the clause of that requirement body |
| Later phases read exported data | Later phases read the exported data and never the dropped resources | the clause of that requirement body |
| A late render is discarded | A render that returns after its deadline is treated as timed out | the clause of that requirement body |
| A timed-out render is logged | The operator logs a timed-out render and an abandoned render that returns | the clause of that requirement body |
| The timeout flag help | The render timeout flag help states that the slot stays held | the clause of that requirement body |

Clause table (old clause, new requirement, sentence of its body; "Edit" lists the words that differ from the old sentence):

| Old requirement | Old clause | New requirement | Sentence | Edit |
| --- | --- | --- | --- | --- |
| Single long-lived library Kernel | The manager SHALL construct one library Kernel for the process life... | The manager constructs one library Kernel for the process... | 1 | none |
| Single long-lived library Kernel | The Kernel is safe for concurrent use across its methods: every ker... | The Kernel is safe for concurrent use and no kernel call ... | 1 | none |
| Single long-lived library Kernel | No lock or ordering gate SHALL be held across a kernel call: a rend... | No lock or ordering gate is held across a kernel call | 1 | none |
| Single long-lived library Kernel | Concurrency is bounded by that pool on the render paths, and to one... | No lock or ordering gate is held across a kernel call | 2 | none |
| Render concurrency is a manager flag ... | The manager SHALL accept `--max-concurrent-renders` (integer, defau... | The manager takes --max-concurrent-renders and builds one... | 1 | none |
| Render concurrency is a manager flag ... | A reconcile SHALL take a slot before it calls its renderer (platfor... | A reconcile holds a render slot from the renderer call to... | 1 | none |
| Render concurrency is a manager flag ... | The manager SHALL accept `--max-concurrent-renders` (integer, defau... | No reconcile holds a rendered CUE value without a slot | 1 | none |
| Render concurrency is a manager flag ... | So no reconcile holds a rendered CUE value without a slot, and the ... | No reconcile holds a rendered CUE value without a slot | 2 | none |
| Render concurrency is a manager flag ... | The manager SHALL accept `--max-concurrent-renders` (integer, defau... | Each controller uses the render flag as its maximum concu... | 1 | none |
| Render concurrency is a manager flag ... | Each of the two controllers SHALL also use the flag as its maximum ... | Each controller uses the render flag as its maximum concu... | 2 | none |
| Render concurrency is a manager flag ... | A reconcile whose wait for a slot ends because its context is cance... | A reconcile whose wait for a slot is cancelled returns th... | 1 | none |
| Render concurrency is a manager flag ... | The Platform controller SHALL stay serial and SHALL take no slot. | The Platform controller stays serial and takes no slot | 1 | none |
| Render concurrency is a manager flag ... | The manager SHALL accept `--max-concurrent-renders` (integer, defau... | The flag help states the memory sizing rule | 1 | none |
| Render concurrency is a manager flag ... | The flag's help SHALL state the memory sizing rule (per-render cost... | The flag help states the memory sizing rule | 2 | none |
| Rendered values are dropped after con... | A rendered resource carries the CUE value it was read from, and a h... | The reconcilers drop the rendered CUE-backed resources af... | 1 | none |
| Rendered values are dropped after con... | The ModuleInstance and ModulePackage reconcilers SHALL therefore dr... | The reconcilers drop the rendered CUE-backed resources af... | 2 | none |
| Rendered values are dropped after con... | Every later phase (shrink judgment, apply, prune, inventory and sta... | Later phases read the exported data and never the dropped... | 1 | none |
| A render is bounded by a manager time... | The manager SHALL accept `--render-timeout` (a duration, default `1... | The manager takes --render-timeout and refuses a negative... | 1 | none |
| A render is bounded by a manager time... | The manager SHALL accept `--render-timeout` (a duration, default `1... | A render runs under the render deadline while it holds it... | 1 | sentence 0 is cut at "and SHALL refuse" so that "that timeout" has its referent; the head is kept; added: "." |
| A render is bounded by a manager time... | While a ModuleInstance or ModulePackage reconcile holds its render ... | A render runs under the render deadline while it holds it... | 2 | none |
| A render is bounded by a manager time... | The wait for a slot SHALL NOT count against it, and the reconcile's... | A render runs under the render deadline while it holds it... | 3 | none |
| A render is bounded by a manager time... | The manager SHALL accept `--render-timeout` (a duration, default `1... | A render timeout of zero disables the deadline | 1 | none |
| A render is bounded by a manager time... | A value of `0` SHALL disable the deadline, and the render then runs... | A render timeout of zero disables the deadline | 2 | none |
| A render is bounded by a manager time... | When the deadline passes before the render returns, the reconcile S... | A reconcile stops waiting for a render that passes its de... | 1 | none |
| A render is bounded by a manager time... | When the deadline passes before the render returns, the reconcile S... | An abandoned render keeps its slot until it returns | 1 | none |
| A render is bounded by a manager time... | The render SHALL keep its slot until it really returns (at the next... | An abandoned render keeps its slot until it returns | 2 | none |
| A render is bounded by a manager time... | A render that returns after its deadline SHALL be treated as timed ... | A render that returns after its deadline is treated as ti... | 1 | none |
| A render is bounded by a manager time... | A reconcile SHALL NOT read anything the abandoned render writes, an... | A reconcile does not read what an abandoned render writes... | 1 | none |
| A render is bounded by a manager time... | While a timed-out render of an object is still running, a reconcile... | A timed-out render that is still running blocks a second ... | 1 | none |
| A render is bounded by a manager time... | One object whose render hangs therefore holds at most one slot, wha... | A timed-out render that is still running blocks a second ... | 2 | none |
| A render is bounded by a manager time... | A panic in the render before its deadline SHALL be logged at error ... | A panic in the render before its deadline reaches the rec... | 1 | none |
| A render is bounded by a manager time... | A panic after the reconcile stopped waiting SHALL be logged at erro... | A panic after the reconcile stopped waiting does not stop... | 1 | none |
| A render is bounded by a manager time... | The operator SHALL log a timed-out render with its timeout, and SHA... | The operator logs a timed-out render and an abandoned ren... | 1 | none |
| A render is bounded by a manager time... | The manager SHALL accept `--render-timeout` (a duration, default `1... | The render timeout flag help states that the slot stays held | 1 | none |
| A render is bounded by a manager time... | The flag's help SHALL state that the slot stays held until the rend... | The render timeout flag help states that the slot stays held | 2 | none |

Lengths (old body, new bodies; characters by the guide's rule):

| Old requirement | Old | New |
| --- | --- | --- |
| Single long-lived library Kernel | 881 | 104, 378, 397 |
| Render concurrency is a manager flag bounded by memory | 1422 | 287, 363, 402, 458, 223, 65, 482 |
| Rendered values are dropped after conversion | 752 | 468, 283 |
| A render is bounded by a manager timeout that keeps its slot | 2491 | 117, 483, 223, 182, 412, 106, 339, 349, 248, 153, 129, 231 |

### Ledger: module-instance-synthesis

Form c2 for every requirement below: REMOVED under the old name, ADDED under new names, new blocks at the end of the spec. `grep -rn --exclude-dir=archive` for each old header found no reference to fix. Old scenario text is copied by script from `origin/main`; item 3 of the guide checklist shows `same` for each. The scenario "Unresolved platform demands reported as resolution failure" holds an em-dash in its WHEN line on `origin/main`; it moves byte for byte, because item 3 compares bytes (see the question in the report).

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| Env fallback when flag empty | The controller accepts a --registry flag and reads OPM_REGISTRY when the flag is empty, Env fallback when flag empty | no (byte equal, checked by the item 3 command) |
| Flag value wins | The controller sets CUE_REGISTRY and OPM_REGISTRY before any CUE module evaluation, Flag value wins | no (byte equal, checked by the item 3 command) |
| Built-in default applied | The controller ships a built-in default registry mapping, Built-in default applied | no (byte equal, checked by the item 3 command) |
| CUE default resolution | The registry precedence is the flag, then OPM_REGISTRY, then the CUE default, CUE default resolution | no (byte equal, checked by the item 3 command) |
| Success reported | The status source and the Ready condition report a resolved module, Success reported | no (byte equal, checked by the item 3 command) |
| Resolution failure reported | A module that cannot be acquired or fetched is a ResolutionFailed, Resolution failure reported | no (byte equal, checked by the item 3 command) |
| Identity mismatch reported as resolution failure | A module with a wrong identity, an invalid artifact or unmet platform demands is a ResolutionFailed, Identity mismatch reported as resolution failure | no (byte equal, checked by the item 3 command) |
| Unresolved platform demands reported as resolution failure | A module with a wrong identity, an invalid artifact or unmet platform demands is a ResolutionFailed, Unresolved platform demands reported as resolution failure | no (byte equal, checked by the item 3 command) |
| Render failure reported | A synthesis, evaluation or render failure of another cause is a RenderFailed, Render failure reported | no (byte equal, checked by the item 3 command) |
| Registry failure after acquisition reported as transient resolution failure | A failure that is not transient sets Stalled and a typed registry fetch failure does not, Registry failure after acquisition reported as transient resolution failure | no (byte equal, checked by the item 3 command) |
| Unset required value that no component reads | An unset required config value is refused by instance synthesis and the controller adds no check, Unset required value that no component reads | no (byte equal, checked by the item 3 command) |
| No spec.values at all | An unset required config value is refused by instance synthesis and the controller adds no check, No spec.values at all | no (byte equal, checked by the item 3 command) |
| Two unset required values | An unset required config value is refused by instance synthesis and the controller adds no check, Two unset required values | no (byte equal, checked by the item 3 command) |
| The value is set | An unset required config value is refused by instance synthesis and the controller adds no check, The value is set | no (byte equal, checked by the item 3 command) |
| Three unset values, one read by a component, one nested | A value a component reads is named like a value no component reads, Three unset values, one read by a component, one nested | no (byte equal, checked by the item 3 command) |
| Setting one value removes its finding only | A value a component reads is named like a value no component reads, Setting one value removes its finding only | no (byte equal, checked by the item 3 command) |
| Two values of the wrong type | A values failure message carries every finding while it fits the condition limit, Two values of the wrong type | no (byte equal, checked by the item 3 command) |
| A value of the wrong type that a component reads | A finding caused by a value in spec.values names its position, A value of the wrong type that a component reads | no (byte equal, checked by the item 3 command) |
| A value of the wrong type that no component reads | A finding caused by a value in spec.values names its position, A value of the wrong type that no component reads | no (byte equal, checked by the item 3 command) |
| A constraint a value does not meet | A finding caused by a value in spec.values names its position, A constraint a value does not meet | no (byte equal, checked by the item 3 command) |
| A field the module does not declare | A finding caused by a value in spec.values names its position, A field the module does not declare | no (byte equal, checked by the item 3 command) |
| A registry fetch failure during synthesis keeps its class | The wording of a values failure does not change how it is classified, A registry fetch failure during synthesis keeps its class | no (byte equal, checked by the item 3 command) |

New scenarios (restate a clause of the old body, with words of the old text only):

| New scenario | Requirement | Old clause restated |
| --- | --- | --- |
| The message form | The refusal message has a fixed form with one finding per unset required value | the clause of that requirement body |
| The findings come from the typed error tree | The controller words the findings from the kernel typed error tree | the clause of that requirement body |
| The unset value reports RenderFailed and Stalled | A ModuleInstance with an unset required value reports RenderFailed and Stalled | the clause of that requirement body |
| The event carries the same text | The event carries the text of the values failure | the clause of that requirement body |
| A registry fetch failure keeps its text | A registry fetch failure during synthesis keeps its own text | the clause of that requirement body |

Clause table (old clause, new requirement, sentence of its body; "Edit" lists the words that differ from the old sentence):

| Old requirement | Old clause | New requirement | Sentence | Edit |
| --- | --- | --- | --- | --- |
| Registry configuration | The controller MUST accept a `--registry` flag for CUE registry con... | The controller accepts a --registry flag and reads OPM_RE... | 1 | none |
| Registry configuration | The controller MUST read the `OPM_REGISTRY` environment variable as... | The controller accepts a --registry flag and reads OPM_RE... | 2 | none |
| Registry configuration | The controller MUST set `CUE_REGISTRY` and `OPM_REGISTRY` to the re... | The controller sets CUE_REGISTRY and OPM_REGISTRY before ... | 1 | none |
| Registry configuration | The controller MUST accept a `--registry` flag for CUE registry con... | The controller ships a built-in default registry mapping | 1 | none |
| Registry configuration | The controller ships a built-in default `--registry` value routing ... | The controller ships a built-in default registry mapping | 2 | none |
| Registry configuration | Operators override by passing an explicit `--registry=<mapping>`, o... | The controller ships a built-in default registry mapping | 3 | none |
| Registry configuration | The controller MUST accept a `--registry` flag for CUE registry con... | The registry precedence is the flag, then OPM_REGISTRY, t... | 1 | the three precedence items keep their numbers as list lines |
| Status reporting | The `status.source` field MAY be updated to reflect:
- The CUE modu... | The status source and the Ready condition report a resolv... | 1 | the list lines keep their dashes |
| Status reporting | The `status.conditions` MUST report:
- `Ready=False` with reason `R... | A module that cannot be acquired or fetched is a Resoluti... | 1 | the list is cut at its items; the head is repeated |
| Status reporting | The `status.conditions` MUST report:
- `Ready=False` with reason `R... | A module with a wrong identity, an invalid artifact or un... | 1 | the list is cut at its items; the head is repeated |
| Status reporting | The `status.conditions` MUST report:
- `Ready=False` with reason `R... | A synthesis, evaluation or render failure of another caus... | 1 | the list keeps its dash |
| Status reporting | The `status.conditions` MUST report:
- `Stalled=True` when the fail... | A failure that is not transient sets Stalled and a typed ... | 1 | the list lines keep their dashes |
| An unset required config value is ref... | A ModuleInstance whose `spec.values` leave a required `#config` val... | An unset required config value is refused by instance syn... | 1 | none |
| An unset required config value is ref... | The controller SHALL report the kernel's refusal and SHALL NOT run ... | An unset required config value is refused by instance syn... | 2 | none |
| An unset required config value is ref... | The message SHALL be `synthesizing release: Kernel.SynthesizeInstan... | The refusal message has a fixed form with one finding per... | 1 | none |
| An unset required config value is ref... | A nested field SHALL be named by its full path (`values.db.host`). | The refusal message has a fixed form with one finding per... | 2 | none |
| An unset required config value is ref... | The message SHALL be `synthesizing release: Kernel.SynthesizeInstan... | A value a component reads is named like a value no compon... | 1 | none |
| An unset required config value is ref... | A value a component reads SHALL be named like a value no component ... | A value a component reads is named like a value no compon... | 2 | none |
| An unset required config value is ref... | The message SHALL be `synthesizing release: Kernel.SynthesizeInstan... | The controller words the findings from the kernel typed e... | 1 | none |
| An unset required config value is ref... | The controller SHALL word the findings from the kernel's typed erro... | The controller words the findings from the kernel typed e... | 2 | none |
| An unset required config value is ref... | A ModuleInstance whose `spec.values` leave a required `#config` val... | A ModuleInstance with an unset required value reports Ren... | 1 | none |
| An unset required config value is ref... | The ModuleInstance SHALL report `Ready=False` and `Stalled=True` wi... | A ModuleInstance with an unset required value reports Ren... | 2 | none |
| A values failure reports every findin... | When instance synthesis refuses the values of a ModuleInstance, the... | A values failure message carries every finding while it f... | 1 | none |
| A values failure reports every findin... | As listed in requirement "A values failure message carries every fi... | The event carries the text of the values failure | 1 | keyword added where the sentence has none; "the same text" is the text of the requirement named by its exact header, which holds the old sentence whole; added: "As listed in requirement "A values failure message carries every finding while it fits the condition limit": " / "SHALL carry" |
| A values failure reports every findin... | A registry fetch failure during synthesis SHALL keep its own text. | A registry fetch failure during synthesis keeps its own text | 1 | none |
| A values failure reports every findin... | A finding caused by a value in `spec.values` SHALL name its positio... | A finding caused by a value in spec.values names its posi... | 1 | none |
| A values failure reports every findin... | The wording SHALL NOT change how the failure is classified: a value... | The wording of a values failure does not change how it is... | 1 | none |

Lengths (old body, new bodies; characters by the guide's rule):

| Old requirement | Old | New |
| --- | --- | --- |
| Registry configuration | 1014 | 200, 122, 441, 406 |
| Status reporting | 1375 | 282, 322, 462, 214, 392 |
| An unset required config value is refused by instance synthesis | 1065 | 319, 365, 443, 440, 310 |
| A values failure reports every finding with its positions | 1015 | 492, 204, 66, 102, 260 |

### Ledger: modulepackage-reconcile-loop

Form b for "No-op detection", because `reconcile-loop-assembly` cites its header: it stays as a MODIFIED parent with its short body and its three scenarios, and the two ADDED siblings carry new scenarios that restate their clause. Form c2 for the other three requirements below: REMOVED under the old name, ADDED under new names, new blocks at the end of the spec. `grep -rn --exclude-dir=archive` for each of the other old headers found no reference to fix. Old scenario text is copied by script from `origin/main`; item 3 of the guide checklist shows `same` for each. The citation `reconcile-loop-assembly`, "A reconcile the apply verdict refuses is retried and not stalled" names a requirement of `reconcile-loop-assembly`, which a later unit rewrites: that unit must keep the quotation valid. `reconcile-loop-assembly` also quotes the header "An apply is judged by the ownership verdict before its first write", which stays (form b). A body that cannot hold the condition of a sibling next to its `Source:` line names the sibling by its exact header ("as listed in requirement"): the two failure reasons.

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| All digests match | No-op detection, All digests match | no (byte equal, checked by the item 3 command) |
| All digests match and an object was adopted by another instance | No-op detection, All digests match and an object was adopted by another instance | no (byte equal, checked by the item 3 command) |
| Matching digests, an object to take in and an inventoried object being deleted | No-op detection, Matching digests, an object to take in and an inventoried object being deleted | no (byte equal, checked by the item 3 command) |
| Status updated on failure | The ReleaseReconciler patches the status at the end of every reconcile attempt, Status updated on failure | no (byte equal, checked by the item 3 command) |
| A skipped interval patches nothing | A reconcile that skips its render is not an attempt and patches only a changed Healthy condition, A skipped interval patches nothing | no (byte equal, checked by the item 3 command) |
| A skipped interval records a changed health judgement | A reconcile that skips its render is not an attempt and patches only a changed Healthy condition, A skipped interval records a changed health judgement | no (byte equal, checked by the item 3 command) |
| Successful reconcile status | The patch of a successful outcome and of a NoOp carries the judged Healthy condition, Successful reconcile status | no (byte equal, checked by the item 3 command) |
| The identity is recorded | A ModulePackage carries the identity of the instance it renders and the identity before an unsettled change, The identity is recorded | no (byte equal, checked by the item 3 command) |
| An existing package gains the field | A ModulePackage stored before the fields existed gains the identity, An existing package gains the field | no (byte equal, checked by the item 3 command) |
| A package does not delete another instance's object | The prune and the deletion cleanup of a ModulePackage judge with the recorded identities, A package does not delete another instance's object | no (byte equal, checked by the item 3 command) |
| Deletion of a package judges with the recorded identity | The prune and the deletion cleanup of a ModulePackage judge with the recorded identities, Deletion of a package judges with the recorded identity | no (byte equal, checked by the item 3 command) |
| A package that never rendered is deleted as before | A ModulePackage with no recorded identity deletes what it deleted before the field existed, A package that never rendered is deleted as before | no (byte equal, checked by the item 3 command) |
| The first render after the upgrade also changes the identity | A ModulePackage with no recorded identity deletes what it deleted before the field existed, The first render after the upgrade also changes the identity | no (byte equal, checked by the item 3 command) |
| A package is refused on another instance's object | A ModulePackage reconcile judges every rendered object with the instance's identity, A package is refused on another instance's object | no (byte equal, checked by the item 3 command) |
| A package adopts an annotated object | A ModulePackage reconcile judges every rendered object with the instance's identity, A package adopts an annotated object | no (byte equal, checked by the item 3 command) |
| A skipped render reads nothing | The read is made by the client that applies and a skipped render reads nothing, A skipped render reads nothing | no (byte equal, checked by the item 3 command) |
| Matching digests and an unreadable object | The failure reason is ImpersonationFailed when the ServiceAccount is missing or the read is Forbidden, Matching digests and an unreadable object | no (byte equal, checked by the item 3 command) |
| Matching digests and a missing ServiceAccount | The failure reason is ImpersonationFailed when the ServiceAccount is missing or the read is Forbidden, Matching digests and a missing ServiceAccount | no (byte equal, checked by the item 3 command) |

New scenarios (restate a clause of the old body, with words of the old text only):

| New scenario | Requirement | Old clause restated |
| --- | --- | --- |
| A reconcile whose digests match is not a no-op | A reconcile whose digests match applies the objects the verdict allows when one is outside the inventory or another instance adopts an object | the clause of that requirement body |
| Matching digests apply the rendered set | A reconcile whose digests match applies the rendered set and an inventoried object being deleted refuses it | the clause of that requirement body |
| Two cases of a missing deletion-cleanup ServiceAccount are not the DeletionSAMissing stall | "Two cases SHALL NOT be this stall, because ..." | "Two cases are not this stall, because ..." (serviceaccount-impersonation, "Deletion-cleanup SA-missing stalls with distinct reason") |
| Both cases remove the finalizer and warn DeletionUnconfirmed | "In both the finalizer SHALL be removed, the controller's own client SHALL still never be used, and a `Warning` event with reason `DeletionUnconfirmed` SHALL say what could not be read." | "In both the finalizer is removed, the controller's own client is still never used, and a `Warning` event with reason `DeletionUnconfirmed` says what could not be read." (same requirement) |
| The restore of ServiceAccount RBAC alone is not a trigger | "The restore of a ServiceAccount's RBAC alone (a Role or a binding) SHALL NOT be a trigger; ..." | "The restore of a ServiceAccount's RBAC alone (a Role or a binding) is not a trigger; ..." (serviceaccount-impersonation, "A created ServiceAccount wakes the instances that impersonate it") |
| The materialize-era reasons are retired | "The materialize-era reasons (`Materialized`, `MaterializeFailed`) SHALL be retired." | "The materialize-era reasons (`Materialized`, `MaterializeFailed`) are retired." (platform-module-generation, "The generated module is validated by building it") |
| Registry fetch failures are transient wherever they occur | "A failure SHALL be transient when its error chain holds ..." | "A failure is transient when its error chain holds ..." (reconcile-backoff, same requirement) |
| A second identity change is refused | The identity fields are written, kept and cleared as on a ModuleInstance | the clause of that requirement body |
| The earlier identity is not asked with | A ModulePackage reconcile never asks the verdict with the earlier identity of an unsettled change | the clause of that requirement body |
| No empty identity | A ModulePackage reconcile never asks the verdict with an empty identity | the clause of that requirement body |
| Judged as a ModuleInstance reconcile does | A ModulePackage reconcile judges exactly as a ModuleInstance reconcile does under the ssa-apply requirement | the clause of that requirement body |
| A ModulePackage reconcile judges as a ModuleInstance reconcile does under the loop assembly | A ModulePackage reconcile judges exactly as a ModuleInstance reconcile does under the reconcile-loop-assembly requirement | the clause of that requirement body |
| Matching digests and a verdict that could not be asked | A reconcile that renders fails when the verdict could not be asked | the clause of that requirement body |
| Another failure to ask the verdict is ApplyFailed | The failure reason is ApplyFailed on the backoff except when the ServiceAccount is missing or the read is Forbidden | the clause of that requirement body |
| A failed reconcile writes nothing | A failed ModulePackage reconcile writes nothing and keeps the inventory | the clause of that requirement body |

Clause table (old clause, new requirement, sentence of its body; "Edit" lists the words that differ from the old sentence):

| Old requirement | Old clause | New requirement | Sentence | Edit |
| --- | --- | --- | --- | --- |
| No-op detection | The `ReleaseReconciler` MUST detect no-op reconciliations when sour... | No-op detection | 1 | none |
| No-op detection | The `ReleaseReconciler` MUST detect no-op reconciliations when sour... | A reconcile whose digests match applies the objects the v... | 1 | none |
| No-op detection | A reconcile whose digests match is not a no-op when a rendered obje... | A reconcile whose digests match applies the objects the v... | 2 | none |
| No-op detection | A reconcile whose digests match is not a no-op when a rendered obje... | A reconcile whose digests match applies the rendered set ... | 1 | keyword added where the sentence has none; the condition of "Such a reconcile" is copied from sentence 1; added: ". " / "SHALL apply" |
| Status always patched | The `ReleaseReconciler` MUST patch `Release.status` at the end of e... | The ReleaseReconciler patches the status at the end of ev... | 1 | none |
| Status always patched | The status shape mirrors ModuleRelease: conditions, digests, invent... | The ReleaseReconciler patches the status at the end of ev... | 2 | none |
| Status always patched | A reconcile that skips its render because its inputs are unchanged ... | A reconcile that skips its render is not an attempt and p... | 1 | none |
| Status always patched | It judges health (`instance-health`) and MUST NOT patch anything bu... | A reconcile that skips its render is not an attempt and p... | 2 | none |
| Status always patched | It requeues after `spec.interval`, or sooner when health asks for it. | A reconcile that skips its render is not an attempt and p... | 3 | none |
| Status always patched | The `ReleaseReconciler` MUST patch `Release.status` at the end of e... | The patch of a successful outcome and of a NoOp carries t... | 1 | none |
| Status always patched | On a successful outcome and on `NoOp`, the patch also carries the j... | The patch of a successful outcome and of a NoOp carries t... | 2 | none |
| A ModulePackage records its instance ... | A ModulePackage MUST carry the identity of the instance it renders ... | A ModulePackage carries the identity of the instance it r... | 1 | none |
| A ModulePackage records its instance ... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | A ModulePackage carries the identity of the instance it r... | 2 | none |
| A ModulePackage records its instance ... | `status.instanceUUID` and `status.previousInstanceUUID` MUST be wri... | The identity fields are written, kept and cleared as on a... | 1 | the subject "Both fields" is replaced by the two field names of sentence 0; added: "`status.instanceUUID` and `status.previousInstanceUUID`" |
| A ModulePackage records its instance ... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | The identity fields are written, kept and cleared as on a... | 2 | none |
| A ModulePackage records its instance ... | A ModulePackage stored before the fields existed MUST be valid unch... | A ModulePackage stored before the fields existed gains th... | 1 | none |
| A ModulePackage records its instance ... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | A ModulePackage stored before the fields existed gains th... | 2 | none |
| A ModulePackage records its instance ... | The prune of stale resources and the deletion cleanup of a ModulePa... | The prune and the deletion cleanup of a ModulePackage jud... | 1 | none |
| A ModulePackage records its instance ... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | The prune and the deletion cleanup of a ModulePackage jud... | 2 | none |
| A ModulePackage records its instance ... | While no identity is recorded, a ModulePackage MUST delete what it ... | A ModulePackage with no recorded identity deletes what it... | 1 | none |
| A ModulePackage records its instance ... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | A ModulePackage with no recorded identity deletes what it... | 2 | none |
| A ModulePackage guards every apply by... | A ModulePackage reconcile that renders MUST judge every rendered ob... | A ModulePackage reconcile judges every rendered object wi... | 1 | list cut inside the parenthesis at the semicolon; the head is repeated; added: ")" / "." |
| A ModulePackage guards every apply by... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | A ModulePackage reconcile judges every rendered object wi... | 2 | none |
| A ModulePackage guards every apply by... | A ModulePackage reconcile that renders MUST judge every rendered ob... | A ModulePackage reconcile never asks the verdict with the... | 1 | list cut inside the parenthesis at the semicolon; the head is repeated; added: "(" / ")" / "." |
| A ModulePackage guards every apply by... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | A ModulePackage reconcile never asks the verdict with the... | 2 | none |
| A ModulePackage guards every apply by... | A ModulePackage reconcile that renders MUST judge every rendered ob... | A ModulePackage reconcile never asks the verdict with an ... | 1 | list cut inside the parenthesis at the semicolon; the head is repeated; added: "(" / ")" / "." |
| A ModulePackage guards every apply by... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | A ModulePackage reconcile never asks the verdict with an ... | 2 | none |
| A ModulePackage guards every apply by... | A ModulePackage reconcile that renders MUST judge every rendered ob... | A ModulePackage reconcile judges exactly as a ModuleInsta... | 1 | list cut at the two citations; the head and "exactly as" are repeated; added: "(" / ")" / "." |
| A ModulePackage guards every apply by... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | A ModulePackage reconcile judges exactly as a ModuleInsta... | 2 | none |
| A ModulePackage guards every apply by... | A ModulePackage reconcile that renders MUST judge every rendered ob... | A ModulePackage reconcile judges exactly as a ModuleInsta... | 1 | list cut at the two citations; the head and "exactly as" are repeated; added: "(" / ")" / "." |
| A ModulePackage guards every apply by... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | A ModulePackage reconcile judges exactly as a ModuleInsta... | 2 | none |
| A ModulePackage guards every apply by... | A ModulePackage reconcile that renders MUST judge every rendered ob... | The read is made by the client that applies and a skipped... | 1 | sentence 0 is cut at ", with the identity" so that "The read" and "the objects" have their referent; the head is kept; added: "." |
| A ModulePackage guards every apply by... | The read MUST be made by the client that applies. | The read is made by the client that applies and a skipped... | 2 | none |
| A ModulePackage guards every apply by... | A reconcile that skips its render MUST NOT read the objects. | The read is made by the client that applies and a skipped... | 3 | none |
| A ModulePackage guards every apply by... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | The read is made by the client that applies and a skipped... | 4 | none |
| A ModulePackage guards every apply by... | A ModulePackage has no drift check to report a verdict that could n... | A reconcile that renders fails when the verdict could not... | 1 | none |
| A ModulePackage guards every apply by... | So a reconcile that renders MUST fail, also when every digest match... | A reconcile that renders fails when the verdict could not... | 2 | the head of sentence 5; added: "." |
| A ModulePackage guards every apply by... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | A reconcile that renders fails when the verdict could not... | 3 | none |
| A ModulePackage guards every apply by... | So a reconcile that renders MUST fail, as listed in requirement "A ... | The failure reason is ImpersonationFailed when the Servic... | 1 | list cut between the two reasons; the head is kept; the condition of sentence 5 is named by the exact header of the requirement that holds it; added: ", as listed in requirement "A reconcile that renders fails when the verdict could not be asked"" / ": " / "." |
| A ModulePackage guards every apply by... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | The failure reason is ImpersonationFailed when the Servic... | 2 | none |
| A ModulePackage guards every apply by... | So a reconcile that renders MUST fail, as listed in requirement "A ... | The failure reason is ApplyFailed on the backoff except w... | 1 | list cut between the two reasons; the head is kept; the condition of sentence 5 is named by the exact header of the requirement that holds it; added: ", as listed in requirement "A reconcile that renders fails when the verdict could not be asked"" / ": " / ", except when " / "." |
| A ModulePackage guards every apply by... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | The failure reason is ApplyFailed on the backoff except w... | 2 | none |
| A ModulePackage guards every apply by... | So a reconcile that renders MUST fail, also when every digest match... | A failed ModulePackage reconcile writes nothing and keeps... | 1 | the head of sentence 5 is copied as context; added: "." |
| A ModulePackage guards every apply by... | It MUST write nothing and MUST keep `status.inventory`. | A failed ModulePackage reconcile writes nothing and keeps... | 2 | none |
| A ModulePackage guards every apply by... | Source: enhancement 0012 ("The apply guard runs on every apply, ref... | A failed ModulePackage reconcile writes nothing and keeps... | 3 | none |

Lengths (old body, new bodies; characters by the guide's rule):

| Old requirement | Old | New |
| --- | --- | --- |
| No-op detection | 741 | 161, 459, 482 |
| Status always patched | 750 | 245, 410, 200 |
| A ModulePackage records its instance identity | 1309 | 477, 419, 479, 347, 463 |
| A ModulePackage guards every apply by ownership | 1429 | 418, 444, 459, 465, 477, 448, 499, 486, 500, 475 |

### Ledger: digest-computation

Form c2 for every requirement below: REMOVED under the old name, ADDED under new names, new blocks at the end of the spec. `grep -rn --exclude-dir=archive` for each old header found no reference to fix. Old scenario text is copied by script from `origin/main`; item 3 of the guide checklist shows `same` for each.

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| The runtime name does not count | The render digest is RenderDigest from the library's inventory package, The runtime name does not count | no (byte equal, checked by the item 3 command) |
| The stored render digest changes once on upgrade | The render digest is RenderDigest from the library's inventory package, The stored render digest changes once on upgrade | no (byte equal, checked by the item 3 command) |
| Order-independent | The render digest is computed from the exported JSON after the rendered CUE values are dropped, Order-independent | no (byte equal, checked by the item 3 command) |
| Content sensitivity | The render digest is computed from the exported JSON after the rendered CUE values are dropped, Content sensitivity | no (byte equal, checked by the item 3 command) |

New scenarios (restate a clause of the old body, with words of the old text only):

| New scenario | Requirement | Old clause restated |
| --- | --- | --- |
| No render digest of the status package | The internal status package declares no render digest | the clause of that requirement body |
| A digest failure is a render failure | A failure to compute the render digest is reported as a render failure | the clause of that requirement body |

Clause table (old clause, new requirement, sentence of its body; "Edit" lists the words that differ from the old sentence):

| Old requirement | Old clause | New requirement | Sentence | Edit |
| --- | --- | --- | --- | --- |
| The render digest is the library's | The render digest the reconcilers record (`lastAttemptedRenderDiges... | The render digest is RenderDigest from the library's inve... | 1 | none |
| The render digest is the library's | Source: enhancement 0012 ("Labels are stamped by the CUE render; th... | The render digest is RenderDigest from the library's inve... | 2 | none |
| The render digest is the library's | The render digest the reconcilers record (`lastAttemptedRenderDiges... | The render digest is computed from the exported JSON afte... | 1 | none |
| The render digest is the library's | It MUST be computed from the exported JSON after the rendered CUE v... | The render digest is computed from the exported JSON afte... | 2 | none |
| The render digest is the library's | Source: enhancement 0012 ("Labels are stamped by the CUE render; th... | The render digest is computed from the exported JSON afte... | 3 | none |
| The render digest is the library's | The `internal/status` package MUST NOT declare a render digest of i... | The internal status package declares no render digest | 1 | none |
| The render digest is the library's | Source: enhancement 0012 ("Labels are stamped by the CUE render; th... | The internal status package declares no render digest | 2 | none |
| The render digest is the library's | The render digest the reconcilers record (`lastAttemptedRenderDiges... | A failure to compute the render digest is reported as a r... | 1 | none |
| The render digest is the library's | A failure to compute it MUST be reported as a render failure (`Rend... | A failure to compute the render digest is reported as a r... | 2 | none |
| The render digest is the library's | Source: enhancement 0012 ("Labels are stamped by the CUE render; th... | A failure to compute the render digest is reported as a r... | 3 | none |

Lengths (old body, new bodies; characters by the guide's rule):

| Old requirement | Old | New |
| --- | --- | --- |
| The render digest is the library's | 628 | 348, 434, 204, 467 |

### Ledger: serviceaccount-impersonation

Form c2 for every requirement below: REMOVED under the old name, ADDED under new names, new blocks at the end of the spec. `grep -rn --exclude-dir=archive` for each old header found no reference to fix. Old scenario text is copied by script from `origin/main`; item 3 of the guide checklist shows `same` for each.

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| Identity seen by the apiserver | An impersonated request reaches the apiserver as the ServiceAccount user with its groups, Identity seen by the apiserver | no (byte equal, checked by the item 3 command) |
| No group is sent | The controller names no group on an impersonated request, No group is sent | no (byte equal, checked by the item 3 command) |
| Group-subject RoleBinding authorizes apply | An RBAC binding on a ServiceAccount group authorises the impersonated request, Group-subject RoleBinding authorizes apply | no (byte equal, checked by the item 3 command) |
| ServiceAccount without a binding is refused | An RBAC binding on a ServiceAccount group authorises the impersonated request, ServiceAccount without a binding is refused | no (byte equal, checked by the item 3 command) |
| SA deleted before finalizer can prune | A missing deletion-cleanup ServiceAccount stalls with DeletionSAMissing, SA deleted before finalizer can prune | no (byte equal, checked by the item 3 command) |
| Other impersonation errors keep existing behavior | A missing deletion-cleanup ServiceAccount stalls with DeletionSAMissing, Other impersonation errors keep existing behavior | no (byte equal, checked by the item 3 command) |
| ServiceAccount removed while the cleanup waits | Two cases of a missing deletion-cleanup ServiceAccount are not the DeletionSAMissing stall, ServiceAccount removed while the cleanup waits | no (byte equal, checked by the item 3 command) |
| Only kept claims and no ServiceAccount | Both cases remove the finalizer and warn DeletionUnconfirmed, Only kept claims and no ServiceAccount | no (byte equal, checked by the item 3 command) |
| Orphan annotation releases a stuck deletion | The orphan annotation skips prune and removes the finalizer, Orphan annotation releases a stuck deletion | no (byte equal, checked by the item 3 command) |
| Orphan annotation ignored for non-SA-missing stall | The orphan annotation takes effect only for a missing ServiceAccount and only with the value true, Orphan annotation ignored for non-SA-missing stall | no (byte equal, checked by the item 3 command) |
| Annotation value other than "true" is treated as absent | The orphan annotation takes effect only for a missing ServiceAccount and only with the value true, Annotation value other than "true" is treated as absent | no (byte equal, checked by the item 3 command) |
| A created ServiceAccount enqueues the instances that name it | A created ServiceAccount enqueues the instances stalled on it, A created ServiceAccount enqueues the instances that name it | no (byte equal, checked by the item 3 command) |
| The flag default counts as the effective ServiceAccount | A created ServiceAccount enqueues the instances stalled on it, The flag default counts as the effective ServiceAccount | no (byte equal, checked by the item 3 command) |
| A healthy instance is not enqueued | A created ServiceAccount enqueues no instance that is not stalled on it, A healthy instance is not enqueued | no (byte equal, checked by the item 3 command) |
| Instances the ServiceAccount cannot change are not enqueued | A created ServiceAccount enqueues no instance it cannot change, Instances the ServiceAccount cannot change are not enqueued | no (byte equal, checked by the item 3 command) |

New scenarios (restate a clause of the old body, with words of the old text only):

| New scenario | Requirement | Old clause restated |
| --- | --- | --- |
| The stall message names the ServiceAccount and lists the recovery options | The DeletionSAMissing message names the ServiceAccount and the recovery options | the clause of that requirement body |
| The stalled release keeps its finalizer and warns once | The DeletionSAMissing stall keeps the finalizer and warns once | the clause of that requirement body |
| Both deletion paths stall alike | The DeletionSAMissing stall applies to both deletion paths | the clause of that requirement body |
| The orphan exit warns and clears the inventory | The orphan annotation emits OrphanedOnDeletion and clears the inventory | the clause of that requirement body |
| An instance being deleted and one stalled on the apply path are both enqueued | A created ServiceAccount enqueues an instance being deleted and one stalled on the apply path | the clause of that requirement body |
| An update or a deletion of a ServiceAccount enqueues nothing | A created ServiceAccount enqueues no instance it cannot change | the clause of that requirement body |
| A restored Role or binding enqueues nothing | The restore of ServiceAccount RBAC alone is not a trigger | the clause of that requirement body |

Clause table (old clause, new requirement, sentence of its body; "Edit" lists the words that differ from the old sentence):

| Old requirement | Old clause | New requirement | Sentence | Edit |
| --- | --- | --- | --- | --- |
| Impersonated identity has the Service... | An impersonated request MUST reach the apiserver as the user `syste... | An impersonated request reaches the apiserver as the Serv... | 1 | none |
| Impersonated identity has the Service... | An impersonated request MUST reach the apiserver as the user `syste... | The controller names no group on an impersonated request | 1 | none |
| Impersonated identity has the Service... | The controller MUST NOT name a group on an impersonated request; th... | The controller names no group on an impersonated request | 2 | none |
| Impersonated identity has the Service... | An impersonated request MUST reach the apiserver as the user `syste... | An RBAC binding on a ServiceAccount group authorises the ... | 1 | none |
| Impersonated identity has the Service... | An RBAC binding whose subject is one of these groups MUST authorise... | An RBAC binding on a ServiceAccount group authorises the ... | 2 | none |
| Deletion-cleanup SA-missing stalls wi... | When the deletion cleanup path attempts to build an impersonated cl... | A missing deletion-cleanup ServiceAccount stalls with Del... | 1 | none |
| Deletion-cleanup SA-missing stalls wi... | - NOT silently fall back to the controller's own client. | A missing deletion-cleanup ServiceAccount stalls with Del... | 2 | none |
| Deletion-cleanup SA-missing stalls wi... | - Stall the release with `Ready=False` and condition reason `Deleti... | A missing deletion-cleanup ServiceAccount stalls with Del... | 3 | none |
| Deletion-cleanup SA-missing stalls wi... | - Requeue on the stalled-recheck interval, not the tight transient ... | A missing deletion-cleanup ServiceAccount stalls with Del... | 4 | none |
| Deletion-cleanup SA-missing stalls wi... | When the deletion cleanup path attempts to build an impersonated cl... | The DeletionSAMissing message names the ServiceAccount an... | 1 | none |
| Deletion-cleanup SA-missing stalls wi... | - Populate the condition message with the SA's namespaced name and ... | The DeletionSAMissing message names the ServiceAccount an... | 2 | none |
| Deletion-cleanup SA-missing stalls wi... | When the deletion cleanup path attempts to build an impersonated cl... | The DeletionSAMissing stall keeps the finalizer and warns... | 1 | none |
| Deletion-cleanup SA-missing stalls wi... | - Retain the finalizer. | The DeletionSAMissing stall keeps the finalizer and warns... | 2 | none |
| Deletion-cleanup SA-missing stalls wi... | - Emit a `Warning` event with reason `DeletionSAMissing` on the Rea... | The DeletionSAMissing stall keeps the finalizer and warns... | 3 | none |
| Deletion-cleanup SA-missing stalls wi... | This behavior MUST apply symmetrically to both `ModuleRelease` and ... | The DeletionSAMissing stall applies to both deletion paths | 1 | added the named reference: As in requirements "A missing deletion-cleanup ServiceAccount stalls with DeletionSAMissing", "The DeletionSAMissing message names the ServiceAccount and the recovery options" and "The DeletionSAMissing stall keeps the finalizer and warns once": |
| Deletion-cleanup SA-missing stalls wi... | Two cases | Two cases of a missing deletion-cleanup ServiceAccount ar... | 1 | added the named reference: As in requirement "A missing deletion-cleanup ServiceAccount stalls with DeletionSAMissing":; added "SHALL NOT be" (Strength table) |
| Deletion-cleanup SA-missing stalls wi... |  this stall, because the missing ServiceAccount has nothing left to... | Two cases of a missing deletion-cleanup ServiceAccount ar... | 2 | none |
| Deletion-cleanup SA-missing stalls wi... | In both the finalizer | Both cases remove the finalizer and warn DeletionUnconfirmed | 1 | added the named reference: As in requirement "Two cases of a missing deletion-cleanup ServiceAccount are not the DeletionSAMissing stall":; added "SHALL be" (Strength table) |
| Deletion-cleanup SA-missing stalls wi... |  removed, the controller's own client | Both cases remove the finalizer and warn DeletionUnconfirmed | 2 | added "SHALL" (Strength table) |
| Deletion-cleanup SA-missing stalls wi... |  still never | Both cases remove the finalizer and warn DeletionUnconfirmed | 3 | none |
| Deletion-cleanup SA-missing stalls wi... |  used, and a `Warning` event with reason `DeletionUnconfirmed` | Both cases remove the finalizer and warn DeletionUnconfirmed | 4 | added "be"; added "SHALL say" (Strength table) |
| Deletion-cleanup SA-missing stalls wi... |  what could not be read. | Both cases remove the finalizer and warn DeletionUnconfirmed | 5 | none |
| Orphan-exit annotation removes finali... | When a release is in the `DeletionSAMissing` stall state AND the an... | The orphan annotation skips prune and removes the finalizer | 1 | none |
| Orphan-exit annotation removes finali... | - Skip prune entirely. | The orphan annotation skips prune and removes the finalizer | 2 | none |
| Orphan-exit annotation removes finali... | - Remove the finalizer. | The orphan annotation skips prune and removes the finalizer | 3 | none |
| Orphan-exit annotation removes finali... | When a release is in the `DeletionSAMissing` stall state AND the an... | The orphan annotation emits OrphanedOnDeletion and clears... | 1 | none |
| Orphan-exit annotation removes finali... | - Emit a `Warning` event with reason `OrphanedOnDeletion` naming th... | The orphan annotation emits OrphanedOnDeletion and clears... | 2 | none |
| Orphan-exit annotation removes finali... | The annotation MUST only take effect when the deletion cleanup impe... | The orphan annotation takes effect only for a missing Ser... | 1 | added the named reference: As in requirement "The orphan annotation skips prune and removes the finalizer": |
| Orphan-exit annotation removes finali... | Any annotation value other than the literal string `"true"` MUST be... | The orphan annotation takes effect only for a missing Ser... | 2 | none |
| A created ServiceAccount wakes the in... | When a ServiceAccount is created, the controller MUST enqueue every... | A created ServiceAccount enqueues the instances stalled o... | 1 | none |
| A created ServiceAccount wakes the in... | Such an instance then recovers without waiting for its stalled rech... | A created ServiceAccount enqueues the instances stalled o... | 2 | none |
| A created ServiceAccount wakes the in... | When a ServiceAccount is created, the controller MUST enqueue every... | A created ServiceAccount enqueues an instance being delet... | 1 | none |
| A created ServiceAccount wakes the in... | This holds for an instance being deleted and for one stalled on the... | A created ServiceAccount enqueues an instance being delet... | 2 | none |
| A created ServiceAccount wakes the in... | The controller MUST NOT enqueue an instance that is not stalled on ... | A created ServiceAccount enqueues no instance that is not... | 1 | none |
| A created ServiceAccount wakes the in... | The controller MUST NOT enqueue an instance the ServiceAccount cann... | A created ServiceAccount enqueues no instance it cannot c... | 1 | added the named reference: As in requirement "A created ServiceAccount enqueues the instances stalled on it": |
| A created ServiceAccount wakes the in... | The restore of a ServiceAccount's RBAC alone (a Role or a binding) | The restore of ServiceAccount RBAC alone is not a trigger | 1 | added the named reference: As in requirement "A created ServiceAccount enqueues the instances stalled on it":; added "SHALL NOT be" (Strength table) |
| A created ServiceAccount wakes the in... |  a trigger; an instance stalled with `ImpersonationFailed` on a for... | The restore of ServiceAccount RBAC alone is not a trigger | 2 | none |

Lengths (old body, new bodies; characters by the guide's rule):

| Old requirement | Old | New |
| --- | --- | --- |
| Impersonated identity has the ServiceAccount's groups | 547 | 258, 393, 412 |
| Deletion-cleanup SA-missing stalls with distinct reason | 1386 | 472, 355, 314, 338, 451, 296 |
| Orphan-exit annotation removes finalizer on SA-missing | 794 | 208, 419, 408 |
| A created ServiceAccount wakes the instances that impersonate it | 1263 | 430, 438, 256, 386, 280 |

### Ledger: platform-crd

Form c2 for every requirement below: REMOVED under the old name, ADDED under new names, new blocks at the end of the spec. `grep -rn --exclude-dir=archive` for each old header found no reference to fix. Old scenario text is copied by script from `origin/main`; item 3 of the guide checklist shows `same` for each.

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| Skew policy defaults to Warn | PlatformSpec carries type, registry and skewPolicy and generates the platform module from type and registry, Skew policy defaults to Warn | no (byte equal, checked by the item 3 command) |
| Minimal valid platform spec | PlatformSpec carries type, registry and skewPolicy and generates the platform module from type and registry, Minimal valid platform spec | no (byte equal, checked by the item 3 command) |
| type is required | PlatformSpec carries type, registry and skewPolicy and generates the platform module from type and registry, type is required | no (byte equal, checked by the item 3 command) |
| Invalid skew policy is rejected at admission | skewPolicy answers a catalog requirement newer than the pinned build and is validated at admission, Invalid skew policy is rejected at admission | no (byte equal, checked by the item 3 command) |
| Scalar subscription accepted | The registry maps catalog paths to subscriptions of one published version, Scalar subscription accepted | no (byte equal, checked by the item 3 command) |
| Filter rejected | The registry maps catalog paths to subscriptions of one published version, Filter rejected | no (byte equal, checked by the item 3 command) |
| Stored legacy singleton remains serviceable | A subscription without a version never reaches a platform build, Stored legacy singleton remains serviceable | no (byte equal, checked by the item 3 command) |
| Status subresource present | PlatformStatus carries conditions, observedGeneration and operatorVersion, Status subresource present | no (byte equal, checked by the item 3 command) |
| Operator printcolumn | PlatformStatus carries conditions, observedGeneration and operatorVersion, Operator printcolumn | no (byte equal, checked by the item 3 command) |
| Status reflects the generate-and-build outcome | The Ready and ContractsFulfilled conditions carry their reasons, Status reflects the generate-and-build outcome | no (byte equal, checked by the item 3 command) |
| Status carries the contract report | The Ready and ContractsFulfilled conditions carry their reasons, Status carries the contract report | no (byte equal, checked by the item 3 command) |

New scenarios (restate a clause of the old body, with words of the old text only):

| New scenario | Requirement | Old clause restated |
| --- | --- | --- |
| The skewPolicy is not part of the generated module | PlatformSpec carries type, registry and skewPolicy and generates the platform module from type and registry | the clause of that requirement body |
| The field documentation names every reason | The documentation of the conditions field names every reason | the clause of that requirement body |

Clause table (old clause, new requirement, sentence of its body; "Edit" lists the words that differ from the old sentence):

| Old requirement | Old clause | New requirement | Sentence | Edit |
| --- | --- | --- | --- | --- |
| PlatformSpec projects core #Platform | `PlatformSpec` SHALL carry `type` (required, informational discrimi... | PlatformSpec carries type, registry and skewPolicy and ge... | 1 | none |
| PlatformSpec projects core #Platform | The operator generates a platform CUE module from `type` and `regis... | PlatformSpec carries type, registry and skewPolicy and ge... | 2 | none |
| PlatformSpec projects core #Platform | `skewPolicy` is the operator's response when a module's declared ca... | skewPolicy answers a catalog requirement newer than the p... | 1 | none |
| Subscription shape | `PlatformSpec.Registry` SHALL map major-suffixed catalog module pat... | The registry maps catalog paths to subscriptions of one p... | 1 | none |
| Subscription shape | A subscription without a version SHALL never reach a platform build. | A subscription without a version never reaches a platform... | 1 | none |
| Subscription shape | Whether `version` is schema-`required` or schema-optional follows t... | A subscription without a version never reaches a platform... | 2 | the em-dash is written as ":" |
| Subscription shape |  admission in the required posture, the reconciler before any regis... | A subscription without a version never reaches a platform... | 3 | none |
| PlatformStatus carries conditions and... | `PlatformStatus` SHALL carry `conditions` (a `metav1.Condition` lis... | PlatformStatus carries conditions, observedGeneration and... | 1 | none |
| PlatformStatus carries conditions and... | The Ready condition SHALL summarise module generation: `Ready=True`... | The Ready and ContractsFulfilled conditions carry their r... | 1 | none |
| PlatformStatus carries conditions and... | A `ContractsFulfilled` condition SHALL report the effective package... | The Ready and ContractsFulfilled conditions carry their r... | 2 | none |
| PlatformStatus carries conditions and... | The field's documentation SHALL name every reason. | The documentation of the conditions field names every reason | 1 | added the named reference: As in requirement "PlatformStatus carries conditions, observedGeneration and operatorVersion": |

Lengths (old body, new bodies; characters by the guide's rule):

| Old requirement | Old | New |
| --- | --- | --- |
| PlatformSpec projects core #Platform | 564 | 305, 258 |
| Subscription shape | 726 | 351, 373 |
| PlatformStatus carries conditions and observedGeneration | 609 | 128, 429, 145 |

### Ledger: platform-module-generation

Form c2 for every requirement below: REMOVED under the old name, ADDED under new names, new blocks at the end of the spec. `grep -rn --exclude-dir=archive` for each old header found no reference to fix. Old scenario text is copied by script from `origin/main`; item 3 of the guide checklist shows `same` for each.

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| Core pin follows the library | The operator carries no core-pin constant of its own, Core pin follows the library | no (byte equal, checked by the item 3 command) |
| A two-catalog CR generates a two-entry module | The registry has one entry per authored subscription, A two-catalog CR generates a two-entry module | no (byte equal, checked by the item 3 command) |
| An active claim's catalog is pinned and imported | The registry has one entry per accepted and active claim, An active claim's catalog is pinned and imported | no (byte equal, checked by the item 3 command) |
| An authored subscription wins over a claim naming its catalog | An authored subscription wins over a claim naming its catalog, An authored subscription wins over a claim naming its catalog | no (byte equal, checked by the item 3 command) |
| A transitive dependency is pinned in the closure | The dependency list is the full closure of the pinned modules, A transitive dependency is pinned in the closure | no (byte equal, checked by the item 3 command) |
| A disabled subscription is generated disabled | A disabled subscription is generated disabled and still pinned and imported, A disabled subscription is generated disabled | no (byte equal, checked by the item 3 command) |
| Deterministic regeneration | Regeneration is deterministic, Deterministic regeneration | no (byte equal, checked by the item 3 command) |
| A claimless platform generates what its spec alone implies | A Platform with no active claim generates the module its spec alone produces, A claimless platform generates what its spec alone implies | no (byte equal, checked by the item 3 command) |
| Clean build sets Ready | The reconciler builds the generated module, reads its inventory and sets Ready=True when the gate passes, Clean build sets Ready | no (byte equal, checked by the item 3 command) |
| Nonexistent pin surfaces on the CR | A build that fails sets Ready=False with reason BuildFailed, Nonexistent pin surfaces on the CR | no (byte equal, checked by the item 3 command) |
| A failed build keeps the last good module | A failed or refused reconcile leaves the previously recorded module in place, A failed build keeps the last good module | no (byte equal, checked by the item 3 command) |
| A built but refused module keeps the last good module | A failed or refused reconcile leaves the previously recorded module in place, A built but refused module keeps the last good module | no (byte equal, checked by the item 3 command) |
| Two identities for one generation get separate directories | Each package's module is written to a directory named by the package identity, Two identities for one generation get separate directories | no (byte equal, checked by the item 3 command) |
| A crash mid-write leaves no current module | Generation writes into a staging directory and renames it into place, A crash mid-write leaves no current module | no (byte equal, checked by the item 3 command) |
| Superseded packages are pruned | After a successful build the platform directory keeps only the current and the leased packages, Superseded packages are pruned | no (byte equal, checked by the item 3 command) |
| A leased package survives the prune | After a successful build the platform directory keeps only the current and the leased packages, A leased package survives the prune | no (byte equal, checked by the item 3 command) |
| Restart regenerates from the CR | The platform directory is emptied at manager start, Restart regenerates from the CR | no (byte equal, checked by the item 3 command) |

New scenarios (restate a clause of the old body, with words of the old text only):

| New scenario | Requirement | Old clause restated |
| --- | --- | --- |
| The module carries a cue.mod at the reserved path | The platform module has a cue.mod/module.cue at the reserved path | the clause of that requirement body |
| The cue.mod pins core and every catalog the registry names | The generated cue.mod pins core at the library's release and every catalog the registry names | the clause of that requirement body |
| The platform.cue carries one importing entry per catalog | The generated platform.cue embeds core.#Platform with one registry entry per catalog | the clause of that requirement body |
| The module-file source comes from the operator's configuration | The module-file source is constructed from the operator's configuration | the clause of that requirement body |
| A module that could not be written is GenerateFailed | A module that could not be written sets Ready=False with reason GenerateFailed | the clause of that requirement body |
| A refused module sets the refusal's reason | A module that built and was refused by the gate sets Ready=False with the refusal's reason | the clause of that requirement body |
| The materialize-era reasons are not used | The materialize-era reasons are retired | the clause of that requirement body |

Clause table (old clause, new requirement, sentence of its body; "Edit" lists the words that differ from the old sentence):

| Old requirement | Old clause | New requirement | Sentence | Edit |
| --- | --- | --- | --- | --- |
| The reconciler generates a platform m... | On reconciling the singleton Platform CR, the operator SHALL genera... | The platform module has a cue.mod/module.cue at the reserved path | 1 | none |
| The reconciler generates a platform m... | a `cue.mod/module.cue` under the fixed reserved-unpublished module ... | The platform module has a cue.mod/module.cue at the reserved path | 2 | added "." |
| The reconciler generates a platform m... | On reconciling the singleton Platform CR, the operator SHALL genera... (head; a `cue.mod/module.cue` under the fixed reserved-unpublished module path ...) | The generated cue.mod pins core at the library's release ... | 1 | the head sentence replaced by the named reference As in requirement "The platform module has a cue.mod/module.cue at the reserved path": and the words "the operator SHALL generate" and the list item's opening "a `cue.mod/module.cue` under the fixed reserved-unpublished module path `opmodel.dev/platforms/cluster@v0`," copied from the old head; no cut inside the list item |
| The reconciler generates a platform m... | pinning core at the library's verified core release (the version ca... | The generated cue.mod pins core at the library's release ... | 1 | none |
| The reconciler generates a platform m... | On reconciling the singleton Platform CR, the operator SHALL genera... | The generated platform.cue embeds core.#Platform with one... | 1 | none |
| The reconciler generates a platform m... | a `platform.cue` embedding `core.#Platform` with one `#registry` en... | The generated platform.cue embeds core.#Platform with one... | 2 | none |
| The reconciler generates a platform m... | The operator SHALL carry no core-pin constant of its own. | The operator carries no core-pin constant of its own | 1 | none |
| The reconciler generates a platform m... | The registry SHALL be resolved from BOTH of the two paths transform... | The registry has one entry per authored subscription | 1 | none |
| The reconciler generates a platform m... | one entry per subscription the CR's `spec.registry` authored, with ... | The registry has one entry per authored subscription | 2 | added "." |
| The reconciler generates a platform m... | The registry SHALL be resolved from BOTH of the two paths transform... | The registry has one entry per accepted and active claim | 1 | none |
| The reconciler generates a platform m... | one entry per accepted-and-active `TransformerRegistration` claim, ... | The registry has one entry per accepted and active claim | 2 | none |
| The reconciler generates a platform m... | A catalog is one `#registry` key, so a catalog named by an authored... | An authored subscription wins over a claim naming its cat... | 1 | none |
| The reconciler generates a platform m... | The dependency list SHALL be the full closure: beyond those roots i... | The dependency list is the full closure of the pinned mod... | 1 | added the named reference: As in requirements "The registry has one entry per authored subscription" and "The registry has one entry per accepted and active claim":; added "." |
| The reconciler generates a platform m... | the module-file source SHALL be constructed from the operator's con... | The module-file source is constructed from the operator's... | 1 | added the named reference: As in requirement "The dependency list is the full closure of the pinned modules": |
| The reconciler generates a platform m... | A disabled subscription (`enable: false`) SHALL be generated with `... | A disabled subscription is generated disabled and still p... | 1 | none |
| The reconciler generates a platform m... | Regeneration SHALL be deterministic: the same inputs produce byte-i... | Regeneration is deterministic | 1 | none |
| The reconciler generates a platform m... | A Platform with no active claim SHALL generate the module its spec ... | A Platform with no active claim generates the module its ... | 1 | none |
| The generated module is validated by ... | After writing the module, the reconciler SHALL build it through the... | The reconciler builds the generated module, reads its inv... | 1 | none |
| The generated module is validated by ... | The Ready condition SHALL reflect the outcome: | The reconciler builds the generated module, reads its inv... | 2 | none |
| The generated module is validated by ... | Ready=True with reason `Generated` when the build succeeds and the ... | The reconciler builds the generated module, reads its inv... | 3 | added "." |
| The generated module is validated by ... | The Ready condition SHALL reflect the outcome: | A build that fails sets Ready=False with reason BuildFailed | 1 | none |
| The generated module is validated by ... | Ready=False with reason `BuildFailed`, with the error naming the fa... | A build that fails sets Ready=False with reason BuildFailed | 2 | added "." |
| The generated module is validated by ... | The Ready condition SHALL reflect the outcome: | A module that could not be written sets Ready=False with ... | 1 | none |
| The generated module is validated by ... | Ready=False with reason `GenerateFailed` when the module could not ... | A module that could not be written sets Ready=False with ... | 2 | added "." |
| The generated module is validated by ... | The Ready condition SHALL reflect the outcome: | A module that built and was refused by the gate sets Read... | 1 | none |
| The generated module is validated by ... | Ready=False with reason `ContractCollisions`, `OverSubscribedContra... | A module that built and was refused by the gate sets Read... | 2 | none |
| The generated module is validated by ... | The materialize-era reasons (`Materialized`, `MaterializeFailed`) | The materialize-era reasons are retired | 1 | added "SHALL be" (Strength table) |
| The generated module is validated by ... |  retired. | The materialize-era reasons are retired | 2 | none |
| The generated module is validated by ... | A failed or refused reconcile SHALL leave the previously recorded m... | A failed or refused reconcile leaves the previously recor... | 1 | none |
| The generated module lives in a per-i... | The operator SHALL write each package's module under a manager-conf... | Each package's module is written to a directory named by ... | 1 | none |
| The generated module lives in a per-i... | Generation SHALL write into a staging directory and rename it into ... | Generation writes into a staging directory and renames it... | 1 | none |
| The generated module lives in a per-i... | After a successful build the operator SHALL keep the current packag... | After a successful build the platform directory keeps onl... | 1 | none |
| The generated module lives in a per-i... | At manager start the operator SHALL empty the platform directory; t... | The platform directory is emptied at manager start | 1 | none |

Lengths (old body, new bodies; characters by the guide's rule):

| Old requirement | Old | New |
| --- | --- | --- |
| The reconciler generates a platform module from the CR | 2750 | 300, 500, 308, 57, 480, 472, 314, 432, 239, 160, 379, 82 |
| The generated module is validated by building it | 1189 | 474, 358, 133, 191, 83, 91 |
| The generated module lives in a per-identity directory and is swapped whole | 1243 | 459, 188, 355, 238 |

### Ledger: reconcile-backoff

Form b for both requirements below: MODIFIED keeps the old header, every old scenario and a shorter body; the splits are ADDED under new names. The headers cannot be renamed, because `grep -rn --exclude-dir=archive` finds "Registry fetch failures are transient wherever they occur" in `platform-reconciler`, `module-instance-synthesis`, `modulepackage-artifact-loading` and `modulepackage-reconcile-loop`, and "A render timeout retries on the backoff" in `library-kernel-runtime`. A live scenario never moves from a MODIFIED parent to an ADDED block, so the parent keeps all of them and each ADDED block carries a restatement scenario of its own. Old scenario text is copied by script from `origin/main`; item 3 of the guide checklist shows `same` for each. The kept header "Registry fetch failures are transient wherever they occur" says where the rule holds; its body states the rule by the error chain, which does not depend on the call, and the ADDED requirement "A render failure is classified by its type, never by its message text" keeps the list of kernel calls.

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| Registry outage during ModuleInstance acquisition retries on the backoff | Registry fetch failures are transient wherever they occur, Registry outage during ModuleInstance acquisition retries on the backoff | no (byte equal, checked by the item 3 command) |
| Registry outage while loading a ModulePackage retries on the backoff | Registry fetch failures are transient wherever they occur, Registry outage while loading a ModulePackage retries on the backoff | no (byte equal, checked by the item 3 command) |
| Registry failure during synthesis or render retries on the backoff | Registry fetch failures are transient wherever they occur, Registry failure during synthesis or render retries on the backoff | no (byte equal, checked by the item 3 command) |
| Not-found and refused credentials retry on the backoff | Registry fetch failures are transient wherever they occur, Not-found and refused credentials retry on the backoff | no (byte equal, checked by the item 3 command) |
| Repeated fetch failures cap at five minutes | Registry fetch failures are transient wherever they occur, Repeated fetch failures cap at five minutes | no (byte equal, checked by the item 3 command) |
| Recovery once the registry answers | Registry fetch failures are transient wherever they occur, Recovery once the registry answers | no (byte equal, checked by the item 3 command) |
| Typed terminal causes stay stalled | Registry fetch failures are transient wherever they occur, Typed terminal causes stay stalled | no (byte equal, checked by the item 3 command) |
| Author defect in a package stalls | Registry fetch failures are transient wherever they occur, Author defect in a package stalls | no (byte equal, checked by the item 3 command) |
| Unclassified acquisition failure stalls | Registry fetch failures are transient wherever they occur, Unclassified acquisition failure stalls | no (byte equal, checked by the item 3 command) |
| Message text does not classify | Registry fetch failures are transient wherever they occur, Message text does not classify | no (byte equal, checked by the item 3 command) |
| A ModuleInstance render times out | A render timeout retries on the backoff, A ModuleInstance render times out | no (byte equal, checked by the item 3 command) |
| A ModulePackage render times out | A render timeout retries on the backoff, A ModulePackage render times out | no (byte equal, checked by the item 3 command) |
| Repeated timeouts cap at five minutes | A render timeout retries on the backoff, Repeated timeouts cap at five minutes | no (byte equal, checked by the item 3 command) |
| A render that fails at its deadline is still a timeout | A render timeout retries on the backoff, A render that fails at its deadline is still a timeout | no (byte equal, checked by the item 3 command) |
| Recovery after a timeout | A render timeout retries on the backoff, Recovery after a timeout | no (byte equal, checked by the item 3 command) |

New scenarios (restate a clause of the old body, with words of the old text only):

| New scenario | Requirement | Old clause restated |
| --- | --- | --- |
| A render failure is classified by its type | A render failure is classified by its type, never by its message text | the clause of that requirement body |
| A transient failure retries on the backoff | A transient failure sets ResolutionFailed and produces FailedTransient | the clause of that requirement body |
| A failure that is not transient stalls | Every other failure stalls and requeues on the 30-minute recheck | the clause of that requirement body |
| A render that does not finish is a timeout | A render timeout is classified before any other render failure | the clause of that requirement body |
| The timeout message names the timeout | The timeout message names the timeout and the render slot | the clause of that requirement body |
| A timeout applies and prunes nothing | A timeout keeps inventory and digests and applies and prunes nothing | the clause of that requirement body |
| A deadline error of another source keeps its classification | A deadline error that is not the render deadline keeps the registry fetch classification | the clause of that requirement body |

Clause table (old clause, new requirement, sentence of its body; "Edit" lists the words that differ from the old sentence):

| Old requirement | Old clause | New requirement | Sentence | Edit |
| --- | --- | --- | --- | --- |
| Registry fetch failures are transient... | A failure | Registry fetch failures are transient wherever they occur | 1 | added the named reference: As in requirement "A render failure is classified by its type, never by its message text":; added "SHALL be" (Strength table) |
| Registry fetch failures are transient... |  transient when its error chain holds the library's typed registry ... | Registry fetch failures are transient wherever they occur | 2 | none |
| Registry fetch failures are transient... | When a render fails, the ModuleInstance and ModulePackage reconcile... | A render failure is classified by its type, never by its ... | 1 | none |
| Registry fetch failures are transient... | This holds for every kernel call the render makes: module acquisiti... | A render failure is classified by its type, never by its ... | 2 | none |
| Registry fetch failures are transient... | A transient failure SHALL set `Ready=False` with reason `Resolution... | A transient failure sets ResolutionFailed and produces Fa... | 1 | added the named reference: As in requirement "Registry fetch failures are transient wherever they occur": |
| Registry fetch failures are transient... | Every other failure SHALL set `Ready=False` and `Stalled=True` and ... | Every other failure stalls and requeues on the 30-minute ... | 1 | added the named reference: As in requirement "Registry fetch failures are transient wherever they occur": |
| A render timeout retries on the backoff | A timeout SHALL set `Ready=False` with reason `RenderTimedOut` and ... | A render timeout retries on the backoff | 1 | added the named reference: As in requirement "A render timeout is classified before any other render failure": |
| A render timeout retries on the backoff | When a ModuleInstance or ModulePackage render does not finish withi... | A render timeout is classified before any other render fa... | 1 | none |
| A render timeout retries on the backoff | The message SHALL name the timeout and say that the render's slot s... | The timeout message names the timeout and the render slot | 1 | added the named reference: As in requirement "A render timeout retries on the backoff": |
| A render timeout retries on the backoff | Inventory and the last-applied digests SHALL keep the last success,... | A timeout keeps inventory and digests and applies and pru... | 1 | added the named reference: As in requirement "A render timeout retries on the backoff": |
| A render timeout retries on the backoff | A render error that carries `context.DeadlineExceeded` from anythin... | A deadline error that is not the render deadline keeps th... | 1 | none |

Lengths (old body, new bodies; characters by the guide's rule):

| Old requirement | Old | New |
| --- | --- | --- |
| Registry fetch failures are transient wherever they occur | 1412 | 475, 338, 442, 409 |
| A render timeout retries on the backoff | 1270 | 494, 310, 227, 168, 273 |

### Ledger: registration-driven-regeneration

Form c2 for every requirement below: REMOVED under the old name, ADDED under new names, new blocks at the end of the spec. `grep -rn --exclude-dir=archive` for each old header found no reference to fix. Old scenario text is copied by script from `origin/main`; item 3 of the guide checklist shows `same` for each.

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| An active claim's catalog reaches the generated package | The platform package is generated from the Platform spec and the accepted and active claims, An active claim's catalog reaches the generated package | no (byte equal, checked by the item 3 command) |
| No active claims leaves generation unchanged | The platform package is generated from the Platform spec and the accepted and active claims, No active claims leaves generation unchanged | no (byte equal, checked by the item 3 command) |
| Generation ignores the waking event's content | Generation computes from the current state, never from the waking event, Generation ignores the waking event's content | no (byte equal, checked by the item 3 command) |
| A blocked claim's catalog survives regeneration | A claim whose deletion is blocked remains in the active set, A blocked claim's catalog survives regeneration | no (byte equal, checked by the item 3 command) |
| A released claim leaves the active set | A claim whose deletion is blocked remains in the active set, A released claim leaves the active set | no (byte equal, checked by the item 3 command) |
| The union shows both sources | The Platform status carries the union of the authored subscriptions and the active claims' catalogs with the source of each entry, The union shows both sources | no (byte equal, checked by the item 3 command) |
| The union follows the active set | The Platform status carries the union of the authored subscriptions and the active claims' catalogs with the source of each entry, The union follows the active set | no (byte equal, checked by the item 3 command) |
| A disabled subscription is listed as not registering | The Platform status shows the pinned version and whether the transformers of an entry register, A disabled subscription is listed as not registering | no (byte equal, checked by the item 3 command) |

Clause table (old clause, new requirement, sentence of its body; "Edit" lists the words that differ from the old sentence):

| Old requirement | Old clause | New requirement | Sentence | Edit |
| --- | --- | --- | --- | --- |
| The generated package is a function o... | The operator SHALL generate the platform package from exactly two i... | The platform package is generated from the Platform spec ... | 1 | none |
| The generated package is a function o... | A platform with no active claims SHALL generate the package it woul... | The platform package is generated from the Platform spec ... | 2 | none |
| The generated package is a function o... | The operator SHALL generate the platform package from exactly two i... | Generation computes from the current state, never from th... | 1 | none |
| The generated package is a function o... | A claim whose deletion is blocked SHALL remain in the active set un... | A claim whose deletion is blocked remains in the active set | 1 | none |
| The effective registry is readable fr... | `Platform` status SHALL carry the resolved union of the registry th... | The Platform status carries the union of the authored sub... | 1 | none |
| The effective registry is readable fr... | `Platform` status SHALL carry the resolved union of the registry th... | The Platform status shows the pinned version and whether ... | 1 | none |
| The effective registry is readable fr... | An entry SHALL also carry the version the package pins it at and wh... | The Platform status shows the pinned version and whether ... | 2 | the two em-dashes are written as "," |
| The effective registry is readable fr... |  it is still pinned and imported | The Platform status shows the pinned version and whether ... | 3 | the two em-dashes are written as "," |
| The effective registry is readable fr... |  but contributes no transformer, and a union read without that dist... | The Platform status shows the pinned version and whether ... | 4 | none |

Lengths (old body, new bodies; characters by the guide's rule):

| Old requirement | Old | New |
| --- | --- | --- |
| The generated package is a function of the platform spec and the active-claim set | 865 | 261, 324, 435 |
| The effective registry is readable from the Platform | 607 | 309, 469 |

### Ledger: registration-removal-guard

Form b for "Shrinking provides while dependents exist is refused before it takes effect" (MODIFIED keeps the header, the first paragraph as the body and all five scenarios; the two splits are ADDED under new names with a restatement scenario each) and form c2 for "A refused upgrade is reported on the provider it came from" (REMOVED, two ADDED). The ADDED blocks name the kept header as the requirement whose subject they continue. `grep -rn --exclude-dir=archive` for each old header found no reference to fix. Old scenario text is copied by script from `origin/main`; item 3 of the guide checklist shows `same` for each.

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| A shrinking upgrade is refused | Shrinking provides while dependents exist is refused before it takes effect, A shrinking upgrade is refused | no (byte equal, checked by the item 3 command) |
| The previously accepted claim keeps serving | Shrinking provides while dependents exist is refused before it takes effect, The previously accepted claim keeps serving | no (byte equal, checked by the item 3 command) |
| A non-shrinking upgrade passes untouched | Shrinking provides while dependents exist is refused before it takes effect, A non-shrinking upgrade passes untouched | no (byte equal, checked by the item 3 command) |
| Dropping a contract nobody demands is allowed | Shrinking provides while dependents exist is refused before it takes effect, Dropping a contract nobody demands is allowed | no (byte equal, checked by the item 3 command) |
| A provider that renders no claim is unaffected | Shrinking provides while dependents exist is refused before it takes effect, A provider that renders no claim is unaffected | no (byte equal, checked by the item 3 command) |
| The instance names what was refused | A refused shrink is reported on the ModuleInstance whose render produced the claim, The instance names what was refused | no (byte equal, checked by the item 3 command) |
| The refusal clears when the cause does | A refused shrink is reported on the ModuleInstance whose render produced the claim, The refusal clears when the cause does | no (byte equal, checked by the item 3 command) |
| The claim's own verdict is untouched | A refused shrink does not write the claim's own status, The claim's own verdict is untouched | no (byte equal, checked by the item 3 command) |

New scenarios (restate a clause of the old body, with words of the old text only):

| New scenario | Requirement | Old clause restated |
| --- | --- | --- |
| A refusal names the dropped contracts and the dependent count | The refusal of a shrink names the dropped contracts and the dependent count | the clause of that requirement body |
| A claim written by another actor is outside the guarantee | The guarantee against shrinking is bounded to writes the operator makes | the clause of that requirement body |

Clause table (old clause, new requirement, sentence of its body; "Edit" lists the words that differ from the old sentence):

| Old requirement | Old clause | New requirement | Sentence | Edit |
| --- | --- | --- | --- | --- |
| Shrinking provides while de... | The operator SHALL refuse to apply a rendered `TransformerRegistrat... | Shrinking provides while dependents exist is refused befo... | 1 | none |
| Shrinking provides while de... | The refusal SHALL name the dropped contracts and the dependent coun... | The refusal of a shrink names the dropped contracts and t... | 1 | added the named reference: As in requirement "Shrinking provides while dependents e...": |
| Shrinking provides while de... | This guarantee is bounded to writes the operator makes. A `Transfor... | The guarantee against shrinking is bounded to writes the ... | 1 | added the named reference: As in requirement "Shrinking provides while dependents e...":; the em-dashes are written as "," |
| A refused upgrade is report... | The operator SHALL report a refused shrink on the `ModuleInstance` ... | A refused shrink is reported on the ModuleInstance whose ... | 1 | none |
| A refused upgrade is report... | The claim's own status SHALL NOT be written by this refusal. Accept... | A refused shrink does not write the claim's own status | 1 | added the named reference: As in requirement "A refused shrink is reported on the M...": |

Lengths (old body, new bodies; characters by the guide's rule):

| Old requirement | Old | New |
| --- | --- | --- |
| Shrinking provides while dependents exist is refused before it takes effect | 887 | 473, 256, 348 |
| A refused upgrade is reported on the provider it came from | 530 | 353, 280 |

### Ledger: test-registry-lifecycle

The block "The render result carries no inventory entries and full apply is covered by the Kind-backed e2e suite" holds two sentences of the old body that carry no keyword of their own beside one sentence that gained SHALL; they are kept together so that no new keyword is added to the render-result sentence.

Form c2 for every requirement below: REMOVED under the old name, ADDED under new names, new blocks at the end of the spec. `grep -rn --exclude-dir=archive` for each old header found no reference to fix. Old scenario text is copied by script from `origin/main`; item 3 of the guide checklist shows `same` for each.

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| Catalog present enables the platform build | Tests resolve fixtures and their dependencies through the configured mapping, Catalog present enables the platform build | no (byte equal, checked by the item 3 command) |
| Fixture prefix redirected alone | A local-iteration mapping redirects only the fixture prefix, Fixture prefix redirected alone | no (byte equal, checked by the item 3 command) |
| Real-renderer pipeline validated against the registry | An integration test exercises the real renderer against the local registry, Real-renderer pipeline validated against the registry | no (byte equal, checked by the item 3 command) |
| Catalog resolved from the generated platform | The integration test resolves the catalog from the generated platform, Catalog resolved from the generated platform | no (byte equal, checked by the item 3 command) |

New scenarios (restate a clause of the old body, with words of the old text only):

| New scenario | Requirement | Old clause restated |
| --- | --- | --- |
| Exporting only CUE_REGISTRY leaves the CLI on the caller's personal configuration | Publish tooling exports OPM_REGISTRY as well as CUE_REGISTRY | the clause of that requirement body |
| The render result carries no inventory entries | The render result carries no inventory entries and full apply is covered by the Kind-backed e2e suite | the clause of that requirement body |
| Full apply to Ready is left to the Kind-backed suite | The render result carries no inventory entries and full apply is covered by the Kind-backed e2e suite | the clause of that requirement body |

Clause table (old clause, new requirement, sentence of its body; "Edit" lists the words that differ from the old sentence):

| Old requirement | Old clause | New requirement | Sentence | Edit |
| --- | --- | --- | --- | --- |
| CUE_REGISTRY configuration ... | Tests SHALL resolve the fixture modules from whatever registry the ... | Tests resolve fixtures and their dependencies through the... | 1 | none |
| CUE_REGISTRY configuration ... | Because the fixtures and their dependencies now occupy SEPARATE dom... | A local-iteration mapping redirects only the fixture prefix | 1 | the em-dashes are written as "," |
| CUE_REGISTRY configuration ... | Tooling that invokes the publish CLI SHALL export `OPM_REGISTRY` as... | Publish tooling exports OPM_REGISTRY as well as CUE_REGISTRY | 1 | none |
| End-to-end integration tests | At least one integration test MUST exercise the real renderer (`ren... | An integration test exercises the real renderer against t... | 1 | none |
| End-to-end integration tests | The test MUST resolve the catalog from the generated platform (the ... | The integration test resolves the catalog from the genera... | 1 | added the named reference: As in requirement "An integration test exercises the rea...": |
| End-to-end integration tests | The render result carries no inventory entries; the reconciler buil... | The render result carries no inventory entries and full a... | 1 | added the named reference: As in requirement "An integration test exercises the rea...": |
| End-to-end integration tests | Full apply → `Ready=True` on a live cluster is covered by the Kind-... | The render result carries no inventory entries and full a... | 2 | added "SHALL" (SHALL be covered by; Strength table) |

Lengths (old body, new bodies; characters by the guide's rule):

| Old requirement | Old | New |
| --- | --- | --- |
| CUE_REGISTRY configuration for tests | 927 | 265, 382, 276 |
| End-to-end integration tests | 851 | 396, 313, 338 |

### Ledger: workflow-hardening

Form c2 for every requirement below: REMOVED under the old name, ADDED under new names, new blocks at the end of the spec. `grep -rn --exclude-dir=archive` for each old header found no reference to fix. Old scenario text is copied by script from `origin/main`; item 3 of the guide checklist shows `same` for each.

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| Every key reader runs in the release Environment | Every job that reads the release App key runs in the release Environment, Every key reader runs in the release Environment | no (byte equal, checked by the item 3 command) |
| Key not yet moved | Every job that reads the release App key runs in the release Environment, Key not yet moved | no (byte equal, checked by the item 3 command) |
| Module image PR call passes no key | A reusable-workflow call does not pass the release App key, Module image PR call passes no key | no (byte equal, checked by the item 3 command) |
| Cascade PR image build skipped | A write-scope pull_request job does not run for a bot head, Cascade PR image build skipped | no (byte equal, checked by the item 3 command) |
| Dependabot and operator module bot PRs hold no write token | A write-scope pull_request job does not run for a bot head, Dependabot and operator module bot PRs hold no write token | no (byte equal, checked by the item 3 command) |
| Human PR e2e publishes in its own job | A write-scope pull_request job does not run for a bot head, Human PR e2e publishes in its own job | no (byte equal, checked by the item 3 command) |
| Cascade PR e2e runs read-only | A job that runs bot-head code holds only read scopes, Cascade PR e2e runs read-only | no (byte equal, checked by the item 3 command) |
| E2E without a publish still runs the podinfo spec | A job that runs bot-head code holds only read scopes, E2E without a publish still runs the podinfo spec | no (byte equal, checked by the item 3 command) |

Clause table (old clause, new requirement, sentence of its body; "Edit" lists the words that differ from the old sentence):

| Old requirement | Old clause | New requirement | Sentence | Edit |
| --- | --- | --- | --- | --- |
| Release App key read only i... | Every job that reads `secrets.RELEASE_APP_PRIVATE_KEY` SHALL declar... | Every job that reads the release App key runs in the rele... | 1 | none |
| Release App key read only i... | A reusable-workflow call SHALL NOT pass the key through `secrets:` ... | A reusable-workflow call does not pass the release App key | 1 | added the named reference: As in requirement "Every job that reads the release App ...": |
| Bot-head code never holds a... | A job triggered by `pull_request` that holds any write scope or `id... | A write-scope pull_request job does not run for a bot head | 1 | none |
| Bot-head code never holds a... | The code on those heads comes from the release cascade, release-ple... | A job that runs bot-head code holds only read scopes | 1 | added the named reference: As in requirement "A write-scope pull_request job does n...": |

Lengths (old body, new bodies; characters by the guide's rule):

| Old requirement | Old | New |
| --- | --- | --- |
| Release App key read only in the release Environment | 651 | 354, 390 |
| Bot-head code never holds a write token or OIDC | 561 | 192, 448 |

### Ledger: docs-bundle

Form c2 for every requirement below: REMOVED under the old name, ADDED under new names, new blocks at the end of the spec. `grep -rn --exclude-dir=archive` for each old header found no reference to fix. Old scenario text is copied by script from `origin/main`; item 3 of the guide checklist shows `same` for each.

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| The test fixture is not shown as an example | The opm-operator docs project builds from a markdown source and a crd source, The test fixture is not shown as an example | no (byte equal, checked by the item 3 command) |
| An enhancement citation stays plain text | The opm-operator docs project sets the section weight, the reconciledBy entries and plain-text citations, An enhancement citation stays plain text | no (byte equal, checked by the item 3 command) |
| The bundle holds both kinds of page | The bundle carries the authored pages and the generated reference and is built by task docs:bundle, The bundle holds both kinds of page | no (byte equal, checked by the item 3 command) |
| A CRD that is not a CustomResourceDefinition is refused | The bundle carries the authored pages and the generated reference and is built by task docs:bundle, A CRD that is not a CustomResourceDefinition is refused | no (byte equal, checked by the item 3 command) |

New scenarios (restate a clause of the old body, with words of the old text only):

| New scenario | Requirement | Old clause restated |
| --- | --- | --- |
| The project is placed in the docs tree | The opm-operator docs project is placed in the docs tree and versioned from tags | the clause of that requirement body |

Clause table (old clause, new requirement, sentence of its body; "Edit" lists the words that differ from the old sentence):

| Old requirement | Old clause | New requirement | Sentence | Edit |
| --- | --- | --- | --- | --- |
| The operator publishes one ... | The repository SHALL declare one docs-kit project, `opm-operator`, ... | The opm-operator docs project is placed in the docs tree ... | 1 | none |
| The operator publishes one ... | placed in a site version's `/docs/` tree, owning `reference/operato... | The opm-operator docs project is placed in the docs tree ... | 2 | final punctuation written as "." |
| The operator publishes one ... | The repository SHALL declare one docs-kit project, `opm-operator`, ... | The opm-operator docs project builds from a markdown sour... | 1 | none |
| The operator publishes one ... | built from a `markdown` source over `docs/site` and a `crd` source ... | The opm-operator docs project builds from a markdown sour... | 2 | final punctuation written as "." |
| The operator publishes one ... | The repository SHALL declare one docs-kit project, `opm-operator`, ... | The opm-operator docs project sets the section weight, th... | 1 | none |
| The operator publishes one ... | the section weight 7, one `reconciledBy` entry per kind, and enhanc... | The opm-operator docs project sets the section weight, th... | 2 | none |
| The operator publishes one ... | The bundle SHALL carry the authored pages under `docs/site/` and th... | The bundle carries the authored pages and the generated r... | 1 | added the named reference: As in requirement "The opm-operator docs project is plac...": |

Lengths (old body, new bodies; characters by the guide's rule):

| Old requirement | Old | New |
| --- | --- | --- |
| The operator publishes one docs bundle | 884 | 198, 387, 199, 371 |

### Ledger: inventory-bridge

Form c2 for every requirement below: REMOVED under the old name, ADDED under new names, new blocks at the end of the spec. `grep -rn --exclude-dir=archive` for each old header found no reference to fix. Old scenario text is copied by script from `origin/main`; item 3 of the guide checklist shows `same` for each.

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| The stored digest changes once on upgrade | The inventory digest is inventory.Digest over the rendered set the last apply recorded, The stored digest changes once on upgrade | no (byte equal, checked by the item 3 command) |
| The digest after an expired Job left the entries | The inventory digest is inventory.Digest over the rendered set the last apply recorded, The digest after an expired Job left the entries | no (byte equal, checked by the item 3 command) |
| Order does not matter | A reader compares the stored digest with the digest of a render, Order does not matter | no (byte equal, checked by the item 3 command) |
| Content does | A reader compares the stored digest with the digest of a render, Content does | no (byte equal, checked by the item 3 command) |

New scenarios (restate a clause of the old body, with words of the old text only):

| New scenario | Requirement | Old clause restated |
| --- | --- | --- |
| A reader does not recompute the digest from the stored entries | A reader compares the stored digest with the digest of a render | the clause of that requirement body |

Clause table (old clause, new requirement, sentence of its body; "Edit" lists the words that differ from the old sentence):

| Old requirement | Old clause | New requirement | Sentence | Edit |
| --- | --- | --- | --- | --- |
| The inventory digest is the... | `status.inventory.digest` MUST be `inventory.Digest` from `opm/k8s/... | The inventory digest is inventory.Digest over the rendere... | 1 | none |
| The inventory digest is the... | It equals the digest of `status.inventory.entries`, with one except... | The inventory digest is inventory.Digest over the rendere... | 2 | none |
| The inventory digest is the... | Source: enhancement 0012. | The inventory digest is inventory.Digest over the rendere... | 3 | none |
| The inventory digest is the... | A reader MUST compare the stored digest with the digest of a render... | A reader compares the stored digest with the digest of a ... | 1 | added the named reference: As in requirement "The inventory digest is inventory.Dig...": |
| The inventory digest is the... | That digest hashes a versioned canonical encoding of each entry's f... | A reader compares the stored digest with the digest of a ... | 2 | none |
| The inventory digest is the... | Source: enhancement 0012. | A reader compares the stored digest with the digest of a ... | 3 | none |

Lengths (old body, new bodies; characters by the guide's rule):

| Old requirement | Old | New |
| --- | --- | --- |
| The inventory digest is the library's canonical digest | 752 | 431, 454 |

### Ledger: live-flux-e2e

Form c2 for every requirement below: REMOVED under the old name, ADDED under new names, new blocks at the end of the spec. `grep -rn --exclude-dir=archive` for each old header found no reference to fix. Old scenario text is copied by script from `origin/main`; item 3 of the guide checklist shows `same` for each.

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| A bare-version claim is accepted, activated and built | The suite asserts the claim, the registry entry and the Platform state, A bare-version claim is accepted, activated and built | no (byte equal, checked by the item 3 command) |
| A v-prefixed claim is accepted and built | The suite asserts the Platform registry entry under a new package identity, A v-prefixed claim is accepted and built | no (byte equal, checked by the item 3 command) |
| Teardown releases the claim | The last step deletes the provider instance and teardown removes any remaining claim, Teardown releases the claim | no (byte equal, checked by the item 3 command) |

New scenarios (restate a clause of the old body, with words of the old text only):

| New scenario | Requirement | Old clause restated |
| --- | --- | --- |
| The suite applies the sample Platform and the provider fixture | The suite drives a claim into a platform build | the clause of that requirement body |
| The expected version is not a literal | The suite reads the expected catalog from the fixture's identity package | the clause of that requirement body |
| The provider is suspended and the claim version gets a v prefix | The v-prefixed claim is covered | the clause of that requirement body |
| The claim is re-judged at its new generation | The claim is re-judged at its new generation | the clause of that requirement body |
| A moved package identity shows the regeneration | The package identity shows the platform was regenerated for the edited claim | the clause of that requirement body |
| The spec needs no registry credentials on the public default registry | The spec deploys and tears down its own controller | the clause of that requirement body |
| The spec uses the LOCAL_REGISTRY and OPERATOR_DOCKER_CONFIG overrides | The spec deploys and tears down its own controller | the clause of that requirement body |

Clause table (old clause, new requirement, sentence of its body; "Edit" lists the words that differ from the old sentence):

| Old requirement | Old clause | New requirement | Sentence | Edit |
| --- | --- | --- | --- | --- |
| Deployed-controller registr... | Against the deployed controller, the e2e suite SHALL drive a claim ... | The suite drives a claim into a platform build | 1 | none |
| Deployed-controller registr... | The suite SHALL read the expected catalog path and version from the... | The suite reads the expected catalog from the fixture's i... | 1 | added the named reference: As in requirement "The suite drives a claim into a platf...": |
| Deployed-controller registr... | The suite SHALL assert three things: the claim reports `accepted: t... | The suite asserts the claim, the registry entry and the P... | 1 | added the named reference: As in requirement "The suite drives a claim into a platf...": |
| Deployed-controller registr... | The suite SHALL then cover the `v`-prefixed spelling. It suspends t... | The v-prefixed claim is covered | 1 | added the named reference: As in requirement "The suite drives a claim into a platf...": |
| Deployed-controller registr... | The suite SHALL assert that the claim is re-judged at its new gener... | The claim is re-judged at its new generation | 1 | added the named reference: As in requirement "The v-prefixed claim is covered": |
| Deployed-controller registr... | Only after that verdict SHALL the suite assert that the Platform's ... | The suite asserts the Platform registry entry under a new... | 1 | added the named reference: As in requirement "The claim is re-judged at its new gen...": |
| Deployed-controller registr... | The package identity is a function of the claims' coordinates, so i... | The package identity shows the platform was regenerated f... | 1 | added the named reference: As in requirement "The suite asserts the Platform regist...":; added "SHALL" (SHALL be a function; Strength table) |
| Deployed-controller registr... | The spec SHALL deploy and tear down its own controller. It SHALL ne... | The spec deploys and tears down its own controller | 1 | added the named reference: As in requirement "The suite drives a claim into a platf...": |
| Deployed-controller registr... | Its last ordered step SHALL delete the provider instance and assert... | The last step deletes the provider instance and teardown ... | 1 | added the named reference: As in requirement "The spec deploys and tears down its o...": |

Lengths (old body, new bodies; characters by the guide's rule):

| Old requirement | Old | New |
| --- | --- | --- |
| Deployed-controller registration claim proof | 2782 | 384, 324, 448, 290, 500, 417, 241, 335, 400 |

### Ledger: module-acquisition

Form c2 for every requirement below: REMOVED under the old name, ADDED under new names, new blocks at the end of the spec. `grep -rn --exclude-dir=archive` for each old header found no reference to fix. Old scenario text is copied by script from `origin/main`; item 3 of the guide checklist shows `same` for each.

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| Published module is acquired | The operator acquires a module from the registry as a library Module, Published module is acquired | no (byte equal, checked by the item 3 command) |
| Unresolvable module surfaces a load error | The operator acquires a module from the registry as a library Module, Unresolvable module surfaces a load error | no (byte equal, checked by the item 3 command) |
| Acquisition delegates to the kernel, not a wrapper package | Acquisition delegates to the library kernel's registry module loader, Acquisition delegates to the kernel, not a wrapper package | no (byte equal, checked by the item 3 command) |
| Core@v0 self-referential metadata is preserved | Acquisition synthesizes no wrapper package and writes no temporary directory, Core@v0 self-referential metadata is preserved | no (byte equal, checked by the item 3 command) |
| Module with catalog imports is acquired | Acquisition synthesizes no wrapper package and writes no temporary directory, Module with catalog imports is acquired | no (byte equal, checked by the item 3 command) |

New scenarios (restate a clause of the old body, with words of the old text only):

| New scenario | Requirement | Old clause restated |
| --- | --- | --- |
| Acquisition writes no temporary directory | Acquisition synthesizes no wrapper package and writes no temporary directory | the clause of that requirement body |

Clause table (old clause, new requirement, sentence of its body; "Edit" lists the words that differ from the old sentence):

| Old requirement | Old clause | New requirement | Sentence | Edit |
| --- | --- | --- | --- | --- |
| Acquire a module from the r... | The operator SHALL acquire a `ModuleRelease`'s target module from t... | The operator acquires a module from the registry as a lib... | 1 | none |
| Acquire a module from the r... | Registry configuration and the absence of process-environment mutat... | The operator acquires a module from the registry as a lib... | 2 | none |
| Acquire a module from the r... | Acquisition SHALL delegate to the library kernel's registry module ... | Acquisition delegates to the library kernel's registry mo... | 1 | the em-dashes are written as "," |
| Acquire a module from the r... | Delegating to the library is load-bearing: the wrapper approach re-... | Acquisition delegates to the library kernel's registry mo... | 2 | none |
| Acquire a module from the r... | Acquisition SHALL NOT synthesize a wrapper CUE package that imports... | Acquisition synthesizes no wrapper package and writes no ... | 1 | added the named reference: As in requirement "Acquisition delegates to the library ...": |
| Acquire a module from the r... | The library loader loads the module as the main module, its own `cu... | Acquisition synthesizes no wrapper package and writes no ... | 2 | the em-dashes are written as "," |

Lengths (old body, new bodies; characters by the guide's rule):

| Old requirement | Old | New |
| --- | --- | --- |
| Acquire a module from the registry as a library Module | 1229 | 371, 466, 476 |

### Ledger: operator-resource-reference

Form c2 for every requirement below: REMOVED under the old name, ADDED under new names, new blocks at the end of the spec. `grep -rn --exclude-dir=archive` for each old header found no reference to fix. Old scenario text is copied by script from `origin/main`; item 3 of the guide checklist shows `same` for each.

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| Each kind has its own page | The operator resource reference is generated by docs-kit, Each kind has its own page | no (byte equal, checked by the item 3 command) |
| The authored intro completes the page | The committed index holds only authored front matter and intro, The authored intro completes the page | no (byte equal, checked by the item 3 command) |
| A generated heading in the authored page is refused | The committed index holds only authored front matter and intro, A generated heading in the authored page is refused | no (byte equal, checked by the item 3 command) |

Clause table (old clause, new requirement, sentence of its body; "Edit" lists the words that differ from the old sentence):

| Old requirement | Old clause | New requirement | Sentence | Edit |
| --- | --- | --- | --- | --- |
| The resource reference is g... | The operator resource reference SHALL be generated by docs-kit's `c... | The operator resource reference is generated by docs-kit | 1 | none |
| The resource reference is g... | The committed `docs/site/reference/operator/_index.md` SHALL hold o... | The committed index holds only authored front matter and ... | 1 | added the named reference: As in requirement "The operator resource reference is ge...": |

Lengths (old body, new bodies; characters by the guide's rule):

| Old requirement | Old | New |
| --- | --- | --- |
| The resource reference is generated into the docs bundle | 732 | 491, 318 |

### Ledger: platform-gated-rendering

Form c2 for every requirement below: REMOVED under the old name, ADDED under new names, new blocks at the end of the spec. `grep -rn --exclude-dir=archive` for each old header found no reference to fix. Old scenario text is copied by script from `origin/main`; item 3 of the guide checklist shows `same` for each.

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| Blocked instances retry when the platform is generated | Fields that re-enqueue, Blocked instances retry when the platform is generated | no (byte equal, checked by the item 3 command) |
| Instances re-render under a new pin set | Fields that re-enqueue, Instances re-render under a new pin set | no (byte equal, checked by the item 3 command) |
| An operator upgrade into an already-Ready Platform recovers blocked instances | The Platform reconciler's own status update and a change to operatorVersion re-enqueue, An operator upgrade into an already-Ready Platform recovers blocked instances | no (byte equal, checked by the item 3 command) |
| A message-only status write enqueues nothing | A Platform update that changes none of the fields enqueues nothing, A message-only status write enqueues nothing | no (byte equal, checked by the item 3 command) |
| A spec edit renders once the regenerated package lands | A bump of the Platform generation alone enqueues nothing, A spec edit renders once the regenerated package lands | no (byte equal, checked by the item 3 command) |
| CLI-owned and suspended instances are not enqueued | A Platform event enqueues only the operator-managed, unsuspended ModuleInstances, CLI-owned and suspended instances are not enqueued | no (byte equal, checked by the item 3 command) |

Clause table (old clause, new requirement, sentence of its body; "Edit" lists the words that differ from the old sentence):

| Old requirement | Old clause | New requirement | Sentence | Edit |
| --- | --- | --- | --- | --- |
| Re-enqueue ModuleInstances ... | The ModuleInstance reconciler SHALL watch the `Platform`. It SHALL ... | Fields that re-enqueue | 1 | none |
| Re-enqueue ModuleInstances ... | - the `Ready` condition's status | Fields that re-enqueue | 2 | none |
| Re-enqueue ModuleInstances ... | - the pin set: `status.packageIdentity` or `status.registry` | Fields that re-enqueue | 3 | none |
| Re-enqueue ModuleInstances ... | - `spec.skewPolicy` | Fields that re-enqueue | 4 | none |
| Re-enqueue ModuleInstances ... | - `status.observedGeneration` | Fields that re-enqueue | 5 | none |
| Re-enqueue ModuleInstances ... | `status.packageIdentity` is the field that identifies the pin set a... | Fields that re-enqueue | 6 | none |
| Re-enqueue ModuleInstances ... | The trigger includes the Platform reconciler's own status update, w... | The Platform reconciler's own status update and a change ... | 1 | added the named reference: As in requirement "Fields that re-enqueue": |
| Re-enqueue ModuleInstances ... | A Platform update that changes none of these fields SHALL NOT enque... | A Platform update that changes none of the fields enqueue... | 1 | added the named reference: As in requirements "Fields that re-enqueue and The Platfo...": |
| Re-enqueue ModuleInstances ... | Examples of such updates are | A Platform update that changes none of the fields enqueue... | 2 | none |
| Re-enqueue ModuleInstances ... | a change to the `Ready` message or to its reason while it stays `Fa... | A Platform update that changes none of the fields enqueue... | 3 | final punctuation written as "." |
| Re-enqueue ModuleInstances ... | A Platform update that changes none of these fields SHALL NOT enque... | A bump of the Platform generation alone enqueues nothing | 1 | added the named reference: As in requirements "Fields that re-enqueue and The Platfo...": |
| Re-enqueue ModuleInstances ... | Examples of such updates are | A bump of the Platform generation alone enqueues nothing | 2 | none |
| Re-enqueue ModuleInstances ... | a bump of `metadata.generation` alone: the platform has not been re... | A bump of the Platform generation alone enqueues nothing | 3 | none |
| Re-enqueue ModuleInstances ... | Platform create and delete events SHALL re-enqueue. | A Platform event enqueues only the operator-managed, unsu... | 1 | added the named reference: As in requirement "Fields that re-enqueue": |
| Re-enqueue ModuleInstances ... | The reconciler SHALL enqueue only the ModuleInstances that render a... | A Platform event enqueues only the operator-managed, unsu... | 2 | none |

Lengths (old body, new bodies; characters by the guide's rule):

| Old requirement | Old | New |
| --- | --- | --- |
| Re-enqueue ModuleInstances when the platform becomes ready | 1745 | 471, 426, 360, 479, 485 |

### Ledger: registration-activation

Form c2 for every requirement below: REMOVED under the old name, ADDED under new names, new blocks at the end of the spec. `grep -rn --exclude-dir=archive` for each old header found no reference to fix. Old scenario text is copied by script from `origin/main`; item 3 of the guide checklist shows `same` for each.

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| A provider going unready leaves the claim active | An active claim stays active when its provider's readiness regresses, A provider going unready leaves the claim active | no (byte equal, checked by the item 3 command) |
| A provider recovering does not re-run the gate | An active claim stays active when its provider's readiness regresses, A provider recovering does not re-run the gate | no (byte equal, checked by the item 3 command) |

New scenarios (restate a clause of the old body, with words of the old text only):

| New scenario | Requirement | Old clause restated |
| --- | --- | --- |
| The readiness gate does not govern an active claim | The readiness gate governs initial activation only | the clause of that requirement body |

Clause table (old clause, new requirement, sentence of its body; "Edit" lists the words that differ from the old sentence):

| Old requirement | Old clause | New requirement | Sentence | Edit |
| --- | --- | --- | --- | --- |
| Activation latches against ... | Once active, a claim SHALL NOT be deactivated because its provider'... | An active claim stays active when its provider's readines... | 1 | none |
| Activation latches against ... | This is deliberate and not an oversight: the active-claim set is an... | An active claim stays active when its provider's readines... | 2 | none |
| Activation latches against ... | The readiness gate SHALL govern initial activation only. | The readiness gate governs initial activation only | 1 | added the named reference: As in requirement "An active claim stays active when its...": |
| Activation latches against ... | A claim leaves the active state by deletion, never by its provider ... | The readiness gate governs initial activation only | 2 | none |
| Activation latches against ... | Established CRDs outlive a provider's pods, so an active claim stay... | The readiness gate governs initial activation only | 3 | none |

Lengths (old body, new bodies; characters by the guide's rule):

| Old requirement | Old | New |
| --- | --- | --- |
| Activation latches against provider health | 623 | 363, 349 |

### Ledger: transformer-registration-crd

Form c2 for every requirement below: REMOVED under the old name, ADDED under new names, new blocks at the end of the spec. `grep -rn --exclude-dir=archive` for each old header found no reference to fix. Old scenario text is copied by script from `origin/main`; item 3 of the guide checklist shows `same` for each.

| Old scenario name | New place (requirement, scenario name) | Changed text? |
| --- | --- | --- |
| A claim missing its catalog is rejected | TransformerRegistrationSpec requires catalog, version, provides and providerRef, A claim missing its catalog is rejected | no (byte equal, checked by the item 3 command) |
| A claim with an empty provides list is accepted | TransformerRegistrationSpec requires catalog, version, provides and providerRef, A claim with an empty provides list is accepted | no (byte equal, checked by the item 3 command) |

New scenarios (restate a clause of the old body, with words of the old text only):

| New scenario | Requirement | Old clause restated |
| --- | --- | --- |
| A claim can reach the cluster with a field absent | Claim validation is not delegated to the authoring catalog | the clause of that requirement body |

Clause table (old clause, new requirement, sentence of its body; "Edit" lists the words that differ from the old sentence):

| Old requirement | Old clause | New requirement | Sentence | Edit |
| --- | --- | --- | --- | --- |
| Every claim field is requir... | `TransformerRegistrationSpec` SHALL require `catalog` (the provider... | TransformerRegistrationSpec requires catalog, version, pr... | 1 | none |
| Every claim field is requir... | `provides` is required but MAY be empty: the CRD sets no minimum on... | TransformerRegistrationSpec requires catalog, version, pr... | 2 | none |
| Every claim field is requir... | This validation SHALL NOT be delegated to the authoring catalog: a ... | Claim validation is not delegated to the authoring catalog | 1 | added the named reference: As in requirement "TransformerRegistrationSpec requires ...": |

Lengths (old body, new bodies; characters by the guide's rule):

| Old requirement | Old | New |
| --- | --- | --- |
| Every claim field is required at admission | 733 | 491, 342 |

## Strength table

Every sentence that gained SHALL, MUST or MAY in this unit.
Old text from `origin/main`. All other keywords are the old ones.

| Requirement (new) | Sentence | Old sentence, no keyword |
| --- | --- | --- |
| No ClaimsKept event when nothing was kept | "The controller MUST NOT emit a `ClaimsKept` event when nothing was kept: ..." | "No `ClaimsKept` event is emitted when nothing was kept: ..." (events-emission, "Events emitted when claims are kept") |
| No LeftBehind event for an absent object or a kept claim | "The controller MUST NOT emit a `LeftBehind` event for an object that was already absent or for a kept PersistentVolumeClaim, which has its own event." | "No `LeftBehind` event is emitted for an object that was already absent or for a kept PersistentVolumeClaim, which has its own event." (events-emission, "Events emitted when objects are left behind") |
| No AdoptedElsewhere event while the count is unchanged | "A reconcile that finds the same number as before, or that skips its render, MUST emit no `AdoptedElsewhere` event: ..." | "A reconcile that finds the same number as before, or that skips its render, emits none: ..." |
| Recovery never moves a tag or edits a release by hand | "Recovery SHALL never move, delete or re-create a tag and SHALL never edit a release by hand." | "Recovery never moves, deletes or re-creates a tag and never edits a release by hand" |
| The recovery dispatch publishes, attaches and releases in order | "A recovery dispatch SHALL run, in order: the publish, ..." | "In order: the publish, ..." (operator-module-release, "A stranded module draft is recovered by dispatch"); the subject "A recovery dispatch" comes from the sentence before it |
| The module version follows the 0.x bump rule on the v0 path | none gained; "A module release is breaking when ..." stays without a keyword, beside a SHALL | (no change) |
| Healthy is independent of Ready | "`Healthy` SHALL be independent of `Ready`: no helper that sets `Ready`, `Reconciling` or `Stalled` changes it, except that ..." | "`Healthy` is independent of `Ready`: no helper that sets `Ready`, `Reconciling` or `Stalled` changes it, except that ..." (status-conditions, "The Healthy condition vocabulary") |
| A failure is transient when its error chain holds a registry fetch failure | "A failure SHALL be transient when its error chain holds ..." | "A failure is transient when its error chain holds ..." (platform-reconciler, "Transient failures retry faster than semantic failures") |
| The demand is read only from a successful render | "The demand SHALL be read only from a successful render; this requirement does not change when the reconciler writes it ..." | "The demand is read only from a successful render; this requirement does not change when the reconciler writes it ..." (kernel-module-renderer, "The render reports the instance's contract demand from the kernel") |
| The health package covers applied inventory objects only | "... the operator's own resources (a `dependsOn` ModulePackage, a TransformerRegistration's provider, the `cluster` Platform) SHALL keep their own condition checks." | "... keep their own condition checks." (kubernetes-tier-adoption, "The operator judges the health of inventory objects only through the library's health package") |
| An object whose kind a CRD of the same apply list defines counts as an object that does not exist | "An object whose kind the API server does not serve yet, while a CustomResourceDefinition of the same apply list defines that kind, SHALL count as an object that does not exist." | "An object whose kind the API server does not serve yet, while a CustomResourceDefinition of the same apply list defines that kind, counts as an object that does not exist." (ssa-apply, "An apply is judged by the ownership verdict before its first write") |
| The identity of the apply verdict is the same for a ModuleInstance and a ModulePackage | "The identity the verdict is asked with MUST be the same for a ModuleInstance and a ModulePackage, and requirements ... are its one definition." | "The identity the verdict is asked with is the same for a ModuleInstance and a ModulePackage, and this paragraph is its one definition." (ssa-apply, "An apply is judged by the ownership verdict before its first write") |
| A refused claim update applies, creates and deletes no object of the set | "In that case no object of the resource set SHALL be applied, created or deleted." | "In that case no object of the resource set is applied, created or deleted." (ssa-apply, "A forced recreate keeps PersistentVolumeClaims unless data deletion is allowed") |
| The event carries the text of the values failure | "The event SHALL carry the same text within the event note limit (capability `events-emission`)." | "The event carries the same text within the event note limit (capability `events-emission`)." (module-instance-synthesis, "A values failure reports every finding with its positions") |
| A reconcile whose digests match applies the rendered set and an inventoried object being deleted refuses it | "Such a reconcile SHALL apply the rendered set, so an inventoried object that is being deleted refuses it ..." | "Such a reconcile applies the rendered set, so an inventoried object that is being deleted refuses it ..." (modulepackage-reconcile-loop, "No-op detection") |
| The render result carries no inventory entries and full apply is covered by the Kind-backed e2e suite | "Full apply → `Ready=True` on a live cluster SHALL be covered by the Kind-backed `test/e2e` suite, not this integration-tier test." | "Full apply → `Ready=True` on a live cluster is covered by the Kind-backed `test/e2e` suite, not this integration-tier test." (test-registry-lifecycle, "End-to-end integration tests") |
| The package identity shows the platform was regenerated for the edited claim | "The package identity SHALL be a function of the claims' coordinates, so its move is what shows the platform was regenerated for the edited claim." | "The package identity is a function of the claims' coordinates, so its move is what shows the platform was regenerated for the edited claim." (live-flux-e2e, "Deployed-controller registration claim proof") |

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

## The pinned gate (Taskfile and CI): moved out

The pinned OpenSpec 1.14.1 gate (`openspec:install` and `openspec:check` tasks and a `Lint` step)
is not part of this change. It becomes a small separate change and PR after the open changes
`stop-operator-install-manifest` and `add-cross-namespace-source-grants` land, because a step
that runs `--all --strict` would fail its own PR while they carry warnings and errors
That change installs OpenSpec 1.14.1 pinned by version with
`--ignore-scripts` and no lockfile, and takes the `cli` repo's tasks as the model.

## Risks / Trade-offs

- A rewrite that loses a rule while the validator is green. Counter: ledger, byte-equality
  command, clause table, and a review per unit (guide section 7).
- The two open changes move text before B1, B3 and B5. Counter: the pause and the fresh
  inventory. An "Archive would refuse" line is reported, never repaired by editing the other
  change.
- Section count: 11 task groups (9 units plus a pause and the final gate), against the config's "at most about seven sections" rule. The
  groups are mechanical units of one kind, each ends green, and the swarm sets one change per repo run as several units. The proposal gate may ask instead for two
  changes (before and after the open changes land).
- No new trust boundary, secret or personal data.

## Open Questions

- Does the proposal gate accept 11 groups in one change, or want two changes?
- The counts of B1, B3 and B5 are known only after the fresh inventory.
