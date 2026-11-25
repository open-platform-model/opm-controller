## RENAMED Requirements

- FROM: `### Requirement: The suite asserts the Platform registry entry under a new package identity`
- TO: `### Requirement: The suite asserts the Platform registry entry under the bare-version package identity`
- FROM: `### Requirement: The package identity shows the platform was regenerated for the edited claim`
- TO: `### Requirement: The package identity does not move on a spelling-only edit`

## MODIFIED Requirements

### Requirement: The suite asserts the Platform registry entry under the bare-version package identity
As in requirement "The claim is re-judged at its new generation": Only after that verdict SHALL the suite assert that the Platform's `status.registry` entry for the backup catalog names the same build in its bare spelling, under the same package identity as the bare-version claim, stable on a second read, and that the Platform reports `Ready=True` with reason `Generated`.

#### Scenario: A v-prefixed claim is accepted and built

- **WHEN** the provider instance is suspended and the live claim's `spec.version` is set to the same build with a `v` prefix
- **THEN** the claim's `status.observedGeneration` SHALL reach its new `metadata.generation` with `accepted: true` and `active: true`
- **AND** the claim's Ready condition SHALL be `True` with reason `Accepted` and a message containing `at v` followed by the build
- **AND**, once that verdict is recorded, the Platform's `status.registry` SHALL list the backup catalog at the same build without the `v`, under the package identity of the bare-version claim, stable on a second read
- **AND** the Platform SHALL report `Ready=True` with reason `Generated`

### Requirement: The package identity does not move on a spelling-only edit
As in requirement "The suite asserts the Platform registry entry under the bare-version package identity": The package identity SHALL be a function of the claims' coordinates in their bare spelling. A move of the identity shows that the platform was regenerated, and a spelling-only edit of a claim's version SHALL NOT move it.

#### Scenario: A moved package identity shows the regeneration
- **WHEN** the package identity moves
- **THEN** that move shows the platform was regenerated, because the package identity is a function of the claims' bare coordinates

#### Scenario: A spelling-only edit keeps the package identity
- **WHEN** the live claim's `spec.version` changes from the bare build to the same build with a `v` prefix
- **THEN** the package identity equals the one of the bare-version claim, so the platform is not regenerated
