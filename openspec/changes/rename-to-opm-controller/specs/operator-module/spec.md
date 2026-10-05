## REMOVED Requirements

### Requirement: The module names the one operator release it deploys

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-019); this requirement is restated under the controller's name.

**Migration**: See the requirement "The module names the one controller release it deploys" in `controller-module`.

### Requirement: The module deploys only an operator that refuses its own instance

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-019); this requirement is restated under the controller's name.

**Migration**: See the requirement "The module deploys only a controller that refuses its own instance" in `controller-module`.

### Requirement: Every object renders through a catalog resource

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-019); this requirement is restated under the controller's name.

**Migration**: See the requirement "Every object renders through a catalog resource" in `controller-module`.

### Requirement: Names stay those of the earlier manifest

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-019); this requirement is restated under the controller's name.

**Migration**: See the requirement "The module renders the fixed names" in `controller-module`.

### Requirement: The module renders only for the operator's fixed instance coordinates

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-019); this requirement is restated under the controller's name.

**Migration**: See the requirement "The module renders only for the controller's fixed instance coordinates" in `controller-module`.

### Requirement: One Deployment selector across every module version

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-019); this requirement is restated under the controller's name.

**Migration**: See the requirement "One Deployment selector across every module version" in `controller-module`.

### Requirement: The CRDs are the controller's generated CRDs

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-019); this requirement is restated under the controller's name.

**Migration**: See the requirement "The CRDs are the controller's generated CRDs" in `controller-module`.

### Requirement: The controller's RBAC is the RBAC it declares

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-019); this requirement is restated under the controller's name.

**Migration**: See the requirement "The controller's RBAC is the RBAC it declares" in `controller-module`.

### Requirement: Drift between the module and the controller fails the check

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-019); this requirement is restated under the controller's name.

**Migration**: See the requirement "Drift between the module and the controller fails the check" in `controller-module`.

### Requirement: The operator's pods satisfy Pod Security restricted

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-019); this requirement is restated under the controller's name.

**Migration**: See the requirement "The controller's pods satisfy Pod Security restricted" in `controller-module`.

### Requirement: The controller's pod and Service stay those of the kustomize tree

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-019); this requirement is restated under the controller's name.

**Migration**: See the requirement "The controller's pod and Service stay those of the kustomize tree" in `controller-module`.

### Requirement: The operator's tuning is typed instance values

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-019); this requirement is restated under the controller's name.

**Migration**: See the requirement "The controller's tuning is typed instance values" in `controller-module`.

### Requirement: The configured registry mapping has one place

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-019); this requirement is restated under the controller's name.

**Migration**: See the requirement "The configured registry mapping has one place" in `controller-module`.

### Requirement: A mirror needs only the repository value

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-019); this requirement is restated under the controller's name.

**Migration**: See the requirement "A mirror needs only the repository value" in `controller-module`.

### Requirement: The image's tag and digest are not values

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-019); this requirement is restated under the controller's name.

**Migration**: See the requirement "The image's tag and digest are not values" in `controller-module`.

### Requirement: The module renders against its own pins with no cluster

**Reason**: The project is renamed from opm-operator to opm-controller (ADR-019); this requirement is restated under the controller's name.

**Migration**: See the requirement "The module renders against its own pins with no cluster" in `controller-module`.
