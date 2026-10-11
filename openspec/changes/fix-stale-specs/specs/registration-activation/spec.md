## REMOVED Requirements

### Requirement: Activation is reported and inert

**Reason**: "Inert" is false. The claim reconciler writes only the claim when it activates it, but the active-claim set is an input to the generated platform module, so activation changes what the Platform reconciler builds and what dependent instances render against. OpenSpec refuses to drop or rename a scenario inside a MODIFIED requirement, so the requirement is retired and added again under a new name.

**Migration**: See the requirement `Activation is reported on the claim and feeds the generated platform`, added below.

## ADDED Requirements

### Requirement: Activation is reported on the claim and feeds the generated platform

Activation SHALL be visible in the claim's status. The claim reconciler's activation SHALL write nothing but the claim itself: it SHALL NOT render a workload, change a `Platform` field or write another claim. Activation is not inert, though: an accepted and active claim is an input to the generated platform module, and the Platform reconciler regenerates the module from the set of such claims (capability `registration-driven-regeneration`).

#### Scenario: Activation writes nothing but the claim

- **WHEN** a claim becomes active
- **THEN** the claim reconciler writes the claim's own status and no other object in the cluster
- **AND** no workload is re-rendered by the claim reconciler as a result

#### Scenario: An active claim feeds the generated platform

- **WHEN** a claim is both accepted and active
- **THEN** the Platform reconciler includes the claim's catalog and version in the active-claim set it generates the platform module from
- **AND** a claim that is refused, or accepted but not active, is not in that set
