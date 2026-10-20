## MODIFIED Requirements

### Requirement: The controller's pod and Service stay those of the kustomize tree
While `config/manager` and `config/default` still install the operator for development (`task operator:installer`, `task operator:controller:install`) and in e2e, the controller Deployment's pod spec and the metrics Service the module renders SHALL equal those of a kustomize build of `config/default`, apart from the image, labels, the selector, the fields the catalog sets to Kubernetes API defaults, and the bindings' names.

#### Scenario: A manifest-only edit fails the test

- **WHEN** a pull request changes the manager container's arguments, probes, environment, volumes or security context in `config/manager` and not in the module
- **THEN** the module's render test fails naming the differing field

#### Scenario: An object added to the kustomize tree only fails the test

- **WHEN** a pull request adds an object to `config/default`, such as a binding in `config/rbac` that grants the controller's ServiceAccount an existing role, and not to the module
- **THEN** the module's render test fails naming the object, a binding by its role and subjects

### Requirement: A test compares the module's pod and Service with the kustomize build on every pull request
As in requirement "The controller's pod and Service stay those of the kustomize tree" (the controller Deployment's pod spec and the metrics Service the module renders): A test SHALL compare them on every pull request, so e2e never tests a controller pod the module does not ship. The operator's releases no longer publish a manifest built from `config/default`; the module's render is the only install manifest.

#### Scenario: A pull request is compared with the kustomize build
- **WHEN** a pull request is opened
- **THEN** a test compares the set of objects the module renders, the controller Deployment's spec and pod spec, and the metrics Service with those of a kustomize build of `config/default`
