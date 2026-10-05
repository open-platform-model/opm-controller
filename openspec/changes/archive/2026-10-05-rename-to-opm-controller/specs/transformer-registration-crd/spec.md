## MODIFIED Requirements

### Requirement: Cluster-scoped TransformerRegistration resource

The controller SHALL define a `TransformerRegistration` custom resource in group `opmodel.dev`, version `v1alpha1`, with `scope: Cluster` and a status subresource — `apiVersion: opmodel.dev/v1alpha1`, `kind: TransformerRegistration`. Nothing in a catalog can derive a CRD this repo owns, so the `catalog_opm` renderer SHALL emit these three literals exactly; the controller's group is authoritative. The types SHALL be registered in the runtime scheme.

#### Scenario: The kind is installable and cluster-scoped

- **WHEN** the CRD is installed
- **THEN** `TransformerRegistration` is registered with `scope: Cluster`
- **AND** a `TransformerRegistration` object carries no namespace

#### Scenario: A rendered claim is accepted by the API server

- **WHEN** an object carrying the group, version and kind the catalog renderer emits is applied, with every required field present
- **THEN** the API server accepts and stores it

### Requirement: Creating a claim requires platform-admin RBAC

The controller SHALL ship a ClusterRole granting create, update and delete on `transformerregistrations`, and SHALL NOT grant those verbs to any tenant-facing role it ships. A module applied under an impersonated tenant ServiceAccount therefore cannot create a claim unless a cluster administrator has deliberately bound that role. The controller's own manager role SHALL carry the read and status verbs it needs.

#### Scenario: A tenant ServiceAccount cannot create a claim

- **WHEN** a rendered `TransformerRegistration` is applied through the impersonated tenant ServiceAccount of a ModulePackage whose subject holds no platform-admin binding
- **THEN** the API server refuses the create as forbidden, and the refusal surfaces on the package as an impersonation failure

#### Scenario: A platform-admin subject may create a claim

- **WHEN** the same object is applied by a subject bound to the shipped platform-admin ClusterRole
- **THEN** the API server accepts the create
