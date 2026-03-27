// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/comgunner/picoclaw/pkg/utils"
)

// ============== Notion Tools ==============

type NotionCreatePageTool struct {
	apiKey string
}

func NewNotionCreatePageTool() *NotionCreatePageTool {
	return NewNotionCreatePageToolFromConfig("")
}

func NewNotionCreatePageToolFromConfig(configAPIKey string) *NotionCreatePageTool {
	apiKey := strings.TrimSpace(os.Getenv(utils.EnvNotionAPIKey))
	if apiKey == "" {
		apiKey = strings.TrimSpace(configAPIKey)
	}
	return &NotionCreatePageTool{
		apiKey: apiKey,
	}
}

func (t *NotionCreatePageTool) Name() string {
	return "notion_create_page"
}

func (t *NotionCreatePageTool) Description() string {
	return "Crear una nueva página en un database de Notion. Requiere database_id y properties en formato JSON."
}

func (t *NotionCreatePageTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"database_id": map[string]any{
				"type":        "string",
				"description": "Database ID donde crear la página",
			},
			"properties": map[string]any{
				"type":        "object",
				"description": "Propiedades de la página en formato JSON. Ejemplo: {\"Name\": {\"title\": [{\"text\": {\"content\": \"Mi Página\"}}]}}",
				"properties":  map[string]any{},
			},
		},
		"required": []string{"database_id", "properties"},
	}
}

func (t *NotionCreatePageTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	databaseID, _ := args["database_id"].(string)
	return executeNotionTool(
		ctx,
		t.apiKey,
		databaseID,
		"database_id",
		args["properties"],
		func(client *utils.NotionClient, cCtx context.Context, id string, props map[string]any) (*utils.NotionPage, error) {
			return client.NotionCreatePage(cCtx, id, props)
		},
		"creada",
	)
}

// ============== Notion Query Database Tool ==============

type NotionQueryDatabaseTool struct {
	apiKey string
}

func NewNotionQueryDatabaseTool() *NotionQueryDatabaseTool {
	return NewNotionQueryDatabaseToolFromConfig("")
}

func NewNotionQueryDatabaseToolFromConfig(configAPIKey string) *NotionQueryDatabaseTool {
	apiKey := strings.TrimSpace(os.Getenv(utils.EnvNotionAPIKey))
	if apiKey == "" {
		apiKey = strings.TrimSpace(configAPIKey)
	}
	return &NotionQueryDatabaseTool{
		apiKey: apiKey,
	}
}

func (t *NotionQueryDatabaseTool) Name() string {
	return "notion_query_database"
}

func (t *NotionQueryDatabaseTool) Description() string {
	return "Consultar un data source (database) de Notion. Requiere data_source_id. Opcionalmente acepta filter y sorts."
}

func (t *NotionQueryDatabaseTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"data_source_id": map[string]any{
				"type":        "string",
				"description": "Data Source ID del database a consultar",
			},
			"filter": map[string]any{
				"type":        "object",
				"description": "Filtro opcional. Ejemplo: {\"property\": \"Status\", \"select\": {\"equals\": \"Active\"}}",
				"properties":  map[string]any{},
			},
			"sorts": map[string]any{
				"type":        "array",
				"description": "Ordenamiento opcional. Ejemplo: [{\"property\": \"Date\", \"direction\": \"descending\"}]",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"property":  map[string]any{"type": "string"},
						"direction": map[string]any{"type": "string"},
					},
				},
			},
		},
		"required": []string{"data_source_id"},
	}
}

func (t *NotionQueryDatabaseTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	dataSourceID, _ := args["data_source_id"].(string)
	filterRaw, hasFilter := args["filter"]
	sortsRaw, hasSorts := args["sorts"]

	// Validar credenciales
	if t.apiKey == "" {
		return UserResult(
			"Notion API Key no configurada. " +
				"Configura en config.json (tools.notion) o usa variable de entorno:\n" +
				"  PICOCLAW_TOOLS_NOTION_API_KEY",
		)
	}

	dataSourceID = strings.TrimSpace(dataSourceID)
	if dataSourceID == "" {
		return ErrorResult("data_source_id es requerido")
	}

	// Construir query
	query := &utils.NotionQueryRequest{
		PageSize: 10,
	}

	if hasFilter {
		switch v := filterRaw.(type) {
		case string:
			var filter any
			if err := json.Unmarshal([]byte(v), &filter); err != nil {
				return ErrorResult(fmt.Sprintf("filter JSON inválido: %v", err))
			}
			query.Filter = filter
		case map[string]any:
			query.Filter = v
		}
	}

	if hasSorts {
		switch v := sortsRaw.(type) {
		case string:
			var sorts []utils.Sort
			if err := json.Unmarshal([]byte(v), &sorts); err != nil {
				return ErrorResult(fmt.Sprintf("sorts JSON inválido: %v", err))
			}
			query.Sorts = sorts
		case []any:
			for _, s := range v {
				if sm, ok := s.(map[string]any); ok {
					sort := utils.Sort{}
					if prop, ok := sm["property"].(string); ok {
						sort.Property = prop
					}
					if dir, ok := sm["direction"].(string); ok {
						sort.Direction = dir
					}
					query.Sorts = append(query.Sorts, sort)
				}
			}
		}
	}

	callCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	client := utils.NewNotionClient(t.apiKey)
	result, err := client.NotionQueryDatabase(callCtx, dataSourceID, query)
	if err != nil {
		return ErrorResult(fmt.Sprintf("notion query database falló: %v", err)).WithError(err)
	}

	if len(result.Results) == 0 {
		return UserResult("No se encontraron resultados en el database.")
	}

	// Formatear resultados
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Se encontraron %d resultados:\n\n", len(result.Results)))
	for i, page := range result.Results {
		sb.WriteString(fmt.Sprintf("%d. ", i+1))
		// Extraer título si existe
		if titleProp, ok := page.Properties["Name"]; ok && len(titleProp.Title) > 0 {
			sb.WriteString(titleProp.Title[0].PlainText)
		} else {
			sb.WriteString(page.ID)
		}
		sb.WriteString(fmt.Sprintf(" - %s\n", page.URL))
	}

	return UserResult(sb.String())
}

// ============== Notion Search Tool ==============

type NotionSearchTool struct {
	apiKey string
}

func NewNotionSearchTool() *NotionSearchTool {
	return NewNotionSearchToolFromConfig("")
}

func NewNotionSearchToolFromConfig(configAPIKey string) *NotionSearchTool {
	apiKey := strings.TrimSpace(os.Getenv(utils.EnvNotionAPIKey))
	if apiKey == "" {
		apiKey = strings.TrimSpace(configAPIKey)
	}
	return &NotionSearchTool{
		apiKey: apiKey,
	}
}

func (t *NotionSearchTool) Name() string {
	return "notion_search"
}

func (t *NotionSearchTool) Description() string {
	return "Buscar páginas y data sources en Notion por título o contenido."
}

func (t *NotionSearchTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"query": map[string]any{
				"type":        "string",
				"description": "Término de búsqueda",
			},
		},
		"required": []string{"query"},
	}
}

func (t *NotionSearchTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	query, _ := args["query"].(string)

	// Validar credenciales
	if t.apiKey == "" {
		return UserResult(
			"Notion API Key no configurada. " +
				"Configura en config.json (tools.notion) o usa variable de entorno:\n" +
				"  PICOCLAW_TOOLS_NOTION_API_KEY",
		)
	}

	query = strings.TrimSpace(query)
	if query == "" {
		return ErrorResult("query es requerido")
	}

	callCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	client := utils.NewNotionClient(t.apiKey)
	result, err := client.NotionSearch(callCtx, query)
	if err != nil {
		return ErrorResult(fmt.Sprintf("notion search falló: %v", err)).WithError(err)
	}

	if len(result.Results) == 0 {
		return UserResult("No se encontraron resultados para la búsqueda.")
	}

	// Formatear resultados
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Se encontraron %d resultados:\n\n", len(result.Results)))
	for i, item := range result.Results {
		sb.WriteString(fmt.Sprintf("%d. ", i+1))
		if page, ok := item.(map[string]any); ok {
			if title, ok := page["title"].([]any); ok && len(title) > 0 {
				if t, ok := title[0].(map[string]any); ok {
					if plain, ok := t["plain_text"].(string); ok {
						sb.WriteString(plain)
					}
				}
			}
			if id, ok := page["id"].(string); ok {
				sb.WriteString(fmt.Sprintf(" (ID: %s)", id))
			}
		}
		sb.WriteString("\n")
	}

	return UserResult(sb.String())
}

// ============== Notion Update Page Tool ==============

type NotionUpdatePageTool struct {
	apiKey string
}

func NewNotionUpdatePageTool() *NotionUpdatePageTool {
	return NewNotionUpdatePageToolFromConfig("")
}

func NewNotionUpdatePageToolFromConfig(configAPIKey string) *NotionUpdatePageTool {
	apiKey := strings.TrimSpace(os.Getenv(utils.EnvNotionAPIKey))
	if apiKey == "" {
		apiKey = strings.TrimSpace(configAPIKey)
	}
	return &NotionUpdatePageTool{
		apiKey: apiKey,
	}
}

func (t *NotionUpdatePageTool) Name() string {
	return "notion_update_page"
}

func (t *NotionUpdatePageTool) Description() string {
	return "Actualizar propiedades de una página existente en Notion."
}

func (t *NotionUpdatePageTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"page_id": map[string]any{
				"type":        "string",
				"description": "Page ID de la página a actualizar",
			},
			"properties": map[string]any{
				"type":        "object",
				"description": "Propiedades a actualizar en formato JSON",
				"properties":  map[string]any{},
			},
		},
		"required": []string{"page_id", "properties"},
	}
}

func (t *NotionUpdatePageTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	pageID, _ := args["page_id"].(string)
	return executeNotionTool(
		ctx,
		t.apiKey,
		pageID,
		"page_id",
		args["properties"],
		func(client *utils.NotionClient, cCtx context.Context, id string, props map[string]any) (*utils.NotionPage, error) {
			return client.NotionUpdatePage(cCtx, id, props)
		},
		"actualizada",
	)
}

func executeNotionTool(
	ctx context.Context,
	apiKey string,
	resourceID string,
	resourceName string,
	propertiesRaw any,
	action func(*utils.NotionClient, context.Context, string, map[string]any) (*utils.NotionPage, error),
	pastVerb string,
) *ToolResult {
	if propertiesRaw == nil {
		return ErrorResult("properties es requerido")
	}

	// Validar credenciales
	if apiKey == "" {
		return UserResult(
			"Notion API Key no configurada. " +
				"Configura en config.json (tools.notion) o usa variable de entorno:\n" +
				"  PICOCLAW_TOOLS_NOTION_API_KEY",
		)
	}

	// Convertir properties
	properties, errResult := parseNotionProperties(propertiesRaw)
	if errResult != nil {
		return errResult
	}

	resourceID = strings.TrimSpace(resourceID)
	if resourceID == "" {
		return ErrorResult(fmt.Sprintf("%s es requerido", resourceName))
	}

	callCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	client := utils.NewNotionClient(apiKey)
	page, err := action(client, callCtx, resourceID, properties)
	if err != nil {
		return ErrorResult(fmt.Sprintf("notion action falló: %v", err)).WithError(err)
	}

	return UserResult(fmt.Sprintf("Página %s exitosamente. URL: %s", pastVerb, page.URL))
}

func parseNotionProperties(raw any) (map[string]any, *ToolResult) {
	var properties map[string]any
	switch v := raw.(type) {
	case string:
		if err := json.Unmarshal([]byte(v), &properties); err != nil {
			return nil, ErrorResult(fmt.Sprintf("properties JSON inválido: %v", err))
		}
	case map[string]any:
		properties = v
	default:
		return nil, ErrorResult("properties debe ser un objeto JSON o string JSON")
	}
	return properties, nil
}
