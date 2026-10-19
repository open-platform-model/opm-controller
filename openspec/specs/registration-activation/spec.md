## Purpose

Define when an accepted `TransformerRegistration` becomes active, what keeps it active, and what does not. Activation is the second of the three states enhancement 0015 ("A transformer registration is a cluster-scoped CR, gated by the RBAC the operator already enforces") defines; acceptance is the first and is a separate capability.

## Requirements

### Requirement: An accepted claim activates when its provider is ready

An accepted claim SHALL become active when the `ModuleInstance` named by `spec.providerRef` reports `Ready=True`, and SHALL remain inactive until then. The operator SHALL watch provider readiness so activation does not wait for an unrelated reconcile. A claim that is not accepted SHALL NOT activate, whatever its provider reports.

#### Scenario: A claim activates when its provider becomes ready

- **WHEN** an accepted claim's provider `ModuleInstance` transitions to `Ready=True`
- **THEN** the claim reports `status.active` true without waiting for its own spec to change

#### Scenario: An accepted claim whose provider is not ready stays inactive

- **WHEN** a claim is accepted and its provider is not `Ready=True`
- **THEN** `status.active` is false and a condition says the claim is waiting on its provider

#### Scenario: A refused claim never activates

- **WHEN** a refused claim's provider reports `Ready=True`
- **THEN** `status.active` stays false

#### Scenario: A refused claim reports no waiting-on-provider state

- **WHEN** a claim that was accepted and waiting on its provider is later refused
- **THEN** it no longer reports that it is waiting on its provider, because the refusal is why it is inactive
- **AND** a claim refused while already active keeps reporting that it is active

### Requirement: Activation is reported on the claim and feeds the generated platform

Activation SHALL be visible in the claim's status. The claim reconciler's activation SHALL write nothing but the claim itself: it SHALL NOT render a workload, change a `Platform` field or write another claim. Activation is not inert, though: an accepted and active claim is an input to the generated platform module, and the Platform reconciler regenerates the module from the set of such claims (capability `registration-driven-regeneration`).

#### Scenario: Activation writes nothing but the claim

- **WHEN** a claim becomes active
- **THEN** the claim reconciler writes the claim's own status and no other object in the cluster
- **AND** no workload is re-rendered by the claim reconciler as a result

#### Scenario: An active claim feeds the generated platform

- **WHEN** a claim is both accepted and active
- **THEN** the Platform reconciler includes the claim in the active-claim set it generates the platform module from
- **AND** a claim that is refused, or accepted but not active, is not in that set

### Requirement: An active claim stays active when its provider's readiness regresses
Once active, a claim SHALL NOT be deactivated because its provider's readiness regressed. This is deliberate and not an oversight: the active-claim set is an input to platform-package regeneration, so a flapping provider that toggled the set would regenerate the platform without the provider's catalog and fail every dependent instance's render for the duration.

#### Scenario: A provider going unready leaves the claim active

- **WHEN** an active claim's provider transitions from `Ready=True` to not ready
- **THEN** the claim stays active and its status does not flap

#### Scenario: A provider recovering does not re-run the gate

- **WHEN** that provider becomes ready again
- **THEN** the claim is still active and no activation transition is recorded, because it never left

### Requirement: The readiness gate governs initial activation only
As in requirement "An active claim stays active when its provider's readiness regresses": The readiness gate SHALL govern initial activation only. A claim leaves the active state by deletion, never by its provider becoming unhealthy. Established CRDs outlive a provider's pods, so an active claim staying active through a provider outage is correct.

#### Scenario: The readiness gate does not govern an active claim
- **WHEN** a claim is already active
- **THEN** the readiness gate does not govern it, because a claim leaves the active state by deletion, never by its provider becoming unhealthy
