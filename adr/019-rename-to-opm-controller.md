# ADR-019: Rename opm-operator to opm-controller

## Status

Accepted

## Context

The project has been renamed twice. It started from the kubebuilder scaffold as `poc-controller`, a proof of concept. When its container image was first published (ADR-014), it took the public name `opm-operator`, and a follow-up rename carried that name to the repository, the Go module, the namespace, the name prefix and the labels. That rename considered `opm-controller` and rejected it, for two reasons: "operator" was taken to be the term of art for a Kubernetes custom controller that ships CRDs, and `opm-controller` was already the runtime actor identity in the `app.kubernetes.io/managed-by` label and the server-side apply field manager (ADR-010), which the project name was not meant to re-couple.

Both reasons have weakened. In the CoreOS sense, an operator is a controller that carries the day-2 knowledge of one application: the Postgres or Prometheus operators back up, fail over and upgrade the one thing they know. This project knows no application. It renders any OPM module against a platform and applies the result, which makes it a generic delivery controller in the family of Flux's controllers and Argo CD. Calling it an operator told readers to expect application knowledge it does not have, and it collided with the human role of the platform operator, the person who runs the platform, in every page that addresses that person.

The second reason turned around once the controller became an OPM module of its own. The cli installs it as a ModuleInstance with fixed coordinates, and the controller refuses to reconcile that instance. The name a human sees on the instance, the namespace and the Deployment is now the name the controller writes on everything it applies, so one name serves both and nothing has to be explained.

## Decision

The project is renamed from `opm-operator` to `opm-controller` across every identifier that is not immutable history, with no migration, alias or compatibility path, because the controller has no production installations to carry over.

- Repository `open-platform-model/opm-controller`, Go module `github.com/open-platform-model/opm-controller`, image `ghcr.io/open-platform-model/opm-controller`.
- Namespace `opm-controller-system`, name prefix `opm-controller-`, `app.kubernetes.io/name: opm-controller`. The kubebuilder base name `controller-manager` becomes `manager`, so the Deployment and its ServiceAccount are `opm-controller-manager` rather than `opm-controller-controller-manager`. The pod label `control-plane: controller-manager` stays.
- The controller's own module is `opmodel.dev/modules/opm_controller@v0`, with the deployed release in the sub-package `controller`, on a new tag train `opm_controller-vX.Y.Z` that starts at 0.1.0. Its instance is `opm-controller` in `opm-controller-system`.
- The API says controller: `spec.owner` takes `cli` or `controller`, `Platform.status.controllerVersion` replaces `operatorVersion`, and the Platform printer column is `Controller`. The value `operator` is rejected.
- The docs bundle is `opm-controller`, with its reference section at `/docs/reference/controller/`.
- The test fixtures move to `testing.opmodel.dev/{modules,releases,catalogs}/controller/*`.

The project name and the runtime actor identity are now the same string by design: `opm-controller` is the instance that deploys the controller, the value of `app.kubernetes.io/managed-by` on every object it applies, and its field manager. The instance and the managed-by label still mean different things, and the controller tells its own instance apart by its coordinates, its module path and the CRDs its inventory records, never by the label.

The vocabulary follows the name. "The controller" is the product: the binary, its Deployment and its releases. The per-kind loops inside it, one for each of ModuleInstance, ModulePackage, Platform and TransformerRegistration, are "reconcilers". "Operator" means only the human role of the platform or infrastructure operator, the Kubernetes match `operator` field, or a third-party operator such as cert-manager or the Prometheus operator.

**Alternatives considered and rejected:**

- *Keep `opm-operator`.* Rejected: it keeps promising application knowledge the project does not have, and keeps one word for both the product and the person who runs the platform.
- *Rename with a transition period, accepting `owner: operator` and the old status field alongside the new ones.* Rejected: a transition is worth its cost only for installations to carry over, and there are none.
- *Keep the kubebuilder base name `controller-manager`.* Rejected: with the new prefix it gives `opm-controller-controller-manager`, a stutter that would last in every pod name, cli constant and RBAC subject.

## Consequences

**Positive:** The name says what the project is, and "operator" in the docs means a person again. The product name, the instance name and the actor identity agree, so the controller has one name in every place a human looks.

**Negative:** Nothing carries over from an `opm-operator` install. The namespace, the Deployment and its selector, the RBAC names, the owner value and the status field all change, so a cluster running an `opm-operator` release is moved by uninstalling it and installing `opm-controller`. The new image, docs bundle, module and fixtures are new packages on GHCR.

**Trade-off:** The old name stays in what cannot change: git tags, published releases, CHANGELOG entries and archived OpenSpec changes. GitHub redirects the old repository URL, so links to that history keep working as long as no new repository takes the old name. This ADR replaces the earlier rename record, which is removed, and states both renames.
