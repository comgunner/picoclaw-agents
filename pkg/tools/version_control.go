// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package tools

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/comgunner/picoclaw/pkg/logger"
)

// VersionControlTool provides semantic version management.
type VersionControlTool struct {
	workspace string
}

// NewVersionControlTool creates a new VersionControlTool instance.
func NewVersionControlTool(workspace string) *VersionControlTool {
	return &VersionControlTool{
		workspace: workspace,
	}
}

// Name returns the tool name.
func (t *VersionControlTool) Name() string {
	return "version_control"
}

// Description returns the tool description.
func (t *VersionControlTool) Description() string {
	return "Semantic version management tool. Compare versions, validate format, bump versions (major/minor/patch), check constraints, and parse version strings. Use action='compare', 'validate', 'bump', 'constraint', or 'parse'."
}

// Parameters returns the JSON schema for tool parameters.
func (t *VersionControlTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"action": map[string]any{
				"type":        "string",
				"description": "Action: 'compare', 'validate', 'bump', 'constraint', or 'parse'",
				"enum":        []string{"compare", "validate", "bump", "constraint", "parse"},
			},
			"version": map[string]any{
				"type":        "string",
				"description": "Version string (e.g., '1.2.3', 'v1.2.3', '1.2.3-beta.1')",
			},
			"other_version": map[string]any{
				"type":        "string",
				"description": "Second version for comparison (required for action='compare')",
			},
			"bump_type": map[string]any{
				"type":        "string",
				"description": "Bump type: 'major', 'minor', or 'patch' (required for action='bump')",
				"enum":        []string{"major", "minor", "patch"},
			},
			"constraint": map[string]any{
				"type":        "string",
				"description": "Version constraint (e.g., '>=1.0.0', '^2.0.0', '~1.2.3') (required for action='constraint')",
			},
		},
		"required": []string{"action"},
	}
}

// Execute runs the version control tool.
func (t *VersionControlTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	action, ok := args["action"].(string)
	if !ok {
		return ErrorResult("action is required and must be one of: compare, validate, bump, constraint, parse")
	}

	switch action {
	case "compare":
		v1, _ := args["version"].(string)
		v2, _ := args["other_version"].(string)
		return t.compare(v1, v2)
	case "validate":
		version, _ := args["version"].(string)
		return t.validate(version)
	case "bump":
		version, _ := args["version"].(string)
		bumpType, _ := args["bump_type"].(string)
		return t.bump(version, bumpType)
	case "constraint":
		version, _ := args["version"].(string)
		constraint, _ := args["constraint"].(string)
		return t.checkConstraint(version, constraint)
	case "parse":
		version, _ := args["version"].(string)
		return t.parse(version)
	default:
		return ErrorResult(fmt.Sprintf("unknown action: %s. Valid options: compare, validate, bump, constraint, parse", action))
	}
}

// compare compares two semantic versions.
func (t *VersionControlTool) compare(v1Str, v2Str string) *ToolResult {
	if v1Str == "" || v2Str == "" {
		return ErrorResult("both 'version' and 'other_version' are required for action='compare'")
	}

	v1, err := semver.NewVersion(normalizeVersion(v1Str))
	if err != nil {
		logger.ErrorCF("tool", "Failed to parse version",
			map[string]any{
				"tool":    "version_control",
				"error":   err.Error(),
				"version": v1Str,
			})
		return ErrorResult(fmt.Sprintf("invalid version '%s': %v", v1Str, err))
	}

	v2, err := semver.NewVersion(normalizeVersion(v2Str))
	if err != nil {
		logger.ErrorCF("tool", "Failed to parse version",
			map[string]any{
				"tool":    "version_control",
				"error":   err.Error(),
				"version": v2Str,
			})
		return ErrorResult(fmt.Sprintf("invalid version '%s': %v", v2Str, err))
	}

	cmp := v1.Compare(v2)
	var result string
	if cmp < 0 {
		result = fmt.Sprintf("%s < %s", v1Str, v2Str)
	} else if cmp > 0 {
		result = fmt.Sprintf("%s > %s", v1Str, v2Str)
	} else {
		result = fmt.Sprintf("%s = %s", v1Str, v2Str)
	}

	resultData := map[string]any{
		"version_1": v1.String(),
		"version_2": v2.String(),
		"comparison": cmp,
		"result":    result,
	}

	_ = resultData // For future structured output
	return SilentResult(result)
}

// validate checks if a version string is valid semantic versioning.
func (t *VersionControlTool) validate(version string) *ToolResult {
	if version == "" {
		return ErrorResult("version is required for action='validate'")
	}

	_, err := semver.NewVersion(normalizeVersion(version))
	if err != nil {
		result := map[string]any{
			"version": version,
			"valid":   false,
			"error":   err.Error(),
		}
		_ = result // For future structured output
		return ErrorResult(fmt.Sprintf("invalid version '%s': %v", version, err))
	}

	result := map[string]any{
		"version": version,
		"valid":   true,
	}
	_ = result // For future structured output
	return SilentResult(fmt.Sprintf("'%s' is a valid semantic version", version))
}

// bump increments a version by major, minor, or patch.
func (t *VersionControlTool) bump(version string, bumpType string) *ToolResult {
	if version == "" {
		return ErrorResult("version is required for action='bump'")
	}
	if bumpType == "" {
		return ErrorResult("bump_type is required for action='bump' (major, minor, or patch)")
	}

	v, err := semver.NewVersion(normalizeVersion(version))
	if err != nil {
		logger.ErrorCF("tool", "Failed to parse version",
			map[string]any{
				"tool":    "version_control",
				"error":   err.Error(),
				"version": version,
			})
		return ErrorResult(fmt.Sprintf("invalid version '%s': %v", version, err))
	}

	var bumped semver.Version
	switch bumpType {
	case "major":
		bumped = v.IncMajor()
	case "minor":
		bumped = v.IncMinor()
	case "patch":
		bumped = v.IncPatch()
	default:
		return ErrorResult(fmt.Sprintf("invalid bump_type '%s': must be major, minor, or patch", bumpType))
	}

	result := map[string]any{
		"original":  v.String(),
		"bump_type": bumpType,
		"bumped":    bumped.String(),
	}

	_ = result // For future structured output
	return SilentResult(fmt.Sprintf("Bumped %s to %s (%s)", v.String(), bumped.String(), bumpType))
}

// checkConstraint checks if a version satisfies a constraint.
func (t *VersionControlTool) checkConstraint(version, constraint string) *ToolResult {
	if version == "" {
		return ErrorResult("version is required for action='constraint'")
	}
	if constraint == "" {
		return ErrorResult("constraint is required for action='constraint'")
	}

	v, err := semver.NewVersion(normalizeVersion(version))
	if err != nil {
		logger.ErrorCF("tool", "Failed to parse version",
			map[string]any{
				"tool":    "version_control",
				"error":   err.Error(),
				"version": version,
			})
		return ErrorResult(fmt.Sprintf("invalid version '%s': %v", version, err))
	}

	c, err := semver.NewConstraint(constraint)
	if err != nil {
		logger.ErrorCF("tool", "Failed to parse constraint",
			map[string]any{
				"tool":       "version_control",
				"error":      err.Error(),
				"constraint": constraint,
			})
		return ErrorResult(fmt.Sprintf("invalid constraint '%s': %v", constraint, err))
	}

	valid := c.Check(v)
	result := map[string]any{
		"version":    v.String(),
		"constraint": constraint,
		"satisfies":  valid,
	}

	_ = result // For future structured output
	if valid {
		return SilentResult(fmt.Sprintf("Version %s satisfies constraint %s", v.String(), constraint))
	}
	return SilentResult(fmt.Sprintf("Version %s does NOT satisfy constraint %s", v.String(), constraint))
}

// parse extracts detailed information from a version string.
func (t *VersionControlTool) parse(version string) *ToolResult {
	if version == "" {
		return ErrorResult("version is required for action='parse'")
	}

	v, err := semver.NewVersion(normalizeVersion(version))
	if err != nil {
		logger.ErrorCF("tool", "Failed to parse version",
			map[string]any{
				"tool":    "version_control",
				"error":   err.Error(),
				"version": version,
			})
		return ErrorResult(fmt.Sprintf("invalid version '%s': %v", version, err))
	}

	result := map[string]any{
		"original":      version,
		"major":         v.Major(),
		"minor":         v.Minor(),
		"patch":         v.Patch(),
		"prerelease":    v.Prerelease(),
		"metadata":      v.Metadata(),
		"normalized":    v.String(),
		"original_repr": v.Original(),
	}

	_ = result // For future structured output
	return SilentResult(fmt.Sprintf("Parsed: v%d.%d.%d (prerelease: %s, metadata: %s)",
		v.Major(), v.Minor(), v.Patch(), v.Prerelease(), v.Metadata()))
}

// normalizeVersion removes leading 'v' or 'V' if present for semver compatibility.
func normalizeVersion(version string) string {
	if len(version) > 0 && (version[0] == 'v' || version[0] == 'V') {
		return version[1:]
	}
	return version
}
