## MODIFIED Requirements

### Requirement: An absent batch Job forces a render because only a render tells whether it is expired
A ModuleInstance or ModulePackage reconcile SHALL skip the render when, and only when, all of these hold, as the scenario "The skip conditions" lists them. The last condition exists because only a render tells whether an absent Job is expired or must be restored (`drift-detection`, "An expired Job leaves the inventory").

#### Scenario: An absent Job is judged by a render
- **WHEN** only a render tells whether an absent Job is expired or must be restored
- **THEN** the reconcile renders

### Requirement: A reconcile that renders for an absent batch Job judges health twice
A ModuleInstance or ModulePackage reconcile SHALL skip the render when, and only when, all of these hold, as the scenario "The skip conditions" lists them. The last condition exists because only a render tells whether an absent Job is expired or must be restored (`drift-detection`, "An expired Job leaves the inventory"). A reconcile that renders for it has judged health once without a patch; the render then judges again.

#### Scenario: A render for an absent Job judges health again
- **WHEN** a reconcile renders for an absent inventory Job
- **THEN** it has judged health once without a patch, and the render then judges again
