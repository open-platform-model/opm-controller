## MODIFIED Requirements

### Requirement: deps:cascade moves upstream pins and reports by exit code
The repo SHALL provide `task deps:cascade`, which moves opm-controller's upstream pins to the newest published versions in the working tree only. It SHALL NOT commit, push or create a branch. Run as `task -x deps:cascade`, it SHALL exit 0 when the working tree changed, 3 when there was nothing to do, and any other code on error, and it SHALL never turn a failure into 0 or 3. Every version decision SHALL come from the shared cascade resolver of `open-platform-model/.github` (workspace RELEASING.md, section "The cascade"). The resolver is found through the absolute path in `CASCADE_RESOLVER`, or otherwise at `../.github/.github/scripts/cascade/cascade-resolve.sh` beside the main checkout. When the resolver is missing, the task SHALL fail with a message naming both options.

#### Scenario: Up-to-date tree
- **WHEN** every pin already equals the newest published version the resolver reports
- **THEN** `task -x deps:cascade` exits 3 and `git status --porcelain` is empty

#### Scenario: A pin is behind
- **WHEN** the resolver reports a newer published library than `go.mod` pins
- **THEN** `task -x deps:cascade` exits 0 and `go.mod` and `go.sum` name the newer library

#### Scenario: Resolver not found
- **WHEN** `CASCADE_RESOLVER` is unset and no `.github` checkout sits beside the repo
- **THEN** the task fails before running and says to check out `open-platform-model/.github` or set `CASCADE_RESOLVER`

#### Scenario: A predicate's "no" is not "nothing to do"
- **WHEN** a resolver call other than `newest`, such as `pin-of` for the target catalog, exits 3
- **THEN** the task exits with a code other than 0 or 3 and `git status --porcelain` is empty

#### Scenario: Relative resolver path
- **WHEN** `CASCADE_RESOLVER` is set to a relative path
- **THEN** the task fails and says the path must be absolute

### Requirement: deps:cascade moves library, the catalog and core, and the opm CLI
The task SHALL move exactly these pins:
- **library.** It runs `go get github.com/open-platform-model/library@<exact version>` and then `go mod tidy`, only when library moved.
- **The opm catalog.** It moves `opmodel.dev/catalogs/opm@v4` to one target `K` in every file that pins it and is below `K`:
  - the sample Platform's `version:` (stored bare);
  - `CatalogVersion()` in `test/fixtures/catalog.go` (stored bare);
  - the four `test/fixtures/modules/*/cue.mod/module.cue` files.

  A file whose catalog is above `K` SHALL NOT be lowered.
- **Core.** It moves core in those modules and in `test/fixtures/catalogs/provider/cue.mod/module.cue` to the version that the file's catalog after the move pins (the higher of the file's catalog and `K`; `K` for the provider fixture), and only when that is greater than the file's own core. Core SHALL never come from the newest published core directly.
- **The opm CLI.** It writes `.opm-cli-version` last.
- **No regenerated reference.** It SHALL NOT edit `docs/site/reference/controller/_index.md`: the docs bundle generates the resource reference from `config/samples` when it is built.

`cue mod get` SHALL name only `opmodel.dev/*` and `testing.opmodel.dev/*` modules, each with an exact version, and SHALL run, followed by one `cue mod tidy`, only in a module where a pin moved. Third-party pins SHALL never be named. A third-party pin that `tidy` raises, adds or removes, a change to `go.mod`'s `go` or `toolchain` directive, and a dep a fixture module gains that its modulepackage lacks SHALL each be reported as a warning. The task SHALL never edit an import path or a `@vN` key; a new major SHALL appear only as the resolver's warning.

#### Scenario: Catalog and core move as a consistent set
- **WHEN** the newest published catalog is newer than the sample's, and that catalog pins a core newer than the fixtures'
- **THEN** the sample, `catalog.go` and the four fixture modules name the new catalog, and the four fixture modules and the provider catalog fixture name the core that catalog pins

#### Scenario: Core already ahead of the catalog's core
- **WHEN** a fixture pins a core newer than the core catalog `K` pins
- **THEN** that fixture's core is left as it is and a warning names both versions

#### Scenario: Newer core published but not pinned by the catalog
- **WHEN** core `v2.0.0-beta.2` is published and the newest catalog pins core `v2.0.0-beta.1`
- **THEN** the fixtures' core stays at `v2.0.0-beta.1`

#### Scenario: Sample Platform moved
- **WHEN** the task moves the catalog `version:` in the sample Platform
- **THEN** `docs/site/reference/controller/_index.md` is unchanged, and `task docs:bundle` on the resulting tree shows the moved version in the example of `reference/controller/platform.md`

#### Scenario: Tidy adds a dependency
- **WHEN** `go mod tidy` or `cue mod tidy` adds or removes a third-party dependency, or `go get` raises the `go` or `toolchain` directive
- **THEN** the task writes a warning naming the dependency or the directives, and a fixture module's new dep that its modulepackage lacks is warned about too

#### Scenario: Third-party pins untouched
- **WHEN** the task runs `cue mod get` in a fixture module
- **THEN** the command names no module outside `opmodel.dev/` and `testing.opmodel.dev/`

### Requirement: Title and body come from the shared resolver and the contract's path classes
`task deps:cascade:title` and `task deps:cascade:body` SHALL call the resolver's `title` and `body` subcommands with `.tasks/cascade/classes` and `.tasks/cascade/pins.sh`.
- **`classes`** SHALL classify `.opm-cli-version` as release-tool and `config/samples/`, `test/`, any `testdata/` directory and `*_test.go` as test, the contract's opm-controller block and nothing more. Any other path, including `go.mod` and `go.sum`, SHALL be shipped.
- **`pins.sh <ref>`** SHALL print one row per logical pin for the working tree (`WORKTREE`) or a git ref, as `<pin-key>`, `<display>`, `<class>`, `<v-prefixed version>`, `<labels>` separated by tabs. The rows are library (shipped), the opm catalog (test, from the sample Platform), core (test, from `test/fixtures/modules/hello`) and the opm CLI (release-tool).

#### Scenario: Library and catalog moved
- **WHEN** the diff against `origin/main` moves library to `v1.0.0-beta.3` and the catalog to `v4.5.1`
- **THEN** `task -x deps:cascade:title` prints `fix(deps): bump library to v1.0.0-beta.3 and opm catalog to v4.5.1`

#### Scenario: Only fixtures moved
- **WHEN** the diff changes only paths under `test/` and `config/samples/`
- **THEN** the title type is `test(fixtures)`

#### Scenario: Only the sample moved, no page regenerated
- **WHEN** the task moves only the sample Platform's catalog `version:`, so the diff changes only paths under `config/samples/` and leaves `docs/site/` untouched
- **THEN** the title type is `test(fixtures)`

#### Scenario: Only the opm CLI moved
- **WHEN** the diff changes only `.opm-cli-version`
- **THEN** the title type is `ci(deps)`

#### Scenario: Pin report agrees between tree and HEAD
- **WHEN** the tree is clean
- **THEN** `pins.sh WORKTREE` and `pins.sh HEAD` print the same rows

## REMOVED Requirements

### Requirement: deps:cascade never touches release, workflow or shared files

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-019); this requirement is restated under the controller's name.

**Migration**: See the requirement "deps:cascade leaves release, workflow and shared files alone" in this capability.

### Requirement: deps:cascade:module moves only the operator module's catalog and core

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-019); this requirement is restated under the controller's name.

**Migration**: See the requirement "deps:cascade:module moves only the controller module's catalog and core" in this capability.

### Requirement: The operator module's pin move opens its own pull request

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-019); this requirement is restated under the controller's name.

**Migration**: See the requirement "The controller module's pin move opens its own pull request" in this capability.

## ADDED Requirements

### Requirement: deps:cascade leaves release, workflow and shared files alone
The task SHALL NOT modify:
- `.cascade-frozen` or `.cascade-hold`;
- `.release-please-manifest.json`, `release-please-config.json` or `CHANGELOG.md`;
- anything under `.github/`;
- `.opm-docs-version`, `docs-kit.cue` or any docs-kit ref;
- `hack/fixtures.sh` or `.tasks/examples.yaml`;
- any CUE `language.version`;
- the Jellyfin sample or `ocirepository.yaml`;
- anything under the controller module's directory `modules/opm_controller/`, whose pins `deps:cascade:module` moves in a PR of its own.

It SHALL only read from registries and SHALL never publish or seed. When an upstream's `language.version` is newer than the `CUE_VERSION` in `.github/workflows/test.yml`, it SHALL warn. When the example output in `docs/site/start/install-the-controller.md` names a catalog other than the sample Platform's, or a core other than `test/fixtures/modules/hello`'s, it SHALL warn and SHALL NOT edit that page.

#### Scenario: Release files untouched
- **WHEN** the task moves every pin
- **THEN** `git diff --name-only` names no path under `.github/`, no release-please file and no `.cascade-*` file

#### Scenario: The controller module is left to its own task
- **WHEN** the newest published catalog is newer than the one the controller module pins
- **THEN** `task deps:cascade` leaves every file under `modules/opm_controller/` unchanged

#### Scenario: Install page example output drifted
- **WHEN** `docs/site/start/install-the-controller.md` prints a catalog or core version the tree no longer pins
- **THEN** the task writes a warning naming the page and both versions, and leaves the page unchanged

#### Scenario: Newer CUE language warned
- **WHEN** the target catalog's `language.version` is newer than `.github/workflows/test.yml`'s `CUE_VERSION`
- **THEN** the task still moves the pin and writes a warning naming both versions

### Requirement: deps:cascade:module moves only the controller module's catalog and core
`task deps:cascade:module` SHALL move the controller module's pins in `modules/opm_controller/cue.mod/module.cue` and nothing else, with the rules `deps:cascade` applies to a fixture module: the opm catalog `opmodel.dev/catalogs/opm@v4` to the newest published catalog `K` when it is below `K`, and core `opmodel.dev/core@v2` to the core the module's catalog pins after the move when that is greater. The module's pins are shipped pins: every user who installs the module receives them. Holds and frozen paths SHALL apply under the same pin keys, `opmodel.dev/catalogs/opm@v4` and `opmodel.dev/core@v2`, so a hold on either holds the module too, and a `.cascade-frozen` entry for `modules/opm_controller` freezes only the module. It SHALL report by the same exit codes as `deps:cascade` (0 changed, 3 nothing to do). It SHALL NOT edit any file outside `modules/opm_controller/`, nor the module's `controller/controller.cue`, its `identity.Version`, its generated CRD and RBAC files, its `CHANGELOG.md` or its `RELEASE`: the image moves only through the image-bump PR, the version only through the identity advance, the generated files only through the generator, and the release files only through release-please. The offline test set SHALL cover it.

#### Scenario: Module pins move, image and version stay
- **WHEN** the newest published catalog is newer than the controller module's
- **THEN** the module's `cue.mod/module.cue` names the new catalog and the core it pins, `git diff --name-only` names only that file and the module's `cue.mod` files `tidy` writes, and the module's image reference, `identity.Version`, generated files, `CHANGELOG.md` and `RELEASE` are unchanged

#### Scenario: Catalog hold holds the module
- **WHEN** an in-date hold on `opmodel.dev/catalogs/opm@v4` caps the catalog below the newest
- **THEN** the module's catalog moves no further than the hold

#### Scenario: Nothing to move
- **WHEN** the module already pins the newest catalog and the core it pins
- **THEN** the task exits 3 and the tree is unchanged

### Requirement: The controller module's pin move opens its own pull request
`.github/workflows/module-deps.yml` SHALL run `task -x deps:cascade:module` on `repository_dispatch` of type `upstream-released` and on `workflow_dispatch`. When the task changed the tree, it SHALL push the branch `module/deps` with the release App's token and open or update one pull request from it, so the PR's CI runs; while the branch carries only bot commits a later run SHALL rebuild it from `main`, and once a human commit is on it a run SHALL add its own commit and never rewrite the branch. Its title and body SHALL come from the resolver's `title` and `body` with `.tasks/cascade/module-classes`, which classifies every path under `modules/opm_controller/` as shipped, and `.tasks/cascade/module-pins.sh`, which prints two rows keyed `opmodel.dev/catalogs/opm@v4` and `opmodel.dev/core@v2`, displayed as the controller module's opm catalog and core. Because that report is its own, each key appears once. The PR SHALL change no file outside `modules/opm_controller/`, so its squash commit proposes a module release and no controller release. While the repository variable `CASCADE_DRY_RUN` is not exactly `false`, the run SHALL write the diff to the job summary and push nothing.

#### Scenario: Catalog published
- **WHEN** this repository receives `upstream-released` for catalog `opm-v4.6.0`, `CASCADE_DRY_RUN` is `false`, and the module pins `v4.5.2`
- **THEN** a PR from `module/deps` moves the module to `v4.6.0` and the core it pins, titled `fix(deps): bump the controller module's opm catalog to v4.6.0 and core to <v>`, and merging it proposes a module release and no controller release

#### Scenario: Dry run
- **WHEN** `CASCADE_DRY_RUN` is unset
- **THEN** the run writes the diff to the job summary and `module/deps` is not created or moved

#### Scenario: Pin report keys are unique
- **WHEN** the resolver reads `.tasks/cascade/module-pins.sh WORKTREE`
- **THEN** it accepts the report, and `pins.sh WORKTREE` for the repository's cascade lists no row for the module
