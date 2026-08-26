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

    this.add.rectangle(GAME_SPACE.width / 2, GAME_SPACE.height / 2, GAME_SPACE.width, GAME_SPACE.height, 0x0a4d2e);
    this.add.rectangle(
      s.TableArea.x + s.TableArea.width / 2,
      s.TableArea.y + s.TableArea.height / 2,
      s.TableArea.width,
      s.TableArea.height,
      0x0f5d2e,
      0.4
    ).setStrokeStyle(1, 0x1a6b3a, 0.2);

    drawDebugSubspaces(this, s, 0.025);

    this.add.text(s.TopBar.x + s.TopBar.width / 2, s.TopBar.y + s.TopBar.height / 2, "RO  •  Stock  •  Turn  •  Table", {
      fontFamily: "Inter",
      fontSize: "13px",
      color: "#e8f5e9",
      fontStyle: "bold",
      backgroundColor: "#0a2e1a",
      padding: { x: 14, y: 8 },
    }).setOrigin(0.5);

    this.add.text(s.MeldArea.x + s.MeldArea.width / 2, s.MeldArea.y + s.MeldArea.height / 2, "MELDS  —  runs & sets", {
      fontFamily: "Inter",
      fontSize: "11px",
      color: "#c8eac8",
      fontStyle: "bold",
      backgroundColor: "#0a2e1a99",
      padding: { x: 10, y: 6 },
    }).setOrigin(0.5);

    this.add.text(s.DiscardRowArea.x + s.DiscardRowArea.width / 2, s.DiscardRowArea.y + s.DiscardRowArea.height / 2, "DISCARD  ROW", {
      fontFamily: "Inter",
      fontSize: "11px",
      color: "#c8eac8",
      fontStyle: "bold",
      backgroundColor: "#0a2e1a99",
      padding: { x: 10, y: 6 },
    }).setOrigin(0.5);

    this.add.text(s.PlayerRackArea.x + s.PlayerRackArea.width / 2, s.PlayerRackArea.y + s.PlayerRackArea.height / 2, "YOUR  RACK  —  14 tiles", {
      fontFamily: "Inter",
      fontSize: "11px",
      color: "#ffecb3",
      fontStyle: "bold",
      backgroundColor: "#3e2723cc",
      padding: { x: 12, y: 6 },
    }).setOrigin(0.5);

    this.add.text(s.ActionButtonsArea.x + s.ActionButtonsArea.width / 2, s.ActionButtonsArea.y + s.ActionButtonsArea.height / 2, "Draw  •  Prev  •  Meld  •  Discard", {
      fontFamily: "JetBrains Mono",
      fontSize: "10px",
      color: "#a5d6a7",
      fontStyle: "bold",
      backgroundColor: "#0a2e1a",
      padding: { x: 10, y: 6 },
    }).setOrigin(0.5);

    this.add.text(s.TableArea.x + s.TableArea.width / 2, s.TableArea.y + s.TableArea.height / 2 - 36, "Romanian Tile Rummy", {
      fontFamily: "Inter",
      fontSize: "26px",
      color: "#ffffff",
      fontStyle: "bold",
      align: "center",
      shadow: { offsetX: 0, offsetY: 2, color: "#0a2e1a", blur: 6, fill: true },
    }).setOrigin(0.5);

    this.add.text(s.TableArea.x + s.TableArea.width / 2, s.TableArea.y + s.TableArea.height / 2 + 14, "Day 4  •  LayoutManager 1000×1000  •  6 Subspaces  •  FIT + CENTER_BOTH", {
      fontFamily: "Inter",
      fontSize: "11px",
      color: "#e0f2e0",
      fontStyle: "bold",
      backgroundColor: "#0a2e1aee",
      padding: { x: 12, y: 7 },
    }).setOrigin(0.5);

    this.scale.on("resize", () => {
      const ns = getSubspaces();
      console.log("Resize — recalculated subspaces", ns);
    });
  }
}
