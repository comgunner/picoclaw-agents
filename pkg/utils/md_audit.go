package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// LinkIssue describes a broken or suspicious link in a Markdown file.
type LinkIssue struct {
	File string
	Line int
	Link string
	Type string
}

var mdLinkRe = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)

// AuditMarkdown scans a directory for broken internal links.
func AuditMarkdown(root string) ([]LinkIssue, error) {
	var issues []LinkIssue
	mdFiles := make(map[string]bool)

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err == nil && strings.HasSuffix(path, ".md") {
			rel, _ := filepath.Rel(root, path)
			mdFiles[rel] = true
		}
		return nil
	})

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || !strings.HasSuffix(path, ".md") {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		dir := filepath.Dir(rel)
		lines := strings.Split(string(content), "\n")

		for i, line := range lines {
			for _, m := range mdLinkRe.FindAllStringSubmatch(line, -1) {
				link := m[2]
				if strings.HasPrefix(link, "http") || strings.HasPrefix(link, "#") {
					continue
				}
				resolved := filepath.Join(dir, link)
				resolved = filepath.Clean(resolved)
				if !mdFiles[resolved] && !fileExists(filepath.Join(root, resolved)) {
					issues = append(issues, LinkIssue{File: rel, Line: i + 1, Link: link, Type: "internal"})
				}
			}
		}
		return nil
	})
	return issues, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// PrintLinkIssues formats link issues for display.
func PrintLinkIssues(issues []LinkIssue) {
	if len(issues) == 0 {
		fmt.Println("✅ No broken internal links found.")
		return
	}
	fmt.Printf("Found %d link issue(s):\n", len(issues))
	for _, iss := range issues {
		fmt.Printf("  %s:%d — %s\n", iss.File, iss.Line, iss.Link)
	}
}
