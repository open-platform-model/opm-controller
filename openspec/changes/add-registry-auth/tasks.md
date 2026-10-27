## 1. Spike: a login works, and what a refused pull does

- [ ] 1.1 Add an in-process OCI registry behind HTTP basic auth to `test/integration/reconcile/` (the `ociserver` and `ocimem` packages of the pinned `cuelabs.dev/go/oci/ociregistry`, wrapped in an auth handler on `httptest`), a helper that writes a `config.json` for its host into a temp directory, and a published `testing.opmodel.dev` fixture module in it. Check that the `ociregistry` packages are importable from `go.mod` without a new dependency; if a new `go.mod` line is needed, stop and report it.
- [ ] 1.2 Test: a module pull through the Kernel is refused without a login and succeeds with `DOCKER_CONFIG` set to the temp directory (t.Setenv; the mapping uses `+insecure` for the host). Add one catalog pull through `AcquireCatalogFromRegistry`.
- [ ] 1.3 Test: wrong login refused, then replace `config.json` (write a new file, rename it over the old one), then pull again with the same Kernel and no restart; assert the observed result. Add the case "pull succeeded, then the registry revokes the login": record whether the next pull fails. Assert what is observed.
- [ ] 1.4 Re-check at the pinned commit that the Platform closure client is built from `os.Environ()` (`internal/controller/platform_controller.go`) and write the result into design.md "Verified".
- [ ] 1.5 Write the observed reload result into design.md ("Behaviour after a refused pull") and the requirement "A pull after a refused pull reads the logins again". If the result differs from the expected one, rewrite that requirement and its scenarios to the observed behaviour. Do not change any other requirement.
- [ ] 1.6 `task dev:fmt dev:vet dev:lint dev:test` green, then commit `test(controller): prove a registry login and the pull after a refused pull`

## 2. The module value

- [ ] 2.1 `modules/opm_controller/module.cue`: add `registryAuth?: secretName: string & !="" & =~"^[a-z0-9]([-a-z0-9.]*[a-z0-9])?$" & strings.MaxRunes(253)` to `#config` with a doc comment for the README reader. Do not touch `extraArgs` or its guards.
- [ ] 2.2 `modules/opm_controller/components.cue`: inside the `manager` component, when `#config.registryAuth` is set, add the volume `registry-auth` (`secretRef` with the name and the item `.dockerconfigjson` to `config.json`, `readOnly: true`, no `optional`, no `defaultMode`), the mount at `/etc/opm/docker` (the `spec.volumes.<name> & {mountPath}` pattern of `tmp`) and the env `DOCKER_CONFIG`. Add nothing when it is unset. Confirm the pinned catalog in `cue.mod/module.cue` renders `secretRef`; if it does not, stop and report (a pin bump is a separate PR and `task deps:update`, never a hand edit).
- [ ] 2.3 `modules/opm_controller/README.md`: add `registryAuth.secretName` to the values table, with the Secret command and the restart note.
- [ ] 2.4 Render the module with the value unset and set (the module's own render task or `cue export`) and compare with the design's YAML. Confirm `git status` shows only files under `modules/opm_controller/`.
- [ ] 2.5 `task dev:fmt dev:vet dev:lint dev:test controller-module:drift` green, then commit `feat(module): take registry logins from a Secret through registryAuth.secretName`

## 3. Render tests

- [ ] 3.1 `test/integration/controllermodule/`: unset renders no `DOCKER_CONFIG`, no `registry-auth` mount, no `registry-auth` volume.
- [ ] 3.2 Set renders the env, the read-only mount and the volume with the item, as in design.md.
- [ ] 3.3 The volume has no `optional: true` and no `defaultMode` below `0444`.
- [ ] 3.4 The render has no `Secret` object, and no rendered role has a rule on `secrets`.
- [ ] 3.5 Bad names are refused (empty, uppercase, over 253 characters), naming `registryAuth.secretName`; add them to the existing refusal table.
- [ ] 3.6 With every value set together, the selector and arguments are unchanged and the pods are admitted under Pod Security restricted (extend `everyValue` and the restricted spec).
- [ ] 3.7 `task dev:fmt dev:vet dev:lint dev:test` green, then commit `test(module): cover the registryAuth render, set and unset`

## 4. End to end on Kind

- [ ] 4.1 Read how the e2e suite deploys the controller. Add a registry that requires a login (a `registry` image with htpasswd auth) as a Deployment and Service in a test namespace, publish the `testing.opmodel.dev` fixture module to it, and create the `dockerconfigjson` Secret in `opm-system`. Use no real token and no GHCR.
- [ ] 4.2 Test: render the module with `registryAuth` set, apply the three parts to the controller Deployment, point `--registry` at the test registry with `+insecure`, and show a ModuleInstance of the fixture module reaches Ready. A control instance with no Secret mounted shows the resolution failure.
- [ ] 4.3 Test: delete the Secret, restart the pod, and show it stays in `ContainerCreating` with a `FailedMount` event naming the Secret.
- [ ] 4.4 `task dev:fmt dev:vet dev:lint dev:test` green; run `task dev:e2e` and record the result (say if Kind was not available), then commit `test(e2e): pull a module from a registry that needs a login`

## 5. Docs

- [ ] 5.1 `docs/site/start/install-the-controller.md`: a step for private registries after the registry step: create the Secret, set `registryAuth.secretName` in the values file, install, rotate (update the Secret, restart `deployment/opm-controller-manager` in `opm-system`; check the Deployment name in the module), what a missing Secret looks like, and the sentence on the refused-pull result from section 1. State that the logins are shared by every instance and that the controller image pull uses `imagePullSecrets`.
- [ ] 5.2 Update the module-values sentence of the install page so it lists `registryAuth`.
- [ ] 5.3 `task docs:bundle:check` green and `task cite-lint` green, then commit `docs(site): document registry logins for the controller`. Merge this PR after the module release that carries the value is published.
