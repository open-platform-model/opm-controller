## Why

The controller pulls core, catalogs and modules anonymously. Its module renders no env, volume or Secret, and `#config` is closed. A private registry, such as a private GHCR package, cannot be used. The module owner needs one typed value that gives the controller its registry logins, with no change to the controller binary.

## What Changes

- A new optional value `registryAuth.secretName` on the module `modules/opm_controller`. When set, the controller Deployment gets the env `DOCKER_CONFIG=/etc/opm/docker`, a read-only mount at `/etc/opm/docker`, and a volume from the named Secret with the item `.dockerconfigjson` mapped to `config.json`.
- When the value is unset, the render is unchanged, byte for byte.
- The administrator creates the Secret (type `kubernetes.io/dockerconfigjson`) in `opm-system`. The module does not own it and does not mark it optional, so a missing Secret keeps the pod in `ContainerCreating` with a mount event.
- One set of logins serves the whole controller: every module and every catalog pull. A new login takes effect after a restart of the controller Deployment.
- Tests: render with the value set and unset, a missing Secret is not optional, a pull that needs a login works, and the behaviour after a refused pull. The result of the last test decides one sentence in the docs.
- Docs: the module README and the controller install page.

Delivery: one PR per section (a PR that changes `modules/opm_controller` changes nothing else)

Decisions that stay as the owner approved them: one credential set for the whole controller, the administrator creates the Secret, the value name `registryAuth.secretName`, restart to rotate, a missing Secret blocks the pod.

## Capabilities

### New Capabilities
- `controller-registry-credentials`: how the controller takes registry logins from a mounted docker config, which pulls they serve, the exposure rules, and the behaviour after a refused pull.

### Modified Capabilities
- `controller-module`: the module gains the optional `registryAuth` value and its render rules (added requirements only; no existing requirement or scenario changes).

## Impact

- **Module** (`modules/opm_controller`): `module.cue` (`#config.registryAuth`), `components.cue` (env, mount, volume), `README.md` (values table). Release class MINOR for the module (`feat`), shipped by the module's own release. Until GA the module ships it as the next `0.x` release.
- **Controller binary, CRDs, RBAC, library, catalog:** none. No `+kubebuilder:rbac` marker, no regenerated `zz_generated_*`, no Go change. Verified in design.md.
- **Tests:** `test/integration/controllermodule/` (render), `test/integration/reconcile/` (login and refused pull against an in-process registry), `test/e2e/` (a Kind registry that needs a login).
- **Docs:** `docs/site/start/install-the-controller.md` (new step on private registries), `modules/opm_controller/README.md`. The later registry page that covers publishing is not part of this change.
- **Compatibility:** existing instances and values are unaffected. Not breaking.
- **Complexity:** one optional value and three conditional render parts. No new flag, no new code path in the controller.

## Enhancement

None.
