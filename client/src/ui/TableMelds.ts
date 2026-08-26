import Phaser from "phaser";
import type { TableMeld } from "../state/snapshot";

export function renderTableMelds(
  scene: Phaser.Scene,
  melds: TableMeld[],
  opts?: { x?: number; y0?: number; rowHeight?: number; tileSpacing?: number }
): Phaser.GameObjects.Container[] {
  const x0 = opts?.x ?? 80;
  const y0 = opts?.y0 ?? 110;
  const rowHeight = opts?.rowHeight ?? 70;
  const tileSpacing = opts?.tileSpacing ?? 50;
  const existing = scene.children.list.filter((c) =>
    (c as unknown as { getData?: (k: string) => unknown }).getData?.("isTableMeld")
  );
  for (const c of existing) c.destroy();
  const containers: Phaser.GameObjects.Container[] = [];
  for (let i = 0; i < melds.length; i++) {
    const meld = melds[i];
    const y = y0 + i * rowHeight;
    const container = scene.add.container(x0, y);
    container.setData("isTableMeld", true);
    container.setData("meldId", meld.ID);
    const bgWidth = meld.Tiles.length * tileSpacing + 20;
    const bg = scene.add.rectangle(0, 0, bgWidth, 50, 0x000000, 0.3).setOrigin(0, 0.5);
    bg.setData("isTableMeld", true);
    container.add(bg);
    const title = scene.add.text(0, -20, `${meld.ID} ${meld.Kind}`, {
      fontFamily: "monospace",
      fontSize: "8px",
      color: "#ffff00",
    });
    title.setData("isTableMeld", true);
    container.add(title);
    for (let j = 0; j < meld.Tiles.length; j++) {
      const tl = meld.Tiles[j];
      const tx = 10 + j * tileSpacing;
      const tileW = 36;
      const tileH = 48;
      const bg2 = scene.add.rectangle(tx, 0, tileW, tileH, 0xffffff).setStrokeStyle(
        1.5,
        tl.IsJoker
          ? 0xf57f17
          : (() => {
              switch (tl.Colour) {
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
            })(),
        1
      );
      if (tl.IsJoker) bg2.setFillStyle(0xfff9c4);
      bg2.setData("isTableMeld", true);
      bg2.setData("tileId", tl.ID);
      container.add(bg2);
      const rankLabel = tl.IsJoker
        ? "J"
        : (() => {
            if (tl.Rank === 1) return "A";
            if (tl.Rank === 11) return "J";
            if (tl.Rank === 12) return "Q";
            if (tl.Rank === 13) return "K";
            return String(tl.Rank);
          })();
      const colourHex = (() => {
        switch (tl.Colour) {
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
      const txt = scene.add.text(tx, -2, rankLabel, {
        fontFamily: "Inter, monospace",
        fontSize: "14px",
        color: tl.IsJoker ? "#f57f17" : colourHex,
        fontStyle: "bold",
      });
      txt.setOrigin(0.5);
      txt.setData("isTableMeld", true);
      txt.setData("tileId", tl.ID);
      container.add(txt);
      if (!tl.IsJoker) {
        const dot = scene.add.text(tx, 12, "●", {
          fontFamily: "monospace",
          fontSize: "7px",
          color: colourHex,
        });
        dot.setOrigin(0.5);
        dot.setData("isTableMeld", true);
        container.add(dot);
      }
    }
    containers.push(container);
  }
  return containers;
}
