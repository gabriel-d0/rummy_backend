import Phaser from "phaser";

// Day 15: Stock and turn indicator — draws StockCount and CurrentSeat/GamePhase/TurnPhase at x=800 y=50

export function renderStockCount(
  scene: Phaser.Scene,
  count: number,
  opts?: { x?: number; y?: number },
): Phaser.GameObjects.Text {
  const x = opts?.x ?? 880;
  const y = opts?.y ?? 40;
  // Clear previous
  const existing = scene.children.list.filter((c) => (c as any).getData?.("isStockCount"));
  for (const c of existing) c.destroy();

  const text = scene.add.text(x, y, `Stock: ${count}`, {
    fontFamily: "monospace",
    fontSize: "14px",
    color: "#ffffff",
    backgroundColor: "#000000aa",
    padding: { x: 8, y: 4 },
  });
  text.setOrigin(0.5);
  text.setData("isStockCount", true);
  // Add a small stock pile visualization (stacked rectangles) below text, not overlapping turn indicator
  for (let i = 0; i < Math.min(3, Math.ceil(count / 20)); i++) {
    const pile = scene.add.rectangle(x, y + 28 + i * 4, 44, 22, 0x4444aa).setStrokeStyle(1, 0x222266);
    pile.setData("isStockCount", true);
  }
  return text;
}

export function renderTurnIndicator(
  scene: Phaser.Scene,
  currentSeat: number,
  gamePhase: string,
  turnPhase: string,
  opts?: { x?: number; y?: number },
): Phaser.GameObjects.Text {
  const x = opts?.x ?? 880;
  const y = opts?.y ?? 105;
  const existing = scene.children.list.filter((c) => (c as any).getData?.("isTurnIndicator"));
  for (const c of existing) c.destroy();

  const text = scene.add.text(x, y, `Current: seat-${currentSeat}\n${gamePhase}/${turnPhase}`, {
    fontFamily: "monospace",
    fontSize: "12px",
    color: "#ffff00",
    backgroundColor: "#00000066",
    padding: { x: 8, y: 4 },
    align: "center",
  });
  text.setOrigin(0.5);
  text.setData("isTurnIndicator", true);
  return text;
}
