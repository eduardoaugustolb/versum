#!/usr/bin/env bash
# Blocks commits that credit an AI coding agent as a co-author or include
# AI-generated attribution footers. Human co-authors are always allowed.
#
# Usage: check-ai-coauthor.sh <commit-message-file>
set -euo pipefail

# Name/email fragments identifying known AI coding agents and assistants.
# Keep case-insensitive; matched against "Co-authored-by:" trailers.
AGENT_PATTERN='claude|anthropic|codex|openai|chatgpt|gpt-[0-9]|copilot|gemini[^a-z]|google ai|cursor ai|devin ai|windsurf|sourcegraph cody|codeium|tabnine|amazon ?q|junie|jetbrains ai|replit ai|cody ai|ai agent|coding agent'

# Generic AI-generated attribution footers some tools add outside trailers.
FOOTER_PATTERN='generated with .*(claude|codex|copilot|chatgpt|gpt|gemini|ai)|🤖 generated'

message_file="${1:?usage: check-ai-coauthor.sh <commit-message-file>}"
text="$(cat "$message_file")"

bad_coauthor=$(printf '%s\n' "$text" | grep -iE '^co-authored-by:' | grep -iE "$AGENT_PATTERN" || true)
if [ -n "$bad_coauthor" ]; then
  echo "error: commit credits an AI agent as co-author:" >&2
  printf '%s\n' "$bad_coauthor" | sed 's/^/  /' >&2
  echo "Human co-authors are welcome; AI agents are not." >&2
  exit 1
fi

bad_footer=$(printf '%s\n' "$text" | grep -iE "$FOOTER_PATTERN" || true)
if [ -n "$bad_footer" ]; then
  echo "error: commit message has an AI-generated attribution footer:" >&2
  printf '%s\n' "$bad_footer" | sed 's/^/  /' >&2
  exit 1
fi
