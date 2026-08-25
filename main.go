// Rummy Backend — Nakama Go runtime entrypoint
// Phase 1 Day 3 (rewritten): Minimal InitModule skeleton in Go.
// This replaces the previous TypeScript runtime per user request to migrate
// to Go. It intentionally contains no game rules yet — only proves the
// Go plugin toolchain (go build --buildmode=plugin → backend.so → Nakama) works.
// All gameplay logic will live under internal packages per AGENTS.md:133-159.
// See docs/project-baseline.md for language decision amendment.

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/heroiclabs/nakama-common/runtime"
)

// InitModule is the entrypoint Nakama calls once at startup.
// Must be exported and match runtime.InitModule signature.
func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	logger.Info("Rummy backend InitModule — Romanian Tile Rummy v0.1.0 Go Day 3 skeleton starting")
	logger.Info("Context env: %v", ctx.Value(runtime.RUNTIME_CTX_ENV)) // env may be nil in some Nakama versions

	// Register trivial RPCs so we can prove the runtime is loaded via API explorer or curl.
	// Client can call: POST /v2/rpc/health or /v2/rpc/version with session token.
	if err := initializer.RegisterRpc("health", rpcHealth); err != nil {
		logger.Error("Failed to register health RPC: %v", err)
		return err
	}
	if err := initializer.RegisterRpc("version", rpcVersion); err != nil {
		logger.Error("Failed to register version RPC: %v", err)
		return err
	}

	logger.Info("Rummy backend RPCs registered: health, version")
	logger.Info("Rummy backend skeletal Go module loaded — no match handlers yet (Day 5)")

	// No match handlers registered yet — Day 5 will add authoritative match.

	return nil
}

// rpcHealth returns JSON with status, timestamp and caller echo. Payload is ignored.
func rpcHealth(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	userID, _ := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	username, _ := ctx.Value(runtime.RUNTIME_CTX_USERNAME).(string)

	logger.Info("rpcHealth called by userId=%v username=%v payload=%v", userID, username, payload)

	resp := map[string]interface{}{
		"status":    "ok",
		"service":   "rummy_backend",
		"version":   "0.1.0-go-day3-skeleton",
		"timestamp": time.Now().UnixMilli(),
		"caller": map[string]interface{}{
			"userId":   userID,
			"username": username,
		},
	}
	// If anonymous (healthcheck via runonce), caller may be empty — keep null-like
	if userID == "" {
		resp["caller"] = nil
	}

	b, err := json.Marshal(resp)
	if err != nil {
		logger.Error("rpcHealth marshal failed: %v", err)
		return "", runtime.NewError("marshal failed", 13)
	}
	return string(b), nil
}

func rpcVersion(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, _payload string) (string, error) {
	userID, _ := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	logger.Debug("rpcVersion called by %v", userID)

	resp := map[string]interface{}{
		"name":    "rummy_backend",
		"version": "0.1.0-go-day3-skeleton",
		"nakama":  "3.26.0",
		"runtime": "go",
		"phase":   "1-day3-go",
		"docs":    "docs/project-baseline.md",
		"note":    "migrated from TypeScript per user request 2026-08-25",
	}

	b, err := json.Marshal(resp)
	if err != nil {
		return "", runtime.NewError("marshal failed", 13)
	}
	return string(b), nil
}
