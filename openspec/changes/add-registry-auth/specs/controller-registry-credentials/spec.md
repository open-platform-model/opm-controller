## ADDED Requirements

### Requirement: One set of logins serves every module and catalog pull

The controller SHALL take registry logins from the docker config file at `$DOCKER_CONFIG/config.json`, one file for the whole process, and SHALL use a matching login for every module, catalog and core pull it makes through the library Kernel and for the Platform module closure. The controller SHALL NOT take a login from a flag, from the registry mapping, or from a custom resource. No controller code SHALL be needed beyond passing the environment through.

#### Scenario: A module pull that needs a login

- **WHEN** a module is held only in a registry that requires a login, and `DOCKER_CONFIG` names a directory whose `config.json` holds a login for that registry host
- **THEN** the Kernel's module pull succeeds

#### Scenario: A module pull without a login

- **WHEN** the same pull runs with no login for the host
- **THEN** the pull is refused and the error reaches the instance as the existing resolution failure

#### Scenario: A catalog pull

- **WHEN** a catalog is held only in a registry that requires a login, and the login is in `config.json`
- **THEN** the catalog pull succeeds through the same client

### Requirement: A pull after a refused pull reads the logins again

After a pull is refused by a registry, the controller SHALL read the docker config file again for the next pull, without a restart, so a file replaced with a correct login is used by that pull. The controller SHALL NOT claim that a rotated login takes effect before a pull has failed: while an old login still works, the controller keeps using it until the controller restarts. The documentation SHALL state this result.

#### Scenario: The file is replaced after a refused pull

- **WHEN** a pull is refused with a wrong login, the file is then replaced with the right login, and the same Kernel pulls again
- **THEN** the second pull succeeds

#### Scenario: A rotated login while the old one works

- **WHEN** a pull has succeeded and the file is then replaced, while the registry still accepts the old login
- **THEN** the controller keeps using the old login until a pull fails or the controller restarts, and the documentation says to restart the Deployment to rotate

### Requirement: The login is a mounted file, never an argument

No login SHALL appear in the controller's arguments, its startup log line, the registry mapping, the instance's recorded values or an environment variable. The environment SHALL carry only the path of the directory that holds the file.

#### Scenario: The rendered arguments and environment

- **WHEN** the module is rendered with `registryAuth` set
- **THEN** no argument or environment value contains the Secret's content, and the only added environment value is the directory path `/etc/opm/docker`

### Requirement: The shared login is documented as shared

The documentation SHALL state that the logins serve every ModuleInstance and catalog the controller resolves, so anyone allowed to create a ModuleInstance can make the controller pull any module that the logins reach, that the token SHOULD be read-only and limited to the packages users may use, and that the controller's own image pull is configured with `imagePullSecrets`, not with this value.

#### Scenario: The install page

- **WHEN** a reader opens the controller install page
- **THEN** it states how to create the Secret, how to set `registryAuth.secretName`, how to rotate (update the Secret, then restart the Deployment), what a missing Secret looks like, and that the logins are shared
