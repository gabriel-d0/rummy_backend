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

export function renderRack(
  scene: Phaser.Scene,
  tiles: TileInstance[],
  seat: number,
  opts?: { x?: number; y?: number; spacing?: number },
): Phaser.GameObjects.Image[] {
  const spacing = opts?.spacing ?? 54;
  const y = opts?.y ?? 700;
  // Center the rack within the wood rack image (800x120 at 512,680)
  // Rack left edge is at 512-400=112, but we center tiles at 512
  const totalWidth = tiles.length > 0 ? (tiles.length - 1) * spacing : 0;
  const x0 = (opts?.x ?? 512) - totalWidth / 2;
  // Clear previous rack images if any (tagged with "rack-tile")
  const existing = scene.children.list.filter((c) => (c as any).getData?.("isRackTile"));
  for (const c of existing) c.destroy();

  const images: Phaser.GameObjects.Image[] = [];
  for (let i = 0; i < tiles.length; i++) {
    const tl = tiles[i];
    const key = tl.IsJoker ? "joker" : "tile";
    const img = scene.add.image(x0 + i * spacing, y, key).setScale(0.9);
    img.setData("isRackTile", true);
    img.setData("tileId", tl.ID);
    img.setData("seat", seat);
    img.setInteractive({ useHandCursor: true });
    // Day 17: make draggable and log dragstart (no drop yet)
    if ((scene.input as any).setDraggable) {
      scene.input.setDraggable(img);
    }
    // Day 16: tint if selected
    if (isSelected(tl.ID)) {
      img.setTint(0xffff00);
    } else {
      img.clearTint();
    }
    img.on("pointerdown", () => {
      onTileClicked(tl.ID);
      // Retint after toggle
      if (isSelected(tl.ID)) {
        img.setTint(0xffff00);
      } else {
        img.clearTint();
      }
      console.log(`onTileClicked ${tl.ID} selected: ${isSelected(tl.ID)}`);
    });
    // Add a small text label for debugging (tile ID short)
    const label = scene.add.text(x0 + i * spacing, y + 36, tl.IsJoker ? "J" : `${tl.Colour}-${tl.Rank}`, {
      fontFamily: "monospace",
      fontSize: "8px",
      color: "#ffffff",
      align: "center",
    });
    label.setOrigin(0.5);
    label.setData("isRackTile", true);
    images.push(img);
  }
  return images;
}

export function sortRack(tiles: TileInstance[]): TileInstance[] {
  // Day 12: sort by Colour then Rank (for later)
  return [...tiles].sort((a, b) => {
    if (a.Colour !== b.Colour) return a.Colour - b.Colour;
    if (a.Rank !== b.Rank) return a.Rank - b.Rank;
    return a.ID.localeCompare(b.ID);
  });
}
