## 1. Rules and the mechanical rewrite

- [x] 1.1 Update `openspec/config.yaml` (`references:`, `Delivers`, comment rule) and `AGENTS.md`
- [x] 1.2 Run `enh cite-rewrite --write` with a worklist; spot-read at least 40 edits
- [x] 1.3 Delete the 21 `enhancement.yaml` files and `enhancements/000-TEMPLATE/`; drop `citations: "link"` from `docs-kit.cue`
- [x] 1.4 Gates green (`task dev:fmt dev:vet`), then commit `docs: align the operator with the enhancements store` (one commit per place if the diff is large)

## 2. The worklist by hand

- [x] 2.1 Rewrite code rows (inline, possessive, string, retired ids) in `api/`, `internal/`, `test/`, `config/rbac`
- [x] 2.2 Rewrite live spec, archived and docs rows
- [x] 2.3 Regenerate CRDs, install manifest and the CUE copy
- [x] 2.4 Gates green (`task dev:manifests dev:generate dev:fmt dev:vet`), then commit `docs: rewrite the remaining enhancement citations by hand`

## 3. Verify

- [x] 3.1 `cite-lint` 0 findings and `cite-rewrite --dry-run` 0 edits; `git grep` of the two retired ids empty outside excluded paths
- [x] 3.2 `task dev:fmt dev:vet`, `task dev:manifests dev:generate docs:bundle:check`, `task operator:installer`, `task operator-module:drift`, `task dev:lint`, unit tests of touched packages, `openspec validate --all --strict --no-interactive`
- [ ] 3.3 Review, fix, archive the change, then commit `docs(openspec): archive align-enhancement-store`
