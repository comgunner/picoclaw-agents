// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package util

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newCodegenCommand() *cobra.Command {
	var codeType string
	var pkgName string
	var outputDir string
	var name string

	cmd := &cobra.Command{
		Use:   "codegen",
		Short: "Generate Go boilerplate code",
		Long:  "Generate Go boilerplate code from templates (api, service, handler, model, config).",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runCodegen(codeType, name, pkgName, outputDir)
		},
	}

	cmd.Flags().StringVarP(&codeType, "type", "t", "", "Type of code to generate: api, service, handler, model, config")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Name of the component (e.g., UserService, OrderHandler)")
	cmd.Flags().StringVarP(&pkgName, "package", "p", "", "Package name (default: derived from name)")
	cmd.Flags().StringVarP(&outputDir, "output", "o", "./generated", "Output directory")
	cmd.MarkFlagRequired("type")
	cmd.MarkFlagRequired("name")

	return cmd
}

func runCodegen(codeType, name, pkgName, outputDir string) error {
	validTypes := []string{"api", "service", "handler", "model", "config"}
	isValid := false
	for _, t := range validTypes {
		if codeType == t {
			isValid = true
			break
		}
	}

	if !isValid {
		return fmt.Errorf("invalid type '%s'. Valid types: %v", codeType, validTypes)
	}

	fmt.Printf("Code generation requested:\n")
	fmt.Printf("  Type: %s\n", codeType)
	fmt.Printf("  Name: %s\n", name)
	fmt.Printf("  Package: %s\n", pkgName)
	fmt.Printf("  Output: %s\n", outputDir)
	fmt.Println("\nNote: For full code generation with templates, use the agent:")
	fmt.Println("  picoclaw agent -m \"Generate a UserService for order management\"")
	fmt.Println("\nThe agent will use the codegen tool to create boilerplate code.")

	return nil
}
