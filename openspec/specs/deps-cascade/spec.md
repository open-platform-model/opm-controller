# deps-cascade Specification

## Purpose
The controller's `task deps:cascade`: how it moves the repo's upstream pins (library, the opm catalog and core with fixtures and consumers, and the opm CLI) to the newest published versions, which pins it leaves alone, how it reports by exit code, and how it is tested.

## Requirements

### Requirement: deps:cascade moves pins in the working tree and reports by exit code
The repo SHALL provide `task deps:cascade`, which moves opm-controller's upstream pins to the newest published versions in the working tree only. It SHALL NOT commit, push or create a branch. Run as `task -x deps:cascade`, it SHALL exit 0 when the working tree changed, 3 when there was nothing to do, and any other code on error, and it SHALL never turn a failure into 0 or 3. It SHALL move exactly the library, opm catalog, core and opm CLI pins.

#### Scenario: Up-to-date tree
- **WHEN** every pin already equals the newest published version the resolver reports
- **THEN** `task -x deps:cascade` exits 3 and `git status --porcelain` is empty

#### Scenario: A pin is behind
- **WHEN** the resolver reports a newer published library than `go.mod` pins
- **THEN** `task -x deps:cascade` exits 0 and `go.mod` and `go.sum` name the newer library

#### Scenario: A predicate's "no" is not "nothing to do"
- **WHEN** a resolver call other than `newest`, such as `pin-of` for the target catalog, exits 3
- **THEN** the task exits with a code other than 0 or 3 and `git status --porcelain` is empty

### Requirement: deps:cascade takes every version decision from the shared cascade resolver
Every version decision of `task deps:cascade` SHALL come from the shared cascade resolver of `open-platform-model/.github` (workspace RELEASING.md, section "The cascade"). The resolver is found through the absolute path in `CASCADE_RESOLVER`, or otherwise at `../.github/.github/scripts/cascade/cascade-resolve.sh` beside the main checkout. When the resolver is missing, the task SHALL fail with a message naming both options.

#### Scenario: Resolver not found
- **WHEN** `CASCADE_RESOLVER` is unset and no `.github` checkout sits beside the repo
- **THEN** the task fails before running and says to check out `open-platform-model/.github` or set `CASCADE_RESOLVER`

#### Scenario: Relative resolver path
- **WHEN** `CASCADE_RESOLVER` is set to a relative path
- **THEN** the task fails and says the path must be absolute

### Requirement: deps:cascade refuses a dirty working tree
`task deps:cascade` SHALL refuse to run (exit 1) on a working tree with any change or untracked file, unless `CASCADE_ALLOW_DIRTY=1`. In that case it SHALL judge its result by comparing a snapshot of the tree from before and after the run.

#### Scenario: Dirty tree refused
- **WHEN** the working tree holds an untracked file and `CASCADE_ALLOW_DIRTY` is unset
- **THEN** the task exits 1 and changes nothing

### Requirement: deps:cascade resolves every target before it edits
`task deps:cascade` SHALL validate `.cascade-frozen` and `.cascade-hold` through the resolver first. It SHALL resolve every target before editing any file.

#### Scenario: Resolver error leaves the tree untouched
- **WHEN** the resolver fails on the first pin the task resolves
- **THEN** the task exits with a code other than 0 or 3 and `git status --porcelain` is empty

### Requirement: deps:cascade sets its own registry environment
`task deps:cascade` SHALL set its own `CUE_REGISTRY` and `OPM_REGISTRY` to `testing.opmodel.dev=ghcr.io/open-platform-model,opmodel.dev=ghcr.io/open-platform-model,registry.cue.works`, and SHALL never inherit the Taskfile's local `localhost:5000` default.

#### Scenario: Own registry mapping
- **WHEN** the Taskfile's default registry is the local `localhost:5000` and `task deps:cascade` runs
- **THEN** `CUE_REGISTRY` and `OPM_REGISTRY` are set to `testing.opmodel.dev=ghcr.io/open-platform-model,opmodel.dev=ghcr.io/open-platform-model,registry.cue.works`, not to `localhost:5000`

### Requirement: deps:cascade moves the library pin
`task deps:cascade` SHALL move the library pin by running `go get github.com/open-platform-model/library@<exact version>` and then `go mod tidy`, only when library moved.

#### Scenario: Library moved
- **WHEN** the resolver reports a newer published library than `go.mod` pins
- **THEN** the task runs `go get github.com/open-platform-model/library@<exact version>` and then `go mod tidy`

#### Scenario: Library not moved
- **WHEN** library did not move
- **THEN** the task runs neither the library `go get` nor the `go mod tidy` that follows it

### Requirement: deps:cascade moves the opm catalog pin to one target
`task deps:cascade` SHALL move `opmodel.dev/catalogs/opm@v4` to one target `K` in every file that pins it and is below `K`: the sample Platform's `version:` (stored bare); `CatalogVersion()` in `test/fixtures/catalog.go` (stored bare); the four `test/fixtures/modules/*/cue.mod/module.cue` files. A file whose catalog is above `K` SHALL NOT be lowered.

#### Scenario: File below the target catalog moves
- **WHEN** a file that pins `opmodel.dev/catalogs/opm@v4` is below the target `K`
- **THEN** the file pins `K`

#### Scenario: File above the target catalog is not lowered
- **WHEN** a file pins a catalog above the target `K`
- **THEN** the file keeps its catalog and is not lowered

### Requirement: deps:cascade moves core to the core its catalog pins
`task deps:cascade` SHALL move core in the four `test/fixtures/modules/*/cue.mod/module.cue` files and in `test/fixtures/catalogs/provider/cue.mod/module.cue` to the version that the file's catalog after the move pins (the higher of the file's catalog and the target catalog `K`; `K` for the provider fixture), and only when that is greater than the file's own core. Core SHALL never come from the newest published core directly.

#### Scenario: Catalog and core move as a consistent set
- **WHEN** the newest published catalog is newer than the sample's, and that catalog pins a core newer than the fixtures'
- **THEN** the sample, `catalog.go` and the four fixture modules name the new catalog, and the four fixture modules and the provider catalog fixture name the core that catalog pins

#### Scenario: Core already ahead of the catalog's core
- **WHEN** a fixture pins a core newer than the core catalog `K` pins
- **THEN** that fixture's core is left as it is and a warning names both versions

#### Scenario: Newer core published but not pinned by the catalog
- **WHEN** core `v2.0.0-beta.2` is published and the newest catalog pins core `v2.0.0-beta.1`
- **THEN** the fixtures' core stays at `v2.0.0-beta.1`

### Requirement: deps:cascade writes the opm CLI pin last
`task deps:cascade` SHALL move the opm CLI pin, `.opm-cli-version`, and SHALL write it last.

#### Scenario: CLI pin written last
- **WHEN** the task moves the opm CLI pin and other pins in one run
- **THEN** `.opm-cli-version` is written after every other pin

### Requirement: deps:cascade leaves the generated controller reference alone
`task deps:cascade` SHALL NOT edit `docs/site/reference/controller/_index.md`: the docs bundle generates the resource reference from `config/samples` when it is built.

#### Scenario: Sample Platform moved
- **WHEN** the task moves the catalog `version:` in the sample Platform
- **THEN** `docs/site/reference/controller/_index.md` is unchanged, and `task docs:bundle` on the resulting tree shows the moved version in the example of `reference/controller/platform.md`

### Requirement: deps:cascade runs cue mod get on OPM modules only
In `task deps:cascade`, `cue mod get` SHALL name only `opmodel.dev/*` and `testing.opmodel.dev/*` modules, each with an exact version, and SHALL run, followed by one `cue mod tidy`, only in a module where a pin moved. Third-party pins SHALL never be named. The task SHALL never edit an import path or a `@vN` key; a new major SHALL appear only as the resolver's warning.

#### Scenario: Third-party pins untouched
- **WHEN** the task runs `cue mod get` in a fixture module
- **THEN** the command names no module outside `opmodel.dev/` and `testing.opmodel.dev/`

### Requirement: deps:cascade warns about changes beyond the pins
`task deps:cascade` SHALL report each of these as a warning: a third-party pin that `tidy` raises, adds or removes; a change to `go.mod`'s `go` or `toolchain` directive; a dep a fixture module gains that its modulepackage lacks.

#### Scenario: Tidy adds a dependency
- **WHEN** `go mod tidy` or `cue mod tidy` adds or removes a third-party dependency, or `go get` raises the `go` or `toolchain` directive
- **THEN** the task writes a warning naming the dependency or the directives, and a fixture module's new dep that its modulepackage lacks is warned about too

### Requirement: deps:cascade skips frozen pins
Before editing a version literal in a file, or naming a key in `cue mod get` for a module, `task deps:cascade` SHALL ask the resolver whether that exact file is frozen for that pin key. When it is, the task SHALL skip that file for that key. After `cue mod tidy`, every frozen key in that module SHALL be byte-unchanged; otherwise the task SHALL exit 1 and name the file and the key.

#### Scenario: Frozen fixture left alone
- **WHEN** `.cascade-frozen` freezes `test/fixtures/modules/hello/cue.mod/module.cue` for the catalog and core, and both are behind
- **THEN** that file is byte-unchanged, the other pins still move, and the task exits 0

#### Scenario: Frozen catalog keeps its core
- **WHEN** `.cascade-frozen` freezes a fixture module file for the catalog only, its catalog is behind the newest, and its core is what its own catalog pins
- **THEN** core is chosen from the file's current catalog, so that file is byte-unchanged

#### Scenario: MVS raises a frozen key
- **WHEN** `cue mod get` of an unfrozen key raises a frozen key in the same module
- **THEN** the task exits 1 and names the module file and the frozen key

### Requirement: deps:cascade honours a hold on core
In `task deps:cascade`, a hold on `opmodel.dev/core@v2` whose `max` is below the core the target catalog pins SHALL keep both the catalog and core at their current versions, with a warning.

#### Scenario: Core hold holds the catalog
- **WHEN** an in-date hold caps core below the core the newest catalog pins
- **THEN** neither the catalog nor core moves, and a warning names the catalog, the core it needs, and the hold

### Requirement: Fixture versions advance once per PR
For each fixture module under `test/fixtures/modules/` and for the provider catalog fixture, `task deps:cascade` SHALL set the declared version in `identity/identity.cue` as follows. Let `B` be the declared version at the merge-base of `CASCADE_BASE` (default `origin/main`) and `HEAD`:
- `B` itself, when the fixture did not change against that merge-base;
- the next patch of `B`, when it did change and `B` is published;
- `B` itself, when it changed and `B` is not yet published.

#### Scenario: Second run does not advance again
- **WHEN** the task has advanced a fixture to the next patch, the result is committed, and the task runs again with the same `CASCADE_BASE`
- **THEN** it exits 3 and the fixture's declared version is unchanged

#### Scenario: Pending version not bumped twice
- **WHEN** a fixture changed against the merge-base and its merge-base version is not yet published
- **THEN** the declared version stays at the merge-base version

### Requirement: Fixture versions are written with the opm CLI
`task deps:cascade` SHALL write a fixture's declared version with `opm module version set` or `opm catalog version set`, using an opm CLI the task installed from the unmodified `.opm-cli-version` under the git directory, and only when it differs.

#### Scenario: Version written only when it differs
- **WHEN** a fixture's declared version already equals the version the task computes
- **THEN** the task does not write it

### Requirement: Fixture consumers follow the advanced version
In the same run that sets a fixture module's declared version, `task deps:cascade` SHALL make each of the module's consumers follow that version:
- the modulepackage `cue.mod/module.cue` (the fixture's `v:`, catalog and core, as a text edit with no `tidy`);
- its `moduleinstance.yaml`;
- `config/samples/opmodel.dev_v1alpha1_moduleinstance.yaml` for the module it instantiates.

#### Scenario: Consumers name the advanced version
- **WHEN** the task advances `test/fixtures/modules/hello` from `0.0.12` to `0.0.13`
- **THEN** `test/fixtures/modulepackages/hello/cue.mod/module.cue` pins `v0.0.13`, and `test/fixtures/modules/hello/moduleinstance.yaml` and the sample ModuleInstance name `v0.0.13`

### Requirement: deps:cascade never modifies release, workflow or shared files
`task deps:cascade` SHALL NOT modify:
- `.cascade-frozen` or `.cascade-hold`;
- `.release-please-manifest.json`, `release-please-config.json` or `CHANGELOG.md`;
- anything under `.github/`;
- `.opm-docs-version`, `docs-kit.cue` or any docs-kit ref;
- `hack/fixtures.sh` or `.tasks/examples.yaml`;
- any CUE `language.version`;
- the Jellyfin sample or `ocirepository.yaml`.

#### Scenario: Release files untouched
- **WHEN** the task moves every pin
- **THEN** `git diff --name-only` names no path under `.github/`, no release-please file and no `.cascade-*` file

### Requirement: deps:cascade leaves the controller module to its own task
`task deps:cascade` SHALL NOT modify anything under the controller module's directory `modules/opm_controller/`, whose pins `deps:cascade:module` moves in a PR of its own.

#### Scenario: The controller module is left to its own task
- **WHEN** the newest published catalog is newer than the one the controller module pins
- **THEN** `task deps:cascade` leaves every file under `modules/opm_controller/` unchanged

### Requirement: deps:cascade only reads from registries
`task deps:cascade` SHALL only read from registries and SHALL never publish or seed.

#### Scenario: Registries are only read
- **WHEN** `task deps:cascade` runs
- **THEN** it reads from registries and publishes or seeds nothing

### Requirement: deps:cascade warns when an upstream CUE language is newer
When an upstream's `language.version` is newer than the `CUE_VERSION` in `.github/workflows/test.yml`, `task deps:cascade` SHALL warn.

#### Scenario: Newer CUE language warned
- **WHEN** the target catalog's `language.version` is newer than `.github/workflows/test.yml`'s `CUE_VERSION`
- **THEN** the task still moves the pin and writes a warning naming both versions

### Requirement: deps:cascade warns about install page drift and leaves the page unchanged
When the example output in `docs/site/start/install-the-controller.md` names a catalog other than the sample Platform's, or a core other than `test/fixtures/modules/hello`'s, `task deps:cascade` SHALL warn and SHALL NOT edit that page.

#### Scenario: Install page example output drifted
- **WHEN** `docs/site/start/install-the-controller.md` prints a catalog or core version the tree no longer pins
- **THEN** the task writes a warning naming the page and both versions, and leaves the page unchanged

### Requirement: deps:cascade:test reports each scenario against a stub resolver
`task deps:cascade:test` SHALL run `deps:cascade` in throwaway copies of the tree against the contract's stub resolver, and SHALL report `PASS` or `FAIL` per scenario. It SHALL exit 0 only when every scenario passes. It SHALL fail when the stub's `sha256sum` differs from the contract's checksum.

#### Scenario: Stub drift caught
- **WHEN** `.tasks/cascade/testdata/stub-resolve.sh` differs by one byte from the contract's text
- **THEN** `task -x deps:cascade:test` fails naming that file

### Requirement: The offline cascade test set needs no network and runs in the Lint job
With `CASCADE_TEST_SET=offline`, `task deps:cascade:test` SHALL need no GHCR or module proxy access beyond a warm Go module cache, and it SHALL run in the `Lint` job of `.github/workflows/lint.yml`, the check workspace RELEASING.md makes required. It SHALL NOT need the real resolver or a `.github` checkout, except for the title and body check.

#### Scenario: Offline set in PR CI
- **WHEN** a pull request runs the `Lint` job
- **THEN** the job runs `task -x deps:cascade:test` with `CASCADE_TEST_SET=offline`, and the job fails if a scenario fails

### Requirement: The full cascade test set runs apart from required CI
The full set of `task deps:cascade:test` SHALL run in a separate, non-required workflow on changes to the cascade files, by hand, and weekly. That set SHALL include older pins restored to the tree's versions with the expected version-advance diff, frozen pins, and the title and body checked against the real resolver.

#### Scenario: Older pins produce the expected diff
- **WHEN** the full set lowers every pin to its `older.tsv` version and runs the task
- **THEN** the task exits 0, the tree differs from the original only in the fixture version-advance paths, and a second run exits 3

### Requirement: Title and body come from the shared resolver
`task deps:cascade:title` and `task deps:cascade:body` SHALL call the resolver's `title` and `body` subcommands with `.tasks/cascade/classes` and `.tasks/cascade/pins.sh`.

#### Scenario: Library and catalog moved
- **WHEN** the diff against `origin/main` moves library to `v1.0.0-beta.3` and the catalog to `v4.5.1`
- **THEN** `task -x deps:cascade:title` prints `fix(deps): bump library to v1.0.0-beta.3 and opm catalog to v4.5.1`

### Requirement: The classes file holds the contract's path classes
`.tasks/cascade/classes` SHALL classify `.opm-cli-version` as release-tool and `config/samples/`, `test/`, any `testdata/` directory and `*_test.go` as test, the contract's opm-controller block and nothing more. Any other path, including `go.mod` and `go.sum`, SHALL be shipped.

#### Scenario: Only fixtures moved
- **WHEN** the diff changes only paths under `test/` and `config/samples/`
- **THEN** the title type is `test(fixtures)`

#### Scenario: Only the sample moved, no page regenerated
- **WHEN** the task moves only the sample Platform's catalog `version:`, so the diff changes only paths under `config/samples/` and leaves `docs/site/` untouched
- **THEN** the title type is `test(fixtures)`

#### Scenario: Only the opm CLI moved
- **WHEN** the diff changes only `.opm-cli-version`
- **THEN** the title type is `ci(deps)`

### Requirement: pins.sh prints one row per logical pin
`.tasks/cascade/pins.sh <ref>` SHALL print one row per logical pin for the working tree (`WORKTREE`) or a git ref, as `<pin-key>`, `<display>`, `<class>`, `<v-prefixed version>`, `<labels>` separated by tabs. The rows are library (shipped), the opm catalog (test, from the sample Platform), core (test, from `test/fixtures/modules/hello`) and the opm CLI (release-tool).

#### Scenario: Pin report agrees between tree and HEAD
- **WHEN** the tree is clean
- **THEN** `pins.sh WORKTREE` and `pins.sh HEAD` print the same rows

### Requirement: deps:cascade:module moves the module's catalog and core
`task deps:cascade:module` SHALL move the controller module's pins in `modules/opm_controller/cue.mod/module.cue` and nothing else, with the rules `deps:cascade` applies to a fixture module: the opm catalog `opmodel.dev/catalogs/opm@v4` to the newest published catalog `K` when it is below `K`, and core `opmodel.dev/core@v2` to the core the module's catalog pins after the move when that is greater. The module's pins are shipped pins: every user who installs the module receives them.

#### Scenario: Module pins move, image and version stay
- **WHEN** the newest published catalog is newer than the controller module's
- **THEN** the module's `cue.mod/module.cue` names the new catalog and the core it pins, `git diff --name-only` names only that file and the module's `cue.mod` files `tidy` writes, and the module's image reference, `identity.Version`, generated files, `CHANGELOG.md` and `RELEASE` are unchanged

### Requirement: deps:cascade:module honours holds and frozen paths
`task deps:cascade:module` SHALL apply holds and frozen paths under the same pin keys, `opmodel.dev/catalogs/opm@v4` and `opmodel.dev/core@v2`, so a hold on either holds the module too, and a `.cascade-frozen` entry for `modules/opm_controller` freezes only the module.

#### Scenario: Catalog hold holds the module
- **WHEN** an in-date hold on `opmodel.dev/catalogs/opm@v4` caps the catalog below the newest
- **THEN** the module's catalog moves no further than the hold

### Requirement: deps:cascade:module reports by the exit codes of deps:cascade
`task deps:cascade:module` SHALL report by the same exit codes as `deps:cascade` (0 changed, 3 nothing to do).

#### Scenario: Nothing to move
- **WHEN** the module already pins the newest catalog and the core it pins
- **THEN** the task exits 3 and the tree is unchanged

### Requirement: deps:cascade:module edits only the module's pins
`task deps:cascade:module` SHALL NOT edit any file outside `modules/opm_controller/`, nor the module's `controller/controller.cue`, its `identity.Version`, its generated CRD and RBAC files, its `CHANGELOG.md` or its `RELEASE`: the image moves only through the image-bump PR, the version only through the identity advance, the generated files only through the generator, and the release files only through release-please.

#### Scenario: Files the module task never edits
- **WHEN** `task deps:cascade:module` runs and a newer catalog exists
- **THEN** it edits no file outside `modules/opm_controller/` and leaves the module's `controller/controller.cue`, `identity.Version`, generated CRD and RBAC files, `CHANGELOG.md` and `RELEASE` unchanged

### Requirement: The offline test set covers deps:cascade:module
The offline test set of `task deps:cascade:test` SHALL cover `task deps:cascade:module`.

#### Scenario: Offline set covers the module task
- **WHEN** `task deps:cascade:test` runs with `CASCADE_TEST_SET=offline`
- **THEN** its scenarios include `task deps:cascade:module`

### Requirement: module-deps.yml opens one pull request for the module's pin move
`.github/workflows/module-deps.yml` SHALL run `task -x deps:cascade:module` on `repository_dispatch` of type `upstream-released` and on `workflow_dispatch`. Unless the repository variable `CASCADE_DRY_RUN` is not exactly `false`, when the task changed the tree, it SHALL push the branch `module/deps` with the release App's token and open or update one pull request from it, so the PR's CI runs.

#### Scenario: Catalog published
- **WHEN** this repository receives `upstream-released` for catalog `opm-v4.6.0`, `CASCADE_DRY_RUN` is `false`, and the module pins `v4.5.2`
- **THEN** a PR from `module/deps` moves the module to `v4.6.0` and the core it pins, titled `fix(deps): bump the controller module's opm catalog to v4.6.0 and core to <v>`, and merging it proposes a module release and no controller release

### Requirement: The module/deps branch is rebuilt only while it holds bot commits
While the branch `module/deps` carries only bot commits, a later run of `module-deps.yml` SHALL rebuild it from `main`. Once a human commit is on it, a run SHALL add its own commit and never rewrite the branch.

#### Scenario: Bot-only branch is rebuilt
- **WHEN** `module/deps` carries only bot commits and `module-deps.yml` runs again
- **THEN** the run rebuilds the branch from `main`

#### Scenario: Human commit on the branch
- **WHEN** a human commit is on `module/deps` and `module-deps.yml` runs again
- **THEN** the run adds its own commit and does not rewrite the branch

### Requirement: The module pull request title and body come from the resolver
The title and body of the `module-deps.yml` pull request SHALL come from the resolver's `title` and `body` with `.tasks/cascade/module-classes`, which classifies every path under `modules/opm_controller/` as shipped, and `.tasks/cascade/module-pins.sh`, which prints two rows keyed `opmodel.dev/catalogs/opm@v4` and `opmodel.dev/core@v2`, displayed as the controller module's opm catalog and core. Because that report is its own, each key appears once.

#### Scenario: Pin report keys are unique
- **WHEN** the resolver reads `.tasks/cascade/module-pins.sh WORKTREE`
- **THEN** it accepts the report, and `pins.sh WORKTREE` for the repository's cascade lists no row for the module

### Requirement: The module pull request changes only the module directory
The pull request of `module-deps.yml` SHALL change no file outside `modules/opm_controller/`, so its squash commit proposes a module release and no controller release.

#### Scenario: Module pull request touches only the module
- **WHEN** `module-deps.yml` opens its pull request
- **THEN** the pull request changes no file outside `modules/opm_controller/`, and its squash commit proposes a module release and no controller release

### Requirement: module-deps.yml writes only a summary in a dry run
While the repository variable `CASCADE_DRY_RUN` is not exactly `false`, a run of `module-deps.yml` SHALL write the diff to the job summary and push nothing.

#### Scenario: Dry run
- **WHEN** `CASCADE_DRY_RUN` is unset
- **THEN** the run writes the diff to the job summary and `module/deps` is not created or moved
