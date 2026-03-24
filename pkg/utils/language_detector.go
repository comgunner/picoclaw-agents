// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package utils

import (
	"strings"
	"unicode"
)

// DetectLanguage detecta automáticamente el idioma de un texto
// Retorna "es" para español, "en" para inglés (default)
func DetectLanguage(text string) string {
	words := strings.Fields(strings.ToLower(text))

	if len(words) == 0 {
		return "en" // Default a inglés
	}

	// Palabras comunes en español (high frequency)
	spanishWords := map[string]bool{
		"el": true, "la": true, "los": true, "las": true,
		"de": true, "que": true, "en": true, "es": true,
		"un": true, "una": true, "unos": true, "unas": true,
		"del": true, "al": true, "por": true, "para": true,
		"con": true, "sin": true, "sobre": true, "entre": true,
		"este": true, "esta": true, "estos": true, "estas": true,
		"ser": true, "estar": true, "tener": true, "hacer": true,
		"como": true, "cuando": true, "donde": true, "porque": true,
		"muy": true, "más": true, "menos": true, "todo": true,
		"todos": true, "todas": true, "algo": true, "nada": true,
		"qué": true, "quién": true, "cuál": true, "cómo": true,
		"también": true, "solo": true, "pero": true, "aunque": true,
		"mientras": true, "hasta": true, "desde": true, "hacia": true,
	}

	// Palabras comunes en inglés (high frequency)
	englishWords := map[string]bool{
		"the": true, "a": true, "an": true, "and": true,
		"or": true, "but": true, "in": true, "on": true,
		"at": true, "to": true, "for": true, "of": true,
		"with": true, "by": true, "from": true, "as": true,
		"is": true, "was": true, "are": true, "were": true,
		"been": true, "be": true, "have": true, "has": true,
		"had": true, "do": true, "does": true, "did": true,
		"will": true, "would": true, "could": true, "should": true,
		"this": true, "that": true, "these": true, "those": true,
		"i": true, "you": true, "he": true, "she": true,
		"it": true, "we": true, "they": true, "what": true,
		"which": true, "who": true, "when": true, "where": true,
		"why": true, "how": true, "all": true, "each": true,
		"every": true, "both": true, "few": true, "more": true,
		"most": true, "other": true, "some": true, "any": true,
	}

	spanishCount := 0
	englishCount := 0

	// Contar palabras en los primeros 20 words (suficiente para detección)
	maxWords := 20
	if len(words) < maxWords {
		maxWords = len(words)
	}

	for i := 0; i < maxWords; i++ {
		word := stripPunctuation(words[i])
		if spanishWords[word] {
			spanishCount++
		}
		if englishWords[word] {
			englishCount++
		}
	}

	// Determinar idioma basado en conteo
	if spanishCount > englishCount {
		return "es"
	}

	// Si hay más palabras en inglés o empate, default a inglés
	return "en"
}

// stripPunctuation elimina puntuación de una palabra
func stripPunctuation(word string) string {
	return strings.TrimFunc(word, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

// IsSpanish verifica si un texto está en español
func IsSpanish(text string) bool {
	return DetectLanguage(text) == "es"
}

// IsEnglish verifica si un texto está en inglés
func IsEnglish(text string) bool {
	return DetectLanguage(text) == "en"
}
