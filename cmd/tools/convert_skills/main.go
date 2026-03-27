// Tool to convert skills from local_work/skills_import/ to embedded format
// Usage: go run cmd/tools/convert_skills/main.go

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Skills to skip (already implemented as native Go structs)
var skipFiles = map[string]bool{
	"backend_developer.md":  true,
	"frontend_developer.md": true,
	"devops_engineer.md":    true,
	"security_engineer.md":  true,
	"qa_engineer.md":        true,
	"data_engineer.md":      true,
	"ml_engineer.md":        true,
}

// Categories to skip
var skipCategories = map[string]bool{
	"examples": true,
	"strategy": true,
}

func main() {
	sourceDir := "local_work/skills_import"
	targetDir := "pkg/skills/data"

	// Clean target directory
	fmt.Println("Cleaning target directory...")
	if err := os.RemoveAll(targetDir); err != nil {
		fmt.Printf("Error cleaning target: %v\n", err)
		os.Exit(1)
	}

	// Read source directory
	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		fmt.Printf("Error reading source directory: %v\n", err)
		os.Exit(1)
	}

	totalConverted := 0
	totalSkipped := 0

	for _, categoryEntry := range entries {
		if !categoryEntry.IsDir() {
			continue
		}

		category := categoryEntry.Name()

		// Skip excluded categories
		if skipCategories[category] {
			fmt.Printf("Skipping category: %s\n", category)
			continue
		}

		categoryPath := filepath.Join(sourceDir, category)
		skillFiles, err := os.ReadDir(categoryPath)
		if err != nil {
			fmt.Printf("Error reading category %s: %v\n", category, err)
			continue
		}

		for _, skillFile := range skillFiles {
			if skillFile.IsDir() {
				continue
			}

			if !strings.HasSuffix(skillFile.Name(), ".md") {
				continue
			}

			// Skip if it's a native Go skill
			if skipFiles[skillFile.Name()] {
				fmt.Printf("Skipping (native Go): %s/%s\n", category, skillFile.Name())
				totalSkipped++
				continue
			}

			// Convert skill
			skillName := strings.TrimSuffix(skillFile.Name(), ".md")
			sourcePath := filepath.Join(categoryPath, skillFile.Name())
			targetPath := filepath.Join(targetDir, category, skillName, "SKILL.md")

			if err := convertSkill(sourcePath, targetPath, category, skillName); err != nil {
				fmt.Printf("Error converting %s/%s: %v\n", category, skillFile.Name(), err)
				continue
			}

			totalConverted++
			fmt.Printf("Converted: %s/%s\n", category, skillName)
		}
	}

	fmt.Printf("\n=== Conversion Complete ===\n")
	fmt.Printf("Total converted: %d\n", totalConverted)
	fmt.Printf("Total skipped (native Go): %d\n", totalSkipped)
	fmt.Printf("Target directory: %s\n", targetDir)
}

func convertSkill(sourcePath, targetPath, category, skillName string) error {
	// Read source file
	content, err := os.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("reading source: %w", err)
	}

	contentStr := string(content)

	// Extract content after the first sections (skip metadata header)
	body := extractBody(contentStr)

	// Extract description from content
	description := extractDescription(body)

	// Create frontmatter
	// Use filename-based name (with hyphens) for compatibility
	frontmatter := fmt.Sprintf(`---
name: %s
description: %s
category: %s
version: 1.0.0
---

`, skillName, description, category)

	// Create target directory
	targetDir := filepath.Dir(targetPath)
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	// Write target file
	finalContent := frontmatter + body
	if err := os.WriteFile(targetPath, []byte(finalContent), 0o644); err != nil {
		return fmt.Errorf("writing target: %w", err)
	}

	return nil
}

// extractBody extracts the main content after metadata headers
func extractBody(content string) string {
	lines := strings.Split(content, "\n")
	startIndex := 0
	foundFirstHeader := false
	foundSecondHeader := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Look for the first # header (skill title)
		if strings.HasPrefix(trimmed, "# ") && !foundFirstHeader {
			foundFirstHeader = true
			continue
		}

		// Skip metadata lines after first header
		if foundFirstHeader && !foundSecondHeader {
			if strings.HasPrefix(trimmed, "**Skill:**") ||
				strings.HasPrefix(trimmed, "**Categoría:**") ||
				strings.HasPrefix(trimmed, "**Versión:**") ||
				strings.HasPrefix(trimmed, "**Original:**") ||
				strings.HasPrefix(trimmed, "**Issues encontrados:**") ||
				strings.HasPrefix(trimmed, "- Paths absolutos") ||
				strings.HasPrefix(trimmed, "- ") ||
				trimmed == "" ||
				strings.HasPrefix(trimmed, "---") {
				continue
			}

			// Found start of actual content
			if strings.HasPrefix(trimmed, "# ") {
				// Another header - this is the actual skill title
				startIndex = i
				foundSecondHeader = true
				break
			}
		}
	}

	if foundSecondHeader {
		return strings.Join(lines[startIndex:], "\n")
	}

	// Fallback: return content after first ---
	parts := strings.SplitN(content, "---", 3)
	if len(parts) >= 3 {
		return strings.TrimSpace(parts[2])
	}

	return content
}

// extractDescription extracts a one-line description from content
func extractDescription(content string) string {
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip empty lines
		if trimmed == "" {
			continue
		}

		// Skip headers
		if strings.HasPrefix(trimmed, "#") {
			continue
		}

		// Skip bullet points
		if strings.HasPrefix(trimmed, "-") || strings.HasPrefix(trimmed, "*") {
			continue
		}

		// Skip lines that are too short or too long
		if len(trimmed) < 20 || len(trimmed) > 120 {
			continue
		}

		// Skip lines with code blocks
		if strings.Contains(trimmed, "```") {
			continue
		}

		// Found a good description line
		// Remove markdown formatting
		cleaned := cleanMarkdown(trimmed)
		if len(cleaned) >= 20 && len(cleaned) <= 120 {
			return cleaned
		}
	}

	// Fallback: generate from skill name
	return "Specialized AI skill for professional tasks"
}

// cleanMarkdown removes basic markdown formatting
func cleanMarkdown(text string) string {
	// Remove bold
	text = strings.ReplaceAll(text, "**", "")
	text = strings.ReplaceAll(text, "__", "")

	// Remove italic
	text = strings.ReplaceAll(text, "*", "")
	text = strings.ReplaceAll(text, "_", "")

	// Remove links [text](url) -> text
	linkRe := regexp.MustCompile(`\[([^\]]+)\]\([^)]+\)`)
	text = linkRe.ReplaceAllString(text, "$1")

	// Remove inline code
	text = strings.ReplaceAll(text, "`", "")

	return strings.TrimSpace(text)
}
