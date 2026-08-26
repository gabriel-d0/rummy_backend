import Phaser from "phaser";
import type { DiscardEntry } from "../state/snapshot";

// Day 14: Discard row rendering (static) — draws DiscardRow at x = 100 + i*40 y = 300, mock with IsOpeningDiscard flagged
export function renderDiscardRow(
  scene: Phaser.Scene,
  discardRow: DiscardEntry[],
  opts?: { x?: number; y?: number; spacing?: number },
): Phaser.GameObjects.Container[] {
  const x0 = opts?.x ?? 100;
  const y = opts?.y ?? 300;
  const spacing = opts?.spacing ?? 40;

  // Clear previous discard row (tagged with "discard-row")
  const existing = scene.children.list.filter((c) => (c as any).getData?.("isDiscardRow"));
  for (const c of existing) c.destroy();

  const containers: Phaser.GameObjects.Container[] = [];
  for (let i = 0; i < discardRow.length; i++) {
    const entry = discardRow[i];
    const x = x0 + i * spacing;
    const container = scene.add.container(x, y);
    container.setData("isDiscardRow", true);
    container.setData("discardIndex", i);
    container.setData("isOpeningDiscard", entry.IsOpeningDiscard);

    // Tile image
    const key = entry.Tile.IsJoker ? "joker" : "tile";
    const img = scene.add.image(0, 0, key).setScale(0.6);
    img.setData("isDiscardRow", true);
    // Red border for opening discard
    if (entry.IsOpeningDiscard) {
      const border = scene.add.rectangle(0, 0, 38, 52, 0xff0000, 0).setStrokeStyle(2, 0xff0000);
      border.setData("isDiscardRow", true);
      container.add(border);
    }
    container.add(img);

    // Label: tile ID and index
    const label = scene.add.text(0, 22, entry.IsOpeningDiscard ? `disc-open` : `${entry.Tile.Colour}-${entry.Tile.Rank}`, {
      fontFamily: "monospace",
      fontSize: "6px",
      color: entry.IsOpeningDiscard ? "#ff0000" : "#ffffff",
      align: "center",
    });
    label.setOrigin(0.5);
    label.setData("isDiscardRow", true);
    container.add(label);

    // Index label
    const idxLabel = scene.add.text(0, -22, `#${i}`, {
      fontFamily: "monospace",
      fontSize: "6px",
      color: "#ffff00",
      align: "center",
    });
    idxLabel.setOrigin(0.5);
    idxLabel.setData("isDiscardRow", true);
    container.add(idxLabel);

    containers.push(container);
  }
  return containers;
}
