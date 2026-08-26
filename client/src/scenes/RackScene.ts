import Phaser from "phaser";
import { discardSelected, meldSelected, renderRack, sortRack } from "../ui/Rack";
import { getLayout } from "../ui/Layout";
import { subscribePrivateSnapshot } from "../state/sync";
import type { PrivateSnapshot } from "../state/snapshot";

// Day 11 + Day 19: RackScene renders PrivateView.OwnRack only (redaction) with subspace layout
// Day 30: subscribes to PrivateSnapshot and re-renders OwnRack only (never foreign rack)
export class RackScene extends Phaser.Scene {
  constructor() {
    super("RackScene");
  }

  private renderPrivateRack(snap: PrivateSnapshot): void {
    const layout = getLayout(this.scale.width, this.scale.height);
    const sorted = sortRack(snap.ownRack);
    const totalW = sorted.length > 0 ? (sorted.length - 1) * 62 : 0;
    const rackCenterX = layout.rack.x + layout.rack.w / 2;
    const startX = rackCenterX - totalW / 2;
    renderRack(this, sorted, snap.ownSeat, {
      x: startX,
      y: layout.rack.y + layout.rack.h / 2,
      spacing: 62,
    });
  }

  create() {
    const layout = getLayout(this.scale.width, this.scale.height);

    // Rack background fills rack subspace (wood 800x120 or full width on mobile, centered)
    this.add
      .image(layout.rack.x + layout.rack.w / 2, layout.rack.y + layout.rack.h / 2, "rack")
      .setDisplaySize(layout.rack.w, layout.rack.h);

    // Draw slot outlines for visual clarity (14 slots)
    for (const slot of layout.rackSlots) {
      this.add
        .rectangle(slot.x + slot.w / 2, slot.y + slot.h / 2, slot.w, slot.h, 0x3d2817, 0)
        .setStrokeStyle(1, 0x5a3d1a, 0.5);
    }

    // Day 12 + Day 16: sorted rack with selection (mock for initial render before any PrivateSnapshot)
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

    // Day 30: subscribe to PrivateSnapshot — re-render OwnRack only (redaction)
    const unsubscribe = subscribePrivateSnapshot((snap) => {
      this.renderPrivateRack(snap);
    });
    this.events.once("shutdown", unsubscribe);
    this.events.once("destroy", unsubscribe);

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
    // Day 20: Meld buttons — validate selected.size>=3 and log MELD_INITIAL or MELD_NEW
    const meldRunBtn = this.add
      .text(btnX - 110, btnY, "[Meld Run]", {
        fontFamily: "monospace",
        fontSize: "12px",
        color: "#00ff00",
        backgroundColor: "#1a3d2e",
        padding: { x: 8, y: 4 },
      })
      .setOrigin(0.5)
      .setInteractive({ useHandCursor: true });
    meldRunBtn.on("pointerdown", () => {
      const res = meldSelected("run");
      if (res) console.log(`meldSelected success: run ${res.tileIds.join(",")}`);
    });
    const meldSetBtn = this.add
      .text(btnX - 210, btnY, "[Meld Set]", {
        fontFamily: "monospace",
        fontSize: "12px",
        color: "#00ff00",
        backgroundColor: "#1a3d2e",
        padding: { x: 8, y: 4 },
      })
      .setOrigin(0.5)
      .setInteractive({ useHandCursor: true });
    meldSetBtn.on("pointerdown", () => {
      const res = meldSelected("set");
      if (res) console.log(`meldSelected success: set ${res.tileIds.join(",")}`);
    });

    this.add
      .text(
        layout.rack.x + layout.rack.w / 2,
        layout.info.y + 8,
        "Rack — Day 20 meldSelected + Day 19 discardSelected",
        {
          fontFamily: "monospace",
          fontSize: "10px",
          color: "#ffff00",
          backgroundColor: "#00000066",
          padding: { x: 6, y: 2 },
        }
      )
      .setOrigin(0.5);

    this.scale.on("resize", () => {
      this.scene.restart();
    });
  }
}
