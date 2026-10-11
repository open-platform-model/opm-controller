## 1. Doc comments and generated copies

- [x] 1.1 Correct the `Provides` doc comment in `api/v1alpha1/transformerregistration_types.go`
- [x] 1.2 Regenerate with `task dev:manifests dev:generate` and `task operator:installer`; check `task operator-module:drift`
- [x] 1.3 `task dev:fmt dev:vet` green and `gofmt -l` clean, then commit `docs: say that acceptance compares an empty provides list like any other`

## 2. Live specs

- [x] 2.1 Edit the Purpose of `platform-crd` by hand
- [x] 2.2 `openspec validate fix-stale-specs --strict` passes
- [ ] 2.3 Review, fix, archive the change (the archive syncs the deltas into `openspec/specs/`), run the gates again, then commit `docs(openspec): archive fix-stale-specs`
