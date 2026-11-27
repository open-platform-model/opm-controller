## MODIFIED Requirements

### Requirement: No permanent Skip stubs in the live tier

e2e spec files SHALL contain only executable specs or explicitly-recorded future items; stubs whose coverage exists at a lower tier and whose live-tier value is subsumed by this capability's specs SHALL be deleted, with the parallel-instances scenario remaining as the sole recorded future item.

#### Scenario: Stub census after landing

- **WHEN** the e2e suite is inspected after this change lands
- **THEN** `lifecycle_test.go`'s stubs SHALL be replaced by live specs, `prune_test.go`/`finalizer_test.go` stub files SHALL be gone, the controller-restart scenario SHALL be an executable spec, and only the parallel-instances scenario SHALL remain recorded

## ADDED Requirements

### Requirement: The live tier proves a controller restart converges

The live tier SHALL prove that a deployed controller that stops between its apply and its status commit leaves a state that the next controller converges. After the restart, the instance SHALL reach `Ready=True` with `status.observedGeneration` equal to its generation; each rendered object SHALL keep the UID it had before the stop; `status.inventory` SHALL list exactly the rendered objects; and no object that carries the instance's UUID label SHALL exist outside `status.inventory`.

#### Scenario: A crash after the first apply

- **GIVEN** a ModuleInstance whose rendered objects exist in the cluster, while its status holds no inventory, no `lastApplied*` digests and no `lastAppliedInputs`, as a controller that died before its status commit leaves it
- **WHEN** the controller starts again
- **THEN** the objects SHALL be taken into `status.inventory` with their UIDs unchanged, and the instance SHALL reach `Ready=True`

#### Scenario: A crash with a pending prune

- **GIVEN** a ModuleInstance with `spec.prune: true` whose status inventory lists an object that its current render no longer produces, with that object live and carrying the instance's ownership labels, and no `lastApplied*` digests
- **WHEN** the controller starts again
- **THEN** that object SHALL be deleted, it SHALL be absent from `status.inventory`, and the rendered objects SHALL keep their UIDs

#### Scenario: The controller pod is deleted during a reconcile

- **GIVEN** a Ready ModuleInstance whose values were just changed
- **WHEN** the controller pod is deleted while that reconcile runs, and the Deployment starts a new pod
- **THEN** the new controller SHALL apply the new values, and the instance SHALL meet the same checks: same UIDs, inventory equal to the render, no labeled object outside the inventory, `Ready=True` on the current generation
