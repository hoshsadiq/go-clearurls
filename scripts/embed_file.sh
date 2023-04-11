#!/usr/bin/env bash

set -euo pipefail

embed() {
  local target in_embed file syntax line

  target="$1"

  in_embed=false
  while IFS="" read -r line || [[ -n "$line" ]]; do
    if [[ "$line" =~ \<!--\ BEGIN\ EMBED\ FILE:\ ([^\ :;]+)(\;([a-z0-9_]+))?\ --\> ]]; then
      file="${BASH_REMATCH[1]}"
      syntax="${BASH_REMATCH[3]}"

      printf "%s\n" "$line"
      printf '```%s\n' "$syntax"
      cat "$file"
      printf '```\n'
      in_embed=true
    elif [[ "$line" =~ \<!--\ END\ EMBED\ FILE\ --\> ]]; then
      printf "%s\n" "$line"
      in_embed=false
    elif [[ $in_embed == false ]]; then
      printf "%s\n" "$line"
    fi
  done <"$target" >"$target.tmp"
  mv "$target.tmp" "$target"
}

main() {
  local target

  for target in "$@"; do
    embed "$target"
  done
}

main "$@"
