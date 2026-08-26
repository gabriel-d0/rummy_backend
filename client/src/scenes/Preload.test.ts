import { describe, it, expect } from "vitest";
import { ASSET_MANIFEST } from "./assets";
import * as fs from "fs";
import * as path from "path";

describe("Preload manifest — Day 10", () => {
  it("has 4 entries (tile, joker, table, rack)", () => {
    expect(Object.keys(ASSET_MANIFEST)).toHaveLength(4);
    expect(ASSET_MANIFEST.tile).toBe("assets/tile.png");
    expect(ASSET_MANIFEST.joker).toBe("assets/joker.png");
    expect(ASSET_MANIFEST.table).toBe("assets/table.png");
    expect(ASSET_MANIFEST.rack).toBe("assets/rack.png");
  });

  it("all 4 asset files exist on disk (no 404)", () => {
    for (const [key, relPath] of Object.entries(ASSET_MANIFEST)) {
      const fullPath = path.resolve(__dirname, "../../public", relPath);
      expect(fs.existsSync(fullPath), `asset ${key} at ${relPath} should exist`).toBe(true);
      const stat = fs.statSync(fullPath);
      expect(stat.size).toBeGreaterThan(0);
    }
  });

  it("all 4 asset files are under public/assets (Vite publicDir)", () => {
    for (const relPath of Object.values(ASSET_MANIFEST)) {
      expect(relPath.startsWith("assets/")).toBe(true);
    }
  });
});
