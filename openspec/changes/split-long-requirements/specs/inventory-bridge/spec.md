## REMOVED Requirements

### Requirement: The inventory digest is the library's canonical digest
**Reason**: Split into requirements that each state one behavior.
**Migration**: See "The inventory digest is inventory.Digest over the rendered set the last apply recorded" and "A reader compares the stored digest with the digest of a render".

## ADDED Requirements

### Requirement: The inventory digest is inventory.Digest over the rendered set the last apply recorded
`status.inventory.digest` MUST be `inventory.Digest` from `opm/k8s/inventory` over the entries of the rendered set that the last apply recorded. It equals the digest of `status.inventory.entries`, with one exception: after a ModuleInstance reconcile removed the entry of an expired Job (`drift-detection`, "A missing object is restored"), the entries are fewer and the digest still names the rendered set. Source: enhancement 0012.

#### Scenario: The stored digest changes once on upgrade
- **WHEN** an instance whose `status.inventory.digest` was written by an earlier operator is applied by this operator with the same entries
- **THEN** the new `status.inventory.digest` equals `inventory.Digest` of those entries and differs from the stored value

#### Scenario: The digest after an expired Job left the entries
- **WHEN** a ModuleInstance reconcile with unchanged digests removes an expired Job from `status.inventory.entries`
- **THEN** `status.inventory.digest` keeps its value, which is `inventory.Digest` of the rendered entries, the Job included

### Requirement: A reader compares the stored digest with the digest of a render
As in requirement "The inventory digest is inventory.Digest over the rendered set the last apply recorded": A reader MUST compare the stored digest with the digest of a render, and MUST NOT recompute it from the stored entries to decide whether to apply. That digest hashes a versioned canonical encoding of each entry's field values, so it depends neither on the order of the entries nor on the JSON tags of the CRD entry type. Source: enhancement 0012.

#### Scenario: Order does not matter
- **WHEN** the same entries are digested in a different order
- **THEN** the digest is identical

#### Scenario: Content does
- **WHEN** an entry's `Name` changes
- **THEN** the digest differs

#### Scenario: A reader does not recompute the digest from the stored entries
- **WHEN** a reader decides whether to apply
- **THEN** it compares the stored digest with the digest of a render and does not recompute it from the stored entries
