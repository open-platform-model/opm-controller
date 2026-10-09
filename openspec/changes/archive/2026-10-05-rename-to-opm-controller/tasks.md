# Tasks: rename-to-opm-controller

Built on the branch `refactor/rename-to-opm-controller`, and merged together with the other repositories' rename changes as one PR titled `refactor!:`. Steps that need artifacts which do not exist until after the merge (the first `opm-controller` release, its module image PR, module 0.1.0) are follow-ups outside this change; design.md, "Image pin of the module", describes them.

## 1. Go module path

- [x] 1.1 `go.mod`, every import and `PROJECT` name `github.com/open-platform-model/opm-controller`; `.golangci.yml` follows.
- [x] 1.2 `go build ./...` and `go vet ./...` green, then commit `refactor!: move the Go module to github.com/open-platform-model/opm-controller`.

## 2. API, runtime names and install shape

- [x] 2.1 `OwnerController` with `+kubebuilder:validation:Enum=cli;controller`; `Platform.status.controllerVersion` (`ControllerVersion`) with the printer column `Controller`; `RenderInputKey.ControllerVersion` and the reconciler fields; doc comments say controller.
- [x] 2.2 Self-skip constants in `internal/reconcile/owninstance.go` (`opm-controller`, `opm-system`, `opmodel.dev/modules/opm_controller`) and their tests, including the `opm_controller_dashboard` prefix traps; condition text names the controller.
- [x] 2.3 Registry `ClientType` and the startup log `Starting opm-controller`; per-kind loops called reconcilers in prose.
- [x] 2.4 Kustomize namespace, `namePrefix`, `app.kubernetes.io/name` and placeholder image; base name `manager` for the Deployment, ServiceAccount, bindings, metrics Service and monitor.
- [x] 2.5 Fixtures on `testing.opmodel.dev/{modules,releases,catalogs}/controller/*` in `test/fixtures`, `config/samples`, Go fixture strings and `hack/fixtures.sh`; `test/integration/controllermodule`; e2e constants and the e2e deploy through `task controller:*`.
- [x] 2.6 `task dev:manifests dev:generate` and `task controller:installer` regenerate the CRDs, deepcopy, `modules/opm_controller/zz_generated_*` and `dist/install.yaml`.
- [x] 2.7 `task dev:fmt dev:vet dev:lint`, and `task dev:test:seeded` against a throwaway registry seeded from the tree, green; then commit `refactor(api)!: name the controller in the owner enum, status and install`.

## 3. The controller module

- [x] 3.1 `git mv modules/opm_operator modules/opm_controller`; sub-package `controller/controller.cue` (`package controller`); module path, package, identity, `_instance`, `_namespace`, `let P`, `let SA`, description and README; image repository `ghcr.io/open-platform-model/opm-controller` with the deployed version and digest marked `TODO(cutover)`.
- [x] 3.2 `RELEASE` 0.1.0, a fresh `CHANGELOG.md`; release-please component and package `opm_controller`, no manifest key, root `package-name: opm-controller`, both module paths excluded, `bootstrap-sha` kept.
- [x] 3.3 `hack/controller-module` (with `min-controller-version`) and `hack/testdata/controller-module-*`; image and branch regexes and messages.
- [x] 3.4 `task controller-module:release-check:test`, `publish:test`, `image:test`, `bot-pr:test`, `fmt:check` and `drift`, and `cue vet` of the module, green; then commit `refactor(module)!: move the controller's module to opmodel.dev/modules/opm_controller`.

## 4. The install manifest only from the module release

- [x] 4.1 `release.yml` `image-release`: drop "Render digest-pinned install manifest", "Upload install.yaml to the draft release" and the Task install; the job holds `contents: read`.
- [x] 4.2 `release-guard.sh publish` requires only `opm-examples.tar.gz` for a `v*` tag and `install.yaml` for an `opm_controller-v*` tag.
- [x] 4.3 actionlint and `task cascade:wiring:check` green; then commit `feat(release)!: publish the install manifest only from the module release`.

## 5. Workflows, tasks and the docs bundle

- [x] 5.1 `controller-module.yml`; every workflow literal (`lint.yml` release branch and manifest key, `release.yml` outputs, branch, repository guards and docs project, `module-image.yml` and `module-deps.yml` tar allow-paths, bot branch `module/controller-image`), `image-pr.yml`, `docs.yml`, `test*.yml`, CODEOWNERS.
- [x] 5.2 Task namespaces `controller` and `controller-module`, `controller:install` and `controller:uninstall` flattened, `controller:deploy` added; kind cluster and buildx builder names; `.tasks/cascade` (the `.github`-owned `wiring-check.sh` is left to the `.github` change).
- [x] 5.3 `docs-kit.cue`: bundle `opm-controller`, owning `reference/controller/` titled "Controller Reference", `stripLabels` and the reconciler wording; the `reconciledBy` test reads the new key.
- [x] 5.4 Delete the `Makefile`.
- [x] 5.5 actionlint, `task cascade:wiring:check` and `task docs:bundle:check` green; then commit `ci!: rename the workflows, tasks and docs bundle to the controller`.

## 6. Docs and ADRs

- [x] 6.1 `docs/site` pages renamed to `start/install-the-controller`, `diagnostics/controller-conditions` and `reference/controller/`; product prose says controller, per-kind loops reconcilers; the install page drops the manifest-migration note and points kubectl at a module release.
- [x] 6.2 ADR-021 "Rename opm-operator to opm-controller"; delete ADR-015; dated ADR bodies, `docs/RENDERING.md`, the design doc and the kind test plans follow the new names.
- [x] 6.3 README, AGENTS, CONSTITUTION and `openspec/config.yaml`.
- [x] 6.4 `task docs:bundle:check` green; then commit `docs!: name the controller in the docs, ADRs and agent guides`.

## 7. Specs and archive

- [x] 7.1 Recreate `operator-module`, `operator-module-release`, `operator-resource-reference` and `operator-version-identity` as `controller-*`; REMOVED + ADDED retitles and MODIFIED wording in every other spec that names the operator; fold the still-valid deltas of `stop-operator-install-manifest` and remove that change.
- [x] 7.2 `openspec validate rename-to-opm-controller --strict` and `openspec validate --all --strict` (no new failures) green, the residue search clean apart from history and the generic word; then commit `docs(openspec): propose rename-to-opm-controller`.
- [x] 7.3 Verify the change against its artifacts, archive it with `retire_capabilities`, restore the requirement order and Purpose wording in the touched main specs, and commit `docs(openspec): archive rename-to-opm-controller`.

## 8. Owner decisions of 2026-10-06 (folded after the archive)

- [x] 8.1 Rebase onto main after `add-healthy-condition`: rename its product words, imports and `testControllerVersion`, regenerate, move its rows into `controller-conditions.md`, re-derive the archived `module-instance-ownership`, `reconcile-loop-assembly` and `render-input-key` deltas from the new main specs, and add the `kubernetes-tier-adoption` retitle and the `instance-health` delta.
- [x] 8.2 Delete the own-instance finalizer release: the finalizer stripping in `handleOwnInstance`, its three tests, its docs sentences, and the requirement (REMOVED with no successor); commit `refactor(reconcile)!: delete the own-instance finalizer guard`. Narrowed by 8.6.
- [x] 8.3 Select the controller's pods by `manager`: `control-plane: manager` in `config/`, the module component key `manager`, the fixed selector in `test/integration/controllermodule`, the e2e selectors, the `controller-module` spec and ADR-021; regenerate `dist/install.yaml`; commit `refactor(module)!: select the controller's pods by manager`.
- [x] 8.4 Restore the recorded run results in `hack/kind-opm-dev-test` to their recorded wording, with a note that they predate the rename.
- [x] 8.5 Each spec edit is made in the archived delta and the main spec alike; every ADDED or MODIFIED block of the archived deltas equals its main-spec block, and no REMOVED title is left in a main spec.
- [x] 8.6 Narrow the finalizer release instead of deleting it: `handleOwnInstance` removes the finalizer from a deleting own instance without pruning, writes no status and emits a `SelfManagementRefused` event; a live own instance keeps the finalizer and the refusal; tests for an instance recognised by its inventory and for a module-path edit; the requirement restated as "The controller releases its finalizer from a deleting own instance without pruning"; the refusal requirement no longer claims an own instance never carries the finalizer; the `controller-conditions.md` and `deletion-and-pruning.md` sentences; commit `fix(reconcile): release the finalizer from a deleting own instance`.
- [x] 8.7 Release a deleting own instance before the owner-skip gate, whatever its `spec.owner`, so a refused own instance handed to the CLI as the refusal advises does not hang on delete; the owner-skip requirement names the exception; a test for refuse, set owner cli, delete that also checks the deleting reconcile writes no status; commit `fix(reconcile): release a deleting own instance whatever its owner`.
