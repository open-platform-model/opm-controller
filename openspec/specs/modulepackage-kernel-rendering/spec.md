# modulepackage-kernel-rendering

## Purpose

The `ModulePackage` reconciler renders its Flux-fetched package through the kernel-backed `KernelPackageRenderer` against the generated platform: a `#ModuleInstance` package that carries its own values is acquired through the kernel's on-disk acquisition and rendered through the single-build render with the leased platform record and its skew policy; rendering blocks inertly when no platform is generated, retries promptly when the platform changes, and packages of any other kind are rejected.

## Requirements

### Requirement: ModulePackage renders through the kernel against the generated platform

`KernelPackageRenderer` SHALL acquire the extracted package as a source-carrying instance through the kernel's on-disk acquisition and render it through the single-build render with the leased platform record and its skew policy.

#### Scenario: ModuleInstance package renders and applies

- **WHEN** a `ModulePackage` artifact holding a `#ModuleInstance` package is rendered while a generated platform is recorded
- **THEN** the rendered resources are applied and recorded as before

### Requirement: Block ModulePackage when no platform is generated

When the store holds no generated-module record, the package renderer SHALL return `ErrPlatformNotReady` after kind detection and before any build, and the reconciler SHALL set `Ready=False` reason `PlatformNotReady`.

#### Scenario: No platform present blocks the package inertly

- **WHEN** a `ModulePackage` is reconciled while no platform is recorded
- **THEN** its status carries `PlatformNotReady` and nothing is applied

### Requirement: Packages of any other kind are rejected

For a fetched package whose `kind` is anything other than `ModuleInstance`, the renderer SHALL return `ErrUnsupportedKind` and the reconciler SHALL surface `Ready=False` with reason `UnsupportedKind` and `Stalled=True`. The rejection SHALL NOT name speculative kinds: the kernel's `#ModuleInstance` shape gate (`oerrors.ErrWrongKind`, the sentinel the library's `opm/errors` package declares) is the detection mechanism, and the resulting error is generic.

#### Scenario: Wrong-kind package is rejected

- **WHEN** a `ModulePackage` whose fetched package has a `kind` other than `ModuleInstance` is reconciled
- **THEN** rendering returns an unsupported-kind error
- **AND** the status reflects `UnsupportedKind` and nothing is applied

### Requirement: The ModulePackage reconciler watches the Platform
The `ModulePackage` reconciler SHALL watch the `Platform` resource. On a Platform update, it SHALL re-enqueue all `ModulePackages` only when a field that a package render consumes differs between the old and the new object. Platform create and delete events SHALL re-enqueue.

#### Scenario: Blocked package retries when the platform is generated

- **WHEN** a `ModulePackage` is blocked with `PlatformNotReady` and a `Platform` is then applied and generated
- **THEN** the reconciler re-enqueues the `ModulePackage`
- **AND** on the next reconcile it renders and applies against the generated platform

### Requirement: ModulePackages are re-enqueued on the fields a package render consumes
As in requirement "The ModulePackage reconciler watches the Platform": A package render reads the same generated platform record as a ModuleInstance render, so the fields SHALL be the same ones the ModuleInstance Platform watch compares:
- the `Ready` condition's status
- the pin set: `status.packageIdentity` or `status.registry`
- `spec.skewPolicy`
- `status.observedGeneration`
- `status.operatorVersion`

#### Scenario: Packages re-render under a new pin set

- **WHEN** the Platform's `status.packageIdentity` changes because a claim became active, and the Platform's generation is unchanged
- **THEN** the reconciler re-enqueues every `ModulePackage`

### Requirement: The Platform status update and an operatorVersion change re-enqueue ModulePackages
As in requirement "ModulePackages are re-enqueued on the fields a package render consumes": The trigger SHALL include the Platform reconciler's own status update, which does not bump the Platform's generation. This lets packages blocked on `PlatformNotReady` retry promptly, not only on their interval requeue. A `status.operatorVersion` change SHALL re-enqueue because, after an operator upgrade, it is the only status change the regenerated Platform writes.

#### Scenario: An operator upgrade into an already-Ready Platform recovers blocked packages

- **GIVEN** a Platform that is already `Ready=True` with reason `Generated`, and an operator that restarts under a new version with an empty platform store
- **WHEN** the packages render before the platform is regenerated and report `PlatformNotReady`, and the Platform reconciler then regenerates the platform and writes a status whose only change is `status.operatorVersion`
- **THEN** the reconciler re-enqueues every `ModulePackage` and renders it on the next reconcile, without waiting for the interval requeue

### Requirement: A Platform update that changes none of the consumed fields enqueues no ModulePackage
As in requirement "ModulePackages are re-enqueued on the fields a package render consumes": A Platform update that changes none of these fields SHALL NOT enqueue any `ModulePackage`. Examples are a change to the `Ready` message, a change to its reason while it stays `False`, a `ContractsFulfilled` update, or a bump of `metadata.generation` alone. A `spec.skewPolicy` edit can render each package once under the previous policy before the `observedGeneration` write renders it under the new one.

#### Scenario: A message-only status write enqueues no package

- **WHEN** the Platform reconciler rewrites the `Ready` condition's message or updates `ContractsFulfilled`, and the `Ready` status, `packageIdentity`, `registry`, `skewPolicy`, `observedGeneration` and `operatorVersion` are unchanged
- **THEN** no `ModulePackage` is enqueued and no package renders

### Requirement: Unresolved demands are ResolutionFailed and a skew refusal is SkewRefused
The reconciler SHALL classify a refused render by its typed cause: unresolved platform demands and unmatched components SHALL be `ResolutionFailed`; a catalog-skew refusal SHALL be `SkewRefused` with a message naming the module path, the module's required build and the platform's build.
(Requirement: A duplicate identity is DuplicateIdentities and any other refusal is RenderFailed)

#### Scenario: Package demands contracts the platform does not provide

- **WHEN** a package renders against a platform whose catalogs do not provide contracts the module demands
- **THEN** the ModulePackage reports `Ready=False` with reason `ResolutionFailed` and `Stalled=True`, and a Warning event carries the unresolved-demands message

#### Scenario: Skew refusal is distinct

- **WHEN** the platform's policy is `Refuse` and the package's module requires a newer catalog build than the platform pins
- **THEN** the package reports `Ready=False` reason `SkewRefused` naming the path and both versions, applies nothing, and requeues on the stalled recheck interval

### Requirement: A duplicate identity is DuplicateIdentities and any other refusal is RenderFailed
The reconciler SHALL classify a refused render by its typed cause: a render whose compiled objects share one apply identity SHALL be `DuplicateIdentities` with the library's message naming each identity and every producing component and transformer (enhancement 0015, "One registration per module"); a transform failure, an over-subscribed provider contract, or any other refusal SHALL be `RenderFailed`.
(Requirement: Unresolved demands are ResolutionFailed and a skew refusal is SkewRefused)

#### Scenario: Duplicate identities are a refusal of their own

- **WHEN** a render's compiled objects share one apply identity
- **THEN** the object reports `Ready=False` with reason `DuplicateIdentities` and `Stalled=True`, the message names the identity and both producing components, nothing is applied, and the same reason is produced whether the render was a ModuleInstance's or a ModulePackage's

#### Scenario: Over-subscribed provider contract is a render failure

- **WHEN** the build refuses because two enabled catalogs provide the same provider-fulfilled contract
- **THEN** the package reports `RenderFailed` with the kernel's message naming the contract key and both catalogs

#### Scenario: Ordinary evaluation error keeps RenderFailed

- **WHEN** the render fails for a cause that is neither a resolution-class failure, a registry fetch failure, a skew refusal nor a duplicate identity
- **THEN** the ModulePackage reports `Ready=False` with reason `RenderFailed`

### Requirement: Every classified refusal stalls the object, emits a Warning event and requeues
Every classified refusal SHALL set `Ready=False` and `Stalled=True`, emit a Warning event and requeue on the stalled recheck interval; the ModuleInstance and ModulePackage classifiers SHALL route typed causes through one shared function so they cannot drift.

#### Scenario: A classified refusal stalls the object
- **WHEN** a render is refused with a classified cause
- **THEN** the object reports `Ready=False` and `Stalled=True`, a Warning event is emitted and the object requeues on the stalled recheck interval

### Requirement: A registry fetch failure is not a refusal and retries on the bounded backoff
A registry fetch failure (the library's typed `*oerrors.FetchError`) with no typed terminal cause is not a refusal: in any phase of the render it SHALL report `ResolutionFailed` without `Stalled` and retry on the bounded backoff, per the `reconcile-backoff` capability.

#### Scenario: Registry fetch failure during the render retries

- **WHEN** the render build fails with the library's typed registry fetch failure and no typed terminal cause
- **THEN** the ModulePackage reports `Ready=False` with reason `ResolutionFailed` and no `Stalled` condition, and requeues on the bounded backoff

### Requirement: A ModulePackage that leaves a required config value unset is refused when it loads
A ModulePackage whose instance values leave a required `#config` value of its module unset SHALL be refused when the package loads, whether or not a component reads the value. A required value is a field declared `foo!`, or a field of a bare type with no default.

#### Scenario: Unset required value that no component reads

- **WHEN** a ModulePackage's module declares `#config: note: string` with no default, no component reads `note`, and the package's values do not set it
- **THEN** the ModulePackage reports `Ready=False` and `Stalled=True` with reason `ResolutionFailed`, the message contains `not fully concrete: values.note: incomplete value string`, and no resource is applied

#### Scenario: The value is set

- **WHEN** the same package's values set `note`
- **THEN** the package loads and the render proceeds to the platform gate

### Requirement: A package value default that does not unify with the config default is refused
As in requirement "A ModulePackage that leaves a required config value unset is refused when it loads": The same SHALL hold for a value that the package's values give a default other than the `#config` default: the two do not unify to a concrete value.

#### Scenario: A default that disagrees with the config default

- **WHEN** the module declares `#config: tier: string | *"a"`, no component reads `tier`, and the package's values hold `tier: string | *"b"`
- **THEN** the ModulePackage reports `Ready=False` and `Stalled=True` with reason `ResolutionFailed`, and the message contains `not fully concrete: values.tier: incomplete value`

### Requirement: The refusal of an unset required value is a package load failure, not retried on the registry backoff
As in requirements "A ModulePackage that leaves a required config value unset is refused when it loads" and "A package value default that does not unify with the config default is refused": The refusal SHALL be the kernel's: the operator SHALL report it as a package load failure, with reason `ResolutionFailed`, `Ready=False` and `Stalled=True`, and SHALL NOT retry it on the registry backoff. Nothing SHALL be applied or pruned.

#### Scenario: A refused package is a load failure that is not retried on the registry backoff
- **WHEN** a ModulePackage is refused because its values leave a required `#config` value unset
- **THEN** the operator reports a package load failure with reason `ResolutionFailed`, `Ready=False` and `Stalled=True`, does not retry it on the registry backoff, and applies and prunes nothing

### Requirement: The message of the refusal carries the prefix, the kernel error text and the findings
As in requirements "A ModulePackage that leaves a required config value unset is refused when it loads" and "A package value default that does not unify with the config default is refused": The message SHALL be the prefix `loading package: `, the kernel's error text up to its first finding, and the kernel's findings, each with its positions; it names the unset value as `values.<field>`.

#### Scenario: The message names the unset value
- **WHEN** a ModulePackage is refused because its values leave a required `#config` value unset
- **THEN** the message is the prefix `loading package: `, the kernel's error text up to its first finding, and the kernel's findings, each with its positions, and it names the unset value as `values.<field>`

### Requirement: The message of a refused ModulePackage names every unset required value with its position
When a ModulePackage is refused because its values leave required `#config` values unset, the message SHALL name every one of them, one finding for each, as `values.<field>: incomplete value <type> (<file>:<line>:<column>)` with the position of its `#config` declaration, joined by `; ` after `loading package: Kernel.AcquireInstanceFromDir: instance "<name>": not fully concrete: `. A nested field SHALL be named by its full path.

#### Scenario: Three values are unset, one nested

- **WHEN** the module declares `#config: {greeting: string, db: host: string, note: string}` and the package's values set none of them
- **THEN** the message names `values.greeting`, `values.db.host` and `values.note`

### Requirement: A value a component reads is named like a value no component reads
As in requirement "The message of a refused ModulePackage names every unset required value with its position": A value a component reads SHALL be named like a value no component reads, and the message SHALL NOT name the place in a component that reads it.

#### Scenario: The component reads the unset value

- **WHEN** a ModulePackage's module declares `#config: {greeting: string, note: string}` with no defaults, a component reads `greeting`, and the package's values set `note` only
- **THEN** the ModulePackage reports `Ready=False` and `Stalled=True` with reason `ResolutionFailed` and the message `loading package: Kernel.AcquireInstanceFromDir: instance "<name>": not fully concrete: values.greeting: incomplete value string (instance.cue:<line>:<column>)`

### Requirement: The message replaces no finding with a count and the reason stays ResolutionFailed
As in requirements "The message of a refused ModulePackage names every unset required value with its position" and "A value a component reads is named like a value no component reads": The message SHALL NOT replace a finding with `(and N more errors)`. The controller SHALL word the findings from the kernel's typed error tree. The reason SHALL stay `ResolutionFailed` with `Ready=False` and `Stalled=True`.

#### Scenario: Two values are unset

- **WHEN** the same package's values set neither
- **THEN** the message names `values.greeting` and `values.note`, each with its position, and holds no `more errors`

### Requirement: A package load failure lists every finding with its positions
When the load of a ModulePackage fails with findings of the kernel (a value of the wrong type, a constraint a value does not meet, a field the module's `#config` does not allow, a value left unset), the message SHALL list every finding with the positions the kernel attributed it to, and SHALL NOT replace findings after the first with a count while the message is inside the condition's limit of 32768 characters.
(Requirement: Past the limit the message ends with a count of findings)

#### Scenario: A value of the wrong type

- **WHEN** a package's module declares `#config: note: string` and the package's values hold `note: 7`
- **THEN** the message contains `#module.#config.note: conflicting values string and 7 (mismatched types string and int) (instance.cue:<line>:<column>, instance.cue:<line>:<column>)`, with reason `ResolutionFailed` and `Stalled=True`

#### Scenario: A constraint a value does not meet

- **WHEN** the module declares `#config: port: int & >0 | *80` and the package's values hold `port: -1`
- **THEN** the message lists three findings for `#module.#config.port`, among them `invalid value -1 (out of bound >0)` with its positions, and holds no `more errors`

#### Scenario: A field the module does not declare

- **WHEN** the package's values hold a field the module's `#config` does not allow
- **THEN** the message contains `field not allowed (instance.cue:<line>:<column>)`

### Requirement: Past the limit the message ends with a count of findings
As in requirement "A package load failure lists every finding with its positions": past that limit it SHALL keep whole findings from the first and end with `; and <N> more findings`.

#### Scenario: A message past the limit keeps whole findings
- **WHEN** the message of a package load failure would pass the condition's limit of 32768 characters
- **THEN** it keeps whole findings from the first and ends with `; and <N> more findings`

### Requirement: A position in a file of the package is relative to the CUE module root
As in requirement "A package load failure lists every finding with its positions": A position in a file of the package SHALL be written relative to the package's CUE module root, so no position names the temporary directory the operator extracts the package to, which is a new one on every reconcile. A position in any other file SHALL be written as the kernel reports it.

#### Scenario: The same package in another directory

- **WHEN** a package refused under the frame `Kernel.AcquireInstanceFromDir: instance "<name>": ` (an unset value, a wrong type no component reads, a constraint, a field not allowed) is loaded from two different extraction directories
- **THEN** both messages are equal and neither names an extraction directory

#### Scenario: A position in any other file is written as the kernel reports it
- **WHEN** a finding of a package load failure has a position in a file that is not a file of the package
- **THEN** the position is written as the kernel reports it

### Requirement: The text in front of the findings is the kernel's
The text in front of the findings is the kernel's and SHALL NOT be changed: for a failure the kernel reports under `Kernel.AcquireInstanceFromDir: instance "<name>": ` it holds no directory, so the whole message of an unchanged package is the same on every reconcile; for a failure of the package build (for example a value of the wrong type that a component reads) the kernel's text is `building instance package from <directory>: ...` and still names the extraction directory.

#### Scenario: A build failure keeps the kernel's frame

- **WHEN** a package's values give a value of the wrong type to a field a component reads, so the package build fails
- **THEN** the message starts `loading package: building instance package from <directory>: `, the findings follow with positions in the package relative to its CUE module root, and the reason is `ResolutionFailed` with `Stalled=True`

### Requirement: The wording does not change the classification and a registry fetch failure keeps its own text
As in requirement "A package load failure lists every finding with its positions": The wording SHALL NOT change how the failure is classified, and a registry fetch failure SHALL keep its own text.

#### Scenario: A registry fetch failure keeps its text

- **WHEN** the package load fails with the library's typed registry fetch failure
- **THEN** the message is `loading package: ` followed by the kernel's error text unchanged, and the failure retries on the backoff as before
