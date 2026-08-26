import Phaser from "phaser";

// Day 11-12: RackScene will render PrivateView.OwnRack only (redaction)
export class RackScene extends Phaser.Scene {
  constructor() {
    super("RackScene");
  }

  create() {
    this.add.rectangle(512, 680, 800, 120, 0x5a2a0a).setStrokeStyle(2, 0x3e1f00);
    this.add.text(512, 680, "RackScene — Day 11\n(PrivateView.OwnRack only)", {
      fontFamily: "monospace",
      fontSize: "12px",
      color: "#ffffff",
      align: "center",
    }).setOrigin(0.5);
    // Day 11: renderRack(tiles: TileInstance[], seat: Seat)
    // Day 12: sortRack by Colour then Rank
    // Day 16-20: onTileClicked, dragstart, discardSelected, meldSelected
  }
}
