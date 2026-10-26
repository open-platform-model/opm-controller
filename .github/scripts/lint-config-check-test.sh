#!/usr/bin/env bash
# Scenario test of lint-config-check.sh: the pass case on a copy of this
# tree, then one defect per scenario, each of which the check must refuse
# with the named message. This script is offline; it needs the linter binary
# (the Lint job and task dev:lint:config:test pass bin/golangci-lint as
# GOLANGCI_LINT). The task is not offline: it installs the linter first.
#
# Usage: [GOLANGCI_LINT=path] lint-config-check-test.sh
set -euo pipefail

repo=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)
check=$repo/.github/scripts/lint-config-check.sh
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT

real=${GOLANGCI_LINT:-golangci-lint}
case "$real" in
  /*) ;;
  */*) real=$PWD/$real ;;
  *) real=$(command -v "$real" 2>/dev/null || true) ;;
esac
[ -n "$real" ] && [ -x "$real" ] || { echo "lint-config test: no linter binary; run task dev:lint:config:test" >&2; exit 2; }

failed=0
count=0

# stub_linter DIR VERSION: a golangci-lint in DIR that reports VERSION and,
# for config verify, prints the proxy it was given and fails.
stub_linter() {
  mkdir -p "$1"
  cat >"$1/golangci-lint" <<STUB
#!/bin/sh
if [ "\$1" = version ]; then echo "$2"; exit 0; fi
echo "stub: HTTPS_PROXY=\$HTTPS_PROXY HTTP_PROXY=\$HTTP_PROXY NO_PROXY=[\$NO_PROXY]"
exit 1
STUB
  chmod +x "$1/golangci-lint"
}

# fresh NAME: a copy of the files the check reads, in $work/NAME.
fresh() {
  local dir=$work/$1
  mkdir -p "$dir/.github/workflows" "$dir/.tasks"
  cp "$repo/.golangci.yml" "$repo/.golangci-lint-version" "$repo/.custom-gcl.yml" \
    "$repo/Taskfile.yml" "$dir/"
  # The repository has no Makefile. The check still judges one that comes
  # back, so the scenarios run against this stand-in.
  printf '%s\n' "GOLANGCI_LINT_VERSION ?= \$(shell tr -d '[:space:]' <.golangci-lint-version)" >"$dir/Makefile"
  cp "$repo"/.tasks/*.yaml "$dir/.tasks/"
  cp "$repo"/.github/workflows/*.yml "$dir/.github/workflows/"
  echo "$dir"
}

# expect NAME STATUS PATTERN: the check on $work/NAME exits STATUS and its
# output holds PATTERN (a fixed string; empty for none). A non-empty
# $use_linter replaces the linter for that run.
use_linter=""
expect() {
  local name=$1 status=$2 pattern=$3 out rc=0
  count=$((count + 1))
  out=$(GOLANGCI_LINT="${use_linter:-$real}" bash "$check" --root "$work/$name" 2>&1) || rc=$?
  if [ "$rc" -ne "$status" ]; then
    echo "FAIL $name: exit $rc, want $status"; echo "$out"; failed=$((failed + 1)); return
  fi
  if [ -n "$pattern" ] && ! grep -Fq -- "$pattern" <<<"$out"; then
    echo "FAIL $name: output lacks '$pattern'"; echo "$out"; failed=$((failed + 1)); return
  fi
  echo "ok   $name"
}

version=$(tr -d '[:space:]' <"$repo/.golangci-lint-version")

d=$(fresh pass)
expect pass 0 "lint-config: ok"

# The embedded schema is the one of the pinned line: a value that the
# v2.14 line added passes. The v2.13 schema refuses the revive rule name.
d=$(fresh value-of-this-line)
sed -i 's|^\( *\)- name: comment-spacings$|\1- name: comment-spacings\n\1- name: use-slices-concat|' "$d/.golangci.yml"
grep -q -- '- name: use-slices-concat$' "$d/.golangci.yml"
expect value-of-this-line 0 "lint-config: ok"

# A schema file or a checksum file in the tree is refused: the linter
# embeds its schema.
d=$(fresh schema-file-returns)
mkdir -p "$d/.github/golangci-lint"
echo '{}' >"$d/.github/golangci-lint/golangci.v2.14.jsonschema.json"
expect schema-file-returns 1 "golangci.v2.14.jsonschema.json: the linter embeds the schema"

d=$(fresh checksum-file-returns)
mkdir -p "$d/.github/golangci-lint"
echo 'x  y' >"$d/.github/golangci-lint/SHA256SUMS"
expect checksum-file-returns 1 "SHA256SUMS: the linter embeds the schema"

d=$(fresh no-check-step)
sed -i '/^ *run: task dev:lint:config$/d' "$d/.github/workflows/lint.yml"
expect no-check-step 1 "lint.yml has no step 'run: task dev:lint:config'"

d=$(fresh commented-check-step)
sed -i 's|^\( *\)run: task dev:lint:config$|\1# run: task dev:lint:config|' "$d/.github/workflows/lint.yml"
expect commented-check-step 1 "lint.yml has no step 'run: task dev:lint:config'"

d=$(fresh action-returns)
printf '      - uses: GolangCI/golangci-lint-action@0123456789abcdef0123456789abcdef01234567\n' >>"$d/.github/workflows/test.yml"
expect action-returns 1 "test.yml:"
expect action-returns 1 "uses golangci/golangci-lint-action"

d=$(fresh taskfile-literal)
sed -i '/^  GOLANGCI_LINT_VERSION:$/{n;d}' "$d/Taskfile.yml"
sed -i 's|^  GOLANGCI_LINT_VERSION:$|  GOLANGCI_LINT_VERSION: v2.9.0|' "$d/Taskfile.yml"
expect taskfile-literal 1 "Taskfile.yml:"
expect taskfile-literal 1 "carries a linter version of its own"
expect taskfile-literal 1 "Taskfile.yml does not read .golangci-lint-version"

d=$(fresh makefile-literal)
sed -i 's|^GOLANGCI_LINT_VERSION ?= .*$|GOLANGCI_LINT_VERSION ?= v2.9.0|' "$d/Makefile"
expect makefile-literal 1 "Makefile:"
expect makefile-literal 1 "carries a linter version of its own"
expect makefile-literal 1 "Makefile does not read .golangci-lint-version"

d=$(fresh install-literal)
sed -i 's|golangci-lint@{{.GOLANGCI_LINT_VERSION}}|golangci-lint@v2.9.0|' "$d/.tasks/tools.yaml"
grep -q 'golangci-lint@v2.9.0' "$d/.tasks/tools.yaml"
expect install-literal 1 ".tasks/tools.yaml:"
expect install-literal 1 "carries a linter version of its own"

d=$(fresh custom-file-other-version)
sed -i 's|^version: .*$|version: v2.7.2|' "$d/.custom-gcl.yml"
expect custom-file-other-version 1 ".custom-gcl.yml names version 'v2.7.2' but .golangci-lint-version names $version"

d=$(fresh other-minor)
echo v2.99.0 >"$d/.golangci-lint-version"
expect other-minor 1 ".custom-gcl.yml names version '$version' but .golangci-lint-version names v2.99.0"

d=$(fresh bad-version)
echo latest >"$d/.golangci-lint-version"
expect bad-version 1 "want vX.Y.Z"

d=$(fresh no-version-file)
rm "$d/.golangci-lint-version"
expect no-version-file 1 ".golangci-lint-version is missing"

# The linter binary is another version than the file names: another minor,
# another patch, and a plugin build of another patch.
d=$(fresh other-version-installed)
stub_linter "$work/stub-other-minor" "v2.99.0"
use_linter="$work/stub-other-minor/golangci-lint"
expect other-version-installed 1 "golangci-lint is 'v2.99.0' but .golangci-lint-version names $version"
stub_linter "$work/stub-other-patch" "${version%.*}.999-custom-gcl-abc"
use_linter="$work/stub-other-patch/golangci-lint"
expect other-version-installed 1 "golangci-lint is '${version%.*}.999-custom-gcl-abc' but"
use_linter=""

# The linter runs with its proxy variables at a closed local port. A stub
# stands for a linter that would reach for the network.
d=$(fresh proxy-guard)
stub_linter "$work/stub-same" "$version-custom-gcl-abc"
use_linter="$work/stub-same/golangci-lint"
expect proxy-guard 1 "stub: HTTPS_PROXY=http://127.0.0.1:9 HTTP_PROXY=http://127.0.0.1:9 NO_PROXY=[]"
use_linter=""

d=$(fresh no-linter)
use_linter="$work/no-such-dir/golangci-lint"
expect no-linter 1 "golangci-lint not found"
use_linter=""

d=$(fresh invalid-config)
printf '\nnot-a-golangci-key: true\n' >>"$d/.golangci.yml"
expect invalid-config 1 "does not pass golangci-lint config verify"
expect invalid-config 1 "not-a-golangci-key"

count=$((count + 1))
if out=$(bash "$check" --no-such-flag 2>&1); then rc=0; else rc=$?; fi
if [ "$rc" -eq 2 ]; then echo "ok   usage"; else echo "FAIL usage: exit $rc, want 2: $out"; failed=$((failed + 1)); fi

if [ "$failed" -gt 0 ]; then
  echo "lint-config test: $failed of $count failed"
  exit 1
fi
echo "lint-config test: $count passed"
