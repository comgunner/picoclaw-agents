package agent

import (
	"encoding/json"
	"os"
	"time"

	"github.com/comgunner/picoclaw/pkg/fileutil"
)

// AgentState persiste el estado del agente entre reinicios.
// Se guarda en disco de forma atómica para evitar corrupción.
type AgentState struct {
	CurrentModel string    `json:"current_model"`
	AgentID      string    `json:"agent_id"`
	Channel      string    `json:"channel"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Save guarda el estado en disco de forma atómica (temp + rename).
func (s *AgentState) Save(path string) error {
	s.UpdatedAt = time.Now()

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	return fileutil.WriteFileAtomic(path, data, 0o644)
}

// LoadAgentState carga el estado desde disco.
// Si el archivo no existe, retorna un estado vacío sin error.
func LoadAgentState(path string) (*AgentState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &AgentState{}, nil
		}
		return nil, err
	}

	var state AgentState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	return &state, nil
}
