// Rummy Backend — Nakama TypeScript runtime entrypoint
// Phase 1 Day 3: Minimal InitModule skeleton.
// This file intentionally does NOT contain game rules yet; it only proves
// the toolchain (TypeScript → JS bundle → Nakama modules) works end-to-end.
// All gameplay logic will live under src/rules, src/match, src/setup per AGENTS.md:133-159.

const RPC_HEALTH_ID = "health";
const RPC_VERSION_ID = "version";

// InitModule is the global entrypoint Nakama calls once per startup.
// Signature must match nkruntime.InitModule — do not rename or export as ES module.
function InitModule(
  ctx: nkruntime.Context,
  logger: nkruntime.Logger,
  nk: nkruntime.Nakama,
  initializer: nkruntime.Initializer
): void {
  logger.info("Rummy backend InitModule — Romanian Tile Rummy v0.1.0 Day 3 skeleton starting");
  logger.info("Context env: %v executionMode: %v", ctx.env, ctx.executionMode);

  // Register a trivial RPC so we can prove the runtime is loaded via console API explorer or curl.
  // Client can call: POST /v2/rpc/health or /v2/rpc/version with session token.
  initializer.registerRpc(RPC_HEALTH_ID, rpcHealth);
  initializer.registerRpc(RPC_VERSION_ID, rpcVersion);

  logger.info("Rummy backend RPCs registered: %v, %v", RPC_HEALTH_ID, RPC_VERSION_ID);
  logger.info("Rummy backend skeletal module loaded — no match handlers yet (Day 5)");
}

// Simple health RPC: returns JSON with status and timestamp. No auth needed beyond valid session.
let rpcHealth: nkruntime.RpcFunction = (
  ctx: nkruntime.Context,
  logger: nkruntime.Logger,
  nk: nkruntime.Nakama,
  payload: string
): string => {
  logger.info("rpcHealth called by userId=%v username=%v payload=%v", ctx.userId, ctx.username, payload);
  const response = {
    status: "ok",
    service: "rummy_backend",
    version: "0.1.0-day3-skeleton",
    timestamp: Date.now(),
    // Echo presence to prove ctx wiring without leaking secrets
    caller: ctx.userId ? { userId: ctx.userId, username: ctx.username } : null,
  };
  return JSON.stringify(response);
};

let rpcVersion: nkruntime.RpcFunction = (
  ctx: nkruntime.Context,
  logger: nkruntime.Logger,
  nk: nkruntime.Nakama,
  _payload: string
): string => {
  logger.debug("rpcVersion called by %v", ctx.userId);
  return JSON.stringify({
    name: "rummy_backend",
    version: "0.1.0-day3-skeleton",
    nakama: "3.26.0",
    runtime: "typescript",
    phase: "1-day3",
    docs: "docs/project-baseline.md",
  });
};
