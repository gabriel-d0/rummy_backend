import Phaser from "phaser";

/**
 * Day 20+: LayoutManager — mathematically stable, fully responsive layout system.
 *
 * Purpose: Solve alignment problems across devices and ensure the game looks
 * and feels good on all screen sizes, including mobile, by dividing the
 * entire game space into logical subspaces with proportional math.
 *
 * No pixel-based positioning. No device-specific branches. No magic numbers.
 */

// 1. Virtual GameSpace — fixed logical coordinate system 1000×1000
// Phaser scales this GameSpace to any screen using Scale.FIT + CENTER_BOTH
export const GAME_SPACE = {
  width: 1000,
  height: 1000,
  // Scale config for Phaser (used in main.ts)
  scale: {
    mode: Phaser.Scale.FIT,
    autoCenter: Phaser.Scale.CENTER_BOTH,
    width: 1000,
    height: 1000,
  },
} as const;

export interface Rect {
  x: number;
  y: number;
  width: number;
  height: number;
}

// 2. Named subspaces — mathematically partition GameSpace into functional regions
// All values are in logical units (0..1000), not pixels
export interface Subspaces {
  /** TopBar: Stock, Turn, Winner — top 8% */
  TopBar: Rect;
  /** TableArea: green felt background — from y=80 to y=680 */
  TableArea: Rect;
  /** MeldArea: inside TableArea — for TableMelds */
  MeldArea: Rect;
  /** DiscardRowArea: ordered discard row */
  DiscardRowArea: Rect;
  /** PlayerRackArea: wood rack with 14 slots */
  PlayerRackArea: Rect;
  /** ActionButtonsArea: Discard, Meld, Extend, etc. */
  ActionButtonsArea: Rect;
}

/**
 * Defines subspaces as proportions of GameSpace.
 * This is the single source of truth — no other file should hardcode x/y.
 */
export function getSubspaces(): Subspaces {
  // Using proportional math: each subspace is a fraction of GAME_SPACE
  // Heights sum to 1000 with gutters, widths are full with outer margins
  const outerMargin = 12; // logical units
  const gutter = 16;

  // TopBar: 80px high at top
  const TopBar: Rect = {
    x: outerMargin,
    y: outerMargin,
    width: GAME_SPACE.width - outerMargin * 2,
    height: 80,
  };

  // TableArea: occupies middle, from below TopBar to above PlayerRackArea
  // Modern: compact rack 132px + action bar 56px for cleaner bottom
  const ActionButtonsArea: Rect = {
    x: outerMargin,
    y: GAME_SPACE.height - outerMargin - 56,
    width: GAME_SPACE.width - outerMargin * 2,
    height: 56,
  };

  const PlayerRackArea: Rect = {
    x: outerMargin,
    y: ActionButtonsArea.y - gutter - 132,
    width: GAME_SPACE.width - outerMargin * 2,
    height: 132,
  };

  const TableArea: Rect = {
    x: outerMargin,
    y: TopBar.y + TopBar.height + gutter,
    width: GAME_SPACE.width - outerMargin * 2,
    height: PlayerRackArea.y - (TopBar.y + TopBar.height + gutter) - gutter,
  };

  // Inside TableArea, define MeldArea and DiscardRowArea as proportional subdivisions
  // MeldArea: top 55% of TableArea
  const MeldArea: Rect = {
    x: TableArea.x + gutter,
    y: TableArea.y + gutter,
    width: TableArea.width - gutter * 2,
    height: TableArea.height * 0.55,
  };

  // DiscardRowArea: bottom 35% of TableArea (leaving gutter between)
  const DiscardRowArea: Rect = {
    x: TableArea.x + gutter,
    y: MeldArea.y + MeldArea.height + gutter,
    width: TableArea.width - gutter * 2,
    height: TableArea.height * 0.35,
  };

  return {
    TopBar,
    TableArea,
    MeldArea,
    DiscardRowArea,
    PlayerRackArea,
    ActionButtonsArea,
  };
}

// 3. Proportional placement helpers — no pixel-based positioning

/** Center X of a rect */
export function centerX(rect: Rect): number {
  return rect.x + rect.width * 0.5;
}

/** Center Y of a rect */
export function centerY(rect: Rect): number {
  return rect.y + rect.height * 0.5;
}

/** Even spacing for N tiles inside an area with padding */
export function spacing(
  areaWidth: number,
  tileCount: number,
  tileWidth: number,
  minGap = 6
): number {
  if (tileCount <= 1) return 0;
  const totalTileWidth = tileCount * tileWidth;
  const available = areaWidth - totalTileWidth;
  // Distribute remaining space as gaps, but respect minGap
  const gap = available / (tileCount + 1);
  return Math.max(minGap, gap + tileWidth);
}

/** X position for tile i in an area with even spacing */
export function tileX(area: Rect, i: number, tileCount: number, tileWidth: number): number {
  const gap = spacing(area.width, tileCount, tileWidth);
  // Start with gap padding, then tile + gap
  return area.x + gap + tileWidth * 0.5 + i * (tileWidth + gap);
}

/** Y position centered in area */
export function tileY(area: Rect): number {
  return centerY(area);
}

/** Helper to get X for meld row i (stacked vertically inside MeldArea) */
export function meldRowY(meldArea: Rect, rowIndex: number, rowHeight: number): number {
  return meldArea.y + rowIndex * rowHeight + rowHeight * 0.5;
}

/** Helper to get positions for tiles in a meld row */
export function tilesInRow(
  rowRect: Rect,
  tileCount: number,
  tileWidth: number
): { x: number; y: number }[] {
  const positions: { x: number; y: number }[] = [];
  const gap = spacing(rowRect.width, tileCount, tileWidth);
  const startX = rowRect.x + gap;
  const y = centerY(rowRect);
  for (let i = 0; i < tileCount; i++) {
    positions.push({ x: startX + i * (tileWidth + gap) + tileWidth * 0.5, y });
  }
  return positions;
}

// 4. Responsiveness — recalculates layout on resize events

export function onGameResize(scene: Phaser.Scene, cb: (subspaces: Subspaces) => void): void {
  scene.scale.on("resize", (gameSize: Phaser.Structs.Size) => {
    // GameSpace is fixed logical 1000x1000, but we recalculate subspaces
    // in case we ever make them depend on actual width/height (currently they are fixed proportions,
    // so this is just for future: if we make subspaces responsive to aspect ratio, this will handle it)
    const subspaces = getSubspaces();
    // Also update GameSpace scaling if needed
    void gameSize; // for future use with actual width/height
    cb(subspaces);
  });
}

// 5. Debug helper — draw subspace bounds for visual verification
export function drawDebugSubspaces(scene: Phaser.Scene, subspaces: Subspaces, alpha = 0.04): void {
  const colors: Record<keyof Subspaces, number> = {
    TopBar: 0xff0000,
    TableArea: 0x00ff00,
    MeldArea: 0x0000ff,
    DiscardRowArea: 0xffff00,
    PlayerRackArea: 0xff00ff,
    ActionButtonsArea: 0x00ffff,
  };
  for (const [name, rect] of Object.entries(subspaces) as [keyof Subspaces, Rect][]) {
    const r = scene.add
      .rectangle(
        rect.x + rect.width / 2,
        rect.y + rect.height / 2,
        rect.width,
        rect.height,
        colors[name],
        alpha
      )
      .setStrokeStyle(1, colors[name], 0.5);
    r.setData("isDebugSubspace", true);
    const label = scene.add.text(rect.x + 4, rect.y + 4, name, {
      fontFamily: "monospace",
      fontSize: "8px",
      color: "#ffffff",
    });
    label.setData("isDebugSubspace", true);
  }
}

export function clearDebugSubspaces(scene: Phaser.Scene): void {
  const existing = scene.children.list.filter((c) => (c as any).getData?.("isDebugSubspace"));
  for (const c of existing) c.destroy();
}
