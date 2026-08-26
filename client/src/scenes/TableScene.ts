import Phaser from "phaser";
import { renderDiscardRow } from "../ui/DiscardRow";
import { renderStockCount, renderTurnIndicator } from "../ui/StockCount";
import { renderTableMelds } from "../ui/TableMelds";

// Day 15: TableScene renders PublicView.TableMelds/DiscardRow/StockCount/CurrentSeat
export class TableScene extends Phaser.Scene {
  constructor() {
    super("TableScene");
  }

  create() {
    // Day 8: table background (green felt 1024x768) behind tiles
    this.add.image(512, 384, "table").setDisplaySize(1024, 768);
    // Day 18: drop zone at y=200 for local drag drop (no server call yet)
    const dropZone = this.add
      .rectangle(512, 200, 800, 80, 0xffffff, 0.05)
      .setStrokeStyle(1, 0xffff00, 0.5)
      .setInteractive({ dropZone: true });
    dropZone.setData("isDropZone", true);
    this.add
      .text(512, 200, "Drop zone — Day 18 (y=200) — drag tile here", {
        fontFamily: "monospace",
        fontSize: "10px",
        color: "#ffff00",
        align: "center",
      })
      .setOrigin(0.5);
    this.input.on("drop", (_pointer: any, gameObject: any, dropZoneObj: any) => {
      if (dropZoneObj.getData("isDropZone")) {
        const tileId = gameObject.getData("tileId");
        console.log(`drop ${tileId} at y=200`);
      }
    });
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
    // Day 14: renderDiscardRow at x = 100 + i*40 y = 300 with mock IsOpeningDiscard flagged
    const mockDiscardRow = [
      { Tile: { ID: "disc-open", Colour: 1, Rank: 7, IsJoker: false }, IsOpeningDiscard: true, Index: 0 },
      { Tile: { ID: "disc-1", Colour: 2, Rank: 3, IsJoker: false }, IsOpeningDiscard: false, Index: 1 },
      { Tile: { ID: "disc-2", Colour: 3, Rank: 9, IsJoker: false }, IsOpeningDiscard: false, Index: 2 },
    ];
    renderDiscardRow(this, mockDiscardRow, { x: 100, y: 300, spacing: 40 });
    // Day 15: renderStockCount and TurnIndicator at x=800 y=50 with mock PublicView
    renderStockCount(this, 77, { x: 800, y: 50 });
    renderTurnIndicator(this, 0, "Playing", "MustDraw", { x: 800, y: 80 });
    this.add
      .text(512, 340, "TableScene — Day 18 drop zone at y=200", {
        fontFamily: "monospace",
        fontSize: "10px",
        color: "#ffff00",
        align: "center",
      })
      .setOrigin(0.5);
  }
}
