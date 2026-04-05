// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)

package security

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ScanResult representa un secret encontrado durante el escaneo.
type ScanResult struct {
	File    string // ruta absoluta del archivo
	Line    int    // número de línea (1-indexed)
	Pattern string // nombre del patrón que coincidió (ej: "OpenAI")
	Match   string // valor redactado: primeros 4 chars + "****"
}

// String devuelve una representación legible del resultado.
func (r ScanResult) String() string {
	return fmt.Sprintf("%s:%d [%s] %s", r.File, r.Line, r.Pattern, r.Match)
}

// patternEntry asocia un nombre legible con su expresión regular.
type patternEntry struct {
	Name    string
	Pattern *regexp.Regexp
}

// Scanner escanea archivos buscando secrets hardcodeados.
// Reutiliza los patrons del redactor para consistencia.
type Scanner struct {
	patterns []patternEntry
}

// NewScanner crea un Scanner con los patrons estándar de PicoClaw.
// Los patrons son los mismos que usa redactor.go para garantizar consistencia.
func NewScanner() *Scanner {
	rawPatterns := []struct {
		name    string
		pattern string
	}{
		{"Anthropic", `sk-ant-[a-zA-Z0-9\-]{32,}`},
		{"Stripe Secret", `sk_live_[0-9a-zA-Z]{24}`},
		{"DeepSeek", `sk-[a-f0-9]{32,}`},
		{"OpenAI", `sk-[a-zA-Z0-9\-]{40,}`},
		{"Google API", `AIza[0-9A-Za-z\-_]{35}`},
		{"GitHub Token", `ghp_[a-zA-Z0-9]{36}`},
		{"GitHub OAuth", `gho_[a-zA-Z0-9]{36}`},
		{"AWS Access Key", `AKIA[0-9A-Z]{16}`},
		{"AWS Secret", `[0-9a-zA-Z/+]{40}`},
		{"Slack Token", `xox[baprs]-[0-9a-zA-Z\-]{10,}`},
		{"Telegram Bot", `[0-9]{8,10}:[a-zA-Z0-9_\-]{35}`},
		{"JWT", `eyJ[a-zA-Z0-9_\-]+\.eyJ[a-zA-Z0-9_\-]+\.[a-zA-Z0-9_\-]+`},
	}

	patterns := make([]patternEntry, 0, len(rawPatterns))
	for _, rp := range rawPatterns {
		compiled, err := regexp.Compile(rp.pattern)
		if err != nil {
			continue
		}
		patterns = append(patterns, patternEntry{Name: rp.name, Pattern: compiled})
	}

	return &Scanner{patterns: patterns}
}

// redactMatch devuelve los primeros 4 caracteres del match seguidos de "****".
// Si el match tiene menos de 4 caracteres, devuelve "****" completo.
func redactMatch(match string) string {
	if len(match) <= 4 {
		return "****"
	}
	return match[:4] + "****"
}

// ScanFile escanea un archivo línea por línea buscando secrets.
// Devuelve un slice de ScanResult (vacío si no hay coincidencias).
// Retorna error si el archivo no puede abrirse.
func (s *Scanner) ScanFile(path string) ([]ScanResult, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("scanner: open %q: %w", path, err)
	}
	defer f.Close()

	var results []ScanResult
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		for _, pe := range s.patterns {
			match := pe.Pattern.FindString(line)
			if match == "" {
				continue
			}
			// Ignorar placeholders comunes en archivos de ejemplo
			if isPlaceholder(match) {
				continue
			}
			results = append(results, ScanResult{
				File:    path,
				Line:    lineNum,
				Pattern: pe.Name,
				Match:   redactMatch(match),
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return results, fmt.Errorf("scanner: read %q: %w", path, err)
	}

	return results, nil
}

// ScanDir escanea recursivamente un directorio buscando secrets.
// Omite directors ocultos (.git, .picoclaw, etc.) y binarios.
// Devuelve todos los resultados encontrados en el árbol.
func (s *Scanner) ScanDir(dir string) ([]ScanResult, error) {
	var allResults []ScanResult

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err // continuar ante errores de permisos individuals
		}
		if d.IsDir() {
			name := d.Name()
			// Omitir directors que no tienen código fuente relevante
			if strings.HasPrefix(name, ".") || name == "vendor" || name == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		// Solo escanear archivos de texto conocidos
		if !isTextFile(path) {
			return nil
		}

		results, err := s.ScanFile(path)
		if err != nil {
			return err // no abortar por un archivo inaccesible
		}
		allResults = append(allResults, results...)
		return nil
	})

	return allResults, err
}

// isPlaceholder detecta si un match es un valor de ejemplo/placeholder que no es un secret real.
func isPlaceholder(s string) bool {
	lower := strings.ToLower(s)
	placeholders := []string{
		"your_key", "your-key", "your_token", "placeholder",
		"xxxx", "example", "changeme", "sk-test", "sk_test",
		"1234567890", "abcdefgh",
	}
	for _, p := range placeholders {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}

// isTextFile determina si un archivo es probable que sea texto plano escaneable.
func isTextFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	textExtensions := map[string]bool{
		".go": true, ".py": true, ".ts": true, ".js": true,
		".json": true, ".yaml": true, ".yml": true, ".toml": true,
		".env": true, ".sh": true, ".md": true, ".txt": true,
		".conf": true, ".config": true, ".ini": true,
	}
	// Archivos sin extensión también son candidatos (Makefile, Dockerfile, etc.)
	if ext == "" {
		return true
	}
	return textExtensions[ext]
}
