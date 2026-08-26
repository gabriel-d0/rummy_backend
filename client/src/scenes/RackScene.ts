import Phaser from "phaser";
import { renderRack, sortRack } from "../ui/Rack";

// Day 11: RackScene renders PrivateView.OwnRack only (redaction) — 14 this.add.image per tile at x = 100 + i*50
export class RackScene extends Phaser.Scene {
  constructor() {
    super("RackScene");
  }

  create() {
    // Day 9: rack background (wood 800x120 with 14 slots)
    this.add.image(512, 680, "rack");
    // Day 12: sortRack by Colour then Rank — mock red-13, red-1, blue-5 → red-1, red-13, blue-5
    const unsorted = [
      { ID: "mock-red-13", Colour: 1, Rank: 13, IsJoker: false },
      { ID: "mock-red-1", Colour: 1, Rank: 1, IsJoker: false },
      { ID: "mock-blue-5", Colour: 3, Rank: 5, IsJoker: false },
    ];
    const mockRack = sortRack(unsorted);
    renderRack(this, mockRack, 0);
    this.add
      .text(512, 620, "RackScene — Day 12 sortRack(red-13,red-1,blue-5 → red-1,red-13,blue-5)", {
        fontFamily: "monospace",
        fontSize: "10px",
        color: "#ffff00",
        align: "center",
      })
      .setOrigin(0.5);
    // Day 16-20: onTileClicked, dragstart, discardSelected, meldSelected
  }
}
