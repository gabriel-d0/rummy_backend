import Phaser from "phaser";
import { renderDiscardRow } from "../ui/DiscardRow";
import { getSubspaces, GAME_SPACE } from "../ui/LayoutManager";
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

// Day 15+ TableScene — now uses LayoutManager GAME_SPACE 1000x1000 with proportional subspaces
// Day 31: subscribes to PublicSnapshot and re-renders TableMelds/DiscardRow/StockCount/CurrentSeat (not OwnRack)
// Fix: unified to single GAME_SPACE layout, separate TopBar Stock/Turn, no overlap

export class TableScene extends Phaser.Scene {
  constructor() {
    super("TableScene");
  }

  private renderPublicSnapshot(snap: PublicSnapshot): void {
    const s = getSubspaces();
    // TopBar: Stock left, Turn right, no overlap
    renderStockCount(this, snap.stockCount, {
      x: s.TopBar.x + 80,
      y: s.TopBar.y + s.TopBar.height / 2,
    });
    renderTurnIndicator(this, snap.currentSeat, snap.gamePhase, snap.turnPhase, {
      x: s.TopBar.x + s.TopBar.width - 120,
      y: s.TopBar.y + s.TopBar.height / 2,
    });
    // MeldArea and DiscardRowArea are proportional subdivisions of TableArea
    renderTableMelds(this, snap.tableMelds, {
      x: s.MeldArea.x + 10,
      y0: s.MeldArea.y + 10,
      rowHeight: 70,
      tileSpacing: 48,
    });
    renderDiscardRow(this, snap.discardRow, {
      x: s.DiscardRowArea.x + 10,
      y: s.DiscardRowArea.y + 16,
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
    const s = getSubspaces();

    // Background fills entire GAME_SPACE (will be scaled via Scale.FIT)
    const bg = this.add
      .image(GAME_SPACE.width / 2, GAME_SPACE.height / 2, "table")
      .setDisplaySize(GAME_SPACE.width, GAME_SPACE.height);
    bg.setData("isBg", true);

    // Outer border for debug (subtle) - use TableArea outer
    this.add
      .rectangle(
        s.TableArea.x + s.TableArea.width / 2,
        s.TableArea.y + s.TableArea.height / 2,
        s.TableArea.width,
        s.TableArea.height,
        0x000000,
        0
      )
      .setStrokeStyle(2, 0x8b7355, 0.4);

    // TopBar: Stock left, Turn right, Start button left-center — no overlap
    renderStockCount(this, 77, {
      x: s.TopBar.x + 80,
      y: s.TopBar.y + s.TopBar.height / 2,
    });
    renderTurnIndicator(this, 0, "Playing", "MustDraw", {
      x: s.TopBar.x + s.TopBar.width - 120,
      y: s.TopBar.y + s.TopBar.height / 2,
    });

    // Day 37: Start button — visible only if Waiting + host Seat 0 + >=2 players
    const startBtn = this.add
      .text(s.TopBar.x + s.TopBar.width / 2, s.TopBar.y + s.TopBar.height / 2, "[Start]", {
        fontFamily: "monospace",
        fontSize: "14px",
        color: "#ffffff",
        backgroundColor: "#1a5c2e",
        padding: { x: 10, y: 6 },
      })
      .setOrigin(0.5)
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

    // Table melds: inside MeldArea (proportional, not tableMelds subspace)
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
      x: s.MeldArea.x + 10,
      y0: s.MeldArea.y + 10,
      rowHeight: 70,
      tileSpacing: 48,
    });

    // Discard row: inside DiscardRowArea
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
      x: s.DiscardRowArea.x + 10,
      y: s.DiscardRowArea.y + 16,
      spacing: 50,
    });

    // Day 31: subscribe to PublicSnapshot — re-render TableMelds/DiscardRow/StockCount/CurrentSeat (not OwnRack)
    const unsubscribe = subscribePublicSnapshot((snap) => {
      this.renderPublicSnapshot(snap);
    });
    this.events.once("shutdown", unsubscribe);
    this.events.once("destroy", unsubscribe);

    // Drop zone: centered between MeldArea and DiscardRowArea, or at bottom of TableArea
    // Use a dedicated area below DiscardRowArea but above PlayerRackArea
    const dropY = s.DiscardRowArea.y + s.DiscardRowArea.height + 16;
    const dropH = 44;
    const dropZone = this.add
      .rectangle(GAME_SPACE.width / 2, dropY + dropH / 2, 600, dropH, 0xffffff, 0.04)
      .setStrokeStyle(1, 0xffff00, 0.4)
      .setInteractive({ dropZone: true });
    dropZone.setData("isDropZone", true);
    this.add
      .text(GAME_SPACE.width / 2, dropY + dropH / 2, "Drop zone — drag tile here", {
        fontFamily: "monospace",
        fontSize: "10px",
        color: "#ffff00",
        align: "center",
      })
      .setOrigin(0.5);
    // Note: drag is handled in RackScene via global game event to support cross-scene
    // TableScene also listens for drop for when drag starts in TableScene (e.g., meld tiles)
    this.input.on("drop", (_pointer: unknown, gameObject: unknown, dropZoneObj: unknown) => {
      const g = gameObject as { getData?: (k: string) => unknown };
      const d = dropZoneObj as { getData?: (k: string) => unknown };
      if (d?.getData?.("isDropZone")) {
        const tileId = g?.getData?.("tileId");
        console.log(`drop ${String(tileId)} at dropZone — TableScene`);
      }
    });

    // Info text at bottom of TableArea, above PlayerRackArea
    this.add
      .text(
        GAME_SPACE.width / 2,
        s.TableArea.y + s.TableArea.height - 12,
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

    // Handle resize — GAME_SPACE is fixed 1000x1000 with Scale.FIT, so we just restart to re-apply proportional layout
    this.scale.on("resize", () => {
      this.scene.restart();
    });
  }
}
