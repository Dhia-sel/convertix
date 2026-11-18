#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
THRESHOLD="${COVERAGE_THRESHOLD:-80}"
go test -race -covermode=atomic -coverprofile=coverage.out ./...
total=$(go tool cover -func=coverage.out | awk '/^total:/ {gsub("%","",$3); print $3}')
echo "Couverture totale : ${total}% (seuil : ${THRESHOLD}%)"
awk -v t="$total" -v m="$THRESHOLD" 'BEGIN { exit !(t+0 >= m+0) }' || { echo "ECHEC : couverture sous le seuil" >&2; exit 1; }