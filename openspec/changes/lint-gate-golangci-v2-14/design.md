## Context

Today `.github/scripts/lint-config-check.sh` runs `config verify --schema <committed v2.8 schema>`, with a checksum file and a scenario test. Archived history: `lint-without-network-schema`. At v2.14.0 the flag is gone: `config verify` validates against the schema embedded in the binary (`pkg/commands/config_verify.go` loads `jsonsch.NextSchema` through an embedded loader) and makes no request.

## Goals / Non-Goals

**Goals:** the lint-config check still proves `.golangci.yml` valid and still runs offline on v2.14.0; the version file, the plugin build file and the check agree; the effective rule set stays as today; `ctrl.Result{Requeue: true}` is replaced.
**Non-Goals:** the Go version, any module, new lint rules, the install path (it needs the network, as before), unpinned `version: latest` of the logtools plugin in `.custom-gcl.yml` (existing, noted for the owner).

## Research & Decisions

### The check without `--schema`

**Explored**: on a stock v2.14.0 and on the custom build with `logcheck`, `unshare -rn <bin> config verify` exits 0 on the repo config. A config with an extra top-level key exits 3 and names it (`additional properties 'not-a-golangci-key' not allowed`). `--schema` fails with `unknown flag`. v2.8.0 in the same namespace fails with a dial error: it downloads. The linter names are not enum-checked by the schema (`nosuchlinter` passes verify; `run` rejects it).
**Decision**: the script runs `"$linter" config verify` with the proxy variables at the closed port, after the exact-version check. It drops the schema file, `SHA256SUMS`, the schema-line rules, and the "verify without --schema" line rule.
**Rationale**: the schema is part of the linter build. `go install` fetches the stock binary under the Go checksum database. The binary that the check runs is the plugin build: `golangci-lint custom` clones the git tag `v2.14.0` of the linter, and no checksum database covers that clone. The schema therefore has the provenance of that tag, not the provenance of a verified module. The deleted `SHA256SUMS` bound only the bytes of the committed file, so no control of the repository is lost. A linter old enough to download fails the closed-port guard at once, so the offline guarantee needs no extra file. The rule "no bare `config verify`" has no purpose now and is removed, not weakened.
**Reversibility**: cheap.

### The scenario test

Scenarios for the schema file, checksum, stale schema and missing minor schema end. `value-of-this-line` stays by name and uses the revive rule name `use-slices-concat`, which the v2.14.0 schema accepts and the v2.13.0 schema does not. A scenario for a committed schema file or checksum file that returns refuses the tree. `direct-verify-*` scenarios end. Kept: version file, literals, plugin file, workflow step, action, installed version, closed-port guard, invalid config, usage.

### Keeping the rule set

Findings of v2.14.0 on the tree, 891 on Go 1.27 (484 on 1.26.10), v2.8.0 reports 0. Verified by running a custom v2.14.0 build with the proposed config: 18 findings on Go 1.26.10 and on Go 1.27.3.

| Linter | Count | Cause | Handling |
| --- | --- | --- | --- |
| goconst | 461 | goconst 1.11 counts strings in composite literals (new type `CompositeLit`); 1.8.2 did not | `goconst.exclude-types: [Call, CompositeLit]`; `Call` is the default and stays |
| modernize | 412 on 1.27 (5 on 1.26) | `embedlit` (406), `errorsastype`, `slicesbackward`, `stringscut` | add to `modernize.disable`; also the new analyzers that stay silent now: `atomictypes`, `importcomment`, `reflecttypeassert`, `slicesclip` |
| prealloc | 1 | prealloc 1.1.0; no setting reverts it (`simple: false` gives 15) | fix in code: one test line |
| staticcheck | 17 | SA1019 on `Result.Requeue` | fix in code |

`errorsastype` is new in the modernize suite of v2.14.0; why `stringscut` reports now was not traced. The v2.14.0 schema checks analyzer names against an enum, so `config verify` refuses a typo in `modernize.disable`.

**Rule-set choice of this change.** The rule set of this move is: every check that v2.14.0 adds or widens and that reports on the tree is turned off by name, except SA1019 and the one prealloc finding, which are fixed in code. Turned off: goconst on composite literals, and the modernize analyzers `embedlit`, `errorsastype`, `slicesbackward`, `stringscut`, plus the silent new ones `atomictypes`, `importcomment`, `reflecttypeassert`, `slicesclip`. Lost with the move and not keepable by any setting: the modernize analyzers `bloop` and `fmtappendf`, which the x/tools version of v2.14.0 removes. Left on: staticcheck SA9010, which honnef.co/go/tools adds and which golangci-lint enables with the default `all`; it reports nothing on the tree. This is a decision of this move, not a standing rule for later moves.

### Replacing `Requeue: true`

controller-runtime 0.25.2 (`pkg/reconcile/reconcile.go`) says: use `RequeueAfter`. There is no zero-delay form: `RequeueAfter > 0` is required. `Requeue: true` enqueues rate limited (`pkg/internal/controller/controller.go`). The controller's limiter is exponential from 1 second (`moduleinstance_controller.go`, `modulepackage_controller.go`), and a key with no failures waits one second. So `RequeueAfter: time.Second` has the same first delay.
**Differences**: the delay is a fixed one second, with no backoff. It equals the old delay only for a key with no earlier failures. A key with earlier failures (for example, one failed finalizer patch) waited 2 seconds, up to 5 minutes, with `Requeue: true`; now it waits one second and its failure count is forgotten. A finalizer that never sticks loops at one second and no longer backs off. The reconcile is never slower. The metric `controller_runtime_reconcile_total` counts these returns under `result="requeue_after"` in place of `result="requeue"`. The owner accepted both differences.
**Alternative rejected**: `//nolint:staticcheck` keeps the field; the field is deprecated and the project removes deprecated use in code.
The live requirement "The finalizer-add reconcile requests an immediate requeue" keeps its heading and scenario name; its body now names the one-second delay.

## Risks / Trade-offs

A disable name that the schema accepts but that does nothing lowers coverage unseen; the counts are compared at implement. The Go 1.27 move needs no further lint config after this change.
