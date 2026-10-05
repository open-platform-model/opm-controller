## MODIFIED Requirements

### Requirement: An accepted claim activates when its provider is ready

An accepted claim SHALL become active when the `ModuleInstance` named by `spec.providerRef` reports `Ready=True`, and SHALL remain inactive until then. The controller SHALL watch provider readiness so activation does not wait for an unrelated reconcile. A claim that is not accepted SHALL NOT activate, whatever its provider reports.

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
