#!/usr/bin/env bash
# Linter configuration check, offline. golangci-lint config verify validates
# .golangci.yml against the JSON schema embedded in the linter binary, so the
# command makes no request. The check runs it with the proxy variables at a
# closed local port: a linter old enough to download its schema (before the
# schema was embedded) fails at once.
#
# One file, .golangci-lint-version, names the linter version that
# .tasks/tools.yaml installs. This check refuses a tree where that file, the
# plugin build file and the use of the check disagree:
#
#   - .golangci-lint-version is vX.Y.Z;
#   - .github/golangci-lint/ holds no schema and no checksum file (the
#     linter embeds the schema of its own version);
#   - .custom-gcl.yml (read by golangci-lint custom, which cannot read
#     another file) names the same vX.Y.Z;
#   - Taskfile.yml and Makefile read the version file, and no workflow, task
#     file or Makefile line carries a linter version of its own or uses
#     golangci/golangci-lint-action (its own configuration check downloads
#     the schema);
#   - .github/workflows/lint.yml has a step 'run: task dev:lint:config', and
#     .tasks/dev.yaml runs this script from it;
#   - the linter binary is exactly vX.Y.Z (the plugin build reports
#     vX.Y.Z-custom-gcl-<hash>);
#   - golangci-lint config verify passes for .golangci.yml.
#
# The file and workflow rules are line patterns, and comment lines are
# skipped. They do not see a step that is present but disabled (if: false),
# a command split over two lines, or an install of @latest. The binary
# version compare catches a second version that reaches the linter binary.
#
# The binary that runs is the plugin build: golangci-lint custom clones the
# git tag of the version and adds the plugin modules. No checksum database
# covers that clone. To move the linter: references/testing-and-lint.md,
# "Moving the golangci-lint version".
#
# Usage: [GOLANGCI_LINT=path] lint-config-check.sh [--root DIR]
#
#   GOLANGCI_LINT  the linter binary (task dev:lint:config passes
#                  bin/golangci-lint); default: golangci-lint on PATH.
#   --root DIR     check the tree at DIR (default: this repository). The
#                  scenario test (lint-config-check-test.sh) uses it.
#
# Exit: 0 ok; 1 a refusal, one line each on stderr; 2 usage.
set -euo pipefail

root=""
while [ "$#" -gt 0 ]; do
  case "$1" in
    --root)
      [ "$#" -ge 2 ] && [ -n "$2" ] || { echo "lint-config: --root needs a directory" >&2; exit 2; }
      root=$2; shift 2 ;;
    *) echo "lint-config: unknown argument: $1" >&2; exit 2 ;;
  esac
done

# The linter, as an absolute path before the directory changes.
linter=${GOLANGCI_LINT:-golangci-lint}
case "$linter" in
  /*) ;;
  */*) linter=$PWD/$linter ;;
  *) linter=$(command -v "$linter" 2>/dev/null || true) ;;
esac

if [ -z "$root" ]; then
  root=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)
fi
cd "$root"

version_file=.golangci-lint-version
schema_dir=.github/golangci-lint
custom_file=.custom-gcl.yml
lint_workflow=.github/workflows/lint.yml
dev_tasks=.tasks/dev.yaml
script=.github/scripts/lint-config-check.sh
check_task=dev:lint:config

problems=0
fail() {
  echo "lint-config: $*" >&2
  problems=$((problems + 1))
}

# The version file.
version=""
minor=""
if [ ! -f "$version_file" ]; then
  fail "$version_file is missing"
else
  version=$(tr -d '[:space:]' <"$version_file")
  if [[ "$version" =~ ^v([0-9]+\.[0-9]+)\.[0-9]+$ ]]; then
    minor=${BASH_REMATCH[1]}
  else
    fail "$version_file holds '$version'; want vX.Y.Z"
  fi
fi

# The linter embeds its schema: no committed schema or checksum file.
if [ -d "$schema_dir" ]; then
  for f in "$schema_dir"/*; do
    [ -e "$f" ] || continue
    fail "$f: the linter embeds the schema of its version; remove the file"
  done
fi

# The plugin build file: golangci-lint custom builds the linter at the
# version this file names.
if [ -f "$custom_file" ] && [ -n "$minor" ]; then
  custom=$(awk '/^version:/ {v = $2; gsub(/["\047]/, "", v); print v; exit}' "$custom_file")
  [ "$custom" = "$version" ] ||
    fail "$custom_file names version '$custom' but $version_file names $version; they must be equal"
fi

# Every place that could name a second version or bring the download back.
# Comment lines are skipped.
for f in Taskfile.yml Makefile; do
  [ -f "$f" ] || continue
  grep -v '^[[:space:]]*#' "$f" | grep -Fq -- "$version_file" ||
    fail "$f does not read $version_file; only that file names the linter version"
done
for f in .github/workflows/*.yml .github/workflows/*.yaml Taskfile.yml .tasks/*.yaml .tasks/*.yml Makefile; do
  [ -f "$f" ] || continue
  while IFS=$'\t' read -r line kind; do
    case "$kind" in
      action) fail "$f:$line: uses golangci/golangci-lint-action; its own configuration check downloads the schema. Run 'task $check_task'" ;;
      literal) fail "$f:$line: carries a linter version of its own; only $version_file names the version" ;;
    esac
  done < <(awk '
    /^[[:space:]]*#/ { next }
    {
      l = tolower($0)
      if (index(l, "golangci/golangci-lint-action")) printf "%d\taction\n", NR
      if (l ~ /golangci-lint[^[:space:]]*@v?[0-9]/ ||
          l ~ /golangci_lint_version[^{$]*[:=][[:space:]]*["\047]?v?[0-9]+\.[0-9]+/) printf "%d\tliteral\n", NR
    }
  ' "$f")
done

# The use of the check: the Lint job runs the task, and the task runs this
# script.
if [ ! -f "$lint_workflow" ]; then
  fail "$lint_workflow is missing"
elif ! grep -Eq "^[[:space:]]*(-[[:space:]]+)?run:[[:space:]]*task[[:space:]]+${check_task}[[:space:]]*$" "$lint_workflow"; then
  fail "$lint_workflow has no step 'run: task $check_task'"
fi
if [ ! -f "$dev_tasks" ]; then
  fail "$dev_tasks is missing"
elif ! grep -v '^[[:space:]]*#' "$dev_tasks" | grep -Fq -- "bash $script"; then
  fail "$dev_tasks does not run 'bash $script'"
fi

# The linter binary, then the configuration itself.
if [ -z "$linter" ] || [ ! -x "$linter" ]; then
  fail "golangci-lint not found (${GOLANGCI_LINT:-PATH}); run 'task $check_task', which installs ${version:-the version in $version_file}"
elif [ -n "$minor" ]; then
  have=$("$linter" version --short 2>/dev/null || true)
  if [ "$have" != "$version" ] && [ "$have" != "${version#v}" ] && [[ "$have" != "$version"-custom-gcl-* ]]; then
    fail "golangci-lint is '$have' but $version_file names $version; remove bin/golangci-lint and run 'task $check_task'"
  else
    HTTPS_PROXY=http://127.0.0.1:9 HTTP_PROXY=http://127.0.0.1:9 \
      https_proxy=http://127.0.0.1:9 http_proxy=http://127.0.0.1:9 NO_PROXY='' no_proxy='' \
      "$linter" config verify ||
      fail ".golangci.yml does not pass golangci-lint config verify"
  fi
fi

if [ "$problems" -gt 0 ]; then
  echo "lint-config: $problems problem(s)" >&2
  exit 1
fi
echo "lint-config: ok ($version, no network)"
