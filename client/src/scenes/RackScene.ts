import Phaser from "phaser";
import { discardSelected, renderRack, sortRack } from "../ui/Rack";
import { getLayout } from "../ui/Layout";

// Day 11 + Day 19: RackScene renders PrivateView.OwnRack only (redaction) with subspace layout
export class RackScene extends Phaser.Scene {
  constructor() {
    super("RackScene");
  }

  create() {
    const layout = getLayout(this.scale.width, this.scale.height);

    // Rack background fills rack subspace (wood 800x120 or full width on mobile, centered)
    this.add.image(layout.rack.x + layout.rack.w / 2, layout.rack.y + layout.rack.h / 2, "rack").setDisplaySize(layout.rack.w, layout.rack.h);

    // Draw slot outlines for visual clarity (14 slots)
    for (const slot of layout.rackSlots) {
      this.add.rectangle(slot.x + slot.w / 2, slot.y + slot.h / 2, slot.w, slot.h, 0x3d2817, 0).setStrokeStyle(1, 0x5a3d1a, 0.5);
    }

    // Day 12 + Day 16: sorted rack with selection
    const unsorted = [
      { ID: "mock-red-13", Colour: 1, Rank: 13, IsJoker: false },
      { ID: "mock-red-1", Colour: 1, Rank: 1, IsJoker: false },
      { ID: "mock-blue-5", Colour: 3, Rank: 5, IsJoker: false },
    ];
    const mockRack = sortRack(unsorted);
    // Render centered within rack subspace, not at fixed 100,700
    const totalW = mockRack.length > 0 ? (mockRack.length - 1) * 62 : 0;
    const rackCenterX = layout.rack.x + layout.rack.w / 2;
    const startX = rackCenterX - totalW / 2;
    renderRack(this, mockRack, 0, { x: startX, y: layout.rack.y + layout.rack.h / 2, spacing: 62 });

    // Day 17: dragstart
    this.input.on("dragstart", (_pointer: any, gameObject: any) => {
      const tileId = gameObject.getData("tileId");
      console.log(`dragstart ${tileId}`);
    });

    // Day 19: Discard button inside rack subspace top-right, not overlapping
    const btnX = layout.rack.x + layout.rack.w - 50;
    const btnY = layout.rack.y - 14;
    const discardBtn = this.add
      .text(btnX, btnY, "[Discard]", {
        fontFamily: "monospace",
        fontSize: "12px",
        color: "#00ff00",
        backgroundColor: "#1a3d2e",
        padding: { x: 8, y: 4 },
      })
      .setOrigin(0.5)
      .setInteractive({ useHandCursor: true });
    discardBtn.on("pointerdown", () => {
      const res = discardSelected();
      if (res) {
        console.log(`discardSelected success: ${res.tileId}`);
      }
    });

    this.add
      .text(layout.rack.x + layout.rack.w / 2, layout.info.y + 8, "Rack — Day 19 discardSelected + dragstart (click tile, then [Discard])", {
        fontFamily: "monospace",
        fontSize: "10px",
        color: "#ffff00",
        backgroundColor: "#00000066",
        padding: { x: 6, y: 2 },
      })
      .setOrigin(0.5);

    this.scale.on("resize", () => {
      this.scene.restart();
    });
  }
}
