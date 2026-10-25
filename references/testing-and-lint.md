# Tests, lint and the linter version

Covers the three test tiers, how to run one test, what the lint config enforces, and how to move the golangci-lint version. Fixtures and registries are not covered here.

## Contents

- Test tiers
- Run one test
- Lint expectations
- Moving the golangci-lint version
- Pitfalls
- Check before you finish

## Test tiers

Pick the lightest tier that proves the behaviour.

| Tier | Where | What it proves | Runs in |
| --- | --- | --- | --- |
| Unit | `internal/*/` | Single-package logic | `task dev:test` |
| Integration | `test/integration/` (`apply`, `crdvalidation`, `controllermodule`, `reconcile`) | Cross-package behaviour against a real API server (envtest) | `task dev:test` |
| E2E | `test/e2e/` (`e2e` build tag) | The deployed controller on a Kind cluster | `task dev:e2e` |

- `task dev:test` runs unit and integration tests and excludes `/test/e2e`. No Kind test runs in the default path.
- E2E needs Kind. It may install cert-manager unless `CERT_MANAGER_INSTALL_SKIP=true`. Say so in your report when you skipped e2e for lack of Kind or a cluster.
- Style: Ginkgo v2 and Gomega. Give `Describe`, `Context` and `It` descriptive text so `-ginkgo.focus` works. Use `Eventually` for async Kubernetes behaviour, never sleeps. Assert errors with `Expect(err).NotTo(HaveOccurred())` or `Expect(...).To(Succeed())`. Put repeated assertions in package-local helpers.
- `test/integration/controllermodule` renders the controller module (`modules/opm_controller`) against a platform generated from its own catalog pin. It needs core and the catalog from `CUE_REGISTRY` (GHCR under `task dev:test`), the `cue` CLI and the kustomize binary that `task dev:test` passes as `KUSTOMIZE`. Without them its specs skip; under `OPM_TEST_REGISTRY_FORCE=1` (PR CI) they fail. Its minimum-controller-version spec needs none of them.
- `test/integration/crdvalidation` proves that every sample in `config/samples` is admitted by the CRDs.

## Run one test

`task dev:test:one` runs one package or suite entrypoint (vars `TEST_PKG`, `RUN`, `FOCUS`) with the envtest assets set. The raw `go test` forms follow, for e2e and for hand runs. The envtest binaries are installed as a dependency of `task dev:test`: run it once to provision them. The Kubernetes version is derived from `k8s.io/api` in `go.mod` by `ENVTEST_K8S_VERSION` in `Taskfile.yml`.

```bash
export KUBEBUILDER_ASSETS="$(./bin/setup-envtest use <k8s-minor> --bin-dir "$PWD/bin" -p path)"
go test ./internal/controller                                   # one package
go test -tags=e2e ./test/e2e -v -ginkgo.v                       # e2e only
KIND_CLUSTER=opm-controller-test-e2e go test -tags=e2e ./test/e2e -run <suite-func> -ginkgo.focus="should run successfully" -v -ginkgo.v
```

`<k8s-minor>` is the `1.N` that `task dev:test` prints in its envtest line. `KUBEBUILDER_ASSETS` must be an absolute path.

## Lint expectations

- `.golangci.yml` enables, among others: `errcheck`, `ginkgolinter`, `gocyclo`, `govet`, `misspell`, `modernize`, `revive`, `staticcheck`, `unused`, `depguard`, `logcheck`. Formatters: `gofmt` and `goimports`.
- `logcheck` (a custom plugin built from the custom-gcl config) validates structured logging and balanced key/value parameters.
- `lll` and `dupl` are relaxed in `api/*` and `internal/*`. Do not rely on the exclusions without need.
- The `depguard` rule keeps the controller from copying the library's Kubernetes tier decisions (the comment in `.golangci.yml` says why).
- `ginkgolinter` flags non-idiomatic Ginkgo.

## Moving the golangci-lint version

`.golangci-lint-version` is the only place that names the linter version. `Taskfile.yml` and `Makefile` read it. The custom-gcl config at the repo root repeats it because `golangci-lint custom` reads only that file; `task dev:lint:config` refuses a tree where the two differ.

The linter embeds the JSON schema of its own version, so `golangci-lint config verify` makes no request and the repository holds no schema file and no checksum file. The check runs it with the proxy variables at a closed local port. The binary that runs is the plugin build, which `golangci-lint custom` clones from the git tag of the version.

- **Any move** (patch, minor or major): edit `.golangci-lint-version` and the `version:` line of the custom-gcl config. Then run the commands below and read the new findings.
- `task dev:lint:config dev:lint:config:test dev:lint` must pass.
- A new linter version can report on code the old one passed. Turn off the new check by name in `.golangci.yml` with a comment, or fix the code; do not leave a red `task dev:lint`.
- The scenario `value-of-this-line` in the lint-config test in `.github/scripts/` proves the embedded schema is the one of the pinned line. On a minor move, change it to a value the new line added.
- If a release changes `config verify` so that it needs the network or drops its offline check, the check fails on the bump PR and needs a new design, never a dropped check.
- Installing the linter needs the network (the Go module proxy and the plugin build), and `task dev:lint*` reinstalls it on every run. Only the configuration check is offline: with the linter in `bin/`, `GOLANGCI_LINT=bin/golangci-lint bash .github/scripts/lint-config-check.sh` runs without a network.

## Pitfalls

- A relative `KUBEBUILDER_ASSETS` makes envtest suites fail to find the binaries.
- Editing `.golangci-lint-version` without the custom-gcl config fails `task dev:lint:config`.
- `config verify` does not check linter names: `nosuchlinter` in `enable` passes it, and `golangci-lint run` refuses it.
- Registry-backed specs skip silently without `CUE_REGISTRY` outside PR CI; check the skip count.

## Check before you finish

- [ ] `task dev:fmt dev:vet dev:test` pass after a Go change.
- [ ] `task dev:lint` passes for a non-trivial change.
- [ ] After a linter move, `task dev:lint:config dev:lint:config:test dev:lint` pass.
- [ ] You said whether e2e ran or was skipped.
