## Why

"Operator" in the Kubernetes sense names a controller that carries the day-2 knowledge of one application. This project knows no application: it renders any OPM module against a platform and applies the result, a generic delivery controller like Flux's controllers or Argo CD. The name promised knowledge the project does not have, and it shared a word with the platform operator, the person the docs address. The owner decided on 2026-10-05 to rename the project to `opm-controller` everywhere, with no migration, alias or compatibility path, because no installation exists to carry over. ADR-019 records the decision and the two earlier names.

The blocked change `stop-operator-install-manifest` is folded in (the owner's 2026-10-05 decision to fold it into the rename): its gates now hold (the cli installs from the module and no longer downloads a controller release's `install.yaml`, and the cascade resolver checks an asset only when asked), so a controller release stops attaching `install.yaml` here. Its legacy install-page task is dropped, and its plan to stop the controller release's cascade notification is dropped because the module release has no notification to take over.

## What Changes

- **BREAKING** Go module `github.com/open-platform-model/opm-controller`; image `ghcr.io/open-platform-model/opm-controller`; namespace `opm-controller-system`; name prefix `opm-controller-`; `app.kubernetes.io/name: opm-controller`. The kubebuilder base name `controller-manager` becomes `manager`, so the Deployment and ServiceAccount are `opm-controller-manager`. The pod label `control-plane` and the controller module's component key follow it: the module's fixed Deployment selector is `manager` (`control-plane`, `app.kubernetes.io/name`, `component.opmodel.dev/name`) before module 0.1.0 publishes it.
- **BREAKING** API: `spec.owner` takes `cli` or `controller` (`OwnerController`); `Platform.status.controllerVersion` replaces `operatorVersion`, with the printer column `Controller`; the render input key carries `ControllerVersion`. The value `operator` is rejected.
- **BREAKING** The self-skip coordinates are the instance `opm-controller` in `opm-controller-system` and the module path `opmodel.dev/modules/opm_controller`.
- The controller no longer strips a cleanup finalizer from its own instance (the owner's 2026-10-06 decision). That guard only handled state an opm-operator release could leave, and no install carries over (0032:D1:R6). The refusal before finalizer registration stays.
- **BREAKING** The controller's module moves to `modules/opm_controller` (`opmodel.dev/modules/opm_controller@v0`, sub-package `controller`), on a new tag train `opm_controller-v*` that starts at 0.1.0 with a fresh changelog and no manifest key. Its tooling moves to `hack/controller-module` with `min-controller-version`.
- **BREAKING** A controller release no longer attaches `install.yaml`; the install manifest is the module release's asset only (folded from `stop-operator-install-manifest`).
- Workflows, release guards, release-please (`package-name: opm-controller`, both module paths excluded for one root release), Taskfile namespaces (`controller:*`, `controller-module:*`, with `controller:install` flattened), kind and buildx names, the registry client type and the startup log follow the new name.
- The docs bundle is `opm-controller`, owning `reference/controller/` ("Controller Reference"); the pages move to `start/install-the-controller` and `diagnostics/controller-conditions`. The fixtures move to `testing.opmodel.dev/{modules,releases,catalogs}/controller/*`.
- Prose calls the product the controller and its per-kind loops reconcilers. The human role keeps the word operator. Recorded run results (the kind-opm-dev FINDINGS header and test-plan sign-off) keep the names they were recorded under (the owner's 2026-10-06 decision).
- ADR-019 is added and ADR-015 is deleted. The `Makefile` is deleted; the e2e suite deploys through `task controller:deploy`, `controller:crds` and their removals.

Release class: MAJOR after GA (renamed API values and fields, a new module path, a removed release asset); on the beta line it ships as the next `1.0.0-beta.N` through a `!` in the PR title. The module starts at 0.1.0 on its new path.

Out of scope: the cli (`opm controller`, its own change), the `.github` cascade maps, the workspace, the site and every other consumer; the annotation `opm.dev/force-delete-orphan`; the active change `add-cross-namespace-source-grants`.

## Capabilities

### New Capabilities

- `controller-module`: recreates `operator-module` under the controller's name; the fixed names are those `config/` gives under the new prefix, and the pod-shape comparison with `config/default` is kept for development and e2e, not for a release manifest.
- `controller-module-release`: recreates `operator-module-release`.
- `controller-resource-reference`: recreates `operator-resource-reference`; `reconciledBy` names reconcilers.
- `controller-version-identity`: recreates `operator-version-identity`.

### Modified Capabilities

- `operator-module`, `operator-module-release`, `operator-resource-reference`, `operator-version-identity`: every requirement removed (retired; successors above).
- `container-image-publish`: the release install manifest requirement is removed and the release image job holds `contents: read`.
- `release-automation`: a controller release requires only `opm-examples.tar.gz`; requirements whose scenario named the operator are restated under new titles.
- `module-instance-ownership`: the requirement "The operator releases its finalizer from its own instance without pruning" is removed and not restated.
- `instance-health`: the `Healthy` requirement's body names the controller.
- `module-instance-ownership`, `platform-crd`, `render-input-key`, `workflow-hardening`, `cascade-receiver`, `docs-bundle`, `deps-cascade`, `example-module-publishing`, `kernel-module-renderer`, `kubernetes-tier-adoption`, `library-kernel-runtime`, `modulepackage-kernel-rendering`, `platform-gated-rendering`, `platform-reconciler`, `viewer-roles`: requirements or scenarios that named the operator are removed and restated under new titles; body-only wording moves in these and the remaining specs that name the operator.

## Impact

- API types `ModuleInstance` (owner enum) and `Platform` (status field, printer column); all four reconcilers' prose; `internal/reconcile` self-skip; `internal/status` render input key and condition text; `cmd/main.go`.
- `config/`, `dist/install.yaml`, `modules/opm_controller`, `hack/`, `.tasks/`, `.github/`, `docs/`, `adr/`, README, AGENTS, CONSTITUTION, `openspec/config.yaml`.
- No reconcile phase changes behavior apart from the renamed values and the removed finalizer release on the controller's own instance.
