## Why

`.github` main is now `6df4000f6cabbf460a5cad6431b6bc9044fa3027`. Since this repo's pin `0f9c6ac` it carries `.github` PR 20: the cascade maps in `wiring/lib.sh` (`notify_targets`, `receiver_sources`, `g3_upstreams`, `publish_paths`, `mirror_sources`, `changelog_source` and the others) accept the repo name `opm-controller` beside `opm-operator`, and the resolver accepts both module trains `opm_operator-v*` and `opm_controller-v*`. GitHub renames a repo in place and every cascade job derives the repo name from `GITHUB_REPOSITORY`, so the name a run sees flips the moment the rename happens, whatever `.github` commit the caller pins. The maps fail closed on a name they do not know. This repo still pins `0f9c6ac`, so after the rename its cascade would refuse to run.

PR 20 changes `wiring/lib.sh`, which both `cascade-notify` and `cascade-publish` run, so the README's "Both actions" rule would make the library move first as the canary. The owner waived the canary rule (decision 37) once for this `.github` change, because no library release is pending. The library, this repo and the cli each move in their own pull request; the owner merges the library first.

## What Changes

- All six `open-platform-model/.github` references (`cascade-notify` in `release.yml`, `cascade-receive.yml` and `cascade-publish` in `deps-cascade.yml`, the `cascade-gates.yml` call, and the resolver checkout `ref:` in `cascade-task.yml` and `module-deps.yml`) move to `6df4000f6cabbf460a5cad6431b6bc9044fa3027 # .github main`.
- `.tasks/cascade/wiring-check.sh` stays as it is. PR 20 does not change it between `0f9c6ac` and `6df4000`, so the copy is already byte-identical; `--pin-on-main` compares it against the new pin.
- Nothing else. The mirrored `.tasks/cascade/{pins.sh,lib.sh,classes,cascade.sh}` on `origin/main` (`cb97ffb0`) hash to the four values `mirror_sources opm-operator` records at `6df4000`, so publish does not refuse this repo for mirror drift. The values in `.tasks/cascade/wiring-check.yaml` and `lint.yml` need no change: PR 20 changes no caller shape.
- The two `cascade-receiver` requirements that cite the `.github` README at `0f9c6ac` cite it at `6df4000`.

## Capabilities

### Modified Capabilities

- `cascade-receiver`: the publish-job and pin requirements cite the `.github` README at the new pin; behavior is unchanged.

## Impact

CI only (`ci`/`chore` commits; no operator or module release). Affected files: `.github/workflows/{release,deps-cascade,cascade-gates,cascade-task,module-deps}.yml`. `CASCADE_DRY_RUN` stays as it is. This change does not rename the repo and does not flip the target lists, which keep `opm-operator` until a later `.github` change.
