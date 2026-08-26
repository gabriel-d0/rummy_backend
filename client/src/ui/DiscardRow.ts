import Phaser from "phaser";
import type { DiscardEntry } from "../state/snapshot";

export function renderDiscardRow(
  scene: Phaser.Scene,
  discardRow: DiscardEntry[],
  opts?: { x?: number; y?: number; spacing?: number }
): Phaser.GameObjects.Container[] {
  const x0 = opts?.x ?? 80;
  const y = opts?.y ?? 280;
  const spacing = opts?.spacing ?? 50;
  const existing = scene.children.list.filter((c) =>
    (c as unknown as { getData?: (k: string) => unknown }).getData?.("isDiscardRow")
  );
  for (const c of existing) c.destroy();
  const containers: Phaser.GameObjects.Container[] = [];
  for (let i = 0; i < discardRow.length; i++) {
    const entry = discardRow[i];
    const x = x0 + i * spacing;
    const container = scene.add.container(x, y);
    container.setData("isDiscardRow", true);
    container.setData("discardIndex", i);
    container.setData("isOpeningDiscard", entry.IsOpeningDiscard);
    const tileW = 36;
    const tileH = 48;
    const bgColor = entry.Tile.IsJoker ? 0xfff9c4 : 0xffffff;
    const strokeColor = entry.Tile.IsJoker
      ? 0xf57f17
      : (() => {
          switch (entry.Tile.Colour) {
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
        })();
    const bg = scene.add
      .rectangle(0, 0, tileW, tileH, bgColor)
      .setStrokeStyle(
        entry.IsOpeningDiscard ? 2 : 1.5,
        entry.IsOpeningDiscard ? 0xff0000 : strokeColor,
        1
      );
    bg.setData("isDiscardRow", true);
    container.add(bg);
    const rankLabel = entry.Tile.IsJoker
      ? "J"
      : (() => {
          if (entry.Tile.Rank === 1) return "A";
          if (entry.Tile.Rank === 11) return "J";
          if (entry.Tile.Rank === 12) return "Q";
          if (entry.Tile.Rank === 13) return "K";
          return String(entry.Tile.Rank);
        })();
    const colourStr = entry.Tile.IsJoker
      ? "#f57f17"
      : (() => {
          switch (entry.Tile.Colour) {
            case 1:
              return "#e53935";
            case 2:
              return "#f9a825";
            case 3:
              return "#1e88e5";
            case 4:
              return "#212121";
            default:
              return "#757575";
          }
        })();
    const rankTxt = scene.add.text(0, -2, rankLabel, {
      fontFamily: "Inter, monospace",
      fontSize: "14px",
      color: colourStr,
      fontStyle: "bold",
    });
    rankTxt.setOrigin(0.5);
    rankTxt.setData("isDiscardRow", true);
    container.add(rankTxt);
    const label = scene.add.text(0, 18, entry.IsOpeningDiscard ? "open" : "", {
      fontFamily: "monospace",
      fontSize: "6px",
      color: "#ff0000",
      align: "center",
    });
    label.setOrigin(0.5);
    label.setData("isDiscardRow", true);
    container.add(label);
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
