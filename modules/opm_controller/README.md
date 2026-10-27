# opm_controller

The opm-controller's own OPM module, `opmodel.dev/modules/opm_controller@v0`. It renders the controller's whole install through the first-party catalog (`opmodel.dev/catalogs/opm@v4`), on a version train of its own: the module version (`identity/identity.cue`) is not the controller version.

Each module release, tagged `opm_controller-vX.Y.Z`, publishes the module to GHCR and attaches its render at default values as `install.yaml`; a controller release attaches no install manifest. `config/` still installs the controller in development and e2e (`task controller:installer`, `task controller:install`), and the render test keeps the two equal.

## What it renders

19 objects, for the instance `opm-controller` in `opm-system` only (any other name or namespace is refused, naming the expected coordinates):

| Objects | Catalog resource |
| --- | --- |
| Namespace `opm-system` | `resources/v1alpha1 #Namespaces` |
| The four `opmodel.dev` CRDs | `resources/v1beta1 #CRDs`, each generated `spec` embedded whole |
| Deployment and ServiceAccount `opm-controller-manager`, Service `opm-controller-manager-metrics-service` | `blueprints/v1beta1 #StatelessWorkload` with `#ServiceAccount`, `#Volumes` and traits |
| ClusterRoles `opm-controller-manager-role`, `opm-controller-metrics-auth-role`, Role `opm-controller-leader-election-role`, and a binding for each | `resources/v1beta1 #Role` |
| ClusterRoles `opm-controller-metrics-reader`, `opm-controller-moduleinstance-{admin,editor,viewer}-role`, `opm-controller-{platform,modulepackage,transformerregistration}-viewer-role`, `opm-controller-transformerregistration-admin-role`, unbound | `resources/v1alpha1 #Objects` (raw objects) |

The eight administrator ClusterRoles render as raw objects with no binding until a catalog release can render a role with no subjects (catalog_opm `add-subjectless-roles`); their rules still come from the generated RBAC data. Nothing else uses raw objects.

Every name is a constant, the name `config/` gives the object with the prefix `opm-controller-`; only the bindings take the catalog's names (the role's own). The Deployment's selector is fixed across module versions, and the render test fails if anything changes it.

## Vetting and building it with the cli

The cli names a module's synthetic instance `<name>-debug` in `default`, which this module refuses. Give it the module's coordinates (`module vet` takes `--instance-name`, `module build` takes `--name`):

```bash
opm module vet . --instance-name opm-controller -n opm-system
opm module build . --name opm-controller -n opm-system
```

## Values (`#config`)

| Field | Default | Effect |
| --- | --- | --- |
| `image.repository` | `ghcr.io/open-platform-model/opm-controller` | The image repository, for a mirror. The tag and digest stay the module's. |
| `registry` | unset | `--registry=<value>`: the CUE registry mapping the controller resolves modules through. |
| `defaultServiceAccount` | unset | `--default-service-account=<value>`: the identity the controller applies as when a ModuleInstance names none. |
| `registryAuth.secretName` | unset | The name of a `kubernetes.io/dockerconfigjson` Secret in `opm-system`. The module mounts its `.dockerconfigjson` key as `/etc/opm/docker/config.json` (read-only, one container) and sets `DOCKER_CONFIG=/etc/opm/docker`, so the controller finds registry logins for the modules, catalogs and core it pulls. Unset renders none of this. See [Registry logins](#registry-logins). |
| `resources` | requests `100m` CPU, `256Mi`; limits `2` CPU, `4Gi` | The manager container's resources. Each quantity defaults on its own. CPU is a number of cores (`4`, `0.5`) or a millicore string (`"500m"`); a string of cores such as `"4"` is refused. The memory limit is given as `<n>Mi` or `<n>Gi`; `GOMEMLIMIT` is derived as 80 percent of it, in MiB. |
| `replicas` | `1` | The Deployment's replicas. |
| `extraArgs` | `[]` | Further controller arguments, after all of the above. `--registry`, `--default-service-account` and the module's own `--metrics-bind-address`, `--leader-elect` and `--health-probe-bind-address` are refused in either flag spelling. |

The arguments render in this order: `--metrics-bind-address=:8443`, `--leader-elect`, `--health-probe-bind-address=:8081`, `--registry`, `--default-service-account`, then `extraArgs`.

## Registry logins

Create the Secret yourself, in the namespace `opm-system`. The module renders that namespace, so do not create it by hand: `opm controller install` refuses a namespace that OPM does not manage. Install with `registryAuth.secretName` set, then create the Secret. The controller pod waits in `ContainerCreating` until the Secret exists, and the kubelet then mounts it:

```bash
kubectl -n opm-system create secret docker-registry opm-registry-auth \
  --docker-server=ghcr.io --docker-username=<user> --docker-password=<read-only token>
```

Set `registryAuth.secretName` to `opm-registry-auth` in the values file. The module creates no Secret and adds no permission to read Secrets: the kubelet mounts the file. The login is in a file, never in an argument or an environment variable.

- One set of logins serves every ModuleInstance and catalog. Give the token read access to the packages that users may pull, and no more.
- A missing Secret, or one without the key `.dockerconfigjson`, keeps the controller pod in `ContainerCreating` with a `FailedMount` event.
- To rotate, update the Secret and restart `deployment/opm-controller-manager` in `opm-system`. The controller reads the file again only after a pull fails.
- The controller's own image pull is `imagePullSecrets`, not this value.

The image tag, digest and pull policy are not values. `controller/controller.cue` names the one controller release this module version deploys, by version, tag and digest; read it without a render with `cue eval ./controller -e Version`. It never falls below `hack/controller-module/min-controller-version`, the first controller release that refuses to reconcile the instance deploying it.

## Regenerating

`zz_generated_crds.cue` and `zz_generated_rbac.cue` are generated from `config/` with `cue import`. Do not edit them. `task dev:manifests` regenerates them together with `config/`, so after an API type or RBAC marker change run:

```bash
task dev:manifests
```

`task controller-module:drift` (run in CI) fails when `config/` is stale against `task dev:manifests`, then when these files are stale against `config/`; `task controller-module:drift REF=<tag>` checks them against a release tag's `config/` instead. A Role or ClusterRole added to `config/rbac/kustomization.yaml` lands in the generated data, and the render test fails until `components.cue` renders it.

The render test is `test/integration/controllermodule`, run by `task dev:test`. It renders this directory against a platform generated from the module's own catalog pin, so it needs the pins from a registry (GHCR by default).
