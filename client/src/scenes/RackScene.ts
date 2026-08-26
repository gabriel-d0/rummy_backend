import Phaser from "phaser";
import { discardSelected, meldSelected, renderRack, sortRack, clearSelected } from "../ui/Rack";
import { getLayout } from "../ui/Layout";
import { subscribePrivateSnapshot, canPlayerAct } from "../state/sync";
import type { PrivateSnapshot } from "../state/snapshot";
import {
  OpClientDiscard,
  OpClientDrawStock,
  OpClientMeldInitial,
  OpClientMeldNew,
} from "../net/protocol";
import { sendMatchState } from "../net/protocol";
import { createSocket, getStoredMatchId } from "../net/nakama";

// Day 11 + Day 19: RackScene renders PrivateView.OwnRack only (redaction) with subspace layout
// Day 30: subscribes to PrivateSnapshot and re-renders OwnRack only (never foreign rack)
export class RackScene extends Phaser.Scene {
  constructor() {
    super("RackScene");
  }

  private renderPrivateRack(snap: PrivateSnapshot): void {
    const layout = getLayout(this.scale.width, this.scale.height);
    const sorted = sortRack(snap.ownRack);
    const totalW = sorted.length > 0 ? (sorted.length - 1) * 62 : 0;
    const rackCenterX = layout.rack.x + layout.rack.w / 2;
    const startX = rackCenterX - totalW / 2;
    renderRack(this, sorted, snap.ownSeat, {
      x: startX,
      y: layout.rack.y + layout.rack.h / 2,
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
    drawBtn?: Phaser.GameObjects.Text
  ): void {
    const canDiscard = canPlayerAct(snap, OpClientDiscard);
    const canMeldInitial = canPlayerAct(snap, OpClientMeldInitial);
    const canMeldNew = canPlayerAct(snap, OpClientMeldNew);
    const canMeld = canMeldInitial || canMeldNew;
    const canDraw = snap ? canPlayerAct(snap, OpClientDrawStock) : false;
    this.setButtonEnabled(discardBtn, canDiscard);
    this.setButtonEnabled(meldRunBtn, canMeld);
    this.setButtonEnabled(meldSetBtn, canMeld);
    if (drawBtn) this.setButtonEnabled(drawBtn, canDraw);
    // Also log for e2e
    console.log(
      `Day 36 button states: discard=${canDiscard} meld=${canMeld} draw=${canDraw} phase=${snap.gamePhase}/${snap.turnPhase} seat=${snap.ownSeat} current=${snap.currentSeat} opened=${snap.players.find((p) => p.seat === snap.ownSeat)?.hasOpened}`
    );
  }

  create() {
    const layout = getLayout(this.scale.width, this.scale.height);

    // Rack background fills rack subspace (wood 800x120 or full width on mobile, centered)
    this.add
      .image(layout.rack.x + layout.rack.w / 2, layout.rack.y + layout.rack.h / 2, "rack")
      .setDisplaySize(layout.rack.w, layout.rack.h);

    // Draw slot outlines for visual clarity (14 slots)
    for (const slot of layout.rackSlots) {
      this.add
        .rectangle(slot.x + slot.w / 2, slot.y + slot.h / 2, slot.w, slot.h, 0x3d2817, 0)
        .setStrokeStyle(1, 0x5a3d1a, 0.5);
    }

    // Day 12 + Day 16: sorted rack with selection (mock for initial render before any PrivateSnapshot)
    const unsorted = [
      { ID: "mock-red-13", Colour: 1, Rank: 13, IsJoker: false },
      { ID: "mock-red-1", Colour: 1, Rank: 1, IsJoker: false },
      { ID: "mock-blue-5", Colour: 3, Rank: 5, IsJoker: false },
    ];
    const mockRack = sortRack(unsorted);
    // Render centered within rack subspace, not at fixed 100,700
    const totalW = mockRack.length > 0 ? (mockRack.length - 1) * 62 : 0;
    const rackCenterX = layout.rack.x + layout.rack.w / 2;
    const startX = rackCenterX - totalW / 2;
    renderRack(this, mockRack, 0, { x: startX, y: layout.rack.y + layout.rack.h / 2, spacing: 62 });

    // Day 17: dragstart
    this.input.on("dragstart", (_pointer: any, gameObject: any) => {
      const tileId = gameObject.getData("tileId");
      console.log(`dragstart ${tileId}`);
    });

    // Day 19: Discard button inside rack subspace top-right, not overlapping
    // Day 38: Opening discard — when OpeningDiscard + ownSeat==currentSeat + ownRack 15, sends OpClientDiscard 2 {"tileId": selectedId}
    const btnX = layout.rack.x + layout.rack.w - 50;
    const btnY = layout.rack.y - 14;
    const discardBtn = this.add
      .text(btnX, btnY, "[Discard]", {
        fontFamily: "monospace",
        fontSize: "12px",
        color: "#00ff00",
        backgroundColor: "#1a3d2e",
        padding: { x: 8, y: 4 },
      })
      .setOrigin(0.5)
      .setInteractive({ useHandCursor: true });
    // Track latest private snapshot for opening discard handler
    let latestPrivate: PrivateSnapshot | null = null;
    discardBtn.on("pointerdown", async () => {
      const res = discardSelected();
      if (!res) return;
      // Day 38: if in OpeningDiscard and ownRack 15, send to server
      if (
        latestPrivate &&
        latestPrivate.gamePhase === "OpeningDiscard" &&
        latestPrivate.ownSeat === latestPrivate.currentSeat &&
        latestPrivate.ownRack.length === 15
      ) {
        // Validate tile is in ownRack (discardSelected already ensures selected is from rack via UI, but we double-check via snapshot)
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
      // For non-opening (normal discard) just log for now (Day 42 will send)
      console.log(
        `discardSelected success: ${res.tileId} (not OpeningDiscard — Day 38 only logs, Day 42 will send)`
      );
    });
    // Day 39: Draw button — visible only if Playing MustDraw and ownSeat==currentSeat, sends OpClientDrawStock 3 {}
    const drawBtn = this.add
      .text(btnX - 310, btnY, "[Draw]", {
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

    // Day 20: Meld buttons — validate selected.size>=3 and log MELD_INITIAL or MELD_NEW
    const meldRunBtn = this.add
      .text(btnX - 110, btnY, "[Meld Run]", {
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
    const meldSetBtn = this.add
      .text(btnX - 210, btnY, "[Meld Set]", {
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

    // Day 30: subscribe to PrivateSnapshot — re-render OwnRack only (redaction)
    // Day 36: also update button states per CurrentSeat/TurnPhase/HasOpened
    // Day 38: track latestPrivate for opening discard handler
    // Day 39: also update Draw button per MustDraw
    const unsubscribe = subscribePrivateSnapshot((snap) => {
      latestPrivate = snap;
      this.renderPrivateRack(snap);
      this.updateButtonStates(snap, discardBtn, meldRunBtn, meldSetBtn, drawBtn);
    });
    this.events.once("shutdown", unsubscribe);
    this.events.once("destroy", unsubscribe);

    this.add
      .text(
        layout.rack.x + layout.rack.w / 2,
        layout.info.y + 8,
        "Rack — Day 20 meldSelected + Day 19 discardSelected",
        {
          fontFamily: "monospace",
          fontSize: "10px",
          color: "#ffff00",
          backgroundColor: "#00000066",
          padding: { x: 6, y: 2 },
        }
      )
      .setOrigin(0.5);

    this.scale.on("resize", () => {
      this.scene.restart();
    });
  }
}
