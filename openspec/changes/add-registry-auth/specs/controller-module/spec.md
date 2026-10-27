## ADDED Requirements

### Requirement: The module accepts an optional registryAuth value

The module's `#config` SHALL accept an optional value `registryAuth` with one field, `secretName`, a Kubernetes object name: lowercase letters, digits, `-` and `.`, starting and ending alphanumeric, at most 253 characters.

#### Scenario: A bad name is refused at the value

- **WHEN** the module is rendered with an empty `secretName`, an uppercase name, or a name longer than 253 characters
- **THEN** the render fails and names `registryAuth.secretName`

### Requirement: A set registryAuth mounts the Secret the administrator owns

When `registryAuth` is set, the controller Deployment SHALL gain exactly three things: the env `DOCKER_CONFIG=/etc/opm/docker`; a read-only mount of the volume `registry-auth` at `/etc/opm/docker`; and the volume `registry-auth`, a Secret volume of the named Secret with the one item `.dockerconfigjson` at `config.json`.

#### Scenario: The value is set

- **WHEN** the module is rendered with `registryAuth: secretName: "opm-registry-auth"`
- **THEN** the controller container has the env `DOCKER_CONFIG=/etc/opm/docker` and a read-only mount `registry-auth` at `/etc/opm/docker`, and the pod has a volume `registry-auth` from the Secret `opm-registry-auth` with the item `.dockerconfigjson` at `config.json`

### Requirement: The module owns no Secret and grants no Secret access

The module SHALL NOT render a Secret object for `registryAuth` and SHALL NOT add a rule on Secrets to any role it renders.

#### Scenario: The module owns no Secret and grants no Secret access

- **WHEN** the module is rendered with `registryAuth` set
- **THEN** the render contains no object of kind `Secret`, and no rendered Role or ClusterRole has a rule on `secrets`

### Requirement: An unset registryAuth changes nothing

When `registryAuth` is unset, the module SHALL render the same objects as before the value existed: no `DOCKER_CONFIG` env, no `registry-auth` mount and no `registry-auth` volume. The Deployment SHALL still equal the one the kustomize tree builds, the container SHALL keep a read-only root filesystem, and the manager role SHALL gain no rule.

#### Scenario: The value is unset

- **WHEN** the module is rendered with no `registryAuth`
- **THEN** the controller Deployment has no `DOCKER_CONFIG` env, no mount named `registry-auth` and no volume named `registry-auth`, and the kustomize parity check passes

### Requirement: A missing Secret keeps the pod from starting

The module SHALL NOT mark the `registry-auth` volume's Secret optional: the rendered volume SHALL carry no `optional` field, or `optional: false`. The volume SHALL NOT set a file mode that the controller's non-root user cannot read. A Secret that does not exist, or that lacks the key `.dockerconfigjson`, SHALL therefore keep the controller pod from starting, with a mount event that names the Secret.

#### Scenario: The Secret is not optional

- **WHEN** the module is rendered with `registryAuth` set
- **THEN** the `registry-auth` volume has no `optional: true`, and any `defaultMode` it sets is at least `0444`

#### Scenario: The Secret does not exist

- **WHEN** the controller is deployed with `registryAuth` set and the named Secret does not exist
- **THEN** the pod stays in `ContainerCreating` and an event of reason `FailedMount` names the Secret

### Requirement: registryAuth leaves the other values and guards as they are

Setting `registryAuth` SHALL NOT change the Deployment selector, the image, the arguments, the resources, the Pod Security posture or the read-only root filesystem. `extraArgs` SHALL keep refusing the arguments it refuses today.

#### Scenario: Every value set together

- **WHEN** the module is rendered with every other value set and `registryAuth` set
- **THEN** the selector equals the fixed selector, the arguments equal those rendered without `registryAuth`, and the pods are admitted under Pod Security restricted
