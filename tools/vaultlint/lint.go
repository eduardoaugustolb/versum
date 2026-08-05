// Command vaultlint checks that Obsidian Vault notes follow the metadata and
// linking standard defined in "Rules/03 - Padrão do Vault.md": required
// frontmatter fields and resolvable wikilinks.
package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// requiredKeys are the frontmatter fields every note must declare, per
// "Toda nota adicionada ao vault possui frontmatter com title, section,
// type, status, tags e up."
var requiredKeys = []string{"title", "section", "type", "status", "tags", "up"}

// trailKeys must be present (their value may be the literal null) because
// every note sits in a prev/next trail, even when it has no neighbor yet.
var trailKeys = []string{"prev", "next"}

var (
	frontmatterKeyPattern = regexp.MustCompile(`^([A-Za-z]+):\s*(.*)$`)
	// The alias separator may be escaped as \| when the link sits inside a
	// Markdown table cell, where an unescaped | would split the column.
	wikilinkPattern = regexp.MustCompile(`\[\[([^\]|\\]+)(\\?\|[^\]]*)?\]\]`)
)

// rootIndexPath is the only note allowed to have up: null, since it is the
// entry point of the vault and has nothing above it.
const rootIndexPath = "_Index.md"

// skippedDirs are not linted: .obsidian holds Obsidian config and templates,
// whose frontmatter contains placeholders such as {{title}} on purpose.
var skippedDirs = map[string]bool{".obsidian": true}

// Issue is a single lint failure found in a note.
type Issue struct {
	Path    string
	Message string
}

func (i Issue) String() string {
	return fmt.Sprintf("%s: %s", i.Path, i.Message)
}

// Lint walks root and returns every issue found, sorted by file path.
func Lint(root string) ([]Issue, error) {
	var issues []Issue
	var checked int

	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if skippedDirs[entry.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".md" {
			return nil
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		relPath = filepath.ToSlash(relPath)

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		checked++
		issues = append(issues, lintNote(root, relPath, string(data))...)
		return nil
	})
	if err != nil {
		return nil, err
	}
	if checked == 0 {
		return nil, fmt.Errorf("no notes found under %s", root)
	}

	sort.Slice(issues, func(i, j int) bool {
		if issues[i].Path != issues[j].Path {
			return issues[i].Path < issues[j].Path
		}
		return issues[i].Message < issues[j].Message
	})
	return issues, nil
}

func lintNote(root, relPath, content string) []Issue {
	var issues []Issue

	frontmatter, body, ok := splitFrontmatter(content)
	if !ok {
		return []Issue{{Path: relPath, Message: "missing frontmatter (must start with a `---` block)"}}
	}

	fields := parseFrontmatter(frontmatter)

	upMayBeNull := relPath == rootIndexPath
	for _, key := range requiredKeys {
		value, present := fields[key]
		nullAllowed := key == "up" && upMayBeNull
		if !present || strings.TrimSpace(value) == "" || (value == "null" && !nullAllowed) {
			issues = append(issues, Issue{relPath, fmt.Sprintf("missing required field %q", key)})
		}
	}

	for _, key := range trailKeys {
		if _, present := fields[key]; !present {
			issues = append(issues, Issue{relPath, fmt.Sprintf("missing required field %q (use null if there is no neighbor)", key)})
		}
	}

	if value, present := fields["tags"]; present {
		if !isNonEmptyList(value) {
			issues = append(issues, Issue{relPath, "field \"tags\" must be a non-empty list, e.g. [versum, docs]"})
		}
	}

	issues = append(issues, checkWikilinks(root, relPath, frontmatter)...)
	issues = append(issues, checkWikilinks(root, relPath, body)...)

	return issues
}

// splitFrontmatter separates the leading `---` YAML block from the rest of
// the note. It returns ok=false when the block is missing or unterminated.
func splitFrontmatter(content string) (frontmatter, body string, ok bool) {
	const delimiter = "---"

	if !strings.HasPrefix(content, delimiter+"\n") {
		return "", "", false
	}
	rest := content[len(delimiter)+1:]

	end := strings.Index(rest, "\n"+delimiter)
	if end == -1 {
		return "", "", false
	}
	return rest[:end], rest[end+len(delimiter)+1:], true
}

// parseFrontmatter reads simple `key: value` lines. It does not support
// multi-line YAML values, which the vault's frontmatter never uses.
func parseFrontmatter(frontmatter string) map[string]string {
	fields := make(map[string]string)
	for line := range strings.SplitSeq(frontmatter, "\n") {
		match := frontmatterKeyPattern.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		fields[match[1]] = strings.TrimSpace(match[2])
	}
	return fields
}

func isNonEmptyList(value string) bool {
	if !strings.HasPrefix(value, "[") || !strings.HasSuffix(value, "]") {
		return false
	}
	inner := strings.TrimSpace(value[1 : len(value)-1])
	return inner != ""
}

// checkWikilinks resolves every [[Path]] or [[Path|Label]] reference in text
// against a file at <root>/<Path>.md.
func checkWikilinks(root, relPath, text string) []Issue {
	var issues []Issue
	for _, match := range wikilinkPattern.FindAllStringSubmatch(text, -1) {
		target := strings.TrimSpace(match[1])
		if target == "" {
			issues = append(issues, Issue{relPath, "empty wikilink target [[]]"})
			continue
		}

		targetPath := filepath.Join(root, filepath.FromSlash(target)+".md")
		if _, err := os.Stat(targetPath); err != nil {
			issues = append(issues, Issue{relPath, fmt.Sprintf("broken wikilink [[%s]]", target)})
		}
	}
	return issues
}
