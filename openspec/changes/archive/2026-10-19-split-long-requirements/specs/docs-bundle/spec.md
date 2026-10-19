## REMOVED Requirements

### Requirement: The operator publishes one docs bundle
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "The opm-operator docs project is placed in the docs tree and versioned from tags", "The opm-operator docs project builds from a markdown source and a crd source", "The opm-operator docs project sets the section weight, the reconciledBy entries and plain-text citations" and "The bundle carries the authored pages and the generated reference and is built by task docs:bundle".

## ADDED Requirements

### Requirement: The opm-operator docs project is placed in the docs tree and versioned from tags
The repository SHALL declare one docs-kit project, `opm-operator`, in `docs-kit.cue`: placed in a site version's `/docs/` tree, owning `reference/operator/`, versioned from tags with the prefix `v`.

#### Scenario: The project is placed in the docs tree
- **WHEN** `docs-kit.cue` is read
- **THEN** it declares the project `opm-operator`, placed in a site version's `/docs/` tree, owning `reference/operator/`, versioned from tags with the prefix `v`

### Requirement: The opm-operator docs project builds from a markdown source and a crd source
The repository SHALL declare one docs-kit project, `opm-operator`, in `docs-kit.cue`: built from a `markdown` source over `docs/site` and a `crd` source over `config/crd/bases` with samples from `config/samples` (picked by kubebuilder file name, those referencing `testing.opmodel.dev` hidden, the two kubebuilder scaffold labels stripped) in the section layout at `reference/operator/`.

#### Scenario: The test fixture is not shown as an example

- **WHEN** the ModuleInstance sample `config/samples/opmodel.dev_v1alpha1_moduleinstance.yaml` references `testing.opmodel.dev`
- **THEN** the ModuleInstance page has no Example part, and `opmodel.dev_v1alpha1_moduleinstance_jellyfin.yaml` is not read

### Requirement: The opm-operator docs project sets the section weight, the reconciledBy entries and plain-text citations
The repository SHALL declare one docs-kit project, `opm-operator`, in `docs-kit.cue`: the section weight 7, one `reconciledBy` entry per kind, and enhancement citations left as plain text, not links.

#### Scenario: An enhancement citation stays plain text

- **WHEN** a CRD description cites `(enhancement 0015)`
- **THEN** the generated page shows `(enhancement 0015)` as plain text, with no link

### Requirement: The bundle carries the authored pages and the generated reference and is built by task docs:bundle
As in requirement "The opm-operator docs project is placed in the docs tree and versioned from tags": The bundle SHALL carry the authored pages under `docs/site/` and the generated resource reference together (docs-kit DESIGN decision 20). `task docs:bundle` SHALL build it into `out/opm-operator/` and `task docs:bundle:check` SHALL build and lint it without publishing.

#### Scenario: The bundle holds both kinds of page

- **WHEN** `task docs:bundle` runs on a clean checkout
- **THEN** `out/opm-operator/content/` holds the authored pages under `start/`, `operating/` and `diagnostics/`, and under `reference/operator/` the index `_index.md` and one page for each of the four kinds

#### Scenario: A CRD that is not a CustomResourceDefinition is refused

- **WHEN** a file in `config/crd/bases` holds a document of another kind
- **THEN** `task docs:bundle:check` fails, naming the file
