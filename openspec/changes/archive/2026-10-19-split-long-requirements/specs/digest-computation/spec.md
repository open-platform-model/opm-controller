## REMOVED Requirements

### Requirement: The render digest is the library's
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "The render digest is RenderDigest from the library's inventory package", "The render digest is computed from the exported JSON after the rendered CUE values are dropped", "The internal status package declares no render digest", "A failure to compute the render digest is reported as a render failure".

## ADDED Requirements

### Requirement: The render digest is RenderDigest from the library's inventory package
The render digest the reconcilers record (`lastAttemptedRenderDigest`, `lastAppliedRenderDigest` and history entries) MUST be `RenderDigest` from the library's `opm/k8s/inventory` over the render's one `object.Export`. Source: enhancement 0012 ("Labels are stamped by the CUE render; the shared render digest ignores the runtime-name label value").

#### Scenario: The runtime name does not count
- **WHEN** two renders differ only in the value of the `app.kubernetes.io/managed-by` label (`opm-controller` and `opm-cli`)
- **THEN** their render digests are equal

#### Scenario: The stored render digest changes once on upgrade
- **WHEN** an instance whose `lastAppliedRenderDigest` was written by an earlier operator renders the same objects under this operator
- **THEN** the new render digest differs from the stored value, and the reconcile is not a no-op

### Requirement: The render digest is computed from the exported JSON after the rendered CUE values are dropped
The render digest the reconcilers record (`lastAttemptedRenderDigest`, `lastAppliedRenderDigest` and history entries) MUST be `RenderDigest` from the library's `opm/k8s/inventory` over the render's one `object.Export`. It MUST be computed from the exported JSON after the rendered CUE values are dropped. Source: enhancement 0012 ("Labels are stamped by the CUE render; the shared render digest ignores the runtime-name label value").

#### Scenario: Order-independent
- **WHEN** the same resources are rendered in a different order
- **THEN** the recorded render digest is identical

#### Scenario: Content sensitivity
- **WHEN** a resource's name, spec or any label other than the managed-by label changes
- **THEN** the recorded render digest differs

### Requirement: The internal status package declares no render digest
The `internal/status` package MUST NOT declare a render digest of its own. Source: enhancement 0012 ("Labels are stamped by the CUE render; the shared render digest ignores the runtime-name label value").

#### Scenario: No render digest of the status package
- **WHEN** the `internal/status` package is read
- **THEN** it declares no render digest of its own

### Requirement: A failure to compute the render digest is reported as a render failure
The render digest the reconcilers record (`lastAttemptedRenderDigest`, `lastAppliedRenderDigest` and history entries) MUST be `RenderDigest` from the library's `opm/k8s/inventory` over the render's one `object.Export`. A failure to compute it MUST be reported as a render failure (`RenderFailed`) with the step "computing render digest". Source: enhancement 0012 ("Labels are stamped by the CUE render; the shared render digest ignores the runtime-name label value").

#### Scenario: A digest failure is a render failure
- **WHEN** the render digest cannot be computed
- **THEN** it is reported as a render failure (`RenderFailed`) with the step "computing render digest"
