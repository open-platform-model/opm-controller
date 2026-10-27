## Context

The controller resolves core, catalogs and modules through CUE's registry client. That client finds registry logins in `$DOCKER_CONFIG/config.json`. The controller Deployment sets no `DOCKER_CONFIG`, has no HOME, and has a read-only root filesystem with `/tmp` as the only writable path. So the client finds no login. The approved design (owner decision) adds one module value that mounts a Secret the administrator created and points `DOCKER_CONFIG` at it. This document records how, what was verified, the test plan and the PR split.

## Research & Decisions

### Where the logins come from
**Context**: A credential must reach the controller process under a read-only root filesystem and without a new RBAC rule.
**Explored**: Map of every registry path in the controller (two client builds, one shared Kernel client and one Platform closure client rebuilt per reconcile) and how other controllers take pull credentials (Flux, Argo CD, Crossplane, Kyverno).
**Decision**: A kubelet-mounted `kubernetes.io/dockerconfigjson` Secret, read through `DOCKER_CONFIG`. One credential set for the whole controller. The administrator creates the Secret. The owner approved this.
**Rationale**: The kubelet mounts the Secret, so the controller needs no Secret RBAC and the manager role stays as it is. The login is a file, not an env var and not an argument, so it appears in no process listing, no startup log line and no instance value. A Secret read through the API, a per-instance `secretRef`, and a per-host mapping were left out by the owner. A per-instance reference can come later without a break.
**Reversibility**: high (hours). One optional value; removing it later is a breaking change only for users who set it.

### The value name and shape
**Decision**: `registryAuth.secretName`, a string. The owner approved the name; the design left it as a proposal. `registryAuth` is an object so that a later field (for example a key name) can join without a new top-level value.
**Rule**: the name must be a valid Kubernetes object name, so a typo fails at render with the value named, not at the kubelet. The pattern is the DNS-1123 subdomain: lowercase letters, digits, `-` and `.`, start and end alphanumeric, at most 253 characters. The empty string is refused like `registry`.

### What the module renders
When `registryAuth` is set, and only then, on the `manager` component:

```yaml
env:
  - {name: DOCKER_CONFIG, value: /etc/opm/docker}
volumeMounts:
  - {name: registry-auth, mountPath: /etc/opm/docker, readOnly: true}
volumes:
  - name: registry-auth
    secret:
      secretName: <registryAuth.secretName>
      items:
        - {key: .dockerconfigjson, path: config.json}
```

How, in the module: the catalog volume source `secretRef` (name, items) renders `secret.secretName` and `secret.items`, and renders `optional` only when it is set (catalog `container_helpers.cue`, the `secretRef` branch). The module never sets `optional`, so the field is absent and the Kubernetes default (`false`) applies. A missing Secret, or a Secret without the key `.dockerconfigjson`, therefore keeps the pod in `ContainerCreating` with a `FailedMount` event. Section 2 confirms the rendered Deployment shows no `optional`.

Module sketch (not final code):

```cue
// #config
registryAuth?: secretName: string & !="" & =~"^[a-z0-9]([-a-z0-9.]*[a-z0-9])?$" & strings.MaxRunes(253)

// components.cue, inside "manager"
if #config.registryAuth != _|_ {
    spec: volumes: "registry-auth": {
        secretRef: {name: #config.registryAuth.secretName, items: [{key: ".dockerconfigjson", path: "config.json"}]}
        readOnly: true
    }
    // container: env: DOCKER_CONFIG: value: "/etc/opm/docker"
    // container: volumeMounts: "registry-auth": spec.volumes."registry-auth" & {mountPath: "/etc/opm/docker"}
}
```

The existing `tmp` mount uses the same `spec.volumes.<name> & {mountPath}` pattern. `extraArgs` and its guards do not change. The module sets no `defaultMode`: the Pod runs as a non-root user and the default file mode `0644` keeps the file readable by it. A `0400` mode would make the file unreadable, because the file is owned by root and the Pod has no `fsGroup`. A test reads the mode (see tasks).

Mount rules from the secrets guidance, all met: the Secret is mounted as a read-only file, in the one container that needs it, as a whole volume (not a `subPath`, which would never update), and no env var carries the secret.

### Verified: no Go change is needed
- `cueenv.Override` returns nil when the mapping and cache directory are empty (the process environment is used), else a copy of `os.Environ()` in which only `CUE_REGISTRY` and `CUE_CACHE_DIR` are replaced: `library/opm/internal/cueenv/cueenv.go:33-41`. `DOCKER_CONFIG` passes through in both cases. A search of `library/opm` for `DOCKER` finds no Go code that touches it.
- The Platform closure client is built from `os.Environ()`: `opm-controller/internal/controller/platform_controller.go:418-427` (taken from exploration; section 1 re-checks it at the pinned commit).
- CUE v0.17.1 `modconfig` loads logins with `ociauth.LoadWithEnv(nil, t.cfg.Env)`: `cuelang.org/go@v0.17.1/mod/modconfig/modconfig.go:331`. The pinned `ociauth` (`cuelabs.dev/go/oci/ociregistry@v0.0.0-20260717083115-5eb5795f322a`) reads `$DOCKER_AUTH_CONFIG`, then `$DOCKER_CONFIG/config.json`: `ociauth/authfile.go:136-169`.
- Not verified by running: that a mounted Secret authenticates in the controller. Section 1 and section 4 prove it.

### Behaviour after a refused pull
**Trace**: The Kernel's shared registry client reads the logins once, when it is built. A failed registry call drops the client, and the next operation builds a new one, which reads the file again (`library/opm/internal/cueenv/registry.go:20-34`, `failed` and `drop`, lines 113-122 and 158-162). The Platform closure client is rebuilt for each reconcile.
**Expected result**: after a refused pull, the next pull reads the (possibly updated) file again. The kubelet updates a mounted Secret volume in place within about a minute (not verified here; section 1 tests the file replacement, not the kubelet).
**Consequence if true**: restart is needed only while the old login still works. A revoked login stops working, the first pull fails, and the next pull uses the new file. The docs still say: restart to rotate, because that is the approved and deterministic way.
**Test (section 1)**: see Test plan. The test result is the source for the docs sentence. If the test disproves the expectation, section 1 rewrites the requirement "A pull after a refused pull reads the logins again" and this subsection to match, before section 2 starts.

### Exposure
- Any ModuleInstance author can make the controller pull any path the shared login reaches (`spec.module.path` is user-controlled). This is a decision of the owner, stated in the docs: the token should have `read:packages` on the packages that the cluster's users may use, and no more.
- The Secret is readable by whoever has `get` on Secrets in `opm-system`. The controller does not need that right: it gets the file from the kubelet.
- Base64 in the Secret is not encryption. Encryption at rest is a cluster matter. The docs say so in one sentence.
- The mapping string is logged at startup and is in the Deployment args; it must never carry a login. The new value keeps logins out of it.
- The controller's own image pull is `imagePullSecrets`, not this value.

### Release and PR split
The repo rule (`AGENTS.md`): "A PR that changes `modules/opm_controller` changes nothing else. Otherwise its squash commit releases both units." And: "Commit type decides the release. `docs`, `chore`, `test`, `ci` and `build` release nothing." The module has its own release component in `release-please-config.json`, and the controller package excludes `modules/opm_controller`, so a `feat` that touches only the module releases only the module.

So the change ships as five PRs, one per section, in this order (proposal line `Delivery: one PR per section`). Each section is also one commit and one PR title.

| Section | Touches | PR title type | Releases |
| --- | --- | --- | --- |
| 1 Spike: login and refused pull | `test/integration/reconcile/`, this change folder | `test(controller)` | nothing |
| 2 The module value | `modules/opm_controller/` only: `module.cue`, `components.cue`, `README.md` | `feat(module)` | the module, MINOR |
| 3 Render tests | `test/integration/controllermodule/` | `test(module)` | nothing |
| 4 End to end on Kind | `test/e2e/`, `.tasks/` if a task is needed | `test(e2e)` | nothing |
| 5 Docs | `docs/site/` | `docs(site)` | nothing |

Why this order: section 1 runs first because design.md carries unverified assumptions (the mount authenticates, the refused-pull behaviour). Section 2 must stand alone because of the repo rule. Section 3 follows it because its set-value specs fail without the value. Docs come last because they state the result of section 1 and name the released value; the docs PR merges after the module release is published, so the page never names a value no release has. The module README is in section 2 (it lives under `modules/opm_controller/`, and the README of a module that ships a value belongs to the same PR).

Each section ends green: section 2 leaves every existing spec green, because the unset render does not change; the kustomize parity spec compares the unset render.

`controller-module:drift` stays green in every section: no `config/` file changes and no generated file changes.

Note: the approved design says "a controller release ships it". The value lives in the module, so a module release ships it; the controller binary does not change and no controller release is needed. Reported as a deviation.

### Test plan

| Test | Section | Where | What it proves |
| --- | --- | --- | --- |
| Unset renders nothing new | 3 | `test/integration/controllermodule/` | With no `registryAuth`, the Deployment has no `DOCKER_CONFIG` env, no `registry-auth` volume and no `registry-auth` mount, and the existing parity specs stay green. |
| Set renders the three parts | 3 | same | With `registryAuth: secretName: "opm-registry-auth"`: env `DOCKER_CONFIG=/etc/opm/docker`; one read-only mount at `/etc/opm/docker`; one volume from the Secret with the item `.dockerconfigjson` to `config.json`. |
| A missing Secret is not optional | 3 | same | The rendered volume has no `optional` field, or `optional` is `false`. |
| No module-owned Secret, no RBAC | 3 | same | The render has no `Secret` object, and the rendered ClusterRoles gain no `secrets` rule. |
| The file is readable by the Pod user | 3 | same | The rendered volume has no `defaultMode` below `0444`. |
| Bad names are refused | 3 | same | The empty string, an uppercase name and a name over 253 characters fail at `registryAuth.secretName`. |
| Guards stay | 3 | same | `extraArgs` refusals and Pod Security restricted specs pass with the value set. |
| A pull that needs a login works | 1 | `test/integration/reconcile/` | An in-process OCI registry behind HTTP basic auth (the `ociserver` and `ocimem` packages that the pinned `cuelabs.dev/go/oci/ociregistry` ships, wrapped in an auth handler). Without a login, the Kernel's module pull is refused. With `DOCKER_CONFIG` set to a directory holding `config.json` for the registry host, the same pull succeeds. Catalog pull uses the same client; one catalog pull is asserted too. |
| After a refused pull | 1 | same | Start with a wrong login in `config.json`. The pull is refused. Replace the file with the right login (write a new file and rename it, as the kubelet does). Pull again with the same Kernel, with no restart. Record whether it succeeds. Also record a pull that succeeded, then a revoked login: the first later pull fails or not. The result becomes the assertion and the docs sentence. |
| The mounted Secret authenticates | 4 | `test/e2e/` | On Kind, a registry that requires a login runs in the cluster. The test creates a `dockerconfigjson` Secret in `opm-system`, applies the three parts rendered by the module with the value set to the running controller Deployment, and checks a ModuleInstance whose module is only in that registry reaches Ready. A second instance, with the Secret deleted and the pod restarted, shows the pod stuck in `ContainerCreating` with a `FailedMount` event (the missing Secret blocks the pod). |

Fixtures use `testing.opmodel.dev`, never an `opmodel.dev/*` path. The test registry host uses `+insecure` in the mapping. No test publishes to the real GHCR and no test needs a token.

### Reconcile phase impact
None. Source, Render, Apply, Prune and Status are unchanged. A refused pull is the existing `ResolutionFailed` path with backoff.

## Risks / Trade-offs

- A reviewer of the module PR sees no test in it (the repo rule forces tests into section 3). Mitigation: the module PR stays small, section 3 follows at once, and section 1 already proves the mechanism.
- Section 4 applies rendered parts to a Deployment that `config/default` installed, not an install made by the module. That is the closest test the existing e2e suite allows without moving it to a module install. If task 4.1 finds the suite can install through the module, it uses that.
- A shared login serves every tenant (see Exposure).
- A Secret of the wrong type or without the key keeps the pod pending. That is the wanted block.

## Open Questions

- None for the owner. Section 1 may change the sentence about reload; that is a result, not a decision.
