## MODIFIED Requirements

### Requirement: The finalizer-add reconcile requests an immediate requeue
The reconcile call that adds the finalizer MUST return `ctrl.Result{RequeueAfter: time.Second}` and MUST NOT set the deprecated `Requeue` field. The finalizer patch does not bump `metadata.generation`, so `predicate.GenerationChangedPredicate` filters its watch event. The delay is fixed, with no backoff; it equals the earlier `Requeue: true` delay only for a key with no earlier failures.

#### Scenario: Finalizer-add reconcile requests immediate requeue
- **GIVEN** a freshly-created ModuleInstance or ModulePackage without the `opmodel.dev/cleanup` finalizer
- **WHEN** the controller's first `Reconcile` call adds the finalizer
- **THEN** the reconcile returns a result with `RequeueAfter` of one second so the workqueue re-enqueues the request directly (bypassing the predicate)
- **AND** the next reconcile starts about one second later, not in the periodic resync window
