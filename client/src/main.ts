import Phaser from "phaser";
import { Preload } from "./scenes/Preload";

const config: Phaser.Types.Core.GameConfig = {
  type: Phaser.AUTO,
  parent: "game",
  width: 1024,
  height: 768,
  backgroundColor: "#0a4d2e",
  scene: [Preload],
};

new Phaser.Game(config);
console.log("Phaser 3 Rummy — Day 2 Vite + TypeScript + Phaser scaffold — 1024x768 Preload");
