# Build report: add-controller-restart-e2e

Base: main

Working branch: `add-controller-restart-e2e`. Both slices are tests only (Delivers starts with `Tests only:`), so no code writer ran and the test-lock check did not apply. Each slice has a tests commit and a commit that ticks its boxes.

## Slices and commits

- S1: Crash-window cases at the reconcile function
  - `f93d4e36` test(reconcile): cover crash recovery before the status commit (tests commit, `Tests-For: S1`; also carries the change folder)
  - `9846cff7` test(reconcile): close the crash recovery slice (`tasks.md`)
- S2: Controller restart specs on Kind
  - `cd19bd1b` test(e2e): cover a controller restart after a crash and a pod delete (tests commit, `Tests-For: S2`; also carries the owner's spec, design and task 2.5 edits)
  - `551d8292` test(e2e): close the controller restart slice (`tasks.md`)

## Guard tests

Every new test is a guard test, because the behaviour exists already.

- S1, "should converge after a crash between the first apply and the status commit": guards "A crash after the first apply". Break: in `internal/apply/guard.go` the take-in was removed. Result: `Ready` was `False` where `True` was expected.
- S1, "should prune a pending stale object after a crash": guards "A crash with a pending prune". Break: `staleEntries` in `internal/reconcile/converted.go` got an empty stale set. Result: `stale-object` was not deleted within 10s.
- S2, "should converge after a crash between the first apply and the status commit": guards "A crash after the first apply". Break: take-in refused in `internal/apply/guard.go`. Result: `ApplyRefused`, and no apply after the restart.
- S2, "should prune a pending stale object after a crash": guards "A crash with a pending prune". Break: `staleEntries` returns `nil`. Result: `hello-stale` was still there after 180s.
- S2, "should converge correctly after a controller restart mid-reconcile": guards "The controller pod is deleted during a reconcile". There was no break check. The delete may land after the apply, so no one-line break makes the spec fail on its window. A regression is caught by the ConfigMap message check (the new values are applied) and by the converged check (same UIDs, inventory equal to the render, no orphan, Ready on the current generation).

## Assumptions

- S1: How are the crash states made? By a reconcile to Ready, then a status merge patch through the status subresource. `instanceUUID`, the conditions and `observedGeneration` stay.
- S1: Which component does the stale inventory entry get? A copy of the rendered `test-module` entry with only the name changed, and `count: 2`.
- S1: Which objects count for "no labeled object outside the inventory"? ConfigMaps in the suite namespace with the UUID label, without a deletion timestamp.
- S1: In which order are the checks? The outcome checks (Ready, inventory, UID) run before the requeue check.
- S1: How is cleanup done? With the file's inline pattern (delete the ConfigMaps, `cleanupInstance`), not `DeferCleanup`.
- S2: Where does the expected render come from? From the inventory at the first Ready, checked to be one ConfigMap in `default` with `spec.prune: true`.
- S2: Which objects count for "no labeled object outside the inventory"? ConfigMaps in all namespaces with the UUID label, not being deleted.
- S2: What does the staged pending-prune inventory hold? `revision` and `digest` null, `count: 2`, and the rendered entry plus `hello-stale`.
- S2: How is `lastAppliedAt` compared, given its 1s resolution? It must not be before the reference time cut to the second. For the pod delete spec the reference is the values change.
- S2: What is the teardown order? The order of `registry_auth_test.go`. `opm-system` is not deleted.
- S2: How does the spec know the pod is new? It waits for one Ready controller pod whose name differs from the pod before the delete.

## Bugs, fix slices and leads

- Bugs fixed: none. Fix slices: none.
- Lead: "wait for the Platform Ready again" in the restart helper passes at once. The Platform keeps `Ready=True` in its status while the controller is stopped. The real wait is the converged check (`lastAppliedAt` after the staging time, within 3 minutes).
- Lead: the repo lint config sets no build tags, so `task dev:lint` does not lint the e2e files. `golangci-lint run --build-tags=e2e ./test/e2e/...` shows 10 findings that were there before, such as `modernize: plusbuild` on each `// +build e2e` line.

## Disputed tests

None.

## Base failures

None.

## Failed checks

- S1, before the tests commit: `task dev:test` failed in `dev:manifests:config` with "CRD for ModuleInstance.opmodel.dev has no storage version". Cause: two git worktrees under `.claude/worktrees/` inside the repo. A clean checkout of `main` passed that step. The owner moved both worktrees out of the repo, and then `task dev:test` passed.
- Every run of `task dev:lint`, and so of `task check:commit`: the plugin build of golangci-lint picked Go 1.26.5 and refused the Go 1.27.3 target. With `GOTOOLCHAIN=go1.27.3`, the Go version that `go.mod` names, the build works and lint gives `0 issues`. The check command ran with that variable at each commit.

## Stops

- S1, step (c), test command: the test command failed because of the worktrees (above). Owner's answer: the owner moves the worktrees. After that the build went on.
- S2, step (b), owner decision: a ModuleInstance reconcile never sets `Reconciling=True`, so the pod delete spec's 10s wait always ran out. Owner's answer: delete the pod right after the values change, with no wait, and say so in the text. The spec scenario, `design.md` and task 2.5 were edited, `openspec validate` passed, and a fresh test writer changed the spec.

## e2e result (task 2.6)

`task dev:e2e` on Kind, on the final test file: exit 0. 14 passed, 0 failed, 0 pending, 10 skipped. The skips existed before this change: the parallel-instances stub, and 9 podinfo specs that need `LOCAL_REGISTRY` or `CONTROLLER_DOCKER_CONFIG`. `concurrent_test.go` holds one `Skip`.

Project rules: `## Commit rules` and `## Release rules` (from `RULES.md`), `## Text rules` (from `AGENTS.md`).

Check command: `task check:commit`, run as `GOTOOLCHAIN=go1.27.3 task check:commit`: passed at both closing commits.

## Fix round

Base: main

Verify found two findings. `task dev:test` passed, `task dev:e2e` passed on Kind (14 passed, 0 failed, 10 skipped; the three restart specs passed), and the test lock found no offence.

### Findings

- The check command failed in its lint step: "the Go language version (go1.26) used to build golangci-lint is lower than the targeted Go version (1.27.3)". Route: Owner, as a base failure. A clean worktree of `main` fails the same way. Cause: this machine has Go 1.26.5 installed, and the linter is built outside the repo, so the build does not use the Go 1.27.3 that `go.mod` names. With `GOTOOLCHAIN=go1.27.3` the command passes. Outcome: the owner answered Keep it out, so the finding is closed and a separate change fixes the tooling.
- `design.md` held no RFC 2119 keyword, against `rules.design` in `openspec/config.yaml`. Route: Fix. Outcome: fixed. A fresh code writer stated the decisions that the design already made with `MUST` and `MUST NOT`. No decision changed, and the `## Slices` section is as it was. No other place of the diff breaks the same rule.

### Commits

- `6fbc7ddd` docs(openspec): state the restart e2e design with RFC 2119 keywords

### Questions to the owner

- Does the lint failure stay out of this change? Answer: A, keep it out.

### Assumptions of the code writer

- Which keyword? `MUST` and `MUST NOT` only, because the design's points are firm decisions.
- Do the two lowercase "may" stay? No. They became "can", so the design holds no lowercase keyword.
- Does the design get a line that names RFC 2119? No. That would add text beyond the finding.

### Leads

- The pod delete spec cannot tell which controller applied the new values. The old pod can finish before the delete lands. The design accepts this as a smoke check.
- `test/e2e/concurrent_test.go` copies three label strings that the library owns (`labels.ModuleInstanceUUID`, `labels.ManagedBy`, `labels.ManagedByController`). The name `ownerUUIDLabel` says "owner" where the glossary says Instance UUID.
- The deploy and teardown block of the restart Describe is the fifth copy in `test/e2e/`.
- `liveUIDs` handles API groups, but the fixture renders one ConfigMap.
- The envtest cases create, read and delete ConfigMaps inline where sibling helpers exist.
- `design.md` cites "test 6.2" of an archived change.
- The two closing commits of the build use the scopes `reconcile` and `e2e` but change only `tasks.md`.
- `CONSTITUTION.md` says a section must close with a commit task, and `openspec/config.yaml` says to write none. That text is outside this change.

Project rules: `## Commit rules` and `## Release rules` (from `RULES.md`), `## Text rules` (from `AGENTS.md`).

Check command: `task check:commit`, run as `GOTOOLCHAIN=go1.27.3 task check:commit`: passed.

### Answers

- Answer: base | 81681835af6fab7c73a8dd70dab221bc6b09eabe | keep | task dev:lint in task check:commit, golangci-lint built with go1.26 refuses the Go 1.27.3 target | 2026-11-27T17:54Z
