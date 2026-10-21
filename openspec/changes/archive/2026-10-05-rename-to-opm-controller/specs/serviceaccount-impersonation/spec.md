## MODIFIED Requirements

### Requirement: Controller may impersonate ServiceAccounts only
The controller's own ClusterRole MUST grant the `impersonate` verb on `serviceaccounts` and on no other resource. It MUST NOT grant `impersonate` on `users` or on `groups`. This MUST hold in every form the role ships in: the kustomize role, the install manifest and the controller module's rendered RBAC. With that role alone the apply path MUST still work.

#### Scenario: Shipped role names serviceaccounts only
- **WHEN** the rules of the shipped manager ClusterRole are read
- **THEN** the only resource with the `impersonate` verb is `serviceaccounts`

#### Scenario: Apply works under the shipped role
- **GIVEN** a controller identity that holds exactly the rules of the shipped manager ClusterRole
- **AND** a ServiceAccount in namespace `team-a` that a RoleBinding allows to manage ConfigMaps
- **WHEN** the controller applies a ConfigMap in `team-a` as that ServiceAccount
- **THEN** the apply succeeds

#### Scenario: Impersonating a user is refused
- **GIVEN** a controller identity that holds exactly the rules of the shipped manager ClusterRole
- **WHEN** it sends a request that impersonates the user `jane`, with or without a group
- **THEN** the apiserver refuses the request as forbidden

#### Scenario: Claiming a group is refused
- **GIVEN** a controller identity that holds exactly the rules of the shipped manager ClusterRole
- **WHEN** it sends a request that impersonates a ServiceAccount and names the group `system:masters`
- **THEN** the apiserver refuses the request as forbidden

### Requirement: Controller may list and watch ServiceAccounts

The controller's own ClusterRole MUST grant `list` and `watch` on `serviceaccounts`, beside `get` and `impersonate`, in every form the role ships in: the kustomize role, the install manifest and the controller module's rendered RBAC. The controller MUST watch ServiceAccount metadata only. The role MUST NOT grant a write verb on `serviceaccounts`.

#### Scenario: Shipped role lets the controller watch ServiceAccounts

- **WHEN** the rules of the shipped manager ClusterRole are read
- **THEN** the verbs on `serviceaccounts` are exactly `get`, `impersonate`, `list` and `watch`
