#!/bin/sh

#/ input [--wait] dayN
#/
#/ when --wait is set, input waits until the next hourly boundary before downloading

usage() {
  grep '^#' "$0" | cut -c4- >&2
  exit 1
}

WAIT=""
while [ "$#" -gt 1 ]; do
  case "$1" in
    --wait) WAIT="$1"; shift ;;
    *) usage ;;
  esac
done

if [ "$#" -ne 1 ]; then
  usage
fi

N=${1#day}

[ -f "aoc.session" ] || { echo "missing aoc session ID; cannot continue. put bare session cookie value in aoc.session" >&2; exit 1; }

if [ "$WAIT" ]; then
  while [ $(date +%M) -ne 0 ]; do
    printf '\r\e[2K%s' "$(date)"
    sleep 1
  done
fi

until curl -L -f \
        "https://adventofcode.com/2022/day/$N/input" \
        -H "Cookie: session=$(cat aoc.session)" \
        > "day$(printf "%02d" "$N")/a/input"; do
  sleep 10
done
cp "day$(printf "%02d" "$N")/a/input" \
   "day$(printf "%02d" "$N")/b/input"
