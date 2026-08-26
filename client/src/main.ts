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
  scene: [Preload, TableScene, RackScene],
};

new Phaser.Game(config);
console.log("Phaser 3 Rummy — Day 1 scaffolding — see client/docs/roadmap.md");
