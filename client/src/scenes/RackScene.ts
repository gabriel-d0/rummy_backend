import Phaser from "phaser";
import { renderRack } from "../ui/Rack";

// Day 11: RackScene renders PrivateView.OwnRack only (redaction) — 14 this.add.image per tile at x = 100 + i*50
export class RackScene extends Phaser.Scene {
  constructor() {
    super("RackScene");
  }

  create() {
    // Day 9: rack background (wood 800x120 with 14 slots)
    this.add.image(512, 680, "rack");
    // Day 11: renderRack with mock PrivateView.OwnRack (2 tiles red-5, blue-7) — npm run dev shows rack
    const mockRack = [
      { ID: "mock-red-5", Colour: 1, Rank: 5, IsJoker: false },
      { ID: "mock-blue-7", Colour: 3, Rank: 7, IsJoker: false },
    ];
    renderRack(this, mockRack, 0);
    this.add
      .text(512, 620, "RackScene — Day 11 renderRack(2 tiles red-5, blue-7)", {
        fontFamily: "monospace",
        fontSize: "10px",
        color: "#ffff00",
        align: "center",
      })
      .setOrigin(0.5);
    // Day 12: sortRack by Colour then Rank will be used here
    // Day 16-20: onTileClicked, dragstart, discardSelected, meldSelected
  }
}
