import Phaser from "phaser";
import { getSubspaces, drawDebugSubspaces, GAME_SPACE } from "../ui/LayoutManager";

export class Preload extends Phaser.Scene {
  constructor() {
    super("Preload");
  }

  preload() {
    console.log("Preload: Day 4 LayoutManager — GameSpace 1000x1000 + 6 subspaces");
  }

  create() {
    console.log("Preload complete — Day 4 LayoutManager 1000×1000 GameSpace + TopBar/TableArea/MeldArea/DiscardRowArea/PlayerRackArea/ActionButtonsArea");
    const s = getSubspaces();

    // Draw GameSpace background
    this.add.rectangle(GAME_SPACE.width / 2, GAME_SPACE.height / 2, GAME_SPACE.width, GAME_SPACE.height, 0x0a4d2e);

    // Draw subspaces for visual verification (debug)
    drawDebugSubspaces(this, s, 0.06);

    // Example: place elements using proportional math (no pixel values)
    // TopBar: centered text
    this.add
      .text(s.TopBar.x + s.TopBar.width / 2, s.TopBar.y + s.TopBar.height / 2, "TopBar — Stock • Turn • Winner", {
        fontFamily: "monospace",
        fontSize: "10px",
        color: "#ffffff",
        backgroundColor: "#00000066",
        padding: { x: 8, y: 4 },
      })
      .setOrigin(0.5);

    // MeldArea: centered tiles
    this.add
      .text(s.MeldArea.x + s.MeldArea.width / 2, s.MeldArea.y + s.MeldArea.height / 2, "MeldArea — 3 tiles centered", {
        fontFamily: "monospace",
        fontSize: "9px",
        color: "#ffffff",
      })
      .setOrigin(0.5);

    // DiscardRowArea: centered
    this.add
      .text(s.DiscardRowArea.x + s.DiscardRowArea.width / 2, s.DiscardRowArea.y + s.DiscardRowArea.height / 2, "DiscardRowArea", {
        fontFamily: "monospace",
        fontSize: "9px",
        color: "#ffffff",
      })
      .setOrigin(0.5);

    // PlayerRackArea
    this.add
      .text(s.PlayerRackArea.x + s.PlayerRackArea.width / 2, s.PlayerRackArea.y + s.PlayerRackArea.height / 2, "PlayerRackArea — 14 slots", {
        fontFamily: "monospace",
        fontSize: "9px",
        color: "#ffffff",
      })
      .setOrigin(0.5);

    // ActionButtonsArea
    this.add
      .text(s.ActionButtonsArea.x + s.ActionButtonsArea.width / 2, s.ActionButtonsArea.y + s.ActionButtonsArea.height / 2, "ActionButtonsArea — Draw • Prev • Meld • Discard", {
        fontFamily: "monospace",
        fontSize: "9px",
        color: "#ffffff",
      })
      .setOrigin(0.5);

    // Title at center of TableArea
    this.add
      .text(s.TableArea.x + s.TableArea.width / 2, s.TableArea.y + s.TableArea.height / 2 - 60, "Rummy — Day 4\nLayoutManager 1000×1000\n6 Subspaces", {
        fontFamily: "monospace",
        fontSize: "14px",
        color: "#ffffff",
        align: "center",
        backgroundColor: "#00000088",
        padding: { x: 12, y: 8 },
      })
      .setOrigin(0.5);

    // Handle resize — recalculate layout
    this.scale.on("resize", () => {
      const ns = getSubspaces();
      console.log("Resize — recalculated subspaces", ns);
    });
  }
}
