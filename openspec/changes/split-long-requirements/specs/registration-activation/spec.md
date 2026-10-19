## REMOVED Requirements

### Requirement: Activation latches against provider health
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "An active claim stays active when its provider's readiness regresses" and "The readiness gate governs initial activation only".

## ADDED Requirements

### Requirement: An active claim stays active when its provider's readiness regresses
Once active, a claim SHALL NOT be deactivated because its provider's readiness regressed. This is deliberate and not an oversight: the active-claim set is an input to platform-package regeneration, so a flapping provider that toggled the set would regenerate the platform without the provider's catalog and fail every dependent instance's render for the duration.

#### Scenario: A provider going unready leaves the claim active

- **WHEN** an active claim's provider transitions from `Ready=True` to not ready
- **THEN** the claim stays active and its status does not flap

### Requirement: The readiness gate governs initial activation only
As in requirement "An active claim stays active when its provider's readiness regresses": The readiness gate SHALL govern initial activation only. A claim leaves the active state by deletion, never by its provider becoming unhealthy. Established CRDs outlive a provider's pods, so an active claim staying active through a provider outage is correct.

#### Scenario: A provider recovering does not re-run the gate

- **WHEN** that provider becomes ready again
- **THEN** the claim is still active and no activation transition is recorded, because it never left
