import Phaser from "phaser";

export class Preload extends Phaser.Scene {
  constructor() {
    super("Preload");
  }

  preload() {
    // Day 2: minimal preload, no assets yet — just log
    console.log("Preload: Day 2 scaffold — no assets");
  }

  create() {
    console.log("Preload complete — Day 2 scaffold");
    this.add.text(512, 384, "Rummy — Day 2\nVite + TS + Phaser\n1024×768", {
      fontFamily: "monospace",
      fontSize: "20px",
      color: "#ffffff",
      align: "center",
    }).setOrigin(0.5);
  }
}
