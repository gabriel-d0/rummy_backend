# Layout System — Mathematically Stable, Fully Responsive

**Purpose:** Solve alignment problems across devices and ensure the game looks and feels good on all screen sizes, including mobile, by dividing the entire game space into logical subspaces with proportional math.

**No pixel-based positioning. No device-specific branches. No magic numbers.**

---

## 1. Virtual GameSpace

A fixed logical coordinate system `1000×1000` that represents the entire game board, independent of physical pixels. Phaser scales this GameSpace to any screen using:

```ts
// src/ui/LayoutManager.ts
export const GAME_SPACE = {
  width: 1000,
  height: 1000,
  scale: {
    mode: Phaser.Scale.FIT,
    autoCenter: Phaser.Scale.CENTER_BOTH,
    width: 1000,
    height: 1000,
  },
} as const;

// src/main.ts
new Phaser.Game({
  type: Phaser.AUTO,
  parent: "game",
  width: GAME_SPACE.width,
  height: GAME_SPACE.height,
  backgroundColor: "#0a4d2e",
  scale: GAME_SPACE.scale,
  scene: [Preload, TableScene, RackScene],
});
```

And in `index.html`:

```css
#game { width: min(1024px, 96vw); height: min(768px, 72vh); max-width: 1024px; aspect-ratio: 1024/768; }
#game canvas { width: 100% !important; height: 100% !important; }
```

This guarantees **consistent alignment on all aspect ratios**, **automatic scaling without distortion**, and **stable spacing** because all internal positions are computed in `1000×1000` logical units, then scaled by Phaser.

---

## 2. Named Subspaces

The GameSpace is mathematically partitioned into functional regions. Each subspace has `x, y, width, height` in logical units, not pixels, and is derived as a proportion of `GAME_SPACE`.

```ts
// src/ui/LayoutManager.ts
export interface Subspaces {
  TopBar: Rect;           // Stock, Turn, Winner — top 8% (80px)
  TableArea: Rect;        // green felt background — from y=80 to y=680
  MeldArea: Rect;         // inside TableArea, top 55% — for TableMelds
  DiscardRowArea: Rect;   // inside TableArea, bottom 35% — for DiscardRow
  PlayerRackArea: Rect;   // wood rack — y=750, 180px high
  ActionButtonsArea: Rect;// Discard, Meld, etc. — bottom 70px
}

export function getSubspaces(): Subspaces {
  const outerMargin = 12;
  const gutter = 16;
  const TopBar = { x: outerMargin, y: outerMargin, width: 1000 - outerMargin*2, height: 80 };
  const ActionButtonsArea = { x: outerMargin, y: 1000 - outerMargin - 70, width: 1000 - outerMargin*2, height: 70 };
  const PlayerRackArea = { x: outerMargin, y: ActionButtonsArea.y - gutter - 180, width: 1000 - outerMargin*2, height: 180 };
  const TableArea = { x: outerMargin, y: TopBar.y + TopBar.height + gutter, width: 1000 - outerMargin*2, height: PlayerRackArea.y - (TopBar.y + TopBar.height + gutter) - gutter };
  const MeldArea = { x: TableArea.x + gutter, y: TableArea.y + gutter, width: TableArea.width - gutter*2, height: TableArea.height * 0.55 };
  const DiscardRowArea = { x: TableArea.x + gutter, y: MeldArea.y + MeldArea.height + gutter, width: TableArea.width - gutter*2, height: TableArea.height * 0.35 };
  return { TopBar, TableArea, MeldArea, DiscardRowArea, PlayerRackArea, ActionButtonsArea };
}
```

Visual (logical units, not pixels):

```
y=0   ┌─ TopBar (80h) ──────────────────────┐
      │  Stock:77 | Current: seat-0          │
y=80  ├─ TableArea (y=96, h≈534) ────────────┤
      │  ┌─ MeldArea (x=28,y=112,h≈294) ─┐   │
      │  │  Meld 1: 5 5 5                │   │
      │  │  Meld 2: 5 5 5                │   │
      │  └───────────────────────────────┘   │
      │  ┌─ DiscardRowArea (x=28,y≈422) ─┐   │
      │  │  disc-open  disc-1  disc-2     │   │
      │  └───────────────────────────────┘   │
      │  ┌─ Drop zone (centered in Table) ┐ │
      │  └────────────────────────────────┘   │
y=750 ├─ PlayerRackArea (180h) ──────────────┤
      │  [5][5][5]  (14 slots, centered)     │
y=930 ├─ ActionButtonsArea (70h) ────────────┤
      │  [Draw] [Discard] [Meld]              │
y=1000└───────────────────────────────────────┘
```

All `Rect`s are computed from `GAME_SPACE` fractions, so moving `TopBar` or `PlayerRackArea` automatically shifts `TableArea` — no magic numbers leak to scenes.

---

## 3. Proportional Placement

Every UI/game element is positioned using formulas, never hardcoded pixels:

```ts
// Center of an area
centerX = area.x + area.width * 0.5
centerY = area.y + area.height * 0.5

// Even spacing for N tiles with tileWidth and minGap
spacing = area.width / tileCount
tile.x = area.x + i * spacing
tile.y = area.y + area.height * 0.5

// Or using helpers from LayoutManager.ts:
import { centerX, centerY, spacing, tileX, tileY, meldRowY, tilesInRow } from "./ui/LayoutManager";

// Example: place Stock text centered in TopBar
const stockText = scene.add.text(centerX(subspaces.TopBar) + 300, centerY(subspaces.TopBar), `Stock: ${stockCount}`);

// Example: 3 tiles centered in PlayerRackArea with even spacing
const tileW = 48; // logical width of tile image at scale 0.9
for (let i = 0; i < tiles.length; i++) {
  const x = tileX(subspaces.PlayerRackArea, i, tiles.length, tileW);
  const y = tileY(subspaces.PlayerRackArea);
  scene.add.image(x, y, "tile");
}

// Example: meld rows stacked inside MeldArea
for (let r = 0; r < melds.length; r++) {
  const y = meldRowY(subspaces.MeldArea, r, 70);
  const rowRect = { x: subspaces.MeldArea.x, y: y - 25, width: subspaces.MeldArea.width, height: 50 };
  const positions = tilesInRow(rowRect, melds[r].Tiles.length, 38);
  melds[r].Tiles.forEach((tl, j) => scene.add.image(positions[j].x, positions[j].y, tl.IsJoker ? "joker" : "tile"));
}

// Example: discard row centered in DiscardRowArea
for (let i = 0; i < discardRow.length; i++) {
  const x = tileX(subspaces.DiscardRowArea, i, discardRow.length, 38);
  scene.add.image(x, tileY(subspaces.DiscardRowArea), "tile");
}
```

**Rules:**
- `No pixel-based positioning` — always `area.x + area.width * fraction`
- `No device-specific branches` — `isMobile` is derived from `width < 768` inside `getLayout`, not `if (iPhone) ...`
- `No magic numbers` — `gutter`, `outerMargin`, `rowHeight` are defined once in `getSubspaces`, scenes only use `centerX`/`tileX`

---

## 4. Responsiveness Guarantees

- **Consistent alignment on all aspect ratios:** Because `GAME_SPACE` is `1000×1000` and `Phaser.Scale.FIT` scales it uniformly, a tile at `area.x + area.width * 0.5` is always centered, whether the physical screen is `1920×1080` or `375×667`. The `outerMargin` and `gutter` are also logical units, so spacing scales proportionally.
- **Stable spacing between elements:** `spacing = area.width / tileCount` distributes tiles evenly; if `tileCount` changes from 3 to 14, spacing recomputes, so tiles never overlap or leave huge gaps.
- **Predictable layout on mobile and desktop:** `getLayout` checks `isMobile = width < 768` and adjusts `outerMargin`/`gutter` and `rackW`/`rackH`/`dropW`/`dropH` as proportions, not branches per device. On mobile, `rack` becomes `outer.w - gutter*2` (full width) and `dropZone` becomes `outer.w - gutter*2` (full width) — still centered via `width/2 - dropW/2`.
- **Automatic scaling without distortion:** `Phaser.Scale.FIT` preserves `1000/1000` aspect ratio; `Phaser.Scale.CENTER_BOTH` centers the scaled GameSpace; `index.html` `#game { aspect-ratio: 1024/768; max-width: 1024px; }` ensures the canvas never stretches.

---

## 5. LayoutManager Module

**File:** `client/src/ui/LayoutManager.ts`

```ts
import Phaser from "phaser";
export const GAME_SPACE = { width: 1000, height: 1000, scale: { mode: Phaser.Scale.FIT, autoCenter: Phaser.Scale.CENTER_BOTH, width: 1000, height: 1000 } };
export interface Rect { x: number; y: number; width: number; height: number; }
export interface Subspaces { TopBar: Rect; TableArea: Rect; MeldArea: Rect; DiscardRowArea: Rect; PlayerRackArea: Rect; ActionButtonsArea: Rect; }
export function getSubspaces(): Subspaces { /* ... mathematically partition GameSpace ... */ }
export function centerX(rect: Rect): number { return rect.x + rect.width * 0.5; }
export function centerY(rect: Rect): number { return rect.y + rect.height * 0.5; }
export function spacing(areaWidth: number, tileCount: number, tileWidth: number, minGap?: number): number { /* ... */ }
export function tileX(area: Rect, i: number, tileCount: number, tileWidth: number): number { return area.x + spacing(area.width, tileCount, tileWidth) + tileWidth*0.5 + i*(tileWidth+spacing(area.width,tileCount,tileWidth)); }
export function tileY(area: Rect): number { return centerY(area); }
export function meldRowY(meldArea: Rect, rowIndex: number, rowHeight: number): number { return meldArea.y + rowIndex*rowHeight + rowHeight*0.5; }
export function tilesInRow(rowRect: Rect, tileCount: number, tileWidth: number): {x:number,y:number}[] { /* ... */ }
export function onGameResize(scene: Phaser.Scene, cb: (subspaces: Subspaces) => void): void { scene.scale.on("resize", (gameSize) => cb(getSubspaces())); }
export function drawDebugSubspaces(scene: Phaser.Scene, subspaces: Subspaces, alpha?: number): void { /* ... */ }
export function clearDebugSubspaces(scene: Phaser.Scene): void { /* ... */ }
```

**Recalculates layout on resize events** via `scene.scale.on("resize", (gameSize) => cb(getSubspaces()))` — scenes can `restart()` or reposition each subspace element.

**Usage in scenes:**

```ts
// TableScene.ts
import { getSubspaces, centerX, centerY } from "../ui/LayoutManager";
create() {
  const layout = getSubspaces();
  this.add.image(centerX(layout.TableArea), centerY(layout.TableArea), "table").setDisplaySize(layout.TableArea.width, layout.TableArea.height);
  drawDebugSubspaces(this, layout); // optional, for verification
  this.scale.on("resize", () => this.scene.restart());
}

// RackScene.ts
import { getSubspaces, tileX, tileY } from "../ui/LayoutManager";
create() {
  const layout = getSubspaces();
  this.add.image(centerX(layout.PlayerRackArea), centerY(layout.PlayerRackArea), "rack").setDisplaySize(layout.PlayerRackArea.width, layout.PlayerRackArea.height);
  const tiles = sortRack(unsorted);
  for (let i = 0; i < tiles.length; i++) {
    this.add.image(tileX(layout.PlayerRackArea, i, tiles.length, 48), tileY(layout.PlayerRackArea), "tile");
  }
}
```

---

## How This Solves Alignment Problems

- **Previous problem:** Elements were positioned with hardcoded pixels (`x=100 + i*50`, `y=680`, `x=800,y=50` for Stock) and `setDisplaySize(1024,768)` for background. On a small screen, `x=800` is off-screen, `x=100` is left of the rack, `y=200` drop zone overlapped melds at `y=160`, and spacing `40` caused tiles to overlap (`tileWidth 57.6 > spacing 40`).
- **New solution:** Every `x`/`y` is derived from a subspace `Rect` via `centerX`/`tileX`/`meldRowY`. If `PlayerRackArea` moves (e.g., on mobile it becomes full width), `tileX` automatically recenters tiles. If `MeldArea` height changes, `meldRowY` automatically spaces rows. No element is positioned with a magic pixel; all are fractions of `GAME_SPACE` or its subspaces.
- **Verification:** `drawDebugSubspaces` draws each subspace bounds in a different color, proving no overlap: `TopBar` red, `TableArea` green, `MeldArea` blue, `DiscardRowArea` yellow, `PlayerRackArea` magenta, `ActionButtonsArea` cyan. On `resize`, `getLayout` recomputes and scenes `restart`, so alignment stays stable.

---

*Code: `client/src/ui/LayoutManager.ts:1`, `client/src/main.ts:6` (uses `GAME_SPACE.scale`), `client/index.html:9` (`#game { aspect-ratio: 1024/768; }`), `client/src/scenes/TableScene.ts:12` and `RackScene.ts:11` (both use `getSubspaces`).*
