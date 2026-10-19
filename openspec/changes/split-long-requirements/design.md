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

`openspec archive` ignores a delta Purpose for an existing spec (DECISIONS.md 2026-10-19 14:00,
verified in the CLI's `specs-apply`). A Purpose is therefore a direct edit of `## Purpose` in
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

Ledgers written so far: unit 1 (opm-operator-B2) and unit 2 (opm-operator-B4), below.

Unit 2 (opm-operator-B4) rules of text, for the five ledgers below:

- Em-dash characters of the old bodies (four places: `module-instance-ownership` "The operator acknowledges a CLI-owned instance with a single condition", `platform-reconciler` "Generate and build the singleton Platform on reconcile", `kernel-module-renderer` "Adapt compiled output to operator resources", twice) are written as `:` or `,` in the new bodies, because the workspace bans the character. The word and the order stay.
- A sentence whose subject was "It" gets the subject of the sentence before it, copied (the "Edit" column shows it).
- Where a body relies on a sibling for a condition, the condition is copied from the old text in front of its sentence (the "Edit" column shows it).
- A sentence over 500 on its own (`Reason constants`, the outcome list of `Surface build outcome on status`, the identity sentence of `Single-slot identity-keyed store for concurrent read`) is cut at a list boundary; the parts are verbatim and the head is repeated.
- `Source:` lines: the whole line stays where it fits with the rule it supports, and bodies are split further until it fits (the enhancement id is never alone, `openspec/config.yaml`). The old `Source:` line covers the whole old requirement, so every split block of that requirement carries it (the library-owns-the-deletion-sequence title for the prune and cleanup deletes; the apply guard title for the Ready message). A block keeps the short form `Source: enhancement 0012.` only where one sentence of the old requirement, with its title, passes 500 characters and has no clause boundary that keeps its condition with its rule: the Ready message body of `status-conditions` ("The Ready message counts the objects adopted by another instance", 335 plus 208) and the two AdoptedElsewhere bodies of `events-emission` (319 plus 208; 278 plus 65 plus 208).

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
| An object is let go | An AdoptedElsewhere event is emitted when the count of adopted objects changes, An object is let go | no (byte equal, checked by the item 3 command) |
| No event while nothing changes | No AdoptedElsewhere event while the count is unchanged, No event while nothing changes | no (byte equal, checked by the item 3 command) |
| A second object is let go | The AdoptedElsewhere event states the count and the library messages, A second object is let go | no (byte equal, checked by the item 3 command) |
| Nothing adopted elsewhere, no event | An AdoptedElsewhere event is emitted when the count of adopted objects changes, Nothing adopted elsewhere, no event | no (byte equal, checked by the item 3 command) |
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
| AdoptedElsewhere event on a no-op reconcile | The AdoptedElsewhere event is emitted whether the reconcile applies, restores or does nothing | the clause of that requirement body |
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
| Events emitted when an object is adopted by another instance | one Warning `AdoptedElsewhere`, action `Apply`, "when a reconcile that renders finds a number of rendered objects that the apply verdict refuses as `adopted-elsewhere`, that number is not zero, and it differs from the number the `Ready` message stated before the reconcile" (whole old condition, copied unchanged) | An AdoptedElsewhere event is emitted when the count of adopted objects changes | 1 |
| Events emitted when an object is adopted by another instance | citation `status-conditions`, "The Ready message counts the objects adopted by another instance" (moves to the no-event requirement, where the Ready message carries the state) | No AdoptedElsewhere event while the count is unchanged | 1 |
| Events emitted when an object is adopted by another instance | "This holds whether the reconcile applies, restores or is a no-op." (old sentence, copied) with the same old condition as the row above, copied unchanged, as its subject | The AdoptedElsewhere event is emitted whether the reconcile applies, restores or does nothing | 1 |
| Events emitted when an object is adopted by another instance | states the count; the library's message for each object unchanged; within the `ApplyRefused` limits (limits restated) | The AdoptedElsewhere event states the count and the library messages | 1 |
| Events emitted when an object is adopted by another instance | same number as before, or a skipped render, emits none; `Ready` carries the state between events | No AdoptedElsewhere event while the count is unchanged | 1 |
| Events emitted when an object is adopted by another instance | `Source:` line (enhancement 0012): whole on the count, the message and the no-event requirements; `Source: enhancement 0012.` (title dropped, no issue or PR reference lost) on the two that carry the old condition, because the whole line does not fit 500 with the whole condition | An AdoptedElsewhere event is emitted when the count of adopted objects changes | 2 |
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
| Events emitted when an object is adopted by another instance | 969 | 345, 370, 492, 467 |
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
| The Ready message counts the objects adopted by another instance | 746 | 477, 319 |

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

## Strength table

Every sentence that gained SHALL, MUST or MAY in this unit (DECISIONS.md 2026-10-19 16:20 and 17:05).
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
(DECISIONS.md 2026-10-19 14:05). That change installs OpenSpec 1.14.1 pinned by version with
`--ignore-scripts` and no lockfile, and takes the `cli` repo's tasks as the model.

## Risks / Trade-offs

- A rewrite that loses a rule while the validator is green. Counter: ledger, byte-equality
  command, clause table, and a review per unit (guide section 7).
- The two open changes move text before B1, B3 and B5. Counter: the pause and the fresh
  inventory. An "Archive would refuse" line is reported, never repaired by editing the other
  change.
- Section count: 11 task groups (9 units plus a pause and the final gate), against the config's "at most about seven sections" rule. The
  groups are mechanical units of one kind, each ends green, and DECISIONS.md 2026-10-19 09:50
  sets one change per repo run as several units. The proposal gate may ask instead for two
  changes (before and after the open changes land).
- No new trust boundary, secret or personal data.

## Open Questions

- Does the proposal gate accept 11 groups in one change, or want two changes?
- The counts of B1, B3 and B5 are known only after the fresh inventory.
