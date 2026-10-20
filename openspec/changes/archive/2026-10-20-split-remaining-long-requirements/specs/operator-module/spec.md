## REMOVED Requirements

### Requirement: The module names the one operator release it deploys
**Reason**: Split into two requirements that each state one behavior.
**Migration**: See "The module carries its own version and names one operator image" and "Every rendered controller container runs the image the module names".

### Requirement: Every object renders through a catalog resource
**Reason**: Split into two requirements that each state one behavior.
**Migration**: See "Operator objects render through the catalog resource made for their kind" and "Administrator ClusterRoles render unbound through the raw-objects resource".

### Requirement: Names stay those of the earlier manifest
**Reason**: Split into four requirements that each state one behavior.
**Migration**: See "Fixed names: namespaced objects", "Fixed names: first ClusterRoles", "Fixed names: other ClusterRoles" and "Fixed names: CRDs".

### Requirement: The operator's tuning is typed instance values
**Reason**: Split into three requirements that each state one behavior.
**Migration**: See "The operator's tuning values are typed in #config and rendered into the Deployment", "Defaults reproduce the earlier manifest's Deployment and each resource quantity defaults on its own" and "The memory limit always resolves and sets GOMEMLIMIT".

## ADDED Requirements

### Requirement: The module carries its own version and names one operator image
The repository SHALL carry the operator's OPM module at the path `opmodel.dev/modules/opm_operator` on the v0 major, with a version of its own that is independent of the operator binary's. The module's source SHALL name the operator image it deploys by a published operator release's version tag and content digest, and SHALL expose that operator version as a value readable from the module's source without rendering it and without a cluster.

#### Scenario: The operator version is read without a render

- **WHEN** a tool evaluates the module's operator-version value from the source tree, with no cluster and no instance
- **THEN** it reads one operator release version, such as `1.0.0-beta.6`, and the tag and digest the module renders the image with

### Requirement: Every rendered controller container runs the image the module names
As in requirement "The module carries its own version and names one operator image": Every rendered controller container SHALL run exactly that image, so one module version deploys exactly one operator version.

#### Scenario: The rendered image matches the named release

- **WHEN** the module is rendered with default values
- **THEN** the controller container's image is `ghcr.io/open-platform-model/opm-operator:v<operator version>@sha256:<digest>` with the version and digest the module names

### Requirement: Operator objects render through the catalog resource made for their kind
The module SHALL render the operator's CRDs, its Namespace, the controller Deployment, its ServiceAccount, the metrics Service, and every Role, ClusterRole, RoleBinding and ClusterRoleBinding through the first-party catalog resource made for that kind, apart from the exception below. The module SHALL render the operator's Namespace itself, so an install records it as an object of the instance. (Requirement: Administrator ClusterRoles render unbound through the raw-objects resource)

#### Scenario: The render contains the whole install shape

- **WHEN** the module is rendered with default values for the instance `opm-operator` in `opm-operator-system`
- **THEN** the result holds 22 objects: 4 CRDs, 1 Namespace, 1 ServiceAccount, 1 Role, 10 ClusterRoles, 1 RoleBinding, 2 ClusterRoleBindings, 1 Service and 1 Deployment

### Requirement: Administrator ClusterRoles render unbound through the raw-objects resource
The administrator ClusterRoles the operator ships for users to bind, its viewer roles among them, SHALL render with no binding. Those administrator ClusterRoles SHALL render through the catalog's raw-objects resource, in one component that holds them and nothing else, with their rules taken from the generated RBAC data; no other object SHALL use the raw-objects resource.

#### Scenario: Raw objects only for the administrator roles

- **WHEN** the module's components are listed
- **THEN** exactly one of them carries the catalog's raw-objects resource, and it renders exactly `opm-operator-metrics-reader`, the three `opm-operator-moduleinstance-*-role` ClusterRoles, `opm-operator-platform-viewer-role`, `opm-operator-modulepackage-viewer-role`, `opm-operator-transformerregistration-viewer-role` and `opm-operator-transformerregistration-admin-role`

#### Scenario: Administrator roles render unbound

- **WHEN** the module is rendered
- **THEN** `opm-operator-metrics-reader`, the three `opm-operator-moduleinstance-*-role` ClusterRoles, `opm-operator-platform-viewer-role`, `opm-operator-modulepackage-viewer-role`, `opm-operator-transformerregistration-viewer-role` and `opm-operator-transformerregistration-admin-role` render with no binding naming them

### Requirement: Fixed names: namespaced objects
The module SHALL render the Namespace `opm-operator-system`, the Deployment `opm-operator-controller-manager`, the ServiceAccount `opm-operator-controller-manager`, the Service `opm-operator-controller-manager-metrics-service`, the Role `opm-operator-leader-election-role`, the names the operator's own install manifest gives them. (Requirement: Fixed names: first ClusterRoles) (Requirement: Fixed names: other ClusterRoles) (Requirement: Fixed names: CRDs)

#### Scenario: The fixed names render

- **WHEN** the module is rendered with default values
- **THEN** every name above appears with its kind, and every namespaced object is in `opm-operator-system`

### Requirement: Fixed names: first ClusterRoles
The module SHALL render the ClusterRoles `opm-operator-manager-role`, `opm-operator-metrics-auth-role`, `opm-operator-metrics-reader`, `opm-operator-moduleinstance-admin-role`, `opm-operator-moduleinstance-editor-role`, `opm-operator-moduleinstance-viewer-role`, the names the operator's own install manifest gives them. (Requirement: Fixed names: namespaced objects) (Requirement: Fixed names: other ClusterRoles) (Requirement: Fixed names: CRDs)

#### Scenario: The first ClusterRole names render
- **WHEN** the module is rendered with default values
- **THEN** the ClusterRoles `opm-operator-manager-role`, `opm-operator-metrics-auth-role`, `opm-operator-metrics-reader`, `opm-operator-moduleinstance-admin-role`, `opm-operator-moduleinstance-editor-role` and `opm-operator-moduleinstance-viewer-role` appear with their kind

### Requirement: Fixed names: other ClusterRoles
The module SHALL render the ClusterRoles `opm-operator-platform-viewer-role`, `opm-operator-modulepackage-viewer-role`, `opm-operator-transformerregistration-viewer-role` and `opm-operator-transformerregistration-admin-role`, the names the operator's own install manifest gives them. Only a role binding MAY take the name the catalog derives from its role. (Requirement: Fixed names: namespaced objects) (Requirement: Fixed names: first ClusterRoles) (Requirement: Fixed names: CRDs)

#### Scenario: Bindings take catalog names

- **WHEN** the module is rendered
- **THEN** the three bindings are named by the catalog from their roles, and each binds its role to the ServiceAccount `opm-operator-controller-manager` in `opm-operator-system`

#### Scenario: The other ClusterRole names render
- **WHEN** the module is rendered with default values
- **THEN** the ClusterRoles `opm-operator-platform-viewer-role`, `opm-operator-modulepackage-viewer-role`, `opm-operator-transformerregistration-viewer-role` and `opm-operator-transformerregistration-admin-role` appear with their kind

### Requirement: Fixed names: CRDs
The module SHALL render the CRDs `moduleinstances.opmodel.dev`, `modulepackages.opmodel.dev`, `platforms.opmodel.dev` and `transformerregistrations.opmodel.dev`, the names the operator's own install manifest gives them. (Requirement: Fixed names: namespaced objects) (Requirement: Fixed names: first ClusterRoles) (Requirement: Fixed names: other ClusterRoles)

#### Scenario: The CRD names render
- **WHEN** the module is rendered with default values
- **THEN** the CRDs `moduleinstances.opmodel.dev`, `modulepackages.opmodel.dev`, `platforms.opmodel.dev` and `transformerregistrations.opmodel.dev` appear with their kind

### Requirement: The operator's tuning values are typed in #config and rendered into the Deployment
The module's `#config` SHALL type these values and render each into the controller Deployment: the operator image's repository, the registry mapping the operator resolves modules through (rendered as `--registry`), the default service account the operator applies as (rendered as `--default-service-account`), the controller container's resources, the replica count, and additional controller arguments appended after every typed argument. Each unset optional value SHALL leave its argument out.

#### Scenario: Every value reaches the Deployment

- **WHEN** the module is rendered with a repository, a registry mapping, a default service account, resources, two replicas and two extra arguments
- **THEN** the Deployment runs two replicas of the image from that repository, the container carries `--registry=<mapping>`, `--default-service-account=<name>` and the two extra arguments after them, and the given resources

#### Scenario: A CPU string the catalog cannot render is refused at the value

- **WHEN** the module is rendered with a CPU request or limit given as a string of cores, such as `"4"` or `"0.5"`
- **THEN** the module refuses the value, naming the field and the accepted forms (cores as a number, millicores as `"<n>m"`), instead of failing later in the catalog's Deployment transformer

### Requirement: Defaults reproduce the earlier manifest's Deployment and each resource quantity defaults on its own
As in requirement "The operator's tuning values are typed in #config and rendered into the Deployment": Defaults SHALL reproduce the earlier manifest's Deployment: one replica, requests `100m` CPU and `256Mi` memory, limits `2` CPU and `4Gi` memory. Each of those four resource quantities SHALL default on its own, so a value that sets one of them keeps the defaults of the others.

#### Scenario: One resources field keeps the other defaults

- **WHEN** the module is rendered with only the memory limit set, to `8Gi`
- **THEN** the container keeps the default CPU limit `2` and the default requests `100m` CPU and `256Mi` memory, and its `GOMEMLIMIT` is `6553MiB`

### Requirement: The memory limit always resolves and sets GOMEMLIMIT
As in requirement "The operator's tuning values are typed in #config and rendered into the Deployment": The memory limit SHALL always resolve to a value, default or given, so the container's `GOMEMLIMIT` is always defined; it SHALL be derived from that memory limit, not set as a value.

#### Scenario: The Go memory limit follows the memory limit

- **WHEN** the module is rendered with a memory limit
- **THEN** the container's `GOMEMLIMIT` is the floor of 80 percent of that limit in MiB (`3276MiB` for `4Gi`), so the soft limit never exceeds the hard one

#### Scenario: An unsupported memory unit is refused

- **WHEN** the module is rendered with a memory limit given as a plain number of bytes, such as `4294967296`, in decimal units, such as `"4G"`, or as a fraction, such as `"1.5Gi"`
- **THEN** the render fails with a message naming the accepted forms, `<n>Mi` and `<n>Gi`

### Requirement: A test compares the module's pod and Service with the kustomize build on every pull request
As in requirement "The controller's pod and Service stay those of the kustomize tree" (While `config/manager` and `config/default` still produce the operator's install manifest, the set of objects the module renders, the controller Deployment's spec and pod spec, and the metrics Service): A test SHALL compare them on every pull request, so the two sources of the install shape cannot drift apart before the module becomes the only one.

#### Scenario: A pull request is compared with the kustomize build
- **WHEN** a pull request is opened
- **THEN** a test compares the set of objects the module renders, the controller Deployment's spec and pod spec, and the metrics Service with those of a kustomize build of `config/default`

## MODIFIED Requirements

### Requirement: The controller's pod and Service stay those of the kustomize tree
While `config/manager` and `config/default` still produce the operator's install manifest, the set of objects the module renders, the controller Deployment's spec and pod spec, and the metrics Service SHALL equal those of a kustomize build of `config/default`, apart from the image, labels, the selector, the fields the catalog sets to Kubernetes API defaults, and the bindings' names (a binding is matched by its role and subjects).

#### Scenario: A manifest-only edit fails the test

- **WHEN** a pull request changes the manager container's arguments, probes, environment, volumes or security context in `config/manager` and not in the module
- **THEN** the module's render test fails naming the differing field

#### Scenario: An object added to the kustomize tree only fails the test

- **WHEN** a pull request adds an object to `config/default`, such as a binding in `config/rbac` that grants the controller's ServiceAccount an existing role, and not to the module
- **THEN** the module's render test fails naming the object, a binding by its role and subjects
