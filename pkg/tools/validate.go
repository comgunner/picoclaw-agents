// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package tools

import (
	"fmt"
	"strings"
)

// ValidateToolArgs validates tool arguments against the tool's parameter schema.
// This prevents type confusion, argument injection, and missing-field errors.
// Returns nil if validation passes, or an error describing the validation failure.
func ValidateToolArgs(schema map[string]any, args map[string]any) error {
	if schema == nil {
		return nil // No schema to validate against
	}

	properties, _ := schema["properties"].(map[string]any)
	required, _ := schema["required"].([]any)

	// Check required fields
	if err := checkRequired(required, args); err != nil {
		return err
	}

	// Validate each provided argument
	for key, value := range args {
		propSchema, ok := properties[key].(map[string]any)
		if !ok {
			// Property not in schema - check if additionalProperties is allowed
			if !allowsAdditional(schema) {
				return fmt.Errorf("unexpected argument %q not in schema", key)
			}
			continue
		}

		if err := checkType(key, value, propSchema); err != nil {
			return err
		}
	}

	return nil
}

// checkRequired validates that all required fields are present
func checkRequired(required []any, args map[string]any) error {
	for _, req := range required {
		reqKey, ok := req.(string)
		if !ok {
			continue
		}
		if _, exists := args[reqKey]; !exists {
			return fmt.Errorf("missing required argument %q", reqKey)
		}
	}
	return nil
}

// allowsAdditional checks if the schema allows additional properties
func allowsAdditional(schema map[string]any) bool {
	if addProps, ok := schema["additionalProperties"]; ok {
		if isBool, ok := addProps.(bool); ok {
			return isBool
		}
		if isMap, ok := addProps.(map[string]any); ok {
			// If it's a schema object, treat as allowed
			return len(isMap) > 0
		}
	}
	return false // Default: deny additional properties
}

// checkType validates the type of a value against its schema
func checkType(key string, value any, schema map[string]any) error {
	typeStr, ok := schema["type"].(string)
	if !ok {
		return nil // No type specified
	}

	// Check enum first (applies to any type)
	if enum, ok := schema["enum"].([]any); ok {
		if !checkEnum(value, enum) {
			return fmt.Errorf(
				"argument %q has invalid value %v (must be one of %v)",
				key,
				value,
				enumValuesToString(enum),
			)
		}
	}

	switch typeStr {
	case "string":
		if _, ok := value.(string); !ok {
			return fmt.Errorf("argument %q must be a string, got %T", key, value)
		}

	case "integer":
		if _, ok := value.(float64); !ok {
			return fmt.Errorf("argument %q must be an integer, got %T", key, value)
		}
		// Check min/max if specified
		if minVal, ok := schema["minimum"].(float64); ok {
			if value.(float64) < minVal {
				return fmt.Errorf("argument %q must be >= %v, got %v", key, minVal, value)
			}
		}
		if maxVal, ok := schema["maximum"].(float64); ok {
			if value.(float64) > maxVal {
				return fmt.Errorf("argument %q must be <= %v, got %v", key, maxVal, value)
			}
		}

	case "number":
		if _, ok := value.(float64); !ok {
			return fmt.Errorf("argument %q must be a number, got %T", key, value)
		}

	case "boolean":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("argument %q must be a boolean, got %T", key, value)
		}

	case "array":
		arr, ok := value.([]any)
		if !ok {
			return fmt.Errorf("argument %q must be an array, got %T", key, value)
		}
		// Validate array items if schema specified
		if itemSchema, ok := schema["items"].(map[string]any); ok {
			for i, item := range arr {
				if err := checkType(fmt.Sprintf("%s[%d]", key, i), item, itemSchema); err != nil {
					return err
				}
			}
		}

	case "object":
		obj, ok := value.(map[string]any)
		if !ok {
			return fmt.Errorf("argument %q must be an object, got %T", key, value)
		}
		// Recursively validate nested object
		return ValidateToolArgs(schema, obj)

	case "null":
		if value != nil {
			return fmt.Errorf("argument %q must be null, got %T", key, value)
		}

	default:
		// Unknown type, skip validation
	}

	return nil
}

// checkEnum validates that a value is one of the allowed enum values
func checkEnum(value any, enum []any) bool {
	for _, e := range enum {
		if e == value {
			return true
		}
		// Handle numeric comparisons
		if ev, ok := e.(float64); ok {
			if fv, ok := value.(float64); ok && ev == fv {
				return true
			}
		}
	}
	return false
}

// enumValuesToString converts enum values to a readable string
func enumValuesToString(enum []any) string {
	values := make([]string, 0, len(enum))
	for _, e := range enum {
		values = append(values, fmt.Sprintf("%v", e))
	}
	return strings.Join(values, ", ")
}
