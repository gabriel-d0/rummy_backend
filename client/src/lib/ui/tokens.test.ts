import { describe, it, expect } from "vitest";
import { colors, fonts, tile } from "./tokens";

describe("Design tokens — Day 3", () => {
  it("colors felt and wood are defined", () => {
    expect(colors.felt).toBe("#0a2e1a");
    expect(colors.wood).toBe("#5d4037");
  });

  it("tile colors are correct hex", () => {
    expect(colors.tile.red).toBe("#e53935");
    expect(colors.tile.yellow).toBe("#f9a825");
    expect(colors.tile.blue).toBe("#1e88e5");
    expect(colors.tile.black).toBe("#212121");
  });

  it("fonts are Inter and JetBrains Mono", () => {
    expect(fonts.sans).toContain("Inter");
    expect(fonts.mono).toContain("JetBrains Mono");
  });

  it("tile sizes are 64x90 rack and 52x72 table", () => {
    expect(tile.rack.w).toBe(64);
    expect(tile.rack.h).toBe(90);
    expect(tile.table.w).toBe(52);
    expect(tile.table.h).toBe(72);
  });
});
