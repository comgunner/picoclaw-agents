// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package channels

import (
	"context"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimiterConfig configuración de rate limiting por canal
type RateLimiterConfig struct {
	RequestsPerSecond rate.Limit // Peticiones por segundo
	BurstSize         int        // Tamaño máximo de burst
	Enabled           bool       // Habilitar/deshabilitar rate limiting
}

// DefaultRateLimitConfig configuración por defecto: 10 mensajes por minuto
var DefaultRateLimitConfig = RateLimiterConfig{
	RequestsPerSecond: rate.Limit(10.0 / 60.0), // 10 mensajes por minuto
	BurstSize:         5,                       // Permitir bursts de hasta 5 mensajes
	Enabled:           true,
}

// RateLimiter controla la frecuencia de mensajes por usuario
type RateLimiter struct {
	mu       sync.RWMutex
	limiters map[string]*rate.Limiter // userID -> limiter
	config   RateLimiterConfig
}

// NewRateLimiter crea un nuevo rate limiter con la configuración dada
func NewRateLimiter(config RateLimiterConfig) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		config:   config,
	}
}

// Allow verifica si un usuario puede enviar un mensaje
// Retorna true si el mensaje está permitido, false si excede el límite
func (rl *RateLimiter) Allow(userID string) bool {
	if !rl.config.Enabled {
		return true // Rate limiting deshabilitado
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[userID]
	if !exists {
		limiter = rate.NewLimiter(rl.config.RequestsPerSecond, rl.config.BurstSize)
		rl.limiters[userID] = limiter
	}

	return limiter.Allow()
}

// Wait espera hasta que el usuario pueda enviar un mensaje
// Bloquea hasta que el rate limiter lo permita
func (rl *RateLimiter) Wait(userID string) error {
	if !rl.config.Enabled {
		return nil
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[userID]
	if !exists {
		limiter = rate.NewLimiter(rl.config.RequestsPerSecond, rl.config.BurstSize)
		rl.limiters[userID] = limiter
	}

	return limiter.Wait(context.Background())
}

// AllowN verifica si un usuario puede enviar N mensajes
// Útil para mensajes que consumen múltiples tokens
func (rl *RateLimiter) AllowN(userID string, n int) bool {
	if !rl.config.Enabled {
		return true
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[userID]
	if !exists {
		limiter = rate.NewLimiter(rl.config.RequestsPerSecond, rl.config.BurstSize)
		rl.limiters[userID] = limiter
	}

	return limiter.AllowN(time.Now(), n)
}

// Reset resetea el rate limiter para un usuario específico
// Útil para testing o para resetear límites manualmente
func (rl *RateLimiter) Reset(userID string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	delete(rl.limiters, userID)
}

// ResetAll resetea todos los rate limiters
func (rl *RateLimiter) ResetAll() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.limiters = make(map[string]*rate.Limiter)
}

// GetUserLimitersCount retorna el número de usuarios con rate limiters activos
func (rl *RateLimiter) GetUserLimitersCount() int {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	return len(rl.limiters)
}

// GetLimiter retorna el limiter para un usuario (para testing)
func (rl *RateLimiter) GetLimiter(userID string) *rate.Limiter {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	return rl.limiters[userID]
}

// RateLimitedError error cuando se excede el rate limit
type RateLimitedError struct {
	UserID        string
	RetryAfter    time.Duration
	Limit         rate.Limit
	BurstSize     int
	CurrentTokens float64
}

func (e *RateLimitedError) Error() string {
	return "rate limit exceeded"
}

// CheckWithRetry verifica si un usuario puede enviar un mensaje
// Retorna error con información de retry si está limitado
func (rl *RateLimiter) CheckWithRetry(userID string) error {
	if !rl.config.Enabled {
		return nil
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[userID]
	if !exists {
		limiter = rate.NewLimiter(rl.config.RequestsPerSecond, rl.config.BurstSize)
		rl.limiters[userID] = limiter
	}

	if limiter.Allow() {
		return nil
	}

	// Calculator tiempo de retry
	retryAfter := limiter.Reserve().Delay()

	return &RateLimitedError{
		UserID:        userID,
		RetryAfter:    retryAfter,
		Limit:         rl.config.RequestsPerSecond,
		BurstSize:     rl.config.BurstSize,
		CurrentTokens: limiter.Tokens(),
	}
}
