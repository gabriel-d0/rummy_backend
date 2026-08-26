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
    btn.setAlpha(enabled ? 1 : 0.4);
    if (enabled) {
      btn.setInteractive({ useHandCursor: true });
      btn.setStyle({ backgroundColor: "#1a5c2e" });
    } else {
      btn.disableInteractive();
      btn.setStyle({ backgroundColor: "#2a2a2a" });
    }
  }

  private createModernButton(x: number, y: number, label: string): Phaser.GameObjects.Text {
    const btn = this.add.text(x, y, label, {
      fontFamily: "Inter, system-ui, monospace",
      fontSize: "12px",
      color: "#ffffff",
      backgroundColor: "#1a5c2e",
      padding: { x: 14, y: 8 },
      fontStyle: "600",
    });
    btn.setOrigin(0.5);
    btn.setInteractive({ useHandCursor: true });
    // Subtle shadow via stroke
    btn.setStroke("#0f2e1a", 1);
    btn.on("pointerover", () => {
      if (btn.alpha === 1) btn.setStyle({ backgroundColor: "#228a3e" });
    });
    btn.on("pointerout", () => {
      if (btn.alpha === 1) btn.setStyle({ backgroundColor: "#1a5c2e" });
      else btn.setStyle({ backgroundColor: "#2a2a2a" });
    });
    return btn;
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

    // Modern rack background — rounded wood with subtle shadow, not stretched image
    const rackBg = this.add.rectangle(
      s.PlayerRackArea.x + s.PlayerRackArea.width / 2,
      s.PlayerRackArea.y + s.PlayerRackArea.height / 2,
      s.PlayerRackArea.width - 4,
      s.PlayerRackArea.height - 4,
      0x5d4037
    );
    rackBg.setStrokeStyle(2, 0x3e2723, 0.8);
    if ((rackBg as unknown as { setRounded?: (r: number) => void }).setRounded) {
      (rackBg as unknown as { setRounded: (r: number) => void }).setRounded(12);
    }
    // Inner highlight
    this.add
      .rectangle(
        s.PlayerRackArea.x + s.PlayerRackArea.width / 2,
        s.PlayerRackArea.y + s.PlayerRackArea.height / 2,
        s.PlayerRackArea.width - 8,
        s.PlayerRackArea.height - 8,
        0x6d4c41,
        0.5
      )
      .setStrokeStyle(1, 0x8d6e63, 0.3);

    // Draw slot outlines — modern, subtle, rounded, centered
    const slotW = 46;
    const slotH = 68;
    const n = 14;
    const totalSlotW = n * slotW + (n - 1) * 8;
    const slotStartX = s.PlayerRackArea.x + (s.PlayerRackArea.width - totalSlotW) / 2;
    const slotY = s.PlayerRackArea.y + (s.PlayerRackArea.height - slotH) / 2;
    for (let i = 0; i < n; i++) {
      const slot = this.add.rectangle(
        slotStartX + i * (slotW + 8) + slotW / 2,
        slotY + slotH / 2,
        slotW,
        slotH,
        0x3e2723,
        0.4
      );
      slot.setStrokeStyle(1, 0x4e342e, 0.6);
      if ((slot as unknown as { setRounded?: (r: number) => void }).setRounded) {
        (slot as unknown as { setRounded: (r: number) => void }).setRounded(6);
      }
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
      const p = pointer as { x: number; y: number; worldX?: number; worldY?: number };
      const px = p.worldX ?? p.x;
      const py = p.worldY ?? p.y;
      const dropH = 36;
      const dropY = s.TableArea.y + s.TableArea.height - dropH - 8;
      const dropX = GAME_SPACE.width / 2;
      const dropW = 500;
      const inDropZone =
        px >= dropX - dropW / 2 && px <= dropX + dropW / 2 && py >= dropY && py <= dropY + dropH;
      if (inDropZone) {
        console.log(`drop ${String(tileId)} at dropZone — RackScene dragend world ${px},${py}`);
        this.game.events.emit("rummy:drop", { tileId, x: px, y: py });
        dragOrig.delete(tileId);
        return;
      }
      if (orig) {
        // Snap back with tween for modern feel
        this.tweens.add({
          targets: g,
          x: orig.x,
          y: orig.y,
          duration: 180,
          ease: "Back.easeOut",
        });
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

    // Day 39: Draw button — at btnXs[0] — modern pill
    const drawBtn = this.createModernButton(btnXs[0], btnY, "Draw");
    drawBtn.setData("isDrawBtn", true);
    drawBtn.setAlpha(0.5);
    drawBtn.disableInteractive();
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
    const drawPrevBtn = this.createModernButton(btnXs[1], btnY, "Prev");
    drawPrevBtn.setData("isDrawPrevBtn", true);
    drawPrevBtn.setAlpha(0.5);
    drawPrevBtn.disableInteractive();
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

    // Day 19: Discard button — at btnXs[4] (rightmost) — modern, red accent when enabled
    const discardBtn = this.createModernButton(btnXs[4], btnY, "Discard");
    discardBtn.setData("isDiscardBtn", true);
    // Override to red for discard
    discardBtn.setStyle({ backgroundColor: "#5a1a1a" });
    discardBtn.on("pointerover", () => {
      if (discardBtn.alpha === 1) discardBtn.setStyle({ backgroundColor: "#7a2a2a" });
    });
    discardBtn.on("pointerout", () => {
      if (discardBtn.alpha === 1) discardBtn.setStyle({ backgroundColor: "#5a1a1a" });
      else discardBtn.setStyle({ backgroundColor: "#2a2a2a" });
    });
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

    // Day 20: Meld buttons — at btnXs[2] and btnXs[3] — modern
    const meldSetBtn = this.createModernButton(btnXs[2], btnY, "Meld Set");
    meldSetBtn.setData("isMeldSetBtn", true);
    meldSetBtn.on("pointerdown", () => {
      const res = meldSelected("set");
      if (res) console.log(`meldSelected success: set ${res.tileIds.join(",")}`);
    });

    const meldRunBtn = this.createModernButton(btnXs[3], btnY, "Meld Run");
    meldRunBtn.setData("isMeldRunBtn", true);
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
