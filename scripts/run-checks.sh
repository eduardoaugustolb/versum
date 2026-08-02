#!/usr/bin/env bash
# Runs the same checks CI runs, scoped to what actually changed, so a commit
# that reaches the remote is already known-good.
set -euo pipefail

repo_root="$(git rev-parse --show-toplevel)"
cd "$repo_root"

changed() {
  git diff --cached --name-only --diff-filter=ACMR | grep -qE "$1"
}

status=0

if changed '^api/' && [ -n "$(find api -name '*.go' -print -quit)" ]; then
  echo "==> api: build, vet, test"
  (cd api && go build ./... && go vet ./... && go test ./...) || status=1
fi

if changed '^bible/'; then
  echo "==> bible/tools: build, vet, test, integrity check"
  (cd bible/tools && go build ./... && go vet ./... && go test ./... && go run ./cmd/bible-integrity -root ../corpus/v1) || status=1
fi

if changed '^tools/vaultlint/' || changed '^Obsidian Vault/'; then
  echo "==> vaultlint: test, lint vault"
  (cd tools/vaultlint && go build ./... && go vet ./... && go test ./... && go run . -root "../../Obsidian Vault") || status=1
fi

exit "$status"
