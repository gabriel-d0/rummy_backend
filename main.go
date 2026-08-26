// Rummy Backend — Nakama Go runtime entrypoint
// Phase 1 Day 3 (rewritten) + Phase 4 Day 22: Minimal InitModule skeleton in Go
// with authoritative match registration.
// This replaces the previous TypeScript runtime per user request to migrate
// to Go. It intentionally contains minimal game rules yet — only proves the
// Go plugin toolchain (go build --buildmode=plugin → backend.so → Nakama) works.
// All gameplay logic will live under internal packages per AGENTS.md:133-159.
// See docs/project-baseline.md for language decision amendment.

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/gabriel-d0/rummy_backend/internal/match"
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

	// RPC for authoritative match creation: client calls rpc("create_match") -> nk.MatchCreate("rummy")
	if err := initializer.RegisterRpc("create_match", rpcCreateMatch); err != nil {
		logger.Error("Failed to register create_match RPC: %v", err)
		return err
	}
	logger.Info("Rummy backend RPC registered: create_match")

	// Register authoritative match handler (Day 22). The match name "rummy"
	// is the stable ID clients use via nk.matchCreate / match join.
	if err := initializer.RegisterMatch("rummy", match.NewRummyMatch); err != nil {
		logger.Error("Failed to register rummy match: %v", err)
		return err
	}
	logger.Info("Rummy backend match registered: rummy")
	logger.Info("Rummy backend skeletal Go module loaded — Day 22 match skeleton registered (full lobby Day 23)")

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

func rpcCreateMatch(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, _payload string) (string, error) {
	userID, _ := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	logger.Info("rpcCreateMatch called by %v", userID)
	matchId, err := nk.MatchCreate(ctx, "rummy", map[string]interface{}{})
	if err != nil {
		logger.Error("MatchCreate rummy failed: %v", err)
		return "", runtime.NewError("match create failed", 13)
	}
	logger.Info("Created rummy match %s for %s", matchId, userID)
	b, _ := json.Marshal(map[string]string{"matchId": matchId})
	return string(b), nil
}
