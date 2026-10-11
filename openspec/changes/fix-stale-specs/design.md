## Context

The specs follow the code. Every changed statement below was read from the code first. The code and its tests stay as they are, except for two comments that repeat the false statements.

## Decisions

### Retire three requirements instead of modifying them

OpenSpec refuses a MODIFIED requirement that drops or renames a scenario (`openspec validate` reports the omitted scenario). The scenarios `A stored claim is inert`, `An active claim changes nothing else` and `Types registered, no reconcile` state something false in their names, so they cannot stay. Each owner requirement is REMOVED and added again under a new name. The scenarios that were true keep their names inside the new requirements (`A stored claim receives a verdict`). The requirement `Every claim field is required at admission` keeps its scenarios, so it is MODIFIED.

### What the code does, with the proof

- **Active claims feed the platform.** `PlatformReconciler.activeClaims` (`internal/controller/platform_controller.go:327-344`) keeps each claim with `Status.Accepted` and `Status.Active`. `Reconcile` (line 172-185) builds the package identity and the registry entries from that set. `SetupWithManager` (line 554-564) watches `TransformerRegistration`. The claim reconciler itself writes only the claim: the test `writes nothing but the claim itself when the claim activates` (`internal/controller/transformerregistration_activation_test.go:298-319`) proves it. So "inert" is true of the claim reconciler's own writes and false of the system.
- **Empty `provides`.** The CRD has no `MinItems` and the integration test `an empty provides list is accepted` (`test/integration/crdvalidation/transformerregistration_test.go:156-162`) creates one. At acceptance, `providesDrift` (`internal/controller/transformerregistration_controller.go:332-349`) sorts both lists and returns "accepted" when `slices.Equal` holds. An empty claim matches a catalog whose derived set is empty. Otherwise the reason is `ProvidesMismatch` (tests `refuses a claim omitting a contract the catalog does implement`, `internal/controller/transformerregistration_catalog_test.go:238-251`). No test pins the empty-and-empty case; the spec states it from `providesDrift`.
- **Platform has a reconciler.** `cmd/main.go:389-401` constructs `PlatformReconciler` and calls `SetupWithManager`. The reconciler spec is `platform-reconciler`.

### Purpose of platform-crd

A Purpose is not part of a delta, so the archive does not sync it. The Purpose line "without a reconciler" is edited by hand in the live spec in the same change.

### Generated copies

The `Provides` doc comment is the source of the CRD description, `dist/install.yaml` and `modules/opm_operator/zz_generated_crds.cue`. They are regenerated with `task dev:manifests dev:generate` and `task operator:installer`, never edited by hand. `task operator-module:drift` proves the CUE copy.

## Risks

- A requirement name that other text links to breaks. A search of the repo found no such link.
- The empty-and-empty acceptance is read from code, not from a test. If a later test shows otherwise, the spec is wrong and must change.
