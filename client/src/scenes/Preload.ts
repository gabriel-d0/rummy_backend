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
    // Day 8-10 will add table.png, rack.png
  }

  create() {
    console.log("Preload complete — Day 7 tile + joker loaded");
    this.scene.start("TableScene");
    this.scene.launch("RackScene");
  }
}
