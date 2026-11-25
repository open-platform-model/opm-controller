# Tasks

## 1. Bare version spelling in identity and status

- [x] 1.1 Add an envtest case to `internal/controller/platform_claims_test.go` (registry-backed context): a subscription and an active claim, both v-prefixed; assert bare versions in `status.registry` and an identity equal to the bare tuple's identity. Verify: the case fails before 1.3.
- [x] 1.2 In `test/e2e/registration_test.go`, the v-prefixed claim step asserts `identityV == identityBare` and reads `status.registry` version without `TrimPrefix`; update its comment. Verify: compiles under `go vet -tags e2e ./test/e2e/...`.
- [x] 1.3 In `internal/controller/platform_controller.go`, add the bare-version helper and use it in `claimCoordinates` and `platformEntries`. Verify: the case from 1.1 passes with `CUE_REGISTRY` exported.
- [x] 1.4 Reword the doc comments of `Subscription.Version` and `ResolvedRegistryEntry.Version` in `api/v1alpha1/platform_types.go` (either spelling accepted; status shows the bare form), then run `task dev:manifests dev:generate`. Verify: `task controller-module:drift` passes.
- [x] 1.5 `task dev:fmt dev:vet dev:lint dev:test` and `task docs:bundle:check` green; run `task dev:e2e` or state that it was skipped.
