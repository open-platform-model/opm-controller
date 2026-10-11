## REMOVED Requirements

### Requirement: Platform registered in the scheme without a reconciler

**Reason**: The requirement and its scenario `Types registered, no reconcile` are false. The operator registers a `PlatformReconciler` with the manager, and a `Platform` named `cluster` triggers a reconcile that generates and builds the platform module. OpenSpec refuses to drop or rename a scenario inside a MODIFIED requirement, so the requirement is retired and added again under a new name.

**Migration**: See the requirement `Platform is registered in the scheme and reconciled`, added below. What the reconciler does is specified by capability `platform-reconciler`.

## ADDED Requirements

### Requirement: Platform is registered in the scheme and reconciled

The operator SHALL register `Platform` and `PlatformList` with the runtime scheme so the types are serializable and installable. The manager SHALL also run a controller for `Platform` that reconciles the singleton named `cluster` (capability `platform-reconciler`), and it SHALL be woken by a change to the `Platform` spec and by a change to the set of accepted and active `TransformerRegistration` claims.

#### Scenario: Types are registered in the scheme

- **WHEN** the manager starts
- **THEN** `Platform` and `PlatformList` are registered in the scheme

#### Scenario: A Platform named cluster is reconciled

- **WHEN** the manager runs and a `Platform` named `cluster` is applied
- **THEN** the Platform controller reconciles it and records the outcome on `status.conditions`
