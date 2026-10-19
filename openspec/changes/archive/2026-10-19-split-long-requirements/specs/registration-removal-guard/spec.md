## MODIFIED Requirements

### Requirement: Shrinking provides while dependents exist is refused before it takes effect
The operator SHALL refuse to apply a rendered `TransformerRegistration` that removes a contract from `spec.provides` while instances still demand that contract, and the refusal SHALL take effect while the previously accepted claim is still the stored one. A refusal that arrives after the new spec has replaced the accepted claim SHALL NOT be considered to satisfy this requirement: by then the old claim is gone and dependents are broken regardless of the reason reported.

#### Scenario: A shrinking upgrade is refused

- **WHEN** a provider module upgrade re-renders its claim with a contract dropped from `provides`, and instances still demand that contract
- **THEN** the claim stored in the cluster is unchanged, and the refusal names the dropped contracts and the dependent count

#### Scenario: The previously accepted claim keeps serving

- **WHEN** such an upgrade is refused
- **THEN** the claim that was accepted before the upgrade is still the effective one, and the instances demanding its contracts continue to render

#### Scenario: A non-shrinking upgrade passes untouched

- **WHEN** a provider upgrade re-renders its claim with the same contract set at a new catalog version
- **THEN** this requirement does not refuse it, and the new claim is applied

#### Scenario: Dropping a contract nobody demands is allowed

- **WHEN** an upgrade removes a contract from `provides` that no instance demands
- **THEN** it is applied

#### Scenario: A provider that renders no claim is unaffected

- **WHEN** a module that renders no `TransformerRegistration` is upgraded
- **THEN** this requirement refuses nothing

## REMOVED Requirements

### Requirement: A refused upgrade is reported on the provider it came from
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "A refused shrink is reported on the ModuleInstance whose render produced the claim" and "A refused shrink does not write the claim's own status".

## ADDED Requirements

### Requirement: The refusal of a shrink names the dropped contracts and the dependent count
As in requirement "Shrinking provides while dependents exist is refused before it takes effect": The refusal SHALL name the dropped contracts and the dependent count, matching the blocked delete's shape, since both are the same act through different doors.

#### Scenario: A refusal names the dropped contracts and the dependent count
- **WHEN** a rendered `TransformerRegistration` removes a contract from `spec.provides` while instances still demand that contract
- **THEN** the refusal names the dropped contracts and the dependent count, matching the blocked delete's shape

### Requirement: The guarantee against shrinking is bounded to writes the operator makes
As in requirement "Shrinking provides while dependents exist is refused before it takes effect": This guarantee is bounded to writes the operator makes. A `TransformerRegistration` written by another actor, which requires platform-admin RBAC the operator ships unbound, is outside it, and the operator SHALL NOT report a guarantee it does not hold.

#### Scenario: A claim written by another actor is outside the guarantee
- **WHEN** a `TransformerRegistration` is written by another actor
- **THEN** it is outside the guarantee, and the operator does not report a guarantee it does not hold

### Requirement: A refused shrink is reported on the ModuleInstance whose render produced the claim
The operator SHALL report a refused shrink on the `ModuleInstance` whose render produced the claim, naming the claim, the dropped contracts and the dependent count. The instance SHALL NOT report ready while one of its rendered resources is being withheld, so a refused upgrade is visible as an unconverged instance rather than as a silent partial apply.

#### Scenario: The instance names what was refused

- **WHEN** an upgrade is refused
- **THEN** the provider's `ModuleInstance` reports not ready, naming the claim, the dropped contracts and the dependent count

#### Scenario: The refusal clears when the cause does

- **WHEN** the instances demanding a dropped contract stop demanding it, and the provider is reconciled again
- **THEN** the upgrade applies and the instance converges, with no operator action on the claim

### Requirement: A refused shrink does not write the claim's own status
As in requirement "A refused shrink is reported on the ModuleInstance whose render produced the claim": The claim's own status SHALL NOT be written by this refusal. Acceptance owns the claim's conditions, and a claim whose stored spec was never replaced has nothing new to report.

#### Scenario: The claim's own verdict is untouched

- **WHEN** an upgrade is refused
- **THEN** the claim's conditions still describe the accepted claim, unchanged by the refusal
