# OPM controller

The Kubernetes controller that renders OPM modules and applies the result to a cluster. It owns four CRDs (`ModuleInstance`, `ModulePackage`, `Platform`, `TransformerRegistration`) and calls the library Kernel for every render.

## Language

### Resources

**ModuleInstance**:
The namespaced custom resource that asks the controller to render one module at a pinned version with given values and apply the result. The CUE form of the same idea is `#ModuleInstance`; the two are not the same object.
_Avoid_: ModuleRelease, Release, release, deployment

**ModulePackage**:
The namespaced custom resource that points at a Flux source artifact holding a CUE package whose `instance.cue` evaluates to `#ModuleInstance`. It is the second way to get an instance rendered, beside a registry pull.
_Avoid_: bundle, release package

**Platform**:
The cluster-scoped singleton custom resource, named `cluster`, that declares the catalogs the controller renders against. It projects the core `#Platform` author surface.
_Avoid_: environment, cluster config

**TransformerRegistration**:
The cluster-scoped custom resource a provider module ships to claim that its catalog implements platform contracts. Its name is the dot-joined `<namespace>.<name>` of the claiming instance.
_Avoid_: registration (alone), provider CR

**Module**:
A CUE module resolved from an OCI registry by `spec.module.path` and `spec.module.version`; it is what an instance renders. A module is not a catalog.
_Avoid_: chart, package

**Component**:
The OPM unit inside a module. The inventory records the component that produced each owned resource.
_Avoid_: service, workload

### Platform and catalogs

**Catalog**:
A CUE module, addressed by a major-suffixed path such as `opmodel.dev/catalogs/opm@v4`, that carries transformers and contracts for a platform to use.
_Avoid_: module (for a catalog), plugin

**Contract**:
A fully qualified name that a component declares it needs and that a catalog's transformers fulfil. An instance lists the contracts it demands in `status.requiredContracts`.
_Avoid_: capability, interface

**Subscription**:
An entry in `Platform.spec.registry` by which a platform admin pins one catalog build and lists the contracts it provides. It is live the moment it is written.
_Avoid_: dependency, install

**Registry** (platform registry):
The set of catalogs a platform runs: the subscriptions plus the catalogs of active claims. `Platform.status.registry` records which source each entry has (`Subscription` or `Registration`).
_Avoid_: OCI registry (that is where modules are pulled from)

**Claim**:
A TransformerRegistration seen as a statement: a provider says its catalog fulfils named contracts. A claim is accepted, then active.
_Avoid_: registration (alone), PVC claim (see Data claim)

**Provider**:
The ModuleInstance that rendered a claim, named by the claim's `providerRef`. The controller stamps it; a module author cannot set it.
_Avoid_: plugin, operator

**Acceptance**:
The verdict that a claim passed: its catalog resolved, its listed contracts equal those re-derived from the catalog, and no other instance holds the same provider. It is separate from activation.
_Avoid_: approval, validation

**Activation**:
The latch that turns an accepted claim on once its provider reports Ready. It never turns off again; only deleting the claim ends it.
_Avoid_: enable (that is the subscription field), health mirroring

**Platform module**:
The CUE module the controller generates from the Platform on its own disk, one directory per generation, and builds once. Renders consume its built form.
_Avoid_: platform package (the status field `packageIdentity` names the built form's identity), generated config

**Package identity**:
The value in `Platform.status.packageIdentity` that identifies the generated platform by its two inputs: the Platform generation and the sorted set of active claims' catalog coordinates.
_Avoid_: platform version, platform hash

**Skew policy**:
The Platform setting, `Warn` (default) or `Refuse`, for a module that requires a newer build of an OPM-namespace path than the platform pins. `Refuse` stops the render with reason `SkewRefused`.
_Avoid_: version check, compatibility mode

### Rendering

**Kernel**:
The single library Kernel of the process. Every verb builds in its own CUE context, so the controller calls it without any lock.
_Avoid_: engine, evaluator

**Render**:
The step that turns an instance and the platform module into the objects to apply, by one Kernel call. Apply, prune and status work follow it.
_Avoid_: compile, template, build

**Render slot**:
One of the `--max-concurrent-renders` places a reconcile takes before it renders and holds until the result is exported. ModuleInstances and ModulePackages share the pool.
_Avoid_: worker, lock, semaphore

**Render inputs**:
The record in `status.lastAppliedInputs` of what the last applied render consumed: a digest over module source, values, platform package identity, skew policy and the controller and library versions, plus the render time.
_Avoid_: render cache key, config

### Ownership and deletion

**Owner**:
The `spec.owner` field saying who manages an instance: `controller` (the default) or `cli`. A `cli`-owned instance gets no render, apply, prune or finalizer; the controller records only a ManagedExternally acknowledgement.
_Avoid_: managed-by (that is a label), manager

**Inventory**:
The list in `status.inventory` of the resources an instance owns. It alone decides what prune may delete, and it is replaced only after a fully successful reconcile.
_Avoid_: manifest, resource list, label selector

**Prune**:
Deleting inventory entries that the new render no longer produces, when `spec.prune` is true. It runs only after apply succeeds.
_Avoid_: garbage collection, cleanup (cleanup is the finalizer path)

**Safety exclusion**:
The rule that prune never deletes a Namespace or a CustomResourceDefinition, whatever `spec.prune` says.
_Avoid_: protected resource, allowlist

**Data claim**:
A PersistentVolumeClaim, kept by default under `spec.dataPolicy: Keep`. `Delete` lets prune and a forced recreate remove it.
_Avoid_: claim (alone, which means a TransformerRegistration), volume

**Adopt**:
Handing an existing live object to an instance by setting the `opmodel.dev/adopt` annotation to the instance UUID. The apply takes the object in and never deletes and recreates it.
_Avoid_: import, take over

**Instance UUID**:
The identity in `status.instanceUUID` that matches the `module-instance.opmodel.dev/uuid` label on an instance's objects. Prune and deletion use it to tell the instance's objects from another's.
_Avoid_: owner reference, release id

**Drift**:
A managed resource that differs from desired state. The controller detects it and does not correct it in v1alpha1.
_Avoid_: out of sync, divergence

### Failure

**Transient failure**:
A failure that may pass without a spec change, such as a network error or a temporary API error. The controller keeps the last good inventory and requeues with backoff.
_Avoid_: soft failure, retryable error

**Stalled failure**:
A failure that repeating the reconcile cannot fix, such as invalid input or an unsupported artifact. The controller sets `Stalled=True`, waits for a spec or source change, and rechecks every 30 minutes.
_Avoid_: permanent failure, hard failure, fatal
