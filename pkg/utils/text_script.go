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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Environment variables for image generation
const (
	EnvImageGenProvider   = "PICOCLAW_TOOLS_IMAGE_GEN_PROVIDER"
	EnvGeminiAPIKey       = "PICOCLAW_TOOLS_IMAGE_GEN_GEMINI_API_KEY"
	EnvGeminiTextModel    = "PICOCLAW_TOOLS_IMAGE_GEN_GEMINI_TEXT_MODEL"
	EnvGeminiImageModel   = "PICOCLAW_TOOLS_IMAGE_GEN_GEMINI_IMAGE_MODEL"
	EnvIdeogramAPIKey     = "PICOCLAW_TOOLS_IMAGE_GEN_IDEOGRAM_API_KEY"
	EnvIdeogramAPIURL     = "PICOCLAW_TOOLS_IMAGE_GEN_IDEOGRAM_API_URL"
	EnvImageGenOutputDir  = "PICOCLAW_TOOLS_IMAGE_GEN_OUTPUT_DIR"
	EnvImageScriptPath    = "PICOCLAW_TOOLS_IMAGE_GEN_IMAGE_SCRIPT_PATH"
	EnvImageGenScriptPath = "PICOCLAW_TOOLS_IMAGE_GEN_IMAGE_GEN_SCRIPT_PATH"
	EnvAspectRatio        = "PICOCLAW_TOOLS_IMAGE_GEN_ASPECT_RATIO"
)

// System prompts embedded in English (defaults)
const (
	// DEFAULT_IMAGE_SCRIPT - Prompt to generate social post scripts
	DEFAULT_IMAGE_SCRIPT = `You are a professional social media copywriter.

Write a brief, clear, and engaging post script about the following topic.

Requirements:
- Length: concise (roughly 80-180 words)
- Hook in the first sentence
- Clear structure: intro, body, call-to-action
- Natural, conversational tone
- No markdown formatting
- Language: Match the language of the input topic

Topic: {topic}

Write the script:`

	// DEFAULT_IMAGE_GEN_SCRIPT - Prompt to generate visual prompts
	DEFAULT_IMAGE_GEN_SCRIPT = `You are an expert visual prompt engineer for AI image generation.

Create a detailed, cinematic image prompt based on the script content below.

Requirements:
- Highly descriptive, sensory language
- Specify lighting, mood, composition, style
- Include quality modifiers: "8k, ultra-detailed, cinematic lighting"
- End with: "No text, no typography, no watermarks"
- Format: Single paragraph, comma-separated descriptors
- Language: English (for best image generation results)

Script content:
{script_content}

Topic: {topic}

Create the visual prompt:`

	// DEFAULT_IMAGE_SCRIPT_EN - User's super-prompt for script generation (English version)
	DEFAULT_IMAGE_SCRIPT_EN = `Role and Mission:
You are a Screenwriter, Script Doctor, Storytelling/Narrative Designer, and Creative Director specialized in viral scripts and expanded storytelling. You master structure, pacing, tension, emotional progression, multiplatform narrative, and narrative marketing.
Your mission is to create brief viral scripts for Facebook that stop the scroll from the first line, maintain suspense until the end, and close with an emotional or cognitive impact.

Essential Context:
- Objective: Create irresistible micro-stories that generate retention, conversation, and virality.
- Audience: Facebook users (18-45 years old) who love intriguing stories, unusual facts, or unexpected twists.
- Scope/deliverable: Narrative text of maximum 1000 characters (without title).
- Restrictions: Meta-appropriate language, no emojis or headers, following indicated tone, structure, and pacing.
- Style/Tone: Direct, theatrical, sarcastic, and with intelligent humor.
- Output format: Only the final script text, without titles or additional notes.

---

CATEGORIES (detect before writing):
Before writing, identify the topic category and apply the corresponding tone:
- **Mystery, horror, legends, or myths →** Dark tone, with constant tension and expectation.
- **Conspiracy theories or UFO science →** Intriguing and provocative tone.
- **Empires, epics, and mega constructions →** Epic, majestic, and surprising tone.
- **Science in general →** Curious, stimulating, and informative tone, with powerful facts.
- **General news →** Curious, stimulating, and informative tone, make sure to give context about what the news is about.

---

WRITING INSTRUCTIONS:
Generate a brief narrative script, maximum 1200 characters (not counting the title), with a direct, sarcastic style and spicy humor.
It should feel like a self-contained mini-story, with:
1. Intriguing or disruptive beginning.
2. Development with growing tension or disturbing detail.
3. Impactful, ironic, or unexpected closing that leaves an echo in the viewer.

The final text must be optimized to engage the Facebook audience, with fast pacing, short sentences, and twists that force them to keep reading.

---

REAL FACTS RULE (mandatory):
If the {topic} deals with real facts, real people, or historical events, explicitly include the location and year within the script (without headers).
Recommended format: "City, Country, 19XX/20XX."
If the place is broad, use region/country. If only decade is available, indicate it: "United States, 1970s."
Do not invent data: if you are not sure, formulate the phrase as approximate location + period.

---

ENGAGEMENT:
Include phrases that invite commenting and sharing.
Do not repeat the same CTA between consecutive stories.

---

PACING AND TONE:
- Maintain an agile and cinematographic pace; avoid long paragraphs.
- Use conversational tone with theatrical moments at the climax.
- Prioritize humor and pacing over swear words; use only those that add value.

WORDS TO SOFTEN (mandatory, for Meta/Facebook):
Systematically substitute (respecting gender and number):
- "murder / murdered / killed" → "unalive / unalived"
- "murderer / murderers" → "unaliver / unalivers"
- "suicide / to commit suicide" → "unalive oneself"
Avoid explicit terms of direct violence; use periphrasis when necessary.

---

RETENTION OPTIMIZATION (Expanded Storytelling + Neuropsychology):
- Create an *information gap* from the first line.
- Apply the *Zeigarnik effect*: leave ideas open until the end.
- Introduce *micro-rewards* every 2-3 sentences: surprising facts or twists.
- Alternate tension and relief with variable pacing (short and descriptive sentences).
- Use sensory or evocative language (fear, amazement, curiosity, humor).
- Add *direct questions* ("Can you imagine being there?", "What would you have done?") to activate participation.
- Close all open narrative loops, releasing the Zeigarnik effect.
- End with an *emotional or cognitive echo* (surprise, irony, or reflection).

---

RECOMMENDED COGNITIVE STRUCTURE:
1. **Striking initial hook:** Generates immediate curiosity (stops the scroll).
2. **Intriguing introduction:** Poses an unsolved mystery or conflict.
3. **Progressive development:** Increases curiosity with facts and twists.
4. **Climax or revelation:** Presents the most powerful fact or main twist.
5. **Resolution or emotional closing:** Solves the mystery or leaves open reflection.

---

ACCEPTANCE CRITERIA:
- Must include: setup, development, twist, closing, and unique CTA.
- Must meet: cinematographic pacing, measured humor, narrative coherence, and appropriate language.
- Must not: repeat CTAs or use explicit language.
- Must maintain: intriguing tone, naturalness, and complete structure.

---
- HASHTAG CRITERIA:
- Add 1 one-word hashtag about the script topic or category; it should be the most relevant and, preferably, within the narration.

The topic to develop is: {topic}`

	// DEFAULT_IMAGE_GEN_SCRIPT_EN - User's super-prompt for visual prompt generation (English version)
	DEFAULT_IMAGE_GEN_SCRIPT_EN = `Read the script carefully:

{script_content}

Analyze and define based on the script:

Main theme (e.g., forbidden ritual, political conspiracy, scientific discovery, epic battle).

Narrative genre (horror, mystery, conspiracy, science, epic).

Country, city, and historical year of the narrated context (if applicable).

Key location (e.g., cursed forest, abandoned circus, futuristic laboratory, ancestral temple, ruined city).

Main character(s) with physical features, clothing, intense emotions, and specific cultural/national elements according to country and era (e.g., Mexican peasants of the 19th century, Roman soldiers in battle, American scientists in the 60s, etc.).

Central action or moment (the narrative climax, the revelation, or an instant of tension frozen in time).

Dominant atmosphere (oppressive, magical, inquieting, epic, futuristic, melancholic).

Precise lighting (e.g., moonlight over pale faces, flickering torches, neon light reflected in rain, golden rays on the horizon, volumetric light through a window, dramatic rim light, split lighting).

Writing the prompt in English

Write the prompt as if it were the description of a Hollywood cinematographic storyboard, in present tense, dynamic and visual:

Start with the action and main setting, including country, city, and historical era if applicable.

Describe in detail the characters (face, clothing, posture, emotional expression), taking care of their nationality, social class, and cultural-historical context.

Add a rich environment coherent with the place and year (architecture, objects, landscapes, or technology of the era).

Intensify the atmosphere and lighting to generate suspense, drama, or amazement, using cinematography and VFX terms.

Reinforce the exact moment that must be frozen as an epic frame from a blockbuster production.

Mandatory visual style

"Cinematic still from a Hollywood blockbuster, photorealistic, ultra-detailed, shot on film, anamorphic aspect ratio, immersive atmosphere, hyper-realistic textures, pores, fabric texture, believable imperfections, no stylization, no surreal distortions".

Color palette coherent with the genre

Horror → Dark, cold, desaturated tones, dense fog.

Mystery → Nebulous blues and grays, deep shadows.

Conspiracy → Somber contrasts, artificial lights, tense interiors.

Science → Bright metallics, vibrant neons, overflowing energy, colossal cosmic phenomena.

Epic → Golds, intense reds, dramatic and luminous skies.

Negative Prompt (mandatory at the end)

Negative prompt:
"IMAGES WITHOUT TEXT, subtitles, typography, inscriptions, written words, posters, signs, numbers, watermark, logo, UI, captions, blurry, low quality, oversaturated colors, distorted anatomy, duplicated faces, deformed hands, cartoonish style"

📌 Key note:
The output must be only the final visual prompt in English, without explanations or labels.
The result should feel like a high-budget cinematographic frame, with fidelity to the era, powerful visual design, and intense emotions.`
)

// TextScriptRequest defines parameters to generate a script
type TextScriptRequest struct {
	Topic        string
	Category     string // "news", "history", "tutorial", "announcement", etc.
	Language     string // "en" or "es" (auto-detected if empty)
	Duration     string // "30s", "60s", "5min"
	Tone         string // "professional", "casual", "engaging"
	TemplatePath string // Optional path to template file
}

// TextScriptResult contains the generated script
type TextScriptResult struct {
	Script            string
	WordCount         int
	EstimatedDuration string
	Language          string
	Template          string
}

// GeminiAPIRequest represents a request to Gemini API
type GeminiAPIRequest struct {
	Contents []GeminiContent `json:"contents"`
}

// GeminiContent represents the content of a request
type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

// GeminiPart represents a part of the content
type GeminiPart struct {
	Text string `json:"text"`
}

// GeminiAPIResponse represents a response from Gemini API
type GeminiAPIResponse struct {
	Candidates []GeminiCandidate `json:"candidates"`
}

// GeminiCandidate represents a response candidate
type GeminiCandidate struct {
	Content GeminiContent `json:"content"`
}

// LoadPromptTemplate loads a template from file or uses embedded default
func LoadPromptTemplate(customPath string, templateType string) (string, error) {
	// If custom path provided, try to load from file
	if customPath != "" {
		data, err := os.ReadFile(customPath)
		if err == nil {
			return strings.TrimSpace(string(data)), nil
		}
		// If fails, continue with embedded default
	}

	// Return embedded default by type
	switch templateType {
	case "script":
		return DEFAULT_IMAGE_SCRIPT, nil
	case "script_en":
		return DEFAULT_IMAGE_SCRIPT_EN, nil
	case "image_prompt":
		return DEFAULT_IMAGE_GEN_SCRIPT, nil
	case "image_prompt_en":
		return DEFAULT_IMAGE_GEN_SCRIPT_EN, nil
	default:
		return DEFAULT_IMAGE_SCRIPT, nil
	}
}

// GenerateTextScript generates a text script from a topic using Gemini API
func GenerateTextScript(
	ctx context.Context,
	apiKey string,
	model string,
	req TextScriptRequest,
) (*TextScriptResult, error) {
	if model == "" {
		model = "gemini-2.5-flash"
	}
	// Detect language if not specified
	if req.Language == "" {
		req.Language = DetectLanguage(req.Topic)
	}

	// Load template
	templatePath := req.TemplatePath
	var template string
	var err error

	// Script generation uses a single base template that already instructs
	// the model to match the topic language.
	template, err = LoadPromptTemplate(templatePath, "script")
	if err != nil {
		return nil, fmt.Errorf("error loading template: %v", err)
	}

	// Replace placeholders
	prompt := strings.ReplaceAll(template, "{topic}", req.Topic)
	prompt = strings.ReplaceAll(prompt, "{category}", req.Category)
	prompt = strings.ReplaceAll(prompt, "{duration}", req.Duration)
	prompt = strings.ReplaceAll(prompt, "{tone}", req.Tone)

	// Build request for Gemini API
	apiRequest := GeminiAPIRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{Text: prompt},
				},
			},
		},
	}

	jsonData, err := json.Marshal(apiRequest)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	// Gemini API URL
	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
		model,
		apiKey,
	)

	// Generate with retry
	var script string
	maxRetries := 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonData))
		if err != nil {
			return nil, fmt.Errorf("error creating request: %v", err)
		}
		httpReq.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 120 * time.Second}
		resp, err := client.Do(httpReq)
		if err != nil {
			if attempt < maxRetries {
				time.Sleep(time.Duration(attempt*2) * time.Second)
				continue
			}
			return nil, fmt.Errorf("error calling Gemini API: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			if attempt < maxRetries {
				time.Sleep(time.Duration(attempt*2) * time.Second)
				continue
			}
			return nil, fmt.Errorf("Gemini API error (%d): %s", resp.StatusCode, string(body))
		}

		var apiResp GeminiAPIResponse
		if err := json.Unmarshal(body, &apiResp); err != nil {
			return nil, fmt.Errorf("error unmarshaling response: %v", err)
		}

		if len(apiResp.Candidates) > 0 {
			for _, part := range apiResp.Candidates[0].Content.Parts {
				script += part.Text
			}
		}

		if script != "" {
			break
		}

		if attempt < maxRetries {
			time.Sleep(time.Duration(attempt*2) * time.Second)
		}
	}

	if script == "" {
		return nil, fmt.Errorf("empty response from Gemini")
	}

	script = strings.TrimSpace(script)

	// Calculate metrics
	wordCount := len(strings.Fields(script))
	estimatedDuration := "60s" // Default
	if wordCount < 100 {
		estimatedDuration = "30s"
	} else if wordCount > 200 {
		estimatedDuration = "90s"
	}

	return &TextScriptResult{
		Script:            script,
		WordCount:         wordCount,
		EstimatedDuration: estimatedDuration,
		Language:          req.Language,
		Template:          templatePath,
	}, nil
}

// BuildVisualPromptFromScript builds a visual prompt from a script
func BuildVisualPromptFromScript(
	ctx context.Context,
	apiKey, model, script, topic, aspectRatio, language, templatePath string,
) (string, error) {
	if model == "" {
		model = "gemini-2.5-flash"
	}
	// Load template
	var template string
	var err error

	// Visual prompt generation should always use the image prompt template
	// (or a custom file via templatePath), independent of topic language.
	template, err = LoadPromptTemplate(templatePath, "image_prompt")
	if err != nil {
		return "", fmt.Errorf("error loading template: %v", err)
	}

	// Replace placeholders - use {script_content} to match user's prompt
	prompt := strings.ReplaceAll(template, "{script_content}", script)
	prompt = strings.ReplaceAll(prompt, "{contenido_del_txt}", script)    // fallback for Spanish prompts
	prompt = strings.ReplaceAll(prompt, "{contenido_del_script}", script) // additional fallback
	prompt = strings.ReplaceAll(prompt, "{topic}", topic)

	// Build request for Gemini API
	apiRequest := GeminiAPIRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{Text: prompt},
				},
			},
		},
	}

	jsonData, err := json.Marshal(apiRequest)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	// Gemini API URL
	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
		model,
		apiKey,
	)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("error calling Gemini API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		// Fallback: build basic prompt
		visualPrompt := fmt.Sprintf(
			"Cinematic shot of %s, high quality, 8k, professional lighting. No text. Image aspect ratio %s.",
			topic,
			aspectRatio,
		)
		return visualPrompt, nil
	}

	var apiResp GeminiAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %v", err)
	}

	var visualPrompt string
	if len(apiResp.Candidates) > 0 {
		for _, part := range apiResp.Candidates[0].Content.Parts {
			visualPrompt += part.Text
		}
	}

	if visualPrompt == "" {
		// Fallback: build basic prompt
		visualPrompt = fmt.Sprintf("Cinematic shot of %s, high quality, 8k, professional lighting. No text.", topic)
	}

	// Add aspect ratio
	visualPrompt = strings.TrimSpace(visualPrompt)
	if !strings.Contains(strings.ToLower(visualPrompt), "aspect ratio") {
		visualPrompt += fmt.Sprintf(" . Image aspect ratio %s.", aspectRatio)
	}

	return visualPrompt, nil
}
