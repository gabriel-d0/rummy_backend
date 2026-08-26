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

const game = new Phaser.Game(config);
// Expose for Playwright e2e — allows inspection of scenes and layout without private leak
if (typeof window !== "undefined") {
  (window as unknown as Record<string, unknown>).__GAME__ = game;
  (window as unknown as Record<string, unknown>).__GAME_SPACE__ = GAME_SPACE;
}
console.log("Phaser 3 Rummy — Day 20+ LayoutManager 1000×1000 GameSpace — see client/docs/layout.md and src/ui/LayoutManager.ts");
