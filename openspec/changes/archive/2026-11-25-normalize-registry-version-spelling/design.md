# Design

## Context

The Platform reconciler reads each catalog version twice from the spec values. `claimCoordinates` reads it from each active claim to build the package identity. `platformEntries` reads it from each subscription and each claim to build the generator's entries, and `resolvedRegistry` copies those entries into `status.registry`. Both read `spec.version` as written. The library's `platformmodule.Generate` already pins the canonical form (`v0.1.0`) in `cue.mod` and stamps the bare form (`0.1.0`) in the registry entry, so it accepts either spelling. The library's spelling helpers (`modversion.Bare`, `modversion.Canonical`) are in an `internal` package, so the controller cannot import them.

## Goals / Non-Goals

**Goals:**

- One spelling of a catalog version inside the controller: bare SemVer, the spelling the generated module stamps.
- A spelling-only claim edit does not move `status.packageIdentity`.

**Non-Goals:**

- No admission rule on either spelling. Both stay accepted.
- No change to how a claim is judged. The acceptance message keeps the spelling the claim wrote.
- No version validation in the controller. CUE's module system stays the validator.

## Decisions

### Normalise where the spec values are read

`claimCoordinates` and `platformEntries` MUST pass each version through one unexported helper before they use it:

```go
// bareVersion returns v without a leading "v": the spelling the generated
// platform module stamps. It validates nothing.
func bareVersion(v string) string { return strings.TrimPrefix(v, "v") }
```

So the identity, the store record, the prune keep set and `status.registry` all see the same value. `resolvedRegistry` needs no change, because it copies the entries.

Alternative: normalise inside `resolvedRegistry` and `NewPackageIdentity`. Rejected, because the generator would then still receive the raw spelling, and two rules would describe one value.

## Research & Decisions

### Fix direction

**Context**: The API doc says `status.registry[].version` is bare SemVer, but status records the spelling that was written. The identity moves on a spelling-only edit.
**Explored**: `internal/controller/platform_controller.go` (claim coordinates, entries, status union); the library `opm/helper/platformmodule/generate.go` (writes the bare form); `test/e2e/registration_test.go` (the v-prefixed claim step, which already names this issue).
**Decision**: Normalise to bare SemVer (owner). Changing the doc only was rejected, because the extra rebuilds would stay. An admission rule was rejected, because it breaks the API and the providers that write `v`.
**Rationale**: The cost is paid once at upgrade, instead of on every spelling-only edit, and status then matches the generated module.

### Scope

**Context**: `platformEntries` builds the status entries from both subscriptions and claims.
**Explored**: The Subscription doc and the `platform-crd` spec both say "bare SemVer", and nothing enforces it.
**Decision**: Normalise both sources (owner).
**Rationale**: The field doc covers every entry of `status.registry`, whatever its source.

## Reconcile phase impact

- Source, Apply, Prune: none.
- Render: a ModuleInstance render key hashes `status.packageIdentity`. On the first reconcile after the upgrade, a platform with an active v-prefixed claim gets a new identity once, so its instances re-render once. After that, a spelling-only claim edit causes no re-render.
- Status: `status.registry[].version` is bare for every entry.

## Risks / Trade-offs

- [One fleet re-render at upgrade for platforms with a v-prefixed active claim] → stated in the proposal as the release note; it happens once.
- [The new envtest case needs a reachable registry (`CUE_REGISTRY`), and it skips without one, like the other build cases in the same file] → the build runs it with `CUE_REGISTRY` exported, and the e2e step covers the identity on a cluster.
- [The doc comment edit regenerates the CRD and the controller module data in the same PR] → allowed by the repo's CRD exception; `task controller-module:drift` checks it.

## Slices

### S1: Bare version spelling in identity and status

- **Delivers**: the added requirement "A catalog version is read in its bare spelling" (both scenarios), and the reworded subscription requirement of `platform-crd`.
- **Seam**: an envtest reconcile case in `internal/controller/platform_claims_test.go` that reads `status.registry` and `status.packageIdentity` from the API, plus the flipped assertion in the e2e v-prefixed claim step in `test/e2e/registration_test.go`.
- **Depends on**: none.
- **Done when**: with a subscription at `v<fixture version>` and an active claim at `v<claim version>`, one reconcile writes both versions bare to `status.registry`, and `status.packageIdentity` equals the identity that the same tuple written bare produces; the e2e step asserts `identityV == identityBare`.
