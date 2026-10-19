#!/usr/bin/env bash
# release-guard.sh assert-draft|assert-published|publish TAG
#
#   assert-draft TAG      fails unless exactly one release carries TAG and it is a draft
#   assert-published TAG  read-only: fails unless exactly one release carries TAG and it
#                         is published (the controller module names only published releases)
#   publish TAG           publishes the draft once its required assets are attached:
#                         install.yaml for a controller module tag (opm_controller-v*),
#                         opm-examples.tar.gz for a controller tag (v*): the install
#                         manifest is the module's render, published by the module release.
#                         A controller module release is published with make_latest=false:
#                         it never takes the repository's "latest" mark (enhancement 0021, "The controller's install artifact is a registry module on its own version train").
set -euo pipefail
[ $# -eq 2 ] || { echo "::error::usage: $0 assert-draft|assert-published|publish TAG"; exit 2; }
mode=$1 tag=$2
: "${GH_REPO:?GH_REPO must name the repository}"
rels=$(gh api --paginate "repos/${GH_REPO}/releases?per_page=100" \
  --jq ".[] | select(.tag_name == \"${tag}\") | {id, draft, prerelease, assets: [.assets[].name]}" | jq -s .)
n=$(jq length <<<"$rels")
[ "$n" -eq 1 ] || { echo "::error::expected exactly one release for ${tag}, found ${n}"; exit 1; }
draft=$(jq -r '.[0].draft' <<<"$rels")
case $mode in
assert-draft)
  [ "$draft" = true ] || { echo "::error::release ${tag} is published; release the next version"; exit 1; } ;;
assert-published)
  [ "$draft" = false ] || { echo "::error::release ${tag} is still a draft"; exit 1; } ;;
publish)
  [ "$draft" = true ] || { echo "release ${tag} already published, nothing to do"; exit 0; }
  latest=()
  case $tag in
  opm_controller-v*) required=(install.yaml) latest=(-f make_latest=false) ;;
  v*) required=(opm-examples.tar.gz) ;;
  *) echo "::error::${tag} is neither a controller nor a controller module tag"; exit 1 ;;
  esac
  for a in "${required[@]}"; do
    jq -e --arg a "$a" '.[0].assets | index($a) != null' <<<"$rels" >/dev/null \
      || { echo "::error::draft ${tag} lacks ${a}"; exit 1; }
  done
  id=$(jq -r '.[0].id' <<<"$rels")
  gh api -X PATCH "repos/${GH_REPO}/releases/${id}" -F draft=false "${latest[@]}" >/dev/null ;;
*) echo "::error::unknown mode ${mode}"; exit 2 ;;
esac
