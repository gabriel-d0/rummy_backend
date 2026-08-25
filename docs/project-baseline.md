# Project Baseline — Phase 1 Day 1 Audit

**Date:** 2026-08-25 (amended 2026-08-25 — Go migration, see §13)  
**Phase / Day:** Phase 1 — Day 1: Repository and environment audit  
**Author:** Lead gameplay/server engineer  
**Branch:** `main` (no commits yet at audit start)  
**Repo:** `git@github.com:gabriel-d0/rummy_backend.git`

This document records the factual starting state of the repository as required by `AGENTS.md:226` (Day 1 — Repository, Nakama, Docker, and developer baseline), expanded to the granular roadmap Phase 1 Day 1 deliverable: *Review existing project, language, tooling, branching, CI, and deployment assumptions. Record findings in `docs/project-baseline.md`.*

> **Amendment 2026-08-25 — Language migrated TypeScript → Go:** After Day 3 the project was rewritten from TypeScript to Go per explicit user request (`refactor: migrate Nakama runtime from TypeScript to Go` `55c7f3b`). `AGENTS.md:124` originally mandated TypeScript; this document now records Go as the authoritative runtime language (see §9 and new §13). All future `internal/*` packages, `main.go`, `Dockerfile` plugin build, and `go test`/`go vet` replace the previous `src/main.ts`/`tsconfig` path. Historical TS decisions remain for audit trail.

---

## 1. Repository Inventory

| Item | Finding |
|---|---|
| **Working directory** | `/Users/gabriel/work/rummy/rummy_backend` (`AGENTS.md:1-3`) |
| **Git status (2026-08-25)** | No commits yet on `main`. `git log --oneline -10` → `fatal: your current branch 'main' does not have any commits yet`. |
| **Tracked files** | None (pre-audit). Only untracked file is `AGENTS.md` (80 470 bytes, ~887 lines including embedded roadmap). |
| **Untracked / untracked dirs** | `docs/` freshly created for this audit. No other top-level entries. |
| **`.gitignore`** | **Does not exist** (`cat .gitignore` → `No such file or directory`). |
| **`.env`, `.env.example`** | Does not exist. |
| **`/src`, `/tests`, `/docker`** | Do not exist. Architecture skeleton from `AGENTS.md:133-159` not yet materialized. |
| **`/docs`** | Only this file after audit. No `docs/rules-decisions.md`, `docs/protocol.md`, etc. yet. Expected per roadmap Phase 1. |
| **README** | Does not exist. `AGENTS.md:238-240` requires one but none present. |
| **License** | None observed. |
| **Git remote** | `origin  git@github.com:gabriel-d0/rummy_backend.git (fetch/push)` — SSH remote configured, `branch.main.remote = origin`, `branch.main.merge = refs/heads/main`. No other remotes. |
| **Branches** | Only local `main` (unborn). No remote branches fetched (`git branch -a` empty beyond HEAD pointer). This is a brand-new repo. |
| **Working tree cleanliness** | Clean aside from untracked `AGENTS.md` (+ new `docs/` after this audit). |

**Evidence snapshot (truncated `ls -la`, 2026-08-25):**
```
total 160
drwxr-xr-x   5 gabriel  staff   160 Aug 25 16:23 .
drwxr-xr-x   3 gabriel  staff    96 Aug 25 16:12 ..
drwxr-xr-x  10 gabriel  staff   320 Aug 25 16:17 .git
-rw-r--r--@  1 gabriel  staff 80470 Aug 25 16:15 AGENTS.md
drwxr-xr-x@  2 gabriel  staff    64 Aug 25 16:23 docs
```

**Parent directory:** `/Users/gabriel/work/rummy/` contains only `rummy_backend/` (no sibling frontend, no monorepo). `/Users/gabriel/work/` also contains `tinybot/` (unrelated app using `postgres:15-alpine` + `app-api` on `8080` — confirms Docker is functional but not related to this project).

---

## 2. Language and Tooling Assessment

### 2.1 Authoritative constraint from AGENTS.md

`AGENTS.md:124-126` mandates:

> TypeScript for Nakama runtime code unless the repository already has a clearly established supported language; if another language is already configured, preserve that choice and explain it.

**Finding (2026-08-25 audit):** No language is established. No `package.json`, `tsconfig.json`, `go.mod`, `Dockerfile`, `nakama/` dir, `Makefile`, or runtime config exists. Therefore **TypeScript is the selected language** for all Nakama runtime work going forward. This decision is recorded as the Day 1 language choice.

For tests, `AGENTS.md:127` requires a test runner appropriate for the language. For TypeScript the expected runner is **Jest or Vitest** (to be decided Day 3/4 when `package.json` is introduced). No runner exists yet.

**Finding (amended 2026-08-25):** User requested migration to Go and Day 3 was rewritten (`55c7f3b`). **Go is now the authoritative runtime language** (`go 1.23.5`, `nakama-common v1.36.0`, `main.go` with `InitModule` at `main.go:22`). The test runner is `go test`/`go vet`; formatting is `gofmt`/`go fmt`. Historical TS decision kept above for traceability — see §13.

### 2.2 Local toolchain (probed 2026-08-25, `darwin 25.6.0 arm64`)

| Tool | Version / Path | Status | Relevance |
|---|---|---|---|
| **Node** | `v26.7.0` (`/usr/local/bin/node`) | ✅ Very recent (node 26). Will pin project to `^20` or `^22 LTS` in `package.json` for stability; 26 is bleeding-edge and not LTS. | Runtime for TypeScript build, Nakama JS runtime bundling, test runner. |
| **npm** | `11.19.0` (`/usr/local/bin/npm`) | ✅ | Package manager. Will add `package-lock.json`. |
| **npx** | available (`/usr/local/bin/npx`) | ✅ | |
| **Go** | `go1.27.0 darwin/arm64` (`/opt/homebrew/bin/go`) | ✅ installed, not required for TS path but available. Nakama itself is Go; useful for tooling/plugins if ever needed. | Not primary language. No `go.mod` in repo, so Go is *not* the established language. |
| **Docker** | `Docker version 29.7.2, build a7dcaa6` (`/usr/local/bin/docker`) | ✅ | Required for `AGENTS.md` Docker Compose workflow. |
| **Docker Compose** | `Docker Compose version v5.4.0` | ✅ (also available as `docker compose` plugin) | `docker compose up --build` command expected from Day 1 acceptance criteria. |
| **Docker daemon** | Running via Docker Desktop (`desktop-linux` context, socket `unix:///Users/gabriel/.docker/run/docker.sock` → `/var/run/docker.sock`) | ✅ Verified with `docker ps` showing live containers `tinybot-api` + `tinybot-db`. | Confirms local Docker workflow is operational. |
| **psql** | Not found (`which psql` → empty) | ⚠️ Not installed locally. Not blocking — Postgres will run in Docker. For debugging, `docker compose exec` or `docker exec` will be used. Could `brew install libpq` / `postgresql@15` later for convenience. | |
| **Git** | System git (via Xcode CLT), `main` unborn | ✅ | |
| **Formatting / lint / type-check** | None configured (no `prettier`, `eslint`, `tsconfig.json`) | ❌ Gap to close Day 2–4 | `AGENTS.md:242,333` expects `format / lint / typecheck / test` scripts. **Amended 2026-08-25:** For Go, use `gofmt`/`go fmt`, `go vet ./...`, `go test ./...` (see `README.md` Dev commands). |
| **CI** | None (no `.github/workflows/`) | ❌ | Roadmap Day 6. |

### 2.3 Implicit assumptions confirmed

- No frontend framework present (`AGENTS.md:129` — do not add one in first phase; confirmed: none to remove).
- No Nakama installation locally or in repo — to be added via Docker Compose `heroiclabs/nakama` image.
- Tile set rules, meld validation, deck, etc. not yet implemented — correctly deferred per Handmade Hero incremental plan.

---

## 3. Branching Strategy

| Aspect | Current State | Recommendation |
|---|---|---|
| **Default branch** | `main` (unborn, tracks `origin/main`) | Keep as default. Roadmap's daily commits push to active branch; for MVP a single `main` with small daily commits is sufficient. Introduce `feature/*` branches only if parallel work is needed. |
| **Existing branches** | None local, none fetched. `git branch -a` empty. | No migration needed. |
| **Commit history** | Zero commits at audit start. First commit will be this baseline (`docs: add Phase 1 Day 1 project baseline audit`). | Enforce *one focused commit per day* per `AGENTS.md:13,211-219` and roadmap Daily Definition of Done §5. Commit style: `chore: / docs: / feat: / test: / fix: / refactor:` — never `update`/`wip`. |
| **Remote sync** | `origin` SSH reachable (not probed with network push yet; will verify on first `git push`). | Ensure `git push -u origin main` on first commit to establish upstream. |
| **Protection / PR workflow** | No `.github` branch protection observed. | Defer until CI exists (Day 6). Keep pushes direct to `main` for daily slices until team grows. |
| **Tag / release** | None. Roadmap Day 135 anticipates a release-candidate tag. | Not applicable yet. |

---

## 4. CI / Automation

| Check | Present? | Notes |
|---|---|---|
| `.github/workflows/*` | ❌ No | Roadmap Days 6 & 135. To add: format, lint, typecheck, unit tests, build verification, Docker smoke. |
| Pre-commit hooks (`.husky/`, `lefthook`, etc.) | ❌ No | Optional; not required by `AGENTS.md`. Could add later after formatter choice. |
| `Makefile` / `justfile` | ❌ No | `AGENTS.md` Day 4 expects developer scripts (start/stop/logs/clean/build/test/lint/typecheck/format). To be added as `package.json` scripts + optionally `Makefile`. |
| Dependabot / Renovate | ❌ No | Recommend enabling Dependabot for `npm` + `docker` after `package.json` lands. |
| Code coverage gate | ❌ No | Consider after rule modules stabilize. |

**Conclusion:** CI is a deliberate Day 6 deliverable, not a Day 1 gap to fix now. This audit simply records its absence.

---

## 5. Deployment Assumptions

| Area | Current | Assumption / Decision for MVP |
|---|---|---|
| **Nakama version** | Not yet pinned. No `docker-compose.yml` to inspect. | Will pin to latest stable `heroiclabs/nakama:3.24.x` (or newest `3.26.x` if verified) with matching `cockroachdb/cockroach` or `postgres` image as required by that Nakama release. `AGENTS.md:123` requires a database dependency matching the Nakama version — decision to be recorded in `docker-compose.yml` comments Day 2. |
| **Database** | No container yet for this project. Local Docker has `tinybot-db` (`postgres:15-alpine`) for a *different* project — proves `postgres:15` works on this host but does not imply choice for Rummy. | Nakama historically supports CockroachDB and Postgres; recent Nakama docs favor Postgres for local dev. Will verify against chosen Nakama image's `migrate` command Day 2 and document in `README.md` + `docker-compose.yml`. |
| **Hosting / production** | None defined. Local-only per `AGENTS.md:122,129`. | Local Docker Compose = sole deployment for Phases 1–15. No cloud (AWS/GCP/Fly) provisioning until Phase 24 (Production Deployment). No secrets, TLS, or IaC yet. |
| **Domains / TLS** | None. | Not needed for MVP. |
| **Backups / migrations** | None. | Phase 24 scope. |
| **Port assumptions** | Unknown until `docker-compose.yml` lands. TinyBot uses `8080` (`app-api`) + `5432` (`postgres`) mapped to `0.0.0.0`. To avoid collision, Rummy's Nakama may use `7350`/`7351` (Nakama default) + `5432` offset or distinct project network. Will document Day 2. | `AGENTS.md:238-240` promises console URL + start/stop/logs instructions — to be written Day 5. |
| **Secrets** | None. No `.env` committed. | Will add `.env.example` with non-secret local defaults Day 1/2 per `AGENTS.md:235`. Never commit real secrets; use managed store only in Phase 24. |
| **Scaling** | Single Nakama instance + single DB in Docker Compose. | Correct for authoritative match loop at MVP scale (2–4 players per match). Horizontal scaling deferred. |

---

## 6. Docker / Local Dev Baseline Check

- Docker Desktop running, context `desktop-linux`, socket at `~/.docker/run/docker.sock` (also symlinked to `/var/run/docker.sock`).
- `docker ps` 2026-08-25 shows 2 healthy containers from unrelated `tinybot` project, proving daemon, networking, and volume mounts work.
- `docker compose version v5.4.0` exists.
- No `docker-compose.yml` / `compose.yml` / `Dockerfile` in `rummy_backend` yet — expected; Day 2 deliverable.
- No `.dockerignore` yet.
- Resource headroom: not measured, but typical 4-player tile state is trivial (<100 KB per match snapshot) so Docker defaults are sufficient. Resource limits can be added later if load tests (Phase 15) warrant.

**Smoke readiness:** Host can run `docker compose up --build` once Day 2 lands; today's check merely confirms prerequisites are met (`AGENTS.md:244-247` acceptance criteria *not yet* expected to pass).

---

## 7. Documentation Baseline

| Doc | Expected per roadmap | Exists now? |
|---|---|---|
| `docs/project-baseline.md` | Phase 1 Day 1 (this file) | ✅ Created this audit |
| `docs/rules-decisions.md` | Phase 1 Day 8 / AGENTS Day 2 | ❌ |
| `docs/architecture.md` | Phase 14 Day 121 | ❌ |
| `docs/state-machine.md` | Phase 14 Day 122 / AGENTS Day 22 | ❌ |
| `docs/protocol.md` | Phase 14 Day 123 / AGENTS Day 22 | ❌ |
| `docs/testing.md` | Phase 14 Day 125 / AGENTS Day 22 | ❌ |
| `README.md` | AGENTS Day 1 | ❌ |
| `AGENTS.md` | Input constraint (source of truth) | ✅ Present, not yet committed |

---

## 8. Gaps and Risks Identified

| # | Gap / Risk | Severity | Mitigation / Owner |
|---|---|---|---|
| 1 | **No `.gitignore`** — risks committing `node_modules/`, `.env`, `dist/`, Nakama data dirs | Medium | Add Day 2 alongside `package.json`/`docker-compose.yml`. Include `node_modules/`, `dist/`, `build/`, `.env`, `.DS_Store`, `data/`, `*.log`. |
| 2 | **No `README.md`** — new developer cannot bootstrap | Medium | Day 5 deliverable. Template to cover prerequisites, `docker compose up --build`, console URL, logs, test commands. |
| 3 | **No formatter/linter/type-check** | Low now, Medium once code lands | Day 4. Recommend: `prettier` + `eslint` (flat config, `typescript-eslint`) + `tsc --noEmit`. Add `npm run format / lint / typecheck / test`. |
| 4 | **`psql` absent locally** | Low | Use `docker exec -it <db> psql -U nakama` or `docker compose exec` for DB inspection. Optionally `brew install libpq` without linking postgres server. |
| 5 | **Node 26 non-LTS** | Low | Pin project engine to `>=20 <27` and test on 20/22 LTS in CI; local 26 is okay for dev but CI should matrix LTS. |
| 6 | **Port collision risk** with `tinybot` if same `5432`/`8080` chosen | Low | Use distinct compose project name (`rummy_backend`) and Nakama defaults `7350`/`7351` + isolated network; document ports in README. |
| 7 | **No CI** | Low now | Day 6. GitHub Actions with `actions/setup-node`, `docker/build-push-action` smoke, caching. |
| 8 | **Empty commit history** | Info | First commit will establish baseline; ensure `git config user.name/email` set before commit (check via `git config --list`). |
| 9 | **Nakama version not yet pinned** | Medium | Day 2 decision must be tested with `docker compose up` against selected `heroiclabs/nakama` + `heroiclabs/nakama-migrate` behavior. Record in `docker-compose.yml` comments + `docs/project-baseline.md` addendum. |

---

## 9. Decisions Made This Day

1. **Language: TypeScript** — per `AGENTS.md:125`, no existing language, so TS is authoritative for Nakama JS runtime modules. All future `src/match/*`, `src/rules/*`, `src/setup/*`, `src/protocol/*` will be TS. **Amended 2026-08-25:** Migrated to Go per user request; now `main.go` + `internal/*` with `go 1.23.5` is authoritative (see §13).
2. **Test runner: TBD (Vitest preferred, Jest fallback)** — to be decided when `package.json` is created Day 2/3. Criteria: fast ESM-friendly, deterministic seed support for shuffle tests, no Nakama runtime mock needed for pure rule modules. **Amended:** For Go, `go test` with deterministic seedable deck/shuffle helpers.
3. **Docker project name: `rummy_backend`** — to be set via `name:` in `compose.yml` or `COMPOSE_PROJECT_NAME` to isolate from `tinybot`. Prevents container/volume name collision.
4. **Branch strategy: trunk-based `main`, daily small commits** — no `develop` branch until needed. Commit style enforced per `AGENTS.md:212-219`.
5. **Out-of-scope for Day 1:** No `docker-compose.yml`, no Nakama runtime skeleton, no `.env.example`, no `README.md`, no game behavior. This audit is intentionally read-only + single-doc deliverable to honor *Handmade Hero* slice size.

---

## 10. Immediate Next Steps (Phase 1 Day 2 — Docker Compose)

Per roadmap Phase 1 table, Day 2 is **Docker Compose**:

- [ ] Add `compose.yml` (or `docker-compose.yml`) with `nakama` + `postgres`/`cockroachdb` pinned versions. Include healthchecks, named volumes, project isolation, and local `nakama.yml` mount.
- [ ] Add `nakama/data/` or `nakama/` local config with minimal `config.yml` for dev (`name: rummy_backend`, `http_key`, console credentials).
- [ ] Add `.env.example` with non-secret defaults (`POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `NAKAMA_PORT`, `NAKAMA_CONSOLE_PORT`).
- [ ] Add `.gitignore` and `.dockerignore`.
- [ ] Add `README.md` skeleton or at least `docs/` note for `docker compose up --build` verification.
- [ ] Verify `docker compose up --build` starts both DB and Nakama, DB migrates, console reachable at `http://localhost:7351` (or chosen port), runtime loads.
- [ ] `docker compose logs nakama --tail=100` shows healthy startup without crash.
- **Blocked on:** None — ready to proceed. Only potential blocker is Nakama version selection needing a quick test pull (`docker pull heroiclabs/nakama:<version>`).

---

## 11. Verification Commands Executed (evidence)

```bash
git log --oneline -10               # fatal: no commits yet
git status                           # On branch main, No commits yet, untracked: AGENTS.md
git remote -v                        # origin git@github.com:gabriel-d0/rummy_backend.git
git branch -a                        # (empty)
ls -la                               # only AGENTS.md + .git
cat .gitignore                       # No such file or directory
which docker && docker --version     # Docker version 29.7.2
docker compose version               # v5.4.0
which node && node --version         # v26.7.0
which npm && npm --version           # 11.19.0
which go && go version               # go1.27.0 darwin/arm64
which psql                           # not found
docker ps                            # 2 containers (tinybot) — daemon healthy
ls -la /Users/gabriel/.docker/run/docker.sock
docker context ls
```

All outputs captured 2026-08-25 before writing this file. Raw logs available in shell history for the audit session.

---

## 12. Acceptance Checklist — Phase 1 Day 1

- [x] Repository inspected (files, git, remotes, branches).
- [x] Language/tooling identified (empty → TypeScript chosen per AGENTS.md, Node 26 / npm 11 / Docker 29.7.2 / Compose v5.4.0 verified).
- [x] Branching model reviewed (new repo, `main` unborn, SSH origin).
- [x] CI/deployment assumptions reviewed (no CI yet — Day 6; no deployment — local Docker only until Phase 24).
- [x] Findings recorded in `docs/project-baseline.md` (this file) with decisions, gaps, and next-day plan.
- [x] No game behavior implemented (audit is read-only; respects Handmade Hero slice).
- [ ] Peer review / push — to be done immediately after this file is written (commit + push).

---

*End of Phase 1 Day 1 audit. Next is Phase 1 Day 2: Docker Compose for Nakama + PostgreSQL local development.*

---

## 13. Amendment 2026-08-25 — Migration from TypeScript to Go

**Trigger:** User directive “let change typescript to go, rewrite the entire backend” + confirmed migration-commit strategy.

**Commits affected:**
- `9354b59 feat: add Nakama TypeScript runtime skeleton with InitModule` (Day 3 TS, now superseded but kept in history)
- `55c7f3b refactor: migrate Nakama runtime from TypeScript to Go` (replaces TS with Go — `package.json`/`tsconfig.json`/`src/main.ts`/`build/` removed, `go.mod`/`go.sum`/`main.go`/`Dockerfile` added)

**What changed:**
- Runtime language: TypeScript (`nakama-runtime` JS, `tsc` → `build/index.js` → `nakama/data/modules/index.js` via volume mount) → **Go** (`go 1.23.5`, `nakama-common v1.36.0`, `main.go:22` `InitModule`, `Dockerfile` multi-stage `heroiclabs/nakama-pluginbuilder:3.26.0` → `backend.so` baked as `/nakama/data/modules/rummy_backend.so` in `rummy_backend:local` image).
- Build: `npm run build` (`tsc`) → `go vet ./...` + `docker compose build` (plugin `CGO_ENABLED=1` `--buildmode=plugin`, 9.5s) then `docker compose up -d`. Volume changed from `./nakama/data:/nakama/data` to `./nakama/data/local.yml:ro` so baked `.so` is not hidden (TS volume hid image; Go requires baked).
- Gitignore: added Go `vendor/*.so/backend.so`, kept legacy Node entries annotated.
- Documentation: `AGENTS.md:124` now satisfied via Go (`go.mod` exists), `main.go:12` is the established language source. Future `internal/rules`, `internal/match`, etc., will be Go packages.

**Why the version pin:** Go plugin must match Nakama’s Go toolchain (`runtime:go1.23.5` in `docker compose logs`). Initial `protobuf v1.36.6` caused `plugin was built with a different version of package google.golang.org/protobuf/internal/pragma` crash; pinned to `v1.36.4` (same as `nakama:v1.36.0` requires) resolves it. `go 1.23.5` in `go.mod` aligns with builder.

**Verification after migration (copied from commit `55c7f3b`):**
- `go vet ./...` ok, `docker compose build nakama` 9.5s, `docker compose up -d` both healthy
- Logs: `Rummy backend InitModule Go Day 3 skeleton starting` (`main.go:23`), `Found runtime modules count 1 [rummy_backend.so]`, `Registered Go RPC health/version`
- RPC: `POST /v2/account/authenticate/device` via `defaultkey` → JWT, `POST /v2/rpc/health` → `{"status":"ok","version":"0.1.0-go-day3-skeleton"}` (Go), `POST /v2/rpc/version` → `runtime:go` — console `200` at `http://127.0.0.1:7351`.

**Impact on roadmap:**
- Phase 1 Days 1–2 unchanged (Docker Compose still works, now builds Go image).
- Day 3 deliverable still “Nakama runtime skeleton” but in Go.
- Future days (rules, deck, match) will be Go pure packages under `internal/` instead of `src/`; test runner will be `go test`.

**Preserved:** History of TS implementation remains in git for audit; this amendment keeps `docs/project-baseline.md` as the single source of truth for language decisions per `AGENTS.md:125`.
