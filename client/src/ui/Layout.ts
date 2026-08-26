import Phaser from "phaser";

// Day 20+: Subspace layout — divides the entire game space into
// mathematically defined subspaces so every component knows its bounds
// and the design is responsive (desktop/tablet/mobile).

export interface Rect {
  x: number;
  y: number;
  w: number;
  h: number;
}

export interface GameLayout {
  width: number;
  height: number;
  isMobile: boolean;
  outer: Rect; // outer border
  topBar: Rect; // Stock + Turn
  tableMelds: Rect; // melds area
  discardRow: Rect; // discard area
  dropZone: Rect; // drop zone
  rack: Rect; // rack background
  rackSlots: Rect[]; // 14 slots inside rack
  info: Rect; // bottom info text
}

export function getLayout(width: number, height: number): GameLayout {
  const isMobile = width < 768;
  const outerMargin = isMobile ? 8 : 12;
  const gutter = isMobile ? 8 : 16;

  const outer: Rect = {
    x: outerMargin,
    y: outerMargin,
    w: width - outerMargin * 2,
    h: height - outerMargin * 2,
  };

  // Top bar: Stock (left) + Turn (right) in top bar area
  const topBarH = 70;
  const topBar: Rect = {
    x: outer.x + gutter,
    y: outer.y + gutter,
    w: outer.w - gutter * 2,
    h: topBarH,
  };

  // Rack at bottom: 800x120 centered, or full width on mobile
  const rackW = isMobile ? outer.w - gutter * 2 : 800;
  const rackH = isMobile ? 100 : 120;
  const rack: Rect = {
    x: width / 2 - rackW / 2,
    y: height - outerMargin - gutter - rackH,
    w: rackW,
    h: rackH,
  };

  // Info text below drop zone, above rack
  const info: Rect = {
    x: outer.x,
    y: rack.y - 30,
    w: outer.w,
    h: 20,
  };

  // Drop zone: centered, 600x50 on desktop, full width minus gutters on mobile, y = above info
  const dropW = isMobile ? outer.w - gutter * 2 : 600;
  const dropH = isMobile ? 44 : 50;
  const dropZone: Rect = {
    x: width / 2 - dropW / 2,
    y: info.y - gutter - dropH,
    w: dropW,
    h: dropH,
  };

  // Discard row: above drop zone
  const discardH = 70;
  const discardRow: Rect = {
    x: outer.x + gutter,
    y: dropZone.y - gutter - discardH,
    w: outer.w - gutter * 2,
    h: discardH,
  };

  // Table melds: between topBar and discardRow
  const tableMelds: Rect = {
    x: outer.x + gutter,
    y: topBar.y + topBar.h + gutter,
    w: outer.w - gutter * 2,
    h: discardRow.y - (topBar.y + topBar.h + gutter) - gutter,
  };

  return {
    width,
    height,
    isMobile,
    outer,
    topBar,
    tableMelds,
    discardRow,
    dropZone,
    rack,
    rackSlots: getRackSlots(rack),
    info,
  };
}

function getRackSlots(rack: Rect): Rect[] {
  const slots: Rect[] = [];
  const slotW = 48;
  const slotH = 80;
  const n = 14;
  const totalW = n * slotW + (n - 1) * 6; // 6px gap
  const startX = rack.x + (rack.w - totalW) / 2;
  const y = rack.y + (rack.h - slotH) / 2;
  for (let i = 0; i < n; i++) {
    slots.push({
      x: startX + i * (slotW + 6),
      y,
      w: slotW,
      h: slotH,
    });
  }
  return slots;
}

export function drawDebugRect(
  scene: Phaser.Scene,
  r: Rect,
  color: number,
  alpha: number = 0.08,
  tag?: string,
): Phaser.GameObjects.Rectangle {
  const rect = scene.add.rectangle(r.x + r.w / 2, r.y + r.h / 2, r.w, r.h, color, alpha).setStrokeStyle(1, color, 0.3);
  if (tag) rect.setData(tag, true);
  return rect;
}

export function onResize(scene: Phaser.Scene, cb: (layout: GameLayout) => void) {
  scene.scale.on("resize", (gameSize: Phaser.Structs.Size) => {
    const layout = getLayout(gameSize.width, gameSize.height);
    cb(layout);
  });
}
