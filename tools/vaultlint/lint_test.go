package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeNote(t *testing.T, root, relPath, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(relPath))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestLintValidVault(t *testing.T) {
	root := t.TempDir()
	writeNote(t, root, "_Index.md", `---
title: "Home"
section: Home
type: index
status: active
tags: [versum]
up: null
prev: null
next: "[[Guide]]"
---

# Home
`)
	writeNote(t, root, "Guide.md", `---
title: "Guide"
section: Docs
type: guide
status: active
tags: [versum, docs]
up: "[[_Index|Home]]"
prev: "[[_Index]]"
next: null
related: []
---

# Guide

See [[_Index|Home]] for the map.
`)

	issues, err := Lint(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}

func TestLintMissingFrontmatter(t *testing.T) {
	root := t.TempDir()
	writeNote(t, root, "Note.md", "# No frontmatter\n")

	issues, err := Lint(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(issues) != 1 || issues[0].Message != "missing frontmatter (must start with a `---` block)" {
		t.Fatalf("unexpected issues: %v", issues)
	}
}

func TestLintMissingRequiredFields(t *testing.T) {
	root := t.TempDir()
	writeNote(t, root, "Note.md", `---
title: "Note"
section: Docs
---

# Note
`)

	issues, err := Lint(root)
	if err != nil {
		t.Fatal(err)
	}

	want := map[string]bool{
		`Note.md: missing required field "type"`:   true,
		`Note.md: missing required field "status"`: true,
		`Note.md: missing required field "tags"`:   true,
		`Note.md: missing required field "up"`:     true,
		`Note.md: missing required field "prev" (use null if there is no neighbor)`: true,
		`Note.md: missing required field "next" (use null if there is no neighbor)`: true,
	}
	if len(issues) != len(want) {
		t.Fatalf("expected %d issues, got %d: %v", len(want), len(issues), issues)
	}
	for _, issue := range issues {
		if !want[issue.String()] {
			t.Errorf("unexpected issue: %s", issue)
		}
	}
}

func TestLintUpNullOutsideRoot(t *testing.T) {
	root := t.TempDir()
	writeNote(t, root, "Note.md", `---
title: "Note"
section: Docs
type: guide
status: active
tags: [versum]
up: null
prev: null
next: null
---

# Note
`)

	issues, err := Lint(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(issues) != 1 || issues[0].Message != `missing required field "up"` {
		t.Fatalf("unexpected issues: %v", issues)
	}
}

func TestLintEmptyTags(t *testing.T) {
	root := t.TempDir()
	writeNote(t, root, "Note.md", `---
title: "Note"
section: Docs
type: guide
status: active
tags: []
up: "[[_Index]]"
prev: null
next: null
---

# Note
`)
	writeNote(t, root, "_Index.md", `---
title: "Home"
section: Home
type: index
status: active
tags: [versum]
up: null
prev: null
next: null
---
`)

	issues, err := Lint(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(issues) != 1 || issues[0].Message != "field \"tags\" must be a non-empty list, e.g. [versum, docs]" {
		t.Fatalf("unexpected issues: %v", issues)
	}
}

func TestLintBrokenWikilink(t *testing.T) {
	root := t.TempDir()
	writeNote(t, root, "_Index.md", `---
title: "Home"
section: Home
type: index
status: active
tags: [versum]
up: null
prev: null
next: null
---

See [[Missing Note]].
`)

	issues, err := Lint(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(issues) != 1 || issues[0].Message != "broken wikilink [[Missing Note]]" {
		t.Fatalf("unexpected issues: %v", issues)
	}
}

func TestLintSkipsObsidianDir(t *testing.T) {
	root := t.TempDir()
	writeNote(t, root, "_Index.md", `---
title: "Home"
section: Home
type: index
status: active
tags: [versum]
up: null
prev: null
next: null
---
`)
	writeNote(t, root, ".obsidian/templates/ADR.md", `---
title: "{{title}}"
section: Docs
type: adr
status: proposed
up: "[[Docs/Decisions/_Index|Decisões]]"
prev: null
next: null
related: []
---
`)

	issues, err := Lint(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}

func TestLintNoNotesFound(t *testing.T) {
	root := t.TempDir()
	if _, err := Lint(root); err == nil {
		t.Fatal("expected an error for an empty vault")
	}
}
