# Agents - Read Before Working

**IMPORTANT:** Before making any changes to this project, you MUST read the documentation in `.ai/docs/`.

Run the following to understand project conventions:
- `.ai/docs/naming-convention.md` - Database, API, and code naming patterns
- `.ai/docs/api-development.md` - How to create API endpoints

See also:
- `.ai/prd.md` - Product Requirements Document
- `.ai/rules/` - Development rules and best practices

---

## Context for Each App

When working on a specific app, you **MUST** read the context from that app's `.ai` folder before making any changes. Each app may have its own rules, conventions, and documentation.

**Always do this:**
1. Find the app you'll be editing (e.g., `apps/client`, `apps/api`, `apps/landing-page`)
2. Read all files in that app's `.ai/` folder - this includes rules, docs, and any AGENT.md
3. Use that context to understand the app-specific conventions

Example: If you're editing `apps/client`, read:
- `apps/client/.ai/rules/*.md`
- `apps/client/.ai/docs/*.md`
- `apps/client/AGENT.md`

The app-specific context may contain rules that differ from or extend the global rules. When in doubt, follow the app-specific rules.
