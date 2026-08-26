# Roadmap — Days to be Implemented

**Scope:** This file lists the *future* days that are **planned but not yet implemented** after the current MVP. For what is already done, see `docs/IMPLEMENTED.md` (Phases 1–16 complete, `rummy-mvp-rc1` at `cfce62e` / `41e794b`, 82 commits `36c2c59..41e794b`). This roadmap follows the Handmade Hero incremental style: one small vertical slice per day, `make check` + `docker compose build` green before next day.

**Current HEAD:** `main@41e794b` — all 24-day `AGENTS.md` days 1–24 done (foundation through minimal CLI `5ade046` + deterministic simulation `01f6a3c` + hardening `3e1b92a` + win `2d278a5`). Tagged `rummy-mvp-rc1`.

**How to read:** Day numbers are per the extended roadmap (≈Day 25–820) that follows the 24-day MVP. Each day is one focused commit; phases end in a runnable, tested state.

---

## Phase 17 — Future product features (still deferred per `AGENTS.md:100`)

These were intentionally deferred until the baseline round is playable and tested. They are the *first* days to be implemented post-MVP if product approves.

| Day range | Focus | Deliverables | Status |
|---:|---|---|---|
| 147–155 | Multi-round match scoring | Dealer rotation, accumulated scores, bonuses, match winner, round history per player. | ⏳ Planned — not started |
| 156–165 | Social features | Table chat, player profiles, friends, private invitations, presence. | ⏳ Planned |
| 166–175 | Matchmaking and lobbies | Public tables, private rooms, skill filters, reconnect UX, rematch. | ⏳ Planned |
| 176–190 | Tournaments | Scheduled tournaments, brackets, standings, prizes (only if legally appropriate). | ⏳ Planned |
| 191–205 | Bots (deterministic) | Test bots for offline, then AI opponents (see Phases 18–23 below for full AI plan). | ⏳ Planned |
| 206–220 | Production operations | Metrics, monitoring, alerting, backups, deployment, load testing. | ⏳ Planned |
| 221–240 | Polish | Animations, sound, accessibility, mobile UX, localization. | ⏳ Planned |

---

## Phase 18 — AI Bots: Foundation and Safe Integration

Bots must use the same authoritative command protocol as humans; never give bots hidden state.

| Day range | Focus | Deliverables | Status |
|---:|---|---|---|
| 241–243 | Bot architecture decision | Document bot lifecycle, identities, difficulty levels, scheduling, how bot actions enter Nakama matches. | ⏳ Planned |
| 244–246 | Bot player model | Bot user/profile metadata, display names, avatar placeholders, bot flags, safe seat assignment. | ⏳ Planned |
| 247–250 | Bot match participation | Host can add 1–3 bots to private/local match, 2–4 total seats. | ⏳ Planned |
| 251–254 | Bot action adapter | Adapter that turns bot decisions into validated `DRAW_STOCK`/`DISCARD`/`MELD_INITIAL`/`MELD_NEW` etc. | ⏳ Planned |
| 255–258 | Bot timing | Configurable think delays, deterministic test mode (no delay), cancellation on match end. | ⏳ Planned |
| 259–262 | Bot visibility boundaries | Bot only receives own rack + public state; tests prove no opponent rack access. | ⏳ Planned |
| 263–266 | Bot error recovery | Invalid bot action → log, recover, fallback legal move without corrupting match. | ⏳ Planned |
| 267–270 | Bot observability | Structured logs: match ID, bot ID, turn, action, duration, fallback, rejection reason. | ⏳ Planned |

---

## Phase 19 — AI Bots: Legal Move Generation and Beginner Difficulty

| Day range | Focus | Deliverables | Status |
|---:|---|---|---|
| 271–274 | Legal move generator | Pure module enumerating legal actions from bot-visible state. | ⏳ Planned |
| 275–278 | Draw decision generation | Available draws: stock, previous discard, earlier pickup. | ⏳ Planned |
| 279–282 | Discard candidate generation | Every legal discard candidate after draw/meld. | ⏳ Planned |
| 283–286 | Initial meld candidate generation | Opening batches satisfying 50-point + at-least-one-run. | ⏳ Planned |
| 287–290 | New meld candidate generation | Valid sets/runs from opened rack. | ⏳ Planned |
| 291–294 | Meld extension candidate generation | Legal extensions to public melds (including others’). | ⏳ Planned |
| 295–298 | Joker replacement candidate generation | Legal joker replacements without illegal reuse. | ⏳ Planned |
| 299–302 | Beginner strategy | Heuristic: open when possible, obvious melds, extend, discard low-value. | ⏳ Planned |
| 303–306 | Beginner bot tests | Bot always legal, never leaks hidden info. | ⏳ Planned |
| 307–310 | Bot turn simulation | Deterministic games with 1–3 beginner bots, verify progress. | ⏳ Planned |

---

## Phase 20 — AI Bots: Intermediate Strategy and Difficulty Tiers

| Day range | Focus | Deliverables | Status |
|---:|---|---|---|
| 311–314 | Evaluation model | Explainable hand-evaluation: meld potential, deadwood, joker value, discard risk. | ⏳ Planned |
| 315–318 | Candidate scoring | Score all legal candidates. | ⏳ Planned |
| 319–322 | Better discard heuristic | Avoid discarding needed tiles. | ⏳ Planned |
| 323–326 | Opening strategy | Preserve vs open with 50-point batch. | ⏳ Planned |
| 327–330 | Pickup strategy | Stock vs latest vs earlier pickup scoring. | ⏳ Planned |
| 331–334 | Table exploitation | Prefer safe extensions that reduce rack or improve options. | ⏳ Planned |
| 335–338 | Joker strategy | Retain/meld/replace/expose jokers. | ⏳ Planned |
| 339–342 | Difficulty configuration | `beginner`/`normal`/`advanced` profiles with search depth/randomness/delay. | ⏳ Planned |
| 343–346 | Controlled randomness | Seeded weighted choice, reproducible. | ⏳ Planned |
| 347–350 | Strategy benchmark suite | Fixed board-state scenarios, compare difficulties. | ⏳ Planned |

---

## Phase 21 — AI Bots: Search, Simulation, and Training Infrastructure

| Day range | Focus | Deliverables | Status |
|---:|---|---|---|
| 351–354 | Search budget design | Max CPU time, candidate count, simulation count, cancellation. | ⏳ Planned |
| 355–358 | One-turn lookahead | Bounded evaluation of immediate actions. | ⏳ Planned |
| 359–362 | Two-turn approximate planning | Limited future-turn estimation without hidden info. | ⏳ Planned |
| 363–366 | Monte Carlo rollout prototype | Plausible unseen stock/opponent distributions from public info only. | ⏳ Planned |
| 367–370 | Simulation state cloning | Efficient deterministic copies for simulations. | ⏳ Planned |
| 371–374 | Search cancellation | Cancel on match end / turn loss / disconnect / timeout. | ⏳ Planned |
| 375–378 | Fallback policy | Legal fallback if simulation fails/times out. | ⏳ Planned |
| 379–382 | Offline self-play runner | Script to run thousands of bot-vs-bot rounds with seeds. | ⏳ Planned |
| 383–386 | Bot telemetry | Win rates, turn length, invalid-action count, timeout rate, etc. | ⏳ Planned |
| 387–390 | Performance tests | Bot turns meet latency/CPU budgets under concurrent load. | ⏳ Planned |

---

## Phase 22 — Matchmaking, Ratings, and Bot-Assisted Player Experience

| Day range | Focus | Deliverables | Status |
|---:|---|---|---|
| 391–394 | Match mode design | Casual, private, practice, ranked, bot-only modes. | ⏳ Planned |
| 395–398 | Practice mode | Private practice match with configurable bots. | ⏳ Planned |
| 399–402 | Bot transparency UX/protocol | Mark bot players clearly in lists, events, logs, payloads. | ⏳ Planned |
| 403–406 | Matchmaking queues | Queue metadata: player count, region, mode, difficulty. | ⏳ Planned |
| 407–410 | Queue timeout policy | When to offer bots after waiting threshold; opt-in where appropriate. | ⏳ Planned |
| 411–414 | Backfill bots | Fill empty seats with bots only before start or per explicit rules. | ⏳ Planned |
| 415–418 | Player rating model | Initial rating for human-vs-human ranked games; keep bot games separate. | ⏳ Planned |
| 419–422 | Rating protection | Bot matches do not inflate/deflate ranked ratings unless explicitly designed. | ⏳ Planned |
| 423–426 | Bot calibration | Adjust difficulty via self-play and human practice results. | ⏳ Planned |
| 427–430 | Matchmaking tests | Queue formation, bot insertion, rating separation, transparency. | ⏳ Planned |

---

## Phase 23 — Production AI Operations, Safety, and Continuous Improvement

| Day range | Focus | Deliverables | Status |
|---:|---|---|---|
| 431–434 | AI feature flags | Server flags for bot availability, strategy version, delays, search depth. | ⏳ Planned |
| 435–438 | Strategy versioning | Version bot policies; live match records strategy that generated decisions. | ⏳ Planned |
| 439–442 | Safe rollout | Staged: local → dev → staging → limited prod cohort → full release. | ⏳ Planned |
| 443–446 | Bot health dashboards | Bot action latency, failures, invalid decisions, fallback rate, queue impact. | ⏳ Planned |
| 447–450 | Alerts and safeguards | Alerts for timeouts, elevated invalid command rate, stuck matches, win-rate shifts. | ⏳ Planned |
| 451–454 | Fairness audits | Verify bots use only player-visible/public info, cannot inspect hidden racks/stock order. | ⏳ Planned |
| 455–458 | Replay and audit logs | Privacy-safe action traces for debugging disputed games. | ⏳ Planned |
| 459–462 | Human feedback loop | Post-match bot feedback: difficulty, suspicious behavior, fun rating. | ⏳ Planned |
| 463–466 | Automated regression corpus | Convert production bugs and strange bot behavior into deterministic regression scenarios. | ⏳ Planned |
| 467–470 | Cost and capacity management | Concurrency limits, execution budgets, worker scaling, cost monitoring. | ⏳ Planned |
| 471–474 | Disaster recovery | Behavior if bot service fails: pause/remove/substitute fallback/safely end match. | ⏳ Planned |
| 475–480 | AI release candidate | Load tests, fairness tests, replay audits, security review, staged-release checklist. | ⏳ Planned |

---

## Phase 24 — Production Deployment and Infrastructure

| Day range | Focus | Deliverables | Status |
|---:|---|---|---|
| 481–484 | Production architecture | Document target hosting, regions, domains, TLS, database topology, backups, scaling, failure boundaries. | ⏳ Planned |
| 485–488 | Infrastructure as Code | Terraform/Pulumi for networking, DB, container runtime, secrets, observability. | ⏳ Planned |
| 489–492 | Container hardening | Production Docker images pinned, non-root, minimal layers, health checks, image scanning. | ⏳ Planned |
| 493–496 | Environment management | Separate local/dev/staging/production configs; prevent dev settings reaching prod. | ⏳ Planned |
| 497–500 | Secrets management | Managed secret store for DB credentials, Nakama keys, signing secrets; never commit. | ⏳ Planned |
| 501–504 | TLS and edge security | HTTPS, certificates, secure headers, CORS, rate limits, WAF/CDN. | ⏳ Planned |
| 505–508 | Database migrations | Versioned migrations, rollback strategy, migration checks, backup-before-migrate. | ⏳ Planned |
| 509–512 | Backup and restore drills | Automate PostgreSQL backups and test restoration into isolated environment. | ⏳ Planned |
| 513–516 | Staging environment | Deploy production-like staging and run automated smoke tests against it. | ⏳ Planned |
| 517–520 | Production release process | Approvals, release notes, rollback process, incident contacts, deployment checklist. | ⏳ Planned |

---

## Phase 25 — Observability, Analytics, and Game Telemetry

| Day range | Focus | Deliverables | Status |
|---:|---|---|---|
| 521–524 | Logging standard | Structured fields: request ID, match ID, player ID hash, bot ID, opcode, transition result, duration, error code. | ⏳ Planned |
| 525–528 | Metrics foundation | Service metrics: requests, match count, active users, command volume, errors, latency, memory, CPU, DB health. | ⏳ Planned |
| 529–532 | Match metrics | Track creation, joins, starts, disconnects, completion, abandonment, average duration, player-count distribution. | ⏳ Planned |
| 533–536 | Gameplay metrics | Track draw type, discard pickup use, opening-meld timing, joker usage, win conditions, rule-validation failures. | ⏳ Planned |
| 537–540 | Bot metrics | Track bot difficulty, decision latency, search depth, fallback rate, win rate, invalid-action attempts. | ⏳ Planned |
| 541–544 | Product analytics privacy design | Event retention, player consent, anonymization, aggregation, PII boundaries. | ⏳ Planned |
| 545–548 | Dashboards | Infrastructure health, live matches, gameplay funnel, bot behavior, error rates. | ⏳ Planned |
| 549–552 | Alerting | Alerts for DB failures, Nakama restart loops, elevated command errors, stuck matches, slow bot actions, unusual abandonment. | ⏳ Planned |
| 553–556 | Distributed tracing | Request/match tracing where feasible, especially for external services and AI workers. | ⏳ Planned |
| 557–560 | Analytics validation | Verify event correctness using controlled test matches and ensure no hidden rack details sent to analytics. | ⏳ Planned |

---

## Phase 26 — Security, Privacy, and Anti-Cheat

| Day range | Focus | Deliverables | Status |
|---:|---|---|---|
| 561–564 | Threat model | Document threats: account abuse, token theft, malformed input, collusion, automation abuse, data leakage, DDoS, admin compromise. | ⏳ Planned |
| 565–568 | Authentication review | Harden account/session handling, token expiry, refresh flow, password policies, recovery. | ⏳ Planned |
| 569–572 | Authorization review | Verify all RPCs, match commands, admin endpoints, storage access enforce least privilege. | ⏳ Planned |
| 573–576 | Protocol fuzzing | Fuzz malformed, oversized, duplicated, reordered, adversarial match commands. | ⏳ Planned |
| 577–580 | Rate limiting | Limits for login, registration, match creation, joining, command spam, chat, bot-related endpoints. | ⏳ Planned |
| 581–584 | Anti-cheat audit | Verify server authority across all actions (racks, stock, meld validation, discard pickup, joker replacement, win detection). | ⏳ Planned |
| 585–588 | Replay integrity | Sign or hash match event streams to detect tampering in stored replays. | ⏳ Planned |
| 589–592 | Privacy controls | Account-data export/delete workflows and retention policies per applicable privacy laws. | ⏳ Planned |
| 593–596 | Dependency security | Dependency scanning, container scanning, security update policy, critical-CVE response. | ⏳ Planned |
| 597–600 | Security test and review | Penetration testing, tabletop incident exercises, security regression tests. | ⏳ Planned |

---

## Phase 27 — Player Accounts, Profiles, and Social Foundation

| Day range | Focus | Deliverables | Status |
|---:|---|---|---|
| 601–604 | Account model | Define guest, registered, linked-provider, verified, suspended, deleted states. | ⏳ Planned |
| 605–608 | Registration and login UX/API | Secure email/password and/or platform-provider login flows. | ⏳ Planned |
| 609–612 | Guest accounts | Low-friction guest play with upgrade/link path. | ⏳ Planned |
| 613–616 | Player profile | Display name, avatar selection/upload policy, language, region, safe public fields. | ⏳ Planned |
| 617–620 | Profile moderation controls | Validation, reserved names, report flow, avatar limits, admin moderation tools. | ⏳ Planned |
| 621–624 | Friends model | Friend requests, accepted friends, blocked users, privacy controls. | ⏳ Planned |
| 625–628 | Presence and invitations | Opt-in online presence and private match invitations. | ⏳ Planned |
| 629–632 | Recent players | Store recent opponents with add friend/block/report/invite again. | ⏳ Planned |
| 633–636 | Social privacy tests | Verify blocks, privacy settings, invitations, profile visibility. | ⏳ Planned |
| 637–640 | Account lifecycle documentation | Recovery, deletion, moderation, support workflows, audit trails. | ⏳ Planned |

---

## Phase 28 — Communication, Moderation, and Community Safety

| Day range | Focus | Deliverables | Status |
|---:|---|---|---|
| 641–644 | Communication scope | Define channels: table chat, private messages, friends chat, lobby chat, system announcements. | ⏳ Planned |
| 645–648 | Table chat MVP | Implement match/table chat with limits, ordering, reconnection behavior. | ⏳ Planned |
| 649–652 | Chat persistence decision | Document whether messages are ephemeral, match-scoped, retained, and retention. | ⏳ Planned |
| 653–656 | Spam protection | Rate limits, duplicate-message controls, cooldowns, flood detection. | ⏳ Planned |
| 657–660 | Player blocking | Ensure blocked users cannot message/invite/interact. | ⏳ Planned |
| 661–664 | Reporting workflow | Report categories, evidence capture, moderation queue, player feedback states. | ⏳ Planned |
| 665–668 | Automated content filtering | Configurable profanity/spam filters with language-aware policies. | ⏳ Planned |
| 669–672 | Moderator tools | Review queues, mutes, warnings, suspensions, bans, audit logs. | ⏳ Planned |
| 673–676 | Appeals and enforcement policy | Document moderation standards, appeal process, staff permissions, evidence retention. | ⏳ Planned |
| 677–680 | Community safety tests | Test rate limits, blocking, reporting, moderation actions, auditability. | ⏳ Planned |

---

## Phase 29 — Economy, Progression, and Cosmetics

*Strictly non-gambling. No cash wagering, real-money stakes, loot boxes, or gambling-like mechanics.*

| Day range | Focus | Deliverables | Status |
|---:|---|---|---|
| 681–684 | Economy principles | Document non-gambling rules, fair-play constraints, virtual-currency policy, regional compliance, parental controls. | ⏳ Planned |
| 685–688 | Player progression | Experience, levels, badges, non-competitive achievement tracking. | ⏳ Planned |
| 689–692 | Match rewards | Non-monetary progression rewards for completed matches, good sportsmanship, tutorials. | ⏳ Planned |
| 693–696 | Cosmetics model | Cosmetic-only items: table themes, tile skins, avatars, emotes, profile frames, victory effects. | ⏳ Planned |
| 697–700 | Inventory service | Ownership, unlocks, equipped items, grants, revocations, audit history. | ⏳ Planned |
| 701–704 | Daily/weekly challenges | Optional goals rewarding engagement without pressuring risky spending. | ⏳ Planned |
| 705–708 | Achievement system | Rule-based achievements, progress tracking, notifications, anti-exploit checks. | ⏳ Planned |
| 709–712 | Store design, if applicable | If monetization approved, cosmetic purchases with transparent pricing, receipts, refunds, parental controls, no pay-to-win. | ⏳ Planned |
| 713–716 | Economy abuse prevention | Detect duplicate grants, replayed purchase events, reward farming, inventory inconsistencies. | ⏳ Planned |
| 717–720 | Economy tests and compliance review | Test transactions, entitlement handling; legal/compliance review before release. | ⏳ Planned |

---

## Phase 30 — Competitive Play, Rankings, and Tournaments

| Day range | Focus | Deliverables | Status |
|---:|---|---|---|
| 721–724 | Competitive ruleset | Define ranked match eligibility, player count, disconnect policy, bot policy, scoring format, seasonal structure. | ⏳ Planned |
| 725–728 | Rating algorithm | Select and implement Elo/Glicko/TrueSkill-like logic for multiplayer. | ⏳ Planned |
| 729–732 | Rating event pipeline | Record rating-relevant results atomically and make recalculation/recovery possible. | ⏳ Planned |
| 733–736 | Leaderboards | Global, regional, friends, seasonal, game-mode leaderboards. | ⏳ Planned |
| 737–740 | Ranked matchmaking | Match via rating, region, latency, queue duration, anti-smurf controls. | ⏳ Planned |
| 741–744 | Season system | Season start/end, soft reset policy, rewards, historical ranks, announcement flow. | ⏳ Planned |
| 745–748 | Tournament design | Formats: scheduled tables, Swiss, bracket, elimination, points-based. | ⏳ Planned |
| 749–752 | Tournament engine | Registration, seating, round creation, advancement, standings, tie-breaking. | ⏳ Planned |
| 753–756 | Tournament integrity | Check-in, disconnect rules, bot restrictions, anti-collusion signals, staff intervention, replay review. | ⏳ Planned |
| 757–760 | Tournament rewards | Cosmetic/badge rewards and transparent eligibility. Avoid real-money prizes unless legal/licensed. | ⏳ Planned |
| 761–764 | Competitive testing | Simulate rating updates, tournament progression, dropouts, ties, rollback. | ⏳ Planned |
| 765–770 | Competitive launch readiness | Load test queues/leaderboards, audit integrity, publish rules, run limited beta season. | ⏳ Planned |

---

## Phase 31 — Client Platforms, UX, Accessibility, and Localization

| Day range | Focus | Deliverables | Status |
|---:|---|---|---|
| 771–774 | Client platform strategy | Decide web/desktop/Android/iOS priorities; define shared protocol/client SDK strategy. | ⏳ Planned |
| 775–778 | Design system | Reusable typography, colours, spacing, buttons, modals, notifications, tile components, table layouts. | ⏳ Planned |
| 779–782 | Rummy table UX | Core table interaction: rack sorting, drag/drop, tile selection, undo-before-submit, turn feedback. | ⏳ Planned |
| 783–786 | Meld-building UX | Accessible creation of runs/sets, explicit joker assignment, invalid-state feedback, server-confirmed actions. | ⏳ Planned |
| 787–790 | Responsive layouts | Desktop, tablet, mobile layouts without reducing rule clarity. | ⏳ Planned |
| 791–794 | Accessibility baseline | Keyboard control, screen-reader labels, focus handling, colour-blind-safe tile identifiers, scalable text, reduced motion. | ⏳ Planned |
| 795–798 | Localization foundation | Externalize UI text, support Romanian and English first, pluralization/date/number rules. | ⏳ Planned |
| 799–802 | Connection UX | Reconnecting states, resync UI, action-pending feedback, network error handling, match recovery. | ⏳ Planned |
| 803–806 | Onboarding and tutorial | Guided rules introduction, practice match with beginner bots, contextual help, rule glossary. | ⏳ Planned |
| 807–810 | Notifications | Opt-in invitations, turn reminders, friend activity, tournaments, moderation actions. | ⏳ Planned |
| 811–814 | Client QA automation | End-to-end client tests for login, match creation, joining, draw/discard, melding, reconnection, error states. | ⏳ Planned |
| 815–820 | UX beta and polish | Usability testing, fix high-impact friction, improve performance, prepare release notes. | ⏳ Planned |

---

## Extended Roadmap Summary (Future)

| Phase | Days | Primary Outcome | Status |
|---:|---|---|---|
| 17 | 147–240 | Future product features (scoring, social, matchmaking, bots, production ops, polish) | ⏳ Not started |
| 18 | 241–480 | AI Bots: foundation through production AI operations (Phases 18–23 above) | ⏳ Not started |
| 24 | 481–520 | Production infrastructure and deployment | ⏳ Not started |
| 25 | 521–560 | Observability, analytics, telemetry | ⏳ Not started |
| 26 | 561–600 | Security, privacy, anti-cheat | ⏳ Not started |
| 27 | 601–640 | Accounts, profiles, social foundation | ⏳ Not started |
| 28 | 641–680 | Chat, moderation, community safety | ⏳ Not started |
| 29 | 681–720 | Safe progression, cosmetics, optional economy (non-gambling) | ⏳ Not started |
| 30 | 721–770 | Ranked play, leaderboards, tournaments | ⏳ Not started |
| 31 | 771–820 | Multi-platform client, UX, accessibility, localization | ⏳ Not started |

By **Day 820**, the project evolves from the current `rummy-mvp-rc1` authoritative backend (`Phases 1–16` done) into a production-ready, AI-supported, secure, social, competitive, accessible multiplayer platform. Execution remains Handmade Hero: one vertical slice per day, `make check` + `docker compose build` green, focused commit, push.

---

*This roadmap lists days that **will be** implemented. For what **is** implemented, see `docs/IMPLEMENTED.md`. Last updated: 2026-08-26 at `main@41e794b` (and `rummy-mvp-rc1` `cfce62e`).*
