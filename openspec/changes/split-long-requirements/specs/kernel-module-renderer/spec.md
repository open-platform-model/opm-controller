## REMOVED Requirements

### Requirement: Render a ModuleRelease through the kernel
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "KernelModuleRenderer renders an instance through the kernel and shares nothing between kernel calls", "Values reach synthesis as a stack of values sources with an origin", "No acquisition or synthesis call passes a per-call registry or load-options argument".

### Requirement: Adapt compiled output to operator resources
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "The renderer adapts compiled objects to the library's object resources", "The refusal of duplicate identities uses the library check and carries its error", "The refusal of duplicate identities builds no resource", "The renderer composes the result's warnings from the advisory rows".

### Requirement: The render reports the instance's contract demand from the kernel
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "Both kernel renderers report the contract demand of the kernel's render", "An empty demand is reported as an empty, non-nil list", "The demand is read only from a successful render".

## ADDED Requirements

### Requirement: KernelModuleRenderer renders an instance through the kernel and shares nothing between kernel calls
`KernelModuleRenderer` SHALL acquire the module from the registry, synthesize the instance (source-carrying) with the supplied values, and render it through the single-build render with the leased platform record and the record's skew policy.
Every kernel call SHALL share nothing: acquisition, synthesis and the build each evaluate in a context of their own, so renders of different objects overlap with no gate.

#### Scenario: Renders resources from a generated platform

- **WHEN** `RenderModule` is called for a resolvable module while a generated platform is recorded
- **THEN** the result carries the rendered resources and any warnings the render reported, and no inventory entries: the reconciler builds those from its one export of the resources

### Requirement: Values reach synthesis as a stack of values sources with an origin
Values SHALL reach synthesis as a stack of values sources, each carrying an origin that identifies where the operator read it from, so a values error names that origin rather than an anonymous filename.

#### Scenario: Values are applied when supplied

- **WHEN** non-nil `RawValues` are passed
- **THEN** they are supplied to instance synthesis as a values source whose origin names the CR field they came from
- **AND** when no values are supplied the module's `#config` defaults apply

#### Scenario: A values error names its origin

- **WHEN** the supplied raw values violate the module's `#config` schema
- **THEN** the render fails with an error naming the origin the operator gave that source, not an anonymous filename

### Requirement: No acquisition or synthesis call passes a per-call registry or load-options argument
No acquisition or synthesis call SHALL pass a per-call registry or load-options argument: the registry mapping is the one the shared Kernel was constructed with.

#### Scenario: A call passes no per-call registry argument
- **WHEN** the renderer makes an acquisition or synthesis call
- **THEN** the call passes no per-call registry or load-options argument, and the registry mapping is the one the shared Kernel was constructed with

### Requirement: The renderer adapts compiled objects to the library's object resources
The renderer SHALL adapt the single-build render's compiled objects, the kernel's own compiled-output type, declared beside the render verb, to the library's Kubernetes object resources (`opm/k8s/object`) exactly as before, after first refusing a render whose compiled objects share one Kubernetes apply identity (apiVersion, kind, namespace and name).

#### Scenario: Compiled item maps to a resource

- **WHEN** a library `Compiled` with a value and provenance is adapted
- **THEN** the resulting library `object.Resource` carries the same value, instance, component, and transformer
- **AND** the renderer does not export it: the result carries no inventory entry, and the reconciler's one `object.Export` is the only export of the resource

#### Scenario: One compiled-output type

- **WHEN** the operator's adapter names the type it converts from
- **THEN** it SHALL be the type the kernel package declares, converted by the library's `object.Resources`, and the operator SHALL NOT declare a resource type of its own for it

#### Scenario: Distinct identities adapt as before

- **WHEN** every compiled object has a distinct apiVersion, kind, namespace and name
- **THEN** adaptation succeeds and produces one resource per object

### Requirement: The refusal of duplicate identities uses the library check and carries its error
The refusal of a render whose compiled objects share one Kubernetes apply identity SHALL use the library's duplicate-identity check and SHALL carry its error unchanged, naming each shared identity once and every component and transformer that produced it.

#### Scenario: A refusal carries the library's error unchanged
- **WHEN** a render's compiled objects share one Kubernetes apply identity
- **THEN** the refusal is the error of the library's duplicate-identity check, unchanged, naming each shared identity once and every component and transformer that produced it

### Requirement: The refusal of duplicate identities builds no resource
The refusal of a render whose compiled objects share one Kubernetes apply identity: when it fires, the renderer SHALL build no resource, so nothing of that render can reach apply (enhancement 0015 ("One registration per module"): a second registration in one module is refused before apply naming both carrying components, and enhancement 0015 ("The rendered CR's name is instance-derived: dot-joined `namespace.name`") gives both the same instance-derived name).

#### Scenario: Two registrations in one module are refused before apply

- **WHEN** a render's compiled objects include two `TransformerRegistration` values with the same name, produced by two components
- **THEN** adaptation fails with the library's duplicate-identity error naming that identity and both components with their transformers, and there is no render result, so no resource reaches conversion

### Requirement: The renderer composes the result's warnings from the advisory rows
The render SHALL report no message strings.
The renderer SHALL compose the result's warnings itself from the render's advisory diagnostic rows: the unhandled-trait table, and the resolved-versions rows marked newer.
Each composed warning SHALL name the same facts as before: for skew, the OPM-namespace path, the version the module requires and the version the platform carries; for an unhandled trait, the component and the trait.

#### Scenario: Warnings reach the reconciler

- **WHEN** the render reports an unhandled optional trait, or a path whose resolved-versions row is marked newer under the `Warn` policy
- **THEN** the render result carries one operator-composed warning naming that fact, and the render result carries no library-authored message string

### Requirement: Both kernel renderers report the contract demand of the kernel's render
Both kernel renderers (`KernelModuleRenderer` for a ModuleInstance, `KernelPackageRenderer` for a ModulePackage) SHALL report on `RenderResult.RequiredContracts` the contract demand the kernel's render reports on its diagnostics: every `#resources` and `#traits` key of every component of the instance, sorted and deduplicated.
The operator SHALL NOT compute the demand itself by reading the instance's components.
Source: enhancement 0013.

#### Scenario: A rendered module reports its components' contracts

- **WHEN** `RenderModule` renders a fixture module whose components declare resources and traits
- **THEN** the result's `RequiredContracts` is exactly the sorted, deduplicated set of those contract FQNs

#### Scenario: Both renderers fill the demand the same way

- **WHEN** a ModuleInstance and a ModulePackage render instances of the same module
- **THEN** both results carry the demand the kernel reported for their render, read through the same adapter

#### Scenario: No operator walk of the components

- **WHEN** the operator's render package is searched for a reader of a component's `#resources` or `#traits`
- **THEN** none exists; the demand comes from the kernel's render diagnostics only

### Requirement: An empty demand is reported as an empty, non-nil list
An empty demand SHALL be reported as an empty, non-nil list.
Source: enhancement 0013 ("The kernel reports every contract a render required, synthesised components included").

#### Scenario: An empty demand is an empty list

- **WHEN** the kernel's render reports no contracts, or reports the list as absent
- **THEN** the result's `RequiredContracts` is an empty, non-nil list

### Requirement: The demand is read only from a successful render
The demand SHALL be read only from a successful render; this requirement does not change when the reconciler writes it (see `reconcile-loop-assembly`, "The instance's contract demand is recorded on its status").

#### Scenario: A failed render yields no demand
- **WHEN** a render does not succeed
- **THEN** no contract demand is read from it
