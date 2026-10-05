## MODIFIED Requirements

### Requirement: Viewer roles ship unbound

The controller SHALL ship no RoleBinding or ClusterRoleBinding naming a viewer role. Binding one is the cluster administrator's decision.

#### Scenario: A fresh install grants no reads

- **WHEN** the controller's install manifest is applied to a cluster
- **THEN** no binding in the manifest names `opm-controller-platform-viewer-role`, `opm-controller-modulepackage-viewer-role` or `opm-controller-transformerregistration-viewer-role`

### Requirement: Viewer roles do not aggregate into built-in roles

No viewer role SHALL carry an `rbac.authorization.k8s.io/aggregate-to-view`, `aggregate-to-edit` or `aggregate-to-admin` label, so binding a built-in role grants no read of OPM kinds through them. Whether they aggregate is open (0030:OQ6); a later change MAY add the label once that question is answered.

#### Scenario: The built-in view role is unchanged

- **WHEN** the controller's install manifest is applied and a user is bound only to the built-in `view` ClusterRole
- **THEN** the API server refuses that user a list of Platforms, ModulePackages and TransformerRegistrations

## REMOVED Requirements

### Requirement: The operator ships a viewer role for every OPM kind

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-019); this requirement is restated under the controller's name.

**Migration**: See the requirement "The controller ships a viewer role for every OPM kind" in this capability.

## ADDED Requirements

### Requirement: The controller ships a viewer role for every OPM kind

The controller's install manifest SHALL carry one viewer ClusterRole per OPM kind it serves: `opm-controller-moduleinstance-viewer-role`, `opm-controller-platform-viewer-role`, `opm-controller-modulepackage-viewer-role` and `opm-controller-transformerregistration-viewer-role`. A cluster administrator SHALL be able to grant a non-admin read access to Platforms, ModulePackages and TransformerRegistrations by binding these roles and nothing else. Source: 0030:D11:R4.

#### Scenario: A bound subject reads every Platform kind

- **WHEN** a cluster administrator binds `opm-controller-platform-viewer-role` and `opm-controller-transformerregistration-viewer-role` to a user with a ClusterRoleBinding, and `opm-controller-modulepackage-viewer-role` with a RoleBinding in namespace `shop`
- **THEN** the API server allows that user to get, list and watch Platforms and TransformerRegistrations cluster-wide and ModulePackages in `shop`, and to get their `status` subresources

#### Scenario: A namespace binding does not reach other namespaces

- **WHEN** `opm-controller-modulepackage-viewer-role` is bound to a user only by a RoleBinding in namespace `shop`
- **THEN** the API server refuses that user a list of ModulePackages in any other namespace

#### Scenario: An unbound user reads nothing

- **WHEN** a user holds no binding to any of the viewer roles
- **THEN** the API server refuses that user get and list on Platforms, ModulePackages and TransformerRegistrations
