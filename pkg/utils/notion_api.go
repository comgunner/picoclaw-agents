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

const (
	EnvNotionAPIKey = "PICOCLAW_TOOLS_NOTION_API_KEY"
	NotionVersion   = "2025-09-03"
)

// NotionClient representa un cliente para la API de Notion
type NotionClient struct {
	APIKey     string
	HTTPClient *http.Client
}

// NotionPage representa una página de Notion
type NotionPage struct {
	Object             string              `json:"object"`
	ID                 string              `json:"id"`
	CreatedTime        string              `json:"created_time"`
	LastEditedTime     string              `json:"last_edited_time"`
	CreatedBy          *NotionUser         `json:"created_by,omitempty"`
	LastEditedBy       *NotionUser         `json:"last_edited_by,omitempty"`
	Cover              any                 `json:"cover,omitempty"`
	Icon               any                 `json:"icon,omitempty"`
	Parent             NotionParent        `json:"parent"`
	Archived           bool                `json:"archived"`
	Properties         map[string]Property `json:"properties"`
	URL                string              `json:"url"`
	PublicURL          string              `json:"public_url,omitempty"`
	DataSourceParentID string              `json:"data_source_parent_id,omitempty"`
	DataSourceID       string              `json:"data_source_id,omitempty"`
}

// NotionUser representa un usuario de Notion
type NotionUser struct {
	Object string `json:"object"`
	ID     string `json:"id"`
}

// NotionParent representa el padre de una página o database
type NotionParent struct {
	Type       string `json:"type"`
	DatabaseID string `json:"database_id,omitempty"`
	PageID     string `json:"page_id,omitempty"`
}

// Property representa una propiedad de página
type Property struct {
	ID          string     `json:"id"`
	Type        string     `json:"type"`
	Title       []RichText `json:"title,omitempty"`
	RichText    []RichText `json:"rich_text,omitempty"`
	Select      *Select    `json:"select,omitempty"`
	MultiSelect []Select   `json:"multi_select,omitempty"`
	Date        *Date      `json:"date,omitempty"`
	Checkbox    bool       `json:"checkbox,omitempty"`
	Number      float64    `json:"number,omitempty"`
	URL         string     `json:"url,omitempty"`
	Email       string     `json:"email,omitempty"`
	Relation    []Relation `json:"relation,omitempty"`
}

// RichText representa texto enriquecido
type RichText struct {
	Type        string       `json:"type"`
	Text        *TextContent `json:"text,omitempty"`
	Annotations *Annotations `json:"annotations,omitempty"`
	PlainText   string       `json:"plain_text,omitempty"`
}

// TextContent representa el contenido de texto
type TextContent struct {
	Content string `json:"content"`
	Link    *Link  `json:"link,omitempty"`
}

// Link representa un enlace
type Link struct {
	URL string `json:"url"`
}

// Annotations representa anotaciones de formato
type Annotations struct {
	Bold          bool   `json:"bold"`
	Italic        bool   `json:"italic"`
	Strikethrough bool   `json:"strikethrough"`
	Underline     bool   `json:"underline"`
	Code          bool   `json:"code"`
	Color         string `json:"color"`
}

// Select representa una opción de select
type Select struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name"`
	Color string `json:"color,omitempty"`
}

// Date representa una fecha
type Date struct {
	Start string `json:"start"`
	End   string `json:"end,omitempty"`
}

// Relation representa una relación
type Relation struct {
	ID string `json:"id"`
}

// NotionDataSource representa un data source (database)
type NotionDataSource struct {
	Object       string              `json:"object"`
	ID           string              `json:"id"`
	DataSourceID string              `json:"data_source_id"`
	Title        []RichText          `json:"title"`
	Description  []RichText          `json:"description"`
	Properties   map[string]Property `json:"properties"`
	Parent       NotionParent        `json:"parent"`
	URL          string              `json:"url"`
	PublicURL    string              `json:"public_url,omitempty"`
	Archived     bool                `json:"archived"`
	IsInline     bool                `json:"is_inline"`
}

// NotionSearchRequest representa una petición de búsqueda
type NotionSearchRequest struct {
	Query       string        `json:"query,omitempty"`
	Filter      *SearchFilter `json:"filter,omitempty"`
	Sort        *SearchSort   `json:"sort,omitempty"`
	StartCursor string        `json:"start_cursor,omitempty"`
	PageSize    int           `json:"page_size,omitempty"`
}

// SearchFilter representa un filtro de búsqueda
type SearchFilter struct {
	Value string `json:"value"`
}

// SearchSort representa ordenamiento de búsqueda
type SearchSort struct {
	Direction string `json:"direction"`
	Timestamp string `json:"timestamp"`
}

// NotionSearchResponse representa una respuesta de búsqueda
type NotionSearchResponse struct {
	Object     string `json:"object"`
	Results    []any  `json:"results"`
	NextCursor string `json:"next_cursor"`
	HasMore    bool   `json:"has_more"`
}

// NotionQueryRequest representa una petición de query a database
type NotionQueryRequest struct {
	Filter      any    `json:"filter,omitempty"`
	Sorts       []Sort `json:"sorts,omitempty"`
	StartCursor string `json:"start_cursor,omitempty"`
	PageSize    int    `json:"page_size,omitempty"`
}

// Sort representa ordenamiento
type Sort struct {
	Property  string `json:"property"`
	Direction string `json:"direction"`
}

// NotionQueryResponse representa una respuesta de query
type NotionQueryResponse struct {
	Object     string       `json:"object"`
	Results    []NotionPage `json:"results"`
	NextCursor string       `json:"next_cursor"`
	HasMore    bool         `json:"has_more"`
}

// NewNotionClient crea un nuevo cliente de Notion
func NewNotionClient(apiKey string) *NotionClient {
	return &NotionClient{
		APIKey:     apiKey,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// NewNotionClientFromEnv crea un cliente de Notion desde variables de entorno
func NewNotionClientFromEnv() (*NotionClient, error) {
	apiKey := strings.TrimSpace(os.Getenv(EnvNotionAPIKey))
	if apiKey == "" {
		return nil, fmt.Errorf("NOTION_API_KEY no configurada")
	}
	return NewNotionClient(apiKey), nil
}

// createHeaders devuelve los headers necesarios para las peticiones a Notion
func (c *NotionClient) createHeaders() http.Header {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+c.APIKey)
	headers.Set("Notion-Version", NotionVersion)
	headers.Set("Content-Type", "application/json")
	return headers
}

// NotionCreatePage crea una nueva página en un database
func (c *NotionClient) NotionCreatePage(
	ctx context.Context,
	databaseID string,
	properties map[string]any,
) (*NotionPage, error) {
	payload := map[string]any{
		"parent": map[string]string{
			"database_id": databaseID,
			"type":        "database_id",
		},
		"properties": properties,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializando payload: %v", err)
	}

	url := "https://api.notion.com/v1/pages"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creando request: %v", err)
	}
	req.Header = c.createHeaders()

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error de red creando página: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error creando página (%d): %s", resp.StatusCode, string(body))
	}

	var page NotionPage
	if err := json.Unmarshal(body, &page); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta: %v", err)
	}

	return &page, nil
}

// NotionQueryDatabase consulta un data source (database)
func (c *NotionClient) NotionQueryDatabase(
	ctx context.Context,
	dataSourceID string,
	query *NotionQueryRequest,
) (*NotionQueryResponse, error) {
	var payload any = map[string]any{}
	if query != nil {
		payload = query
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializando payload: %v", err)
	}

	url := fmt.Sprintf("https://api.notion.com/v1/data_sources/%s/query", dataSourceID)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creando request: %v", err)
	}
	req.Header = c.createHeaders()

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error de red consultando database: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error consultando database (%d): %s", resp.StatusCode, string(body))
	}

	var result NotionQueryResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta: %v", err)
	}

	return &result, nil
}

// NotionSearch busca páginas y data sources en Notion
func (c *NotionClient) NotionSearch(ctx context.Context, searchQuery string) (*NotionSearchResponse, error) {
	payload := NotionSearchRequest{
		Query:    searchQuery,
		PageSize: 10,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializando payload: %v", err)
	}

	url := "https://api.notion.com/v1/search"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creando request: %v", err)
	}
	req.Header = c.createHeaders()

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error de red buscando: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error buscando (%d): %s", resp.StatusCode, string(body))
	}

	var result NotionSearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta: %v", err)
	}

	return &result, nil
}

// NotionGetPage obtiene una página por su ID
func (c *NotionClient) NotionGetPage(ctx context.Context, pageID string) (*NotionPage, error) {
	var page NotionPage
	url := fmt.Sprintf("https://api.notion.com/v1/pages/%s", pageID)
	if err := c.doGet(ctx, url, "página", &page); err != nil {
		return nil, err
	}
	return &page, nil
}

// NotionUpdatePage actualiza las propiedades de una página
func (c *NotionClient) NotionUpdatePage(
	ctx context.Context,
	pageID string,
	properties map[string]any,
) (*NotionPage, error) {
	payload := map[string]any{
		"properties": properties,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializando payload: %v", err)
	}

	url := fmt.Sprintf("https://api.notion.com/v1/pages/%s", pageID)
	req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creando request: %v", err)
	}
	req.Header = c.createHeaders()

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error de red actualizando página: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error actualizando página (%d): %s", resp.StatusCode, string(body))
	}

	var page NotionPage
	if err := json.Unmarshal(body, &page); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta: %v", err)
	}

	return &page, nil
}

// NotionGetDataSource obtiene un data source por su ID
func (c *NotionClient) NotionGetDataSource(ctx context.Context, dataSourceID string) (*NotionDataSource, error) {
	var ds NotionDataSource
	url := fmt.Sprintf("https://api.notion.com/v1/data_sources/%s", dataSourceID)
	if err := c.doGet(ctx, url, "data source", &ds); err != nil {
		return nil, err
	}
	return &ds, nil
}

func (c *NotionClient) doGet(ctx context.Context, url string, resourceName string, result any) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creando request: %v", err)
	}
	req.Header = c.createHeaders()

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error de red obteniendo %s: %v", resourceName, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error leyendo respuesta: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error obteniendo %s (%d): %s", resourceName, resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("error decodificando respuesta: %v", err)
	}

	return nil
}

// NotionCreateBlock añade bloques de contenido a una página
func (c *NotionClient) NotionCreateBlock(ctx context.Context, pageID string, blocks []map[string]any) error {
	payload := map[string]any{
		"children": blocks,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error serializando payload: %v", err)
	}

	url := fmt.Sprintf("https://api.notion.com/v1/blocks/%s/children", pageID)
	req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("error creando request: %v", err)
	}
	req.Header = c.createHeaders()

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error de red creando bloques: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error leyendo respuesta: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error creando bloques (%d): %s", resp.StatusCode, string(body))
	}

	return nil
}
