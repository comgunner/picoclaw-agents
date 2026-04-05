// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package migrate

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

// PrintConfigDiff reads source config, converts it, and prints a field-level diff.
// Used by --show-diff in dry-run mode.
func PrintConfigDiff(sourcePath string) {
	raw, err := LoadOpenClawConfig(sourcePath)
	if err != nil {
		fmt.Printf("  [cannot read source config: %v]\n", err)
		return
	}

	converted, warnings, _ := ConvertConfig(raw)

	// Marshal converted config back to map for unified diff comparison.
	convertedBytes, err := json.Marshal(converted)
	if err != nil {
		fmt.Printf("  [cannot marshal converted config: %v]\n", err)
		return
	}
	var convertedMap map[string]any
	if err := json.Unmarshal(convertedBytes, &convertedMap); err != nil {
		fmt.Printf("  [cannot unmarshal converted config: %v]\n", err)
		return
	}

	fmt.Println()
	fmt.Println("  ── Config diff: source → converted ─────────────────────")
	printFlatDiff(raw, convertedMap)

	if len(warnings) > 0 {
		fmt.Println()
		for _, w := range warnings {
			fmt.Printf("  ⚠  %s\n", w)
		}
	}
	fmt.Println("  ─────────────────────────────────────────────────────────")
}

// PrintNanoClawConfigDiff reads a nanoclaw config and shows what would be generated.
func PrintNanoClawConfigDiff(sourcePath string) {
	data, err := os.ReadFile(sourcePath)
	if err != nil {
		fmt.Printf("  [cannot read nanoclaw config: %v]\n", err)
		return
	}
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		fmt.Printf("  [invalid JSON in nanoclaw config: %v]\n", err)
		return
	}

	converted, warnings, _ := ConvertNanoClawConfig(raw)

	fmt.Println()
	fmt.Println("  ── NanoClaw config diff: source → converted ─────────────")
	printFlatDiff(raw, converted)

	if len(warnings) > 0 {
		fmt.Println()
		for _, w := range warnings {
			fmt.Printf("  ⚠  %s\n", w)
		}
	}
	fmt.Println("  ─────────────────────────────────────────────────────────")
}

// printFlatDiff flattens both maps to dot-notation keys and shows removed/added/changed lines.
func printFlatDiff(before, after map[string]any) {
	flatBefore := flattenMap("", before)
	flatAfter := flattenMap("", after)

	// Collect all keys
	allKeys := map[string]bool{}
	for k := range flatBefore {
		allKeys[k] = true
	}
	for k := range flatAfter {
		allKeys[k] = true
	}

	keys := make([]string, 0, len(allKeys))
	for k := range allKeys {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	changed := 0
	for _, k := range keys {
		bv, bExists := flatBefore[k]
		av, aExists := flatAfter[k]

		bStr := jsonVal(bv)
		aStr := jsonVal(av)

		switch {
		case bExists && !aExists:
			fmt.Printf("  \033[31m- %-40s  %s\033[0m\n", k+":", bStr)
			changed++
		case !bExists && aExists:
			fmt.Printf("  \033[32m+ %-40s  %s\033[0m\n", k+":", aStr)
			changed++
		case bStr != aStr:
			fmt.Printf("  \033[31m- %-40s  %s\033[0m\n", k+":", bStr)
			fmt.Printf("  \033[32m+ %-40s  %s\033[0m\n", k+":", aStr)
			changed++
		}
	}

	if changed == 0 {
		fmt.Println("  (no differences — config is already compatible)")
	} else {
		fmt.Printf("\n  %d field(s) would change\n", changed)
	}
}

// flattenMap converts a nested map to dot-notation keys, masking secrets.
func flattenMap(prefix string, m map[string]any) map[string]string {
	result := map[string]string{}
	for k, v := range m {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		switch val := v.(type) {
		case map[string]any:
			for fk, fv := range flattenMap(key, val) {
				result[fk] = fv
			}
		case []any:
			result[key] = fmt.Sprintf("[%d items]", len(val))
		default:
			result[key] = maskSecret(key, jsonVal(v))
		}
	}
	return result
}

// maskSecret hides values for keys that look like secrets.
func maskSecret(key, val string) string {
	lower := strings.ToLower(key)
	if strings.Contains(lower, "api_key") || strings.Contains(lower, "token") ||
		strings.Contains(lower, "secret") || strings.Contains(lower, "password") {
		if len(val) > 6 {
			return val[:3] + "***" + val[len(val)-3:]
		}
		return "***"
	}
	return val
}

func jsonVal(v any) string {
	if v == nil {
		return "null"
	}
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	return string(b)
}
