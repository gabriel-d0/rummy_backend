import Phaser from "phaser";
import type { TileInstance } from "../state/snapshot";

// Day 11: Rack rendering (static) — draws 14 this.add.image per tile at x = 100 + i*50
// Mirrors Go TileInstance{ID,Colour,Rank,IsJoker} and PrivateView.OwnRack (redaction)

export function renderRack(
  scene: Phaser.Scene,
  tiles: TileInstance[],
  seat: number,
  opts?: { x?: number; y?: number; spacing?: number },
): Phaser.GameObjects.Image[] {
  const x0 = opts?.x ?? 100;
  const y = opts?.y ?? 680;
  const spacing = opts?.spacing ?? 50;
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
    // Day 12 will add sortRack; for now we just render in given order
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
