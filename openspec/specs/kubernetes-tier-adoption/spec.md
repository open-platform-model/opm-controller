# kubernetes-tier-adoption Specification

## Purpose
Defines how the operator adopts the library's Kubernetes tier: it takes resource conversion and the OPM label vocabulary from `opm/k8s/object` and `opm/k8s/labels`, keeps no copy of them, and its lint refuses one reintroduced at the old path.

## Requirements

### Requirement: The operator takes the Kubernetes object code from the library and declares no copy
The operator SHALL take the Kubernetes resource wrapper over the kernel's compiled output, its JSON and unstructured conversion and the duplicate-identity check from the library's Kubernetes tier (`opm/k8s/object` and `opm/k8s/labels`). The operator SHALL declare no type, constant or function of its own for any of them and no alias to them. Source: enhancement 0012 ("The Kubernetes runtime decisions live in a fenced tier, `opm/k8s`, beside the kernel and inside the library module").

#### Scenario: No local copy

- **WHEN** the operator's Go source is searched for the package `pkg/core`, a `Resource` type wrapping a CUE value with instance, component and transformer provenance, or a declaration of the `app.kubernetes.io/managed-by` or `module-instance.opmodel.dev/uuid` label key
- **THEN** none is found; every reader imports `opm/k8s/object` or `opm/k8s/labels`

### Requirement: The operator takes the OPM label code from the library and declares no copy
The operator SHALL take the OPM label keys and values, and the recognition of an OPM manager label value from the library's Kubernetes tier (`opm/k8s/object` and `opm/k8s/labels`). The operator SHALL declare no type, constant or function of its own for any of them and no alias to them. Source: enhancement 0012 ("The Kubernetes runtime decisions live in a fenced tier, `opm/k8s`, beside the kernel and inside the library module").

#### Scenario: The label code comes from the library
- **WHEN** the operator needs an OPM label key or value, or the recognition of an OPM manager label value
- **THEN** it takes them from `opm/k8s/object` or `opm/k8s/labels` and declares no type, constant, function or alias of its own for them

### Requirement: Lint refuses an import of the operator's earlier object copy
The operator's lint SHALL refuse an import of the package path that held its earlier copy, `github.com/open-platform-model/opm-operator/pkg/core`. Source: enhancement 0012 ("The Kubernetes runtime decisions live in a fenced tier, `opm/k8s`, beside the kernel and inside the library module").

#### Scenario: Lint refuses the old path

- **WHEN** a package exists at `github.com/open-platform-model/opm-operator/pkg/core` and a Go file in the repository imports it
- **THEN** `task dev:lint` fails on the rule `no-local-kubernetes-tier-copy`, with a message naming the library packages to use instead

### Requirement: The reconciler exports a render once and takes the digest, the objects and the entries from it
The reconciler SHALL convert a render's resources with one call to the library's `object.Export`, which exports each resource from CUE once. The reconciler SHALL take the render digest, the unstructured objects it applies and the inventory entries from that one result, inside the render slot. The reconciler SHALL drop the rendered CUE values as soon as the export returns, so the digest and the entries read exported data only. No other code path SHALL export a rendered resource.

#### Scenario: Digest and apply objects from one export

- **WHEN** a render of N resources is converted for apply
- **THEN** the render digest equals the library's `inventory.RenderDigest` over `object.Export` of the same resources, and the N apply objects are that export's `Exported.Object` values in input order

#### Scenario: Inventory entries from the same export

- **WHEN** a render of N resources is converted for apply
- **THEN** the conversion carries N inventory entries, one per `Exported.Object` in input order, and the render result carries none

### Requirement: A failure of the export is reported by the conversion
A failure of the export SHALL be reported by the conversion: a value that will not export is a render failure named "computing render digest", and exported JSON that will not decode to an object is an apply failure named "converting resources" with the resource. With the renderer's export gone, the conversion is the first place such a failure is reported.

#### Scenario: Export failure keeps its reason

- **WHEN** one rendered value fails to export
- **THEN** the conversion fails with reason `RenderFailed` and a message that begins "computing render digest", and the rendered resources are dropped

#### Scenario: Decode failure keeps its reason

- **WHEN** one rendered value exports to JSON that does not decode to an object
- **THEN** the conversion fails with reason `ApplyFailed` and a message that begins "converting resources", and the rendered resources are dropped

### Requirement: The operator computes the inventory stale set and the digests only with the library
The operator SHALL compute the inventory stale set, the inventory digest and the render digest only with the library's `opm/k8s/inventory` (`StaleSet`, `Digest`, `RenderDigest`). The operator SHALL build inventory entries with that package's `NewEntry`. Source: enhancement 0012 ("The Kubernetes runtime decisions live in a fenced tier, `opm/k8s`, beside the kernel and inside the library module").

#### Scenario: The stale set and the digests come from the library
- **WHEN** the operator computes the inventory stale set, the inventory digest or the render digest
- **THEN** it uses `StaleSet`, `Digest` or `RenderDigest` of `opm/k8s/inventory`, and it builds inventory entries with `NewEntry`

### Requirement: The operator declares no inventory code beyond the conversion
The operator SHALL declare no entry constructor, identity relation, stale-set function, inventory digest or render digest of its own. Its only inventory code is the conversion between the API's `InventoryEntry` and `inventory.Entry`. Source: enhancement 0012 ("The Kubernetes runtime decisions live in a fenced tier, `opm/k8s`, beside the kernel and inside the library module").

#### Scenario: No local copy

- **WHEN** the operator's Go source is searched for `ComputeStaleSet`, `ComputeDigest`, `NewEntryFromResource`, `IdentityEqual`, `K8sIdentityEqual` or a render digest function in `internal/status`
- **THEN** none is found, and every reader calls `opm/k8s/inventory`

### Requirement: A test keeps the inventory conversion package closed
A test of `internal/inventory` SHALL refuse any further export from that package, so a copy cannot come back there unnoticed. Source: enhancement 0012 ("The Kubernetes runtime decisions live in a fenced tier, `opm/k8s`, beside the kernel and inside the library module").

#### Scenario: A copy in the conversion package fails the tests

- **WHEN** a function is exported from `internal/inventory` beyond `Current`, `FromEntries`, `FromEntry`, `ToEntries` and `ToEntry`
- **THEN** the package's closed-surface test fails

### Requirement: The operator judges the readiness of an inventory object only with the library's health package
The operator SHALL judge the readiness of an object in `status.inventory` only with the library's `opm/k8s/health` (`Evaluate`, `IsHealthy`, `Aggregate`, `ProgressDeadlineExceeded`). The operator fetches each object with its own client and passes it in, because the library reads no cluster. Source: enhancement 0012 ("The Kubernetes runtime decisions live in a fenced tier, `opm/k8s`, beside the kernel and inside the library module").

#### Scenario: The operator passes the fetched object to the library
- **WHEN** the operator judges the readiness of an object in `status.inventory`
- **THEN** the operator fetches the object with its own client and passes it in, because the library reads no cluster

### Requirement: The Healthy condition code declares no readiness rule of its own
The code that computes the `Healthy` condition SHALL read no status field of a fetched object and SHALL declare no readiness rule of its own for any Kubernetes kind: no rollout, replica-count, Job, PersistentVolumeClaim or `Ready`-condition check, and no alias to the library's. Source: enhancement 0012 ("The Kubernetes runtime decisions live in a fenced tier, `opm/k8s`, beside the kernel and inside the library module").

#### Scenario: No local evaluator

- **WHEN** the non-test files of `internal/reconcile` are parsed by the package's unit test
- **THEN** none calls an `unstructured.Nested*` accessor and `health.go` holds no `"status"` string literal, so every verdict comes from `health.Evaluate`, `health.IsHealthy`, `health.Aggregate` and `health.ProgressDeadlineExceeded`

### Requirement: The health package covers applied inventory objects only
This requirement covers applied inventory objects only; the operator's own resources (a `dependsOn` ModulePackage, a TransformerRegistration's provider, the `cluster` Platform) SHALL keep their own condition checks. Source: enhancement 0012 ("The Kubernetes runtime decisions live in a fenced tier, `opm/k8s`, beside the kernel and inside the library module").

#### Scenario: The operator's own resources keep their own condition checks
- **WHEN** the operator judges a `dependsOn` ModulePackage, a TransformerRegistration's provider or the `cluster` Platform
- **THEN** it uses the condition checks of those resources, not the library's `opm/k8s/health`

### Requirement: A prune or cleanup delete follows a delete verdict of the library's ownership package
Every delete of a cluster object that the operator makes in a prune of stale resources or in a deletion cleanup, for a ModuleInstance or a ModulePackage, MUST follow a delete verdict of the library's `opm/k8s/ownership` package on a live read of that object. The operator MUST take the kinds that are never deleted from the same package. Source: enhancement 0012 ("The library owns the whole deletion sequence but only the per-object verdict and the order for apply").

#### Scenario: The prune's outcome is the verdict's
- **GIVEN** one stale object for each answer of the library's delete verdict
- **WHEN** the prune runs
- **THEN** it deletes exactly the objects for which the verdict says proceed, and names each other object with the verdict's reason

### Requirement: The operator sets no install admission input in a delete verdict
Every delete of a cluster object that the operator makes in a prune of stale resources or in a deletion cleanup, for a ModuleInstance or a ModulePackage, MUST follow a delete verdict of the library's `opm/k8s/ownership` package on a live read of that object. The operator MUST NOT set the verdict's install admission input. Source: enhancement 0012 ("The library owns the whole deletion sequence but only the per-object verdict and the order for apply").

#### Scenario: The install admission input is not set
- **WHEN** the operator asks the library's `opm/k8s/ownership` package for a delete verdict
- **THEN** the operator does not set the verdict's install admission input

### Requirement: Prune and cleanup deletes are lifecycle actions
Every delete of a cluster object that the operator makes in a prune of stale resources or in a deletion cleanup, for a ModuleInstance or a ModulePackage, MUST be an action the library's `opm/k8s/lifecycle` transition named, performed with the propagation policy and the precondition of that action. Source: enhancement 0012 ("The library owns the whole deletion sequence but only the per-object verdict and the order for apply").

#### Scenario: A delete is an action of the lifecycle transition
- **WHEN** the operator deletes a cluster object in a prune of stale resources or in a deletion cleanup
- **THEN** the delete is an action the library's `opm/k8s/lifecycle` transition named, performed with the propagation policy and the precondition of that action

### Requirement: The operator builds each deletion plan with the library's constructor
The operator MUST build each plan with the library's constructor and MUST NOT order the entries or set a propagation policy with code of its own. Source: enhancement 0012 ("The library owns the whole deletion sequence but only the per-object verdict and the order for apply").

#### Scenario: The plan's order and propagation reach the API server unchanged
- **GIVEN** a plan over a ConfigMap and a Deployment
- **WHEN** the operator performs it
- **THEN** the DELETE requests arrive in the plan's order, each with the action's propagation policy and UID precondition

### Requirement: One function of the operator performs the deletion plan
The operator MUST build each plan with the library's constructor and MUST NOT order the entries or set a propagation policy with code of its own. One function of the operator MUST perform the plan's reads and deletes. Source: enhancement 0012 ("The library owns the whole deletion sequence but only the per-object verdict and the order for apply").

#### Scenario: One function performs the plan's reads and deletes
- **WHEN** the operator performs a deletion plan
- **THEN** one function of the operator performs the plan's reads and deletes

### Requirement: The forced recreate delete asks no verdict
Every delete of a cluster object that the operator makes in a prune of stale resources or in a deletion cleanup, for a ModuleInstance or a ModulePackage, MUST follow a delete verdict of the library's `opm/k8s/ownership` package on a live read of that object. The one other delete the operator makes, the delete inside a forced recreate, asks no verdict. Source: enhancement 0012 ("The library owns the whole deletion sequence but only the per-object verdict and the order for apply").

#### Scenario: The forced recreate delete asks no verdict
- **WHEN** the operator deletes an object inside a forced recreate
- **THEN** the delete asks no verdict

### Requirement: The forced recreate delete carries the UID precondition
The one other delete the operator makes, the delete inside a forced recreate, MUST carry the UID precondition (capability `ssa-apply`). Source: enhancement 0012 ("The library owns the whole deletion sequence but only the per-object verdict and the order for apply").

#### Scenario: The forced recreate is not a plan step
- **GIVEN** an apply that is allowed to recreate objects, and an object the API server refuses to update
- **WHEN** the object is deleted and created again
- **THEN** the DELETE carries the UID precondition and no deletion plan is built for it

### Requirement: A test keeps the places that may delete a cluster object closed
A test MUST keep the list of places that may delete a cluster object closed: the function that performs the plan, and the forced recreate. Source: enhancement 0012 ("The library owns the whole deletion sequence but only the per-object verdict and the order for apply").

#### Scenario: A new delete outside the listed places fails the tests
- **GIVEN** a change that adds a call that deletes a cluster object in a file the test does not list
- **WHEN** the unit tests run
- **THEN** the call-site test fails and names the file

### Requirement: A test finds a delete by the type of the receiver
A test MUST find a delete by the type of the receiver, a Kubernetes client, and not by the method name alone, so that a method of the same name on another type is not counted, and a second test MUST show that the matcher finds a delete of each form the client offers. Source: enhancement 0012 ("The library owns the whole deletion sequence but only the per-object verdict and the order for apply").

#### Scenario: A method of the same name on another type is not a delete
- **GIVEN** the operator's calls that remove a status condition
- **WHEN** the call-site test runs
- **THEN** it does not count them

### Requirement: A write follows an apply verdict of the library's ownership package
Every write of a cluster object that the operator makes on behalf of a ModuleInstance or a ModulePackage MUST follow an apply verdict of the library's `opm/k8s/ownership` package on a live read of that object. Source: enhancement 0012 ("The apply guard runs on every apply, refusing foreign objects outside the inventory and any object adopted by another instance; a per-object adopt annotation is its only override").

#### Scenario: The apply's outcome is the verdict's
- **GIVEN** one rendered object for each answer of the library's apply verdict
- **WHEN** the guard runs
- **THEN** it allows exactly the objects for which the verdict says apply, and names each other object with the verdict's reason and message

### Requirement: The operator declares no ownership rule of its own for an apply
The operator MUST declare no ownership rule of its own for an apply: it compares no managed-by label, no UUID label and no adopt annotation itself. Source: enhancement 0012 ("The apply guard runs on every apply, refusing foreign objects outside the inventory and any object adopted by another instance; a per-object adopt annotation is its only override").

#### Scenario: No ownership rule of the operator's own
- **WHEN** the operator decides a write of a cluster object on behalf of a ModuleInstance or a ModulePackage
- **THEN** it compares no managed-by label, no UUID label and no adopt annotation itself

### Requirement: The operator never writes the adopt annotation
The operator MUST NOT set, change or remove the `opmodel.dev/adopt` annotation on any object. Source: enhancement 0012 ("The library owns the whole deletion sequence but only the per-object verdict and the order for apply" and "The apply guard runs on every apply, refusing foreign objects outside the inventory and any object adopted by another instance; a per-object adopt annotation is its only override").

#### Scenario: No write of the adopt annotation
- **WHEN** the non-test Go files of the operator are parsed by a unit test
- **THEN** none names the library's adopt annotation key, so the annotation is read only inside the library's verdicts and is never written

### Requirement: A test keeps the places that may write a cluster object closed
A test MUST keep the list of places that may write a cluster object closed: the staged apply, the dry runs of drift detection and of the claim check, and the status and finalizer patches of the operator's own kinds. Source: enhancement 0012 ("The apply guard runs on every apply, refusing foreign objects outside the inventory and any object adopted by another instance; a per-object adopt annotation is its only override").

#### Scenario: A new write outside the listed places fails the tests
- **GIVEN** a change that adds a call that creates, updates, patches or applies a cluster object in a file the test does not list
- **WHEN** the unit tests run
- **THEN** the call-site test fails and names the file

### Requirement: A test finds a write by the type of the receiver
A test MUST find a write by the type of the receiver, as the test of the delete call sites does, and a second test MUST show that the matcher finds a write of each form the client and the resource manager offer. Source: enhancement 0012 ("The apply guard runs on every apply, refusing foreign objects outside the inventory and any object adopted by another instance; a per-object adopt annotation is its only override").

#### Scenario: A write is found by the type of the receiver
- **WHEN** a second test runs the matcher over a write of each form the client and the resource manager offer
- **THEN** the matcher finds each of them
