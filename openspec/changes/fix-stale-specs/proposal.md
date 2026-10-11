## Why

Three live specs say things the code does not do. Seeding the enhancements store found them (see the stale-spec rows of 2026-10-11 in the refactoring-enhancements swarm). The owner chose to fix the specs now: the specs follow the code, and the code does not change.

- `transformer-registration-crd` and `registration-activation` say an active claim is inert. It is not: the Platform reconciler reads every accepted and active claim and builds the platform module from it (`internal/controller/platform_controller.go`, `activeClaims`).
- `transformer-registration-crd`, the Go doc comment on `TransformerRegistrationSpec.Provides` and the generated CRD description say acceptance refuses an empty `provides` list on its merits. It does not. Acceptance compares the list with the re-derived contract set for exact equality (`providesDrift`), so an empty list passes for a catalog that implements nothing and fails with `ProvidesMismatch` for a catalog that implements something.
- `platform-crd` says `Platform` has no reconciler. It has one (`PlatformReconciler`, started in `cmd/main.go`).

## What Changes

- `transformer-registration-crd`: the requirement `Every claim field is required at admission` says what acceptance does with an empty list. The requirement `Claim status separates acceptance from activation` is retired and added again as `Claim status reports acceptance and activation as separate states`, with the false scenario replaced.
- `registration-activation`: the requirement `Activation is reported and inert` is retired and added again as `Activation is reported on the claim and feeds the generated platform`.
- `platform-crd`: the requirement `Platform registered in the scheme without a reconciler` is retired and added again as `Platform is registered in the scheme and reconciled`. The Purpose line of the spec is edited by hand, because a Purpose is not part of a delta.
- Go doc comment of `Provides` in `api/v1alpha1/transformerregistration_types.go`, and the copies generated from it (CRD, `dist/install.yaml`, the CUE copy), say the same as the spec.
- Two comments that repeat the false statements are corrected: the reconciler doc in `transformerregistration_controller.go` and the note in `test/integration/crdvalidation/transformerregistration_test.go`.

No code behaviour changes.

## Enhancement

None.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `transformer-registration-crd`: empty `provides` wording; the claim status requirement is replaced.
- `registration-activation`: the inert requirement is replaced.
- `platform-crd`: the no-reconciler requirement is replaced.

## Impact

- SemVer: none. Commits use `docs`. The CRD schema does not change, only its `description` text.
- Three requirements change their names. OpenSpec refuses a MODIFIED requirement that drops or renames a scenario, and the three false scenario names (`A stored claim is inert`, `An active claim changes nothing else`, `Types registered, no reconcile`) must go. Anything that links to the old requirement names (a docs page or a tool) needs the new name. A search of this repo found none.
- Stacks on the branch `docs/align-enhancement-store`; merge after its PR.
