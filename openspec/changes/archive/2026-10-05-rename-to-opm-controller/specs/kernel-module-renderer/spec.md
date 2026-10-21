## MODIFIED Requirements

### Requirement: Render a ModuleRelease through the kernel

`KernelModuleRenderer` SHALL acquire the module from the registry, synthesize the instance (source-carrying) with the supplied values, and render it through the single-build render with the leased platform record and the record's skew policy. Every kernel call SHALL share nothing: acquisition, synthesis and the build each evaluate in a context of their own, so renders of different objects overlap with no gate.

Values SHALL reach synthesis as a stack of values sources, each carrying an origin that identifies where the controller read it from, so a values error names that origin rather than an anonymous filename. No acquisition or synthesis call SHALL pass a per-call registry or load-options argument: the registry mapping is the one the shared Kernel was constructed with.

#### Scenario: Renders resources from a generated platform

- **WHEN** `RenderModule` is called for a resolvable module while a generated platform is recorded
- **THEN** the result carries the rendered resources and any warnings the render reported, and no inventory entries: the reconciler builds those from its one export of the resources

#### Scenario: Values are applied when supplied

- **WHEN** non-nil `RawValues` are passed
- **THEN** they are supplied to instance synthesis as a values source whose origin names the CR field they came from
- **AND** when no values are supplied the module's `#config` defaults apply

#### Scenario: A values error names its origin

- **WHEN** the supplied raw values violate the module's `#config` schema
- **THEN** the render fails with an error naming the origin the controller gave that source, not an anonymous filename

## REMOVED Requirements

### Requirement: Adapt compiled output to operator resources

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-021); this requirement is restated under the controller's name.

**Migration**: See the requirement "Adapt compiled output to controller resources" in this capability.

### Requirement: The render reports the instance's contract demand from the kernel

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-021); this requirement is restated under the controller's name.

**Migration**: See the requirement "The render reads the instance's contract demand from the kernel" in this capability.

## ADDED Requirements

### Requirement: Adapt compiled output to controller resources

The renderer SHALL adapt the single-build render's compiled objects — the kernel's own compiled-output type, declared beside the render verb — to the library's Kubernetes object resources (`opm/k8s/object`) exactly as before, after first refusing a render whose compiled objects share one Kubernetes apply identity (apiVersion, kind, namespace and name). The refusal SHALL use the library's duplicate-identity check and SHALL carry its error unchanged, naming each shared identity once and every component and transformer that produced it; when it fires, the renderer SHALL build no resource, so nothing of that render can reach apply (enhancement 0015 ("One registration per module"): a second registration in one module is refused before apply naming both carrying components, and enhancement 0015 ("The rendered CR's name is instance-derived: dot-joined `namespace.name`") gives both the same instance-derived name).

The render SHALL report no message strings. The renderer SHALL compose the result's warnings itself from the render's advisory diagnostic rows: the unhandled-trait table, and the resolved-versions rows marked newer. Each composed warning SHALL name the same facts as before — for skew, the OPM-namespace path, the version the module requires and the version the platform carries; for an unhandled trait, the component and the trait.

#### Scenario: Compiled item maps to a resource

- **WHEN** a library `Compiled` with a value and provenance is adapted
- **THEN** the resulting library `object.Resource` carries the same value, instance, component, and transformer
- **AND** the renderer does not export it: the result carries no inventory entry, and the reconciler's one `object.Export` is the only export of the resource

#### Scenario: One compiled-output type

- **WHEN** the controller's adapter names the type it converts from
- **THEN** it SHALL be the type the kernel package declares, converted by the library's `object.Resources`, and the controller SHALL NOT declare a resource type of its own for it

#### Scenario: Warnings reach the reconciler

- **WHEN** the render reports an unhandled optional trait, or a path whose resolved-versions row is marked newer under the `Warn` policy
- **THEN** the render result carries one controller-composed warning naming that fact, and the render result carries no library-authored message string

#### Scenario: Two registrations in one module are refused before apply

- **WHEN** a render's compiled objects include two `TransformerRegistration` values with the same name, produced by two components
- **THEN** adaptation fails with the library's duplicate-identity error naming that identity and both components with their transformers, and there is no render result, so no resource reaches conversion

#### Scenario: Distinct identities adapt as before

- **WHEN** every compiled object has a distinct apiVersion, kind, namespace and name
- **THEN** adaptation succeeds and produces one resource per object

### Requirement: The render reads the instance's contract demand from the kernel

Both kernel renderers (`KernelModuleRenderer` for a ModuleInstance, `KernelPackageRenderer` for a ModulePackage) SHALL report on `RenderResult.RequiredContracts` the contract demand the kernel's render reports on its diagnostics: every `#resources` and `#traits` key of every component of the instance, sorted and deduplicated. The controller SHALL NOT compute the demand itself by reading the instance's components. An empty demand SHALL be reported as an empty, non-nil list. Source: enhancement 0013 ("The kernel reports every contract a render required, synthesised components included").

The demand is read only from a successful render; this requirement does not change when the reconciler writes it (see `reconcile-loop-assembly`, "The instance's contract demand is recorded on its status").

#### Scenario: A rendered module reports its components' contracts

- **WHEN** `RenderModule` renders a fixture module whose components declare resources and traits
- **THEN** the result's `RequiredContracts` is exactly the sorted, deduplicated set of those contract FQNs

#### Scenario: Both renderers fill the demand the same way

- **WHEN** a ModuleInstance and a ModulePackage render instances of the same module
- **THEN** both results carry the demand the kernel reported for their render, read through the same adapter

#### Scenario: An empty demand is an empty list

- **WHEN** the kernel's render reports no contracts, or reports the list as absent
- **THEN** the result's `RequiredContracts` is an empty, non-nil list

#### Scenario: No controller walk of the components

- **WHEN** the controller's render package is searched for a reader of a component's `#resources` or `#traits`
- **THEN** none exists; the demand comes from the kernel's render diagnostics only
