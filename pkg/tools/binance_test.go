// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)

package tools

import (
	"context"
	"strings"
	"testing"

	"github.com/comgunner/picoclaw/pkg/utils"
)

func TestNewBinanceSpotBalanceToolFromConfig_UsesConfigWhenEnvEmpty(t *testing.T) {
	t.Setenv(utils.EnvBinanceAPIKey, "")
	t.Setenv(utils.EnvBinanceSecretKey, "")

	tool := NewBinanceSpotBalanceToolFromConfig("cfg-api", "cfg-secret")
	if tool.apiKey != "cfg-api" {
		t.Fatalf("apiKey = %q, want %q", tool.apiKey, "cfg-api")
	}
	if tool.secretKey != "cfg-secret" {
		t.Fatalf("secretKey = %q, want %q", tool.secretKey, "cfg-secret")
	}
}

func TestNewBinanceSpotBalanceToolFromConfig_EnvOverridesConfig(t *testing.T) {
	t.Setenv(utils.EnvBinanceAPIKey, "env-api")
	t.Setenv(utils.EnvBinanceSecretKey, "env-secret")

	tool := NewBinanceSpotBalanceToolFromConfig("cfg-api", "cfg-secret")
	if tool.apiKey != "env-api" {
		t.Fatalf("apiKey = %q, want %q", tool.apiKey, "env-api")
	}
	if tool.secretKey != "env-secret" {
		t.Fatalf("secretKey = %q, want %q", tool.secretKey, "env-secret")
	}
}

func TestNormalizeFuturesDirection(t *testing.T) {
	got, err := normalizeFuturesDirection("LONG")
	if err != nil || got != "BUY" {
		t.Fatalf("normalize LONG = %q, err=%v", got, err)
	}
	got, err = normalizeFuturesDirection("SHORT")
	if err != nil || got != "SELL" {
		t.Fatalf("normalize SHORT = %q, err=%v", got, err)
	}
	_, err = normalizeFuturesDirection("INVALID")
	if err == nil {
		t.Fatal("expected error for invalid side")
	}
}

func TestParsePositiveDecimalArg(t *testing.T) {
	args := map[string]any{"quantity": "0.001"}
	got, err := parsePositiveDecimalArg(args, "quantity")
	if err != nil || got != "0.001" {
		t.Fatalf("string quantity parse failed: got=%q err=%v", got, err)
	}

	args = map[string]any{"quantity": float64(2.5)}
	got, err = parsePositiveDecimalArg(args, "quantity")
	if err != nil || got != "2.5" {
		t.Fatalf("float quantity parse failed: got=%q err=%v", got, err)
	}

	args = map[string]any{"quantity": float64(0)}
	_, err = parsePositiveDecimalArg(args, "quantity")
	if err == nil {
		t.Fatal("expected error for non-positive quantity")
	}
}

func TestBinanceFuturesTools_MissingCreds(t *testing.T) {
	openTool := NewBinanceFuturesOpenPositionToolFromConfig("", "")
	closeTool := NewBinanceFuturesClosePositionToolFromConfig("", "")

	openRes := openTool.Execute(context.Background(), map[string]any{
		"symbol":   "BTCUSDT",
		"side":     "LONG",
		"quantity": "0.001",
		"confirm":  true,
	})
	if !strings.Contains(openRes.ForUser, "Trading operations are disabled") {
		t.Fatalf("unexpected open missing-creds message: %q", openRes.ForUser)
	}

	closeRes := closeTool.Execute(context.Background(), map[string]any{
		"symbol":  "BTCUSDT",
		"confirm": true,
	})
	if !strings.Contains(closeRes.ForUser, "Trading operations are disabled") {
		t.Fatalf("unexpected close missing-creds message: %q", closeRes.ForUser)
	}
}

func TestBinanceFuturesTools_ConfirmRequired(t *testing.T) {
	t.Setenv(utils.EnvBinanceAPIKey, "env-api")
	t.Setenv(utils.EnvBinanceSecretKey, "env-secret")

	openTool := NewBinanceFuturesOpenPositionToolFromConfig("cfg-api", "cfg-secret")
	closeTool := NewBinanceFuturesClosePositionToolFromConfig("cfg-api", "cfg-secret")

	openRes := openTool.Execute(context.Background(), map[string]any{
		"symbol":   "BTCUSDT",
		"side":     "LONG",
		"quantity": "0.001",
		"confirm":  false,
	})
	if !openRes.IsError || !strings.Contains(openRes.ForLLM, "requires confirm: true") {
		t.Fatalf("expected confirm error for open tool, got: isError=%v msg=%q", openRes.IsError, openRes.ForLLM)
	}

	closeRes := closeTool.Execute(context.Background(), map[string]any{
		"symbol":  "BTCUSDT",
		"confirm": false,
	})
	if !closeRes.IsError || !strings.Contains(closeRes.ForLLM, "requires confirm: true") {
		t.Fatalf("expected confirm error for close tool, got: isError=%v msg=%q", closeRes.IsError, closeRes.ForLLM)
	}
}

func TestBinanceFuturesBalanceTool_MissingCreds(t *testing.T) {
	tool := NewBinanceFuturesBalanceToolFromConfig("", "")
	res := tool.Execute(context.Background(), map[string]any{})
	if !strings.Contains(res.ForUser, "Futures balance queries require API credentials") {
		t.Fatalf("unexpected futures balance missing-creds message: %q", res.ForUser)
	}
}

func TestBinanceOrderBookTool_Validation(t *testing.T) {
	tool := NewBinanceOrderBookTool()

	res := tool.Execute(context.Background(), map[string]any{})
	if !res.IsError || !strings.Contains(res.ForLLM, "symbol is required") {
		t.Fatalf("expected symbol validation error, got isError=%v msg=%q", res.IsError, res.ForLLM)
	}

	res = tool.Execute(context.Background(), map[string]any{
		"symbol": "BTCUSDT",
		"limit":  "abc",
	})
	if !res.IsError || !strings.Contains(res.ForLLM, "limit must be an integer") {
		t.Fatalf("expected limit validation error, got isError=%v msg=%q", res.IsError, res.ForLLM)
	}
}

func TestBinanceFuturesOrderBookTool_Validation(t *testing.T) {
	tool := NewBinanceFuturesOrderBookTool()

	res := tool.Execute(context.Background(), map[string]any{})
	if !res.IsError || !strings.Contains(res.ForLLM, "symbol is required") {
		t.Fatalf("expected symbol validation error, got isError=%v msg=%q", res.IsError, res.ForLLM)
	}

	res = tool.Execute(context.Background(), map[string]any{
		"symbol": "BTCUSDT",
		"limit":  "abc",
	})
	if !res.IsError || !strings.Contains(res.ForLLM, "limit must be an integer") {
		t.Fatalf("expected limit validation error, got isError=%v msg=%q", res.IsError, res.ForLLM)
	}
}

func TestBinanceListFuturesVolumeTool_Validation(t *testing.T) {
	tool := NewBinanceListFuturesVolumeTool()

	res := tool.Execute(context.Background(), map[string]any{
		"top": "abc",
	})
	if !res.IsError || !strings.Contains(res.ForLLM, "top must be an integer") {
		t.Fatalf("expected top validation error, got isError=%v msg=%q", res.IsError, res.ForLLM)
	}
}
