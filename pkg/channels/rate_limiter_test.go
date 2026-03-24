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
	"fmt"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestRateLimiter_Allow(t *testing.T) {
	config := RateLimiterConfig{
		RequestsPerSecond: rate.Limit(1.0), // 1 mensaje por segundo
		BurstSize:         3,
		Enabled:           true,
	}

	rl := NewRateLimiter(config)
	userID := "test_user"

	// Los primeros 3 mensajes deberían ser permitidos (burst)
	for i := 0; i < 3; i++ {
		if !rl.Allow(userID) {
			t.Errorf("Allow() = false, want true (burst message %d)", i+1)
		}
	}

	// El 4to mensaje debería ser rechazado (excede burst)
	if rl.Allow(userID) {
		t.Error("Allow() = true, want false (exceeded burst)")
	}
}

func TestRateLimiter_AllowDisabled(t *testing.T) {
	config := RateLimiterConfig{
		RequestsPerSecond: rate.Limit(1.0),
		BurstSize:         1,
		Enabled:           false,
	}

	rl := NewRateLimiter(config)
	userID := "test_user"

	// Debería permitir todos los mensajes cuando está deshabilitado
	for i := 0; i < 10; i++ {
		if !rl.Allow(userID) {
			t.Errorf("Allow() = false, want true (disabled rate limiter)")
		}
	}
}

func TestRateLimiter_MultipleUsers(t *testing.T) {
	config := RateLimiterConfig{
		RequestsPerSecond: rate.Limit(1.0),
		BurstSize:         2,
		Enabled:           true,
	}

	rl := NewRateLimiter(config)

	user1 := "user1"
	user2 := "user2"

	// Cada usuario tiene su propio limiter
	// User1 consume su burst
	rl.Allow(user1)
	rl.Allow(user1)

	// User2 debería poder enviar mensajes independientemente
	if !rl.Allow(user2) {
		t.Error("Allow() = false for user2, want true (independent limiters)")
	}
}

func TestRateLimiter_Wait(t *testing.T) {
	config := RateLimiterConfig{
		RequestsPerSecond: rate.Limit(10.0), // 10 mensajes por segundo
		BurstSize:         1,
		Enabled:           true,
	}

	rl := NewRateLimiter(config)
	userID := "test_user"

	// Consumir el burst
	rl.Allow(userID)

	// Wait debería bloquear brevemente y luego permitir
	done := make(chan error, 1)
	go func() {
		done <- rl.Wait(userID)
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Wait() error = %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Error("Wait() timed out")
	}
}

func TestRateLimiter_AllowN(t *testing.T) {
	config := RateLimiterConfig{
		RequestsPerSecond: rate.Limit(1.0),
		BurstSize:         5,
		Enabled:           true,
	}

	rl := NewRateLimiter(config)
	userID := "test_user"

	// AllowN con N=3 debería ser permitido (dentro del burst)
	if !rl.AllowN(userID, 3) {
		t.Error("AllowN(3) = false, want true")
	}

	// AllowN con N=3 debería ser rechazado (excede burst restante)
	if rl.AllowN(userID, 3) {
		t.Error("AllowN(3) = true, want false")
	}
}

func TestRateLimiter_Reset(t *testing.T) {
	config := RateLimiterConfig{
		RequestsPerSecond: rate.Limit(1.0),
		BurstSize:         2,
		Enabled:           true,
	}

	rl := NewRateLimiter(config)
	userID := "test_user"

	// Consumir burst
	rl.Allow(userID)
	rl.Allow(userID)

	// Resetear
	rl.Reset(userID)

	// Debería permitir mensajes nuevamente
	if !rl.Allow(userID) {
		t.Error("Allow() = false after Reset(), want true")
	}
}

func TestRateLimiter_ResetAll(t *testing.T) {
	config := RateLimiterConfig{
		RequestsPerSecond: rate.Limit(1.0),
		BurstSize:         2,
		Enabled:           true,
	}

	rl := NewRateLimiter(config)

	// Crear múltiples usuarios
	rl.Allow("user1")
	rl.Allow("user2")
	rl.Allow("user3")

	// Resetear todos
	rl.ResetAll()

	// Verificar que se resetearon
	count := rl.GetUserLimitersCount()
	if count != 0 {
		t.Errorf("GetUserLimitersCount() = %d, want 0", count)
	}
}

func TestRateLimiter_GetUserLimitersCount(t *testing.T) {
	config := RateLimiterConfig{
		RequestsPerSecond: rate.Limit(1.0),
		BurstSize:         2,
		Enabled:           true,
	}

	rl := NewRateLimiter(config)

	// Inicialmente debería estar vacío
	if count := rl.GetUserLimitersCount(); count != 0 {
		t.Errorf("GetUserLimitersCount() = %d, want 0", count)
	}

	// Crear 5 usuarios diferentes
	for i := 0; i < 5; i++ {
		userID := fmt.Sprintf("user%d", i)
		rl.Allow(userID)
	}

	if count := rl.GetUserLimitersCount(); count != 5 {
		t.Errorf("GetUserLimitersCount() = %d, want 5", count)
	}
}

func TestRateLimiter_CheckWithRetry(t *testing.T) {
	config := RateLimiterConfig{
		RequestsPerSecond: rate.Limit(1.0),
		BurstSize:         1,
		Enabled:           true,
	}

	rl := NewRateLimiter(config)
	userID := "test_user"

	// Primer mensaje debería ser permitido
	err := rl.CheckWithRetry(userID)
	if err != nil {
		t.Errorf("CheckWithRetry() error = %v, want nil", err)
	}

	// Segundo mensaje debería retornar error
	err = rl.CheckWithRetry(userID)
	if err == nil {
		t.Error("CheckWithRetry() error = nil, want *RateLimitedError")
	}

	// Verificar tipo de error
	rateErr, ok := err.(*RateLimitedError)
	if !ok {
		t.Errorf("Error type = %T, want *RateLimitedError", err)
	} else {
		if rateErr.UserID != userID {
			t.Errorf("RateLimitedError.UserID = %s, want %s", rateErr.UserID, userID)
		}
		if rateErr.RetryAfter <= 0 {
			t.Error("RateLimitedError.RetryAfter should be > 0")
		}
	}
}

func TestRateLimiter_CheckWithRetryDisabled(t *testing.T) {
	config := RateLimiterConfig{
		RequestsPerSecond: rate.Limit(1.0),
		BurstSize:         1,
		Enabled:           false,
	}

	rl := NewRateLimiter(config)
	userID := "test_user"

	// Debería permitir todos los mensajes cuando está deshabilitado
	for i := 0; i < 10; i++ {
		err := rl.CheckWithRetry(userID)
		if err != nil {
			t.Errorf("CheckWithRetry() error = %v, want nil (disabled)", err)
		}
	}
}

func TestRateLimitedError_Error(t *testing.T) {
	err := &RateLimitedError{
		UserID:     "test_user",
		RetryAfter: 5 * time.Second,
		Limit:      rate.Limit(1.0),
		BurstSize:  5,
	}

	expected := "rate limit exceeded"
	if err.Error() != expected {
		t.Errorf("Error() = %q, want %q", err.Error(), expected)
	}
}

func TestDefaultRateLimitConfig(t *testing.T) {
	// Verificar que la configuración por defecto es razonable
	if DefaultRateLimitConfig.RequestsPerSecond <= 0 {
		t.Error("DefaultRateLimitConfig.RequestsPerSecond should be > 0")
	}

	if DefaultRateLimitConfig.BurstSize <= 0 {
		t.Error("DefaultRateLimitConfig.BurstSize should be > 0")
	}

	if !DefaultRateLimitConfig.Enabled {
		t.Error("DefaultRateLimitConfig.Enabled should be true")
	}

	// 10 mensajes por minuto = 10/60 = 0.166... mensajes por segundo
	expected := rate.Limit(10.0 / 60.0)
	if DefaultRateLimitConfig.RequestsPerSecond != expected {
		t.Errorf("DefaultRateLimitConfig.RequestsPerSecond = %v, want %v",
			DefaultRateLimitConfig.RequestsPerSecond, expected)
	}
}

func TestRateLimiter_ConcurrentAccess(t *testing.T) {
	config := RateLimiterConfig{
		RequestsPerSecond: rate.Limit(100.0),
		BurstSize:         50,
		Enabled:           true,
	}

	rl := NewRateLimiter(config)
	userID := "test_user"

	// Acceder concurrentemente desde múltiples goroutines
	done := make(chan bool, 100)
	for i := 0; i < 100; i++ {
		go func() {
			rl.Allow(userID)
			done <- true
		}()
	}

	// Esperar a que todas las goroutines terminen
	for i := 0; i < 100; i++ {
		<-done
	}

	// Si llegamos aquí sin panic, el test pasó
}
