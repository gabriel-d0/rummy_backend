import Phaser from "phaser";
import { discardSelected, renderRack, sortRack } from "../ui/Rack";

// Day 11: RackScene renders PrivateView.OwnRack only (redaction) — 14 this.add.image per tile at x = 100 + i*50
export class RackScene extends Phaser.Scene {
  constructor() {
    super("RackScene");
  }

  create() {
    // Day 9: rack background (wood 800x120 with 14 slots)
    this.add.image(512, 680, "rack");
    // Day 16: onTileClicked toggles selected Set and tints 0xffff00 — mock red-13, red-1, blue-5 → red-1, red-13, blue-5
    const unsorted = [
      { ID: "mock-red-13", Colour: 1, Rank: 13, IsJoker: false },
      { ID: "mock-red-1", Colour: 1, Rank: 1, IsJoker: false },
      { ID: "mock-blue-5", Colour: 3, Rank: 5, IsJoker: false },
    ];
    const mockRack = sortRack(unsorted);
    renderRack(this, mockRack, 0);
    // Day 17: dragstart logs tileId, no drop yet
    this.input.on("dragstart", (_pointer: any, gameObject: any) => {
      const tileId = gameObject.getData("tileId");
      console.log(`dragstart ${tileId}`);
    });
    // Day 19: discardSelected validates exactly 1 selected and logs DISCARD {tileId}, no server call yet
    const discardBtn = this.add
      .text(900, 620, "[Discard]", {
        fontFamily: "monospace",
        fontSize: "12px",
        color: "#00ff00",
        backgroundColor: "#333333",
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
      .text(512, 620, "RackScene — Day 19 discardSelected(1) + Day 17 dragstart", {
        fontFamily: "monospace",
        fontSize: "10px",
        color: "#ffff00",
        align: "center",
      })
      .setOrigin(0.5);
    // Day 20: meldSelected
  }
}
