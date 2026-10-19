## REMOVED Requirements

### Requirement: Acquire a module from the registry as a library Module
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "The operator acquires a module from the registry as a library Module", "Acquisition delegates to the library kernel's registry module loader" and "Acquisition synthesizes no wrapper package and writes no temporary directory".

## ADDED Requirements

### Requirement: The operator acquires a module from the registry as a library Module
The operator SHALL acquire a `ModuleRelease`'s target module from the OCI registry as a library `*module.Module`, given the module path, version, and registry mapping. Registry configuration and the absence of process-environment mutation are the library's responsibility (the kernel's configured registry); acquisition remains safe to call concurrently from reconcilers.

#### Scenario: Published module is acquired

- **WHEN** acquisition is called with the path and version of a module published in the registry
- **THEN** it returns a `*module.Module` whose decoded metadata (name, version) matches the published module

#### Scenario: Unresolvable module surfaces a load error

- **WHEN** acquisition is called with a path/version not present in the registry
- **THEN** it returns an error identifying the failed module load

### Requirement: Acquisition delegates to the library kernel's registry module loader
Acquisition SHALL delegate to the library kernel's registry module loader, `Kernel.LoadModuleFromRegistry(ctx, modPath, version)`, and decode the returned value via `Kernel.NewModuleFromValue`. Delegating to the library is load-bearing: the wrapper approach re-embedded core@v0's `#Module` (collapsing its self-referential `metadata.modulePath`/`version`), could not resolve the target's transitive catalog/core dependencies, and failed the loader's root shape gate.

#### Scenario: Acquisition delegates to the kernel, not a wrapper package

- **WHEN** the `internal/moduleacquire` package is inspected
- **THEN** `Acquire` calls `Kernel.LoadModuleFromRegistry` and there is no generated wrapper package, no `writeShim`, and no temporary-directory staging

### Requirement: Acquisition synthesizes no wrapper package and writes no temporary directory
As in requirement "Acquisition delegates to the library kernel's registry module loader": Acquisition SHALL NOT synthesize a wrapper CUE package that imports and embeds the target module, and SHALL NOT write the module's source to a temporary directory. The library loader loads the module as the main module, its own `cue.mod/module.cue` resolves transitive deps and its `kind`/`metadata` sit at the package root, so author-set metadata is preserved and dependencies resolve.

#### Scenario: Core@v0 self-referential metadata is preserved

- **WHEN** acquisition loads a module authored against `opmodel.dev/core@v0` (whose `#Module` derives `modulePath`/`version` from themselves)
- **THEN** the decoded module's `metadata.modulePath` and `metadata.version` equal the author-set values, with no `field not allowed` admission error

#### Scenario: Module with catalog imports is acquired

- **WHEN** acquisition loads a module whose components import catalog resource definitions (for example `opmodel.dev/catalogs/opm/resources`)
- **THEN** those transitive dependencies resolve and acquisition succeeds without the operator declaring them

#### Scenario: Acquisition writes no temporary directory
- **WHEN** acquisition loads a module
- **THEN** it does not write the module's source to a temporary directory
