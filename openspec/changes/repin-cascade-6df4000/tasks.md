## 1. Move the pin

- [x] 1.1 Replace `0f9c6ac2c9b752a79f4874f637ef9955bcf00c13` with `6df4000f6cabbf460a5cad6431b6bc9044fa3027` on all six `.github` references in `.github/workflows`, keeping ` # .github main`. Verify: `grep -rn -A1 'open-platform-model/.github' .github/workflows` shows only the new SHA.
- [x] 1.2 Confirm `.tasks/cascade/wiring-check.sh` is the file at the new SHA: `gh api -H 'Accept: application/vnd.github.raw' "repos/open-platform-model/.github/contents/.github/scripts/cascade/wiring-check.sh?ref=6df4000f6cabbf460a5cad6431b6bc9044fa3027" | cmp - .tasks/cascade/wiring-check.sh`.
- [x] 1.3 Confirm `.tasks/cascade/wiring-check.yaml` and `lint.yml` need no change against the README at the new SHA, and that `.tasks/cascade/{pins.sh,lib.sh,classes,cascade.sh}` on `origin/main` hash to `mirror_sources opm-operator` in `wiring/lib.sh` at the new SHA.
- [x] 1.4 Verify: `task cascade:wiring:check`, `actionlint` and `task dev:lint` pass; with network, `GH_TOKEN=$(gh auth token) bash .tasks/cascade/wiring-check.sh --pin-on-main` prints `cascade wiring: ok`. Commit `ci(deps): pin the cascade to .github 6df4000`.

## 2. Verify and archive

- [ ] 2.1 `openspec validate repin-cascade-6df4000 --strict` passes and the verify skill reports no CRITICAL finding.
- [ ] 2.2 `openspec archive repin-cascade-6df4000 -y`; `openspec validate --all --strict` fails only where it fails on `origin/main`. Commit `chore(openspec): archive repin-cascade-6df4000`.
