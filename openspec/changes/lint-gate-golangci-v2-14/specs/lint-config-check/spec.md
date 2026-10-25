## MODIFIED Requirements

### Requirement: Linter configuration check needs no network
The linter configuration check (`task dev:lint:config`) SHALL validate `.golangci.yml` with the linter's own `config verify` command, which validates against the JSON schema embedded in the linter binary, and the check itself SHALL NOT make a network request. The check SHALL run the linter with its proxy variables set to a closed local port. Installing the linter, which the task does before the check, is outside this requirement. The `Lint` job SHALL run this check on every run.

#### Scenario: Valid configuration without a network
- **WHEN** the check script that `task dev:lint:config` runs is started in a network namespace with no route, with the linter already installed
- **THEN** it exits 0 and reports the version and that no network was used

#### Scenario: Invalid configuration
- **WHEN** `.golangci.yml` holds a key the schema does not allow
- **THEN** the check exits non-zero and its output names the key

#### Scenario: Linter reaches for the network
- **WHEN** the linter's `config verify` tries to open a connection
- **THEN** the connection goes to a closed local port and the check fails at once

### Requirement: The check has a scenario test in the Lint job
A scenario test SHALL run the check on a copy of the tree (pass) and on one copy per defect (refusal with the expected message). The test itself SHALL make no network request. It SHALL run in the `Lint` job, with the linter that the configuration check step installed, and through `task dev:lint:config:test`, which installs the linter first.

#### Scenario: Schema is the one of the named line
- **WHEN** `.golangci.yml` uses a value that the pinned linter line added
- **THEN** the check exits 0

#### Scenario: Scenario test passes on the tree
- **WHEN** `task dev:lint:config:test` runs on an unmodified tree
- **THEN** every scenario reports `ok` and the task exits 0

## ADDED Requirements

### Requirement: Version and use agree
The check SHALL refuse, with one line per problem and a non-zero exit, a tree where any of these does not hold: the version file holds `vX.Y.Z`; the installed linter is the version the file names; `lint.yml` runs the check through `task dev:lint:config`; no workflow, task file or `Makefile` uses the linter's GitHub action. The repository SHALL hold no committed linter schema and no schema checksum file.

#### Scenario: Workflow stops running the check
- **WHEN** `lint.yml` has no step that runs `task dev:lint:config`
- **THEN** the check reports it

#### Scenario: Installed linter is another version
- **WHEN** the linter binary reports another version than the file names
- **THEN** the check reports both versions

#### Scenario: A committed schema or checksum file returns
- **WHEN** `.github/golangci-lint/` holds a schema or a checksum file
- **THEN** the check names the file

#### Scenario: The linter action returns
- **WHEN** a workflow, a task file or `Makefile` uses `golangci/golangci-lint-action`
- **THEN** the check names the file and line

## REMOVED Requirements

### Requirement: Version, schema, checksum and use agree
**Reason**: The linter embeds the schema of its own version, so no schema file or checksum is committed. The version-to-schema agreement cannot be wrong.
**Migration**: Replaced by "Version and use agree". The workflow, installed-version and action rules carry over under the same scenario names where the rule remains. The schema, checksum and `config verify` without `--schema` scenarios end: a bare `config verify` is now offline, and a linter old enough to download the schema fails the closed-port guard of "Linter configuration check needs no network".
