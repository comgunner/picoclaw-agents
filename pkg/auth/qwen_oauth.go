// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package auth

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"time"
)

const (
	// QwenOAuthTimeout timeout máximo para completar el flujo OAuth
	QwenOAuthTimeout = 15 * time.Minute

	// QwenOAuthPollInterval intervalo entre polls para verificar completación
	QwenOAuthPollInterval = 2 * time.Second

	// QwenOAuthWaitInitial tiempo inicial de espera para que aparezca output
	QwenOAuthWaitInitial = 5 * time.Second

	// QwenOAuthMaxWaitOutput tiempo máximo para esperar output inicial
	QwenOAuthMaxWaitOutput = 15 * time.Second

	// QwenDeviceCodeLength longitud del device code (8 caracteres alfanuméricos)
	QwenDeviceCodeLength = 8
)

// OAuthResult contiene el resultado del flujo OAuth
type OAuthResult struct {
	Success    bool
	AuthURL    string
	DeviceCode string
	Error      error
}

// QwenOAuth gestiona el flujo de autenticación OAuth con tmux
// Permite extracción no bloqueante de URLs y device codes
type QwenOAuth struct {
	sessionManager *SessionManager
	tmuxSession    string
	outputFile     string
	workspace      string
}

// NewQwenOAuth crea una instancia de QwenOAuth
// workspace: directorio base del workspace (ej: ~/.picoclaw/workspace)
func NewQwenOAuth(workspace string) *QwenOAuth {
	timestamp := time.Now().Unix()
	return &QwenOAuth{
		sessionManager: NewSessionManager(workspace),
		tmuxSession:    fmt.Sprintf("qwen-oauth-%d", timestamp),
		outputFile:     fmt.Sprintf("/tmp/qwen-oauth-%d.txt", timestamp),
		workspace:      workspace,
	}
}

// StartOAuth inicia el flujo OAuth usando tmux para captura no bloqueante
// Ejecuta: picoclaw-agents auth login --provider qwen
// Extrae: OAuth URL y device code del output
func (o *QwenOAuth) StartOAuth(ctx context.Context) (*OAuthResult, error) {
	// Verificar dependencies
	if err := o.checkDependencies(); err != nil {
		return nil, err
	}

	// Registrar cleanup
	defer o.cleanup()

	// 1. Crear sesión tmux en background
	// El commando se ejecuta dentro de tmux para captura de output
	cmd := exec.CommandContext(ctx, "tmux", "new-session", "-d", "-s", o.tmuxSession,
		fmt.Sprintf("picoclaw-agents auth login --provider qwen 2>&1 | tee '%s'", o.outputFile))

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to create tmux session: %w", err)
	}

	// 2. Esperar output inicial (max QwenOAuthMaxWaitOutput)
	time.Sleep(QwenOAuthWaitInitial)

	// Polling para esperar output
	maxPolls := int(QwenOAuthMaxWaitOutput / time.Second)
	raw := make([]byte, 0, 4096) // Preallocate with reasonable capacity
	for i := 0; i < maxPolls; i++ {
		time.Sleep(time.Second)
		raw, _ = os.ReadFile(o.outputFile)
		if bytes.Contains(raw, []byte("authorize")) {
			break
		}
	}

	// 3. Capturar output completo del buffer tmux
	captureCmd := exec.Command("tmux", "capture-pane", "-t", o.tmuxSession, "-p")
	paneOut, err := captureCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to capture tmux output: %w", err)
	}

	// Append al archivo
	f, err := os.OpenFile(o.outputFile, os.O_APPEND|os.O_WRONLY, 0o600)
	if err == nil {
		f.Write(paneOut)
		f.Close()
	}

	// Combinar output
	raw = append(raw, paneOut...)
	output := string(raw)

	// 4. Extraer información OAuth
	result := o.extractOAuthInfo(output)
	if result.Error != nil {
		return result, nil
	}

	return result, nil
}

// extractOAuthInfo extrae URL y device code del output de tmux
// Usa regex para patrons específicos de Qwen OAuth
func (o *QwenOAuth) extractOAuthInfo(output string) *OAuthResult {
	result := &OAuthResult{Success: false}

	// Regex para URL de autorización
	// Pattern: https://chat.qwen.ai/authorize?user_code=XXXXX
	urlRegex := regexp.MustCompile(`https://chat\.qwen\.ai/authorize[^\s"]+`)
	authURL := urlRegex.FindString(output)

	if authURL == "" {
		result.Error = fmt.Errorf("no OAuth URL found in output")
		return result
	}

	// Regex para device code (8 caracteres mayúsculas/dígitos)
	// Pattern: Code: ABCDEFGH o user_code=ABCDEFGH en URL
	// Primero intentar con "Code: XXX" formato
	codeRegex := regexp.MustCompile(`Code:\s*([A-Z0-9]{8})(?:\s|$)`)
	codeMatch := codeRegex.FindStringSubmatch(output)

	deviceCode := ""
	if len(codeMatch) > 1 && codeMatch[1] != "" {
		deviceCode = codeMatch[1]
	} else {
		// Si no, intentar con user_code=XXX en URL
		urlRegex := regexp.MustCompile(`user_code=([A-Z0-9]{8})(?:&|$)`)
		urlMatch := urlRegex.FindStringSubmatch(output)
		if len(urlMatch) > 1 && urlMatch[1] != "" {
			deviceCode = urlMatch[1]
		}
	}

	// Validar formato del device code
	if deviceCode != "" {
		deviceCodeRegex := regexp.MustCompile(`^[A-Z0-9]{8}$`)
		if !deviceCodeRegex.MatchString(deviceCode) {
			deviceCode = "" // Invalidar si no matchea patrón exacto
		}
	}

	result.Success = true
	result.AuthURL = authURL
	result.DeviceCode = deviceCode

	return result
}

// WaitForCompletion espera a que el usuario complete la autorización en el navegador
// Polling cada QwenOAuthPollInterval para verificar si la sesión fue creada
func (o *QwenOAuth) WaitForCompletion(ctx context.Context, timeout time.Duration) error {
	if timeout == 0 {
		timeout = QwenOAuthTimeout
	}

	ticker := time.NewTicker(QwenOAuthPollInterval)
	defer ticker.Stop()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	startTime := time.Now()

	for {
		select {
		case <-ctx.Done():
			elapsed := time.Since(startTime)
			return fmt.Errorf("OAuth timeout after %v (max: %v)", elapsed, timeout)

		case <-ticker.C:
			session, err := o.sessionManager.Load()
			if err != nil {
				continue
			}

			if session != nil && session.AccessToken != "" {
				// Sesión creada exitosamente
				return nil
			}
		}
	}
}

// PollSession verifica el estado actual de la sesión
// Útil para mostrar progreso al usuario
func (o *QwenOAuth) PollSession() (*SessionInfo, error) {
	return o.sessionManager.GetSessionInfo()
}

// GetDeviceCode retorna el device code actual (si está disponible)
// Útil para mostrar al usuario durante el flujo
func (o *QwenOAuth) GetDeviceCode() (string, error) {
	// Leer output file y extraer device code
	data, err := os.ReadFile(o.outputFile)
	if err != nil {
		return "", err
	}

	result := o.extractOAuthInfo(string(data))
	if !result.Success {
		return "", fmt.Errorf("device code not found")
	}

	return result.DeviceCode, nil
}

// GetAuthURL retorna la URL de autorización (si está disponible)
func (o *QwenOAuth) GetAuthURL() (string, error) {
	// Leer output file y extraer URL
	data, err := os.ReadFile(o.outputFile)
	if err != nil {
		return "", err
	}

	result := o.extractOAuthInfo(string(data))
	if !result.Success {
		return "", fmt.Errorf("auth URL not found")
	}

	return result.AuthURL, nil
}

// GetSessionManager retorna el SessionManager para acceso externo
func (o *QwenOAuth) GetSessionManager() *SessionManager {
	return o.sessionManager
}

// checkDependencies verifica que tmux y picoclaw-agents estén instalados
func (o *QwenOAuth) checkDependencies() error {
	// Verificar tmux
	if _, err := exec.LookPath("tmux"); err != nil {
		return &DependencyError{
			Name:    "tmux",
			Message: "tmux not found in PATH (required for OAuth extraction)",
			Install: "Install tmux: brew install tmux (macOS) | sudo apt install tmux (Linux)",
		}
	}

	// Verificar picoclaw-agents
	if _, err := exec.LookPath("picoclaw-agents"); err != nil {
		return &DependencyError{
			Name:    "picoclaw-agents",
			Message: "picoclaw-agents not found in PATH",
			Install: "Build with: make build && sudo make install",
		}
	}

	return nil
}

// cleanup libera recursos (tmux session y archivo temporal)
func (o *QwenOAuth) cleanup() {
	// Matar sesión tmux si existe
	exec.Command("tmux", "kill-session", "-t", o.tmuxSession).Run()

	// Eliminar archivo temporal
	os.Remove(o.outputFile)
}

// DependencyError error por dependencia faltante
type DependencyError struct {
	Name    string
	Message string
	Install string
}

func (e *DependencyError) Error() string {
	return fmt.Sprintf("%s: %s", e.Message, e.Install)
}

// OpenBrowser abre una URL en el navegador predeterminado del sistema
// Soporta macOS, Linux y Windows
func OpenBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "linux":
		cmd = "xdg-open"
		args = []string{url}
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	execCmd := exec.Command(cmd, args...)
	execCmd.Stdout = nil
	execCmd.Stderr = nil
	return execCmd.Start()
}
