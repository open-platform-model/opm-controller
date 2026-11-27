# Proposal

## Why

The e2e spec "should converge correctly after a controller restart mid-reconcile" is a bare `Skip` in `test/e2e/concurrent_test.go`. Nothing in the suite proves that a controller which dies after it applied objects, but before it wrote the status, leaves a state that the next controller converges: no object recreated, none left outside the inventory, no prune lost. The code supports that claim (the inventory and the `lastApplied*` digests are written only in the deferred status commit; an allowed object outside the inventory is taken in; the stale set comes from the old inventory), but it has never run.

## What Changes

- Two envtest cases in `test/integration/reconcile/state_recovery_test.go` stage the crash state on the reconcile function: a first apply whose status was never written, and an update whose prune was pending.
- An Ordered e2e Describe in `test/e2e/concurrent_test.go` replaces the restart `Skip` with three specs on Kind: the staged first-apply crash, the staged pending prune, and a graceful controller pod delete while a reconcile runs. It holds the pod stop and restart fixture.
- The parallel-instances stub stays a `Skip`, as the one recorded future item.
- Spec `live-flux-e2e`: the Skip-stub requirement names only the parallel-instances scenario as a future item, and a new requirement states what the restart specs prove.

No product code changes. Release class: none. The change ships as `test` commits, which release nothing; after GA it would be no release either.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `live-flux-e2e`: the restart scenario leaves the list of recorded future items, and the live tier gains the restart convergence requirement.

## Impact

- API types and reconcilers: none changed. The tests exercise the ModuleInstance reconciler (`internal/reconcile`) and the apply guard and prune (`internal/apply`).
- Test files: `test/integration/reconcile/state_recovery_test.go`, `test/e2e/concurrent_test.go`.
- The e2e job gains about two minutes (one controller deploy, two lease waits of about 15s, three `hello` reconciles).
- Fixture: `hello` (`testing.opmodel.dev/modules/controller/hello@v0` v0.0.14), already published and readable on GHCR without a login. No fixture bump.

## Enhancement

None.
