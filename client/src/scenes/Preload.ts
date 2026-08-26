import Phaser from "phaser";

export class Preload extends Phaser.Scene {
  constructor() {
    super("Preload");
  }

  preload() {
    // Day 6: Preload scene loads a single 1x1 tile sprite (placeholder)
    // Day 7: Add tile.png for red-1 and joker.png for Joly
    this.load.image("tile", "assets/tile.png");
    this.load.image("joker", "assets/joker.png");
    // Day 8: Add table.png (green felt 1024x768)
    this.load.image("table", "assets/table.png");
    // Day 9: Add rack.png (wood 800x120 with 14 slots)
    this.load.image("rack", "assets/rack.png");
  }

  create() {
    console.log("Preload complete — Day 9 tile + joker + table + rack loaded");
    this.scene.start("TableScene");
    this.scene.launch("RackScene");
  }
}
