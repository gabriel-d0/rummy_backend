import Phaser from "phaser";

// Day 13: TableScene will render PublicView.TableMelds/DiscardRow/StockCount/CurrentSeat
// Day 15: Stock and turn indicator
export class TableScene extends Phaser.Scene {
  constructor() {
    super("TableScene");
  }

  create() {
    this.add.rectangle(512, 384, 1024, 768, 0x0a4d2e);
    this.add
      .text(512, 384, "TableScene — Day 13\n(see client/docs/roadmap.md Phase 3)", {
        fontFamily: "monospace",
        fontSize: "14px",
        color: "#ffffff",
        align: "center",
      })
      .setOrigin(0.5);
    // Day 13-15 will add renderTableMelds, renderDiscardRow, renderStockCount, TurnIndicator
  }
}
