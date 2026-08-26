import Phaser from "phaser";

// Day 11-12: RackScene will render PrivateView.OwnRack only (redaction)
export class RackScene extends Phaser.Scene {
  constructor() {
    super("RackScene");
  }

  create() {
    // Day 9: rack background (wood 800x120 with 14 slots)
    this.add.image(512, 680, "rack");
    this.add
      .text(512, 680, "RackScene — Day 9 rack background\n(PrivateView.OwnRack only)", {
        fontFamily: "monospace",
        fontSize: "12px",
        color: "#ffffff",
        align: "center",
      })
      .setOrigin(0.5);
    // Day 11: renderRack(tiles: TileInstance[], seat: Seat)
    // Day 12: sortRack by Colour then Rank
    // Day 16-20: onTileClicked, dragstart, discardSelected, meldSelected
  }
}
