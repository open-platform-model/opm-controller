## MODIFIED Requirements

### Requirement: The inventory digest is the library's canonical digest
`status.inventory.digest` MUST be `inventory.Digest` from `opm/k8s/inventory` over the inventory's entries. That digest hashes a versioned canonical encoding of each entry's field values, so it depends neither on the order of the entries nor on the JSON tags of the CRD entry type. Source: 0012:D7:R2/R3.

#### Scenario: Order does not matter
- **WHEN** the same entries are digested in a different order
- **THEN** the digest is identical

#### Scenario: Content does
- **WHEN** an entry's `Name` changes
- **THEN** the digest differs

#### Scenario: The stored digest changes once on upgrade
- **WHEN** an instance whose `status.inventory.digest` was written by an earlier controller is applied by this controller with the same entries
- **THEN** the new `status.inventory.digest` equals `inventory.Digest` of those entries and differs from the stored value
