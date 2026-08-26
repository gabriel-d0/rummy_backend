import Phaser from "phaser";
import { renderDiscardRow } from "../ui/DiscardRow";
import { getLayout } from "../ui/Layout";
import { renderStockCount, renderTurnIndicator } from "../ui/StockCount";
import { renderTableMelds } from "../ui/TableMelds";
import {
  subscribePublicSnapshot,
  subscribePrivateSnapshot,
  shouldShowStartButton,
} from "../state/sync";
import type { PublicSnapshot, PrivateSnapshot } from "../state/snapshot";
import { OpClientStart } from "../net/protocol";
import { sendMatchState } from "../net/protocol";
import { createSocket, getStoredMatchId } from "../net/nakama";

// Day 15+ TableScene — now uses subspace Layout for mathematically defined bounds and responsive design
// Day 31: subscribes to PublicSnapshot and re-renders TableMelds/DiscardRow/StockCount/CurrentSeat (not OwnRack)
export class TableScene extends Phaser.Scene {
  constructor() {
    super("TableScene");
  }

  private renderPublicSnapshot(snap: PublicSnapshot): void {
    const layout = getLayout(this.scale.width, this.scale.height);
    // Stock and turn — topBar right-aligned
    renderStockCount(this, snap.stockCount, {
      x: layout.topBar.x + layout.topBar.w - 60,
      y: layout.topBar.y + 18,
    });
    renderTurnIndicator(this, snap.currentSeat, snap.gamePhase, snap.turnPhase, {
      x: layout.topBar.x + layout.topBar.w - 60,
      y: layout.topBar.y + 45,
    });
    // Table melds — inside tableMelds subspace
    renderTableMelds(this, snap.tableMelds, {
      x: layout.tableMelds.x + 10,
      y0: layout.tableMelds.y + 10,
      rowHeight: 70,
      tileSpacing: 48,
    });
    // Discard row — inside discardRow subspace
    renderDiscardRow(this, snap.discardRow, {
      x: layout.discardRow.x + 10,
      y: layout.discardRow.y + 18,
      spacing: 50,
    });
    // Update info text
    const info = this.children.list.find((c) =>
      (c as unknown as { getData?: (k: string) => unknown }).getData?.("isInfo")
    ) as Phaser.GameObjects.Text | undefined;
    if (info) {
      info.setText(
        `Table — Stock:${snap.stockCount} • seat-${snap.currentSeat} ${snap.gamePhase}/${snap.turnPhase}`
      );
    }
  }

  create() {
    const layout = getLayout(this.scale.width, this.scale.height);

    // Background fills entire game (will be resized on resize)
    const bg = this.add
      .image(layout.width / 2, layout.height / 2, "table")
      .setDisplaySize(layout.width, layout.height);
    bg.setData("isBg", true);

    // Outer border for debug (subtle)
    this.add
      .rectangle(
        layout.outer.x + layout.outer.w / 2,
        layout.outer.y + layout.outer.h / 2,
        layout.outer.w,
        layout.outer.h,
        0x000000,
        0
      )
      .setStrokeStyle(2, 0x8b7355, 0.4);

    // Top bar: Stock + Turn are inside topBar, right-aligned, not overlapping melds (mock for initial render before PublicSnapshot)
    renderStockCount(this, 77, {
      x: layout.topBar.x + layout.topBar.w - 60,
      y: layout.topBar.y + 18,
    });
    renderTurnIndicator(this, 0, "Playing", "MustDraw", {
      x: layout.topBar.x + layout.topBar.w - 60,
      y: layout.topBar.y + 45,
    });

    // Day 37: Start button — visible only if Waiting + host Seat 0 + >=2 players
    const startBtn = this.add
      .text(layout.topBar.x + 12, layout.topBar.y + 18, "[Start]", {
        fontFamily: "monospace",
        fontSize: "14px",
        color: "#ffffff",
        backgroundColor: "#1a5c2e",
        padding: { x: 10, y: 6 },
      })
      .setOrigin(0, 0.5)
      .setInteractive({ useHandCursor: true })
      .setData("isStartBtn", true)
      .setVisible(false)
      .setAlpha(0.5);

    startBtn.on("pointerdown", async () => {
      if (!startBtn.visible || startBtn.alpha < 1) return;
      startBtn.disableInteractive();
      startBtn.setAlpha(0.5);
      try {
        const sock = await createSocket();
        const matchId = getStoredMatchId();
        if (!matchId) {
          console.log("Start failed: no matchId");
          return;
        }
        await sendMatchState(sock, matchId, OpClientStart, {}, `req-start-${Date.now()}`);
        console.log("sent OpClientStart — Day 37");
      } catch (e) {
        console.log("Start failed", e);
        startBtn.setInteractive({ useHandCursor: true });
        startBtn.setAlpha(1);
      }
    });

    const updateStartBtn = (snap: PrivateSnapshot) => {
      const shouldShow = shouldShowStartButton(snap);
      startBtn.setVisible(shouldShow);
      if (shouldShow) {
        startBtn.setAlpha(1);
        startBtn.setInteractive({ useHandCursor: true });
      } else {
        startBtn.disableInteractive();
        startBtn.setAlpha(0.5);
      }
    };
    const unsubStart = subscribePrivateSnapshot(updateStartBtn);
    this.events.once("shutdown", unsubStart);
    this.events.once("destroy", unsubStart);

    // Table melds: inside tableMelds subspace (mock)
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
    renderTableMelds(this, mockMelds, {
      x: layout.tableMelds.x + 10,
      y0: layout.tableMelds.y + 10,
      rowHeight: 70,
      tileSpacing: 48,
    });

    // Discard row: inside discardRow subspace (mock)
    const mockDiscardRow = [
      {
        Tile: { ID: "disc-open", Colour: 1, Rank: 7, IsJoker: false },
        IsOpeningDiscard: true,
        Index: 0,
      },
      {
        Tile: { ID: "disc-1", Colour: 2, Rank: 3, IsJoker: false },
        IsOpeningDiscard: false,
        Index: 1,
      },
      {
        Tile: { ID: "disc-2", Colour: 3, Rank: 9, IsJoker: false },
        IsOpeningDiscard: false,
        Index: 2,
      },
    ];
    renderDiscardRow(this, mockDiscardRow, {
      x: layout.discardRow.x + 10,
      y: layout.discardRow.y + 18,
      spacing: 50,
    });

    // Day 31: subscribe to PublicSnapshot — re-render TableMelds/DiscardRow/StockCount/CurrentSeat (not OwnRack)
    const unsubscribe = subscribePublicSnapshot((snap) => {
      this.renderPublicSnapshot(snap);
    });
    this.events.once("shutdown", unsubscribe);
    this.events.once("destroy", unsubscribe);

    // Drop zone: inside dropZone subspace
    const dz = layout.dropZone;
    const dropZone = this.add
      .rectangle(dz.x + dz.w / 2, dz.y + dz.h / 2, dz.w, dz.h, 0xffffff, 0.04)
      .setStrokeStyle(1, 0xffff00, 0.4)
      .setInteractive({ dropZone: true });
    dropZone.setData("isDropZone", true);
    this.add
      .text(dz.x + dz.w / 2, dz.y + dz.h / 2, "Drop zone — drag tile here", {
        fontFamily: "monospace",
        fontSize: "10px",
        color: "#ffff00",
        align: "center",
      })
      .setOrigin(0.5);
    this.input.on("drop", (_pointer: any, gameObject: any, dropZoneObj: any) => {
      if (dropZoneObj.getData("isDropZone")) {
        const tileId = gameObject.getData("tileId");
        console.log(`drop ${tileId} at dropZone`);
      }
    });

    // Info text inside info subspace (below drop, above rack)
    this.add
      .text(
        layout.info.x + layout.info.w / 2,
        layout.info.y + 8,
        "Table — Stock:77 • seat-0 Playing/MustDraw",
        {
          fontFamily: "monospace",
          fontSize: "10px",
          color: "#ffff00",
          backgroundColor: "#00000066",
          padding: { x: 6, y: 2 },
        }
      )
      .setOrigin(0.5)
      .setData("isInfo", true);

    // Handle resize — recompute layout and re-render static mocks
    this.scale.on("resize", (gameSize: Phaser.Structs.Size) => {
      const nl = getLayout(gameSize.width, gameSize.height);
      // For MVP we just re-create the scene on resize (simpler than moving all objects)
      // In production, we would reposition each subspace element
      this.scene.restart();
      // Note: RackScene is separate but shares same scale event; it will also recompute
      void nl;
    });
  }
}
