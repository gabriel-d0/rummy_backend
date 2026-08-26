import Phaser from "phaser";

export class Preload extends Phaser.Scene {
  constructor() {
    super("Preload");
  }

  preload() {
    // Day 6: Preload scene loads a single 1x1 tile sprite (placeholder)
    // Day 7-10 will add tile.png, joker.png, table.png, rack.png
    this.load.image(
      "tile",
      "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8BwHwAF/gL+6XhIAAAAAElFTkSuQmCC"
    );
  }

  create() {
    console.log("Preload complete — Day 6");
    this.scene.start("TableScene");
    this.scene.launch("RackScene");
  }
}
