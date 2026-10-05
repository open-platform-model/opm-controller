## MODIFIED Requirements

### Requirement: The render digest is the library's
The render digest the reconcilers record (`lastAttemptedRenderDigest`, `lastAppliedRenderDigest` and history entries) MUST be `RenderDigest` from the library's `opm/k8s/inventory` over the render's one `object.Export`. It MUST be computed from the exported JSON after the rendered CUE values are dropped. The `internal/status` package MUST NOT declare a render digest of its own. A failure to compute it MUST be reported as a render failure (`RenderFailed`) with the step "computing render digest". Source: 0012:D6:R2/R3.

#### Scenario: Order-independent
- **WHEN** the same resources are rendered in a different order
- **THEN** the recorded render digest is identical

#### Scenario: Content sensitivity
- **WHEN** a resource's name, spec or any label other than the managed-by label changes
- **THEN** the recorded render digest differs

#### Scenario: The runtime name does not count
- **WHEN** two renders differ only in the value of the `app.kubernetes.io/managed-by` label (`opm-controller` and `opm-cli`)
- **THEN** their render digests are equal

#### Scenario: The stored render digest changes once on upgrade
- **WHEN** an instance whose `lastAppliedRenderDigest` was written by an earlier controller renders the same objects under this controller
- **THEN** the new render digest differs from the stored value, and the reconcile is not a no-op
