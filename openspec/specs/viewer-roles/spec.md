# viewer-roles Specification

## Purpose

Define the read-only ClusterRoles the controller ships for a cluster administrator to bind, so a non-admin can read OPM state without a hand-written role: what each grants, that none is bound by default, and that none aggregates into a built-in role.

## Requirements

### Requirement: The controller ships a viewer role for every OPM kind

The controller's install manifest SHALL carry one viewer ClusterRole per OPM kind it serves: `opm-controller-moduleinstance-viewer-role`, `opm-controller-platform-viewer-role`, `opm-controller-modulepackage-viewer-role` and `opm-controller-transformerregistration-viewer-role`. A cluster administrator SHALL be able to grant a non-admin read access to Platforms, ModulePackages and TransformerRegistrations by binding these roles and nothing else.

#### Scenario: A bound subject reads every Platform kind

- **WHEN** a cluster administrator binds `opm-controller-platform-viewer-role` and `opm-controller-transformerregistration-viewer-role` to a user with a ClusterRoleBinding, and `opm-controller-modulepackage-viewer-role` with a RoleBinding in namespace `shop`
- **THEN** the API server allows that user to get, list and watch Platforms and TransformerRegistrations cluster-wide and ModulePackages in `shop`, and to get their `status` subresources

#### Scenario: A namespace binding does not reach other namespaces

- **WHEN** `opm-controller-modulepackage-viewer-role` is bound to a user only by a RoleBinding in namespace `shop`
- **THEN** the API server refuses that user a list of ModulePackages in any other namespace

#### Scenario: An unbound user reads nothing

- **WHEN** a user holds no binding to any of the viewer roles
- **THEN** the API server refuses that user get and list on Platforms, ModulePackages and TransformerRegistrations

### Requirement: Viewer roles grant reads only

Each viewer role SHALL grant only `get`, `list` and `watch` on its own kind and `get` on that kind's `status` subresource, in the `opmodel.dev` API group. No viewer role SHALL grant any verb that creates, changes or deletes an object, any access to another kind or subresource, or any wildcard.

#### Scenario: A viewer cannot change what it reads

- **WHEN** a user bound to the three viewer roles of Platforms, ModulePackages and TransformerRegistrations attempts to create, update, patch or delete an object of those kinds, or to update its status
- **THEN** the API server refuses each request as forbidden

### Requirement: Viewer roles ship unbound

The controller SHALL ship no RoleBinding or ClusterRoleBinding naming a viewer role. Binding one is the cluster administrator's decision.

#### Scenario: A fresh install grants no reads

- **WHEN** the controller's install manifest is applied to a cluster
- **THEN** no binding in the manifest names `opm-controller-platform-viewer-role`, `opm-controller-modulepackage-viewer-role` or `opm-controller-transformerregistration-viewer-role`

### Requirement: Viewer roles do not aggregate into built-in roles

No viewer role SHALL carry an `rbac.authorization.k8s.io/aggregate-to-view`, `aggregate-to-edit` or `aggregate-to-admin` label, so binding a built-in role grants no read of OPM kinds through them. Whether they aggregate is an open question; a later change MAY add the label once that question is answered.

#### Scenario: The built-in view role is unchanged

- **WHEN** the controller's install manifest is applied and a user is bound only to the built-in `view` ClusterRole
- **THEN** the API server refuses that user a list of Platforms, ModulePackages and TransformerRegistrations
