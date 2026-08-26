# Rummy Client — Phaser 3 (Handmade Hero)

Modern, maintainable Phaser 3 web client for Romanian Tile Rummy (`106` tiles, `2` jokers, `50`-point opening, `anticlockwise`, `Version 1`).

**Stack:** `phaser@3.80` + `vite@5` + `typescript@5` + `@heroiclabs/nakama-js@2.9` + `eslint` + `prettier` + `vitest` + `playwright@1.62` — same `Handmade Hero` incremental style as backend: one small vertical slice per day, `npm run dev` green, `npx playwright test` green, one focused commit.

**Backend:** Go `1.23.5` `rummy-mvp-rc1` at `127.0.0.1:7350` `defaultkey` via `docker compose up --build -d` + `make smoke` `SMOKE PASSED`.

**Docs:** `client/docs/roadmap.md` is the source of truth — 115 days from empty `client/` to `client-mvp-rc1` with Playwright after each slice.

## Quick Start (Day 1)

```bash
cd client
npm install
npm run dev # http://127.0.0.1:5173 — blank Phaser canvas
npx playwright test # Day 1 scaffold
```

See `client/docs/roadmap.md` Phase 1 for `VITE_NAKAMA_*` env and `docker compose up --build -d` prerequisites.
