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

Ledgers written so far: unit 1 (opm-operator-B2), below.

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

## Strength table

Every sentence that gained SHALL, MUST or MAY in this unit (DECISIONS.md 2026-10-19 16:20 and 17:05).
Old text from `origin/main`. All other keywords are the old ones.

| Requirement (new) | Sentence | Old sentence, no keyword |
| --- | --- | --- |
| No ClaimsKept event when nothing was kept | "The controller MUST NOT emit a `ClaimsKept` event when nothing was kept: ..." | "No `ClaimsKept` event is emitted when nothing was kept: ..." (events-emission, "Events emitted when claims are kept") |
| No LeftBehind event for an absent object or a kept claim | "The controller MUST NOT emit a `LeftBehind` event for an object that was already absent or for a kept PersistentVolumeClaim, ..." | "none for an absent object or a kept PVC, which has its own event" |
| No AdoptedElsewhere event while the count is unchanged | "A reconcile that finds the same number as before, or that skips its render, MUST emit no `AdoptedElsewhere` event: ..." | "A reconcile that finds the same number as before, or that skips its render, emits none: ..." |
| Recovery never moves a tag or edits a release by hand | "Recovery SHALL never move, delete or re-create a tag and SHALL never edit a release by hand." | "Recovery never moves, deletes or re-creates a tag and never edits a release by hand" |
| The module version follows the 0.x bump rule on the v0 path and the sentences kept without a keyword | none gained; "A module release is breaking when ..." and the recovery sentences stay without a keyword, beside a SHALL | (no change) |

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
