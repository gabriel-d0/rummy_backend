import Phaser from "phaser";

// Day 13: TableScene will render PublicView.TableMelds/DiscardRow/StockCount/CurrentSeat
// Day 15: Stock and turn indicator
export class TableScene extends Phaser.Scene {
  constructor() {
    super("TableScene");
  }

  create() {
    // Day 8: table background (green felt 1024x768) behind tiles
    this.add.image(512, 384, "table").setDisplaySize(1024, 768);
    // Day 7: render one tile to prove asset pipeline
    this.add.image(100, 100, "tile").setScale(2);
    this.add.image(140, 100, "joker").setScale(2);
    this.add
      .text(512, 384, "TableScene — Day 8 table background\n(see client/docs/roadmap.md Phase 2)", {
        fontFamily: "monospace",
        fontSize: "14px",
        color: "#ffffff",
        align: "center",
      })
      .setOrigin(0.5);
    // Day 13-15 will add renderTableMelds, renderDiscardRow, renderStockCount, TurnIndicator
  }
}
