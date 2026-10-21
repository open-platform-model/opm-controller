## MODIFIED Requirements

### Requirement: A Healthy condition reports whether the applied objects have rolled out

The controller SHALL write a `Healthy` condition on a ModuleInstance and on a ModulePackage. It reports whether the objects in `status.inventory` have rolled out, as judged by the library's `opm/k8s/health`. It SHALL be one of:

- `True` with reason `RolledOut`: every inventory object was read and `health.IsHealthy` holds for its `health.Evaluate` status;
- `False` with reason `ProgressDeadlineExceeded`: at least one object is a Deployment for which `health.ProgressDeadlineExceeded` holds;
- `False` with reason `NotRolledOut`: no Deployment is stalled, and at least one object is not healthy or does not exist;
- `Unknown` with reason `HealthUnknown`: no object is known to be unhealthy, and an object could not be read, the reader could not be built, or the inventory is empty.

When more than one applies, `ProgressDeadlineExceeded` wins over `NotRolledOut`, which wins over `HealthUnknown`. `RolledOut` holds only when none of the others applies. The message SHALL begin with the ready count and the total `health.Aggregate` returns. For `ProgressDeadlineExceeded` and `NotRolledOut` it SHALL name up to five objects that are not healthy, each with its kind, namespace, name and status (`Missing` for an object that does not exist), followed by how many more there are. For `ProgressDeadlineExceeded` the stalled Deployments come first, in inventory order, then the other objects that are not healthy, in inventory order; for `NotRolledOut` every named object is in inventory order. For a read failure it SHALL name the error of the first unreadable object in inventory order. A read refused because the object's kind is no longer served (the API server has no mapping for it) counts as a missing object, not an unreadable one, because such an object cannot exist. Two judgements of the same cluster state SHALL produce the same message.

`Ready` SHALL keep its meaning: the render was applied. No reason, status or timing of `Ready` depends on `Healthy`. `spec.dependsOn` SHALL keep waiting on `Ready`, and a TransformerRegistration SHALL keep activating on its provider's `Ready`.

#### Scenario: Every object healthy

- **WHEN** an instance's inventory holds a ConfigMap and a Deployment whose observed generation, updated, available and total replicas match its spec
- **THEN** `Healthy` is `True` with reason `RolledOut` and a message beginning "2/2 objects ready"

#### Scenario: A rollout in progress

- **WHEN** the Deployment's controller has not yet observed its current generation
- **THEN** `Healthy` is `False` with reason `NotRolledOut`, and the message names `Deployment <namespace>/<name> (NotReady)`
- **AND** `Ready` is `True` with reason `ReconciliationSucceeded`

#### Scenario: A stalled Deployment

- **WHEN** the Deployment has observed its generation and reports `Progressing` with reason `ProgressDeadlineExceeded`, and another object is not ready
- **THEN** `Healthy` is `False` with reason `ProgressDeadlineExceeded`, and the message names the stalled Deployment first, whatever its place in the inventory

#### Scenario: A missing object

- **WHEN** an inventory object does not exist on the cluster
- **THEN** `Healthy` is `False` with reason `NotRolledOut`, and the message names the object with status `Missing`

#### Scenario: A kind that is no longer served

- **WHEN** an inventory object's kind has no mapping on the API server, because its CustomResourceDefinition was deleted
- **THEN** the object counts as `Missing`, and `Healthy` is `False` with reason `NotRolledOut`

#### Scenario: An unreadable object

- **WHEN** every readable object is healthy and one object's read fails with an error other than not found
- **THEN** `Healthy` is `Unknown` with reason `HealthUnknown`, and the message names the error

#### Scenario: An empty inventory

- **WHEN** the inventory holds no objects
- **THEN** `Healthy` is `Unknown` with reason `HealthUnknown`

#### Scenario: Ready and its readers are unchanged

- **WHEN** a ModulePackage is `Ready=True` and `Healthy=False`
- **THEN** a ModulePackage whose `spec.dependsOn` names it proceeds
- **AND** a TransformerRegistration provided by a `Ready=True`, `Healthy=False` ModuleInstance activates as before

### Requirement: An expired Job does not make a ModuleInstance unhealthy

A `batch/v1` Job that a ModuleInstance renders with `ttlSecondsAfterFinished` and that the cluster removed SHALL NOT make `Healthy` `False`. The reconcile that finds it absent with unchanged digests removes it from `status.inventory` (`drift-detection`, "A missing object is restored"), and `Healthy` SHALL be judged over the entries that remain. The ready count and the total in the message SHALL count those entries only. A reconcile that would skip its render SHALL NOT report an absent Job as `Missing` before a render classified it, unless the last drift detection failed (`render-input-key`).

The controller does not know how the Job ended. A Job with a TTL that failed makes `Healthy` `False` while it exists; once the cluster removed it, it is an expired Job like any other and no longer counts. When the removal leaves `status.inventory` with no entries, `Healthy` SHALL be `Unknown` with reason `HealthUnknown`, as for any empty inventory.

A Job without a TTL that does not exist, and every other missing object, SHALL be judged as before. This requirement does not apply to a ModulePackage.

#### Scenario: An instance with an expired Job stays healthy

- **GIVEN** a ModuleInstance with `Healthy=True` whose inventory holds a ConfigMap and a Job with a TTL that completed
- **WHEN** the cluster removes the Job and the instance is reconciled on its interval
- **THEN** `Healthy` stays `True` with reason `RolledOut` and the message counts the ConfigMap only
- **AND** the reconcile requeues on the instance reconcile interval, not on the health requeue

#### Scenario: The only inventory object was an expired Job

- **GIVEN** a ModuleInstance whose render is one Job with a TTL
- **WHEN** the cluster removes the Job and the controller reconciles with unchanged digests
- **THEN** `status.inventory.entries` is empty and `Healthy` is `Unknown` with reason `HealthUnknown`
- **AND** the reconcile requeues on the instance reconcile interval

#### Scenario: An absent Job after a failed drift detection

- **GIVEN** a ModuleInstance whose inventory Job does not exist
- **WHEN** drift detection fails on the reconcile that renders for it
- **THEN** the Job stays in `status.inventory`, and `Healthy` is `False` with reason `NotRolledOut` and names the Job as `Missing`
