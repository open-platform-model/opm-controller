---
title: "Install the controller"
description: "Install the OPM controller so the cluster reconciles instances on its own."
type: how-to
weight: 16
---

The OPM controller runs in the namespace `opm-controller-system`. It renders and applies every ModuleInstance whose `spec.owner` is absent or `controller`, such as the ones you create with kubectl or through GitOps, and every ModulePackage. Install it when the cluster should keep instances as declared without anyone running `opm`. A cluster where only the CLI deploys needs only the controller's resource definitions, which `opm controller install --crds-only` installs.

The controller is itself an OPM module, `opmodel.dev/modules/opm_controller`, published to the registry on a version train of its own. A module version is not a controller version: each module version deploys one controller release. `opm controller install` pulls that module and deploys it as the ModuleInstance `opm-controller` in `opm-controller-system`, owned by the CLI. The controller never reconciles that instance, whatever its owner. So re-running the install repairs the controller, as long as the instance stays owned by the CLI. If its `spec.owner` is set to anything other than `cli`, the install refuses and prints the `kubectl patch` that sets it back. Without the CLI, `kubectl apply` installs the same objects from the `install.yaml` that each module release attaches.

## Before you begin

- `kubectl` access with cluster-admin rights. The module renders the namespace `opm-controller-system`, four resource definitions under `opmodel.dev` (`moduleinstances`, `modulepackages`, `platforms` and `transformerregistrations`), and ClusterRoles and ClusterRoleBindings.
- To install with `opm`: the CLI, configured with `opm config init`. See [Install the CLI](/docs/start/install-the-cli/). The CLI pulls the controller module and the OPM modules it depends on from its configured registry, `ghcr.io/open-platform-model` by default, and looks up the newest catalog release there before it changes anything on the cluster. The module also depends on the Kubernetes schemas `cue.dev/x/k8s.io`, which the default configuration pulls from `registry.cue.works`, so the host that runs `opm` must reach both registries.
- Network access from the cluster to `ghcr.io`. The nodes pull the controller image from it, and the controller resolves core, the catalogs and modules from `ghcr.io/open-platform-model`. It exits at startup when it cannot resolve core.
- For an air-gapped cluster: a registry mirror that holds the controller module and every module it depends on, and an image mirror the nodes can pull from. Step 5 sets both.
- To use ModulePackages: Flux's source-controller, installed before the controller. A ModulePackage loads its instance from a Flux OCIRepository, GitRepository or Bucket, and the controller looks for those resource definitions only when it starts. ModuleInstances need no Flux.

## Steps

1. Install the controller.

   With the `opm` CLI, run:

   ```sh
   opm controller install
   ```

   The output should look similar to this, shortened. The CLI names the module version it installed and the controller version that module deploys:

   ```text
   INFO controller module opmodel.dev/modules/opm_controller 0.Y.Z (pinned; deploys opm-controller v1.0.0-beta.N)
   ...
   ✔ opm-controller v1.0.0-beta.N installed from module 0.Y.Z
   ```

   Each opm release pins one module version, which the install uses by default. To install another, add `--version` with a module version: `0.2.0` pins that release, and `v0` takes the newest release of major 0. `--version` does not take a controller release tag such as `v1.0.0-beta.N`.

   The install first runs every check that can refuse it, and a refused install changes nothing on the cluster. Among them, it refuses a module whose controller has a higher major or minor version than the CLI, and values the module does not accept. It also refuses any object the module renders that already exists and is not managed by OPM. Then it applies the resource definitions and waits until the API server serves them. Then it applies the whole render as the instance, records every object it applied in the instance's inventory, and waits until the controller has rolled out. Last, it creates the cluster Platform (step 3). `--timeout` bounds every wait together, 5 minutes by default. The install first waits out objects that an earlier uninstall left terminating. `--crds-only` applies only the resource definitions this opm binary carries for the module version it pins. It needs no registry, creates no instance and no Platform, and refuses `--version`.

   To upgrade the controller, run `opm controller install` again: with a newer opm release, or with `--version`. The install applies what changed and deletes the objects the new module version no longer renders, except the resource definitions and the namespace. Re-running it with the same version and values changes nothing.

   Without the CLI, apply the `install.yaml` of a module release. Pick a tag `opm_controller-vX.Y.Z` from the [controller releases](https://github.com/open-platform-model/opm-controller/releases); this page uses `opm_controller-v0.1.0`.

   To apply the module release's manifest:

   ```sh
   kubectl apply --server-side -f https://github.com/open-platform-model/opm-controller/releases/download/opm_controller-v0.1.0/install.yaml
   ```

   The manifest is the module rendered with its default values for the instance `opm-controller` in `opm-controller-system`. It pins the controller image by tag and digest, `ghcr.io/open-platform-model/opm-controller:v1.0.0-beta.N@sha256:...`. To check the image's signature, run this with the digest the manifest names:

   ```sh
   cosign verify ghcr.io/open-platform-model/opm-controller@sha256:<digest> \
     --certificate-identity-regexp='^https://github.com/open-platform-model/opm-controller/\.github/workflows/release\.yml@refs/heads/main$' \
     --certificate-oidc-issuer=https://token.actions.githubusercontent.com
   ```

   A kubectl install writes no instance record, so `opm controller uninstall` refuses on it. Running `opm controller install` later records the running controller as its instance.

   Do not use the `releases/latest/download/install.yaml` link. The repository publishes two release trains, and the module's releases are never marked Latest. Every v1.0.0 controller release is marked Pre-release, which GitHub's Latest link skips, so that link serves a retired manifest.

2. Wait for the controller to roll out.

   `opm controller install` has already waited. If you applied the manifest with kubectl, run:

   ```sh
   kubectl -n opm-controller-system rollout status deployment/opm-controller-manager --timeout=180s
   ```

   The output should look similar to this:

   ```text
   deployment "opm-controller-manager" successfully rolled out
   ```

3. Create the cluster Platform.

   The controller renders against the cluster's Platform, a cluster-wide resource that pins one build of each catalog. Until the controller has generated the Platform, every ModuleInstance and ModulePackage waits with `Ready=False` and reason `PlatformNotReady`. The controller's own instance does not depend on it: the install renders the controller module against the module's own catalog pins, never against the cluster Platform.

   If you installed with `opm controller install`, it has already created the Platform `cluster`, subscribed to the newest release of `opmodel.dev/catalogs/opm@v4`. It leaves an existing Platform untouched, and with `--skip-platform` it creates none.

   Otherwise, write the Platform to a file, `platform.yaml`:

   ```yaml
   apiVersion: opmodel.dev/v1alpha1
   kind: Platform
   metadata:
     name: cluster
   spec:
     type: kubernetes
     registry:
       opmodel.dev/catalogs/opm@v4:
         version: "4.Y.Z"
   ```

   Replace `4.Y.Z` with a published build of the opm catalog, such as the newest on the [opm catalog](/catalogs/opm/4/) page.

   Apply it:

   ```sh
   kubectl apply -f platform.yaml
   ```

   `cluster` is the only name a Platform can have. Each key under `spec.registry` is a catalog path with its major version, and `version` names exactly one published build of that major. `spec.type` is informational. The optional `spec.skewPolicy` decides what happens when a module requires a newer build of core or a catalog than the Platform pins: `Warn`, the default, renders against the Platform's build and reports the skew, and `Refuse` refuses the render. See [Platforms and catalogs](/docs/concepts/platforms-and-catalogs/).

4. Give the controller an identity to apply with.

   The controller applies a module's resources as a ServiceAccount in the instance's own namespace, by impersonation. Its own ClusterRole, `opm-controller-manager-role`, grants nothing on workload kinds such as Deployments and Services. With no ServiceAccount named, the controller applies as itself, the apply is forbidden, and the instance reports reason `ApplyFailed`.

   In each namespace that holds instances, create a ServiceAccount and bind it to a role that covers what the modules render. Then name it in the `spec.serviceAccountName` of each ModuleInstance or ModulePackage. For example, the built-in `edit` ClusterRole, bound in one namespace, covers common namespaced kinds such as Deployments, StatefulSets, Services, ConfigMaps and Secrets:

   ```yaml
   apiVersion: v1
   kind: ServiceAccount
   metadata:
     name: opm-applier
     namespace: shop
   ---
   apiVersion: rbac.authorization.k8s.io/v1
   kind: RoleBinding
   metadata:
     name: opm-applier
     namespace: shop
   roleRef:
     apiGroup: rbac.authorization.k8s.io
     kind: ClusterRole
     name: edit
   subjects:
     - kind: ServiceAccount
       name: opm-applier
       namespace: shop
   ```

   A module that renders cluster-scoped objects needs a ClusterRole bound with a ClusterRoleBinding instead. Bind `cluster-admin` only on a test cluster: the controller can then apply anything a module renders.

   To use one ServiceAccount name in every namespace instead of naming it on each instance, set the module value `defaultServiceAccount`. Write it to a values file, `controller-values.cue`:

   ```cue
   values: defaultServiceAccount: "opm-applier"
   ```

   Then install with it:

   ```sh
   opm controller install -f controller-values.cue
   ```

   The module renders the value as the controller's `--default-service-account` argument. The install records the values on the controller's instance, and a later install keeps every recorded value it does not change, so you pass each value once. `--reset-values` starts again from the module's defaults. `spec.serviceAccountName` still wins where it is set. The ServiceAccount must exist in every namespace that holds an instance, or that instance stalls with reason `ImpersonationFailed`.

   The module's other values are `registry` (step 5), `image.repository`, `replicas`, `resources` for the controller container, and `extraArgs` for further controller arguments. The [module's README](https://github.com/open-platform-model/opm-controller/tree/main/modules/opm_controller#values-config) lists their defaults. The install refuses a value the module does not accept and names it.

5. Point the controller at your module registry.

   The controller resolves core, the catalogs and modules through a CUE registry mapping. By default it maps `opmodel.dev` and `testing.opmodel.dev` to `ghcr.io/open-platform-model`, and every other path to `registry.cue.works`. If your modules are published somewhere else, set the module value `registry` in CUE registry syntax. It replaces the default, so keep the `opmodel.dev` entry:

   ```cue
   values: registry: "example.com=registry.example.com,opmodel.dev=ghcr.io/open-platform-model,registry.cue.works"
   ```

   Install with the file as in step 4. The install layers it over the recorded values, so a `defaultServiceAccount` set earlier stays.

   On an air-gapped cluster, three things read from a registry, and each needs your mirror:

   - The CLI pulls the controller module and its dependencies through its own registry mapping: the global `--registry` flag, the `OPM_REGISTRY` environment variable or the CLI's config file. Point it at a mirror that holds every module the controller module depends on, for example `OPM_REGISTRY=mirror.example.com/cue opm controller install`.
   - The controller resolves modules through the `registry` value above. Set it to the same mirror, as the cluster reaches it.
   - The nodes pull the controller image by its digest. Mirror the image with its digest unchanged, then either set the module value `image.repository` to the mirror's repository, or configure the nodes' container runtime to pull `ghcr.io` through the mirror. The tag and digest are not values; they stay the module's.

6. Grant users access to ModuleInstances.

   People who create and edit controller-managed instances need rights on ModuleInstances only; the controller applies the workloads. The module renders three ClusterRoles to bind:

   - `opm-controller-moduleinstance-admin-role`: every verb on ModuleInstances.
   - `opm-controller-moduleinstance-editor-role`: create, delete, get, list, patch, update and watch.
   - `opm-controller-moduleinstance-viewer-role`: get, list and watch.

   Each of them can also read a ModuleInstance's status. For example, to let the group `shop-team` edit instances in the namespace `shop`:

   ```sh
   kubectl -n shop create rolebinding shop-team-instances \
     --clusterrole=opm-controller-moduleinstance-editor-role --group=shop-team
   ```

   For people who use the `opm` CLI, `opm controller install --rbac --user <name>` (or `--group <name>`) also creates the ClusterRole `opm-cli-user` and binds it cluster-wide. It grants every verb on ModuleInstances, get, patch and update on their status, and get and list on Platforms. `--user` and `--group` need `--rbac`, and only one of them can be given. The role is not part of the controller's instance.

7. Grant read-only access to the platform.

   People who only look, in a terminal or a dashboard, may also need to read the Platform, the ModulePackages and the TransformerRegistrations. The module renders one viewer ClusterRole for each, granting get, list and watch on the kind and get on its status:

   - `opm-controller-platform-viewer-role`: Platforms. The kind is cluster-scoped, so bind it with a ClusterRoleBinding.
   - `opm-controller-modulepackage-viewer-role`: ModulePackages. Bind it with a RoleBinding to grant one namespace, or a ClusterRoleBinding to grant all of them.
   - `opm-controller-transformerregistration-viewer-role`: TransformerRegistrations. Cluster-scoped, so bind it with a ClusterRoleBinding.

   A RoleBinding to a role for a cluster-scoped kind is accepted but grants nothing. For example, to let the group `platform-readers` read the Platform and the registrations:

   ```sh
   kubectl create clusterrolebinding platform-readers-platform \
     --clusterrole=opm-controller-platform-viewer-role --group=platform-readers
   kubectl create clusterrolebinding platform-readers-registrations \
     --clusterrole=opm-controller-transformerregistration-viewer-role --group=platform-readers
   ```

   None of these roles is bound when the controller is installed, and none of them is added to the built-in `view` role: binding `view` grants no read of OPM kinds.

## Check that it worked

```sh
kubectl get platform cluster
```

The output should look similar to this:

<!-- x-release-please-start-version -->

```text
NAME      TYPE         READY   REASON      CONTROLLER
cluster   kubernetes   True    Generated   v1.0.0-beta.8
```

<!-- x-release-please-end -->

`Generated` means the controller generated the platform from the subscribed catalogs and built it. `CONTROLLER` is the version of the running controller. `READY` `False` with reason `BuildFailed` usually means a subscribed catalog version did not resolve: it is not published, or the cluster cannot reach the registry. The condition's message names the cause. See [Controller conditions](/docs/diagnostics/controller-conditions/).

If you installed with `opm`, the controller's own instance is owned by the CLI:

```sh
kubectl -n opm-controller-system get moduleinstance opm-controller -o jsonpath='{.spec.owner}{"\n"}'
```

The output should be:

```text
cli
```

To read the controller's log:

```sh
kubectl -n opm-controller-system logs deploy/opm-controller-manager
```

The first lines should look similar to this; the second names the core version it resolved:

<!-- x-release-please-start-version -->

```text
INFO	setup	Starting opm-controller	{"version": "v1.0.0-beta.8"}
INFO	setup	OPM core schema resolved	{"version": "..."}
INFO	Flux source CRDs not installed; ModulePackage source watches disabled	{"controller": "modulepackage", "kinds": "OCIRepository,GitRepository,Bucket"}
INFO	setup	Starting manager
```

<!-- x-release-please-end -->

The Flux line appears on a cluster without Flux. If you install Flux after the controller, restart the controller so that it watches Flux sources:

```sh
kubectl -n opm-controller-system rollout restart deployment/opm-controller-manager
```

## Related

- [Controller Reference](/docs/reference/controller/)
- [Platforms and catalogs](/docs/concepts/platforms-and-catalogs/)
- [Who owns an instance](/docs/concepts/who-owns-an-instance/)
- [Delete an instance safely](/docs/operating/delete-an-instance-safely/). To remove the controller, run `opm controller uninstall`. It deletes the objects the controller's instance records, then the instance, and keeps the resource definitions and the namespace. It refuses while any ModuleInstance still carries the controller's cleanup finalizer; `--remove-finalizers` strips that finalizer and proceeds, leaving those instances' resources unmanaged. On a cluster with no controller instance, it deletes nothing and names `opm controller install`.
