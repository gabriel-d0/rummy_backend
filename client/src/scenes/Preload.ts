import Phaser from "phaser";
import { ASSET_MANIFEST } from "./assets";

export class Preload extends Phaser.Scene {
  constructor() {
    super("Preload");
  }

  preload() {
    // Day 10: Asset manifest — single source of truth for all 4 assets
    for (const [key, path] of Object.entries(ASSET_MANIFEST)) {
      this.load.image(key, path);
    }
  }

  create() {
    console.log("Preload complete — Day 10 manifest tile + joker + table + rack loaded");
    this.scene.start("TableScene");
    this.scene.launch("RackScene");
  }
}
