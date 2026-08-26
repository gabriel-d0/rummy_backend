import Phaser from "phaser";
import { getSubspaces, drawDebugSubspaces, GAME_SPACE } from "../ui/LayoutManager";
import { loadFonts, createText } from "../ui/fonts";

export class Preload extends Phaser.Scene {
  constructor() {
    super("Preload");
  }

  async preload() {
    console.log("Preload: Day 4 LayoutManager — GameSpace 1000x1000 + 6 subspaces");
    await loadFonts();
  }

  async create() {
    await loadFonts();
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

    createText(this, s.TopBar.x + s.TopBar.width / 2, s.TopBar.y + s.TopBar.height / 2, "RO  •  Stock  •  Turn  •  Table", "label");
    createText(this, s.MeldArea.x + s.MeldArea.width / 2, s.MeldArea.y + s.MeldArea.height / 2, "MELDS  —  runs & sets", "subtitle");
    createText(this, s.DiscardRowArea.x + s.DiscardRowArea.width / 2, s.DiscardRowArea.y + s.DiscardRowArea.height / 2, "DISCARD  ROW", "subtitle");
    createText(this, s.PlayerRackArea.x + s.PlayerRackArea.width / 2, s.PlayerRackArea.y + s.PlayerRackArea.height / 2, "YOUR  RACK  —  14 tiles", "subtitle");
    createText(this, s.ActionButtonsArea.x + s.ActionButtonsArea.width / 2, s.ActionButtonsArea.y + s.ActionButtonsArea.height / 2, "Draw  •  Prev  •  Meld  •  Discard", "mono");
    createText(this, s.TableArea.x + s.TableArea.width / 2, s.TableArea.y + s.TableArea.height / 2 - 36, "Romanian Tile Rummy", "title");
    createText(this, s.TableArea.x + s.TableArea.width / 2, s.TableArea.y + s.TableArea.height / 2 + 14, "Day 4  •  LayoutManager 1000×1000  •  6 Subspaces  •  FIT + CENTER_BOTH", "subtitle");

    this.scale.on("resize", () => {
      const ns = getSubspaces();
      console.log("Resize — recalculated subspaces", ns);
    });
  }
}
