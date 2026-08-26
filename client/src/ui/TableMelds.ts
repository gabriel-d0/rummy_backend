import Phaser from "phaser";
import type { TableMeld } from "../state/snapshot";

// Day 13: Table melds rendering (static) — draws each TableMeld at y = 100 + i*80 with Tiles as this.add.image per tile
export function renderTableMelds(
  scene: Phaser.Scene,
  melds: TableMeld[],
  opts?: { x?: number; y0?: number; rowHeight?: number; tileSpacing?: number },
): Phaser.GameObjects.Container[] {
  const x0 = opts?.x ?? 100;
  const y0 = opts?.y0 ?? 100;
  const rowHeight = opts?.rowHeight ?? 80;
  const tileSpacing = opts?.tileSpacing ?? 40;

  // Clear previous melds (tagged with "table-meld")
  const existing = scene.children.list.filter((c) => (c as any).getData?.("isTableMeld"));
  for (const c of existing) c.destroy();

  const containers: Phaser.GameObjects.Container[] = [];
  for (let i = 0; i < melds.length; i++) {
    const meld = melds[i];
    const y = y0 + i * rowHeight;
    // Create a container for the meld
    const container = scene.add.container(x0, y);
    container.setData("isTableMeld", true);
    container.setData("meldId", meld.ID);

    // Background for meld
    const bgWidth = meld.Tiles.length * tileSpacing + 20;
    const bg = scene.add.rectangle(0, 0, bgWidth, 50, 0x000000, 0.3).setOrigin(0, 0.5);
    bg.setData("isTableMeld", true);
    container.add(bg);

    // Title: meld ID and kind
    const title = scene.add.text(0, -20, `${meld.ID} ${meld.Kind}`, {
      fontFamily: "monospace",
      fontSize: "8px",
      color: "#ffff00",
    });
    title.setData("isTableMeld", true);
    container.add(title);

    // Tiles
    for (let j = 0; j < meld.Tiles.length; j++) {
      const tl = meld.Tiles[j];
      const key = tl.IsJoker ? "joker" : "tile";
      const img = scene.add.image(10 + j * tileSpacing, 0, key).setScale(0.6);
      img.setData("isTableMeld", true);
      img.setData("tileId", tl.ID);
      img.setData("meldId", meld.ID);
      // Small label for debugging
      const label = scene.add.text(10 + j * tileSpacing, 18, tl.IsJoker ? "J" : `${tl.Colour}-${tl.Rank}`, {
        fontFamily: "monospace",
        fontSize: "6px",
        color: "#ffffff",
        align: "center",
      });
      label.setOrigin(0.5);
      label.setData("isTableMeld", true);
      container.add(img);
      container.add(label);
    }
    containers.push(container);
  }
  return containers;
}
