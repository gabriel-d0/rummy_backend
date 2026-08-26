import Phaser from "phaser";
import { Preload } from "./scenes/Preload";
import { TableScene } from "./scenes/TableScene";
import { RackScene } from "./scenes/RackScene";

const config: Phaser.Types.Core.GameConfig = {
  type: Phaser.AUTO,
  parent: "game",
  width: 1024,
  height: 768,
  backgroundColor: "#0a4d2e",
  scale: {
    mode: Phaser.Scale.FIT,
    autoCenter: Phaser.Scale.CENTER_BOTH,
    width: 1024,
    height: 768,
  },
  scene: [Preload, TableScene, RackScene],
};

new Phaser.Game(config);
console.log("Phaser 3 Rummy — Day 20+ subspace Layout — see client/docs/roadmap.md and src/ui/Layout.ts");
