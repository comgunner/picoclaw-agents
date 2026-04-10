// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"

	"github.com/comgunner/picoclaw/pkg/logger"
)

// StdioTransport communicates with an MCP subprocess via stdin/stdout (JSON-RPC 2.0).
type StdioTransport struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Reader
	stderr io.ReadCloser
	mu     sync.Mutex
	id     int64 // simple int64 is fine under mutex
	done   chan struct{}
	cancel context.CancelFunc // signals stderr goroutine to stop
}

// NewStdioTransport spawns a subprocess for MCP communication.
func NewStdioTransport(command string, args []string, env map[string]string) (*StdioTransport, error) {
	cmd := exec.Command(command, args...)

	// Inherit current env + override with config env vars
	cmdEnv := os.Environ()
	for k, v := range env {
		cmdEnv = append(cmdEnv, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Env = cmdEnv

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("stdin pipe: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start %s: %w", command, err)
	}

	// Cancelable context for stderr goroutine
	ctx, cancel := context.WithCancel(context.Background())

	t := &StdioTransport{
		cmd:    cmd,
		stdin:  stdin,
		stdout: bufio.NewReader(stdout),
		stderr: stderr,
		done:   make(chan struct{}),
		cancel: cancel,
	}

	// Capture stderr in background for debugging
	go t.captureStderr(ctx)

	return t, nil
}

// Call sends a JSON-RPC request. Respects context cancellation.
// CRITICAL: Uses goroutine + select on ctx.Done() to prevent blocking reads.
func (t *StdioTransport) Call(ctx context.Context, method string, params map[string]any) (*json.RawMessage, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	atomic.AddInt64(&t.id, 1)
	req := jsonRPCRequest{
		JSONRPC: "2.0",
		ID:      t.id,
		Method:  method,
		Params:  params,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}
	data = append(data, '\n')

	// Use a channel to receive the read result, so we can select on ctx.Done()
	type readResult struct {
		line string
		err  error
	}
	readCh := make(chan readResult, 1)
	go func() {
		line, err := t.stdout.ReadString('\n')
		readCh <- readResult{line, err}
	}()

	// Write request (not cancellable — already sent to subprocess)
	if _, err := t.stdin.Write(data); err != nil {
		return nil, fmt.Errorf("write stdin: %w", err)
	}

	// Wait for response OR context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case result := <-readCh:
		if result.err != nil {
			return nil, fmt.Errorf("read stdout: %w", result.err)
		}

		line := result.line
		if len(line) > MAX_LINE_BYTES {
			return nil, fmt.Errorf("response line exceeds MAX_LINE_BYTES (%d)", MAX_LINE_BYTES)
		}

		var resp jsonRPCResponse
		if err := json.Unmarshal([]byte(line), &resp); err != nil {
			return nil, fmt.Errorf("parse response: %w", err)
		}
		if resp.Error != nil {
			return nil, fmt.Errorf("MCP error %d: %s", resp.Error.Code, resp.Error.Message)
		}
		return &resp.Result, nil
	}
}

// Close terminates the subprocess and waits for cleanup.
func (t *StdioTransport) Close() error {
	t.cancel() // signal stderr goroutine to stop
	t.stdin.Close()
	<-t.done // wait for stderr goroutine
	return t.cmd.Wait()
}

func (t *StdioTransport) captureStderr(ctx context.Context) {
	defer close(t.done)
	scanner := bufio.NewScanner(t.stderr)
	scanner.Buffer(make([]byte, MAX_LINE_BYTES), MAX_LINE_BYTES)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
			logger.DebugCF("mcp", "MCP stderr",
				map[string]any{"cmd": t.cmd.Path, "line": scanner.Text()})
		}
	}
}
