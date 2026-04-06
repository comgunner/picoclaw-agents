package utils

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// ImportRule defines a forbidden import pattern.
type ImportRule struct {
	From    string // Package that must not import...
	MustNot string // ...this package
}

// DefaultRules are the default forbidden import rules for PicoClaw.
var DefaultRules = []ImportRule{
	{"pkg/agent", "pkg/channels"},
	{"pkg/tools", "pkg/agent"},
	{"pkg/mcp", "pkg/agent"},
	{"pkg/mcp", "pkg/providers"},
}

// Violation describes a single import rule breach.
type Violation struct {
	Rule ImportRule
	File string
}

// CheckImports scans the Go files in root and reports import violations.
func CheckImports(root string, rules []ImportRule) ([]Violation, error) {
	if rules == nil {
		rules = DefaultRules
	}
	fset := token.NewFileSet()
	var violations []Violation

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		node, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return nil
		}
		for _, imp := range node.Imports {
			importPath := strings.Trim(imp.Path.Value, `"`)
			for _, rule := range rules {
				if strings.Contains(importPath, rule.MustNot) {
					rel, _ := filepath.Rel(root, path)
					if strings.HasPrefix(filepath.ToSlash(rel), rule.From) {
						violations = append(violations, Violation{Rule: rule, File: rel})
					}
				}
			}
		}
		return nil
	})
	return violations, err
}

// PrintViolations formats violations for display.
func PrintViolations(violations []Violation) {
	if len(violations) == 0 {
		fmt.Println("✅ No import violations found.")
		return
	}
	fmt.Printf("Found %d violation(s):\n", len(violations))
	for _, v := range violations {
		fmt.Printf("  %s → must not import %s\n    File: %s\n",
			v.Rule.From, v.Rule.MustNot, v.File)
	}
}
