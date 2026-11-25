# Proposal

## Why

A claim or a subscription may spell a catalog version as `0.1.0` or `v0.1.0`. Both name the same build, and the library already writes the bare form into the generated platform module. The controller keeps the spelling as written in two places. `Platform.status.registry[].version` then disagrees with the module and with its own API documentation, which says "bare SemVer" (issue opm-controller#234). `status.packageIdentity` also moves on a spelling-only claim edit, so the controller rebuilds the platform and every ModuleInstance re-renders for a build that did not change.

## What Changes

- The Platform reconciler strips a leading `v` from each claim's and each subscription's version before it builds the claim coordinates of the package identity and the entries of `status.registry`.
- `status.registry[].version` always shows the bare SemVer spelling, whatever the source spelled.
- A claim edit that changes only the spelling of its version keeps `status.packageIdentity`, so it causes no regeneration and no fleet re-render.
- The doc comments of `Subscription.version` and `ResolvedRegistryEntry.version` say that either spelling is accepted and that status shows the bare form. The CRD descriptions and the controller module data are regenerated from them.

Release class: PATCH after GA (`fix(controller)`). Beta ships it as the next `-beta.N`. The API shape does not change.

Upgrade cost: on the first reconcile after the upgrade, each platform with an active v-prefixed claim gets a new package identity once. That causes one rebuild and one re-render of its ModuleInstances. Platforms whose claims are all bare keep their identity.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `registration-driven-regeneration`: adds a requirement that every catalog version is read in its bare spelling, for the package identity and for the status union.
- `platform-crd`: a subscription version is accepted in either spelling and names the same build.
- `live-flux-e2e`: the v-prefixed claim step expects the bare-version package identity, and the identity requirement says that a spelling-only edit does not move it.

## Impact

- Reconciler: `PlatformReconciler` in `internal/controller/platform_controller.go` (claim coordinates, registry entries).
- API types: doc comments on `Subscription` and `ResolvedRegistryEntry` in `api/v1alpha1/platform_types.go`; generated CRD in `config/crd/bases/` and the module data in `modules/opm_controller/` (CRD exception to the module-only PR rule).
- Tests: `internal/controller/platform_claims_test.go` (envtest), `test/e2e/registration_test.go` (the v-prefixed claim step).
- Consumers: every ModuleInstance render key hashes `status.packageIdentity`; see the upgrade cost above.

## Enhancement

None.
