## REMOVED Requirements

### Requirement: Re-enqueue ModulePackages when the platform becomes ready

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-021); this requirement is restated under the controller's name.

**Migration**: See the requirement "ModulePackages re-enqueue when the platform becomes ready" in this capability.

## MODIFIED Requirements

### Requirement: A package that leaves a required config value unset is refused

A ModulePackage whose instance values leave a required `#config` value of its module unset SHALL be refused when the package loads, whether or not a component reads the value. A required value is a field declared `foo!`, or a field of a bare type with no default. The same SHALL hold for a value that the package's values give a default other than the `#config` default: the two do not unify to a concrete value. The refusal SHALL be the kernel's: the controller SHALL report it as a package load failure, with reason `ResolutionFailed`, `Ready=False` and `Stalled=True`, and SHALL NOT retry it on the registry backoff. The message SHALL be the prefix `loading package: `, the kernel's error text up to its first finding, and the kernel's findings, each with its positions; it names the unset value as `values.<field>`. Nothing SHALL be applied or pruned.

#### Scenario: Unset required value that no component reads

- **WHEN** a ModulePackage's module declares `#config: note: string` with no default, no component reads `note`, and the package's values do not set it
- **THEN** the ModulePackage reports `Ready=False` and `Stalled=True` with reason `ResolutionFailed`, the message contains `not fully concrete: values.note: incomplete value string`, and no resource is applied

#### Scenario: A default that disagrees with the config default

- **WHEN** the module declares `#config: tier: string | *"a"`, no component reads `tier`, and the package's values hold `tier: string | *"b"`
- **THEN** the ModulePackage reports `Ready=False` and `Stalled=True` with reason `ResolutionFailed`, and the message contains `not fully concrete: values.tier: incomplete value`

#### Scenario: The value is set

- **WHEN** the same package's values set `note`
- **THEN** the package loads and the render proceeds to the platform gate

### Requirement: A package values failure reports every finding with stable positions

When the load of a ModulePackage fails with findings of the kernel (a value of the wrong type, a constraint a value does not meet, a field the module's `#config` does not allow, a value left unset), the message SHALL list every finding with the positions the kernel attributed it to, and SHALL NOT replace findings after the first with a count while the message is inside the condition's limit of 32768 characters; past that limit it SHALL keep whole findings from the first and end with `; and <N> more findings`. A position in a file of the package SHALL be written relative to the package's CUE module root, so no position names the temporary directory the controller extracts the package to, which is a new one on every reconcile. A position in any other file SHALL be written as the kernel reports it. The text in front of the findings is the kernel's and is not changed: for a failure the kernel reports under `Kernel.AcquireInstanceFromDir: instance "<name>": ` it holds no directory, so the whole message of an unchanged package is the same on every reconcile; for a failure of the package build (for example a value of the wrong type that a component reads) the kernel's text is `building instance package from <directory>: ...` and still names the extraction directory. The wording SHALL NOT change how the failure is classified, and a registry fetch failure SHALL keep its own text.

#### Scenario: A value of the wrong type

- **WHEN** a package's module declares `#config: note: string` and the package's values hold `note: 7`
- **THEN** the message contains `#module.#config.note: conflicting values string and 7 (mismatched types string and int) (instance.cue:<line>:<column>, instance.cue:<line>:<column>)`, with reason `ResolutionFailed` and `Stalled=True`

#### Scenario: A constraint a value does not meet

- **WHEN** the module declares `#config: port: int & >0 | *80` and the package's values hold `port: -1`
- **THEN** the message lists three findings for `#module.#config.port`, among them `invalid value -1 (out of bound >0)` with its positions, and holds no `more errors`

#### Scenario: A field the module does not declare

- **WHEN** the package's values hold a field the module's `#config` does not allow
- **THEN** the message contains `field not allowed (instance.cue:<line>:<column>)`

#### Scenario: The same package in another directory

- **WHEN** a package refused under the frame `Kernel.AcquireInstanceFromDir: instance "<name>": ` (an unset value, a wrong type no component reads, a constraint, a field not allowed) is loaded from two different extraction directories
- **THEN** both messages are equal and neither names an extraction directory

#### Scenario: A build failure keeps the kernel's frame

- **WHEN** a package's values give a value of the wrong type to a field a component reads, so the package build fails
- **THEN** the message starts `loading package: building instance package from <directory>: `, the findings follow with positions in the package relative to its CUE module root, and the reason is `ResolutionFailed` with `Stalled=True`

#### Scenario: A registry fetch failure keeps its text

- **WHEN** the package load fails with the library's typed registry fetch failure
- **THEN** the message is `loading package: ` followed by the kernel's error text unchanged, and the failure retries on the backoff as before

## ADDED Requirements

### Requirement: ModulePackages re-enqueue when the platform becomes ready

The `ModulePackage` reconciler SHALL watch the `Platform` resource. On a Platform update, it SHALL re-enqueue all `ModulePackages` only when a field that a package render consumes differs between the old and the new object. A package render reads the same generated platform record as a ModuleInstance render, so the fields are the same ones the ModuleInstance Platform watch compares:

- the `Ready` condition's status
- the pin set: `status.packageIdentity` or `status.registry`
- `spec.skewPolicy`
- `status.observedGeneration`
- `status.controllerVersion`

The trigger includes the Platform reconciler's own status update, which does not bump the Platform's generation. This lets packages blocked on `PlatformNotReady` retry promptly, not only on their interval requeue. A `status.controllerVersion` change re-enqueues because, after a controller upgrade, it is the only status change the regenerated Platform writes. A Platform update that changes none of these fields SHALL NOT enqueue any `ModulePackage`. Examples are a change to the `Ready` message, a change to its reason while it stays `False`, a `ContractsFulfilled` update, or a bump of `metadata.generation` alone. Platform create and delete events SHALL re-enqueue. A `spec.skewPolicy` edit can render each package once under the previous policy before the `observedGeneration` write renders it under the new one.

#### Scenario: Blocked package retries when the platform is generated

- **WHEN** a `ModulePackage` is blocked with `PlatformNotReady` and a `Platform` is then applied and generated
- **THEN** the reconciler re-enqueues the `ModulePackage`
- **AND** on the next reconcile it renders and applies against the generated platform

#### Scenario: A message-only status write enqueues no package

- **WHEN** the Platform reconciler rewrites the `Ready` condition's message or updates `ContractsFulfilled`, and the `Ready` status, `packageIdentity`, `registry`, `skewPolicy`, `observedGeneration` and `controllerVersion` are unchanged
- **THEN** no `ModulePackage` is enqueued and no package renders

#### Scenario: Packages re-render under a new pin set

- **WHEN** the Platform's `status.packageIdentity` changes because a claim became active, and the Platform's generation is unchanged
- **THEN** the reconciler re-enqueues every `ModulePackage`

#### Scenario: A controller upgrade into an already-Ready Platform recovers blocked packages

- **GIVEN** a Platform that is already `Ready=True` with reason `Generated`, and a controller that restarts under a new version with an empty platform store
- **WHEN** the packages render before the platform is regenerated and report `PlatformNotReady`, and the Platform reconciler then regenerates the platform and writes a status whose only change is `status.controllerVersion`
- **THEN** the reconciler re-enqueues every `ModulePackage` and renders it on the next reconcile, without waiting for the interval requeue
