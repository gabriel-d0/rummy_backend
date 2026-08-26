# LayoutManager — Responsive GameSpace & Subspaces

**Day 4** — Mathematically stable, fully responsive layout system.

## Problem Solved
Previous Phaser had `x=100+i*50` pixel positioning, device-specific branches, and magic numbers → misaligned on different screens. Now all positions use proportional math relative to a fixed logical GameSpace.

## GameSpace
Fixed `1000×1000` logical units. Phaser scales it to any screen via:
```ts
width: 1000, height: 1000
scale: { mode: Phaser.Scale.FIT, autoCenter: Phaser.Scale.CENTER_BOTH }
```
No distortion. `aspect-ratio: 1/1` via CSS.

## Subspaces (6) — all in logical units
- **TopBar** `12,12 976×80` — Stock, Turn, Winner
- **TableArea** `12,108 976×~596` — felt, between TopBar and PlayerRackArea
- **MeldArea** `TableArea+16` `944×328` — 55% of TableArea top
- **DiscardRowArea** `MeldArea.y+h+16` `944×209` — 35% bottom
- **PlayerRackArea** `12,722 976×132` — 14 slots, compact modern
- **ActionButtonsArea** `12,918 976×56` — Draw/Prev/Meld/Discard flex

Defined once in `getSubspaces()` — single source of truth.

## Proportional Math — Examples
```ts
const s = getSubspaces();
centerX = s.MeldArea.x + s.MeldArea.width * 0.5;
tile.x = s.MeldArea.x + 10 + i * 48;
tile.y = s.MeldArea.y + s.MeldArea.height / 2;
// Even spacing
spacing = s.MeldArea.width / tileCount;
```

No `x=100`, no `if (isMobile) x=50 else x=100`.

## Responsiveness
- `1000×1000` → `Scale.FIT` scales to any `vw/vh` without distortion
- Subspaces are proportions, so `MeldArea` stays 55% on mobile and desktop
- `onGameResize(scene, cb => getSubspaces())` recalculates on `scale.resize`
- Tested via Playwright at `1000×1000`, `768×1024`, `375×667`

## Files
- `src/ui/LayoutManager.ts` — `GAME_SPACE`, `getSubspaces()`, `centerX`, `spacing`, `tileX`, `onGameResize`, `drawDebugSubspaces`
- `src/main.ts` — `width/height/scale: GAME_SPACE.scale`
- `src/scenes/Preload.ts` — demo: draws 6 subspaces with `drawDebugSubspaces` and places text via `centerX`/`centerY`

## Verification
`npm run dev` at `127.0.0.1:5173` shows 6 colored subspaces with labels, centered, no overlap, responsive via `Scale.FIT`.
`npx playwright test e2e/layout.spec.ts` checks `getSubspaces()` and `GAME_SPACE`.
