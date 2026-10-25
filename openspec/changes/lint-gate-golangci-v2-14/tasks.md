## 1. Gate and rule set

- [ ] 1.1 Set `.golangci-lint-version` and `.custom-gcl.yml` to `v2.14.0`.
- [ ] 1.2 Rewrite `.github/scripts/lint-config-check.sh`: run `config verify` with the closed-port proxies; drop the schema, checksum, minor-line and "verify without --schema" rules; keep the other rules; update its header comment.
- [ ] 1.3 Delete `.github/golangci-lint/` (schema and `SHA256SUMS`).
- [ ] 1.4 Update `.github/scripts/lint-config-check-test.sh`: remove the schema and direct-verify scenarios, change `value-of-this-line` to the revive rule `use-slices-concat`, add a scenario that refuses a returned schema or checksum file, drop the schema copy in `fresh`.
- [ ] 1.5 Update `.golangci.yml`: `goconst.exclude-types: [Call, CompositeLit]` and the `modernize.disable` additions of the design, each with a comment.
- [ ] 1.6 Update `references/testing-and-lint.md` (move procedure, pitfalls, the `--schema` note) and the lint lines of `AGENTS.md` and the `dev:lint:config` task descriptions.
- [ ] 1.7 Rewrite the Purpose line of `openspec/specs/lint-config-check/spec.md` (it names the committed schema) when the change is archived.
- [ ] 1.8 Run `task dev:lint:config dev:lint:config:test`, and `unshare -rn` on the check; run `task dev:lint` and expect only SA1019 and prealloc findings.
- [ ] 1.9 Commit as `ci(lint): move golangci-lint to v2.14.0 and verify the config without a schema file`.

## 2. Requeue and prealloc

- [ ] 2.1 Add a constant next to `FinalizerName` (`internal/reconcile`) and use `ctrl.Result{RequeueAfter: <constant>}` at `moduleinstance.go` and `modulepackage.go`.
- [ ] 2.2 Update the 15 test assertions of `Requeue: true`; fix the prealloc site in `internal/apply/guard_test.go`.
- [ ] 2.3 Run `task dev:fmt dev:vet dev:test` and `task dev:lint`; expect 0 issues.
- [ ] 2.4 Commit as `refactor(controller): request the finalizer requeue with RequeueAfter`.
