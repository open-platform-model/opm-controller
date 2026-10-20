## MODIFIED Requirements

### Requirement: The inventory digest is inventory.Digest over the rendered set the last apply recorded
`status.inventory.digest` MUST be `inventory.Digest` from `opm/k8s/inventory` over the entries of the rendered set that the last apply recorded. It equals the digest of `status.inventory.entries`, with one exception: after a ModuleInstance reconcile removed the entry of an expired Job (`drift-detection`, "An expired Job leaves the inventory"), the entries are fewer and the digest still names the rendered set. Source: enhancement 0012.

#### Scenario: The stored digest changes once on upgrade
- **WHEN** an instance whose `status.inventory.digest` was written by an earlier operator is applied by this operator with the same entries
- **THEN** the new `status.inventory.digest` equals `inventory.Digest` of those entries and differs from the stored value

#### Scenario: The digest after an expired Job left the entries
- **WHEN** a ModuleInstance reconcile with unchanged digests removes an expired Job from `status.inventory.entries`
- **THEN** `status.inventory.digest` keeps its value, which is `inventory.Digest` of the rendered entries, the Job included
