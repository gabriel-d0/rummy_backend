# Rummy Client — SvelteKit

SvelteKit + Tailwind + Nakama JS client for Romanian Tile Rummy (106 tiles, 2 jokers, 50-point opening, Version 1). Day 5 — Quick Start with `VITE_NAKAMA_*` env.

## Quick Start

```bash
# 1. Backend up (from repo root)
docker compose up --build -d && make smoke

# 2. Client (from client/)
cd client
npm install
cp .env.example .env # optional — defaults are 127.0.0.1:7350 defaultkey
npm run dev -- --open # http://localhost:5173

# 3. Check
npm run check # svelte-check
npm run lint # prettier --check + eslint
npm run test:unit -- --run # vitest
npm run test:e2e # playwright
npm run build # vite build
```

## Env

`VITE_NAKAMA_HOST` `VITE_NAKAMA_PORT` `VITE_NAKAMA_KEY` `VITE_NAKAMA_USE_SSL` — see `.env.example`. `VITE_NAKAMA_KEY` must be `defaultkey` for local Nakama (`defaultkey:` not `defaulthttpkey:`).

## Docs

- `client/docs/roadmap.md` — 95 days Day 1-95, no table scroll, 52px vs 64px, Playwright after each slice
- `../docs/protocol.md` — opcodes 1..9/100..199, Version 1, envelope
- `../docs/state-machine.md` — GamePhase/TurnPhase AllowedOps

```


## Next Steps (Day 5)

Per `client/docs/roadmap.md` Phase 1 Day 5 done — Quick Start now covers `npm i`, `npm run dev`, `VITE_NAKAMA_*`, `docker compose up --build -d` prerequisites.
Next is **Day 6 — Smoke test** (`client/docs/roadmap.md:6`): `tests/smoke.spec.ts` verifies `npm run dev` serves `200` at `5173`.
```
