## MODIFIED Requirements

### Requirement: The finalizer-add reconcile requests an immediate requeue
The reconcile call that adds the finalizer MUST request a requeue that does not wait for the periodic resync, with `ctrl.Result{RequeueAfter: time.Second}`, rather than relying on the watch event produced by the finalizer patch: that watch event is filtered by the controller's `predicate.GenerationChangedPredicate`, because finalizer patches modify `metadata.finalizers` but do not bump `metadata.generation`. The delay is a fixed one second with no backoff. It equals the delay of the earlier `Requeue: true` only for a key with no earlier failures. The result MUST NOT set the deprecated `Requeue` field.

#### Scenario: Finalizer-add reconcile requests immediate requeue
- **GIVEN** a freshly-created ModuleRelease without the `releases.opmodel.dev/cleanup` finalizer
- **WHEN** the controller's first `Reconcile` call adds the finalizer
- **THEN** the reconcile returns a result with `RequeueAfter` of one second so the workqueue re-enqueues the request directly (bypassing the predicate)
- **AND** the next reconcile starts about one second later, not in the periodic resync window
