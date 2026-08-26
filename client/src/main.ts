import Phaser from "phaser";
import { Preload } from "./scenes/Preload";
import { TableScene } from "./scenes/TableScene";
import { RackScene } from "./scenes/RackScene";
import { GAME_SPACE } from "./ui/LayoutManager";

const config: Phaser.Types.Core.GameConfig = {
  type: Phaser.AUTO,
  parent: "game",
  width: GAME_SPACE.width,
  height: GAME_SPACE.height,
  backgroundColor: "#0a4d2e",
  scale: GAME_SPACE.scale,
  scene: [Preload, TableScene, RackScene],
};

new Phaser.Game(config);
console.log("Phaser 3 Rummy — Day 20+ LayoutManager 1000×1000 GameSpace — see client/docs/layout.md and src/ui/LayoutManager.ts");
