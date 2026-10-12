## 1. Move the pin

- [x] 1.1 Replace `6df4000f6cabbf460a5cad6431b6bc9044fa3027` with `2b66da51fede4332969c324b069012632d8edea9` on every `.github` reference in `.github/workflows`, keeping ` # .github main`. Replace `.tasks/cascade/wiring-check.sh` with the file at the new SHA. Commit `ci(deps): pin the cascade to .github 2b66da5`.

## 2. Wire the check

- [x] 2.1 Add `.github/workflows/cite-lint.yml` and the `extra-references` item in `.tasks/cascade/wiring-check.yaml`. Commit `ci: call the shared cite-lint workflow`.
- [x] 2.2 Add `task cite-lint` to `Taskfile.yml` and list it in `AGENTS.md`. Commit `build: add task cite-lint`.

## 3. Verify and archive

- [ ] 3.1 `task cite-lint`, `task cascade:wiring:check`, `openspec validate wire-cite-lint --strict`, actionlint and zizmor on the changed workflows pass.
- [ ] 3.2 Archive the change. Commit `chore(openspec): archive wire-cite-lint`.
