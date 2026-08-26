import Phaser from "phaser";
import { discardSelected, meldSelected, renderRack, sortRack, clearSelected } from "../ui/Rack";
import { getSubspaces, GAME_SPACE } from "../ui/LayoutManager";
import { subscribePrivateSnapshot, canPlayerAct, canDrawPrevious } from "../state/sync";
import type { PrivateSnapshot } from "../state/snapshot";
import {
  OpClientDiscard,
  OpClientDrawStock,
  OpClientDrawPreviousDiscard,
  OpClientMeldInitial,
  OpClientMeldNew,
} from "../net/protocol";
import { sendMatchState } from "../net/protocol";
import { createSocket, getStoredMatchId } from "../net/nakama";

// Day 11 + Day 19: RackScene renders PrivateView.OwnRack only (redaction) with proportional layout
// Day 30: subscribes to PrivateSnapshot and re-renders OwnRack only (never foreign rack)
// Fix: unified to GAME_SPACE getSubspaces, buttons in ActionButtonsArea with flex, no overlap

export class RackScene extends Phaser.Scene {
  constructor() {
    super("RackScene");
  }

  private renderPrivateRack(snap: PrivateSnapshot): void {
    const s = getSubspaces();
    const sorted = sortRack(snap.ownRack);
    // Center within PlayerRackArea, using slot logic or centered
    const area = s.PlayerRackArea;
    const totalW = sorted.length > 0 ? (sorted.length - 1) * 62 : 0;
    const centerX = area.x + area.width / 2;
    const startX = centerX - totalW / 2;
    const y = area.y + area.height / 2;
    renderRack(this, sorted, snap.ownSeat, {
      x: startX,
      y,
      spacing: 62,
    });
  }

  private setButtonEnabled(btn: Phaser.GameObjects.Text, enabled: boolean): void {
    btn.setAlpha(enabled ? 1 : 0.35);
    if (enabled) {
      btn.setInteractive({ useHandCursor: true });
    } else {
      btn.disableInteractive();
    }
  }

  private updateButtonStates(
    snap: PrivateSnapshot,
    discardBtn: Phaser.GameObjects.Text,
    meldRunBtn: Phaser.GameObjects.Text,
    meldSetBtn: Phaser.GameObjects.Text,
    drawBtn?: Phaser.GameObjects.Text,
    drawPrevBtn?: Phaser.GameObjects.Text
  ): void {
    const canDiscard = canPlayerAct(snap, OpClientDiscard);
    const canMeldInitial = canPlayerAct(snap, OpClientMeldInitial);
    const canMeldNew = canPlayerAct(snap, OpClientMeldNew);
    const canMeld = canMeldInitial || canMeldNew;
    const canDraw = snap ? canPlayerAct(snap, OpClientDrawStock) : false;
    const canPrev = snap ? canDrawPrevious(snap) : false;
    this.setButtonEnabled(discardBtn, canDiscard);
    this.setButtonEnabled(meldRunBtn, canMeld);
    this.setButtonEnabled(meldSetBtn, canMeld);
    if (drawBtn) this.setButtonEnabled(drawBtn, canDraw);
    if (drawPrevBtn) this.setButtonEnabled(drawPrevBtn, canPrev);
    console.log(
      `Day 36 button states: discard=${canDiscard} meld=${canMeld} draw=${canDraw} prev=${canPrev} phase=${snap.gamePhase}/${snap.turnPhase} seat=${snap.ownSeat} current=${snap.currentSeat} opened=${snap.players.find((p) => p.seat === snap.ownSeat)?.hasOpened}`
    );
  }

  create() {
    const s = getSubspaces();

    // Rack background fills PlayerRackArea (wood)
    this.add
      .image(
        s.PlayerRackArea.x + s.PlayerRackArea.width / 2,
        s.PlayerRackArea.y + s.PlayerRackArea.height / 2,
        "rack"
      )
      .setDisplaySize(s.PlayerRackArea.width, s.PlayerRackArea.height);

    // Draw slot outlines for visual clarity (14 slots) centered in PlayerRackArea
    const slotW = 48;
    const slotH = 80;
    const n = 14;
    const totalSlotW = n * slotW + (n - 1) * 6;
    const slotStartX = s.PlayerRackArea.x + (s.PlayerRackArea.width - totalSlotW) / 2;
    const slotY = s.PlayerRackArea.y + (s.PlayerRackArea.height - slotH) / 2;
    for (let i = 0; i < n; i++) {
      this.add
        .rectangle(
          slotStartX + i * (slotW + 6) + slotW / 2,
          slotY + slotH / 2,
          slotW,
          slotH,
          0x3d2817,
          0
        )
        .setStrokeStyle(1, 0x5a3d1a, 0.5);
    }

    // Day 12 + Day 16: sorted rack with selection (mock for initial render before any PrivateSnapshot)
    const unsorted = [
      { ID: "mock-red-13", Colour: 1, Rank: 13, IsJoker: false },
      { ID: "mock-red-1", Colour: 1, Rank: 1, IsJoker: false },
      { ID: "mock-blue-5", Colour: 3, Rank: 5, IsJoker: false },
    ];
    const mockRack = sortRack(unsorted);
    const mockArea = s.PlayerRackArea;
    const mockTotalW = mockRack.length > 0 ? (mockRack.length - 1) * 62 : 0;
    const mockCenterX = mockArea.x + mockArea.width / 2;
    const mockStartX = mockCenterX - mockTotalW / 2;
    renderRack(this, mockRack, 0, {
      x: mockStartX,
      y: mockArea.y + mockArea.height / 2,
      spacing: 62,
    });

    // Day 17: dragstart — store dragged tileId and original position for snap-back
    const dragOrig = new Map<string, { x: number; y: number }>();
    this.input.on("dragstart", (_pointer: unknown, gameObject: unknown) => {
      const g = gameObject as { getData?: (k: string) => unknown; x: number; y: number };
      const tileId = String(g?.getData?.("tileId") ?? "");
      dragOrig.set(tileId, { x: g.x, y: g.y });
      // Bring to top
      const img = gameObject as Phaser.GameObjects.Image;
      this.children.bringToTop(img);
      console.log(`dragstart ${String(tileId)}`);
    });

    // Also handle drag for cross-scene: on drag, bring tile to top, on dragend check drop zone
    this.input.on(
      "drag",
      (_pointer: unknown, gameObject: unknown, dragX: number, dragY: number) => {
        const g = gameObject as Phaser.GameObjects.Image;
        g.x = dragX;
        g.y = dragY;
      }
    );

    this.input.on("dragend", (pointer: unknown, gameObject: unknown) => {
      const g = gameObject as {
        getData?: (k: string) => unknown;
        x: number;
        y: number;
        setPosition?: (x: number, y: number) => void;
      };
      const tileId = String(g?.getData?.("tileId") ?? "");
      const orig = dragOrig.get(tileId);
      // Check if dropped in drop zone (which is in TableScene at GAME_SPACE.width/2, dropY)
      // For cross-scene, we check pointer position against drop zone bounds in GAME_SPACE
      const p = pointer as { x: number; y: number };
      // Drop zone is at bottom of TableArea (see TableScene) — centered, 500x36
      const dropH = 36;
      const dropY = s.TableArea.y + s.TableArea.height - dropH - 8;
      const dropX = GAME_SPACE.width / 2;
      const dropW = 500;
      const inDropZone =
        p.x >= dropX - dropW / 2 &&
        p.x <= dropX + dropW / 2 &&
        p.y >= dropY &&
        p.y <= dropY + dropH;
      if (inDropZone) {
        console.log(`drop ${String(tileId)} at dropZone — RackScene dragend`);
        this.game.events.emit("rummy:drop", { tileId, x: p.x, y: p.y });
        dragOrig.delete(tileId);
        return;
      }
      // Snap back to rack if not dropped in zone — restore original position
      if (orig && typeof (g as unknown as { setPosition?: unknown }).setPosition === "function") {
        (g as unknown as { setPosition: (x: number, y: number) => void }).setPosition(
          orig.x,
          orig.y
        );
      } else if (orig) {
        g.x = orig.x;
        g.y = orig.y;
      }
      dragOrig.delete(tileId);
    });

    // ActionButtonsArea: flex layout for 5 buttons without overlap, centered, no info overlap
    const ab = s.ActionButtonsArea;
    const btnY = ab.y + ab.height / 2;
    // Flex: distribute 5 buttons evenly in ActionButtonsArea
    const btnCount = 5;
    const btnW = 110;
    const gap = (ab.width - btnCount * btnW) / (btnCount + 1);
    const btnXs = Array.from(
      { length: btnCount },
      (_, i) => ab.x + gap + btnW / 2 + i * (btnW + gap)
    );

    // Track latest private snapshot for opening discard handler
    let latestPrivate: PrivateSnapshot | null = null;

    // Day 39: Draw button — at btnXs[0]
    const drawBtn = this.add
      .text(btnXs[0], btnY, "[Draw]", {
        fontFamily: "monospace",
        fontSize: "12px",
        color: "#00ff00",
        backgroundColor: "#1a3d2e",
        padding: { x: 8, y: 4 },
      })
      .setOrigin(0.5)
      .setInteractive({ useHandCursor: true })
      .setAlpha(0.5)
      .disableInteractive();
    drawBtn.setData("isDrawBtn", true);
    drawBtn.on("pointerdown", async () => {
      if (drawBtn.alpha < 1) return;
      try {
        const sock = await createSocket();
        const matchId = getStoredMatchId();
        if (!matchId) {
          console.log("draw failed: no matchId");
          return;
        }
        drawBtn.disableInteractive();
        drawBtn.setAlpha(0.5);
        await sendMatchState(sock, matchId, OpClientDrawStock, {}, `req-draw-${Date.now()}`);
        console.log("sent OpClientDrawStock — Day 39");
      } catch (e) {
        console.log("draw failed", e);
      }
    });

    // Day 40: DrawPrevious button — at btnXs[1]
    const drawPrevBtn = this.add
      .text(btnXs[1], btnY, "[Prev]", {
        fontFamily: "monospace",
        fontSize: "12px",
        color: "#00ff00",
        backgroundColor: "#1a3d2e",
        padding: { x: 8, y: 4 },
      })
      .setOrigin(0.5)
      .setInteractive({ useHandCursor: true })
      .setAlpha(0.5)
      .disableInteractive();
    drawPrevBtn.setData("isDrawPrevBtn", true);
    drawPrevBtn.on("pointerdown", async () => {
      if (drawPrevBtn.alpha < 1) return;
      try {
        const sock = await createSocket();
        const matchId = getStoredMatchId();
        if (!matchId) {
          console.log("drawPrev failed: no matchId");
          return;
        }
        drawPrevBtn.disableInteractive();
        drawPrevBtn.setAlpha(0.5);
        await sendMatchState(
          sock,
          matchId,
          OpClientDrawPreviousDiscard,
          {},
          `req-prev-${Date.now()}`
        );
        console.log("sent OpClientDrawPreviousDiscard — Day 40");
      } catch (e) {
        console.log("drawPrev failed", e);
      }
    });

    // Day 19: Discard button — at btnXs[4] (rightmost)
    const discardBtn = this.add
      .text(btnXs[4], btnY, "[Discard]", {
        fontFamily: "monospace",
        fontSize: "12px",
        color: "#00ff00",
        backgroundColor: "#1a3d2e",
        padding: { x: 8, y: 4 },
      })
      .setOrigin(0.5)
      .setInteractive({ useHandCursor: true });
    discardBtn.on("pointerdown", async () => {
      const res = discardSelected();
      if (!res) return;
      if (
        latestPrivate &&
        latestPrivate.gamePhase === "OpeningDiscard" &&
        latestPrivate.ownSeat === latestPrivate.currentSeat &&
        latestPrivate.ownRack.length === 15
      ) {
        const inRack = latestPrivate.ownRack.some((t) => t.ID === res.tileId);
        if (!inRack) {
          console.log(`discard failed: tile ${res.tileId} not in ownRack 15`);
          return;
        }
        try {
          const sock = await createSocket();
          const matchId = getStoredMatchId();
          if (!matchId) {
            console.log("discard failed: no matchId");
            return;
          }
          await sendMatchState(
            sock,
            matchId,
            OpClientDiscard,
            { tileId: res.tileId },
            `req-discard-${Date.now()}`
          );
          console.log(`sent OpClientDiscard opening ${res.tileId} — Day 38`);
          clearSelected();
        } catch (e) {
          console.log("discard failed", e);
        }
        return;
      }
      console.log(
        `discardSelected success: ${res.tileId} (not OpeningDiscard — Day 38 only logs, Day 42 will send)`
      );
    });

    // Day 20: Meld buttons — at btnXs[2] and btnXs[3]
    const meldSetBtn = this.add
      .text(btnXs[2], btnY, "[Meld Set]", {
        fontFamily: "monospace",
        fontSize: "12px",
        color: "#00ff00",
        backgroundColor: "#1a3d2e",
        padding: { x: 8, y: 4 },
      })
      .setOrigin(0.5)
      .setInteractive({ useHandCursor: true });
    meldSetBtn.on("pointerdown", () => {
      const res = meldSelected("set");
      if (res) console.log(`meldSelected success: set ${res.tileIds.join(",")}`);
    });

    const meldRunBtn = this.add
      .text(btnXs[3], btnY, "[Meld Run]", {
        fontFamily: "monospace",
        fontSize: "12px",
        color: "#00ff00",
        backgroundColor: "#1a3d2e",
        padding: { x: 8, y: 4 },
      })
      .setOrigin(0.5)
      .setInteractive({ useHandCursor: true });
    meldRunBtn.on("pointerdown", () => {
      const res = meldSelected("run");
      if (res) console.log(`meldSelected success: run ${res.tileIds.join(",")}`);
    });

    // Set initial button states based on mock snapshot (Playing MustDraw, ownSeat 0) so Draw is enabled
    const mockSnap: PrivateSnapshot = {
      v: 1,
      gamePhase: "Playing",
      turnPhase: "MustDraw",
      currentSeat: 0,
      players: [
        { id: "alice", seat: 0, hasOpened: false, rackCount: 3 },
        { id: "bob", seat: 1, hasOpened: false, rackCount: 14 },
      ],
      stockCount: 77,
      discardRow: [],
      tableMelds: [],
      winner: -1,
      ownRack: mockRack,
      ownSeat: 0,
    };
    this.updateButtonStates(mockSnap, discardBtn, meldRunBtn, meldSetBtn, drawBtn, drawPrevBtn);

    // Day 30: subscribe to PrivateSnapshot — re-render OwnRack only (redaction)
    // Day 36: also update button states per CurrentSeat/TurnPhase/HasOpened
    // Day 38: track latestPrivate for opening discard handler
    // Day 39: also update Draw button per MustDraw
    // Day 40: also update Prev button per HasOpened+DiscardRow
    const unsubscribe = subscribePrivateSnapshot((snap) => {
      latestPrivate = snap;
      this.renderPrivateRack(snap);
      this.updateButtonStates(snap, discardBtn, meldRunBtn, meldSetBtn, drawBtn, drawPrevBtn);
    });
    this.events.once("shutdown", unsubscribe);
    this.events.once("destroy", unsubscribe);

    // Info text just above PlayerRackArea, left-aligned, not overlapping buttons
    this.add
      .text(
        s.PlayerRackArea.x + 8,
        s.PlayerRackArea.y - 10,
        "Rack — Day 20 meldSelected + Day 19 discardSelected",
        {
          fontFamily: "monospace",
          fontSize: "9px",
          color: "#ffff00",
          backgroundColor: "#00000055",
          padding: { x: 4, y: 2 },
        }
      )
      .setOrigin(0, 1)
      .setData("isRackInfo", true);

    this.scale.on("resize", () => {
      this.scene.restart();
    });
  }
}
