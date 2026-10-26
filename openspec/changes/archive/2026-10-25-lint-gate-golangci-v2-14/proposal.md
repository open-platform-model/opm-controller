## Why

The controller must move to golangci-lint v2.14.0, the oldest release that reads Go 1.27 export data. The move is blocked in three places:

- The lint-config check runs `golangci-lint config verify --schema <committed file>`. The `--schema` flag no longer exists: since before v2.12.0 the schema is embedded in the binary and `config verify` makes no request. The repo rule says a dropped flag needs a new design, never a dropped check.
- v2.14.0 reports findings on today's code that v2.8.0 does not: `goconst`, `modernize`, `prealloc` and `staticcheck`.
- One finding is real: `ctrl.Result{Requeue: true}` is deprecated (SA1019).

## What Changes

- The check keeps its job and its offline guarantee. It drops the committed schema and `SHA256SUMS`, and runs `config verify` on the linter the version file names, with the proxy variables at a closed port.
- `.golangci-lint-version` and `.custom-gcl.yml` move to `v2.14.0`.
- `.golangci.yml` turns off, by name, each check that v2.14.0 adds or widens, so the effective rule set stays as it is today, except the modernize analyzers `bloop` and `fmtappendf`, which v2.14.0 removes.
- The two finalizer-add returns, and the tests that assert them, use `RequeueAfter: time.Second` in place of `Requeue: true`. One prealloc finding in a test is fixed in code.

## Capabilities

### Modified Capabilities
- `lint-config-check`: the check validates against the schema embedded in the pinned linter, not a committed file.
- `finalizer-and-deletion`: the finalizer-add reconcile asks for a requeue after one second, a fixed one second.

## Impact

`.golangci-lint-version`, `.custom-gcl.yml`, `.golangci.yml`, `.github/scripts/lint-config-check.sh` and its test, `.github/golangci-lint/` (deleted), `references/testing-and-lint.md`, `AGENTS.md`, two reconcile files and their tests. No Go version, module, workflow or other controller behaviour changes. One observable change: the controller-runtime metric `controller_runtime_reconcile_total` counts these two requeues under `result="requeue_after"` in place of `result="requeue"`.
