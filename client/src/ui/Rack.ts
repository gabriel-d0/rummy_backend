import Phaser from "phaser";
import type { TileInstance } from "../state/snapshot";

// Day 11: Rack rendering (static) — draws 14 this.add.image per tile at x = 100 + i*50
// Mirrors Go TileInstance{ID,Colour,Rank,IsJoker} and PrivateView.OwnRack (redaction)

// Day 16: Tile selection — toggles selected Set<TileInstanceId> and tints selected tile 0xffff00
export const selected: Set<string> = new Set();

export function onTileClicked(tileId: string): void {
  if (selected.has(tileId)) {
    selected.delete(tileId);
  } else {
    selected.add(tileId);
  }
  console.log(`selected: ${Array.from(selected).join(",")}`);
}

export function isSelected(tileId: string): boolean {
  return selected.has(tileId);
}

export function clearSelected(): void {
  selected.clear();
}

export function getSelectedIds(): string[] {
  return Array.from(selected);
}

// Day 19: discardSelected validates exactly 1 selected and logs DISCARD {tileId}, no server call yet
export function discardSelected(): { tileId: string } | null {
  if (selected.size !== 1) {
    console.log(`discardSelected failed: selected.size ${selected.size} want 1`);
    return null;
  }
  const tileId = Array.from(selected)[0];
  console.log(`DISCARD {tileId: ${tileId}}`);
  return { tileId };
}

// Day 20: meldSelected validates selected.size>=3 and logs MELD_INITIAL or MELD_NEW, no server call yet
export function meldSelected(kind: string): { kind: string; tileIds: string[] } | null {
  if (selected.size < 3) {
    console.log(`meldSelected failed: selected.size ${selected.size} want >=3`);
    return null;
  }
  if (kind !== "run" && kind !== "set") {
    console.log(`meldSelected failed: kind ${kind} want run or set`);
    return null;
  }
  const tileIds = Array.from(selected);
  console.log(
    `MELD_${kind.toUpperCase()} {kind: ${kind}, tileIds: [${tileIds.join(",")}]} (no HasOpened check, no server call yet — Day 20)`
  );
  return { kind, tileIds };
}

function colourToHex(colour: number): number {
  switch (colour) {
    case 1:
      return 0xe53935;
    case 2:
      return 0xf9a825;
    case 3:
      return 0x1e88e5;
    case 4:
      return 0x212121;
    default:
      return 0x757575;
  }
}

function rankToLabel(rank: number): string {
  if (rank === 1) return "A";
  if (rank === 11) return "J";
  if (rank === 12) return "Q";
  if (rank === 13) return "K";
  return String(rank);
}

function createTileContainer(
  scene: Phaser.Scene,
  tl: TileInstance,
  x: number,
  y: number
): Phaser.GameObjects.Container {
  const container = scene.add.container(x, y);
  container.setData("isRackTile", true);
  container.setData("tileId", tl.ID);
  container.setSize(48, 64);
  const bg = scene.add
    .rectangle(0, 0, 48, 64, 0xffffff)
    .setStrokeStyle(2, colourToHex(tl.Colour), 1);
  bg.setData("isRackTile", true);
  container.add(bg);
  if (tl.IsJoker) {
    bg.setFillStyle(0xfff9c4);
    bg.setStrokeStyle(2, 0xf57f17, 1);
    const jokerText = scene.add.text(0, 0, "J", {
      fontFamily: "Inter, monospace",
      fontSize: "22px",
      color: "#f57f17",
      fontStyle: "bold",
    });
    jokerText.setOrigin(0.5);
    jokerText.setData("isRackTile", true);
    container.add(jokerText);
    const sub = scene.add.text(0, 14, "Joly", {
      fontFamily: "monospace",
      fontSize: "7px",
      color: "#795548",
    });
    sub.setOrigin(0.5);
    sub.setData("isRackTile", true);
    container.add(sub);
  } else {
    const rankLabel = rankToLabel(tl.Rank);
    const colourHex = colourToHex(tl.Colour);
    const colourStr = `#${colourHex.toString(16).padStart(6, "0")}`;
    const rankText = scene.add.text(0, -4, rankLabel, {
      fontFamily: "Inter, monospace",
      fontSize: "20px",
      color: colourStr,
      fontStyle: "bold",
    });
    rankText.setOrigin(0.5);
    rankText.setData("isRackTile", true);
    container.add(rankText);
    const suit = scene.add.text(0, 18, "●", {
      fontFamily: "monospace",
      fontSize: "10px",
      color: colourStr,
    });
    suit.setOrigin(0.5);
    suit.setData("isRackTile", true);
    container.add(suit);
    const idLabel = scene.add.text(0, 26, tl.ID.slice(0, 4), {
      fontFamily: "monospace",
      fontSize: "5px",
      color: "#999999",
    });
    idLabel.setOrigin(0.5);
    idLabel.setAlpha(0.6);
    idLabel.setData("isRackTile", true);
    container.add(idLabel);
  }
  // Modern selection: blue glow + checkmark, not yellow
  if (isSelected(tl.ID)) {
    bg.setStrokeStyle(2.5, 0x1e88e5, 1);
    bg.setFillStyle(0xe3f2fd);
    container.setScale(1.06);
    const glow = scene.add.rectangle(0, 0, 52, 68, 0x1e88e5, 0.08);
    glow.setData("isRackTile", true);
    container.addAt(glow, 1);
    const check = scene.add.text(14, -22, "✓", {
      fontFamily: "Inter, monospace",
      fontSize: "10px",
      color: "#1e88e5",
      fontStyle: "bold",
      backgroundColor: "#ffffff",
      padding: { x: 3, y: 1 },
    });
    check.setOrigin(0.5);
    check.setData("isRackTile", true);
    container.add(check);
  }
  container.setInteractive(
    new Phaser.Geom.Rectangle(-24, -32, 48, 64),
    Phaser.Geom.Rectangle.Contains
  );
  if ((scene.input as unknown as { setDraggable?: (obj: unknown) => void }).setDraggable) {
    (scene.input as unknown as { setDraggable: (obj: unknown) => void }).setDraggable(container);
  }
  container.on("pointerdown", () => {
    onTileClicked(tl.ID);
    // Visual update will happen on next render, but also do immediate
    if (isSelected(tl.ID)) {
      bg.setStrokeStyle(2.5, 0x1e88e5, 1);
      bg.setFillStyle(0xe3f2fd);
      container.setScale(1.06);
    } else {
      bg.setStrokeStyle(2, colourToHex(tl.Colour), 1);
      bg.setFillStyle(0xffffff);
      container.setScale(1);
    }
    console.log(`onTileClicked ${tl.ID} selected: ${isSelected(tl.ID)}`);
  });
  return container;
}

export function renderRack(
  scene: Phaser.Scene,
  tiles: TileInstance[],
  seat: number,
  opts?: { x?: number; y?: number; spacing?: number }
): Phaser.GameObjects.Container[] {
  const spacing = opts?.spacing ?? 62;
  const y = opts?.y ?? 700;
  const totalWidth = tiles.length > 0 ? (tiles.length - 1) * spacing : 0;
  const x0 = (opts?.x ?? 512) - totalWidth / 2;
  const existing = scene.children.list.filter((c) =>
    (c as unknown as { getData?: (k: string) => unknown }).getData?.("isRackTile")
  );
  for (const c of existing) c.destroy();
  const containers: Phaser.GameObjects.Container[] = [];
  for (let i = 0; i < tiles.length; i++) {
    const tl = tiles[i];
    const container = createTileContainer(scene, tl, x0 + i * spacing, y);
    container.setData("seat", seat);
    containers.push(container);
  }
  return containers;
}

export function sortRack(tiles: TileInstance[]): TileInstance[] {
  return [...tiles].sort((a, b) => {
    if (a.Colour !== b.Colour) return a.Colour - b.Colour;
    if (a.Rank !== b.Rank) return a.Rank - b.Rank;
    return a.ID.localeCompare(b.ID);
  });
}
