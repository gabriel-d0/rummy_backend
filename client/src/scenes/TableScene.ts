import Phaser from "phaser";
import { renderTableMelds } from "../ui/TableMelds";

// Day 13: TableScene renders PublicView.TableMelds/DiscardRow/StockCount/CurrentSeat
export class TableScene extends Phaser.Scene {
  constructor() {
    super("TableScene");
  }

  create() {
    // Day 8: table background (green felt 1024x768) behind tiles
    this.add.image(512, 384, "table").setDisplaySize(1024, 768);
    // Day 13: renderTableMelds with mock PublicView.TableMelds (1 run 1-2-3 red and 1 set 7 red/yellow/blue)
    const mockMelds = [
      {
        ID: "mock-run-1-2-3",
        Kind: "run",
        Tiles: [
          { ID: "m1-t1", Colour: 1, Rank: 1, IsJoker: false },
          { ID: "m1-t2", Colour: 1, Rank: 2, IsJoker: false },
          { ID: "m1-t3", Colour: 1, Rank: 3, IsJoker: false },
        ],
        JokerReps: {},
        OwnerSeat: 0,
      },
      {
        ID: "mock-set-7",
        Kind: "set",
        Tiles: [
          { ID: "m2-t1", Colour: 1, Rank: 7, IsJoker: false },
          { ID: "m2-t2", Colour: 2, Rank: 7, IsJoker: false },
          { ID: "m2-t3", Colour: 3, Rank: 7, IsJoker: false },
        ],
        JokerReps: {},
        OwnerSeat: 1,
      },
    ];
    renderTableMelds(this, mockMelds, { x: 100, y0: 160, rowHeight: 80, tileSpacing: 40 });
    this.add
      .text(512, 320, "TableScene — Day 13 renderTableMelds(1 run 1-2-3 red, 1 set 7)", {
        fontFamily: "monospace",
        fontSize: "10px",
        color: "#ffff00",
        align: "center",
      })
      .setOrigin(0.5);
    // Day 14-15 will add renderDiscardRow, renderStockCount, TurnIndicator
  }
}
