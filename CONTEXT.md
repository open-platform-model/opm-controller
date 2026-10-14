# opm-operator: where do I go?

Router only. Facts and rules are in `AGENTS.md`.

| Task | Read first | Then |
| --- | --- | --- |
| Change a CRD type, marker or RBAC marker (`api/v1alpha1`) | `CONSTITUTION.md` | `references/docs-bundles.md`, "Generated files that feed the reference" |
| Change reconcile logic, or anything that calls the Kernel | `CONSTITUTION.md` | `docs/RENDERING.md`, `openspec/config.yaml` |
| Write or fix a Go test, run one test, run e2e | `references/testing-and-lint.md` | `Taskfile.yml`, `.tasks/dev.yaml` |
| Fix a lint error, or move the golangci-lint version | `references/testing-and-lint.md` | `.golangci.yml`, `.custom-gcl.yml` |
| Add, bump or publish a test fixture, or a test needs a registry | `references/fixtures.md` | `.tasks/examples.yaml`, `hack/fixtures.sh` |
| Choose a commit type, break the API, or change the operator module | `references/release-units.md` | `release-please-config.json` |
| Recover a failed release or a stranded draft | `references/release-units.md` | `.github/workflows/release.yml` |
| Edit a workflow, a cascade pin, `.opm-cli-version` or `deps:cascade` | `references/release-cascade.md` | `.tasks/cascade/wiring-check.yaml`, `RELEASING.md` at the workspace root |
| Edit a site page, `docs-kit.cue` or a docs-kit pin | `references/docs-bundles.md` | `docs/site/`, workspace `STYLE.md` and `VOICE.md` |
| Add a controller or a new kind | `references/docs-bundles.md`, "`reconciledBy`" | `internal/controller/` |
| Record an architecture decision | `adr/TEMPLATE.md` | `adr/` |
| Start a feature change | `openspec/config.yaml` | `CONSTITUTION.md` |
| Deliver or cite an enhancement | `references/enhancements.md` at the workspace root | `openspec/config.yaml`, the `references:` entry |
| Commit, open a PR or write GitHub text | `references/git-text.md` at the workspace root | |
| Tag or publish by hand, or change tag rules | `references/releases.md` at the workspace root | |
| Run `cue` or touch a registry | `references/registry.md` at the workspace root | |
